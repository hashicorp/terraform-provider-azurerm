#!/usr/bin/env bash
# ValidationEngine Module for Terraform AzureRM Provider AI Setup (Bash)
# Handles repository validation, workspace detection, and system requirements

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
    
    if declare -f write_plain >/dev/null 2>&1; then
        write_plain "Repository validated: ${repo_dir}"
    else
        echo "Repository validated: ${repo_dir}"
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

# Function to verify AI infrastructure installation
verify_installation() {
    local workspace_root="${1:-$(get_workspace_root)}"
    
    if declare -f write_section >/dev/null 2>&1; then
        write_section "Verifying AI infrastructure files"
    else
        echo "============================================================"
        echo " Verifying AI infrastructure files"
        echo "============================================================"
        echo ""
    fi
    
    local all_good=true
    local manifest_file="${HOME}/.terraform-ai-installer/file-manifest.config"
    
    # Check main files
    local main_files
    main_files=$(get_manifest_files "MAIN_FILES" "${manifest_file}")
    if [[ $? -eq 0 && -n "${main_files}" ]]; then
        while IFS= read -r file; do
            [[ -z "${file}" ]] && continue
            local full_path="${workspace_root}/${file}"
            if [[ -f "${full_path}" ]]; then
                printf "\033[0;32m  [FOUND] %s\033[0m\n" "${file}"
            else
                printf "\033[0;31m  [MISSING] %s\033[0m\n" "${file}"
                all_good=false
            fi
        done <<< "${main_files}"
    fi
    
    # Check instruction files
    local instruction_files
    instruction_files=$(get_manifest_files "INSTRUCTION_FILES" "${manifest_file}")
    if [[ $? -eq 0 && -n "${instruction_files}" ]]; then
        # Check if instructions directory exists
        local instructions_dir="${workspace_root}/.github/instructions"
        if [[ -d "${instructions_dir}" ]]; then
            printf "\033[0;32m  [FOUND] .github/instructions/\033[0m\n"
            
            while IFS= read -r file; do
                [[ -z "${file}" ]] && continue
                local full_path="${workspace_root}/${file}"
                local filename=$(basename "${file}")
                if [[ -f "${full_path}" ]]; then
                    printf "\033[0;32m    [FOUND] .github/instructions/%s\033[0m\n" "${filename}"
                else
                    printf "\033[0;31m    [MISSING] .github/instructions/%s\033[0m\n" "${filename}"
                    all_good=false
                fi
            done <<< "${instruction_files}"
        else
            printf "\033[0;31m  [MISSING] .github/instructions/\033[0m\n"
            all_good=false
        fi
    fi
    
    # Check prompt files
    local prompt_files
    prompt_files=$(get_manifest_files "PROMPT_FILES" "${manifest_file}")
    if [[ $? -eq 0 && -n "${prompt_files}" ]]; then
        # Check if prompts directory exists
        local prompts_dir="${workspace_root}/.github/prompts"
        if [[ -d "${prompts_dir}" ]]; then
            printf "\033[0;32m  [FOUND] .github/prompts/\033[0m\n"
            
            while IFS= read -r file; do
                [[ -z "${file}" ]] && continue
                local full_path="${workspace_root}/${file}"
                local filename=$(basename "${file}")
                if [[ -f "${full_path}" ]]; then
                    printf "\033[0;32m    [FOUND] .github/prompts/%s\033[0m\n" "${filename}"
                else
                    printf "\033[0;31m    [MISSING] .github/prompts/%s\033[0m\n" "${filename}"
                    all_good=false
                fi
            done <<< "${prompt_files}"
        else
            printf "\033[0;31m  [MISSING] .github/prompts/\033[0m\n"
            all_good=false
        fi
    fi
    
    # Check universal files
    local universal_files
    universal_files=$(get_manifest_files "UNIVERSAL_FILES" "${manifest_file}")
    if [[ $? -eq 0 && -n "${universal_files}" ]]; then
        while IFS= read -r file; do
            [[ -z "${file}" ]] && continue
            local full_path="${workspace_root}/${file}"
            local dir_path=$(dirname "${file}")
            local filename=$(basename "${file}")
            
            if [[ -f "${full_path}" ]]; then
                printf "\033[0;32m  [FOUND] %s\033[0m\n" "${dir_path}/"
                printf "\033[0;32m    [FOUND] %s\033[0m\n" "${filename}"
            else
                printf "\033[0;31m  [MISSING] %s\033[0m\n" "${dir_path}/"
                printf "\033[0;31m    [MISSING] %s\033[0m\n" "${filename}"
                all_good=false
            fi
        done <<< "${universal_files}"
    fi
    
    echo ""
    if [[ "${all_good}" == "true" ]]; then
        printf "\033[0;32mAll AI infrastructure files are present in the source repository!\033[0m\n"
    else
        printf "\033[0;31mSome AI infrastructure files are missing!\033[0m\n"
        echo ""
        echo "Run the installer to restore missing files."
    fi
    echo ""
}

# Export functions for use in other scripts
export -f validate_repository test_system_requirements
export -f test_git_repository get_current_branch test_workspace_valid test_internet_connectivity verify_installation
