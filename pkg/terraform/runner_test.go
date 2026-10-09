package terraform_test

import (
	"testing"

	"github.com/tfout/tfout/pkg/terraform"
)

func TestRunnerNonExistentDir(t *testing.T) {
	runner := terraform.NewRunner()
	_, err := runner.RunOutputs("/path/that/does/not/exist/9999")
	if err == nil {
		t.Error("Expected error when running in non-existent directory")
	}
}

func TestRunnerIsInstalled(t *testing.T) {
	runner := terraform.NewRunner()
	// Should return boolean without panic
	_ = runner.IsInstalled()
}
