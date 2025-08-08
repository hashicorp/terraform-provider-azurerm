#!/bin/bash

#==============================================================================
# Terraform AzureRM Provider AI Setup - Bash Installation Script
#==============================================================================
#
# DESCRIPTION:
#   Installs GitHub Copilot instruction files and AI prompts for enhanced
#   development of the Terraform AzureRM Provider. This script configures
#   VS Code with specialized AI instructions, coding patterns, and best
#   practices for Azure resource development.
#
# USAGE:
#   ./install-copilot-setup.sh [OPTIONS]
#
# OPTIONS:
#   -repository-path <path>   Path to terraform-provider-azurerm repository
#   -clean                    Remove all AI setup and restore backups
#   -help                     Show detailed help information
#
# EXAMPLES:
#   ./install-copilot-setup.sh
#   ./install-copilot-setup.sh -repository-path "/home/user/terraform-provider-azurerm"
#   ./install-copilot-setup.sh -clean
#
# AUTHOR:
#   HashiCorp Terraform AzureRM Provider Team
#
# VERSION:
#   1.0.0
#
# REQUIREMENTS:
#   - Bash 3.2+ (compatible with macOS default bash)
#   - VS Code installed
#   - Git repository access
#
# NOTES:
#   This script creates backups of existing VS Code settings and provides
#   a clean uninstall option to restore the original configuration.
#
#==============================================================================

# Script parameters (optional)
REPOSITORY_PATH=""
CLEAN_MODE=""
YES_TO_ALL=""

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -repository-path|-repo-path|-r)
            REPOSITORY_PATH="$2"
            shift # past argument
            shift # past value
            ;;
        -clean|-c)
            CLEAN_MODE="true"
            shift # past argument
            ;;
        -auto-approve)
            YES_TO_ALL="true"
            shift # past argument
            ;;
        -help|-h)
            echo "Terraform AzureRM Provider AI Setup"
            echo ""
            echo "USAGE:"
            echo "  ./install-copilot-setup.sh [OPTIONS]"
            echo ""
            echo "OPTIONS:"
            echo "  -repository-path <path>   Path to terraform-provider-azurerm repository"
            echo "  -clean                    Remove all installed files and restore backups"
            echo "  -auto-approve             Skip interactive approval"
            echo "  -help                     Show this help message"
            echo ""
            echo "EXAMPLES:"
            echo "  ./install-copilot-setup.sh                                   # Auto-discover repository"
            echo "  ./install-copilot-setup.sh -repository-path /path/to/repo    # Use specific path"
            echo "  ./install-copilot-setup.sh -auto-approve                     # Non-interactive install"
            echo "  ./install-copilot-setup.sh -clean                            # Remove installation"
            echo ""
            echo "AI FEATURES:"
            echo "  - Context-aware code generation and review"
            echo "  - Azure-specific implementation patterns"
            echo "  - Testing guidelines and best practices"
            echo "  - Documentation standards enforcement"
            echo "  - Error handling and debugging assistance"
            echo ""
            exit 0
            ;;
        --*|-*)
            echo "Unknown option: $1"
            exit 1
            ;;
        *)
            echo "Error: Unexpected argument '$1'"
            echo "Use -help for usage information"
            exit 1
            ;;
    esac
done

# Function to get most recent backup file
get_most_recent_backup() {
    local user_dir="$1"

    # Find the most recent backup file
    local backup_file
    backup_file=$(find "$user_dir" -name "settings.json.backup.*" -type f 2>/dev/null | sort | tail -n 1)

    if [[ -n "$backup_file" && -f "$backup_file" ]]; then
        echo "$backup_file"
    fi
}

