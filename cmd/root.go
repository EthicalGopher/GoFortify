/*
Copyright © 2026 GoFortify
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)



// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gofortify",
	Short: "GoFortify: A High-Performance Security Reverse Proxy & Traffic Inspector",
	Long: `GoFortify is a sophisticated security reverse proxy designed to protect backend 
services from common web-based attacks. It sits in front of your application, 
intercepting and analyzing incoming traffic in real-time.

Key Features:
- SQL Injection (SQLi) Detection & Mitigation
- Cross-Site Scripting (XSS) Prevention
- Intelligent Rate Limiting & Brute Force Protection
- Interactive Terminal User Interface (TUI) for Real-Time Monitoring
- Detailed Security Event Logging

Usage:
  gofortify [command]`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Persistent flags that will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.GoFortify.yaml)")
}


