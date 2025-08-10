# Installation Detection Module
# Functions for detecting previous installations and managing installation state

function Test-PreviousInstallation {
    <#
    .SYNOPSIS
        Detects previous installation (including partial/failed states)
    #>
    param([string]$RepositoryPath)
    
    $installationMarkers = @{
        VSCodeSettings = Join-Path $env:APPDATA "Code\User\settings.json"
        CopilotInstructions = Join-Path $RepositoryPath ".github\copilot-instructions.md"
        InstructionsDir = Join-Path $RepositoryPath ".github\instructions"
        BackupMarker = Join-Path $env:APPDATA "Code\User\.terraform-azurerm-backups"
    }
    
    $installationState = @{
        HasVSCodeSettings = $false
        HasCopilotInstructions = $false
        HasInstructions = $false
        HasBackups = $false
        IsPartialInstall = $false
        BackupTimestamp = $null
    }
    
    # Check VS Code settings for our markers
    if (Test-Path $installationMarkers.VSCodeSettings) {
        try {
            $settings = Get-Content $installationMarkers.VSCodeSettings -Raw
            $installationState.HasVSCodeSettings = $settings -match '"terraform_azurerm_'
        } catch {
            Write-StatusMessage "Could not read VS Code settings" "Warning"
        }
    }
    
    # Check for copilot instructions
    if (Test-Path $installationMarkers.CopilotInstructions) {
        $installationState.HasCopilotInstructions = Test-FileIntegrity -FilePath $installationMarkers.CopilotInstructions
    }
    
    # Check for instructions directory
    if (Test-Path $installationMarkers.InstructionsDir) {
        $instructionFiles = Get-ChildItem -Path $installationMarkers.InstructionsDir -Filter "*.instructions.md" -ErrorAction SilentlyContinue
        $installationState.HasInstructions = $instructionFiles.Count -gt 0
    }
    
    # Check for backup directory
    if (Test-Path $installationMarkers.BackupMarker) {
        $installationState.HasBackups = $true
        try {
            $backupInfo = Get-Item $installationMarkers.BackupMarker
            $installationState.BackupTimestamp = $backupInfo.LastWriteTime
        } catch {
            # Ignore timestamp errors
        }
    }
    
    # Determine if this is a partial installation
    $hasAnyMarkers = $installationState.HasVSCodeSettings -or 
                     $installationState.HasCopilotInstructions -or 
                     $installationState.HasInstructions -or 
                     $installationState.HasBackups
    
    $hasAllMarkers = $installationState.HasVSCodeSettings -and 
                     $installationState.HasCopilotInstructions -and 
                     $installationState.HasInstructions
    
    $installationState.IsPartialInstall = $hasAnyMarkers -and -not $hasAllMarkers
    
    return $installationState
}

function Get-InstallationPaths {
    <#
    .SYNOPSIS
        Gets standard installation paths
    #>
    param([string]$RepositoryPath)
    
    return @{
        VSCodeSettings = Join-Path $env:APPDATA "Code\User\settings.json"
        VSCodeBackupDir = Join-Path $env:APPDATA "Code\User\.terraform-azurerm-backups"
        CopilotInstructions = Join-Path $RepositoryPath ".github\copilot-instructions.md"
        InstructionsSource = Join-Path $RepositoryPath ".github\instructions"
        InstructionsTarget = Join-Path $RepositoryPath ".github\instructions"
    }
}

function Test-InstallationHealth {
    <#
    .SYNOPSIS
        Tests the health of an existing installation
    #>
    param([string]$RepositoryPath)
    
    $paths = Get-InstallationPaths -RepositoryPath $RepositoryPath
    $health = @{
        IsHealthy = $true
        Issues = @()
        Warnings = @()
    }
    
    # Check file integrity
    if (Test-Path $paths.CopilotInstructions) {
        if (-not (Test-FileIntegrity -FilePath $paths.CopilotInstructions)) {
            $health.Issues += "Copilot instructions file is corrupted"
            $health.IsHealthy = $false
        }
    } else {
        $health.Issues += "Copilot instructions file is missing"
        $health.IsHealthy = $false
    }
    
    # Check instructions directory
    if (Test-Path $paths.InstructionsSource) {
        $instructionFiles = Get-ChildItem -Path $paths.InstructionsSource -Filter "*.instructions.md"
        if ($instructionFiles.Count -eq 0) {
            $health.Issues += "No instruction files found"
            $health.IsHealthy = $false
        }
    } else {
        $health.Issues += "Instructions directory is missing"
        $health.IsHealthy = $false
    }
    
    # Check VS Code settings
    if (Test-Path $paths.VSCodeSettings) {
        try {
            $settings = Get-Content $paths.VSCodeSettings -Raw
            if ($settings -notmatch '"terraform_azurerm_') {
                $health.Warnings += "VS Code settings may not be properly configured"
            }
        } catch {
            $health.Warnings += "Could not verify VS Code settings"
        }
    }
    
    return $health
}
