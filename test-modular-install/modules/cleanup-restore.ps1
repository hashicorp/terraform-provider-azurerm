<#
.SYNOPSIS
    Cleanup and restore functionality for Terraform AzureRM Provider AI Setup

.DESCRIPTION
    Handles complete removal of AI setup and restoration from backups,
    providing a clean uninstall experience.
#>

function Remove-AISetup {
    param(
        [string]$TargetRepoPath,
        [string]$BackupDir
    )
    
    Write-Host "[INFO] Starting AI setup removal..." -ForegroundColor Blue
    
    # Import required modules
    Import-Module (Join-Path $PSScriptRoot "vscode-setup.ps1") -Force
    Import-Module (Join-Path $PSScriptRoot "copilot-install.ps1") -Force
    
    $success = $true
    
    # Remove VS Code settings
    $settingsPath = Get-VSCodeUserSettingsPath
    if (-not (Remove-CopilotSettings -SettingsPath $settingsPath)) {
        $success = $false
    }
    
    # Remove Copilot instructions
    if (-not (Remove-CopilotInstructions -TargetRepoPath $TargetRepoPath -BackupDir $BackupDir)) {
        $success = $false
    }
    
    if ($success) {
        Write-Host "[SUCCESS] AI setup removal completed" -ForegroundColor Green
    } else {
        Write-Host "[WARNING] AI setup removal completed with some errors" -ForegroundColor Yellow
    }
    
    return $success
}

function Restore-FromBackup {
    param(
        [string]$TargetRepoPath,
        [string]$BackupDir
    )
    
    Write-Host "[INFO] Starting restoration from backup..." -ForegroundColor Blue
    
    if (-not (Test-Path $BackupDir)) {
        Write-Host "[ERROR] Backup directory not found: $BackupDir" -ForegroundColor Red
        return $false
    }
    
    # Import required modules
    Import-Module (Join-Path $PSScriptRoot "vscode-setup.ps1") -Force
    Import-Module (Join-Path $PSScriptRoot "copilot-install.ps1") -Force
    
    $success = $true
    
    # Restore VS Code settings
    $settingsPath = Get-VSCodeUserSettingsPath
    $settingsBackups = Get-ChildItem $BackupDir -Filter "settings.json.backup_*" | Sort-Object LastWriteTime -Descending
    
    if ($settingsBackups.Count -gt 0) {
        $latestBackup = $settingsBackups[0]
        try {
            Copy-Item $latestBackup.FullName $settingsPath -Force
            Write-Host "[SUCCESS] Restored VS Code settings from backup" -ForegroundColor Green
        }
        catch {
            Write-Error "Failed to restore VS Code settings: $_"
            $success = $false
        }
    } else {
        Write-Host "[INFO] No VS Code settings backup found" -ForegroundColor Yellow
    }
    
    # Restore Copilot instructions
    if (-not (Restore-CopilotInstructions -TargetRepoPath $TargetRepoPath -BackupDir $BackupDir)) {
        $success = $false
    }
    
    if ($success) {
        Write-Host "[SUCCESS] Restoration completed" -ForegroundColor Green
    } else {
        Write-Host "[WARNING] Restoration completed with some errors" -ForegroundColor Yellow
    }
    
    return $success
}

function Get-BackupInfo {
    param([string]$BackupDir)
    
    if (-not (Test-Path $BackupDir)) {
        Write-Host "[INFO] No backup directory found" -ForegroundColor Yellow
        return @()
    }
    
    $backups = Get-ChildItem $BackupDir -Filter "*.backup_*" | Sort-Object LastWriteTime -Descending
    
    $backupInfo = @()
    foreach ($backup in $backups) {
        $backupInfo += [PSCustomObject]@{
            FileName = $backup.Name
            OriginalName = ($backup.Name -split '\.backup_')[0]
            BackupDate = $backup.LastWriteTime
            Size = $backup.Length
            Path = $backup.FullName
        }
    }
    
    return $backupInfo
}

function Show-BackupSummary {
    param([string]$BackupDir)
    
    Write-Host "[INFO] Backup Summary" -ForegroundColor Blue
    Write-Host "===================" -ForegroundColor Blue
    
    $backups = Get-BackupInfo -BackupDir $BackupDir
    
    if ($backups.Count -eq 0) {
        Write-Host "[INFO] No backups found" -ForegroundColor Yellow
        return
    }
    
    foreach ($backup in $backups) {
        $sizeKB = [math]::Round($backup.Size / 1024, 2)
        Write-Host "[INFO] $($backup.OriginalName)" -ForegroundColor Cyan
        Write-Host "       Backed up: $($backup.BackupDate)" -ForegroundColor Gray
        Write-Host "       Size: $sizeKB KB" -ForegroundColor Gray
        Write-Host "" 
    }
    
    Write-Host "[INFO] Total backups: $($backups.Count)" -ForegroundColor Blue
}

function Remove-OldBackups {
    param(
        [string]$BackupDir,
        [int]$RetainCount = 5
    )
    
    if (-not (Test-Path $BackupDir)) {
        Write-Host "[INFO] No backup directory found" -ForegroundColor Yellow
        return $true
    }
    
    Write-Host "[INFO] Cleaning up old backups (retaining $RetainCount most recent)..." -ForegroundColor Blue
    
    $backups = Get-ChildItem $BackupDir -Filter "*.backup_*" | Sort-Object LastWriteTime -Descending
    $toRemove = $backups | Select-Object -Skip $RetainCount
    
    if ($toRemove.Count -eq 0) {
        Write-Host "[INFO] No old backups to remove" -ForegroundColor Green
        return $true
    }
    
    foreach ($backup in $toRemove) {
        try {
            Remove-Item $backup.FullName -Force
            Write-Host "[INFO] Removed old backup: $($backup.Name)" -ForegroundColor Yellow
        }
        catch {
            Write-Warning "Failed to remove backup $($backup.Name): $_"
        }
    }
    
    Write-Host "[SUCCESS] Backup cleanup completed" -ForegroundColor Green
    return $true
}

# Cleanup functions are loaded via dot sourcing - no export needed
