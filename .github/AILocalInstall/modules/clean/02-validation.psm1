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
    
    # Use Test-AIInstallation to check existing installation (correct pattern)
    $installationCheck = Test-AIInstallation -RepositoryPath $RepositoryPath
    
    return @{
        Success = $installationCheck.Success
        InstructionFiles = $installationCheck.InstructionFiles
        PromptFiles = $installationCheck.PromptFiles
        MainFiles = $installationCheck.MainFiles
        ExpectedInstructionFiles = $installationCheck.ExpectedInstructionFiles
        ExpectedPromptFiles = $installationCheck.ExpectedPromptFiles
        ExpectedMainFiles = $installationCheck.ExpectedMainFiles
        SettingsConfigured = $installationCheck.SettingsConfigured
        Errors = $installationCheck.Errors
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
