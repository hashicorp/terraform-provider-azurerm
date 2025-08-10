#!/bin/bash

#==============================================================================
# AI Installation Module
#==============================================================================
#
# This module handles the actual installation of the AI agent configuration,
# including copying files and updating VS Code settings.
#

# Install AI agent configuration
install_ai_agent() {
    local repo_path="$1"
    local settings_path
    settings_path="$(get_vscode_settings_path)"
    
    log_info "Installing AI agent configuration..."
    
    # Check VS Code installation
    if ! check_vscode_installation; then
        return 1
    fi
    
    # Create backup
    if ! create_backup "$settings_path"; then
        log_error "Failed to create backup"
        return 1
    fi
    
    # Copy copilot instructions file
    if ! copy_copilot_instructions "$repo_path"; then
        log_error "Failed to copy copilot instructions"
        return 1
    fi
    
    # Update VS Code settings
    if ! update_vscode_settings "$settings_path" "$repo_path"; then
        log_error "Failed to update VS Code settings"
        return 1
    fi
    
    log_success "AI agent configuration installed successfully"
    
    # Show summary
    echo
    get_installation_summary "$repo_path"
    
    return 0
}

# Copy copilot instructions file
copy_copilot_instructions() {
    local repo_path="$1"
    local instructions_file="$repo_path/.github/copilot-instructions.md"
    local github_dir="$repo_path/.github"
    
    # Create .github directory if it doesn't exist
    if [[ ! -d "$github_dir" ]]; then
        if ! mkdir -p "$github_dir"; then
            log_error "Failed to create .github directory: $github_dir"
            return 1
        fi
        log_info "Created .github directory: $github_dir"
    fi
    
    # Copy from the script's directory
    local script_dir
    script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    local source_file="$script_dir/../copilot-instructions.md"
    
    if [[ ! -f "$source_file" ]]; then
        log_error "Source copilot instructions file not found: $source_file"
        return 1
    fi
    
    if cp "$source_file" "$instructions_file"; then
        log_success "Copied copilot instructions to: $instructions_file"
        return 0
    else
        log_error "Failed to copy copilot instructions file"
        return 1
    fi
}

# Update VS Code settings
update_vscode_settings() {
    local settings_path="$1"
    local repo_path="$2"
    local instructions_file="$repo_path/.github/copilot-instructions.md"
    
    # Create settings file if it doesn't exist
    if [[ ! -f "$settings_path" ]]; then
        local settings_dir
        settings_dir="$(dirname "$settings_path")"
        if ! mkdir -p "$settings_dir"; then
            log_error "Failed to create settings directory: $settings_dir"
            return 1
        fi
        echo '{}' > "$settings_path"
        log_info "Created new settings file: $settings_path"
    fi
    
    # Validate existing settings
    if ! is_valid_json "$settings_path"; then
        log_warning "Existing settings file is not valid JSON, creating new one"
        echo '{}' > "$settings_path"
    fi
    
    # Update settings using jq
    local temp_file
    temp_file=$(mktemp)
    
    if jq --arg instructions_file "$instructions_file" \
        '. += {"github.copilot.chat.codeGeneration.instructions": $instructions_file}' \
        "$settings_path" > "$temp_file"; then
        
        if mv "$temp_file" "$settings_path"; then
            log_success "Updated VS Code settings: $settings_path"
            return 0
        else
            log_error "Failed to move updated settings file"
            rm -f "$temp_file"
            return 1
        fi
    else
        log_error "Failed to update settings with jq"
        rm -f "$temp_file"
        return 1
    fi
}

# Verify installation
verify_installation() {
    local repo_path="$1"
    local settings_path
    settings_path="$(get_vscode_settings_path)"
    
    log_info "Verifying installation..."
    
    # Check copilot instructions file
    local instructions_status
    instructions_status=$(check_copilot_instructions_integrity "$repo_path")
    if [[ "$instructions_status" != "valid" ]]; then
        log_error "Copilot instructions file verification failed: $instructions_status"
        return 1
    fi
    
    # Check VS Code settings
    local vscode_status
    vscode_status=$(detect_azurerm_settings "$settings_path")
    if [[ "$vscode_status" != "present" ]]; then
        log_error "VS Code settings verification failed: $vscode_status"
        return 1
    fi
    
    log_success "Installation verification passed"
    return 0
}

# Repair installation
repair_installation() {
    local repo_path="$1"
    
    log_info "Attempting to repair installation..."
    
    # Get current status
    local status
    status=$(detect_installation_status "$repo_path")
    
    case "$status" in
        "corrupted")
            log_info "Repairing corrupted installation"
            install_ai_agent "$repo_path"
            ;;
        "partial")
            log_info "Completing partial installation"
            install_ai_agent "$repo_path"
            ;;
        "complete")
            log_info "Installation appears complete, verifying..."
            if verify_installation "$repo_path"; then
                log_success "Installation is working correctly"
                return 0
            else
                log_info "Verification failed, reinstalling..."
                install_ai_agent "$repo_path"
            fi
            ;;
        *)
            log_error "Unknown installation status: $status"
            return 1
            ;;
    esac
}

# Show installation instructions
show_installation_instructions() {
    cat << 'INSTRUCTIONS_EOF'

AI Agent Installation Complete!
===============================

The Terraform AzureRM Provider AI setup has been installed successfully.

Next Steps:
1. Restart VS Code to ensure the new settings are loaded
2. Open the terraform-provider-azurerm repository in VS Code
3. Start using GitHub Copilot with context-aware assistance!

AI Features Available:
- Context-aware code generation and review
- Azure-specific implementation patterns
- Testing guidelines and best practices
- Documentation standards enforcement
- Error handling and debugging assistance

Usage Tips:
- Ask Copilot to help implement new Azure resources
- Request guidance on testing patterns
- Get assistance with documentation formatting
- Ask for code reviews with provider-specific feedback

For more information, see the copilot-instructions.md file in the repository.

INSTRUCTIONS_EOF
}

# Handle upgrade from previous version
handle_upgrade() {
    local repo_path="$1"
    
    log_info "Handling upgrade from previous installation..."
    
    # Check if this is an upgrade scenario
    if has_previous_installation "$repo_path"; then
        log_info "Previous installation detected, upgrading..."
        
        # Perform upgrade
        install_ai_agent "$repo_path"
        
        if verify_installation "$repo_path"; then
            log_success "Upgrade completed successfully"
            return 0
        else
            log_error "Upgrade verification failed"
            return 1
        fi
    else
        log_info "No previous installation found, performing fresh install"
        install_ai_agent "$repo_path"
    fi
}
