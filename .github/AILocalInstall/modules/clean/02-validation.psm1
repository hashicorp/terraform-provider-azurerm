# Validation Module - Clean Architecture  
# Repository and installation validation using real patterns

function Test-RepositoryStructure {
    <#
    .SYNOPSIS
        Validates repository structure (real implementation pattern)
    #>
    param([string]$RepositoryPath)
    
    Write-StatusMessage "Validating repository structure..." "Info"
    
    $required = @(
        ".github\instructions",
        ".github\copilot-instructions.md",
        "internal"
    )
    
    $missing = @()
    foreach ($path in $required) {
        $fullPath = Join-Path $RepositoryPath $path
        if (-not (Test-Path $fullPath)) {
            $missing += $path
        }
    }
    
    $isValid = ($missing.Count -eq 0)
    
    return @{
        IsValidRepository = $isValid
        MissingPaths = $missing
        RepositoryPath = $RepositoryPath
    }
}

function Test-CurrentInstallation {
    <#
    .SYNOPSIS
        Tests current installation status (real pattern from working system)
    #>
    param([string]$RepositoryPath)
    
    Write-StatusMessage "Checking current installation status..." "Info"
    
    # Use real patterns - check files in repository and VS Code settings
    $instructionCheck = Install-InstructionFiles -RepositoryPath $RepositoryPath
    $promptCheck = Install-PromptFiles -RepositoryPath $RepositoryPath
    $mainCheck = Install-MainFiles -RepositoryPath $RepositoryPath
    
    # Check VS Code settings (real pattern)
    $vsCodeUserPath = Join-Path $env:APPDATA "Code\User"
    $settingsPath = Join-Path $vsCodeUserPath "settings.json"
    $vsCodeConfigured = $false
    
    if (Test-Path $settingsPath) {
        try {
            $settings = Get-Content $settingsPath -Raw | ConvertFrom-Json
            $vsCodeConfigured = $null -ne $settings.terraform_azurerm_provider_mode
        } catch {
            $vsCodeConfigured = $false
        }
    }
    
    $totalSuccess = $instructionCheck.Success -and $promptCheck.Success -and $mainCheck.Success -and $vsCodeConfigured
    
    return @{
        Success = $totalSuccess
        InstructionFiles = $instructionCheck
        PromptFiles = $promptCheck
        MainFiles = $mainCheck
        VSCodeSettings = $vsCodeConfigured
        Errors = @()
    }
}

function Test-Prerequisites {
    <#
    .SYNOPSIS
        Checks if system prerequisites are met for AI installation
    #>
    
    # Check PowerShell version (5.1+ required) 
    if ($PSVersionTable.PSVersion.Major -lt 5) {
        Write-StatusMessage "PowerShell 5.1 or later is required, found version $($PSVersionTable.PSVersion)" "Error"
        return $false
    }
    
    # Use existing config infrastructure to get VS Code path
    try {
        $config = Get-InstallationConfig -RepositoryPath (Get-Location)
        $vsCodeUserPath = $config.VSCodeUserPath
        
        # Test if we can access/create VS Code User directory
        if (-not (Test-Path $vsCodeUserPath)) {
            New-Item -Path $vsCodeUserPath -ItemType Directory -Force | Out-Null
        }
        
        # Test write access with minimal file
        $testFile = Join-Path $vsCodeUserPath ".install-test"
        "test" | Out-File -FilePath $testFile -Encoding UTF8
        Remove-Item $testFile -Force
        
        Write-StatusMessage "Prerequisites check passed" "Success"
        return $true
        
    } catch {
        Write-StatusMessage "Cannot access VS Code User directory: $($_.Exception.Message)" "Error"
        return $false
    }
}

Export-ModuleMember -Function Test-RepositoryStructure, Test-CurrentInstallation, Test-Prerequisites
