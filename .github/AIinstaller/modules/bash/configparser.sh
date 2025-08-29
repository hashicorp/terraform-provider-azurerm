#!/usr/bin/env bash
# ConfigParser Module for Terraform AzureRM Provider AI Setup (Bash)
# Handles configuration parsing, user profiles, and installation settings

# Function to get user profile directory
get_user_profile() {
    echo "${HOME}/.terraform-ai-installer"
}

# Function to get source repository URL
get_source_repository() {
    echo "${SOURCE_REPOSITORY:-https://raw.githubusercontent.com/hashicorp/terraform-provider-azurerm}"
}

# Function to get current branch
get_source_branch() {
    echo "${BRANCH:-exp/terraform_copilot}"
}

# Function to read manifest configuration
read_manifest_config() {
    local manifest_file="${1:-${HOME}/.terraform-ai-installer/file-manifest.config}"
    
    if [[ ! -f "${manifest_file}" ]]; then
        if declare -f write_error >/dev/null 2>&1; then
            write_error "Manifest file not found: ${manifest_file}"
        else
            echo -e "\033[0;31m[ERROR]\033[0m Manifest file not found: ${manifest_file}"
        fi
        echo "Please run with --bootstrap first to set up the installer."
        return 1
    fi
    
    # For now, return success - in future could parse actual manifest
    return 0
}

