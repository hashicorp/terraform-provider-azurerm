# TerraformAzureRMSetup.psd1 - Module manifest for Terraform AzureRM Provider AI Setup

@{
    # Script module or binary module file associated with this manifest.
    RootModule = 'TerraformAzureRMSetup.psm1'
    
    # Version number of this module.
    ModuleVersion = '1.0.0'
    
    # Supported PSEditions
    CompatiblePSEditions = @('Desktop', 'Core')
    
    # ID used to uniquely identify this module
    GUID = 'f8a6c8d0-1234-5678-9abc-def012345678'
    
    # Author of this module
    Author = 'HashiCorp'
    
    # Company or vendor of this module
    CompanyName = 'HashiCorp'
    
    # Copyright statement for this module
    Copyright = '(c) HashiCorp. All rights reserved.'
    
    # Description of the functionality provided by this module
    Description = 'AI setup tools for Terraform AzureRM Provider development with GitHub Copilot integration'
    
    # Minimum version of the PowerShell engine required by this module
    PowerShellVersion = '5.1'
    
    # Functions to export from this module, for best performance, do not use wildcards and do not delete the entry, use an empty array if there are no functions to export.
    FunctionsToExport = @(
        # Core Functions
        'Find-RepositoryRoot',
        'Test-Prerequisites', 
        'New-SafeBackup',
        'Test-FileIntegrity',
        
        # VS Code Setup
        'Get-VSCodeUserSettingsPath',
        'Test-PreviousInstallation',
        'Backup-VSCodeSettings',
        'Install-VSCodeSettings',
        'Get-CopilotSettings',
        
        # Copilot Installation
        'Test-CopilotInstallation',
        'Install-CopilotInstructions',
        'Remove-CopilotInstructions',
        'Get-CopilotInstallationInfo',
        
        # Main Installation
        'Install-TerraformAzureRMAI',
        'Remove-TerraformAzureRMAI',
        'Show-InstallationStatus'
    )
    
    # Cmdlets to export from this module, for best performance, do not use wildcards and do not delete the entry, use an empty array if there are no cmdlets to export.
    CmdletsToExport = @()
    
    # Variables to export from this module
    VariablesToExport = @()
    
    # Aliases to export from this module, for best performance, do not use wildcards and do not delete the entry, use an empty array if there are no aliases to export.
    AliasesToExport = @()
    
    # Private data to pass to the module specified in RootModule/ModuleToProcess. This may also contain a PSData hashtable with additional module metadata used by PowerShell.
    PrivateData = @{
        PSData = @{
            # Tags applied to this module. These help with module discovery in online galleries.
            Tags = @('Terraform', 'Azure', 'AzureRM', 'AI', 'Copilot', 'Development')
            
            # A URL to the license for this module.
            LicenseUri = 'https://github.com/hashicorp/terraform-provider-azurerm/blob/main/LICENSE'
            
            # A URL to the main website for this project.
            ProjectUri = 'https://github.com/hashicorp/terraform-provider-azurerm'
            
            # ReleaseNotes of this module
            ReleaseNotes = 'Initial release of AI setup tools for Terraform AzureRM Provider'
        }
    }
}
