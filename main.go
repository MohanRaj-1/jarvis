package main

import (
	"context"
	"log"
	"os"
	"strings"

	internalai "jarvis/internal/ai"
	aigit "jarvis/internal/ai/git"
	internalgit "jarvis/internal/git"
	internalworkspace "jarvis/internal/workspace"
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
	workspaceRoot := strings.TrimSpace(os.Getenv("JARVIS_WORKSPACE_ROOT"))
	if workspaceRoot == "" {
		var err error
		workspaceRoot, err = os.Getwd()
		if err != nil {
			log.Fatalf("determine workspace root: %v", err)
		}
	}
	ws, err := internalworkspace.New(workspaceRoot)
	if err != nil {
		log.Fatalf("initialize workspace: %v", err)
	}
	readFileTool := workspacetools.NewReadFileTool(ws)
	listDirectoryTool := workspacetools.NewListDirectoryTool(ws)
	findFilesTool := workspacetools.NewFindFilesTool(ws)
	fileInfoTool := workspacetools.NewFileInfoTool(ws)
	analyzeImportsTool := analyzertools.NewAnalyzeImportsTool(ws)
	analyzeGoFileTool := analyzertools.NewAnalyzeGoFileTool(ws)
	repository := internalgit.NewDefaultRepository(ws)
	currentBranchTool := gittools.NewCurrentBranchTool(repository)
	statusTool := gittools.NewStatusTool(repository)
	diffTool := gittools.NewDiffTool(repository)
	branchesTool := gittools.NewBranchesTool(repository)
	logTool := gittools.NewLogTool(repository)
	showCommitTool := gittools.NewShowCommitTool(repository)

	aiClient, err := internalai.NewClient(context.Background(), internalai.Config{
		Provider: internalai.ProviderGemini,
		APIKey:   os.Getenv("GEMINI_API_KEY"),
		Model:    geminiModel(),
	})
	if err != nil {
		log.Fatalf("initialize AI client: %v", err)
	}

	commitMessageService := aigit.CommitMessageService{
		Git: repository,
		AI:  aiClient,
	}
	explainCommitService := aigit.ExplainCommitService{
		Git: repository,
		AI:  aiClient,
	}
	reviewDiffService := aigit.ReviewDiffService{
		Git: repository,
		AI:  aiClient,
	}
	releaseNotesService := aigit.ReleaseNotesService{
		Git: repository,
		AI:  aiClient,
	}
	repositorySummaryService := aigit.RepositorySummaryService{
		Git: repository,
		AI:  aiClient,
	}

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "Jarvis",
			Version: "0.7.0",
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
		readFileTool.Handle,
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "list_directory",
			Description: "Lists the immediate children of a directory",
		},
		listDirectoryTool.Handle,
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "find_files",
			Description: "Recursively finds files matching a filepath glob pattern",
		},
		findFilesTool.Handle,
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "file_info",
			Description: "Returns metadata for a file or directory",
		},
		fileInfoTool.Handle,
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "analyze_imports",
			Description: "Extracts import paths from a Go source file",
		},
		analyzeImportsTool.Handle,
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "analyze_go_file",
			Description: "Extracts imports and functions from a Go source file",
		},
		analyzeGoFileTool.Handle,
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "git_current_branch",
			Description: "Returns the current branch of a Git repository",
		},
		currentBranchTool.Handle,
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "git_status",
			Description: "Returns the current branch and file status of a Git repository",
		},
		statusTool.Handle,
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "git_diff",
			Description: "Returns the current working tree diff of a Git repository",
		},
		diffTool.Handle,
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "git_branches",
			Description: "Returns the current branch and all local branches of a Git repository",
		},
		branchesTool.Handle,
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "git_log",
			Description: "Returns the most recent commits in a Git repository",
		},
		logTool.Handle,
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "git_show_commit",
			Description: "Returns detailed information about a Git commit",
		},
		showCommitTool.Handle,
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "git_generate_commit_message",
			Description: "Generates a Conventional Commit message from a Git working tree diff",
		},
		gittools.NewGenerateCommitMessage(commitMessageService),
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "git_explain_commit",
			Description: "Explains what a Git commit changed and why it matters",
		},
		gittools.NewExplainCommit(explainCommitService),
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "git_review_diff",
			Description: "Reviews the current working tree diff and returns a structured report",
		},
		gittools.NewReviewDiff(reviewDiffService),
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "git_generate_release_notes",
			Description: "Generates structured release notes from a Git commit range",
		},
		gittools.NewGenerateReleaseNotes(releaseNotesService),
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "git_repository_summary",
			Description: "Summarizes the current branch, recent work, and working changes of a Git repository",
		},
		gittools.NewRepositorySummary(repositorySummaryService),
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "explain_go_file",
			Description: "Explains a Go source file using the configured AI provider",
		},
		aitools.NewExplainGoFile(aiClient),
	)

	log.Println("Jarvis starting...")

	if err := server.Run(
		context.Background(),
		&mcp.StdioTransport{},
	); err != nil {
		log.Fatal(err)
	}
}

func geminiModel() string {
	if model := strings.TrimSpace(os.Getenv("GEMINI_MODEL")); model != "" {
		return model
	}

	return "gemini-3.6-flash"
}
