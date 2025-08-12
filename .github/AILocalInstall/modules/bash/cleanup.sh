#!/bin/bash

#==============================================================================
# Cleanup Module
#==============================================================================

# Global cleanup flag
CLEANUP_MODE=false

# Main cleanup function
cleanup_ai_setup() {
    local repo_path="$1"
    
    log_info "Starting AI setup cleanup..."
    
    if [[ ! -d "$repo_path" ]]; then
        log_error "Repository path not found: $repo_path"
        return 1
    fi
    
    # Remove copied files from VS Code user directory
    remove_vscode_files
    
    # Clean VS Code settings (includes backup restoration)
    clean_vscode_settings
    
    # Clean up temporary files
    cleanup_temp_files
    
    log_success "AI setup cleanup completed!"
}

# Note: We don't remove the .github/AILocalInstall directory as it's part of the repository
# Only the VS Code configuration is considered "installed"

# Remove copied files from VS Code user directory
remove_vscode_files() {
    # Define VS Code user directory paths
    local vscode_user_dir
    if [[ "$OSTYPE" == "darwin"* ]]; then
        vscode_user_dir="$HOME/Library/Application Support/Code/User"
    elif [[ -n "${WSL_DISTRO_NAME:-}" ]] || grep -q microsoft /proc/version 2>/dev/null; then
        # WSL detected - first check for native Linux VS Code installation
        if [[ -d "$HOME/.config/Code/User" ]]; then
            # Native Linux VS Code installation found in WSL
            vscode_user_dir="$HOME/.config/Code/User"
        else
            # Fallback to Windows VS Code path
            local windows_user
            windows_user=$(powershell.exe -c "echo \$env:USERNAME" 2>/dev/null | tr -d '\r\n' || echo "")
            if [[ -n "$windows_user" ]]; then
                # Convert Windows path to WSL path
                vscode_user_dir="/mnt/c/Users/${windows_user}/AppData/Roaming/Code/User"
            else
                # Final fallback to Linux path
                vscode_user_dir="$HOME/.config/Code/User"
            fi
        fi
    else
        vscode_user_dir="$HOME/.config/Code/User"
    fi
    
    local target_instructions_dir="$vscode_user_dir/instructions/terraform-azurerm"
    local target_prompts_dir="$vscode_user_dir/prompts"
    local target_copilot_file="$vscode_user_dir/copilot-instructions.md"
    
    log_info "Removing copied files from VS Code user directory..."
    
    # Remove instruction files
    if [[ -d "$target_instructions_dir" ]]; then
        rm -rf "$target_instructions_dir"
        log_success "Removed instruction files from VS Code user directory"
    else
        log_info "No instruction files found in VS Code user directory"
    fi
    
    # Remove terraform-azurerm prompt files from root prompts directory
    if [[ -d "$target_prompts_dir" ]]; then
        local removed_prompts=0
        local terraform_prompts=(
            "add-unit-tests.prompt.md"
            "code-review-committed-changes.prompt.md"
            "code-review-local-changes.prompt.md"
            "setup-go-dev-environment.prompt.md"
            "summarize-repo-deep-dive.prompt.md"
            "summarize-repo.prompt.md"
        )
        
        for prompt_file in "${terraform_prompts[@]}"; do
            local prompt_path="$target_prompts_dir/$prompt_file"
            if [[ -f "$prompt_path" ]]; then
                rm -f "$prompt_path"
                removed_prompts=$((removed_prompts + 1))
            fi
        done
        
        if [[ $removed_prompts -gt 0 ]]; then
            log_success "Removed $removed_prompts terraform-azurerm prompt files from VS Code user directory"
        else
            log_info "No terraform-azurerm prompt files found in VS Code user directory"
        fi
    else
        log_info "No prompts directory found in VS Code user directory"
    fi
    
    # Remove main copilot instructions file
    if [[ -f "$target_copilot_file" ]]; then
        rm -f "$target_copilot_file"
        log_success "Removed copilot instructions from VS Code user directory"
    else
        log_info "No copilot instructions file found in VS Code user directory"
    fi
    
    # Clean up empty parent directories if they exist
    local instructions_parent="$vscode_user_dir/instructions"
    
    if [[ -d "$instructions_parent" ]] && [[ -z "$(find "$instructions_parent" -mindepth 1)" ]]; then
        rmdir "$instructions_parent" 2>/dev/null || true
        log_info "Removed empty instructions parent directory"
    fi
    
    # Note: We don't remove the prompts directory since we only removed individual files
    # and other prompt files from other projects may exist
}

