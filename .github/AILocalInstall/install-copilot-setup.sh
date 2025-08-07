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
#   --repository-path <path>  Path to terraform-provider-azurerm repository
#   --clean                   Remove all AI setup and restore backups
#   --help                    Show detailed help information
#
# EXAMPLES:
#   ./install-copilot-setup.sh
#   ./install-copilot-setup.sh --repository-path "/home/user/terraform-provider-azurerm"
#   ./install-copilot-setup.sh --clean
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

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --repository-path|--repo-path|-r)
            REPOSITORY_PATH="$2"
            shift # past argument
            shift # past value
            ;;
        --clean|-c)
            CLEAN_MODE="true"
            shift # past argument
            ;;
        --help|-h)
            echo "Terraform AzureRM Provider AI Setup"
            echo ""
            echo "USAGE:"
            echo "  ./install-copilot-setup.sh [OPTIONS]"
            echo ""
            echo "OPTIONS:"
            echo "  --repository-path <path>  Path to terraform-provider-azurerm repository"
            echo "  --clean                   Remove all installed files and restore backups"
            echo "  --help                    Show this help message"
            echo ""
            echo "EXAMPLES:"
            echo "  ./install-copilot-setup.sh                                    # Auto-discover repository"
            echo "  ./install-copilot-setup.sh --repository-path /path/to/repo    # Use specific path"
            echo "  ./install-copilot-setup.sh --clean                            # Remove installation"
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
            echo "Unknown option $1"
            exit 1
            ;;
        *)
            echo "Error: Unexpected argument '$1'"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

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
    local previous_install_result=$(test_previous_installation "$USER_DIR")
    local has_previous=$(echo "$previous_install_result" | head -n 1)
    
    if [[ "$has_previous" != "true" ]]; then
        echo "No previous installation detected. Nothing to clean."
        exit 0
    fi
    
    echo "Found previous installation files:"
    echo "$previous_install_result" | grep "^FOUND:" | sed 's/^FOUND:/  - /'
    echo ""
    
    # Enhanced backup detection and handling
    local backup_info
    backup_info=$(get_backup_info "$BACKUP_PATH")
    local backup_status="${backup_info%%:*}"
    local backup_length="${backup_info##*:}"
    
    if [[ "$backup_status" == "NO_BACKUP" ]]; then
        # No backup found - check if settings contain AzureRM settings
        local has_azurerm_settings
        has_azurerm_settings=$(detect_azurerm_settings "$SETTINGS_PATH")
        
        if [[ "$has_azurerm_settings" == "true" ]]; then
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
        else
            echo "No backup found, but no AzureRM Copilot settings detected in settings.json"
            echo "Nothing to clean regarding VS Code settings."
        fi
    elif [[ "$backup_status" == "MARKER" ]]; then
        if [[ "$backup_length" == "-1" ]]; then
            # Fake backup from failed installation
            echo "Found recovery marker from failed installation"
            echo ""
            echo "MANUAL CLEANUP REQUIRED:"
            echo "The previous installation failed, but may have left partial settings."
            echo "Please manually remove these entries from $SETTINGS_PATH:"
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
            echo ""
            echo "After manual cleanup, remove the recovery marker:"
            rm -f "$BACKUP_PATH"
            echo "Recovery marker removed."
            exit 0
        elif [[ "$backup_length" == "0" ]]; then
            # Original file was empty/new
            echo "Restoring original state (no settings.json originally existed)..."
            if rm -f "$SETTINGS_PATH" 2>/dev/null; then
                echo "Removed settings.json (originally was empty)"
            else
                echo "WARNING: Could not remove settings.json"
            fi
            rm -f "$BACKUP_PATH"
            echo "Backup cleaned up."
        else
            # Normal backup with length marker
            echo "Restoring original settings.json from backup..."
            if [[ -f "$BACKUP_PATH" && -f "$SETTINGS_PATH" ]]; then
                if test_file_integrity "$BACKUP_PATH"; then
                    if cp "$BACKUP_PATH" "$SETTINGS_PATH" 2>/dev/null; then
                        echo "Original settings restored successfully!"
                        rm -f "$BACKUP_PATH"
                        echo "Backup cleaned up."
                    else
                        echo "ERROR: Could not restore settings from backup"
                        exit 1
                    fi
                else
                    echo "ERROR: Backup file appears corrupted"
                    echo "Manual restoration required - backup location: $BACKUP_PATH"
                    exit 1
                fi
            else
                echo "ERROR: Backup or current settings file missing"
                exit 1
            fi
        fi
    else
        # SIZE backup (legacy)
        if [[ "$backup_length" -gt 0 ]]; then
            echo "Restoring original settings.json from backup..."
            if [[ -f "$BACKUP_PATH" && -f "$SETTINGS_PATH" ]]; then
                if test_file_integrity "$BACKUP_PATH"; then
                    if cp "$BACKUP_PATH" "$SETTINGS_PATH" 2>/dev/null; then
                        echo "Original settings restored successfully!"
                        rm -f "$BACKUP_PATH"
                        echo "Backup cleaned up."
                    else
                        echo "ERROR: Could not restore settings from backup"
                        exit 1
                    fi
                else
                    echo "ERROR: Backup file appears corrupted"
                    echo "Manual restoration required - backup location: $BACKUP_PATH"
                    exit 1
                fi
            else
                echo "ERROR: Backup or current settings file missing"
                exit 1
            fi
        else
            echo "Backup indicates original file was empty"
            if rm -f "$SETTINGS_PATH" 2>/dev/null; then
                echo "Removed settings.json (originally was empty)"
            else
                echo "WARNING: Could not remove settings.json"
            fi
            rm -f "$BACKUP_PATH"
            echo "Backup cleaned up."
        fi
    fi
    
    # Interactive confirmation for remaining files (if needed)
    if [[ -f "$SETTINGS_PATH" && ! -f "$BACKUP_PATH" ]]; then
        local has_azurerm_settings
        has_azurerm_settings=$(detect_azurerm_settings "$SETTINGS_PATH")
        
        if [[ "$has_azurerm_settings" == "true" ]]; then
            echo ""
            echo "WARNING: settings.json still contains AzureRM settings after restoration!"
            echo "This may indicate a backup/restore issue."
            echo ""
            while true; do
                echo -n "Continue cleanup of instruction files? (C)ontinue/(X)exit: "
                read -r response
                response=$(echo "$response" | tr '[:lower:]' '[:upper:]')
                case "$response" in
                    C) break ;;
                    X) echo "Cleanup cancelled by user."; exit 0 ;;
                    *) echo "Invalid choice. Please enter C or X." ;;
                esac
            done
        fi
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
    local copilot_instructions_path="$USER_DIR/copilot-instructions.md"
    if [[ -f "$copilot_instructions_path" ]]; then
        echo "Removing main copilot instructions..."
        if rm -f "$copilot_instructions_path" 2>/dev/null; then
            echo "Copilot instructions removed"
        else
            echo "ERROR: Failed to remove copilot instructions"
        fi
    fi
    
    echo ""
    echo "Cleanup completed successfully!"
    echo "Restart VS Code to ensure all changes take effect."
    exit 0
