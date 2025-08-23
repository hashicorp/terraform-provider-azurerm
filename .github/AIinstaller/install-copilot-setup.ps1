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

# Suppress verbose output during module loading to keep output clean
$OriginalVerbosePreference = $VerbosePreference
$VerbosePreference = 'SilentlyContinue'

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
        } else {
            Write-Host "INVALID REPOSITORY: The specified directory does not appear to be a terraform-provider-azurerm repository." -ForegroundColor "Red"
            Write-Host "The -RepoDirectory parameter must point to a valid terraform-provider-azurerm repository root." -ForegroundColor "Red"
            exit 1
        }
    } else {
        Write-Host "DIRECTORY NOT FOUND: The specified RepoDirectory does not exist." -ForegroundColor "Red"
        Write-Host "The path '$RepoDirectory' could not be found on this system." -ForegroundColor "Red"
        exit 1
    }
}

# Import required modules from the script location (for bootstrapped installs)
$ScriptDirectory = Split-Path $MyInvocation.MyCommand.Path -Parent
$ModulesPath = Join-Path $ScriptDirectory "modules\powershell"

# If modules don't exist locally, try workspace location (for direct repo execution)
if (-not (Test-Path $ModulesPath)) {
    $ModulesPath = Join-Path $Global:WorkspaceRoot ".github\AIinstaller\modules\powershell"
}

$RequiredModules = @(
    "ConfigParser",
    "FileOperations",
    "ValidationEngine", 
    "UI"
)

foreach ($module in $RequiredModules) {
    $modulePath = Join-Path $ModulesPath "$module.psm1"
    if (Test-Path $modulePath) {
        try {
            # Import module with explicit scope to ensure functions are available
            Import-Module $modulePath -Force -DisableNameChecking -Scope Global 4>$null
            
            # Load configuration immediately after ConfigParser import
            if ($module -eq "ConfigParser") {
                $manifestPath = Join-Path $Global:WorkspaceRoot ".github/AIinstaller/file-manifest.config"
                $Global:ManifestConfig = ConfigParser\Get-ManifestConfig -ManifestPath $manifestPath
                $Global:InstallerConfig = ConfigParser\Get-InstallerConfig -WorkspaceRoot $Global:WorkspaceRoot -ManifestConfig $Global:ManifestConfig
            }
        }
        catch {
            Write-Host "ERROR: Failed to import module $module`: $_" -ForegroundColor Red
            exit 1
        }
    } else {
        Write-Host "Required module '$module' not found at: $modulePath" -ForegroundColor "Red"
        exit 1
    }
}

# Restore original verbose preference after module loading
$VerbosePreference = $OriginalVerbosePreference

# Note: Branch detection will be shown later in the contextual help display
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
            Write-WarningMessage "Repository directory not found: $workspaceRoot"
            return "unknown"
        }
        
        Push-Location $workspaceRoot
        try {
            $currentBranch = git rev-parse --abbrev-ref HEAD 2>$null
            if ($LASTEXITCODE -ne 0) {
                Write-WarningMessage "Could not determine current git branch"
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
        Write-WarningMessage "Git not available or error checking git repository"
        return "unknown"
    }
}

function Invoke-CleanWorkspace {
    param([bool]$AutoApprove, [bool]$DryRun)
    
    Write-Section "Clean Workspace"
    
    if ($DryRun) {
        Write-WarningMessage "DRY RUN - No files will be deleted"
        Write-Host ""
    }
    
    # Use the FileOperations module to properly remove all AI files
    try {
        $result = Remove-AllAIFiles -Force:$AutoApprove -DryRun:$DryRun -WorkspaceRoot $Global:WorkspaceRoot
        
        if ($result.Success) {
            # Use the superior summary function
            $details = @{
                "Files removed" = $result.FilesRemoved
                "Directories cleaned" = $result.DirectoriesCleaned
                "Operation type" = if ($DryRun) { "Dry run (simulation)" } else { "Live cleanup" }
            }
            
            Show-Summary -Title "Clean Operation Results" -Details $details
        } else {
            Write-Host ""
            Write-WarningMessage "Clean operation encountered issues:"
            foreach ($issue in $result.Issues) {
                Write-Host "  - $issue" -ForegroundColor Red
            }
        }
        
        return $result
    }
    catch {
        Write-ErrorMessage "Failed to clean workspace: $($_.Exception.Message)"
        return @{ Success = $false; Issues = @($_.Exception.Message) }
    }
}

