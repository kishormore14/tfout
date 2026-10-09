package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/tfout/tfout/pkg/formatter"
	"github.com/tfout/tfout/pkg/schema"
)

// RunRender executes the `tfout render` command.
func RunRender(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("render", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var filePath string
	var globalFlags GlobalFlags

	fs.StringVar(&filePath, "file", "", "Path to JSON file to render (reads from stdin if omitted)")
	AddGlobalFlags(fs, &globalFlags)

	if err := fs.Parse(args); err != nil {
		return 1
	}

	var inputReader io.Reader
	if filePath != "" {
		f, err := os.Open(filePath)
		if err != nil {
			fmt.Fprintf(stderr, "Error opening JSON file %q: %v\n", filePath, err)
			return 1
		}
		defer f.Close()
		inputReader = f
	} else {
		inputReader = stdin
	}

	doc, err := schema.ParseAndValidate(inputReader)
	if err != nil {
		fmt.Fprintf(stderr, "Error rendering JSON: %v\n", err)
		return 1
	}

	cfg, err := globalFlags.ToFormatterConfig()
	if err != nil {
		fmt.Fprintf(stderr, "Error: %v\n", err)
		return 1
	}

	output := formatter.RenderDocumentToBanner(doc, cfg)
	fmt.Fprint(stdout, output)
	return 0
}
