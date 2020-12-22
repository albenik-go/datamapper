package main

import (
	"go/ast"
	"go/token"
	"log"
	"reflect"
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
		tspec, texpr, ok := getSpecType(spec)
		if !ok || tspec.Name.Name != f.typeName {
			continue
		}

		info := &modelInfo{
			ModelName:    tspec.Name.Name,
			SelectFields: make([]*fieldInfo, 0, len(texpr.Fields.List)),
			InsertFields: make([]*fieldInfo, 0, len(texpr.Fields.List)),
			UpdateFields: make([]*fieldInfo, 0, len(texpr.Fields.List)),
		}

		f.processTypeSpec(texpr, info)

		f.models = append(f.models, info)
	}

	return false
}

func (f *File) processTypeSpec(expr *ast.StructType, info *modelInfo) {
	for _, field := range expr.Fields.List {
		if ident, ok := field.Type.(*ast.Ident); ok && ident.Obj != nil && ident.Obj.Kind == ast.Typ {
			_, texpr, ok := getSpecType(ident.Obj.Decl)
			if !ok {
				continue
			}
			f.processTypeSpec(texpr, info)
			continue
		}

		if len(field.Names) == 0 || field.Names[0].Name == "" {
			log.Fatalf("Field error: %#v", field)
		}

		tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
		optstr, ok := tag.Lookup("map")
		if !ok {
			continue
		}

		opts := strings.Split(optstr, ",")
		if len(opts) == 0 {
			log.Fatalf("Invalid tag %q", tag)
		}

		f := &fieldInfo{FieldName: field.Names[0].Name, ColumnName: opts[0]}
		info.SelectFields = append(info.SelectFields, f)
		if !tagOptionExists(opts[1:], "auto") {
			info.InsertFields = append(info.InsertFields, f)
			info.UpdateFields = append(info.UpdateFields, f)
		}
	}
}

func getSpecType(spec interface{}) (s *ast.TypeSpec, t *ast.StructType, ok bool) {
	if s, ok = spec.(*ast.TypeSpec); !ok {
		return nil, nil, false
	}

	if t, ok = s.Type.(*ast.StructType); !ok || t.Incomplete {
		return nil, nil, false
	}

	return
}

func tagOptionExists(opts []string, name string) bool {
	for _, o := range opts {
		if o == name {
			return true
		}
	}
	return false
}
