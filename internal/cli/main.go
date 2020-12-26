package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/albenik-go/datamapper/codegen"
)

func Main(name string) int {
	const defaultTagName = "col"

	var (
		sourceArg      = flag.String("source", "", "Input file name. Required.")
		destinationArg = flag.String("destination", "", "Output file. Defaults to stdout.")
		packageArg     = flag.String("package", "", "Package of the generated code. Defaults to the package of the input.")
		typesArg       = flag.String("types", "", "Comma-separated list of type names. Required.")
		buildTagsArgs  = flag.String("tags", "", "Comma-separated list of build tags to apply")
		nameTagArg     = flag.String("name_tag", defaultTagName, fmt.Sprintf("Struct tag key to get column name from. Defaults to %q.", defaultTagName))
		optionsTagArg  = flag.String("options_tag", defaultTagName, fmt.Sprintf("Struct tag key to get field options from. Defaults to %q.", defaultTagName))
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", name)
		fmt.Fprintf(os.Stderr, "\t%s [flags] -type T [directory]\n", name)
		fmt.Fprintf(os.Stderr, "\t%s [flags] -type T files... # Must be a single package\n", name)
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	dest := os.Stdout
	if len(*destinationArg) > 0 {
		var err error
		if dest, err = os.OpenFile(*destinationArg, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			return 1
		}
		defer dest.Close()
	}

	var types []string
	if len(*typesArg) > 0 {
		types = strings.Split(*typesArg, ",")
	}

	var tags []string
	if len(*buildTagsArgs) > 0 {
		tags = strings.Split(*buildTagsArgs, ",")
	}

	if len(*sourceArg) == 0 {
		*sourceArg = "."
	}

	if err := codegen.Generate(*packageArg, *sourceArg, tags, types, *nameTagArg, *optionsTagArg, dest); err != nil {
		fmt.Fprintln(os.Stderr, "Generation failed:", err)
		return 1
	}

	return 0
}
