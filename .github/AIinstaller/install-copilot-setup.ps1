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
    Write-Host " FATAL ERROR: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host " Cannot continue without required modules." -ForegroundColor Red
    exit 1
}

# Initialize workspace root after module loading
$Global:WorkspaceRoot = $null

# Configuration will be loaded on-demand in functions that need it
$Global:ManifestConfig = $null
$Global:InstallerConfig = $null

# ============================================================================
# WORKSPACE DETECTION - Simple and reliable
# ============================================================================

function Get-WorkspaceRoot {
    param([string]$RepoDirectory, [string]$ScriptDirectory)

    # If RepoDirectory is provided, use it (validation happens later)
    if ($RepoDirectory) {
        return $RepoDirectory
    }

    # Otherwise, if no RepoDirectory provided, just use current directory
    # This allows fast-fail workspace validation to handle invalid directories
    return Get-Location

    # If no workspace found, return the directory where the script was called from
    # This allows help and other functions to work, with validation happening separately
    return (Get-Location).Path
}

# ============================================================================
# MAIN EXECUTION - Clean and simple
# ============================================================================

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

        # Step 3: Initialize global configuration
        if ($workspaceValidation.Valid) {
            # Only load manifest if workspace is valid
            if ($RepoDirectory) {
                # Running from user profile with -RepoDirectory - manifest is in the current user directory
                $manifestPath = Join-Path (Get-Location) "file-manifest.config"
            } else {
                # Running from source repository - manifest is in the source AIinstaller directory
                $manifestPath = Join-Path $ScriptDirectory "file-manifest.config"
            }
            $Global:ManifestConfig = Get-ManifestConfig -ManifestPath $manifestPath
            $Global:InstallerConfig = Get-InstallerConfig -WorkspaceRoot $Global:WorkspaceRoot -ManifestConfig $Global:ManifestConfig
        } else {
            # Invalid workspace - provide minimal configuration for UI display
            $Global:InstallerConfig = @{ Version = "1.0.0" }
            $Global:ManifestConfig = @{}
        }

        # Step 4: Get branch information - simple and direct
        try {
            $currentBranch = git branch --show-current 2>$null
            if (-not $currentBranch -or $currentBranch.Trim() -eq "") {
                $currentBranch = "Unknown"
            }
        }
        catch {
            $currentBranch = "Unknown"
        }

        # Step 4: Get branch information for UI display and safety checks
        if ($RepoDirectory) {
            # Get current branch of the target repository (only if workspace exists)
            $originalLocation = Get-Location
            $currentBranch = "Unknown"
            try {
                if (Test-Path $Global:WorkspaceRoot) {
                    Set-Location $Global:WorkspaceRoot
                    $currentBranch = git branch --show-current 2>$null
                    if (-not $currentBranch -or $currentBranch.Trim() -eq "") {
                        $currentBranch = "Unknown"
                    }
                }
            }
            catch {
                $currentBranch = "Unknown"
            }
            finally {
                if (Test-Path $originalLocation) {
                    Set-Location $originalLocation
                }
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

        # Check if current branch is a source branch (main, master, or exp/terraform_copilot)
        # Source branches are protected from AI infrastructure installation for safety
        $sourceBranches = @("main", "master", "exp/terraform_copilot")
        $isSourceRepo = ($currentBranch -in $sourceBranches)
        $branchType = if ($isSourceRepo) { "source" } else {
            if ($currentBranch -eq "Unknown") { "Unknown" } else { "feature" }
        }

        # Convert hyphenated parameter names to camelCase variables
        $AutoApprove = ${Auto-Approve}
        $DryRun = ${Dry-Run}

        # CONSISTENT PATTERN: Every operation gets the same header and branch detection
        Write-Header -Title "Terraform AzureRM Provider - AI Infrastructure Installer" -Version $Global:InstallerConfig.Version
        Show-BranchDetection -BranchName $currentBranch -BranchType $branchType

        # SAFETY CHECK - Block operations on source branch when using -RepoDirectory (except Verify, Help, Bootstrap)
        if ($RepoDirectory) {
            if ($currentBranch -in $sourceBranches -and -not ($Verify -or $Help -or $Bootstrap)) {
                Show-SafetyViolation -BranchName $currentBranch -Operation "Install" -FromUserProfile
                exit 1
            }
        }

        # Detect if we're running from user profile directory (needed for all help contexts)
        $currentDir = Get-Location
        $userProfileInstallerDir = Join-Path $env:USERPROFILE ".terraform-ai-installer"
        $isFromUserProfile = $currentDir.Path -eq $userProfileInstallerDir -or [bool]$RepoDirectory

        # Detect what command was attempted (for better error messages)
        $attemptedCommand = ""
        if ($Bootstrap) { $attemptedCommand = "-Bootstrap" }
        elseif ($Verify) { $attemptedCommand = "-Verify" }
        elseif ($Clean) { $attemptedCommand = "-Clean" }
        elseif ($Help) { $attemptedCommand = "-Help" }
        elseif ($DryRun) { $attemptedCommand = "-Dry-Run" }
        elseif ($AutoApprove) { $attemptedCommand = "-Auto-Approve" }
        elseif ($RepoDirectory -and -not ($Help -or $Verify -or $Bootstrap -or $Clean)) {
            $attemptedCommand = "-RepoDirectory `"$RepoDirectory`""
        }

        # Simple parameter handling
        if ($Help) {
            Show-Help -BranchType $branchType -WorkspaceValid $workspaceValidation.Valid -WorkspaceIssue $workspaceValidation.Reason -FromUserProfile $isFromUserProfile -AttemptedCommand $attemptedCommand
            return
        }

        # For all other operations, workspace must be valid
        if (-not $workspaceValidation.Valid) {
            Show-WorkspaceValidationError -Reason $workspaceValidation.Reason -FromUserProfile:$isFromUserProfile

            # Show help menu for guidance
            Show-Help -BranchType $branchType -WorkspaceValid $false -WorkspaceIssue $workspaceValidation.Reason -FromUserProfile $isFromUserProfile -AttemptedCommand $attemptedCommand
            exit 1
        }

        if ($Verify) {
            Invoke-VerifyWorkspace | Out-Null
            return
        }

        if ($Bootstrap) {
            Invoke-Bootstrap -AutoApprove $AutoApprove -DryRun $DryRun | Out-Null
            return
        }

        if ($Clean) {
            Invoke-CleanWorkspace -AutoApprove $AutoApprove -DryRun $DryRun -WorkspaceRoot $Global:WorkspaceRoot -FromUserProfile:([bool]$RepoDirectory) | Out-Null
            return
        }

        # Installation path (when -RepoDirectory is provided and not other specific operations)
        if ($RepoDirectory -and -not ($Help -or $Verify -or $Bootstrap -or $Clean)) {
            # Proceed with installation
            Invoke-InstallInfrastructure -AutoApprove $AutoApprove -DryRun $DryRun -WorkspaceRoot $Global:WorkspaceRoot -ManifestConfig $Global:ManifestConfig -TargetBranch $currentBranch | Out-Null
            return
        }

        # Default: show source branch help and welcome
        Show-SourceBranchHelp
        Show-SourceBranchWelcome -BranchName $currentBranch
    }
    catch {
        Write-Host ""
        Write-Host " ERROR: $($_.Exception.Message)" -ForegroundColor Red
        Write-Host ""
        exit 1
    }
}

# Execute main function
Main
