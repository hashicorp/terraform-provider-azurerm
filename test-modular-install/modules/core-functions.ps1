<#
.SYNOPSIS
    Core utility functions for Terraform AzureRM Provider AI Setup

.DESCRIPTION
    Provides repository discovery, file integrity checking, and backup functionality
    for the modular AI setup installation system.
#>

function Find-RepositoryRoot {
    param([string]$StartPath = (Get-Location).Path)
    
    Write-Host "[INFO] Searching for terraform-provider-azurerm repository..." -ForegroundColor Blue
    
    $currentPath = $StartPath
    $maxDepth = 10
    $depth = 0
    
    while ($depth -lt $maxDepth) {
        $potentialRepoPath = Join-Path $currentPath "terraform-provider-azurerm"
        if (Test-Path $potentialRepoPath -PathType Container) {
            $goModPath = Join-Path $potentialRepoPath "go.mod"
            if (Test-Path $goModPath) {
                $goModContent = Get-Content $goModPath -Raw
                if ($goModContent -match "module.*terraform-provider-azurerm") {
                    Write-Host "[SUCCESS] Found repository at: $potentialRepoPath" -ForegroundColor Green
                    return $potentialRepoPath
                }
            }
        }
        
        # Check if current directory IS the repo
        $goModPath = Join-Path $currentPath "go.mod"
        if (Test-Path $goModPath) {
            $goModContent = Get-Content $goModPath -Raw
            if ($goModContent -match "module.*terraform-provider-azurerm") {
                Write-Host "[SUCCESS] Found repository at: $currentPath" -ForegroundColor Green
                return $currentPath
            }
        }
        
        $parentPath = Split-Path $currentPath -Parent
        if ($parentPath -eq $currentPath) { break }
        $currentPath = $parentPath
        $depth++
    }
    
    Write-Host "[ERROR] Repository not found in search path" -ForegroundColor Red
    return $null
}

function Test-FileIntegrity {
    param(
        [string]$SourcePath,
        [string]$DestinationPath
    )
    
    if (-not (Test-Path $SourcePath)) {
        Write-Warning "Source file not found: $SourcePath"
        return $false
    }
    
    if (-not (Test-Path $DestinationPath)) {
        Write-Warning "Destination file not found: $DestinationPath"
        return $false
    }
    
    $sourceHash = Get-FileHash $SourcePath -Algorithm SHA256
    $destHash = Get-FileHash $DestinationPath -Algorithm SHA256
    
    return $sourceHash.Hash -eq $destHash.Hash
}

function New-SafeBackup {
    param(
        [string]$FilePath,
        [string]$BackupDir
    )
    
    if (-not (Test-Path $FilePath)) {
        Write-Host "[INFO] No existing file to backup: $FilePath" -ForegroundColor Yellow
        return $null
    }
    
    if (-not (Test-Path $BackupDir)) {
        New-Item -ItemType Directory -Path $BackupDir -Force | Out-Null
    }
    
    $timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
    $fileName = [System.IO.Path]::GetFileName($FilePath)
    $backupName = "${fileName}.backup_${timestamp}"
    $backupPath = Join-Path $BackupDir $backupName
    
    try {
        Copy-Item $FilePath $backupPath -Force
        Write-Host "[SUCCESS] Created backup: $backupPath" -ForegroundColor Green
        return $backupPath
    }
    catch {
        Write-Error "Failed to create backup: $_"
        return $null
    }
}

function Test-Prerequisites {
    Write-Host "[INFO] Checking prerequisites..." -ForegroundColor Blue
    
    # Check PowerShell version
    if ($PSVersionTable.PSVersion.Major -lt 5) {
        Write-Error "PowerShell 5.1 or later is required"
        return $false
    }
    
    # Check if VS Code is installed
    $vscodeExe = Get-Command "code" -ErrorAction SilentlyContinue
    if (-not $vscodeExe) {
        Write-Warning "VS Code 'code' command not found in PATH. Install VS Code or add it to PATH."
        return $false
    }
    
    Write-Host "[SUCCESS] Prerequisites satisfied" -ForegroundColor Green
    return $true
}

# Core functions are loaded via dot sourcing - no export needed
