package termio

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func ReadMultiline() (string, error) {
	fmt.Println("Enter text (press Ctrl+D to finish):")

	var input strings.Builder
	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanned := scanner.Scan()
		if err := scanner.Err(); err != nil {
			if err == io.EOF {
				break // Exit the loop on EOF
			}
			return "", fmt.Errorf("error reading input: %v", err)
		}
		if !scanned {
			break // Also exit the loop if Scan returns false
		}
		line := scanner.Text()
		if input.Len() > 0 {
			input.WriteString("\n")
		}
		input.WriteString(line)
	}

	return input.String(), nil
}