# Function to get backup length from settings or backup file
get_backup_length_from_settings() {
    local settings_path="$1"

    if [[ ! -f "$settings_path" ]]; then
        return 1
    fi

    # Try multiple methods to extract backup length
    local backup_length=""

    if command -v jq >/dev/null 2>&1; then
        backup_length=$(jq -r '."// AZURERM_BACKUP_LENGTH" // empty' "$settings_path" 2>/dev/null)
    elif command -v python3 >/dev/null 2>&1; then
        backup_length=$(python3 -c "
import json
try:
    with open('$settings_path', 'r') as f:
        data = json.load(f)
    print(data.get('// AZURERM_BACKUP_LENGTH', ''))
except:
    pass
" 2>/dev/null)
    else
        # Fallback: grep for the marker
        backup_length=$(grep -o '"// AZURERM_BACKUP_LENGTH"[[:space:]]*:[[:space:]]*"[^"]*"' "$settings_path" 2>/dev/null | sed 's/.*"// AZURERM_BACKUP_LENGTH"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/')
    fi

    if [[ -n "$backup_length" ]]; then
        echo "$backup_length"
        return 0
    fi

    return 1
}

# Check bash version compatibility
check_bash_version() {
    # Get bash version (e.g., "3.2.57" -> "3.2")
    local bash_version
    bash_version=$(bash --version | head -n1 | sed 's/.*version \([0-9]\+\.[0-9]\+\).*/\1/')

    # Check if version is 3.2 or higher
    local major minor
    major=$(echo "$bash_version" | cut -d. -f1)
    minor=$(echo "$bash_version" | cut -d. -f2)

    if [[ $major -lt 3 ]] || [[ $major -eq 3 && $minor -lt 2 ]]; then
        echo "Error: This script requires Bash 3.2 or higher."
        echo "Current version: $bash_version"
        echo "Please upgrade your bash or use the PowerShell script instead."
        exit 1
    fi
}

# Verify bash compatibility
check_bash_version

echo "Terraform AzureRM Provider AI Setup"
echo "====================================="

# Function to discover repository root
find_repository_root() {
    local current_path="$PWD"

    while [[ "$current_path" != "/" && "$current_path" != "" ]]; do
        if [[ -d "$current_path/.git" && -f "$current_path/go.mod" ]]; then
            # Check if this is the terraform-provider-azurerm repository
            if grep -q "module github.com/hashicorp/terraform-provider-azurerm" "$current_path/go.mod" 2>/dev/null; then
                echo "$current_path"
                return 0
            fi
        fi
        current_path="$(dirname "$current_path")"
    done

    return 1
}

# Function to validate file content integrity
test_file_integrity() {
    local file_path="$1"
    local expected_pattern="$2"

    if [[ ! -f "$file_path" ]]; then
        return 1
    fi

    # Basic length check (minimum viable content)
    if [[ $(wc -c < "$file_path" 2>/dev/null || echo "0") -lt 50 ]]; then
        return 1
    fi

    # Pattern check if provided
    if [[ -n "$expected_pattern" ]] && ! grep -q "$expected_pattern" "$file_path" 2>/dev/null; then
        return 1
    fi

    return 0
}

# Function to create safe backup with integrity verification
create_safe_backup() {
    local source_path="$1"
    local backup_path="$2"

    if [[ ! -f "$source_path" ]]; then
        return 1
    fi

    # Create backup
    if ! cp "$source_path" "$backup_path" 2>/dev/null; then
        echo "ERROR: Failed to create backup of $source_path"
        return 1
    fi

    # Verify backup integrity
    if ! test_file_integrity "$backup_path"; then
        echo "ERROR: Backup verification failed for $backup_path"
        rm -f "$backup_path" 2>/dev/null
        return 1
    fi

    echo "Backup created successfully: $backup_path"
    return 0
}

# Function to create fake backup marker for failed installations
create_fake_backup_marker() {
    local backup_path="$1"

    # Create fake backup with special marker for cleanup detection
    cat > "$backup_path" << 'EOF'
{
    "// AZURERM_BACKUP_LENGTH": "-1"
}
EOF

    if [[ $? -eq 0 ]]; then
        echo "Created recovery marker for cleanup"
        return 0
    else
        echo "ERROR: Could not create recovery marker"
        return 1
    fi
}

# Function to detect backup type and get original length
get_backup_info() {
    local backup_path="$1"

    if [[ ! -f "$backup_path" ]]; then
        echo "NO_BACKUP"
        return 1
    fi

    # Try to extract backup length marker
    if command -v jq >/dev/null 2>&1; then
        local length_marker
        length_marker=$(jq -r '."// AZURERM_BACKUP_LENGTH" // empty' "$backup_path" 2>/dev/null)

        if [[ -n "$length_marker" ]]; then
            echo "MARKER:$length_marker"
            return 0
        fi
    fi

    # Fallback: check file size
    local file_size
    file_size=$(wc -c < "$backup_path" 2>/dev/null || echo "0")
    echo "SIZE:$file_size"
    return 0
}

# Function to detect AzureRM settings in VS Code settings.json
detect_azurerm_settings() {
    local settings_path="$1"

    if [[ ! -f "$settings_path" ]]; then
        echo "false"
        return 1
    fi

    # Check for AzureRM-specific settings
    local azurerm_patterns=(
        "github.copilot.chat.commitMessageGeneration.instructions"
        "github.copilot.chat.summarizeAgentConversationHistory.enabled"
        "github.copilot.chat.reviewSelection.enabled"
        "github.copilot.chat.reviewSelection.instructions"
        "github.copilot.advanced"
        "// AZURERM_BACKUP_LENGTH"
        "*.instructions.md"
        "*.prompt.md"
        "*.azurerm.md"
    )

    for pattern in "${azurerm_patterns[@]}"; do
        if grep -q "$pattern" "$settings_path" 2>/dev/null; then
            echo "true"
            return 0
        fi
    done

    echo "false"
    return 1
}

# Function to detect previous installation
test_previous_installation() {
    local user_dir="$1"

    local indicators=(
        "$user_dir/instructions/terraform-azurerm"
        "$user_dir/prompts/terraform-azurerm-basic.md"
        "$user_dir/copilot-instructions.md"
    )

    local found_indicators=()
    for indicator in "${indicators[@]}"; do
        if [[ -e "$indicator" ]]; then
            found_indicators+=("$indicator")
        fi
    done

    if [[ ${#found_indicators[@]} -gt 0 ]]; then
        echo "true"
        for file in "${found_indicators[@]}"; do
            echo "FOUND:$file"
        done
    else
        echo "false"
    fi
}

# Function to check JSON validity
is_valid_json() {
    local json_content="$1"

    if command -v jq >/dev/null 2>&1; then
        echo "$json_content" | jq empty 2>/dev/null
        return $?
    elif command -v python3 >/dev/null 2>&1; then
        echo "$json_content" | python3 -m json.tool >/dev/null 2>&1
        return $?
    elif command -v python >/dev/null 2>&1; then
        echo "$json_content" | python -m json.tool >/dev/null 2>&1
        return $?
    else
        # Basic syntax check - look for obvious JSON structure
        if [[ "$json_content" =~ ^\{.*\}$ ]]; then
            return 0
        else
            return 1
        fi
    fi
}

# Function to merge JSON objects (bash implementation)
merge_json() {
    local original="$1"
    local new_settings="$2"

    if command -v jq >/dev/null 2>&1; then
        echo "$original" | jq --argjson new "$new_settings" '. * $new'
    else
        # Fallback: Simple string-based merge for bash without jq
        # This is a basic implementation for the most common cases
        local merged="$original"

        # Remove the closing brace to add new content
        merged="${merged%\}}"

        # Add comma if there are existing settings
        if [[ "$merged" != "{" ]]; then
            merged="$merged,"
        fi

        # Extract content between braces from new settings
        local new_content="${new_settings#\{}"
        new_content="${new_content%\}}"

        # Combine and close
        merged="$merged$new_content}"

        echo "$merged"
    fi
}

# Auto-discover repository root or use provided path
if [[ -n "$REPOSITORY_PATH" ]]; then
    REPO_ROOT="$REPOSITORY_PATH"
    echo "Using provided repository path: $REPO_ROOT"
else
    echo "Auto-discovering repository location..."
    REPO_ROOT=$(find_repository_root)
fi

if [[ -z "$REPO_ROOT" ]]; then
    echo "Repository auto-discovery failed."
    echo "Please provide the path to your terraform-provider-azurerm repository:"
    read -r REPO_ROOT

    if [[ -z "$REPO_ROOT" ]]; then
        echo "ERROR: No repository path provided"
        exit 1
    fi
fi

# Validate repository
if [[ ! -d "$REPO_ROOT" ]]; then
    echo "ERROR: Path does not exist: $REPO_ROOT"
    exit 1
fi

if [[ ! -f "$REPO_ROOT/go.mod" ]]; then
    echo "ERROR: Invalid repository: go.mod not found in $REPO_ROOT"
    echo "Please ensure you're pointing to the root of the terraform-provider-azurerm repository."
    exit 1
fi

if ! grep -q "module github.com/hashicorp/terraform-provider-azurerm" "$REPO_ROOT/go.mod" 2>/dev/null; then
    echo "ERROR: Invalid repository: not terraform-provider-azurerm"
    local found_module=$(head -n 1 "$REPO_ROOT/go.mod" 2>/dev/null | sed 's/module //')
    echo "Found module: $found_module"
    exit 1
fi

echo "Repository validated: $REPO_ROOT"

# Define paths
USER_DIR="$HOME/.vscode"
INSTRUCTIONS_DIR="$USER_DIR/instructions/terraform-azurerm"
PROMPTS_DIR="$USER_DIR/prompts"
SOURCE_DIR="$REPO_ROOT/.github"
SETTINGS_PATH="$USER_DIR/settings.json"
BACKUP_PATH="$USER_DIR/settings.json.azurerm-backup"

# Enhanced clean mode with comprehensive restoration scenarios
if [[ "$CLEAN_MODE" == "true" ]]; then
    echo "Starting cleanup process..."

    # Detect what needs to be cleaned
    previous_install_result=$(test_previous_installation "$USER_DIR")
    has_previous=$(echo "$previous_install_result" | head -n 1)

    if [[ "$has_previous" != "true" ]]; then
        echo "No previous installation detected. Nothing to clean."
        exit 0
    fi

    echo "Found previous installation files:"
    echo "$previous_install_result" | grep "^FOUND:" | sed 's/^FOUND:/  - /'
    echo ""

    # Check for settings.json cleanup options
    if [[ -f "$SETTINGS_PATH" ]]; then
        echo "Handling VS Code settings cleanup..."

        # Check for AzureRM settings in current settings.json
        has_azurerm_settings=$(detect_azurerm_settings "$SETTINGS_PATH")

        if [[ "$has_azurerm_settings" == "true" ]]; then
            echo "Found AzureRM settings in settings.json - attempting to clean..."

            # Try to find a backup or metadata to determine original state
            backup_length=$(get_backup_length_from_settings "$SETTINGS_PATH")
            most_recent_backup=$(get_most_recent_backup "$USER_DIR")

            if [[ -n "$backup_length" ]]; then
                # Found metadata in settings - use it for restoration
                echo "Found backup metadata in settings.json: original length was $backup_length bytes"

                if [[ "$backup_length" == "0" ]]; then
                    # Original was empty/non-existent
                    echo "Removing settings.json (original file was empty or didn't exist)"
                    if rm -f "$SETTINGS_PATH" 2>/dev/null; then
                        echo "Settings.json removed successfully"
                    else
                        echo "WARNING: Could not remove settings.json"
                    fi
                elif [[ -n "$most_recent_backup" ]]; then
                    # Restore from most recent backup
                    echo "Restoring from backup: $most_recent_backup"
                    if cp "$most_recent_backup" "$SETTINGS_PATH" 2>/dev/null; then
                        echo "Settings.json restored from backup"
                        rm -f "$most_recent_backup" 2>/dev/null
                        echo "Backup file cleaned up"
                    else
                        echo "ERROR: Failed to restore from backup"
                        echo "Manual restoration required from: $most_recent_backup"
                    fi
                else
                    echo "WARNING: Backup metadata found but no backup file available"
                    echo "Attempting to clean AzureRM settings from current file..."

                    # Try to clean the settings automatically
                    if clean_azurerm_settings_from_file "$SETTINGS_PATH"; then
                        echo "AzureRM settings cleaned from settings.json"
                    else
                        echo "ERROR: Failed to clean settings automatically"
                        echo "Manual cleanup required"
                    fi
                fi
            elif [[ -n "$most_recent_backup" ]]; then
                # No metadata but found a backup
                echo "Found backup file: $most_recent_backup"
                echo "Restoring from backup..."
                if cp "$most_recent_backup" "$SETTINGS_PATH" 2>/dev/null; then
                    echo "Settings.json restored from backup"
                    rm -f "$most_recent_backup" 2>/dev/null
                    echo "Backup file cleaned up"
                else
                    echo "ERROR: Failed to restore from backup"
                    echo "Manual restoration required from: $most_recent_backup"
                fi
            else
                # No backup or metadata found - manual cleanup required
                echo ""
                echo "ERROR: No backup found, but AzureRM Copilot settings detected in settings.json"
                echo ""
                echo "MANUAL CLEANUP REQUIRED:"
                echo "Since no backup exists, you must manually remove the following entries from"
                echo "$SETTINGS_PATH:"
                echo ""
                echo "Remove these entries:"
                echo "- \"github.copilot.chat.commitMessageGeneration.instructions\""
                echo "- \"github.copilot.chat.summarizeAgentConversationHistory.enabled\""
                echo "- \"github.copilot.chat.reviewSelection.enabled\""
                echo "- \"github.copilot.chat.reviewSelection.instructions\""
                echo "- \"github.copilot.advanced\" (length and temperature settings)"
                echo "- \"github.copilot.enable\" (if added by AzureRM setup)"
                echo "- \"files.associations\" entries for *.instructions.md, *.prompt.md, *.azurerm.md"
                echo "- \"// AZURERM_BACKUP_LENGTH\" (if present)"
                echo ""
                echo "BACKUP RECOMMENDATION:"
                echo "Before making changes, create a backup:"
                echo "  cp \"$SETTINGS_PATH\" \"$SETTINGS_PATH.backup-\$(date +%Y%m%d-%H%M%S)\""
                exit 1
            fi
        else
            echo "No AzureRM settings found in settings.json"

            # Check if settings.json is empty and remove it
            local is_empty=false
            if command -v jq >/dev/null 2>&1; then
                local prop_count
                prop_count=$(jq '. | length' "$SETTINGS_PATH" 2>/dev/null)
                if [[ "$prop_count" == "0" ]]; then
                    is_empty=true
                fi
            elif command -v python3 >/dev/null 2>&1; then
                local is_empty_py
                is_empty_py=$(python3 -c "
import json
try:
    with open('$SETTINGS_PATH', 'r') as f:
        data = json.load(f)
    print('true' if len(data) == 0 else 'false')
except:
    print('false')
" 2>/dev/null)
                if [[ "$is_empty_py" == "true" ]]; then
                    is_empty=true
                fi
            else
                # Fallback: check if file contains only braces and whitespace
                local trimmed_content
                trimmed_content=$(cat "$SETTINGS_PATH" 2>/dev/null | tr -d '[:space:]')
                if [[ "$trimmed_content" == "{}" ]]; then
                    is_empty=true
                fi
            fi

            if [[ "$is_empty" == "true" ]]; then
                echo "Settings.json is empty - removing unnecessary file"
                if rm -f "$SETTINGS_PATH" 2>/dev/null; then
                    echo "Empty settings.json removed (VS Code will recreate when needed)"
                else
                    echo "WARNING: Could not remove empty settings.json"
                fi
            fi
        fi
    else
        echo "No settings.json found - no cleanup needed"
    fi

    # Remove instruction files
    if [[ -d "$INSTRUCTIONS_DIR" ]]; then
        echo "Removing instruction files..."
        if rm -rf "$INSTRUCTIONS_DIR" 2>/dev/null; then
            echo "Instructions removed"
        else
            echo "ERROR: Failed to remove instructions"
        fi
    fi

    # Remove prompt files (only AzureRM prompts)
    if [[ -d "$PROMPTS_DIR" ]]; then
        echo "Removing AI prompt files..."
        find "$PROMPTS_DIR" -name "*azurerm*" -type f -delete 2>/dev/null
        find "$PROMPTS_DIR" -name "*terraform*" -type f -delete 2>/dev/null

        # Remove prompts directory if empty
        if [[ -d "$PROMPTS_DIR" && -z "$(ls -A "$PROMPTS_DIR" 2>/dev/null)" ]]; then
            rmdir "$PROMPTS_DIR" 2>/dev/null
        fi

        echo "AI prompts removed"
    fi

    # Remove main copilot instructions
    copilot_instructions_path="$USER_DIR/copilot-instructions.md"
    if [[ -f "$copilot_instructions_path" ]]; then
        echo "Removing main copilot instructions..."
        if rm -f "$copilot_instructions_path" 2>/dev/null; then
            echo "Copilot instructions removed"
        else
            echo "ERROR: Failed to remove copilot instructions"
        fi
    fi

    # Clean up all backup files
    find "$USER_DIR" -name "settings.json.backup.*" -type f -delete 2>/dev/null

    echo ""
    echo "Cleanup completed!"
    echo "Restart VS Code to ensure all changes take effect."
    exit 0
fi

# Function to clean AzureRM settings from settings.json
clean_azurerm_settings_from_file() {
    local settings_path="$1"

    if [[ ! -f "$settings_path" ]]; then
        return 1
    fi

    # Read the original content
    local original_content
    original_content=$(cat "$settings_path" 2>/dev/null)

    if [[ -z "$original_content" ]]; then
        return 1
    fi

    # Create a temporary file for processing
    local temp_file
    temp_file=$(mktemp)

    if [[ $? -ne 0 ]]; then
        return 1
    fi

    # Try to clean using various methods
    local cleaned_content="$original_content"

    # Remove GitHub Copilot settings (both string and array formats)
    cleaned_content=$(echo "$cleaned_content" | sed -E '
        s/"github\.copilot\.enable"[[:space:]]*:[[:space:]]*\{[^{}]*\}[[:space:]]*,?[[:space:]]*//g
        s/"github\.copilot\.chat\.commitMessageGeneration\.instructions"[[:space:]]*:[[:space:]]*"[^"]*"[[:space:]]*,?[[:space:]]*//g
        s/"github\.copilot\.chat\.summarizeAgentConversationHistory\.enabled"[[:space:]]*:[[:space:]]*(true|false)[[:space:]]*,?[[:space:]]*//g
        s/"github\.copilot\.chat\.reviewSelection\.enabled"[[:space:]]*:[[:space:]]*(true|false)[[:space:]]*,?[[:space:]]*//g
        s/"github\.copilot\.chat\.reviewSelection\.instructions"[[:space:]]*:[[:space:]]*"[^"]*"[[:space:]]*,?[[:space:]]*//g
        s/"github\.copilot\.advanced"[[:space:]]*:[[:space:]]*\{[^}]*\}[[:space:]]*,?[[:space:]]*//g
    ')

    # Clean up file associations
    cleaned_content=$(echo "$cleaned_content" | sed -E '
        s/"\*\.instructions\.md"[[:space:]]*:[[:space:]]*"[^"]*"[[:space:]]*,?[[:space:]]*//g
        s/"\*\.prompt\.md"[[:space:]]*:[[:space:]]*"[^"]*"[[:space:]]*,?[[:space:]]*//g
        s/"\*\.azurerm\.md"[[:space:]]*:[[:space:]]*"[^"]*"[[:space:]]*,?[[:space:]]*//g
    ')

    # Remove backup length marker
    cleaned_content=$(echo "$cleaned_content" | sed -E '
        s/"\/\/ AZURERM_BACKUP_LENGTH"[[:space:]]*:[[:space:]]*"[^"]*"[[:space:]]*,?[[:space:]]*//g
    ')

    # Clean up trailing commas and empty lines
    cleaned_content=$(echo "$cleaned_content" | sed -E '
        s/,([[:space:]]*[}\]])/\1/g
        s/\n[[:space:]]*\n/\n/g
    ')

    # Write cleaned content to temp file
    echo "$cleaned_content" > "$temp_file"

    # Test if cleaned content is valid JSON (if jq is available)
    if command -v jq >/dev/null 2>&1; then
        if ! jq . "$temp_file" >/dev/null 2>&1; then
            # If cleaning resulted in invalid JSON, the file was likely corrupted
            echo "Settings file appears corrupted after cleaning - removing entirely"
            rm -f "$settings_path" 2>/dev/null
            rm -f "$temp_file" 2>/dev/null
            echo "Corrupted settings.json removed (VS Code will recreate when needed)"
            return 0
        fi

        # Check if the result is an empty object
        local prop_count
        prop_count=$(jq '. | length' "$temp_file" 2>/dev/null)
        if [[ "$prop_count" == "0" ]]; then
            echo "Settings file contained only AzureRM metadata - removing empty file"
            rm -f "$settings_path" 2>/dev/null
            rm -f "$temp_file" 2>/dev/null
            echo "Empty settings.json removed (VS Code will recreate when needed)"
            return 0
        fi
    elif command -v python3 >/dev/null 2>&1; then
        # Use Python for JSON validation if jq is not available
        local is_valid
        is_valid=$(python3 -c "
import json
try:
    with open('$temp_file', 'r') as f:
        data = json.load(f)
    print('valid')
    print(len(data))
except:
    print('invalid')
    print('0')
" 2>/dev/null)

        local validity=$(echo "$is_valid" | head -n 1)
        local prop_count=$(echo "$is_valid" | tail -n 1)

        if [[ "$validity" != "valid" ]]; then
            echo "Settings file appears corrupted after cleaning - removing entirely"
            rm -f "$settings_path" 2>/dev/null
            rm -f "$temp_file" 2>/dev/null
            echo "Corrupted settings.json removed (VS Code will recreate when needed)"
            return 0
        fi

        if [[ "$prop_count" == "0" ]]; then
            echo "Settings file contained only AzureRM metadata - removing empty file"
            rm -f "$settings_path" 2>/dev/null
            rm -f "$temp_file" 2>/dev/null
            echo "Empty settings.json removed (VS Code will recreate when needed)"
            return 0
        fi
    fi

    # Move cleaned content back to original file
    if mv "$temp_file" "$settings_path" 2>/dev/null; then
        echo "AzureRM settings cleaned from settings.json successfully!"
        return 0
    else
        echo "ERROR: Failed to write cleaned settings back to file"
        rm -f "$temp_file" 2>/dev/null
        return 1
    fi
}

# Previous installation detection to prevent backup corruption
previous_install_result=$(test_previous_installation "$USER_DIR")
has_previous=$(echo "$previous_install_result" | head -n 1)

if [[ "$has_previous" == "true" ]]; then
    echo ""
    echo "Previous installation detected!"
    echo "Found existing files:"
    echo "$previous_install_result" | grep "^FOUND:" | sed 's/^FOUND:/  - /'
    echo ""
    echo "This will overwrite existing AzureRM AI configuration."
    if [[ "$YES_TO_ALL" == "true" ]]; then
        echo "Auto-confirming installation (-auto-approve flag provided)"
    else
        while true; do
            echo -n "Continue with installation? (Y)es/(N)o: "
            read -r response
            response=$(echo "$response" | tr '[:lower:]' '[:upper:]')
            case "$response" in
                Y) break ;;
                N) echo "Installation cancelled by user."; exit 0 ;;
                *) echo "Invalid choice. Please enter Y or N." ;;
            esac
        done
    fi
