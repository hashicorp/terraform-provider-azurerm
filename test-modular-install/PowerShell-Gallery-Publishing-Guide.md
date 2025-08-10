# PowerShell Gallery Publishing Guide

This guide walks through publishing the TerraformAzureRMSetup PowerShell module to the PowerShell Gallery for easy distribution and installation.

## Prerequisites

1. **PowerShell Gallery Account**
   - Create account at [PowerShell Gallery](https://www.powershellgallery.com/)
   - Get API key from your account settings

2. **Required PowerShell Modules**
   - `PowerShellGet` (usually pre-installed)
   - `PackageManagement` (usually pre-installed)

3. **Module Requirements**
   - Valid module manifest (.psd1)
   - Proper module structure
   - Unique module name
   - Appropriate version numbering

## Step 1: Prepare Module for Publishing

### Update Module Manifest

First, ensure the module manifest has all required fields for PowerShell Gallery:

```powershell
# Edit TerraformAzureRMSetup.psd1 to include:

@{
    # Script module or binary module file associated with this manifest.
    RootModule = 'TerraformAzureRMSetup.psm1'

    # Version number of this module.
    ModuleVersion = '1.0.0'

    # Supported PSEditions
    CompatiblePSEditions = @('Desktop', 'Core')

    # ID used to uniquely identify this module
    GUID = '12345678-1234-1234-1234-123456789012'  # Generate unique GUID

    # Author of this module
    Author = 'HashiCorp Terraform Team'

    # Company or vendor of this module
    CompanyName = 'HashiCorp'

    # Copyright statement for this module
    Copyright = '(c) 2024 HashiCorp. All rights reserved.'

    # Description of the functionality provided by this module
    Description = 'PowerShell module for setting up GitHub Copilot and VS Code for Terraform AzureRM Provider development'

    # Minimum version of the PowerShell engine required by this module
    PowerShellVersion = '5.1'

    # Functions to export from this module
    FunctionsToExport = @(
        'Install-TerraformAzureRMDevelopmentEnvironment',
        'Show-InstallationStatus',
        'Find-RepositoryRoot',
        'Test-Prerequisites',
        'Test-CopilotInstallation',
        'Get-VSCodeUserSettingsPath',
        'Install-VSCodeCopilotConfiguration',
        'Install-CopilotForRepository'
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
            # Tags applied to this module
            Tags = @('Terraform', 'AzureRM', 'Copilot', 'VSCode', 'Development', 'Setup', 'Azure', 'HashiCorp')

            # A URL to the license for this module.
            LicenseUri = 'https://github.com/hashicorp/terraform-provider-azurerm/blob/main/LICENSE'

            # A URL to the main website for this project.
            ProjectUri = 'https://github.com/hashicorp/terraform-provider-azurerm'

            # A URL to an icon representing this module.
            IconUri = 'https://github.com/hashicorp/terraform-provider-azurerm/raw/main/website/img/terraform-logo.png'

            # ReleaseNotes of this module
            ReleaseNotes = @'
## 1.0.0
- Initial release
- GitHub Copilot setup automation
- VS Code configuration management
- Terraform AzureRM Provider development environment setup
- Comprehensive prerequisite checking
- Backup and safety features
'@

            # External dependent modules of this module
            ExternalModuleDependencies = @()

            # Prerelease string of this module
            # Prerelease = 'beta1'
        }
    }
}
```

### Create Proper Module Structure

Ensure your module follows PowerShell Gallery conventions:

```
TerraformAzureRMSetup/
├── TerraformAzureRMSetup.psd1      # Module manifest
├── TerraformAzureRMSetup.psm1      # Main module file
├── CoreFunctions.psm1              # Core functions
├── VSCodeSetup.psm1                # VS Code setup functions
├── CopilotInstall.psm1             # Copilot installation functions
├── README.md                       # Module documentation
├── CHANGELOG.md                    # Version history
└── LICENSE                         # License file
```

## Step 2: Generate Unique GUID

```powershell
# Generate a unique GUID for your module
$guid = [System.Guid]::NewGuid()
Write-Host "Module GUID: $guid"
```

## Step 3: Test Module Locally

Before publishing, thoroughly test your module:

```powershell
# Test module manifest
Test-ModuleManifest .\TerraformAzureRMSetup.psd1

# Import and test functions
Import-Module .\TerraformAzureRMSetup.psd1 -Force

# Test key functions
Test-Prerequisites
Find-RepositoryRoot
Get-Command -Module TerraformAzureRMSetup

# Run your test suite
.\test-modules-simple.ps1 -ShowDetails
```

## Step 4: Prepare for Publishing

### Install Required Modules

```powershell
# Update PowerShellGet if needed
Install-Module PowerShellGet -Force -AllowClobber

# Check current version
Get-Module PowerShellGet -ListAvailable
```

### Set Up API Key

```powershell
# Get API key from PowerShell Gallery account settings
# Store it securely (this command will prompt for the key)
$apiKey = Read-Host "Enter your PowerShell Gallery API Key" -AsSecureString
$credential = New-Object System.Management.Automation.PSCredential("APIKey", $apiKey)

# Or set it directly (less secure - for automation only)
# $apiKey = "your-api-key-here"
```

## Step 5: Publish to PowerShell Gallery

### Option A: Publish from Local Directory

```powershell
# Navigate to your module directory
cd "C:\path\to\your\TerraformAzureRMSetup"

# Publish the module
Publish-Module -Path . -NuGetApiKey $apiKey -Repository PSGallery -Verbose

# Or with credential object
Publish-Module -Path . -Credential $credential -Repository PSGallery -Verbose
```

### Option B: Publish from PowerShell Module Path

```powershell
# Copy module to PowerShell module path
$modulePath = "$env:PSModulePath".Split(';')[0]
$destinationPath = Join-Path $modulePath "TerraformAzureRMSetup"

# Create destination directory
New-Item -Path $destinationPath -ItemType Directory -Force

# Copy all module files
Copy-Item -Path ".\*" -Destination $destinationPath -Recurse -Force

# Publish from module path
Publish-Module -Name "TerraformAzureRMSetup" -NuGetApiKey $apiKey -Repository PSGallery -Verbose
```

## Step 6: Verify Publication

After publishing, verify your module is available:

```powershell
# Search for your module
Find-Module TerraformAzureRMSetup

# Get detailed information
Find-Module TerraformAzureRMSetup | Format-List *

# Test installation in a new PowerShell session
Install-Module TerraformAzureRMSetup -Scope CurrentUser
Get-Module TerraformAzureRMSetup -ListAvailable
```

## Publishing Script Example

Here's a complete script to automate the publishing process:

```powershell
# publish-module.ps1

[CmdletBinding()]
param(
    [Parameter(Mandatory=$true)]
    [string]$ApiKey,
    
    [Parameter()]
    [string]$ModulePath = ".\modules",
    
    [Parameter()]
    [switch]$WhatIf,
    
    [Parameter()]
    [switch]$Force
)

# Function to update module version
function Update-ModuleVersion {
    param([string]$ManifestPath, [string]$NewVersion)
    
    $manifest = Import-PowerShellDataFile $ManifestPath
    $manifestContent = Get-Content $ManifestPath -Raw
    
    $manifestContent = $manifestContent -replace "ModuleVersion\s*=\s*'[^']*'", "ModuleVersion = '$NewVersion'"
    
    Set-Content -Path $ManifestPath -Value $manifestContent -Encoding UTF8
    Write-Host "Updated module version to $NewVersion" -ForegroundColor Green
}

# Function to create module package
function New-ModulePackage {
    param([string]$SourcePath, [string]$PackagePath)
    
    Write-Host "Creating module package..." -ForegroundColor Blue
    
    # Create package directory
    if (Test-Path $PackagePath) {
        Remove-Item $PackagePath -Recurse -Force
    }
    New-Item -Path $PackagePath -ItemType Directory -Force | Out-Null
    
    # Copy module files
    $moduleFiles = @(
        "TerraformAzureRMSetup.psd1",
        "TerraformAzureRMSetup.psm1", 
        "CoreFunctions.psm1",
        "VSCodeSetup.psm1",
        "CopilotInstall.psm1"
    )
    
    foreach ($file in $moduleFiles) {
        $sourcePath = Join-Path $SourcePath $file
        if (Test-Path $sourcePath) {
            Copy-Item $sourcePath -Destination $PackagePath -Force
            Write-Host "  Copied: $file" -ForegroundColor Gray
        } else {
            Write-Warning "Module file not found: $file"
        }
    }
    
    # Copy documentation
    $docFiles = @("README.md", "CHANGELOG.md", "LICENSE")
    foreach ($file in $docFiles) {
        $sourcePath = Join-Path (Split-Path $SourcePath) $file
        if (Test-Path $sourcePath) {
            Copy-Item $sourcePath -Destination $PackagePath -Force
            Write-Host "  Copied: $file" -ForegroundColor Gray
        }
    }
}

# Main publishing logic
try {
    Write-Host "==================================================================" -ForegroundColor Cyan
    Write-Host "  Publishing TerraformAzureRMSetup Module to PowerShell Gallery" -ForegroundColor Cyan
    Write-Host "==================================================================" -ForegroundColor Cyan
    Write-Host ""
    
    # Paths
    $sourceModulePath = Resolve-Path $ModulePath
    $packagePath = Join-Path $env:TEMP "TerraformAzureRMSetup"
    $manifestPath = Join-Path $sourceModulePath "TerraformAzureRMSetup.psd1"
    
    # Validate module manifest
    Write-Host "Validating module manifest..." -ForegroundColor Blue
    $manifest = Test-ModuleManifest $manifestPath -ErrorAction Stop
    Write-Host "  Module: $($manifest.Name)" -ForegroundColor Gray
    Write-Host "  Version: $($manifest.Version)" -ForegroundColor Gray
    Write-Host "  Author: $($manifest.Author)" -ForegroundColor Gray
    
    # Create package
    New-ModulePackage -SourcePath $sourceModulePath -PackagePath $packagePath
    
    # Test package
    Write-Host "Testing package..." -ForegroundColor Blue
    $packageManifest = Join-Path $packagePath "TerraformAzureRMSetup.psd1"
    Test-ModuleManifest $packageManifest -ErrorAction Stop | Out-Null
    Write-Host "  Package validation: PASSED" -ForegroundColor Green
    
    if ($WhatIf) {
        Write-Host "WhatIf: Would publish module from $packagePath" -ForegroundColor Yellow
        return
    }
    
    # Publish module
    Write-Host "Publishing to PowerShell Gallery..." -ForegroundColor Blue
    $publishParams = @{
        Path = $packagePath
        NuGetApiKey = $ApiKey
        Repository = "PSGallery"
        Verbose = $true
    }
    
    if ($Force) {
        $publishParams.Force = $true
    }
    
    Publish-Module @publishParams
    
    Write-Host ""
    Write-Host "Module published successfully!" -ForegroundColor Green
    Write-Host "You can install it with: Install-Module TerraformAzureRMSetup" -ForegroundColor Green
    
} catch {
    Write-Error "Publishing failed: $($_.Exception.Message)"
    exit 1
} finally {
    # Cleanup
    if (Test-Path $packagePath) {
        Remove-Item $packagePath -Recurse -Force -ErrorAction SilentlyContinue
    }
}
```

## Version Management

### Semantic Versioning

Follow semantic versioning (SemVer) for your module versions:

- **Major** (1.0.0): Breaking changes
- **Minor** (1.1.0): New features, backward compatible
- **Patch** (1.0.1): Bug fixes, backward compatible

### Updating Versions

```powershell
# Update version in manifest for new release
function Update-ModuleVersion {
    param(
        [string]$ManifestPath,
        [ValidateSet('Major', 'Minor', 'Patch')]
        [string]$BumpType = 'Patch'
    )
    
    $manifest = Import-PowerShellDataFile $ManifestPath
    $currentVersion = [Version]$manifest.ModuleVersion
    
    switch ($BumpType) {
        'Major' { $newVersion = [Version]::new($currentVersion.Major + 1, 0, 0) }
        'Minor' { $newVersion = [Version]::new($currentVersion.Major, $currentVersion.Minor + 1, 0) }
        'Patch' { $newVersion = [Version]::new($currentVersion.Major, $currentVersion.Minor, $currentVersion.Build + 1) }
    }
    
    # Update manifest file
    $content = Get-Content $ManifestPath -Raw
    $content = $content -replace "ModuleVersion\s*=\s*'[^']*'", "ModuleVersion = '$newVersion'"
    Set-Content -Path $ManifestPath -Value $content -Encoding UTF8
    
    Write-Host "Updated module version from $currentVersion to $newVersion" -ForegroundColor Green
}

# Usage
Update-ModuleVersion -ManifestPath ".\TerraformAzureRMSetup.psd1" -BumpType "Minor"
```

## Post-Publishing Tasks

### Monitor Downloads and Feedback

1. **PowerShell Gallery Statistics**
   - Check download counts
   - Monitor user ratings and feedback

2. **GitHub Integration**
   - Link PowerShell Gallery to GitHub repository
   - Set up automated publishing with GitHub Actions

3. **Documentation Updates**
   - Update README with installation instructions
   - Create examples and tutorials

### Maintenance

```powershell
# Check for module updates needed
Find-Module TerraformAzureRMSetup | Format-List Name, Version, PublishedDate

# Update module if needed
Update-Module TerraformAzureRMSetup

# Uninstall specific version if needed
Uninstall-Module TerraformAzureRMSetup -RequiredVersion "1.0.0"
```

## Troubleshooting Publishing Issues

### Common Publishing Errors

1. **"Module name already exists"**
   - Choose a unique module name
   - Check existing modules: `Find-Module YourModuleName`

2. **"Invalid manifest"**
   - Validate manifest: `Test-ModuleManifest .\YourModule.psd1`
   - Check required fields are present

3. **"API key invalid"**
   - Verify API key from PowerShell Gallery account
   - Check key permissions and expiration

4. **"Version already exists"**
   - Increment version number in manifest
   - Cannot republish same version

### Publishing Checklist

- [ ] Module manifest is valid (`Test-ModuleManifest`)
- [ ] Unique module name
- [ ] Proper version number
- [ ] All required fields in manifest
- [ ] Functions are properly exported
- [ ] Documentation is included
- [ ] API key is valid
- [ ] Module works after installation

## Example Usage After Publishing

Once published, users can install and use your module easily:

```powershell
# Install from PowerShell Gallery
Install-Module TerraformAzureRMSetup -Scope CurrentUser

# Import and use
Import-Module TerraformAzureRMSetup

# Run setup
Install-TerraformAzureRMDevelopmentEnvironment

# Check status
Show-InstallationStatus
```

This makes your Terraform AzureRM Provider development setup tools available to the entire PowerShell community!
