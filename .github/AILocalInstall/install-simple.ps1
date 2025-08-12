#!/usr/bin/env pwsh
<#
.SYNOPSIS
    Simple Terraform AzureRM Provider AI Setup

.DESCRIPTION
    A simple script to install, verify, or clean GitHub Copilot instruction files
    for the Terraform AzureRM Provider. Follows KISS principles.

.PARAMETER Verify
    Lists what AI files are currently installed

.PARAMETER Clean
    Removes installed AI files and restores backups

.PARAMETER RepoPath
    Path to terraform-provider-azurerm repository (auto-discovered if not provided)

.EXAMPLE
    .\install-simple.ps1
    Installs AI setup files

.EXAMPLE
    .\install-simple.ps1 -Verify
    Lists installed AI files

.EXAMPLE
    .\install-simple.ps1 -Clean
    Removes AI setup and restores backups
#>

param(
    [switch]$Verify,
    [switch]$Clean,
    [string]$RepoPath
)

# Auto-discover repository if not provided
if (-not $RepoPath) {
    $currentDir = Get-Location
    $searchPaths = @($currentDir, $currentDir.Parent, $currentDir.Parent.Parent)
    
    foreach ($path in $searchPaths) {
        if (Test-Path "$path\.github\instructions" -PathType Container) {
            $RepoPath = $path.FullName
            break
        }
    }
    
    if (-not $RepoPath) {
        Write-Error "Could not find terraform-provider-azurerm repository. Please specify -RepoPath"
        exit 1
    }
}

# Paths
$manifestPath = "$RepoPath\.github\AILocalInstall\file-manifest.config"
$instructionsDir = "$RepoPath\.github\instructions"
$promptsDir = "$RepoPath\.github\prompts"
$vscodeDirPlatform = if ($IsWindows -or $env:OS -eq "Windows_NT") { "$env:APPDATA\Code\User" } else { "$env:HOME/.config/Code/User" }
$vscodeDir = $vscodeDirPlatform
$instructionsDestDir = "$vscodeDir\instructions\terraform-azurerm"
$promptsDestDir = "$vscodeDir\prompts"
$settingsPath = "$vscodeDir\settings.json"
$backupDir = "$vscodeDir\ai-backup-$(Get-Date -Format 'yyyyMMdd')"

# Read manifest file
function Read-Manifest {
    if (-not (Test-Path $manifestPath)) {
        Write-Error "Manifest file not found: $manifestPath"
        exit 1
    }
    
    $manifest = @{
        InstructionFiles = @()
        PromptFiles = @()
        MainFiles = @()
    }
    
    $currentSection = $null
    Get-Content $manifestPath | ForEach-Object {
        $line = $_.Trim()
        if ($line -match '^\[(.+)\]$') {
            $currentSection = $matches[1]
        }
        elseif ($line -and -not $line.StartsWith('#')) {
            switch ($currentSection) {
                'INSTRUCTION_FILES' { $manifest.InstructionFiles += $line }
                'PROMPT_FILES' { $manifest.PromptFiles += $line }
                'MAIN_FILES' { $manifest.MainFiles += $line }
            }
        }
    }
    
    return $manifest
}

