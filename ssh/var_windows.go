// +build windows

package ssh

import (
	"os"
	"syscall"
)

// Stdout?
var Fd int = int(os.Stdin.Fd())

// SIGWINCH not supported on windows
// TODO: Find out a way to get resize event
const ResizeEvent syscall.Signal = syscall.SIGBUS
