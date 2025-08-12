# Configuration File Management Module for PowerShell
# This module provides functions to read the file-manifest.config file

function Read-ConfigSection {
    param(
        [string]$ConfigFile,
        [string]$Section
    )
    
    if (-not (Test-Path $ConfigFile)) {
        Write-Error "Configuration file not found: $ConfigFile"
        return @()
    }
    
    $inSection = $false
    $files = @()
    
    # Read file line by line - handles both LF and CRLF
    foreach ($line in [System.IO.File]::ReadAllLines($ConfigFile)) {
        # Skip empty lines and comments
        if ([string]::IsNullOrWhiteSpace($line) -or $line.TrimStart().StartsWith('#')) {
            continue
        }
        
        # Check for section headers [SECTION_NAME]
        if ($line -match '^\[(.+)\]$') {
            $currentSection = $matches[1]
            if ($currentSection -eq $Section) {
                $inSection = $true
            } else {
                $inSection = $false
            }
            continue
        }
        
        # If we're in the target section, add the line (trimmed)
        if ($inSection) {
            $trimmedLine = $line.Trim()
            if (-not [string]::IsNullOrWhiteSpace($trimmedLine)) {
                $files += $trimmedLine
            }
        }
    }
    
    return $files
}

function Get-InstructionFiles {
    param(
        [string]$ConfigFile = (Join-Path $PSScriptRoot "..\..\file-manifest.config")
    )
    
    return Read-ConfigSection -ConfigFile $ConfigFile -Section "INSTRUCTION_FILES"
}

function Get-PromptFiles {
    param(
        [string]$ConfigFile = (Join-Path $PSScriptRoot "..\..\file-manifest.config")
    )
    
    return Read-ConfigSection -ConfigFile $ConfigFile -Section "PROMPT_FILES"
}

function Get-MainFiles {
    param(
        [string]$ConfigFile = (Join-Path $PSScriptRoot "..\..\file-manifest.config")
    )
    
    return Read-ConfigSection -ConfigFile $ConfigFile -Section "MAIN_FILES"
}

# Legacy support - keep existing function names working but use config file
function Get-ExpectedInstructionFiles {
    param(
        [string]$ConfigFile = (Join-Path $PSScriptRoot "..\..\file-manifest.config")
    )
    
    return Get-InstructionFiles -ConfigFile $ConfigFile
}

function Get-ExpectedPromptFiles {
    param(
        [string]$ConfigFile = (Join-Path $PSScriptRoot "..\..\file-manifest.config")
    )
    
    return Get-PromptFiles -ConfigFile $ConfigFile
}

Export-ModuleMember -Function @(
    'Read-ConfigSection',
    'Get-InstructionFiles', 
    'Get-PromptFiles',
    'Get-MainFiles',
    'Get-ExpectedInstructionFiles',
    'Get-ExpectedPromptFiles'
)
