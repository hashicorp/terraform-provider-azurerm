# Terraform AzureRM Provider AI Setup - Modular Installation System

A modular, maintainable installation system for GitHub Copilot instruction files and AI-enhanced development environment setup.

## Overview

This modular installation system replaces the monolithic approach with a well-structured, testable, and maintainable solution for setting up AI-enhanced development tools for the Terraform AzureRM Provider.

## Features

- **Modular Architecture**: Separated concerns into focused modules
- **Comprehensive Testing**: Full test suite validates all components
- **Backup & Restore**: Automatic backup of existing configurations
- **Cross-Platform**: Supports Windows, Linux, and macOS
- **Text-Based UI**: No emoji dependencies, better terminal compatibility
- **Error Handling**: Robust error handling and recovery mechanisms

## Directory Structure

```
test-modular-install/
├── modules/                    # Core functionality modules
│   ├── core-functions.ps1     # Repository discovery, backups, prerequisites
│   ├── vscode-setup.ps1       # VS Code settings management
│   ├── copilot-install.ps1    # Copilot instruction file installation
│   └── cleanup-restore.ps1    # Cleanup and restore operations
├── main-installer.ps1         # Main orchestrator script
├── test-runner.ps1            # Comprehensive test suite
├── fake-repo/                 # Test repository structure
│   ├── go.mod                 # Fake go.mod for testing
│   └── .github/               # Test .github directory
└── README.md                  # This documentation
```

## Quick Start

### Installation

```powershell
# Run with auto-discovery
.\main-installer.ps1

# Specify repository path
.\main-installer.ps1 -RepositoryPath "C:\Code\terraform-provider-azurerm"

# Silent installation (no prompts)
.\main-installer.ps1 -AutoApprove
```

### Cleanup

```powershell
# Remove AI setup and restore backups
.\main-installer.ps1 -Clean

# Silent cleanup
.\main-installer.ps1 -Clean -AutoApprove
```

### Information

```powershell
# Show help
.\main-installer.ps1 -Help

# Show backup information
.\main-installer.ps1 -ShowBackups
```

## Testing

Run the comprehensive test suite to validate all components:

```powershell
.\test-runner.ps1
```

The test suite validates:
- Module imports and dependencies
- Core functionality (repository discovery, backups)
- VS Code settings management
- Copilot installation functions
- Cleanup and restore operations
- Main installer integration

## Module Documentation

### core-functions.ps1

**Purpose**: Foundational utilities used by all other modules

**Key Functions**:
- `Find-RepositoryRoot`: Auto-discovers terraform-provider-azurerm repository
- `New-SafeBackup`: Creates timestamped backups with integrity checks
- `Test-Prerequisites`: Validates PowerShell version and VS Code availability
- `Test-FileIntegrity`: Compares file hashes for backup validation

**Dependencies**: None (foundational module)

### vscode-setup.ps1

**Purpose**: Manages VS Code settings.json configuration for Copilot integration

**Key Functions**:
- `Get-VSCodeUserSettingsPath`: Cross-platform settings.json path detection
- `Test-PreviousInstallation`: Detects existing Copilot configurations
- `Install-VSCodeSettings`: Merges new settings with existing configuration
- `Get-CopilotSettings`: Generates optimized Copilot settings for AzureRM development
- `Remove-CopilotSettings`: Clean removal of Copilot-specific settings

**Dependencies**: core-functions.ps1

### copilot-install.ps1

**Purpose**: Handles Copilot instruction file installation and management

**Key Functions**:
- `Install-CopilotInstructions`: Copies instruction files with backup creation
- `Remove-CopilotInstructions`: Clean removal of instruction files
- `Test-CopilotInstallation`: Validates installation completeness
- `Restore-CopilotInstructions`: Restores instruction files from backups

**Dependencies**: core-functions.ps1

### cleanup-restore.ps1

**Purpose**: Comprehensive cleanup and backup restoration functionality

**Key Functions**:
- `Remove-AISetup`: Complete AI setup removal with backup restoration
- `Restore-FromBackup`: Restores all components from backup directory
- `Get-BackupInfo`: Provides detailed backup information
- `Show-BackupSummary`: User-friendly backup summary display
- `Remove-OldBackups`: Automatic cleanup of old backup files