fi

# Create all necessary directories
echo "Creating directories..."
if mkdir -p "$INSTRUCTIONS_DIR" "$PROMPTS_DIR" 2>/dev/null; then
    echo "Directories created"
else
    echo "ERROR: Failed to create directories"
    exit 1
fi

# Copy all instruction files
echo "Copying instruction files..."
instructions_src_dir="$SOURCE_DIR/instructions"
if [[ -d "$instructions_src_dir" ]]; then
    instruction_count=$(find "$instructions_src_dir" -name "*.md" 2>/dev/null | wc -l)

    if [[ "$instruction_count" -eq 0 ]]; then
        echo "WARNING: No instruction files (*.md) found in $instructions_src_dir"
    else
        if cp "$instructions_src_dir"/*.md "$INSTRUCTIONS_DIR/" 2>/dev/null; then
            echo "Copied $instruction_count instruction files"
        else
            echo "ERROR: Failed to copy instruction files"
            exit 1
        fi
    fi
else
    echo "WARNING: Instructions directory not found: $instructions_src_dir"
fi

# Copy AI prompt files
echo "Copying AI prompt files..."
prompts_src_dir="$SOURCE_DIR/prompts"
if [[ -d "$prompts_src_dir" ]]; then
    prompt_count=$(find "$prompts_src_dir" -name "*.md" 2>/dev/null | wc -l)

    if [[ "$prompt_count" -eq 0 ]]; then
        echo "WARNING: No prompt files (*.md) found in $prompts_src_dir"
    else
        if cp "$prompts_src_dir"/*.md "$PROMPTS_DIR/" 2>/dev/null; then
            echo "Copied $prompt_count prompt files"
        else
            echo "ERROR: Failed to copy prompt files"
            exit 1
        fi
    fi
else
    echo "WARNING: Prompts directory not found: $prompts_src_dir"
fi

# Copy main copilot instructions
echo "Copying main copilot instructions..."
copilot_src_path="$SOURCE_DIR/copilot-instructions.md"
if [[ -f "$copilot_src_path" ]]; then
    if cp "$copilot_src_path" "$USER_DIR/" 2>/dev/null; then
        # Verify the copied file
        if test_file_integrity "$USER_DIR/copilot-instructions.md" "Terraform.*Provider.*Azure"; then
            echo "Copilot instructions copied and verified"
        else
            echo "WARNING: Copilot instructions copied but verification failed"
        fi
    else
        echo "ERROR: Failed to copy copilot instructions"
        exit 1
    fi
else
    echo "WARNING: copilot-instructions.md not found: $copilot_src_path"
fi

# Function to detect backup type and get original length
get_backup_info() {
    local backup_path="$1"

    if [[ ! -f "$backup_path" ]]; then
        echo "NO_BACKUP"
        return 1
    fi

    # Try to extract backup length marker using various methods
    if command -v jq >/dev/null 2>&1; then
        local length_marker
        length_marker=$(jq -r '."// AZURERM_BACKUP_LENGTH" // empty' "$backup_path" 2>/dev/null)

        if [[ -n "$length_marker" ]]; then
            echo "MARKER:$length_marker"
            return 0
        fi
    elif command -v python3 >/dev/null 2>&1; then
        local length_marker
        length_marker=$(python3 -c "
import json, sys
try:
    with open('$backup_path', 'r') as f:
        data = json.load(f)
    marker = data.get('// AZURERM_BACKUP_LENGTH', '')
    if marker:
        print(marker)
except:
    pass
" 2>/dev/null)

        if [[ -n "$length_marker" ]]; then
            echo "MARKER:$length_marker"
            return 0
        fi
    else
        # Fallback: grep for the marker
        local length_marker
        length_marker=$(grep -o '"// AZURERM_BACKUP_LENGTH"[[:space:]]*:[[:space:]]*"[^"]*"' "$backup_path" 2>/dev/null | sed 's/.*"// AZURERM_BACKUP_LENGTH"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/')

        if [[ -n "$length_marker" ]]; then
            echo "MARKER:$length_marker"
            return 0
        fi
    fi

    # Fallback: check file size
    local file_size
    file_size=$(wc -c < "$backup_path" 2>/dev/null || echo "0")
    echo "SIZE:$file_size"
    return 0
}

# Function to detect AzureRM settings in VS Code settings.json
detect_azurerm_settings() {
    local settings_path="$1"

    if [[ ! -f "$settings_path" ]]; then
        echo "false"
        return 1
    fi

    # Check for AzureRM-specific settings
    local azurerm_patterns=(
        "github.copilot.chat.commitMessageGeneration.instructions"
        "github.copilot.chat.summarizeAgentConversationHistory.enabled"
        "github.copilot.chat.reviewSelection.enabled"
        "github.copilot.chat.reviewSelection.instructions"
        "github.copilot.advanced"
        "// AZURERM_BACKUP_LENGTH"
        "*.instructions.md"
        "*.prompt.md"
        "*.azurerm.md"
    )

    for pattern in "${azurerm_patterns[@]}"; do
        if grep -q "$pattern" "$settings_path" 2>/dev/null; then
            echo "true"
            return 0
        fi
    done

    echo "false"
    return 1
}

# Function to check if a JSON property exists
json_has_property() {
    local json_file="$1"
    local property_path="$2"

    if [[ ! -f "$json_file" ]]; then
        echo "false"
        return 1
    fi

    if command -v jq >/dev/null 2>&1; then
        local result
        result=$(jq -r "has(\"$property_path\")" "$json_file" 2>/dev/null)
        echo "${result:-false}"
        return $?
    else
        # Fallback: simple grep check
        if grep -q "\"$property_path\"" "$json_file" 2>/dev/null; then
            echo "true"
        else
            echo "false"
        fi
        return 0
    fi
}

# Function to add property to JSON file
json_add_property() {
    local json_file="$1"
    local property_path="$2"
    local property_value="$3"
    local property_type="${4:-string}"

    if [[ ! -f "$json_file" ]]; then
        return 1
    fi

    if command -v jq >/dev/null 2>&1; then
        local temp_file
        temp_file=$(mktemp)

        case "$property_type" in
            "string")
                jq ". + {\"$property_path\": \"$property_value\"}" "$json_file" > "$temp_file"
                ;;
            "boolean")
                jq ". + {\"$property_path\": $property_value}" "$json_file" > "$temp_file"
                ;;
            "number")
                jq ". + {\"$property_path\": $property_value}" "$json_file" > "$temp_file"
                ;;
            "object")
                jq ". + {\"$property_path\": $property_value}" "$json_file" > "$temp_file"
                ;;
            *)
                jq ". + {\"$property_path\": \"$property_value\"}" "$json_file" > "$temp_file"
                ;;
        esac

        if [[ $? -eq 0 ]] && [[ -s "$temp_file" ]]; then
            mv "$temp_file" "$json_file"
            return 0
        else
            rm -f "$temp_file"
            return 1
        fi
    else
        # Fallback: manual JSON editing (basic)
        # This is a simplified approach - not perfect for all JSON structures
        local last_brace_line
        last_brace_line=$(grep -n "^}$" "$json_file" | tail -n 1 | cut -d: -f1)

        if [[ -n "$last_brace_line" ]]; then
            # Insert before the last closing brace
            case "$property_type" in
                "boolean"|"number"|"object")
                    sed -i "${last_brace_line}i\\    \"$property_path\": $property_value," "$json_file"
                    ;;
                *)
                    sed -i "${last_brace_line}i\\    \"$property_path\": \"$property_value\"," "$json_file"
                    ;;
            esac
            return 0
        fi
        return 1
    fi
}

# Function to create initial settings.json structure
create_initial_settings() {
    local json_file="$1"

    cat > "$json_file" << 'EOF'
{
}
EOF

    return $?
}

# Function to merge AzureRM settings into VS Code settings.json
merge_azurerm_settings() {
    local settings_path="$1"
    local backup_length="$2"

    # Create empty settings if file doesn't exist
    if [[ ! -f "$settings_path" ]]; then
        create_initial_settings "$settings_path"
        if [[ $? -ne 0 ]]; then
            echo "ERROR: Failed to create initial settings.json"
            return 1
        fi
    fi

    echo "Adding AzureRM Copilot configuration to settings.json..."

    # Add backup length marker
    if ! json_add_property "$settings_path" "// AZURERM_BACKUP_LENGTH" "$backup_length" "string"; then
        echo "WARNING: Could not add backup length marker"
    fi

    # GitHub Copilot chat settings
    if [[ "$(json_has_property "$settings_path" "github.copilot.chat.commitMessageGeneration.instructions")" == "false" ]]; then
        json_add_property "$settings_path" "github.copilot.chat.commitMessageGeneration.instructions" "Generate commit messages for Terraform AzureRM Provider changes following conventional commit format. Include resource names and Azure service context." "string"
    fi

    if [[ "$(json_has_property "$settings_path" "github.copilot.chat.summarizeAgentConversationHistory.enabled")" == "false" ]]; then
        json_add_property "$settings_path" "github.copilot.chat.summarizeAgentConversationHistory.enabled" "true" "boolean"
    fi

    if [[ "$(json_has_property "$settings_path" "github.copilot.chat.reviewSelection.enabled")" == "false" ]]; then
        json_add_property "$settings_path" "github.copilot.chat.reviewSelection.enabled" "true" "boolean"
    fi

    if [[ "$(json_has_property "$settings_path" "github.copilot.chat.reviewSelection.instructions")" == "false" ]]; then
        json_add_property "$settings_path" "github.copilot.chat.reviewSelection.instructions" "Review code for Terraform AzureRM Provider standards: Azure API integration, error handling patterns, schema validation, and testing requirements. Reference the specialized instruction files." "string"
    fi

    # Advanced Copilot settings
    if [[ "$(json_has_property "$settings_path" "github.copilot.advanced")" == "false" ]]; then
        local advanced_settings='{
        "length": 3000,
        "temperature": 0.1,
        "inlineSuggestCount": 3,
        "top_p": 1
    }'
        json_add_property "$settings_path" "github.copilot.advanced" "$advanced_settings" "object"
    fi

    # Enable Copilot if not explicitly set
    if [[ "$(json_has_property "$settings_path" "github.copilot.enable")" == "false" ]]; then
        json_add_property "$settings_path" "github.copilot.enable" '{"*": true}' "object"
    fi

    # File associations for instruction files
    if [[ "$(json_has_property "$settings_path" "files.associations")" == "false" ]]; then
        local file_associations='{
        "*.instructions.md": "markdown",
        "*.prompt.md": "markdown",
        "*.azurerm.md": "markdown"
    }'
        json_add_property "$settings_path" "files.associations" "$file_associations" "object"
    else
        # Add to existing file associations (simplified approach)
        echo "File associations already exist - manual merge required for:"
        echo "  *.instructions.md -> markdown"
        echo "  *.prompt.md -> markdown"
        echo "  *.azurerm.md -> markdown"
    fi

    echo "AzureRM Copilot settings merged successfully"
    return 0
}

# Function to restore settings from backup
restore_settings_from_backup() {
    local settings_path="$1"
    local backup_path="$2"

    if [[ ! -f "$backup_path" ]]; then
        echo "ERROR: Backup file not found: $backup_path"
        return 1
    fi

    if cp "$backup_path" "$settings_path" 2>/dev/null; then
        echo "Settings restored from backup"
        return 0
    else
        echo "ERROR: Failed to restore settings from backup"
        return 1
    fi
}

# Function to validate JSON file
validate_json_file() {
    local json_file="$1"

    if [[ ! -f "$json_file" ]]; then
        return 1
    fi

    if command -v jq >/dev/null 2>&1; then
        jq empty "$json_file" >/dev/null 2>&1
        return $?
    elif command -v python3 >/dev/null 2>&1; then
        python3 -c "import json; json.load(open('$json_file'))" >/dev/null 2>&1
        return $?
    else
        # Basic syntax check
        if [[ $(head -c 1 "$json_file") == "{" ]] && [[ $(tail -c 2 "$json_file" | head -c 1) == "}" ]]; then
            return 0
        else
            return 1
        fi
    fi
}

# Function to test JSON validity
test_json_validity() {
    local json_path="$1"

    if [[ ! -f "$json_path" ]]; then
        return 1
    fi

    # Try with jq first (most reliable)
    if command -v jq >/dev/null 2>&1; then
        if jq empty "$json_path" >/dev/null 2>&1; then
            return 0
        else
            return 1
        fi
    fi

    # Try with python as fallback
    if command -v python3 >/dev/null 2>&1; then
        if python3 -c "import json; json.load(open('$json_path'))" >/dev/null 2>&1; then
            return 0
        else
            return 1
        fi
    fi

    # Basic syntax check as last resort
    if grep -q "^[[:space:]]*{" "$json_path" && grep -q "}[[:space:]]*$" "$json_path"; then
        return 0
    else
        return 1
    fi
}

# Function to check if JSON property exists (handles nested properties)
json_has_property() {
    local json_path="$1"
    local property_path="$2"

    if [[ ! -f "$json_path" ]]; then
        return 1
    fi

    # Try with jq first
    if command -v jq >/dev/null 2>&1; then
        local result
        result=$(jq -r "has(\"$property_path\")" "$json_path" 2>/dev/null)
        if [[ "$result" == "true" ]]; then
            return 0
        fi

        # Check for nested properties like "github.copilot.advanced"
        if [[ "$property_path" == *"."* ]]; then
            # Convert dot notation to jq path
            local jq_path
            jq_path=$(echo "$property_path" | sed 's/\./\"][\"/'g)
            result=$(jq -r "has(\"$jq_path\")" "$json_path" 2>/dev/null)
            if [[ "$result" == "true" ]]; then
                return 0
            fi
        fi
    fi

    # Fallback to grep
    if grep -q "\"$property_path\"" "$json_path" 2>/dev/null; then
        return 0
    fi

    return 1
}

# Function to get JSON property value
json_get_property() {
    local json_path="$1"
    local property_path="$2"

    if [[ ! -f "$json_path" ]]; then
        echo ""
        return 1
    fi

    # Try with jq first
    if command -v jq >/dev/null 2>&1; then
        local result
        result=$(jq -r ".$property_path // empty" "$json_path" 2>/dev/null)
        if [[ -n "$result" && "$result" != "null" ]]; then
            echo "$result"
            return 0
        fi
    fi

    # Fallback to grep and sed
    local value
    value=$(grep "\"$property_path\"" "$json_path" 2>/dev/null | sed 's/.*: *"\([^"]*\)".*/\1/')
    echo "$value"
}

