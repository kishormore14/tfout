package cli

import (
	"flag"
	"fmt"
	"io"

	"github.com/tfout/tfout/pkg/formatter"
)

// RunBanner executes the `tfout banner` command.
func RunBanner(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("banner", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var title string
	var statusStr string
	var message string
	var details DetailFlag
	var items StringSliceFlag
	var globalFlags GlobalFlags

	fs.StringVar(&title, "title", "", "Banner title (required)")
	fs.StringVar(&statusStr, "status", "info", "Status: success, warning, error, info")
	fs.StringVar(&message, "message", "", "Optional description message")
	fs.Var(&details, "detail", "Key-value pair for banner detail (e.g. --detail \"Domain=example.com\" or --detail \"Domain\" \"example.com\")")
	fs.Var(&items, "item", "ListItem (repeatable flag)")
	AddGlobalFlags(fs, &globalFlags)

	if err := fs.Parse(args); err != nil {
		return 1
	}

	// Support positional pairs if --detail was passed with key and next arg is value
	positional := fs.Args()
	for i := 0; i < len(positional); i++ {
		if title == "" && i == 0 {
			title = positional[0]
			continue
		}
		// If detail had key without value, assign from positional
		for idx := range details {
			if details[idx].Value == "" && i < len(positional) {
				details[idx].Value = positional[i]
				i++
				break
			}
		}
	}

	if title == "" {
		fmt.Fprintln(stderr, "Error: missing required option --title (e.g., tfout banner --title \"ROUTE53 HOSTED ZONE CREATED\")")
		return 1
	}

	cfg, err := globalFlags.ToFormatterConfig()
	if err != nil {
		fmt.Fprintf(stderr, "Error: %v\n", err)
		return 1
	}

	status := formatter.Status(statusStr)
	bannerData := formatter.BannerData{
		Title:         title,
		Status:        status,
		Message:       message,
		Details:       []formatter.DetailPair(details),
		Items:         []string(items),
		SensitiveWarn: false,
	}

	output := formatter.RenderBanner(bannerData, cfg)
	fmt.Fprint(stdout, output)
	return 0
}
