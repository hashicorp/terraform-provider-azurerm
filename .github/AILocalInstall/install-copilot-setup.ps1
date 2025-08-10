#!/usr/bin/env pwsh
<#
.SYNOPSIS
    Terraform AzureRM Provider AI Setup - Installation Script

.DESCRIPTION
    Installs GitHub Copilot instruction files and AI prompts for enhanced development
    of the Terraform AzureRM Provider. This script configures VS Code with specialized
    AI instructions, coding patterns, and best practices for Azure resource development.

.PARAMETER RepositoryPath
    Specifies the path to the terraform-provider-azurerm repository.
    If not provided, the script will attempt to auto-discover the repository location.

.PARAMETER Clean
    Removes all installed AI setup files and restores original VS Code settings from backups.

.PARAMETER AutoApprove
    Skips interactive approval prompts.

.PARAMETER Help
    Displays detailed help information and usage examples.

.EXAMPLE
    .\install-modular.ps1
    Auto-discovers the repository and installs AI setup with default settings.

.EXAMPLE
    .\install-modular.ps1 -RepositoryPath "C:\Projects\terraform-provider-azurerm"
    Installs AI setup using the specified repository path.

.EXAMPLE
    .\install-modular.ps1 -AutoApprove
    Installs AI setup without interactive prompts.

.EXAMPLE
    .\install-modular.ps1 -Clean
    Removes all AI setup and restores original VS Code settings.

.NOTES
    File Name      : install-modular.ps1
    Author         : HashiCorp Terraform AzureRM Provider Team
    Prerequisite   : PowerShell 5.1+, VS Code
    Version        : 2.0.0 (Modular)
    
    This modular version replaces the monolithic installer with a clean,
    maintainable module system while providing identical functionality.
#>

param(
    [string]$RepositoryPath,
    [switch]$Clean,
    [switch]$AutoApprove,
    [switch]$Help
)

# Set strict mode and error action
Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

# Get script directory
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ModulesDir = Join-Path $ScriptDir "modules\powershell"

# Import modules
$modules = @(
    "core-functions.psm1",
    "backup-management.psm1", 
    "installation-detection.psm1",
    "ai-installation.psm1",
    "cleanup.psm1"
)

foreach ($module in $modules) {
    $modulePath = Join-Path $ModulesDir $module
    if (Test-Path $modulePath) {
        Import-Module $modulePath -Force
    } else {
        Write-Error "Required module not found: $modulePath"
        exit 1
    }
}

# Show help if requested
if ($Help) {
    Write-Host "Terraform AzureRM Provider AI Setup (Modular)" -ForegroundColor Cyan
    Write-Host "=============================================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "USAGE:" -ForegroundColor Yellow
    Write-Host "  .\install-modular.ps1 [OPTIONS]"
    Write-Host ""
    Write-Host "OPTIONS:" -ForegroundColor Yellow
    Write-Host "  -RepositoryPath <path>    Path to terraform-provider-azurerm repository"
    Write-Host "  -Clean                    Remove all installed files and restore backups"
    Write-Host "  -AutoApprove              Skip interactive approval prompts"
    Write-Host "  -Help                     Show this help message"
    Write-Host ""
    Write-Host "EXAMPLES:" -ForegroundColor Yellow
    Write-Host "  .\install-modular.ps1                                  # Auto-discover repository"
    Write-Host "  .\install-modular.ps1 -RepositoryPath C:\path\to\repo  # Use specific path"
    Write-Host "  .\install-modular.ps1 -AutoApprove                     # Non-interactive install"
    Write-Host "  .\install-modular.ps1 -Clean                           # Remove installation"
    Write-Host ""
    Write-Host "AI FEATURES:" -ForegroundColor Green
    Write-Host "  - Context-aware code generation and review"
    Write-Host "  - Azure-specific implementation patterns"
    Write-Host "  - Testing guidelines and best practices"
    Write-Host "  - Documentation standards enforcement"
    Write-Host "  - Error handling and debugging assistance"
    Write-Host ""
    Write-Host "MODULAR ARCHITECTURE:" -ForegroundColor Green
    Write-Host "  - Clean separation of concerns"
    Write-Host "  - Maintainable and testable modules"
    Write-Host "  - Enhanced error handling and logging"
    Write-Host "  - Improved backup and restore functionality"
    Write-Host ""
    exit 0
}

# Main execution starts here
Write-Host "Terraform AzureRM Provider AI Setup (Modular)" -ForegroundColor Cyan
Write-Host "=============================================" -ForegroundColor Cyan
Write-Host ""