# Function to create default VS Code settings.json
create_default_settings() {
    local settings_path="$1"

    cat > "$settings_path" << 'EOF'
{
    "// AZURERM_BACKUP_LENGTH": "0",
    "github.copilot.enable": {
        "*": true,
        "plaintext": true,
        "markdown": true,
        "scminput": false
    },
    "github.copilot.advanced": {
        "secret_key": "length",
        "length": 8000,
        "temperature": 0.1,
        "top_p": 1,
        "inlineSuggestCount": 3,
        "listCount": 10,
        "indentationMode": {
            "python": "explicit",
            "javascript": "auto",
            "typescript": "auto",
            "go": "explicit"
        }
    },
    "github.copilot.chat.commitMessageGeneration.instructions": "Instructions for AzureRM Terraform Provider:\n\n1. Follow conventional commit format: TYPE: description\n   - Types: ENHANCEMENT, BUG, DOCS, STYLE, REFACTOR, TEST, CHORE\n   - Keep description under 50 characters\n\n2. For complex changes, include body with:\n   - Breaking changes (if any)\n   - New features highlighted\n   - Azure services/APIs mentioned\n   - Line length under 72 characters\n\n3. Reference issues: \"Closes Issue: #12345\" or \"Fixes Issue: #54321\"\n\n4. Examples:\n   - ENHANCEMENT: add front door profile log scrubbing support\n   - BUG: fix storage account network rules state drift\n   - DOCS: add explicit warning to fieldName for clarity",
    "github.copilot.chat.summarizeAgentConversationHistory.enabled": true,
    "github.copilot.chat.reviewSelection.enabled": true,
    "github.copilot.chat.reviewSelection.instructions": "Focus on Terraform AzureRM Provider patterns:\n\n**Code Quality Priorities:**\n1. ZERO TOLERANCE: No unnecessary comments (refactor instead)\n2. Error messages: Use %+v formatting with field names in backticks\n3. Pointer package: Use pointer.To() and pointer.From() consistently\n4. Azure patterns: Validate PATCH operations and \"None\" value handling\n\n**Implementation Verification:**\n- Typed vs Untyped resource patterns correct for context\n- CustomizeDiff validation tested with acceptance tests\n- Import functionality working with proper ID validation\n- Documentation follows alphabetical field ordering\n\n**Security & Performance:**\n- Input validation prevents injection attacks\n- Azure API calls optimized for efficiency\n- Sensitive data never logged\n- Resource cleanup handled properly",
    "files.associations": {
        "*.instructions.md": "markdown",
        "*.prompt.md": "markdown",
        "*.azurerm.md": "markdown"
    }
}
EOF

    # Verify file was created successfully
    if [[ -f "$settings_path" ]] && test_json_validity "$settings_path"; then
        return 0
    else
        return 1
    fi
}

