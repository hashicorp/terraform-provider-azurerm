#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Main AI Infrastructure Installer for Terraform AzureRM Provider (macOS/Linux)
# Version: 1.0.0
# Description: Interactive installer for AI-powered development infrastructure

#requires bash 4.0+

set -euo pipefail

# ============================================================================
# PARAMETER DEFINITIONS
# ============================================================================

# Global variables
VERSION="1.0.0"
BRANCH="exp/terraform_copilot"
SOURCE_REPOSITORY="https://raw.githubusercontent.com/hashicorp/terraform-provider-azurerm"

# Command line parameters with help text
BOOTSTRAP=false           # Copy installer to user profile for feature branch use
REPO_DIRECTORY=""         # Path to the repository directory for git operations (when running from user profile)
AUTO_APPROVE=false        # Overwrite existing files without prompting
DRY_RUN=false             # Show what would be done without making changes
VERIFY=false              # Check the current state of the workspace
CLEAN=false               # Remove AI infrastructure from the workspace
HELP=false                # Show detailed help information

# ============================================================================
# MODULE LOADING - This must succeed or the script cannot continue
# ============================================================================

# Get script directory with robust detection
get_script_directory() {
    local source="${BASH_SOURCE[0]}"
    while [[ -L "${source}" ]]; do
        local dir="$(cd -P "$(dirname "${source}")" && pwd)"
        source="$(readlink "${source}")"
        [[ ${source} != /* ]] && source="${dir}/${source}"
    done
    echo "$(cd -P "$(dirname "${source}")" && pwd)"
}

get_modules_path() {
    local script_directory="$1"
    
    # Simple logic: modules are always in the same relative location
    local modules_path="${script_directory}/modules/bash"
    
    # If not found, try from workspace root (for direct repo execution)
    if [[ ! -d "${modules_path}" ]]; then
        local current_path="${script_directory}"
        while [[ -n "${current_path}" && "${current_path}" != "$(dirname "${current_path}")" ]]; do
            if [[ -f "${current_path}/go.mod" ]]; then
                modules_path="${current_path}/.github/AIinstaller/modules/bash"
                break
            fi
            current_path="$(dirname "${current_path}")"
        done
    fi
    
    echo "${modules_path}"
}

import_required_modules() {
    local modules_path="$1"
    
    # Define all required modules in dependency order
    local modules=(
        "configparser"
        "ui" 
        "validationengine"
        "fileoperations"
    )
    
    # Load each module cleanly
    for module in "${modules[@]}"; do
        local module_path="${modules_path}/${module}.sh"
        
        if [[ ! -f "${module_path}" ]]; then
            echo ""
            echo "============================================================"
            echo "[ERROR] Required module '${module}' not found at: ${module_path}"
            echo ""
            echo "If running from user profile, run bootstrap first:"
            echo "  $0 -bootstrap"
            echo ""
            return 1
        fi
        
        if ! source "${module_path}"; then
            echo ""
            echo "============================================================"
            echo "[ERROR] Failed to import module '${module}': ${module_path}"
            echo ""
            return 1
        fi
    done
    
    # Verify critical functions are available
    local required_functions=(
        "get_manifest_config"
        "get_installer_config" 
        "write_header"
        "verify_installation"
    )
    
    for func in "${required_functions[@]}"; do
        if ! command -v "${func}" >/dev/null 2>&1; then
            echo ""
            echo "============================================================"
            echo "[ERROR] Required function '${func}' not available after module loading"
            echo ""
            return 1
        fi
    done
    
    return 0
}

# Get script directory and load modules
SCRIPT_DIR="$(get_script_directory)"
MODULES_PATH="$(get_modules_path "${SCRIPT_DIR}")"

# Import all required modules or exit with error
if ! import_required_modules "${MODULES_PATH}"; then
    exit 1
fi
# ============================================================================
# IMPLEMENTATION FUNCTIONS
# ============================================================================

# Function to get user profile directory
get_user_profile() {
    echo "${HOME}/.terraform-ai-installer"
}

# Function to bootstrap installer to user profile
bootstrap_installer() {
    write_section "Bootstrap - Copying Installer to User Profile"
    
    # Validate that we're running from the right location
    local current_location
    current_location="$(pwd)"
    local user_profile
    user_profile="$(get_user_profile)"
    
    # Prevent bootstrap from user profile (circular operation)
    if [[ "${current_location}" == "${user_profile}"* ]]; then
        show_bootstrap_location_error "${current_location}" "terraform-provider-azurerm/.github/AIinstaller"
        return 1
    fi
    
    # Detect if we're in the repo root and adjust SCRIPT_DIR accordingly
    local installer_dir
    if [[ -f ".github/AIinstaller/install-copilot-setup.sh" ]] && [[ -d ".github/AIinstaller/modules" ]]; then
        # Running from repo root - adjust SCRIPT_DIR to point to installer directory
        installer_dir="$(pwd)/.github/AIinstaller"
        SCRIPT_DIR="${installer_dir}"
    elif [[ -f "install-copilot-setup.sh" ]] && [[ -d "modules" ]]; then
        # Running from installer directory - use current SCRIPT_DIR
        installer_dir="${SCRIPT_DIR}"
    else
        show_bootstrap_directory_validation_error "${current_location}"
        return 1
    fi
    
    # Create directory if needed
    if [[ ! -d "${user_profile}" ]]; then
        if [[ "${DRY_RUN}" != "true" ]]; then
            mkdir -p "${user_profile}"
        fi
    fi
    
    # Delegate file operations to the file operations module
    local manifest_file="${SCRIPT_DIR}/file-manifest.config"
    if [[ ! -f "${manifest_file}" ]]; then
        write_error_message "Configuration file not found: ${manifest_file}"
        return 1
    fi
    
    # Perform bootstrap operation using file operations module
    if bootstrap_files_to_profile "${SCRIPT_DIR}" "${user_profile}" "${manifest_file}"; then
        # Calculate total size in KB (simple integer division)
        local total_size_kb
        total_size_kb=$((BOOTSTRAP_STATS_TOTAL_SIZE / 1024))
        
        # Get current branch for intelligent next steps
        local current_branch
        current_branch=$(git branch --show-current 2>/dev/null || echo "unknown")
        
        # Use the enhanced bootstrap completion function for consistent output
        show_bootstrap_completion "${BOOTSTRAP_STATS_FILES_COPIED}" "${total_size_kb} KB" "${user_profile}" "$(get_workspace_root)" "${current_branch}"
    else
        show_bootstrap_failure_error "${BOOTSTRAP_STATS_FILES_FAILED}" "${user_profile}" "${0}"
        return 1
    fi
}

# Function to clean installation
clean_installation() {
    local workspace_root
    if [[ -n "${REPO_DIRECTORY}" ]]; then
        workspace_root="${REPO_DIRECTORY}"
    else
        workspace_root="$(get_workspace_root)"
    fi
    
    write_section "Cleaning AI Infrastructure"
    
    # Use the fileoperations module function for cleanup
    clean_infrastructure "${workspace_root}"
}

# ============================================================================
# COMMAND LINE ARGUMENT PROCESSING
# ============================================================================

# Parse command line arguments
parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -bootstrap)
                BOOTSTRAP=true
                shift
                ;;
            -repo-directory)
                if [[ -z "$2" || "$2" == -* ]]; then
                    write_error_message "Option -repo-directory requires a directory path"
                    exit 1
                fi
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
                echo ""
                echo "Use -help for usage information"
                exit 1
                ;;
        esac
    done
}

# ============================================================================
# MAIN EXECUTION LOGIC
# ============================================================================

# Main function with enhanced branch protection logic
main() {
    parse_arguments "$@"
    
    # Display header and branch detection (matches PowerShell output)
    write_header "Terraform AzureRM Provider - AI Infrastructure Installer" "${VERSION}"
    
    # Get branch and workspace information with proper error handling
    local current_branch workspace_root branch_type is_source_branch
    if [[ -n "${REPO_DIRECTORY}" ]]; then
        workspace_root="${REPO_DIRECTORY}"
        if [[ -d "${workspace_root}/.git" ]]; then
            current_branch=$(cd "${workspace_root}" && git branch --show-current 2>/dev/null || echo "unknown")
        else
            current_branch="unknown"
        fi
    else
        workspace_root="$(get_workspace_root)"
        current_branch=$(git branch --show-current 2>/dev/null || echo "unknown")
    fi
    
    # Determine branch type early (like PowerShell version)
    # Compare against the actual source branch that contains AI infrastructure
    local source_branch_name="${BRANCH}"  # Use the configured branch name
    if [[ "${current_branch}" == "${source_branch_name}" ]]; then
        branch_type="source"
        is_source_branch=true
    else
        branch_type="feature"
        is_source_branch=false
    fi
    
    # Show branch detection with consistent formatting
    show_branch_detection "${current_branch}" "${workspace_root}"
    
    # Handle help parameter first
    if [[ "${HELP}" == "true" ]]; then
        # Determine branch type for dynamic help
        if [[ "${is_source_branch}" == "true" ]]; then
            show_usage "source" "true" ""
        else
            show_usage "feature" "true" ""
        fi
        exit 0
    fi
    
    # Handle bootstrap parameter with proper safety checks
    if [[ "${BOOTSTRAP}" == "true" ]]; then
        # Enhanced location validation - prevent bootstrap from user profile
        current_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
        local user_profile_dir
        user_profile_dir="$(get_user_profile)"
        
        if [[ "${current_dir}" == "${user_profile_dir}" ]]; then
            show_bootstrap_location_error "${current_dir}" "<repo>/.github/AIinstaller/"
            exit 1
        fi
        
        # Verify we're in the source repository for bootstrap
        if [[ ! -f "${workspace_root}/go.mod" ]] || ! grep -q "terraform-provider-azurerm" "${workspace_root}/go.mod" 2>/dev/null; then
            show_bootstrap_repository_validation_error "${workspace_root}"
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
            show_clean_unavailable_on_source_error
            exit 1
        fi
        
        clean_installation
        exit 0
    fi
    
    # Default installation with enhanced safety checks
    if [[ "${is_source_branch}" == "true" ]]; then
        # Show safety error for source repository protection
        show_source_repository_safety_error "./install-copilot-setup.sh"
        exit 1
    fi
    
    # Validate repository if using REPO_DIRECTORY
    if [[ -n "${REPO_DIRECTORY}" ]]; then
        validate_repository "${REPO_DIRECTORY}"
    fi
    
    # Call the install_infrastructure function from fileoperations module
    install_infrastructure "${workspace_root}"
}

# ============================================================================
# SCRIPT EXECUTION
# ============================================================================

# Run main function with all arguments
main "$@"