try {
    # Check prerequisites
    Write-StatusMessage "Checking prerequisites..." "Info"
    if (-not (Test-Prerequisites)) {
        Write-StatusMessage "Prerequisites not met. Please install the required software." "Error"
        exit 1
    }
    Write-StatusMessage "Prerequisites check passed" "Success"

    # Discover repository path
    if (-not $RepositoryPath) {
        Write-StatusMessage "Auto-discovering repository location..." "Info"
        $RepositoryPath = Find-RepositoryRoot
        
        if (-not $RepositoryPath) {
            Write-StatusMessage "Could not find terraform-provider-azurerm repository" "Error"
            Write-Host "Please run this script from within the repository or specify -RepositoryPath" -ForegroundColor Yellow
            exit 1
        }
        
        Write-StatusMessage "Found repository at: $RepositoryPath" "Success"
    } else {
        if (-not (Test-Path $RepositoryPath)) {
            Write-StatusMessage "Specified repository path does not exist: $RepositoryPath" "Error"
            exit 1
        }
        Write-StatusMessage "Using specified repository path: $RepositoryPath" "Info"
    }

    if ($Clean) {
        # Cleanup mode
        Write-StatusMessage "Running in cleanup mode..." "Info"
        
        $installationState = Test-PreviousInstallation -RepositoryPath $RepositoryPath
        
        if (-not ($installationState.HasVSCodeSettings -or $installationState.HasBackups)) {
            Write-StatusMessage "No AI agent installation found to clean up" "Info"
            exit 0
        }
        
        if (-not $AutoApprove) {
            Write-Host ""
            Write-Host "This will remove the AI agent installation and restore VS Code settings from backup." -ForegroundColor Yellow
            $confirm = Read-Host "Continue? (y/N)"
            if ($confirm -notmatch "^[Yy]") {
                Write-StatusMessage "Cleanup cancelled by user" "Info"
                exit 0
            }
        }
        
        $success = Remove-AIAgent -RepositoryPath $RepositoryPath
        Show-CleanupSummary -RepositoryPath $RepositoryPath
        
        if ($success) {
            Write-StatusMessage "Cleanup completed successfully!" "Success"
            exit 0
        } else {
            Write-StatusMessage "Cleanup completed with errors" "Warning"
            exit 1
        }
    } else {
        # Installation mode
        Write-StatusMessage "Running in installation mode..." "Info"
        
        # Check for previous installation
        $installationState = Test-PreviousInstallation -RepositoryPath $RepositoryPath
        
        if ($installationState.HasVSCodeSettings -and $installationState.HasInstructions) {
            # Check installation health
            $health = Test-InstallationHealth -RepositoryPath $RepositoryPath
            
            if ($health.IsHealthy) {
                Write-StatusMessage "AI agent is already installed and healthy" "Success"
                
                if (-not $AutoApprove) {
                    $reinstall = Read-Host "Reinstall anyway? (y/N)"
                    if ($reinstall -notmatch "^[Yy]") {
                        Write-StatusMessage "Installation skipped by user" "Info"
                        exit 0
                    }
                } else {
                    Write-StatusMessage "Skipping installation (already installed)" "Info"
                    exit 0
                }
            } else {
                Write-StatusMessage "Previous installation found but has issues:" "Warning"
                foreach ($issue in $health.Issues) {
                    Write-StatusMessage "  - $issue" "Warning"
                }
                Write-StatusMessage "Will attempt to repair installation..." "Info"
            }
        } elseif ($installationState.IsPartialInstall) {
            Write-StatusMessage "Partial installation detected, completing installation..." "Info"
        }
        
        # Confirm installation
        if (-not $AutoApprove) {
            Write-Host ""
            Write-Host "This will install AI-enhanced development setup for the Terraform AzureRM Provider." -ForegroundColor Yellow
            Write-Host "Your existing VS Code settings will be backed up automatically." -ForegroundColor Yellow
            $confirm = Read-Host "Continue? (Y/n)"
            if ($confirm -match "^[Nn]") {
                Write-StatusMessage "Installation cancelled by user" "Info"
                exit 0
            }
        }
        
        # Perform installation
        $success = Install-AIAgent -RepositoryPath $RepositoryPath -Force:$AutoApprove
        
        if ($success) {
            Write-Host ""
            Write-Host "🎉 Installation completed successfully!" -ForegroundColor Green
            Write-Host ""
            Write-Host "Next steps:" -ForegroundColor Cyan
            Write-Host "1. Restart VS Code to activate the new settings" -ForegroundColor White
            Write-Host "2. Ensure GitHub Copilot extension is installed and authenticated" -ForegroundColor White
            Write-Host "3. Start developing with AI-enhanced assistance!" -ForegroundColor White
            Write-Host ""
            Write-Host "The AI agent is now configured with:" -ForegroundColor Cyan
            Write-Host "- Azure-specific implementation patterns" -ForegroundColor White
            Write-Host "- Testing guidelines and best practices" -ForegroundColor White
            Write-Host "- Documentation standards enforcement" -ForegroundColor White
            Write-Host "- Enhanced error handling assistance" -ForegroundColor White
            Write-Host ""
            exit 0
        } else {
            Write-StatusMessage "Installation failed" "Error"
            exit 1
        }
    }
} catch {
    Write-StatusMessage "Unexpected error: $_" "Error"
    Write-StatusMessage "Stack trace: $($_.ScriptStackTrace)" "Error"
    exit 1
}
