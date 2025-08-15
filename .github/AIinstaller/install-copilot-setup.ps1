# Main AI Infrastructure Installer for Terraform AzureRM Provider
# Version: 1.0.0
# Description: Interactive installer for AI-powered development infrastructure

#requires -version 5.1

[CmdletBinding()]
param(
    [Parameter(HelpMessage = "Copy installer to user profile for feature branch use")]
    [switch]$Bootstrap,
    
    [Parameter(HelpMessage = "Path to the repository directory for git operations (when running from user profile)")]
    [string]$RepoDirectory,
    
    [Parameter(HelpMessage = "Overwrite existing files without prompting")]
    [switch]${Auto-Approve},
    
    [Parameter(HelpMessage = "Show what would be done without making changes")]
    [switch]${Dry-Run},
    
    [Parameter(HelpMessage = "Check the current state of the workspace")]
    [switch]$Verify,
    
    [Parameter(HelpMessage = "Remove AI infrastructure from the workspace")]
    [switch]$Clean,
    
    [Parameter(HelpMessage = "Show detailed help information")]
    [switch]$Help
)

# Import required modules
$ModulesPath = Join-Path $PSScriptRoot "modules\powershell"
$RequiredModules = @(
    "ConfigParser",
    "FileOperations",
    "ValidationEngine", 
    "UI"
)

foreach ($module in $RequiredModules) {
    $modulePath = Join-Path $ModulesPath "$module.psm1"
    if (Test-Path $modulePath) {
        Import-Module $modulePath -Force
    } else {
        Write-Error "Required module '$module' not found at: $modulePath"
        exit 1
    }
}

# Helper function to get workspace root
function Get-WorkspaceRoot {
    <#
    .SYNOPSIS
    Dynamically determines the workspace root directory
    
    .DESCRIPTION
    Finds the root of the terraform-provider-azurerm workspace by looking for key files
    or navigating up from the current script location
    #>
    
    # Start from script directory and look for workspace indicators
    $currentPath = $PSScriptRoot
    
    while ($currentPath -and $currentPath -ne (Split-Path $currentPath -Parent)) {
        # Look for terraform-provider-azurerm indicators
        if ((Test-Path (Join-Path $currentPath "go.mod")) -and 
            (Test-Path (Join-Path $currentPath "main.go")) -and
            (Test-Path (Join-Path $currentPath "internal"))) {
            return $currentPath
        }
        
        # Move up one directory
        $currentPath = Split-Path $currentPath -Parent
    }
    
    # Fallback: assume we're in .github/AIinstaller and go up two levels
    return Split-Path (Split-Path $PSScriptRoot -Parent) -Parent
}

# Global variables
# Initialize workspace root - will be updated if RepoDirectory parameter is provided
$Global:WorkspaceRoot = Get-WorkspaceRoot

# Override workspace root if RepoDirectory parameter is provided
if ($RepoDirectory) {
    if (Test-Path $RepoDirectory) {
        # Validate that this looks like a terraform-provider-azurerm repository
        $goModPath = Join-Path $RepoDirectory "go.mod"
        $mainGoPath = Join-Path $RepoDirectory "main.go"
        $internalPath = Join-Path $RepoDirectory "internal"
        
        if ((Test-Path $goModPath) -and (Test-Path $mainGoPath) -and (Test-Path $internalPath)) {
            $Global:WorkspaceRoot = $RepoDirectory
            Show-RepositoryInfo -Directory $RepoDirectory
        } else {
            Write-Error "INVALID REPOSITORY: The specified directory does not appear to be a terraform-provider-azurerm repository."
            Show-ErrorBlock -Issue "The -RepoDirectory parameter must point to a valid terraform-provider-azurerm repository root." -Solutions @(
                "Ensure you're pointing to the repository ROOT directory",
                "Verify the directory contains terraform-provider-azurerm source code", 
                "Example: -RepoDirectory 'C:\github.com\hashicorp\terraform-provider-azurerm'"
            )
            exit 1
        }
    } else {
        Write-Error "DIRECTORY NOT FOUND: The specified RepoDirectory does not exist."
        Show-ErrorBlock -Issue "The path '$RepoDirectory' could not be found on this system." -Solutions @(
            "Check the path spelling and ensure it exists",
            "Use an absolute path (e.g., 'C:\path\to\repo')",
            "Ensure you have permissions to access the directory"
        ) -ExampleUsage ".\install-copilot-setup.ps1 -RepoDirectory 'C:\github.com\hashicorp\terraform-provider-azurerm'"
        exit 1
    }
}

