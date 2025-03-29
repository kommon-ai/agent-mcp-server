package goosemcp

import (
	"context"
	"fmt"
	"os"

	"github.com/kommon-ai/agent-go/pkg/agent"
	"github.com/kommon-ai/goose-connect/pkg/goose"
	"github.com/strowk/foxy-contexts/pkg/fxctx"
	"github.com/strowk/foxy-contexts/pkg/mcp"
	"github.com/strowk/foxy-contexts/pkg/toolinput"
	"k8s.io/utils/ptr"
)

func NewGooseTool() fxctx.Tool {
	schema := toolinput.NewToolInputSchema(
		toolinput.WithRequiredString("name", "goose agent name"),
		toolinput.WithRequiredString("prompt", "prompt to run"),
		toolinput.WithString("instruction", "instruction to run"),
		toolinput.WithRequiredString("repo", "repository to run"),
	)
	description := "Run a goose agent"
	return fxctx.NewTool(
		&mcp.Tool{
			Name:        "goose-run",
			Description: &description,
			InputSchema: schema.GetMcpToolInputSchema(),
		},
		func(args map[string]any) *mcp.CallToolResult {
			input, err := schema.Validate(args)
			if err != nil {
				return errResponse(fmt.Errorf("invalid input: %w", err))
			}
			name, err := input.String("name")
			if err != nil {
				return errResponse(fmt.Errorf("invalid input: %w", err))
			}
			prompt, err := input.String("prompt")
			if err != nil {
				return errResponse(fmt.Errorf("invalid input: %w", err))
			}
			instruction, _ := input.String("instruction")
			openRouterAPIKey := os.Getenv("OPENROUTER_API_KEY")
			if openRouterAPIKey == "" {
				return errResponse(fmt.Errorf("OPENROUTER_API_KEY is not set"))
			}
			installationToken := os.Getenv("GITHUB_TOKEN")
			if installationToken == "" {
				return errResponse(fmt.Errorf("GITHUB_TOKEN is not set"))
			}
			repo, _ := input.String("repo")
			a, err := goose.NewGooseAgent(goose.GooseOptions{
				SessionID:   name,
				Instruction: instruction,
				GitHub: &goose.GooseGitHub{
					InstallationToken: installationToken,
					Host:              "https://github.com",
					Repo:              repo,
				},
				Provider: &agent.NoopProvider{},
			})
			if err != nil {
				return errResponse(fmt.Errorf("failed to run goose agent: %w", err))
			}
			data, err := a.Execute(context.Background(), prompt)
			if err != nil {
				return errResponse(fmt.Errorf("failed to run goose agent: %w", err))
			}
			content := mcp.TextContent{
				Type: "text",
				Text: string(data),
			}
			return &mcp.CallToolResult{
				Content: []any{content},
				IsError: ptr.To(false),
			}
		},
	)
}

func errResponse(err error) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		IsError: ptr.To(true),
		Meta:    map[string]any{},
		Content: []any{
			mcp.TextContent{
				Type: "text",
				Text: err.Error(),
			},
		},
	}
}
