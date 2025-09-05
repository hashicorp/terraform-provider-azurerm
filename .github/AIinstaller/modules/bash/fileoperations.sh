#!/usr/bin/env bash
# FileOperations Module for Terraform AzureRM Provider AI Setup (Bash)
# Handles file operations, downloads, and installation tasks

# Bash 3.2 compatible function to read command output into array
# Usage: read_into_array array_name "command"
read_into_array() {
    local array_name="$1"
    local command="$2"
    local temp_file
    temp_file=$(mktemp)

    # Execute command and store output in temp file
    eval "$command" > "$temp_file"

    # Read lines into array using a counter
    local counter=0
    while IFS= read -r line; do
        eval "${array_name}[$counter]=\"\$line\""
        counter=$((counter + 1))
    done < "$temp_file"

    # Store array size for iteration
    eval "${array_name}_size=$counter"

    rm -f "$temp_file"
}

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

# Function to validate bootstrap operation prerequisites
validate_bootstrap_prerequisites() {
    local current_dir="$1"
    local user_profile="$2"
    local current_branch="$3"
    local branch_type="$4"

    # Rule 1: CRITICAL - Must NOT be running from user profile directory
    # Bootstrap should ONLY be run from the source repository, never from ~/.terraform-ai-installer
    if [[ "${current_dir}" == *".terraform-ai-installer"* ]]; then
        write_red " BOOTSTRAP VIOLATION: Cannot run bootstrap from user profile directory"
        echo ""
        write_yellow " Bootstrap must be run from the source terraform-provider-azurerm repository."
        write_yellow " You are currently running from: '${current_dir}'"
        echo ""
        print_separator
        echo ""
        write_cyan "SOLUTION:"
        write_cyan "  1. Navigate to the exp/terraform_copilot branch::"
        write_plain "    cd \"<path-to-your-terraform-provider-azurerm>\""
        write_plain "    git checkout exp/terraform_copilot"
        echo ""
        write_cyan "  2. Then run bootstrap from there:"
        write_plain "    ./.github/AIinstaller/install-copilot-setup.sh -bootstrap"
        echo ""
        return 1
    fi

    return 0
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

    # Check if file already exists and compare freshness
    if [[ -f "${target_path}" ]]; then
        local local_timestamp
        local remote_timestamp

        # Get local file modification time
        if command -v stat >/dev/null 2>&1; then
            # Try GNU stat first, then BSD stat
            local_timestamp=$(stat -c%Y "${target_path}" 2>/dev/null || stat -f%m "${target_path}" 2>/dev/null || echo "0")
        else
            local_timestamp="0"
        fi

        # Try to get remote file timestamp using HEAD request
        local is_remote_newer=true  # Default to downloading if we can't determine

        if command -v curl >/dev/null 2>&1; then
            # Get Last-Modified header from remote file
            local last_modified_header
            last_modified_header=$(curl -sI "${url}" 2>/dev/null | grep -i "last-modified:" | cut -d' ' -f2- | tr -d '\r\n')

            if [[ -n "${last_modified_header}" ]]; then
                # Convert HTTP date to timestamp (this is a simplified approach)
                if command -v date >/dev/null 2>&1; then
                    # Try to parse the HTTP date format
                    remote_timestamp=$(date -d "${last_modified_header}" +%s 2>/dev/null || echo "0")

                    if [[ "${remote_timestamp}" != "0" && "${local_timestamp}" != "0" ]]; then
                        if [[ "${remote_timestamp}" -le "${local_timestamp}" ]]; then
                            is_remote_newer=false
                        fi
                    fi
                fi
            fi
        fi

        # If local file is up-to-date, skip download
        if [[ "${is_remote_newer}" == "false" ]]; then
            return 0  # File is up-to-date, no need to download
        fi
    fi

    # Create target directory if it doesn't exist
    mkdir -p "$(dirname "${target_path}")"

    # Try download with proper error handling
    local download_success=false
    local error_msg=""

    if command -v curl >/dev/null 2>&1; then
        # Use curl with proper error handling
        if curl -fsSL "${url}" -o "${target_path}" 2>/dev/null; then
            download_success=true
        else
            error_msg="curl failed to download from ${url}"
        fi
    elif command -v wget >/dev/null 2>&1; then
        # Use wget with proper error handling
        if wget -q "${url}" -O "${target_path}" 2>/dev/null; then
            download_success=true
        else
            error_msg="wget failed to download from ${url}"
        fi
    else
        write_error_message "Neither curl nor wget is available for downloading files"
        return 1
    fi

    if [[ "${download_success}" == "true" ]]; then
        # Verify file was actually created and has content
        if [[ -f "${target_path}" && -s "${target_path}" ]]; then
            return 0
        else
            return 1
        fi
    else
        return 1
    fi
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
    local current_branch="$3"
    local branch_type="$4"

    # CRITICAL: Source branch protection - prevent operations on source branches
    # Check if we're on a protected source branch (main, master, exp/terraform_copilot)
    case "${current_branch}" in
        "main"|"master"|"exp/terraform_copilot")
            show_safety_violation "${current_branch}" "${operation_name}" "false" "${workspace_root}"
            return 1
            ;;
    esac

    # Repository structure validation using validation engine if available
    if declare -f validate_repository >/dev/null 2>&1; then
        # Check repository validation (returns 0 for valid, 1 for invalid)
        if ! validate_repository "${workspace_root}"; then
            echo ""
            write_error_message "WORKSPACE VALIDATION FAILED: Missing required files"
            echo ""
            write_plain "Expected to find terraform-provider-azurerm structure"
            write_plain "Directory: ${workspace_root}"
            write_plain "Reason: Missing required files (go.mod with azurerm content, main.go, or internal/services directory)"
            echo ""
            write_plain "Please ensure you are in the correct terraform-provider-azurerm repository directory."
            echo "${operation_name^} cancelled due to invalid repository structure"
            echo ""
            print_separator
            echo ""

            # Show help menu for guidance (matching PowerShell behavior)
            if declare -f show_usage >/dev/null 2>&1; then
                show_usage "${branch_type:-feature}" "false" "Missing required files"
            fi
            return 1
        fi
    fi

    # Operation is allowed
    return 0
}

