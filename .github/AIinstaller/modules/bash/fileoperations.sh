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

# Function to copy file with progress
copy_file() {
    local source="$1"
    local target="$2"
    local description="$3"
    local max_length="$4"  # Optional max length for formatting

    if [[ "${DRY_RUN:-false}" == "true" ]]; then
        echo "  [DRY-RUN] Would copy: ${description}"
        return 0
    fi

    # Create target directory if it doesn't exist
    mkdir -p "$(dirname "${target}")"

    if cp "${source}" "${target}"; then
        show_file_operation "Copying" "${description}" "OK" "${max_length}"
        return 0
    else
        show_file_operation "Copying" "${description}" "FAILED" "${max_length}"
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
        write_error_message "Neither curl nor wget is available for downloading files"
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

    # Calculate max filename length for consistent formatting (extract filenames only)
    local filenames_only=()
    for file in "${files_array[@]}"; do
        filenames_only+=("$(basename "${file}")")
    done
    local max_length
    max_length=$(calculate_max_filename_length "${filenames_only[@]}")

    for file in "${files_array[@]}"; do
        local source_path="${source_dir}/${file}"
        local target_path="${target_dir}/${file}"
        local filename=$(basename "${file}")

        echo -n ""  # No inline display needed - show_file_operation handles it

        if [[ "${DRY_RUN:-false}" == "true" ]]; then
            show_file_operation "Copying" "${filename}" "DRY-RUN" "${max_length}"
            files_copied=$((files_copied + 1))
        elif [[ -f "${source_path}" ]]; then
            # Create target directory structure
            mkdir -p "$(dirname "${target_path}")"

            if cp "${source_path}" "${target_path}"; then
                # Calculate file size
                local file_size
                file_size=$(get_file_size "${target_path}")
                total_size=$((total_size + file_size))

                show_file_operation "Copying" "${filename}" "OK" "${max_length}"
                files_copied=$((files_copied + 1))
            else
                show_file_operation "Copying" "${filename}" "FAILED" "${max_length}"
                files_failed=$((files_failed + 1))
            fi
        else
            show_file_operation "Copying" "${filename}" "SOURCE NOT FOUND" "${max_length}"
            files_failed=$((files_failed + 1))
        fi
    done

    # Return statistics via global variables (bash doesn't support complex return values)
    COPY_STATS_FILES_COPIED=${files_copied}
    COPY_STATS_FILES_FAILED=${files_failed}
    COPY_STATS_TOTAL_SIZE=${total_size}

    return $([[ ${files_failed} -eq 0 ]] && echo 0 || echo 1)
}

# Function to validate operation is allowed on current branch
validate_operation_allowed() {
    local workspace_root="$1"
    local operation_name="${2:-operation}"

    # CRITICAL: Source branch protection - prevent operations on source branches
    # Check if we're on a protected source branch (main, master, exp/terraform_copilot)
    local current_branch
    current_branch=$(cd "${workspace_root}" && git branch --show-current 2>/dev/null || echo "Unknown")

    case "${current_branch}" in
        "main"|"master"|"exp/terraform_copilot")
            write_error_message "${operation_name^^} BLOCKED: Cannot perform ${operation_name} on source branch '${current_branch}'"
            echo ""
            write_plain "Source branches (main, master, exp/terraform_copilot) are protected from ${operation_name}."
            write_plain "This prevents accidental modification of the source repository."
            echo ""
            write_plain "${YELLOW}REQUIRED ACTIONS:${NC}"
            write_plain "  1. Switch to a feature branch: git checkout -b feature/your-branch-name"
            write_plain "  2. Run ${operation_name} from the feature branch"
            echo ""
            return 1
            ;;
    esac

    # Repository structure validation using validation engine if available
    if declare -f validate_repository >/dev/null 2>&1; then
        local validation_result
        validation_result=$(validate_repository "${workspace_root}")

        # Parse validation results
        local valid=$(echo "${validation_result}" | grep "Valid=" | cut -d'=' -f2)

        if [[ "${valid}" != "true" ]]; then
            local reason=$(echo "${validation_result}" | grep "Reason=" | cut -d'=' -f2-)
            write_warning "Target directory may not be a terraform-provider-azurerm repository"
            echo "Expected to find terraform-provider-azurerm structure"
            echo "Directory: ${workspace_root}"
            echo "Reason: ${reason}"
            echo ""
            echo "Continue anyway? (y/N)"
            read -r response
            if [[ ! "${response}" =~ ^[Yy]$ ]]; then
                echo "${operation_name^} cancelled by user"
                return 1
            fi
        fi
    fi

    # Operation is allowed
    return 0
}

