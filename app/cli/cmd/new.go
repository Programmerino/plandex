package cmd

import (
	"fmt"
	"os"

	"plandex-cli/api"
	"plandex-cli/auth"
	"plandex-cli/lib"
	"plandex-cli/term"
	"plandex-cli/types"

	shared "plandex-shared"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var name string
var contextBaseDir string
var noGlobalContext bool
var globalContextFile string

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:     "new",
	Aliases: []string{"n"},
	Short:   "Start a new plan",
	// Long:  ``,
	Args: cobra.ExactArgs(0),
	Run:  new,
}

func init() {
	RootCmd.AddCommand(newCmd)
	newCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the new plan")
	newCmd.Flags().StringVar(&contextBaseDir, "context-dir", ".", "Base directory to auto-load context from")
	newCmd.Flags().BoolVar(&noGlobalContext, "no-global-context", false, "Don't load the global context for this plan")
	newCmd.Flags().StringVar(&globalContextFile, "global-context-file", "", "Override global context with content from this file")

	AddNewPlanFlags(newCmd)
}

func new(cmd *cobra.Command, args []string) {
	auth.MustResolveAuthWithOrg()
	lib.MustResolveOrCreateProject()

	term.StartSpinner("")

	errCh := make(chan error, 2)

	var planId string
	var config *shared.PlanConfig

	go func() {
		req := shared.CreatePlanRequest{
			Name: name,
			SuppressGlobalContext: noGlobalContext || globalContextFile != "",
		}
		res, apiErr := api.Client.CreatePlan(lib.CurrentProjectId, req)
		if apiErr != nil {
			errCh <- fmt.Errorf("error creating plan: %v", apiErr.Msg)
			return
		}
		planId = res.Id
		errCh <- nil
	}()

	go func() {
		var apiErr *shared.ApiError
		config, apiErr = api.Client.GetDefaultPlanConfig()
		if apiErr != nil {
			errCh <- fmt.Errorf("error getting plan config: %v", apiErr.Msg)
			return
		}
		errCh <- nil
	}()

	for i := 0; i < 2; i++ {
		err := <-errCh
		if err != nil {
			term.OutputErrorAndExit("Error: %v", err)
		}
	}

	err := lib.WriteCurrentPlan(planId)

	if err != nil {
		term.OutputErrorAndExit("Error setting current plan: %v", err)
	}

	err = lib.WriteCurrentBranch("main")
	if err != nil {
		term.OutputErrorAndExit("Error setting current branch: %v", err)
	}

	if name == "" {
		name = "draft"
	}

	term.StopSpinner()

	fmt.Printf("✅ Started new plan %s and set it to current plan\n", color.New(color.Bold, term.ColorHiGreen).Sprint(name))
	fmt.Printf("⚙️  Using default config\n")

	// Handle global context file override
	if globalContextFile != "" {
		content, err := os.ReadFile(globalContextFile)
		if err != nil {
			term.OutputErrorAndExit("Error reading global context file: %v", err)
		}
		// Load the override context as a note
		lib.MustLoadContext(nil, &types.LoadContextParams{
			Note: string(content),
		})
		fmt.Printf("📋 Loaded global context from %s\n", globalContextFile)
	} else if noGlobalContext {
		fmt.Println("🚫 Global context disabled for this plan")
	}

	resolveAutoMode(config)

	resolveModelPack()

	// autoModeLabel := shared.ConfigSettingsByKey["automode"].KeyToLabel(string(config.AutoMode))
	// fmt.Println("⚡️ Auto-mode:", autoModeLabel)

	if config.AutoLoadContext {
		fmt.Println("📥 Automatic context loading is enabled")

		baseDir := contextBaseDir
		if baseDir == "" {
			baseDir = "."
		}

		lib.MustLoadContext([]string{baseDir}, &types.LoadContextParams{
			DefsOnly:          true,
			SkipIgnoreWarning: true,
			AutoLoaded:        true,
		})
	} else {
		fmt.Println()
	}

	var cmds []string
	if term.IsRepl {
		cmds = []string{"config", "plans", "cd", "models"}
	} else {
		cmds = []string{"tell", "chat", "config"}
	}

	if !config.AutoLoadContext {
		cmds = append([]string{"load"}, cmds...)
	}

	fmt.Println()
	term.PrintCmds("", cmds...)
}