# Function to remove deprecated files (equivalent to PowerShell Remove-DeprecatedFiles)
remove_deprecated_files() {
    local workspace_root="$1"
    local manifest_file="${2:-${HOME}/.terraform-ai-installer/file-manifest.config}"
    local dry_run="${3:-false}"
    local quiet="${4:-false}"

    local deprecated_count=0
    local deprecated_files=()

    # Ensure manifest file exists
    if [[ ! -f "${manifest_file}" ]]; then
        if [[ "${quiet}" == "false" ]]; then
            write_error_message "Manifest file not found: ${manifest_file}"
        fi
        return 1
    fi

    # Get current manifest files using bash 3.2 compatible method
    local current_instruction_files current_prompt_files
    read_into_array current_instruction_files "get_manifest_files \"INSTRUCTION_FILES\" \"${manifest_file}\""
    read_into_array current_prompt_files "get_manifest_files \"PROMPT_FILES\" \"${manifest_file}\""

    # Check existing instruction files in workspace
    local instructions_dir="${workspace_root}/.github/instructions"
    if [[ -d "${instructions_dir}" ]]; then
        while IFS= read -r -d '' existing_file; do
            if [[ -f "${existing_file}" && "${existing_file}" == *.instructions.md ]]; then
                local basename_file=$(basename "${existing_file}")
                local is_current=false

                # Check if this file is in current manifest (bash 3.2 compatible)
                local i=0
                while [[ $i -lt $current_instruction_files_size ]]; do
                    local current_file_var="current_instruction_files[$i]"
                    local current_file="${!current_file_var}"
                    if [[ "$(basename "${current_file}")" == "${basename_file}" ]]; then
                        is_current=true
                        break
                    fi
                    i=$((i + 1))
                done

                # If not in current manifest, mark for removal
                if [[ "${is_current}" == "false" ]]; then
                    deprecated_files+=("${existing_file}:Instruction:${basename_file}")

                    if [[ "${dry_run}" == "false" ]]; then
                        if rm -f "${existing_file}" 2>/dev/null; then
                            [[ "${quiet}" == "false" ]] && echo "  Removed deprecated instruction file: ${basename_file}"
                        else
                            [[ "${quiet}" == "false" ]] && write_error_message "  Failed to remove: ${existing_file}"
                        fi
                    else
                        [[ "${quiet}" == "false" ]] && echo "  [DRY-RUN] Would remove instruction file: ${basename_file}"
                    fi
                    ((deprecated_count++))
                fi
            fi
        done < <(find "${instructions_dir}" -name "*.instructions.md" -print0 2>/dev/null)
    fi

    # Check existing prompt files in workspace
    local prompts_dir="${workspace_root}/.github/prompts"
    if [[ -d "${prompts_dir}" ]]; then
        while IFS= read -r -d '' existing_file; do
            if [[ -f "${existing_file}" && "${existing_file}" == *.prompt.md ]]; then
                local basename_file=$(basename "${existing_file}")
                local is_current=false

                # Check if this file is in current manifest (bash 3.2 compatible)
                local j=0
                while [[ $j -lt $current_prompt_files_size ]]; do
                    local current_file_var="current_prompt_files[$j]"
                    local current_file="${!current_file_var}"
                    if [[ "$(basename "${current_file}")" == "${basename_file}" ]]; then
                        is_current=true
                        break
                    fi
                    j=$((j + 1))
                done

                # If not in current manifest, mark for removal
                if [[ "${is_current}" == "false" ]]; then
                    deprecated_files+=("${existing_file}:Prompt:${basename_file}")

                    if [[ "${dry_run}" == "false" ]]; then
                        if rm -f "${existing_file}" 2>/dev/null; then
                            [[ "${quiet}" == "false" ]] && echo "  Removed deprecated prompt file: ${basename_file}"
                        else
                            [[ "${quiet}" == "false" ]] && write_error_message "  Failed to remove: ${existing_file}"
                        fi
                    else
                        [[ "${quiet}" == "false" ]] && echo "  [DRY-RUN] Would remove prompt file: ${basename_file}"
                    fi
                    ((deprecated_count++))
                fi
            fi
        done < <(find "${prompts_dir}" -name "*.prompt.md" -print0 2>/dev/null)
    fi

    # Return count of deprecated files found/removed
    echo "${deprecated_count}"
}