# Get dynamic configuration from manifest FIRST
$ManifestConfig = Get-ManifestConfig

# Initialize installer configuration after workspace root is finalized
$Global:InstallerConfig = @{
    Version = "1.0.0"
    Branch = "exp/terraform_copilot"
    SourceRepository = "https://raw.githubusercontent.com/hashicorp/terraform-provider-azurerm"
    Files = @{
        Instructions = @{
            Source = ".github/copilot-instructions.md"
            Target = (Join-Path $Global:WorkspaceRoot ".github/copilot-instructions.md")
            Description = "Main Copilot instructions for AI-powered development"
        }
        InstructionFiles = @{
            Source = ".github/instructions"
            Target = (Join-Path $Global:WorkspaceRoot ".github/instructions")
            Description = "Detailed implementation guidelines and patterns"
            Files = $ManifestConfig.Sections.INSTRUCTION_FILES
        }
        PromptFiles = @{
            Source = ".github/prompts"
            Target = (Join-Path $Global:WorkspaceRoot ".github/prompts")
            Description = "AI prompt templates for development workflows"
            Files = $ManifestConfig.Sections.PROMPT_FILES
        }
        InstallerFiles = @{
            Source = ".github/AIinstaller"
            Target = "$env:USERPROFILE\.terraform-ai-installer"
            Description = "Installer scripts and modules for bootstrap functionality"
            Files = $ManifestConfig.Sections.INSTALLER_FILES
        }
    }
}

# Note: InstallerConfig will be initialized after RepoDirectory parameter processing

function Test-BranchType {
    <#
    .SYNOPSIS
    Determine if running from source branch or feature branch
    #>
    
    try {
        # Use RepoDirectory parameter if provided, otherwise use current workspace root
        $workspaceRoot = if ($RepoDirectory) { 
            $RepoDirectory 
        } else { 
            $Global:WorkspaceRoot 
        }
        
        if (-not $workspaceRoot -or -not (Test-Path $workspaceRoot)) {
            Write-Warning "Repository directory not found: $workspaceRoot"
            return "unknown"
        }
        
        Push-Location $workspaceRoot
        try {
            $currentBranch = git rev-parse --abbrev-ref HEAD 2>$null
            if ($LASTEXITCODE -ne 0) {
                Write-Warning "Could not determine current git branch"
                return "unknown"
            }
            
            if ($currentBranch -eq $Global:InstallerConfig.Branch) {
                return "source"
            } else {
                return "feature"
            }
        }
        finally {
            Pop-Location
        }
    }
    catch {
        Write-Warning "Git not available or error checking git repository"
        return "unknown"
    }
}

