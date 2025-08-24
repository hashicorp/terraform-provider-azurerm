#!/usr/bin/env bash
# FileOperations Module for Terraform AzureRM Provider AI Setup (Bash)
# Handles file operations, downloads, and installation tasks

# Function to get workspace root directory
get_workspace_root() {
    local current_dir="$(pwd)"
    
    # Look for terraform-provider-azurerm indicators
    while [[ "${current_dir}" != "/" ]]; do
        if [[ -f "${current_dir}/go.mod" ]] && grep -q "terraform-provider-azurerm" "${current_dir}/go.mod" 2>/dev/null; then
            echo "${current_dir}"
            return 0
        fi
        if [[ -f "${current_dir}/main.go" ]] && [[ -d "${current_dir}/internal/services" ]]; then
            echo "${current_dir}"
            return 0
        fi
        current_dir="$(dirname "${current_dir}")"
    done
    
    # Fallback to current directory
    echo "$(pwd)"
}

# Function to check if current directory is source repository
is_source_repository() {
    local dir="${1:-$(pwd)}"
    
    # Check if this is the terraform-provider-azurerm source repository
    if [[ -f "${dir}/go.mod" ]] && grep -q "terraform-provider-azurerm" "${dir}/go.mod" 2>/dev/null; then
        return 0
    fi
    
    if [[ -f "${dir}/main.go" ]] && [[ -d "${dir}/internal/services" ]] && [[ -d "${dir}/.github/AIinstaller" ]]; then
        return 0
    fi
    
    return 1
}

# Function to validate repository directory
validate_repository() {
    local repo_dir="$1"
    
    if [[ ! -d "${repo_dir}" ]]; then
        write_error "Repository directory does not exist: ${repo_dir}"
        return 1
    fi
    
    # Critical safety check: prevent running in source repository
    if is_source_repository "${repo_dir}"; then
        write_error "SAFETY CHECK FAILED: Cannot install to source repository directory"
        echo ""
        echo "This appears to be the terraform-provider-azurerm source repository."
        echo "Installing here would overwrite your local changes with remote files."
        echo ""
        echo "To install to a different repository:"
        echo "  ${0} /path/to/target/repository"
        echo ""
        echo "To bootstrap the installer to your user profile:"
        echo "  ${0} -bootstrap"
        echo ""
        return 1
    fi
    
    # Validate this is a terraform-provider-azurerm repository
    if [[ ! -f "${repo_dir}/go.mod" ]] || ! grep -q "terraform-provider-azurerm" "${repo_dir}/go.mod" 2>/dev/null; then
        write_warning "Target directory may not be a terraform-provider-azurerm repository"
        echo "Expected to find go.mod with terraform-provider-azurerm reference"
        echo "Directory: ${repo_dir}"
        echo ""
        echo "Continue anyway? (y/N)"
        read -r response
        if [[ ! "${response}" =~ ^[Yy]$ ]]; then
            echo "Installation cancelled by user"
            return 1
        fi
    fi
    
    return 0
}

# Function to copy file with progress
copy_file() {
    local source="$1"
    local target="$2"
    local description="$3"
    
    if [[ "${DRY_RUN:-false}" == "true" ]]; then
        echo "  [DRY-RUN] Would copy: ${description}"
        return 0
    fi
    
    # Create target directory if it doesn't exist
    mkdir -p "$(dirname "${target}")"
    
    if cp "${source}" "${target}"; then
        show_file_operation "Copying" "${description}" "OK"
        return 0
    else
        show_file_operation "Copying" "${description}" "FAILED"
        return 1
    fi
}

# Function to download file from repository
download_file() {
    local source_path="$1"
    local target_path="$2"
    local description="$3"
    
    local url="${SOURCE_REPOSITORY:-https://raw.githubusercontent.com/hashicorp/terraform-provider-azurerm}/${BRANCH:-exp/terraform_copilot}/${source_path}"
    
    if [[ "${DRY_RUN:-false}" == "true" ]]; then
        echo "  [DRY-RUN] Would download: ${description}"
        return 0
    fi
    
    # Create target directory if it doesn't exist
    mkdir -p "$(dirname "${target_path}")"
    
    if command -v curl >/dev/null 2>&1; then
        if curl -fsSL "${url}" -o "${target_path}"; then
            echo "  Downloading: ${description} [OK]"
            return 0
        fi
    elif command -v wget >/dev/null 2>&1; then
        if wget -q "${url}" -O "${target_path}"; then
            echo "  Downloading: ${description} [OK]"
            return 0
        fi
    else
        if declare -f write_error >/dev/null 2>&1; then
            write_error "Neither curl nor wget is available for downloading files"
        else
            echo -e "\033[0;31m[ERROR]\033[0m Neither curl nor wget is available for downloading files"
        fi
        return 1
    fi
    
    echo "  Downloading: ${description} [FAILED]"
    return 1
}

