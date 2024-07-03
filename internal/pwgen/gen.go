package pwgen

import (
	"crypto/rand"
	"fmt"
	internalRand "github.com/torbenconto/gopwd/internal/crypt/rand"
	"math/big"
	"strings"
)

var (
	Digits = "0123456789"
	Upper  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Lower  = "abcdefghijklmnopqrstuvwxyz"
	Syms   = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
)

type PasswordGenerator struct {
	config PasswordGeneratorConfig
}

type PasswordGeneratorConfig struct {
	Length    int
	Humanized bool
	Symbols   bool
	Numbers   bool
	Lowercase bool
	Uppercase bool
}

func NewPasswordGenerator(config PasswordGeneratorConfig) *PasswordGenerator {
	return &PasswordGenerator{
		config: config,
	}
}

func (pg *PasswordGenerator) Generate() (string, error) {
	if pg.config.Humanized {
		return pg.generateHumanized()
	} else {
		return pg.generateNonHumanized()
	}
}

func (pg *PasswordGenerator) generateHumanized() (string, error) {
	var password strings.Builder

	for password.Len() < pg.config.Length {
		// Select a random word
		wordIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(wordlist))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random word index: %w", err)
		}
		word := wordlist[wordIndex.Int64()]

		// Optionally capitalize the word
		if pg.config.Uppercase && internalRand.RandomInt(2) == 1 { // Simulate a 50% chance
			word = strings.Title(word)
		}

		// Append the word to the password
		if password.Len()+len(word) <= pg.config.Length {
			password.WriteString(word)
		} else {
			// If adding the whole word would exceed the length, stop adding words
			break
		}

		// Optionally insert a digit
		if pg.config.Numbers && password.Len() < pg.config.Length {
			digitIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(Digits))))
			if err != nil {
				return "", fmt.Errorf("failed to generate random digit index: %w", err)
			}
			password.WriteByte(Digits[digitIndex.Int64()])
		}

		// Optionally insert a symbol
		if pg.config.Symbols && password.Len() < pg.config.Length {
			symbolIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(Syms))))
			if err != nil {
				return "", fmt.Errorf("failed to generate random symbol index: %w", err)
			}
			password.WriteByte(Syms[symbolIndex.Int64()])
		}
	}

	// Trim the password if it exceeds the desired length
	if password.Len() > pg.config.Length {
		return password.String()[:pg.config.Length], nil
	}

	return password.String(), nil
}

func (pg *PasswordGenerator) generateNonHumanized() (string, error) {
	var characters string
	if pg.config.Symbols {
		characters += Syms
	}
	if pg.config.Numbers {
		characters += Digits
	}
	if pg.config.Lowercase {
		characters += Lower
	}
	if pg.config.Uppercase {
		characters += Upper
	}
	if !pg.config.Lowercase && !pg.config.Uppercase {
		characters += Upper + Lower // Add both lower and upper case letters if neither is specified
	}

	var password strings.Builder
	if len(characters) == 0 {
		return "", fmt.Errorf("no characters available for password generation")
	}

	for i := 0; i < pg.config.Length; i++ {
		randomBytes := make([]byte, 1)
		_, err := rand.Read(randomBytes)
		if err != nil {
			return "", fmt.Errorf("failed to generate random bytes: %w", err)
		}
		randomIndex := int(randomBytes[0]) % len(characters)
		password.WriteString(string(characters[randomIndex]))
	}

	return password.String(), nil
}
