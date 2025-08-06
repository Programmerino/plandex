# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Plandex is an AI coding agent terminal-based tool that can plan and execute large coding tasks. It consists of:
- **CLI** (`app/cli/`): Go-based command-line interface
- **Server** (`app/server/`): Go backend with PostgreSQL database
- **Shared** (`app/shared/`): Shared Go code between CLI and server

## Development Setup

### Prerequisites
- Go 1.23.3
- PostgreSQL 14
- Python 3 (for LiteLLM proxy)
- reflex 0.3.1 (`go install github.com/cespare/reflex@v0.3.1`)

### Nix Development Environment

**IMPORTANT: This project uses Nix for development. Always use nix-shell for building and running.**

```bash
# Enter the Nix development shell (from project root)
nix-shell

# Or run commands directly within nix-shell
nix-shell --run "cd app/server && go build"
nix-shell --run "cd app/cli && go build"
```

The `shell.nix` file provides:
- Go 1.23.11
- PostgreSQL 14
- Python 3.13
- All required development tools

### Build Commands

**Note: Run all build commands within nix-shell or use `nix-shell --run "command"`**

**CLI Build:**
```bash
# Within nix-shell:
cd app/cli
go build -o plandex-dev

# Or from outside nix-shell:
nix-shell --run "cd app/cli && go build -o plandex-dev"
```

**Server Build:**
```bash
# Within nix-shell:
cd app/server
go build

# Or from outside nix-shell:
nix-shell --run "cd app/server && go build"
```

**Full Development Mode (with hot reload):**
```bash
# From root directory - builds both CLI and server with file watchers
nix-shell --run "./app/scripts/dev.sh"
```

**Docker Local Mode:**
```bash
# Start local server with PostgreSQL
cd app
./start_local.sh  # or docker compose up
```

### Testing
```bash
# Run tests within nix-shell:
nix-shell --run "./test/smoke_test.sh"

# Other test scripts available:
nix-shell --run "./test/test_custom_models.sh"
nix-shell --run "./test/plan_deletion_test.sh"
```

### Environment Variables
```bash
# Development
export DATABASE_URL="postgres://user:password@localhost:5432/plandex?sslmode=disable"
export GOENV=development
export PLANDEX_ENV=development  # For CLI to connect to dev server
export LOCAL_MODE=1  # For local server mode

# API Keys (if using BYO keys)
export OPENROUTER_API_KEY=...
```

## Architecture Overview

### Core Components

1. **Plan System** (`app/server/model/plan/`): Central execution engine
   - `tell_exec.go`: Main plan execution logic
   - `build_exec.go`: Code generation and building
   - `tell_stream_processor.go`: Streaming response handler
   - Plans support branching, versioning, and rollback

2. **Context Management** (`app/server/db/context_helpers_*.go`): 
   - Handles loading files/directories into plan context
   - Smart context window management (2M tokens effective)
   - Tree-sitter based project mapping for large codebases

3. **Model Integration** (`app/shared/ai_models_*.go`):
   - Supports multiple providers: OpenAI, Anthropic, Google, OpenRouter, Ollama
   - Model packs for different cost/capability tradeoffs
   - Custom model configuration support

4. **File Operations & Diff System** (`app/server/syntax/`):
   - Tree-sitter based structured edits
   - Reliable file modifications with validation
   - Git-style diff generation and application

5. **CLI Command Structure** (`app/cli/cmd/`):
   - Cobra-based command framework
   - REPL mode (`repl.go`) with fuzzy autocomplete
   - Each command in separate file (e.g., `tell.go`, `build.go`, `apply.go`)

### Key Data Flow

1. User issues command → CLI (`app/cli/cmd/`)
2. CLI calls API → Server handlers (`app/server/handlers/`)
3. Server processes with models → Plan execution (`app/server/model/plan/`)
4. Results streamed back → CLI displays updates (`app/cli/stream_tui/`)
5. Changes sandboxed until `apply` command

### Database Schema
- Uses PostgreSQL with migrations (`app/server/migrations/`)
- Key tables: plans, contexts, builds, organizations, users
- Transaction-based operations for consistency

## Important Conventions

- Go modules: CLI and Server are separate modules, both import `plandex-shared`
- File paths in code use absolute paths
- Context updates are atomic - all succeed or none apply
- Streaming responses use custom protocol over HTTP
- Git integration optional but recommended

## Common Development Tasks

**Add a new CLI command:**
1. Create new file in `app/cli/cmd/`
2. Register with Cobra in command init
3. Add API client method if needed in `app/cli/api/`

**Add a new model provider:**
1. Update provider configs in `app/shared/ai_models_providers.go`
2. Add credentials handling in `app/shared/ai_models_credentials.go`
3. Update model client in `app/server/model/client.go`

**Modify plan execution logic:**
- Core logic in `app/server/model/plan/tell_exec.go`
- Stream processing in `tell_stream_processor.go`
- Prompts in `app/server/model/prompts/`

## Debugging Tips

- Logs written to `~/.plandex/plandex.log` (rotating)
- Use `--debug` flag for verbose output
- Server logs to stdout in development mode
- LiteLLM proxy runs on port 4000 for model request debugging