#!/usr/bin/env bash
#==============================================================================
# Installation Verification Module (Bash)
#==============================================================================
#
# This module provides functions to verify the integrity of the AI installation
# against the installation manifest, detecting missing or corrupted files.
#

# Source required modules
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$script_dir/core-functions.sh"

# Verify installation integrity against manifest
test_installation_integrity() {
    local repo_path="$1"
    local config_file="$repo_path/.github/AILocalInstall/file-manifest.config"
    
    if [[ -z "$repo_path" ]]; then
        log_error "Repository path is required"
        return 1
    fi
    
    # Use config file instead of manifest
    if [[ ! -f "$config_file" ]]; then
        log_error "Configuration file not found: $config_file"
        return 1
    fi
    
    log_info "Loading configuration file..."
    
    # Initialize verification results
    local overall_status=true
    local source_files_valid=true
    local target_files_valid=true
    local vscode_settings_valid=true
    local issues=()
    
    # Verify source files
    log_info "Verifying source files..."
    if ! test_source_files "$repo_path" "$config_file"; then
        source_files_valid=false
        overall_status=false
        issues+=("Source file issues detected")
    fi
    
    # Verify target files (copied to VS Code)
    log_info "Verifying target files in VS Code directory..."
    if ! test_target_files "$config_file"; then
        target_files_valid=false
        overall_status=false
        issues+=("Target file issues detected")
    fi
    
    # Verify VS Code settings
    log_info "Verifying VS Code settings..."
    if ! test_vscode_settings "$config_file"; then
        vscode_settings_valid=false
        overall_status=false
        issues+=("VS Code settings issues detected")
    fi
    
    # Display results
    show_verification_results "$overall_status" "$source_files_valid" "$target_files_valid" "$vscode_settings_valid" "${issues[@]}"
    
    if [[ "$overall_status" == true ]]; then
        return 0
    else
        return 1
    fi
}

# Test source files against manifest
test_source_files() {
    local repo_path="$1"
    local manifest_path="$2"
    local all_valid=true
    
    # Check main copilot instructions
    local copilot_path="$repo_path/.github/copilot-instructions.md"
    if ! test_file_integrity "$copilot_path" "copilot-instructions" "$manifest_path"; then
        log_error "Main copilot instructions file issues"
        all_valid=false
    fi
    
    # Check instruction files
    local instructions_dir="$repo_path/.github/instructions"
    local expected_instructions
    mapfile -t expected_instructions < <(get_expected_instruction_files)
    
    for file in "${expected_instructions[@]}"; do
        local file_path="$instructions_dir/$file"
        if ! test_file_integrity "$file_path" "instruction" "$manifest_path"; then
            log_error "Instruction file issues: $file"
            all_valid=false
        fi
    done
    
    # Check prompt files
    local prompts_dir="$repo_path/.github/prompts"
    local expected_prompts
    mapfile -t expected_prompts < <(get_expected_prompt_files)
    
    for file in "${expected_prompts[@]}"; do
        local file_path="$prompts_dir/$file"
        if ! test_file_integrity "$file_path" "prompt" "$manifest_path"; then
            log_error "Prompt file issues: $file"
            all_valid=false
        fi
    done
    
    if [[ "$all_valid" == true ]]; then
        return 0
    else
        return 1
    fi
}

# Test target files (copied to VS Code)
test_target_files() {
    local manifest_path="$1"
    local all_valid=true
    
    # Get VS Code user directory
    local vscode_user_dir
    if [[ "$OSTYPE" == "darwin"* ]]; then
        vscode_user_dir="$HOME/Library/Application Support/Code/User"
    else
        vscode_user_dir="$HOME/.config/Code/User"
    fi
    
    if [[ ! -d "$vscode_user_dir" ]]; then
        log_error "VS Code user directory not found: $vscode_user_dir"
        return 1
    fi
    
    # Check copied copilot instructions
    local copilot_target="$vscode_user_dir/copilot-instructions.md"
    if [[ ! -f "$copilot_target" ]]; then
        log_error "Missing copilot instructions in VS Code directory"
        all_valid=false
    fi
    
    # Check copied instruction files
    local instructions_target="$vscode_user_dir/instructions/terraform-azurerm"
    if [[ ! -d "$instructions_target" ]]; then
        log_error "Missing instructions directory in VS Code"
        all_valid=false
    else
        local instruction_count
        instruction_count=$(find "$instructions_target" -name "*.instructions.md" 2>/dev/null | wc -l)
        if [[ "$instruction_count" -lt 13 ]]; then
            log_error "Missing instruction files in VS Code (found $instruction_count, expected 13)"
            all_valid=false
        fi
    fi
    
    # Check copied prompt files
    local prompts_target="$vscode_user_dir/prompts/terraform-azurerm"
    if [[ ! -d "$prompts_target" ]]; then
        log_error "Missing prompts directory in VS Code"
        all_valid=false
    else
        local prompt_count
        prompt_count=$(find "$prompts_target" -name "*.prompt.md" 2>/dev/null | wc -l)
        if [[ "$prompt_count" -lt 6 ]]; then
            log_error "Missing prompt files in VS Code (found $prompt_count, expected 6)"
            all_valid=false
        fi
    fi
    
    if [[ "$all_valid" == true ]]; then
        return 0
    else
        return 1
    fi
}

