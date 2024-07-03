package gpg

import (
	"bytes"
	"os"
	"os/exec"
)

func (g *GPG) Encrypt(plaintext []byte) ([]byte, error) {
	args := append(g.Args(), "--encrypt")
	args = append(args, "--recipient", g.ID())

	buffer := &bytes.Buffer{}

	cmd := exec.Command(g.Binary(), args...)
	cmd.Stdin = bytes.NewReader(plaintext)
	cmd.Stdout = buffer
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	return buffer.Bytes(), err
}
