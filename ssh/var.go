// +build !windows

package ssh

import (
	"os"
	"syscall"
)

// Fd to get terminal size
var Fd int = int(os.Stdin.Fd())

// Terminal resize event
const ResizeEvent syscall.Signal = syscall.SIGWINCH
