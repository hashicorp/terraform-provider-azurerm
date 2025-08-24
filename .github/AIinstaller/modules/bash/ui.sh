#!/usr/bin/env bash
# UI Module for Terraform AzureRM Provider AI Setup (Bash)
# Handles all user interface, output formatting, and user interaction

# Function to calculate spacing for aligned labels
format_aligned_label_spacing() {
    local label="$1"
    local reference="$2"
    
    local label_length=${#label}
    local reference_length=${#reference}
    local spaces_needed=$((reference_length - label_length))
    
    printf "%${spaces_needed}s" ""
}

# Simple format aligned label function
format_aligned_label() {
    local label="$1"
    echo "${label}:"
}

# Color definitions with cross-platform compatibility
# Check if colors are supported
if [[ -t 1 ]] && command -v tput >/dev/null 2>&1; then
    # Terminal supports colors
    export RED='\033[0;31m'
    export GREEN='\033[0;32m'
    export YELLOW='\033[1;33m'
    export BLUE='\033[0;34m'
    export CYAN='\033[0;36m'
    export BOLD='\033[1m'
    export NC='\033[0m' # No Color
else
    # No color support (pipes, non-interactive, etc.)
    export RED=''
    export GREEN=''
    export YELLOW=''
    export BLUE=''
    export CYAN=''
    export BOLD=''
    export NC=''
fi

# Helper function to print colored separator line
print_separator() {
    echo -e "${CYAN}============================================================${NC}"
}

# Function to display main application header
write_header() {
    local title="${1:-Terraform AzureRM Provider - AI Infrastructure Installer}"
    local version="${2:-1.0.0}"
    
    echo ""
    print_separator
    echo -e "${CYAN} ${title}${NC}"
    echo -e "${CYAN} Version: ${version}${NC}"
    print_separator
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
    print_separator
    echo -e "${CYAN} ${title}${NC}"
    print_separator
    echo ""
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
# Function to display file operations (enhanced to match PowerShell output)
show_file_operation() {
    local operation="$1"
    local filename="$2"
    local status="$3"
    
    case "${status}" in
        "OK"|"SUCCESS")
            echo -e "   ${CYAN}${operation}: ${NC}${filename} ${GREEN}[OK]${NC}"
            ;;
        "FAILED"|"ERROR")
            echo -e "   ${CYAN}${operation}: ${NC}${filename} ${RED}[FAILED]${NC}"
            ;;
        "SKIPPED"|"EXISTS")
            echo -e "   ${CYAN}${operation}: ${NC}${filename} ${YELLOW}[SKIPPED]${NC}"
            ;;
        *)
            echo -e "   ${CYAN}${operation}: ${NC}${filename} [${status}]"
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

# Function to show branch detection display (matches PowerShell output)
show_branch_detection() {
    local current_branch="$1"
    local workspace_root="$2"
    
    # Determine branch type for display  
    local branch_label
    if [[ "${current_branch}" == "main" || "${current_branch}" == "exp/terraform_copilot" ]]; then
        branch_label="SOURCE BRANCH DETECTED"
        # SOURCE BRANCH: Cyan label, Green value (matches PowerShell)
        echo -e "${CYAN}${branch_label}: ${GREEN}${current_branch}${NC}"
    else
        branch_label="FEATURE BRANCH DETECTED"
        # FEATURE BRANCH: Cyan label, Yellow value (matches PowerShell)
        echo -e "${CYAN}${branch_label}: ${YELLOW}${current_branch}${NC}"
    fi
    
    # Calculate spacing to align WORKSPACE with branch label
    local workspace_spacing
    workspace_spacing=$(format_aligned_label_spacing "WORKSPACE" "${branch_label}")
    # WORKSPACE: Cyan label, Green path (matches PowerShell)
    echo -e "${CYAN}WORKSPACE${workspace_spacing}: ${GREEN}${workspace_root}${NC}"
}

