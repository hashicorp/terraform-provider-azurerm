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

function Get-InstallationConfig {
    <#
    .SYNOPSIS
    Get the installation configuration with all file mappings and settings
    #>
    param(
        [string]$Branch = "exp/terraform_copilot"
    )
    
    return @{
        SourceBranch = $Branch
        BaseUrl = "https://raw.githubusercontent.com/hashicorp/terraform-provider-azurerm/$Branch"
        
        Files = @{
            # Main Copilot instructions
            ".github/copilot-instructions.md" = @{
                Url = "/.github/copilot-instructions.md"
                Required = $true
                Type = "Instructions"
                Description = "Main GitHub Copilot instructions file"
            }
            
            # Detailed instruction files
            ".github/instructions/implementation-guide.instructions.md" = @{
                Url = "/.github/instructions/implementation-guide.instructions.md"
                Required = $true
                Type = "Instructions"
                Description = "Complete implementation guide"
            }
            
            ".github/instructions/documentation-guidelines.instructions.md" = @{
                Url = "/.github/instructions/documentation-guidelines.instructions.md"
                Required = $true
                Type = "Instructions"
                Description = "Documentation standards and guidelines"
            }
            
            ".github/instructions/testing-guidelines.instructions.md" = @{
                Url = "/.github/instructions/testing-guidelines.instructions.md"
                Required = $true
                Type = "Instructions"
                Description = "Testing guidelines and patterns"
            }
            
            ".github/instructions/provider-guidelines.instructions.md" = @{
                Url = "/.github/instructions/provider-guidelines.instructions.md"
                Required = $true
                Type = "Instructions"
                Description = "Azure provider guidelines"
            }
            
            ".github/instructions/code-clarity-enforcement.instructions.md" = @{
                Url = "/.github/instructions/code-clarity-enforcement.instructions.md"
                Required = $true
                Type = "Instructions"
                Description = "Code clarity and enforcement guidelines"
            }
            
            ".github/instructions/azure-patterns.instructions.md" = @{
                Url = "/.github/instructions/azure-patterns.instructions.md"
                Required = $true
                Type = "Instructions"
                Description = "Azure-specific implementation patterns"
            }
            
            ".github/instructions/error-patterns.instructions.md" = @{
                Url = "/.github/instructions/error-patterns.instructions.md"
                Required = $true
                Type = "Instructions"
                Description = "Error handling patterns"
            }
            
            ".github/instructions/migration-guide.instructions.md" = @{
                Url = "/.github/instructions/migration-guide.instructions.md"
                Required = $true
                Type = "Instructions"
                Description = "Migration patterns and procedures"
            }
            
            ".github/instructions/schema-patterns.instructions.md" = @{
                Url = "/.github/instructions/schema-patterns.instructions.md"
                Required = $true
                Type = "Instructions"
                Description = "Schema design patterns"
            }
            
            ".github/instructions/performance-optimization.instructions.md" = @{
                Url = "/.github/instructions/performance-optimization.instructions.md"
                Required = $true
                Type = "Instructions"
                Description = "Performance optimization guidelines"
            }
            
            ".github/instructions/security-compliance.instructions.md" = @{
                Url = "/.github/instructions/security-compliance.instructions.md"
                Required = $true
                Type = "Instructions"
                Description = "Security and compliance patterns"
            }
            
            ".github/instructions/troubleshooting-decision-trees.instructions.md" = @{
                Url = "/.github/instructions/troubleshooting-decision-trees.instructions.md"
                Required = $true
                Type = "Instructions"
                Description = "Troubleshooting decision trees"
            }
            
            ".github/instructions/api-evolution-patterns.instructions.md" = @{
                Url = "/.github/instructions/api-evolution-patterns.instructions.md"
                Required = $true
                Type = "Instructions"
                Description = "API evolution and versioning patterns"
            }
            
            # Prompt files
            ".github/prompts/add-unit-tests.prompt.md" = @{
                Url = "/.github/prompts/add-unit-tests.prompt.md"
                Required = $true
                Type = "Prompts"
                Description = "Prompt for adding unit tests"
            }
            
            ".github/prompts/code-review-committed-changes.prompt.md" = @{
                Url = "/.github/prompts/code-review-committed-changes.prompt.md"
                Required = $true
                Type = "Prompts"
                Description = "Prompt for reviewing committed changes"
            }
            
            ".github/prompts/code-review-local-changes.prompt.md" = @{
                Url = "/.github/prompts/code-review-local-changes.prompt.md"
                Required = $true
                Type = "Prompts"
                Description = "Prompt for reviewing local changes"
            }
            
            ".github/prompts/setup-go-dev-environment.prompt.md" = @{
                Url = "/.github/prompts/setup-go-dev-environment.prompt.md"
                Required = $true
                Type = "Prompts"
                Description = "Prompt for setting up Go development environment"
            }
            
            ".github/prompts/summarize-repo-deep-dive.prompt.md" = @{
                Url = "/.github/prompts/summarize-repo-deep-dive.prompt.md"
                Required = $true
                Type = "Prompts"
                Description = "Prompt for deep repository analysis"
            }
            
            ".github/prompts/summarize-repo.prompt.md" = @{
                Url = "/.github/prompts/summarize-repo.prompt.md"
                Required = $true
                Type = "Prompts"
                Description = "Prompt for repository summary"
            }
            
            # VS Code settings
            ".vscode/settings.json" = @{
                Url = "/.vscode/settings.json"
                Required = $true
                Type = "Configuration"
                Description = "VS Code workspace settings"
            }
        }
        
        GitIgnoreEntries = @(
            "# AI Infrastructure (auto-generated by install-copilot-setup.ps1)",
            ".github/copilot-instructions.md",
            ".github/instructions/",
            ".github/prompts/",
            ".vscode/settings.json"
        )
        
        RequiredDirectories = @(
            ".github",
            ".github/instructions",
            ".github/prompts",
            ".vscode"
        )
    }
}

function Get-FileDownloadUrl {
    <#
    .SYNOPSIS
    Get the full download URL for a specific file
    #>
    param(
        [Parameter(Mandatory)]
        [string]$FilePath,
        
        [string]$Branch = "exp/terraform_copilot"
    )
    
    $config = Get-InstallationConfig -Branch $Branch
    
    if ($config.Files.ContainsKey($FilePath)) {
        $fileInfo = $config.Files[$FilePath]
        return $config.BaseUrl + $fileInfo.Url
    }
    
    return $null
}

function ConvertTo-RelativePath {
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

#endregion

#region Export Module Members

Export-ModuleMember -Function @(
    'Get-ManifestConfig',
    'Get-InstallationConfig',
    'Get-FileDownloadUrl',
    'ConvertTo-RelativePath'
)

#endregion