# Backup existing files
function Backup-Files {
    param($filesToBackup)
    
    if ($filesToBackup.Count -eq 0) { return }
    
    New-Item -ItemType Directory -Path $backupDir -Force | Out-Null
    
    foreach ($file in $filesToBackup) {
        if (Test-Path $file) {
            $relativePath = $file.Replace($vscodeDir, '').TrimStart('\', '/')
            $backupPath = Join-Path $backupDir $relativePath
            $backupParent = Split-Path $backupPath -Parent
            
            if ($backupParent) {
                New-Item -ItemType Directory -Path $backupParent -Force | Out-Null
            }
            
            Copy-Item $file $backupPath -Force
            Write-Host "Backed up: $relativePath"
        }
    }
}

# Update settings.json with Terraform AzureRM Provider specific settings
function Update-SettingsJson {
    # Read existing settings or create new
    $settings = @{}
    if (Test-Path $settingsPath) {
        try {
            $existingContent = Get-Content $settingsPath -Raw
            if ($existingContent.Trim()) {
                $tempSettings = $existingContent | ConvertFrom-Json
                
                # Convert to hashtable for PowerShell 5.1 compatibility
                if ($tempSettings -is [PSCustomObject]) {
                    $settings = @{}
                    $tempSettings.PSObject.Properties | ForEach-Object {
                        $settings[$_.Name] = $_.Value
                    }
                } else {
                    $settings = $tempSettings
                }
            }
        } catch {
            Write-Warning "Could not parse existing VS Code settings, creating new configuration"
            $settings = @{}
        }
    }
    
    # Add Terraform AzureRM Provider specific settings
    $terraformSettings = @{
        "terraform_azurerm_provider_mode" = $true
        "terraform_azurerm_ai_enhanced" = $true
        "terraform_azurerm_installation_date" = (Get-Date -Format "yyyy-MM-dd HH:mm:ss")
        "github.copilot.advanced" = @{
            "debug.overrideEngine" = "copilot-chat"
            "debug.testOverrideProxyUrl" = "https://api.github.com"
        }
        "github.copilot.chat.localeOverride" = "auto"
        "github.copilot.enable" = @{
            "*" = $true
            "plaintext" = $false
            "markdown" = $true
            "scminput" = $false
            "go" = $true
            "terraform" = $true
            "hcl" = $true
        }
    }
    
    # Merge settings
    foreach ($key in $terraformSettings.Keys) {
        $settings[$key] = $terraformSettings[$key]
    }
    
    # Write updated settings
    try {
        $settings | ConvertTo-Json -Depth 10 | Set-Content $settingsPath -Encoding UTF8
        Write-Host "Updated VS Code settings.json"
    } catch {
        Write-Warning "Failed to update VS Code settings: $($_.Exception.Message)"
    }
}

# Install files
function Install-Files {
    $manifest = Read-Manifest
    
    # Ensure VS Code directories exist
    New-Item -ItemType Directory -Path $vscodeDir -Force | Out-Null
    New-Item -ItemType Directory -Path $instructionsDestDir -Force | Out-Null
    New-Item -ItemType Directory -Path $promptsDestDir -Force | Out-Null
    
    # Collect files that will be overwritten for backup
    $filesToBackup = @()
    
    # Check instruction files
    foreach ($file in $manifest.InstructionFiles) {
        $destPath = "$instructionsDestDir\$file"
        if (Test-Path $destPath) { $filesToBackup += $destPath }
    }
    
    # Check prompt files  
    foreach ($file in $manifest.PromptFiles) {
        $destPath = "$promptsDestDir\$file"
        if (Test-Path $destPath) { $filesToBackup += $destPath }
    }
    
    # Check main files
    foreach ($file in $manifest.MainFiles) {
        $destPath = "$vscodeDir\$file"
        if (Test-Path $destPath) { $filesToBackup += $destPath }
    }
    
    # Check settings.json
    if (Test-Path $settingsPath) { $filesToBackup += $settingsPath }
    
    # Backup existing files
    if ($filesToBackup.Count -gt 0) {
        Write-Host "Creating backup of existing files..."
        Backup-Files $filesToBackup
    }
    
    # Copy instruction files to instructions/terraform-azurerm/
    foreach ($file in $manifest.InstructionFiles) {
        $sourcePath = "$instructionsDir\$file"
        $destPath = "$instructionsDestDir\$file"
        
        if (Test-Path $sourcePath) {
            Copy-Item $sourcePath $destPath -Force
            Write-Host "Installed: instructions/terraform-azurerm/$file"
        } else {
            Write-Warning "Source file not found: $file"
        }
    }
    
    # Copy prompt files to prompts/
    foreach ($file in $manifest.PromptFiles) {
        $sourcePath = "$promptsDir\$file"
        $destPath = "$promptsDestDir\$file"
        
        if (Test-Path $sourcePath) {
            Copy-Item $sourcePath $destPath -Force
            Write-Host "Installed: prompts/$file"
        } else {
            Write-Warning "Source file not found: $file"
        }
    }
    
    # Copy main files to root User directory
    foreach ($file in $manifest.MainFiles) {
        $sourcePath = "$RepoPath\.github\$file"
        $destPath = "$vscodeDir\$file"
        
        if (Test-Path $sourcePath) {
            Copy-Item $sourcePath $destPath -Force
            Write-Host "Installed: $file"
        } else {
            Write-Warning "Source file not found: $file"
        }
    }
    
    # Update settings.json
    Update-SettingsJson
    
    Write-Host "`nInstallation completed successfully!"
    Write-Host "Backup created at: $backupDir"
}

# Verify installation
function Test-Installation {
    $manifest = Read-Manifest
    $totalFiles = 0
    $installedFiles = 0
    
    Write-Host "Terraform AzureRM Provider AI Setup - Installation Status"
    Write-Host ("=" * 57)
    
    Write-Host "`nInstruction Files:"
    foreach ($file in $manifest.InstructionFiles) {
        $destPath = "$instructionsDestDir\$file"
        $totalFiles++
        if (Test-Path $destPath) {
            Write-Host "  [OK] instructions/terraform-azurerm/$file" -ForegroundColor Green
            $installedFiles++
        } else {
            Write-Host "  [ERROR] instructions/terraform-azurerm/$file" -ForegroundColor Red
        }
    }
    
    Write-Host "`nPrompt Files:"
    foreach ($file in $manifest.PromptFiles) {
        $destPath = "$promptsDestDir\$file"
        $totalFiles++
        if (Test-Path $destPath) {
            Write-Host "  [OK] prompts/$file" -ForegroundColor Green
            $installedFiles++
        } else {
            Write-Host "  [ERROR] prompts/$file" -ForegroundColor Red
        }
    }
    
    Write-Host "`nMain Files:"
    foreach ($file in $manifest.MainFiles) {
        $destPath = "$vscodeDir\$file"
        $totalFiles++
        if (Test-Path $destPath) {
            Write-Host "  [OK] $file" -ForegroundColor Green
            $installedFiles++
        } else {
            Write-Host "  [ERROR] $file" -ForegroundColor Red
        }
    }
    
    Write-Host "`nConfiguration:"
    if (Test-Path $settingsPath) {
        try {
            $settings = Get-Content $settingsPath -Raw | ConvertFrom-Json
            if ($settings.terraform_azurerm_provider_mode -eq $true) {
                Write-Host "  [OK] VS Code settings.json configured" -ForegroundColor Green
                $installedFiles++
            } else {
                Write-Host "  [ERROR] VS Code settings.json not configured for Terraform AzureRM" -ForegroundColor Red
            }
            $totalFiles++
        } catch {
            Write-Host "  [ERROR] VS Code settings.json invalid JSON" -ForegroundColor Red
            $totalFiles++
        }
    } else {
        Write-Host "  [ERROR] VS Code settings.json not found" -ForegroundColor Red
        $totalFiles++
    }
    
    $percentage = [math]::Round(($installedFiles / $totalFiles) * 100, 1)
    Write-Host ""
    Write-Host "Installation Status: $installedFiles of $totalFiles files ($percentage%)"
    
    if ($installedFiles -eq $totalFiles) {
        Write-Host "Status: Complete [OK]" -ForegroundColor Green
    } else {
        Write-Host "Status: Incomplete [ERROR]" -ForegroundColor Yellow
    }
}

# Clean installation
function Remove-Installation {
    $manifest = Read-Manifest
    
    Write-Host "Removing AI setup files..."
    
    # Remove installed files
    foreach ($file in ($manifest.InstructionFiles + $manifest.PromptFiles + $manifest.MainFiles)) {
        $destPath = "$vscodeDir\$file"
        if (Test-Path $destPath) {
            Remove-Item $destPath -Force
            Write-Host "Removed: $file"
        }
    }
    
    # Restore from most recent backup
    $backupDirs = Get-ChildItem "$vscodeDirPlatform" -Directory | Where-Object { $_.Name -match '^\.vscode\\ai-backup-\d{8}$' } | Sort-Object Name -Descending
    
    if ($backupDirs.Count -gt 0) {
        $latestBackup = $backupDirs[0].FullName
        Write-Host "`nRestoring from backup: $latestBackup"
        
        Get-ChildItem $latestBackup -Recurse -File | ForEach-Object {
            $relativePath = $_.FullName.Replace($latestBackup, '').TrimStart('\', '/')
            $restorePath = "$vscodeDir\$relativePath"
            $restoreParent = Split-Path $restorePath -Parent
            
            if ($restoreParent) {
                New-Item -ItemType Directory -Path $restoreParent -Force | Out-Null
            }
            
            Copy-Item $_.FullName $restorePath -Force
            Write-Host "Restored: $relativePath"
        }
    }
    
    Write-Host "`nCleanup completed successfully!"
}

# Main execution
if ($Verify) {
    Test-Installation
}
elseif ($Clean) {
    Remove-Installation
}
else {
    Install-Files
}
