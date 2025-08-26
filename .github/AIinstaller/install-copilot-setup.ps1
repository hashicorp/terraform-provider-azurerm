# Main AI Infrastructure Installer for Terraform AzureRM Provider
# Version: 1.0.0
# Description: Interactive installer for AI-powered development infrastructure

#requires -version 5.1

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

# ============================================================================
# MODULE LOADING - This must succeed or the script cannot continue
# ============================================================================

function Get-ModulesPath {
    param([string]$ScriptDirectory)
    
    # Simple logic: modules are always in the same relative location
    $ModulesPath = Join-Path $ScriptDirectory "modules\powershell"
    
    # If not found, try from workspace root (for direct repo execution)
    if (-not (Test-Path $ModulesPath)) {
        $currentPath = $ScriptDirectory
        while ($currentPath -and $currentPath -ne (Split-Path $currentPath -Parent)) {
            if (Test-Path (Join-Path $currentPath "go.mod")) {
                $ModulesPath = Join-Path $currentPath ".github\AIinstaller\modules\powershell"
                break
            }
            $currentPath = Split-Path $currentPath -Parent
        }
    }
    
    return $ModulesPath
}

function Import-RequiredModules {
    param([string]$ModulesPath)
    
    # Define all required modules in dependency order
    $modules = @(
        "ConfigParser",
        "UI", 
        "ValidationEngine",
        "FileOperations"
    )
    
    # Load each module cleanly
    foreach ($module in $modules) {
        $modulePath = Join-Path $ModulesPath "$module.psm1"
        
        if (-not (Test-Path $modulePath)) {
            throw "Required module '$module' not found at: $modulePath"
        }
        
        try {
            Remove-Module $module -Force -ErrorAction SilentlyContinue
            Import-Module $modulePath -Force -DisableNameChecking -Global -ErrorAction Stop
        }
        catch {
            throw "Failed to import module '$module': $_"
        }
    }
    
    # Verify critical functions are available
    $requiredFunctions = @("Get-ManifestConfig", "Get-InstallerConfig", "Write-Header", "Invoke-VerifyWorkspace")
    foreach ($func in $requiredFunctions) {
        if (-not (Get-Command $func -ErrorAction SilentlyContinue)) {
            throw "Required function '$func' not available after module loading"
        }
    }
}

# Get script directory with robust detection
$ScriptDirectory = if ($PSScriptRoot) { 
    $PSScriptRoot 
} elseif ($MyInvocation.MyCommand.Path) { 
    Split-Path $MyInvocation.MyCommand.Path -Parent 
} else { 
    # Fallback: assume we're in the AIinstaller directory
    Get-Location | ForEach-Object { $_.Path }
}

# Load modules with clear error handling
try {
    $ModulesPath = Get-ModulesPath -ScriptDirectory $ScriptDirectory
    Import-RequiredModules -ModulesPath $ModulesPath
}
catch {
    Write-Host "FATAL ERROR: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Cannot continue without required modules." -ForegroundColor Red
    exit 1
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

# Initialize workspace root after module loading
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
            Write-Host ""
            Write-Host "INVALID REPOSITORY: The specified directory does not appear to be a terraform-provider-azurerm repository." -ForegroundColor "Red"
            Write-Host "The -RepoDirectory parameter must point to a valid terraform-provider-azurerm repository root." -ForegroundColor "Red"
            # Clear RepoDirectory to show help instead of exiting
            $RepoDirectory = ""
        }
    } else {
        Write-Host ""
        Write-Host "DIRECTORY NOT FOUND: The specified RepoDirectory does not exist." -ForegroundColor "Red"
        Write-Host "The path '$RepoDirectory' could not be found on this system." -ForegroundColor "Red"
        # Clear RepoDirectory to show help instead of exiting
        $RepoDirectory = ""
    }
}

# Configuration will be loaded on-demand in functions that need it
$Global:ManifestConfig = $null
$Global:InstallerConfig = $null

# ============================================================================
# WORKSPACE DETECTION - Simple and reliable
# ============================================================================

function Get-WorkspaceRoot {
    param([string]$RepoDirectory, [string]$ScriptDirectory)
    
    # If RepoDirectory is provided, use it
    if ($RepoDirectory) {
        if (Test-Path $RepoDirectory) {
            return $RepoDirectory
        } else {
            throw "Specified RepoDirectory does not exist: $RepoDirectory"
        }
    }
    
    # Otherwise, find workspace root from script location
    $currentPath = $ScriptDirectory
    while ($currentPath -and $currentPath -ne (Split-Path $currentPath -Parent)) {
        if (Test-Path (Join-Path $currentPath "go.mod")) {
            return $currentPath
        }
        $currentPath = Split-Path $currentPath -Parent
    }
    
    # If no workspace found, return the directory where the script was called from
    # This allows help and other functions to work, with validation happening separately
    return (Get-Location).Path
}

