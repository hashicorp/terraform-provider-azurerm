<#
.SYNOPSIS
    VS Code configuration and settings management for Terraform AzureRM Provider AI Setup

.DESCRIPTION
    Handles VS Code settings.json configuration, extension management, and workspace setup
    for the AI-powered development environment.
#>

function Get-VSCodeUserSettingsPath {
    if ($IsWindows -or $env:OS -eq "Windows_NT") {
        return Join-Path $env:APPDATA "Code\User\settings.json"
    } elseif ($IsLinux) {
        return Join-Path $env:HOME ".config/Code/User/settings.json"
    } elseif ($IsMacOS) {
        return Join-Path $env:HOME "Library/Application Support/Code/User/settings.json"
    } else {
        Write-Warning "Unknown operating system, trying Windows path"
        return Join-Path $env:APPDATA "Code\User\settings.json"
    }
}

function Test-PreviousInstallation {
    param([string]$SettingsPath)
    
    if (-not (Test-Path $SettingsPath)) {
        return $false
    }
    
    $settingsContent = Get-Content $SettingsPath -Raw -ErrorAction SilentlyContinue
    if (-not $settingsContent) {
        return $false
    }
    
    try {
        $settings = $settingsContent | ConvertFrom-Json
        return $settings.PSObject.Properties.Name -contains "github.copilot.chat.welcomeMessage"
    }
    catch {
        return $false
    }
}

function Backup-VSCodeSettings {
    param(
        [string]$SettingsPath,
        [string]$BackupDir
    )
    
    Import-Module (Join-Path $PSScriptRoot "core-functions.ps1") -Force
    
    if (-not (Test-Path $SettingsPath)) {
        Write-Host "[INFO] No existing VS Code settings to backup" -ForegroundColor Yellow
        return $null
    }
    
    $backupPath = New-SafeBackup -FilePath $SettingsPath -BackupDir $BackupDir
    if ($backupPath) {
        Write-Host "[SUCCESS] VS Code settings backed up to: $backupPath" -ForegroundColor Green
    }
    
    return $backupPath
}

function Install-VSCodeSettings {
    param(
        [string]$SettingsPath,
        [hashtable]$NewSettings
    )
    
    Write-Host "[INFO] Installing VS Code settings..." -ForegroundColor Blue
    
    $existingSettings = @{}
    if (Test-Path $SettingsPath) {
        try {
            $settingsContent = Get-Content $SettingsPath -Raw
            $existingSettings = $settingsContent | ConvertFrom-Json -AsHashtable
        }
        catch {
            Write-Warning "Could not parse existing settings.json, creating new one"
            $existingSettings = @{}
        }
    }
    
    # Merge settings (new settings take precedence)
    foreach ($key in $NewSettings.Keys) {
        $existingSettings[$key] = $NewSettings[$key]
    }
    
    try {
        $settingsDir = Split-Path $SettingsPath -Parent
        if (-not (Test-Path $settingsDir)) {
            New-Item -ItemType Directory -Path $settingsDir -Force | Out-Null
        }
        
        $existingSettings | ConvertTo-Json -Depth 10 | Set-Content $SettingsPath -Encoding UTF8
        Write-Host "[SUCCESS] VS Code settings installed" -ForegroundColor Green
        return $true
    }
    catch {
        Write-Error "Failed to install VS Code settings: $_"
        return $false
    }
}

function Get-CopilotSettings {
    return @{
        "github.copilot.enable" = @{
            "*" = $true
            "yaml" = $true
            "plaintext" = $true
            "markdown" = $true
            "go" = $true
            "powershell" = $true
            "shellscript" = $true
        }
        "github.copilot.chat.welcomeMessage" = "never"
        "github.copilot.advanced" = @{
            "secret_key" = "terraform-azurerm-provider"
            "length" = 1000
            "temperature" = ""
            "top_p" = ""
            "stops" = @{
                "*" = @(
                    "class",
                    "def",
                    "if",
                    "else",
                    "elif",
                    "except",
                    "finally"
                )
                "go" = @(
                    "func",
                    "type",
                    "var",
                    "const",
                    "package",
                    "import"
                )
            }
        }
        "github.copilot.chat.localeOverride" = "en"
        "github.copilot.renameSuggestions.triggerAutomatically" = $true
    }
}

function Remove-CopilotSettings {
    param([string]$SettingsPath)
    
    Write-Host "[INFO] Removing Copilot settings from VS Code..." -ForegroundColor Blue
    
    if (-not (Test-Path $SettingsPath)) {
        Write-Host "[INFO] No settings file found" -ForegroundColor Yellow
        return $true
    }
    
    try {
        $settingsContent = Get-Content $SettingsPath -Raw
        $settings = $settingsContent | ConvertFrom-Json -AsHashtable
        
        # Remove Copilot-related settings
        $copilotKeys = @(
            "github.copilot.enable",
            "github.copilot.chat.welcomeMessage",
            "github.copilot.advanced",
            "github.copilot.chat.localeOverride",
            "github.copilot.renameSuggestions.triggerAutomatically"
        )
        
        foreach ($key in $copilotKeys) {
            if ($settings.ContainsKey($key)) {
                $settings.Remove($key)
                Write-Host "[INFO] Removed setting: $key" -ForegroundColor Yellow
            }
        }
        
        $settings | ConvertTo-Json -Depth 10 | Set-Content $SettingsPath -Encoding UTF8
        Write-Host "[SUCCESS] Copilot settings removed" -ForegroundColor Green
        return $true
    }
    catch {
        Write-Error "Failed to remove Copilot settings: $_"
        return $false
    }
}

# VS Code functions are loaded via dot sourcing - no export needed
