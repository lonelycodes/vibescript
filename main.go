package main

import (
	"os"

	"github.com/lonelycodes/vibescript/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