function Invoke-Bootstrap {
    <#
    .SYNOPSIS
    Copy installer files to user profile for feature branch use
    #>
    
    try {
        Write-Section "Bootstrap - Copying Installer to User Profile"
        
        # Create target directory
        $targetDirectory = Join-Path $env:USERPROFILE ".terraform-ai-installer"
        if (-not (Test-Path $targetDirectory)) {
            New-Item -ItemType Directory -Path $targetDirectory -Force | Out-Null
            Show-DirectoryOperation -Directory $targetDirectory -Status "Created"
        } else {
            Show-DirectoryOperation -Directory $targetDirectory -Status "Existing"
        }
        
        # Files to bootstrap from configuration
        $filesToBootstrap = $Global:InstallerConfig.Files.InstallerFiles.Files
        
        # Statistics
        $statistics = @{
            "Files Copied" = 0
            "Files Downloaded" = 0
            "Files Failed" = 0
            "Total Size" = 0
        }
        
        # Determine if we should copy locally or download from remote
        $isSourceRepo = Test-SourceRepository
        $aiInstallerSourcePath = Join-Path $Global:WorkspaceRoot ".github/AIinstaller"
        
        if ($isSourceRepo -and (Test-Path $aiInstallerSourcePath)) {
            Write-OperationStatus -Message "  Copying installer files from local source repository..." -Type "Progress"
            Write-Host ""
            
            # Copy files locally from source repository
            foreach ($file in $filesToBootstrap) {
                try {
                    $sourcePath = Join-Path $aiInstallerSourcePath $file
                    $fileName = Split-Path $file -Leaf
                    
                    # Create subdirectory if needed
                    $fileDir = Split-Path $file -Parent
                    if ($fileDir) {
                        $targetSubDir = Join-Path $targetDirectory $fileDir
                        if (-not (Test-Path $targetSubDir)) {
                            New-Item -ItemType Directory -Path $targetSubDir -Force | Out-Null
                        }
                        $targetPath = Join-Path $targetSubDir $fileName
                    } else {
                        $targetPath = Join-Path $targetDirectory $fileName
                    }
                    
                    Show-FileOperation -Operation "Copying" -FileName $fileName -NoNewLine
                    
                    if (Test-Path $sourcePath) {
                        Copy-Item $sourcePath $targetPath -Force
                        
                        if (Test-Path $targetPath) {
                            $fileSize = (Get-Item $targetPath).Length
                            $statistics["Files Copied"]++
                            $statistics["Total Size"] += $fileSize
                            
                            Write-Host " [OK]" -ForegroundColor "Green"
                        } else {
                            Write-Host " [FAILED]" -ForegroundColor "Red"
                            $statistics["Files Failed"]++
                        }
                    } else {
                        Write-Host " [SOURCE NOT FOUND]" -ForegroundColor "Red"
                        $statistics["Files Failed"]++
                    }
                }
                catch {
                    Write-Host " [ERROR] ($($_.Exception.Message))" -ForegroundColor "Red"
                    $statistics["Files Failed"]++
                }
            }
        } else {
            Write-Host "  Downloading installer files from remote source branch..." -ForegroundColor "Cyan"
            Write-Host ""
            
            # Download files from remote repository
            $baseUri = "$($Global:InstallerConfig.SourceRepository)/$($Global:InstallerConfig.Branch)/.github/AIinstaller"
            
            foreach ($file in $filesToBootstrap) {
                try {
                    $uri = "$baseUri/$file"
                    $fileName = Split-Path $file -Leaf
                    
                    # Create subdirectory if needed
                    $fileDir = Split-Path $file -Parent
                    if ($fileDir) {
                        $targetSubDir = Join-Path $targetDirectory $fileDir
                        if (-not (Test-Path $targetSubDir)) {
                            New-Item -ItemType Directory -Path $targetSubDir -Force | Out-Null
                        }
                        $targetPath = Join-Path $targetSubDir $fileName
                    } else {
                        $targetPath = Join-Path $targetDirectory $fileName
                    }
                    
                    Write-Host "    Downloading: $fileName" -ForegroundColor "Gray" -NoNewline
                    
                    # Download with progress
                    Invoke-WebRequest -Uri $uri -OutFile $targetPath -UseBasicParsing | Out-Null
                    
                    if (Test-Path $targetPath) {
                        $fileSize = (Get-Item $targetPath).Length
                        $statistics["Files Downloaded"]++
                        $statistics["Total Size"] += $fileSize
                        
                        Write-Host " [OK]" -ForegroundColor "Green"
                    } else {
                        Write-Host " [FAILED]" -ForegroundColor "Red"
                        $statistics["Files Failed"]++
                    }
                }
                catch {
                    Write-Host " [ERROR] ($($_.Exception.Message))" -ForegroundColor "Red"
                    $statistics["Files Failed"]++
                }
            }
        }
        
        Write-Host ""
        
        if ($statistics["Files Failed"] -eq 0) {
            $totalSizeKB = [math]::Round($statistics["Total Size"] / 1KB, 1)
            
            Write-Success "Bootstrap completed successfully!"
            Write-Host ""
            
            if ($statistics["Files Copied"] -gt 0) {
                Write-Host "  Files copied: $($statistics["Files Copied"])" -ForegroundColor "Green"
            }
            if ($statistics["Files Downloaded"] -gt 0) {
                Write-Host "  Files downloaded: $($statistics["Files Downloaded"])" -ForegroundColor "Green"
            }
            
            Write-Host "  Total size: $totalSizeKB KB" -ForegroundColor "Green"
            Write-Host "  Location: $targetDirectory" -ForegroundColor "Green"
            Write-Host ""
            Write-Host "NEXT STEPS:" -ForegroundColor "Cyan"
            Write-Host "  1. Switch to your feature branch:" -ForegroundColor "White"
            Write-Host "     git checkout feature/your-branch-name" -ForegroundColor "Gray"
            Write-Host ""
            Write-Host "  2. Run the installer from your user profile:" -ForegroundColor "White"
            Write-Host "     & `"$targetDirectory\install-copilot-setup.ps1`" -RepoDirectory `"$($Global:WorkspaceRoot)`"" -ForegroundColor "Gray"
            Write-Host ""
            Write-Host "  Note: The -RepoDirectory parameter tells the installer where to find the git repository" -ForegroundColor "Yellow"
            Write-Host "        for branch detection when running from your user profile." -ForegroundColor "Yellow"
            Write-Host ""
            
            return @{
                Success = $true
                TargetDirectory = $targetDirectory
                Statistics = $statistics
            }
        } else {
            Write-Error "Bootstrap failed: $($statistics["Files Failed"]) files could not be processed"
            return @{
                Success = $false
                Statistics = $statistics
            }
        }
    }
    catch {
        Write-Error "Bootstrap failed: $($_.Exception.Message)"
        return @{
            Success = $false
            Error = $_.Exception.Message
        }
    }
}

