#!/bin/bash

#==============================================================================
# User Interface Module
#==============================================================================
#
# This module handles user interactions, including prompts, confirmations,
# and status displays.
#

# Show welcome message
show_welcome() {
    local script_name="$1"
    
    cat << 'WELCOME_EOF'

🤖 Terraform AzureRM Provider AI Setup
======================================

This script installs an AI agent configuration that enhances GitHub Copilot
with deep knowledge of the Terraform AzureRM Provider codebase.

Features:
- Context-aware code generation and review
- Azure-specific implementation patterns  
- Testing guidelines and best practices
- Documentation standards enforcement
- Error handling and debugging assistance

The installation will:
1. Check your VS Code installation
2. Create a backup of your current settings
3. Install the AI agent configuration files
4. Update your VS Code settings

WELCOME_EOF
}

# Show help information
show_help() {
    local script_name="$1"
    
    cat << HELP_EOF
Usage: $script_name [OPTIONS]

Install AI agent configuration for the Terraform AzureRM Provider.

OPTIONS:
    -h, --help              Show this help message
    -r, --repository PATH   Specify repository path (default: current directory)
    -f, --force            Force installation without prompts
    -s, --status           Show installation status and exit
    -u, --uninstall        Uninstall the AI agent configuration
    -b, --backup           Create backup only (no installation)
    -R, --restore          Restore from most recent backup
    --list-backups         List available backups
    --verify               Verify current installation
    --repair               Repair corrupted installation
    -v, --verbose          Enable verbose output
    -q, --quiet            Suppress non-error output

EXAMPLES:
    $script_name                           # Install in current directory
    $script_name -r /path/to/repo          # Install in specified directory
    $script_name --status                  # Check installation status
    $script_name --backup                  # Create backup only
    $script_name --restore                 # Restore from backup
    $script_name --verify                  # Verify installation
    $script_name --repair                  # Repair installation

For more information, visit:
https://github.com/hashicorp/terraform-provider-azurerm

HELP_EOF
}

# Prompt for user confirmation
prompt_confirmation() {
    local message="$1"
    local default="${2:-n}"
    
    local prompt
    if [[ "$default" == "y" ]]; then
        prompt="$message [Y/n]: "
    else
        prompt="$message [y/N]: "
    fi
    
    while true; do
        printf "%s" "$prompt"
        read -r response
        
        # Use default if empty response
        if [[ -z "$response" ]]; then
            response="$default"
        fi
        
        case "$response" in
            [Yy]|[Yy][Ee][Ss])
                return 0
                ;;
            [Nn]|[Nn][Oo])
                return 1
                ;;
            *)
                echo "Please answer yes or no."
                ;;
        esac
    done
}

# Show installation status in human-readable format
show_status() {
    local repo_path="$1"
    
    echo "🔍 Checking installation status..."
    echo
    
    get_installation_summary "$repo_path"
    
    echo
    
    local status
    status=$(detect_installation_status "$repo_path")
    
    case "$status" in
        "clean")
            echo "✅ Status: Ready for fresh installation"
            ;;
        "partial")
            echo "⚠️  Status: Partial installation detected"
            echo "   Some components may be missing or outdated"
            ;;
        "complete")
            echo "✅ Status: Installation appears complete"
            ;;
        "corrupted")
            echo "❌ Status: Installation is corrupted"
            echo "   Files may be damaged or incomplete"
            ;;
        *)
            echo "❓ Status: Unknown ($status)"
            ;;
    esac
    
    echo
    
    # Show next steps based on status
    case "$status" in
        "clean")
            echo "💡 Next steps:"
            echo "   Run the installer to set up the AI agent configuration"
            ;;
        "partial")
            echo "💡 Next steps:"
            echo "   Run the installer to complete the setup"
            echo "   Or use --repair to fix the current installation"
            ;;
        "complete")
            echo "💡 Next steps:"
            echo "   Open VS Code and start using GitHub Copilot!"
            echo "   Use --verify to confirm everything is working"
            ;;
        "corrupted")
            echo "💡 Next steps:"
            echo "   Use --repair to fix the installation"
            echo "   Or run a fresh installation"
            ;;
    esac
}

