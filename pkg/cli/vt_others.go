//go:build !windows

package cli

func enableWindowsVT() {
	// No-op on non-Windows platforms
}
