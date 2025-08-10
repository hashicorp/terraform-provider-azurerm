# Core Functions Module
# Functions for repository discovery, file operations, and basic utilities

function Find-RepositoryRoot {
    <#
    .SYNOPSIS
        Finds the terraform-provider-azurerm repository root directory
    #>
    param([string]$StartPath = (Get-Location).Path)
    
    $currentPath = $StartPath
    while ($currentPath -and $currentPath -ne [System.IO.Path]::GetPathRoot($currentPath)) {
        $gitPath = Join-Path $currentPath ".git"
        $goModPath = Join-Path $currentPath "go.mod"
        
        # Check if this is the terraform-provider-azurerm repository
        if ((Test-Path $gitPath) -and (Test-Path $goModPath)) {
            $goModContent = Get-Content $goModPath -Raw -ErrorAction SilentlyContinue
            if ($goModContent -match "module github\.com/hashicorp/terraform-provider-azurerm") {
                return $currentPath
            }
        }
        
        $currentPath = Split-Path $currentPath -Parent
    }
    
    return $null
}

function Test-FileIntegrity {
    <#
    .SYNOPSIS
        Validates file content integrity
    #>
    param([string]$FilePath, [string]$ExpectedPattern = "")
    
    if (-not (Test-Path $FilePath)) {
        return $false
    }
    
    try {
        $content = Get-Content $FilePath -Raw -ErrorAction Stop
        
        # Basic length check (minimum viable content)
        if ($content.Length -lt 50) {
            return $false
        }
        
        # Pattern check if provided
        if ($ExpectedPattern -and $content -notmatch $ExpectedPattern) {
            return $false
        }
        
        return $true
    } catch {
        return $false
    }
}

function Test-Prerequisites {
    <#
    .SYNOPSIS
        Checks if system prerequisites are met
    #>
    
    # Check PowerShell version
    if ($PSVersionTable.PSVersion.Major -lt 5) {
        Write-Warning "PowerShell 5.1 or later is required"
        return $false
    }
    
    # Check VS Code installation
    $vscodeInstalled = $false
    $vscodeLocations = @(
        "${env:LOCALAPPDATA}\Programs\Microsoft VS Code\Code.exe",
        "${env:ProgramFiles}\Microsoft VS Code\Code.exe",
        "${env:ProgramFiles(x86)}\Microsoft VS Code\Code.exe"
    )
    
    foreach ($location in $vscodeLocations) {
        if (Test-Path $location) {
            $vscodeInstalled = $true
            break
        }
    }
    
    if (-not $vscodeInstalled) {
        Write-Warning "VS Code installation not found"
        return $false
    }
    
    return $true
}

function Write-StatusMessage {
    <#
    .SYNOPSIS
        Writes formatted status messages
    #>
    param(
        [string]$Message,
        [ValidateSet("Info", "Success", "Warning", "Error")]
        [string]$Type = "Info"
    )
    
    $colors = @{
        "Info" = "Cyan"
        "Success" = "Green"
        "Warning" = "Yellow"
        "Error" = "Red"
    }
    
    $icons = @{
        "Info" = "[INFO]"
        "Success" = "[SUCCESS]"
        "Warning" = "[WARNING]"
        "Error" = "[ERROR]"
    }
    
    Write-Host "$($icons[$Type]) $Message" -ForegroundColor $colors[$Type]
}
