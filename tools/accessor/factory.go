package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/tools/go/packages"
)

type factory struct {
	dir        string
	file       string
	pkg        string
	targets    []string
	generators map[string]*Generator

	curPkg          *packages.Package
	curFileName     string
	curFileSet      *token.FileSet
	curImports      []string
	curGenerator    *Generator
	curNamedImports map[string]string
}

// NewGenerators .
func NewGenerators(targets []string, dir string, file string) ([]*Generator, error) {
	if len(targets) == 0 || dir == "" {
		return nil, fmt.Errorf("these fields must be non-empty, targets %s, f.dir %s", targets, dir)
	}

	f := &factory{dir: dir, file: file, targets: targets, generators: make(map[string]*Generator)}

	if err := f.walkDir(f.inspectPkg); err != nil {
		return nil, fmt.Errorf("inspectPkg: %w", err)
	}

	if err := f.walkDir(f.inspectDeclaration); err != nil {
		return nil, fmt.Errorf("inspectDeclaration: %w", err)
	}

	generators := lo.MapToSlice(f.generators, func(k string, v *Generator) *Generator { return v })

	return generators, nil
}

func (f *factory) walkDir(fn func(file *ast.File) error) error {
	if f.dir == "" {
		return fmt.Errorf("no dir specified")
	}
	return filepath.Walk(f.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != f.dir {
			return filepath.SkipDir
		}
		if filepath.Ext(info.Name()) != ".go" {
			return nil
		}
		if strings.HasSuffix(info.Name(), "_accessor.go") {
			return nil
		}
		curFileName := strings.TrimSuffix(info.Name(), ".go")
		if f.file != "" && f.file != curFileName {
			return nil
		}
		f.curFileName = curFileName
		f.curFileSet = token.NewFileSet()
		file, err := parser.ParseFile(f.curFileSet, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("parser.ParseFile %s: %w", path, err)
		}
		slog.Debug("begin to parse file %s\n", path)

		cfg := &packages.Config{
			Mode:  packages.LoadAllSyntax,
			Tests: false,
		}
		pkgs, err := packages.Load(cfg, "file="+path)
		if err != nil {
			return fmt.Errorf("failed to load packages: %v", err)
		}
		if packages.PrintErrors(pkgs) > 0 {
			return fmt.Errorf("errors found while loading packages")
		}

		f.curPkg = pkgs[0]

		return fn(file)
	})
}

func (f *factory) inspectPkg(file *ast.File) error {
	switch pkg := file.Name.Name; {
	case f.pkg == "":
		f.pkg = pkg
	case f.pkg != pkg:
		return fmt.Errorf("package name mismatch: %s!= %s", f.pkg, pkg)
	}
	return nil
}

func (f *factory) getGenerator(fileName string) *Generator {
	g, ok := f.generators[fileName]
	if !ok {
		g = NewGenerator(f.dir, f.pkg, fileName)
	}

	return g
}

func (f *factory) inspectDeclaration(node *ast.File) error {
	slog.Info("inspectDeclaration", "file", f.curFileName)

	f.curImports = make([]string, 0)
	f.curNamedImports = make(map[string]string, 0)
	f.curGenerator = f.getGenerator(f.curFileName)

	var err error
	ast.Inspect(node, func(n ast.Node) bool {
		switch ts := n.(type) {
		case *ast.ImportSpec:
			if err = f.inspectImport(ts); err != nil {
				slog.Error("ImportSpec", "ts", ts, "error", err)
				return false
			}
		case *ast.TypeSpec:
			if err = f.inspectTypeSpec(f.curPkg.Types, ts); err != nil {
				slog.Error("TypeSpec", "ts", ts, "error", err)
				return false
			}
		}

		return true
	})

	if len(f.curGenerator.Structs) > 0 {
		f.curGenerator.InspectImports(f.curImports, f.curNamedImports)
		f.generators[f.curGenerator.FileName] = f.curGenerator
	}

	return err
}

func (f *factory) inspectImport(decl *ast.ImportSpec) error {
	path, err := parseNode(decl.Path)
	if err != nil {
		return fmt.Errorf("parseNode: %w", err)
	}

	name := ""
	if decl.Name != nil {
		name = decl.Name.String()
	}
	if name == "_" || name == "." {
		// FIXME: explicit cases should be considered, but we ignore them for simplicity.
		return nil
	}

	if name != "" {
		f.curNamedImports[name] = fmt.Sprintf("%s %s", name, path)
		return nil
	}

	f.curImports = append(f.curImports, path)
	return nil
}

func (f *factory) inspectTypeSpec(pkg *types.Package, spec *ast.TypeSpec) error {
	if !lo.Contains(f.targets, spec.Name.Name) {
		slog.Debug("struct not in targets, skip", spec.Name.Name)
		return nil
	}

	g := f.curGenerator

	fields, err := f.inspectFields(pkg, spec.Name.Name)
	if err != nil {
		return err
	}

	if len(fields) > 0 {
		s := NewStruct(spec.Name.Name, fields)
		g.Structs = append(g.Structs, s)
	}

	return nil
}

func (f *factory) inspectFields(pkg *types.Package, name string) ([]Field, error) {
	fields := []Field{}
	ts := pkg.Scope().Lookup(name).Type().Underlying().(*types.Struct)

	for i := 0; i < ts.NumFields(); i++ {
		tagStr := ts.Tag(i)
		inline := false
		var tag Tag
		if len(tagStr) > 0 {
			// Trim the backticks from the tag value
			tagValue := strings.Trim(tagStr, "`")
			tags, err := ParseTag(tagValue)
			if err != nil {
				return nil, err
			}

			if t, err := tags.Get(tagName); err == nil {
				// skip
				if t.Name == tagIgnore {
					continue
				}
				if t.HasValue(tagInline) {
					inline = true
				}

				tag = *t
			}
		}

		tsField := ts.Field(i)

		// skip embedded if not tag as inline
		if tsField.Embedded() {
			if !inline {
				continue
			}
			p := pkg
			switch t := tsField.Type().(type) {
			case *types.Named:
				p = t.Obj().Pkg()
			case *types.Pointer:
				if named, ok := t.Elem().(*types.Named); ok {
					p = named.Obj().Pkg()
				}
			}

			f.inspectPackageImport(p)

			if items, err := f.inspectFields(p, tsField.Name()); err == nil {
				fields = append(fields, items...)
			} else {
				return nil, err
			}
			continue
		}

		fieldType := tsField.Type()
		comparable := isComparable(fieldType)
		typeStr, err := typeToString(fieldType, f.pkg)
		if err != nil {
			return nil, fmt.Errorf("parseNode: %w", err)
		}

		field := NewField(tsField.Name(), typeStr, comparable, tag)
		fields = append(fields, field)
	}

	return fields, nil
}

func (f *factory) inspectPackageImport(p *types.Package) {
	for _, ims := range p.Imports() {
		name := ims.Name()
		path := ims.Path()

		last, _ := lo.Last(strings.Split(path, "/"))
		if last == name {
			f.curImports = append(f.curImports, strconv.Quote(path))
		} else {
			f.curNamedImports[name] = strconv.Quote(path)
		}
	}
}
