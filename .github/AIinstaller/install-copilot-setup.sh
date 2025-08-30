#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Main AI Infrastructure Installer for Terraform AzureRM Provider (macOS/Linux)
# Version: 1.0.0
# Description: Interactive installer for AI-powered development tools
# Requires bash 4.0+

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
# WORKSPACE DETECTION - Simple and reliable
# ============================================================================

get_workspace_root() {
    local repo_directory="$1"
    local script_directory="$2"

    # If repo_directory is provided, use it (validation happens later)
    if [[ -n "${repo_directory}" ]]; then
        echo "${repo_directory}"
        return
    fi

    # Otherwise, find workspace root from script location
    local current_path="${script_directory}"
    while [[ -n "${current_path}" && "${current_path}" != "$(dirname "${current_path}")" ]]; do
        if [[ -f "${current_path}/go.mod" ]]; then
            echo "${current_path}"
            return
        fi
        current_path="$(dirname "${current_path}")"
    done

    # If no workspace found, return the directory where the script was called from
    # This allows help and other functions to work, with validation happening separately
    pwd
}

# ============================================================================
# MAIN EXECUTION - Clean and simple
# ============================================================================

main() {
    #
    # Main entry point for the installer - matches PowerShell structure
    #

    # STEP 1: Parse command line arguments
    parse_arguments "$@"

    # STEP 2: Early safety check - fail fast if on source branch with repo directory
    if [[ -n "${REPO_DIRECTORY}" ]]; then
        # Get current branch of the target repository quickly
        local current_branch
        if [[ -d "${REPO_DIRECTORY}/.git" ]]; then
            current_branch=$(cd "${REPO_DIRECTORY}" && git branch --show-current 2>/dev/null || echo "unknown")
        else
            current_branch="unknown"
        fi

        # Block operations on source branch immediately (except verify, help, bootstrap)
        # Source branches: main, master, exp/terraform_copilot
        local source_branches=("main" "master" "exp/terraform_copilot")
        local is_source_branch=false
        for branch in "${source_branches[@]}"; do
            if [[ "${current_branch}" == "${branch}" ]]; then
                is_source_branch=true
                break
            fi
        done

        if [[ "${is_source_branch}" == "true" ]] && [[ "${VERIFY}" != "true" ]] && [[ "${HELP}" != "true" ]] && [[ "${BOOTSTRAP}" != "true" ]]; then
            show_safety_violation "${current_branch}" "Install" "true"
            exit 1
        fi
    fi

    # STEP 3: Initialize workspace and validate it's a proper terraform-provider-azurerm repo
    local workspace_root
    workspace_root="$(get_workspace_root "${REPO_DIRECTORY}" "${SCRIPT_DIR}")"

    # STEP 4: Early workspace validation before doing anything else
    local workspace_valid workspace_reason
    if validate_repository "${workspace_root}"; then
        workspace_valid=true
        workspace_reason=""
    else
        workspace_valid=false
        workspace_reason="Missing required files"
    fi

    # STEP 5: Get branch information for consistent display
    local current_branch branch_type
    if [[ -n "${REPO_DIRECTORY}" ]]; then
        if [[ -d "${workspace_root}/.git" ]]; then
            current_branch=$(cd "${workspace_root}" && git branch --show-current 2>/dev/null || echo "unknown")
        else
            current_branch="unknown"
        fi
    else
        current_branch=$(git branch --show-current 2>/dev/null || echo "unknown")
    fi

    # Determine branch type - be explicit about what we know vs don't know
    case "${current_branch}" in
        "main"|"master"|"exp/terraform_copilot")
            branch_type="source"
            ;;
        "unknown"|"")
            branch_type="unknown"
            ;;
        *)
            # Any other valid branch name is a feature branch
            branch_type="feature"
            ;;
    esac

    # STEP 6: CONSISTENT PATTERN - Every operation gets the same header and branch detection
    write_header "Terraform AzureRM Provider - AI Infrastructure Installer" "${VERSION}"
    show_branch_detection "${current_branch}" "${workspace_root}"

    # STEP 7: Simple parameter handling (like PowerShell)
    if [[ "${HELP}" == "true" ]]; then
        show_usage "${branch_type}" "${workspace_valid}" "${workspace_reason}"
        exit 0
    fi

    # STEP 8: For all other operations, workspace must be valid
    if [[ "${workspace_valid}" != "true" ]]; then
        echo ""
        write_error_message "WORKSPACE VALIDATION FAILED: ${workspace_reason}"
        echo ""

        # Context-aware error message based on how the script was invoked
        if [[ -n "${REPO_DIRECTORY}" ]]; then
            write_plain " Please ensure the -repo-directory argument is pointing to a valid GitHub terraform-provider-azurerm repository."
        else
            write_plain " Please ensure you are running this script from within a terraform-provider-azurerm repository."
        fi
        echo ""
        print_separator

        # Show help menu for guidance
        show_usage "${branch_type}" "false" "${workspace_reason}"
        exit 1
    fi

    # STEP 9: Execute single operation based on parameters (like PowerShell)
    if [[ "${VERIFY}" == "true" ]]; then
        verify_installation "${workspace_root}"
        exit 0
    fi

    if [[ "${BOOTSTRAP}" == "true" ]]; then
        # Show operation title (main header already displayed)
        write_section "Bootstrap - Copying Installer to User Profile"

        # Execute the bootstrap operation
        if bootstrap_files_to_profile "${SCRIPT_DIR}" "$(get_user_profile)" "${SCRIPT_DIR}/file-manifest.config"; then
            # Show detailed summary with next steps
            local user_profile
            user_profile=$(get_user_profile)
            local size_kb=$((BOOTSTRAP_STATS_TOTAL_SIZE / 1024))
            show_operation_summary "Bootstrap" "true" "false" \
                "Files Copied:${BOOTSTRAP_STATS_FILES_COPIED}" \
                "Total Size:${size_kb} KB" \
                "Location:${user_profile}" \
                --next-steps \
                "1. Switch to your feature branch:" \
                "   git checkout feature/your-branch-name" \
                "" \
                "2. Run the installer from your user profile:" \
                "   cd ~/.terraform-ai-installer" \
                "   ./install-copilot-setup.sh -repo-directory \"<path-to-your-terraform-provider-azurerm>\""
        else
            write_error_message "Bootstrap operation failed"
            exit 1
        fi
        exit 0
    fi

    if [[ "${CLEAN}" == "true" ]]; then
        clean_infrastructure "${workspace_root}"
        exit 0
    fi

    # STEP 10: Installation path (when -repo-directory is provided and not other specific operations)
    if [[ -n "${REPO_DIRECTORY}" ]] && [[ "${HELP}" != "true" ]] && [[ "${VERIFY}" != "true" ]] && [[ "${BOOTSTRAP}" != "true" ]] && [[ "${CLEAN}" != "true" ]]; then
        # Proceed with installation
        install_infrastructure "${workspace_root}"
        exit 0
    fi

    # STEP 11: Default - show source branch help and welcome
    show_source_repository_safety_error "./install-copilot-setup.sh"
    exit 0
}

# ============================================================================
# COMMAND LINE ARGUMENT PROCESSING
# ============================================================================

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
# SCRIPT EXECUTION
# ============================================================================

# Run main function with all arguments - single entry point like PowerShell
main "$@"
