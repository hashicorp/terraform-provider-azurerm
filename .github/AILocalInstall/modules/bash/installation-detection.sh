#!/bin/bash

#==============================================================================
# Installation Detection Module
#==============================================================================
#
# This module handles detection of existing AI installations and VS Code
# settings to determine the appropriate installation strategy.
#

# Detect AzureRM settings in VS Code settings.json
detect_azurerm_settings() {
    local settings_path="$1"
    
    if [[ ! -f "$settings_path" ]]; then
        echo "none"
        return 0
    fi
    
    if ! is_valid_json "$settings_path"; then
        echo "corrupted"
        return 0
    fi
    
    # Check for copilot instructions
    if json_property_exists "$settings_path" "github.copilot.chat.codeGeneration.instructions"; then
        local instructions
        instructions=$(get_json_property "$settings_path" "github.copilot.chat.codeGeneration.instructions")
        if [[ "$instructions" == *"terraform-provider-azurerm"* ]]; then
            echo "present"
            return 0
        fi
    fi
    
    echo "none"
    return 0
}

# Check copilot instructions file integrity
check_copilot_instructions_integrity() {
    local repo_path="$1"
    local instructions_file="$repo_path/.github/copilot-instructions.md"
    
    if [[ ! -f "$instructions_file" ]]; then
        echo "missing"
        return 0
    fi
    
    # Check for required headers that indicate a valid instructions file
    local required_patterns=(
        "terraform-provider-azurerm"
        "Azure Resource Manager"
        "Custom instructions"
    )
    
    for pattern in "${required_patterns[@]}"; do
        if ! grep -q "$pattern" "$instructions_file"; then
            echo "corrupted"
            return 0
        fi
    done
    
    echo "valid"
    return 0
}

# Detect installation status
detect_installation_status() {
    local repo_path="$1"
    local settings_path
    settings_path="$(get_vscode_settings_path)"
    
    # Check VS Code settings
    local vscode_status
    vscode_status=$(detect_azurerm_settings "$settings_path")
    
    # Check copilot instructions file
    local instructions_status
    instructions_status=$(check_copilot_instructions_integrity "$repo_path")
    
    # Determine overall status
    case "$vscode_status:$instructions_status" in
        "none:missing")
            echo "clean"
            ;;
        "none:valid")
            echo "partial"
            ;;
        "none:corrupted")
            echo "partial"
            ;;
        "present:missing")
            echo "partial"
            ;;
        "present:valid")
            echo "complete"
            ;;
        "present:corrupted")
            echo "corrupted"
            ;;
        "corrupted:*")
            echo "corrupted"
            ;;
        *)
            echo "unknown"
            ;;
    esac
}

# Check if previous installation exists
has_previous_installation() {
    local repo_path="$1"
    local status
    status=$(detect_installation_status "$repo_path")
    
    case "$status" in
        "clean")
            return 1
            ;;
        *)
            return 0
            ;;
    esac
}

# Get installation summary
get_installation_summary() {
    local repo_path="$1"
    local settings_path
    settings_path="$(get_vscode_settings_path)"
    
    local vscode_status
    vscode_status=$(detect_azurerm_settings "$settings_path")
    
    local instructions_status
    instructions_status=$(check_copilot_instructions_integrity "$repo_path")
    
    local overall_status
    overall_status=$(detect_installation_status "$repo_path")
    
    cat << SUMMARY_EOF
Installation Summary:
  Repository Path: $repo_path
  VS Code Settings: $vscode_status
  Instructions File: $instructions_status
  Overall Status: $overall_status
  Settings Path: $settings_path
SUMMARY_EOF
    
    # Add backup information if available
    local backup_dir
    backup_dir="$(get_backup_directory)"
    if [[ -d "$backup_dir" ]]; then
        echo "  Backup Directory: $backup_dir"
        local backup_count
        backup_count=$(find "$backup_dir" -name "settings_backup_*.json" -type f | wc -l)
        echo "  Available Backups: $backup_count"
    else
        echo "  Backup Directory: none"
        echo "  Available Backups: 0"
    fi
}

# Validate VS Code settings file
validate_settings_file() {
    local settings_path="$1"
    
    if [[ ! -f "$settings_path" ]]; then
        log_info "Settings file does not exist: $settings_path"
        return 0
    fi
    
    if ! is_valid_json "$settings_path"; then
        log_error "Settings file is not valid JSON: $settings_path"
        return 1
    fi
    
    log_success "Settings file is valid JSON: $settings_path"
    return 0
}

# Check VS Code installation
check_vscode_installation() {
    local settings_path
    settings_path="$(get_vscode_settings_path)"
    local settings_dir
    settings_dir="$(dirname "$settings_path")"
    
    if [[ ! -d "$settings_dir" ]]; then
        log_error "VS Code settings directory not found: $settings_dir"
        log_error "Please ensure VS Code is installed and has been run at least once"
        return 1
    fi
    
    log_success "VS Code settings directory found: $settings_dir"
    return 0
}

# Get copilot instructions file path
get_copilot_instructions_path() {
    local repo_path="$1"
    echo "$repo_path/.github/copilot-instructions.md"
}

# Check if copilot instructions file exists
has_copilot_instructions_file() {
    local repo_path="$1"
    local instructions_file
    instructions_file=$(get_copilot_instructions_path "$repo_path")
    
    [[ -f "$instructions_file" ]]
}

# Get file sizes for comparison
get_file_info() {
    local file_path="$1"
    
    if [[ ! -f "$file_path" ]]; then
        echo "File not found: $file_path"
        return 1
    fi
    
    local size
    local modified
    size=$(wc -c < "$file_path")
    modified=$(date -r "$file_path" '+%Y-%m-%d %H:%M:%S')
    
    echo "Size: $size bytes, Modified: $modified"
}
