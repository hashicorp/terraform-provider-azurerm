#!/usr/bin/env bash
# ValidationEngine Module for Terraform AzureRM Provider AI Setup (Bash)
# Handles comprehensive validation, dependency checking, and system requirements
# STREAMLINED VERSION - Contains only functions actually used by main script and dependencies

# Private Functions

# Function to find workspace root by looking for go.mod file
find_workspace_root() {
    local start_path="$1"
    local current_path="${start_path}"
    local max_depth=10
    local depth=0

    while [[ ${depth} -lt ${max_depth} && -n "${current_path}" ]]; do
        local go_mod_path="${current_path}/go.mod"
        if [[ -f "${go_mod_path}" ]]; then
            echo "${current_path}"
            return 0
        fi

        # Move to parent directory
        local parent_path=$(dirname "${current_path}")
        if [[ "${parent_path}" == "${current_path}" ]]; then
            # Reached root directory
            break
        fi
        current_path="${parent_path}"
        ((depth++))
    done

    return 1
}

# Function to test bash version (equivalent to PowerShell version test)
test_bash_version() {
    local bash_version="${BASH_VERSION}"
    local version_major=$(echo "${bash_version}" | cut -d. -f1)

    if [[ ${version_major} -ge 4 ]]; then
        echo "Valid=true"
        echo "Version=${bash_version}"
        echo "Reason=Bash version ${bash_version} meets requirements"
    else
        echo "Valid=false"
        echo "Version=${bash_version}"
        echo "Reason=Bash version ${bash_version} is too old. Minimum version 4.0 required"
    fi
}

# Function to test required commands
test_required_commands() {
    local required_commands=("git" "curl" "mkdir" "cp" "rm" "dirname" "realpath")
    local missing_commands=()
    local valid=true

    for cmd in "${required_commands[@]}"; do
        if ! command -v "${cmd}" >/dev/null 2>&1; then
            missing_commands+=("${cmd}")
            valid=false
        fi
    done

    echo "Valid=${valid}"
    if [[ ${valid} == "true" ]]; then
        echo "Reason=All required commands available"
    else
        echo "Reason=Missing commands: ${missing_commands[*]}"
        echo "MissingCommands=${missing_commands[*]}"
    fi
}

# Function to validate repository structure (pure validation logic)
validate_repository() {
    local repo_dir="$1"

    # Check if directory exists
    if [[ ! -d "${repo_dir}" ]]; then
        return 1
    fi

    # Check for terraform-provider-azurerm repository (comprehensive validation)
    local has_go_mod=false
    local has_azurerm_content=false
    local has_main_go=false
    local has_services_dir=false

    # Check for go.mod and its content
    if [[ -f "${repo_dir}/go.mod" ]]; then
        has_go_mod=true
        if grep -q "terraform-provider-azurerm" "${repo_dir}/go.mod" 2>/dev/null; then
            has_azurerm_content=true
        fi
    fi

    # Check for expected structure
    if [[ -f "${repo_dir}/main.go" ]]; then
        has_main_go=true
    fi

    if [[ -d "${repo_dir}/internal/services" ]]; then
        has_services_dir=true
    fi

    # Valid only if we have ALL expected characteristics
    if [[ "${has_go_mod}" == "true" && "${has_azurerm_content}" == "true" &&
          "${has_main_go}" == "true" && "${has_services_dir}" == "true" ]]; then
        return 0
    fi

    return 1
}

# Function to test system requirements (comprehensive version)
test_system_requirements() {
    local bash_result=$(test_bash_version)
    local commands_result=$(test_required_commands)
    local internet_result=""

    # Test internet connectivity
    if test_internet_connectivity; then
        internet_result="Connected=true"$'\n'"Reason=Internet connectivity verified"
    else
        internet_result="Connected=false"$'\n'"Reason=No internet connectivity detected. Check network connection and firewall settings."
    fi

    # Parse results
    local bash_valid=$(echo "${bash_result}" | grep "Valid=" | cut -d= -f2)
    local commands_valid=$(echo "${commands_result}" | grep "Valid=" | cut -d= -f2)
    local internet_connected=$(echo "${internet_result}" | grep "Connected=" | cut -d= -f2)

    # Overall validation
    if [[ "${bash_valid}" == "true" && "${commands_valid}" == "true" && "${internet_connected}" == "true" ]]; then
        echo "OverallValid=true"
    else
        echo "OverallValid=false"
    fi

    echo "Bash=${bash_result}"
    echo "Commands=${commands_result}"
    echo "Internet=${internet_result}"
}

