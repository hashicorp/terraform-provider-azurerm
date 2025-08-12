# Cross-Platform Support

The Terraform AzureRM Provider AI Setup installer now supports multiple platforms with automatic VS Code User directory detection.

## Supported Platforms

### Windows
- **VS Code User Directory**: `%APPDATA%\Code\User`
- **Path Example**: `C:\Users\username\AppData\Roaming\Code\User`
- **Detection**: `$IsWindows` or `$env:OS -eq "Windows_NT"`

### macOS
- **VS Code User Directory**: `~/Library/Application Support/Code/User`
- **Path Example**: `/Users/username/Library/Application Support/Code/User`
- **Detection**: `$IsMacOS` or file exists `/System/Library/CoreServices/SystemVersion.plist`

### Linux/Unix
- **VS Code User Directory**: `~/.config/Code/User`
- **Path Example**: `/home/username/.config/Code/User`
- **Detection**: Default fallback for non-Windows/non-macOS systems

## Installation Manifest Changes

The simple manifest now uses:
- `"relativeTo": "VSCODE_USER_DIR"` - Cross-platform VS Code User directory
- `"relativeTo": "APPDATA"` - Legacy Windows-only support (backward compatibility)

## How Platform Detection Works

```powershell
function Get-VSCodeUserDirectory {
    if ($IsWindows -or $env:OS -eq "Windows_NT") {
        return Join-Path $env:APPDATA "Code\User"
    }
    elseif ($IsMacOS -or (Test-Path "/System/Library/CoreServices/SystemVersion.plist")) {
        $homeDir = if ($env:HOME) { $env:HOME } else { "~" }
        return Join-Path $homeDir "Library/Application Support/Code/User"
    }
    else {
        $homeDir = if ($env:HOME) { $env:HOME } else { "~" }
        return Join-Path $homeDir ".config/Code/User"
    }
}
```

## File Structure After Installation

### All Platforms
```
VS Code User Directory/
├── copilot-instructions.md
├── instructions/
│   └── terraform-azurerm/
│       ├── implementation-guide.instructions.md
│       ├── azure-patterns.instructions.md
│       ├── testing-guidelines.instructions.md
│       └── ... (other instruction files)
└── prompts/
    └── terraform-azurerm/
        ├── add-unit-tests.prompt.md
        ├── code-review-committed-changes.prompt.md
        └── ... (other prompt files)
```

## Verification

The verification now works across all platforms:

```bash
# Windows PowerShell
.\install-copilot-setup.ps1 -Verify

# macOS/Linux PowerShell Core
pwsh ./install-copilot-setup.ps1 -Verify
```

The verification automatically detects the platform and checks files in the correct VS Code User directory location.
