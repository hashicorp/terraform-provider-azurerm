# CoreFunctions.psm1 - Core utility functions for Terraform AzureRM Provider AI Setup

function Find-RepositoryRoot {
    <#
    .SYNOPSIS
    Finds the terraform-provider-azurerm repository root directory
    
    .DESCRIPTION
    Searches for the repository by looking for key markers like go.mod and main.go
    
    .OUTPUTS
    String path to repository root, or $null if not found
    #>
    
    param(
        [string]$StartPath = (Get-Location).Path
    )
    
    $currentPath = $StartPath
    $maxDepth = 10
    $depth = 0
    
    while ($depth -lt $maxDepth) {
        # Check for terraform-provider-azurerm markers
        $goModPath = Join-Path $currentPath "go.mod"
        $mainGoPath = Join-Path $currentPath "main.go"
        $internalPath = Join-Path $currentPath "internal"
        
        if ((Test-Path $goModPath) -and (Test-Path $mainGoPath) -and (Test-Path $internalPath)) {
            # Verify it's actually the azurerm provider
            $goModContent = Get-Content $goModPath -Raw
            if ($goModContent -match "terraform-provider-azurerm") {
                Write-Verbose "Found terraform-provider-azurerm repository at: $currentPath"
                return $currentPath
            }
        }
        
        # Move up one directory
        $parentPath = Split-Path $currentPath -Parent
        if ($parentPath -eq $currentPath) {
            # Reached filesystem root
            break
        }
        
        $currentPath = $parentPath
        $depth++
    }
    
    return $null
}

function Test-Prerequisites {
    <#
    .SYNOPSIS
    Tests if all prerequisites are met for AI setup installation
    
    .DESCRIPTION
    Checks for VS Code, PowerShell version, and other requirements
    
    .OUTPUTS
    Boolean - True if all prerequisites are met
    #>
    
    # Check PowerShell version
    if ($PSVersionTable.PSVersion.Major -lt 5) {
        throw "PowerShell 5.0 or later is required"
    }
    
    # Check if VS Code is installed (optional but recommended)
    $vsCodePaths = @(
        "${env:LOCALAPPDATA}\Programs\Microsoft VS Code\Code.exe",
        "${env:PROGRAMFILES}\Microsoft VS Code\Code.exe",
        "${env:PROGRAMFILES(X86)}\Microsoft VS Code\Code.exe"
    )
    
    $vsCodeFound = $false
    foreach ($path in $vsCodePaths) {
        if (Test-Path $path) {
            $vsCodeFound = $true
            break
        }
    }
    
    if (-not $vsCodeFound) {
        Write-Warning "VS Code not found in standard locations. Some features may not work."
    }
    
    Write-Verbose "Prerequisites check completed successfully"
    return $true
}

function New-SafeBackup {
    <#
    .SYNOPSIS
    Creates a safe backup of a file with timestamp
    
    .PARAMETER SourcePath
    Path to the file to backup
    
    .PARAMETER BackupDirectory
    Directory where backup should be stored
    
    .OUTPUTS
    String path to created backup file
    #>
    
    param(
        [Parameter(Mandatory = $true)]
        [string]$SourcePath,
        
        [Parameter(Mandatory = $true)]
        [string]$BackupDirectory
    )
    
    if (-not (Test-Path $SourcePath)) {
        throw "Source file not found: $SourcePath"
    }
    
    if (-not (Test-Path $BackupDirectory)) {
        New-Item -ItemType Directory -Path $BackupDirectory -Force | Out-Null
    }
    
    $timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
    $fileName = Split-Path $SourcePath -Leaf
    $backupFileName = "$timestamp-$fileName"
    $backupPath = Join-Path $BackupDirectory $backupFileName
    
    Copy-Item $SourcePath $backupPath -Force
    Write-Verbose "Created backup: $backupPath"
    
    return $backupPath
}

function Test-FileIntegrity {
    <#
    .SYNOPSIS
    Tests if a file exists and is readable
    
    .PARAMETER Path
    Path to test
    
    .OUTPUTS
    Boolean indicating if file is accessible
    #>
    
    param(
        [Parameter(Mandatory = $true)]
        [string]$Path
    )
    
    if (-not (Test-Path $Path)) {
        return $false
    }
    
    try {
        $null = Get-Content $Path -TotalCount 1 -ErrorAction Stop
        return $true
    }
    catch {
        return $false
    }
}

# Export module members
Export-ModuleMember -Function Find-RepositoryRoot, Test-Prerequisites, New-SafeBackup, Test-FileIntegrity
