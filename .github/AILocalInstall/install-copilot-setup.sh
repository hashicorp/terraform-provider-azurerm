#!/bin/bash
# Clean Architecture Main Installer for Terraform AzureRM Provider AI Setup (Bash)
# Cross-platform support for macOS, Linux, and other Unix-like systems

set -euo pipefail

# Help function
show_usage() {
    cat << EOF
SYNOPSIS
    Clean Architecture Main Installer for Terraform AzureRM Provider AI Setup

DESCRIPTION
    This script uses the new clean modular architecture with real implementations
    from the working system. Features progress reporting and proper error handling.

USAGE
    ./install-copilot-setup.sh [OPTIONS]

OPTIONS
    -r, --repository-path <path>  Path to terraform-provider-azurerm repository
                                  If not provided, auto-discovery will be attempted
    -c, --clean                   Remove all installed AI setup files and restore backups
    -a, --auto-approve            Skip interactive approval prompts
    -v, --verify                  Run verification only without installing
    -h, --help                    Display this help message

EXAMPLES
    ./install-copilot-setup.sh                      # Auto-discover repository and install
    ./install-copilot-setup.sh -auto-approve         # Non-interactive installation
    ./install-copilot-setup.sh -verify              # Verify existing installation
    ./install-copilot-setup.sh -clean               # Remove installation and restore backups
    ./install-copilot-setup.sh -repository-path /path/to/repo  # Use specific repository path

PLATFORM SUPPORT
    - macOS: ~/Library/Application Support/Code/User
    - Linux: ~/.config/Code/User  
    - Windows (Git Bash/MSYS2): \$APPDATA/Code/User

EOF
}

# Parse command line arguments
REPOSITORY_PATH=""
CLEAN_MODE=false
AUTO_APPROVE=false
VERIFY_MODE=false
SHOW_HELP=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -repository-path)
            REPOSITORY_PATH="$2"
            shift 2
            ;;
        -clean)
            CLEAN_MODE=true
            shift
            ;;
        -auto-approve)
            AUTO_APPROVE=true
            shift
            ;;
        -verify)
            VERIFY_MODE=true
            shift
            ;;
        -help)
            SHOW_HELP=true
            shift
            ;;
        *)
            echo "Unknown option: $1" >&2
            echo "Use --help for usage information." >&2
            exit 1
            ;;
    esac
done

# Get the script directory and import all clean modules
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MODULES_DIR="$SCRIPT_DIR/modules"
BASH_MODULES_DIR="$MODULES_DIR/bash"

# Import clean architecture modules - completely self-contained
if [[ ! -f "$BASH_MODULES_DIR/core.sh" ]] || 
   [[ ! -f "$BASH_MODULES_DIR/validation.sh" ]] || 
   [[ ! -f "$BASH_MODULES_DIR/installation.sh" ]] || 
   [[ ! -f "$BASH_MODULES_DIR/ui.sh" ]]; then
    echo "ERROR: Required modules not found in $BASH_MODULES_DIR" >&2
    exit 1
fi

# Source modules in dependency order
# Note: installation.sh sources the other modules internally, so we don't need to source them all here
source "$BASH_MODULES_DIR/installation.sh"

# Initialize global variables
init_globals

# Show help if requested
if [[ "$SHOW_HELP" == "true" ]]; then
    show_help
    exit 0
fi

# Handle clean mode
handle_clean_mode() {
    write_status_message "Starting cleanup mode..." "Info"
    
    # Auto-detect repository if not provided
    if [[ -z "$REPOSITORY_PATH" ]]; then
        REPOSITORY_PATH=$(find_repository_root)
        if [[ -z "$REPOSITORY_PATH" ]]; then
            write_status_message "Could not auto-detect repository path. Please specify -r parameter." "Error"
            exit 1
        fi
    fi
    
    # Validate repository
    local validation_result
    validation_result=$(test_repository_structure "$REPOSITORY_PATH")
    
    if [[ "$(echo "$validation_result" | grep "IS_VALID_REPOSITORY" | cut -d'=' -f2)" != "true" ]]; then
        write_status_message "Repository validation failed. Please ensure you're in the correct Terraform AzureRM Provider repository." "Error"
        exit 1
    fi
    
    # Request confirmation unless auto-approve is used
    if [[ "$AUTO_APPROVE" != "true" ]]; then
        echo ""
        echo -e "${YELLOW}This will remove all AI setup files and restore backups.${NC}"
        echo -e "${YELLOW}Are you sure you want to continue? (y/N): ${NC}"
        read -r response
        if [[ "$response" != "y" && "$response" != "Y" ]]; then
            write_status_message "Cleanup cancelled." "Info"
            exit 0
        fi
    fi
    
    # Perform cleanup
    write_status_message "Removing AI setup files..." "Info"
    local cleanup_results
    cleanup_results=$(remove_ai_installation "$REPOSITORY_PATH")
    
    local success
    success=$(echo "$cleanup_results" | grep "SUCCESS" | cut -d'=' -f2)
    
    if [[ "$success" == "true" ]]; then
        write_status_message "Cleanup completed successfully!" "Success"
    else
        write_status_message "Cleanup completed with some issues." "Warning"
        # Parse and display errors
        echo "$cleanup_results" | grep "ERROR_" | while IFS='=' read -r key value; do
            write_status_message "  - $value" "Error"
        done
    fi
    
    exit 0
}

