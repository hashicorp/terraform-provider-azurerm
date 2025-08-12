#!/bin/bash

# Cleanup script for AI installation test temporary files (Bash version)
# Removes all temporary files and directories created during AI installation testing

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Colo

# Default values
FORCE=false

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -f|--force)
            FORCE=true
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Cleanup temporary files created during AI installation testing"
            echo ""
            echo "OPTIONS:"
            echo "  -f, --force    Skip confirmation prompts"
            echo "  -h, --help     Show this help message"
            echo ""
            echo "EXAMPLES:"
            echo "  $0              # Interactive cleanup"
            echo "  $0 --force      # Automatic cleanup"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

echo -e "${CYAN}AI Installation Test Cleanup Utility (Bash)${NC}"
echo -e "${CYAN}===========================================${NC}"
echo ""

# Function to clean bash temporary files
cleanup_bash_temps() {
    echo -e "${YELLOW}Checking bash temporary files...${NC}"
    
    # Find all terraform-azurerm test directories
    TEMP_DIRS=$(find /tmp -maxdepth 1 -name "*terraform-azurerm*" -type d 2>/dev/null || true)
    
    if [[ -z "$TEMP_DIRS" ]]; then
        echo -e "${GREEN}  No bash temporary directories found${NC}"
        return
    fi
    
    TEMP_COUNT=$(echo "$TEMP_DIRS" | wc -l)
    echo -e "${YELLOW}  Found $TEMP_COUNT temporary directories${NC}"
    
    if [[ "$FORCE" != true ]]; then
        read -p "Remove $TEMP_COUNT bash temporary directories? (y/N): " -n 1 -
        echo ""
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo -e "${YELLOW}  Skipped bash cleanup${NC}"
            return
        fi
    fi
    
    # Remove directories
    REMOVED=0
    while IFS= read -r dir; do
        if [[ -n "$dir" ]]; then
            if rm -rf "$dir" 2>/dev/null; then
                echo -e "${GREEN}  Removed: $(basename "$dir")${NC}"
                ((REMOVED++))
            else
                echo -e "${RED}  Failed to remove: $(basename "$dir")${NC}"
            fi
        fi
    done <<< "$TEMP_DIRS"
    
    echo -e "${GREEN}  Successfully removed $REMOVED/$TEMP_COUNT directories${NC}"
}

# Function to clean Windows temp files (if running in WSL)
cleanup_windows_temps() {
    if [[ ! -d "/mnt/c" ]]; then
        echo -e "${YELLOW}Not running in WSL, skipping Windows cleanup${NC}"
        return
    fi
    
    echo -e "${YELLOW}Checking Windows temporary files...${NC}"
    
    # Check Windows temp directory
    WIN_TEMP="/mnt/c/Users/$(whoami)/AppData/Local/Temp"
    if [[ ! -d "$WIN_TEMP" ]]; then
        # Try alternative path
        WIN_TEMP="/mnt/c/Users/*/AppData/Local/Temp"
        if ! ls $WIN_TEMP/AzureRM* 2>/dev/null >/dev/null; then
            echo -e "${YELLOW}  Could not access Windows temp directory${NC}"
            return
        fi
    fi
    
    # Count Windows temp directories
    WIN_TEMP_DIRS=$(find $WIN_TEMP -maxdepth 1 -name "*AzureRM*" -type d 2>/dev/null || true)
    
    if [[ -z "$WIN_TEMP_DIRS" ]]; then
        echo -e "${GREEN}  No Windows temporary directories found${NC}"
        return
    fi
    
    WIN_TEMP_COUNT=$(echo "$WIN_TEMP_DIRS" | wc -l)
    echo -e "${YELLOW}  Found $WIN_TEMP_COUNT Windows temporary directories${NC}"
    
    if [[ "$FORCE" != true ]]; then
        read -p "Remove $WIN_TEMP_COUNT Windows temporary directories? (y/N): " -n 1 -
        echo ""
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo -e "${YELLOW}  Skipped Windows cleanup${NC}"
            return
        fi
    fi
    
    # Remove Windows directories
    WIN_REMOVED=0
    while IFS= read -r dir; do
        if [[ -n "$dir" ]]; then
            if rm -rf "$dir" 2>/dev/null; then
                echo -e "${GREEN}  Removed: $(basename "$dir")${NC}"
                ((WIN_REMOVED++))
            else
                echo -e "${RED}  Failed to remove: $(basename "$dir")${NC}"
            fi
        fi
    done <<< "$WIN_TEMP_DIRS"
    
    echo -e "${GREEN}  Successfully removed $WIN_REMOVED/$WIN_TEMP_COUNT Windows directories${NC}"
}

# Main execution
main() {
    cleanup_bash_temps
    echo ""
    
    cleanup_windows_temps
    echo ""
    
    echo -e "${GREEN}Cleanup completed successfully!${NC}"
    echo ""
    echo -e "${CYAN}To run this script automatically in the future:${NC}"
    echo -e "${NC}  $0 --force${NC}"
}

# Run main function
main "$@"
