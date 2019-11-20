package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	tree, err := ParseTree(os.Stdin)
	if err != nil {
		log.Fatalf("Could not parse tree: %v", err)
	}
	connDetails, err := tree.Flatten()
	if err != nil {
		log.Fatalf("Could not flatten tree into connection details: %v", err)
	}

	selections := connDetails.Selections()
	choice, ok, err := fzf(selections, ">")
	if err != nil {
		log.Fatalf("could not fzf: %v", err)
	}
	if !ok {
		return
	}
	for i, sel := range selections {
		if choice == sel {
			fmt.Println(connDetails[i].ConnectionString())
		}
	}
}
