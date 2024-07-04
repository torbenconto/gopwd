//go:build darwin

package gpg

import "os"

func TTY() string {
	dest, err := os.Readlink("/dev/fd/0")
	if err != nil {
		return ""
	}

	return dest
}
