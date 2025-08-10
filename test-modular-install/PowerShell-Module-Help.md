# PowerShell Module Help Documentation

## Overview

This document provides comprehensive help information for the Terraform AzureRM Provider AI Setup PowerShell modules. These modules automate the setup of GitHub Copilot and VS Code for developing with the Terraform AzureRM provider.

## Module Architecture

### Main Modules

1. **TerraformAzureRMSetup.psm1** - Main orchestration module
2. **CoreFunctions.psm1** - Core utility functions
3. **VSCodeSetup.psm1** - VS Code configuration functions
4. **CopilotInstall.psm1** - GitHub Copilot installation functions

### Module Manifest

- **TerraformAzureRMSetup.psd1** - Module manifest with metadata and dependencies

## Function Reference

### TerraformAzureRMSetup Module

#### Install-TerraformAzureRMDevelopmentEnvironment

**Synopsis:** Main function to set up the complete development environment

**Syntax:**
```powershell
Install-TerraformAzureRMDevelopmentEnvironment [[-RepositoryPath] <string>] [-SkipBackup] [-Force] [<CommonParameters>]
```

**Description:**
Orchestrates the complete setup of GitHub Copilot and VS Code for Terraform AzureRM provider development. This is the main entry point for the setup process.

**Parameters:**

- **RepositoryPath** (Optional)
  - Type: String
  - Default: Current directory
  - Description: Path to the Terraform AzureRM provider repository

- **SkipBackup** (Optional)
  - Type: Switch
  - Description: Skip creating backups of existing configuration files

- **Force** (Optional)
  - Type: Switch
  - Description: Overwrite existing configurations without prompting

**Examples:**
```powershell
# Basic setup in current directory
Install-TerraformAzureRMDevelopmentEnvironment

# Setup with specific repository path
Install-TerraformAzureRMDevelopmentEnvironment -RepositoryPath "C:\dev\terraform-provider-azurerm"

# Force setup without backups
Install-TerraformAzureRMDevelopmentEnvironment -SkipBackup -Force
```

**Outputs:**
- Installation progress messages
- Configuration file paths
- Success/failure status

---

#### Show-InstallationStatus

**Synopsis:** Display the current installation status and configuration

**Syntax:**
```powershell
Show-InstallationStatus [[-RepositoryPath] <string>] [<CommonParameters>]
```

**Description:**
Analyzes and displays the current state of the development environment setup, including GitHub Copilot installation status, VS Code configuration, and repository setup.

**Parameters:**

- **RepositoryPath** (Optional)
  - Type: String
  - Default: Current directory
  - Description: Path to the repository to analyze

**Examples:**
```powershell
# Check status in current directory
Show-InstallationStatus

# Check status for specific repository
Show-InstallationStatus -RepositoryPath "C:\dev\terraform-provider-azurerm"
```

**Outputs:**
- GitHub Copilot installation status
- VS Code configuration status
- Repository setup status
- Recommended next steps

---

### CoreFunctions Module

#### Find-RepositoryRoot

**Synopsis:** Find the root directory of a Git repository

**Syntax:**
```powershell
Find-RepositoryRoot [[-StartPath] <string>] [<CommonParameters>]
```

**Description:**
Searches upward from the specified path to find the root of a Git repository by looking for the .git directory.

**Parameters:**

- **StartPath** (Optional)
  - Type: String
  - Default: Current directory
  - Description: Starting path for the search

**Examples:**
```powershell
# Find repository root from current directory
$repoRoot = Find-RepositoryRoot

# Find repository root from specific path
$repoRoot = Find-RepositoryRoot -StartPath "C:\dev\terraform-provider-azurerm\internal"
```

**Returns:**
- String: Path to repository root
- $null: If no repository found

---

#### Test-Prerequisites

**Synopsis:** Check if required software is installed

**Syntax:**
```powershell
Test-Prerequisites [<CommonParameters>]
```

**Description:**
Validates that all required software components are installed and accessible, including Git, VS Code, and PowerShell.

**Examples:**
```powershell
# Check all prerequisites
$prereqStatus = Test-Prerequisites

if ($prereqStatus) {
    Write-Host "All prerequisites are met"
} else {
    Write-Host "Some prerequisites are missing"
}
```

**Returns:**
- Boolean: $true if all prerequisites are met, $false otherwise

---

#### New-SafeBackup

**Synopsis:** Create a safe backup of a file or directory

**Syntax:**
```powershell
New-SafeBackup [-SourcePath] <string> [[-BackupDirectory] <string>] [<CommonParameters>]
```

