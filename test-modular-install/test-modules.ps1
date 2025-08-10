# test-modules.ps1 - Comprehensive test suite for PowerShell modules

[CmdletBinding()]
param(
    [switch]$ShowDetails,
    [switch]$ShowHelp
)

if ($ShowHelp) {
    Write-Host @"
PowerShell Module Test Suite for Terraform AzureRM Provider AI Setup

USAGE:
  .\test-modules.ps1 [OPTIONS]

OPTIONS:
  -ShowDetails Show detailed test output
  -ShowHelp    Show this help message

DESCRIPTION:
  Tests all PowerShell modules (.psm1) for proper functionality,
  help documentation, and module structure compliance.

"@ -ForegroundColor Green
    return
}

# Test configuration
$ModulesPath = Join-Path $PSScriptRoot "modules"
$TestResults = @{
    Passed = 0
    Failed = 0
    Tests = @()
}

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
    
    $TestResults.Tests += @{
        Name = $TestName
        Success = $Success
        Details = $Details
    }
    
    if ($Success) {
        $TestResults.Passed++
    } else {
        $TestResults.Failed++
    }
}

function Test-ModuleStructure {
    param([string]$ModulePath)
    
    $moduleName = [System.IO.Path]::GetFileNameWithoutExtension($ModulePath)
    
    # Test 1: Module file exists
    $exists = Test-Path $ModulePath
    Write-TestResult "Module file exists: $moduleName" $exists "Path: $ModulePath"
    
    if (-not $exists) { return }
    
    # Test 2: Module can be parsed
    try {
        $null = [System.Management.Automation.Language.Parser]::ParseFile($ModulePath, [ref]$null, [ref]$null)
        Write-TestResult "Module syntax valid: $moduleName" $true
    }
    catch {
        Write-TestResult "Module syntax valid: $moduleName" $false $_.Exception.Message
        return
    }
    
    # Test 3: Module can be imported
    try {
        Import-Module $ModulePath -Force -ErrorAction Stop
        Write-TestResult "Module imports successfully: $moduleName" $true
        
        # Test 4: Check exported functions
        $exportedFunctions = Get-Command -Module $moduleName -CommandType Function
        $functionCount = $exportedFunctions.Count
        Write-TestResult "Module exports functions: $moduleName" ($functionCount -gt 0) "Exported: $functionCount functions"
        
        # Test 5: Test function help
        $functionsWithHelp = 0
        foreach ($func in $exportedFunctions) {
            $help = Get-Help $func.Name -ErrorAction SilentlyContinue
            if ($help.Synopsis -and $help.Synopsis -ne $func.Name) {
                $functionsWithHelp++
            }
        }
        
        $hasHelp = $functionsWithHelp -eq $functionCount
        Write-TestResult "All functions have help: $moduleName" $hasHelp "Help: $functionsWithHelp/$functionCount functions"
        
        Remove-Module $moduleName -Force -ErrorAction SilentlyContinue
    }
    catch {
        Write-TestResult "Module imports successfully: $moduleName" $false $_.Exception.Message
    }
}

function Test-ModuleManifest {
    $manifestPath = Join-Path $ModulesPath "TerraformAzureRMSetup.psd1"
    
    # Test 1: Manifest exists
    $exists = Test-Path $manifestPath
    Write-TestResult "Module manifest exists" $exists "Path: $manifestPath"
    
    if (-not $exists) { return }
    
    # Test 2: Manifest is valid
    try {
        $manifest = Test-ModuleManifest $manifestPath -ErrorAction Stop
        Write-TestResult "Module manifest is valid" $true "Version: $($manifest.Version)"
        
        # Test 3: Required fields present
        $requiredFields = @('ModuleVersion', 'Author', 'Description', 'FunctionsToExport')
        $manifestData = Import-PowerShellDataFile $manifestPath
        
        foreach ($field in $requiredFields) {
            $hasField = $manifestData.ContainsKey($field) -and $manifestData[$field]
            Write-TestResult "Manifest has $field" $hasField
        }
        
        # Test 4: Functions match exports
        $exportedInManifest = $manifestData.FunctionsToExport
        Write-TestResult "Manifest exports functions" ($exportedInManifest.Count -gt 0) "Exports: $($exportedInManifest.Count) functions"
        
    }
    catch {
        Write-TestResult "Module manifest is valid" $false $_.Exception.Message
    }
}

function Test-FunctionFunctionality {
    try {
        # Import main module
        $mainModulePath = Join-Path $ModulesPath "TerraformAzureRMSetup.psd1"
        Import-Module $mainModulePath -Force
        
        # Test Find-RepositoryRoot
        try {
            $repoPath = Find-RepositoryRoot -StartPath $PSScriptRoot
            $foundRepo = $null -ne $repoPath
            Write-TestResult "Find-RepositoryRoot works" $foundRepo "Found: $repoPath"
        }
        catch {
            Write-TestResult "Find-RepositoryRoot works" $false $_.Exception.Message
        }
        
        # Test Test-Prerequisites
        try {
            $prereqResult = Test-Prerequisites
            Write-TestResult "Test-Prerequisites works" $true "Result: $prereqResult"
        }
        catch {
            Write-TestResult "Test-Prerequisites works" $false $_.Exception.Message
        }
        
        # Test Get-VSCodeUserSettingsPath
        try {
            $vsCodePath = Get-VSCodeUserSettingsPath
            $pathValid = -not [string]::IsNullOrEmpty($vsCodePath)
            Write-TestResult "Get-VSCodeUserSettingsPath works" $pathValid "Path: $vsCodePath"
        }
        catch {
            Write-TestResult "Get-VSCodeUserSettingsPath works" $false $_.Exception.Message
        }
        
        # Test Test-CopilotInstallation
        try {
            $testPath = $PSScriptRoot
            $copilotTest = Test-CopilotInstallation -RepositoryPath $testPath
            Write-TestResult "Test-CopilotInstallation works" $true "Result: $copilotTest"
        }
        catch {
            Write-TestResult "Test-CopilotInstallation works" $false $_.Exception.Message
        }
        
        # Test Show-InstallationStatus
        try {
            $null = Show-InstallationStatus -RepositoryPath $PSScriptRoot 2>&1
            Write-TestResult "Show-InstallationStatus works" $true "Output captured"
        }
        catch {
            Write-TestResult "Show-InstallationStatus works" $false $_.Exception.Message
        }
        
        Remove-Module TerraformAzureRMSetup -Force -ErrorAction SilentlyContinue
    }
    catch {
        Write-TestResult "Module functionality test" $false $_.Exception.Message
    }
}

