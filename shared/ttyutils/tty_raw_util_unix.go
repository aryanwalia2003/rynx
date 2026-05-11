//go:build !windows
package ttyutils

import (
	"syscall"
	"unsafe"
)

type tty_termios_struct struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Line   uint8
	Cc     [32]uint8
	_      [3]byte
	Ispeed uint32
	Ospeed uint32
}

func Tty_raw_enable_util() (*tty_termios_struct, error) {
	var old tty_termios_struct
	// TCGETS = 0x5401
	if _, _, e := syscall.Syscall6(syscall.SYS_IOCTL, 0, 0x5401, uintptr(unsafe.Pointer(&old)), 0, 0, 0); e != 0 {
		return nil, e
	}
	raw := old
	raw.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.ISIG | syscall.IEXTEN
	raw.Iflag &^= syscall.IXON | syscall.ICRNL | syscall.BRKINT | syscall.INPCK | syscall.ISTRIP
	raw.Cflag |= syscall.CS8
	raw.Cc[6] = 1 // VMIN
	raw.Cc[5] = 0 // VTIME
	// TCSETS = 0x5402
	if _, _, e := syscall.Syscall6(syscall.SYS_IOCTL, 0, 0x5402, uintptr(unsafe.Pointer(&raw)), 0, 0, 0); e != 0 {
		return nil, e
	}
	return &old, nil
}

func Tty_raw_restore_util(old *tty_termios_struct) {
	if old == nil {
		return
	}
	// TCSETS = 0x5402
	syscall.Syscall6(syscall.SYS_IOCTL, 0, 0x5402, uintptr(unsafe.Pointer(old)), 0, 0, 0)
}
