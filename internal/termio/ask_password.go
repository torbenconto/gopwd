package termio

import (
	"fmt"
	"golang.org/x/term"
	"syscall"
)

func PromptPassword() (string, error) {
	fmt.Print("Enter password: ")
	password, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()

	if err != nil {
		return "", err
	}

	fmt.Print("Confirm password: ")
	confirm, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()

	if err != nil {
		return "", err
	}

	if string(password) != string(confirm) {
		return "", fmt.Errorf("passwords do not match")
	}

	return string(password), nil
}
