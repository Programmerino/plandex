{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  name = "plandex-dev";
  
  buildInputs = with pkgs; [
    # Go development
    go_1_23
    gopls
    delve
    gotools
    
    # Build tools
    gnumake
    gcc
    
    # Git for version control
    git
    
    # Tree-sitter for server
    tree-sitter
    
    # PostgreSQL for database
    postgresql_14
    
    # Python for LiteLLM proxy
    (python3.withPackages (ps: with ps; [
      litellm
      uvicorn
      fastapi
      aiohttp
      backoff
      orjson
      pyjwt
      pyyaml
      rich
    ]))
    
    # For testing apply scripts
    bash
    zsh
    
    # Utilities
    jq
    curl
    wget
  ];

  shellHook = ''
    echo "Plandex development environment"
    echo "Go version: $(go version)"
    echo "PostgreSQL version: $(postgres --version)"
    echo "Python version: $(python3 --version)"
    echo ""
    echo "Build CLI: cd app/cli && go build -o plandex"
    echo "Build Server: cd app/server && go build"
    echo "Run tests: ./test/smoke_test.sh"
    echo ""
    
    # Set up Go environment
    export GOPATH=$HOME/go
    export PATH=$GOPATH/bin:$PATH
    
    # Install reflex for hot reload if not already installed
    if ! command -v reflex &> /dev/null; then
      echo "Installing reflex for hot reload..."
      go install github.com/cespare/reflex@v0.3.1
    fi
  '';
}