#!/bin/bash
# Validation Module - Clean Architecture  
# Repository and installation validation using real patterns

test_repository_structure() {
    # Validates repository structure (real implementation pattern)
    local repository_path="$1"
    
    write_status_message "Validating repository structure..." "Info"
    
    local required_paths=(
        ".github/instructions"
        ".github/copilot-instructions.md"
        "internal"
    )
    
    local missing_paths=()
    local is_valid=true
    
    for path in "${required_paths[@]}"; do
        local full_path="$repository_path/$path"
        if [[ ! -e "$full_path" ]]; then
            missing_paths+=("$path")
            is_valid=false
        fi
    done
    
    # Output results as key=value pairs
    if [[ "$is_valid" == "true" ]]; then
        echo "IS_VALID_REPOSITORY=true"
    else
        echo "IS_VALID_REPOSITORY=false"
        printf "MISSING_PATHS=%s\n" "$(IFS=,; echo "${missing_paths[*]}")"
    fi
    echo "REPOSITORY_PATH=$repository_path"
}

test_current_installation() {
    # Tests current installation status (real pattern from working system)
    local repository_path="$1"
    
    write_status_message "Checking current installation status..." "Info"
    
    # Use test_ai_installation to check existing installation (correct pattern)
    test_ai_installation "$repository_path"
}

test_prerequisites() {
    # Checks if system prerequisites are met for AI installation
    
    # Check bash version (4.0+ recommended)
    if [[ ${BASH_VERSION%%.*} -lt 4 ]]; then
        write_status_message "Warning: Bash 4.0 or later is recommended, found version $BASH_VERSION" "Warning"
    fi
    
    # Use existing config infrastructure to get VS Code path
    local config_output
    config_output=$(get_installation_config "$(pwd)")
    
    if [[ $? -ne 0 ]]; then
        write_status_message "Failed to get installation configuration" "Error"
        return 1
    fi
    
    # Parse config output
    eval "$config_output"
    
    # Test if we can access/create VS Code User directory
    if [[ ! -d "$VSCODE_USER_PATH" ]]; then
        mkdir -p "$VSCODE_USER_PATH" 2>/dev/null
        if [[ $? -ne 0 ]]; then
            write_status_message "Cannot create VS Code User directory: $VSCODE_USER_PATH" "Error"
            return 1
        fi
    fi
    
    # Test write access with minimal file
    local test_file="$VSCODE_USER_PATH/.install-test"
    echo "test" > "$test_file" 2>/dev/null
    if [[ $? -eq 0 ]]; then
        rm -f "$test_file" 2>/dev/null
        write_status_message "Prerequisites check passed" "Success"
        return 0
    else
        write_status_message "Cannot write to VS Code User directory: $VSCODE_USER_PATH" "Error"
        return 1
    fi
}

# Export functions by creating a list of available functions
__VALIDATION_FUNCTIONS="test_repository_structure test_current_installation test_prerequisites"
