#!/bin/bash

# Terraform AzureRM Provider AI Installation Module (Bash)
# Cross-platform support for macOS, Linux, and other Unix-like systems

set -euo pipefail

# Include dependencies (core.sh must be first as it contains write_status_message)
# Handle both direct execution and sourcing from other directories
INSTALL_SCRIPT_DIR=""
if [[ -n "${BASH_SOURCE[0]}" ]]; then
    INSTALL_SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
elif [[ -n "$0" ]]; then
    INSTALL_SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
else
    INSTALL_SCRIPT_DIR="$(pwd)"
fi

source "$INSTALL_SCRIPT_DIR/core.sh"
source "$INSTALL_SCRIPT_DIR/validation.sh"
source "$INSTALL_SCRIPT_DIR/ui.sh"

# Determine VS Code user directory based on OS
get_vscode_user_dir() {
    case "$(uname -s)" in
        Darwin*)
            echo "$HOME/Library/Application Support/Code/User"
            ;;
        Linux*)
            if [[ -n "${XDG_CONFIG_HOME:-}" ]]; then
                echo "$XDG_CONFIG_HOME/Code/User"
            else
                echo "$HOME/.config/Code/User"
            fi
            ;;
        CYGWIN*|MINGW32*|MSYS*|MINGW*)
            echo "$APPDATA/Code/User"
            ;;
        *)
            # Fallback for other Unix-like systems
            if [[ -n "${XDG_CONFIG_HOME:-}" ]]; then
                echo "$XDG_CONFIG_HOME/Code/User"
            else
                echo "$HOME/.config/Code/User"
            fi
            ;;
    esac
}

# Function to read file manifest from repository (simplified version)
read_file_manifest_simple() {
    local repository_path="$1"
    
    # Count instruction files
    local instructions_dir="$repository_path/.github/instructions"
    local instruction_count=0
    if [[ -d "$instructions_dir" ]]; then
        instruction_count=$(find "$instructions_dir" -name "*.instructions.md" -type f 2>/dev/null | wc -l)
    fi
    
    # Count prompt files
    local prompts_dir="$repository_path/.github/prompts"
    local prompt_count=0
    if [[ -d "$prompts_dir" ]]; then
        prompt_count=$(find "$prompts_dir" -name "*.prompt.md" -type f 2>/dev/null | wc -l)
    fi
    
    # Count main files
    local main_count=0
    local main_file="$repository_path/.github/copilot-instructions.md"
    if [[ -f "$main_file" ]]; then
        main_count=1
    fi
    
    # Output counts
    echo "INSTRUCTION_COUNT=$instruction_count"
    echo "PROMPT_COUNT=$prompt_count"
    echo "MAIN_COUNT=$main_count"
}

# Function to test AI installation
test_ai_installation() {
    local repository_path="$1"
    local vscode_user_dir
    vscode_user_dir=$(get_vscode_user_dir)
    
    # Read manifest to get expected file counts
    local manifest
    manifest=$(read_file_manifest_simple "$repository_path")
    local expected_instruction_count=0
    local expected_prompt_count=0
    local expected_main_count=0
    
    while IFS='=' read -r key value; do
        case "$key" in
            "INSTRUCTION_COUNT") expected_instruction_count="$value" ;;
            "PROMPT_COUNT") expected_prompt_count="$value" ;;
            "MAIN_COUNT") expected_main_count="$value" ;;
        esac
    done <<< "$manifest"
    
    # Count actual installed files
    local instruction_count=0
    local prompt_count=0
    local main_count=0
    local settings_configured=false
    local error_count=0
    
    # Check instruction files
    local instructions_path="$vscode_user_dir/instructions/terraform-azurerm"
    if [[ -d "$instructions_path" ]]; then
        instruction_count=$(find "$instructions_path" -name "*.instructions.md" -type f 2>/dev/null | wc -l)
    fi
    
    # Check prompt files  
    local prompts_path="$vscode_user_dir/prompts"
    if [[ -d "$prompts_path" ]]; then
        prompt_count=$(find "$prompts_path" -name "*.md" -type f 2>/dev/null | wc -l)
    fi
    
    # Check main files
    local main_file="$vscode_user_dir/copilot-instructions.md"
    if [[ -f "$main_file" ]]; then
        main_count=1
    fi
    
    # Check settings.json
    local settings_path="$vscode_user_dir/settings.json"
    if [[ -f "$settings_path" ]]; then
        if grep -q "terraform-azurerm" "$settings_path" 2>/dev/null || grep -q "AZURERM_" "$settings_path" 2>/dev/null; then
            settings_configured=true
        fi
    fi
    
    # Determine overall success
    local success=false
    if [[ $instruction_count -eq $expected_instruction_count ]] && 
       [[ $prompt_count -eq $expected_prompt_count ]] && 
       [[ $main_count -eq $expected_main_count ]] && 
       [[ "$settings_configured" == "true" ]]; then
        success=true
    fi
    
    # Output results in key=value format
    echo "SUCCESS=$success"
    echo "INSTRUCTION_FILES_COUNT=$instruction_count"
    echo "EXPECTED_INSTRUCTION_FILES_COUNT=$expected_instruction_count"
    echo "PROMPT_FILES_COUNT=$prompt_count"
    echo "EXPECTED_PROMPT_FILES_COUNT=$expected_prompt_count"
    echo "MAIN_FILES_COUNT=$main_count"
    echo "EXPECTED_MAIN_FILES_COUNT=$expected_main_count"
    echo "SETTINGS_CONFIGURED=$settings_configured"
    echo "ERROR_COUNT=$error_count"
}