function Test-SourceRepository {
    <#
    .SYNOPSIS
    Determines if we're running on the source repository vs a target repository
    
    .DESCRIPTION
    Checks various indicators to determine if this is the source repository where
    AI infrastructure files are maintained vs a target repository where they
    would be installed.
    
    .OUTPUTS
    Boolean - True if this is the source repository, False if target
    #>
    
    # Check if we're on the exp/terraform_copilot branch (source branch)
    try {
        Push-Location $Global:WorkspaceRoot
        $currentBranch = git rev-parse --abbrev-ref HEAD 2>$null
        if ($currentBranch -eq "exp/terraform_copilot") {
            return $true
        }
    } catch {
        # Git not available or not in a git repo
    } finally {
        Pop-Location
    }
    
    # Check if AIinstaller directory exists (only in source)
    $aiInstallerPath = Join-Path $Global:WorkspaceRoot ".github/AIinstaller"
    if (Test-Path $aiInstallerPath) {
        return $true
    }
    
    # Check if this directory structure looks like the source
    $copilotInstructionsPath = Join-Path $Global:WorkspaceRoot ".github/copilot-instructions.md"
    $instructionsPath = Join-Path $Global:WorkspaceRoot ".github/instructions"
    $promptsPath = Join-Path $Global:WorkspaceRoot ".github/prompts"
    
    if ((Test-Path $copilotInstructionsPath) -and 
        (Test-Path $instructionsPath) -and 
        (Test-Path $promptsPath) -and
        (Test-Path $aiInstallerPath)) {
        return $true
    }
    
    return $false
}

