package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"strings"
)

type File struct {
	pkg      *packageInfo
	file     *ast.File
	typeName string

	models []*modelInfo
}

func (f *File) genDecl(node ast.Node) bool {
	decl, ok := node.(*ast.GenDecl)
	if !ok || decl.Tok != token.TYPE {
		return true
	}

	for _, spec := range decl.Specs {
		tspec := spec.(*ast.TypeSpec) // Guaranteed to succeed as this is TYPE.
		if tspec.Name.Name != f.typeName {
			continue
		}

		texpr, ok := tspec.Type.(*ast.StructType)
		if !ok || texpr.Incomplete {
			continue
		}

		info := &modelInfo{
			ModelName:    tspec.Name.Name,
			SelectFields: make([]*fieldInfo, 0, len(texpr.Fields.List)),
			InsertFields: make([]*fieldInfo, 0, len(texpr.Fields.List)),
			UpdateFields: make([]*fieldInfo, 0, len(texpr.Fields.List)),
		}

		for _, field := range texpr.Fields.List {
			if len(field.Names) == 0 || field.Names[0].Name == "" {
				log.Fatalf("Field error: %#v", field)
			}

			tags, err := parseTags(strings.Trim(field.Tag.Value, "`"))
			if err != nil {
				log.Fatal("Tag parse error:", err)
			}

			tag, ok := tags["map"]
			if !ok {
				continue
			}

			fmt.Println("TAG:", tag.Name)

			info.SelectFields = append(info.SelectFields, &fieldInfo{
				FieldName:  field.Names[0].Name,
				ColumnName: tag.Name,
			})

			if !tag.HasOption("autogen") {
				info.InsertFields = append(info.InsertFields, &fieldInfo{
					FieldName:  field.Names[0].Name,
					ColumnName: tag.Name,
				})

				info.UpdateFields = append(info.UpdateFields, &fieldInfo{
					FieldName:  field.Names[0].Name,
					ColumnName: tag.Name,
				})
			}
		}

		f.models = append(f.models, info)
	}

	return false
}
