# Core Module - Clean Architecture
# Core functionality using real patterns from working system

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

function Initialize-CleanEnvironment {
    <#
    .SYNOPSIS
        Initializes the clean environment (minimal real implementation)
    #>
    param([string]$RepositoryPath)
    
    Write-StatusMessage "Initializing clean environment..." "Info"
    
    # Basic validation that repo path exists
    if (-not (Test-Path $RepositoryPath)) {
        throw "Repository path does not exist: $RepositoryPath"
    }
    
    return $true
}

function Get-InstallationConfig {
    <#
    .SYNOPSIS
        Gets installation configuration (real pattern from working system)
    #>
    param([string]$RepositoryPath)
    
    # Real pattern - read from config files that exist in repository
    $configPath = Join-Path $RepositoryPath ".github\AILocalInstall"
    
    return @{
        RepositoryPath = $RepositoryPath
        ConfigPath = $configPath
        BackupEnabled = $true
        VSCodeUserPath = Join-Path $env:APPDATA "Code\User"
    }
}

Export-ModuleMember -Function Find-RepositoryRoot, Initialize-CleanEnvironment, Get-InstallationConfig
