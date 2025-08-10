<#
.SYNOPSIS
    Terraform AzureRM Provider AI Setup - Modular Installation Script

.DESCRIPTION
    Installs GitHub Copilot instruction files and AI prompts for enhanced development
    of the Terraform AzureRM Provider. This modular version provides better maintainability
    and testing capabilities.

.PARAMETER RepositoryPath
    Specifies the path to the terraform-provider-azurerm repository.
    If not provided, the script will attempt to auto-discover the repository location.

.PARAMETER Clean
    Removes all installed AI setup files and restores original VS Code settings from backups.

.PARAMETER AutoApprove
    Skips interactive approval prompts.

.PARAMETER Help
    Displays detailed help information and usage examples.

.PARAMETER ShowBackups
    Shows information about existing backups without making changes.

.EXAMPLE
    .\main-installer.ps1
    Auto-discovers the repository and installs AI setup with default settings.

.EXAMPLE
    .\main-installer.ps1 -RepositoryPath "C:\Projects\terraform-provider-azurerm"
    Installs AI setup using the specified repository path.

.EXAMPLE
    .\main-installer.ps1 -Clean
    Removes all AI setup and restores original VS Code settings.

.EXAMPLE
    .\main-installer.ps1 -ShowBackups
    Displays information about existing backups.

.NOTES
    File Name      : main-installer.ps1
    Author         : HashiCorp Terraform AzureRM Provider Team
    Prerequisite   : PowerShell 5.1+, VS Code
    Version        : 2.0.0 (Modular)
    
    This modular version separates concerns for better maintainability and testing.
#>

param(
    [string]$RepositoryPath,
    [switch]$Clean,
    [switch]$AutoApprove,
    [switch]$Help,
    [switch]$ShowBackups
)

# Module import paths
$ModulesDir = Join-Path $PSScriptRoot "modules"
$CoreFunctionsPath = Join-Path $ModulesDir "core-functions.ps1"
$VSCodeSetupPath = Join-Path $ModulesDir "vscode-setup.ps1"
$CopilotInstallPath = Join-Path $ModulesDir "copilot-install.ps1"
$CleanupRestorePath = Join-Path $ModulesDir "cleanup-restore.ps1"

# Check if running in test mode (detect if we're in test-modular-install directory)
$IsTestMode = $PSScriptRoot -like "*test-modular-install*"
if ($IsTestMode) {
    Write-Host "[TEST MODE] Running modular installer in test mode" -ForegroundColor Magenta
}

function Show-Help {
    Write-Host "Terraform AzureRM Provider AI Setup (Modular Version)" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "USAGE:" -ForegroundColor Yellow
    Write-Host "  .\main-installer.ps1 [OPTIONS]"
    Write-Host ""
    Write-Host "OPTIONS:" -ForegroundColor Yellow
    Write-Host "  -RepositoryPath <path>   Specify repository path"
    Write-Host "  -Clean                   Remove AI setup and restore backups"
    Write-Host "  -AutoApprove             Skip interactive prompts"
    Write-Host "  -ShowBackups             Display backup information"
    Write-Host "  -Help                    Show this help message"
    Write-Host ""
    Write-Host "EXAMPLES:" -ForegroundColor Yellow
    Write-Host "  .\main-installer.ps1"
    Write-Host "  .\main-installer.ps1 -RepositoryPath C:\Code\terraform-provider-azurerm"
    Write-Host "  .\main-installer.ps1 -Clean"
    Write-Host "  .\main-installer.ps1 -ShowBackups"
    Write-Host ""
}

function Import-RequiredModules {
    $modules = @($CoreFunctionsPath, $VSCodeSetupPath, $CopilotInstallPath, $CleanupRestorePath)
    
    foreach ($module in $modules) {
        if (-not (Test-Path $module)) {
            Write-Host "[ERROR] Required module not found: $module" -ForegroundColor Red
            return $false
        }
        
        try {
            . $module  # Use dot sourcing instead of Import-Module for .ps1 files
            Write-Host "[INFO] Loaded module: $(Split-Path $module -Leaf)" -ForegroundColor Gray
        }
        catch {
            Write-Host "[ERROR] Failed to load module $module`: $($_.Exception.Message)" -ForegroundColor Red
            return $false
        }
    }
    
    return $true
}

