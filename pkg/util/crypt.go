package util

import (
	"os"
	"strings"
)

func ReadGPGID(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	gpgIDBytes := make([]byte, 0)
	buffer := make([]byte, 4096)
	for {
		n, err := file.Read(buffer)
		if n > 0 {
			gpgIDBytes = append(gpgIDBytes, buffer[:n]...)
		}
		if err != nil {
			break
		}
	}

	return strings.TrimSpace(string(gpgIDBytes)), nil
}
