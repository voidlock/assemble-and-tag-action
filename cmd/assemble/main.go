package main

import (
	"context"

	"github.com/sethvargo/go-githubactions"
	"github.com/voidlock/assemble-and-tag/pkg/assemble"
)

func main() {
	ctx := context.Background()
	action := githubactions.New()

	cmd, err := assemble.NewFromInputs(action)
	if err != nil {
		githubactions.Fatalf("%v", err)
	}

	if err := cmd.Run(ctx); err != nil {
		githubactions.Fatalf("%v", err)
	}
}