# Function to test system requirements (original version for compatibility)
test_system_requirements_basic() {
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
        write_error_message "Missing required system tools: ${missing_tools[*]}"
        return 1
    fi

    return 0
}

# Function to test Git repository with branch safety checks
test_git_repository() {
    local repo_dir="$1"
    local allow_bootstrap_on_source="${2:-false}"

    # Initialize result variables
    local valid=false
    local is_git_repo=false
    local has_remote=false
    local current_branch="Unknown"
    local is_source_branch=false
    local reason=""

    if [[ ! -d "${repo_dir}/.git" ]]; then
        reason="Not a Git repository: ${repo_dir}"
            write_warning_message "${reason}"
    else
        is_git_repo=true

        # Check if git command is available
        if ! command -v git >/dev/null 2>&1; then
            reason="Git command not available"
            write_warning_message "${reason}"
        else
            # Get current branch
            current_branch=$(cd "${repo_dir}" && git branch --show-current 2>/dev/null || echo "Unknown")

            # Check for remote
            if cd "${repo_dir}" && git remote -v >/dev/null 2>&1; then
                has_remote=true
            fi

            # Check if on source branch (main, master, or exp/terraform_copilot)
            case "${current_branch}" in
                "main"|"master"|"exp/terraform_copilot")
                    is_source_branch=true
                    ;;
            esac

            # Validate based on branch safety rules
            if [[ "${is_source_branch}" == "true" ]] && [[ "${allow_bootstrap_on_source}" != "true" ]]; then
                valid=false
                reason="Cannot install on source branch '${current_branch}' without explicit permission. Use feature branch for safety."
            else
                valid=true
                reason="Git repository validation passed"
            fi
        fi
    fi

    # Output results in a structured format
    echo "Valid=${valid}"
    echo "IsGitRepo=${is_git_repo}"
    echo "HasRemote=${has_remote}"
    echo "CurrentBranch=${current_branch}"
    echo "IsSourceBranch=${is_source_branch}"
    echo "Reason=${reason}"

    [[ "${valid}" == "true" ]]
}

# Function to test workspace validity
test_workspace_valid() {
    local workspace_path="${1:-$(pwd)}"

    # Initialize result variables
    local valid=false
    local is_terraform_provider=false
    local is_azurerm_provider=false
    local workspace_root=""
    local reason=""

    # Find workspace root
    workspace_root=$(find_workspace_root "${workspace_path}")

    if [[ -n "${workspace_root}" ]]; then
        # Check if it's a Terraform provider
        if [[ -f "${workspace_root}/main.go" ]] && [[ -d "${workspace_root}/internal" ]]; then
            is_terraform_provider=true

            # Check if it's specifically the AzureRM provider
            if grep -q "terraform-provider-azurerm" "${workspace_root}/go.mod" 2>/dev/null; then
                is_azurerm_provider=true
                valid=true
                reason="Valid AzureRM provider workspace detected"
            else
                reason="Terraform provider detected but not AzureRM provider"
            fi
        else
            reason="Directory contains go.mod but is not a Terraform provider"
        fi
    else
        reason="No workspace root found (missing go.mod)"
    fi

    # Output results in structured format
    echo "Valid=${valid}"
    echo "IsTerraformProvider=${is_terraform_provider}"
    echo "IsAzureRMProvider=${is_azurerm_provider}"
    echo "WorkspaceRoot=${workspace_root}"
    echo "Reason=${reason}"

    [[ "${valid}" == "true" ]]
}

