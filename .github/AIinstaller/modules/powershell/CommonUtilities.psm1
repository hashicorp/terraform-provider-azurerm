# CommonUtilities Module for Terraform AzureRM Provider AI Setup
# Shared utility functions used across multiple modules

#region Cross-Platform Utilities

function Get-UserHomeDirectory {
    <#
    .SYNOPSIS
    Get the user's home directory in a cross-platform way

    .DESCRIPTION
    Returns the appropriate home directory path for the current platform:
    - Windows: Uses $env:USERPROFILE
    - macOS/Linux: Uses $env:HOME

    Handles PowerShell 5.1 compatibility where $IsWindows variable doesn't exist.

    .OUTPUTS
    String - The path to the user's home directory

    .EXAMPLE
    $homeDir = Get-UserHomeDirectory
    $installerPath = Join-Path $homeDir ".terraform-ai-installer"
    #>
    if ($IsWindows -or $env:OS -eq "Windows_NT" -or (-not $PSVersionTable.Platform)) {
        # Windows (including PowerShell 5.1 which doesn't have $IsWindows)
        return $env:USERPROFILE
    } else {
        # macOS and Linux
        return $env:HOME
    }
}

function Get-CrossPlatformInstallerPath {
    <#
    .SYNOPSIS
    Get the installer path in a cross-platform way for display purposes

    .DESCRIPTION
    Returns a formatted path string for the installer directory that's appropriate
    for display in user messages and documentation.

    .OUTPUTS
    String - Formatted installer path with quotes for display

    .EXAMPLE
    Write-Host "Navigate to: $(Get-CrossPlatformInstallerPath)"
    #>
    $homeDir = Get-UserHomeDirectory
    if ($IsWindows -or $env:OS -eq "Windows_NT" -or (-not $PSVersionTable.Platform)) {
        # Windows style path
        return "`"$homeDir\.terraform-ai-installer`""
    } else {
        # Unix style path
        return "`"$homeDir/.terraform-ai-installer`""
    }
}

#endregion

# Export public functions
Export-ModuleMember -Function Get-UserHomeDirectory, Get-CrossPlatformInstallerPath
