/*
Copyright © 2026 GoFortify
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// citeCmd represents the cite command
var citeCmd = &cobra.Command{
	Use:   "cite",
	Short: "How to cite GoFortify in your research or project",
	Long: `Display citation information for GoFortify in various formats.
Use this information if you reference GoFortify in academic papers, 
security reports, or other formal documentation.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("How to cite GoFortify:")
		fmt.Println("----------------------")
		fmt.Println("\nPlain Text:")
		fmt.Println("EthicalGopher. (2026). GoFortify: A High-Performance Security Reverse Proxy & Traffic Inspector. https://github.com/EthicalGopher/GoFortify")
		
		fmt.Println("\nBibTeX:")
		bibtex := `@software{GoFortify_2026,
  author = {EthicalGopher},
  title = {{GoFortify: A High-Performance Security Reverse Proxy \& Traffic Inspector}},
  url = {https://github.com/EthicalGopher/GoFortify},
  year = {2026},
  version = {1.0.0}
}`
		fmt.Println(bibtex)
	},
}

func init() {
	rootCmd.AddCommand(citeCmd)
}
