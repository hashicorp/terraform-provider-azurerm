#!/usr/bin/env pwsh
<#
.SYNOPSIS
    Clean Architecture Main Installer for Terraform AzureRM Provider AI Setup

.DESCRIPTION
    This script uses the new clean modular architecture with real implementations
    from the working system. Features a progress bar and proper error handling.

.PARAMETER RepositoryPath
    Specifies the path to the terraform-provider-azurerm repository.
    If not provided, the script will attempt to auto-discover the repository location.

.PARAMETER Clean
    Removes all installed AI setup files and restores original VS Code settings from backups.

.PARAMETER Auto-Approve
    Skips interactive approval prompts.

.PARAMETER Verify
    Run verification only without installing

.PARAMETER Help
    Displays detailed help information and usage examples.

.EXAMPLE
    .\install-clean.ps1
    
.EXAMPLE
    .\install-clean.ps1 -Auto-Approve
    
.EXAMPLE
    .\install-clean.ps1 -Verify

.EXAMPLE
    .\install-clean.ps1 -Clean
#>

param(
    [string]$RepositoryPath,
    [switch]$Clean,
    [switch]${Auto-Approve},
    [switch]$Verify,
    [switch]$Help
)

# Get the script directory and import all clean modules
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ModulesDir = Join-Path $ScriptDir "modules"

# Import clean modules with real implementations
$CleanModulesDir = Join-Path $ModulesDir "clean"

# Also import working modules for functions we depend on
$WorkingModulesDir = Join-Path $ModulesDir "powershell"

try {
    # Import clean architecture modules only - no dependencies on old modules
    Import-Module (Join-Path $CleanModulesDir "01-core.psm1") -Force -ErrorAction Stop
    Import-Module (Join-Path $CleanModulesDir "02-validation.psm1") -Force -ErrorAction Stop
    Import-Module (Join-Path $CleanModulesDir "06-installation.psm1") -Force -ErrorAction Stop
    Import-Module (Join-Path $CleanModulesDir "07-ui.psm1") -Force -ErrorAction Stop
    
} catch {
    Write-Host "Failed to load required modules: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Show help if requested
if ($Help) {
    Show-Help
    exit 0
}

# Handle clean mode
if ($Clean) {
    try {
        Write-Host "Starting cleanup mode..." -ForegroundColor Yellow
        
        # Auto-detect repository if not provided
        if (-not $RepositoryPath) {
            $RepositoryPath = Find-RepositoryRoot
            if (-not $RepositoryPath) {
                throw "Could not auto-detect repository path. Please specify -RepositoryPath parameter."
            }
        }
        
        # Validate repository
        $validation = Test-RepositoryStructure -RepositoryPath $RepositoryPath
        if (-not $validation.IsValidRepository) {
            throw "Repository validation failed. Please ensure you're in the correct Terraform AzureRM Provider repository."
        }
        
        # Request confirmation unless Auto-Approve is used
        if (-not ${Auto-Approve}) {
            Write-Host ""
            Write-Host "This will remove all AI setup files and restore backups." -ForegroundColor Yellow
            Write-Host "Are you sure you want to continue? (y/N): " -NoNewline -ForegroundColor Yellow
            $response = Read-Host
            if ($response -ne 'y' -and $response -ne 'Y') {
                Write-Host "Cleanup cancelled." -ForegroundColor Green
                exit 0
            }
        }
        
        # Perform cleanup
        Write-Host "Removing AI setup files..." -ForegroundColor Yellow
        $cleanupResults = Remove-AIInstallation -RepositoryPath $RepositoryPath
        
        if ($cleanupResults.Success) {
            Write-Host "Cleanup completed successfully!" -ForegroundColor Green
        } else {
            Write-Host "Cleanup completed with some issues." -ForegroundColor Yellow
            foreach ($errorMsg in $cleanupResults.Errors) {
                Write-Host "  - $errorMsg" -ForegroundColor Red
            }
        }
        
        exit 0
        
    } catch {
        Write-Host "Cleanup failed: $($_.Exception.Message)" -ForegroundColor Red
        exit 1
    }
}

function Main {
    try {
        # Show welcome banner with progress bar style
        Show-WelcomeBanner -Version "2.0.0"
        
        # Auto-detect repository if not provided
        if (-not $RepositoryPath) {
            Write-StatusMessage "Auto-detecting repository path..." "Info"
            $RepositoryPath = Find-RepositoryRoot
            if (-not $RepositoryPath) {
                throw "Could not auto-detect repository path. Please specify -RepositoryPath parameter."
            }
        }
        
        Write-InstallationProgress -Stage "Validating Repository" -StageNumber 1 -TotalStages 5
        
        # Validate repository
        $validation = Test-RepositoryStructure -RepositoryPath $RepositoryPath
        Show-RepositoryInfo -RepositoryPath $RepositoryPath -ValidationResults $validation
        
        if (-not $validation.IsValidRepository) {
            throw "Repository validation failed. Please ensure you're in the correct Terraform AzureRM Provider repository."
        }
        
        # If verify mode, run verification and exit
        if ($Verify) {
            Write-InstallationProgress -Stage "Running Verification" -StageNumber 5 -TotalStages 5
            
            Write-StatusMessage "Running installation verification..." "Info"
            $verifyResult = Test-AIInstallation -RepositoryPath $RepositoryPath
            
            Show-InstallationSummary -Results $verifyResult
            
            if ($verifyResult.Success) {
                Write-StatusMessage "Verification completed successfully!" "Success"
                exit 0
            } else {
                Write-StatusMessage "Verification found issues" "Warning"
                exit 1
            }
        }
        
        Write-InstallationProgress -Stage "Checking Prerequisites" -StageNumber 2 -TotalStages 5
        
        # Check prerequisites
        if (-not (Test-Prerequisites)) {
            throw "Prerequisites check failed"
        }
        
        # Request confirmation
        if (-not (Request-UserConfirmation -Message "Ready to install AI setup. Continue?" -Force:${Auto-Approve})) {
            Write-StatusMessage "Installation cancelled by user" "Info"
            exit 0
        }
        
        Write-InstallationProgress -Stage "Installing Components" -StageNumber 3 -TotalStages 5
        
        # Run installation
        Write-StatusMessage "Starting installation with clean architecture..." "Info"
        $results = Start-CompleteInstallation -RepositoryPath $RepositoryPath -Force:${Auto-Approve}
        
        Write-InstallationProgress -Stage "Finalizing Installation" -StageNumber 4 -TotalStages 5
        
        # Show results
        Show-InstallationSummary -Results $results
        
        Write-InstallationProgress -Stage "Complete!" -StageNumber 5 -TotalStages 5
        
        if ($results.Success) {
            Show-NextSteps
            Write-StatusMessage "Clean architecture installation completed successfully!" "Success"
            exit 0
        } else {
            Write-StatusMessage "Installation completed with issues. Check the summary above." "Warning" 
            exit 1
        }
        
    } catch {
        Write-StatusMessage "Installation failed: $($_.Exception.Message)" "Error"
        Write-Host ""
        Write-Host "For troubleshooting help, check the logs above or run with -Verify to diagnose issues." -ForegroundColor Yellow
        exit 1
    }
}

# Run main function
Main
