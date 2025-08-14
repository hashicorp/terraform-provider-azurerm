# Main AI Infrastructure Installer for Terraform AzureRM Provider
# Version: 1.0.0
# Description: Interactive installer for AI-powered development infrastructure

#requires -version 5.1

[CmdletBinding()]
param(
    [Parameter(HelpMessage = "Copy installer to user profile for feature branch use")]
    [switch]$Bootstrap,
    
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
$Global:WorkspaceRoot = Get-WorkspaceRoot

# Get dynamic configuration from manifest
$ManifestConfig = Get-ManifestConfig

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

function Test-BranchType {
    <#
    .SYNOPSIS
    Determine if running from source branch or feature branch
    #>
    
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
    catch {
        Write-Warning "Git not available or not in a git repository"
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
            Write-Host "  Created directory: $targetDirectory" -ForegroundColor "Green"
        } else {
            Write-Host "  Using existing directory: $targetDirectory" -ForegroundColor "Yellow"
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
            Write-Host "  Copying installer files from local source repository..." -ForegroundColor "Cyan"
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
                    
                    Write-Host "    Copying: $fileName" -ForegroundColor "Gray" -NoNewline
                    
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
            
            Show-Success "Bootstrap completed successfully!"
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
            Write-Host "     & `"$targetDirectory\install-copilot-setup.ps1`"" -ForegroundColor "Gray"
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
        Write-Host "DRY RUN - No files will be deleted" -ForegroundColor Yellow
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
        Write-Host "DRY RUN - No files will be created or removed" -ForegroundColor Yellow
        Write-Host ""
    }
    
    # Step 1: Clean up deprecated files first (automatic part of installation)
    Write-Host "Checking for deprecated files..." -ForegroundColor Gray
    $deprecatedFiles = Remove-DeprecatedFiles -ManifestConfig $ManifestConfig -WorkspaceRoot $Global:WorkspaceRoot -DryRun $DryRun -Quiet $true
    
    if ($deprecatedFiles.Count -gt 0) {
        Write-Host "  Removed $($deprecatedFiles.Count) deprecated files" -ForegroundColor Green
    } else {
        Write-Host "  No deprecated files found" -ForegroundColor Gray
    }
    Write-Host ""
    
    # Step 2: Install/update current files
    Write-Host "Installing current AI infrastructure files..." -ForegroundColor White
    
    # TODO: Implement actual file installation logic here
    # This should copy files from the manifest configuration to target locations
    Write-Host "Note: Full installation logic would be implemented here" -ForegroundColor Yellow
    
    return @{ Success = $true }
}

function Main {
    <#
    .SYNOPSIS
    Main entry point for the installer
    #>
    
    try {
        # Determine branch type
        $branchType = Test-BranchType
        $isSourceBranch = $branchType -eq "source"
        
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
            Write-Host "🚀 WELCOME TO AI-POWERED TERRAFORM DEVELOPMENT!" -ForegroundColor "Cyan"
            Write-Host ""
            Write-Host "You're running from the source branch ($($Global:InstallerConfig.Branch))." -ForegroundColor "White"
            Write-Host "To get started with AI infrastructure on your feature branch:" -ForegroundColor "White"
            Write-Host ""
            
            Write-Host "📋 QUICK START (Recommended):" -ForegroundColor "Green"
            Write-Host "  Run the bootstrap command to set up the installer:" -ForegroundColor "White"
            Write-Host "  .\install-copilot-setup.ps1 -Bootstrap" -ForegroundColor "Yellow"
            Write-Host ""
            
            Write-Host "📝 MANUAL WORKFLOW:" -ForegroundColor "Cyan"
            Write-Host "  1. Bootstrap: .\install-copilot-setup.ps1 -Bootstrap" -ForegroundColor "Gray"
            Write-Host "  2. Switch branch: git checkout feature/your-branch-name" -ForegroundColor "Gray"
            Write-Host "  3. Install: Run installer from user profile" -ForegroundColor "Gray"
            Write-Host ""
            
            Write-Host "🔍 OTHER OPTIONS:" -ForegroundColor "White"
            Write-Host "  -Verify       Check current workspace status" -ForegroundColor "Gray"
            Write-Host "  -Help         Show detailed help information" -ForegroundColor "Gray"
            Write-Host ""
            
            # Interactive prompt for better UX
            Write-Host "❓ Would you like to run bootstrap now? [Y/n]: " -NoNewline -ForegroundColor "Yellow"
            $response = Read-Host
            
            if ($response -eq "" -or $response -eq "y" -or $response -eq "Y") {
                Write-Host ""
                Write-Host "🔄 Running bootstrap automatically..." -ForegroundColor "Green"
                
                $result = Invoke-Bootstrap
                if (-not $result.Success) {
                    exit 1
                }
                return
            } else {
                Write-Host ""
                Write-Host "ℹ️  No problem! Run with -Bootstrap when you're ready." -ForegroundColor "Cyan"
                Write-Host "   Or use -Help for more information." -ForegroundColor "Cyan"
            }
        } else {
            Write-Host "FEATURE BRANCH DETECTED" -ForegroundColor "Cyan"
            Write-Host ""
            Write-Host "Starting AI infrastructure installation..." -ForegroundColor "White"
            
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