fi

# Previous installation detection to prevent backup corruption
local previous_install_result=$(test_previous_installation "$USER_DIR")
local has_previous=$(echo "$previous_install_result" | head -n 1)

if [[ "$has_previous" == "true" ]]; then
    echo ""
    echo "Previous installation detected!"
    echo "Found existing files:"
    echo "$previous_install_result" | grep "^FOUND:" | sed 's/^FOUND:/  - /'
    echo ""
    echo "This will overwrite existing AzureRM AI configuration."
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
local instructions_src_dir="$SOURCE_DIR/instructions"
if [[ -d "$instructions_src_dir" ]]; then
    local instruction_count=$(find "$instructions_src_dir" -name "*.md" 2>/dev/null | wc -l)
    
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
local prompts_src_dir="$SOURCE_DIR/prompts"
if [[ -d "$prompts_src_dir" ]]; then
    local prompt_count=$(find "$prompts_src_dir" -name "*.md" 2>/dev/null | wc -l)
    
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
local copilot_src_path="$SOURCE_DIR/copilot-instructions.md"
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

# Smart merge VS Code settings.json with enhanced backup and recovery
echo "Configuring VS Code settings (preserving existing settings)..."

local existing_settings_valid=true
local settings_backup_created=false

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
                echo "  1. Run the clean command first: ./install-copilot-setup.sh --clean"
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
                echo "  2. Run: ./install-copilot-setup.sh --clean (will warn about missing backup)"
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
                echo "  $0 --clean"
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
local original_length="0"
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
    echo "  $0 --clean"
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
echo "For help: ./install-copilot-setup.sh --help"
echo "To remove: ./install-copilot-setup.sh --clean"
