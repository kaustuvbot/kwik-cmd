package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Start interactive mode",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInteractive()
	},
}

func runInteractive() error {
	reader := bufio.NewReader(os.Stdin)
	
	fmt.Println("=== kwik-cmd Interactive Mode ===")
	fmt.Println("Type 'help' for available commands, 'exit' to quit")
	fmt.Println()

	for {
		fmt.Print("kwik-cmd> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		if input == "exit" || input == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		if input == "help" {
			printHelp()
			continue
		}

		// Parse and execute simple commands
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "track":
			if len(parts) > 1 {
				trackInput := strings.Join(parts[1:], " ")
				trackInput = strings.Trim(trackInput, "\"")
				fmt.Printf("Tracking: %s\n", trackInput)
				// Would call tracker here
			}
		case "suggest":
			fmt.Println("Use 'kwik-cmd suggest <partial>'")
		case "stats":
			fmt.Println("Use 'kwik-cmd stats'")
		default:
			fmt.Printf("Unknown command: %s\n", parts[0])
		}
	}

	return nil
}

func printHelp() {
	fmt.Println(`Available commands:
  track <command>  - Track a command
  suggest <partial> - Get suggestions
  stats             - Show statistics  
  analyze           - Analyze patterns
  export            - Export history
  help              - Show this help
  exit              - Quit`)
}

func init() {
	rootCmd.AddCommand(interactiveCmd)
}
