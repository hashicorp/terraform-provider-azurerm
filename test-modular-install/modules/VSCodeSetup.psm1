# VSCodeSetup.psm1 - VS Code configuration management for Terraform AzureRM Provider AI Setup

function Get-VSCodeUserSettingsPath {
    <#
    .SYNOPSIS
    Gets the path to VS Code user settings.json
    
    .OUTPUTS
    String path to settings.json file
    #>
    
    $vsCodeConfigPath = switch ($true) {
        $IsWindows { 
            Join-Path $env:APPDATA "Code\User\settings.json" 
        }
        $IsMacOS { 
            Join-Path $env:HOME "Library/Application Support/Code/User/settings.json" 
        }
        $IsLinux { 
            Join-Path $env:HOME ".config/Code/User/settings.json" 
        }
        default { 
            # Fallback for older PowerShell versions
            Join-Path $env:APPDATA "Code\User\settings.json" 
        }
    }
    
    return $vsCodeConfigPath
}

function Test-PreviousInstallation {
    <#
    .SYNOPSIS
    Tests if AI setup has been previously installed
    
    .PARAMETER RepositoryPath
    Path to the repository to check
    
    .OUTPUTS
    Boolean indicating if previous installation exists
    #>
    
    param(
        [Parameter(Mandatory = $true)]
        [string]$RepositoryPath
    )
    
    $instructionsPath = Join-Path $RepositoryPath ".github\copilot-instructions.md"
    return Test-Path $instructionsPath
}

function Backup-VSCodeSettings {
    <#
    .SYNOPSIS
    Creates a backup of VS Code settings
    
    .PARAMETER BackupDirectory
    Directory where backup should be stored
    
    .OUTPUTS
    String path to backup file, or $null if no settings found
    #>
    
    param(
        [Parameter(Mandatory = $true)]
        [string]$BackupDirectory
    )
    
    $settingsPath = Get-VSCodeUserSettingsPath
    
    if (-not (Test-Path $settingsPath)) {
        Write-Verbose "No VS Code settings found to backup"
        return $null
    }
    
    if (-not (Test-Path $BackupDirectory)) {
        New-Item -ItemType Directory -Path $BackupDirectory -Force | Out-Null
    }
    
    $timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
    $backupFileName = "$timestamp-vscode-settings.json"
    $backupPath = Join-Path $BackupDirectory $backupFileName
    
    Copy-Item $settingsPath $backupPath -Force
    Write-Verbose "VS Code settings backed up to: $backupPath"
    
    return $backupPath
}

function Install-VSCodeSettings {
    <#
    .SYNOPSIS
    Installs or updates VS Code settings for optimal AI development
    
    .PARAMETER Force
    Force installation even if settings already exist
    
    .OUTPUTS
    Boolean indicating success
    #>
    
    param(
        [switch]$Force
    )
    
    $settingsPath = Get-VSCodeUserSettingsPath
    $settingsDir = Split-Path $settingsPath -Parent
    
    # Ensure settings directory exists
    if (-not (Test-Path $settingsDir)) {
        New-Item -ItemType Directory -Path $settingsDir -Force | Out-Null
    }
    
    # Read existing settings or create new
    $settings = @{}
    if ((Test-Path $settingsPath) -and -not $Force) {
        try {
            $existingContent = Get-Content $settingsPath -Raw
            $settings = $existingContent | ConvertFrom-Json -AsHashtable
        }
        catch {
            Write-Warning "Could not parse existing settings, creating new configuration"
            $settings = @{}
        }
    }
    
    # Add AI-optimized settings
    $aiSettings = @{
        "github.copilot.enable" = @{
            "*" = $true
            "markdown" = $true
            "go" = $true
            "powershell" = $true
        }
        "github.copilot.advanced" = @{
            "inlineSuggestEnable" = $true
        }
        "editor.suggestSelection" = "first"
        "editor.tabCompletion" = "on"
        "editor.quickSuggestions" = @{
            "comments" = $true
            "strings" = $true
            "other" = $true
        }
        "go.formatTool" = "goimports"
        "go.lintTool" = "golangci-lint"
    }
    
    # Merge settings
    foreach ($key in $aiSettings.Keys) {
        $settings[$key] = $aiSettings[$key]
    }
    
    # Write settings
    $settingsJson = $settings | ConvertTo-Json -Depth 10
    Set-Content -Path $settingsPath -Value $settingsJson -Encoding UTF8
    
    Write-Verbose "VS Code settings updated successfully"
    return $true
}

function Get-CopilotSettings {
    <#
    .SYNOPSIS
    Generates GitHub Copilot workspace settings for the repository
    
    .OUTPUTS
    Hashtable containing Copilot workspace configuration
    #>
    
    return @{
        "github.copilot.enable" = @{
            "*" = $true
            "markdown" = $true
            "go" = $true
            "yaml" = $true
            "json" = $true
            "powershell" = $true
        }
        "github.copilot.advanced" = @{
            "inlineSuggestEnable" = $true
            "length" = 500
        }
        "github.copilot.chat.localeOverride" = "en"
        "github.copilot.chat.welcomeMessage" = "never"
    }
}

# Export module members
Export-ModuleMember -Function Get-VSCodeUserSettingsPath, Test-PreviousInstallation, Backup-VSCodeSettings, Install-VSCodeSettings, Get-CopilotSettings
