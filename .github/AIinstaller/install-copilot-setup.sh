#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Main AI Infrastructure Installer for Terraform AzureRM Provider (macOS/Linux)
# Version: 1.0.0
# Description: Interactive installer for AI-powered development infrastructure

set -euo pipefail

# Global variables
VERSION="1.0.0"
BRANCH="exp/terraform_copilot"
SOURCE_REPOSITORY="https://raw.githubusercontent.com/hashicorp/terraform-provider-azurerm"

# Load modules
SCRIPT_DIR="$(dirname "$(realpath "${0}")")"

# Simple color function for pre-module-load errors
print_separator_early() {
    if [[ -t 1 ]] && command -v tput >/dev/null 2>&1; then
        echo -e "\033[0;36m============================================================\033[0m"
    else
        echo "============================================================"
    fi
}

print_error_early() {
    local message="$1"
    if [[ -t 1 ]] && command -v tput >/dev/null 2>&1; then
        echo -e "\033[0;31m[ERROR]\033[0m ${message}"
    else
        echo "[ERROR] ${message}"
    fi
}

# Load all bash modules
for module in "${SCRIPT_DIR}/modules/bash"/*.sh; do
    if [[ -f "${module}" ]]; then
        source "${module}"
    fi
done

# Verify essential modules are loaded
if ! command -v show_usage >/dev/null 2>&1 || ! command -v write_section >/dev/null 2>&1 || ! command -v get_user_profile >/dev/null 2>&1 || ! command -v copy_file >/dev/null 2>&1 || ! command -v verify_installation >/dev/null 2>&1; then
    echo ""
    print_separator_early
    print_error_early "Required modules are missing from ${SCRIPT_DIR}/modules/bash/"
    echo "Please ensure the following modules exist:"
    echo "  - ui.sh"
    echo "  - configparser.sh" 
    echo "  - fileoperations.sh"
    echo "  - validationengine.sh"
    echo ""
    echo "If running from user profile, run bootstrap first:"
    echo "  ./install-copilot-setup.sh -bootstrap"
    exit 1
fi

# Command line options

# Command line options
BOOTSTRAP=false
REPO_DIRECTORY=""
AUTO_APPROVE=false
DRY_RUN=false
VERIFY=false
CLEAN=false
HELP=false

# Function to bootstrap installer
bootstrap_installer() {
    write_section "Bootstrap - Copying Installer to User Profile"
    
    local user_profile
    user_profile="$(get_user_profile)"
    
    # Show path information like PowerShell
    show_path_info "${user_profile}"
    
    # Create directory if needed
    if [[ ! -d "${user_profile}" ]]; then
        if [[ "${DRY_RUN}" != "true" ]]; then
            mkdir -p "${user_profile}"
        fi
    fi
    
    echo -e "${CYAN}Copying installer files from local source repository...${NC}"
    echo ""
    
    local script_dir
    script_dir="$(dirname "$(realpath "${0}")")"
    
    # Statistics tracking
    local files_copied=0
    local files_failed=0
    local total_size=0
    
    # Get bootstrap files from configuration (not hardcoded)
    local manifest_file="${script_dir}/file-manifest.config"
    local bootstrap_files_list
    
    if [[ -f "${manifest_file}" ]]; then
        # Use configuration parser to get INSTALLER_FILES_BOOTSTRAP section
        bootstrap_files_list=$(get_manifest_files "INSTALLER_FILES_BOOTSTRAP" "${manifest_file}")
    else
        write_error_message "Configuration file not found: ${manifest_file}"
        return 1
    fi
    
    # Convert to array and process each file
    local -a bootstrap_files
    while IFS= read -r line; do
        if [[ -n "${line}" ]]; then
            bootstrap_files+=("${line}")
        fi
    done <<< "${bootstrap_files_list}"
    
    # Copy each file according to configuration
    for file_path in "${bootstrap_files[@]}"; do
        # Remove .github/AIinstaller/ prefix to get relative path within installer directory
        local relative_path="${file_path#.github/AIinstaller/}"
        local source_path="${script_dir}/${relative_path}"
        local filename=$(basename "${relative_path}")
        
        # Determine target path based on file type and maintain directory structure
        local target_path
        if [[ "${filename}" == *.psm1 ]]; then
            # PowerShell modules go in modules/powershell/ subdirectory
            if [[ "${DRY_RUN}" != "true" ]]; then
                mkdir -p "${user_profile}/modules/powershell"
            fi
            target_path="${user_profile}/modules/powershell/${filename}"
        elif [[ "${filename}" == *.sh ]] && [[ "${relative_path}" == modules/bash/* ]]; then
            # Bash modules go in modules/bash/ subdirectory
            if [[ "${DRY_RUN}" != "true" ]]; then
                mkdir -p "${user_profile}/modules/bash"
            fi
            target_path="${user_profile}/modules/bash/${filename}"
        else
            # Main files (config, scripts) go directly in target directory
            target_path="${user_profile}/${filename}"
        fi
        
        if copy_file "${source_path}" "${target_path}" "${filename}"; then
            files_copied=$((files_copied + 1))
            # Calculate file size
            if [[ -f "${target_path}" ]]; then
                local file_size
                file_size=$(get_file_size "${target_path}")
                total_size=$((total_size + file_size))
            fi
        else
            files_failed=$((files_failed + 1))
        fi
    done
    
    # Make installer script executable
    if [[ "${DRY_RUN}" != "true" ]]; then
        chmod +x "${user_profile}/install-copilot-setup.sh"
    fi
    
    if [[ "${files_failed}" -eq 0 ]]; then
        # Calculate total size in KB (simple integer division)
        local total_size_kb
        total_size_kb=$((total_size / 1024))
        
        # Use the enhanced bootstrap completion function for consistent output
        show_bootstrap_completion "${files_copied}" "${total_size_kb} KB" "${user_profile}" "$(get_workspace_root)"
    else
        write_error_message "Bootstrap failed with ${files_failed} file(s) failing to copy"
        return 1
    fi
}

# Function to clean installation
clean_installation() {
    local workspace_root
    workspace_root="$(get_workspace_root)"
    
    write_section "Cleaning AI Infrastructure"
    
    # Get all user-facing files from manifest sections
    local cleanup_files
    cleanup_files=($(get_files_for_cleanup "${workspace_root}"))
    
    if [[ ${#cleanup_files[@]} -eq 0 ]]; then
        write_error_message "No files found in manifest to clean"
        return 1
    fi
    
    local files_to_remove=()
    for file in "${cleanup_files[@]}"; do
        files_to_remove+=("${workspace_root}/${file}")
    done
    
    for file in "${files_to_remove[@]}"; do
        if [[ -e "${file}" ]]; then
            if [[ "${DRY_RUN}" == "true" ]]; then
                write_operation_status "[DRY-RUN] Would remove: ${file}" "Info"
            else
                rm -rf "${file}"
                write_operation_status "Removed: ${file}" "Success"
            fi
        fi
    done
    
    show_completion "AI Infrastructure cleanup completed!"
}

# Parse command line arguments
parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -bootstrap)
                BOOTSTRAP=true
                shift
                ;;
            -repo-directory)
                REPO_DIRECTORY="$2"
                shift 2
                ;;
            -auto-approve)
                AUTO_APPROVE=true
                shift
                ;;
            -dry-run)
                DRY_RUN=true
                shift
                ;;
            -verify)
                VERIFY=true
                shift
                ;;
            -clean)
                CLEAN=true
                shift
                ;;
            -help)
                HELP=true
                shift
                ;;
            *)
                write_error_message "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done
}

# Main function
main() {
    parse_arguments "$@"
    
    # Display header and branch detection (matches PowerShell output)
    write_header "Terraform AzureRM Provider - AI Infrastructure Installer" "1.0.0"
    
    # Get branch and workspace information
    local current_branch workspace_root branch_type is_source_branch
    if [[ -n "${REPO_DIRECTORY}" ]]; then
        workspace_root="${REPO_DIRECTORY}"
        current_branch=$(cd "${workspace_root}" && git branch --show-current 2>/dev/null || echo "unknown")
    else
        workspace_root="$(get_workspace_root)"
        current_branch=$(git branch --show-current 2>/dev/null || echo "unknown")
    fi
    
    # Determine branch type early (like PowerShell version)
    # Compare against the actual source branch that contains AI infrastructure
    local source_branch_name="exp/terraform_copilot"  # This is the actual source branch for AI infrastructure
    if [[ "${current_branch}" == "${source_branch_name}" ]]; then
        branch_type="source"
        is_source_branch=true
    else
        branch_type="feature"
        is_source_branch=false
    fi
    
    # Show branch detection
    show_branch_detection "${current_branch}" "${workspace_root}"
    
    # Handle help parameter
    if [[ "${HELP}" == "true" ]]; then
        show_usage
        exit 0
    fi
    
    # Handle bootstrap parameter
    if [[ "${BOOTSTRAP}" == "true" ]]; then
        # Detect if running from user profile directory (incorrect)
        current_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
        if [[ "${current_dir}" == "${HOME}/.terraform-ai-installer" ]]; then
            show_bootstrap_location_error "${current_dir}" "<repo>/.github/AIinstaller/"
            exit 1
        fi
        
        bootstrap_installer
        exit 0
    fi
    
    # Handle verify parameter
    if [[ "${VERIFY}" == "true" ]]; then
        verify_installation "${workspace_root}"
        exit 0
    fi
    
    # Handle clean parameter (feature branch only)
    if [[ "${CLEAN}" == "true" ]]; then
        if [[ "${is_source_branch}" == "true" ]]; then
            echo ""
            print_separator
            echo ""
            write_error_message "Clean operation not available on source branch. This would remove development files."
            echo ""
            exit 1
        fi
        
        clean_installation
        exit 0
    fi
    
    # Default installation - check if this is safe
    if [[ "${is_source_branch}" == "true" ]]; then
        print_separator
        echo ""
        write_error_message "SAFETY CHECK FAILED: Cannot install to source repository directory"
        echo ""
        write_plain "This appears to be the terraform-provider-azurerm source repository."
        write_plain "Installing here would overwrite your local changes with remote files."
        echo ""
        write_plain "${YELLOW}SAFE OPTIONS:${NC}"
        write_plain "  1. Bootstrap installer to user profile:"
        write_plain "     ${0} -bootstrap"
        echo ""
        write_plain "  2. Install to a different repository:"
        write_plain "     ${0} -repo-directory /path/to/target/repository"
        echo ""
        write_plain "For help: ${0} -help"
        echo ""
        exit 1
    fi
    
    # Validate repository if using REPO_DIRECTORY
    if [[ -n "${REPO_DIRECTORY}" ]]; then
        validate_repository "${REPO_DIRECTORY}"
    fi
    
    # Call the install_infrastructure function from fileoperations module
    install_infrastructure "${workspace_root}"
}

# Run main function with all arguments
main "$@"
