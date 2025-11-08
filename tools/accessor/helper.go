package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"regexp"
	"strings"

	"golang.org/x/tools/go/packages"
)

func concat(strs ...string) (s string) {
	for _, str := range strs {
		if str == "" {
			continue
		}
		s += upper(str)
	}
	return
}

func contains[T comparable](src []T, dst T) bool {
	for _, v := range src {
		if v == dst {
			return true
		}
	}
	return false
}

func upper(str string) string {
	return strings.ToUpper(str[:1]) + str[1:]
}

func fillTypeArguments(t, param, arg string) string {
	re := regexp.MustCompile(`(^|[^_0-9\p{L}])` + regexp.QuoteMeta(param) + `($|[^_0-9\p{L}])`)
	return re.ReplaceAllStringFunc(t, func(s string) string {
		if strings.HasPrefix(s, param) {
			return arg + s[len(param):]
		} else if strings.HasSuffix(s, param) {
			return s[:len(s)-len(param)] + arg
		} else {
			return strings.ReplaceAll(s, param, arg)
		}
	})
}

func getPackageNameFromType(typeName string) string {
	if !strings.Contains(typeName, ".") {
		return ""
	}

	// hacky way to handle anonymous types and generic types
	if strings.Contains(typeName, "{") || strings.Contains(typeName, "[") {
		return ""
	}

	pkgName := strings.Split(typeName, ".")[0]

	// handle pointer type
	pkgName = strings.TrimLeft(pkgName, "*")
	return pkgName
}

// isComparable checks if a given type is comparable.
func isComparable(t types.Type) bool {
	switch t := t.(type) {
	case *types.Basic:
		return t.Info()&types.IsOrdered != 0 || t.Kind() == types.Bool || t.Kind() == types.String
	case *types.Pointer, *types.Interface, *types.Chan:
		return true
	case *types.Struct:
		for i := 0; i < t.NumFields(); i++ {
			if !isComparable(t.Field(i).Type()) {
				return false
			}
		}
		return true
	case *types.Array:
		return isComparable(t.Elem())
	}
	return false
}

// findFieldByVar finds the corresponding ast.Node for a given types.Var.
func findFieldByVar(pkg *packages.Package, v *types.Var) *ast.Field {
	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}
				for _, field := range structType.Fields.List {
					for _, fieldName := range field.Names {
						if pkg.Fset.Position(fieldName.Pos()).String() == pkg.Fset.Position(v.Pos()).String() {
							return field
						}
					}
				}
			}
		}
	}
	return nil
}

func parseNode(node any) (string, error) {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, token.NewFileSet(), node)
	if err != nil {
		return "", fmt.Errorf("printer.Fprint: %w", err)
	}
	return buf.String(), nil
}

// typeToExpr converts a types.Type to an ast.Expr.
func typeToExpr(t types.Type) ast.Expr {
	switch t := t.(type) {
	case *types.Basic:
		return &ast.Ident{Name: t.Name()}
	case *types.Array:
		return &ast.ArrayType{
			Len: &ast.BasicLit{Kind: token.INT, Value: fmt.Sprint(t.Len())},
			Elt: typeToExpr(t.Elem()),
		}
	case *types.Slice:
		return &ast.ArrayType{
			Elt: typeToExpr(t.Elem()),
		}
	case *types.Struct:
		fields := &ast.FieldList{}
		for i := 0; i < t.NumFields(); i++ {
			field := t.Field(i)
			fields.List = append(fields.List, &ast.Field{
				Names: []*ast.Ident{{Name: field.Name()}},
				Type:  typeToExpr(field.Type()),
			})
		}
		return &ast.StructType{
			Fields: fields,
		}
	case *types.Pointer:
		return &ast.StarExpr{
			X: typeToExpr(t.Elem()),
		}
	case *types.Named:
		obj := t.Obj()
		if obj.Pkg() != nil {
			return &ast.SelectorExpr{
				X:   &ast.Ident{Name: obj.Pkg().Name()},
				Sel: &ast.Ident{Name: obj.Name()},
			}
		}
		return &ast.Ident{Name: obj.Name()}
	case *types.Alias:
		obj := t.Obj()
		if obj.Pkg() != nil {
			return &ast.SelectorExpr{
				X:   &ast.Ident{Name: obj.Pkg().Name()},
				Sel: &ast.Ident{Name: obj.Name()},
			}
		}
		return &ast.Ident{Name: obj.Name()}
		//return &ast.Ident{Name: "aa"}
	case *types.Interface:
		// Handle Interface{} as a special case
		if t.NumMethods() == 0 {
			return &ast.InterfaceType{
				Methods: &ast.FieldList{},
			}
		}
		methods := &ast.FieldList{}
		for i := 0; i < t.NumMethods(); i++ {
			method := t.Method(i)
			methods.List = append(methods.List, &ast.Field{
				Names: []*ast.Ident{{Name: method.Name()}},
				Type:  typeToExpr(method.Type()),
			})
		}
		return &ast.InterfaceType{
			Methods: methods,
		}
	case *types.Map:
		return &ast.MapType{
			Key:   typeToExpr(t.Key()),
			Value: typeToExpr(t.Elem()),
		}
	case *types.Chan:
		var dir ast.ChanDir
		switch t.Dir() {
		case types.SendRecv:
			dir = ast.SEND | ast.RECV
		case types.SendOnly:
			dir = ast.SEND
		case types.RecvOnly:
			dir = ast.RECV
		}
		return &ast.ChanType{
			Dir:   dir,
			Value: typeToExpr(t.Elem()),
		}
	case *types.Signature:
		params := &ast.FieldList{}
		for i := 0; i < t.Params().Len(); i++ {
			params.List = append(params.List, &ast.Field{
				Type: typeToExpr(t.Params().At(i).Type()),
			})
		}
		results := &ast.FieldList{}
		for i := 0; i < t.Results().Len(); i++ {
			results.List = append(results.List, &ast.Field{
				Type: typeToExpr(t.Results().At(i).Type()),
			})
		}
		return &ast.FuncType{
			Params:  params,
			Results: results,
		}
	default:
		return &ast.BadExpr{}
	}
}

// typeToString converts a types.Type to its string representation in Go code.
func typeToString(t types.Type, curPkg string) (string, error) {
	expr := typeToExpr(t)
	var sb strings.Builder
	fset := token.NewFileSet()
	if err := printer.Fprint(&sb, fset, expr); err != nil {
		return "", err
	}
	res := sb.String()

	var a string

	switch t.(type) {
	case *types.Pointer:
		a, _ = strings.CutPrefix(res, "*"+curPkg+".")
		if a != res {
			a = "*" + a
		}
	default:
		a, _ = strings.CutPrefix(res, curPkg+".")
	}
	return a, nil
}
