# install-with-modules.ps1 - Example of using proper PowerShell modules

param(
    [string]$RepositoryPath,
    [switch]$Clean,
    [switch]$Status,
    [switch]$Force,
    [switch]$Help
)

if ($Help) {
    Write-Host @"
Terraform AzureRM Provider AI Setup (Proper PowerShell Modules)

USAGE:
  .\install-with-modules.ps1 [OPTIONS]

OPTIONS:
  -RepositoryPath <path>   Specify repository path
  -Clean                   Remove AI setup
  -Status                  Show installation status
  -Force                   Force installation/reinstallation
  -Help                    Show this help message

EXAMPLES:
  .\install-with-modules.ps1
  .\install-with-modules.ps1 -RepositoryPath C:\Code\terraform-provider-azurerm
  .\install-with-modules.ps1 -Clean
  .\install-with-modules.ps1 -Status

"@ -ForegroundColor Green
    return
}

# Import the main module
$ModulePath = Join-Path $PSScriptRoot "modules\TerraformAzureRMSetup.psd1"

if (-not (Test-Path $ModulePath)) {
    Write-Error "Module not found: $ModulePath"
    return
}

try {
    Import-Module $ModulePath -Force -Verbose:$false
    
    if ($Status) {
        Show-InstallationStatus -RepositoryPath $RepositoryPath
    }
    elseif ($Clean) {
        Remove-TerraformAzureRMAI -RepositoryPath $RepositoryPath
    }
    else {
        Install-TerraformAzureRMAI -RepositoryPath $RepositoryPath -Force:$Force
    }
}
catch {
    Write-Error "Operation failed: $($_.Exception.Message)"
}
finally {
    # Clean up module
    Remove-Module TerraformAzureRMSetup -ErrorAction SilentlyContinue
}
