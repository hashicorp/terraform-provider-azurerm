#!/bin/bash

#==============================================================================
# Terraform AzureRM Provider AI Setup - Modular Bash Installation Script
#==============================================================================
#
# DESCRIPTION:
#   Modular version of the Terraform AzureRM Provider AI Setup installer.
#   Uses a clean modular architecture with separate modules for different
#   concerns: core functions, backup management, installation detection,
#   AI installation, and cleanup operations.
#
# USAGE:
#   ./install-modular.sh [OPTIONS]
#
# OPTIONS:
#   -repository-path <path>   Path to terraform-provider-azurerm repository
#   -clean                    Remove all AI setup and restore backups
#   -auto-approve             Skip interactive approval prompts
#   -help                     Show detailed help information
#
# EXAMPLES:
#   ./install-modular.sh
#   ./install-modular.sh -repository-path "/home/user/terraform-provider-azurerm"
#   ./install-modular.sh -auto-approve
#   ./install-modular.sh -clean
#
# MODULAR ARCHITECTURE:
#   - Clean separation of concerns
#   - Maintainable and testable modules
#   - Enhanced error handling and logging
#   - Improved backup and restore functionality
#
#==============================================================================

# Set strict error handling
set -euo pipefail

# Script parameters
REPOSITORY_PATH=""
CLEAN_MODE=""
AUTO_APPROVE=""
HELP_MODE=""

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MODULES_DIR="${SCRIPT_DIR}/modules/bash"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -repository-path)
            REPOSITORY_PATH="$2"
            shift 2
            ;;
        -clean)
            CLEAN_MODE="true"
            shift
            ;;
        -auto-approve)
            AUTO_APPROVE="true"
            shift
            ;;
        -help)
            HELP_MODE="true"
            shift
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use -help for usage information"
            exit 1
            ;;
    esac
done

# Source modules
modules=(
    "core-functions.sh"
    "backup-management.sh"
    "installation-detection.sh"
    "ai-installation.sh"
    "cleanup.sh"
)

for module in "${modules[@]}"; do
    module_path="${MODULES_DIR}/${module}"
    if [[ -f "$module_path" ]]; then
        source "$module_path"
    else
        echo "ERROR: Required module not found: $module_path"
        exit 1
    fi
done

# Show help if requested
if [[ "$HELP_MODE" == "true" ]]; then
    show_help
    exit 0
fi

# Main execution
main() {
    echo "Terraform AzureRM Provider AI Setup (Modular)"
    echo "============================================="
    echo

    # Check prerequisites
    check_prerequisites

    # Detect repository path
    if [[ -z "$REPOSITORY_PATH" ]]; then
        REPOSITORY_PATH=$(discover_repository)
        if [[ -z "$REPOSITORY_PATH" ]]; then
            echo "ERROR: Could not find terraform-provider-azurerm repository"
            echo "Please specify -repository-path or run from within the repository"
            exit 1
        fi
        log_info "Auto-discovered repository path: $REPOSITORY_PATH"
    else
        log_info "Using specified repository path: $REPOSITORY_PATH"
    fi

    # Validate repository
    if ! validate_repository "$REPOSITORY_PATH"; then
        echo "ERROR: Invalid repository path: $REPOSITORY_PATH"
        exit 1
    fi

    # Execute based on mode
    if [[ "$CLEAN_MODE" == "true" ]]; then
        log_info "Running in cleanup mode..."
        cleanup_installation "$REPOSITORY_PATH"
    else
        log_info "Running in installation mode..."
        
        # Check for existing installation
        local install_status
        install_status=$(detect_installation_status "$REPOSITORY_PATH")
        
        case "$install_status" in
            "clean")
                log_info "No previous installation detected, proceeding with fresh install..."
                ;;
            "partial")
                log_info "Partial installation detected, completing installation..."
                ;;
            "corrupted")
                log_warning "Previous installation found but has issues:"
                log_warning "  - Copilot instructions file is corrupted"
                log_info "Will attempt to repair installation..."
                ;;
            "complete")
                log_warning "Previous installation found but has issues:"
                log_warning "  - Copilot instructions file is corrupted"
                log_info "Will attempt to repair installation..."
                ;;
        esac
        
        # Perform installation
        install_ai_agent "$REPOSITORY_PATH"
    fi
}

# Run main function
main "$@"
