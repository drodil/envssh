// +build !windows

package ssh

import (
	"os"
	"syscall"
)

// Fd is file descriptor to get terminal size.
var Fd int = int(os.Stdin.Fd())

// ResizeEvent is syscall resize event.
const ResizeEvent syscall.Signal = syscall.SIGWINCH
