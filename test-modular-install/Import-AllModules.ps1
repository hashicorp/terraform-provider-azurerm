#Requires -Version 5.1

<#
.SYNOPSIS
    Import helper script for local development of Terraform AzureRM Provider testing modules.

.DESCRIPTION
    This script simplifies importing all PowerShell modules for local development and testing.
    Designed for use within the Terraform AzureRM Provider repository.

.PARAMETER Force
    Force reimport of modules (useful during development).

.PARAMETER ModuleName
    Import only the specified module.

.EXAMPLE
    .\Import-AllModules.ps1
    Import all modules for local use.

.EXAMPLE
    .\Import-AllModules.ps1 -Force
    Force reimport all modules (useful when developing).

.EXAMPLE
    .\Import-AllModules.ps1 -ModuleName "TerraformAzureRM.Testing"
    Import only the Testing module.
#>

[CmdletBinding()]
param(
    [switch]$Force,
    [ValidateSet("TerraformAzureRM.Testing", "TerraformAzureRM.Environment", "TerraformAzureRM.Validation", "TerraformAzureRM.Deployment")]
    [string]$ModuleName
)

$ErrorActionPreference = "Stop"

Write-Host "Terraform AzureRM Provider - Module Import Helper" -ForegroundColor Cyan
Write-Host "=" * 50

# Get the script directory
$ScriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path

# Define available modules
$AvailableModules = @(
    "TerraformAzureRM.Testing",
    "TerraformAzureRM.Environment", 
    "TerraformAzureRM.Validation",
    "TerraformAzureRM.Deployment"
)

# Determine which modules to import
$ModulesToImport = if ($ModuleName) {
    @($ModuleName)
} else {
    $AvailableModules
}

# Import modules
foreach ($Module in $ModulesToImport) {
    $ModulePath = Join-Path $ScriptPath "$Module\$Module.psm1"
    
    if (Test-Path $ModulePath) {
        Write-Host "Importing module: $Module" -ForegroundColor Green
        
        try {
            if ($Force) {
                Import-Module $ModulePath -Force -Global
            } else {
                Import-Module $ModulePath -Global
            }
            
            # Display module info
            $ModuleInfo = Get-Module $Module
            if ($ModuleInfo) {
                Write-Host "  Version: $($ModuleInfo.Version)" -ForegroundColor Gray
                Write-Host "  Functions: $($ModuleInfo.ExportedFunctions.Count)" -ForegroundColor Gray
            }
        }
        catch {
            Write-Warning "Failed to import module $Module`: $_"
        }
    } else {
        Write-Warning "Module not found: $ModulePath"
    }
}

Write-Host "`nModule import complete!" -ForegroundColor Cyan

# Display available functions
Write-Host "`nAvailable functions:" -ForegroundColor Yellow
$ImportedModules = Get-Module | Where-Object { $_.Name -like "TerraformAzureRM.*" }

foreach ($Module in $ImportedModules) {
    Write-Host "`n$($Module.Name):" -ForegroundColor Cyan
    $Module.ExportedFunctions.Keys | Sort-Object | ForEach-Object {
        Write-Host "  - $_" -ForegroundColor Gray
    }
}

Write-Host "`nTip: Use 'Get-Help <FunctionName> -Detailed' for function documentation" -ForegroundColor Yellow
