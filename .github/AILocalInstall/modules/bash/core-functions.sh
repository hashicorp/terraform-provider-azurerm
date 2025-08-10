#!/bin/bash

#==============================================================================
# Core Functions Module
#==============================================================================
# 
# This module provides core utility functions used throughout the AI setup
# system, including logging, status messages, and basic system checks.
#

# ANSI color codes for terminal output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m' # No Color

# Logging functions
log_info() {
    local message="$1"
    echo -e "${BLUE}[INFO]${NC} $message"
}

log_success() {
    local message="$1"
    echo -e "${GREEN}[SUCCESS]${NC} $message"
}

log_warning() {
    local message="$1"
    echo -e "${YELLOW}[WARNING]${NC} $message"
}

log_error() {
    local message="$1"
    echo -e "${RED}[ERROR]${NC} $message" >&2
}

# Status message function with different types
write_status_message() {
    local message="$1"
    local type="${2:-Info}"
    
    case "$type" in
        "Info")
            log_info "$message"
            ;;
        "Success")
            log_success "$message"
            ;;
        "Warning")
            log_warning "$message"
            ;;
        "Error")
            log_error "$message"
            ;;
        *)
            echo "$message"
            ;;
    esac
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Check for bash version (3.2+)
    if [[ ${BASH_VERSION%%.*} -lt 3 ]] || [[ ${BASH_VERSION%%.*} -eq 3 && ${BASH_VERSION#*.} -lt 2 ]]; then
        log_error "Bash 3.2 or later is required (current: $BASH_VERSION)"
        return 1
    fi
    
    # Check for required commands
    local required_commands=("git" "jq")
    local missing_commands=()
    
    for cmd in "${required_commands[@]}"; do
        if ! command_exists "$cmd"; then
            missing_commands+=("$cmd")
        fi
    done
    
    if [[ ${#missing_commands[@]} -gt 0 ]]; then
        log_error "Missing required commands: ${missing_commands[*]}"
        log_error "Please install the missing commands and try again"
        return 1
    fi
    
    log_success "Prerequisites check passed"
    return 0
}

# Discover repository root
discover_repository() {
    local current_dir="$PWD"
    
    # Start from current directory and walk up
    while [[ "$current_dir" != "/" ]]; do
        if [[ -f "$current_dir/go.mod" ]] && [[ -d "$current_dir/.git" ]]; then
            # Check if this is the terraform-provider-azurerm repository
            if grep -q "github.com/hashicorp/terraform-provider-azurerm" "$current_dir/go.mod" 2>/dev/null; then
                echo "$current_dir"
                return 0
            fi
        fi
        current_dir="$(dirname "$current_dir")"
    done
    
    # Not found
    return 1
}

# Validate repository structure
validate_repository() {
    local repo_path="$1"
    
    # Check if directory exists
    if [[ ! -d "$repo_path" ]]; then
        log_error "Repository path does not exist: $repo_path"
        return 1
    fi
    
    # Check for go.mod
    if [[ ! -f "$repo_path/go.mod" ]]; then
        log_error "No go.mod found in repository: $repo_path"
        return 1
    fi
    
    # Check if it's the correct repository
    if ! grep -q "github.com/hashicorp/terraform-provider-azurerm" "$repo_path/go.mod" 2>/dev/null; then
        log_error "Not a terraform-provider-azurerm repository: $repo_path"
        return 1
    fi
    
    # Check for .git directory
    if [[ ! -d "$repo_path/.git" ]]; then
        log_error "Not a git repository: $repo_path"
        return 1
    fi
    
    return 0
}

# Get VS Code settings path based on OS
get_vscode_settings_path() {
    case "$(uname -s)" in
        Darwin*)
            echo "$HOME/Library/Application Support/Code/User/settings.json"
            ;;
        Linux*)
            echo "$HOME/.config/Code/User/settings.json"
            ;;
        CYGWIN*|MINGW32*|MSYS*|MINGW*)
            echo "$APPDATA/Code/User/settings.json"
            ;;
        *)
            # Default to Linux path
            echo "$HOME/.config/Code/User/settings.json"
            ;;
    esac
}

# Get backup directory path
get_backup_directory() {
    case "$(uname -s)" in
        Darwin*)
            echo "$HOME/Library/Application Support/Code/User/.terraform-azurerm-backups"
            ;;
        Linux*)
            echo "$HOME/.config/Code/User/.terraform-azurerm-backups"
            ;;
        CYGWIN*|MINGW32*|MSYS*|MINGW*)
            echo "$APPDATA/Code/User/.terraform-azurerm-backups"
            ;;
        *)
            # Default to Linux path
            echo "$HOME/.config/Code/User/.terraform-azurerm-backups"
            ;;
    esac
}

# Show help information
show_help() {
    cat << 'HELP_EOF'
Terraform AzureRM Provider AI Setup (Modular)
=============================================

USAGE:
  ./install-modular.sh [OPTIONS]

OPTIONS:
  -repository-path <path>    Path to terraform-provider-azurerm repository
  -clean                     Remove all installed files and restore backups
  -auto-approve              Skip interactive approval prompts
  -help                      Show this help message

EXAMPLES:
  ./install-modular.sh                                    # Auto-discover repository
  ./install-modular.sh -repository-path /path/to/repo     # Use specific path
  ./install-modular.sh -auto-approve                      # Non-interactive install
  ./install-modular.sh -clean                             # Remove installation

AI FEATURES:
  - Context-aware code generation and review
  - Azure-specific implementation patterns
  - Testing guidelines and best practices
  - Documentation standards enforcement
  - Error handling and debugging assistance

MODULAR ARCHITECTURE:
  - Clean separation of concerns
  - Maintainable and testable modules
  - Enhanced error handling and logging
  - Improved backup and restore functionality

HELP_EOF
}

# JSON validation function
is_valid_json() {
    local file="$1"
    jq empty "$file" >/dev/null 2>&1
}

# Read JSON property
get_json_property() {
    local file="$1"
    local property="$2"
    
    if [[ -f "$file" ]] && is_valid_json "$file"; then
        jq -r ".$property // empty" "$file" 2>/dev/null
    fi
}

# Check if JSON property exists
json_property_exists() {
    local file="$1"
    local property="$2"
    
    if [[ -f "$file" ]] && is_valid_json "$file"; then
        jq -e "has(\"$property\")" "$file" >/dev/null 2>&1
    else
        return 1
    fi
}
