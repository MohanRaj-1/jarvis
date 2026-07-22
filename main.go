package main

import (
	"context"
	"log"

	aitools "jarvis/tools/ai"
	analyzertools "jarvis/tools/analyzer"
	gittools "jarvis/tools/git"
	"jarvis/tools/greeting"
	mathtools "jarvis/tools/math"
	systemtools "jarvis/tools/system"
	workspacetools "jarvis/tools/workspace"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "Jarvis",
			Version: "0.1.0",
		},
		nil,
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "hello",
			Description: "Greets a person",
		},
		greeting.Hello,
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "bye",
			Description: "Farewell a person",
		},
		greeting.Bye,
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "add",
			Description: "Add two numbers",
		},
		mathtools.Add,
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "weekday",
			Description: "Gives the current day of the week",
		},
		systemtools.Weekday,
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "read_file",
			Description: "Reads the contents and size of a file up to 1 MB",
		},
		workspacetools.ReadFile,
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "list_directory",
			Description: "Lists the immediate children of a directory",
		},
		workspacetools.ListDirectory,
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "find_files",
			Description: "Recursively finds files matching a filepath glob pattern",
		},
		workspacetools.FindFiles,
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "file_info",
			Description: "Returns metadata for a file or directory",
		},
		workspacetools.FileInfo,
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "analyze_imports",
			Description: "Extracts import paths from a Go source file",
		},
		analyzertools.AnalyzeImports,
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "analyze_go_file",
			Description: "Extracts imports and functions from a Go source file",
		},
		analyzertools.AnalyzeGoFile,
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "git_current_branch",
			Description: "Returns the current branch of a Git repository",
		},
		gittools.CurrentBranch,
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "git_status",
			Description: "Returns the current branch and file status of a Git repository",
		},
		gittools.RepositoryStatus,
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "git_branches",
			Description: "Returns the current branch and all local branches of a Git repository",
		},
		gittools.RepositoryBranches,
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "explain_go_file",
			Description: "Explains a Go source file using Gemini",
		},
		aitools.ExplainGoFile,
	)

	log.Println("Jarvis starting...")

	if err := server.Run(
		context.Background(),
		&mcp.StdioTransport{},
	); err != nil {
		log.Fatal(err)
	}
}
