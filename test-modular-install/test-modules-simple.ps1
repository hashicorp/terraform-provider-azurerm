# test-modules-simple.ps1 - Simple test suite for PowerShell modules

[CmdletBinding()]
param(
    [switch]$ShowDetails,
    [switch]$ShowHelp
)

if ($ShowHelp) {
    Write-Host @"
PowerShell Module Test Suite for Terraform AzureRM Provider AI Setup

USAGE:
  .\test-modules-simple.ps1 [OPTIONS]

OPTIONS:
  -ShowDetails Show detailed test output
  -ShowHelp    Show this help message

"@ -ForegroundColor Green
    return
}

# Initialize counters
$script:PassedTests = 0
$script:FailedTests = 0

function Write-TestResult {
    param(
        [string]$TestName,
        [bool]$Success,
        [string]$Details = ""
    )
    
    $status = if ($Success) { "PASS" } else { "FAIL" }
    $color = if ($Success) { "Green" } else { "Red" }
    
    Write-Host "[$status] $TestName" -ForegroundColor $color
    if ($Details -and ($ShowDetails -or -not $Success)) {
        Write-Host "       $Details" -ForegroundColor Gray
    }
    
    if ($Success) {
        $script:PassedTests++
    } else {
        $script:FailedTests++
    }
}

Write-Host "=================================================================" -ForegroundColor Cyan
Write-Host "  PowerShell Module Test Suite for Terraform AzureRM AI Setup" -ForegroundColor Cyan
Write-Host "=================================================================" -ForegroundColor Cyan
Write-Host ""

$ModulesPath = Join-Path $PSScriptRoot "modules"

# Test 1: Check if modules directory exists
$modulesExist = Test-Path $ModulesPath
Write-TestResult "Modules directory exists" $modulesExist "Path: $ModulesPath"

if (-not $modulesExist) {
    Write-Host "Cannot continue without modules directory" -ForegroundColor Red
    return
}

# Test 2: Check individual module files
$moduleFiles = @(
    "CoreFunctions.psm1",
    "VSCodeSetup.psm1", 
    "CopilotInstall.psm1",
    "TerraformAzureRMSetup.psm1"
)

Write-Host ""
Write-Host "[INFO] Testing individual modules..." -ForegroundColor Blue

foreach ($moduleFile in $moduleFiles) {
    $modulePath = Join-Path $ModulesPath $moduleFile
    $exists = Test-Path $modulePath
    Write-TestResult "Module file exists: $moduleFile" $exists "Path: $modulePath"
    
    if ($exists) {
        # Test if module can be parsed
        try {
            $null = [System.Management.Automation.Language.Parser]::ParseFile($modulePath, [ref]$null, [ref]$null)
            Write-TestResult "Module syntax valid: $moduleFile" $true
        }
        catch {
            Write-TestResult "Module syntax valid: $moduleFile" $false $_.Exception.Message
            continue
        }
        
        # Test if module can be imported
        try {
            $moduleName = [System.IO.Path]::GetFileNameWithoutExtension($moduleFile)
            Import-Module $modulePath -Force -ErrorAction Stop
            Write-TestResult "Module imports: $moduleFile" $true
            
            # Check exported functions
            $exportedFunctions = Get-Command -Module $moduleName -CommandType Function -ErrorAction SilentlyContinue
            $functionCount = if ($exportedFunctions) { $exportedFunctions.Count } else { 0 }
            Write-TestResult "Module exports functions: $moduleFile" ($functionCount -gt 0) "Exported: $functionCount functions"
            
            Remove-Module $moduleName -Force -ErrorAction SilentlyContinue
        }
        catch {
            Write-TestResult "Module imports: $moduleFile" $false $_.Exception.Message
        }
    }
}

# Test 3: Check module manifest
Write-Host ""
Write-Host "[INFO] Testing module manifest..." -ForegroundColor Blue

$manifestPath = Join-Path $ModulesPath "TerraformAzureRMSetup.psd1"
$manifestExists = Test-Path $manifestPath
Write-TestResult "Module manifest exists" $manifestExists "Path: $manifestPath"

if ($manifestExists) {
    try {
        $manifest = Test-ModuleManifest $manifestPath -ErrorAction Stop
        Write-TestResult "Module manifest is valid" $true "Version: $($manifest.Version)"
        
        # Check required fields
        $manifestData = Import-PowerShellDataFile $manifestPath
        $requiredFields = @('ModuleVersion', 'Author', 'Description')
        
        foreach ($field in $requiredFields) {
            $hasField = $manifestData.ContainsKey($field) -and $manifestData[$field]
            Write-TestResult "Manifest has $field" $hasField
        }
        
    }
    catch {
        Write-TestResult "Module manifest is valid" $false $_.Exception.Message
    }
}

# Test 4: Test main module functionality
Write-Host ""
Write-Host "[INFO] Testing main module functionality..." -ForegroundColor Blue

try {
    Import-Module $manifestPath -Force -ErrorAction Stop
    
    # Test key functions
    $testFunctions = @(
        "Find-RepositoryRoot",
        "Test-Prerequisites", 
        "Get-VSCodeUserSettingsPath",
        "Test-CopilotInstallation",
        "Show-InstallationStatus"
    )
    
    foreach ($funcName in $testFunctions) {
        $command = Get-Command $funcName -ErrorAction SilentlyContinue
        $exists = $null -ne $command
        Write-TestResult "Function available: $funcName" $exists
        
        if ($exists) {
            # Test function help
            $help = Get-Help $funcName -ErrorAction SilentlyContinue
            $hasHelp = $help.Synopsis -and $help.Synopsis -ne $funcName
            Write-TestResult "Function has help: $funcName" $hasHelp
        }
    }
    
    # Test Find-RepositoryRoot functionality
    try {
        $repoPath = Find-RepositoryRoot -StartPath $PSScriptRoot
        $foundRepo = $null -ne $repoPath
        Write-TestResult "Find-RepositoryRoot works" $foundRepo "Found: $repoPath"
    }
    catch {
        Write-TestResult "Find-RepositoryRoot works" $false $_.Exception.Message
    }
    
    Remove-Module TerraformAzureRMSetup -Force -ErrorAction SilentlyContinue
}
catch {
    Write-TestResult "Main module imports" $false $_.Exception.Message
}

# Display results
Write-Host ""
Write-Host "=================================================================" -ForegroundColor Cyan
Write-Host "  Test Results Summary" -ForegroundColor Cyan
Write-Host "=================================================================" -ForegroundColor Cyan

$totalTests = $script:PassedTests + $script:FailedTests
$passRate = if ($totalTests -gt 0) { [math]::Round(($script:PassedTests / $totalTests) * 100, 1) } else { 0 }

Write-Host "Total Tests: $totalTests" -ForegroundColor White
Write-Host "Passed: $($script:PassedTests)" -ForegroundColor Green  
Write-Host "Failed: $($script:FailedTests)" -ForegroundColor Red
Write-Host "Pass Rate: $passRate%" -ForegroundColor $(if ($passRate -ge 90) { "Green" } elseif ($passRate -ge 70) { "Yellow" } else { "Red" })

Write-Host ""
$overallStatus = if ($script:FailedTests -eq 0) { "SUCCESS" } else { "FAILED" }
$statusColor = if ($script:FailedTests -eq 0) { "Green" } else { "Red" }
Write-Host "[$overallStatus] PowerShell module test suite completed" -ForegroundColor $statusColor
