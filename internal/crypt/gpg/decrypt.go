package gpg

import (
	"bytes"
	"os"
	"os/exec"
)

func (g *GPG) Decrypt(ciphertext []byte) ([]byte, error) {
	args := append(g.Args(), "--decrypt")

	cmd := exec.Command(g.Binary(), args...)
	cmd.Stdin = bytes.NewReader(ciphertext)
	cmd.Stderr = os.Stderr

	return cmd.Output()
}
