// Package main is the entry point for the GoFortify application.
// It initializes and executes the root command for the CLI.
package main

import (
	"github.com/EthicalGopher/GoFortify/cmd"
)

// main is the primary entry point which delegates execution to the cmd package.
func main() {
	cmd.Execute()
}
