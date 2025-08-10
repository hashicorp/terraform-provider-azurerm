# PowerShell Module Quick Reference

## Quick Start Commands

```powershell
# 1. Import the module
Import-Module .\modules\TerraformAzureRMSetup.psd1

# 2. Run complete setup
Install-TerraformAzureRMDevelopmentEnvironment

# 3. Check status
Show-InstallationStatus
```

## Core Functions Quick Reference

| Function | Purpose | Example |
|----------|---------|---------|
| `Install-TerraformAzureRMDevelopmentEnvironment` | Main setup function | `Install-TerraformAzureRMDevelopmentEnvironment -Force` |
| `Show-InstallationStatus` | Check current status | `Show-InstallationStatus -Verbose` |
| `Find-RepositoryRoot` | Find Git repository root | `$repo = Find-RepositoryRoot` |
| `Test-Prerequisites` | Check required software | `if (Test-Prerequisites) { ... }` |
| `Test-CopilotInstallation` | Check Copilot status | `$status = Test-CopilotInstallation` |
| `Get-VSCodeUserSettingsPath` | Get VS Code settings path | `$path = Get-VSCodeUserSettingsPath` |

## Common Parameter Combinations

```powershell
# Force installation without prompts
Install-TerraformAzureRMDevelopmentEnvironment -Force -SkipBackup

# Detailed status with verbose output
Show-InstallationStatus -Verbose

# Test specific repository
Test-CopilotInstallation -RepositoryPath "C:\dev\terraform-provider-azurerm"

# Find repository from subdirectory
Find-RepositoryRoot -StartPath ".\internal\services"
```

## Troubleshooting Commands

```powershell
# Check what's missing
Test-Prerequisites

# Get detailed Copilot status
Test-CopilotInstallation | Format-List

# Find VS Code installation
Get-VSCodeUserSettingsPath

# Test file integrity
Test-FileIntegrity -Path "settings.json"

# Get help for any function
Get-Help <FunctionName> -Full
```

## Error Recovery

```powershell
# If settings are corrupted
$settingsPath = Get-VSCodeUserSettingsPath
if (Test-Path $settingsPath) {
    New-SafeBackup -SourcePath $settingsPath
    # Then re-run installation
}

# If Copilot isn't working
Install-CopilotForRepository -RepositoryPath (Get-Location) -Force

# Check prerequisites again
if (-not (Test-Prerequisites)) {
    Write-Host "Install missing software and try again"
}
```
