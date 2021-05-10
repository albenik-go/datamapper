package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/albenik-go/datamapper/codegen"
)

func Main(name string) int {
	const defaultTagName = "db"

	var (
		sourceArg       = flag.String("src", "", "Input file name. Required.")
		destinationArg  = flag.String("dst", "", "Output file. Defaults to stdout.")
		pkgNameArg      = flag.String("pkg", "", "Package of the generated code. (Required!)")
		typesArg        = flag.String("types", "", "Comma-separated list of types names.")
		excludeTypesArg = flag.String("types_exclude", "", "Comma-separated list of types names to exclude.")
		nameTagArg      = flag.String("nametag", defaultTagName, fmt.Sprintf("Struct tag key to get column name from. Defaults to %q.", defaultTagName))
		optsTagArg      = flag.String("optstag", defaultTagName, fmt.Sprintf("Struct tag key to get field options from. Defaults to %q.", defaultTagName))
	)
	flag.Parse()

	if *pkgNameArg == "" {
		flag.Usage()
		os.Exit(1)
	}

	dest := os.Stdout
	if *destinationArg != "" {
		var err error
		if dest, err = os.OpenFile(*destinationArg, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			return 1
		}
		defer dest.Close()
	}

	exclude := false
	var types []string
	if *typesArg != "" {
		types = strings.Split(*typesArg, ",")
	}

	if len(*excludeTypesArg) > 0 {
		if len(*typesArg) > 0 {
			fmt.Fprintln(os.Stderr, "-types & -types_exclude cannot be used together")
			return 2
		}
		types = strings.Split(*excludeTypesArg, ",")
		exclude = true
	}

	if len(*sourceArg) == 0 {
		*sourceArg = "."
	}

	if err := codegen.SimplifiedGenerate(*sourceArg, *pkgNameArg, *nameTagArg, *optsTagArg, types, exclude, dest); err != nil {
		fmt.Fprintln(os.Stderr, "Generation failed:", err)
		return 1
	}

	return 0
}