**Description:**
Creates a timestamped backup of the specified file or directory in a safe location, preventing data loss during configuration changes.

**Parameters:**

- **SourcePath** (Required)
  - Type: String
  - Description: Path to the file or directory to backup

- **BackupDirectory** (Optional)
  - Type: String
  - Default: "backup" subdirectory
  - Description: Directory where backups should be stored

**Examples:**
```powershell
# Backup a VS Code settings file
$backupPath = New-SafeBackup -SourcePath "$env:APPDATA\Code\User\settings.json"

# Backup to specific directory
$backupPath = New-SafeBackup -SourcePath "settings.json" -BackupDirectory "C:\Backups"
```

**Returns:**
- String: Path to the created backup file/directory

---

#### Test-FileIntegrity

**Synopsis:** Test if a file can be read and is not corrupted

**Syntax:**
```powershell
Test-FileIntegrity [-Path] <string> [<CommonParameters>]
```

**Description:**
Performs basic integrity checks on a file to ensure it can be read and is not obviously corrupted.

**Parameters:**

- **Path** (Required)
  - Type: String
  - Description: Path to the file to test

**Examples:**
```powershell
# Test a configuration file
$isValid = Test-FileIntegrity -Path "settings.json"

if ($isValid) {
    Write-Host "File is readable"
} else {
    Write-Host "File may be corrupted"
}
```

**Returns:**
- Boolean: $true if file is readable, $false otherwise

---

### VSCodeSetup Module

#### Install-VSCodeCopilotConfiguration

**Synopsis:** Install GitHub Copilot configuration for VS Code

**Syntax:**
```powershell
Install-VSCodeCopilotConfiguration [[-RepositoryPath] <string>] [-Force] [<CommonParameters>]
```

**Description:**
Configures VS Code with optimized settings for GitHub Copilot usage in Terraform AzureRM provider development.

**Parameters:**

- **RepositoryPath** (Optional)
  - Type: String
  - Default: Current directory
  - Description: Path to the repository

- **Force** (Optional)
  - Type: Switch
  - Description: Overwrite existing configuration without prompting

**Examples:**
```powershell
# Install basic Copilot configuration
Install-VSCodeCopilotConfiguration

# Force install with specific repository
Install-VSCodeCopilotConfiguration -RepositoryPath "C:\dev\terraform-provider-azurerm" -Force
```

**Outputs:**
- Configuration file paths
- Installation status messages

---

#### Get-VSCodeUserSettingsPath

**Synopsis:** Get the path to VS Code user settings file

**Syntax:**
```powershell
Get-VSCodeUserSettingsPath [<CommonParameters>]
```

**Description:**
Returns the file system path to the VS Code user settings.json file, handling different installation types and operating systems.

**Examples:**
```powershell
# Get settings path
$settingsPath = Get-VSCodeUserSettingsPath
Write-Host "VS Code settings: $settingsPath"
```

**Returns:**
- String: Path to VS Code settings.json file
- $null: If VS Code is not installed or settings path cannot be determined

---

#### Update-VSCodeSettings

**Synopsis:** Update VS Code settings with new configuration

**Syntax:**
```powershell
Update-VSCodeSettings [-Settings] <hashtable> [-SettingsPath] <string> [<CommonParameters>]
```

**Description:**
Safely updates VS Code settings by merging new settings with existing ones, preserving user customizations while adding required Copilot configurations.

**Parameters:**

- **Settings** (Required)
  - Type: Hashtable
  - Description: Settings to add or update

- **SettingsPath** (Required)
  - Type: String
  - Description: Path to the settings.json file

**Examples:**
```powershell
# Update specific settings
$newSettings = @{
    "github.copilot.enable" = $true
    "editor.suggestSelection" = "first"
}
Update-VSCodeSettings -Settings $newSettings -SettingsPath $settingsPath
```

**Outputs:**
- Settings update status
- Backup information if changes are made

---

### CopilotInstall Module

#### Test-CopilotInstallation

**Synopsis:** Check if GitHub Copilot is properly installed and configured

**Syntax:**
```powershell
Test-CopilotInstallation [[-RepositoryPath] <string>] [<CommonParameters>]
```

**Description:**
Performs comprehensive checks to determine if GitHub Copilot is installed, configured, and working properly in the development environment.

**Parameters:**

- **RepositoryPath** (Optional)
  - Type: String
  - Default: Current directory
  - Description: Repository path for context-specific checks