**Dependencies**: core-functions.ps1, vscode-setup.ps1, copilot-install.ps1

## Error Handling

The modular system implements comprehensive error handling:

- **Graceful Degradation**: Continues operation when non-critical components fail
- **Detailed Logging**: Text-based status indicators for better compatibility
- **Backup Recovery**: Automatic backup creation and restoration on failure
- **Validation**: File integrity checks and installation verification
- **User Feedback**: Clear status messages and actionable error information

## Backup System

### Automatic Backups

The system automatically creates timestamped backups:
- VS Code settings.json → `settings.json.backup_YYYYMMDD_HHMMSS`
- Copilot instructions → `copilot-instructions.md.backup_YYYYMMDD_HHMMSS`
- Instruction directories → Full recursive backup with timestamp

### Backup Management

```powershell
# View backup information
.\main-installer.ps1 -ShowBackups

# Manual cleanup of old backups (keeps 5 most recent)
Import-Module .\modules\cleanup-restore.ps1
Remove-OldBackups -BackupDir "\.github\backups" -RetainCount 5
```

## Configuration

### VS Code Settings

The system installs optimized Copilot settings for AzureRM development:

```json
{
    "github.copilot.enable": {
        "*": true,
        "go": true,
        "powershell": true,
        "markdown": true
    },
    "github.copilot.chat.welcomeMessage": "never",
    "github.copilot.advanced": {
        "secret_key": "terraform-azurerm-provider",
        "length": 1000,
        "stops": {
            "go": ["func", "type", "var", "const", "package", "import"]
        }
    }
}
```

### File Locations

- **Backup Directory**: `.github/backups/` (within repository)
- **VS Code Settings**: Platform-specific user settings.json
- **Instruction Files**: `.github/copilot-instructions.md` and `.github/instructions/`

## Migration from Monolithic System

To migrate from the original monolithic installation script:

1. **Test the modular system**: Run `.\test-runner.ps1` to validate
2. **Backup existing setup**: Use `-ShowBackups` to see current backups
3. **Clean install**: Run `.\main-installer.ps1 -Clean` then `.\main-installer.ps1`
4. **Verify functionality**: Test VS Code Copilot integration
5. **Replace scripts**: Update deployment/documentation to use modular system

## Troubleshooting

### Common Issues

**Module Import Errors**:
```powershell
# Verify all modules exist
Get-ChildItem .\modules\*.ps1

# Test individual module import
Import-Module .\modules\core-functions.ps1 -Force
```

**Repository Discovery Fails**:
```powershell
# Manual path specification
.\main-installer.ps1 -RepositoryPath "C:\Path\To\terraform-provider-azurerm"
```

**VS Code Settings Corruption**:
```powershell
# Restore from backup
.\main-installer.ps1 -Clean
# Check backup directory for manual restoration
.\main-installer.ps1 -ShowBackups
```

### Test Mode

The system automatically detects when running in test mode and adjusts behavior:
- Uses fake repository for testing
- Creates isolated test backups
- Provides test-specific logging
- Skips real file modifications

## Development

### Adding New Modules

1. Create module in `modules/` directory
2. Follow naming convention: `feature-action.ps1`
3. Use approved PowerShell verbs
4. Include proper module exports: `Export-ModuleMember -Function ...`
5. Add tests to `test-runner.ps1`
6. Update this documentation

### Testing Guidelines

- Each module should be testable in isolation
- Use the fake-repo structure for testing
- Implement both positive and negative test cases
- Validate error handling and edge cases
- Test cross-platform compatibility

## Support

For issues with the modular installation system:

1. Run the test suite: `.\test-runner.ps1`
2. Check backup status: `.\main-installer.ps1 -ShowBackups`
3. Review error messages for specific module failures
4. Test individual modules in isolation
5. Report issues with full error output and system information

## Version History

- **v2.0.0**: Initial modular release
  - Separated monolithic script into focused modules
  - Added comprehensive test suite
  - Improved error handling and backup system
  - Removed emoji dependencies for better compatibility
  - Enhanced cross-platform support
