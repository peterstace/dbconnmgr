package main

import (
	"fmt"
	"os"
)

func main() {
	tree, err := ParseTree(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse tree: %v\n", err)
		os.Exit(1)
	}
	connDetails, err := tree.Flatten()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not flatten tree into connection details: %v\n", err)
		os.Exit(1)
	}

	selections := connDetails.Selections()
	choice, ok, err := fzf(selections, ">")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not fzf: %v\n", err)
		os.Exit(1)
	}
	if !ok {
		fmt.Fprintf(os.Stderr, "cancelled\n")
		os.Exit(130)
	}
	for i, sel := range selections {
		if choice == sel {
			fmt.Println(connDetails[i].ConnectionString())
		}
	}
}
