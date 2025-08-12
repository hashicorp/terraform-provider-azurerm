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

Export-ModuleMember -Function Test-RepositoryStructure, Test-CurrentInstallation
