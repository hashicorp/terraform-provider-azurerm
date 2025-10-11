#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Main AI Infrastructure Installer for Terraform AzureRM Provider (macOS/Linux)
# Version: 1.0.0
# Description: Interactive installer for AI-powered development tools
# Requires bash 3.2+ (compatible with macOS default bash)

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
DRY_RUN=false             # Show what would be done without making changes
VERIFY=false              # Check the current state of the workspace
CLEAN=false               # Remove AI infrastructure from the workspace
HELP=false                # Show detailed help information

# Export variables that need to be accessible to modules as global variables
# Note: Other variables are passed as function parameters, so they don't need to be exported
export DRY_RUN

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

    # STEP 1: Parse command line arguments first
    parse_arguments "$@"

    # STEP 2: Show header immediately for consistent user experience
    write_header "Terraform AzureRM Provider - AI Infrastructure Installer" "${VERSION}"

    # STEP 3: Get workspace root and branch information for display and safety checks
    local workspace_root
    workspace_root="$(get_workspace_root "${REPO_DIRECTORY}" "${SCRIPT_DIR}")"

    local current_branch
    if [[ -n "${REPO_DIRECTORY}" ]]; then
        if [[ -d "${workspace_root}/.git" ]]; then
            current_branch=$(cd "${workspace_root}" && git branch --show-current 2>/dev/null || echo "unknown")
        else
            current_branch="unknown"
        fi
    else
        current_branch=$(git branch --show-current 2>/dev/null || echo "unknown")
    fi

    # STEP 4: Show branch detection immediately after getting branch info
    show_branch_detection "${current_branch}" "${workspace_root}"

    # STEP 5: Early safety check - fail fast if on source branch with repo directory
    if [[ -n "${REPO_DIRECTORY}" ]]; then
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
            # Safety violation - header and branch detection already shown above
            show_safety_violation "${current_branch}" "Install" "true"
            exit 1
        fi
    fi

    # STEP 6: Initialize workspace validation (workspace_root already set above)
    local workspace_valid workspace_reason
    if validate_repository "${workspace_root}"; then
        workspace_valid=true
        workspace_reason=""
    else
        workspace_valid=false
        workspace_reason="Missing required files"
    fi

    # STEP 7: Determine branch type - be explicit about what we know vs don't know
    local branch_type
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

    # STEP 8: Detect what command was attempted (for better error messages)
    local attempted_command=""
    if [[ "${BOOTSTRAP}" == "true" ]]; then
        attempted_command="-bootstrap"
    elif [[ "${VERIFY}" == "true" ]]; then
        attempted_command="-verify"
    elif [[ "${CLEAN}" == "true" ]]; then
        attempted_command="-clean"
    elif [[ "${HELP}" == "true" ]]; then
        attempted_command="-help"
    elif [[ "${DRY_RUN}" == "true" ]]; then
        attempted_command="-dry-run"
    elif [[ -n "${REPO_DIRECTORY}" && "${HELP}" != "true" && "${VERIFY}" != "true" && "${BOOTSTRAP}" != "true" && "${CLEAN}" != "true" ]]; then
        attempted_command="-repo-directory \"${REPO_DIRECTORY}\""
    fi

    # STEP 9: Simple parameter handling (like PowerShell)
    if [[ "${HELP}" == "true" ]]; then
        show_usage "${branch_type}" "${workspace_valid}" "${workspace_reason}" "${attempted_command}"
        exit 0
    fi

    # STEP 10: For all operations, workspace must be valid (each operation handles its own specific validation)
    if [[ "${workspace_valid}" != "true" ]]; then
        show_workspace_validation_error "${workspace_reason}" "$([[ -n "${REPO_DIRECTORY}" ]] && echo "true" || echo "false")"

        # Show help menu for guidance
        show_usage "${branch_type}" "false" "${workspace_reason}" "${attempted_command}"
        exit 1
    fi

    # STEP 11: Execute single operation based on parameters (like PowerShell)
    if [[ "${VERIFY}" == "true" ]]; then
        verify_installation "${workspace_root}"
        exit 0
    fi

    if [[ "${BOOTSTRAP}" == "true" ]]; then
        # Show operation title (main header already displayed)
        write_section "Bootstrap - Copying Installer to User Profile"

        # Execute the bootstrap operation with built-in validation
        if bootstrap_files_to_profile "$(pwd)" "$(get_user_profile)" "${SCRIPT_DIR}/file-manifest.config" "${current_branch}" "${branch_type}" "${SCRIPT_DIR}"; then
            # Show detailed summary with next steps
            local user_profile
            user_profile=$(get_user_profile)
            local size_kb=$((BOOTSTRAP_STATS_TOTAL_SIZE / 1024))
            show_operation_summary "Bootstrap" "true" "false" \
                "Files Copied:${BOOTSTRAP_STATS_FILES_COPIED}" \
                "Total Size:${size_kb} KB" \
                "Location:${user_profile}" \
                --next-steps \
                " 1. Switch to your feature branch:" \
                "    git checkout feature/your-branch-name" \
                "" \
                " 2. Run the installer from your user profile:" \
                "    cd ~/.terraform-ai-installer" \
                "    ./install-copilot-setup.sh -repo-directory \"<path-to-your-terraform-provider-azurerm>\""

            # Show welcome message after successful bootstrap
            show_source_branch_welcome "${current_branch}"
        else
            exit 1
        fi
        exit 0
    fi

    if [[ "${CLEAN}" == "true" ]]; then
        clean_infrastructure "${workspace_root}" "${current_branch}" "${branch_type}"
        exit 0
    fi

    # STEP 11: Installation path (when -repo-directory is provided and not other specific operations)
    if [[ -n "${REPO_DIRECTORY}" ]] && [[ "${HELP}" != "true" ]] && [[ "${VERIFY}" != "true" ]] && [[ "${BOOTSTRAP}" != "true" ]] && [[ "${CLEAN}" != "true" ]]; then
        # Proceed with installation
        install_infrastructure "${workspace_root}" "${current_branch}" "${branch_type}"
        exit 0
    fi

    # STEP 12: Default - show source branch help and welcome
    show_source_branch_help
    show_source_branch_welcome "${current_branch}"
    exit 0
}

