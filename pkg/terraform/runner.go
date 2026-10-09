package terraform

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Runner handles execution of the terraform CLI binary.
type Runner struct {
	BinaryPath string
}

// NewRunner creates a new Runner, attempting to find `terraform` in PATH.
func NewRunner() *Runner {
	path, err := exec.LookPath("terraform")
	if err != nil {
		path = "terraform" // fallback
	}
	return &Runner{BinaryPath: path}
}

// IsInstalled returns true if the terraform binary can be found in PATH.
func (r *Runner) IsInstalled() bool {
	_, err := exec.LookPath(r.BinaryPath)
	return err == nil
}

// RunOutputs executes `terraform output -json` in the specified directory safely.
func (r *Runner) RunOutputs(dir string) ([]byte, error) {
	if !r.IsInstalled() {
		return nil, errors.New("terraform binary not found in PATH. Please install HashiCorp Terraform to use 'tfout outputs'")
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("invalid working directory %q: %w", dir, err)
	}

	info, err := os.Stat(absDir)
	if err != nil || !info.IsDir() {
		return nil, fmt.Errorf("working directory does not exist or is not a directory: %s", absDir)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Execute `terraform output -json` directly without shell wrapper to prevent command injection
	cmd := exec.CommandContext(ctx, r.BinaryPath, "output", "-json")
	cmd.Dir = absDir

	out, err := cmd.CombinedOutput()
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("terraform command timed out after 30 seconds")
		}
		return nil, fmt.Errorf("terraform output failed: %w\nOutput: %s", err, string(out))
	}

	return out, nil
}