function Invoke-VerifyWorkspace {
    <#
    .SYNOPSIS
    Verifies the presence of AI infrastructure files in the workspace
    
    .DESCRIPTION
    Checks for all required AI infrastructure files including:
    - Main copilot instructions
    - Detailed instruction files
    - Prompts directory
    - VS Code settings
    
    .OUTPUTS
    Returns verification results and displays status to console
    #>
    
    Write-Section "Workspace Verification"
    
    # Use the dynamically determined workspace root
    $workspaceRoot = $Global:WorkspaceRoot
    Push-Location $workspaceRoot
    
    try {
        # Check if we're on the source branch/repository
        $isSourceRepo = Test-SourceRepository
        
        $results = @{
            Success = $true
            Files = @()
            Issues = @()
            IsSourceRepo = $isSourceRepo
        }
        
        Write-Host "Checking workspace: $workspaceRoot" -ForegroundColor Gray
        
        if ($isSourceRepo) {
            Write-Host "NOTE: Running on SOURCE repository - AI infrastructure files are part of this repository" -ForegroundColor Cyan
        } else {
            Write-Host "NOTE: Running on TARGET repository - checking if AI infrastructure files have been installed" -ForegroundColor Yellow
        }
        Write-Host ""
        
        # Check main instructions file
        $instructionsFile = $Global:InstallerConfig.Files.Instructions.Target
        if (Test-Path $instructionsFile) {
            $results.Files += @{
                Path = $instructionsFile
                Status = "Present"
                Description = "Main Copilot instructions"
            }
            Write-Host "  [FOUND] $(Resolve-Path $instructionsFile -Relative)" -ForegroundColor Green
        } else {
            $results.Files += @{
                Path = $instructionsFile
                Status = "Missing"
                Description = "Main Copilot instructions"
            }
            $results.Issues += "Missing: $instructionsFile"
            Write-Host "  [MISSING] $(Resolve-Path $instructionsFile -Relative -ErrorAction SilentlyContinue)" -ForegroundColor Red
        }
        
        # Check instructions directory
        $instructionsDir = $Global:InstallerConfig.Files.InstructionFiles.Target
        if (Test-Path $instructionsDir -PathType Container) {
            $results.Files += @{
                Path = $instructionsDir
                Status = "Present"
                Description = "Instructions directory"
            }
            Write-Host "  [FOUND] $(Resolve-Path $instructionsDir -Relative)/" -ForegroundColor Green
            
            # Check specific instruction files
            $requiredFiles = $Global:InstallerConfig.Files.InstructionFiles.Files
            
            foreach ($file in $requiredFiles) {
                $filePath = Join-Path $instructionsDir $file
                if (Test-Path $filePath) {
                    Write-Host "    [FOUND] $file" -ForegroundColor Green
                } else {
                    Write-Host "    [MISSING] $file" -ForegroundColor Red
                    $results.Issues += "Missing: $filePath"
                }
            }
        } else {
            $results.Files += @{
                Path = $instructionsDir
                Status = "Missing"
                Description = "Instructions directory"
            }
            $results.Issues += "Missing: $instructionsDir"
            Write-Host "  [MISSING] $(Resolve-Path $instructionsDir -Relative -ErrorAction SilentlyContinue)/" -ForegroundColor Red
        }
        
        # Check prompts directory
        $promptsDir = $Global:InstallerConfig.Files.PromptFiles.Target
        if (Test-Path $promptsDir -PathType Container) {
            $results.Files += @{
                Path = $promptsDir
                Status = "Present"
                Description = "Prompts directory"
            }
            Write-Host "  [FOUND] $(Resolve-Path $promptsDir -Relative)/" -ForegroundColor Green
            
            # Check specific prompt files
            $requiredPrompts = $Global:InstallerConfig.Files.PromptFiles.Files
            
            foreach ($file in $requiredPrompts) {
                $filePath = Join-Path $promptsDir $file
                if (Test-Path $filePath) {
                    Write-Host "    [FOUND] $file" -ForegroundColor Green
                } else {
                    Write-Host "    [MISSING] $file" -ForegroundColor Red
                    $results.Issues += "Missing: $filePath"
                }
            }
        } else {
            $results.Files += @{
                Path = $promptsDir
                Status = "Missing"
                Description = "Prompts directory"
            }
            $results.Issues += "Missing: $promptsDir"
            Write-Host "  [MISSING] $(Resolve-Path $promptsDir -Relative -ErrorAction SilentlyContinue)/" -ForegroundColor Red
        }
        
        # Check .vscode directory and settings
        $vscodeDir = Join-Path $workspaceRoot ".vscode"
        if (Test-Path $vscodeDir -PathType Container) {
            Write-Host "  [FOUND] .vscode/" -ForegroundColor Green
            
            $settingsFile = Join-Path $vscodeDir "settings.json"
            if (Test-Path $settingsFile) {
                Write-Host "    [FOUND] settings.json" -ForegroundColor Green
            } else {
                Write-Host "    [MISSING] settings.json" -ForegroundColor Red
                $results.Issues += "Missing: $settingsFile"
            }
        } else {
            Write-Host "  [MISSING] .vscode/" -ForegroundColor Red
            $results.Issues += "Missing: $vscodeDir"
        }
        
        Write-Host ""
        
        # Check for deprecated files (files that exist but are no longer in manifest)
        if (-not $isSourceRepo) {
            $deprecatedFiles = Remove-DeprecatedFiles -ManifestConfig $ManifestConfig -WorkspaceRoot $workspaceRoot -DryRun $true -Quiet $true
            if ($deprecatedFiles.Count -gt 0) {
                $results.Issues += "Found $($deprecatedFiles.Count) deprecated files"
                Write-Host ""
                Write-Host "TIP: Deprecated files will be automatically removed during installation" -ForegroundColor Cyan
            }
        }
        
        if ($results.Issues.Count -gt 0) {
            $results.Success = $false
            Write-Host "Issues found:" -ForegroundColor Yellow
            foreach ($issue in $results.Issues) {
                Write-Host "  - $issue" -ForegroundColor Red
            }
            
            if (-not $isSourceRepo) {
                Write-Host ""
                Write-Host "TIP: To install missing files, run: .\install-copilot-setup.ps1" -ForegroundColor Cyan
            }
        } else {
            if ($isSourceRepo) {
                Write-Host "All AI infrastructure files are present in the source repository!" -ForegroundColor Green
            } else {
                Write-Host "All AI infrastructure files have been successfully installed!" -ForegroundColor Green
            }
        }
        
        return $results
    }
    finally {
        Pop-Location
    }
}