# Function to install instruction files
install_instruction_files() {
    local repository_path="$1"
    local force="${2:-false}"
    
    local manifest
    manifest=$(read_file_manifest "$repository_path")
    
    local success=$(echo "$manifest" | jq -r '.Success')
    if [[ "$success" != "true" ]]; then
        write_status_message "Failed to read file manifest" "Error"
        echo '{"Success": false, "Errors": ["Failed to read file manifest"]}'
        return 1
    fi
    
    local instruction_files
    instruction_files=$(echo "$manifest" | jq -r '.InstructionFiles[]?' 2>/dev/null || echo "")
    
    if [[ -z "$instruction_files" ]]; then
        write_status_message "No instruction files found in repository" "Warning"
        echo '{"Success": true, "InstalledCount": 0, "Errors": []}'
        return 0
    fi
    
    local vscode_user_dir
    vscode_user_dir=$(get_vscode_user_dir)
    local target_dir="$vscode_user_dir/terraform-azurerm"
    
    # Create target directory
    mkdir -p "$target_dir"
    
    local installed_count=0
    local errors=()
    local source_dir="$repository_path/.github/instructions"
    
    while IFS= read -r filename; do
        [[ -z "$filename" ]] && continue
        
        local source_file="$source_dir/$filename"
        local target_file="$target_dir/$filename"
        
        if [[ ! -f "$source_file" ]]; then
            errors+=("Source file not found: $filename")
            write_status_message "Source file not found: $filename" "Warning"
            continue
        fi
        
        if [[ -f "$target_file" ]] && [[ "$force" != "true" ]]; then
            write_status_message "Skipping existing file: $filename (use -Force to overwrite)" "Info"
            continue
        fi
        
        if cp "$source_file" "$target_file" 2>/dev/null; then
            write_status_message "Installed instruction file: $filename" "Success"
            ((installed_count++))
        else
            errors+=("Failed to copy $filename")
            write_status_message "Failed to copy $filename" "Error"
        fi
    done <<< "$instruction_files"
    
    write_status_message "Installed $installed_count instruction files to VS Code" "Success"
    
    local result="{\"Success\": true, \"InstalledCount\": $installed_count, \"Errors\": []}"
    if [[ ${#errors[@]} -gt 0 ]]; then
        local errors_json
        errors_json=$(printf '%s\n' "${errors[@]}" | jq -R . | jq -s .)
        result=$(echo "$result" | jq --argjson errors "$errors_json" '.Errors = $errors')
        if [[ $installed_count -eq 0 ]]; then
            result=$(echo "$result" | jq '.Success = false')
        fi
    fi
    
    echo "$result"
}

# Function to install prompt files
install_prompt_files() {
    local repository_path="$1"
    local force="${2:-false}"
    
    local manifest
    manifest=$(read_file_manifest "$repository_path")
    
    local success=$(echo "$manifest" | jq -r '.Success')
    if [[ "$success" != "true" ]]; then
        write_status_message "Failed to read file manifest" "Error"
        echo '{"Success": false, "Errors": ["Failed to read file manifest"]}'
        return 1
    fi
    
    local prompt_files
    prompt_files=$(echo "$manifest" | jq -r '.PromptFiles[]?' 2>/dev/null || echo "")
    
    if [[ -z "$prompt_files" ]]; then
        write_status_message "No prompt files found in repository" "Warning"
        echo '{"Success": true, "InstalledCount": 0, "Errors": []}'
        return 0
    fi
    
    local vscode_user_dir
    vscode_user_dir=$(get_vscode_user_dir)
    local target_dir="$vscode_user_dir/prompts/terraform-azurerm"
    
    # Create target directory
    mkdir -p "$target_dir"
    
    local installed_count=0
    local errors=()
    local source_dir="$repository_path/.github/copilot-instructions"
    
    while IFS= read -r filename; do
        [[ -z "$filename" ]] && continue
        
        local source_file="$source_dir/$filename"
        local target_file="$target_dir/$filename"
        
        if [[ ! -f "$source_file" ]]; then
            errors+=("Source file not found: $filename")
            write_status_message "Source file not found: $filename" "Warning"
            continue
        fi
        
        if [[ -f "$target_file" ]] && [[ "$force" != "true" ]]; then
            write_status_message "Skipping existing file: $filename (use -Force to overwrite)" "Info"
            continue
        fi
        
        if cp "$source_file" "$target_file" 2>/dev/null; then
            write_status_message "Installed prompt file: $filename" "Success"
            ((installed_count++))
        else
            errors+=("Failed to copy $filename")
            write_status_message "Failed to copy $filename" "Error"
        fi
    done <<< "$prompt_files"
    
    write_status_message "Installed $installed_count prompt files to VS Code" "Success"
    
    local result="{\"Success\": true, \"InstalledCount\": $installed_count, \"Errors\": []}"
    if [[ ${#errors[@]} -gt 0 ]]; then
        local errors_json
        errors_json=$(printf '%s\n' "${errors[@]}" | jq -R . | jq -s .)
        result=$(echo "$result" | jq --argjson errors "$errors_json" '.Errors = $errors')
        if [[ $installed_count -eq 0 ]]; then
            result=$(echo "$result" | jq '.Success = false')
        fi
    fi
    
    echo "$result"
}

# Function to install main files
install_main_files() {
    local repository_path="$1"
    local force="${2:-false}"
    
    local manifest
    manifest=$(read_file_manifest "$repository_path")
    
    local success=$(echo "$manifest" | jq -r '.Success')
    if [[ "$success" != "true" ]]; then
        write_status_message "Failed to read file manifest" "Error"
        echo '{"Success": false, "Errors": ["Failed to read file manifest"]}'
        return 1
    fi
    
    local main_files
    main_files=$(echo "$manifest" | jq -r '.MainFiles[]?' 2>/dev/null || echo "")
    
    if [[ -z "$main_files" ]]; then
        write_status_message "No main files found in repository" "Warning"
        echo '{"Success": true, "InstalledCount": 0, "Errors": []}'
        return 0
    fi
    
    local vscode_user_dir
    vscode_user_dir=$(get_vscode_user_dir)
    
    local installed_count=0
    local errors=()
    
    while IFS= read -r filename; do
        [[ -z "$filename" ]] && continue
        
        local source_file="$repository_path/.github/$filename"
        local target_file="$vscode_user_dir/$filename"
        
        if [[ ! -f "$source_file" ]]; then
            errors+=("Source file not found: $filename")
            write_status_message "Source file not found: $filename" "Warning"
            continue
        fi
        
        if [[ -f "$target_file" ]] && [[ "$force" != "true" ]]; then
            write_status_message "Skipping existing file: $filename (use -Force to overwrite)" "Info"
            continue
        fi
        
        if cp "$source_file" "$target_file" 2>/dev/null; then
            write_status_message "Installed main file: $filename" "Success"
            ((installed_count++))
        else
            errors+=("Failed to copy $filename")
            write_status_message "Failed to copy $filename" "Error"
        fi
    done <<< "$main_files"
    
    write_status_message "Installed $installed_count main files to VS Code" "Success"
    
    local result="{\"Success\": true, \"InstalledCount\": $installed_count, \"Errors\": []}"
    if [[ ${#errors[@]} -gt 0 ]]; then
        local errors_json
        errors_json=$(printf '%s\n' "${errors[@]}" | jq -R . | jq -s .)
        result=$(echo "$result" | jq --argjson errors "$errors_json" '.Errors = $errors')
        if [[ $installed_count -eq 0 ]]; then
            result=$(echo "$result" | jq '.Success = false')
        fi
    fi
    
    echo "$result"
}

# Function to create safe backup
create_safe_backup() {
    local source_path="$1"
    local backup_reason="${2:-backup}"
    
    # Create timestamped backup path
    local timestamp=$(date +"%Y%m%d_%H%M%S")
    local source_filename=$(basename "$source_path" .json)
    local backup_filename="${source_filename}_backup_${timestamp}.json"
    
    # Use standard backup directory
    local vscode_user_dir
    vscode_user_dir=$(get_vscode_user_dir)
    local backup_dir="$vscode_user_dir/.terraform-azurerm-backups"
    mkdir -p "$backup_dir"
    
    local backup_path="$backup_dir/$backup_filename"
    local installation_date=$(date +"%Y-%m-%d %H:%M:%S")
    
    if [[ ! -f "$source_path" ]]; then
        # File doesn't exist - create backup indicating this (length = 0)
        local no_file_backup=$(cat <<EOF
{
    "// AZURERM_BACKUP_LENGTH": 0,
    "// AZURERM_INSTALLATION_DATE": "$installation_date"
}
EOF
)
        echo "$no_file_backup" > "$backup_path"
        write_status_message "Created backup (no original file): $backup_filename" "Success"
        
        echo "{\"Success\": true, \"BackupPath\": \"$backup_path\", \"BackupFileName\": \"$backup_filename\", \"OriginalFileLength\": 0}"
        return 0
    fi
    
    # File exists - try to read it
    local original_content
    if ! original_content=$(cat "$source_path" 2>/dev/null); then
        write_status_message "Cannot read existing settings.json file - manual merge required" "Warning"
        echo "{\"Success\": false, \"Error\": \"Manual merge required: Cannot read settings.json file\"}"
        return 1
    fi
    
    local original_length=${#original_content}
    
    # Validate JSON and merge metadata
    if ! echo "$original_content" | jq . >/dev/null 2>&1; then
        write_status_message "Settings.json contains invalid JSON syntax - manual merge required" "Warning"
        echo "{\"Success\": false, \"Error\": \"Manual merge required: Settings.json contains invalid JSON syntax\"}"
        return 1
    fi
    
    # Merge metadata into JSON
    local merged_content
    merged_content=$(echo "$original_content" | jq \
        --arg length "$original_length" \
        --arg date "$installation_date" \
        '. + {"// AZURERM_BACKUP_LENGTH": ($length | tonumber), "// AZURERM_INSTALLATION_DATE": $date}')
    
    echo "$merged_content" > "$backup_path"
    write_status_message "Created backup: $backup_filename (original length: $original_length)" "Success"
    
    echo "{\"Success\": true, \"BackupPath\": \"$backup_path\", \"BackupFileName\": \"$backup_filename\", \"OriginalFileLength\": $original_length}"
}

# Function to update VS Code settings
update_vscode_settings() {
    local repository_path="$1"
    local force="${2:-false}"
    local is_reinstall="${3:-false}"
    
    local vscode_user_dir
    vscode_user_dir=$(get_vscode_user_dir)
    local settings_path="$vscode_user_dir/settings.json"
    
    write_status_message "Updating VS Code settings..." "Info"
    
    # Create backup if settings exist and not forcing or reinstalling
    if [[ -f "$settings_path" ]] && [[ "$force" != "true" ]] && [[ "$is_reinstall" != "true" ]]; then
        local backup_result
        backup_result=$(create_safe_backup "$settings_path" "AI installation")
        local backup_success=$(echo "$backup_result" | jq -r '.Success')
        
        if [[ "$backup_success" != "true" ]]; then
            local backup_error=$(echo "$backup_result" | jq -r '.Error // "Unknown error"')
            if [[ "$backup_error" == *"Manual merge required"* ]]; then
                write_status_message "MANUAL MERGE REQUIRED: Cannot read existing settings.json" "Warning"
                write_status_message "File exists but contains invalid JSON or is corrupted" "Warning"
                write_status_message "Creating minimal backup marker and proceeding with file-only installation" "Info"
                
                # Create manual merge scenario backup
                local timestamp=$(date +"%Y%m%d_%H%M%S")
                local backup_dir="$vscode_user_dir/.terraform-azurerm-backups"
                mkdir -p "$backup_dir"
                local backup_path="$backup_dir/settings_backup_${timestamp}.json"
                
                local manual_merge_backup=$(cat <<EOF
{
    "// AZURERM_BACKUP_LENGTH": -1,
    "// AZURERM_INSTALLATION_DATE": "$(date +"%Y-%m-%d %H:%M:%S")",
    "// MANUAL_MERGE_NOTICE": "This was a manual merge scenario - user must manually add AI settings to their settings.json"
}
EOF
)
                echo "$manual_merge_backup" > "$backup_path" 2>/dev/null || true
                write_status_message "Created manual merge marker: $(basename "$backup_path")" "Info"
                
                write_status_message "INSTALLATION INCOMPLETE: Settings.json requires manual configuration" "Warning"
                write_status_message "Instruction and prompt files have been installed successfully" "Success"
                write_status_message "" "Info"
                write_status_message "=== MANUAL MERGE INSTRUCTIONS ===" "Warning"
                write_status_message "Your VS Code settings.json file contains invalid JSON and could not be automatically updated." "Info"
                write_status_message "To complete the AI setup, you must manually add the following settings to your settings.json:" "Info"
                write_status_message "" "Info"
                write_status_message "1. Fix any JSON syntax errors in your current settings.json file" "Info"
                write_status_message "2. Add the AI settings from this repository's .vscode/settings.json file" "Info"
                write_status_message "3. The required settings include GitHub Copilot configuration and file associations" "Info"
                write_status_message "" "Info"
                write_status_message "Repository settings file location: $repository_path/.vscode/settings.json" "Info"
                write_status_message "Your settings file location: $settings_path" "Info"
                write_status_message "" "Info"
                write_status_message "After manually merging settings, the AI features will be fully functional." "Success"
                write_status_message "================================" "Warning"
                
                echo '{"Success": false, "ManualMergeRequired": true}'
                return 1
            fi
            
            write_status_message "Could not create backup - proceeding with manual merge scenario" "Warning"
        fi
    fi
    
    # Read existing settings or create new
    local settings='{}'
    local manual_merge_required=false
    
    if [[ -f "$settings_path" ]]; then
        local existing_content
        if existing_content=$(cat "$settings_path" 2>/dev/null); then
            # Remove comments for JSON parsing
            local clean_json
            clean_json=$(echo "$existing_content" | sed 's|//.*$||' | grep -v '^[[:space:]]*$' | tr '\n' ' ')
            
            if [[ -n "$clean_json" ]] && echo "$clean_json" | jq . >/dev/null 2>&1; then
                settings="$clean_json"
            fi
        else
            manual_merge_required=true
        fi
    fi
    
    # Exit early if manual merge is required
    if [[ "$manual_merge_required" == "true" ]]; then
        echo '{"Success": false, "ManualMergeRequired": true}'
        return 1
    fi
    
    # Add AI system settings
    local terraform_settings=$(cat <<'EOF'
{
    "github.copilot.enable": {
        "*": true,
        "plaintext": true,
        "markdown": true,
        "scminput": false
    },
    "github.copilot.advanced": {
        "debug.overrideEngine": "copilot-gpt-4o",
        "debug.testOverrideProxyUrl": "",
        "debug.overrideProxyUrl": "",
        "length": 500,
        "temperature": "",
        "top_p": "",
        "stops": {
            "*": [],
            "python": [
                "\ndef ",
                "\nclass ",
                "\nif ",
                "\n\n#"
            ]
        },
        "indentationMode": {
            "*": "automatic",
            "python": "fill",
            "markdown": "fill"
        },
        "inlineSuggestCount": 3,
        "listCount": 10,
        "secret_key": "",
        "api_url": ""
    },
    "github.copilot.chat.commitMessageGeneration.instructions": [
        {
            "text": "Write conventional commit message with scope when applicable"
        }
    ],
    "github.copilot.chat.summarizeAgentConversationHistory.enabled": true,
    "github.copilot.chat.reviewSelection.enabled": true,
    "github.copilot.chat.reviewSelection.instructions": [
        {
            "text": "Review for Terraform AzureRM provider patterns, Azure API best practices, and Go coding standards"
        }
    ],
    "files.associations": {
        "*.instructions.md": "markdown",
        ".github/*.md": "markdown"
    }
}
EOF
)
    
    # Merge with existing settings
    local updated_settings
    updated_settings=$(echo "$settings" | jq --argjson terraform "$terraform_settings" '. * $terraform')
    
    # Ensure directory exists
    mkdir -p "$(dirname "$settings_path")"
    
    # Write updated settings
    echo "$updated_settings" | jq . > "$settings_path"
    
    write_status_message "VS Code settings updated successfully" "Success"
    echo '{"Success": true}'
}

# Function to start complete installation
start_complete_installation() {
    local repository_path="$1"
    local force="${2:-false}"
    local is_reinstall="${3:-false}"
    
    write_status_message "Starting Terraform AzureRM AI installation..." "Info"
    write_status_message "Repository: $repository_path" "Info"
    
    local overall_success=true
    local all_errors=()
    
    # Install instruction files
    write_status_message "Installing instruction files..." "Info"
    local instruction_result
    instruction_result=$(install_instruction_files "$repository_path" "$force")
    local instruction_success=$(echo "$instruction_result" | jq -r '.Success')
    
    if [[ "$instruction_success" == "true" ]]; then
        local count=$(echo "$instruction_result" | jq -r '.InstalledCount')
        write_status_message "✓ Instruction files installed ($count files)" "Success"
    else
        write_status_message "✗ Failed to install instruction files" "Error"
        overall_success=false
        local errors
        errors=$(echo "$instruction_result" | jq -r '.Errors[]?' 2>/dev/null || echo "")
        while IFS= read -r error; do
            [[ -n "$error" ]] && all_errors+=("$error")
        done <<< "$errors"
    fi
    
    # Install prompt files
    write_status_message "Installing prompt files..." "Info"
    local prompt_result
    prompt_result=$(install_prompt_files "$repository_path" "$force")
    local prompt_success=$(echo "$prompt_result" | jq -r '.Success')
    
    if [[ "$prompt_success" == "true" ]]; then
        local count=$(echo "$prompt_result" | jq -r '.InstalledCount')
        write_status_message "✓ Prompt files installed ($count files)" "Success"
    else
        write_status_message "✗ Failed to install prompt files" "Error"
        overall_success=false
        local errors
        errors=$(echo "$prompt_result" | jq -r '.Errors[]?' 2>/dev/null || echo "")
        while IFS= read -r error; do
            [[ -n "$error" ]] && all_errors+=("$error")
        done <<< "$errors"
    fi
    
    # Install main files
    write_status_message "Installing main files..." "Info"
    local main_result
    main_result=$(install_main_files "$repository_path" "$force")
    local main_success=$(echo "$main_result" | jq -r '.Success')
    
    if [[ "$main_success" == "true" ]]; then
        local count=$(echo "$main_result" | jq -r '.InstalledCount')
        write_status_message "✓ Main files installed ($count files)" "Success"
    else
        write_status_message "✗ Failed to install main files" "Error"
        overall_success=false
        local errors
        errors=$(echo "$main_result" | jq -r '.Errors[]?' 2>/dev/null || echo "")
        while IFS= read -r error; do
            [[ -n "$error" ]] && all_errors+=("$error")
        done <<< "$errors"
    fi
    
    # Update VS Code settings
    write_status_message "Updating VS Code settings..." "Info"
    local settings_result
    settings_result=$(update_vscode_settings "$repository_path" "$force" "$is_reinstall")
    local settings_success=$(echo "$settings_result" | jq -r '.Success')
    
    if [[ "$settings_success" == "true" ]]; then
        write_status_message "✓ VS Code settings updated" "Success"
    else
        local manual_merge=$(echo "$settings_result" | jq -r '.ManualMergeRequired // false')
        if [[ "$manual_merge" == "true" ]]; then
            write_status_message "⚠ Manual merge required for VS Code settings" "Warning"
            overall_success=false
            all_errors+=("Manual merge required for VS Code settings")
        else
            write_status_message "✗ Failed to update VS Code settings" "Error"
            overall_success=false
            all_errors+=("Failed to update VS Code settings")
        fi
    fi
    
    # Summary
    if [[ "$overall_success" == "true" ]]; then
        write_status_message "" "Info"
        write_status_message "🎉 AI installation completed successfully!" "Success"
        write_status_message "" "Info"
        write_status_message "Next steps:" "Info"
        write_status_message "1. Restart VS Code to activate the AI features" "Info"
        write_status_message "2. Open a Terraform file to test GitHub Copilot" "Info"
        write_status_message "3. Use @workspace in Copilot Chat for context-aware assistance" "Info"
    else
        write_status_message "" "Info"
        write_status_message "⚠ AI installation completed with some issues" "Warning"
        
        if [[ ${#all_errors[@]} -gt 0 ]]; then
            write_status_message "" "Info"
            write_status_message "Issues encountered:" "Warning"
            for error in "${all_errors[@]}"; do
                write_status_message "• $error" "Warning"
            done
        fi
    fi
    
    local result="{\"Success\": $overall_success}"
    if [[ ${#all_errors[@]} -gt 0 ]]; then
        local errors_json
        errors_json=$(printf '%s\n' "${all_errors[@]}" | jq -R . | jq -s .)
        result=$(echo "$result" | jq --argjson errors "$errors_json" '.Errors = $errors')
    fi
    
    echo "$result"
}

# Function to remove AI installation
remove_ai_installation() {
    local repository_path="$1"
    
    write_status_message "Starting Terraform AzureRM AI removal..." "Info"
    
    local vscode_user_dir
    vscode_user_dir=$(get_vscode_user_dir)
    
    local result='{"Success": true, "Errors": []}'
    local errors=()
    
    # Paths
    local settings_path="$vscode_user_dir/settings.json"
    local instructions_dir="$vscode_user_dir/terraform-azurerm"
    local prompts_dir="$vscode_user_dir/prompts/terraform-azurerm"
    local backup_dir="$vscode_user_dir/.terraform-azurerm-backups"
    
    # Handle settings.json restoration
    if [[ -f "$settings_path" ]]; then
        # Check if we created this file from scratch by reading our metadata
        local should_delete=false
        local is_manual_merge=false
        
        if settings_content=$(cat "$settings_path" 2>/dev/null) && echo "$settings_content" | jq . >/dev/null 2>&1; then
            local backup_length
            backup_length=$(echo "$settings_content" | jq -r '.["// AZURERM_BACKUP_LENGTH"] // null')
            
            if [[ "$backup_length" == "0" ]]; then
                # We created this file from scratch - just delete it
                should_delete=true
                write_status_message "No original settings.json existed - will delete our created file" "Info"
            elif [[ "$backup_length" == "-1" ]]; then
                # Manual merge scenario - user manually merged settings, just clean our entries
                is_manual_merge=true
                write_status_message "Manual merge scenario detected - will clean AI settings but preserve user's merged content" "Info"
            elif [[ "$backup_length" != "null" ]] && [[ "$backup_length" -gt 0 ]]; then
                # Original file existed - try to restore from backup
                write_status_message "Original settings.json existed (length: $backup_length) - will restore from backup" "Info"
            fi
        fi
        
        if [[ "$should_delete" == "true" ]]; then
            # Delete the entire file since we created it from scratch
            if rm -f "$settings_path" 2>/dev/null; then
                write_status_message "Deleted settings.json (no original file existed)" "Success"
            else
                write_status_message "Failed to delete settings.json" "Error"
                result=$(echo "$result" | jq '.Success = false')
                errors+=("Failed to delete settings.json")
            fi
        elif [[ "$is_manual_merge" == "true" ]]; then
            # Manual merge scenario - DO NOT touch settings.json, user manually merged and must manually clean
            write_status_message "Manual merge scenario detected - user must manually clean AI settings from settings.json" "Warning"
            write_status_message "We will NOT modify your settings.json file (you manually merged, you manually clean)" "Info"
            write_status_message "Please remove Terraform AzureRM AI settings from your settings.json manually" "Warning"
            
            # Remove the manual merge backup marker file
            if [[ -d "$backup_dir" ]]; then
                while IFS= read -r -d '' backup_file; do
                    if backup_content=$(cat "$backup_file" 2>/dev/null) && echo "$backup_content" | jq . >/dev/null 2>&1; then
                        local marker_backup_length
                        marker_backup_length=$(echo "$backup_content" | jq -r '.["// AZURERM_BACKUP_LENGTH"] // null')
                        if [[ "$marker_backup_length" == "-1" ]]; then
                            rm -f "$backup_file" 2>/dev/null || true
                            write_status_message "Removed manual merge marker: $(basename "$backup_file")" "Success"
                        fi
                    fi
                done < <(find "$backup_dir" -name "settings_backup_*.json" -type f -print0 2>/dev/null || true)
            fi
            
            write_status_message "Manual merge cleanup completed - settings.json left untouched" "Success"
        else
            # Try to restore from backup or clean manually (normal scenarios only)
            local restored=false
            
            if [[ -d "$backup_dir" ]]; then
                # Find most recent backup
                local latest_backup
                latest_backup=$(find "$backup_dir" -name "settings_backup_*.json" -type f -printf '%T@ %p\n' 2>/dev/null | sort -nr | head -1 | cut -d' ' -f2- || echo "")
                
                if [[ -n "$latest_backup" ]] && [[ -f "$latest_backup" ]]; then
                    # Verify backup integrity and restore
                    if backup_content=$(cat "$latest_backup" 2>/dev/null) && echo "$backup_content" | jq . >/dev/null 2>&1; then
                        # Remove metadata from restored content
                        local cleaned_content
                        cleaned_content=$(echo "$backup_content" | jq 'del(.["// AZURERM_BACKUP_LENGTH"]) | del(.["// AZURERM_INSTALLATION_DATE"])')
                        
                        if echo "$cleaned_content" > "$settings_path" 2>/dev/null; then
                            write_status_message "VS Code settings restored from backup" "Success"
                            restored=true
                        else
                            write_status_message "Failed to restore VS Code settings" "Error"
                            result=$(echo "$result" | jq '.Success = false')
                            errors+=("Failed to restore VS Code settings from backup")
                        fi
                    fi
                fi
                
                if [[ "$restored" != "true" ]]; then
                    # No backup found or restore failed, remove our settings manually
                    if remove_terraform_settings_from_vscode "$settings_path"; then
                        write_status_message "Terraform settings removed from VS Code" "Success"
                    else
                        write_status_message "Failed to clean Terraform settings from VS Code" "Warning"
                        errors+=("Failed to clean Terraform settings from VS Code")
                    fi
                fi
                
                write_status_message "Backup directory preserved at: $backup_dir" "Info"
                write_status_message "Backups contain your original VS Code settings - keep them safe!" "Info"
            else
                # No backup directory, just clean manually
                if remove_terraform_settings_from_vscode "$settings_path"; then
                    write_status_message "Terraform settings removed from VS Code" "Success"
                else
                    write_status_message "Failed to clean Terraform settings from VS Code" "Warning"
                    errors+=("Failed to clean Terraform settings from VS Code")
                fi
            fi
        fi
    else
        write_status_message "No settings.json found - nothing to clean" "Info"
    fi
    
    # Remove instruction files
    if [[ -d "$instructions_dir" ]]; then
        write_status_message "Removing only installed instruction files..." "Info"
        
        local manifest
        manifest=$(read_file_manifest "$repository_path")
        local manifest_success=$(echo "$manifest" | jq -r '.Success')
        
        if [[ "$manifest_success" == "true" ]]; then
            local removed_count=0
            local instruction_files
            instruction_files=$(echo "$manifest" | jq -r '.InstructionFiles[]?' 2>/dev/null || echo "")
            
            while IFS= read -r filename; do
                [[ -z "$filename" ]] && continue
                
                local target_file="$instructions_dir/$filename"
                if [[ -f "$target_file" ]]; then
                    if rm -f "$target_file" 2>/dev/null; then
                        write_status_message "Removed instruction file: $filename" "Success"
                        ((removed_count++))
                    else
                        errors+=("Failed to remove $filename")
                        write_status_message "Failed to remove $filename" "Warning"
                    fi
                fi
            done <<< "$instruction_files"
            
            write_status_message "Removed $removed_count instruction files from VS Code" "Success"
            
            # Only remove directory if it's empty
            if [[ $(find "$instructions_dir" -type f 2>/dev/null | wc -l) -eq 0 ]]; then
                rmdir "$instructions_dir" 2>/dev/null && write_status_message "Removed empty terraform-azurerm instructions directory" "Success"
            else
                write_status_message "Preserved terraform-azurerm instructions directory (contains user files)" "Info"
            fi
        else
            write_status_message "Could not read manifest for removal - skipping instruction files" "Warning"
            errors+=("Failed to read manifest for instruction file removal")
        fi
    fi
    
    # Remove prompt files
    write_status_message "Removing only installed prompt files..." "Info"
    
    local manifest
    manifest=$(read_file_manifest "$repository_path")
    local manifest_success=$(echo "$manifest" | jq -r '.Success')
    
    if [[ "$manifest_success" == "true" ]]; then
        local removed_count=0
        local prompt_files
        prompt_files=$(echo "$manifest" | jq -r '.PromptFiles[]?' 2>/dev/null || echo "")
        
        # Check both possible locations: subdirectory and root
        local prompt_locations=("$prompts_dir" "$vscode_user_dir/prompts")
        
        while IFS= read -r filename; do
            [[ -z "$filename" ]] && continue
            
            local file_removed=false
            for location in "${prompt_locations[@]}"; do
                local target_file="$location/$filename"
                if [[ -f "$target_file" ]]; then
                    if rm -f "$target_file" 2>/dev/null; then
                        write_status_message "Removed prompt file: $filename from $location" "Success"
                        ((removed_count++))
                        file_removed=true
                        break
                    else
                        errors+=("Failed to remove $filename from $location")
                        write_status_message "Failed to remove $filename from $location" "Warning"
                    fi
                fi
            done
            
            if [[ "$file_removed" != "true" ]]; then
                write_status_message "Prompt file not found (may have been manually deleted): $filename" "Info"
            fi
        done <<< "$prompt_files"
        
        write_status_message "Removed $removed_count prompt files from VS Code" "Success"
        
        # Only remove terraform-azurerm subdirectory if it exists and is empty
        if [[ -d "$prompts_dir" ]] && [[ $(find "$prompts_dir" -type f 2>/dev/null | wc -l) -eq 0 ]]; then
            rmdir "$prompts_dir" 2>/dev/null && write_status_message "Removed empty prompts directory" "Success"
        else
            write_status_message "Preserved prompts directory (contains user files)" "Info"
        fi
    else
        write_status_message "Could not read manifest for removal - skipping prompt files" "Warning"
        errors+=("Failed to read manifest for prompt file removal")
    fi
    
    # Remove main files
    write_status_message "Removing only installed main files..." "Info"
    
    if [[ "$manifest_success" == "true" ]]; then
        local removed_count=0
        local main_files
        main_files=$(echo "$manifest" | jq -r '.MainFiles[]?' 2>/dev/null || echo "")
        
        while IFS= read -r filename; do
            [[ -z "$filename" ]] && continue
            
            local target_file="$vscode_user_dir/$filename"
            if [[ -f "$target_file" ]]; then
                if rm -f "$target_file" 2>/dev/null; then
                    write_status_message "Removed main file: $filename" "Success"
                    ((removed_count++))
                else
                    errors+=("Failed to remove $filename")
                    write_status_message "Failed to remove $filename" "Warning"
                fi
            else
                write_status_message "Main file not found (may have been manually deleted): $filename" "Info"
            fi
        done <<< "$main_files"
        
        write_status_message "Removed $removed_count main files from VS Code" "Success"
    else
        write_status_message "Could not read manifest for removal - skipping main files" "Warning"
        errors+=("Failed to read manifest for main file removal")
    fi
    
    # Final result
    if [[ ${#errors[@]} -gt 0 ]]; then
        local errors_json
        errors_json=$(printf '%s\n' "${errors[@]}" | jq -R . | jq -s .)
        result=$(echo "$result" | jq --argjson errors "$errors_json" '.Errors = $errors')
        
        if [[ $(echo "$result" | jq -r '.Success') == "true" ]]; then
            write_status_message "AI installation removal completed successfully!" "Success"
        else
            result=$(echo "$result" | jq '.Success = false')
            write_status_message "AI installation removal completed with some errors" "Warning"
        fi
    else
        write_status_message "AI installation removal completed successfully!" "Success"
    fi
    
    write_status_message "Your original VS Code settings backups are preserved for safety" "Info"
    
    echo "$result"
}

# Helper function to remove terraform settings from VS Code
remove_terraform_settings_from_vscode() {
    local settings_path="$1"
    
    if [[ ! -f "$settings_path" ]]; then
        return 0
    fi
    
    local content
    if ! content=$(cat "$settings_path" 2>/dev/null); then
        return 1
    fi
    
    if ! echo "$content" | jq . >/dev/null 2>&1; then
        return 1
    fi
    
    # Remove Terraform AzureRM specific settings and our metadata
    local cleaned_settings
    cleaned_settings=$(echo "$content" | jq '
        del(.["terraform_azurerm_"] // empty) |
        del(.["// AZURERM_BACKUP_LENGTH"]) |
        del(.["// AZURERM_INSTALLATION_DATE"]) |
        del(.["github.copilot.chat.commitMessageGeneration.instructions"]) |
        del(.["github.copilot.chat.summarizeAgentConversationHistory.enabled"]) |
        del(.["github.copilot.chat.reviewSelection.enabled"]) |
        del(.["github.copilot.chat.reviewSelection.instructions"]) |
        del(.["github.copilot.advanced"]) |
        del(.["github.copilot.enable"]) |
        if .["files.associations"] then
            .["files.associations"] |= del(.["*.instructions.md"]) | del(.[".github/*.md"])
        else . end |
        if (.["files.associations"] // {} | length) == 0 then
            del(.["files.associations"])
        else . end
    ')
    
    echo "$cleaned_settings" > "$settings_path"
    write_status_message "Cleaned Terraform settings from VS Code" "Success"
    return 0
}