**Examples:**
```powershell
# Test Copilot installation
$copilotStatus = Test-CopilotInstallation

# Test with specific repository context
$copilotStatus = Test-CopilotInstallation -RepositoryPath "C:\dev\terraform-provider-azurerm"
```

**Returns:**
- PSCustomObject with properties:
  - VSCodeInstalled: Boolean
  - CopilotExtensionInstalled: Boolean
  - CopilotConfigured: Boolean
  - AuthenticationStatus: String
  - RecommendedActions: Array of strings

---

#### Install-CopilotForRepository

**Synopsis:** Install and configure GitHub Copilot for a specific repository

**Syntax:**
```powershell
Install-CopilotForRepository [-RepositoryPath] <string> [-Force] [<CommonParameters>]
```

**Description:**
Performs repository-specific GitHub Copilot installation and configuration, including workspace settings and language-specific optimizations.

**Parameters:**

- **RepositoryPath** (Required)
  - Type: String
  - Description: Path to the repository to configure

- **Force** (Optional)
  - Type: Switch
  - Description: Overwrite existing configuration files

**Examples:**
```powershell
# Install Copilot for current repository
Install-CopilotForRepository -RepositoryPath (Get-Location)

# Force install with overwrites
Install-CopilotForRepository -RepositoryPath "C:\dev\terraform-provider-azurerm" -Force
```

**Outputs:**
- Installation progress messages
- Configuration file locations
- Success/failure status

---

## Common Parameters

All functions support PowerShell common parameters:

- **-Verbose**: Show detailed operation information
- **-Debug**: Show debug information
- **-ErrorAction**: How to handle errors (Stop, Continue, SilentlyContinue, etc.)
- **-ErrorVariable**: Variable to store error information
- **-WarningAction**: How to handle warnings
- **-WarningVariable**: Variable to store warning information

## Usage Examples

### Quick Setup

```powershell
# Import the main module
Import-Module .\modules\TerraformAzureRMSetup.psd1

# Run complete setup
Install-TerraformAzureRMDevelopmentEnvironment

# Check status
Show-InstallationStatus
```

### Advanced Configuration

```powershell
# Import module with specific repository
Import-Module .\modules\TerraformAzureRMSetup.psd1

# Find repository root
$repoPath = Find-RepositoryRoot -StartPath "C:\dev\terraform-provider-azurerm\internal"

# Check prerequisites first
if (Test-Prerequisites) {
    Write-Host "Prerequisites OK, proceeding with installation"
    
    # Install with custom settings
    Install-TerraformAzureRMDevelopmentEnvironment -RepositoryPath $repoPath -Force
    
    # Verify installation
    $status = Test-CopilotInstallation -RepositoryPath $repoPath
    if ($status.CopilotConfigured) {
        Write-Host "Setup completed successfully" -ForegroundColor Green
    } else {
        Write-Host "Setup needs attention: $($status.RecommendedActions -join ', ')" -ForegroundColor Yellow
    }
}
```

### Troubleshooting

```powershell
# Check detailed status
Show-InstallationStatus -Verbose

# Test specific components
$prereq = Test-Prerequisites
$copilot = Test-CopilotInstallation
$vscodePath = Get-VSCodeUserSettingsPath

Write-Host "Prerequisites: $prereq"
Write-Host "Copilot Status: $($copilot.CopilotConfigured)"
Write-Host "VS Code Settings: $vscodePath"

# Re-run specific installation steps
if (-not $copilot.CopilotExtensionInstalled) {
    Install-CopilotForRepository -RepositoryPath (Get-Location) -Force
}
```

## Error Handling

All functions include comprehensive error handling and will return meaningful error messages. Common error scenarios:

1. **Missing Prerequisites**: Functions will report which software needs to be installed
2. **Permission Issues**: Clear messages about file/directory access problems
3. **Configuration Conflicts**: Guidance on resolving existing configuration conflicts
4. **Network Issues**: Helpful messages for connectivity problems during extension installation

## Getting Help

Use PowerShell's built-in help system for detailed information:

```powershell
# Get help for specific functions
Get-Help Install-TerraformAzureRMDevelopmentEnvironment -Full
Get-Help Find-RepositoryRoot -Examples
Get-Help Test-CopilotInstallation -Parameter RepositoryPath

# List all available functions
Get-Command -Module TerraformAzureRMSetup

# Get function syntax
Get-Command Install-TerraformAzureRMDevelopmentEnvironment -Syntax
```
