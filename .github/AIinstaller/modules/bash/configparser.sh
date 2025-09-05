#!/usr/bin/env bash
# ConfigParser Module for Terraform AzureRM Provider AI Setup (Bash)
# STREAMLINED VERSION - Contains only functions actually used by main script

# Function to get manifest configuration - equivalent to PowerShell Get-ManifestConfig
get_manifest_config() {
    local manifest_path="${1:-}"
    local branch="${2:-exp/terraform_copilot}"

    # Find manifest file if not specified
    if [[ -z "${manifest_path}" ]]; then
        # Try multiple locations (equivalent to PowerShell logic)
        local possible_paths=(
            "${HOME}/.terraform-ai-installer/file-manifest.config"
            "$(dirname "$(dirname "$0")")/file-manifest.config"
            "$(dirname "$(dirname "$(dirname "$0")")")/file-manifest.config"
        )

        for path in "${possible_paths[@]}"; do
            if [[ -f "${path}" ]]; then
                manifest_path="${path}"
                break
            fi
        done

        # Fallback
        if [[ -z "${manifest_path}" ]]; then
            local script_root="$(dirname "$(dirname "$0")")"
            manifest_path="${script_root}/file-manifest.config"
        fi
    fi

    if [[ ! -f "${manifest_path}" ]]; then
        write_error_message "Manifest file not found: ${manifest_path}"
        return 1
    fi

    # Parse manifest sections (simplified - just check if file exists and is readable)
    local has_instructions=false
    local has_universal=false

    if grep -q "\[INSTRUCTION_FILES\]" "${manifest_path}" 2>/dev/null; then
        has_instructions=true
    fi

    if grep -q "\[UNIVERSAL_FILES\]" "${manifest_path}" 2>/dev/null; then
        has_universal=true
    fi

    # Return structured data (equivalent to PowerShell PSCustomObject)
    echo "ManifestPath=${manifest_path}"
    echo "Branch=${branch}"
    echo "HasInstructions=${has_instructions}"
    echo "HasUniversal=${has_universal}"
    echo "SourceRepository=https://raw.githubusercontent.com/hashicorp/terraform-provider-azurerm"

    return 0
}

# Function to get installer configuration - equivalent to PowerShell Get-InstallerConfig
get_installer_config() {
    local workspace_root="${1:-}"

    if [[ -z "${workspace_root}" ]]; then
        workspace_root="$(get_workspace_root)"
    fi

    if [[ -z "${workspace_root}" ]]; then
        write_error_message "Could not determine workspace root"
        return 1
    fi

    # Return installer configuration (equivalent to PowerShell logic)
    echo "WorkspaceRoot=${workspace_root}"
    echo "Repository=hashicorp/terraform-provider-azurerm"
    echo "InstallationPath=${HOME}/.terraform-ai-installer"

    return 0
}

# Helper function to get workspace root (used by get_installer_config)
get_workspace_root() {
    # Try to find .git directory to determine workspace root
    local current_dir="$(pwd)"
    local search_dir="${current_dir}"

    # Search up to 10 levels for .git directory
    local i=1
    while [[ $i -le 10 ]]; do
        if [[ -d "${search_dir}/.git" ]]; then
            echo "${search_dir}"
            return 0
        fi

        local parent_dir="$(dirname "${search_dir}")"
        if [[ "${parent_dir}" == "${search_dir}" ]]; then
            break  # Reached filesystem root
        fi
        search_dir="${parent_dir}"
        i=$((i + 1))
    done

    # Fallback to current directory
    echo "${current_dir}"
    return 0
}

# Function to parse manifest section (used by get_manifest_files)
parse_manifest_section() {
    local manifest_file="$1"
    local section_name="$2"

    if [[ ! -f "${manifest_file}" ]]; then
        write_error_message "Manifest file not found: ${manifest_file}"
        return 1
    fi

    local in_section=false

    while IFS= read -r line; do
        # Remove leading/trailing whitespace
        line=$(echo "${line}" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')

        # Skip empty lines and comments
        if [[ -z "${line}" ]] || echo "${line}" | grep -q '^#'; then
            continue
        fi

        # Check for section headers [SECTION_NAME]
        if echo "${line}" | grep -q '^\[[^]]*\]$'; then
            local current_section="$(echo "${line}" | sed 's/^\[\([^]]*\)\]$/\1/')"
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

# Function to get files from manifest by section (used by fileoperations module)
get_manifest_files() {
    local section_name="$1"
    local manifest_file="${2:-${HOME}/.terraform-ai-installer/file-manifest.config}"

    # Require manifest file - no fallback
    if [[ ! -f "${manifest_file}" ]]; then
        write_error_message "Manifest file not found: ${manifest_file}"
        echo "Please run with -bootstrap first to set up the installer."
        return 1
    fi

    # Parse from manifest
    parse_manifest_section "${manifest_file}" "${section_name}"
}

# Export functions used by the installer and fileoperations modules
export -f get_manifest_config get_installer_config get_workspace_root parse_manifest_section get_manifest_files