# Test VS Code settings
test_vscode_settings() {
    local manifest_path="$1"
    local all_valid=true
    
    # Get VS Code user directory
    local vscode_user_dir
    if [[ "$OSTYPE" == "darwin"* ]]; then
        vscode_user_dir="$HOME/Library/Application Support/Code/User"
    else
        vscode_user_dir="$HOME/.config/Code/User"
    fi
    
    local settings_path="$vscode_user_dir/settings.json"
    
    if [[ ! -f "$settings_path" ]]; then
        log_error "VS Code settings.json not found"
        return 1
    fi
    
    if ! is_valid_json "$settings_path"; then
        log_error "VS Code settings.json is not valid JSON"
        return 1
    fi
    
    # Check required properties
    local required_properties=(
        "github.copilot.chat.commitMessageGeneration.instructions"
        "github.copilot.chat.reviewSelection.enabled"
        "github.copilot.chat.reviewSelection.instructions"
        "github.copilot.advanced"
        "github.copilot.enable"
    )
    
    for property in "${required_properties[@]}"; do
        if ! jq -e --arg prop "$property" 'has($prop)' "$settings_path" >/dev/null 2>&1; then
            log_error "Missing required VS Code setting: $property"
            all_valid=false
        fi
    done
    
    # Check metadata markers
    local metadata_markers=(
        "// AZURERM_INSTALLATION_DATE"
        "// AZURERM_BACKUP_LENGTH"
    )
    
    for marker in "${metadata_markers[@]}"; do
        if ! jq -e --arg marker "$marker" 'has($marker)' "$settings_path" >/dev/null 2>&1; then
            log_error "Missing metadata marker: $marker"
            all_valid=false
        fi
    done
    
    if [[ "$all_valid" == true ]]; then
        return 0
    else
        return 1
    fi
}

# Test individual file integrity
test_file_integrity() {
    local file_path="$1"
    local file_type="$2"
    local manifest_path="$3"
    local valid=true
    
    if [[ ! -f "$file_path" ]]; then
        log_error "File does not exist: $file_path"
        return 1
    fi
    
    # Check file size
    local file_size
    file_size=$(wc -c < "$file_path" 2>/dev/null || echo 0)
    
    # Set minimum size based on file type
    local min_size=0
    case "$file_type" in
        "copilot-instructions") min_size=1000 ;;
        "instruction") min_size=1500 ;;
        "prompt") min_size=300 ;;
    esac
    
    if [[ "$file_size" -lt "$min_size" ]]; then
        log_error "File size ($file_size bytes) is below minimum ($min_size bytes): $file_path"
        valid=false
    fi
    
    # Check for expected content based on file type
    case "$file_type" in
        "copilot-instructions")
            if ! grep -q "# Custom instructions" "$file_path" 2>/dev/null; then
                log_error "Missing expected header in copilot instructions"
                valid=false
            fi
            ;;
        "instruction")
            if ! grep -q "\.instructions\.md" "$file_path" 2>/dev/null; then
                log_error "File does not appear to be a valid instruction file: $file_path"
                valid=false
            fi
            ;;
        "prompt")
            if ! grep -q "## Instructions" "$file_path" 2>/dev/null; then
                log_error "Missing expected header in prompt file: $file_path"
                valid=false
            fi
            ;;
    esac
    
    if [[ "$valid" == true ]]; then
        return 0
    else
        return 1
    fi
}

