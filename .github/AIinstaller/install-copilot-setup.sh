#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Main AI Infrastructure Installer for Terraform AzureRM Provider (macOS/Linux)
# Version: 1.0.0
# Description: Interactive installer for AI-powered development infrastructure

set -euo pipefail

# Global variables
VERSION="1.0.0"
BRANCH="exp/terraform_copilot"
SOURCE_REPOSITORY="https://raw.githubusercontent.com/hashicorp/terraform-provider-azurerm"

# Load modules
SCRIPT_DIR="$(dirname "$(realpath "${0}")")"

# Load all bash modules
for module in "${SCRIPT_DIR}/modules/bash"/*.sh; do
    if [[ -f "${module}" ]]; then
        source "${module}"
    fi
done

# Verify essential modules are loaded
if ! command -v show_usage >/dev/null 2>&1 || ! command -v write_section >/dev/null 2>&1 || ! command -v get_user_profile >/dev/null 2>&1 || ! command -v copy_file >/dev/null 2>&1; then
    echo -e "\033[0;31m[ERROR]\033[0m Required modules are missing from ${SCRIPT_DIR}/modules/bash/"
    echo "Please ensure the following modules exist:"
    echo "  - ui.sh"
    echo "  - configparser.sh" 
    echo "  - fileoperations.sh"
    echo ""
    echo "If running from user profile, run bootstrap first:"
    echo "  ./install-copilot-setup.sh -bootstrap"
    exit 1
fi

# Command line options

# Command line options
BOOTSTRAP=false
REPO_DIRECTORY=""
AUTO_APPROVE=false
DRY_RUN=false
VERIFY=false
CLEAN=false
HELP=false

# Function to bootstrap installer
bootstrap_installer() {
    write_section "Bootstrap - Copying Installer to User Profile"
    
    local user_profile
    user_profile="$(get_user_profile)"
    
    if [[ -d "${user_profile}" ]]; then
        write_info "Using existing directory: ${user_profile}"
    else
        write_info "Creating directory: ${user_profile}"
        if [[ "${DRY_RUN}" != "true" ]]; then
            mkdir -p "${user_profile}"
        fi
    fi
    
    write_info "Copying installer files from local source repository..."
    echo ""
    
    local script_dir
    script_dir="$(dirname "$(realpath "${0}")")"
    
    # Statistics tracking
    local files_copied=0
    local files_failed=0
    local total_size=0
    
    # Files to bootstrap
    local bootstrap_files=(
        "file-manifest.config"
        "install-copilot-setup.sh"
    )
    
    # Copy each file individually
    for file in "${bootstrap_files[@]}"; do
        local source_path="${script_dir}/${file}"
        local target_path="${user_profile}/${file}"
        local filename=$(basename "${file}")
        
        if copy_file "${source_path}" "${target_path}" "${filename}"; then
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
    
    # Copy bash modules if they exist
    if [[ -d "${script_dir}/modules/bash" ]]; then
        # Create modules directory
        if [[ "${DRY_RUN}" != "true" ]]; then
            mkdir -p "${user_profile}/modules/bash"
        fi
        
        # Copy each module file individually
        for module_file in "${script_dir}/modules/bash"/*.sh; do
            if [[ -f "${module_file}" ]]; then
                local module_filename=$(basename "${module_file}")
                local target_module="${user_profile}/modules/bash/${module_filename}"
                
                if copy_file "${module_file}" "${target_module}" "${module_filename}"; then
                    files_copied=$((files_copied + 1))
                    # Calculate file size
                    if [[ -f "${target_module}" ]]; then
                        local file_size
                        file_size=$(get_file_size "${target_module}")
                        total_size=$((total_size + file_size))
                    fi
                else
                    files_failed=$((files_failed + 1))
                fi
            fi
        done
    fi
    
    echo ""
    
    if [[ "${files_failed}" -eq 0 ]]; then
        # Calculate total size in KB (simple integer division)
        local total_size_kb
        total_size_kb=$((total_size / 1024))
        
        write_success "Bootstrap completed successfully!"
        echo ""
        echo -e "  Files copied: ${GREEN}${files_copied}${NC}"
        echo -e "  Total size: ${GREEN}${total_size_kb} KB${NC}"
        echo -e "  Location: ${GREEN}${user_profile}${NC}"
        echo ""
        echo -e "${BLUE}NEXT STEPS:${NC}"
        echo -e "  ${BOLD}1. Switch to your feature branch:${NC}"
        echo -e "     ${YELLOW}git checkout feature/your-branch-name${NC}"
        echo ""
        echo -e "  ${BOLD}2. Run the installer from your user profile:${NC}"
        echo -e "     ${YELLOW}${user_profile}/install-copilot-setup.sh -repo-directory \"$(get_workspace_root)\"${NC}"
        echo ""
        echo -e "  ${YELLOW}Note: The -repo-directory parameter tells the installer where to find the git repository${NC}"
        echo -e "  ${YELLOW}      for branch detection when running from your user profile.${NC}"
        echo ""
    else
        write_error "Bootstrap failed with ${files_failed} file(s) failing to copy"
        return 1
    fi
}

# Function to install AI infrastructure
install_infrastructure() {
    local workspace_root="$1"
    
    log_section "Installing AI Infrastructure"
    
    # Read manifest configuration
    local manifest_file="${HOME}/.terraform-ai-installer/file-manifest.config"
    if [[ ! -f "${manifest_file}" ]]; then
        log_error "Manifest file not found: ${manifest_file}"
        echo "Please run with --bootstrap first to set up the installer."
        exit 1
    fi
    
    log_info "Installing to workspace: ${workspace_root}"
    
    # List of instruction files (this should ideally be read from manifest)
    local instruction_files=(
        "api-evolution-patterns.instructions.md"
        "azure-patterns.instructions.md"
        "code-clarity-enforcement.instructions.md"
        "documentation-guidelines.instructions.md"
        "error-patterns.instructions.md"
        "implementation-guide.instructions.md"
        "migration-guide.instructions.md"
        "performance-optimization.instructions.md"
        "provider-guidelines.instructions.md"
        "schema-patterns.instructions.md"
        "security-compliance.instructions.md"
        "testing-guidelines.instructions.md"
        "troubleshooting-decision-trees.instructions.md"
    )
    
    # Calculate total files to install (main copilot instructions + instruction files)
    local total_files=$((1 + ${#instruction_files[@]}))
    local current_file=0
    
    # Install main copilot instructions
    current_file=$((current_file + 1))
    show_completion ${current_file} ${total_files} "Installing main copilot instructions"
    download_file ".github/copilot-instructions.md" "${workspace_root}/.github/copilot-instructions.md" "Main Copilot instructions"
    
    # Install instruction files
    local instructions_dir="${workspace_root}/.github/instructions"
    mkdir -p "${instructions_dir}"
    
    for file in "${instruction_files[@]}"; do
        current_file=$((current_file + 1))
        show_completion ${current_file} ${total_files} "Installing ${file}"
        download_file ".github/instructions/${file}" "${instructions_dir}/${file}" "instructions/${file}"
    done
    
    log_success "AI Infrastructure installation completed!"
}

# Function to verify installation
verify_installation() {
    local workspace_root
    workspace_root="$(get_workspace_root)"
    
    log_section "Verifying AI Infrastructure"
    
    local files_to_check=(
        "${workspace_root}/.github/copilot-instructions.md"
        "${workspace_root}/.github/instructions"
    )
    
    local all_good=true
    
    for file in "${files_to_check[@]}"; do
        if [[ -e "${file}" ]]; then
            echo "  [OK] ${file}"
        else
            echo "  [MISSING] ${file}"
            all_good=false
        fi
    done
    
    if [[ "${all_good}" == "true" ]]; then
        log_success "All AI infrastructure files are present"
    else
        log_warning "Some AI infrastructure files are missing"
        echo "Run the installer to restore missing files."
    fi
}

# Function to clean installation
clean_installation() {
    local workspace_root
    workspace_root="$(get_workspace_root)"
    
    log_section "Cleaning AI Infrastructure"
    
    local files_to_remove=(
        "${workspace_root}/.github/copilot-instructions.md"
        "${workspace_root}/.github/instructions"
    )
    
    for file in "${files_to_remove[@]}"; do
        if [[ -e "${file}" ]]; then
            if [[ "${DRY_RUN}" == "true" ]]; then
                echo "  [DRY-RUN] Would remove: ${file}"
            else
                rm -rf "${file}"
                echo "  Removed: ${file}"
            fi
        fi
    done
    
    log_success "AI Infrastructure cleanup completed!"
}

# Parse command line arguments
parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -bootstrap)
                BOOTSTRAP=true
                shift
                ;;
            -repo-directory)
                REPO_DIRECTORY="$2"
                shift 2
                ;;
            -auto-approve)
                AUTO_APPROVE=true
                shift
                ;;
            -dry-run)
                DRY_RUN=true
                shift
                ;;
            -verify)
                VERIFY=true
                shift
                ;;
            -clean)
                CLEAN=true
                shift
                ;;
            -help)
                HELP=true
                shift
                ;;
            *)
                log_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done
}

# Main function
main() {
    parse_arguments "$@"
    
    if [[ "${HELP}" == "true" ]]; then
        show_usage
        exit 0
    fi
    
    if [[ "${BOOTSTRAP}" == "true" ]]; then
        bootstrap_installer
        exit 0
    fi
    
    if [[ "${VERIFY}" == "true" ]]; then
        verify_installation
        exit 0
    fi
    
    if [[ "${CLEAN}" == "true" ]]; then
        clean_installation
        exit 0
    fi
    
    # Default: install infrastructure
    local workspace_root
    if [[ -n "${REPO_DIRECTORY}" ]]; then
        validate_repository "${REPO_DIRECTORY}"
        workspace_root="${REPO_DIRECTORY}"
    else
        workspace_root="$(get_workspace_root)"
    fi
    
    install_infrastructure "${workspace_root}"
}

# Run main function with all arguments
main "$@"
