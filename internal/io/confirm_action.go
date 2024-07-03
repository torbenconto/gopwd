package io

import "fmt"

func ConfirmAction() (bool, error) {
	fmt.Print("Are you sure? [y/N] ")
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return false, err
	}

	if input == "y" || input == "Y" {
		return true, nil
	}

	return false, nil
}
