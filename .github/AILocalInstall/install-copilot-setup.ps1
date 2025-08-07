<#
.SYNOPSIS
    Terraform AzureRM Provider AI Setup - PowerShell Installation Script

.DESCRIPTION
    Installs GitHub Copilot instruction files and AI prompts for enhanced development
    of the Terraform AzureRM Provider. This script configures VS Code with specialized
    AI instructions, coding patterns, and best practices for Azure resource development.

.PARAMETER RepositoryPath
    Specifies the path to the terraform-provider-azurerm repository.
    If not provided, the script will attempt to auto-discover the repository location.

.PARAMETER Clean
    Removes all installed AI setup files and restores original VS Code settings from backups.

.PARAMETER Help
    Displays detailed help information and usage examples.

.EXAMPLE
    .\install-copilot-setup.ps1
    Auto-discovers the repository and installs AI setup with default settings.

.EXAMPLE
    .\install-copilot-setup.ps1 -RepositoryPath "C:\Projects\terraform-provider-azurerm"
    Installs AI setup using the specified repository path.

.EXAMPLE
    .\install-copilot-setup.ps1 -Clean
    Removes all AI setup and restores original VS Code settings.

.NOTES
    File Name      : install-copilot-setup.ps1
    Author         : HashiCorp Terraform AzureRM Provider Team
    Prerequisite   : PowerShell 5.1+, VS Code
    Version        : 1.0.0
    
    This script creates backups of existing VS Code settings and provides
    a clean uninstall option to restore the original configuration.
#>

param(
    [string]$RepositoryPath,
    [switch]$Clean,
    [switch]$Help
)

if ($Help) {
    Write-Host "Terraform AzureRM Provider AI Setup" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "USAGE:" -ForegroundColor Yellow
    Write-Host "  .\install-copilot-setup.ps1 [OPTIONS]"
    Write-Host ""
    Write-Host "OPTIONS:" -ForegroundColor Yellow
    Write-Host "  -RepositoryPath <path>    Path to terraform-provider-azurerm repository"
    Write-Host "  -Clean                    Remove all installed files and restore backups"
    Write-Host "  -Help                     Show this help message"
    Write-Host ""
    Write-Host "EXAMPLES:" -ForegroundColor Yellow
    Write-Host "  .\install-copilot-setup.ps1                                  # Auto-discover repository"
    Write-Host "  .\install-copilot-setup.ps1 -RepositoryPath C:\path\to\repo  # Use specific path"
    Write-Host "  .\install-copilot-setup.ps1 -Clean                           # Remove installation"
    Write-Host ""
    Write-Host "AI FEATURES:" -ForegroundColor Green
    Write-Host "  - Context-aware code generation and review"
    Write-Host "  - Azure-specific implementation patterns"
    Write-Host "  - Testing guidelines and best practices"
    Write-Host "  - Documentation standards enforcement"
    Write-Host "  - Error handling and debugging assistance"
    Write-Host ""
    exit 0
}

Write-Host "Terraform AzureRM Provider AI Setup" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan

# Function to find terraform-provider-azurerm repository
function Find-RepositoryRoot {
    param([string]$StartPath = (Get-Location).Path)
    
    $currentPath = $StartPath
    while ($currentPath -and $currentPath -ne [System.IO.Path]::GetPathRoot($currentPath)) {
        $gitPath = Join-Path $currentPath ".git"
        $goModPath = Join-Path $currentPath "go.mod"
        
        # Check if this is the terraform-provider-azurerm repository
        if ((Test-Path $gitPath) -and (Test-Path $goModPath)) {
            $goModContent = Get-Content $goModPath -Raw -ErrorAction SilentlyContinue
            if ($goModContent -match "module github\.com/hashicorp/terraform-provider-azurerm") {
                return $currentPath
            }
        }
        
        $currentPath = Split-Path $currentPath -Parent
    }
    
    return $null
}

# Function to validate file content integrity
function Test-FileIntegrity {
    param([string]$FilePath, [string]$ExpectedPattern = "")
    
    if (-not (Test-Path $FilePath)) {
        return $false
    }
    
    try {
        $content = Get-Content $FilePath -Raw -ErrorAction Stop
        
        # Basic length check (minimum viable content)
        if ($content.Length -lt 50) {
            return $false
        }
        
        # Pattern check if provided
        if ($ExpectedPattern -and $content -notmatch $ExpectedPattern) {
            return $false
        }
        
        return $true
    } catch {
        return $false
    }
}

# Function to create safe backup with integrity verification
function Create-SafeBackup {
    param([string]$SourcePath, [string]$BackupPath)
    
    if (-not (Test-Path $SourcePath)) {
        return $false
    }
    
    try {
        # Create backup
        Copy-Item -Path $SourcePath -Destination $BackupPath -Force -ErrorAction Stop
        
        # Verify backup integrity
        if (-not (Test-FileIntegrity -FilePath $BackupPath)) {
            Write-Host "ERROR: Backup verification failed for $BackupPath" -ForegroundColor Red
            if (Test-Path $BackupPath) {
                Remove-Item -Path $BackupPath -Force -ErrorAction SilentlyContinue
            }
            return $false
        }
        
        # Store backup metadata for integrity tracking (using local variable now)
        $sourceLength = (Get-Item $SourcePath).Length
        
        Write-Host "Backup created and verified: $(Split-Path $BackupPath -Leaf) ($sourceLength bytes)" -ForegroundColor Green
        return $true
        
    } catch {
        Write-Host "ERROR: Failed to create backup: $($_.Exception.Message)" -ForegroundColor Red
        return $false
    }
}

