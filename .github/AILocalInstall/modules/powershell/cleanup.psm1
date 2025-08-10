# Cleanup Module
# Functions for removing AI agent installation and restoring backups

function Remove-AIAgent {
    <#
    .SYNOPSIS
        Removes AI agent installation and restores backups
    #>
    param([string]$RepositoryPath)
    
    Write-StatusMessage "Starting AI agent removal..." "Info"
    
    $paths = Get-InstallationPaths -RepositoryPath $RepositoryPath
    $success = $true
    
    # Restore VS Code settings from backup
    if (Test-Path $paths.VSCodeBackupDir) {
        $latestBackup = Get-MostRecentBackup -BackupDirectory $paths.VSCodeBackupDir -FilePattern "settings_backup_*.json"
        
        if ($latestBackup) {
            if (Restore-FromBackup -BackupPath $latestBackup -TargetPath $paths.VSCodeSettings) {
                Write-StatusMessage "VS Code settings restored from backup" "Success"
            } else {
                Write-StatusMessage "Failed to restore VS Code settings" "Error"
                $success = $false
            }
        } else {
            # No backup found, remove our settings manually
            if (Test-Path $paths.VSCodeSettings) {
                if (Remove-TerraformSettingsFromVSCode -SettingsPath $paths.VSCodeSettings) {
                    Write-StatusMessage "Terraform settings removed from VS Code" "Success"
                } else {
                    Write-StatusMessage "Failed to clean Terraform settings from VS Code" "Warning"
                }
            }
        }
        
        # Clean up backup directory
        Remove-BackupFiles -BackupDirectory $paths.VSCodeBackupDir
    }
    
    # Note: We don't remove the .github/instructions directory as it's part of the repository
    # Only the VS Code configuration is considered "installed"
    
    if ($success) {
        Write-StatusMessage "AI agent removal completed successfully!" "Success"
    } else {
        Write-StatusMessage "AI agent removal completed with some errors" "Warning"
    }
    
    return $success
}

function Remove-TerraformSettingsFromVSCode {
    <#
    .SYNOPSIS
        Removes Terraform-specific settings from VS Code settings.json
    #>
    param([string]$SettingsPath)
    
    if (-not (Test-Path $SettingsPath)) {
        return $true
    }
    
    try {
        $content = Get-Content $SettingsPath -Raw
        $settings = $content | ConvertFrom-Json
        
        # Convert to hashtable if needed for PowerShell 5.1 compatibility
        if ($settings -is [PSCustomObject]) {
            $hashtable = @{}
            $settings.PSObject.Properties | ForEach-Object {
                $hashtable[$_.Name] = $_.Value
            }
            $settings = $hashtable
        }
        
        # Remove Terraform AzureRM specific settings
        $terraformKeys = $settings.Keys | Where-Object { $_ -like "terraform_azurerm_*" }
        foreach ($key in $terraformKeys) {
            $settings.Remove($key)
        }
        
        # Write back the cleaned settings
        $settings | ConvertTo-Json -Depth 10 | Set-Content -Path $SettingsPath -Encoding UTF8
        
        Write-StatusMessage "Cleaned Terraform settings from VS Code" "Success"
        return $true
    } catch {
        Write-StatusMessage "Failed to clean VS Code settings: $_" "Error"
        return $false
    }
}

function Test-CleanupSuccess {
    <#
    .SYNOPSIS
        Verifies that cleanup was successful
    #>
    param([string]$RepositoryPath)
    
    $installationState = Test-PreviousInstallation -RepositoryPath $RepositoryPath
    
    # Check if VS Code settings still contain our markers
    if ($installationState.HasVSCodeSettings) {
        Write-StatusMessage "VS Code still contains Terraform settings" "Warning"
        return $false
    }
    
    # Check if backup directory still exists
    if ($installationState.HasBackups) {
        Write-StatusMessage "Backup directory still exists" "Warning"
        return $false
    }
    
    Write-StatusMessage "Cleanup verification successful" "Success"
    return $true
}

function Show-CleanupSummary {
    <#
    .SYNOPSIS
        Shows a summary of cleanup operations
    #>
    param([string]$RepositoryPath)
    
    Write-Host ""
    Write-Host "Cleanup Summary:" -ForegroundColor Cyan
    Write-Host "===============" -ForegroundColor Cyan
    
    $installationState = Test-PreviousInstallation -RepositoryPath $RepositoryPath
    
    Write-Host "VS Code Settings: " -NoNewline
    if ($installationState.HasVSCodeSettings) {
        Write-Host "Still contains Terraform settings" -ForegroundColor Yellow
    } else {
        Write-Host "Cleaned" -ForegroundColor Green
    }
    
    Write-Host "Backup Directory: " -NoNewline
    if ($installationState.HasBackups) {
        Write-Host "Still exists" -ForegroundColor Yellow
    } else {
        Write-Host "Removed" -ForegroundColor Green
    }
    
    Write-Host "Instructions: " -NoNewline
    if ($installationState.HasInstructions) {
        Write-Host "Preserved (part of repository)" -ForegroundColor Green
    } else {
        Write-Host "Not found" -ForegroundColor Yellow
    }
    
    Write-Host ""
}
