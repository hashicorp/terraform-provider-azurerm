#!/usr/bin/env bash
# FileOperations Module for Terraform AzureRM Provider AI Setup (Bash)
# Handles file operations, downloads, and installation tasks

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
        echo "  Copying: ${description} [OK]"
        return 0
    else
        echo "  Copying: ${description} [FAILED]"
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
        
        echo -n "    Copying: ${filename}"
        
        if [[ "${DRY_RUN:-false}" == "true" ]]; then
            echo " [DRY-RUN]"
            files_copied=$((files_copied + 1))
        elif [[ -f "${source_path}" ]]; then
            # Create target directory structure
            mkdir -p "$(dirname "${target_path}")"
            
            if cp "${source_path}" "${target_path}"; then
                # Calculate file size
                local file_size
                file_size=$(get_file_size "${target_path}")
                total_size=$((total_size + file_size))
                
                echo -e " \033[0;32m[OK]\033[0m"
                files_copied=$((files_copied + 1))
            else
                echo -e " \033[0;31m[FAILED]\033[0m"
                files_failed=$((files_failed + 1))
            fi
        else
            echo -e " \033[0;31m[SOURCE NOT FOUND]\033[0m"
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
export -f copy_file download_file get_file_size get_directory_size create_directory_structure
export -f remove_path backup_file verify_file make_executable copy_files_with_stats