# Function to detect previous installation (including partial/failed states)
function Test-PreviousInstallation {
    param([string]$UserDir)
    
    $foundIndicators = @()
    
    # Check for instruction directory (even if empty)
    $instructionsDir = "$UserDir\instructions\terraform-azurerm"
    if (Test-Path $instructionsDir) {
        $instructionFiles = Get-ChildItem -Path $instructionsDir -File -ErrorAction SilentlyContinue
        if ($instructionFiles.Count -gt 0) {
            $foundIndicators += $instructionsDir
        } else {
            $foundIndicators += "$instructionsDir (empty directory - partial cleanup)"
        }
    }
    
    # Also check for empty parent instructions directory (from partial cleanup)
    $parentInstructionsDir = "$UserDir\instructions"
    if (Test-Path $parentInstructionsDir) {
        $parentItems = Get-ChildItem -Path $parentInstructionsDir -ErrorAction SilentlyContinue
        if ($parentItems.Count -eq 0) {
            $foundIndicators += "$parentInstructionsDir (empty directory - needs cleanup)"
        }
    }
    
    # Check for main copilot instructions
    $copilotInstructions = "$UserDir\copilot-instructions.md"
    if (Test-Path $copilotInstructions) {
        $foundIndicators += $copilotInstructions
    }
    
    # Check for prompts directory with any AzureRM-related files
    $promptsDir = "$UserDir\prompts"
    if (Test-Path $promptsDir) {
        # First check for obvious filename matches
        $azureRMPrompts = @(Get-ChildItem -Path $promptsDir -Filter "*azurerm*" -File -ErrorAction SilentlyContinue)
        $terraformPrompts = @(Get-ChildItem -Path $promptsDir -Filter "*terraform*" -File -ErrorAction SilentlyContinue)
        
        # Also check for content-based matches in .prompt.md files
        $allPromptFiles = @(Get-ChildItem -Path $promptsDir -Filter "*.prompt.md" -File -ErrorAction SilentlyContinue)
        foreach ($promptFile in $allPromptFiles) {
            try {
                $content = Get-Content $promptFile.FullName -Head 10 -ErrorAction SilentlyContinue | Out-String
                if ($content -match "terraform.*azurerm|azurerm.*provider|terraform.*provider.*azurerm") {
                    $azureRMPrompts += $promptFile
                }
            } catch {
                # If we can't read the file, skip it
            }
        }
        
        $allPrompts = $azureRMPrompts + $terraformPrompts | Sort-Object FullName | Get-Unique
        foreach ($prompt in $allPrompts) {
            if ($foundIndicators -notcontains $prompt.FullName) {
                $foundIndicators += $prompt.FullName
            }
        }
    }
    
    return @{
        HasPrevious = ($foundIndicators.Count -gt 0)
        FoundFiles = $foundIndicators
    }
}

# Function to extract backup length from settings.json
function Get-BackupLengthFromSettings {
    param([string]$SettingsPath)
    
    if (-not (Test-Path $SettingsPath)) {
        return $null
    }
    
    try {
        $content = Get-Content $SettingsPath -Raw -ErrorAction Stop
        $settings = $content | ConvertFrom-Json -ErrorAction Stop
        
        if ($settings.PSObject.Properties.Name -contains "// AZURERM_BACKUP_LENGTH") {
            return $settings."// AZURERM_BACKUP_LENGTH"
        }
    } catch {
        # If we can't parse JSON or find the field, return null
        return $null
    }
    
    return $null
}

# Function to get the most recent backup file
function Get-MostRecentBackup {
    param([string]$UserDir)
    
    try {
        $backupFiles = Get-ChildItem "$UserDir\settings.json.backup.*" -ErrorAction SilentlyContinue
        if ($backupFiles) {
            $mostRecent = ($backupFiles | Sort-Object LastWriteTime -Descending | Select-Object -First 1).FullName
            return $mostRecent
        }
    } catch {
        # If we can't access backup files, return null
        return $null
    }
    
    return $null
}