function Invoke-CleanWorkspace {
    param([bool]$AutoApprove, [bool]$DryRun)
    
    Write-Section "Clean Workspace"
    
    if ($DryRun) {
        Write-Warning "DRY RUN - No files will be deleted"
        Write-Host ""
    }
    
    $filesToRemove = @()
    
    # Add main instructions file
    if (Test-Path $Global:InstallerConfig.Files.Instructions.Target) {
        $filesToRemove += $Global:InstallerConfig.Files.Instructions.Target
    }
    
    # Add instruction files directory and its contents
    $instructionsDir = $Global:InstallerConfig.Files.InstructionFiles.Target
    if (Test-Path $instructionsDir) {
        $filesToRemove += $instructionsDir
    }
    
    # Add prompt files directory and its contents
    $promptsDir = $Global:InstallerConfig.Files.PromptFiles.Target
    if (Test-Path $promptsDir) {
        $filesToRemove += $promptsDir
    }
    
    # Add VS Code settings (if they were installed by this tool)
    # Note: We don't remove .vscode/settings.json as it may contain user settings
    
    if ($filesToRemove.Count -eq 0) {
        Write-Host "No AI infrastructure files found to remove." -ForegroundColor Yellow
        return @{ Success = $true }
    }
    
    Write-Host "Files to remove:" -ForegroundColor White
    foreach ($file in $filesToRemove) {
        Write-Host "  - $file" -ForegroundColor Gray
    }
    Write-Host ""
    
    if (-not $AutoApprove -and -not $DryRun) {
        $confirm = Read-Host "Remove these files? (y/N)"
        if ($confirm -ne 'y' -and $confirm -ne 'Y') {
            Write-Host "Operation cancelled." -ForegroundColor Yellow
            return @{ Success = $false }
        }
    }
    
    if (-not $DryRun) {
        foreach ($file in $filesToRemove) {
            try {
                Remove-Item $file -Recurse -Force
                Write-Host "  Removed: $file" -ForegroundColor Green
            }
            catch {
                Write-Host "  Failed to remove: $file - $($_.Exception.Message)" -ForegroundColor Red
            }
        }
    }
    
    return @{ Success = $true }
}

