//go:build linux

package editor

import "os"

func Editor() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	return "vi"
}
