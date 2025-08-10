# AI Local Install - Test Environment Cleanup
# Cleans up all temporary test directories and files

param(
    [switch]$Force = $false,
    [switch]$DryRun = $false
)

$ErrorActionPreference = "Stop"

$TestBase = "C:\Users\$env:USERNAME\AppData\Local\Temp"
$TestPatterns = @(
    "AzureRM-*",
    "TerraformTest-*", 
    "InstallTest-*"
)

function Write-CleanupHeader {
    Write-Host ""
    Write-Host "============================================================" -ForegroundColor Yellow
    Write-Host "AI LOCAL INSTALL - TEST CLEANUP" -ForegroundColor Cyan
    Write-Host "============================================================" -ForegroundColor Yellow
    Write-Host ""
}

function Get-TestDirectories {
    $testDirs = @()
    
    foreach ($pattern in $TestPatterns) {
        $dirs = Get-ChildItem $TestBase -Directory -Name | Where-Object { $_ -like $pattern }
        $testDirs += $dirs | ForEach-Object { Join-Path $TestBase $_ }
    }
    
    return $testDirs | Sort-Object
}

function Show-DirectorySize {
    param([string]$Path)
    
    try {
        $size = (Get-ChildItem $Path -Recurse -File | Measure-Object -Property Length -Sum).Sum
        if ($size -lt 1KB) { 
            return "$size bytes" 
        } elseif ($size -lt 1MB) { 
            return "{0:N1} KB" -f ($size / 1KB) 
        } else { 
            return "{0:N1} MB" -f ($size / 1MB) 
        }
    } catch {
        return "Unknown"
    }
}

function Remove-TestDirectories {
    param([string[]]$Directories, [bool]$DryRun = $false)
    
    $totalRemoved = 0
    
    foreach ($dir in $Directories) {
        if (Test-Path $dir) {
            $size = Show-DirectorySize $dir
            $itemCount = (Get-ChildItem $dir -Recurse).Count
            
            if ($DryRun) {
                Write-Host "WOULD REMOVE: $dir ($size, $itemCount items)" -ForegroundColor Yellow
            } else {
                try {
                    Write-Host "REMOVING: $dir ($size, $itemCount items)" -ForegroundColor Red
                    Remove-Item $dir -Recurse -Force
                    $totalRemoved++
                } catch {
                    Write-Host "ERROR: Failed to remove $dir - $($_.Exception.Message)" -ForegroundColor Red
                }
            }
        }
    }
    
    return $totalRemoved
}

# Main execution
Write-CleanupHeader

Write-Host "Scanning for test directories..." -ForegroundColor Gray
$testDirectories = Get-TestDirectories

if ($testDirectories.Count -eq 0) {
    Write-Host "No test directories found. Environment is clean." -ForegroundColor Green
    exit 0
}

Write-Host "Found $($testDirectories.Count) test directories:" -ForegroundColor Cyan
foreach ($dir in $testDirectories) {
    $size = Show-DirectorySize $dir
    $age = (Get-Date) - (Get-Item $dir).CreationTime
    $ageHours = $age.TotalHours.ToString('F1')
    $dirName = Split-Path $dir -Leaf
    Write-Host "  * $dirName - $size - $ageHours hours old" -ForegroundColor White
}

$totalSize = 0
$testDirectories | ForEach-Object { 
    try { 
        $totalSize += (Get-ChildItem $_ -Recurse -File | Measure-Object -Property Length -Sum).Sum 
    } catch { 
        # Ignore errors calculating size
    }
}

if ($totalSize -lt 1KB) { 
    $totalSizeStr = "$totalSize bytes" 
} elseif ($totalSize -lt 1MB) { 
    $totalSizeStr = "{0:N1} KB" -f ($totalSize / 1KB) 
} else { 
    $totalSizeStr = "{0:N1} MB" -f ($totalSize / 1MB) 
}

Write-Host ""
Write-Host "Total space used: $totalSizeStr" -ForegroundColor Magenta

if ($DryRun) {
    Write-Host ""
    Write-Host "DRY RUN MODE - No files will be deleted" -ForegroundColor Yellow
    Remove-TestDirectories $testDirectories -DryRun $true
    Write-Host ""
    Write-Host "To actually remove these directories, run without -DryRun" -ForegroundColor Gray
    exit 0
}

if (-not $Force) {
    Write-Host ""
    Write-Host "WARNING: This will permanently delete all test directories!" -ForegroundColor Red
    $response = Read-Host "Continue? (Y)es/(N)o"
    
    if ($response -notmatch '^[Yy]') {
        Write-Host "Cleanup cancelled." -ForegroundColor Yellow
        exit 0
    }
}

Write-Host ""
Write-Host "Starting cleanup..." -ForegroundColor Yellow

$removed = Remove-TestDirectories $testDirectories

Write-Host ""
if ($removed -eq $testDirectories.Count) {
    Write-Host "SUCCESS: Successfully removed all $removed test directories" -ForegroundColor Green
    Write-Host "DISK SPACE: Freed up $totalSizeStr of disk space" -ForegroundColor Green
} else {
    $failed = $testDirectories.Count - $removed
    Write-Host "WARNING: Removed $removed/$($testDirectories.Count) directories ($failed failed)" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Cleanup completed!" -ForegroundColor Cyan
