package cli

import (
	"fmt"
	"io"
	"os"
)

// PrintUsage prints help text to output.
func PrintUsage(w io.Writer) {
	usage := `tfout — A Reusable Open-Source Terminal Formatter for HashiCorp Terraform

Usage:
  tfout <command> [flags]

Available Commands:
  banner    Print a formatted status banner
  render    Render a structured JSON document into a terminal presentation
  outputs   Read and format Terraform outputs from working directory
  doctor    Check system readiness and environment diagnostics
  version   Display version and build information

Global Flags:
  --color         Color mode: auto, always, never (default "auto")
  --theme         Theme style: default, minimal, mono (default "default")
  --border        Border style: single, double, none (default "single")
  --show-sensitive Reveal sensitive values with a warning badge

Examples:
  # Print a Route 53 DNS banner
  tfout banner --title "ROUTE53 HOSTED ZONE CREATED" \
    --status success \
    --detail "Domain" "example.com" \
    --detail "Name Servers" "ns-123.example.net, ns-456.example.net"

  # Render a structured JSON document
  tfout render --file summary.json
  cat summary.json | tfout render

  # Format Terraform outputs
  tfout outputs --dir ./infrastructure

Run 'tfout <command> --help' for command-specific flags.
`
	fmt.Fprint(w, usage)
}

// Execute parses command line arguments and routes execution to the appropriate command.
func Execute(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		PrintUsage(stderr)
		return 1
	}

	subCmd := args[0]
	subArgs := args[1:]

	switch subCmd {
	case "banner":
		return RunBanner(subArgs, stdout, stderr)
	case "render":
		return RunRender(subArgs, stdin, stdout, stderr)
	case "outputs":
		return RunOutputs(subArgs, stdout, stderr)
	case "version", "-v", "--version":
		return RunVersion(stdout)
	case "doctor":
		return RunDoctor(stdout)
	case "help", "-h", "--help":
		PrintUsage(stdout)
		return 0
	default:
		fmt.Fprintf(stderr, "Unknown command %q\n\n", subCmd)
		PrintUsage(stderr)
		return 1
	}
}

// Main is the standard entrypoint wrapper for os.Args.
func Main() {
	exitCode := Execute(os.Args[1:], os.Stdin, os.Stdout, os.Stderr)
	os.Exit(exitCode)
}
