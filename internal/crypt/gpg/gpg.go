package gpg

import (
	"os"
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

	// Set GPG_TTY environment variable
	if gt := os.Getenv("GPG_TTY"); gt == "" {
		if t := TTY(); t != "" {
			_ = os.Setenv("GPG_TTY", t)
		}
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

func (g *GPG) Binary() string {
	return g.binaryPath
}

func (g *GPG) ID() string {
	return g.id
}

func (g *GPG) Args() []string {
	return g.args
}
