package editor

import (
	"fmt"
	"os"
	"os/exec"
)

// tries to read using the default editor
func Read() (string, error) {
	tempfile, err := os.CreateTemp("", "input.txt")
	if err != nil {
		fmt.Println("error creating temp file:", err)
		return "", err
	}

	defer func(tempfie *os.File) {
		err := tempfie.Close()
		if err != nil {
			fmt.Println("error closing temp file:", err)
		}
		err = os.Remove(tempfie.Name())
		if err != nil {
			fmt.Println("error removing temp file:", err)
		}
	}(tempfile)

	cmd := exec.Command("vim", tempfile.Name())

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("error running editor:", err)
		return "", err
	}

	content, err := os.ReadFile(tempfile.Name())
	if err != nil {
		return "", err
	}
	return string(content), nil
}
