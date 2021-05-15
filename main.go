package main

import (
	"chia-bouncer/bouncer"
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	_, err := bouncer.Run()
	if err != nil {
		return err
	}
	return nil
}
