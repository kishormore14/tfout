package cli

import (
	"flag"
	"fmt"
	"io"

	"github.com/tfout/tfout/pkg/formatter"
	"github.com/tfout/tfout/pkg/terraform"
)

// RunOutputs executes the `tfout outputs` command.
func RunOutputs(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("outputs", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var dir string
	var globalFlags GlobalFlags

	fs.StringVar(&dir, "dir", ".", "Terraform working directory")
	fs.StringVar(&dir, "C", ".", "Terraform working directory (short flag)")
	fs.StringVar(&dir, "chdir", ".", "Terraform working directory")
	AddGlobalFlags(fs, &globalFlags)

	if err := fs.Parse(args); err != nil {
		return 1
	}

	cfg, err := globalFlags.ToFormatterConfig()
	if err != nil {
		fmt.Fprintf(stderr, "Error: %v\n", err)
		return 1
	}

	runner := terraform.NewRunner()
	rawJSON, err := runner.RunOutputs(dir)
	if err != nil {
		fmt.Fprintf(stderr, "Error reading terraform outputs: %v\n", err)
		return 1
	}

	outputsMap, err := terraform.ParseOutputs(rawJSON)
	if err != nil {
		fmt.Fprintf(stderr, "Error parsing terraform output JSON: %v\n", err)
		return 1
	}

	bannerData := outputsMap.ConvertToBanner(dir, cfg)
	output := formatter.RenderBanner(bannerData, cfg)
	fmt.Fprint(stdout, output)
	return 0
}
