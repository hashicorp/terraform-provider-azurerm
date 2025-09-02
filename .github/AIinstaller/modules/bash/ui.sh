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

# Color helper functions for consistent output
write_cyan() {
    local message="$1"
    echo -e "${CYAN}${message}${NC}"
}

write_green() {
    local message="$1"
    echo -e "${GREEN}${message}${NC}"
}

write_yellow() {
    local message="$1"
    echo -e "${YELLOW}${message}${NC}"
}

write_white() {
    local message="$1"
    echo -e "${WHITE}${message}${NC}"
}

write_red() {
    local message="$1"
    echo -e "${RED}${message}${NC}"
}

write_blue() {
    local message="$1"
    echo -e "${BLUE}${message}${NC}"
}

write_gray() {
    local message="$1"
    echo -e "${GRAY}${message}${NC}"
}

# Helper functions for common label patterns
write_label() {
    local label="$1"
    local value="$2"
    echo -e "${CYAN}${label}: ${NC}${value}"
}

write_colored_label() {
    local label="$1"
    local value="$2"
    local value_color="$3"
    echo -e "${CYAN}${label}: ${value_color}${value}${NC}"
}

write_section_header() {
    local header="$1"
    echo -e "${CYAN}${header}:${NC}"
}

write_file_operation_status() {
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

# Function to show operation summary (generic function for all operations - matches PowerShell Show-OperationSummary)
show_operation_summary() {
    local operation_name="$1"
    local success="$2"
    local dry_run="${3:-false}"
    shift 3

    # Parse arguments for details and next steps
    local -A details_hash
    local -a ordered_keys  # Array to preserve order
    local longest_key=""
    local next_steps=()
    local parsing_steps=false

    # Process all remaining arguments
    while [[ $# -gt 0 ]]; do
        if [[ "$1" == "--next-steps" ]]; then
            parsing_steps=true
            shift
            continue
        fi

        if [[ "$parsing_steps" == true ]]; then
            next_steps+=("$1")
        elif [[ "$1" =~ ^([^:]+):[[:space:]]*(.+)$ ]]; then
            local key="${BASH_REMATCH[1]}"  # Preserve spaces in key names
            local value="${BASH_REMATCH[2]}"
            details_hash["$key"]="$value"
            ordered_keys+=("$key")  # Preserve insertion order

            # Track longest key for alignment
            if [[ ${#key} -gt ${#longest_key} ]]; then
                longest_key="$key"
            fi
        fi
        shift
    done

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

    echo ""
    echo -e "${color}${completion_message}${NC}"
    echo ""

    # Display details if any exist
    if [[ ${#details_hash[@]} -gt 0 ]]; then
        # Summary section with cyan headers (matches PowerShell structure)
        print_separator 60 "${CYAN}" "="
        write_cyan " ${operation_name^^} SUMMARY:"
        print_separator 60 "${CYAN}" "="
        echo ""
        write_section_header "DETAILS"

        # Define expected order based on operation type
        local -a expected_order
        if [[ "${operation_name,,}" == "verification" ]]; then
            expected_order=("Branch Type" "Target Branch" "Files Verified" "Issues Found" "Location")
        else
            # Installation operation order
            expected_order=("Branch Type" "Target Branch" "Files Installed" "Items Successful" "Total Size" "Location")
        fi

        # Create ordered display array
        local -a display_keys
        # First add keys that are in the expected order
        for expected_key in "${expected_order[@]}"; do
            for key in "${ordered_keys[@]}"; do
                if [[ "$key" == "$expected_key" ]]; then
                    display_keys+=("$key")
                    break
                fi
            done
        done
        # Then add any remaining keys not in expected order
        for key in "${ordered_keys[@]}"; do
            local found=false
            for display_key in "${display_keys[@]}"; do
                if [[ "$key" == "$display_key" ]]; then
                    found=true
                    break
                fi
            done
            if [[ "$found" == false ]]; then
                display_keys+=("$key")
            fi
        done

        # Display each detail with consistent alignment (use ordered display)
        for key in "${display_keys[@]}"; do
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
                write_green "${value}"
            else
                # Text values in yellow
                write_yellow "${value}"
            fi
        done
        echo ""  # Add newline after DETAILS section
    fi

    # Add next steps if provided
    if [[ ${#next_steps[@]} -gt 0 ]]; then
        write_cyan "NEXT STEPS:"
        echo ""
        write_cyan "  1. Switch to your feature branch:"
        write_white "     git checkout feature/your-branch-name"
        echo ""
        write_cyan "  2. Run the installer from your user profile:"
        write_white "     cd \"\$HOME/.terraform-ai-installer\""
        write_white "     ./install-copilot-setup.sh -repo-directory \"<path-to-your-terraform-provider-azurerm>\""
        echo ""
    fi
}

# Function to display next steps after successful bootstrap operation
show_bootstrap_next_steps() {
    local target_directory="${1:-$HOME/.terraform-ai-installer}"

    write_cyan "NEXT STEPS:"
    echo ""
    write_cyan "  1. Switch to your feature branch:"
    write_white "     git checkout feature/your-branch-name"
    echo ""
    write_cyan "  2. Run the installer from your user profile:"
    write_white "     cd \"\$HOME/.terraform-ai-installer\""
    write_white "     ./install-copilot-setup.sh -RepoDirectory \"<path-to-your-terraform-provider-azurerm>\""
    echo ""
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
    write_cyan " ${title}"
    write_cyan " Version: ${version}"
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

    echo ""
    print_separator
}

# Function to display section headers
write_section() {
    local section_title="$1"

    write_cyan " ${section_title}"
    print_separator
    echo ""
}

# Function to display error messages
write_error_message() {
    local message="$1"
    echo -e "${RED} ${message}${NC}" >&2
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
    write_yellow "CORRECT USAGE:"
    write_plain "  Run bootstrap from the source repository:"
    write_plain "  cd ${expected_location}"
    write_plain "  ./install-copilot-setup.sh -bootstrap"
    echo ""
}

# Function to get user profile directory
get_user_profile() {
    # Return the full expanded path (not using ~ shorthand) to match PowerShell behavior
    local expanded_home
    expanded_home="${HOME:-/home/$(whoami)}"
    echo "${expanded_home}/.terraform-ai-installer"
}

# Function to write operation status with consistent formatting
write_operation_status() {
    local message="$1"
    local status="${2:-Info}"

    case "$status" in
        "Success")
            echo -e "${GREEN} [SUCCESS] ${message}${NC}"
            ;;
        "Warning")
            echo -e "${YELLOW} [WARNING] ${message}${NC}"
            ;;
        "Error")
            echo -e "${RED} [ERROR] ${message}${NC}" >&2
            ;;
        "Info"|*)
            echo -e "${BLUE} [INFO] ${message}${NC}"
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
            write_file_operation_status "${operation}" "${formatted_filename}" "OK"
            ;;
        "FAILED"|"ERROR")
            write_file_operation_status "${operation}" "${formatted_filename}" "FAILED"
            ;;
        "SKIPPED"|"EXISTS")
            write_file_operation_status "${operation}" "${formatted_filename}" "SKIPPED"
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
    write_red "ISSUE:"
    write_plain "  ${issue}"
    echo ""

    if [[ -n "${solutions_str}" ]]; then
        write_yellow "SOLUTIONS:"
        # Split solutions by semicolon and display each
        IFS=';' read -ra solutions_array <<< "${solutions_str}"
        for solution in "${solutions_array[@]}"; do
            solution="${solution# }"  # Remove leading space
            write_plain "  - ${solution}"
        done
        echo ""
    fi

    if [[ -n "${example_usage}" ]]; then
        write_green "EXAMPLE:"
        write_plain "  ${example_usage}"
        echo ""
    fi

    if [[ -n "${additional_info}" ]]; then
        write_cyan "ADDITIONAL INFO:"
        write_plain "  ${additional_info}"
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
    write_green "INSTALLATION COMPLETE"
    print_separator 40 "${GREEN}" "="
    echo ""

    # Show branch information if provided
    if [[ -n "${branch_name}" ]]; then
        show_branch_detection "${branch_name}" "${branch_type}"
        echo ""
    fi

    # Show summary statistics
    write_cyan "SUMMARY:"
    write_label "  Files copied" "${files_succeeded}"
    if [[ "${files_failed}" -gt 0 ]]; then
        write_colored_label "  Files failed" "${files_failed}" "${RED}"
    fi
    if [[ -n "${total_size}" ]]; then
        write_label "  Total size" "${total_size}"
    fi
    if [[ -n "${install_location}" ]]; then
        write_label "  Location" "${install_location}"
    fi
    echo ""
}

# Function to show key-value pairs
show_key_value() {
    local key="$1"
    local value="$2"
    write_label "${key}" "${value}"
}

# Function to show next steps (matches PowerShell formatting)
show_next_steps() {
    local steps=("$@")

    if [[ ${#steps[@]} -gt 0 ]]; then
        write_cyan "NEXT STEPS:"
        echo ""

        for i in "${!steps[@]}"; do
            local step_num=$((i + 1))
            write_plain "  ${step_num}. ${steps[i]}"
        done
        echo ""
    fi
}

# Function to show bootstrap completion summary
show_bootstrap_location_error() {
    local current_location="$1"
    local expected_location="$2"

    print_separator
    echo ""
    write_operation_status "Bootstrap must be run from the source repository, not from user profile directory." "Error"
    echo ""
    write_cyan "CORRECT USAGE:"
    write_gray "  cd /path/to/terraform-provider-azurerm"
    write_gray "  ./.github/AIinstaller/install-copilot-setup.sh -bootstrap"
    echo ""
    write_colored_label "CURRENT LOCATION" "${current_location}" "${YELLOW}"
    write_colored_label "EXPECTED LOCATION" "${expected_location}" "${GREEN}"
    echo ""
}

# Function to show divider
show_divider() {
    local char="${1:--}"
    local length="${2:-60}"

    printf "%${length}s\n" | tr ' ' "${char}"
}

# Function to display dynamic help based on branch type and context
show_usage() {
    local branch_type="${1:-feature}"
    local workspace_valid="${2:-true}"
    local workspace_issue="${3:-}"
    local attempted_command="${4:-}"

    echo ""
    write_cyan "DESCRIPTION:"
    write_plain "  Interactive installer for AI-assisted development infrastructure that enhances"
    write_plain "  GitHub Copilot with Terraform-specific knowledge, patterns, and best practices."
    echo ""

    # Dynamic options and examples based on branch type
    case "${branch_type}" in
        "source")
            show_source_branch_help "${attempted_command}"
            ;;
        "feature")
            show_feature_branch_help "${attempted_command}"
            ;;
        *)
            show_unknown_branch_help "${workspace_valid}" "${workspace_issue}" "${attempted_command}"
            ;;
    esac

    write_cyan "For more information, visit: https://github.com/hashicorp/terraform-provider-azurerm"
    echo ""
}

# Function to show source branch specific help
show_source_branch_help() {
    local attempted_command="${1:-}"

    write_cyan "USAGE:"
    write_plain "  ./install-copilot-setup.sh [OPTIONS]"
    echo ""
    write_cyan "AVAILABLE OPTIONS:"
    write_plain "  -bootstrap        Copy installer to user profile (~/.terraform-ai-installer/)"
    write_plain "                    Run this from the source branch to set up for feature branch use"
    write_plain "  -verify           Check current workspace status and validate setup"
    write_plain "  -help             Show this help information"
    echo ""
    write_cyan "EXAMPLES:"
    write_plain "  Bootstrap installer (run from source branch):"
    write_plain "    ./install-copilot-setup.sh -bootstrap"
    echo ""
    write_plain "  Verify setup:"
    write_plain "    ./install-copilot-setup.sh -verify"
    echo ""

    # Show command-specific help if a command was attempted
    if [[ -n "${attempted_command}" ]]; then
        echo ""
        write_yellow "NOTE: You tried to run '${attempted_command}' but this is a source branch."
        write_plain "      Use -bootstrap first to copy the installer to your user profile,"
        write_plain "      then switch to a feature branch for installation operations."
    fi

    write_cyan "BOOTSTRAP WORKFLOW:"
    write_plain "  1. Run -bootstrap from source branch (exp/terraform_copilot) to copy installer to user profile"
    write_plain "  2. Switch to your feature branch: git checkout feature/your-branch-name"
    write_plain "  3. Navigate to user profile: cd ~/.terraform-ai-installer/"
    write_plain "  4. Run installer: ./install-copilot-setup.sh -repo-directory \"/path/to/your/feature/branch\""
    echo ""
}

# Function to show feature branch specific help
show_feature_branch_help() {
    local attempted_command="${1:-}"

    write_cyan "USAGE:"
    write_plain "  ./install-copilot-setup.sh [OPTIONS]"
    echo ""

    write_cyan "AVAILABLE OPTIONS:"
    write_plain "  -repo-directory   Repository path (path to your feature branch directory)"
    write_plain "  -dry-run          Show what would be done without making changes"
    write_plain "  -verify           Check current workspace status and validate setup"
    write_plain "  -clean            Remove AI infrastructure from workspace"
    write_plain "  -help             Show this help information"
    echo ""

    write_cyan "EXAMPLES:"
    write_cyan "  Install AI infrastructure:"
    write_plain "    cd ~/.terraform-ai-installer/"
    write_plain "    ./install-copilot-setup.sh -repo-directory \"/path/to/your/feature/branch\""
    echo ""
    write_cyan "  Dry-Run (preview changes):"
    write_plain "    cd ~/.terraform-ai-installer/"
    write_plain "    ./install-copilot-setup.sh -repo-directory \"/path/to/your/feature/branch\" -dry-run"
    echo ""
    write_cyan "  Clean removal:"
    write_plain "    cd ~/.terraform-ai-installer/"
    write_plain "    ./install-copilot-setup.sh -repo-directory \"/path/to/your/feature/branch\" -clean"
    echo ""

    # Show command-specific help if a command was attempted
    if [[ -n "${attempted_command}" ]]; then
        echo ""
        write_yellow "NOTE: You tried to run '${attempted_command}'."
        case "${attempted_command}" in
            "-repo-directory"*)
                write_plain "      This is correct! You're trying to install AI infrastructure."
                write_plain "      Make sure you're running from ~/.terraform-ai-installer/ directory."
                ;;
            "-bootstrap")
                write_plain "      Bootstrap is for source branches only. Use -repo-directory instead."
                ;;
            *)
                write_plain "      For feature branches, use -repo-directory to specify your workspace."
                ;;
        esac
    fi

    write_cyan "WORKFLOW:"
    write_plain "  1. Navigate to user profile installer directory: cd ~/.terraform-ai-installer/"
    write_plain "  2. Run installer with path to your feature branch"
    write_plain "  3. Start developing with enhanced GitHub Copilot AI features"
    write_plain "  4. Use -clean to remove AI infrastructure when done"
    echo ""
}

# Function to show generic help when branch type cannot be determined
show_unknown_branch_help() {
    local workspace_valid="${1:-true}"
    local workspace_issue="${2:-}"
    local attempted_command="${3:-}"

    # Show workspace issue if detected
    if [[ "${workspace_valid}" != "true" && -n "${workspace_issue}" ]]; then
        write_cyan "WORKSPACE ISSUE DETECTED:"
        write_yellow "  ${workspace_issue}"
        echo ""
        write_cyan "SOLUTION:"

        # Use dynamic command or default to -help
        local command_example="${attempted_command:-"-help"}"

        write_plain "  Navigate to a terraform-provider-azurerm repository, or use the -repo-directory parameter:"
        write_plain "  ./install-copilot-setup.sh -repo-directory \"/path/to/terraform-provider-azurerm\" ${command_example}"
        echo ""
        print_separator
        echo ""
    fi

    write_cyan "GENERAL USAGE:"
    write_plain "  ./install-copilot-setup.sh [OPTIONS]"
    echo ""

    write_cyan "COMMON OPTIONS:"
    write_plain "  -bootstrap        Copy installer to user profile (source branch only)"
    write_plain "  -repo-directory   Repository path (for feature branch operations)"
    write_plain "  -verify           Check current workspace status and validate setup"
    write_plain "  -help             Show this help information"
    echo ""

    write_cyan "GETTING STARTED:"
    write_plain "  1. From source branch: ./install-copilot-setup.sh -bootstrap"
    write_plain "  2. From feature branch: use installer in ~/.terraform-ai-installer/"
    echo ""
}

# Function to display source branch welcome and guidance
show_source_branch_welcome() {
    local branch_name="${1:-exp/terraform_copilot}"

    echo ""
    write_green " WELCOME TO AI-ASSISTED AZURERM TERRAFORM DEVELOPMENT"
    echo ""
    write_cyan "Use the contextual help system above to get started."
    echo ""
}

# Function to show bootstrap failure error with troubleshooting
show_bootstrap_failure_error() {
    local files_failed="$1"
    local user_profile="$2"
    local script_name="$3"

    echo ""
    write_error_message "Bootstrap failed with ${files_failed} file(s) failing to copy"
    echo ""
    write_yellow "TROUBLESHOOTING:"
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
    write_yellow "SOLUTION:"
    write_cyan "Navigate to your terraform-provider-azurerm repository and run:"
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
    write_cyan "Expected structure:"
    write_plain "  From repo root: .github/AIinstaller/install-copilot-setup.sh"
    write_plain "  From installer: install-copilot-setup.sh, modules/"
    echo ""
}

show_repository_directory_required_error() {
    local current_location="$1"

    echo ""
    write_error_message "Repository directory required when running from outside terraform-provider-azurerm repository"
    echo ""
    write_plain "You are running the installer from: ${current_location}"
    write_plain "This is not a terraform-provider-azurerm repository directory."
    echo ""
    write_cyan "Required: Specify target repository with -repo-directory"
    write_plain "  ./install-copilot-setup.sh -repo-directory /path/to/terraform-provider-azurerm"
    echo ""
    write_cyan "Alternative: Bootstrap installer to user profile"
    write_plain "  ./install-copilot-setup.sh -bootstrap"
    echo ""
}

# Function to display safety violation message for source branch operations
show_safety_violation() {
    local branch_name="${1:-exp/terraform_copilot}"
    local operation="${2:-operation}"
    local from_user_profile="${3:-false}"
    local workspace_root="${4:-$PWD}"

    write_red " SAFETY VIOLATION: Cannot perform operations on source branch"
    print_separator 60 "${CYAN}" "="
    echo ""

    if [[ "${from_user_profile}" == "true" ]]; then
        write_yellow " The -RepoDirectory points to the source branch '${branch_name}'."
    else
        write_yellow " You are currently in the source branch '${branch_name}'."
    fi

    write_yellow " Operations other than -Verify, -Help, and -Bootstrap are not allowed on the source branch."
    echo ""
    write_cyan "SOLUTION:"
    write_cyan "  Switch to a feature branch in your target repository:"
    write_plain "    cd \"<path-to-your-terraform-provider-azurerm>\""
    write_plain "    git checkout -b feature/your-branch-name"
    echo ""
    write_cyan "  Then run the installer from your user profile:"
    write_plain "    cd \"\$HOME/.terraform-ai-installer\""
    write_plain "    ./install-copilot-setup.sh -RepoDirectory \"<path-to-your-terraform-provider-azurerm>\""
    echo ""
}

# Function to display workspace validation error message
show_workspace_validation_error() {
    local reason="${1:-Unknown validation error}"
    local from_user_profile="${2:-false}"

    echo ""
    write_error_message "WORKSPACE VALIDATION FAILED: ${reason}"
    echo ""

    # Context-aware error message based on how the script was invoked
    if [[ "${from_user_profile}" == "true" ]]; then
        write_yellow " Please ensure the -repo-directory argument is pointing to a valid GitHub terraform-provider-azurerm repository."
    else
        write_yellow " Please ensure you are running this script from within a terraform-provider-azurerm repository."
    fi
    echo ""
    print_separator
}



# Export all UI functions for use in other scripts
export -f write_cyan write_green write_yellow write_white write_red write_blue write_gray
export -f write_plain write_label write_colored_label write_section_header write_section
export -f write_header write_operation_status
export -f write_error_message write_warning_message write_success_message
export -f write_file_operation_status show_completion_summary show_safety_violation
export -f show_usage show_source_branch_welcome show_workspace_validation_error
export -f show_operation_summary
export -f print_separator get_user_profile format_aligned_label_spacing calculate_max_filename_length