# ============================================================================
# MAIN EXECUTION - Clean and simple
# ============================================================================

function Invoke-CleanWorkspace {
    param([bool]$AutoApprove, [bool]$DryRun)
    
    Write-Separator
    Write-Host "Clean Workspace" -ForegroundColor Cyan
    Write-Separator
    Write-Host ""
    
    if ($DryRun) {
        Write-Host "DRY RUN - No files will be deleted" -ForegroundColor Yellow
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
            
            # Handle dry-run vs actual operation messaging differently
            if ($DryRun) {
                # For dry-run, show positive confirmation that files were verified
                $dryRunIssues = $result.Issues | Where-Object { $_ -match "Dry run - no changes made" }
                $actualIssues = $result.Issues | Where-Object { $_ -notmatch "Dry run - no changes made" }
                
                if ($dryRunIssues.Count -gt 0) {
                    Write-Host "Dry run completed successfully - all $($dryRunIssues.Count) files verified and ready for removal" -ForegroundColor Green
                    Write-Host ""
                    Write-Host "Files that would be removed:" -ForegroundColor Cyan
                    Write-Separator
                    foreach ($issue in $dryRunIssues) {
                        # Extract just the filename from the error message
                        $fileName = ($issue -split ": Dry run")[0] -replace "Failed to remove ", ""
                        Write-Host "  - $fileName" -ForegroundColor Gray
                    }
                }
                
                # Show any actual issues (non-dry-run related)
                if ($actualIssues.Count -gt 0) {
                    Write-Host "Actual Issues Encountered:" -ForegroundColor Cyan
                    Write-Separator
                    foreach ($issue in $actualIssues) {
                        Write-Host "  - $issue" -ForegroundColor Red
                    }
                    Write-Host ""
                }
            } else {
                # For actual operations, show the issues as errors
                Write-Host "Clean Operation Encountered Issues:" -ForegroundColor Cyan
                foreach ($issue in $result.Issues) {
                    Write-Host "  - $issue" -ForegroundColor Red
                }
                Write-Host ""
            }
        }
        
        return $result
    }
    catch {
        Write-Host "Failed to clean workspace: $($_.Exception.Message)" -ForegroundColor Red
        return @{ Success = $false; Issues = @($_.Exception.Message) }
    }
}

