#!/usr/bin/env bash
# UI Module for Terraform AzureRM Provider AI Setup (Bash)
# STREAMLINED VERSION - Contains only functions actually used by main script and dependencies

# Color definitions with cross-platform compatibility
if [[ -t 1 ]] && command -v tput >/dev/null 2>&1; then
    # Terminal supports colors
    export RED='\033[0;31m'
    export GREEN='\033[0;32m'
    export YELLOW='\033[1;33m'
    export BLUE='\033[0;34m'
    export CYAN='\033[0;36m'
    export GRAY='\033[0;37m'
    export WHITE='\033[1;37m'
    export BOLD='\033[1m'
    export NC='\033[0m' # No Color
else
    # No color support (pipes, non-interactive, etc.)
    export RED=''
    export GREEN=''
    export YELLOW=''
    export BLUE=''
    export CYAN=''
    export GRAY=''
    export WHITE=''
    export BOLD=''
    export NC=''
fi

# Global workspace root variable for UI consistency
WORKSPACE_ROOT=""

# Function to show operation summary (generic function for all operations - matches PowerShell Show-OperationSummary)
show_operation_summary() {
    local operation_name="$1"
    local success="$2"
    local dry_run="${3:-false}"
    shift 3
    
    # Parse remaining arguments as details (key:value pairs)
    local -A details_hash
    local longest_key=""
    
    # Process details arguments
    while [[ $# -gt 0 ]]; do
        if [[ "$1" =~ ^([^:]+):[[:space:]]*(.+)$ ]]; then
            local key="${BASH_REMATCH[1]// /}"  # Remove spaces from key
            local value="${BASH_REMATCH[2]}"
            details_hash["$key"]="$value"
            
            # Track longest key for alignment
            if [[ ${#key} -gt ${#longest_key} ]]; then
                longest_key="$key"
            fi
        fi
        shift
    done
    
    echo ""
    
    # Show operation completion with consistent formatting
    local status_text
    local color
    if [[ "$success" == "true" ]]; then
        status_text="completed successfully"
        color="${GREEN}"
    else
        status_text="failed"
        color="${RED}"
    fi
    
    local completion_message
    if [[ "$dry_run" == "true" ]]; then
        completion_message=" ${operation_name} ${status_text} (dry run)"
    else
        completion_message=" ${operation_name} ${status_text}"
    fi
    
    echo -e "${color}${completion_message}${NC}"
    echo ""
    
    # Display details if any exist
    if [[ ${#details_hash[@]} -gt 0 ]]; then
        # Summary section with cyan headers (matches PowerShell structure)
        print_separator 60 "${CYAN}" "="
        echo -e "${CYAN} ${operation_name^^} SUMMARY:${NC}"
        print_separator 60 "${CYAN}" "="
        echo ""
        echo -e "${CYAN}DETAILS:${NC}"
        
        # Display each detail with consistent alignment
        for key in "${!details_hash[@]}"; do
            local value="${details_hash[$key]}"
            
            # Calculate required spacing for alignment
            local label_length=${#key}
            local longest_length=${#longest_key}
            local required_width=$((longest_length - label_length))
            
            if [[ ${required_width} -lt 0 ]]; then
                required_width=0
            fi
            
            # Write key with consistent formatting and proper spacing
            printf "  ${CYAN} %s%*s :${NC} " "${key}" ${required_width} ""
            
            # Determine value color based on content
            if [[ "$value" =~ ^[0-9]+$ ]] || [[ "$value" =~ ^[0-9]+(\.[0-9]+)?[[:space:]]*(KB|MB|GB|TB|B)$ ]]; then
                # Numbers and file sizes in green
                echo -e "${GREEN}${value}${NC}"
            else
                # Text values in yellow
                echo -e "${YELLOW}${value}${NC}"
            fi
        done
    fi
}

# Enhanced show_operation_summary function with next steps support (matches PowerShell exactly)
show_operation_summary_with_steps() {
    local operation_name="$1"
    local success="$2"
    local dry_run="${3:-false}"
    shift 3
    
    # Parse arguments for details and next steps
    local details=()
    local next_steps=()
    local parsing_steps=false
    
    while [[ $# -gt 0 ]]; do
        if [[ "$1" == "--next-steps" ]]; then
            parsing_steps=true
            shift
            continue
        fi
        
        if [[ "$parsing_steps" == true ]]; then
            next_steps+=("$1")
        else
            details+=("$1")
        fi
        shift
    done
    
    # Call the base operation summary function
    show_operation_summary "$operation_name" "$success" "$dry_run" "${details[@]}"
    
    # Add next steps if provided
    if [[ ${#next_steps[@]} -gt 0 ]]; then
        echo ""
        echo -e "${CYAN}NEXT STEPS:${NC}"
        echo ""
        for i in "${!next_steps[@]}"; do
            local step_num=$((i + 1))
            echo -e "  ${CYAN}${step_num}. ${next_steps[i]}${NC}"
            # Add newline between steps, but not after the last step
            if [[ $i -lt $((${#next_steps[@]} - 1)) ]]; then
                echo ""
            fi
        done
        echo ""
    fi
}

# Convenience function for bootstrap completion (uses generic summary with next steps)
show_bootstrap_completion() {
    local files_copied="$1"
    local size_info="$2"
    local user_profile="$3"
    local workspace_root="$4"
    
    # Use the enhanced generic function with next steps
    show_operation_summary_with_steps "Bootstrap" "true" "false" \
        "Files Copied:${files_copied}" \
        "Total Size:${size_info}" \
        "Location:${user_profile}" \
        --next-steps \
        "Switch to your feature branch:\n     ${WHITE}git checkout feature/your-branch-name${NC}" \
        "Run the installer from your user profile:\n     ${WHITE}cd ~/.terraform-ai-installer${NC}\n     ${WHITE}./install-copilot-setup.sh -repo-directory \"<path-to-your-terraform-provider-azurerm>\"${NC}"
}

# Convenience function for installation summary (matches PowerShell pattern)
show_installation_summary() {
    local files_installed="$1"
    local install_location="$2"
    local total_size="${3:-}"
    local success="${4:-true}"
    
    local details=("Files Installed:${files_installed}" "Location:${install_location}")
    if [[ -n "$total_size" ]]; then
        details+=("Total Size:${total_size}")
    fi
    
    show_operation_summary_with_steps "Installation" "$success" "false" \
        "${details[@]}" \
        --next-steps \
        "Restart your development environment to activate AI features" \
        "Check GitHub Copilot extension for enhanced Terraform suggestions"
}

# Convenience function for verification summary (matches PowerShell pattern)
show_verification_summary() {
    local items_checked="$1"
    local items_passed="$2"
    local items_failed="$3"
    local success="$4"
    
    show_operation_summary_with_steps "Verification" "$success" "false" \
        "Items Checked:${items_checked}" \
        "Items Passed:${items_passed}" \
        "Items Failed:${items_failed}" \
        --next-steps \
        "Run installation if components are missing" \
        "Use -clean option to remove installation if needed"
}

# Convenience function for cleanup summary (matches PowerShell pattern)
show_cleanup_summary() {
    local files_removed="$1"
    local cleanup_location="$2"
    local success="${3:-true}"
    
    show_operation_summary_with_steps "Cleanup" "$success" "false" \
        "Files Removed:${files_removed}" \
        "Location:${cleanup_location}" \
        --next-steps \
        "AI infrastructure has been completely removed" \
        "Run bootstrap and installation to restore features"
}

# Helper function to print colored separator line
print_separator() {
    local length="${1:-60}"
    local color="${2:-${CYAN}}"
    local character="${3:-=}"
    
    printf "${color}"
    for ((i=1; i<=length; i++)); do
        printf "${character}"
    done
    printf "${NC}\n"
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

# Function to format aligned labels with proper spacing
format_aligned_label() {
    local label="$1"
    local longest_label="$2"
    
    # Calculate required spacing for alignment (PowerShell style)
    local label_length=${#label}
    local longest_length=${#longest_label}
    local required_width=$((longest_length - label_length))
    
    if [[ ${required_width} -lt 0 ]]; then
        required_width=0
    fi
    
    # Return with leading space and trailing spaces to match PowerShell format
    printf " %s%*s " "${label}" ${required_width} ""
}

# Function to display branch detection with type-based formatting
show_branch_detection() {
    local branch_name="${1:-Unknown}"
    local workspace_root="${2:-}"
    
    # Set global for consistency across UI functions
    WORKSPACE_ROOT="${workspace_root}"
    
    # Determine branch label based on type
    local branch_label
    case "${branch_name}" in
        "exp/terraform_copilot"|"main"|"master")
            branch_label="SOURCE BRANCH DETECTED"
            ;;
        "unknown"|"Unknown")
            branch_label="BRANCH DETECTED"
            ;;
        *)
            branch_label="FEATURE BRANCH DETECTED"
            ;;
    esac
    
    # Use the longest possible label for alignment
    local longest_label="FEATURE BRANCH DETECTED"
    
    # Calculate spacing for branch label
    local branch_label_length=${#branch_label}
    local longest_length=${#longest_label}
    local branch_required_width=$((longest_length - branch_label_length))
    if [[ ${branch_required_width} -lt 0 ]]; then
        branch_required_width=0
    fi
    
    # Display branch information with consistent alignment
    printf "${CYAN} %s%*s : ${NC}${YELLOW}%s${NC}\n" "${branch_label}" ${branch_required_width} "" "${branch_name}"
    
    # Dynamic workspace label with proper alignment and colors
    if [[ -n "${workspace_root}" ]]; then
        local workspace_label="WORKSPACE"
        local workspace_label_length=${#workspace_label}
        local workspace_required_width=$((longest_length - workspace_label_length))
        if [[ ${workspace_required_width} -lt 0 ]]; then
            workspace_required_width=0
        fi
        
        printf "${CYAN} %s%*s : ${NC}${GREEN}%s${NC}\n" "${workspace_label}" ${workspace_required_width} "" "${workspace_root}"
    fi
}

# Function to display section headers
write_section() {
    local section_title="$1"
    
    echo ""
    print_separator
    echo -e "${CYAN} ${section_title}${NC}"
    print_separator
    echo ""
}

# Function to display error messages
write_error_message() {
    local message="$1"
    echo -e "${RED}${message}${NC}" >&2
}

# Function to display warning messages
write_warning_message() {
    local message="$1"
    echo -e "${YELLOW}${message}${NC}"
}

# Function to display success messages
write_success_message() {
    local message="$1"
    echo -e "${GREEN}${message}${NC}"
}

# Function to display plain messages
write_plain() {
    local message="$1"
    echo "${message}"
}

# Function to show path information during bootstrap
show_path_info() {
    local user_profile="$1"
    
    echo "Target Directory: ${user_profile}"
    
    # Show if directory exists and is writable
    if [[ -d "${user_profile}" ]]; then
        echo "Directory status: ${GREEN}exists${NC}"
    else
        echo "Directory status: ${YELLOW}will be created${NC}"
    fi
    
    # Check write permissions
    local parent_dir
    parent_dir="$(dirname "${user_profile}")"
    if [[ -w "${parent_dir}" ]]; then
        echo "Write access: ${GREEN}confirmed${NC}"
    else
        echo "Write access: ${RED}denied${NC}"
    fi
}

# Function to show bootstrap location error
show_bootstrap_location_error() {
    local current_dir="$1"
    local expected_location="$2"
    
    echo ""
    print_separator
    echo ""
    write_error_message "BOOTSTRAP LOCATION ERROR: Bootstrap cannot be run from user profile"
    echo ""
    write_plain "Bootstrap is designed to copy files TO the user profile, not FROM it."
    write_plain "Current location: ${current_dir}"
    echo ""
    write_plain "${YELLOW}CORRECT USAGE:${NC}"
    write_plain "  Run bootstrap from the source repository:"
    write_plain "  cd ${expected_location}"
    write_plain "  ./install-copilot-setup.sh -bootstrap"
    echo ""
}

# Function to get user profile directory
get_user_profile() {
    echo "${HOME}/.terraform-ai-installer"
}

# Function to display success message
write_success() {
    local message="$1"
    local prefix="${2:-[SUCCESS]}"
    
    if [[ "$prefix" == "[SUCCESS]" ]]; then
        write_operation_status "$message" "Success"
    else
        echo -e "${GREEN}${prefix} ${message}${NC}"
    fi
}

# Function to display warning message
write_warning() {
    local message="$1"
    local prefix="${2:-[WARNING]}"
    
    if [[ "$prefix" == "[WARNING]" ]]; then
        write_operation_status "$message" "Warning"
    else
        echo -e "${YELLOW}${prefix} ${message}${NC}"
    fi
}

# Function to display error message
write_error() {
    local message="$1"
    local prefix="${2:-[ERROR]}"
    
    if [[ "$prefix" == "[ERROR]" ]]; then
        write_operation_status "$message" "Error"
    else
        echo -e "${RED}${prefix} ${message}${NC}" >&2
    fi
}

# Function to display info message
write_info() {
    local message="$1"
    local prefix="${2:-[INFO]}"
    
    if [[ "$prefix" == "[INFO]" ]]; then
        write_operation_status "$message" "Info"
    else
        echo -e "${BLUE}${prefix} ${message}${NC}"
    fi
}

# Unified message functions (matches PowerShell Write-WarningMessage, Write-ErrorMessage, etc.)
write_warning_message() {
    local message="$1"
    echo -e "${YELLOW}[WARNING] ${message}${NC}"
}

write_error_message() {
    local message="$1"
    echo -e "${RED}[ERROR] ${message}${NC}" >&2
}

write_verbose_message() {
    local message="$1"
    echo -e "${BLUE}[VERBOSE] ${message}${NC}"
}

# Function to write operation status with consistent formatting
write_operation_status() {
    local message="$1"
    local status="${2:-Info}"
    
    case "$status" in
        "Success")
            echo -e "${GREEN}[SUCCESS] ${message}${NC}"
            ;;
        "Warning")
            echo -e "${YELLOW}[WARNING] ${message}${NC}"
            ;;
        "Error")
            echo -e "${RED}[ERROR] ${message}${NC}" >&2
            ;;
        "Info"|*)
            echo -e "${BLUE}[INFO] ${message}${NC}"
            ;;
    esac
}

# Function to format aligned label spacing (for consistent alignment)
format_aligned_label_spacing() {
    local label="$1"
    local reference_label="$2"
    
    # Calculate spacing needed to align labels
    local label_len=${#label}
    local ref_len=${#reference_label}
    local spaces_needed=$((ref_len - label_len))
    
    if [[ $spaces_needed -gt 0 ]]; then
        printf "%*s" $spaces_needed ""
    fi
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

# Function to calculate maximum filename length for dynamic spacing (matches PowerShell)
calculate_max_filename_length() {
    local -a filenames=("$@")
    local max_length=0
    
    for filename in "${filenames[@]}"; do
        local length=${#filename}
        if [[ $length -gt $max_length ]]; then
            max_length=$length
        fi
    done
    
    echo $max_length
}

# Function to display file operation status
# Function to display file operations (enhanced to match PowerShell output with dynamic padding)
show_file_operation() {
    local operation="$1"
    local filename="$2"
    local status="$3"
    local max_length="$4"  # Required parameter - no default to ensure dynamic calculation
    
    # Align filename to match PowerShell format using dynamic length
    local formatted_filename
    formatted_filename=$(printf "%-${max_length}s" "${filename}")
    
    case "${status}" in
        "OK"|"SUCCESS")
            echo -e "   ${CYAN}${operation}: ${NC}${formatted_filename} ${GREEN}[OK]${NC}"
            ;;
        "FAILED"|"ERROR")
            echo -e "   ${CYAN}${operation}: ${NC}${formatted_filename} ${RED}[FAILED]${NC}"
            ;;
        "SKIPPED"|"EXISTS")
            echo -e "   ${CYAN}${operation}: ${NC}${formatted_filename} ${YELLOW}[SKIPPED]${NC}"
            ;;
        *)
            echo -e "   ${CYAN}${operation}: ${NC}${formatted_filename} [${status}]"
            ;;
    esac
}

# Function to display error block with solutions
show_error_block() {
    local issue="$1"
    local solutions_str="$2"
    local example_usage="${3:-}"
    local additional_info="${4:-}"
    
    echo ""
    echo -e "${RED}ISSUE:${NC}"
    echo "  ${issue}"
    echo ""
    
    if [[ -n "${solutions_str}" ]]; then
        echo -e "${YELLOW}SOLUTIONS:${NC}"
        # Split solutions by semicolon and display each
        IFS=';' read -ra solutions_array <<< "${solutions_str}"
        for solution in "${solutions_array[@]}"; do
            solution="${solution# }"  # Remove leading space
            echo "  â€¢ ${solution}"
        done
        echo ""
    fi
    
    if [[ -n "${example_usage}" ]]; then
        echo -e "${GREEN}EXAMPLE:${NC}"
        echo "  ${example_usage}"
        echo ""
    fi
    
    if [[ -n "${additional_info}" ]]; then
        echo -e "${CYAN}ADDITIONAL INFO:${NC}"
        echo "  ${additional_info}"
        echo ""
    fi
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
    echo -e "${GREEN}INSTALLATION COMPLETE${NC}"
    print_separator 40 "${GREEN}" "="
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
# Function to show bootstrap location error (matches PowerShell Show-BootstrapLocationError)
show_bootstrap_location_error() {
    local current_location="$1"
    local expected_location="$2"
    
    print_separator
    echo ""
    write_operation_status "Bootstrap must be run from the source repository, not from user profile directory." "Error"
    echo ""
    echo -e "${CYAN}CORRECT USAGE:${NC}"
    echo -e "  ${GRAY}cd /path/to/terraform-provider-azurerm${NC}"
    echo -e "  ${GRAY}./.github/AIinstaller/install-copilot-setup.sh -bootstrap${NC}"
    echo ""
    echo -e "${CYAN}CURRENT LOCATION: ${YELLOW}${current_location}${NC}"
    echo -e "${CYAN}EXPECTED LOCATION: ${GREEN}${expected_location}${NC}"
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

# Function to display detailed usage information
show_usage() {
    cat << 'EOF'
Terraform AzureRM Provider - AI Infrastructure Installer (macOS/Linux)

PURPOSE:
Enables AI-powered development features for HashiCorp's Terraform AzureRM Provider,
including GitHub Copilot integration and intelligent code suggestions.

USAGE:
  ./install-copilot-setup.sh [options]

REQUIRED OPERATIONS:
  -bootstrap          Copy installer to user profile for feature branch use
                     Must be run first from source repository
                     Creates: ~/.terraform-ai-installer/

VERIFICATION & CLEANUP:
  -verify            Check current installation state
  -clean             Remove AI infrastructure (feature branches only)

INSTALLER OPTIONS:
  -auto-approve      Overwrite existing files without prompting
  -dry-run           Show what would be done without making changes
  -repo-directory    Path to repository directory (for remote operations)

HELP:
  -help              Display this help message

SECURITY NOTES:
- Branch protection prevents accidental overwrites of source repository
- Bootstrap copies installer files locally for safe feature branch usage
- All operations include verification and rollback capabilities

EXAMPLES:
  # Initial setup (run once from source repository)
  ./install-copilot-setup.sh -bootstrap

  # Install to feature branch (run from user profile)
  ~/.terraform-ai-installer/install-copilot-setup.sh

  # Check installation state
  ~/.terraform-ai-installer/install-copilot-setup.sh -verify

  # Remove from feature branch
  ~/.terraform-ai-installer/install-copilot-setup.sh -clean

For more information, see: .github/AIinstaller/README.md
EOF
}

# Function to show source repository safety error
show_source_repository_safety_error() {
    local script_name="$1"
    
    print_separator
    echo ""
    write_error_message "SAFETY CHECK FAILED: Cannot install to source repository directory"
    echo ""
    write_plain "This appears to be the terraform-provider-azurerm source repository."
    write_plain "Installing here would overwrite your local changes with remote files."
    echo ""
    write_plain "${YELLOW}SAFE OPTIONS:${NC}"
    write_plain "  1. Bootstrap installer to user profile:"
    write_plain "     ${script_name} -bootstrap"
    echo ""
    write_plain "  2. Install to a different repository:"
    write_plain "     ${script_name} -repo-directory /path/to/target/repository"
    echo ""
    write_plain "For help: ${script_name} -help"
    echo ""
}

# Function to show clean operation unavailable error for source branch
show_clean_unavailable_on_source_error() {
    echo ""
    print_separator
    echo ""
    write_error_message "Clean operation not available on source branch"
    echo ""
    write_plain "This would remove development files from the source repository."
    write_plain "Use clean only on feature branches to remove AI infrastructure."
    echo ""
}

# Function to show bootstrap repository validation error
show_bootstrap_repository_validation_error() {
    local workspace_root="$1"
    
    write_error_message "Bootstrap must be run from terraform-provider-azurerm repository"
    echo ""
    write_plain "Expected to find go.mod with terraform-provider-azurerm content"
    write_plain "Current directory: ${workspace_root}"
}

# Function to show bootstrap failure error with troubleshooting
show_bootstrap_failure_error() {
    local files_failed="$1"
    local user_profile="$2"
    local script_name="$3"
    
    echo ""
    write_error_message "Bootstrap failed with ${files_failed} file(s) failing to copy"
    echo ""
    write_plain "${YELLOW}TROUBLESHOOTING:${NC}"
    write_plain "  1. Check file permissions in source directory"
    write_plain "  2. Ensure sufficient disk space in ${user_profile}"
    write_plain "  3. Verify you're running from the correct directory"
    echo ""
    write_plain "For help: ${script_name} -help"
    echo ""
}

# Function to write plain text (no prefix)
write_plain() {
    local message="$1"
    echo -e "${message}"
}

# Function to show bootstrap location error
show_bootstrap_location_error() {
    local current_location="$1"
    local expected_location="$2"
    
    echo ""
    print_separator
    echo ""
    write_error_message "Bootstrap must be run from the source repository, not user profile"
    echo ""
    write_plain "Current location: ${current_location}"
    write_plain "Expected location: ${expected_location}"
    echo ""
    write_plain "${YELLOW}SOLUTION:${NC}"
    write_plain "Navigate to your terraform-provider-azurerm repository and run:"
    write_plain "  cd /path/to/terraform-provider-azurerm/.github/AIinstaller"
    write_plain "  ./install-copilot-setup.sh -bootstrap"
    echo ""
}

# Function to show bootstrap directory validation error
show_bootstrap_directory_validation_error() {
    local current_location="$1"
    
    echo ""
    write_error_message "Bootstrap must be run from terraform-provider-azurerm root or .github/AIinstaller directory"
    echo ""
    write_plain "Current location: ${current_location}"
    write_plain "Expected structure:"
    write_plain "  From repo root: .github/AIinstaller/install-copilot-setup.sh"
    write_plain "  From installer: install-copilot-setup.sh, modules/"
    echo ""
}

# Export functions for use in other scripts
export -f write_header write_operation_status write_success write_warning write_error write_info write_section write_plain
export -f write_warning_message write_error_message write_verbose_message
export -f show_completion show_file_operation show_error_block show_repository_info
export -f prompt_confirmation show_completion_summary show_key_value show_divider show_usage
export -f show_branch_detection show_path_info show_bootstrap_completion show_bootstrap_location_error format_aligned_label
export -f format_aligned_label_spacing show_next_steps print_separator calculate_max_filename_length show_bootstrap_directory_validation_error
export -f show_operation_summary show_operation_summary_with_steps show_installation_summary show_verification_summary show_cleanup_summary
