# Test script to mock different clean operation failure scenarios
# This script temporarily modifies validation results to test UI output

param(
    [Parameter(Mandatory)]
    [ValidateSet("WorkspaceInvalid", "SystemRequirements", "Internet", "PowerShell", "ExecutionPolicy")]
    [string]$FailureType
)

# Import the modules
Import-Module ".\modules\powershell\ValidationEngine.psm1" -Force
Import-Module ".\modules\powershell\FileOperations.psm1" -Force

# Mock the validation function based on failure type
function New-MockValidationResult {
    param([string]$Type)
    
    switch ($Type) {
        "WorkspaceInvalid" {
            return @{
                OverallValid = $false
                Git = @{
                    Valid = $true
                    CurrentBranch = "feature/test-branch"
                    Reason = "Git validation passed"
                }
                Workspace = @{
                    Valid = $false
                    Reason = "Invalid workspace: go.mod file not found in workspace root. Ensure you are in the terraform-provider-azurerm directory."
                }
                SystemRequirements = @{
                    OverallValid = $true
                    PowerShell = @{ Valid = $true }
                    ExecutionPolicy = @{ Valid = $true }
                    Commands = @{ Valid = $true }
                    Internet = @{ Connected = $true }
                }
            }
        }
        "SystemRequirements" {
            return @{
                OverallValid = $false
                Git = @{
                    Valid = $true
                    CurrentBranch = "feature/test-branch"
                    Reason = "Git validation passed"
                }
                Workspace = @{
                    Valid = $true
                    Reason = "Workspace validation passed"
                }
                SystemRequirements = @{
                    OverallValid = $false
                    PowerShell = @{ Valid = $false; Reason = "PowerShell 5.1 or higher required. Current version: 4.0" }
                    ExecutionPolicy = @{ Valid = $false; Reason = "Execution policy must be RemoteSigned or Unrestricted" }
                    Commands = @{ Valid = $true }
                    Internet = @{ Connected = $true }
                }
            }
        }
        "Internet" {
            return @{
                OverallValid = $false
                Git = @{
                    Valid = $true
                    CurrentBranch = "feature/test-branch"
                    Reason = "Git validation passed"
                }
                Workspace = @{
                    Valid = $true
                    Reason = "Workspace validation passed"
                }
                SystemRequirements = @{
                    OverallValid = $false
                    PowerShell = @{ Valid = $true }
                    ExecutionPolicy = @{ Valid = $true }
                    Commands = @{ Valid = $true }
                    Internet = @{ Connected = $false; Reason = "Cannot connect to GitHub. Check internet connection and firewall settings." }
                }
            }
        }
    }
}

# Create mock validation result
$mockValidation = New-MockValidationResult -Type $FailureType

Write-Host "============================================================" -ForegroundColor Green
Write-Host " Testing Clean Operation - $FailureType Failure Scenario" -ForegroundColor Green
Write-Host "============================================================" -ForegroundColor Green
Write-Host ""

# Show what the clean operation would display
if (-not $mockValidation.OverallValid) {
    $errorMessages = @()
    
    if (-not $mockValidation.Workspace.Valid) {
        $errorMessages += $mockValidation.Workspace.Reason
    }
    
    if (-not $mockValidation.SystemRequirements.OverallValid) {
        if (-not $mockValidation.SystemRequirements.PowerShell.Valid) {
            $errorMessages += $mockValidation.SystemRequirements.PowerShell.Reason
        }
        if (-not $mockValidation.SystemRequirements.ExecutionPolicy.Valid) {
            $errorMessages += $mockValidation.SystemRequirements.ExecutionPolicy.Reason
        }
        if (-not $mockValidation.SystemRequirements.Internet.Connected) {
            $errorMessages += $mockValidation.SystemRequirements.Internet.Reason
        }
    }
    
    Write-Host "Validating cleanup prerequisites..." -ForegroundColor Yellow
    Write-Host ""
    Write-Host "To Clean a Target Repository:" -ForegroundColor Cyan
    Write-Host "  1. Switch to a feature branch" -ForegroundColor White
    Write-Host "  2. Or run this command on a cloned repository" -ForegroundColor White
    Write-Host "  3. Or use -RepoDirectory to specify target directory" -ForegroundColor White
    Write-Host ""
    Write-Host "Clean Operation Encountered Issues:" -ForegroundColor Cyan
    foreach ($errorMsg in $errorMessages) {
        Write-Host "  - $errorMsg" -ForegroundColor Red
    }
    Write-Host ""
}
