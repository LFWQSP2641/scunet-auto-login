package main

import (
	"context"
	"fmt"
	"os"
	"scunet-auto-login/cmd"
)

func main() {
	if err := cmd.Execute(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
