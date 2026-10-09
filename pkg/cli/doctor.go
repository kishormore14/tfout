package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/tfout/tfout/pkg/formatter"
	"github.com/tfout/tfout/pkg/terraform"
)

// RunDoctor executes the `tfout doctor` command.
func RunDoctor(stdout io.Writer) int {
	var details []formatter.DetailPair
	var items []string

	details = append(details, formatter.DetailPair{
		Key:   "tfout Core Executable",
		Value: fmt.Sprintf("OK (version %s)", Version),
	})

	details = append(details, formatter.DetailPair{
		Key:   "OS / Architecture",
		Value: fmt.Sprintf("%s / %s", runtime.GOOS, runtime.GOARCH),
	})

	details = append(details, formatter.DetailPair{
		Key:   "Go Runtime",
		Value: runtime.Version(),
	})

	// Environment checks
	noColorVal := os.Getenv("NO_COLOR")
	if noColorVal != "" {
		details = append(details, formatter.DetailPair{
			Key:   "NO_COLOR Env",
			Value: fmt.Sprintf("Set (%s) - Colors disabled by default", noColorVal),
		})
	} else {
		details = append(details, formatter.DetailPair{
			Key:   "NO_COLOR Env",
			Value: "Not set",
		})
	}

	termVal := os.Getenv("TERM")
	if termVal == "" {
		termVal = "Unset"
	}
	details = append(details, formatter.DetailPair{
		Key:   "TERM Env",
		Value: termVal,
	})

	// Optional Terraform Dependency Check
	runner := terraform.NewRunner()
	if runner.IsInstalled() {
		out, err := exec.Command(runner.BinaryPath, "version").Output()
		tfVer := "Installed"
		if err == nil {
			firstLine := strings.Split(string(out), "\n")[0]
			tfVer = strings.TrimSpace(firstLine)
		}
		details = append(details, formatter.DetailPair{
			Key:   "HashiCorp Terraform (Optional)",
			Value: fmt.Sprintf("Detected: %s", tfVer),
		})
		items = append(items, "Terraform integration ('tfout outputs') is fully available.")
	} else {
		details = append(details, formatter.DetailPair{
			Key:   "HashiCorp Terraform (Optional)",
			Value: "Not found in PATH",
		})
		items = append(items, "Core features ('tfout banner', 'tfout render') are fully operational.")
		items = append(items, "'tfout outputs' requires HashiCorp Terraform binary installed on PATH.")
	}

	bannerData := formatter.BannerData{
		Title:   "TFOUT SYSTEM DOCTOR",
		Status:  formatter.StatusSuccess,
		Message: "System readiness check completed.",
		Details: details,
		Items:   items,
	}

	cfg := formatter.Config{
		ColorMode:   formatter.ColorAuto,
		Theme:       formatter.ThemeDefault,
		BorderStyle: formatter.BorderSingle,
		TermWidth:   80,
	}

	output := formatter.RenderBanner(bannerData, cfg)
	fmt.Fprint(stdout, output)
	return 0
}
