#!/usr/bin/env pwsh
<#
.SYNOPSIS
    Cleanup script for AI installation test temporary files

.DESCRIPTION
    Removes all temporary files and directories created during AI installation testing.
    This includes both Windows PowerShell and WSL bash temporary directories.

.PARAMETER Force
    Skip confirmation prompts and force cleanup

.EXAMPLE
    .\cleanup-temp-files.ps1
    
.EXAMPLE
    .\cleanup-temp-files.ps1 -Force
#>

param(
    [switch]$Force
)

Write-Host "AI Installation Test Cleanup Utility" -ForegroundColor Cyan
Write-Host "====================================" -ForegroundColor Cyan
Write-Host ""

# Function to clean Windows temporary files
function Remove-WindowsTempFiles {
    Write-Host "Checking Windows temporary files..." -ForegroundColor Yellow
    
    $tempDirs = Get-ChildItem $env:TEMP -Filter "*AzureRM*" -Directory -ErrorAction SilentlyContinue
    
    if ($tempDirs.Count -eq 0) {
        Write-Host "  No Windows temporary directories found" -ForegroundColor Green
        return
    }
    
    Write-Host "  Found $($tempDirs.Count) temporary directories" -ForegroundColor Yellow
    
    if (-not $Force) {
        $response = Read-Host "Remove $($tempDirs.Count) Windows temporary directories? (y/N)"
        if ($response -ne 'y' -and $response -ne 'Y') {
            Write-Host "  Skipped Windows cleanup" -ForegroundColor Yellow
            return
        }
    }
    
    $removed = 0
    foreach ($dir in $tempDirs) {
        try {
            Remove-Item $dir.FullName -Recurse -Force -ErrorAction Stop
            $removed++
            Write-Host "  Removed: $($dir.Name)" -ForegroundColor Green
        }
        catch {
            Write-Host "  Failed to remove: $($dir.Name) - $($_.Exception.Message)" -ForegroundColor Red
        }
    }
    
    Write-Host "  Successfully removed $removed/$($tempDirs.Count) directories" -ForegroundColor Green
}

# Function to clean WSL/Bash temporary files
function Remove-BashTempFiles {
    Write-Host "Checking WSL/Bash temporary files..." -ForegroundColor Yellow
    
    # Check if WSL is available
    try {
        $null = wsl -- echo "test" 2>$null
    }
    catch {
        Write-Host "  WSL not available, skipping bash cleanup" -ForegroundColor Yellow
        return
    }
    
    # Count bash temp directories
    try {
        $bashTempCount = wsl -- bash -c "find /tmp -name '*terraform-azurerm*' -type d 2>/dev/null | wc -l"
        $bashTempCount = [int]$bashTempCount.Trim()
    }
    catch {
        Write-Host "  Could not check bash temporary files" -ForegroundColor Yellow
        return
    }
    
    if ($bashTempCount -eq 0) {
        Write-Host "  No bash temporary directories found" -ForegroundColor Green
        return
    }
    
    Write-Host "  Found $bashTempCount bash temporary directories" -ForegroundColor Yellow
    
    if (-not $Force) {
        $response = Read-Host "Remove $bashTempCount bash temporary directories? (y/N)"
        if ($response -ne 'y' -and $response -ne 'Y') {
            Write-Host "  Skipped bash cleanup" -ForegroundColor Yellow
            return
        }
    }
    
    try {
        wsl -- bash -c "rm -rf /tmp/terraform-azurerm-test-* 2>/dev/null"
        Write-Host "  Successfully removed bash temporary directories" -ForegroundColor Green
    }
    catch {
        Write-Host "  Failed to remove bash temporary directories: $($_.Exception.Message)" -ForegroundColor Red
    }
}

# Function to clean any VS Code backup files
function Remove-VSCodeBackups {
    Write-Host "Checking for temporary VS Code backup files..." -ForegroundColor Yellow
    
    $backupPatterns = @(
        "$env:TEMP\*vscode*backup*",
        "$env:TEMP\*copilot*backup*"
    )
    
    $foundBackups = $false
    foreach ($pattern in $backupPatterns) {
        $backups = Get-ChildItem $pattern -ErrorAction SilentlyContinue
        if ($backups) {
            $foundBackups = $true
            Write-Host "  Found backup files matching: $pattern" -ForegroundColor Yellow
        }
    }
    
    if (-not $foundBackups) {
        Write-Host "  No temporary backup files found" -ForegroundColor Green
        return
    }
    
    if (-not $Force) {
        $response = Read-Host "Remove temporary backup files? (y/N)"
        if ($response -ne 'y' -and $response -ne 'Y') {
            Write-Host "  Skipped backup cleanup" -ForegroundColor Yellow
            return
        }
    }
    
    foreach ($pattern in $backupPatterns) {
        $backups = Get-ChildItem $pattern -ErrorAction SilentlyContinue
        foreach ($backup in $backups) {
            try {
                Remove-Item $backup.FullName -Recurse -Force -ErrorAction Stop
                Write-Host "  Removed: $($backup.Name)" -ForegroundColor Green
            }
            catch {
                Write-Host "  Failed to remove: $($backup.Name) - $($_.Exception.Message)" -ForegroundColor Red
            }
        }
    }
}

# Main execution
try {
    Remove-WindowsTempFiles
    Write-Host ""
    
    Remove-BashTempFiles
    Write-Host ""
    
    Remove-VSCodeBackups
    Write-Host ""
    
    Write-Host "Cleanup completed successfully!" -ForegroundColor Green
}
catch {
    Write-Host "Cleanup failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "To run this script automatically in the future:" -ForegroundColor Cyan
Write-Host "  .\cleanup-temp-files.ps1 -Force" -ForegroundColor White
