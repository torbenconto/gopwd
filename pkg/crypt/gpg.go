package crypt

import (
	"bytes"
	"os"
	"os/exec"
)

var (
	defaultArgs = []string{"--quiet", "--yes", "--compress-algo=none", "--no-encrypt-to", "--no-auto-check-trustdb"}
)

type GPG struct {
	id         string
	binaryPath string
	args       []string
}

type Config struct {
	BinaryPath string
	Args       []string
}

func NewGPG(id string, config Config) *GPG {
	gpg := &GPG{
		id: id,
	}

	// Set default values
	if config.BinaryPath == "" {
		config.BinaryPath = "gpg"
	}
	if len(config.Args) == 0 {
		config.Args = defaultArgs
	}

	// Set config values
	gpg.binaryPath = config.BinaryPath
	gpg.args = config.Args

	return gpg
}

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

func (g *GPG) Decrypt(ciphertext []byte) ([]byte, error) {
	args := append(g.Args(), "--decrypt")

	buffer := &bytes.Buffer{}

	cmd := exec.Command(g.Binary(), args...)
	cmd.Stdin = bytes.NewReader(ciphertext)
	cmd.Stdout = buffer
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	return buffer.Bytes(), err
}

func (g *GPG) Binary() string {
	return g.binaryPath
}

func (g *GPG) ID() string {
	return g.id
}

func (g *GPG) Args() []string {
	return g.args
}