# Display verification results
show_verification_results() {
    local overall_status="$1"
    local source_files_valid="$2"
    local target_files_valid="$3"
    local vscode_settings_valid="$4"
    shift 4
    local issues=("$@")
    
    echo
    log_info "Installation Verification Results"
    echo "=================================================="
    
    # Overall status
    if [[ "$overall_status" == true ]]; then
        log_success "[PASS] Overall Status: PASSED"
    else
        log_error "[FAIL] Overall Status: FAILED"
    fi
    
    echo
    # Source files
    echo "Source Files:"
    if [[ "$source_files_valid" == true ]]; then
        echo "  [PASS] All source files valid"
    else
        echo "  [FAIL] Source file issues detected"
    fi
    
    # Target files
    echo "📂 Target Files (VS Code):"
    if [[ "$target_files_valid" == true ]]; then
        echo "  [PASS] All target files valid"
    else
        echo "  [FAIL] Target file issues detected"
    fi
    
    # VS Code settings
    echo "⚙️ VS Code Settings:"
    if [[ "$vscode_settings_valid" == true ]]; then
        echo "  [PASS] All settings valid"
    else
        echo "  [FAIL] Settings issues detected"
    fi
    
    # Show issues if any
    if [[ "${#issues[@]}" -gt 0 ]]; then
        echo
        echo "Issues found:"
        for issue in "${issues[@]}"; do
            echo "  - $issue"
        done
    fi
    
    echo
}

# Generate integrity report
generate_integrity_report() {
    local repo_path="$1"
    local output_file="$2"
    
    if [[ -z "$output_file" ]]; then
        output_file="installation-integrity-report.txt"
    fi
    
    log_info "Generating installation integrity report..."
    
    {
        echo "Terraform AzureRM Provider AI Installation Integrity Report"
        echo "Generated: $(date -Iseconds)"
        echo "Repository: $repo_path"
        echo "=========================================================="
        echo
        
        # Capture verification output
        if test_installation_integrity "$repo_path"; then
            echo "VERIFICATION RESULT: PASSED"
        else
            echo "VERIFICATION RESULT: FAILED"
        fi
        
        echo
        echo "Detailed verification output saved above."
        echo "Run 'test_installation_integrity \"$repo_path\"' for interactive verification."
        
    } > "$output_file"
    
    log_success "Integrity report saved to: $output_file"
}

# Check installation completeness
check_installation_completeness() {
    local repo_path="$1"
    
    log_info "Checking installation completeness..."
    
    # Quick check for essential components
    local missing_components=()
    
    # Check source files
    if [[ ! -f "$repo_path/.github/copilot-instructions.md" ]]; then
        missing_components+=("Main copilot instructions")
    fi
    
    if [[ ! -d "$repo_path/.github/instructions" ]]; then
        missing_components+=("Instructions directory")
    fi
    
    if [[ ! -d "$repo_path/.github/prompts" ]]; then
        missing_components+=("Prompts directory")
    fi
    
    # Check VS Code installation
    local vscode_user_dir
    if [[ "$OSTYPE" == "darwin"* ]]; then
        vscode_user_dir="$HOME/Library/Application Support/Code/User"
    else
        vscode_user_dir="$HOME/.config/Code/User"
    fi
    
    if [[ ! -f "$vscode_user_dir/copilot-instructions.md" ]]; then
        missing_components+=("VS Code copilot instructions")
    fi
    
    if [[ ! -d "$vscode_user_dir/instructions/terraform-azurerm" ]]; then
        missing_components+=("VS Code instructions directory")
    fi
    
    if [[ ! -d "$vscode_user_dir/prompts/terraform-azurerm" ]]; then
        missing_components+=("VS Code prompts directory")
    fi
    
    # Report results
    if [[ "${#missing_components[@]}" -eq 0 ]]; then
        log_success "Installation appears complete"
        return 0
    else
        log_warning "Installation appears incomplete. Missing components:"
        for component in "${missing_components[@]}"; do
            echo "  - $component"
        done
        return 1
    fi
}