function Invoke-InstallInfrastructure {
    param([bool]$AutoApprove, [bool]$DryRun)
    
    Write-Section "Installing AI Infrastructure"
    
    if ($DryRun) {
        Write-WarningMessage "DRY RUN - No files will be created or removed"
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
    Write-VerboseMessage "Parameters: AutoApprove=$AutoApprove, DryRun=$DryRun"
    
    # Use the FileOperations module to actually install files
    try {
        $result = Install-AllAIFiles -Force:$AutoApprove -DryRun:$DryRun -WorkspaceRoot $Global:WorkspaceRoot
        
        if ($result.OverallSuccess) {
            # Use the superior completion summary function
            $nextSteps = @()
            if ($result.Skipped -gt 0) {
                $nextSteps += "Review skipped files and use -Auto-Approve if needed"
            }
            $nextSteps += "Start using GitHub Copilot with your new AI-powered infrastructure"
            $nextSteps += "Check the .github/instructions/ folder for detailed guidelines"
            
            # Get branch information for completion summary
            $workspaceRoot = if ($RepoDirectory) { $RepoDirectory } else { $Global:WorkspaceRoot }
            $currentBranch = Get-CurrentBranch -WorkspaceRoot $workspaceRoot
            $isSourceRepo = Test-SourceRepository
            $branchType = if ($isSourceRepo) { "source" } else { "feature" }
            
            Show-CompletionSummary -FilesInstalled $result.Successful -FilesSkipped $result.Skipped -FilesFailed $result.Failed -NextSteps $nextSteps -BranchName $currentBranch -BranchType $branchType
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
            Write-ErrorMessage "REPOSITORY DIRECTORY REQUIRED: When running from user profile installer, -RepoDirectory is required."
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
            Write-ErrorMessage "REPOSITORY DIRECTORY REQUIRED: Cannot determine git branch from current location."
            Show-ErrorBlock -Issue "The installer cannot determine the current git branch, which usually means:" -Solutions @(
                "You're running from outside a git repository",
                "You're running from your user profile location", 
                "Git is not available or the directory is not a git repository"
            ) -ExampleUsage "`"$PSCommandPath`" -RepoDirectory `"C:\github.com\hashicorp\terraform-provider-azurerm`"" -AdditionalInfo "The -RepoDirectory parameter tells the installer where to find the git repository for branch detection when running from outside the repository directory."
            exit 1
        }
        
        # Handle help parameter
        if ($Help) {
            # Get branch information for context
            $workspaceRoot = if ($RepoDirectory) { $RepoDirectory } else { $Global:WorkspaceRoot }
            $currentBranch = Get-CurrentBranch -WorkspaceRoot $workspaceRoot
            $isSourceRepo = Test-SourceRepository
            $branchType = if ($isSourceRepo) { "source" } else { "feature" }
            
            Show-Help -BranchName $currentBranch -BranchType $branchType
            return
        }
        
        # Handle bootstrap parameter (source branch only)
        if ($Bootstrap) {
            if (-not $isSourceBranch) {
                Write-ErrorMessage "Bootstrap can only be run from the source branch ($($Global:InstallerConfig.Branch))"
                return
            }
            
            # Show consistent header with branch detection for bootstrap
            Write-Header -Title "Terraform AzureRM Provider - AI Infrastructure Installer" -Version $Global:InstallerConfig.Version
            
            # Get branch information for context
            $workspaceRoot = if ($RepoDirectory) { $RepoDirectory } else { $Global:WorkspaceRoot }
            $currentBranch = Get-CurrentBranch -WorkspaceRoot $workspaceRoot
            $isSourceRepo = Test-SourceRepository
            $branchType = if ($isSourceRepo) { "source" } else { "feature" }
            
            Show-BranchDetection -BranchName $currentBranch -BranchType $branchType
            
            # Show workspace path for consistency
            $formattedWorkspaceLabel = Format-AlignedLabel -Label "WORKSPACE" -CurrentBranchType $branchType
            Write-Host $formattedWorkspaceLabel -ForegroundColor Cyan -NoNewline
            Write-Host $workspaceRoot -ForegroundColor Green
            Write-Host ""
            
            $result = Invoke-Bootstrap
            if (-not $result.Success) {
                exit 1
            }
            return
        }
        
        # Handle verify parameter
        if ($Verify) {
            # Show consistent header for verify operation
            Write-Header -Title "Terraform AzureRM Provider - AI Infrastructure Installer" -Version $Global:InstallerConfig.Version
            
            # Get branch information for verification context
            $workspaceRoot = if ($RepoDirectory) { $RepoDirectory } else { $Global:WorkspaceRoot }
            $currentBranch = Get-CurrentBranch -WorkspaceRoot $workspaceRoot
            $isSourceRepo = Test-SourceRepository
            $branchType = if ($isSourceRepo) { "source" } else { "feature" }
            
            # Show branch context before verification
            UI\Show-BranchDetection -BranchName $currentBranch -BranchType $branchType
            
            $result = Invoke-VerifyWorkspace -BranchType $branchType
            return
        }
        
        # Handle clean parameter (feature branch only)
        if ($Clean) {
            if ($isSourceBranch) {
                Write-ErrorMessage "Clean operation not available on source branch. This would remove development files."
                return
            }
            
            $result = Invoke-CleanWorkspace -AutoApprove:$AutoApprove -DryRun:$DryRun
            if (-not $result.Success) {
                exit 1
            }
            return
        }
        
        # Default installation flow (when no specific action requested)
        Write-Header -Title "Terraform AzureRM Provider - AI Infrastructure Installer" -Version $Global:InstallerConfig.Version
        
        # Check if any specific action was requested
        $hasExplicitAction = $Bootstrap -or $Verify -or $Clean -or $Help -or $DryRun
        
        if (-not $hasExplicitAction) {
            # No action specified - just show help and exit
            if ($isSourceBranch) {
                Show-SourceBranchHelp -BranchName $Global:InstallerConfig.Branch -WorkspacePath $Global:WorkspaceRoot
                Write-Host ""
                Show-SourceBranchWelcome -BranchName $Global:InstallerConfig.Branch
            } else {
                if ($branchType -eq "unknown") {
                    Show-UnknownBranchError -HasRepoDirectory:$PSBoundParameters.ContainsKey('RepoDirectory') -RepoDirectory $RepoDirectory -ScriptPath $PSCommandPath
                } else {
                    Show-FeatureBranchHelp
                }
            }
            return  # Exit cleanly after showing help
        }
        
        # If we get here, user requested a specific action on source branch that's not allowed
        if ($isSourceBranch) {
            Show-Help -BranchName $Global:InstallerConfig.Branch -BranchType "source" -SkipHeader
            Write-Host ""
            Write-ErrorMessage "Requested operation not available on source branch."
            Write-Host "Source branch is for development only. Use -Bootstrap, -Verify, or -Help." -ForegroundColor Yellow
            return
        }
        
        # If we get here, it means we're on a feature branch or unknown branch with no action specified
        if ($branchType -eq "unknown") {
            # If branch type is unknown, we need to determine why and provide appropriate error
            Show-UnknownBranchError -HasRepoDirectory ([bool]$RepoDirectory) -RepoDirectory $RepoDirectory -ScriptPath $PSCommandPath
            exit 1
        }
        
        # Feature branch with no action - just show help
        Show-FeatureBranchHelp
        return
    }
    catch {
        $errorMessage = $_.Exception.Message
        Write-ErrorMessage "Installer failed with error: $errorMessage"
        
        Write-Host "Stack trace:" -ForegroundColor Red
        Write-Host $_.ScriptStackTrace -ForegroundColor Gray
        
        exit 1
    }
}

# Execute main function
Main
