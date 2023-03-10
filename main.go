package main

import (
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s FILE LINE COL\n", os.Args[0])
		os.Exit(1)
	}

	pathArg := os.Args[1]
	lineArg, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid line number %q\n", os.Args[2])
		os.Exit(1)
	}

	colArg, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid column number %q\n", os.Args[3])
		os.Exit(1)
	}

	err = lookupAndPrintGoDef(pathArg, lineArg, colArg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func lookupAndPrintGoDef(path string, line int, col int) error {
	// Step 1: load the Go package
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	loadMode := (packages.NeedName |
		packages.NeedSyntax |
		packages.NeedDeps |
		packages.NeedTypes |
		packages.NeedTypesInfo)

	cfg := &packages.Config{
		Mode: loadMode,
		Dir:  filepath.Dir(absPath),
	}

	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		return err
	} else if len(pkgs) == 0 {
		return fmt.Errorf("No packages loaded")
	}

	pkg := pkgs[0]

	// Step 2: find the AST identifier
	var astFile *ast.File
	for _, f := range pkg.Syntax {
		if pkg.Fset.Position(f.Pos()).Filename == absPath {
			astFile = f
			break
		}
	}
	if astFile == nil {
		return fmt.Errorf("Could not find AST file for %q", absPath)
	}

	var astIdent *ast.Ident
	ast.Inspect(astFile, func(node ast.Node) bool {
		if node == nil || astIdent != nil {
			return false
		}
		start, end := pkg.Fset.Position(node.Pos()), pkg.Fset.Position(node.End())
		if line < start.Line ||
			line > end.Line ||
			(line == start.Line && col < start.Column) ||
			(line == end.Line && col > end.Column) {
			return false
		}
		if node, ok := node.(*ast.Ident); ok {
			astIdent = node
			return false
		}
		return true
	})
	if astIdent == nil {
		return fmt.Errorf("Could not find AST identifier at %s:%d:%d", path, line, col)
	}

	// Step 3: lookup the definition
	obj, ok := pkg.TypesInfo.Uses[astIdent]
	if !ok {
		obj = pkg.TypesInfo.Defs[astIdent]
	}
	if obj == nil {
		return fmt.Errorf("Could not find type object for ident %q at %s:%d:%d", astIdent.Name, path, line, col)
	} else if !obj.Pos().IsValid() {
		return fmt.Errorf("Invalid position for type object for %q at %s:%d:%d", astIdent.Name, path, line, col)
	}

	defPosition := pkg.Fset.Position(obj.Pos())
	fmt.Printf("%q is defined at %s:%d:%d\n", obj.Name(), normalizePath(defPosition.Filename), defPosition.Line, defPosition.Column)

	return nil
}

func normalizePath(p string) string {
	cwd, err := os.Getwd()
	if err != nil {
		return p
	}

	if !strings.HasPrefix(p, cwd) {
		return p
	}

	relPath, err := filepath.Rel(cwd, p)
	if err != nil {
		return p
	}

	return relPath
}
