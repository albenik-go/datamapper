package codegen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io"
	"reflect"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/albenik-go/datamapper/codegen/tag"
	"github.com/albenik-go/datamapper/codegen/template"
)

// `json:"column_name" col:",auto"` if nameTag=="json"
// `col:"column_name,auto"`
func Generate(targetPkg, srcName string, tags, types []string, exclude bool, nameTag, optTag string, out io.Writer) error {
	buf := bytes.NewBuffer(nil)

	conf := &packages.Config{
		Mode:  packages.NeedName | packages.NeedFiles | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo,
		Tests: false,
	}
	if len(tags) > 0 {
		conf.BuildFlags = []string{"-tags " + strings.Join(tags, ",")}
	}
	pkgs, err := packages.Load(conf, srcName)
	if err != nil {
		return err
	}
	if len(pkgs) != 1 {
		return fmt.Errorf("error: %d packages found", len(pkgs))
	}

	pkg := pkgs[0]

	if len(targetPkg) == 0 {
		targetPkg = pkg.Name
	}

	if err := template.WriteHeader(buf, &template.Header{Package: targetPkg}); err != nil {
		return err
	}

	for _, f := range pkg.Syntax {
		var models []*template.ModelInfo
		ast.Inspect(f, walkFunc(types, exclude, nameTag, optTag, &models))
		for _, m := range models {
			if err := template.WriteModel(buf, m); err != nil {
				return err
			}
		}
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		fmt.Println(string(buf.Bytes()))
		return err
	}

	_, err = io.Copy(out, bytes.NewReader(src))
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
				ModelName:    typeName,
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

func collectModelInfo(info *template.ModelInfo, expr *ast.StructType, nameTag, optTag string) {
	for _, field := range expr.Fields.List {
		// processing embeded types if any
		if ident, ok := field.Type.(*ast.Ident); ok && ident.Obj != nil && ident.Obj.Kind == ast.Typ {
			var stype *ast.StructType
			if _, stype, ok = getStructType(ident.Obj.Decl, false, nil); !ok {
				continue
			}
			collectModelInfo(info, stype, nameTag, optTag)
			continue
		}

		if len(field.Names) == 0 || field.Names[0].Name == "" {
			// unexpected case
			panic(fmt.Errorf("field error: %#v", field))
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
			panic(fmt.Errorf("invalid tag %q", structTag))
		}

		colName := tagOpts[0].Name
		if colName == "-" {
			continue
		}

		if nameTag != optTag {
			nameTagStr, ok := structTag.Lookup(nameTag)
			if !ok {
				panic(fmt.Errorf("nameTag specified but not found for field: %#v", field))
			}

			opts := tag.ParseOptions(nameTagStr)
			if len(opts) == 0 || len(opts[0].Name) == 0 || opts[0].Name == "-" {
				panic(fmt.Errorf("nameTag invalid for field: %#v", field))
			}

			colName = opts[0].Name
		}

		f := &template.FieldInfo{
			FieldName:  field.Names[0].Name,
			ColumnName: colName,
		}

		if tagOpts.Lookup("nullable") != nil {
			f.Wrappers = append(f.Wrappers, "Nullable")
		}

		if tt := tagOpts.Lookup("type"); tt != nil {
			switch tt.Value {
			case "int":
				switch t := field.Type.(type) {
				case *ast.Ident:
					switch t.Name {
					case "bool":
						f.Wrappers = append(f.Wrappers, "IntBool")
					}
				}
			}
		}

		info.SelectFields = append(info.SelectFields, f)
		if tagOpts.Lookup("auto") == nil {
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
