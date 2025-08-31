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
        # Try to find manifest file in multiple locations
        $possiblePaths = @(
            # User profile installer directory (when running from bootstrapped copy)
            (Join-Path $env:USERPROFILE ".terraform-ai-installer\file-manifest.config"),
            # Current script directory (when running from user profile)
            (Join-Path (Split-Path $PSScriptRoot -Parent) "file-manifest.config"),
            # Original repository structure (when running from source)
            (Join-Path (Split-Path (Split-Path $PSScriptRoot -Parent) -Parent) "file-manifest.config")
        )

        $ManifestPath = $null
        foreach ($path in $possiblePaths) {
            if (Test-Path $path) {
                $ManifestPath = $path
                break
            }
        }

        if (-not $ManifestPath) {
            # Fallback to old behavior
            $scriptRoot = Split-Path (Split-Path $PSScriptRoot -Parent) -Parent
            $ManifestPath = Join-Path $scriptRoot "file-manifest.config"
        }
    }

    if (-not (Test-Path $ManifestPath)) {
        throw "Manifest file not found: $ManifestPath"
    }

    $manifest = @{
        BaseUrl = "https://raw.githubusercontent.com/hashicorp/terraform-provider-azurerm/exp/terraform_copilot"
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

    # DOWNLOAD SOURCE: Always use exp/terraform_copilot branch because that's where the AI files exist
    # DOWNLOAD TARGET: Copy files to the local workspace directory (regardless of local branch)

    return @{
        Version = "1.0.0"
        Branch = "exp/terraform_copilot"
        SourceRepository = "https://raw.githubusercontent.com/hashicorp/terraform-provider-azurerm/exp/terraform_copilot"
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
