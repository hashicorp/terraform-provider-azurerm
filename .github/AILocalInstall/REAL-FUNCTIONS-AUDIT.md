# Real Functions Audit
# This documents the ACTUAL functions available in the working modular system

## backup-management.psm1
**Exported Functions:**
- New-SafeBackup
- Get-BackupLengthFromSettings  
- Get-MostRecentBackup
- Restore-FromBackup
- Remove-BackupFiles

## config-management.psm1
**Exported Functions:**
- Read-ConfigSection
- Get-InstructionFiles
- Get-PromptFiles
- Get-MainFiles
- Get-ExpectedInstructionFiles
- Get-ExpectedPromptFiles

## installation-verification.psm1
**Exported Functions:**
- Test-HardcodedInstallationIntegrity
- Test-InstallationIntegrity

## core-functions.psm1
**Available Functions:**
- Find-RepositoryRoot
- Test-FileIntegrity
- Test-Prerequisites
- Write-StatusMessage

## ai-installation.psm1
**Available Functions:**
- Install-CopilotInstructions
- Install-VSCodeSettings
- Install-CompleteSetup

## cleanup.psm1
**Available Functions:**
- Remove-Installation
- Test-InstallationStatus
- Remove-TerraformSettingsFromVSCode

## File Operations Used in Real System:
- Copy-Item (standard PowerShell)
- Test-Path (standard PowerShell)
- Join-Path (standard PowerShell)
- Get-Content/Set-Content (standard PowerShell)
- ConvertFrom-Json/ConvertTo-Json (standard PowerShell)

## VS Code Settings Approach:
- Real config uses: terraform_azurerm_provider_mode, terraform_azurerm_ai_enhanced, etc.
- Uses backup length tracking in settings
- PowerShell 5.1 compatible hashtable conversion

## Installation Pattern:
- Files are expected to be IN the repository already
- No file copying between directories
- Validation checks file existence and integrity
- VS Code settings are the only external modification