# Function to merge VS Code settings with new AzureRM settings
merge_vscode_settings() {
    local settings_path="$1"
    local backup_created="$2"

    # Get original file size for backup length tracking
    local original_size=0
    if [[ -f "$settings_path" ]]; then
        original_size=$(wc -c < "$settings_path" 2>/dev/null || echo "0")
    fi

    # Properties to add/update for AzureRM configuration
    local azurerm_settings=(
        "github.copilot.enable"
        "github.copilot.advanced"
        "github.copilot.chat.commitMessageGeneration.instructions"
        "github.copilot.chat.summarizeAgentConversationHistory.enabled"
        "github.copilot.chat.reviewSelection.enabled"
        "github.copilot.chat.reviewSelection.instructions"
        "files.associations"
    )

    # Test if file exists and is valid JSON
    if [[ -f "$settings_path" ]] && test_json_validity "$settings_path"; then
        echo "Merging with existing settings.json..."

        # Check which properties need to be added
        local needs_merge=false
        for property in "${azurerm_settings[@]}"; do
            if ! json_has_property "$settings_path" "$property"; then
                needs_merge=true
                break
            fi
        done

        if [[ "$needs_merge" == "true" ]]; then
            # Create a temporary merged file
            local temp_file="${settings_path}.tmp"

            if command -v jq >/dev/null 2>&1; then
                # Use jq for precise merging
                jq '. + {
                    "// AZURERM_BACKUP_LENGTH": "'$original_size'",
                    "github.copilot.enable": {
                        "*": true,
                        "plaintext": true,
                        "markdown": true,
                        "scminput": false
                    },
                    "github.copilot.advanced": {
                        "secret_key": "length",
                        "length": 8000,
                        "temperature": 0.1,
                        "top_p": 1,
                        "inlineSuggestCount": 3,
                        "listCount": 10,
                        "indentationMode": {
                            "python": "explicit",
                            "javascript": "auto",
                            "typescript": "auto",
                            "go": "explicit"
                        }
                    },
                    "github.copilot.chat.commitMessageGeneration.instructions": "Instructions for AzureRM Terraform Provider:\n\n1. Follow conventional commit format: TYPE: description\n   - Types: ENHANCEMENT, BUG, DOCS, STYLE, REFACTOR, TEST, CHORE\n   - Keep description under 50 characters\n\n2. For complex changes, include body with:\n   - Breaking changes (if any)\n   - New features highlighted\n   - Azure services/APIs mentioned\n   - Line length under 72 characters\n\n3. Reference issues: \"Closes Issue: #12345\" or \"Fixes Issue: #54321\"\n\n4. Examples:\n   - ENHANCEMENT: add front door profile log scrubbing support\n   - BUG: fix storage account network rules state drift\n   - DOCS: add explicit warning to fieldName for clarity",
                    "github.copilot.chat.summarizeAgentConversationHistory.enabled": true,
                    "github.copilot.chat.reviewSelection.enabled": true,
                    "github.copilot.chat.reviewSelection.instructions": "Focus on Terraform AzureRM Provider patterns:\n\n**Code Quality Priorities:**\n1. ZERO TOLERANCE: No unnecessary comments (refactor instead)\n2. Error messages: Use %+v formatting with field names in backticks\n3. Pointer package: Use pointer.To() and pointer.From() consistently\n4. Azure patterns: Validate PATCH operations and \"None\" value handling\n\n**Implementation Verification:**\n- Typed vs Untyped resource patterns correct for context\n- CustomizeDiff validation tested with acceptance tests\n- Import functionality working with proper ID validation\n- Documentation follows alphabetical field ordering\n\n**Security & Performance:**\n- Input validation prevents injection attacks\n- Azure API calls optimized for efficiency\n- Sensitive data never logged\n- Resource cleanup handled properly",
                    "files.associations": (."files.associations" // {}) + {
                        "*.instructions.md": "markdown",
                        "*.prompt.md": "markdown",
                        "*.azurerm.md": "markdown"
                    }
                }' "$settings_path" > "$temp_file" 2>/dev/null

                if [[ $? -eq 0 ]] && test_json_validity "$temp_file"; then
                    mv "$temp_file" "$settings_path"
                    echo "Settings merged successfully using jq"
                    return 0
                else
                    rm -f "$temp_file"
                    echo "jq merge failed, falling back to manual merge"
                fi
            fi

            # Fallback to manual merge
            echo "Performing manual settings merge..."
            manual_merge_settings "$settings_path" "$original_size"
        else
            echo "All AzureRM settings already present, adding backup length marker..."
            # Just add the backup length marker
            if command -v jq >/dev/null 2>&1; then
                local temp_file="${settings_path}.tmp"
                jq '. + {"// AZURERM_BACKUP_LENGTH": "'$original_size'"}' "$settings_path" > "$temp_file" 2>/dev/null
                if [[ $? -eq 0 ]] && test_json_validity "$temp_file"; then
                    mv "$temp_file" "$settings_path"
                else
                    rm -f "$temp_file"
                fi
            fi
        fi
    else
        echo "Creating new settings.json..."
        create_default_settings "$settings_path"
    fi
}

