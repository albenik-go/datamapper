package codegen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/albenik-go/datamapper/codegen/tag"
	"github.com/albenik-go/datamapper/codegen/template"
)

// SimplifiedGenerate generates mapper go code
// `name_tag:"column_name" opt_tag:",auto"` or `tag:"column_name,auto"` if tags are equal.
func SimplifiedGenerate(filename, pkg, nameTag, optTag string, types []string, exclude bool, out io.Writer) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.AllErrors)
	if err != nil {
		return errors.Wrapf(err, "cannot parse source file %q", filename)
	}

	var dmgenVer string
	if buildinfo, ok := debug.ReadBuildInfo(); ok {
		dmgenVer = buildinfo.Main.Version
	} else {
		dmgenVer = "(NO VERSION INFO)"
	}

	buf := bytes.NewBuffer(nil)
	if err = template.WriteHeader(buf, &template.Header{Package: pkg, DmgenVersion: dmgenVer}); err != nil {
		return err
	}

	var models []*template.ModelInfo
	ast.Inspect(f, walkFunc(types, exclude, nameTag, optTag, &models))
	for _, m := range models {
		if err := template.WriteModel(buf, m); err != nil {
			return err
		}
	}

	var formattedCode []byte
	if formattedCode, err = format.Source(buf.Bytes()); err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		lines := strings.Split(buf.String(), "\n")
		for n, l := range lines {
			fmt.Printf("%03d: %s\n", n+1, l) //nolint:forbidigo
		}

		if _, copyErr := io.Copy(out, buf); copyErr != nil {
			err = multierr.Append(err, copyErr)
		}

		return errors.Wrap(err, "generated code format error")
	}

	_, err = io.Copy(out, bytes.NewReader(formattedCode))
	return err
}

func walkFunc(types []string, exclude bool, nameTag, optTag string, models *[]*template.ModelInfo) func(ast.Node) bool {
	return func(node ast.Node) bool {
		decl, ok := node.(*ast.GenDecl)
		if !ok || decl.Tok != token.TYPE {
			return true
		}

		for _, spec := range decl.Specs {
			typeName, typeInfo, ok := getStructType(spec, exclude, types)
			if !ok {
				continue
			}

			info := &template.ModelInfo{
				EntityType:   typeName,
				SelectFields: make([]*template.FieldInfo, 0, len(typeInfo.Fields.List)),
				InsertFields: make([]*template.FieldInfo, 0, len(typeInfo.Fields.List)),
				UpdateFields: make([]*template.FieldInfo, 0, len(typeInfo.Fields.List)),
			}
			collectModelInfo(info, typeInfo, nameTag, optTag)
			*models = append(*models, info)
		}

		return false
	}
}

func collectModelInfo(info *template.ModelInfo, expr *ast.StructType, nameTag, optTag string) { //nolint:cyclop,gocognit
	for _, field := range expr.Fields.List {
		// processing embedded types if any
		var (
			indent *ast.Ident
			obj    *ast.Object
		)
		switch ftype := field.Type.(type) {
		case *ast.Ident:
			indent = ftype
			obj = ftype.Obj
		case *ast.SelectorExpr:
			if xIdent, ok := ftype.X.(*ast.Ident); ok {
				indent = xIdent
				obj = xIdent.Obj
			}
		case *ast.StarExpr:
			if xIdent, ok := ftype.X.(*ast.Ident); ok {
				indent = xIdent
				obj = xIdent.Obj
			}
		}

		if indent == nil {
			panic(errors.Errorf("unexpected field type descriptor %T", field.Type))
		}

		if obj != nil && obj.Kind == ast.Typ {
			if _, stype, ok := getStructType(obj.Decl, false, nil); ok {
				collectModelInfo(info, stype, nameTag, optTag)
			}
			continue
		}

		if len(field.Names) == 0 || field.Names[0].Name == "" {
			// unexpected case
			panic(errors.Errorf("field error: %#v", field))
		}

		if field.Tag == nil {
			continue
		}
		structTag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
		tagStr, ok := structTag.Lookup(optTag)
		if !ok && optTag == nameTag {
			continue
		}

		tagOpts := tag.ParseOptions(tagStr)
		if len(tagOpts) == 0 {
			panic(errors.Errorf("invalid tag %q", structTag))
		}

		colName := tagOpts[0].Name
		if colName == "-" {
			continue
		}

		if nameTag != optTag {
			nameTagStr, ok := structTag.Lookup(nameTag)
			if !ok {
				panic(errors.Errorf("nameTag specified but not found for field: %#v", field))
			}

			opts := tag.ParseOptions(nameTagStr)
			if len(opts) == 0 || len(opts[0].Name) == 0 || opts[0].Name == "-" {
				panic(errors.Errorf("nameTag invalid for field: %#v", field))
			}

			colName = opts[0].Name
		}

		f := &template.FieldInfo{
			FieldName:  field.Names[0].Name,
			FieldType:  indent.Name,
			ColumnName: colName,
		}

		if tt := tagOpts.Lookup("wrap"); tt != nil {
			f.Wrapper = tt.Value
		}

		info.SelectFields = append(info.SelectFields, f)
		if tagOpts.Lookup("auto") != nil {
			info.AutoincrementField = f
		} else {
			info.InsertFields = append(info.InsertFields, f)
			info.UpdateFields = append(info.UpdateFields, f)
		}
	}
}

func getStructType(spec interface{}, exclude bool, filter []string) (string, *ast.StructType, bool) {
	var (
		s  *ast.TypeSpec
		t  *ast.StructType
		ok bool
	)

	if s, ok = spec.(*ast.TypeSpec); !ok {
		return "", nil, false
	}

	found := !exclude
	if len(filter) > 0 {
		found = !found
		for _, name := range filter {
			if s.Name.Name == name {
				found = !found
				break
			}
		}
	}
	if !found {
		return "", nil, false
	}

	if t, ok = s.Type.(*ast.StructType); !ok || t.Incomplete {
		return "", nil, false
	}

	return s.Name.Name, t, true
}
