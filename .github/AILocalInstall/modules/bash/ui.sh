#!/bin/bash
# User Interface Module  
# Handles user interaction, progress reporting, and output formatting using real patterns

show_welcome_banner() {
    # Displays the installation welcome banner
    local version="${1:-2.0.0}"
    
    echo ""
    echo -e "${CYAN}============================================================================${NC}"
    echo -e "${WHITE} Terraform AzureRM Provider - AI Setup Installer v$version${NC}"
    echo -e "${CYAN}============================================================================${NC}"
    echo ""
    echo -e "${GRAY} This installer will configure VS Code with:${NC}"
    echo -e "${GRAY}   - AI-powered coding instructions for AzureRM provider development${NC}"
    echo -e "${GRAY}   - GitHub Copilot prompt templates and examples${NC}"
    echo -e "${GRAY}   - Optimized VS Code settings for Terraform development${NC}"
    echo ""
}

show_installation_summary() {
    # Displays a comprehensive installation summary
    # Input format: key=value pairs separated by newlines
    local results="$1"
    
    echo ""
    echo -e "${CYAN}============================================================================${NC}"
    echo -e "${WHITE} Installation Summary${NC}"
    echo -e "${CYAN}============================================================================${NC}"
    echo ""
    
    # Parse results
    local success="" instruction_count=0 expected_instruction_count=0 prompt_count=0 expected_prompt_count=0
    local main_count=0 expected_main_count=0 settings_configured="false" error_count=0
    
    while IFS='=' read -r key value; do
        case "$key" in
            "SUCCESS") success="$value" ;;
            "INSTRUCTION_FILES_COUNT") instruction_count="$value" ;;
            "EXPECTED_INSTRUCTION_FILES_COUNT") expected_instruction_count="$value" ;;
            "PROMPT_FILES_COUNT") prompt_count="$value" ;;
            "EXPECTED_PROMPT_FILES_COUNT") expected_prompt_count="$value" ;;
            "MAIN_FILES_COUNT") main_count="$value" ;;
            "EXPECTED_MAIN_FILES_COUNT") expected_main_count="$value" ;;
            "SETTINGS_CONFIGURED") settings_configured="$value" ;;
            "ERROR_COUNT") error_count="$value" ;;
        esac
    done <<< "$results"
    
    if [[ "$success" == "true" ]]; then
        echo -e "${GREEN} Status: INSTALLATION SUCCESSFUL${NC}"
        echo ""
        echo -e "${WHITE} Components Installed:${NC}"
        
        # Show found vs expected counts with appropriate colors
        local instruction_color prompt_color main_color settings_color
        [[ "$instruction_count" -eq "$expected_instruction_count" ]] && instruction_color="$GREEN" || instruction_color="$YELLOW"
        [[ "$prompt_count" -eq "$expected_prompt_count" ]] && prompt_color="$GREEN" || prompt_color="$YELLOW"
        [[ "$main_count" -eq "$expected_main_count" ]] && main_color="$GREEN" || main_color="$YELLOW"
        [[ "$settings_configured" == "true" ]] && settings_color="$GREEN" || settings_color="$RED"
        
        echo -e "   - Instruction Files: ${instruction_color}$instruction_count of $expected_instruction_count files found${NC}"
        echo -e "   - Prompt Files: ${prompt_color}$prompt_count of $expected_prompt_count files found${NC}"
        echo -e "   - Main Files: ${main_color}$main_count of $expected_main_count files found${NC}"
        
        local settings_text
        [[ "$settings_configured" == "true" ]] && settings_text="Correctly configured" || settings_text="Not configured"
        echo -e "   - VS Code Settings: ${settings_color}$settings_text${NC}"
    else
        # Determine if this is no installation or partial installation
        local total_found=$((instruction_count + prompt_count + main_count))
        
        if [[ "$total_found" -eq 0 && "$settings_configured" != "true" ]]; then
            echo -e "${RED} Status: NO INSTALLATION FOUND${NC}"
        else
            echo -e "${YELLOW} Status: INSTALLATION COMPLETED WITH ISSUES${NC}"
        fi
        echo ""
        
        # Show partial installation status
        echo -e "${WHITE} Components Found:${NC}"
        
        # Show each component with appropriate color based on actual status
        local instruction_color prompt_color main_color settings_color
        [[ "$instruction_count" -eq "$expected_instruction_count" ]] && instruction_color="$GREEN" || instruction_color="$RED"
        [[ "$prompt_count" -eq "$expected_prompt_count" ]] && prompt_color="$GREEN" || prompt_color="$RED"
        [[ "$main_count" -eq "$expected_main_count" ]] && main_color="$GREEN" || main_color="$RED"
        [[ "$settings_configured" == "true" ]] && settings_color="$GREEN" || settings_color="$RED"
        
        echo -e "   - Instruction Files: ${instruction_color}$instruction_count of $expected_instruction_count files found${NC}"
        echo -e "   - Prompt Files: ${prompt_color}$prompt_count of $expected_prompt_count files found${NC}"
        echo -e "   - Main Files: ${main_color}$main_count of $expected_main_count files found${NC}"
        
        local settings_text
        [[ "$settings_configured" == "true" ]] && settings_text="Correctly configured" || settings_text="Not configured"
        echo -e "   - VS Code Settings: ${settings_color}$settings_text${NC}"
        echo ""
        
        echo -e "${RED} Errors encountered: $error_count${NC}"
    fi
    
    echo ""
}