# Enhanced clean mode with multiple restoration scenarios
if ($Clean) {
    Write-Host "Starting cleanup process..." -ForegroundColor Yellow
    
    $userDir = "$env:APPDATA\Code\User"
    $instructionsDir = "$userDir\instructions\terraform-azurerm"
    $promptsDir = "$userDir\prompts"
    $settingsPath = "$userDir\settings.json"
    
    # Find the most recent backup file with timestamp pattern
    $backupFiles = Get-ChildItem "$userDir\settings.json.backup.*" -ErrorAction SilentlyContinue
    $backupPath = if ($backupFiles) {
        ($backupFiles | Sort-Object LastWriteTime -Descending | Select-Object -First 1).FullName
    } else {
        $null
    }
    
    # Detect what needs to be cleaned (with aggressive detection)
    $previousInstall = Test-PreviousInstallation -UserDir $userDir
    
    # Also check for AzureRM traces in settings.json
    $hasSettingsTraces = $false
    $hasEmptySettings = $false
    if (Test-Path $settingsPath) {
        try {
            $settingsContent = Get-Content $settingsPath -Raw -ErrorAction SilentlyContinue
            if ($settingsContent -match "AZURERM|terraform-azurerm|copilot.*instructions") {
                $hasSettingsTraces = $true
            }
            # Also check if settings.json is empty (might need cleanup)
            try {
                $settingsJson = $settingsContent | ConvertFrom-Json -ErrorAction Stop
                $propCount = $settingsJson.PSObject.Properties.Count
                if ($settingsJson -and ($propCount -eq $null -or $propCount -eq 0)) {
                    $hasEmptySettings = $true
                }
            } catch {
                # If we can't parse as JSON, check if it's effectively empty (just braces and whitespace)
                $trimmedContent = $settingsContent.Trim() -replace '\s+', ''
                if ($trimmedContent -eq '{}') {
                    $hasEmptySettings = $true
                }
            }
        } catch {
            # Could not read or parse settings.json
        }
    }
    
    if (-not $previousInstall.HasPrevious -and -not $hasSettingsTraces -and -not $hasEmptySettings) {
        Write-Host "No previous installation detected. Nothing to clean." -ForegroundColor Green
        exit 0
    }
    
    if ($previousInstall.HasPrevious) {
        Write-Host "Found previous installation files:" -ForegroundColor Yellow
        foreach ($file in $previousInstall.FoundFiles) {
            Write-Host "  - $file" -ForegroundColor Gray
        }
        Write-Host ""
    }
    
    if ($hasSettingsTraces) {
        Write-Host "Found AzureRM traces in settings.json (will attempt cleanup)" -ForegroundColor Yellow
        Write-Host ""
    }
    
    # Interactive confirmation for manual backup
    if ((Test-Path $settingsPath) -and (-not $backupPath -or -not (Test-Path $backupPath))) {
        Write-Host "WARNING: settings.json exists but no backup found!" -ForegroundColor Red
        Write-Host "This may contain important VS Code settings outside of AzureRM configuration." -ForegroundColor Yellow
        Write-Host ""
        do {
            $response = Read-Host "Continue cleanup? (C)ontinue/(X)exit"
            $response = $response.ToUpper()
        } while ($response -ne "C" -and $response -ne "X")
        
        if ($response -eq "X") {
            Write-Host "Cleanup cancelled by user." -ForegroundColor Yellow
            exit 0
        }
    }
    
    # Remove instruction files
    # Remove instruction files
    if (Test-Path $instructionsDir) {
        Write-Host "Removing instruction files..." -ForegroundColor Red
        try {
            $instructionFiles = Get-ChildItem -Path $instructionsDir -File -ErrorAction SilentlyContinue
            if ($instructionFiles.Count -gt 0) {
                Write-Host "Found $($instructionFiles.Count) instruction files to remove" -ForegroundColor Yellow
            } else {
                Write-Host "Instruction directory is empty (from previous failed cleanup)" -ForegroundColor Yellow
            }
            
            Remove-Item -Path $instructionsDir -Recurse -Force -ErrorAction Stop
            Write-Host "Instructions removed" -ForegroundColor Green
        } catch {
            Write-Host "ERROR: Failed to remove instructions: $($_.Exception.Message)" -ForegroundColor Red
        }
    } else {
        Write-Host "No terraform-azurerm instruction directory found to clean" -ForegroundColor Yellow
    }
    
    # Also check for empty parent instructions directory
    $parentInstructionsDir = "$userDir\instructions"
    if (Test-Path $parentInstructionsDir) {
        try {
            $remainingItems = Get-ChildItem -Path $parentInstructionsDir -ErrorAction SilentlyContinue
            if ($remainingItems.Count -eq 0) {
                Remove-Item -Path $parentInstructionsDir -Force -ErrorAction Stop
                Write-Host "Removed empty instructions directory" -ForegroundColor Green
            } else {
                Write-Host "Instructions directory contains other files - keeping it" -ForegroundColor Yellow
            }
        } catch {
            Write-Host "WARNING: Could not clean empty instructions directory: $($_.Exception.Message)" -ForegroundColor Yellow
        }
    }
    
    # Remove prompt files (only AzureRM prompts)
    if (Test-Path $promptsDir) {
        Write-Host "Removing AI prompt files..." -ForegroundColor Red
        try {
            # First find obvious filename matches
            $azureRMPrompts = @(Get-ChildItem -Path $promptsDir -Filter "*azurerm*" -File -ErrorAction SilentlyContinue)
            $terraformPrompts = @(Get-ChildItem -Path $promptsDir -Filter "*terraform*" -File -ErrorAction SilentlyContinue)
            
            # Also find content-based matches in .prompt.md files
            $allPromptFiles = @(Get-ChildItem -Path $promptsDir -Filter "*.prompt.md" -File -ErrorAction SilentlyContinue)
            foreach ($promptFile in $allPromptFiles) {
                try {
                    $content = Get-Content $promptFile.FullName -Head 10 -ErrorAction SilentlyContinue | Out-String
                    if ($content -match "terraform.*azurerm|azurerm.*provider|terraform.*provider.*azurerm") {
                        $azureRMPrompts += $promptFile
                    }
                } catch {
                    # If we can't read the file, skip it
                }
            }
            
            $allPrompts = $azureRMPrompts + $terraformPrompts | Sort-Object FullName | Get-Unique
            
            if ($allPrompts.Count -gt 0) {
                $allPrompts | Remove-Item -Force -ErrorAction Stop
                Write-Host "Removed $($allPrompts.Count) AzureRM prompt files" -ForegroundColor Green
            } else {
                Write-Host "No AzureRM prompt files found to remove" -ForegroundColor Yellow
            }
            
            # Remove prompts directory if empty
            if ((Get-ChildItem -Path $promptsDir -ErrorAction SilentlyContinue | Measure-Object).Count -eq 0) {
                Remove-Item -Path $promptsDir -Force -ErrorAction SilentlyContinue
                Write-Host "Removed empty prompts directory" -ForegroundColor Green
            }
            
        } catch {
            Write-Host "ERROR: Failed to remove prompts: $($_.Exception.Message)" -ForegroundColor Red
        }
    } else {
        Write-Host "No prompt directory found to clean" -ForegroundColor Yellow
    }
    
    # Remove main copilot instructions
    $copilotInstructionsPath = "$userDir\copilot-instructions.md"
    if (Test-Path $copilotInstructionsPath) {
        Write-Host "Removing main copilot instructions..." -ForegroundColor Red
        try {
            Remove-Item -Path $copilotInstructionsPath -Force -ErrorAction Stop
            Write-Host "Copilot instructions removed" -ForegroundColor Green
        } catch {
            Write-Host "ERROR: Failed to remove copilot instructions: $($_.Exception.Message)" -ForegroundColor Red
        }
    }
    
    # Complex restoration logic with multiple scenarios
    if ($backupPath -and (Test-Path $backupPath)) {
        Write-Host "Restoring original settings.json..." -ForegroundColor Green
        
        # Verify backup integrity using stored metadata from settings.json or backup file
        $storedBackupLength = Get-BackupLengthFromSettings -SettingsPath $settingsPath
        
        # If not found in settings.json, check the most recent backup file for manual merge flag
        if (-not $storedBackupLength) {
            $latestBackup = Get-MostRecentBackup -UserDir $userDir
            if ($latestBackup) {
                $storedBackupLength = Get-BackupLengthFromSettings -SettingsPath $latestBackup
            }
        }
        if ($storedBackupLength) {
            # Handle special condition flags
            if ($storedBackupLength -eq "-1") {
                # Manual merge installation - cannot auto-restore
                Write-Host "Manual merge detected - cannot auto-restore" -ForegroundColor Yellow
                Write-Host ""
                Write-Host "MANUAL CLEANUP REQUIRED:" -ForegroundColor Red
                Write-Host "   This installation was done manually without automatic backup." -ForegroundColor White
                Write-Host "   You must manually remove AzureRM Copilot settings from your settings.json:" -ForegroundColor White
                Write-Host ""
                Write-Host "   Remove these entries from settings.json:" -ForegroundColor White
                Write-Host "   - `"github.copilot.chat.commitMessageGeneration.instructions`"" -ForegroundColor Gray
                Write-Host "   - `"github.copilot.chat.summarizeAgentConversationHistory.enabled`"" -ForegroundColor Gray
                Write-Host "   - `"github.copilot.chat.reviewSelection.enabled`"" -ForegroundColor Gray
                Write-Host "   - `"github.copilot.chat.reviewSelection.instructions`"" -ForegroundColor Gray
                Write-Host "   - `"github.copilot.advanced`" (length and temperature settings)" -ForegroundColor Gray
                Write-Host "   - `"github.copilot.enable`" (if added by AzureRM setup)" -ForegroundColor Gray
                Write-Host "   - `"files.associations`" entries for *.instructions.md, *.prompt.md, *.azurerm.md" -ForegroundColor Gray
                Write-Host "   - `"// AZURERM_BACKUP_LENGTH`" (if present)" -ForegroundColor Gray
                Write-Host ""
                
                # Remove the fake backup file
                if (Test-Path $backupPath) {
                    Remove-Item -Path $backupPath -Force -ErrorAction SilentlyContinue
                    Write-Host "Cleanup marker removed" -ForegroundColor Green
                }
                return
            }
            elseif ($storedBackupLength -eq "0") {
                # Originally empty file - remove settings.json entirely
                Write-Host "Original file was empty - removing settings.json entirely" -ForegroundColor Green
                if (Test-Path $settingsPath) {
                    Remove-Item -Path $settingsPath -Force -ErrorAction Stop
                }
                if (Test-Path $backupPath) {
                    Remove-Item -Path $backupPath -Force -ErrorAction Stop
                }
                Write-Host "Settings.json removed (original was empty)" -ForegroundColor Green
                return
            }
            else {
                # Normal case - verify backup integrity
                $backupLength = (Get-Item $backupPath).Length
                if ($backupLength -ne [int]$storedBackupLength) {
                    Write-Host "WARNING: Backup size mismatch. Expected: $storedBackupLength, Found: $backupLength" -ForegroundColor Yellow
                    Write-Host "Backup may be corrupted. Proceeding with caution..." -ForegroundColor Yellow
                }
            }
        }
        
        try {
            # Test backup validity by attempting to parse as JSON
            $backupContent = Get-Content $backupPath -Raw -ErrorAction Stop
            $null = $backupContent | ConvertFrom-Json -ErrorAction Stop
            
            # Backup is valid JSON, proceed with restoration
            Copy-Item -Path $backupPath -Destination $settingsPath -Force -ErrorAction Stop
            Remove-Item -Path $backupPath -Force -ErrorAction Stop
            
            Write-Host "Original settings restored successfully!" -ForegroundColor Green
            
        } catch {
            Write-Host "ERROR: Backup restoration failed: $($_.Exception.Message)" -ForegroundColor Red
            Write-Host ""
            Write-Host "MANUAL INTERVENTION REQUIRED:" -ForegroundColor Yellow
            Write-Host "  1. Backup location: $backupPath" -ForegroundColor White
            Write-Host "  2. Current settings: $settingsPath" -ForegroundColor White
            Write-Host "  3. Please manually review and restore settings if needed" -ForegroundColor White
            Write-Host ""
            
            # Provide manual merge instructions
            Write-Host "Manual merge instructions:" -ForegroundColor Cyan
            Write-Host "  - Open both files in VS Code for comparison" -ForegroundColor White
            Write-Host "  - Copy important settings from backup to current file" -ForegroundColor White
            Write-Host "  - Remove AzureRM-specific GitHub Copilot settings" -ForegroundColor White
            Write-Host "  - Validate JSON syntax before saving" -ForegroundColor White
        }
        
    } elseif (Test-Path $settingsPath) {
        # No backup found, but settings.json exists - check if AzureRM settings are present
        $settingsContent = Get-Content $settingsPath -Raw -ErrorAction SilentlyContinue
        $hasAzureRMSettings = $false
        
        if ($settingsContent) {
            try {
                $settings = $settingsContent | ConvertFrom-Json -ErrorAction Stop
                $hasAzureRMSettings = (
                    ($settings.PSObject.Properties.Name -contains "github.copilot.chat.commitMessageGeneration.instructions") -or
                    ($settings.PSObject.Properties.Name -contains "github.copilot.chat.summarizeAgentConversationHistory.enabled") -or
                    ($settings.PSObject.Properties.Name -contains "github.copilot.chat.reviewSelection.enabled") -or
                    ($settings.PSObject.Properties.Name -contains "github.copilot.chat.reviewSelection.instructions") -or
                    ($settings.PSObject.Properties.Name -contains "github.copilot.advanced") -or
                    ($settingsContent -match '"//\s*AZURERM_BACKUP_LENGTH"\s*:') -or
                    ($settings."files.associations" -and (
                        $settings."files.associations"."*.instructions.md" -or
                        $settings."files.associations"."*.prompt.md" -or
                        $settings."files.associations"."*.azurerm.md"
                    ))
                )
            } catch {
                Write-Host "WARNING: Unable to parse settings.json - may contain syntax errors" -ForegroundColor Yellow
                # If we can't parse JSON, check raw content for AzureRM traces
                $hasAzureRMSettings = ($settingsContent -match "AZURERM|terraform-azurerm|copilot.*instructions")
            }
        }
        
        if ($hasAzureRMSettings) {
            # Check backup length to determine cleanup strategy
            if ($settingsContent -match '"//\s*AZURERM_BACKUP_LENGTH"\s*:\s*"(-?\d+)"' -or 
                $settingsContent -match '"AZURERM_BACKUP_LENGTH"\s*:\s*(-?\d+)') {
                $backupLength = [int]$matches[1]
                
                switch ($backupLength) {
                    0 {
                        # File was created entirely by our install - safe to remove completely
                        Write-Host "Settings file was created entirely by AzureRM install - removing completely" -ForegroundColor Yellow
                        Remove-Item -Path $settingsPath -Force -ErrorAction Stop
                        Write-Host "Settings.json removed (was created by AzureRM install, VS Code will recreate when needed)" -ForegroundColor Green
                        return
                    }
                    -1 {
                        # Manual merge was required - too risky for automatic cleanup
                        Write-Host "ERROR: Settings were manually merged during install" -ForegroundColor Red
                        Write-Host "Automatic cleanup cannot safely remove AzureRM settings from manually merged file." -ForegroundColor Yellow
                        Write-Host "Please manually remove AzureRM settings from:" -ForegroundColor Yellow
                        Write-Host "  $settingsPath" -ForegroundColor White
                        return
                    }
                    default {
                        # Positive number - original file was backed up, proceed with regex cleanup
                        Write-Host "Cleaning AzureRM settings from settings.json (backup available)..." -ForegroundColor Yellow
                    }
                }
            } else {
                # No backup length found - assume standard cleanup
                Write-Host "Cleaning AzureRM settings from settings.json..." -ForegroundColor Yellow
            }
            
            try {
                # Clean up the settings - comprehensive approach
                $cleanedContent = $settingsContent
                
                # Remove backup metadata (handle both comment style and property style)
                $cleanedContent = $cleanedContent -replace '"//\s*AZURERM_BACKUP_LENGTH"\s*:\s*"[^"]*"\s*,?\s*', ''
                $cleanedContent = $cleanedContent -replace '"AZURERM_BACKUP_LENGTH"\s*:\s*[^,}]*\s*,?\s*', ''
                
                # Remove GitHub Copilot settings (handle both string and array formats)
                $cleanedContent = $cleanedContent -replace '"github\.copilot\.enable"\s*:\s*\{[^{}]*\}\s*,?\s*', ''
                $cleanedContent = $cleanedContent -replace '"github\.copilot\.chat\.commitMessageGeneration\.instructions"\s*:\s*"[^"]*"\s*,?\s*', ''
                # Handle complex array format for commit message instructions (multi-line with nested objects)
                $cleanedContent = $cleanedContent -replace '(?s)"github\.copilot\.chat\.commitMessageGeneration\.instructions"\s*:\s*\[.*?\]\s*,?\s*', ''
                $cleanedContent = $cleanedContent -replace '"github\.copilot\.chat\.summarizeAgentConversationHistory\.enabled"\s*:\s*(true|false)\s*,?\s*', ''
                $cleanedContent = $cleanedContent -replace '"github\.copilot\.chat\.reviewSelection\.enabled"\s*:\s*(true|false)\s*,?\s*', ''
                $cleanedContent = $cleanedContent -replace '"github\.copilot\.chat\.reviewSelection\.instructions"\s*:\s*"[^"]*"\s*,?\s*', ''
                # Handle complex array format for review selection instructions (multi-line with nested objects)
                $cleanedContent = $cleanedContent -replace '(?s)"github\.copilot\.chat\.reviewSelection\.instructions"\s*:\s*\[.*?\]\s*,?\s*', ''
                $cleanedContent = $cleanedContent -replace '"github\.copilot\.advanced"\s*:\s*\{[^}]*\}\s*,?\s*', ''
                
                # Clean up file associations
                $cleanedContent = $cleanedContent -replace '"\*\.instructions\.md"\s*:\s*"[^"]*"\s*,?\s*', ''
                $cleanedContent = $cleanedContent -replace '"\*\.prompt\.md"\s*:\s*"[^"]*"\s*,?\s*', ''
                $cleanedContent = $cleanedContent -replace '"\*\.azurerm\.md"\s*:\s*"[^"]*"\s*,?\s*', ''
                
                # Clean up trailing commas and empty lines
                $cleanedContent = $cleanedContent -replace ',(\s*[}\]])', '$1'
                $cleanedContent = $cleanedContent -replace '\n\s*\n', "`n"
                
                # Test if cleaned content is valid JSON
                try {
                    $parsedJson = $cleanedContent | ConvertFrom-Json -ErrorAction Stop
                } catch {
                    # If cleaning resulted in invalid JSON, the file was likely corrupted
                    # Fall back to removing the file entirely since it contained only AzureRM content
                    Write-Host "Settings file appears corrupted after cleaning - removing entirely" -ForegroundColor Yellow
                    Remove-Item -Path $settingsPath -Force -ErrorAction Stop
                    Write-Host "Corrupted settings.json removed (VS Code will recreate when needed)" -ForegroundColor Green
                    return
                }
                
                # Check if the result is an empty object and handle appropriately
                $propCount = $parsedJson.PSObject.Properties.Count
                if ($propCount -eq $null -or $propCount -eq 0) {
                    Write-Host "Settings file contained only AzureRM metadata - removing empty file" -ForegroundColor Yellow
                    Remove-Item -Path $settingsPath -Force -ErrorAction Stop
                    Write-Host "Empty settings.json removed (VS Code will recreate when needed)" -ForegroundColor Green
                } else {
                    # Write cleaned content back for non-empty settings
                    Set-Content -Path $settingsPath -Value $cleanedContent -ErrorAction Stop
                    Write-Host "AzureRM settings cleaned from settings.json successfully!" -ForegroundColor Green
                }
                
            } catch {
                Write-Host "ERROR: Failed to clean settings.json automatically: $($_.Exception.Message)" -ForegroundColor Red
                Write-Host "Manual cleanup required for settings.json" -ForegroundColor Yellow
            }
        } else {
            Write-Host "No AzureRM settings found in settings.json" -ForegroundColor Green
            
            # Check if settings.json is empty and remove it
            try {
                $settingsJson = $settingsContent | ConvertFrom-Json -ErrorAction Stop
                $propCount = $settingsJson.PSObject.Properties.Count
                if ($settingsJson -and ($propCount -eq $null -or $propCount -eq 0)) {
                    Write-Host "Settings.json is empty - removing unnecessary file" -ForegroundColor Yellow
                    Remove-Item -Path $settingsPath -Force -ErrorAction Stop
                    Write-Host "Empty settings.json removed (VS Code will recreate when needed)" -ForegroundColor Green
                }
            } catch {
                # If we can't parse as JSON, check if it's effectively empty (just braces and whitespace)
                $trimmedContent = $settingsContent.Trim() -replace '\s+', ''
                if ($trimmedContent -eq '{}') {
                    Write-Host "Settings.json contains only empty braces - removing unnecessary file" -ForegroundColor Yellow
                    Remove-Item -Path $settingsPath -Force -ErrorAction Stop
                    Write-Host "Empty settings.json removed (VS Code will recreate when needed)" -ForegroundColor Green
                } else {
                    Write-Host "Settings.json exists but not removing (contains non-standard content)" -ForegroundColor Yellow
                }
            }
        }
        
    } else {
        Write-Host "No settings.json found - no cleanup needed" -ForegroundColor Green
    }
    
    Write-Host ""
    Write-Host "Cleanup completed!" -ForegroundColor Green
    Write-Host "Restart VS Code to ensure all changes take effect." -ForegroundColor Cyan
    exit 0
}