# Hardcoded installation verification (simplified check)
test_hardcoded_installation_integrity() {
    local repo_path="$1"
    
    if [[ -z "$repo_path" ]]; then
        log_error "Repository path is required"
        return 1
    fi
    
    echo "Verifying AI installation (hardcoded check)..."
    
    # Get VS Code user directory
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
    
    local has_vscode_settings=false
    local has_instructions=false
    local has_prompts=false
    
    # Check VS Code settings with backup length validation
    local settings_path="$vscode_user_dir/settings.json"
    echo "  Checking settings at: $settings_path"
    if [[ -f "$settings_path" ]]; then
        if grep -q "AZURERM_BACKUP_LENGTH\|AZURERM_INSTALLATION_DATE" "$settings_path" 2>/dev/null; then
            has_vscode_settings=true
            echo "Found AI settings in VS Code"
            
            # Check backup length values
            local backup_length
            backup_length=$(jq -r '."// AZURERM_BACKUP_LENGTH" // empty' "$settings_path" 2>/dev/null)
            if [[ "$backup_length" == "0" ]]; then
                echo "  - Backup status: No original settings.json existed (fake backup created)"
            elif [[ "$backup_length" == "-1" ]]; then
                echo "  - Backup status: Manual merge scenario detected"
            elif [[ -n "$backup_length" && "$backup_length" -gt 0 ]]; then
                echo "  - Backup status: Original settings.json backed up ($backup_length chars)"
            fi
        else
            echo "VS Code settings exist but no AI markers found"
        fi
    else
        echo "VS Code settings.json not found"
    fi
    
    # Check instruction files with specific file verification
    local instructions_dir="$vscode_user_dir/instructions/terraform-azurerm"
    if [[ -d "$instructions_dir" ]]; then
        # Check for specific instruction files from config
        local expected_instructions
        mapfile -t expected_instructions < <(get_expected_instruction_files)
        
        local found_instructions=0
        local missing_instructions=()
        
        for expected_instruction in "${expected_instructions[@]}"; do
            if [[ -f "$instructions_dir/$expected_instruction" ]]; then
                ((found_instructions++))
            else
                missing_instructions+=("$expected_instruction")
            fi
        done
        
        # Check for core copilot file in user directory
        local root_copilot_file="$vscode_user_dir/copilot-instructions.md"
        local has_copilot_file=false
        if [[ -f "$root_copilot_file" ]]; then
            has_copilot_file=true
            ((found_instructions++))
        fi
        
        if [[ $found_instructions -gt 0 ]]; then
            has_instructions=true
            echo "Found instruction files ($found_instructions of $((${#expected_instructions[@]} + 1)) expected files)"
            
            if [[ ${#missing_instructions[@]} -gt 0 ]]; then
                echo "  - Missing instructions: ${missing_instructions[*]}"
            fi
            
            # Check for core instruction file content
            if [[ "$has_copilot_file" == true ]]; then
                local file_size
                file_size=$(wc -c < "$root_copilot_file" 2>/dev/null || echo "0")
                if [[ $file_size -gt 1000 ]]; then
                    local size_kb=$((file_size / 1024))
                    echo "  - Core copilot instruction file verified (${size_kb} KB)"
                else
                    echo "  - Core copilot instruction file too small or empty"
                fi
            fi
        else
            echo "Instruction directory exists but no expected files found"
        fi
    else
        echo "No instruction files found"
    fi
    
    # Check prompts files with specific file verification
    local prompts_dir="$vscode_user_dir/prompts"
    if [[ -d "$prompts_dir" ]]; then
        # Check for specific prompt files from config
        local expected_prompts
        mapfile -t expected_prompts < <(get_expected_prompt_files)
        
        local found_prompts=0
        local missing_prompts=()
        
        for expected_prompt in "${expected_prompts[@]}"; do
            if [[ -f "$prompts_dir/$expected_prompt" ]]; then
                ((found_prompts++))
            else
                missing_prompts+=("$expected_prompt")
            fi
        done
        
        if [[ $found_prompts -gt 0 ]]; then
            has_prompts=true
            echo "Found prompt files ($found_prompts of ${#expected_prompts[@]} expected files)"
            
            if [[ ${#missing_prompts[@]} -gt 0 ]]; then
                echo "  - Missing prompts: ${missing_prompts[*]}"
            fi
        else
            echo "Prompts directory exists but no expected files found"
        fi
    else
        echo "No prompt files found"
    fi
    
    # Overall assessment with detailed criteria
    local installation_found=false
    if [[ "$has_vscode_settings" == true || "$has_instructions" == true || "$has_prompts" == true ]]; then
        installation_found=true
    fi
    
    # Calculate installation completeness score
    local score=0
    local max_score=3
    
    if [[ "$has_vscode_settings" == true ]]; then ((score++)); fi
    if [[ "$has_instructions" == true ]]; then ((score++)); fi
    if [[ "$has_prompts" == true ]]; then ((score++)); fi
    
    local completeness=$((score * 100 / max_score))
    
    echo
    echo "Installation Completeness: ${completeness}% ($score of $max_score components)"
    
    if [[ "$installation_found" == true ]]; then
        if [[ $score -eq $max_score ]]; then
            echo
            echo "AI installation verification PASSED (Complete)"
            echo "All components are properly installed and verified"
            return 0
        elif [[ $score -ge 2 ]]; then
            echo
            echo "AI installation verification PASSED (Mostly Complete)"
            echo "Most components found, installation appears functional"
            return 0
        else
            echo
            echo "AI installation verification PARTIAL"
            echo "Some components found but installation may be incomplete"
            return 1
        fi
    else
        echo
        echo "AI installation verification FAILED"
        echo "No AI enhancement components found"
        return 1
    fi
}