function Get-BackupDirectory {
    param([string]$RepoPath)
    
    if ($IsTestMode) {
        return Join-Path $PSScriptRoot "test-backups"
    } else {
        return Join-Path $RepoPath ".github\backups"
    }
}

function Main {
    if ($Help) {
        Show-Help
        return 0
    }
    
    # Import all required modules
    if (-not (Import-RequiredModules)) {
        Write-Host "[ERROR] Failed to import required modules" -ForegroundColor Red
        return 1
    }
    
    # Check prerequisites
    try {
        Test-Prerequisites
        Write-Host "[INFO] Prerequisites check passed" -ForegroundColor Green
    }
    catch {
        Write-Host "[ERROR] Prerequisites not met: $($_.Exception.Message)" -ForegroundColor Red
        return 1
    }
    
    # Determine repository path
    if (-not $RepositoryPath) {
        if ($IsTestMode) {
            $RepositoryPath = Join-Path $PSScriptRoot "fake-repo"
            Write-Host "[TEST MODE] Using fake repository: $RepositoryPath" -ForegroundColor Magenta
        } else {
            try {
                $RepositoryPath = Find-RepositoryRoot
                if (-not $RepositoryPath) {
                    Write-Host "[ERROR] Could not find terraform-provider-azurerm repository" -ForegroundColor Red
                    return 1
                }
            }
            catch {
                Write-Host "[ERROR] Error finding repository: $($_.Exception.Message)" -ForegroundColor Red
                return 1
            }
        }
    }
    
    $BackupDir = Get-BackupDirectory -RepoPath $RepositoryPath
    
    # Handle show backups
    if ($ShowBackups) {
        try {
            Show-BackupSummary -BackupDir $BackupDir
            return 0
        }
        catch {
            Write-Host "[ERROR] Failed to show backup summary: $($_.Exception.Message)" -ForegroundColor Red
            return 1
        }
    }
    
    # Handle cleanup
    if ($Clean) {
        Write-Host "[INFO] Starting cleanup process..." -ForegroundColor Blue
        
        if (-not $AutoApprove) {
            $response = Read-Host "This will remove AI setup and restore backups. Continue? (y/N)"
            if ($response -ne 'y' -and $response -ne 'Y') {
                Write-Host "[INFO] Cleanup cancelled by user" -ForegroundColor Yellow
                return 0
            }
        }
        
        try {
            $success = Remove-AISetup -TargetRepoPath $RepositoryPath -BackupDir $BackupDir
            if ($success) {
                Write-Host "[SUCCESS] Cleanup completed successfully" -ForegroundColor Green
                return 0
            } else {
                Write-Host "[ERROR] Cleanup completed with errors" -ForegroundColor Red
                return 1
            }
        }
        catch {
            Write-Host "[ERROR] Cleanup failed: $($_.Exception.Message)" -ForegroundColor Red
            return 1
        }
    }
    
    # Handle installation
    Write-Host "[INFO] Starting AI setup installation..." -ForegroundColor Blue
    Write-Host "[INFO] Repository: $RepositoryPath" -ForegroundColor Blue
    Write-Host "[INFO] Backup directory: $BackupDir" -ForegroundColor Blue
    
    if ($IsTestMode) {
        Write-Host "[TEST MODE] This is a test installation - no real files will be modified" -ForegroundColor Magenta
    }
    
    # Check for previous installation
    try {
        $settingsPath = Get-VSCodeUserSettingsPath
        $previousInstallation = Test-PreviousInstallation -SettingsPath $settingsPath
        
        if ($previousInstallation -and -not $AutoApprove) {
            $response = Read-Host "Previous installation detected. Continue with reinstallation? (y/N)"
            if ($response -ne 'y' -and $response -ne 'Y') {
                Write-Host "[INFO] Installation cancelled by user" -ForegroundColor Yellow
                return 0
            }
        }
    }
    catch {
        Write-Host "[WARNING] Could not check for previous installation: $($_.Exception.Message)" -ForegroundColor Yellow
        $settingsPath = "$env:APPDATA\Code\User\settings.json"  # Default fallback
    }
    
    # Create backup directory
    if (-not (Test-Path $BackupDir)) {
        New-Item -ItemType Directory -Path $BackupDir -Force | Out-Null
        Write-Host "[INFO] Created backup directory: $BackupDir" -ForegroundColor Blue
    }
    
    # Backup VS Code settings
    try {
        $settingsBackup = Backup-VSCodeSettings -SettingsPath $settingsPath -BackupDir $BackupDir
        Write-Host "[INFO] VS Code settings backed up" -ForegroundColor Green
    }
    catch {
        Write-Host "[WARNING] Could not backup VS Code settings: $($_.Exception.Message)" -ForegroundColor Yellow
        $settingsBackup = $false
    }
    
    # Install VS Code settings
    try {
        $copilotSettings = Get-CopilotSettings
        $vsCodeSuccess = Install-VSCodeSettings -SettingsPath $settingsPath -NewSettings $copilotSettings
        if ($vsCodeSuccess) {
            Write-Host "[SUCCESS] VS Code settings installed" -ForegroundColor Green
        }
    }
    catch {
        Write-Host "[ERROR] Failed to install VS Code settings: $($_.Exception.Message)" -ForegroundColor Red
        $vsCodeSuccess = $false
    }
    
    # Install Copilot instructions (only in non-test mode)
    $copilotSuccess = $true
    if ($IsTestMode) {
        Write-Host "[TEST MODE] Skipping Copilot instruction installation" -ForegroundColor Magenta
    } else {
        try {
            # For real installation, we need to find the source repository (current script location)
            $sourceRepoPath = Split-Path (Split-Path $PSScriptRoot -Parent) -Parent
            $copilotSuccess = Install-CopilotInstructions -SourceRepoPath $sourceRepoPath -TargetRepoPath $RepositoryPath -BackupDir $BackupDir
            if ($copilotSuccess) {
                Write-Host "[SUCCESS] Copilot instructions installed" -ForegroundColor Green
            }
        }
        catch {
            Write-Host "[ERROR] Failed to install Copilot instructions: $($_.Exception.Message)" -ForegroundColor Red
            $copilotSuccess = $false
        }
    }
    
    # Verify installation
    if ($vsCodeSuccess -and $copilotSuccess) {
        Write-Host "[SUCCESS] AI setup installation completed successfully!" -ForegroundColor Green
        Write-Host "[INFO] VS Code should now have enhanced Copilot capabilities for AzureRM provider development" -ForegroundColor Blue
        
        if (-not $IsTestMode) {
            Write-Host "[INFO] Restart VS Code to ensure all changes take effect" -ForegroundColor Yellow
        }
        
        # Show backup summary
        try {
            Show-BackupSummary -BackupDir $BackupDir
        }
        catch {
            Write-Host "[WARNING] Could not show backup summary: $($_.Exception.Message)" -ForegroundColor Yellow
        }
        
        return 0
    } else {
        Write-Host "[ERROR] Installation completed with errors" -ForegroundColor Red
        
        # Offer to restore from backup
        if ($settingsBackup -and -not $AutoApprove) {
            $response = Read-Host "Installation failed. Restore from backup? (y/N)"
            if ($response -eq 'y' -or $response -eq 'Y') {
                try {
                    Restore-FromBackup -TargetRepoPath $RepositoryPath -BackupDir $BackupDir
                    Write-Host "[INFO] Restored from backup" -ForegroundColor Green
                }
                catch {
                    Write-Host "[ERROR] Failed to restore from backup: $($_.Exception.Message)" -ForegroundColor Red
                }
            }
        }
        
        return 1
    }
}

# Run main function and exit with appropriate code
$exitCode = Main
exit $exitCode