# Function to manually merge settings when jq is not available
manual_merge_settings() {
    local settings_path="$1"
    local original_size="$2"

    # Read existing settings
    local existing_content
    if [[ -f "$settings_path" ]]; then
        existing_content=$(cat "$settings_path" 2>/dev/null)
    else
        existing_content="{}"
    fi

    # Create temporary merged file
    local temp_file="${settings_path}.tmp"

    # Start building new content
    cat > "$temp_file" << EOF
{
    "// AZURERM_BACKUP_LENGTH": "$original_size",
EOF

    # Add existing settings (excluding opening/closing braces)
    if [[ "$existing_content" != "{}" ]] && [[ -n "$existing_content" ]]; then
        echo "$existing_content" | sed '1d;$d' | sed 's/^[[:space:]]*/    /' >> "$temp_file"
        echo "," >> "$temp_file"
    fi

    # Add AzureRM-specific settings
    cat >> "$temp_file" << 'EOF'
    "github.copilot.enable": {
        "*": true,
        "plaintext": true,
        "markdown": true,
        "scminput": false
    },
    "github.copilot.advanced": {
        "secret_key": "length",
        "length": 8000,
        "temperature": 0.1,
        "top_p": 1,
        "inlineSuggestCount": 3,
        "listCount": 10,
        "indentationMode": {
            "python": "explicit",
            "javascript": "auto",
            "typescript": "auto",
            "go": "explicit"
        }
    },
    "github.copilot.chat.commitMessageGeneration.instructions": "Instructions for AzureRM Terraform Provider:\n\n1. Follow conventional commit format: TYPE: description\n   - Types: ENHANCEMENT, BUG, DOCS, STYLE, REFACTOR, TEST, CHORE\n   - Keep description under 50 characters\n\n2. For complex changes, include body with:\n   - Breaking changes (if any)\n   - New features highlighted\n   - Azure services/APIs mentioned\n   - Line length under 72 characters\n\n3. Reference issues: \"Closes Issue: #12345\" or \"Fixes Issue: #54321\"\n\n4. Examples:\n   - ENHANCEMENT: add front door profile log scrubbing support\n   - BUG: fix storage account network rules state drift\n   - DOCS: add explicit warning to fieldName for clarity",
    "github.copilot.chat.summarizeAgentConversationHistory.enabled": true,
    "github.copilot.chat.reviewSelection.enabled": true,
    "github.copilot.chat.reviewSelection.instructions": "Focus on Terraform AzureRM Provider patterns:\n\n**Code Quality Priorities:**\n1. ZERO TOLERANCE: No unnecessary comments (refactor instead)\n2. Error messages: Use %+v formatting with field names in backticks\n3. Pointer package: Use pointer.To() and pointer.From() consistently\n4. Azure patterns: Validate PATCH operations and \"None\" value handling\n\n**Implementation Verification:**\n- Typed vs Untyped resource patterns correct for context\n- CustomizeDiff validation tested with acceptance tests\n- Import functionality working with proper ID validation\n- Documentation follows alphabetical field ordering\n\n**Security & Performance:**\n- Input validation prevents injection attacks\n- Azure API calls optimized for efficiency\n- Sensitive data never logged\n- Resource cleanup handled properly",
    "files.associations": {
        "*.instructions.md": "markdown",
        "*.prompt.md": "markdown",
        "*.azurerm.md": "markdown"
    }
}
EOF

    # Verify and move the file
    if test_json_validity "$temp_file"; then
        mv "$temp_file" "$settings_path"
        echo "Manual merge completed successfully"
        return 0
    else
        rm -f "$temp_file"
        echo "Manual merge failed - invalid JSON generated"
        return 1
    fi
}

# Manual backup confirmation for extra safety
echo ""
echo "MANUAL BACKUP OPPORTUNITY"
echo "========================="
echo "Your VS Code settings file location:"
echo "   $SETTINGS_PATH"
echo ""
echo "RECOMMENDED: Create your own backup copy before proceeding!"
echo "   Example: Copy to Documents, Desktop, or another safe location"
echo ""
echo "The script will also create automatic backups, but this gives you full control."
echo ""
echo "Actions you can take now:"
echo "   - Open file manager and navigate to the path above"
echo "   - Copy settings.json to your preferred backup location"
echo "   - Choose 'C' to continue or 'X' to exit without changes"

if [[ -f "$SETTINGS_PATH" ]]; then
    local file_size=$(wc -c < "$SETTINGS_PATH" 2>/dev/null || echo "unknown")
    echo ""
    echo "Current settings.json found - size: $file_size bytes"
else
    echo ""
    echo "No existing settings.json found - new file will be created"
fi

echo ""
echo "What would you like to do?"
echo "   [C] Continue with automatic setup"
echo "   [X] Exit script (no changes made)"
echo ""

if [[ "$YES_TO_ALL" == "true" ]]; then
    echo "Auto-selecting Continue (-auto-approve flag provided)"
    choice="C"
else
    while true; do
        echo -n "Choice (C/X): "
        # Use read -n 1 if available, otherwise regular read
        if read -n 1 choice 2>/dev/null; then
            echo ""
        else
            # Fallback for systems without -n support
            read choice
        fi

        case "${choice}" in
            C|c)
                echo "Continuing with setup..."
                break
                ;;
            X|x)
                echo "Setup cancelled by user. No changes were made."
                exit 0
                ;;
            *)
                echo "Invalid choice. Please press C to continue or X to exit."
                ;;
        esac
    done
fi

case "${choice}" in
    C|c)
        echo "Continuing with setup..."
        ;;
esac

# Smart merge VS Code settings.json with enhanced backup and recovery
echo "Configuring VS Code settings (preserving existing settings)..."

existing_settings_valid=true
settings_backup_created=false

# Handle existing settings.json with comprehensive error handling
if [[ -f "$SETTINGS_PATH" ]]; then
    # Check for previous installation integrity
    if grep -q "terraform-azurerm" "$SETTINGS_PATH" 2>/dev/null && \
       grep -q "copilot-instructions.md" "$SETTINGS_PATH" 2>/dev/null; then

        # Verify installation integrity with length
        local length_line=$(grep "AZURERM_BACKUP_LENGTH" "$SETTINGS_PATH" 2>/dev/null)

        if [[ -n "$length_line" ]]; then
            # Extract length value from comment (format: "// AZURERM_BACKUP_LENGTH": "length_value")
            local stored_length=$(echo "$length_line" | sed 's/.*"AZURERM_BACKUP_LENGTH"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/')

            echo ""
            echo "AzureRM Copilot installation detected with length verification."
            echo ""
            echo "This means you've already run this installation script before."
            echo "Running it again would backup your MERGED settings instead of your ORIGINAL settings."
            echo ""

            # Check if backup exists and verify its integrity
            if [[ -f "$BACKUP_PATH" ]]; then
                if [[ "$stored_length" == "0" ]]; then
                    echo "Installation verified: Originally empty file detected"
                else
                    # Verify backup matches stored length
                    local backup_length=$(wc -c < "$BACKUP_PATH" 2>/dev/null || echo "unknown")

                    if [[ "$backup_length" == "$stored_length" ]]; then
                        echo "Installation verified: Backup integrity confirmed"
                    else
                        echo "WARNING: Backup file length mismatch detected"
                        echo "   Stored length: $stored_length bytes"
                        echo "   Backup length: $backup_length bytes"
                    fi
                fi

                echo ""
                echo "To prevent losing your original configuration, please:"
                echo "  1. Run the clean command first: ./install-copilot-setup.sh -clean"
                echo "  2. Then run the installation again: ./install-copilot-setup.sh"
                echo ""
                echo "The clean command will restore your original settings before reinstalling."
            else
                echo "WARNING: Your original settings backup is missing!"
                echo "   Expected location: $(basename "$BACKUP_PATH")"
                echo ""
                echo "This means clean mode cannot restore your original settings."
                echo "You have two options:"
                echo ""
                echo "OPTION 1 (Recommended): Manual backup first"
                echo "  1. Copy your current settings.json to a safe location"
                echo "  2. Run: ./install-copilot-setup.sh -clean (will warn about missing backup)"
                echo "  3. Manually restore your backed-up settings if needed"
                echo "  4. Run: ./install-copilot-setup.sh"
                echo ""
                echo "OPTION 2 (Advanced): Force proceed (RISKY)"
                echo "- Your current settings.json contains merged AI + original settings"
                echo "- Re-running install will merge AI settings again (may cause duplicates)"
                echo "- Only do this if you're confident about your current settings"
                echo ""
            fi
        else
            # Settings detected but no hash - corruption error
            echo ""
            echo "INSTALLATION INTEGRITY ERROR"
            echo "============================"
            echo ""
            echo "We detected AzureRM Copilot settings in your settings.json,"
            echo "but cannot verify the installation integrity."
            echo ""
            echo "This usually means:"
            echo "  - Settings were manually modified"
            echo "  - Hash comment was removed or changed"
            echo "  - Installation was partially corrupted"
            echo ""
            echo "MANUAL RECOVERY REQUIRED:"
            echo "  1. Back up your current settings.json to a safe location"
            echo "  2. Manually remove all AzureRM Copilot entries from settings.json"
            echo "  3. Run this script again for fresh installation"
            echo ""
            echo "Look for and remove lines containing:"
            echo "  - terraform-azurerm"
            echo "  - copilot-instructions.md"
            echo "  - AZURERM_BACKUP_LENGTH"
            echo ""
            echo "For detailed recovery instructions:"
            echo "https://github.com/hashicorp/terraform-provider-azurerm/blob/main/.github/instructions/README.md#troubleshooting"
            echo ""
        fi
        echo ""
        exit 1
    fi

    echo "No previous installation detected - safe to proceed with backup and merge..."

    # Test JSON validity before proceeding
    if command -v jq >/dev/null 2>&1; then
        if ! jq empty < "$SETTINGS_PATH" >/dev/null 2>&1; then
            existing_settings_valid=false
        fi
    elif command -v python3 >/dev/null 2>&1; then
        if ! python3 -c "import json; json.load(open('$SETTINGS_PATH'))" >/dev/null 2>&1; then
            existing_settings_valid=false
        fi
    fi

    if [[ "$existing_settings_valid" == "false" ]]; then
        echo "WARNING: Could not parse existing settings.json as JSON"
        echo ""

        # Create backup even for invalid JSON
        if [[ ! -f "$BACKUP_PATH" ]]; then
            if cp "$SETTINGS_PATH" "$BACKUP_PATH" 2>/dev/null; then
                echo "Backup created for invalid JSON file"
                settings_backup_created=true
            else
                echo "ERROR: Could not create backup of invalid settings.json"
                exit 1
            fi
        fi

        echo ""
        echo "MANUAL INTERVENTION REQUIRED:"
        echo "  1. Fix JSON syntax errors in: $SETTINGS_PATH"
        echo "  2. Re-run this script after fixing the JSON"
        echo "  3. Backup of current file saved to: $BACKUP_PATH"
        echo ""
        echo "Common JSON fixes:"
        echo "  - Remove trailing commas before closing braces }"
        echo "  - Ensure all strings are properly quoted"
        echo "  - Check for unescaped quotes within strings"
        echo "  - Validate brackets and braces are balanced"
        exit 1
    fi

    # Enhanced backup creation with integrity verification
    if [[ ! -f "$BACKUP_PATH" ]]; then
        if create_safe_backup "$SETTINGS_PATH" "$BACKUP_PATH"; then
            echo "Backup created successfully"
        else
            echo "ERROR: Failed to create settings backup. Installation aborted for safety."
            echo ""
            echo "TROUBLESHOOTING:"
            echo "- Ensure VS Code is not running"
            echo "- Check disk space in ~/.vscode directory"
            echo "- Verify write permissions for ~/.vscode directory"
            echo "- Try running with administrator/sudo privileges"

            # Create fake backup marker for cleanup tracking
            if create_fake_backup_marker "$BACKUP_PATH"; then
                echo ""
                echo "Recovery marker created. You can later run:"
                echo "  $0 -clean"
                echo "to clean up any partial installation."
            fi

            exit 1
        fi
    else
        echo "Note: Backup already exists: $(basename "$BACKUP_PATH")"
    fi