# ============================================================================
# COMMAND LINE ARGUMENT PROCESSING
# ============================================================================

check_typos() {
    local param="$1"
    local suggestion=""

    # Handle bare dash edge case
    if [[ "${param}" == "-" ]] || [[ "${param}" == "--" ]]; then
        printf " \033[31mError:\033[0m\033[36m Failed to parse command-line argument:\033[0m\n"
        printf " \033[36mArgument provided but not defined:\033[0m \033[33m${param}\033[0m\n"
        echo ""
        printf " \033[36mFor more help on using this command, run:\033[0m\n"
        printf "   \033[37m$0 -help\033[0m\n"
        echo ""
        exit 1
    fi

    # Remove leading dashes and convert to lowercase
    local clean_param="${param#-}"
    clean_param="${clean_param#-}"
    local lower_param="$(echo "${clean_param}" | tr '[:upper:]' '[:lower:]')"

    # Direct prefix matching (higher priority)
    if echo "${lower_param}" | grep -q '^cl'; then
        suggestion="clean"
    elif echo "${lower_param}" | grep -q '^bo'; then
        suggestion="bootstrap"
    elif echo "${lower_param}" | grep -q '^ve'; then
        suggestion="verify"
    elif echo "${lower_param}" | grep -q '^he'; then
        suggestion="help"
    elif echo "${lower_param}" | grep -q '^dr'; then
        suggestion="dry-run"
    elif echo "${lower_param}" | grep -q '^re'; then
        suggestion="repo-directory"
    # Fuzzy matching (lower priority)
    elif [[ "${lower_param}" == *cle* ]]; then
        suggestion="clean"
    elif [[ "${lower_param}" == *boo* ]]; then
        suggestion="bootstrap"
    elif [[ "${lower_param}" == *ver* ]]; then
        suggestion="verify"
    elif [[ "${lower_param}" == *hel* ]]; then
        suggestion="help"
    elif [[ "${lower_param}" == *dry* ]]; then
        suggestion="dry-run"
    elif [[ "${lower_param}" == *repo* ]]; then
        suggestion="repo-directory"
    fi

    if [[ -n "${suggestion}" ]]; then
        printf " \033[31mError:\033[0m\033[36m Failed to parse command-line argument:\033[0m\n"
        printf " \033[36mArgument provided but not defined:\033[0m \033[33m${param}\033[0m\n"
        printf " \033[36mDid you mean:\033[0m \033[32m-${suggestion}\033[0m\033[36m?\033[0m\n"
        echo ""
        printf " \033[36mFor more help on using this command, run:\033[0m\n"
        printf "   \033[37m$0 -help\033[0m\n"
        echo ""
        exit 1
    fi
}

parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -bootstrap)
                BOOTSTRAP=true
                shift
                ;;
            -repo-directory)
                if [[ $# -lt 2 ]] || [[ "${2:-}" == -* ]]; then
                    write_error_message " Option -repo-directory requires a directory path"
                    exit 1
                fi
                REPO_DIRECTORY="$2"
                shift 2
                ;;
            -dry-run)
                DRY_RUN=true
                export DRY_RUN
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
                # Check for typos before showing generic error
                if [[ "$1" == -* ]]; then
                    check_typos "$1"
                fi

                printf " \033[31mError:\033[0m\033[36m Failed to parse command-line argument:\033[0m\n"
                printf " \033[36mUnknown option:\033[0m \033[33m$1\033[0m\n"
                echo ""
                printf " \033[36mFor more help on using this command, run:\033[0m\n"
                printf "   \033[37m$0 -help\033[0m\n"
                echo ""
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