# Function to get file size
get_file_size() {
    local filepath="$1"
    
    if [[ -f "${filepath}" ]]; then
        local file_size
        if command -v stat >/dev/null 2>&1; then
            # Try GNU stat first, then BSD stat
            file_size=$(stat -c%s "${filepath}" 2>/dev/null || stat -f%z "${filepath}" 2>/dev/null || echo "0")
        else
            # Fallback to wc if stat is not available
            file_size=$(wc -c < "${filepath}" 2>/dev/null || echo "0")
        fi
        echo "${file_size}"
    else
        echo "0"
    fi
}

# Function to calculate total directory size
get_directory_size() {
    local dirpath="$1"
    
    if [[ -d "${dirpath}" ]]; then
        local total_size=0
        while IFS= read -r -d '' file; do
            local file_size
            file_size=$(get_file_size "${file}")
            total_size=$((total_size + file_size))
        done < <(find "${dirpath}" -type f -print0 2>/dev/null)
        echo "${total_size}"
    else
        echo "0"
    fi
}

# Function to create directory structure
create_directory_structure() {
    local target_dir="$1"
    local description="${2:-directory structure}"
    
    if [[ "${DRY_RUN:-false}" == "true" ]]; then
        echo "  [DRY-RUN] Would create: ${description}"
        return 0
    fi
    
    if mkdir -p "${target_dir}"; then
        echo "  Creating: ${description} [OK]"
        return 0
    else
        echo "  Creating: ${description} [FAILED]"
        return 1
    fi
}

# Function to remove file or directory
remove_path() {
    local target_path="$1"
    local description="${2:-$(basename "${target_path}")}"
    
    if [[ "${DRY_RUN:-false}" == "true" ]]; then
        echo "  [DRY-RUN] Would remove: ${description}"
        return 0
    fi
    
    if [[ -e "${target_path}" ]]; then
        if rm -rf "${target_path}"; then
            echo "  Removed: ${description} [OK]"
            return 0
        else
            echo "  Removed: ${description} [FAILED]"
            return 1
        fi
    else
        echo "  Remove: ${description} [NOT FOUND]"
        return 0
    fi
}

# Function to backup existing file
backup_file() {
    local filepath="$1"
    local backup_suffix="${2:-.backup.$(date +%Y%m%d_%H%M%S)}"
    
    if [[ -f "${filepath}" ]]; then
        local backup_path="${filepath}${backup_suffix}"
        
        if [[ "${DRY_RUN:-false}" == "true" ]]; then
            echo "  [DRY-RUN] Would backup: $(basename "${filepath}") -> $(basename "${backup_path}")"
            return 0
        fi
        
        if cp "${filepath}" "${backup_path}"; then
            echo "  Backup: $(basename "${filepath}") -> $(basename "${backup_path}") [OK]"
            return 0
        else
            echo "  Backup: $(basename "${filepath}") [FAILED]"
            return 1
        fi
    fi
    
    return 0
}

# Function to verify file exists and is readable
verify_file() {
    local filepath="$1"
    local description="${2:-$(basename "${filepath}")}"
    
    if [[ -f "${filepath}" ]] && [[ -r "${filepath}" ]]; then
        echo "  [OK] ${description}"
        return 0
    elif [[ -e "${filepath}" ]]; then
        echo "  [UNREADABLE] ${description}"
        return 1
    else
        echo "  [MISSING] ${description}"
        return 1
    fi
}

# Function to make file executable
make_executable() {
    local filepath="$1"
    local description="${2:-$(basename "${filepath}")}"
    
    if [[ "${DRY_RUN:-false}" == "true" ]]; then
        echo "  [DRY-RUN] Would make executable: ${description}"
        return 0
    fi
    
    if [[ -f "${filepath}" ]]; then
        if chmod +x "${filepath}"; then
            echo "  Make executable: ${description} [OK]"
            return 0
        else
            echo "  Make executable: ${description} [FAILED]"
            return 1
        fi
    else
        echo "  Make executable: ${description} [FILE NOT FOUND]"
        return 1
    fi
}

# Function to copy multiple files with error tracking
copy_files_with_stats() {
    local -n files_array=$1
    local target_dir="$2"
    local source_dir="${3:-$(pwd)}"
    
    local files_copied=0
    local files_failed=0
    local total_size=0
    
    for file in "${files_array[@]}"; do
        local source_path="${source_dir}/${file}"
        local target_path="${target_dir}/${file}"
        local filename=$(basename "${file}")
        
        echo -n ""  # No inline display needed - show_file_operation handles it
        
        if [[ "${DRY_RUN:-false}" == "true" ]]; then
            show_file_operation "Copying" "${filename}" "DRY-RUN"
            files_copied=$((files_copied + 1))
        elif [[ -f "${source_path}" ]]; then
            # Create target directory structure
            mkdir -p "$(dirname "${target_path}")"
            
            if cp "${source_path}" "${target_path}"; then
                # Calculate file size
                local file_size
                file_size=$(get_file_size "${target_path}")
                total_size=$((total_size + file_size))
                
                show_file_operation "Copying" "${filename}" "OK"
                files_copied=$((files_copied + 1))
            else
                show_file_operation "Copying" "${filename}" "FAILED"
                files_failed=$((files_failed + 1))
            fi
        else
            show_file_operation "Copying" "${filename}" "SOURCE NOT FOUND"
            files_failed=$((files_failed + 1))
        fi
    done
    
    # Return statistics via global variables (bash doesn't support complex return values)
    COPY_STATS_FILES_COPIED=${files_copied}
    COPY_STATS_FILES_FAILED=${files_failed}
    COPY_STATS_TOTAL_SIZE=${total_size}
    
    return $([[ ${files_failed} -eq 0 ]] && echo 0 || echo 1)
}

