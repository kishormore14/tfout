//go:build windows

package cli

import (
	"syscall"
	"unsafe"
)

func enableWindowsVT() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	procCreateFile := kernel32.NewProc("CreateFileW")
	procGetConsoleMode := kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode := kernel32.NewProc("SetConsoleMode")
	procCloseHandle := kernel32.NewProc("CloseHandle")

	// Open handle to active Windows Console Output screen buffer (CONOUT$)
	conout, err := syscall.UTF16PtrFromString("CONOUT$")
	if err != nil {
		return
	}

	handle, _, _ := procCreateFile.Call(
		uintptr(unsafe.Pointer(conout)),
		uintptr(syscall.GENERIC_READ|syscall.GENERIC_WRITE),
		uintptr(syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE),
		0,
		uintptr(syscall.OPEN_EXISTING),
		0,
		0,
	)

	if handle != uintptr(syscall.InvalidHandle) && handle != 0 {
		var mode uint32
		ret, _, _ := procGetConsoleMode.Call(handle, uintptr(unsafe.Pointer(&mode)))
		if ret != 0 {
			const ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004
			mode |= ENABLE_VIRTUAL_TERMINAL_PROCESSING
			procSetConsoleMode.Call(handle, uintptr(mode))
		}
		procCloseHandle.Call(handle)
	}
}
