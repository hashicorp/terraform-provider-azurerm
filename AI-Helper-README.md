# AI Development Helper

**A single-file, zero-dependency solution for AI-enhanced Terraform AzureRM Provider development.**

## Quick Start

```powershell
# Navigate to the provider repository
cd "c:\github.com\hashicorp\terraform-provider-azurerm"

# Load the helper functions
. .\ai-dev-helper.ps1

# Show available commands
Quick-Commands

# Validate your environment
Validate-Environment
```

## Key Features

- **Zero Dependencies**: Single PowerShell file, no modules or external deps
- **Simple Functions**: Direct, focused functions for common tasks
- **AI Agent Ready**: Perfect for GitHub Copilot integration
- **Local Development**: Designed for local development workflow

## Common Commands

```powershell
# Build and test
Build-Provider
Test-AcceptanceTest -Service "cdn" -Test "basic"

# Development helpers
Find-Resource -Pattern "frontdoor" 
Show-ServiceStructure -Service "cdn"

# Generate test templates
Generate-ResourceTest -Service "cdn" -ResourceName "profile"
```

## Usage Examples

### Test a specific resource
```powershell
Test-AcceptanceTest -Service "cdn" -Test "basic" -Resource "CdnFrontDoorProfile"
```

### Find and explore resources
```powershell
Find-Resource -Pattern "cdn"
Show-ServiceStructure -Service "cdn"
```

### Environment validation
```powershell
Validate-Environment  # Full environment check
Test-Prerequisites    # Just check prerequisites
```

## For AI Agents (GitHub Copilot)

This helper is designed to work seamlessly with AI coding assistants:

1. **Simple Command Structure**: Easy for AI to understand and use
2. **Clear Function Names**: Self-documenting function names
3. **Minimal Dependencies**: No complex module system to navigate
4. **Common Patterns**: Follows standard PowerShell patterns

## File Structure

```
terraform-provider-azurerm/
├── ai-dev-helper.ps1     # This file - all functions included
├── AI-Helper-README.md   # This documentation
└── [rest of provider files...]
```

That's it! No complex directory structures, no dependency management, just one file with everything you need.

## Functions Reference

| Function | Purpose |
|----------|---------|
| `Build-Provider` | Build the Terraform provider |
| `Test-Build` | Quick build validation |
| `Test-AcceptanceTest` | Run acceptance tests |
| `Test-UnitTests` | Run unit tests |
| `Find-Resource` | Search for resource files |
| `Show-ServiceStructure` | Display service directory structure |
| `Generate-ResourceTest` | Generate test templates |
| `Validate-Environment` | Complete environment validation |
| `Test-Prerequisites` | Check development prerequisites |
| `Quick-Commands` | Show available commands |

## Perfect for:

- GitHub Copilot integration
- Local development workflows  
- Quick testing and validation
- AI-assisted development
- Simple automation without complexity
