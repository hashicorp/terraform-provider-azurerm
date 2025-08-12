# AI Installation Test Cleanup Utility (PowerShell)
# Cleans up temporary files created during AI installation testing

param(
    [switch]$Force = $false
)

Write-Host "AI Installation Test Cleanup Utility (PowerShell)" -ForegroundColor Cyan
Write-Host "======================================================" -ForegroundColor Cyan
Write-Host

# Function to safely remove directories
function Remove-SafeDirectory {
    param([string]$Path)
    
    if (Test-Path $Path) {
        try {
            Remove-Item $Path -Recurse -Force -ErrorAction Stop
            Write-Host "  Removed: $Path" -ForegroundColor Green
            return $true
        }
        catch {
            Write-Host "  Failed to remove: $Path - $($_.Exception.Message)" -ForegroundColor Yellow
            return $false
        }
    }
    return $false
}

# Clean up PowerShell test temporary files
Write-Host "Checking PowerShell temporary files..." -ForegroundColor Yellow

$tempPath = [System.IO.Path]::GetTempPath()
$azureRMDirs = Get-ChildItem -Path $tempPath -Directory -Name "AzureRM-*" -ErrorAction SilentlyContinue
$removedCount = 0

if ($azureRMDirs.Count -gt 0) {
    Write-Host "  Found $($azureRMDirs.Count) PowerShell temporary directories" -ForegroundColor White
    
    if ($Force -or (Read-Host "Remove all PowerShell temporary directories? (y/N)") -eq 'y') {
        foreach ($dir in $azureRMDirs) {
            $fullPath = Join-Path $tempPath $dir
            if (Remove-SafeDirectory $fullPath) {
                $removedCount++
            }
        }
        Write-Host "  Removed $removedCount PowerShell temporary directories" -ForegroundColor Green
    }
    else {
        Write-Host "  Cleanup cancelled by user" -ForegroundColor Yellow
    }
}
else {
    Write-Host "  No PowerShell temporary directories found" -ForegroundColor Green
}

# Clean up bash test temporary files (if accessible)
Write-Host "`nChecking bash temporary files..." -ForegroundColor Yellow

try {
    $bashTempDirs = wsl -- find /tmp -maxdepth 1 -name "terraform-azurerm-test-*" -type d 2>/dev/null
    if ($bashTempDirs) {
        $bashDirArray = $bashTempDirs -split "`n" | Where-Object { $_.Trim() -ne "" }
        Write-Host "  Found $($bashDirArray.Count) bash temporary directories" -ForegroundColor White
        
        if ($Force -or (Read-Host "Remove all bash temporary directories? (y/N)") -eq 'y') {
            foreach ($dir in $bashDirArray) {
                wsl -- rm -rf "$dir" 2>/dev/null
                if ($LASTEXITCODE -eq 0) {
                    Write-Host "  Removed: $dir" -ForegroundColor Green
                    $removedCount++
                }
                else {
                    Write-Host "  Failed to remove: $dir" -ForegroundColor Yellow
                }
            }
        }
    }
    else {
        Write-Host "  No bash temporary directories found" -ForegroundColor Green
    }
}
catch {
    Write-Host "  WSL not available or no bash temporary files" -ForegroundColor Gray
}

Write-Host "`nCleanup completed successfully!" -ForegroundColor Green

if (-not $Force) {
    Write-Host "`nTo run this script automatically in the future:" -ForegroundColor Cyan
    Write-Host "  .\cleanup-temp-files.ps1 -Force" -ForegroundColor White
}