# Export functions for use in other scripts
export -f get_workspace_root is_source_repository validate_repository
export -f copy_file download_file get_file_size get_directory_size create_directory_structure
export -f remove_path backup_file verify_file make_executable copy_files_with_stats

# Function to install AI infrastructure (moved from main script)
install_infrastructure() {
    local workspace_root="$1"
    
    if declare -f write_section >/dev/null 2>&1; then
        write_section "Installing AI Infrastructure"
    else
        echo "=== Installing AI Infrastructure ==="
    fi
    
    # Verify manifest file exists
    local manifest_file="${HOME}/.terraform-ai-installer/file-manifest.config"
    if [[ ! -f "${manifest_file}" ]]; then
        if declare -f write_error >/dev/null 2>&1; then
            write_error "Manifest file not found: ${manifest_file}"
        else
            echo -e "\033[0;31m[ERROR]\033[0m Manifest file not found: ${manifest_file}"
        fi
        echo "Please run with --bootstrap first to set up the installer."
        return 1
    fi
    
    if declare -f write_plain >/dev/null 2>&1; then
        write_plain "Installing to workspace: ${workspace_root}"
    else
        echo "Installing to workspace: ${workspace_root}"
    fi
    
    # Get file lists from manifest
    local main_files instruction_files prompt_files universal_files
    
    # Read files from manifest sections
    readarray -t main_files < <(get_manifest_files "MAIN_FILES" "${manifest_file}")
    readarray -t instruction_files < <(get_manifest_files "INSTRUCTION_FILES" "${manifest_file}")
    readarray -t prompt_files < <(get_manifest_files "PROMPT_FILES" "${manifest_file}")
    readarray -t universal_files < <(get_manifest_files "UNIVERSAL_FILES" "${manifest_file}")
    
    # Calculate total files to install
    local total_files=$((${#main_files[@]} + ${#instruction_files[@]} + ${#prompt_files[@]} + ${#universal_files[@]}))
    local current_file=0
    
    # Install main files
    for file in "${main_files[@]}"; do
        [[ -z "${file}" ]] && continue
        current_file=$((current_file + 1))
        
        local filename=$(basename "${file}")
        local target_path="${workspace_root}/${file}"
        
        if declare -f show_completion >/dev/null 2>&1; then
            show_completion ${current_file} ${total_files} "Installing ${filename}"
        fi
        
        download_file "${file}" "${target_path}" "${filename}"
    done
    
    # Install instruction files
    local instructions_dir="${workspace_root}/.github/instructions"
    mkdir -p "${instructions_dir}"
    
    for file in "${instruction_files[@]}"; do
        [[ -z "${file}" ]] && continue
        current_file=$((current_file + 1))
        
        local filename=$(basename "${file}")
        local target_path="${workspace_root}/${file}"
        
        if declare -f show_completion >/dev/null 2>&1; then
            show_completion ${current_file} ${total_files} "Installing ${filename}"
        fi
        
        download_file "${file}" "${target_path}" "instructions/${filename}"
    done
    
    # Install prompt files
    local prompts_dir="${workspace_root}/.github/prompts"
    mkdir -p "${prompts_dir}"
    
    for file in "${prompt_files[@]}"; do
        [[ -z "${file}" ]] && continue
        current_file=$((current_file + 1))
        
        local filename=$(basename "${file}")
        local target_path="${workspace_root}/${file}"
        
        if declare -f show_completion >/dev/null 2>&1; then
            show_completion ${current_file} ${total_files} "Installing ${filename}"
        fi
        
        download_file "${file}" "${target_path}" "prompts/${filename}"
    done
    
    # Install universal files (like .vscode/settings.json)
    for file in "${universal_files[@]}"; do
        [[ -z "${file}" ]] && continue
        current_file=$((current_file + 1))
        
        local filename=$(basename "${file}")
        local target_path="${workspace_root}/${file}"
        
        if declare -f show_completion >/dev/null 2>&1; then
            show_completion ${current_file} ${total_files} "Installing ${filename}"
        fi
        
        download_file "${file}" "${target_path}" "${filename}"
    done
    
    if declare -f write_success >/dev/null 2>&1; then
        write_success "AI Infrastructure installation completed!"
    else
        echo -e "\033[0;32m[SUCCESS]\033[0m AI Infrastructure installation completed!"
    fi
}

# Export the new function
export -f install_infrastructure
