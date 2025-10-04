package main

import (
	"context"
	"fmt"
	"os"

	"github.com/LFWQSP2641/scunet-auto-login/cmd"
)

func main() {
	if err := cmd.Execute(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