else
    echo "No existing settings.json found - creating new one..."
fi

# Intelligent settings merging with multiple JSON processor support
echo "Merging settings with AzureRM AI configuration..."

# Determine original file size for cleanup metadata
original_length="0"
if [[ -f "$SETTINGS_PATH" && ! -f "$BACKUP_PATH" ]]; then
    # No backup exists yet - calculate original file size
    original_length=$(wc -c < "$SETTINGS_PATH" 2>/dev/null || echo "0")
elif [[ -f "$BACKUP_PATH" ]]; then
    # Backup exists - use backup file size
    original_length=$(wc -c < "$BACKUP_PATH" 2>/dev/null || echo "0")
else
    # No original file existed
    original_length="0"
fi

# Try multiple JSON processors in order of preference
if command -v jq >/dev/null 2>&1; then
    echo "Using jq for smart JSON merging..."

    # Create temporary file with AzureRM settings
    local azure_settings_tmp="/tmp/azurerm_settings_$$.json"
    cat > "$azure_settings_tmp" << 'EOF'
{
    "github.copilot.chat.commitMessageGeneration.instructions": [
        {
            "text": "Provide a concise and clear commit message that summarizes the changes made in the code. For complex changes, include the following details: 1) Specify if the change introduces a breaking change and describe its impact. 2) Highlight any new resources or features added. 3) Mention updates to Azure services or APIs. Aim to keep the message under 72 characters per line for readability."
        }
    ],
    "github.copilot.chat.summarizeAgentConversationHistory.enabled": false,
    "github.copilot.chat.reviewSelection.enabled": true,
    "github.copilot.chat.reviewSelection.instructions": [
        {"file": "copilot-instructions.md"},
        {"file": "instructions/terraform-azurerm/implementation-guide.instructions.md"},
        {"file": "instructions/terraform-azurerm/azure-patterns.instructions.md"},
        {"file": "instructions/terraform-azurerm/testing-guidelines.instructions.md"},
        {"file": "instructions/terraform-azurerm/documentation-guidelines.instructions.md"},
        {"file": "instructions/terraform-azurerm/provider-guidelines.instructions.md"},
        {"file": "instructions/terraform-azurerm/error-patterns.instructions.md"},
        {"file": "instructions/terraform-azurerm/code-clarity-enforcement.instructions.md"},
        {"file": "instructions/terraform-azurerm/schema-patterns.instructions.md"},
        {"file": "instructions/terraform-azurerm/migration-guide.instructions.md"},
        {"file": "instructions/terraform-azurerm/performance-optimization.instructions.md"},
        {"file": "instructions/terraform-azurerm/security-compliance.instructions.md"},
        {"file": "instructions/terraform-azurerm/troubleshooting-decision-trees.instructions.md"},
        {"file": "instructions/terraform-azurerm/api-evolution-patterns.instructions.md"}
    ],
    "files.associations": {
        "*.instructions.md": "markdown",
        "*.prompt.md": "markdown",
        "*.azurerm.md": "markdown"
    },
    "github.copilot.advanced": {
        "length": 3000,
        "temperature": 0.1
    },
    "github.copilot.enable": {
        "*": true,
        "terminal": true
    }
}
EOF

    # Use jq for reliable JSON merging (AzureRM settings take precedence)
    if [[ -f "$SETTINGS_PATH" ]]; then
        jq -s '.[0] * .[1]' "$SETTINGS_PATH" "$azure_settings_tmp" > "${SETTINGS_PATH}.tmp"
    else
        cp "$azure_settings_tmp" "${SETTINGS_PATH}.tmp"
    fi

    # Add length comment to merged settings for future verification
    jq --arg length "$original_length" '. + {"// AZURERM_BACKUP_LENGTH": $length}' "${SETTINGS_PATH}.tmp" > "$SETTINGS_PATH"
    rm -f "${SETTINGS_PATH}.tmp" "$azure_settings_tmp"

    echo "Settings merged successfully using jq!"

elif command -v python3 >/dev/null 2>&1; then
    echo "Using Python for JSON merging..."

    python3 -c "
import json
import sys

settings_path = '$SETTINGS_PATH'
original_length = '$original_length'

# Read existing settings
try:
    if '$existing_settings_valid' == 'true':
        with open(settings_path, 'r') as f:
            existing = json.load(f)
    else:
        existing = {}
except:
    existing = {}

# Define comprehensive AzureRM settings
azurerm = {
    'github.copilot.chat.commitMessageGeneration.instructions': [
        {
            'text': 'Provide a concise and clear commit message that summarizes the changes made in the code. For complex changes, include the following details: 1) Specify if the change introduces a breaking change and describe its impact. 2) Highlight any new resources or features added. 3) Mention updates to Azure services or APIs. Aim to keep the message under 72 characters per line for readability.'
        }
    ],
    'github.copilot.chat.summarizeAgentConversationHistory.enabled': False,
    'github.copilot.chat.reviewSelection.enabled': True,
    'github.copilot.chat.reviewSelection.instructions': [
        {'file': 'copilot-instructions.md'},
        {'file': 'instructions/terraform-azurerm/implementation-guide.instructions.md'},
        {'file': 'instructions/terraform-azurerm/azure-patterns.instructions.md'},
        {'file': 'instructions/terraform-azurerm/testing-guidelines.instructions.md'},
        {'file': 'instructions/terraform-azurerm/documentation-guidelines.instructions.md'},
        {'file': 'instructions/terraform-azurerm/provider-guidelines.instructions.md'},
        {'file': 'instructions/terraform-azurerm/error-patterns.instructions.md'},
        {'file': 'instructions/terraform-azurerm/code-clarity-enforcement.instructions.md'},
        {'file': 'instructions/terraform-azurerm/schema-patterns.instructions.md'},
        {'file': 'instructions/terraform-azurerm/migration-guide.instructions.md'},
        {'file': 'instructions/terraform-azurerm/performance-optimization.instructions.md'},
        {'file': 'instructions/terraform-azurerm/security-compliance.instructions.md'},
        {'file': 'instructions/terraform-azurerm/troubleshooting-decision-trees.instructions.md'},
        {'file': 'instructions/terraform-azurerm/api-evolution-patterns.instructions.md'}
    ],
    'github.copilot.advanced': {
        'length': 3000,
        'temperature': 0.1
    },
    'github.copilot.enable': {
        '*': True,
        'terminal': True
    }
}

# Enhanced file associations merging
new_file_associations = {
    '*.instructions.md': 'markdown',
    '*.prompt.md': 'markdown',
    '*.azurerm.md': 'markdown'
}

if 'files.associations' in existing:
    new_file_associations.update(existing['files.associations'])

azurerm['files.associations'] = new_file_associations

# Intelligent merging (preserves user settings, adds AzureRM enhancements)
for key, value in azurerm.items():
    if key in existing:
        if key == 'files.associations':
            existing[key] = value  # Complete replacement for file associations
        elif key == 'github.copilot.enable':
            # Merge enable settings (preserve user's disable choices)
            if isinstance(existing[key], dict):
                for enable_key, enable_value in value.items():
                    if enable_key not in existing[key]:
                        existing[key][enable_key] = enable_value
            else:
                existing[key] = value
        else:
            # For other settings, AzureRM takes precedence
            existing[key] = value
    else:
        existing[key] = value

# Add backup metadata for cleanup
existing['// AZURERM_BACKUP_LENGTH'] = original_length

# Write merged settings
with open(settings_path, 'w') as f:
    json.dump(existing, f, indent=2)

