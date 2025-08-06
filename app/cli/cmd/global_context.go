package cmd

import (
	"fmt"
	"io"
	"os"
	"plandex-cli/api"
	"plandex-cli/auth"
	"plandex-cli/term"

	"github.com/spf13/cobra"
)

var globalContextCmd = &cobra.Command{
	Use:     "global-context",
	Aliases: []string{"gc"},
	Short:   "Manage global user context",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Usage: plandex global-context [set|show|clear]")
	},
}

var globalContextSetCmd = &cobra.Command{
	Use:   "set [file_path]",
	Short: "Set global context from a file or stdin",
	Args:  cobra.MaximumNArgs(1),
	Run:   globalContextSet,
}

var globalContextShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the current global context",
	Run:   globalContextShow,
}

var globalContextClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear the global context",
	Run:   globalContextClear,
}

func init() {
	globalContextCmd.AddCommand(globalContextSetCmd)
	globalContextCmd.AddCommand(globalContextShowCmd)
	globalContextCmd.AddCommand(globalContextClearCmd)
	RootCmd.AddCommand(globalContextCmd)
}

func globalContextSet(cmd *cobra.Command, args []string) {
	auth.MustResolveAuthWithOrg()

	var content []byte
	var err error
	if len(args) > 0 {
		// Read from file
		content, err = os.ReadFile(args[0])
		if err != nil {
			term.OutputErrorAndExit("Error reading file: %v", err)
		}
	} else {
		// Read from stdin
		content, err = io.ReadAll(os.Stdin)
		if err != nil {
			term.OutputErrorAndExit("Error reading from stdin: %v", err)
		}
	}

	// Make API request
	apiErr := api.Client.UpdateGlobalContext(string(content))
	if apiErr != nil {
		term.OutputErrorAndExit("Error updating global context: %v", apiErr.Msg)
	}

	fmt.Println("✓ Global context updated successfully")
}

func globalContextShow(cmd *cobra.Command, args []string) {
	auth.MustResolveAuthWithOrg()

	// Make API request
	globalContext, apiErr := api.Client.GetGlobalContext()
	if apiErr != nil {
		term.OutputErrorAndExit("Error getting global context: %v", apiErr.Msg)
	}

	if globalContext == "" {
		fmt.Println("No global context set")
		return
	}

	fmt.Println("Global Context:")
	fmt.Println("================")
	fmt.Println(globalContext)
}

func globalContextClear(cmd *cobra.Command, args []string) {
	auth.MustResolveAuthWithOrg()

	// Make API request
	apiErr := api.Client.DeleteGlobalContext()
	if apiErr != nil {
		term.OutputErrorAndExit("Error clearing global context: %v", apiErr.Msg)
	}

	fmt.Println("✓ Global context cleared")
}