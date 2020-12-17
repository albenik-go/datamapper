package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"log"
	"strings"

	"golang.org/x/tools/go/packages"
)

type generator struct {
	buf bytes.Buffer // Accumulated output.
	pkg *packageInfo // pkg we are scanning.
}

func (g *generator) printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

func (g *generator) parsePackage(patterns []string, tags []string) {
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedImports |
			packages.NeedTypes |
			packages.NeedTypesSizes |
			packages.NeedSyntax |
			packages.NeedTypesInfo,
		// TODO: Need to think about constants in test files. Maybe write type_string_test.go
		// in a separate pass? For later.
		Tests:      false,
		BuildFlags: []string{fmt.Sprintf("-tags=%s", strings.Join(tags, " "))},
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages found", len(pkgs))
	}
	g.addPackage(pkgs[0])
}

func (g *generator) addPackage(pkg *packages.Package) {
	g.pkg = &packageInfo{
		name:  pkg.Name,
		defs:  pkg.TypesInfo.Defs,
		files: make([]*File, len(pkg.Syntax)),
	}

	for i, file := range pkg.Syntax {
		g.pkg.files[i] = &File{
			file: file,
			pkg:  g.pkg,
		}
	}
}

func (g *generator) generate(typeName string) {
	for _, file := range g.pkg.files {
		file.typeName = typeName
		if file.file != nil {
			ast.Inspect(file.file, file.genDecl)
			if len(file.models) > 0 {
				execTemplate(&g.buf, file.models[0])
			}
		}
	}
}

func (g *generator) format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return g.buf.Bytes()
	}
	return src
}