# Clean VS Code settings with backup restoration
clean_vscode_settings() {
    # Define VS Code user directory path
    local vscode_user_dir
    if [[ "$OSTYPE" == "darwin"* ]]; then
        vscode_user_dir="$HOME/Library/Application Support/Code/User"
    elif [[ -n "${WSL_DISTRO_NAME:-}" ]] || grep -q microsoft /proc/version 2>/dev/null; then
        # WSL detected - first check for native Linux VS Code installation
        if [[ -d "$HOME/.config/Code/User" ]]; then
            # Native Linux VS Code installation found in WSL
            vscode_user_dir="$HOME/.config/Code/User"
        else
            # Fallback to Windows VS Code path
            local windows_user
            windows_user=$(powershell.exe -c "echo \$env:USERNAME" 2>/dev/null | tr -d '\r\n' || echo "")
            if [[ -n "$windows_user" ]]; then
                # Convert Windows path to WSL path
                vscode_user_dir="/mnt/c/Users/${windows_user}/AppData/Roaming/Code/User"
            else
                # Final fallback to Linux path
                vscode_user_dir="$HOME/.config/Code/User"
            fi
        fi
    else
        vscode_user_dir="$HOME/.config/Code/User"
    fi
    
    local settings_path="$vscode_user_dir/settings.json"
    local backup_dir="$vscode_user_dir/azurerm_backups"
    
    if [[ ! -f "$settings_path" ]]; then
        log_info "No VS Code settings file found to clean"
        return 0
    fi
    
    # Try to restore from backup if available
    if [[ -d "$backup_dir" ]]; then
        local backup_file
        backup_file=$(find "$backup_dir" -name "settings_backup_*.json" -type f -print0 | \
                     xargs -0 ls -t 2>/dev/null | head -n 1)
        
        if [[ -n "$backup_file" && -f "$backup_file" ]]; then
            log_info "Restoring VS Code settings from backup..."
            if cp "$backup_file" "$settings_path"; then
                log_success "VS Code settings restored from backup"
                # Clean up backup directory
                rm -rf "$backup_dir"
                log_info "Cleaned up backup directory"
                return 0
            else
                log_error "Failed to restore from backup, attempting manual cleanup"
            fi
        fi
    fi
    
    # Fallback: manually remove Terraform settings
    log_info "Cleaning Terraform settings from VS Code manually..."
    
    # Use jq to remove Terraform-specific settings if available
    if command -v jq >/dev/null 2>&1; then
        local temp_file
        temp_file=$(mktemp)
        
        if jq 'del(
            ."terraform_azurerm_installation_date",
            ."terraform_azurerm_backup_length",
            ."github.copilot.chat.commitMessageGeneration.instructions",
            ."github.copilot.chat.summarizeAgentConversationHistory.enabled",
            ."github.copilot.chat.reviewSelection.enabled",
            ."github.copilot.chat.reviewSelection.instructions",
            ."github.copilot.advanced",
            ."github.copilot.enable",
            ."files.associations"
        )' "$settings_path" > "$temp_file" 2>/dev/null; then
            
            if mv "$temp_file" "$settings_path"; then
                log_success "Cleaned Terraform settings from VS Code"
                return 0
            else
                log_error "Failed to update cleaned settings file"
                rm -f "$temp_file"
            fi
        else
            log_error "Failed to parse settings.json with jq"
            rm -f "$temp_file"
        fi
    fi
    
    log_warning "Could not automatically clean VS Code settings"
    log_info "You may need to manually remove Terraform-related settings from: $settings_path"
    return 1
}

# Clean up temporary files
cleanup_temp_files() {
    log_info "Cleaning up temporary files..."
    
    local temp_dirs=(
        "/tmp/ai_local_install_$$"
        "/tmp/terraform_ai_setup_$$"
    )
    
    for temp_dir in "${temp_dirs[@]}"; do
        if [[ -d "$temp_dir" ]]; then
            rm -rf "$temp_dir"
            log_info "Removed temporary directory: $temp_dir"
        fi
    done
}

# Confirm cleanup operation
confirm_cleanup() {
    echo ""
    echo "⚠️  CLEANUP CONFIRMATION"
    echo "======================================"
    echo ""
    echo "This will:"
    echo "  • Remove AI instruction files from VS Code user directory"
    echo "  • Remove Terraform-specific prompt files from VS Code"
    echo "  • Restore original VS Code settings from backup"
    echo "  • Clean up temporary files"
    echo ""
    echo "Note: Repository files (.github/AILocalInstall) will NOT be removed"
    echo ""
    
    if [[ "$AUTO_APPROVE" != "true" ]]; then
        read -p "Do you want to continue? (y/N): " -r
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "Cleanup cancelled by user"
            return 1
        fi
    fi
    
    return 0
}

# Main cleanup entry point - called by the main script
cleanup_installation() {
    local repo_path="$1"
    perform_cleanup "$repo_path"
}

# Main cleanup entry point
perform_cleanup() {
    local repo_path="$1"
    
    log_info "=== CLEANUP MODE ==="
    
    if ! confirm_cleanup; then
        return 1
    fi
    
    if cleanup_ai_setup "$repo_path"; then
        echo "✅ Cleanup completed successfully!"
        return 0
    else
        log_error "Cleanup failed"
        return 1
    fi
}
