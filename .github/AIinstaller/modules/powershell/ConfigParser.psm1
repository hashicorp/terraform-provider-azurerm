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

function Test-WorkspaceValid {
    <#
    .SYNOPSIS
    Validate that the current directory is a valid Terraform AzureRM Provider workspace
    #>
    
    $requiredFiles = @(
        "go.mod",
        "main.go",
        "internal",
        "website"
    )
    
    $requiredContent = @{
        "go.mod" = "terraform-provider-azurerm"
        "main.go" = "github.com/hashicorp/terraform-provider-azurerm"
    }
    
    # Check required files exist
    foreach ($file in $requiredFiles) {
        if (-not (Test-Path $file)) {
            return @{
                Valid = $false
                Reason = "Missing required file or directory: $file"
            }
        }
    }
    
    # Check required content
    foreach ($file in $requiredContent.Keys) {
        if (Test-Path $file) {
            $content = Get-Content $file -Raw -ErrorAction SilentlyContinue
            if ($content -notmatch [regex]::Escape($requiredContent[$file])) {
                return @{
                    Valid = $false
                    Reason = "File $file does not contain expected content: $($requiredContent[$file])"
                }
            }
        }
    }
    
    return @{
        Valid = $true
        Reason = "Valid Terraform AzureRM Provider workspace"
    }
}

function Test-GitRepository {
    <#
    .SYNOPSIS
    Validate git repository state and branch
    #>
    
    # Check if git repository
    if (-not (Test-Path ".git")) {
        return @{
            Valid = $false
            Reason = "Not a git repository"
        }
    }
    
    # Get current branch
    try {
        $currentBranch = git rev-parse --abbrev-ref HEAD 2>$null
        if ($LASTEXITCODE -ne 0) {
            return @{
                Valid = $false
                Reason = "Unable to determine current git branch"
            }
        }
        
        # Check if on source branch (would cause conflicts)
        if ($currentBranch -eq "exp/terraform_copilot") {
            return @{
                Valid = $false
                Reason = "Cannot run installer from source branch 'exp/terraform_copilot'. Switch to a different branch first."
            }
        }
        
        return @{
            Valid = $true
            CurrentBranch = $currentBranch
            Reason = "Valid git repository on branch '$currentBranch'"
        }
    }
    catch {
        return @{
            Valid = $false
            Reason = "Error checking git status: $($_.Exception.Message)"
        }
    }
}

function Get-WorkspaceStatus {
    <#
    .SYNOPSIS
    Get comprehensive workspace status including installed files
    #>
    
    $config = Get-InstallationConfig
    $status = @{
        Workspace = Test-WorkspaceValid
        Git = Test-GitRepository
        InstalledFiles = @{}
        MissingFiles = @()
        OutdatedFiles = @()
        TotalFiles = $config.Files.Count
        InstalledCount = 0
    }
    
    # Check each file
    foreach ($filePath in $config.Files.Keys) {
        $fileInfo = $config.Files[$filePath]
        $fileStatus = @{
            Path = $filePath
            Exists = Test-Path $filePath
            Required = $fileInfo.Required
            Type = $fileInfo.Type
            Description = $fileInfo.Description
            Size = 0
            LastModified = $null
        }
        
        if ($fileStatus.Exists) {
            $fileDetails = Get-Item $filePath -ErrorAction SilentlyContinue
            if ($fileDetails) {
                $fileStatus.Size = $fileDetails.Length
                $fileStatus.LastModified = $fileDetails.LastWriteTime
            }
            $status.InstalledCount++
        } else {
            if ($fileInfo.Required) {
                $status.MissingFiles += $filePath
            }
        }
        
        $status.InstalledFiles[$filePath] = $fileStatus
    }
    
    return $status
}

function Test-InternetConnectivity {
    <#
    .SYNOPSIS
    Test internet connectivity and GitHub access
    #>
    
    try {
        # Test GitHub connectivity
        $testUrl = "https://api.github.com"
        $response = Invoke-WebRequest -Uri $testUrl -Method HEAD -TimeoutSec 10 -UseBasicParsing
        
        if ($response.StatusCode -eq 200) {
            return @{
                Connected = $true
                Reason = "Internet connectivity confirmed"
            }
        } else {
            return @{
                Connected = $false
                Reason = "GitHub API returned status code: $($response.StatusCode)"
            }
        }
    }
    catch {
        return @{
            Connected = $false
            Reason = "Unable to connect to GitHub: $($_.Exception.Message)"
        }
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

function Test-FileUpToDate {
    <#
    .SYNOPSIS
    Check if a local file is up to date with the remote version
    #>
    param(
        [Parameter(Mandatory)]
        [string]$FilePath,
        
        [string]$Branch = "exp/terraform_copilot"
    )
    
    if (-not (Test-Path $FilePath)) {
        return @{
            UpToDate = $false
            Reason = "File does not exist locally"
        }
    }
    
    try {
        $downloadUrl = Get-FileDownloadUrl -FilePath $FilePath -Branch $Branch
        if (-not $downloadUrl) {
            return @{
                UpToDate = $false
                Reason = "Unable to determine download URL"
            }
        }
        
        # Get remote file info
        $remoteResponse = Invoke-WebRequest -Uri $downloadUrl -Method HEAD -UseBasicParsing -ErrorAction Stop
        $remoteLastModified = $remoteResponse.Headers['Last-Modified']
        
        # Get local file info
        $localFile = Get-Item $FilePath
        
        if ($remoteLastModified) {
            $remoteDate = [DateTime]::Parse($remoteLastModified)
            if ($localFile.LastWriteTime -lt $remoteDate) {
                return @{
                    UpToDate = $false
                    Reason = "Local file is older than remote version"
                    LocalDate = $localFile.LastWriteTime
                    RemoteDate = $remoteDate
                }
            }
        }
        
        return @{
            UpToDate = $true
            Reason = "File appears to be up to date"
        }
    }
    catch {
        return @{
            UpToDate = $false
            Reason = "Unable to check remote file: $($_.Exception.Message)"
        }
    }
}

#endregion

#region Export Module Members

Export-ModuleMember -Function @(
    'Get-ManifestConfig',
    'Get-InstallationConfig',
    'Test-WorkspaceValid',
    'Test-GitRepository',
    'Get-WorkspaceStatus',
    'Test-InternetConnectivity',
    'Get-FileDownloadUrl',
    'ConvertTo-RelativePath',
    'Test-FileUpToDate'
)

#endregion