# Main function
main() {
    # Handle clean mode first
    if [[ "$CLEAN_MODE" == "true" ]]; then
        handle_clean_mode
        return
    fi
    
    # Show welcome banner
    show_welcome_banner "2.0.0"
    
    # Auto-detect repository if not provided
    if [[ -z "$REPOSITORY_PATH" ]]; then
        write_status_message "Auto-detecting repository path..." "Info"
        REPOSITORY_PATH=$(find_repository_root)
        if [[ -z "$REPOSITORY_PATH" ]]; then
            write_status_message "Could not auto-detect repository path. Please specify -r parameter." "Error"
            exit 1
        fi
    fi
    
    write_installation_progress "Validating Repository" 1 5
    
    # Validate repository
    local validation_result
    validation_result=$(test_repository_structure "$REPOSITORY_PATH")
    show_repository_info "$REPOSITORY_PATH" "$validation_result"
    
    if [[ "$(echo "$validation_result" | grep "IS_VALID_REPOSITORY" | cut -d'=' -f2)" != "true" ]]; then
        write_status_message "Repository validation failed. Please ensure you're in the correct Terraform AzureRM Provider repository." "Error"
        exit 1
    fi
    
    # If verify mode, run verification and exit
    if [[ "$VERIFY_MODE" == "true" ]]; then
        write_installation_progress "Running Verification" 5 5
        
        write_status_message "Running installation verification..." "Info"
        local verify_result
        verify_result=$(test_ai_installation "$REPOSITORY_PATH")
        
        show_installation_summary "$verify_result"
        
        local success
        success=$(echo "$verify_result" | grep "SUCCESS" | cut -d'=' -f2)
        
        if [[ "$success" == "true" ]]; then
            write_status_message "Verification completed successfully!" "Success"
            exit 0
        else
            # Check if this is no installation or partial installation
            local instruction_count prompt_count main_count settings_configured
            instruction_count=$(echo "$verify_result" | grep "INSTRUCTION_FILES_COUNT" | cut -d'=' -f2 | tr -d '\r\n' | tr -d ' ')
            prompt_count=$(echo "$verify_result" | grep "PROMPT_FILES_COUNT" | cut -d'=' -f2 | tr -d '\r\n' | tr -d ' ')
            main_count=$(echo "$verify_result" | grep "MAIN_FILES_COUNT" | cut -d'=' -f2 | tr -d '\r\n' | tr -d ' ')
            settings_configured=$(echo "$verify_result" | grep "SETTINGS_CONFIGURED" | cut -d'=' -f2 | tr -d '\r\n' | tr -d ' ')
            
            # Ensure numbers are valid (default to 0 if empty)
            instruction_count=${instruction_count:-0}
            prompt_count=${prompt_count:-0}
            main_count=${main_count:-0}
            
            local total_found=$((instruction_count + prompt_count + main_count))
            
            if [[ "$total_found" -eq 0 && "$settings_configured" != "true" ]]; then
                write_status_message "No AI installation detected locally" "Info"
            else
                write_status_message "Verification found issues with existing installation" "Warning"
            fi
            exit 1
        fi
    fi
    
    write_installation_progress "Checking Prerequisites" 2 5
    
    # Check prerequisites
    if ! test_prerequisites; then
        write_status_message "Prerequisites check failed" "Error"
        exit 1
    fi
    
    # Request confirmation
    local user_approved=false
    if [[ "$AUTO_APPROVE" == "true" ]]; then
        user_approved=true
    else
        if request_user_confirmation "Ready to install AI setup. Continue?"; then
            user_approved=true
        fi
    fi
    
    if [[ "$user_approved" != "true" ]]; then
        write_status_message "Installation cancelled by user" "Info"
        exit 0
    fi
    
    write_installation_progress "Installing Components" 3 5
    
    # Run installation
    write_status_message "Starting installation with clean architecture..." "Info"
    local results
    results=$(start_complete_installation "$REPOSITORY_PATH" "$AUTO_APPROVE")
    
    write_installation_progress "Finalizing Installation" 4 5
    
    # Show results
    show_installation_summary "$results"
    
    write_installation_progress "Complete!" 5 5
    
    local success
    success=$(echo "$results" | grep "SUCCESS" | cut -d'=' -f2)
    
    if [[ "$success" == "true" ]]; then
        show_next_steps
        write_status_message "Clean architecture installation completed successfully!" "Success"
        exit 0
    else
        write_status_message "Installation completed with issues. Check the summary above." "Warning"
        exit 1
    fi
}

# Error handling
trap 'write_status_message "Installation failed: An error occurred on line $LINENO" "Error"; echo ""; echo "For troubleshooting help, check the logs above or run with --verify to diagnose issues."; exit 1' ERR

# Run main function
main "$@"
