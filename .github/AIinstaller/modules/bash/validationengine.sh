#!/usr/bin/env bash
# ValidationEngine Module for Terraform AzureRM Provider AI Setup (Bash)
# Handles repository validation, workspace detection, and system requirements

# Function to get workspace root
get_workspace_root() {
    local current_path
    current_path="$(dirname "$(realpath "${0}")")"
    
    while [[ "${current_path}" != "/" ]]; do
        # Look for terraform-provider-azurerm indicators
        if [[ -f "${current_path}/go.mod" ]] && 
           [[ -f "${current_path}/main.go" ]] && 
           [[ -d "${current_path}/internal" ]]; then
            echo "${current_path}"
            return 0
        fi
        
        # Move up one directory
        current_path="$(dirname "${current_path}")"
    done
    
    # Fallback: assume we're in .github/AIinstaller and go up two levels
    echo "$(dirname "$(dirname "$(dirname "$(realpath "${0}")")")")"
}

# Function to validate repository directory
validate_repository() {
    local repo_dir="$1"
    
    if [[ ! -d "${repo_dir}" ]]; then
        if declare -f write_error >/dev/null 2>&1; then
            write_error "DIRECTORY NOT FOUND: The specified RepoDirectory does not exist."
        else
            echo -e "\033[0;31m[ERROR]\033[0m DIRECTORY NOT FOUND: The specified RepoDirectory does not exist."
        fi
        echo "The path '${repo_dir}' could not be found on this system."
        echo "Solutions:"
        echo "  - Check the path spelling and ensure it exists"
        echo "  - Use an absolute path (e.g., '/Users/username/terraform-provider-azurerm')"
        echo "  - Ensure you have permissions to access the directory"
        exit 1
    fi
    
    # Validate that this looks like a terraform-provider-azurerm repository
    if [[ ! -f "${repo_dir}/go.mod" ]] || 
       [[ ! -f "${repo_dir}/main.go" ]] || 
       [[ ! -d "${repo_dir}/internal" ]]; then
        if declare -f write_error >/dev/null 2>&1; then
            write_error "INVALID REPOSITORY: The specified directory does not appear to be a terraform-provider-azurerm repository."
        else
            echo -e "\033[0;31m[ERROR]\033[0m INVALID REPOSITORY: The specified directory does not appear to be a terraform-provider-azurerm repository."
        fi
        echo "The --repo-directory parameter must point to a valid terraform-provider-azurerm repository root."
        echo "Solutions:"
        echo "  - Ensure you're pointing to the repository ROOT directory"
        echo "  - Verify the directory contains terraform-provider-azurerm source code"
        echo "  - Example: --repo-directory '/Users/username/terraform-provider-azurerm'"
        exit 1
    fi
    
    if declare -f write_info >/dev/null 2>&1; then
        write_info "Repository validated: ${repo_dir}"
    else
        echo -e "\033[0;34m[INFO]\033[0m Repository validated: ${repo_dir}"
    fi
}

# Function to test system requirements
test_system_requirements() {
    local missing_tools=()
    
    # Check for curl or wget
    if ! command -v curl >/dev/null 2>&1 && ! command -v wget >/dev/null 2>&1; then
        missing_tools+=("curl or wget")
    fi
    
    # Check for basic Unix tools
    local required_tools=("bash" "mkdir" "cp" "rm" "dirname" "realpath")
    for tool in "${required_tools[@]}"; do
        if ! command -v "${tool}" >/dev/null 2>&1; then
            missing_tools+=("${tool}")
        fi
    done
    
    if [[ ${#missing_tools[@]} -gt 0 ]]; then
        if declare -f write_error >/dev/null 2>&1; then
            write_error "Missing required system tools: ${missing_tools[*]}"
        else
            echo -e "\033[0;31m[ERROR]\033[0m Missing required system tools: ${missing_tools[*]}"
        fi
        return 1
    fi
    
    return 0
}

# Function to test Git repository
test_git_repository() {
    local repo_dir="$1"
    
    if [[ ! -d "${repo_dir}/.git" ]]; then
        if declare -f write_warning >/dev/null 2>&1; then
            write_warning "Not a Git repository: ${repo_dir}"
        else
            echo -e "\033[1;33m[WARNING]\033[0m Not a Git repository: ${repo_dir}"
        fi
        return 1
    fi
    
    # Check if git command is available
    if ! command -v git >/dev/null 2>&1; then
        if declare -f write_warning >/dev/null 2>&1; then
            write_warning "Git command not available"
        else
            echo -e "\033[1;33m[WARNING]\033[0m Git command not available"
        fi
        return 1
    fi
    
    return 0
}

# Function to get current Git branch
get_current_branch() {
    local repo_dir="$1"
    
    if test_git_repository "${repo_dir}"; then
        local branch
        branch=$(cd "${repo_dir}" && git branch --show-current 2>/dev/null || echo "unknown")
        echo "${branch}"
    else
        echo "unknown"
    fi
}

# Function to test workspace validity
test_workspace_valid() {
    local workspace_root="$1"
    
    # Test basic repository structure
    if ! validate_repository "${workspace_root}" 2>/dev/null; then
        return 1
    fi
    
    # Test Git repository if git is available
    if command -v git >/dev/null 2>&1; then
        test_git_repository "${workspace_root}"
    fi
    
    return 0
}

# Function to test internet connectivity
test_internet_connectivity() {
    local test_url="https://raw.githubusercontent.com"
    
    if command -v curl >/dev/null 2>&1; then
        if curl -fsSL --connect-timeout 10 "${test_url}" >/dev/null 2>&1; then
            return 0
        fi
    elif command -v wget >/dev/null 2>&1; then
        if wget -q --timeout=10 --tries=1 "${test_url}" -O /dev/null 2>/dev/null; then
            return 0
        fi
    fi
    
    return 1
}

# Export functions for use in other scripts
export -f get_workspace_root validate_repository test_system_requirements
export -f test_git_repository get_current_branch test_workspace_valid test_internet_connectivity
