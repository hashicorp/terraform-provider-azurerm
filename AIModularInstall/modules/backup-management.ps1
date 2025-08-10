# Backup Management Module
# Functions for creating, managing, and restoring backups

function Create-SafeBackup {
    <#
    .SYNOPSIS
        Creates a safe backup with integrity verification
    #>
    param([string]$SourcePath, [string]$BackupPath)
    
    if (-not (Test-Path $SourcePath)) {
        return $false
    }
    
    try {
        # Create backup directory if it doesn't exist
        $backupDir = Split-Path $BackupPath -Parent
        if (-not (Test-Path $backupDir)) {
            New-Item -ItemType Directory -Path $backupDir -Force | Out-Null
        }
        
        # Create backup
        Copy-Item -Path $SourcePath -Destination $BackupPath -Force -ErrorAction Stop
        
        # Verify backup integrity
        if (-not (Test-FileIntegrity -FilePath $BackupPath)) {
            Remove-Item $BackupPath -Force -ErrorAction SilentlyContinue
            Write-StatusMessage "Backup verification failed for $SourcePath" "Error"
            return $false
        }
        
        Write-StatusMessage "Created backup: $BackupPath" "Success"
        return $true
    } catch {
        Write-StatusMessage "Failed to create backup for $SourcePath`: $_" "Error"
        return $false
    }
}

function Get-BackupLengthFromSettings {
    <#
    .SYNOPSIS
        Extracts backup length from VS Code settings.json
    #>
    param([string]$SettingsPath)
    
    if (-not (Test-Path $SettingsPath)) {
        return 0
    }
    
    try {
        $content = Get-Content $SettingsPath -Raw
        if ($content -match '"terraform_azurerm_backup_length":\s*(\d+)') {
            return [int]$matches[1]
        }
    } catch {
        Write-StatusMessage "Error reading backup length from settings: $_" "Warning"
    }
    
    return 0
}

function Get-MostRecentBackup {
    <#
    .SYNOPSIS
        Gets the most recent backup file
    #>
    param([string]$BackupDirectory, [string]$FilePattern = "*")
    
    if (-not (Test-Path $BackupDirectory)) {
        return $null
    }
    
    try {
        $backups = Get-ChildItem -Path $BackupDirectory -Filter $FilePattern | 
                   Where-Object { -not $_.PSIsContainer } |
                   Sort-Object LastWriteTime -Descending
        
        if ($backups.Count -gt 0) {
            return $backups[0].FullName
        }
    } catch {
        Write-StatusMessage "Error finding backup files: $_" "Warning"
    }
    
    return $null
}

function Restore-FromBackup {
    <#
    .SYNOPSIS
        Restores a file from backup
    #>
    param([string]$BackupPath, [string]$TargetPath)
    
    if (-not (Test-Path $BackupPath)) {
        Write-StatusMessage "Backup file not found: $BackupPath" "Error"
        return $false
    }
    
    try {
        # Verify backup integrity before restore
        if (-not (Test-FileIntegrity -FilePath $BackupPath)) {
            Write-StatusMessage "Backup file is corrupted: $BackupPath" "Error"
            return $false
        }
        
        # Create target directory if needed
        $targetDir = Split-Path $TargetPath -Parent
        if (-not (Test-Path $targetDir)) {
            New-Item -ItemType Directory -Path $targetDir -Force | Out-Null
        }
        
        Copy-Item -Path $BackupPath -Destination $TargetPath -Force -ErrorAction Stop
        Write-StatusMessage "Restored from backup: $TargetPath" "Success"
        return $true
    } catch {
        Write-StatusMessage "Failed to restore from backup: $_" "Error"
        return $false
    }
}

function Remove-BackupFiles {
    <#
    .SYNOPSIS
        Removes backup files and directories
    #>
    param([string]$BackupDirectory)
    
    if (-not (Test-Path $BackupDirectory)) {
        return $true
    }
    
    try {
        Remove-Item -Path $BackupDirectory -Recurse -Force -ErrorAction Stop
        Write-StatusMessage "Cleaned up backup directory: $BackupDirectory" "Success"
        return $true
    } catch {
        Write-StatusMessage "Failed to clean up backup directory: $_" "Warning"
        return $false
    }
}
