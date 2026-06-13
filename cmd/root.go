package cmd

import (
	"fmt"
	"os"

	"github.com/lansdownian/goalans/internal/app"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var mockFlag bool

var rootCmd = &cobra.Command{
	Use:   "goalans",
	Short: "Live football scores in your terminal",
	Long:  `Goalans is a terminal UI for following football matches — live scores, finished results, and match events.`,
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(app.New(mockFlag), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVar(&mockFlag, "mock", false, "Use sample data instead of fetching from the network")
}