# Note: Function exports moved to end of script

# Function to install AI infrastructure (moved from main script)
install_infrastructure() {
    local workspace_root="$1"

    write_section "Installing AI Infrastructure"

    # Validate operation is allowed on current branch
    if ! validate_operation_allowed "${workspace_root}" "installation"; then
        return 1
    fi

    # Verify manifest file exists
    local manifest_file="${HOME}/.terraform-ai-installer/file-manifest.config"
    if [[ ! -f "${manifest_file}" ]]; then
        write_error_message "Manifest file not found: ${manifest_file}"
        echo "Please run with -bootstrap first to set up the installer."
        return 1
    fi

    echo "Installing to workspace: ${workspace_root}"

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

    write_operation_status "AI Infrastructure installation completed!" "Success"
}

# Function to clean AI infrastructure files with empty directory cleanup
clean_infrastructure() {
    local workspace_root="$1"

    if [[ -z "${workspace_root}" ]]; then
        write_error_message "Workspace root directory not specified"
        return 1
    fi

    # Validate operation is allowed on current branch
    if ! validate_operation_allowed "${workspace_root}" "cleanup"; then
        return 1
    fi

    # Get files from manifest
    local cleanup_files
    cleanup_files=($(get_files_for_cleanup "${workspace_root}"))

    if [[ ${#cleanup_files[@]} -eq 0 ]]; then
        write_error_message "No files found in manifest to clean"
        return 1
    fi

    local files_to_remove=()
    local directories_to_check=()

    # Prepare file and directory lists
    for file in "${cleanup_files[@]}"; do
        local full_path="${workspace_root}/${file}"
        files_to_remove+=("${full_path}")

        # Add parent directory to list for empty directory cleanup
        local parent_dir="$(dirname "${full_path}")"
        if [[ "${parent_dir}" != "${workspace_root}" ]]; then
            directories_to_check+=("${parent_dir}")
        fi
    done

    # Remove files
    for file in "${files_to_remove[@]}"; do
        if [[ -e "${file}" ]]; then
            if [[ "${DRY_RUN}" == "true" ]]; then
                write_operation_status "[DRY-RUN] Would remove: ${file}" "Info"
            else
                rm -rf "${file}"
                write_operation_status "Removed: ${file}" "Success"
            fi
        fi
    done

    # Remove empty parent directories (only if not in dry-run mode)
    if [[ "${DRY_RUN}" != "true" ]]; then
        # Sort directories by depth (deepest first) to remove from bottom up
        local sorted_dirs=($(printf '%s\n' "${directories_to_check[@]}" | sort -u | sort -r))

        for dir in "${sorted_dirs[@]}"; do
            # Only remove if directory exists and is empty
            if [[ -d "${dir}" ]] && [[ -z "$(ls -A "${dir}" 2>/dev/null)" ]]; then
                if rmdir "${dir}" 2>/dev/null; then
                    write_operation_status "Removed empty directory: ${dir}" "Success"
                fi
            fi
        done
    fi

    echo ""
    write_operation_status "AI Infrastructure cleanup completed!" "Success"
    echo ""
}

# Note: Function exports moved to end of script

# Function to perform bootstrap operation (copy installer files to user profile)
bootstrap_files_to_profile() {
    local script_dir="$1"
    local user_profile="$2"
    local manifest_file="$3"

    # Get bootstrap files from configuration
    local bootstrap_files_list
    bootstrap_files_list=$(get_manifest_files "INSTALLER_FILES_BOOTSTRAP" "${manifest_file}")

    if [[ -z "${bootstrap_files_list}" ]]; then
        write_error_message "No bootstrap files found in manifest: ${manifest_file}"
        write_plain "Expected section: [INSTALLER_FILES_BOOTSTRAP]"
        return 1
    fi

    # Convert to array
    local -a bootstrap_files
    while IFS= read -r line; do
        if [[ -n "${line}" ]]; then
            bootstrap_files+=("${line}")
        fi
    done <<< "${bootstrap_files_list}"

    # Count files for progress
    local total_files=${#bootstrap_files[@]}
    echo -e "${CYAN}Copying installer files from current repository...${NC}"
    echo ""

    # Calculate max filename length for formatting
    local -a bootstrap_filenames
    for file in "${bootstrap_files[@]}"; do
        bootstrap_filenames+=("$(basename "${file}")")
    done
    local max_length
    max_length=$(calculate_max_filename_length "${bootstrap_filenames[@]}")

    # Statistics tracking
    local files_copied=0
    local files_failed=0
    local total_size=0

    # Copy each file according to configuration
    for file_path in "${bootstrap_files[@]}"; do
        # Remove .github/AIinstaller/ prefix to get relative path within installer directory
        local relative_path="${file_path#.github/AIinstaller/}"
        local source_path="${script_dir}/${relative_path}"
        local filename=$(basename "${relative_path}")

        # Determine target path based on file type and maintain directory structure
        local target_path
        if [[ "${filename}" == *.psm1 ]]; then
            # PowerShell modules go in modules/powershell/ subdirectory
            if [[ "${DRY_RUN}" != "true" ]]; then
                mkdir -p "${user_profile}/modules/powershell"
            fi
            target_path="${user_profile}/modules/powershell/${filename}"
        elif [[ "${filename}" == *.sh ]] && [[ "${relative_path}" == modules/bash/* ]]; then
            # Bash modules go in modules/bash/ subdirectory
            if [[ "${DRY_RUN}" != "true" ]]; then
                mkdir -p "${user_profile}/modules/bash"
            fi
            target_path="${user_profile}/modules/bash/${filename}"
        else
            # Main files (config, scripts) go directly in target directory
            target_path="${user_profile}/${filename}"
        fi

        if copy_file "${source_path}" "${target_path}" "${filename}" "${max_length}"; then
            files_copied=$((files_copied + 1))
            # Calculate file size
            if [[ -f "${target_path}" ]]; then
                local file_size
                file_size=$(get_file_size "${target_path}")
                total_size=$((total_size + file_size))
            fi
        else
            files_failed=$((files_failed + 1))
        fi
    done

    # Make installer script executable
    if [[ "${DRY_RUN}" != "true" ]]; then
        chmod +x "${user_profile}/install-copilot-setup.sh"
    fi

    # Validate critical files were copied
    local critical_files=(
        "${user_profile}/install-copilot-setup.sh"
        "${user_profile}/file-manifest.config"
    )

    local validation_failed=false
    for critical_file in "${critical_files[@]}"; do
        if [[ ! -f "${critical_file}" ]]; then
            write_error_message "Critical file missing: ${critical_file}"
            validation_failed=true
        fi
    done

    if [[ "${validation_failed}" == "true" ]]; then
        write_error_message "Bootstrap validation failed - critical files missing"
        return 1
    fi

    # Return statistics via global variables
    BOOTSTRAP_STATS_FILES_COPIED=${files_copied}
    BOOTSTRAP_STATS_FILES_FAILED=${files_failed}
    BOOTSTRAP_STATS_TOTAL_SIZE=${total_size}

    return $([[ ${files_failed} -eq 0 ]] && echo 0 || echo 1)
}

# ==============================================================================
# FUNCTION EXPORTS - PowerShell-like pattern (all exports at end)
# ==============================================================================

# Export all functions for use in other scripts
export -f get_workspace_root is_source_repository validate_operation_allowed
export -f copy_file download_file get_file_size get_directory_size create_directory_structure
export -f remove_path backup_file verify_file make_executable copy_files_with_stats
export -f install_infrastructure clean_infrastructure bootstrap_files_to_profile