function Invoke-InstallInfrastructure {
    param([bool]$AutoApprove, [bool]$DryRun)
    
    Write-Section "Installing AI Infrastructure"
    
    if ($DryRun) {
        Write-Warning "DRY RUN - No files will be created or removed"
        Write-Host ""
    }
    
    # Step 1: Clean up deprecated files first (automatic part of installation)
    Write-Host "Checking for deprecated files..." -ForegroundColor Gray
    $deprecatedFiles = Remove-DeprecatedFiles -ManifestConfig $ManifestConfig -WorkspaceRoot $Global:WorkspaceRoot -DryRun $DryRun -Quiet $true
    
    if ($deprecatedFiles.Count -gt 0) {
        Write-OperationStatus -Message "  Removed $($deprecatedFiles.Count) deprecated files" -Type "Success"
    } else {
        Write-OperationStatus -Message "  No deprecated files found" -Type "Info"
    }
    Write-Host ""
    
    # Step 2: Install/update current files
    Write-OperationStatus -Message "Installing current AI infrastructure files..." -Type "Info"
    
    # Debug: Show parameter values
    Write-Verbose "Parameters: AutoApprove=$AutoApprove, DryRun=$DryRun"
    
    # Use the FileOperations module to actually install files
    try {
        $result = Install-AllAIFiles -Force:$AutoApprove -DryRun:$DryRun
        
        if ($result.OverallSuccess) {
            Show-InstallationResults -Results $result
        } else {
            Show-InstallationResults -Results $result
        }
        
        return @{ Success = $result.OverallSuccess; Details = $result }
    }
    catch {
        Write-Host "Installation failed: $($_.Exception.Message)" -ForegroundColor Red
        return @{ Success = $false; Error = $_.Exception.Message }
    }
}

