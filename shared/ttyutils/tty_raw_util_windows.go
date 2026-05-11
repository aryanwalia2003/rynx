//go:build windows
package ttyutils

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode = kernel32.NewProc("SetConsoleMode")
)

type tty_termios_struct uint32

func Tty_raw_enable_util() (*tty_termios_struct, error) {
	h := syscall.Handle(os.Stdin.Fd())
	var mode uint32
	ret, _, err := procGetConsoleMode.Call(uintptr(h), uintptr(unsafe.Pointer(&mode)))
	if ret == 0 {
		return nil, err
	}

	// Disable ENABLE_ECHO_INPUT (0x4), ENABLE_LINE_INPUT (0x2), ENABLE_PROCESSED_INPUT (0x1)
	raw := mode &^ (uint32(0x0001) | uint32(0x0002) | uint32(0x0004))
	ret, _, err = procSetConsoleMode.Call(uintptr(h), uintptr(raw))
	if ret == 0 {
		return nil, err
	}

	oldMode := tty_termios_struct(mode)
	return &oldMode, nil
}

func Tty_raw_restore_util(old *tty_termios_struct) {
	if old == nil {
		return
	}
	h := syscall.Handle(os.Stdin.Fd())
	_, _, _ = procSetConsoleMode.Call(uintptr(h), uintptr(*old))
}
