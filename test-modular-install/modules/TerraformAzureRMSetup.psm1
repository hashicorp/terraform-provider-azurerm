# TerraformAzureRMSetup.psm1 - Root module that imports all sub-modules

# Import sub-modules
$ModulePath = $PSScriptRoot

Import-Module (Join-Path $ModulePath "CoreFunctions.psm1") -Force
Import-Module (Join-Path $ModulePath "VSCodeSetup.psm1") -Force  
Import-Module (Join-Path $ModulePath "CopilotInstall.psm1") -Force

function Install-TerraformAzureRMAI {
    <#
    .SYNOPSIS
    Main installation function for Terraform AzureRM Provider AI setup
    
    .PARAMETER RepositoryPath
    Path to the terraform-provider-azurerm repository
    
    .PARAMETER Force
    Force installation even if already installed
    
    .PARAMETER SkipBackup
    Skip creating backups
    
    .OUTPUTS
    Boolean indicating success
    #>
    
    [CmdletBinding()]
    param(
        [string]$RepositoryPath,
        [switch]$Force,
        [switch]$SkipBackup
    )
    
    try {
        # Find repository if not specified
        if (-not $RepositoryPath) {
            $RepositoryPath = Find-RepositoryRoot
            if (-not $RepositoryPath) {
                throw "Could not find terraform-provider-azurerm repository"
            }
        }
        
        Write-Verbose "Using repository path: $RepositoryPath"
        
        # Test prerequisites
        if (-not (Test-Prerequisites)) {
            throw "Prerequisites not met"
        }
        
        # Create backup directory
        $backupDir = Join-Path $RepositoryPath ".ai-setup-backups"
        
        if (-not $SkipBackup) {
            # Backup VS Code settings
            $vsCodeBackup = Backup-VSCodeSettings -BackupDirectory $backupDir
            if ($vsCodeBackup) {
                Write-Verbose "VS Code settings backed up"
            }
        }
        
        # Install VS Code settings
        if (-not (Install-VSCodeSettings -Force:$Force)) {
            Write-Warning "Failed to install VS Code settings"
        }
        
        # Install Copilot instructions
        if (-not (Install-CopilotInstructions -RepositoryPath $RepositoryPath -Force:$Force)) {
            Write-Warning "Failed to install Copilot instructions"
        }
        
        Write-Host "Terraform AzureRM Provider AI setup completed successfully!" -ForegroundColor Green
        return $true
    }
    catch {
        Write-Error "Installation failed: $($_.Exception.Message)"
        return $false
    }
}

function Remove-TerraformAzureRMAI {
    <#
    .SYNOPSIS
    Removes Terraform AzureRM Provider AI setup
    
    .PARAMETER RepositoryPath
    Path to the terraform-provider-azurerm repository
    
    .OUTPUTS
    Boolean indicating success
    #>
    
    [CmdletBinding()]
    param(
        [string]$RepositoryPath
    )
    
    try {
        # Find repository if not specified
        if (-not $RepositoryPath) {
            $RepositoryPath = Find-RepositoryRoot
            if (-not $RepositoryPath) {
                throw "Could not find terraform-provider-azurerm repository"
            }
        }
        
        # Remove Copilot instructions
        $result = Remove-CopilotInstructions -RepositoryPath $RepositoryPath
        
        if ($result) {
            Write-Host "AI setup removed successfully!" -ForegroundColor Green
        } else {
            Write-Host "No AI setup found to remove." -ForegroundColor Yellow
        }
        
        return $true
    }
    catch {
        Write-Error "Removal failed: $($_.Exception.Message)"
        return $false
    }
}

function Show-InstallationStatus {
    <#
    .SYNOPSIS
    Shows the current installation status
    
    .PARAMETER RepositoryPath
    Path to the terraform-provider-azurerm repository
    #>
    
    [CmdletBinding()]
    param(
        [string]$RepositoryPath
    )
    
    try {
        # Find repository if not specified
        if (-not $RepositoryPath) {
            $RepositoryPath = Find-RepositoryRoot
            if (-not $RepositoryPath) {
                Write-Warning "Could not find terraform-provider-azurerm repository"
                return
            }
        }
        
        Write-Host "=== Terraform AzureRM Provider AI Setup Status ===" -ForegroundColor Cyan
        Write-Host "Repository: $RepositoryPath" -ForegroundColor Gray
        
        # Check Copilot installation
        $copilotInfo = Get-CopilotInstallationInfo -RepositoryPath $RepositoryPath
        
        if ($copilotInfo.IsInstalled) {
            Write-Host "[INSTALLED] GitHub Copilot Instructions" -ForegroundColor Green
            Write-Host "  Path: $($copilotInfo.InstructionsPath)" -ForegroundColor Gray
            Write-Host "  Size: $($copilotInfo.FileSize) bytes" -ForegroundColor Gray
            Write-Host "  Modified: $($copilotInfo.LastModified)" -ForegroundColor Gray
        } else {
            Write-Host "[NOT INSTALLED] GitHub Copilot Instructions" -ForegroundColor Red
        }
        
        # Check VS Code settings
        $vsCodePath = Get-VSCodeUserSettingsPath
        if (Test-Path $vsCodePath) {
            Write-Host "[EXISTS] VS Code User Settings" -ForegroundColor Green
            Write-Host "  Path: $vsCodePath" -ForegroundColor Gray
        } else {
            Write-Host "[NOT FOUND] VS Code User Settings" -ForegroundColor Yellow
        }
        
        # Check for backups
        $backupDir = Join-Path $RepositoryPath ".ai-setup-backups"
        if (Test-Path $backupDir) {
            $backupFiles = Get-ChildItem $backupDir -File
            Write-Host "[FOUND] Backup Directory ($($backupFiles.Count) files)" -ForegroundColor Blue
            Write-Host "  Path: $backupDir" -ForegroundColor Gray
        } else {
            Write-Host "[NOT FOUND] Backup Directory" -ForegroundColor Gray
        }
    }
    catch {
        Write-Error "Status check failed: $($_.Exception.Message)"
    }
}
