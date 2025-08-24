# ConfigParser Module for Terraform AzureRM Provider AI Setup
# Handles configuration file parsing, validation, and management

#region Public Functions

function Get-ManifestConfig {
    <#
    .SYNOPSIS
    Parse the file manifest configuration and return structured data
    
    .PARAMETER ManifestPath
    Path to the manifest file. Defaults to file-manifest.config in the AIinstaller directory
    
    .PARAMETER Branch
    Git branch for remote URLs
    #>
    param(
        [string]$ManifestPath,
        [string]$Branch = "exp/terraform_copilot"
    )
    
    # Find manifest file if not specified
    if (-not $ManifestPath) {
        $scriptRoot = Split-Path (Split-Path $PSScriptRoot -Parent) -Parent
        $ManifestPath = Join-Path $scriptRoot "file-manifest.config"
    }
    
    if (-not (Test-Path $ManifestPath)) {
        throw "Manifest file not found: $ManifestPath"
    }
    
    $manifest = @{
        SourceBranch = $Branch
        BaseUrl = "https://raw.githubusercontent.com/hashicorp/terraform-provider-azurerm/$Branch"
        Sections = @{}
    }
    
    $currentSection = $null
    $content = Get-Content $ManifestPath
    
    foreach ($line in $content) {
        $line = $line.Trim()
        
        # Skip empty lines and comments
        if (-not $line -or $line.StartsWith('#')) {
            continue
        }
        
        # Check for section headers
        if ($line.StartsWith('[') -and $line.EndsWith(']')) {
            $currentSection = $line.Substring(1, $line.Length - 2)
            $manifest.Sections[$currentSection] = @()
            continue
        }
        
        # Add files to current section
        if ($currentSection -and $line) {
            $manifest.Sections[$currentSection] += $line
        }
    }
    
    return $manifest
}

function Get-FileDownloadUrl {
    <#
    .SYNOPSIS
    Get the full download URL for a specific file from the manifest
    #>
    param(
        [Parameter(Mandatory)]
        [string]$FilePath,
        
        [string]$Branch = "exp/terraform_copilot"
    )
    
    $manifestConfig = Get-ManifestConfig -Branch $Branch
    $baseUrl = $manifestConfig.BaseUrl
    
    # Check all sections for the file
    foreach ($section in $manifestConfig.Sections.Keys) {
        if ($manifestConfig.Sections[$section] -contains $FilePath) {
            return "$baseUrl/$FilePath"
        }
    }

    return $null
}

function Get-FileLocalPath {
    <#
    .SYNOPSIS
    Get the correct local path for a file based on its manifest path
    #>
    param(
        [Parameter(Mandatory)]
        [string]$FilePath,
        
        [string]$WorkspaceRoot = "."
    )
    
    # Simply join the workspace root with the file path from manifest
    return Join-Path $WorkspaceRoot $FilePath
}function ConvertTo-RelativePath {
    <#
    .SYNOPSIS
    Convert absolute path to relative path from workspace root
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Path
    )
    
    $workspaceRoot = Get-Location
    $absolutePath = $ExecutionContext.SessionState.Path.GetUnresolvedProviderPathFromPSPath($Path)
    
    try {
        $relativePath = [System.IO.Path]::GetRelativePath($workspaceRoot.Path, $absolutePath)
        return $relativePath -replace '\\', '/'
    }
    catch {
        return $Path
    }
}

function Get-InstallerConfig {
    <#
    .SYNOPSIS
    Get the complete installer configuration with all file mappings and targets
    
    .PARAMETER WorkspaceRoot
    The root directory of the workspace
    
    .PARAMETER ManifestConfig
    The manifest configuration from Get-ManifestConfig
    #>
    param(
        [Parameter(Mandatory)]
        [string]$WorkspaceRoot,
        
        [Parameter(Mandatory)]
        [hashtable]$ManifestConfig
    )
    
    return @{
        Version = "1.0.0"
        Branch = "exp/terraform_copilot"
        SourceRepository = "https://raw.githubusercontent.com/hashicorp/terraform-provider-azurerm"
        Files = @{
            Instructions = @{
                Source = ".github/copilot-instructions.md"
                Target = (Join-Path $WorkspaceRoot ".github/copilot-instructions.md")
                Description = "Main Copilot instructions for AI-powered development"
            }
            InstructionFiles = @{
                Target = (Join-Path $WorkspaceRoot ".github/instructions")
                Description = "Detailed implementation guidelines and patterns"
                Files = $ManifestConfig.Sections.INSTRUCTION_FILES
            }
            PromptFiles = @{
                Target = (Join-Path $WorkspaceRoot ".github/prompts")
                Description = "AI prompt templates for development workflows"
                Files = $ManifestConfig.Sections.PROMPT_FILES
            }
            InstallerFiles = @{
                Target = "$env:USERPROFILE\.terraform-ai-installer"
                Description = "Cross-platform installer scripts and modules (PowerShell + Bash) for bootstrap functionality"
                Files = $ManifestConfig.Sections.INSTALLER_FILES_BOOTSTRAP
            }
            UniversalFiles = @{
                Target = (Join-Path $WorkspaceRoot ".vscode")
                Description = "Platform-independent configuration files"
                Files = $ManifestConfig.Sections.UNIVERSAL_FILES
            }
        }
    }
}

#endregion

#region Export Module Members

Export-ModuleMember -Function @(
    'Get-ManifestConfig',
    'Get-InstallerConfig',
    'Get-FileDownloadUrl',
    'Get-FileLocalPath',
    'ConvertTo-RelativePath'
)

#endregion
