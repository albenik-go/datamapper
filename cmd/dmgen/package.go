package main

import (
	"go/ast"
	"go/types"
)

type packageInfo struct {
	name  string
	defs  map[*ast.Ident]types.Object
	files []*File
}
