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
    
    throw "Could not find terraform-provider-azurerm workspace root. Use -RepoDirectory parameter."
}

# ============================================================================
# MAIN EXECUTION - Clean and simple
# ============================================================================

function Invoke-CleanWorkspace {
    param([bool]$AutoApprove, [bool]$DryRun)
    
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
                if ($currentBranch -eq "unknown") { "unknown" } else { "feature" }
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
        if (-not $workspaceValidation.Valid) {
            Write-Host "WORKSPACE VALIDATION FAILED: $($workspaceValidation.Reason)" -ForegroundColor Red
            Write-Host "Please ensure you're running this script from within a terraform-provider-azurerm repository." -ForegroundColor Red
            exit 1
        }
        
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
        
        # Step 4: Run centralized validation (replaces scattered Test-SourceRepository calls)
        $validation = Test-PreInstallation -AllowBootstrapOnSource:$Bootstrap
        
        # CRITICAL: Check validation results and exit immediately if unsafe
        if (-not $validation.OverallValid) {
            if ($validation.Git.Reason -like "*SAFETY VIOLATION*") {
                Write-Host "SAFETY VIOLATION DETECTED" -ForegroundColor Red
                Write-Separator -Character "-" -Color Red
                Write-Host ""
                Write-Host $validation.Git.Reason -ForegroundColor Red
                Write-Host ""
                Write-Host "SOLUTION:" -ForegroundColor Yellow
                Write-Host "Switch to a feature branch before running the installer:" -ForegroundColor White
                Write-Host "  cd `"$Global:WorkspaceRoot`"" -ForegroundColor Gray
                Write-Host "  git checkout -b feature/your-branch-name" -ForegroundColor Gray
                Write-Host ""
                exit 1
            } else {
                Write-Host "VALIDATION FAILED: $($validation.Git.Reason)" -ForegroundColor Red
                exit 1
            }
        }
        
        $currentBranch = $validation.Git.CurrentBranch
        # Handle empty or null branch names
        if (-not $currentBranch -or $currentBranch.Trim() -eq "") {
            $currentBranch = "unknown"
        }
        $isSourceRepo = ($currentBranch -eq "exp/terraform_copilot")
        $branchType = if ($isSourceRepo) { "source" } else { 
            if ($currentBranch -eq "unknown") { "unknown" } else { "feature" }
        }
        
        # Convert hyphenated parameter names to camelCase variables
        $AutoApprove = ${Auto-Approve}
        $DryRun = ${Dry-Run}
        
        # Simple parameter handling
        if ($Help) {
            Write-Header -Title "Terraform AzureRM Provider - AI Infrastructure Installer" -Version $Global:InstallerConfig.Version
            Show-Help -BranchName $currentBranch -BranchType $branchType -SkipHeader:$true
            return
        }
        
        if ($Verify) {
            Write-Header -Title "Terraform AzureRM Provider - AI Infrastructure Installer" -Version $Global:InstallerConfig.Version
            Show-BranchDetection -BranchName $currentBranch -BranchType $branchType
            $result = Invoke-VerifyWorkspace
            return
        }
        
        if ($Bootstrap) {
            Write-Header -Title "Terraform AzureRM Provider - AI Infrastructure Installer" -Version $Global:InstallerConfig.Version
            Show-BranchDetection -BranchName $currentBranch -BranchType $branchType
            $result = Invoke-Bootstrap -AutoApprove $AutoApprove -DryRun $DryRun
            return
        }
        
        if ($Clean) {
            Write-Header -Title "Terraform AzureRM Provider - AI Infrastructure Installer" -Version $Global:InstallerConfig.Version
            Show-BranchDetection -BranchName $currentBranch -BranchType $branchType
            $result = Invoke-CleanWorkspace -AutoApprove $AutoApprove -DryRun $DryRun
            return
        }
        
        # Installation path (when -RepoDirectory is provided and not other specific operations)
        if ($RepoDirectory -and -not ($Help -or $Verify -or $Bootstrap -or $Clean)) {
            Write-Header -Title "Terraform AzureRM Provider - AI Infrastructure Installer" -Version $Global:InstallerConfig.Version
            Show-BranchDetection -BranchName $currentBranch -BranchType $branchType
            
            # Safety check already handled by Test-PreInstallation above
            # Proceed with installation on validated branch
            $result = Invoke-InstallInfrastructure -AutoApprove $AutoApprove -DryRun $DryRun
            return
        }
        
        # Default: show help
        Write-Header -Title "Terraform AzureRM Provider - AI Infrastructure Installer" -Version $Global:InstallerConfig.Version
        Show-SourceBranchHelp -BranchName $currentBranch -WorkspacePath $Global:WorkspaceRoot
        Write-Host ""
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
