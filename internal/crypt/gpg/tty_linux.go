//go:build linux

package gpg

import "os"

var fd0 = "/proc/self/fd/0"

func TTY() string {
	dest, err := os.Readlink(fd0)
	if err != nil {
		return ""
	}

	return dest
}
