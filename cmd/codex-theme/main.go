package main

import (
	"fmt"
	"os"

	"github.com/ychampion/codex-themes/internal/cli"
)

func main() {
	if err := cli.NewRootCommand(cli.Options{}).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
