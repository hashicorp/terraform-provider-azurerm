# AI Installation Module
# Functions for installing GitHub Copilot instructions and VS Code configuration

function Install-CopilotInstructions {
    <#
    .SYNOPSIS
        Installs GitHub Copilot instruction files
    #>
    param([string]$RepositoryPath)
    
    $sourcePath = Join-Path $RepositoryPath ".github\instructions"
    $targetPath = Join-Path $RepositoryPath ".github\instructions"
    $copilotInstructionsPath = Join-Path $RepositoryPath ".github\copilot-instructions.md"
    
    Write-StatusMessage "Installing GitHub Copilot instructions..." "Info"
    
    # Ensure the instructions are in place (they should be part of the repository)
    if (-not (Test-Path $sourcePath)) {
        Write-StatusMessage "Instructions directory not found: $sourcePath" "Error"
        return $false
    }
    
    # Count instruction files
    $instructionFiles = Get-ChildItem -Path $sourcePath -Filter "*.instructions.md" -ErrorAction SilentlyContinue
    if ($instructionFiles.Count -eq 0) {
        Write-StatusMessage "No instruction files found in $sourcePath" "Error"
        return $false
    }
    
    Write-StatusMessage "Found $($instructionFiles.Count) instruction files" "Success"
    
    # Check for copilot-instructions.md
    if (Test-Path $copilotInstructionsPath) {
        if (Test-FileIntegrity -FilePath $copilotInstructionsPath) {
            Write-StatusMessage "Copilot instructions file verified" "Success"
        } else {
            Write-StatusMessage "Copilot instructions file exists but may be corrupted" "Warning"
        }
    } else {
        Write-StatusMessage "Copilot instructions file not found: $copilotInstructionsPath" "Warning"
    }
    
    return $true
}

function Install-VSCodeSettings {
    <#
    .SYNOPSIS
        Installs VS Code settings for AI development
    #>
    param([string]$RepositoryPath, [switch]$Force)
    
    $settingsPath = Join-Path $env:APPDATA "Code\User\settings.json"
    $backupDir = Join-Path $env:APPDATA "Code\User\.terraform-azurerm-backups"
    
    Write-StatusMessage "Configuring VS Code settings..." "Info"
    
    # Create backup if settings exist
    if (Test-Path $settingsPath) {
        if (-not $Force) {
            $timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
            $backupPath = Join-Path $backupDir "settings_backup_$timestamp.json"
            
            if (-not (New-SafeBackup -SourcePath $settingsPath -BackupPath $backupPath)) {
                Write-StatusMessage "Failed to create backup of VS Code settings" "Error"
                return $false
            }
        }
    }
    
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
            Write-StatusMessage "Could not parse existing VS Code settings, creating new configuration" "Warning"
            $settings = @{}
        }
    }
    
    # Add Terraform AzureRM Provider specific settings
    $terraformSettings = @{
        "terraform_azurerm_provider_mode" = $true
        "terraform_azurerm_ai_enhanced" = $true
        "terraform_azurerm_installation_date" = (Get-Date -Format "yyyy-MM-dd HH:mm:ss")
        "terraform_azurerm_backup_length" = 0
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
    
    # Write settings back
    try {
        $settings | ConvertTo-Json -Depth 10 | Set-Content -Path $settingsPath -Encoding UTF8
        Write-StatusMessage "VS Code settings updated successfully" "Success"
        return $true
    } catch {
        Write-StatusMessage "Failed to update VS Code settings: $_" "Error"
        return $false
    }
}

function Install-AIAgent {
    <#
    .SYNOPSIS
        Complete AI agent installation
    #>
    param([string]$RepositoryPath, [switch]$Force)
    
    Write-StatusMessage "Starting AI agent installation..." "Info"
    
    # Install Copilot instructions
    if (-not (Install-CopilotInstructions -RepositoryPath $RepositoryPath)) {
        Write-StatusMessage "Failed to install Copilot instructions" "Error"
        return $false
    }
    
    # Install VS Code settings
    if (-not (Install-VSCodeSettings -RepositoryPath $RepositoryPath -Force:$Force)) {
        Write-StatusMessage "Failed to install VS Code settings" "Error"
        return $false
    }
    
    # Create installation marker
    $backupDir = Join-Path $env:APPDATA "Code\User\.terraform-azurerm-backups"
    if (-not (Test-Path $backupDir)) {
        New-Item -ItemType Directory -Path $backupDir -Force | Out-Null
    }
    
    Write-StatusMessage "AI agent installation completed successfully!" "Success"
    return $true
}