# Show backup list
show_backup_list() {
    local backup_dir
    backup_dir="$(get_backup_directory)"
    
    echo "📋 Available backups:"
    echo
    
    if [[ ! -d "$backup_dir" ]]; then
        echo "   No backups found (backup directory doesn't exist)"
        return 0
    fi
    
    local backup_files
    backup_files=($(find "$backup_dir" -name "settings_backup_*.json" -type f | sort -r))
    
    if [[ ${#backup_files[@]} -eq 0 ]]; then
        echo "   No backups found"
        return 0
    fi
    
    local count=1
    for backup_file in "${backup_files[@]}"; do
        local basename_file
        basename_file=$(basename "$backup_file")
        local timestamp
        timestamp=$(echo "$basename_file" | sed 's/settings_backup_\(.*\)\.json/\1/' | tr '_' ' ')
        local size
        size=$(wc -c < "$backup_file")
        
        echo "   $count. $basename_file"
        echo "      Created: $timestamp"
        echo "      Size: $size bytes"
        echo "      Path: $backup_file"
        echo
        
        ((count++))
    done
}

# Show progress indicator
show_progress() {
    local current="$1"
    local total="$2"
    local message="$3"
    
    local percent=$((current * 100 / total))
    local filled=$((percent / 5))
    local empty=$((20 - filled))
    
    printf "\r%s [" "$message"
    
    # Show filled portion
    for ((i=0; i<filled; i++)); do
        printf "="
    done
    
    # Show empty portion
    for ((i=0; i<empty; i++)); do
        printf " "
    done
    
    printf "] %d%%" "$percent"
    
    if [[ $current -eq $total ]]; then
        echo " ✅"
    fi
}

# Show error with context
show_error_with_context() {
    local error_msg="$1"
    local context="$2"
    
    echo "❌ Error: $error_msg"
    
    if [[ -n "$context" ]]; then
        echo "   Context: $context"
    fi
    
    echo "   Check the log above for more details"
    echo "   If the problem persists, try running with --verbose for more information"
}

# Show success message with next steps
show_success_message() {
    local operation="$1"
    local repo_path="$2"
    
    case "$operation" in
        "install")
            echo "🎉 Installation completed successfully!"
            echo
            show_installation_instructions
            ;;
        "uninstall")
            echo "🗑️  Uninstallation completed successfully!"
            echo "   Your original VS Code settings have been restored"
            ;;
        "backup")
            echo "💾 Backup created successfully!"
            local backup_dir
            backup_dir="$(get_backup_directory)"
            echo "   Backup location: $backup_dir"
            ;;
        "restore")
            echo "🔄 Restore completed successfully!"
            echo "   Your VS Code settings have been restored from backup"
            ;;
        "verify")
            echo "✅ Verification completed successfully!"
            echo "   Your AI agent configuration is working correctly"
            ;;
        "repair")
            echo "🔧 Repair completed successfully!"
            echo "   Your AI agent configuration has been fixed"
            ;;
        *)
            echo "✅ Operation '$operation' completed successfully!"
            ;;
    esac
}

# Interactive mode for installation
interactive_install() {
    local repo_path="$1"
    
    show_welcome "$(basename "$0")"
    echo
    
    # Show current status
    show_status "$repo_path"
    echo
    
    # Check if already installed
    local status
    status=$(detect_installation_status "$repo_path")
    
    if [[ "$status" == "complete" ]]; then
        if prompt_confirmation "AI agent appears to be already installed. Continue anyway?"; then
            echo "Proceeding with installation..."
        else
            log_info "Installation cancelled by user"
            return 1
        fi
    fi
    
    # Confirm installation
    if prompt_confirmation "Proceed with AI agent installation?"; then
        echo
        install_ai_agent "$repo_path"
        local exit_code=$?
        
        if [[ $exit_code -eq 0 ]]; then
            echo
            show_success_message "install" "$repo_path"
        else
            echo
            show_error_with_context "Installation failed" "See log messages above"
        fi
        
        return $exit_code
    else
        log_info "Installation cancelled by user"
        return 1
    fi
}

# Show version information
show_version() {
    local script_dir
    script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    local version_file="$script_dir/../VERSION"
    
    echo "Terraform AzureRM Provider AI Setup"
    
    if [[ -f "$version_file" ]]; then
        local version
        version=$(cat "$version_file")
        echo "Version: $version"
    else
        echo "Version: unknown"
    fi
    
    echo "GitHub: https://github.com/hashicorp/terraform-provider-azurerm"
}

# Show operation summary
show_operation_summary() {
    local operation="$1"
    local repo_path="$2"
    local start_time="$3"
    local end_time="$4"
    
    local duration=$((end_time - start_time))
    
    echo
    echo "📊 Operation Summary:"
    echo "   Operation: $operation"
    echo "   Repository: $repo_path"
    echo "   Duration: ${duration}s"
    echo "   Timestamp: $(date)"
    echo
}