function Invoke-InstallInfrastructure {
    param([bool]$AutoApprove, [bool]$DryRun)
    
    Write-Separator
    Write-Host "Installing AI Infrastructure" -ForegroundColor Cyan
    Write-Separator
    Write-Host ""
    
    if ($DryRun) {
        Write-Host "DRY RUN - No files will be created or removed" -ForegroundColor Yellow
        Write-Host ""
    }
    
    # Step 1: Clean up deprecated files first (automatic part of installation)
    Write-Host "Checking for deprecated files..." -ForegroundColor Gray
    $deprecatedFiles = Remove-DeprecatedFiles -ManifestConfig $Global:ManifestConfig -WorkspaceRoot $Global:WorkspaceRoot -DryRun $DryRun -Quiet $true
    
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
            $currentBranch = $Global:InstallerConfig.Branch
            $branchType = if (Test-SourceRepository) { "source" } else { 
                if ($currentBranch -eq "Unknown") { "Unknown" } else { "feature" }
            }
            
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
        # Step 1: Initialize workspace and validate it's a proper terraform-provider-azurerm repo
        $Global:WorkspaceRoot = Get-WorkspaceRoot -RepoDirectory $RepoDirectory -ScriptDirectory $ScriptDirectory
        
        # Step 2: Early workspace validation before doing anything else
        $workspaceValidation = Test-WorkspaceValid -WorkspacePath $Global:WorkspaceRoot
        
        # Initialize configuration based on workspace validity
        if ($workspaceValidation.Valid) {
            # Step 3: Initialize configuration (this sets up global branch info)
            # CRITICAL: Manifest file should be in the installer directory, not the target repository
            if ($RepoDirectory) {
                # Running from user profile - manifest is in the installer directory (where this script is)
                $manifestPath = Join-Path $ScriptDirectory "file-manifest.config"
            } else {
                # Running from source repository - manifest is in the repository's AIinstaller directory
                $manifestPath = Join-Path $Global:WorkspaceRoot ".github/AIinstaller/file-manifest.config"
            }
            $Global:ManifestConfig = Get-ManifestConfig -ManifestPath $manifestPath
            $Global:InstallerConfig = Get-InstallerConfig -WorkspaceRoot $Global:WorkspaceRoot -ManifestConfig $Global:ManifestConfig
        } else {
            # Invalid workspace - provide minimal configuration for help
            $Global:InstallerConfig = @{ Version = "1.0.0" }
            $Global:ManifestConfig = @{}
        }
        
        # Step 4: Simple branch safety check for -RepoDirectory operations
        if ($RepoDirectory) {
            # Get current branch of the target repository
            $originalLocation = Get-Location
            try {
                Set-Location $Global:WorkspaceRoot
                $currentBranch = git branch --show-current 2>$null
                if (-not $currentBranch -or $currentBranch.Trim() -eq "") {
                    $currentBranch = "Unknown"
                }
            }
            catch {
                $currentBranch = "Unknown"
            }
            finally {
                Set-Location $originalLocation
            }
        } else {
            # Not using -RepoDirectory, get branch info from current location
            try {
                $currentBranch = git branch --show-current 2>$null
                if (-not $currentBranch -or $currentBranch.Trim() -eq "") {
                    $currentBranch = "Unknown"
                }
            }
            catch {
                $currentBranch = "Unknown"
            }
        }
        
        $isSourceRepo = ($currentBranch -eq "exp/terraform_copilot")
        $branchType = if ($isSourceRepo) { "source" } else { 
            if ($currentBranch -eq "Unknown") { "Unknown" } else { "feature" }
        }
        
        # Convert hyphenated parameter names to camelCase variables
        $AutoApprove = ${Auto-Approve}
        $DryRun = ${Dry-Run}
        
        # CONSISTENT PATTERN: Every operation gets the same header and branch detection
        Write-Header -Title "Terraform AzureRM Provider - AI Infrastructure Installer" -Version $Global:InstallerConfig.Version
        Show-BranchDetection -BranchName $currentBranch -BranchType $branchType
        
        # Simple parameter handling
        if ($Help) {
            Write-Separator
            Show-Help -BranchType $branchType -WorkspaceValid $workspaceValidation.Valid -WorkspaceIssue $workspaceValidation.Reason
            return
        }
        
        # For all other operations, workspace must be valid
        if (-not $workspaceValidation.Valid) {
            Write-Separator
            Write-Host "WORKSPACE VALIDATION FAILED: $($workspaceValidation.Reason)" -ForegroundColor Red
            Write-Host "Please ensure you're running this script from within a terraform-provider-azurerm repository." -ForegroundColor Red
            Write-Separator
            Write-Host ""
            
            # Show help menu for guidance
            Show-Help -BranchType $branchType -WorkspaceValid $false -WorkspaceIssue $workspaceValidation.Reason
            exit 1
        }
        
        if ($Verify) {
            $result = Invoke-VerifyWorkspace
            return
        }
        
        if ($Bootstrap) {
            $result = Invoke-Bootstrap -AutoApprove $AutoApprove -DryRun $DryRun
            return
        }
        
        if ($Clean) {
            # Safety check for source branch operations
            if ($RepoDirectory -and $currentBranch -eq "exp/terraform_copilot") {
                Show-SafetyViolation -BranchName $currentBranch -Operation "Clean" -FromUserProfile
                exit 1
            }
            
            $result = Invoke-CleanWorkspace -AutoApprove $AutoApprove -DryRun $DryRun
            return
        }
        
        # Installation path (when -RepoDirectory is provided and not other specific operations)
        if ($RepoDirectory -and -not ($Help -or $Verify -or $Bootstrap -or $Clean)) {
            # Safety check for source branch operations
            if ($currentBranch -eq "exp/terraform_copilot") {
                Show-SafetyViolation -BranchName $currentBranch -Operation "Install" -FromUserProfile
                exit 1
            }
            
            # Proceed with installation on validated branch
            $result = Invoke-InstallInfrastructure -AutoApprove $AutoApprove -DryRun $DryRun
            return
        }
        
        # Default: show source branch help and welcome
        Write-Separator
        Show-SourceBranchHelp
        Show-SourceBranchWelcome -BranchName $currentBranch
        
    }
    catch {
        Write-Host ""
        Write-Host "ERROR: $($_.Exception.Message)" -ForegroundColor Red
        Write-Host ""
        exit 1
    }
}

# Execute main function
Main