# Note: Function exports moved to end of script

# Function to install AI infrastructure (moved from main script)
install_infrastructure() {
    local workspace_root="$1"
    local current_branch="$2"
    local branch_type="$3"

    write_section "Installing AI Infrastructure"

    # Validate operation is allowed on current branch
    if ! validate_operation_allowed "${workspace_root}" "installation" "${current_branch}" "${branch_type}"; then
        return 1
    fi

    # Verify manifest file exists
    local manifest_file="${HOME}/.terraform-ai-installer/file-manifest.config"
    if [[ ! -f "${manifest_file}" ]]; then
        write_error_message "Manifest file not found: ${manifest_file}"
        echo "Please run with -bootstrap first to set up the installer."
        return 1
    fi

    # Step 1: Show dry run notice if applicable
    if [[ "${DRY_RUN:-false}" == "true" ]]; then
        write_yellow "DRY RUN - No files will be created or removed"
        echo ""
    fi

    # Step 2: Check for deprecated files (automatic part of installation)
    write_cyan "Checking for deprecated files..."

    # Remove deprecated files automatically
    local deprecated_count
    deprecated_count=$(remove_deprecated_files "${workspace_root}" "${manifest_file}" "${DRY_RUN:-false}" false)

    if [[ "${deprecated_count}" -gt 0 ]]; then
        if [[ "${DRY_RUN:-false}" == "true" ]]; then
            write_cyan "  Found ${deprecated_count} deprecated file(s) that would be removed"
        else
            write_cyan "  Removed ${deprecated_count} deprecated file(s)"
        fi
    else
        write_cyan "  No deprecated files found"
    fi
    echo ""

    # Step 3: Installing current AI infrastructure files message
    write_cyan "Installing current AI infrastructure files..."

    # Step 4: Validation phase
    write_cyan "Validating installation prerequisites..."

    # Basic validation - workspace exists and is writable
    if [[ ! -d "${workspace_root}" ]]; then
        write_error_message "Workspace directory does not exist: ${workspace_root}"
        return 1
    fi

    if [[ ! -w "${workspace_root}" ]]; then
        write_error_message "Workspace directory is not writable: ${workspace_root}"
        return 1
    fi

    write_green "All prerequisites validated successfully!"
    echo ""

    # Step 5: Build complete file list (like PowerShell does)
    local all_files=()
    local file_destinations=()

    # Build unified file list with destination mappings
    build_file_list "${manifest_file}" "${workspace_root}" all_files file_destinations

    local total_files=${#all_files[@]}

    # Step 6: Show preparation message
    write_cyan "Preparing to install ${total_files} files..."
    echo ""

    # Step 7: Process all files in a single streamlined loop
    local successful_files=0
    local failed_files=0

    install_all_files all_files file_destinations successful_files failed_files

    # Show comprehensive installation summary like PowerShell (no duplicate success message)
    show_installation_summary "${workspace_root}" "${successful_files}" "${failed_files}" "${total_files}" "${current_branch}" "${branch_type}"
}

# Modular function to build unified file list with destination mappings
build_file_list() {
    local manifest_file="$1"
    local workspace_root="$2"
    local -n all_files_ref="$3"
    local -n file_destinations_ref="$4"

    # Read all file sections from manifest and build unified list (bash 3.2 compatible)
    local main_files instruction_files prompt_files universal_files

    read_into_array main_files "get_manifest_files \"MAIN_FILES\" \"${manifest_file}\""
    read_into_array instruction_files "get_manifest_files \"INSTRUCTION_FILES\" \"${manifest_file}\""
    read_into_array prompt_files "get_manifest_files \"PROMPT_FILES\" \"${manifest_file}\""
    read_into_array universal_files "get_manifest_files \"UNIVERSAL_FILES\" \"${manifest_file}\""

    # Add main files (root level)
    for file in "${main_files[@]}"; do
        [[ -z "${file}" ]] && continue
        all_files_ref+=("${file}")
        file_destinations_ref+=("${workspace_root}/${file}")
    done

    # Add instruction files (.github/instructions/)
    for file in "${instruction_files[@]}"; do
        [[ -z "${file}" ]] && continue
        local filename=$(basename "${file}")
        all_files_ref+=("${file}")
        file_destinations_ref+=("${workspace_root}/.github/instructions/${filename}")
    done

    # Add prompt files (.github/prompts/)
    for file in "${prompt_files[@]}"; do
        [[ -z "${file}" ]] && continue
        local filename=$(basename "${file}")
        all_files_ref+=("${file}")
        file_destinations_ref+=("${workspace_root}/.github/prompts/${filename}")
    done

    # Add universal files (various locations)
    for file in "${universal_files[@]}"; do
        [[ -z "${file}" ]] && continue
        all_files_ref+=("${file}")
        file_destinations_ref+=("${workspace_root}/${file}")
    done
}

# Streamlined function to install all files in a single loop
install_all_files() {
    local -n all_files_ref="$1"
    local -n file_destinations_ref="$2"
    local -n successful_ref="$3"
    local -n failed_ref="$4"

    local total_files=${#all_files_ref[@]}

    # Create necessary directories upfront
    mkdir -p "${workspace_root}/.github/instructions"
    mkdir -p "${workspace_root}/.github/prompts"

    # Temporarily disable exit on error for the download loop
    set +e

    for ((i=0; i<total_files; i++)); do
        local source_file="${all_files_ref[i]}"
        local target_path="${file_destinations_ref[i]}"
        local filename=$(basename "${source_file}")

        # Calculate relative path for display (like PowerShell)
        local relative_path
        relative_path=$(get_relative_display_path "${target_path}" "${workspace_root}")

        # Calculate and display progress with proper colors (matching PowerShell format)
        local percentage=$(( (i + 1) * 100 / total_files ))

        # Use right-aligned 3-digit format like show_completion function (automatically handles padding)
        printf "  ${CYAN}Downloading ${GREEN}[%3d%%]${CYAN}: ${NC}%s\n" "${percentage}" "${relative_path}"

        # Create target directory if needed
        local target_dir=$(dirname "${target_path}")
        mkdir -p "${target_dir}"

        # Determine download category for consistency
        local download_category
        download_category=$(get_download_category "${source_file}")

        # Download the file - handle errors gracefully to continue with other files
        if download_file "${source_file}" "${target_path}" "${download_category}/${filename}"; then
            ((successful_ref++))
        else
            ((failed_ref++))
            # Continue with next file even if this one failed (silent like PowerShell)
        fi
    done

    # Re-enable exit on error
    set -e
}

# Helper function to get relative display path (like PowerShell shows)
get_relative_display_path() {
    local target_path="$1"
    local workspace_root="$2"

    # Remove workspace root prefix to show relative path
    local relative_path="${target_path#${workspace_root}/}"
    echo "${relative_path}"
}

# Helper function to determine download category
get_download_category() {
    local source_file="$1"

    # Use pattern matching instead of regex for bash 3.2 compatibility
    case "${source_file}" in
        *.github/instructions/*)
            echo "instructions"
            ;;
        *.github/prompts/*)
            echo "prompts"
            ;;
        *)
            echo "files"
            ;;
    esac
}

# Function to calculate total size of files in KB
calculate_installed_files_size() {
    local workspace_root="$1"
    local manifest_file="${HOME}/.terraform-ai-installer/file-manifest.config"
    local total_size_bytes=0

    # Return 0 if manifest doesn't exist
    if [[ ! -f "${manifest_file}" ]]; then
        echo "0"
        return
    fi

    # Get all file lists from manifest (bash 3.2 compatible)
    local all_files=()
    local main_files instruction_files prompt_files universal_files

    read_into_array main_files "get_manifest_files \"MAIN_FILES\" \"${manifest_file}\" 2>/dev/null || true"
    read_into_array instruction_files "get_manifest_files \"INSTRUCTION_FILES\" \"${manifest_file}\" 2>/dev/null || true"
    read_into_array prompt_files "get_manifest_files \"PROMPT_FILES\" \"${manifest_file}\" 2>/dev/null || true"
    read_into_array universal_files "get_manifest_files \"UNIVERSAL_FILES\" \"${manifest_file}\" 2>/dev/null || true"

    # Calculate total size by checking each installed file (bash 3.2 compatible)
    local i=0
    while [[ $i -lt $main_files_size ]]; do
        local file_var="main_files[$i]"
        local file="${!file_var}"
        [[ -z "${file}" ]] && { i=$((i + 1)); continue; }
        local target_path="${workspace_root}/${file}"
        if [[ -f "${target_path}" ]]; then
            local file_size
            file_size=$(stat -f%z "${target_path}" 2>/dev/null || stat -c%s "${target_path}" 2>/dev/null || echo "0")
            total_size_bytes=$((total_size_bytes + file_size))
        fi
        i=$((i + 1))
    done

    # Check instruction files
    i=0
    while [[ $i -lt $instruction_files_size ]]; do
        local file_var="instruction_files[$i]"
        local file="${!file_var}"
        [[ -z "${file}" ]] && { i=$((i + 1)); continue; }
        local filename=$(basename "${file}")
        local target_path="${workspace_root}/.github/instructions/${filename}"
        if [[ -f "${target_path}" ]]; then
            local file_size
            file_size=$(stat -f%z "${target_path}" 2>/dev/null || stat -c%s "${target_path}" 2>/dev/null || echo "0")
            total_size_bytes=$((total_size_bytes + file_size))
        fi
        i=$((i + 1))
    done

    # Check prompt files
    i=0
    while [[ $i -lt $prompt_files_size ]]; do
        local file_var="prompt_files[$i]"
        local file="${!file_var}"
        [[ -z "${file}" ]] && { i=$((i + 1)); continue; }
        local filename=$(basename "${file}")
        local target_path="${workspace_root}/.github/prompts/${filename}"
        if [[ -f "${target_path}" ]]; then
            local file_size
            file_size=$(stat -f%z "${target_path}" 2>/dev/null || stat -c%s "${target_path}" 2>/dev/null || echo "0")
            total_size_bytes=$((total_size_bytes + file_size))
        fi
        i=$((i + 1))
    done

    # Check universal files
    i=0
    while [[ $i -lt $universal_files_size ]]; do
        local file_var="universal_files[$i]"
        local file="${!file_var}"
        [[ -z "${file}" ]] && { i=$((i + 1)); continue; }
        local target_path="${workspace_root}/${file}"
        if [[ -f "${target_path}" ]]; then
            local file_size
            file_size=$(stat -f%z "${target_path}" 2>/dev/null || stat -c%s "${target_path}" 2>/dev/null || echo "0")
            total_size_bytes=$((total_size_bytes + file_size))
        fi
        i=$((i + 1))
    done

    # Convert to KB (rounded up)
    local total_size_kb=$((total_size_bytes / 1024))
    if [[ $((total_size_bytes % 1024)) -gt 0 ]]; then
        total_size_kb=$((total_size_kb + 1))
    fi

    echo "${total_size_kb}"
}

# Comprehensive installation summary (matching PowerShell style)
show_installation_summary() {
    local workspace_root="$1"
    local successful_files="$2"
    local failed_files="$3"
    local total_files="$4"
    local current_branch="$5"
    local branch_type="$6"

    # Calculate actual total size of installed files
    local total_size_kb
    total_size_kb=$(calculate_installed_files_size "${workspace_root}")
    local skipped_files=$((total_files - successful_files - failed_files))

    # Show detailed summary using the sophisticated show_operation_summary function
    # Clean branch_type variable to remove any potential line breaks or whitespace
    branch_type=$(echo "${branch_type}" | tr -d '\r\n' | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')

    show_operation_summary "Installation" "true" "false" \
        "Branch Type: ${branch_type}" \
        "Target Branch: ${current_branch}" \
        "Files Installed: ${successful_files}" \
        "Total Size: ${total_size_kb} KB" \
        "Files Skipped: ${skipped_files}" \
        "Location: ${workspace_root}"
}

# Function to get all files that should be cleaned up from manifest
get_files_for_cleanup() {
    local workspace_root="$1"
    local manifest_file="${HOME}/.terraform-ai-installer/file-manifest.config"

    # Check if manifest file exists
    if [[ ! -f "${manifest_file}" ]]; then
        write_error_message "Manifest file not found: ${manifest_file}"
        return 1
    fi

    # Get all file sections from manifest
    local all_files=()

    # Get files from each section (bash 3.2 compatible)
    local main_files instruction_files prompt_files universal_files
    read_into_array main_files "get_manifest_files \"MAIN_FILES\" \"${manifest_file}\" 2>/dev/null || true"
    read_into_array instruction_files "get_manifest_files \"INSTRUCTION_FILES\" \"${manifest_file}\" 2>/dev/null || true"
    read_into_array prompt_files "get_manifest_files \"PROMPT_FILES\" \"${manifest_file}\" 2>/dev/null || true"
    read_into_array universal_files "get_manifest_files \"UNIVERSAL_FILES\" \"${manifest_file}\" 2>/dev/null || true"

    # Combine all files into one list (bash 3.2 compatible)
    # Get all files and combine them in one output
    {
        get_manifest_files "MAIN_FILES" "${manifest_file}" 2>/dev/null || true
        get_manifest_files "INSTRUCTION_FILES" "${manifest_file}" 2>/dev/null || true
        get_manifest_files "PROMPT_FILES" "${manifest_file}" 2>/dev/null || true
        get_manifest_files "UNIVERSAL_FILES" "${manifest_file}" 2>/dev/null || true
    } | while IFS= read -r file; do
        # Skip empty entries
        if [[ -z "${file}" ]]; then
            continue
        fi

        # Output unique files only
        echo "${file}"
    done | sort -u
}

# Function to clean AI infrastructure files with empty directory cleanup
clean_infrastructure() {
    local workspace_root="$1"
    local current_branch="$2"
    local branch_type="$3"

    # CRITICAL: Clean operations are FORBIDDEN on source branches for safety
    # Validate operation is allowed on current branch BEFORE showing section header
    if ! validate_operation_allowed "${workspace_root}" "cleanup" "${current_branch}" "${branch_type}"; then
        return 1
    fi

    # Show section header for cleanup operation
    write_section "Clean Workspace"
    if [[ -z "${workspace_root}" ]]; then
        write_error_message "Workspace root directory not specified"
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
    local existing_files=()
    local directories_to_check=()
    local existing_directories=()

    # Directories that should NEVER be removed (important repository infrastructure)
    local protected_directories=(
        ".github"  # Main GitHub directory contains workflows, issue templates, etc.
        ".vscode"  # VS Code workspace settings and configurations
        "internal" # Core provider code
        "vendor"   # Go dependencies
        "scripts"  # Build and maintenance scripts
        "website"  # Documentation
        "examples" # Example configurations
        "helpers"  # Helper utilities
    )

    # AI directories that are safe to clean up (matches PowerShell version)
    local ai_directories=(
        ".github/instructions"
        ".github/prompts"
    )

    # Prepare file and directory lists
    for file in "${cleanup_files[@]}"; do
        local full_path="${workspace_root}/${file}"
        files_to_remove+=("${full_path}")

        # Track existing files for progress calculation
        if [[ -e "${full_path}" ]]; then
            existing_files+=("${file}")
        fi

        # Add parent directory to list for empty directory cleanup
        local parent_dir="$(dirname "${full_path}")"
        if [[ "${parent_dir}" != "${workspace_root}" ]]; then
            # Only allow cleanup of specific AI infrastructure directories
            # This prevents accidental removal of important repository directories
            local relative_dir="${parent_dir#${workspace_root}/}"
            local allowed_for_cleanup=false

            # Check if directory is in our safe AI directories list
            for ai_dir in "${ai_directories[@]}"; do
                if [[ "${relative_dir}" == "${ai_dir}" ]] || [[ "${relative_dir}" == "${ai_dir}/"* ]]; then
                    allowed_for_cleanup=true
                    break
                fi
            done

            if [[ "${allowed_for_cleanup}" == "true" ]]; then
                directories_to_check+=("${parent_dir}")
            fi
        fi
    done

    # Check for existing directories
    for dir in "${directories_to_check[@]}"; do
        if [[ -d "${dir}" ]]; then
            existing_directories+=("${dir}")
        fi
    done

    # Remove duplicates from directories
    local unique_dirs=($(printf '%s\n' "${existing_directories[@]}" | sort -u))
    existing_directories=("${unique_dirs[@]}")

    # If nothing exists, show clean message and exit early (match PowerShell behavior)
    if [[ ${#existing_files[@]} -eq 0 ]] && [[ ${#existing_directories[@]} -eq 0 ]]; then
        write_green " No AI infrastructure files found to remove."
        write_green " Workspace is already clean!"

        # Show detailed cleanup summary using the sophisticated show_operation_summary function
        # for "already clean" scenario - follows PowerShell master order
        show_operation_summary "Cleanup" "true" "${DRY_RUN:-false}" \
            "Branch Type: ${branch_type}" \
            "Target Branch: ${current_branch}" \
            "Operation Type: $(if [[ "${DRY_RUN:-false}" == "true" ]]; then echo "Dry run cleanup"; else echo "Live cleanup"; fi)" \
            "Files Removed: 0" \
            "Directories Cleaned: 0" \
            "Location: ${workspace_root}"

        return 0
    fi

    # Show scanning message
    write_cyan "Scanning for AI infrastructure files..."
    write_yellow "Found ${#existing_files[@]} AI files and ${#existing_directories[@]} directories to remove."

    # Calculate maximum filename length for aligned output
    local max_name_length=0
    for file in "${existing_files[@]}"; do
        local filename=$(basename "${file}")
        local length=${#filename}
        if [[ $length -gt $max_name_length ]]; then
            max_name_length=$length
        fi
    done

    # Also check directory names for padding
    for dir in "${existing_directories[@]}"; do
        local dirname=$(basename "${dir}")
        local length=${#dirname}
        if [[ $length -gt $max_name_length ]]; then
            max_name_length=$length
        fi
    done

    # Display header
    echo ""
    write_cyan "Removing AI Infrastructure Files"
    print_separator

    # Calculate total work for progress - add safety check
    local total_work=$((${#existing_files[@]} + ${#existing_directories[@]}))
    local work_completed=0

    # Remove files (only process existing ones)
    for file in "${existing_files[@]}"; do
        # Increment work completed
        work_completed=$((work_completed + 1))

        # Safe arithmetic with explicit error handling
        local percent_complete=0
        if [[ $total_work -gt 0 ]]; then
            percent_complete=$(( (work_completed * 100) / total_work ))
        fi

        # Extract just the filename for cleaner display
        local filename=$(basename "${file}")
        local full_path="${workspace_root}/${file}"

        # Calculate padding for filename alignment (same as PowerShell)
        local filename_padding=""
        local padding_needed=$((max_name_length - ${#filename}))
        if [[ $padding_needed -gt 0 ]]; then
            filename_padding=$(printf "%*s" $padding_needed "")
        fi

        # Dynamic padding for progress percentage alignment (matches PowerShell)
        local progress_padding=""
        if [[ $percent_complete -lt 10 ]]; then
            progress_padding="  "  # 1-digit: 2 spaces
        elif [[ $percent_complete -lt 100 ]]; then
            progress_padding=" "   # 2-digit: 1 space
        else
            progress_padding=""    # 3-digit: 0 spaces
        fi

        # Display removal operation using exact PowerShell format
        printf "  ${CYAN}Removing File      ${GREEN}[%s%%%s]${CYAN}: ${NC}%s%s " "${percent_complete}" "${progress_padding}" "${filename}" "${filename_padding}"

        if [[ "${DRY_RUN:-false}" == "true" ]]; then
            echo -e "${YELLOW}[WOULD REMOVE]${NC}"
        else
            if rm -rf "${full_path}" 2>/dev/null; then
                echo -e "${GREEN}[OK]${NC}"
            else
                echo -e "${RED}[FAILED]${NC}"
            fi
        fi
    done

    # Process directories (show progress in both dry-run and live modes)
    # Sort directories by depth (deepest first) to remove from bottom up
    local sorted_dirs=($(printf '%s\n' "${existing_directories[@]}" | sort -r))

    for dir in "${sorted_dirs[@]}"; do
        # Always increment work_completed for progress tracking (whether we remove or skip)
        work_completed=$((work_completed + 1))

        # Safe arithmetic with explicit error handling
        local percent_complete=0
        if [[ $total_work -gt 0 ]]; then
            percent_complete=$(( (work_completed * 100) / total_work ))
        fi

        local dirname=$(basename "${dir}")

        # Calculate padding for directory name alignment (same as files)
        local dirname_padding=""
        local padding_needed=$((max_name_length - ${#dirname}))
        if [[ $padding_needed -gt 0 ]]; then
            dirname_padding=$(printf "%*s" $padding_needed "")
        fi

        # Dynamic padding for progress percentage alignment (matches PowerShell)
        local progress_padding=""
        if [[ $percent_complete -lt 10 ]]; then
            progress_padding="  "  # 1-digit: 2 spaces
        elif [[ $percent_complete -lt 100 ]]; then
            progress_padding=" "   # 2-digit: 1 space
        else
            progress_padding=""    # 3-digit: 0 spaces
        fi

        # Display directory removal using exact PowerShell format
        printf "  ${CYAN}Removing Directory ${GREEN}[%s%%%s]${CYAN}: ${NC}%s%s " "${percent_complete}" "${progress_padding}" "${dirname}" "${dirname_padding}"

        if [[ "${DRY_RUN:-false}" == "true" ]]; then
            echo -e "${YELLOW}[WOULD REMOVE]${NC}"
        else
            # Only actually try to remove if directory exists and is empty (live mode only)
            if [[ -d "${dir}" ]] && [[ -z "$(ls -A "${dir}" 2>/dev/null)" ]]; then
                if rmdir "${dir}" 2>/dev/null; then
                    echo -e "${GREEN}[OK]${NC}"
                else
                    echo -e "${YELLOW}[NOT EMPTY]${NC}"
                fi
            else
                # Directory doesn't exist or isn't empty - show as skipped
                echo -e "${YELLOW}[SKIPPED]${NC}"
            fi
        fi
    done

    # Count actual removed items for summary
    local files_removed=${#existing_files[@]}
    local dirs_removed=${#existing_directories[@]}

    # Determine success status
    local success_status="true"  # Cleanup operations are considered successful unless there was an error

    # Show detailed cleanup summary using the sophisticated show_operation_summary function
    # follows PowerShell master order
    branch_type=$(echo "${branch_type}" | tr -d '\r\n' | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')

    show_operation_summary "Cleanup" "${success_status}" "${DRY_RUN:-false}" \
        "Branch Type: ${branch_type}" \
        "Target Branch: ${current_branch}" \
        "Operation Type: $(if [[ "${DRY_RUN:-false}" == "true" ]]; then echo "Dry run cleanup"; else echo "Live cleanup"; fi)" \
        "Files Removed: ${files_removed}" \
        "Directories Cleaned: ${dirs_removed}" \
        "Location: ${workspace_root}"

    # Disable error tracing
    set +e
    set +x
}

# Note: Function exports moved to end of script

# Function to perform bootstrap operation (copy installer files to user profile)
bootstrap_files_to_profile() {
    local current_dir="$1"
    local user_profile="$2"
    local manifest_file="$3"
    local current_branch="$4"
    local branch_type="$5"
    local script_dir="$6"

    # Validate bootstrap prerequisites before proceeding
    if ! validate_bootstrap_prerequisites "${current_dir}" "${user_profile}" "${current_branch}" "${branch_type}"; then
        return 1
    fi

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
export -f get_workspace_root is_source_repository validate_operation_allowed validate_bootstrap_prerequisites
export -f copy_file download_file get_file_size get_directory_size create_directory_structure
export -f remove_path backup_file verify_file make_executable copy_files_with_stats
export -f remove_deprecated_files install_infrastructure clean_infrastructure bootstrap_files_to_profile
export -f get_files_for_cleanup