function Test-HelperFunctions {
    try {
        # Import CoreFunctions module directly
        $coreFunctionsPath = Join-Path $ModulesPath "CoreFunctions.psm1"
        Import-Module $coreFunctionsPath -Force
        
        # Test New-SafeBackup
        try {
            $testFile = Join-Path $PSScriptRoot "test-file.txt"
            $backupDir = Join-Path $PSScriptRoot "test-backup"
            
            # Create test file
            "Test content" | Out-File $testFile -Encoding UTF8
            
            if (Test-Path $testFile) {
                $backupPath = New-SafeBackup -SourcePath $testFile -BackupDirectory $backupDir
                $backupCreated = Test-Path $backupPath
                Write-TestResult "New-SafeBackup works" $backupCreated "Backup: $backupPath"
                
                # Cleanup
                Remove-Item $testFile -Force -ErrorAction SilentlyContinue
                Remove-Item $backupDir -Recurse -Force -ErrorAction SilentlyContinue
            }
        }
        catch {
            Write-TestResult "New-SafeBackup works" $false $_.Exception.Message
        }
        
        # Test Test-FileIntegrity
        try {
            $testFile = Join-Path $PSScriptRoot "test-integrity.txt"
            "Test content" | Out-File $testFile -Encoding UTF8
            
            $integrityResult = Test-FileIntegrity -Path $testFile
            Write-TestResult "Test-FileIntegrity works" $integrityResult "File is readable"
            
            # Test with non-existent file
            $nonExistentResult = Test-FileIntegrity -Path "non-existent-file.txt"
            Write-TestResult "Test-FileIntegrity handles missing files" (-not $nonExistentResult) "Correctly returns false"
            
            # Cleanup
            Remove-Item $testFile -Force -ErrorAction SilentlyContinue
        }
        catch {
            Write-TestResult "Test-FileIntegrity works" $false $_.Exception.Message
        }
        
        Remove-Module CoreFunctions -Force -ErrorAction SilentlyContinue
    }
    catch {
        Write-TestResult "Helper functions test" $false $_.Exception.Message
    }
}

# Main test execution
Write-Host "=================================================================" -ForegroundColor Cyan
Write-Host "  PowerShell Module Test Suite for Terraform AzureRM AI Setup" -ForegroundColor Cyan
Write-Host "=================================================================" -ForegroundColor Cyan
Write-Host ""

# Test individual modules
$moduleFiles = @(
    "CoreFunctions.psm1",
    "VSCodeSetup.psm1", 
    "CopilotInstall.psm1",
    "TerraformAzureRMSetup.psm1"
)

Write-Host "[INFO] Testing individual modules..." -ForegroundColor Blue
foreach ($moduleFile in $moduleFiles) {
    $modulePath = Join-Path $ModulesPath $moduleFile
    Test-ModuleStructure $modulePath
}

Write-Host ""
Write-Host "[INFO] Testing module manifest..." -ForegroundColor Blue
Test-ModuleManifest

Write-Host ""
Write-Host "[INFO] Testing function functionality..." -ForegroundColor Blue
Test-FunctionFunctionality

Write-Host ""
Write-Host "[INFO] Testing helper functions..." -ForegroundColor Blue
Test-HelperFunctions

# Display results
Write-Host ""
Write-Host "=================================================================" -ForegroundColor Cyan
Write-Host "  Test Results Summary" -ForegroundColor Cyan
Write-Host "=================================================================" -ForegroundColor Cyan

$totalTests = $TestResults.Passed + $TestResults.Failed
$passRate = if ($totalTests -gt 0) { [math]::Round(($TestResults.Passed / $totalTests) * 100, 1) } else { 0 }

Write-Host "Total Tests: $totalTests" -ForegroundColor White
Write-Host "Passed: $($TestResults.Passed)" -ForegroundColor Green  
Write-Host "Failed: $($TestResults.Failed)" -ForegroundColor Red
Write-Host "Pass Rate: $passRate%" -ForegroundColor $(if ($passRate -ge 90) { "Green" } elseif ($passRate -ge 70) { "Yellow" } else { "Red" })

if ($TestResults.Failed -gt 0) {
    Write-Host ""
    Write-Host "Failed Tests:" -ForegroundColor Red
    foreach ($test in $TestResults.Tests) {
        if (-not $test.Success) {
            Write-Host "  - $($test.Name): $($test.Details)" -ForegroundColor Red
        }
    }
}

Write-Host ""
$overallStatus = if ($TestResults.Failed -eq 0) { "SUCCESS" } else { "FAILED" }
$statusColor = if ($TestResults.Failed -eq 0) { "Green" } else { "Red" }
Write-Host "[$overallStatus] PowerShell module test suite completed" -ForegroundColor $statusColor
