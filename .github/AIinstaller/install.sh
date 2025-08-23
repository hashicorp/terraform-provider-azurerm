#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Universal AI Infrastructure Installer for Terraform AzureRM Provider
# Cross-platform launcher that detects the platform and runs the appropriate installer

set -euo pipefail

# Function to detect platform
detect_platform() {
    case "$(uname -s)" in
        Darwin*)
            echo "macos"
            ;;
        Linux*)
            echo "linux"
            ;;
        CYGWIN*|MINGW*|MSYS*)
            echo "windows"
            ;;
        *)
            echo "unknown"
            ;;
    esac
}

# Function to check if PowerShell is available
has_powershell() {
    command -v pwsh >/dev/null 2>&1 || command -v powershell >/dev/null 2>&1
}

# Function to show usage
show_usage() {
    cat << 'EOF'
Universal AI Infrastructure Installer for Terraform AzureRM Provider

This script automatically detects your platform and runs the appropriate installer:
- macOS/Linux: Uses bash installer (install-copilot-setup.sh)
- Windows: Uses PowerShell installer (install-copilot-setup.ps1)

USAGE:
    ./install.sh [OPTIONS]

The script will pass all arguments to the platform-specific installer.

PLATFORM-SPECIFIC USAGE:

macOS/Linux:
    ./install.sh -bootstrap                         # Bootstrap installer
    ./install.sh -repo-directory "/path/to/repo"    # Install AI infrastructure
    ./install.sh -verify                            # Verify installation
    ./install.sh -clean                             # Clean installation

Windows (with PowerShell available):
    ./install.sh -Bootstrap                         # Bootstrap installer
    ./install.sh -RepoDirectory "C:\path\to\repo"   # Install AI infrastructure
    ./install.sh -Verify                            # Verify installation
    ./install.sh -Clean                             # Clean installation

OPTIONS:
    All options are passed through to the platform-specific installer.
    See the respective installer help for detailed options:
    
    - For bash installer: ./install-copilot-setup.sh -help
    - For PowerShell installer: ./install-copilot-setup.ps1 -Help

EOF
}

# Main function
main() {
    local platform
    platform="$(detect_platform)"
    
    local script_dir
    script_dir="$(dirname "$(realpath "${0}")")"
    
    # Check for help flag
    for arg in "$@"; do
        case "${arg}" in
            -help|-Help)
                show_usage
                exit 0
                ;;
        esac
    done
    
    echo "Detected platform: ${platform}"
    echo ""
    
    case "${platform}" in
        macos|linux)
            local bash_installer="${script_dir}/install-copilot-setup.sh"
            
            if [[ ! -f "${bash_installer}" ]]; then
                echo "ERROR: Bash installer not found: ${bash_installer}"
                echo ""
                echo "The bash installer script is missing. Please ensure you have the complete"
                echo "AI installer package with both platform installers."
                exit 1
            fi
            
            echo "Running bash installer for ${platform}..."
            echo "   Command: ${bash_installer} $*"
            echo ""
            
            # Make executable if needed
            chmod +x "${bash_installer}"
            
            # Execute bash installer with all arguments
            exec "${bash_installer}" "$@"
            ;;
            
        windows)
            local ps_installer="${script_dir}/install-copilot-setup.ps1"
            
            if [[ ! -f "${ps_installer}" ]]; then
                echo "ERROR: PowerShell installer not found: ${ps_installer}"
                echo ""
                echo "The PowerShell installer script is missing. Please ensure you have the complete"
                echo "AI installer package with both platform installers."
                exit 1
            fi
            
            echo "Running PowerShell installer for Windows..."
            
            # Check if PowerShell is available
            if has_powershell; then
                # Prefer pwsh (PowerShell Core) over legacy PowerShell
                if command -v pwsh >/dev/null 2>&1; then
                    echo "   Command: pwsh -File \"${ps_installer}\" $*"
                    echo ""
                    exec pwsh -File "${ps_installer}" "$@"
                else
                    echo "   Command: powershell -File \"${ps_installer}\" $*"
                    echo ""
                    exec powershell -File "${ps_installer}" "$@"
                fi
            else
                echo "ERROR: PowerShell not found on this Windows system."
                echo ""
                echo "Please install PowerShell to use the AI installer on Windows:"
                echo "  - Install PowerShell Core: https://github.com/PowerShell/PowerShell"
                echo "  - Or use Windows PowerShell (usually pre-installed)"
                echo ""
                echo "Alternative: You can try running the PowerShell script directly:"
                echo "  powershell -File \"${ps_installer}\" $*"
                exit 1
            fi
            ;;
            
        unknown)
            echo "ERROR: Unsupported platform: $(uname -s)"
            echo ""
            echo "The universal installer supports:"
            echo "  - macOS (Darwin)"
            echo "  - Linux"
            echo "  - Windows (with PowerShell)"
            echo ""
            echo "Your platform '$(uname -s)' is not currently supported."
            echo "You may be able to use the bash installer directly if your platform"
            echo "supports bash and standard Unix tools:"
            echo ""
            echo "  ./install-copilot-setup.sh --help"
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"
