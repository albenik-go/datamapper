package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/albenik-go/datamapper/codegen"
)

func Main() int {
	const defaultTagName = "db"

	var (
		sourceArg      = flag.String("src", "", "Input file name. Required.")
		destinationArg = flag.String("dst", "", "Output file. Defaults to stdout.")
		pkgNameArg     = flag.String("pkg", "", "Package of the generated code. (Required!)")
		includeArg     = flag.String("include", "", "Comma-separated list of types names.")
		excludeArg     = flag.String("exclude", "", "Comma-separated list of types names to exclude.")
		nameTagArg     = flag.String("nametag", defaultTagName, fmt.Sprintf("Struct tag key to get column name from. Defaults to %q.", defaultTagName))
		optsTagArg     = flag.String("optstag", defaultTagName, fmt.Sprintf("Struct tag key to get field options from. Defaults to %q.", defaultTagName))
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

	typesStr := *excludeArg
	exclude := typesStr != ""

	if exclude {
		if *includeArg != "" {
			fmt.Fprintln(os.Stderr, "options -include & -exclude cannot be used together")
			return 2
		}
	} else {
		typesStr = *includeArg
	}

	var types []string
	if typesStr != "" {
		types = strings.Split(typesStr, ",")
	}

	if *sourceArg == "" {
		*sourceArg = "."
	}

	if err := codegen.SimplifiedGenerate(*sourceArg, *pkgNameArg, *nameTagArg, *optsTagArg, types, exclude, dest); err != nil {
		fmt.Fprintln(os.Stderr, "Generation failed:", err)
		return 1
	}

	return 0
}