# Function to parse manifest section (simplified to match PowerShell implementation)
parse_manifest_section() {
    local manifest_file="$1"
    local section_name="$2"
    
    if [[ ! -f "${manifest_file}" ]]; then
        if declare -f write_error >/dev/null 2>&1; then
            write_error "Manifest file not found: ${manifest_file}"
        else
            echo -e "\033[0;31m[ERROR]\033[0m Manifest file not found: ${manifest_file}"
        fi
        return 1
    fi
    
    local in_section=false
    
    while IFS= read -r line; do
        # Remove leading/trailing whitespace
        line=$(echo "${line}" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
        
        # Skip empty lines and comments
        if [[ -z "${line}" || "${line}" =~ ^# ]]; then
            continue
        fi
        
        # Check for section headers [SECTION_NAME]
        if [[ "${line}" =~ ^\[([^]]+)\]$ ]]; then
            local current_section="${BASH_REMATCH[1]}"
            if [[ "${current_section}" == "${section_name}" ]]; then
                in_section=true
            else
                in_section=false
            fi
            continue
        fi
        
        # Output files from the requested section
        if [[ "${in_section}" == "true" && -n "${line}" ]]; then
            echo "${line}"
        fi
    done < "${manifest_file}"
}

# Function to get files from manifest by section
get_manifest_files() {
    local section_name="$1"
    local manifest_file="${2:-${HOME}/.terraform-ai-installer/file-manifest.config}"
    
    # Require manifest file - no fallback
    if [[ ! -f "${manifest_file}" ]]; then
        if declare -f write_error >/dev/null 2>&1; then
            write_error "Manifest file not found: ${manifest_file}"
        else
            echo -e "\033[0;31m[ERROR]\033[0m Manifest file not found: ${manifest_file}"
        fi
        echo "Please run with -bootstrap first to set up the installer."
        return 1
    fi
    
    # Parse from manifest
    parse_manifest_section "${manifest_file}" "${section_name}"
}

# Function to get bootstrap files list
get_bootstrap_files() {
    local -a files=(
        "file-manifest.config"
        "install-copilot-setup.sh"
    )
    
    printf '%s\n' "${files[@]}"
}

# Function to create default configuration
create_default_config() {
    local config_dir="${1:-$(get_user_profile)}"
    local config_file="${config_dir}/installer.config"
    
    if [[ "${DRY_RUN:-false}" == "true" ]]; then
        echo "  [DRY-RUN] Would create default configuration: ${config_file}"
        return 0
    fi
    
    # Create config directory if it doesn't exist
    mkdir -p "${config_dir}"
    
    # Create basic configuration file
    cat > "${config_file}" << 'EOF'
# Terraform AzureRM Provider AI Installer Configuration
# Generated automatically - modify as needed

# Installation settings
DEFAULT_BRANCH=exp/terraform_copilot
SOURCE_REPOSITORY=https://raw.githubusercontent.com/hashicorp/terraform-provider-azurerm
DOWNLOAD_TIMEOUT=30
RETRY_COUNT=3

# Directory settings
BACKUP_SUFFIX=.backup
CREATE_BACKUPS=true
VERIFY_DOWNLOADS=true

# Logging settings
VERBOSE_MODE=false
DEBUG_MODE=false
LOG_FILE_OPERATIONS=true
EOF
    
    if [[ -f "${config_file}" ]]; then
        echo "  Created default configuration: ${config_file} [OK]"
        return 0
    else
        echo "  Created default configuration: ${config_file} [FAILED]"
        return 1
    fi
}

# Function to load configuration from file
load_config() {
    local config_file="${1:-$(get_user_profile)/installer.config}"
    
    if [[ -f "${config_file}" ]]; then
        # Source the configuration file to load variables
        # Note: This is safe because we control the config file content
        if source "${config_file}" 2>/dev/null; then
            if declare -f write_plain >/dev/null 2>&1; then
                write_plain "Loaded configuration from: ${config_file}"
            fi
            return 0
        else
            if declare -f write_warning >/dev/null 2>&1; then
                write_warning "Failed to load configuration from: ${config_file}"
            else
                echo -e "\033[1;33m[WARNING]\033[0m Failed to load configuration from: ${config_file}"
            fi
            return 1
        fi
    fi
    
    return 0
}

# Function to validate configuration
validate_config() {
    local errors=0
    
    # Check required variables
    if [[ -z "${SOURCE_REPOSITORY:-}" ]]; then
        if declare -f write_error >/dev/null 2>&1; then
            write_error "SOURCE_REPOSITORY not configured"
        else
            echo -e "\033[0;31m[ERROR]\033[0m SOURCE_REPOSITORY not configured"
        fi
        errors=$((errors + 1))
    fi
    
    if [[ -z "${BRANCH:-}" ]]; then
        if declare -f write_warning >/dev/null 2>&1; then
            write_warning "BRANCH not configured, using default"
        else
            echo -e "\033[1;33m[WARNING]\033[0m BRANCH not configured, using default"
        fi
        BRANCH="exp/terraform_copilot"
    fi
    
    return ${errors}
}

# Function to get file download URL
get_file_download_url() {
    local file_path="$1"
    local repository="${SOURCE_REPOSITORY:-https://raw.githubusercontent.com/hashicorp/terraform-provider-azurerm}"
    local branch="${BRANCH:-exp/terraform_copilot}"
    
    echo "${repository}/${branch}/${file_path}"
}

# Function to convert to relative path (helper function)
convert_to_relative_path() {
    local full_path="$1"
    local base_path="${2:-$(pwd)}"
    
    # Simple relative path conversion
    if [[ "${full_path}" == "${base_path}"* ]]; then
        echo "${full_path#${base_path}/}"
    else
        echo "${full_path}"
    fi
}

# Function to get installer version
get_installer_version() {
    echo "${VERSION:-1.0.0}"
}

# Function to get temp directory for downloads
get_temp_directory() {
    local temp_base="${TMPDIR:-/tmp}"
    local temp_dir="${temp_base}/terraform-ai-installer-$$"
    
    mkdir -p "${temp_dir}"
    echo "${temp_dir}"
}

# Function to get all user-facing files for clean operations
get_files_for_cleanup() {
    local workspace_root="${1:-$(pwd)}"
    local manifest_file="${HOME}/.terraform-ai-installer/file-manifest.config"
    
    # The manifest file should only exist in the user's home directory
    # The AIinstaller directory is not copied to feature branches
    if [[ ! -f "${manifest_file}" ]]; then
        write_error_message "Manifest file not found: ${manifest_file}"
        write_error_message "Please run with -bootstrap first to set up the installer."
        return 1
    fi
    
    # Get files from all user-facing sections
    local sections=("MAIN_FILES" "INSTRUCTION_FILES" "PROMPT_FILES" "UNIVERSAL_FILES")
    local all_files=()
    
    for section in "${sections[@]}"; do
        local section_files
        section_files=($(parse_manifest_section "${manifest_file}" "${section}"))
        all_files+=("${section_files[@]}")
    done
    
    printf '%s\n' "${all_files[@]}"
}

# Function to cleanup temp directory
cleanup_temp_directory() {
    local temp_dir="$1"
    
    if [[ -n "${temp_dir}" ]] && [[ "${temp_dir}" =~ /tmp/terraform-ai-installer- ]]; then
        rm -rf "${temp_dir}" 2>/dev/null || true
    fi
}

# Function to get manifest configuration (wrapper around read_manifest_config)
get_manifest_config() {
    local branch="${1:-exp/terraform_copilot}"
    local manifest_path="$2"
    
    # Find manifest file if not specified
    if [[ -z "${manifest_path}" ]]; then
        local possible_paths=(
            "${HOME}/.terraform-ai-installer/file-manifest.config"
            "$(dirname "${BASH_SOURCE[0]}")/../file-manifest.config"
            "$(dirname "$(dirname "${BASH_SOURCE[0]}")")/file-manifest.config"
        )
        
        for path in "${possible_paths[@]}"; do
            if [[ -f "${path}" ]]; then
                manifest_path="${path}"
                break
            fi
        done
    fi
    
    if [[ -z "${manifest_path}" ]] || [[ ! -f "${manifest_path}" ]]; then
        if declare -f write_error >/dev/null 2>&1; then
            write_error "Manifest file not found. Run with --bootstrap first."
        else
            echo -e "\033[0;31m[ERROR]\033[0m Manifest file not found. Run with --bootstrap first."
        fi
        return 1
    fi
    
    # Return manifest configuration data
    echo "ManifestPath=${manifest_path}"
    echo "Branch=${branch}"
    echo "Found=true"
    return 0
}

# Function to get installer configuration
get_installer_config() {
    local workspace_root="$1"
    local manifest_output="$2"
    
    if [[ -z "${workspace_root}" ]]; then
        if declare -f write_error >/dev/null 2>&1; then
            write_error "Workspace root is required"
        else
            echo -e "\033[0;31m[ERROR]\033[0m Workspace root is required"
        fi
        return 1
    fi
    
    # Parse manifest output
    local branch="exp/terraform_copilot"
    if [[ -n "${manifest_output}" ]]; then
        branch=$(echo "${manifest_output}" | grep "Branch=" | cut -d'=' -f2 || echo "exp/terraform_copilot")
    fi
    
    # Return installer configuration
    echo "Version=1.0.0"
    echo "Branch=${branch}"
    echo "WorkspaceRoot=${workspace_root}"
    echo "Repository=hashicorp/terraform-provider-azurerm"
    return 0
}

# Export functions for use in other scripts
export -f get_user_profile get_source_repository get_source_branch read_manifest_config
export -f get_manifest_files get_bootstrap_files get_files_for_cleanup create_default_config load_config
export -f validate_config get_file_download_url convert_to_relative_path get_installer_version
export -f get_temp_directory cleanup_temp_directory parse_manifest_section get_manifest_config get_installer_config