print('Settings merged successfully using Python!')
" 2>/dev/null

    if [[ $? -ne 0 ]]; then
        echo "ERROR: Python JSON merge failed"
        echo ""
        echo "This may be due to invalid JSON syntax in your existing settings.json"
        echo ""

        # Create fake backup marker for cleanup detection
        if create_fake_backup_marker "$BACKUP_PATH"; then
            echo "Created recovery marker for cleanup"
        fi

        echo ""
        echo "MANUAL INTERVENTION REQUIRED:"
        echo "  1. Fix JSON syntax errors in: $SETTINGS_PATH"
        echo "  2. Manually add the following to your corrected settings.json:"
        echo ""
        echo '     {'
        echo '       "github.copilot.chat.commitMessageGeneration.instructions": ['
        echo '         {'
        echo '           "text": "Provide a concise and clear commit message..."'
        echo '         }'
        echo '       ],'
        echo '       "github.copilot.chat.summarizeAgentConversationHistory.enabled": false,'
        echo '       "github.copilot.chat.reviewSelection.enabled": true,'
        echo '       "github.copilot.chat.reviewSelection.instructions": ['
        echo "         {\"file\": \"$USER_DIR/copilot-instructions.md\"},"
        echo "         {\"file\": \"$USER_DIR/instructions/terraform-azurerm/implementation-guide.instructions.md\"},"
        echo "         {\"file\": \"$USER_DIR/instructions/terraform-azurerm/azure-patterns.instructions.md\"},"
        echo "         {\"file\": \"$USER_DIR/instructions/terraform-azurerm/testing-guidelines.instructions.md\"},"
        echo "         {\"file\": \"$USER_DIR/instructions/terraform-azurerm/documentation-guidelines.instructions.md\"},"
        echo "         {\"file\": \"$USER_DIR/instructions/terraform-azurerm/provider-guidelines.instructions.md\"},"
        echo "         {\"file\": \"$USER_DIR/instructions/terraform-azurerm/error-patterns.instructions.md\"},"
        echo "         {\"file\": \"$USER_DIR/instructions/terraform-azurerm/code-clarity-enforcement.instructions.md\"},"
        echo "         {\"file\": \"$USER_DIR/instructions/terraform-azurerm/schema-patterns.instructions.md\"},"
        echo "         {\"file\": \"$USER_DIR/instructions/terraform-azurerm/migration-guide.instructions.md\"},"
        echo "         {\"file\": \"$USER_DIR/instructions/terraform-azurerm/performance-optimization.instructions.md\"},"
        echo "         {\"file\": \"$USER_DIR/instructions/terraform-azurerm/security-compliance.instructions.md\"},"
        echo "         {\"file\": \"$USER_DIR/instructions/terraform-azurerm/troubleshooting-decision-trees.instructions.md\"},"
        echo "         {\"file\": \"$USER_DIR/instructions/terraform-azurerm/api-evolution-patterns.instructions.md\"}"
        echo '       ],'
        echo '       "github.copilot.advanced": {'
        echo '         "length": 3000,'
        echo '         "temperature": 0.1'
        echo '       },'
        echo '       "github.copilot.enable": {'
        echo '         "*": true,'
        echo '         "terminal": true'
        echo '       },'
        echo '       "files.associations": {'
        echo '         "*.instructions.md": "markdown",'
        echo '         "*.prompt.md": "markdown",'
        echo '         "*.azurerm.md": "markdown"'
        echo '       }'
        echo '       // ...your existing settings...'
        echo '     }'
        echo ""
        echo "Common JSON fixes:"
        echo "  - Remove trailing commas before closing braces }"
        echo "  - Ensure all strings are properly quoted"
        echo "  - Check for unescaped quotes within strings"
        echo "  - Validate brackets and braces are balanced"
        echo ""
        echo "IMPORTANT: Make a backup copy of your settings.json before editing!"
        echo "  Example: cp '$SETTINGS_PATH' '$SETTINGS_PATH.backup-\$(date +%Y%m%d-%H%M%S)'"
        exit 1
    fi

elif command -v node >/dev/null 2>&1; then
    echo "Using Node.js for JSON merging..."

    node -e "
const fs = require('fs');

const settingsPath = '$SETTINGS_PATH';
const originalLength = '$original_length';

// Read existing settings
let existing = {};
try {
    if ('$existing_settings_valid' === 'true') {
        existing = JSON.parse(fs.readFileSync(settingsPath, 'utf8'));
    }
} catch (e) {
    existing = {};
}

// Define comprehensive AzureRM settings
const azurerm = {
    'github.copilot.chat.commitMessageGeneration.instructions': [
        {
            'text': 'Provide a concise and clear commit message that summarizes the changes made in the code. For complex changes, include the following details: 1) Specify if the change introduces a breaking change and describe its impact. 2) Highlight any new resources or features added. 3) Mention updates to Azure services or APIs. Aim to keep the message under 72 characters per line for readability.'
        }
    ],
    'github.copilot.chat.summarizeAgentConversationHistory.enabled': false,
    'github.copilot.chat.reviewSelection.enabled': true,
    'github.copilot.chat.reviewSelection.instructions': [
        {'file': 'copilot-instructions.md'},
        {'file': 'instructions/terraform-azurerm/implementation-guide.instructions.md'},
        {'file': 'instructions/terraform-azurerm/azure-patterns.instructions.md'},
        {'file': 'instructions/terraform-azurerm/testing-guidelines.instructions.md'},
        {'file': 'instructions/terraform-azurerm/documentation-guidelines.instructions.md'},
        {'file': 'instructions/terraform-azurerm/provider-guidelines.instructions.md'},
        {'file': 'instructions/terraform-azurerm/error-patterns.instructions.md'},
        {'file': 'instructions/terraform-azurerm/code-clarity-enforcement.instructions.md'},
        {'file': 'instructions/terraform-azurerm/schema-patterns.instructions.md'},
        {'file': 'instructions/terraform-azurerm/migration-guide.instructions.md'},
        {'file': 'instructions/terraform-azurerm/performance-optimization.instructions.md'},
        {'file': 'instructions/terraform-azurerm/security-compliance.instructions.md'},
        {'file': 'instructions/terraform-azurerm/troubleshooting-decision-trees.instructions.md'},
        {'file': 'instructions/terraform-azurerm/api-evolution-patterns.instructions.md'}
    ],
    'github.copilot.advanced': {
        'length': 3000,
        'temperature': 0.1
    },
    'github.copilot.enable': {
        '*': true,
        'terminal': true
    }
};

// Enhanced file associations merging
const newFileAssociations = {
    '*.instructions.md': 'markdown',
    '*.prompt.md': 'markdown',
    '*.azurerm.md': 'markdown'
};

if (existing['files.associations']) {
    Object.assign(newFileAssociations, existing['files.associations']);
}

azurerm['files.associations'] = newFileAssociations;

// Intelligent merging
Object.keys(azurerm).forEach(key => {
    if (key in existing) {
        if (key === 'files.associations') {
            existing[key] = azurerm[key];
        } else if (key === 'github.copilot.enable') {
            if (typeof existing[key] === 'object' && existing[key] !== null) {
                Object.keys(azurerm[key]).forEach(enableKey => {
                    if (!(enableKey in existing[key])) {
                        existing[key][enableKey] = azurerm[key][enableKey];
                    }
                });
            } else {
                existing[key] = azurerm[key];
            }
        } else {
            existing[key] = azurerm[key];
        }
    } else {
        existing[key] = azurerm[key];
    }
});

// Add backup metadata for cleanup
existing['// AZURERM_BACKUP_LENGTH'] = originalLength;

// Write merged settings
fs.writeFileSync(settingsPath, JSON.stringify(existing, null, 2));
console.log('Settings merged successfully using Node.js!');
" 2>/dev/null

    if [[ $? -ne 0 ]]; then
        echo "ERROR: Node.js JSON merge failed"
        echo ""
        echo "This may be due to invalid JSON syntax in your existing settings.json"
        echo ""

        # Create fake backup marker for cleanup detection
        if create_fake_backup_marker "$BACKUP_PATH"; then
            echo "Created recovery marker for cleanup"
        fi

        echo ""
        echo "MANUAL INTERVENTION REQUIRED:"
        echo "Please refer to the Python error handling section above for complete manual setup instructions."
        echo ""
        echo "IMPORTANT: Make a backup copy of your settings.json before editing!"
        echo "  Example: cp '$SETTINGS_PATH' '$SETTINGS_PATH.backup-\$(date +%Y%m%d-%H%M%S)'"
        exit 1
    fi

else
    # Fallback: Manual merge required
    echo "WARNING: No JSON processing tool available (jq, python3, or node)..."
    echo ""

    # Create fake backup marker for cleanup detection
    if create_fake_backup_marker "$BACKUP_PATH"; then
        echo "Created recovery marker for cleanup"
    fi

    echo ""
    echo "MANUAL MERGE REQUIRED:"
    echo "Please manually add these settings to your settings.json:"
    echo ""
    echo "Add the comprehensive AzureRM settings listed in the error recovery section above."
    echo ""
    echo "IMPORTANT: Make a backup copy of your settings.json before editing!"
    echo "  Example: cp '$SETTINGS_PATH' '$SETTINGS_PATH.backup-\$(date +%Y%m%d-%H%M%S)'"
    echo ""

    # Provide manual merge template
    echo "Automatic installation requires JSON processing tools. Install one of:"
    echo "  # Ubuntu/Debian: sudo apt install jq"
    echo "  # macOS: brew install jq"
    echo "  # CentOS/RHEL: sudo yum install jq"
    echo "  # Or install python3 or node.js"
    echo ""
    echo "After manual setup, you can run:"
    echo "  $0 -clean"
    echo "to remove the installation when needed."
    echo ""
    exit 1
fi

# Verify the written file
if test_file_integrity "$SETTINGS_PATH" "github\.copilot"; then
    echo "Settings merged and verified successfully!"
else
    echo "WARNING: Settings written but verification failed"

    # Attempt to restore backup if we created one
    if [[ "$settings_backup_created" == "0" && -f "$BACKUP_PATH" ]]; then
        if cp "$BACKUP_PATH" "$SETTINGS_PATH" 2>/dev/null; then
            echo "Settings restored from backup"
        else
            echo "ERROR: Failed to restore backup!"
        fi
    fi
fi

echo ""
echo "Installation Complete!"
echo "========================="
echo ""
echo "Files installed to:"
echo "  Instructions: $INSTRUCTIONS_DIR"
echo "  AI Prompts: $PROMPTS_DIR"
echo "  VS Code Settings: $SETTINGS_PATH"
if [[ "$settings_backup_created" == "0" ]]; then
    echo "  Settings Backup: $BACKUP_PATH"
fi
echo ""
echo "Next Steps:"
echo "  1. Restart VS Code to load new settings"
echo "  2. Open your terraform-provider-azurerm repository"
echo "  3. Start using AI-powered development!"
echo ""
echo "Available Features:"
echo "  - Context-aware code generation with @workspace"
echo "  - Azure-specific implementation patterns"
echo "  - Automated testing and documentation"
echo "  - Code review and quality enforcement"
echo "  - Error handling and debugging assistance"
echo ""
echo "For help: ./install-copilot-setup.sh -help"
echo "To remove: ./install-copilot-setup.sh -clean"