# Auto-discover repository root or use provided path
if ($RepositoryPath) {
    $repoRoot = $RepositoryPath
    Write-Host "Using provided repository path: $repoRoot" -ForegroundColor Cyan
} else {
    Write-Host "Auto-discovering repository location..." -ForegroundColor Blue
    $repoRoot = Find-RepositoryRoot
}

if (-not $repoRoot) {
    Write-Host "Repository auto-discovery failed." -ForegroundColor Yellow
    Write-Host "Please provide the path to your terraform-provider-azurerm repository:" -ForegroundColor Yellow
    $repoRoot = Read-Host "Repository path"
    
    if (-not $repoRoot) {
        Write-Host "ERROR: No repository path provided" -ForegroundColor Red
        exit 1
    }
}

# Validate repository
if (-not (Test-Path $repoRoot)) {
    Write-Host "ERROR: Path does not exist: $repoRoot" -ForegroundColor Red
    exit 1
}

$goModPath = Join-Path $repoRoot "go.mod"
if (-not (Test-Path $goModPath)) {
    Write-Host "ERROR: Invalid repository: go.mod not found in $repoRoot" -ForegroundColor Red
    Write-Host "Please ensure you're pointing to the root of the terraform-provider-azurerm repository." -ForegroundColor Yellow
    exit 1
}

