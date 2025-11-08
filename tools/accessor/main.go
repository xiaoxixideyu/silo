package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var (
	path   string
	getter bool
	setter bool
	update bool
	cast   string
	prefix string
	method string
	column string
	file   string

	rootCmd = &cobra.Command{
		Use:     "accessor [struct to parse]",
		Short:   "accessor",
		Example: "accessor User User2",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			generate(args)
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", ".", "set file to search, use filename without .go suffix")
	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", ".", "set path to search")
	rootCmd.PersistentFlags().BoolVarP(&getter, "getter", "g", false, "generate getter (default false)")
	rootCmd.PersistentFlags().BoolVarP(&setter, "setter", "s", true, "generate setter")
	rootCmd.PersistentFlags().BoolVarP(&update, "update", "u", true, "with update track")
	rootCmd.PersistentFlags().StringVarP(&method, "method", "m", "Update", "update method name, used when update track is true")
	rootCmd.PersistentFlags().StringVarP(&column, "column", "C", "Column", "column middle name")
	rootCmd.PersistentFlags().StringVarP(&cast, "cast", "c", CastTypeSnake, "camel or snake, used when update track is true")
	rootCmd.PersistentFlags().StringVar(&prefix, "prefix", "", "method prefix")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("accessor error", "error", err)
		os.Exit(1)
	}
}

func generate(target []string) {
	if !isDirectory(path) {
		log.Fatalf("path %v is not dir", path)
	}

	generators, err := NewGenerators(target, path, file)
	if err != nil {
		slog.Error("Failed to create generators", "error", err.Error())
	}

	for _, generator := range generators {
		slog.Info(fmt.Sprintf("generate %s ...\n", generator.FileName))
		err := generator.Generate(
			WithGetter(getter),
			WithSetter(setter),
			WithPrefix(prefix),
			WithUpdate(update),
			WithCast(cast),
			WithMethod(method),
			WithColumn(column),
		)
		if err != nil {
			log.Fatalf("Failed to generate, error: %s", err.Error())
		}
	}
}

// isDirectory reports whether the named file is a directory.
func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}