function Main {
    <#
    .SYNOPSIS
    Main entry point for the installer
    #>
    
    try {
        # Convert hyphenated parameter names to camelCase variables
        $AutoApprove = ${Auto-Approve}
        $DryRun = ${Dry-Run}
        
        # Check if running from user profile installer
        $isUserProfileInstaller = $PSScriptRoot -like "*\.terraform-ai-installer*"
        
        # If running from user profile installer, RepoDirectory is required for proper branch detection
        if ($isUserProfileInstaller -and -not $RepoDirectory -and -not $Help) {
            Write-Error "REPOSITORY DIRECTORY REQUIRED: When running from user profile installer, -RepoDirectory is required."
            Show-ErrorBlock -Issue "You're running the installer from your user profile location, but haven't specified where to find the terraform-provider-azurerm repository for branch detection." -Solutions @(
                "Use the -RepoDirectory parameter to specify the repository path"
            ) -ExampleUsage "`"$PSCommandPath`" -RepoDirectory `"C:\github.com\hashicorp\terraform-provider-azurerm`"" -AdditionalInfo "The installer needs to detect your current git branch to determine whether you're working on the main development branch or a feature branch. This affects which operations are available and how files are managed."
            exit 1
        }
        
        # Determine branch type
        $branchType = Test-BranchType
        $isSourceBranch = $branchType -eq "source"
        
        # Handle unknown branch type (additional safety check)
        if ($branchType -eq "unknown" -and -not $RepoDirectory -and -not $Help) {
            Write-Error "REPOSITORY DIRECTORY REQUIRED: Cannot determine git branch from current location."
            Show-ErrorBlock -Issue "The installer cannot determine the current git branch, which usually means:" -Solutions @(
                "You're running from outside a git repository",
                "You're running from your user profile location", 
                "Git is not available or the directory is not a git repository"
            ) -ExampleUsage "`"$PSCommandPath`" -RepoDirectory `"C:\github.com\hashicorp\terraform-provider-azurerm`"" -AdditionalInfo "The -RepoDirectory parameter tells the installer where to find the git repository for branch detection when running from outside the repository directory."
            exit 1
        }
        
        # Handle help parameter
        if ($Help) {
            Show-Help
            return
        }
        
        # Handle bootstrap parameter (source branch only)
        if ($Bootstrap) {
            if (-not $isSourceBranch) {
                Write-Error "Bootstrap can only be run from the source branch ($($Global:InstallerConfig.Branch))"
                return
            }
            
            $result = Invoke-Bootstrap
            if (-not $result.Success) {
                exit 1
            }
            return
        }
        
        # Handle verify parameter
        if ($Verify) {
            $result = Invoke-VerifyWorkspace
            return
        }
        
        # Handle clean parameter (feature branch only)
        if ($Clean) {
            if ($isSourceBranch) {
                Write-Error "Clean operation not available on source branch. This would remove development files."
                return
            }
            
            $result = Invoke-CleanWorkspace -AutoApprove:$AutoApprove -DryRun:$DryRun
            if (-not $result.Success) {
                exit 1
            }
            return
        }
        
        # Default installation flow
        Write-Header -Title "Terraform AzureRM Provider - AI Infrastructure Installer" -Version $Global:InstallerConfig.Version
        
        if ($isSourceBranch) {
            Show-SourceBranchWelcome -BranchName $Global:InstallerConfig.Branch -BootstrapCommand ".\install-copilot-setup.ps1 -Bootstrap"
            
            # Interactive prompt for better UX
            $shouldBootstrap = Get-BootstrapConfirmation
            
            if ($shouldBootstrap) {
                $result = Invoke-Bootstrap
                if (-not $result.Success) {
                    exit 1
                }
                return
            } else {
                return
            }
        } else {
            if ($branchType -eq "unknown") {
                # If branch type is unknown, we need to determine why and provide appropriate error
                Show-UnknownBranchError -HasRepoDirectory ([bool]$RepoDirectory) -RepoDirectory $RepoDirectory -ScriptPath $PSCommandPath
                exit 1
            }
            
            $workspaceRoot = if ($RepoDirectory) { $RepoDirectory } else { $Global:WorkspaceRoot }
            $currentBranch = Get-CurrentBranch -WorkspaceRoot $workspaceRoot
            Show-BranchDetection -BranchName $currentBranch
            Write-Host ""
            Write-OperationStatus -Message "Starting AI infrastructure installation..." -Type "Info"
            
            $result = Invoke-InstallInfrastructure -AutoApprove:$AutoApprove -DryRun:$DryRun
            if (-not $result.Success) {
                exit 1
            }
        }
    }
    catch {
        $errorMessage = $_.Exception.Message
        Write-Error "Installer failed with error: $errorMessage"
        
        Write-Host "Stack trace:" -ForegroundColor Red
        Write-Host $_.ScriptStackTrace -ForegroundColor Gray
        
        exit 1
    }
}

# Execute main function
Main
