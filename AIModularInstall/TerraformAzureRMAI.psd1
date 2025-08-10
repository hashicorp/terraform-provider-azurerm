@{
    # Module manifest for Terraform AzureRM Provider AI Setup
    
    # Script module or binary module file associated with this manifest.
    RootModule = 'TerraformAzureRMAI.psm1'
    
    # Version number of this module.
    ModuleVersion = '2.0.0'
    
    # ID used to uniquely identify this module
    GUID = 'f1e2d3c4-5b6a-7890-abcd-ef1234567890'
    
    # Author of this module
    Author = 'HashiCorp Terraform AzureRM Provider Team'
    
    # Company or vendor of this module
    CompanyName = 'HashiCorp'
    
    # Copyright statement for this module
    Copyright = '(c) HashiCorp. All rights reserved.'
    
    # Description of the functionality provided by this module
    Description = 'Modular installation system for Terraform AzureRM Provider AI development setup'
    
    # Minimum version of the PowerShell engine required by this module
    PowerShellVersion = '5.1'
    
    # Functions to export from this module, for best performance, do not use wildcards and do not delete the entry
    FunctionsToExport = @(
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
    
    # Cmdlets to export from this module
    CmdletsToExport = @()
    
    # Variables to export from this module
    VariablesToExport = @()
    
    # Aliases to export from this module
    AliasesToExport = @()
    
    # Private data to pass to the module specified in RootModule/ModuleToProcess
    PrivateData = @{
        PSData = @{
            # Tags applied to this module for discoverability
            Tags = @('Terraform', 'Azure', 'AzureRM', 'AI', 'Copilot', 'Development')
            
            # A URL to the license for this module.
            LicenseUri = 'https://github.com/hashicorp/terraform-provider-azurerm/blob/main/LICENSE'
            
            # A URL to the main website for this project.
            ProjectUri = 'https://github.com/hashicorp/terraform-provider-azurerm'
            
            # ReleaseNotes of this module
            ReleaseNotes = @'
2.0.0 - Modular Architecture Release
- Replaced monolithic installer with clean modular system
- Enhanced error handling and logging
- Improved backup and restore functionality
- Better separation of concerns
- Maintainable and testable code structure
'@
        }
    }
}
