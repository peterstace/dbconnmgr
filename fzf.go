package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
)

func init() {
	// Check to make sure fzf is installed.
	if err := exec.Command("fzf", "--version").Run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func fzf(input []string, prompt string) (choice string, ok bool, err error) {
	cmd := exec.Command(
		"fzf",
		"--height", strconv.Itoa(1+len(input)),
		"--inline-info",
		"--prompt", prompt,
		"--reverse",
	)
	cmd.Stdout = new(bytes.Buffer)
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", false, err
	}
	go func() {
		defer stdin.Close()
		for _, i := range input {
			fmt.Fprintf(stdin, "%s\n", i)
		}
	}()

	if err := cmd.Run(); err != nil {
		switch err.(type) {
		case *exec.ExitError:
			// Ctrl-C or bad input provided.
			return "", false, nil
		default:
			return "", false, err
		}
	}

	trimmed := cmd.Stdout.(*bytes.Buffer).Bytes()
	trimmed = trimmed[:len(trimmed)-1]
	return string(trimmed), true, nil
}
