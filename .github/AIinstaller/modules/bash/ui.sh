#!/usr/bin/env bash
# UI Module for Terraform AzureRM Provider AI Setup (Bash)
# Handles all user interface, output formatting, and user interaction

# Color definitions
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Function to display main application header
write_header() {
    local title="${1:-Terraform AzureRM Provider - AI Infrastructure Installer}"
    local version="${2:-1.0.0}"
    
    echo ""
    echo -e "${CYAN}============================================================${NC}"
    echo -e "${CYAN} ${title}${NC}"
    echo -e "${CYAN} Version: ${version}${NC}"
    echo -e "${CYAN}============================================================${NC}"
    echo ""
}

# Function to display success message
write_success() {
    local message="$1"
    local prefix="${2:-[SUCCESS]}"
    
    echo -e "${GREEN}${prefix} ${message}${NC}"
}

# Function to display warning message
write_warning() {
    local message="$1"
    local prefix="${2:-[WARNING]}"
    
    echo -e "${YELLOW}${prefix} ${message}${NC}"
}

# Function to display error message
write_error() {
    local message="$1"
    local prefix="${2:-[ERROR]}"
    
    echo -e "${RED}${prefix} ${message}${NC}" >&2
}

# Function to display info message
write_info() {
    local message="$1"
    local prefix="${2:-[INFO]}"
    
    echo -e "${BLUE}${prefix} ${message}${NC}"
}

# Function to display section header
write_section() {
    local title="$1"
    
    echo ""
    echo -e "${BOLD}[SECTION] ${title}${NC}"
    echo "-------------------------------------------------------"
}

# Function to show percentage completion
show_completion() {
    local current="$1"
    local total="$2"
    local description="$3"
    
    local percentage=$(( (current * 100) / total ))
    
    printf "${BLUE}[%3d%%]${NC} %s\n" "${percentage}" "${description}"
}

# Function to display file operation status
show_file_operation() {
    local operation="$1"
    local filename="$2"
    local status="$3"
    
    case "${status}" in
        "OK"|"SUCCESS")
            echo -e "  ${operation}: ${filename} ${GREEN}[${status}]${NC}"
            ;;
        "FAILED"|"ERROR")
            echo -e "  ${operation}: ${filename} ${RED}[${status}]${NC}"
            ;;
        "SKIPPED"|"EXISTS")
            echo -e "  ${operation}: ${filename} ${YELLOW}[${status}]${NC}"
            ;;
        *)
            echo -e "  ${operation}: ${filename} [${status}]"
            ;;
    esac
}

# Function to display error block with solutions
show_error_block() {
    local issue="$1"
    shift
    local solutions=("$@")
    
    echo ""
    echo -e "${RED}${BOLD}ERROR DETAILS${NC}"
    echo -e "${RED}=============${NC}"
    echo -e "${RED}Issue:${NC} ${issue}"
    echo ""
    echo -e "${YELLOW}${BOLD}Solutions:${NC}"
    
    for i in "${!solutions[@]}"; do
        echo -e "${YELLOW}  $((i+1)). ${solutions[i]}${NC}"
    done
    echo ""
}

# Function to show repository information
show_repository_info() {
    local directory="$1"
    
    write_info "Repository Directory: ${directory}"
    
    # Try to get git branch if available
    if command -v git >/dev/null 2>&1 && [[ -d "${directory}/.git" ]]; then
        local branch
        branch=$(cd "${directory}" && git branch --show-current 2>/dev/null || echo "unknown")
        write_info "Current Branch: ${branch}"
    fi
}

# Function to prompt user for confirmation
prompt_confirmation() {
    local message="$1"
    local default="${2:-n}"
    
    local prompt_text
    if [[ "${default}" == "y" ]]; then
        prompt_text="${message} [Y/n]: "
    else
        prompt_text="${message} [y/N]: "
    fi
    
    echo -n -e "${YELLOW}${prompt_text}${NC}"
    read -r response
    
    # Use default if no response
    if [[ -z "${response}" ]]; then
        response="${default}"
    fi
    
    case "${response}" in
        [Yy]|[Yy][Ee][Ss])
            return 0
            ;;
        *)
            return 1
            ;;
    esac
}

# Function to display completion summary
show_completion_summary() {
    local operation="$1"
    local files_processed="$2"
    local files_succeeded="$3"
    local files_failed="$4"
    
    echo ""
    echo -e "${BOLD}OPERATION SUMMARY${NC}"
    echo "================="
    echo "Operation: ${operation}"
    echo "Files Processed: ${files_processed}"
    echo -e "Files Succeeded: ${GREEN}${files_succeeded}${NC}"
    
    if [[ "${files_failed}" -gt 0 ]]; then
        echo -e "Files Failed: ${RED}${files_failed}${NC}"
    else
        echo -e "Files Failed: ${GREEN}0${NC}"
    fi
    echo ""
}

# Function to display key-value pairs in a formatted way
show_key_value() {
    local key="$1"
    local value="$2"
    local max_key_length="${3:-20}"
    
    printf "  %-${max_key_length}s : %s\n" "${key}" "${value}"
}

# Function to show divider
show_divider() {
    local char="${1:--}"
    local length="${2:-60}"
    
    printf "%${length}s\n" | tr ' ' "${char}"
}

# Function to show usage information
show_usage() {
    cat << EOF
AI Infrastructure Installer for Terraform AzureRM Provider (macOS/Linux)

USAGE:
    ./install-copilot-setup.sh [OPTIONS]

OPTIONS:
    -bootstrap          Copy installer to user profile for feature branch use
    -repo-directory     Path to the repository directory for git operations
    -auto-approve       Overwrite existing files without prompting
    -dry-run           Show what would be done without making changes
    -verify            Check the current state of the workspace
    -clean             Remove AI infrastructure from the workspace
    -help              Show this help message

EXAMPLES:
    # Bootstrap installer to user profile (run once)
    ./install-copilot-setup.sh -bootstrap

    # Install AI infrastructure from user profile
    ~/.terraform-ai-installer/install-copilot-setup.sh -repo-directory "/path/to/terraform-provider-azurerm"

    # Verify current installation
    ./install-copilot-setup.sh -verify

    # Clean installation
    ./install-copilot-setup.sh -clean

EOF
}

# Export functions for use in other scripts
export -f write_header write_success write_warning write_error write_info write_section
export -f show_completion show_file_operation show_error_block show_repository_info
export -f prompt_confirmation show_completion_summary show_key_value show_divider show_usage