# Function to run comprehensive pre-installation validation
test_pre_installation() {
    local allow_bootstrap_on_source="${1:-false}"
    local workspace_path="${2:-$(pwd)}"

    # Initialize results
    local overall_valid=true
    local timestamp=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

    # Get workspace root for git operations
    local workspace_root=$(find_workspace_root "${workspace_path}")

    # Test Git repository (CRITICAL: check first for branch safety)
    local git_result=""
    if [[ -n "${workspace_root}" ]]; then
        git_result=$(test_git_repository "${workspace_root}" "${allow_bootstrap_on_source}")
    else
        git_result=$(test_git_repository "${workspace_path}" "${allow_bootstrap_on_source}")
    fi

    local git_valid=$(echo "${git_result}" | grep "Valid=" | cut -d= -f2)
    if [[ "${git_valid}" != "true" ]]; then
        overall_valid=false
    fi

    # Test workspace validity
    local workspace_result=$(test_workspace_valid "${workspace_path}")
    local workspace_valid=$(echo "${workspace_result}" | grep "Valid=" | cut -d= -f2)
    if [[ "${workspace_valid}" != "true" ]]; then
        overall_valid=false
    fi

    # Test system requirements
    local system_result=$(test_system_requirements)
    local system_valid=$(echo "${system_result}" | grep "OverallValid=" | cut -d= -f2)
    if [[ "${system_valid}" != "true" ]]; then
        overall_valid=false
    fi

    # Output comprehensive results
    echo "OverallValid=${overall_valid}"
    echo "Timestamp=${timestamp}"
    echo "Git=${git_result}"
    echo "Workspace=${workspace_result}"
    echo "SystemRequirements=${system_result}"

    [[ "${overall_valid}" == "true" ]]
}

# Public Functions

