package promptutils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Prompt_multiline_util(header string) (string, error) {
	fmt.Printf("\n%s\n", header)
	fmt.Println("  (Enter each line. Leave a blank line to finish.)")

	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("  > ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			break
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return strings.Join(lines, "\n"), nil
}
