package main

import (
	"fmt"
	"os"

	"github.com/EthicalGopher/SentinelShield/server"
	"github.com/EthicalGopher/SentinelShield/tui"
	"github.com/EthicalGopher/SentinelShield/tui/shared"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	go func() {
		if err := server.Server(); err != nil {
			// Send the error to the log channel. The TUI will pick it up when it's ready.
			shared.LogChan <- fmt.Sprintf("FATAL: Server exited with error: %v", err)
		}
	}()

	p := tea.NewProgram(tui.NewRoot())
	if _, err := p.Run(); err != nil {
		fmt.Printf("TUI crashed: %v\n", err)
		os.Exit(1)
	}
}
