/*
Copyright © 2026 GoFortify
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/EthicalGopher/GoFortify/server"
	"github.com/EthicalGopher/GoFortify/shared"
	"github.com/EthicalGopher/GoFortify/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	backendURL    string
	proxyPort     int
	filename_rate string
	filename_xss  string
	filename_sql  string
	ratelimit     int
)

// initialiseCMD represents the command to start the proxy and TUI
var initialiseCMD = &cobra.Command{
	Use:   "init",
	Short: "Initialize and start the security proxy and monitoring TUI",
	Long: `The 'init' command launches the GoFortify security engine. 

It starts a high-performance reverse proxy that:
1. Listens for incoming HTTP traffic on a specified local port.
2. Performs deep packet inspection for SQLi and XSS patterns.
3. Enforces rate-limiting policies based on client IP.
4. Forwards validated, clean traffic to your upstream backend server.
5. Displays a real-time dashboard (TUI) for immediate threat visibility.

All security events are logged to JSON files for further analysis.`,
	Example: `  # Start proxy on 5174, forwarding to localhost:8080
  gofortify init --port 5174 --backend-url http://localhost:8080

  # Start with customized rate limiting (e.g., 50 req/min)
  gofortify init -p 3000 -b http://127.0.0.1:8081 --ratelimit 50`,
	Run: func(cmd *cobra.Command, args []string) {
		go func() {
			if err := server.Server(backendURL, proxyPort, filename_rate, filename_xss, filename_sql, ratelimit); err != nil {
				shared.LogChan <- fmt.Sprintf("FATAL: Server exited with error: %v", err)
			}
		}()
		p := tea.NewProgram(tui.NewRoot())
		if _, err := p.Run(); err != nil {
			fmt.Printf("TUI crashed: %v\n", err)
			os.Exit(1)
		}

	},
}

func init() {
	rootCmd.AddCommand(initialiseCMD)
	initialiseCMD.Flags().IntVarP(&proxyPort, "port", "p", 5174, "Port for the GoFortify proxy to listen on")
	initialiseCMD.Flags().StringVarP(&backendURL, "backend-url", "b", "http://localhost:8080", "URL of the backend server to protect")
	initialiseCMD.Flags().StringVarP(&filename_rate, "ratelimit-file", "f", "vulnerabilities/rate_limit.json", "Filename for rate limit logs")
	initialiseCMD.Flags().StringVarP(&filename_sql, "sql", "s", "vulnerabilities/sqlInjection.json", "Filename for SQL injection logs")
	initialiseCMD.Flags().StringVarP(&filename_xss, "xss", "x", "vulnerabilities/xss.json", "Filename for XSS logs")
	initialiseCMD.Flags().IntVarP(&ratelimit, "ratelimit", "r", 100, "Filename for XSS logs")

}