$goModContent = Get-Content $goModPath -Raw -ErrorAction SilentlyContinue
if ($goModContent -notmatch "module github\.com/hashicorp/terraform-provider-azurerm") {
    Write-Host "ERROR: Invalid repository: not terraform-provider-azurerm" -ForegroundColor Red
    Write-Host "Found module: $(($goModContent -split "`n" | Select-Object -First 1) -replace 'module ', '')" -ForegroundColor Yellow
    exit 1
}

Write-Host "Repository validated: $repoRoot" -ForegroundColor Green

# Define paths
$userDir = "$env:APPDATA\Code\User"
$instructionsDir = "$userDir\instructions\terraform-azurerm"
$promptsDir = "$userDir\prompts"
$sourceDir = "$repoRoot\.github"
$settingsPath = "$userDir\settings.json"

# Create timestamped backup path
$timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
$backupPath = "$userDir\settings.json.backup.$timestamp"

# Previous installation detection to prevent backup corruption
$previousInstall = Test-PreviousInstallation -UserDir $userDir

if ($previousInstall.HasPrevious) {
    Write-Host ""
    Write-Host "Previous installation detected!" -ForegroundColor Yellow
    Write-Host "Found existing files:" -ForegroundColor Yellow
    foreach ($file in $previousInstall.FoundFiles) {
        Write-Host "  - $file" -ForegroundColor Gray
    }
    Write-Host ""
    Write-Host "This will overwrite existing AzureRM AI configuration." -ForegroundColor Yellow
    do {
        $response = Read-Host "Continue with installation? (Y)es/(N)o"
        $response = $response.ToUpper()
    } while ($response -ne "Y" -and $response -ne "N")
    
    if ($response -eq "N") {
        Write-Host "Installation cancelled by user." -ForegroundColor Yellow
        exit 0
    }
}

# Create all necessary directories
Write-Host "Creating directories..." -ForegroundColor Blue
try {
    New-Item -ItemType Directory -Force -Path $instructionsDir -ErrorAction Stop | Out-Null
    New-Item -ItemType Directory -Force -Path $promptsDir -ErrorAction Stop | Out-Null
    Write-Host "Directories created" -ForegroundColor Green
} catch {
    Write-Host "ERROR: Failed to create directories: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Copy all instruction files
Write-Host "Copying instruction files..." -ForegroundColor Blue
$instructionsSrcDir = "$sourceDir\instructions"
if (Test-Path $instructionsSrcDir) {
    try {
        $instructionFiles = Get-ChildItem -Path $instructionsSrcDir -Filter "*.md" -ErrorAction Stop
        
        if ($instructionFiles.Count -eq 0) {
            Write-Host "WARNING: No instruction files (*.md) found in $instructionsSrcDir" -ForegroundColor Yellow
        } else {
            Copy-Item -Path "$instructionsSrcDir\*.md" -Destination $instructionsDir -Force -ErrorAction Stop
            Write-Host "Copied $($instructionFiles.Count) instruction files" -ForegroundColor Green
        }
    } catch {
        Write-Host "ERROR: Failed to copy instruction files: $($_.Exception.Message)" -ForegroundColor Red
        exit 1
    }
} else {
    Write-Host "WARNING: Instructions directory not found: $instructionsSrcDir" -ForegroundColor Yellow
}

# Copy AI prompt files
Write-Host "Copying AI prompt files..." -ForegroundColor Blue
$promptsSrcDir = "$sourceDir\prompts"
if (Test-Path $promptsSrcDir) {
    try {
        $promptFiles = Get-ChildItem -Path $promptsSrcDir -Filter "*.md" -ErrorAction Stop
        
        if ($promptFiles.Count -eq 0) {
            Write-Host "WARNING: No prompt files (*.md) found in $promptsSrcDir" -ForegroundColor Yellow
        } else {
            Copy-Item -Path "$promptsSrcDir\*.md" -Destination $promptsDir -Force -ErrorAction Stop
            Write-Host "Copied $($promptFiles.Count) prompt files" -ForegroundColor Green
        }
    } catch {
        Write-Host "ERROR: Failed to copy prompt files: $($_.Exception.Message)" -ForegroundColor Red
        exit 1
    }
} else {
    Write-Host "WARNING: Prompts directory not found: $promptsSrcDir" -ForegroundColor Yellow
}

# Copy main copilot instructions
Write-Host "Copying main copilot instructions..." -ForegroundColor Blue
$copilotSrcPath = "$sourceDir\copilot-instructions.md"
if (Test-Path $copilotSrcPath) {
    try {
        Copy-Item -Path $copilotSrcPath -Destination $userDir -Force -ErrorAction Stop
        
        # Verify the copied file
        if (Test-FileIntegrity -FilePath "$userDir\copilot-instructions.md" -ExpectedPattern "Terraform.*Provider.*Azure") {
            Write-Host "Copilot instructions copied and verified" -ForegroundColor Green
        } else {
            Write-Host "WARNING: Copilot instructions copied but verification failed" -ForegroundColor Yellow
        }
    } catch {
        Write-Host "ERROR: Failed to copy copilot instructions: $($_.Exception.Message)" -ForegroundColor Red
        exit 1
    }
} else {
    Write-Host "WARNING: copilot-instructions.md not found: $copilotSrcPath" -ForegroundColor Yellow
}

# Manual backup confirmation for extra safety
Write-Host ""
Write-Host "MANUAL BACKUP OPPORTUNITY" -ForegroundColor Cyan
Write-Host "========================="
Write-Host "Your VS Code settings file location:"
Write-Host "   $settingsPath" -ForegroundColor Yellow
Write-Host ""
Write-Host "RECOMMENDED: Create your own backup copy before proceeding!" -ForegroundColor Green
Write-Host "   Example: Copy to Documents, Desktop, or another safe location"
Write-Host ""
Write-Host "The script will also create automatic backups, but this gives you full control."
Write-Host ""
Write-Host "Actions you can take now:"
Write-Host "   - Open File Explorer and navigate to the path above"
Write-Host "   - Copy settings.json to your preferred backup location"
Write-Host "   - Choose 'C' to continue or 'X' to exit without changes"

if (Test-Path $settingsPath) {
    $fileSize = (Get-Item $settingsPath).Length
    Write-Host ""
    Write-Host "Current settings.json found - size: $fileSize bytes" -ForegroundColor Yellow
} else {
    Write-Host ""
    Write-Host "No existing settings.json found - new file will be created" -ForegroundColor Green
}

Write-Host ""
Write-Host "What would you like to do?"
Write-Host "   [C] Continue with automatic setup"
Write-Host "   [X] Exit script (no changes made)"
Write-Host ""

do {
    $choice = Read-Host "Choice (C/X)"
    $choice = $choice.ToUpper()
    
    switch ($choice) {
        "C" {
            Write-Host "Continuing with setup..." -ForegroundColor Green
            break
        }
        "X" {
            Write-Host "Setup cancelled by user. No changes were made." -ForegroundColor Yellow
            exit 0
        }
        default {
            Write-Host "Invalid choice. Please press C to continue or X to exit." -ForegroundColor Red
        }
    }
} while ($choice -ne "C")

# Smart merge VS Code settings.json with enhanced backup and recovery
Write-Host "Configuring VS Code settings (preserving existing settings)..." -ForegroundColor Blue

$existingSettings = @{}
$settingsBackupCreated = $false

# Determine original settings length - will be set based on processing outcome
if (-not (Test-Path $settingsPath)) {
    # Special flag: No original file existed
    $originalSettingsLength = "0"
} else {
    # Normal case: Calculate actual file size (will be overridden to "-1" if parsing fails)
    $originalSettingsLength = (Get-Item $settingsPath).Length.ToString()
}

# Handle existing settings.json with comprehensive error handling
if (Test-Path $settingsPath) {
    try {
        $existingContent = Get-Content $settingsPath -Raw -ErrorAction Stop
        $existingSettings = $existingContent | ConvertFrom-Json -ErrorAction Stop
        
        # Create safe backup with integrity verification (only for non-manual merge)
        # Create backup if it doesn't exist
        if (-not (Test-Path $backupPath)) {
            $settingsBackupCreated = Create-SafeBackup -SourcePath $settingsPath -BackupPath $backupPath
            
            if (-not $settingsBackupCreated) {
                Write-Host "ERROR: Failed to create settings backup. Installation aborted for safety." -ForegroundColor Red
                exit 1
            }
        } else {
            Write-Host "Note: Backup already exists: $(Split-Path $backupPath -Leaf)" -ForegroundColor Yellow
        }
        
    } catch {
        Write-Host "WARNING: Could not parse existing settings.json as JSON" -ForegroundColor Yellow
        Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
        Write-Host ""
        
        # Set error recovery flag for cleanup detection
        $originalSettingsLength = "-1"
        
        # Create fake backup with only the manual merge flag for cleanup detection
        try {
            $fakeBackup = @{
                "// AZURERM_BACKUP_LENGTH" = "-1"
            }
            $fakeBackup | ConvertTo-Json -Depth 32 | Set-Content -Path $backupPath -Force -ErrorAction Stop
            Write-Host "Created recovery marker for cleanup" -ForegroundColor Green
        } catch {
            Write-Host "ERROR: Could not create recovery marker" -ForegroundColor Red
            exit 1
        }
        
        Write-Host ""
        Write-Host "MANUAL INTERVENTION REQUIRED:" -ForegroundColor Yellow
        Write-Host "  1. Fix JSON syntax errors in: $settingsPath" -ForegroundColor White
        Write-Host "  2. Manually add the following to your corrected settings.json:" -ForegroundColor White
        Write-Host ""
        Write-Host "     {" -ForegroundColor Gray
        Write-Host "       `"github.copilot.chat.commitMessageGeneration.instructions`": [`"file:$userDir\\copilot-instructions.md`"]," -ForegroundColor Gray
        Write-Host "       `"github.copilot.chat.summarizeAgentConversationHistory.enabled`": false," -ForegroundColor Gray
        Write-Host "       `"github.copilot.chat.reviewSelection.enabled`": true," -ForegroundColor Gray
        Write-Host "       `"github.copilot.chat.reviewSelection.instructions`": [" -ForegroundColor Gray
        Write-Host "         `"file:$userDir\\copilot-instructions.md`"," -ForegroundColor Gray
        Write-Host "         `"file:$userDir\\instructions\\terraform-azurerm\\implementation-guide.instructions.md`"," -ForegroundColor Gray
        Write-Host "         `"file:$userDir\\instructions\\terraform-azurerm\\azure-patterns.instructions.md`"," -ForegroundColor Gray
        Write-Host "         `"file:$userDir\\instructions\\terraform-azurerm\\testing-guidelines.instructions.md`"," -ForegroundColor Gray
        Write-Host "         `"file:$userDir\\instructions\\terraform-azurerm\\documentation-guidelines.instructions.md`"," -ForegroundColor Gray
        Write-Host "         `"file:$userDir\\instructions\\terraform-azurerm\\provider-guidelines.instructions.md`"," -ForegroundColor Gray
        Write-Host "         `"file:$userDir\\instructions\\terraform-azurerm\\error-patterns.instructions.md`"," -ForegroundColor Gray
        Write-Host "         `"file:$userDir\\instructions\\terraform-azurerm\\code-clarity-enforcement.instructions.md`"," -ForegroundColor Gray
        Write-Host "         `"file:$userDir\\instructions\\terraform-azurerm\\schema-patterns.instructions.md`"," -ForegroundColor Gray
        Write-Host "         `"file:$userDir\\instructions\\terraform-azurerm\\migration-guide.instructions.md`"," -ForegroundColor Gray
        Write-Host "         `"file:$userDir\\instructions\\terraform-azurerm\\performance-optimization.instructions.md`"," -ForegroundColor Gray
        Write-Host "         `"file:$userDir\\instructions\\terraform-azurerm\\security-compliance.instructions.md`"," -ForegroundColor Gray
        Write-Host "         `"file:$userDir\\instructions\\terraform-azurerm\\troubleshooting-decision-trees.instructions.md`"," -ForegroundColor Gray
        Write-Host "         `"file:$userDir\\instructions\\terraform-azurerm\\api-evolution-patterns.instructions.md`"" -ForegroundColor Gray
        Write-Host "       ]," -ForegroundColor Gray
        Write-Host "       `"github.copilot.advanced`": {" -ForegroundColor Gray
        Write-Host "         `"length`": 3000," -ForegroundColor Gray
        Write-Host "         `"temperature`": 0.1" -ForegroundColor Gray
        Write-Host "       }," -ForegroundColor Gray
        Write-Host "       `"github.copilot.enable`": {" -ForegroundColor Gray
        Write-Host "         `"*`": true," -ForegroundColor Gray
        Write-Host "         `"terminal`": true" -ForegroundColor Gray
        Write-Host "       }," -ForegroundColor Gray
        Write-Host "       `"files.associations`": {" -ForegroundColor Gray
        Write-Host "         `"*.instructions.md`": `"markdown`"," -ForegroundColor Gray
        Write-Host "         `"*.prompt.md`": `"markdown`"," -ForegroundColor Gray
        Write-Host "         `"*.azurerm.md`": `"markdown`"" -ForegroundColor Gray
        Write-Host "       }" -ForegroundColor Gray
        Write-Host "       // ...your existing settings..." -ForegroundColor Gray
        Write-Host "     }" -ForegroundColor Gray
        Write-Host "  3. Backup of current file saved to: $backupPath" -ForegroundColor White
        Write-Host ""
        Write-Host "Common JSON fixes:" -ForegroundColor Cyan
        Write-Host "  - Remove trailing commas before closing braces }" -ForegroundColor White
        Write-Host "  - Ensure all strings are properly quoted" -ForegroundColor White
        Write-Host "  - Check for unescaped quotes within strings" -ForegroundColor White
        Write-Host "  - Validate brackets and braces are balanced" -ForegroundColor White
        Write-Host ""
        Write-Host "IMPORTANT: Make a backup copy of your settings.json before editing!" -ForegroundColor Red
        Write-Host "  Example: Copy-Item '$settingsPath' '$settingsPath.manual-backup'" -ForegroundColor Gray
        exit 1
    }
} else {
    Write-Host "No existing settings.json found - creating new one..." -ForegroundColor Green
    $existingSettings = @{}
}

# Define comprehensive AzureRM provider-specific settings
$azureRMSettings = @{
    "github.copilot.chat.commitMessageGeneration.instructions" = @(
        @{
            "text" = "Provide a concise and clear commit message that summarizes the changes made in the code. For complex changes, include the following details: 1) Specify if the change introduces a breaking change and describe its impact. 2) Highlight any new resources or features added. 3) Mention updates to Azure services or APIs. Aim to keep the message under 72 characters per line for readability."
        }
    )
    "github.copilot.chat.summarizeAgentConversationHistory.enabled" = $false
    "github.copilot.chat.reviewSelection.enabled" = $true
    "github.copilot.chat.reviewSelection.instructions" = @(
        @{"file" = "copilot-instructions.md"}
        @{"file" = "instructions/terraform-azurerm/implementation-guide.instructions.md"}
        @{"file" = "instructions/terraform-azurerm/azure-patterns.instructions.md"}
        @{"file" = "instructions/terraform-azurerm/testing-guidelines.instructions.md"}
        @{"file" = "instructions/terraform-azurerm/documentation-guidelines.instructions.md"}
        @{"file" = "instructions/terraform-azurerm/provider-guidelines.instructions.md"}
        @{"file" = "instructions/terraform-azurerm/error-patterns.instructions.md"}
        @{"file" = "instructions/terraform-azurerm/code-clarity-enforcement.instructions.md"}
        @{"file" = "instructions/terraform-azurerm/schema-patterns.instructions.md"}
        @{"file" = "instructions/terraform-azurerm/migration-guide.instructions.md"}
        @{"file" = "instructions/terraform-azurerm/performance-optimization.instructions.md"}
        @{"file" = "instructions/terraform-azurerm/security-compliance.instructions.md"}
        @{"file" = "instructions/terraform-azurerm/troubleshooting-decision-trees.instructions.md"}
        @{"file" = "instructions/terraform-azurerm/api-evolution-patterns.instructions.md"}
    )
    "github.copilot.advanced" = @{
        "length" = 3000
        "temperature" = 0.1
    }
    "github.copilot.enable" = @{
        "*" = $true
        "terminal" = $true
    }
}

# Enhanced file associations merging
$newFileAssociations = @{
    "*.instructions.md" = "markdown"
    "*.prompt.md" = "markdown"
    "*.azurerm.md" = "markdown"
}

if ($existingSettings."files.associations") {
    # Preserve all existing file associations
    $existingSettings."files.associations" | Get-Member -MemberType NoteProperty | ForEach-Object {
        $key = $_.Name
        $value = $existingSettings."files.associations".$key
        if (-not $newFileAssociations.ContainsKey($key)) {
            $newFileAssociations[$key] = $value
        }
    }
}

$azureRMSettings."files.associations" = $newFileAssociations

# Add backup length for future verification and restoration
$azureRMSettings."// AZURERM_BACKUP_LENGTH" = $originalSettingsLength

# Intelligent settings merging (preserves user settings, adds AzureRM enhancements)
try {
    foreach ($key in $azureRMSettings.Keys) {
        if ($existingSettings.PSObject.Properties.Name -contains $key) {
            # Key exists - merge intelligently
            switch ($key) {
                "files.associations" {
                    # Merge file associations
                    $existingSettings.$key = $azureRMSettings.$key
                }
                "github.copilot.enable" {
                    # Merge enable settings (preserve user's disable choices)
                    if ($existingSettings.$key -is [hashtable] -or $existingSettings.$key.GetType().Name -eq "PSCustomObject") {
                        foreach ($enableKey in $azureRMSettings.$key.Keys) {
                            if (-not ($existingSettings.$key.PSObject.Properties.Name -contains $enableKey)) {
                                $existingSettings.$key | Add-Member -MemberType NoteProperty -Name $enableKey -Value $azureRMSettings.$key.$enableKey -Force
                            }
                        }
                    } else {
                        $existingSettings.$key = $azureRMSettings.$key
                    }
                }
                default {
                    # For other settings, AzureRM takes precedence
                    $existingSettings.$key = $azureRMSettings.$key
                }
            }
        } else {
            # Key doesn't exist - add it
            $existingSettings | Add-Member -MemberType NoteProperty -Name $key -Value $azureRMSettings.$key -Force
        }
    }
    
    # Convert merged settings to JSON with proper formatting
    $mergedJson = $existingSettings | ConvertTo-Json -Depth 10 -Compress:$false
    
    # Write to file with error handling
    $mergedJson | Out-File -FilePath $settingsPath -Encoding UTF8 -Force -ErrorAction Stop
    
    # Verify the written file
    if (Test-FileIntegrity -FilePath $settingsPath -ExpectedPattern "github\.copilot") {
        Write-Host "Settings merged and verified successfully!" -ForegroundColor Green
    } else {
        Write-Host "WARNING: Settings written but verification failed" -ForegroundColor Yellow
    }
    
} catch {
    Write-Host "ERROR: Failed to merge settings: $($_.Exception.Message)" -ForegroundColor Red
    
    # Attempt to restore backup if we created one
    if ($settingsBackupCreated -and (Test-Path $backupPath)) {
        try {
            Copy-Item -Path $backupPath -Destination $settingsPath -Force -ErrorAction Stop
            Write-Host "Settings restored from backup" -ForegroundColor Green
        } catch {
            Write-Host "ERROR: Failed to restore backup!" -ForegroundColor Red
        }
    }
    
    exit 1
}

Write-Host ""
Write-Host "Installation Complete!" -ForegroundColor Green
Write-Host "=============================" -ForegroundColor Green
Write-Host ""
Write-Host "Files installed to:" -ForegroundColor Cyan
Write-Host "  Instructions: $instructionsDir" -ForegroundColor Gray
Write-Host "  AI Prompts: $promptsDir" -ForegroundColor Gray
Write-Host "  VS Code Settings: $settingsPath" -ForegroundColor Gray
if ($settingsBackupCreated) {
    Write-Host "  Settings Backup: $backupPath" -ForegroundColor Gray
}
Write-Host ""
Write-Host "Next Steps:" -ForegroundColor Yellow
Write-Host "  1. Restart VS Code to load new settings" -ForegroundColor White
Write-Host "  2. Open your terraform-provider-azurerm repository" -ForegroundColor White
Write-Host "  3. Start using AI-powered development!" -ForegroundColor White
Write-Host ""
Write-Host "Available Features:" -ForegroundColor Green
Write-Host "  - Context-aware code generation with @workspace" -ForegroundColor White
Write-Host "  - Azure-specific implementation patterns" -ForegroundColor White
Write-Host "  - Automated testing and documentation" -ForegroundColor White
Write-Host "  - Code review and quality enforcement" -ForegroundColor White
Write-Host "  - Error handling and debugging assistance" -ForegroundColor White
Write-Host ""
Write-Host "For help: .\install-copilot-setup.ps1 -Help" -ForegroundColor Cyan
Write-Host "To remove: .\install-copilot-setup.ps1 -Clean" -ForegroundColor Cyan
