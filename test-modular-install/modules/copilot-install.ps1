<#
.SYNOPSIS
    Copilot instruction file installation for Terraform AzureRM Provider AI Setup

.DESCRIPTION
    Handles copying and installation of GitHub Copilot instruction files and AI prompts
    for the enhanced development environment.
#>

function Install-CopilotInstructions {
    param(
        [string]$SourceRepoPath,
        [string]$TargetRepoPath,
        [string]$BackupDir
    )
    
    Write-Host "[INFO] Installing Copilot instruction files..." -ForegroundColor Blue
    
    # Import core functions
    Import-Module (Join-Path $PSScriptRoot "core-functions.ps1") -Force
    
    # Define source and target paths
    $sourceInstructionsPath = Join-Path $SourceRepoPath ".github\copilot-instructions.md"
    $targetInstructionsPath = Join-Path $TargetRepoPath ".github\copilot-instructions.md"
    
    $sourceInstructionsDir = Join-Path $SourceRepoPath ".github\instructions"
    $targetInstructionsDir = Join-Path $TargetRepoPath ".github\instructions"
    
    # Backup existing files
    if (Test-Path $targetInstructionsPath) {
        $backup = New-SafeBackup -FilePath $targetInstructionsPath -BackupDir $BackupDir
        if (-not $backup) {
            Write-Error "Failed to backup existing copilot-instructions.md"
            return $false
        }
    }
    
    # Copy main instruction file
    try {
        $targetGithubDir = Split-Path $targetInstructionsPath -Parent
        if (-not (Test-Path $targetGithubDir)) {
            New-Item -ItemType Directory -Path $targetGithubDir -Force | Out-Null
        }
        
        Copy-Item $sourceInstructionsPath $targetInstructionsPath -Force
        Write-Host "[SUCCESS] Copied copilot-instructions.md" -ForegroundColor Green
        
        # Verify file integrity
        if (Test-FileIntegrity -SourcePath $sourceInstructionsPath -DestinationPath $targetInstructionsPath) {
            Write-Host "[SUCCESS] File integrity verified" -ForegroundColor Green
        } else {
            Write-Warning "File integrity check failed for copilot-instructions.md"
        }
    }
    catch {
        Write-Error "Failed to copy copilot-instructions.md: $_"
        return $false
    }
    
    # Copy instruction directory
    if (Test-Path $sourceInstructionsDir) {
        try {
            if (Test-Path $targetInstructionsDir) {
                # Backup existing instructions directory
                $instrBackup = New-SafeBackup -FilePath $targetInstructionsDir -BackupDir $BackupDir
                if (-not $instrBackup) {
                    Write-Warning "Failed to backup existing instructions directory"
                }
            }
            
            Copy-Item $sourceInstructionsDir $targetInstructionsDir -Recurse -Force
            Write-Host "[SUCCESS] Copied instructions directory" -ForegroundColor Green
            
            # Count files copied
            $fileCount = (Get-ChildItem $targetInstructionsDir -Recurse -File).Count
            Write-Host "[INFO] Copied $fileCount instruction files" -ForegroundColor Blue
        }
        catch {
            Write-Error "Failed to copy instructions directory: $_"
            return $false
        }
    }
    
    return $true
}

function Remove-CopilotInstructions {
    param(
        [string]$TargetRepoPath,
        [string]$BackupDir
    )
    
    Write-Host "[INFO] Removing Copilot instruction files..." -ForegroundColor Blue
    
    $targetInstructionsPath = Join-Path $TargetRepoPath ".github\copilot-instructions.md"
    $targetInstructionsDir = Join-Path $TargetRepoPath ".github\instructions"
    
    $success = $true
    
    # Remove main instruction file
    if (Test-Path $targetInstructionsPath) {
        try {
            Remove-Item $targetInstructionsPath -Force
            Write-Host "[SUCCESS] Removed copilot-instructions.md" -ForegroundColor Green
        }
        catch {
            Write-Error "Failed to remove copilot-instructions.md: $_"
            $success = $false
        }
    }
    
    # Remove instructions directory
    if (Test-Path $targetInstructionsDir) {
        try {
            Remove-Item $targetInstructionsDir -Recurse -Force
            Write-Host "[SUCCESS] Removed instructions directory" -ForegroundColor Green
        }
        catch {
            Write-Error "Failed to remove instructions directory: $_"
            $success = $false
        }
    }
    
    return $success
}

function Test-CopilotInstallation {
    param([string]$TargetRepoPath)
    
    $targetInstructionsPath = Join-Path $TargetRepoPath ".github\copilot-instructions.md"
    $targetInstructionsDir = Join-Path $TargetRepoPath ".github\instructions"
    
    $mainFileExists = Test-Path $targetInstructionsPath
    $instructionsDirExists = Test-Path $targetInstructionsDir
    
    if ($mainFileExists -and $instructionsDirExists) {
        $fileCount = (Get-ChildItem $targetInstructionsDir -Recurse -File -ErrorAction SilentlyContinue).Count
        Write-Host "[SUCCESS] Copilot installation verified - $fileCount instruction files found" -ForegroundColor Green
        return $true
    } elseif ($mainFileExists) {
        Write-Host "[WARNING] Main instruction file found but instructions directory missing" -ForegroundColor Yellow
        return $false
    } elseif ($instructionsDirExists) {
        Write-Host "[WARNING] Instructions directory found but main file missing" -ForegroundColor Yellow
        return $false
    } else {
        Write-Host "[INFO] No Copilot installation found" -ForegroundColor Yellow
        return $false
    }
}

function Restore-CopilotInstructions {
    param(
        [string]$TargetRepoPath,
        [string]$BackupDir
    )
    
    Write-Host "[INFO] Restoring Copilot instructions from backup..." -ForegroundColor Blue
    
    # Import core functions
    Import-Module (Join-Path $PSScriptRoot "core-functions.ps1") -Force
    
    if (-not (Test-Path $BackupDir)) {
        Write-Host "[WARNING] No backup directory found" -ForegroundColor Yellow
        return $false
    }
    
    # Find most recent backup
    $backupFiles = Get-ChildItem $BackupDir -Filter "*.backup_*" | Sort-Object LastWriteTime -Descending
    
    foreach ($backupFile in $backupFiles) {
        if ($backupFile.Name -like "copilot-instructions.md.backup_*") {
            $targetPath = Join-Path $TargetRepoPath ".github\copilot-instructions.md"
            
            try {
                Copy-Item $backupFile.FullName $targetPath -Force
                Write-Host "[SUCCESS] Restored copilot-instructions.md from backup" -ForegroundColor Green
                return $true
            }
            catch {
                Write-Error "Failed to restore backup: $_"
                return $false
            }
        }
    }
    
    Write-Host "[WARNING] No copilot-instructions.md backup found" -ForegroundColor Yellow
    return $false
}

# Copilot functions are loaded via dot sourcing - no export needed
