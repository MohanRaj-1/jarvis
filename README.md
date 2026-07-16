# Jarvis

Jarvis is an AI-powered developer assistant built with Go and the Model Context Protocol (MCP). It helps developers analyze Go source code, interact with workspace files, and leverage AI to better understand code.

> 🚧 **Current Version:** v0.5.0 (Beta)

---

## Features

### 🤖 AI Developer Assistant

- AI-powered Go file explanations
- Google Gemini integration
- Provider abstraction for future AI providers

### 🔍 Go Source Analyzer

- Parse Go source code using the Go AST
- Extract imports
- Extract functions
- Extract structs
- Extract interfaces
- Extract methods
- Package analysis

### 📁 Workspace Tools

- Read File
- List Directory
- File Info
- Find Files

### 🛠 Utility Tools

- Greeting Tool
- Math Tool
- Weekday Tool

---

## Tech Stack

- Go
- Model Context Protocol (MCP)
- Google Gemini API
- Git & GitHub

---

## Project Structure

```text
jarvis/
├── internal/
│   ├── ai/
│   ├── analyzer/
│   └── workspace/
├── tools/
├── main.go
├── go.mod
└── README.md
```

---

## Roadmap

- [x] Project setup
- [x] Basic MCP tools
- [x] Workspace tools
- [x] Go source analyzer
- [x] AI-powered Go file explanations
- [ ] Project-wide analysis
- [ ] AI chat improvements
- [ ] Testing and CI/CD
- [ ] Docker support
- [ ] Production release (v1.0.0)

---

## Releases

| Version | Description                   |
| ------- | ----------------------------- |
| v0.1.0  | Project setup                 |
| v0.2.0  | Basic MCP tools               |
| v0.3.0  | Workspace tools               |
| v0.4.0  | Go source analyzer            |
| v0.5.0  | AI Developer Assistant (Beta) |

---

## License

This project is licensed under the MIT License.
