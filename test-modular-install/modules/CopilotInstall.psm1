# CopilotInstall.psm1 - GitHub Copilot instruction installation for Terraform AzureRM Provider

function Test-CopilotInstallation {
    <#
    .SYNOPSIS
    Tests if GitHub Copilot instructions are installed
    
    .PARAMETER RepositoryPath
    Path to the repository to check
    
    .OUTPUTS
    Boolean indicating if Copilot instructions are installed
    #>
    
    param(
        [Parameter(Mandatory = $true)]
        [string]$RepositoryPath
    )
    
    $instructionsPath = Join-Path $RepositoryPath ".github\copilot-instructions.md"
    return Test-Path $instructionsPath
}

function Install-CopilotInstructions {
    <#
    .SYNOPSIS
    Installs GitHub Copilot instruction files
    
    .PARAMETER RepositoryPath
    Path to the repository where instructions should be installed
    
    .PARAMETER Force
    Force installation even if files already exist
    
    .OUTPUTS
    Boolean indicating success
    #>
    
    param(
        [Parameter(Mandatory = $true)]
        [string]$RepositoryPath,
        
        [switch]$Force
    )
    
    $githubDir = Join-Path $RepositoryPath ".github"
    $instructionsPath = Join-Path $githubDir "copilot-instructions.md"
    
    # Check if already exists
    if ((Test-Path $instructionsPath) -and -not $Force) {
        Write-Warning "Copilot instructions already exist. Use -Force to overwrite."
        return $false
    }
    
    # Ensure .github directory exists
    if (-not (Test-Path $githubDir)) {
        New-Item -ItemType Directory -Path $githubDir -Force | Out-Null
    }
    
    # Create the instruction file
    $instructionContent = Get-CopilotInstructionContent
    Set-Content -Path $instructionsPath -Value $instructionContent -Encoding UTF8
    
    Write-Verbose "Copilot instructions installed at: $instructionsPath"
    return $true
}

function Get-CopilotInstructionContent {
    <#
    .SYNOPSIS
    Gets the content for the Copilot instruction file
    
    .OUTPUTS
    String containing the instruction file content
    #>
    
    return @"
---
applyTo: "internal/**/*.go"
description: "This is the official Terraform Provider for Azure (Resource Manager), written in Go. It enables Terraform to manage Azure resources through the Azure Resource Manager APIs."
---
# Custom instructions

This is the official Terraform Provider for Azure (Resource Manager), written in Go. It enables Terraform to manage Azure resources through the Azure Resource Manager APIs.

## Key Guidelines

### Code Standards
- Follow Go coding standards and conventions
- Use proper error handling with context
- Implement comprehensive acceptance tests
- Follow Terraform Plugin SDK patterns
- Use Azure SDK for Go patterns

### Resource Implementation
- Use typed resource implementations for new resources
- Follow CRUD lifecycle patterns
- Implement proper state management
- Use appropriate validation functions
- Handle Azure-specific patterns correctly

### Testing Requirements
- Write comprehensive acceptance tests
- Test both success and failure scenarios
- Include import functionality tests
- Use consistent test naming patterns
- Ensure tests are idempotent

### Documentation Standards
- Include working examples
- Document all resource attributes
- Provide import documentation
- Keep documentation synchronized with code

### Azure Integration
- Use official Azure SDK for Go
- Implement proper polling for long-running operations
- Handle Azure API rate limits appropriately
- Follow Azure resource naming conventions
- Use proper authentication patterns

For detailed implementation guidance, see the instruction files in `.github/instructions/`.
"@
}

function Remove-CopilotInstructions {
    <#
    .SYNOPSIS
    Removes GitHub Copilot instruction files
    
    .PARAMETER RepositoryPath
    Path to the repository where instructions should be removed
    
    .OUTPUTS
    Boolean indicating success
    #>
    
    param(
        [Parameter(Mandatory = $true)]
        [string]$RepositoryPath
    )
    
    $instructionsPath = Join-Path $RepositoryPath ".github\copilot-instructions.md"
    
    if (Test-Path $instructionsPath) {
        Remove-Item $instructionsPath -Force
        Write-Verbose "Copilot instructions removed from: $instructionsPath"
        return $true
    }
    
    Write-Verbose "No Copilot instructions found to remove"
    return $false
}

function Get-CopilotInstallationInfo {
    <#
    .SYNOPSIS
    Gets information about the current Copilot installation
    
    .PARAMETER RepositoryPath
    Path to the repository to check
    
    .OUTPUTS
    Hashtable containing installation information
    #>
    
    param(
        [Parameter(Mandatory = $true)]
        [string]$RepositoryPath
    )
    
    $instructionsPath = Join-Path $RepositoryPath ".github\copilot-instructions.md"
    
    $info = @{
        IsInstalled = $false
        InstructionsPath = $instructionsPath
        FileSize = 0
        LastModified = $null
    }
    
    if (Test-Path $instructionsPath) {
        $fileInfo = Get-Item $instructionsPath
        $info.IsInstalled = $true
        $info.FileSize = $fileInfo.Length
        $info.LastModified = $fileInfo.LastWriteTime
    }
    
    return $info
}

# Export module members
Export-ModuleMember -Function Test-CopilotInstallation, Install-CopilotInstructions, Remove-CopilotInstructions, Get-CopilotInstallationInfo