show_repository_info() {
    # Displays repository validation information
    local repository_path="$1"
    local validation_results="$2"
    
    echo ""
    echo -e "${WHITE}Repository Information:${NC}"
    echo -e "${GRAY}  Path: $repository_path${NC}"
    
    # Parse validation results
    local is_valid
    while IFS='=' read -r key value; do
        case "$key" in
            "IS_VALID_REPOSITORY") is_valid="$value" ;;
        esac
    done <<< "$validation_results"
    
    if [[ "$is_valid" == "true" ]]; then
        echo -e "${GREEN}  Status: Valid Terraform AzureRM Provider repository${NC}"
    else
        echo -e "${RED}  Status: Invalid repository or missing files${NC}"
    fi
    echo ""
}

request_user_confirmation() {
    # Requests user confirmation for installation
    local message="${1:-Proceed with installation?}"
    local force="${2:-false}"
    
    if [[ "$force" == "true" ]]; then
        echo -e "${YELLOW}Force mode enabled - proceeding automatically${NC}"
        return 0
    fi
    
    echo ""
    echo -e "${WHITE}$message${NC}"
    echo -ne "${YELLOW}  [Y] Yes  [N] No  (Default: Y): ${NC}"
    
    read -r response
    if [[ -z "$response" || "$response" =~ ^[Yy]([Ee][Ss])?$ ]]; then
        return 0
    fi
    
    return 1
}

show_next_steps() {
    # Shows next steps after installation
    echo ""
    echo -e "${CYAN}============================================================================${NC}"
    echo -e "${WHITE} Next Steps${NC}"
    echo -e "${CYAN}============================================================================${NC}"
    echo ""
    echo -e "${WHITE} 1. Restart VS Code to ensure all settings take effect${NC}"
    echo -e "${WHITE} 2. Open the Terraform AzureRM Provider repository in VS Code${NC}"
    echo -e "${WHITE} 3. Try using GitHub Copilot with the new AI instructions!${NC}"
    echo ""
    echo -e "${GRAY} The AI is now configured to help you with:${NC}"
    echo -e "${GRAY}   - Writing resource implementations following provider patterns${NC}"
    echo -e "${GRAY}   - Creating comprehensive acceptance tests${NC}"
    echo -e "${GRAY}   - Generating proper documentation${NC}"
    echo -e "${GRAY}   - Following Azure SDK integration best practices${NC}"
    echo ""
    echo -e "${GREEN} Happy coding!${NC}"
    echo ""
}

write_installation_progress() {
    # Shows progress during installation (simplified - no progress bar)
    local stage="$1"
    local stage_number="$2"
    local total_stages="$3"
    
    # Calculate percentage and show status with percentage
    [[ "$total_stages" -eq 0 ]] && total_stages=1
    local percent=$((stage_number * 100 / total_stages))
    
    write_status_message "$stage - Step $stage_number of $total_stages ($percent%)" "Info"
}

show_help() {
    # Displays help information for the installer
    echo -e "${CYAN}Terraform AzureRM Provider AI Setup${NC}"
    echo -e "${CYAN}====================================================${NC}"
    echo ""
    echo -e "${YELLOW}USAGE:${NC}"
    echo "  ./install-copilot-setup.sh [OPTIONS]"
    echo ""
    echo -e "${YELLOW}OPTIONS:${NC}"
    echo "  -repository-path <path>       Path to terraform-provider-azurerm repository"
    echo "  -clean                        Remove all installed files and restore backups"
    echo "  -auto-approve                 Skip interactive approval prompts"
    echo "  -verify                       Run verification only without installing"
    echo "  -help                         Show this help message"
    echo ""
    echo -e "${YELLOW}EXAMPLES:${NC}"
    echo "  ./install-copilot-setup.sh                                       # Auto-discover repository"
    echo "  ./install-copilot-setup.sh -repository-path /path/to/repo        # Use specific path"
    echo "  ./install-copilot-setup.sh -auto-approve                         # Non-interactive install"
    echo "  ./install-copilot-setup.sh -clean                                # Remove installation"
    echo "  ./install-copilot-setup.sh -verify                               # Verify current installation"
    echo ""
    echo -e "${YELLOW}PLATFORM SUPPORT:${NC}"
    echo "  - macOS: ~/Library/Application Support/Code/User"
    echo "  - Linux: ~/.config/Code/User"
    echo "  - Windows (Git Bash/MSYS2): %APPDATA%/Code/User"
    echo ""
}

# Export functions by creating a list of available functions
__UI_FUNCTIONS="show_welcome_banner show_installation_summary show_repository_info request_user_confirmation show_next_steps write_installation_progress show_help"
