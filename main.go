package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/beyang/srclib-ctags/graph"
)

func main() {
	defs, err := graph.DefsForFiles([]string{"app/app.go"})
	if err != nil {
		fmt.Printf("failed due to error: %s\n", err)
		os.Exit(1)
	}
	json.NewEncoder(os.Stdout).Encode(defs)
}
