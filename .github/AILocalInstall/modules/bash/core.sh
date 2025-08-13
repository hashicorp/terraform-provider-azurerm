#!/bin/bash
# Core Module - Clean Architecture
# Core functionality using real patterns from working system

# Color codes for output formatting
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
WHITE='\033[1;37m'
GRAY='\033[0;37m'
NC='\033[0m' # No Color

write_status_message() {
    # Centralized status message function with color coding
    local message="$1"
    local level="${2:-Info}"
    local timestamp="$(date '+%H:%M:%S')"
    
    case "$level" in
        "Success")
            echo -e "${GREEN}[SUCCESS]${NC} $message"
            ;;
        "Error")
            echo -e "${RED}[ERROR]${NC} $message" >&2
            ;;
        "Warning")
            echo -e "${YELLOW}[WARNING]${NC} $message"
            ;;
        "Info")
            echo -e "${BLUE}[INFO]${NC} $message"
            ;;
        *)
            echo "[$level] $message"
            ;;
    esac
}

# Global variables for cross-platform compatibility
SCRIPT_DIR=""
VSCODE_USER_PATH=""

# Cross-platform VS Code user directory detection
get_vscode_user_path() {
    case "$(uname -s)" in
        Darwin)     # macOS
            echo "$HOME/Library/Application Support/Code/User"
            ;;
        Linux)      # Linux
            echo "$HOME/.config/Code/User"
            ;;
        CYGWIN*|MINGW*|MSYS*)  # Windows (Git Bash, MSYS2, etc.)
            echo "$APPDATA/Code/User"
            ;;
        *)
            # Default to Linux path for other Unix-like systems
            echo "$HOME/.config/Code/User"
            ;;
    esac
}

# Initialize global variables
init_globals() {
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    VSCODE_USER_PATH="$(get_vscode_user_path)"
}

find_repository_root() {
    # Finds the terraform-provider-azurerm repository root directory
    local start_path="${1:-$(pwd)}"
    local current_path="$start_path"
    
    while [[ -n "$current_path" && "$current_path" != "/" ]]; do
        local git_path="$current_path/.git"
        local go_mod_path="$current_path/go.mod"
        
        # Check if this is the terraform-provider-azurerm repository
        if [[ -d "$git_path" && -f "$go_mod_path" ]]; then
            if grep -q "module github.com/hashicorp/terraform-provider-azurerm" "$go_mod_path" 2>/dev/null; then
                echo "$current_path"
                return 0
            fi
        fi
        
        current_path="$(dirname "$current_path")"
    done
    
    return 1
}

initialize_clean_environment() {
    # Initializes the clean environment (minimal real implementation)
    local repository_path="$1"
    
    write_status_message "Initializing clean environment..." "Info"
    
    # Basic validation that repo path exists
    if [[ ! -d "$repository_path" ]]; then
        echo "ERROR: Repository path does not exist: $repository_path" >&2
        return 1
    fi
    
    return 0
}

get_installation_config() {
    # Gets installation configuration (real pattern from working system)
    local repository_path="$1"
    
    # Real pattern - read from config files that exist in repository
    local config_path="$repository_path/.github/AILocalInstall"
    
    # Output as key=value pairs that can be sourced
    cat << EOF
REPOSITORY_PATH=$repository_path
CONFIG_PATH=$config_path
BACKUP_ENABLED=true
VSCODE_USER_PATH=$VSCODE_USER_PATH
EOF
}

# Export functions by creating a list of available functions
__CORE_FUNCTIONS="write_status_message find_repository_root initialize_clean_environment get_installation_config get_vscode_user_path init_globals"
