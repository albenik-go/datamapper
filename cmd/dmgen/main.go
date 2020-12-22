package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	typeNames = flag.String("type", "", "comma-separated list of type names; must be set")
	output    = flag.String("output", "", "output file name; default srcdir/<type>_mapper.go")
	buildTags = flag.String("tags", "", "comma-separated list of build tags to apply")
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("dmgen: ")

	flag.Usage = usage
	flag.Parse()

	if len(*typeNames) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	typesList := strings.Split(*typeNames, ",")

	var tags []string
	if len(*buildTags) > 0 {
		tags = strings.Split(*buildTags, ",")
	}

	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}

	var dir string
	if len(args) == 1 && isDir(args[0]) {
		dir = args[0]
	} else {
		if len(tags) != 0 {
			log.Fatal("-tags option applies only to directories, not when files are specified")
		}
		dir = filepath.Dir(args[0])
	}

	g := generator{}
	g.parsePackage(args, tags)
	g.writeHeader(strings.Join(os.Args[1:], " "), g.pkg.name)
	g.generate(typesList)
	src := g.format()

	outputName := *output
	if outputName == "" {
		baseName := fmt.Sprintf("%s_mapper.go", typesList[0])
		outputName = filepath.Join(dir, strings.ToLower(baseName))
	}
	err := ioutil.WriteFile(outputName, src, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

// usage is a replacement usage function for the flags package.
func usage() {
	fmt.Fprintf(os.Stderr, "Usage of dmgen:\n")
	fmt.Fprintf(os.Stderr, "\tdmgen [flags] -type T [directory]\n")
	fmt.Fprintf(os.Stderr, "\tdmgen [flags] -type T files... # Must be a single package\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

// isDir reports whether the named file is a directory.
func isDir(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}
