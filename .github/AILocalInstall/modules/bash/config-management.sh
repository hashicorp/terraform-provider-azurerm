#!/usr/bin/env bash
#==============================================================================
# Config Management Module (Bash)
#==============================================================================
# Matches config-management.psm1 functionality
# Provides config-based installation verification and file management

# Source the core functions that have our robust config parser
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$script_dir/core-functions.sh"

# Config-based verification function (matches PowerShell Test-HardcodedInstallationIntegrity)
test_installation_integrity_config() {
    local repo_path="$1"
    local config_file="$repo_path/.github/AILocalInstall/file-manifest.config"
    
    echo "Verifying AI installation (config-based check)..."
    
    if [[ ! -f "$config_file" ]]; then
        echo "ERROR: Configuration file not found: $config_file"
        return 1
    fi
    
    # Check VS Code directory exists
    local vscode_dir
    if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" ]]; then
        vscode_dir="$APPDATA/Code/User"
    else
        vscode_dir="$HOME/.config/Code/User"
    fi
    
    if [[ -d "$vscode_dir" ]]; then
        echo "Found VS Code directory"
        if [[ -f "$vscode_dir/settings.json" ]]; then
            echo "  - Found VS Code settings.json"
        else
            echo "  - VS Code settings.json not found"
        fi
    else
        echo "VS Code directory not found: $vscode_dir"
    fi
    
    # Count instruction files using our robust config parser
    local instruction_files
    mapfile -t instruction_files < <(read_config_section "$config_file" "INSTRUCTION_FILES")
    echo "Found instruction files (${#instruction_files[@]} expected)"
    
    local found_instructions=0
    for file in "${instruction_files[@]}"; do
        if [[ -f "$repo_path/.github/instructions/$file" ]]; then
            ((found_instructions++))
        fi
    done
    echo "  - Found ${found_instructions} of ${#instruction_files[@]} instruction files"
    
    # Count prompt files
    local prompt_files
    mapfile -t prompt_files < <(read_config_section "$config_file" "PROMPT_FILES")
    echo "Found prompt files (${#prompt_files[@]} expected)"
    
    local found_prompts=0
    for file in "${prompt_files[@]}"; do
        if [[ -f "$repo_path/.github/prompts/$file" ]]; then
            ((found_prompts++))
        fi
    done
    echo "  - Found ${found_prompts} of ${#prompt_files[@]} prompt files"
    
    # Check main copilot file
    if [[ -f "$repo_path/.github/copilot-instructions.md" ]]; then
        local file_size=$(stat -f%z "$repo_path/.github/copilot-instructions.md" 2>/dev/null || stat -c%s "$repo_path/.github/copilot-instructions.md" 2>/dev/null)
        echo "  - Main copilot instruction file verified ($(echo "$file_size" | awk '{printf "%.1f KB", $1/1024}'))"
    else
        echo "  - Main copilot instruction file NOT FOUND"
    fi
    
    # Calculate completeness
    local total_expected=$((${#instruction_files[@]} + ${#prompt_files[@]} + 1)) # +1 for copilot-instructions.md
    local total_found=$((found_instructions + found_prompts))
    if [[ -f "$repo_path/.github/copilot-instructions.md" ]]; then
        ((total_found++))
    fi
    
    local completeness=$((total_found * 100 / total_expected))
    echo ""
    echo "Installation Completeness: ${completeness}% (${total_found} of ${total_expected} components)"
    echo ""
    
    if [[ $completeness -eq 100 ]]; then
        echo "AI installation verification PASSED (Complete)"
        echo "All components are properly installed and verified"
        return 0
    else
        echo "AI installation verification FAILED (Incomplete)"
        echo "Some components are missing or corrupted"
        return 1
    fi
}

# Export the function so it can be called from outside
export -f test_installation_integrity_config
