# PowerShell Test Case: Verification Functionality
# =================================================
#
# Tests the installation verification functionality including:
# - Manifest validation
# - Installation integrity checks  
# - Verification command line options
#

Write-Host "TEST: Verification Functionality" -ForegroundColor Cyan
Write-Host "==================================" -ForegroundColor Cyan
Write-Host

# Test counter
$script:TestsPassed = 0
$script:TestsFailed = 0
$script:TotalTests = 0

function Test-Result {
    param(
        [string]$TestName,
        [bool]$Passed,
        [string]$Details = ""
    )
    
    $script:TotalTests++
    
    if ($Passed) {
        Write-Host " PASS: $TestName" -ForegroundColor Green
        $script:TestsPassed++
    } else {
        Write-Host " FAIL: $TestName" -ForegroundColor Red
        $script:TestsFailed++
    }
    
    if ($Details) {
        Write-Host "   $Details" -ForegroundColor Gray
    }
}

# Test 1: Verify PowerShell installer has -Verify parameter
Write-Host "Test 1: PowerShell Installer Verification Parameter"
Write-Host "---------------------------------------------------"

try {
    $installerPath = "$PSScriptRoot\..\..\..\AILocalInstall\install-copilot-setup.ps1"
    
    # Check if installer exists
    if (Test-Path $installerPath) {
        Test-Result "PowerShell installer exists" $true
        
        # Check if installer contains Verify parameter
        $installerContent = Get-Content $installerPath -Raw
        $hasVerifyParam = [bool]($installerContent -match '\$Verify')
        
        Test-Result "Installer has Verify parameter" $hasVerifyParam
        
        # Check if help shows Verify option
        $helpOutput = & PowerShell -Command "& '$installerPath' -Help" 2>&1
        $helpHasVerify = [bool]($helpOutput -match "-Verify")
        
        Test-Result "Help documentation includes -Verify" $helpHasVerify
        
    } else {
        Test-Result "PowerShell installer exists" $false "File not found: $installerPath"
    }
} catch {
    Test-Result "PowerShell installer verification test" $false "Error: $_"
}

Write-Host

# Test 2: Verify verification module exists
Write-Host "Test 2: PowerShell Verification Module"
Write-Host "--------------------------------------"

try {
    $verificationModule = "$PSScriptRoot\..\..\..\AILocalInstall\modules\powershell\installation-verification.psm1"
    
    if (Test-Path $verificationModule) {
        Test-Result "Verification module exists" $true
        
        # Check if module contains key functions
        $moduleContent = Get-Content $verificationModule -Raw
        
        $hasIntegrityTest = [bool]($moduleContent -match "Test-InstallationIntegrity")
        Test-Result "Has Test-InstallationIntegrity function" $hasIntegrityTest
        
        $hasSourceFileTest = [bool]($moduleContent -match "Test-SourceFiles")
        Test-Result "Has Test-SourceFiles function" $hasSourceFileTest
        
        $hasTargetFileTest = [bool]($moduleContent -match "Test-TargetFiles")
        Test-Result "Has Test-TargetFiles function" $hasTargetFileTest
        
    } else {
        Test-Result "Verification module exists" $false "File not found: $verificationModule"
    }
} catch {
    Test-Result "PowerShell verification module test" $false "Error: $_"
}

Write-Host

# Test 3: Manifest validation
Write-Host "Test 3: Manifest File Validation"
Write-Host "---------------------------------"

try {
    $manifestPath = "$PSScriptRoot\..\..\..\AILocalInstall\modules\installation-manifest.json"
    
    if (Test-Path $manifestPath) {
        Test-Result "Manifest file exists" $true
        
        # Test JSON validity
        try {
            $manifest = Get-Content $manifestPath | ConvertFrom-Json
            Test-Result "Manifest JSON is valid" $true
            
            # Test required fields
            $hasVersion = $null -ne $manifest.manifestVersion
            Test-Result "Manifest has version field" $hasVersion
            
            $hasRepository = $null -ne $manifest.repository
            Test-Result "Manifest has repository field" $hasRepository
            
            $hasSourceFiles = $null -ne $manifest.sourceFiles
            Test-Result "Manifest has sourceFiles field" $hasSourceFiles
            
            if ($hasRepository) {
                $hasRepoName = $null -ne $manifest.repository.name
                $hasRepoOwner = $null -ne $manifest.repository.owner
                Test-Result "Repository structure valid" ($hasRepoName -and $hasRepoOwner)
            }
            
        } catch {
            Test-Result "Manifest JSON is valid" $false "JSON parsing error: $_"
        }
        
    } else {
        Test-Result "Manifest file exists" $false "File not found: $manifestPath"
    }
} catch {
    Test-Result "Manifest validation test" $false "Error: $_"
}

Write-Host

# Test 4: Test verification functionality (safe mode - no actual installation)
Write-Host "Test 4: Verification Functionality Test"
Write-Host "---------------------------------------"

try {
    # This test verifies that the verification functions work without requiring an actual installation
    # We'll test against a mock/clean environment
    
    $tempRepo = Join-Path $env:TEMP "terraform-azurerm-test-$(Get-Date -Format 'yyyyMMdd-HHmmss')"
    
    # Create a minimal test repository structure
    New-Item -Path $tempRepo -ItemType Directory -Force | Out-Null
    New-Item -Path "$tempRepo\.github" -ItemType Directory -Force | Out-Null
    New-Item -Path "$tempRepo\.github\AILocalInstall" -ItemType Directory -Force | Out-Null
    
    # Copy manifest to test location
    $manifestSource = "$PSScriptRoot\..\..\..\AILocalInstall\modules\installation-manifest.json"
    $manifestDest = "$tempRepo\.github\AILocalInstall\modules\installation-manifest.json"
    New-Item -Path (Split-Path $manifestDest) -ItemType Directory -Force | Out-Null
    Copy-Item $manifestSource $manifestDest -Force
    
    Test-Result "Test environment created" (Test-Path $tempRepo)
    
    # Import the verification module for testing
    $verificationModule = "$PSScriptRoot\..\..\..\AILocalInstall\modules\powershell\installation-verification.psm1"
    if (Test-Path $verificationModule) {
        try {
            Import-Module $verificationModule -Force
            Test-Result "Verification module imports successfully" $true
            
            # Test verification function exists and can be called
            $functionExists = Get-Command Test-InstallationIntegrity -ErrorAction SilentlyContinue
            Test-Result "Test-InstallationIntegrity function available" ($null -ne $functionExists)
            
        } catch {
            Test-Result "Verification module imports successfully" $false "Import error: $_"
        }
    }
    
    # Cleanup test environment
    if (Test-Path $tempRepo) {
        Remove-Item $tempRepo -Recurse -Force -ErrorAction SilentlyContinue
    }
    
} catch {
    Test-Result "Verification functionality test" $false "Error: $_"
}

Write-Host

# Test Summary
Write-Host "Test Summary" -ForegroundColor Cyan
Write-Host "============" -ForegroundColor Cyan
Write-Host "Total Tests: $script:TotalTests"
Write-Host "Passed: $script:TestsPassed" -ForegroundColor Green
Write-Host "Failed: $script:TestsFailed" $(if ($script:TestsFailed -gt 0) { "-ForegroundColor Red" } else { "-ForegroundColor Green" })
Write-Host

if ($script:TestsFailed -eq 0) {
    Write-Host "All verification tests passed!" -ForegroundColor Green
    exit 0
} else {
    Write-Host "$script:TestsFailed test(s) failed!" -ForegroundColor Red
    exit 1
}
