# ConfigParser Module for Terraform AzureRM Provider AI Setup
# Handles configuration file parsing, validation, and management
# STREAMLINED VERSION - Contains only functions actually used by main script

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
    
    # Determine if we're in source repository to get the correct branch
    # Push the workspace root to global scope temporarily for Test-SourceRepository
    $originalWorkspaceRoot = $Global:WorkspaceRoot
    $Global:WorkspaceRoot = $WorkspaceRoot
    
    try {
        $isSourceRepo = Test-SourceRepository
        $currentBranch = if ($isSourceRepo) { 
            "exp/terraform_copilot" 
        } else { 
            # For target repositories, use main branch as the source
            "main" 
        }
    }
    finally {
        # Restore original global workspace root
        $Global:WorkspaceRoot = $originalWorkspaceRoot
    }
    
    return @{
        Version = "1.0.0"
        Branch = $currentBranch
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

# Export only the functions actually used by the main script
Export-ModuleMember -Function @(
    'Get-ManifestConfig',
    'Get-InstallerConfig'
)