# Function to show repository information
show_repository_info() {
    local directory="$1"
    
    write_plain "Repository Directory: ${directory}"
    
    # Try to get git branch if available
    if command -v git >/dev/null 2>&1 && [[ -d "${directory}/.git" ]]; then
        local branch
        branch=$(cd "${directory}" && git branch --show-current 2>/dev/null || echo "unknown")
        write_plain "Current Branch: ${branch}"
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

# Function to display completion summary (enhanced to match PowerShell quality)
show_completion_summary() {
    local operation="$1"
    local files_processed="$2"
    local files_succeeded="$3"
    local files_failed="${4:-0}"
    local total_size="${5:-}"
    local install_location="${6:-}"
    local branch_name="${7:-}"
    local branch_type="${8:-feature}"
    
    echo ""
    echo -e "${GREEN}${operation} completed successfully!${NC}"
    echo ""
    
    # Show branch information if provided
    if [[ -n "${branch_name}" ]]; then
        show_branch_detection "${branch_name}" "${branch_type}"
        echo ""
    fi
    
    # Show summary statistics
    echo -e "${CYAN}SUMMARY:${NC}"
    echo -e "  ${GREEN}Files copied${NC} : ${files_succeeded}"
    if [[ "${files_failed}" -gt 0 ]]; then
        echo -e "  ${RED}Files failed${NC} : ${files_failed}"
    fi
    if [[ -n "${total_size}" ]]; then
        echo -e "  ${CYAN}Total size${NC}   : ${total_size}"
    fi
    if [[ -n "${install_location}" ]]; then
        echo -e "  ${CYAN}Location${NC}     : ${install_location}"
    fi
    echo ""
}

# Function to show key-value pairs
show_key_value() {
    local key="$1"
    local value="$2"
    echo -e "${CYAN}${key}: ${NC}${value}"
}

# Function to show next steps (matches PowerShell formatting)
show_next_steps() {
    local steps=("$@")
    
    if [[ ${#steps[@]} -gt 0 ]]; then
        echo -e "${CYAN}NEXT STEPS:${NC}"
        echo ""
        
        for i in "${!steps[@]}"; do
            local step_num=$((i + 1))
            echo -e "  ${step_num}. ${steps[i]}"
        done
        echo ""
    fi
}

# Function to show path information (matches PowerShell output)
show_path_info() {
    local path="$1"
    
    echo -e "${CYAN}PATH: ${YELLOW}${path}${NC}"
    echo ""
}

# Function to show bootstrap completion summary
show_bootstrap_completion() {
    local files_copied="$1"
    local total_size="$2"
    local location="$3"
    local repo_directory="$4"
    
    echo ""
    echo -e "${GREEN}Bootstrap completed successfully!${NC}"
    echo ""
    echo -e "  ${CYAN}Files copied${NC} : ${GREEN}${files_copied}${NC}"
    echo -e "  ${CYAN}Total size${NC}   : ${GREEN}${total_size}${NC}"
    echo -e "  ${CYAN}Location${NC}     : ${YELLOW}${location}${NC}"
    echo ""
    
    echo -e "${CYAN}NEXT STEPS:${NC}"
    echo ""
    echo -e "  ${BOLD}1. Switch to your feature branch:${NC}"
    echo -e "     ${YELLOW}git checkout feature/your-branch-name${NC}"
    echo ""
    echo -e "  ${BOLD}2. Run the installer from your user profile:${NC}"
    echo ""
    echo -e "     ${BOLD}Windows (PowerShell):${NC}"
    echo -e "     ${YELLOW}& \"\$env:USERPROFILE\\.terraform-ai-installer\\install-copilot-setup.ps1\" -RepoDirectory \"${repo_directory}\"${NC}"
    echo ""
    echo -e "     ${BOLD}macOS/Linux (Bash):${NC}"
    echo -e "     ${YELLOW}~/.terraform-ai-installer/install-copilot-setup.sh -repo-directory \"${repo_directory}\"${NC}"
    echo ""
    echo -e "  ${YELLOW}Note: Choose the installer for your platform. Both PowerShell and Bash components${NC}"
    echo -e "        ${YELLOW}are available for maximum flexibility (e.g., VMs, containers, WSL).${NC}"
    echo ""
}

# Function to show divider
show_divider() {
    local char="${1:--}"
    local length="${2:-60}"
    
    printf "%${length}s\n" | tr ' ' "${char}"
}

# Function to show usage information
show_usage() {
    echo ""
    print_separator
    echo ""
    echo "DESCRIPTION:"
    echo "  Interactive installer for AI-powered development infrastructure that enhances"
    echo "  GitHub Copilot with Terraform-specific knowledge, patterns, and best practices."
    echo ""
    echo "USAGE:"
    echo "  ./install-copilot-setup.sh [OPTIONS]"
    echo ""
    echo "AVAILABLE OPTIONS:"
    echo "  -bootstrap        Copy installer to user profile (~/.terraform-ai-installer/)"
    echo "  -verify           Check current workspace status and validate setup"
    echo "  -help             Show this help information"
    echo ""
    echo "EXAMPLES:"
    echo "  Bootstrap installer:"
    echo "    ./install-copilot-setup.sh -bootstrap"
    echo ""
    echo "  Verify setup:"
    echo "    ./install-copilot-setup.sh -verify"
    echo ""
    echo "NEXT STEPS:"
    echo "  1. Run bootstrap to set up the installer in your user profile"
    echo "  2. Switch to your feature branch: git checkout feature/your-branch-name"
    echo "  3. Run installer from user profile to install AI infrastructure"
    echo ""
    echo "For more information, visit: https://github.com/hashicorp/terraform-provider-azurerm"
    echo ""
}

# Function to write plain text (no prefix)
write_plain() {
    local message="$1"
    echo -e "${message}"
}

# Export functions for use in other scripts
export -f write_header write_success write_warning write_error write_info write_section write_plain
export -f show_completion show_file_operation show_error_block show_repository_info
export -f prompt_confirmation show_completion_summary show_key_value show_divider show_usage
export -f show_branch_detection show_path_info show_bootstrap_completion format_aligned_label
export -f format_aligned_label_spacing show_next_steps print_separator
