# Terraform AzureRM Provider Testing Modules

A comprehensive suite of PowerShell modules designed to streamline testing and development workflows for the Terraform AzureRM Provider. These modules are designed for local development use only.

## Overview

This collection provides modular, reusable PowerShell tools that automate common testing scenarios, environment setup, and validation tasks for local Terraform AzureRM Provider development.

## Modules Included

### 🧪 **TerraformAzureRM.Testing**
Core testing functionality for Terraform AzureRM Provider development.

**Key Features:**
- Automated acceptance test execution
- Test result analysis and reporting
- Resource cleanup verification
- Azure credential validation

### 🔧 **TerraformAzureRM.Environment**
Environment management and setup utilities.

**Key Features:**
- Azure environment configuration
- Credential management
- Provider installation and verification
- Development environment validation

### 📊 **TerraformAzureRM.Validation**
Resource validation and compliance checking tools.

**Key Features:**
- Resource state validation
- Azure API compliance verification
- Schema validation
- Configuration drift detection

### 🚀 **TerraformAzureRM.Deployment**
Deployment automation and management utilities.

**Key Features:**
- Automated deployment workflows
- Resource lifecycle management
- Rollback capabilities
- Deployment verification

## Quick Start

### Local Installation

```powershell
# Navigate to the modules directory
Set-Location "c:\github.com\hashicorp\terraform-provider-azurerm\test-modular-install"

# Import all modules
Import-Module .\TerraformAzureRM.Testing -Force
Import-Module .\TerraformAzureRM.Environment -Force
Import-Module .\TerraformAzureRM.Validation -Force
Import-Module .\TerraformAzureRM.Deployment -Force
```

### Basic Usage

```powershell
# Set up testing environment
Initialize-TerraformAzureRMEnvironment

# Run acceptance tests for a specific service
Invoke-TerraformAcceptanceTest -ServiceName "cdn" -TestPattern "TestAccCdnFrontDoorProfile_basic"

# Validate deployment
Test-AzureResourceDeployment -ResourceGroupName "test-rg" -ResourceType "Microsoft.Cdn/profiles"

# Clean up test resources
Remove-TestResources -ResourceGroupPattern "acctest-*"
```

## Module Documentation

Each module includes comprehensive documentation:

- **Module Help**: Use `Get-Help <ModuleName>` for overview
- **Function Help**: Use `Get-Help <FunctionName> -Detailed` for specific functions
- **Examples**: Use `Get-Help <FunctionName> -Examples` for usage examples

## Testing

A comprehensive test suite is included to validate module functionality:

```powershell
# Run all module tests
.\Test-AllModules.ps1

# Run tests for specific module
.\Test-AllModules.ps1 -ModuleName "TerraformAzureRM.Testing"

# Run tests with verbose output
.\Test-AllModules.ps1 -VerboseMode
```

## Development Standards

These modules follow PowerShell best practices:

- **Approved Verbs**: All functions use approved PowerShell verbs
- **Parameter Validation**: Comprehensive input validation
- **Error Handling**: Structured error handling with meaningful messages
- **Help Documentation**: Complete help documentation for all functions
- **Module Manifests**: Proper module manifests with metadata
- **Local Development**: Designed for local development workflows only

## Usage Notes

- These modules are designed for local development and testing only
- No external dependencies or publishing requirements
- All modules are self-contained within this directory
- Perfect for Terraform AzureRM Provider contributors and maintainers

## Local Development Features

### Import Helpers
```powershell
# Quick import all modules
.\Import-AllModules.ps1

# Import with reload (for development)
.\Import-AllModules.ps1 -Force
```

### Development Testing
```powershell
# Test module syntax and structure
Test-ModuleStructure -ModuleName "TerraformAzureRM.Testing"

# Validate help documentation
Test-ModuleHelp -ModuleName "TerraformAzureRM.Environment"

# Check for best practices compliance
Test-ModuleCompliance -ModuleName "TerraformAzureRM.Validation"
```

### Example Workflows

#### Testing New Resource Implementation
```powershell
# Set up environment
Initialize-TerraformAzureRMEnvironment

# Test specific resource
Invoke-TerraformAcceptanceTest -ServiceName "cdn" -ResourceName "frontdoor_profile" -TestType "basic"

# Validate resource compliance
Test-ResourceCompliance -ResourceType "azurerm_cdn_frontdoor_profile"

# Clean up
Remove-TestResources -ServiceName "cdn"
```

#### Environment Validation
```powershell
# Check Azure credentials
Test-AzureCredentials

# Validate Terraform installation
Test-TerraformInstallation

# Check provider compilation
Test-ProviderCompilation
```

## Support

For issues, questions, or contributions related to these testing modules, please refer to the Terraform AzureRM Provider repository documentation and contribution guidelines.

## Architecture

These modules are designed with a clear separation of concerns:

- **Testing Module**: Handles all test execution and result analysis
- **Environment Module**: Manages configuration and setup
- **Validation Module**: Provides compliance and validation checks
- **Deployment Module**: Automates deployment workflows

Each module is independent but can work together for comprehensive workflows.
