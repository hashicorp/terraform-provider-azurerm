# Terraform AzureRM Provider AI Setup Module
# Main module file that loads all submodules

# Get the directory containing this module
$ModuleRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
$ModulesDir = Join-Path $ModuleRoot "modules"

# Import all module files
$moduleFiles = @(
    "core-functions.ps1",
    "backup-management.ps1",
    "installation-detection.ps1", 
    "ai-installation.ps1",
    "cleanup.ps1"
)

foreach ($moduleFile in $moduleFiles) {
    $modulePath = Join-Path $ModulesDir $moduleFile
    if (Test-Path $modulePath) {
        . $modulePath
    } else {
        Write-Warning "Module file not found: $modulePath"
    }
}

# Export module information
Export-ModuleMember -Function @(
    'Find-RepositoryRoot',
    'Test-FileIntegrity', 
    'Test-Prerequisites',
    'Write-StatusMessage',
    'Create-SafeBackup',
    'Get-BackupLengthFromSettings',
    'Get-MostRecentBackup',
    'Restore-FromBackup',
    'Remove-BackupFiles',
    'Test-PreviousInstallation',
    'Get-InstallationPaths',
    'Test-InstallationHealth',
    'Install-CopilotInstructions',
    'Install-VSCodeSettings',
    'Install-AIAgent',
    'Remove-AIAgent',
    'Remove-TerraformSettingsFromVSCode',
    'Test-CleanupSuccess',
    'Show-CleanupSummary'
)