# Function to get workspace root (public wrapper for find_workspace_root)
get_workspace_root() {
    local start_path="${1:-$(pwd)}"
    find_workspace_root "${start_path}"
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
                write_section "Workspace Verification"
    else
        echo "============================================================"
        echo " Workspace Verification"
        echo "============================================================"
        echo ""
    fi

    local all_good=true
    local manifest_file="${HOME}/.terraform-ai-installer/file-manifest.config"
    local files_checked=0
    local files_passed=0
    local files_failed=0

    # Check main files
    local main_files
    main_files=$(get_manifest_files "MAIN_FILES" "${manifest_file}")
    if [[ $? -eq 0 && -n "${main_files}" ]]; then
        while IFS= read -r file; do
            [[ -z "${file}" ]] && continue
            local full_path="${workspace_root}/${file}"
            files_checked=$((files_checked + 1))
            if [[ -f "${full_path}" ]]; then
                write_green "  [FOUND] ${file}"
                files_passed=$((files_passed + 1))
            else
                write_red "  [MISSING] ${file}"
                files_failed=$((files_failed + 1))
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
        files_checked=$((files_checked + 1))
        if [[ -d "${instructions_dir}" ]]; then
            write_green "  [FOUND] .github/instructions/"
            files_passed=$((files_passed + 1))

            while IFS= read -r file; do
                [[ -z "${file}" ]] && continue
                local full_path="${workspace_root}/${file}"
                local filename=$(basename "${file}")
                files_checked=$((files_checked + 1))
                if [[ -f "${full_path}" ]]; then
                    write_green "    [FOUND] .github/instructions/${filename}"
                    files_passed=$((files_passed + 1))
                else
                    write_red "    [MISSING] .github/instructions/${filename}"
                    files_failed=$((files_failed + 1))
                    all_good=false
                fi
            done <<< "${instruction_files}"
        else
            write_red "  [MISSING] .github/instructions/"
            files_failed=$((files_failed + 1))
            all_good=false
        fi
    fi

    # Check prompt files
    local prompt_files
    prompt_files=$(get_manifest_files "PROMPT_FILES" "${manifest_file}")
    if [[ $? -eq 0 && -n "${prompt_files}" ]]; then
        # Check if prompts directory exists
        local prompts_dir="${workspace_root}/.github/prompts"
        files_checked=$((files_checked + 1))
        if [[ -d "${prompts_dir}" ]]; then
            write_green "  [FOUND] .github/prompts/"
            files_passed=$((files_passed + 1))

            while IFS= read -r file; do
                [[ -z "${file}" ]] && continue
                local full_path="${workspace_root}/${file}"
                local filename=$(basename "${file}")
                files_checked=$((files_checked + 1))
                if [[ -f "${full_path}" ]]; then
                    write_green "    [FOUND] .github/prompts/${filename}"
                    files_passed=$((files_passed + 1))
                else
                    write_red "    [MISSING] .github/prompts/${filename}"
                    files_failed=$((files_failed + 1))
                    all_good=false
                fi
            done <<< "${prompt_files}"
        else
            write_red "  [MISSING] .github/prompts/"
            files_failed=$((files_failed + 1))
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

            # Count directory first (like PowerShell does)
            files_checked=$((files_checked + 1))
            if [[ -d "${workspace_root}/${dir_path}" ]]; then
                write_green "  [FOUND] ${dir_path}/"
                files_passed=$((files_passed + 1))
            else
                write_red "  [MISSING] ${dir_path}/"
                files_failed=$((files_failed + 1))
                all_good=false
            fi

            # Count file separately (like PowerShell does)
            files_checked=$((files_checked + 1))
            if [[ -f "${full_path}" ]]; then
                write_green "    [FOUND] ${filename}"
                files_passed=$((files_passed + 1))
            else
                write_red "    [MISSING] ${filename}"
                files_failed=$((files_failed + 1))
                all_good=false
            fi
        done <<< "${universal_files}"
    fi

    echo ""
    if [[ "${all_good}" == "true" ]]; then
        # Show verification summary using dynamic show_operation_summary
        if declare -f show_operation_summary >/dev/null 2>&1; then
            # Get dynamic branch information
            local current_branch="unknown"
            local branch_type="feature"
            if command -v git >/dev/null 2>&1 && [[ -d "${workspace_root}/.git" ]]; then
                current_branch=$(cd "${workspace_root}" && git branch --show-current 2>/dev/null || echo "unknown")
                # Determine branch type
                local source_branches=("main" "master" "exp/terraform_copilot")
                for branch in "${source_branches[@]}"; do
                    if [[ "$current_branch" == "$branch" ]]; then
                        branch_type="source"
                        break
                    fi
                done
            fi
            local issues_found=$((files_checked - files_passed))

            # Match PowerShell order: Branch Type, Target Branch, Issues Found, Files Verified, Items Successful, Location
            show_operation_summary "Verification" "true" "false" \
                "Branch Type:${branch_type}" \
                "Target Branch:${current_branch}" \
                "Issues Found:${issues_found}" \
                "Files Verified:${files_checked}" \
                "Items Successful:${files_passed}" \
                "Location:${workspace_root}"
        fi
    else
        write_red "Some AI infrastructure files are missing!"

        # Show verification summary using dynamic show_operation_summary
        if declare -f show_operation_summary >/dev/null 2>&1; then
            # Get dynamic branch information
            local current_branch="unknown"
            local branch_type="feature"
            if command -v git >/dev/null 2>&1 && [[ -d "${workspace_root}/.git" ]]; then
                current_branch=$(cd "${workspace_root}" && git branch --show-current 2>/dev/null || echo "unknown")
                # Determine branch type
                local source_branches=("main" "master" "exp/terraform_copilot")
                for branch in "${source_branches[@]}"; do
                    if [[ "$current_branch" == "$branch" ]]; then
                        branch_type="source"
                        break
                    fi
                done
            fi
            local issues_found=$((files_checked - files_passed))

            show_operation_summary "Verification" "false" "false" \
                "Branch Type:${branch_type}" \
                "Target Branch:${current_branch}" \
                "Issues Found:${issues_found}" \
                "Files Verified:${files_checked}" \
                "Items Successful:${files_passed}" \
                "Location:${workspace_root}" \
                --next-steps \
                "Run installation if components are missing" \
                "Use -clean option to remove installation if needed"
        else
            echo "Run the installer to restore missing files."
        fi
    fi
}

# Export functions for use in other scripts
export -f test_system_requirements validate_repository
export -f test_git_repository test_workspace_valid test_internet_connectivity verify_installation
