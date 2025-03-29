/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/kommon-ai/agent-mcp-server/pkg/goosemcp"
	"github.com/spf13/cobra"
	"github.com/strowk/foxy-contexts/pkg/app"
	"github.com/strowk/foxy-contexts/pkg/mcp"
	"github.com/strowk/foxy-contexts/pkg/stdio"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"k8s.io/utils/ptr"
)

var (
	agentType     string
	agentProvider string
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the MCP server with specified agent type",
	Long: `Start the MCP (Master Control Program) server with the specified agent type.
This command will initialize and run the server instance configured for the
chosen agent type and provider.

Example:
  agent-mcp-server start --agent-type=chat --agent-provider=goose`,
	Run: func(cmd *cobra.Command, args []string) {
		// Here you would implement the actual server startup logic based on agent type and provider
		foxyApp := getApp()
		err := foxyApp.Run()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	},
}

func getCapabilities() *mcp.ServerCapabilities {
	return &mcp.ServerCapabilities{
		Tools: &mcp.ServerCapabilitiesTools{
			ListChanged: ptr.To(false),
		},
	}
}

func getApp() *app.Builder {
	return app.
		NewBuilder().
		WithName("agent-mcp-server").
		WithVersion("0.0.1").
		WithTransport(stdio.NewTransport()).
		WithServerCapabilities(getCapabilities()).
		WithTool(goosemcp.NewGooseTool).
		WithFxOptions(
			fx.Provide(func() *zap.Logger {
				cfg := zap.NewDevelopmentConfig()
				cfg.Level.SetLevel(zap.ErrorLevel)
				logger, _ := cfg.Build()
				return logger
			}),
			fx.Option(fx.WithLogger(
				func(logger *zap.Logger) fxevent.Logger {
					return &fxevent.ZapLogger{Logger: logger}
				},
			)),
		)
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Add the agent-type flag to the start command
	startCmd.Flags().StringVar(&agentType, "agent-type", "", "Specify the agent type (e.g., chat, copilot)")

	// Add the agent-provider flag to the start command
	// Now only goose is supported
	startCmd.Flags().StringVar(&agentProvider, "agent-provider", "", "Specify the agent provider (e.g., goose, etc)")
}
