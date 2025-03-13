package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error: "+err.Error())
		os.Exit(1)
	}
}
