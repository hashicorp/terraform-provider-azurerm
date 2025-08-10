# AI Local Install - Comprehensive Test Runner
# Tests all aspects of the install-copilot-setup.ps1 script

param(
    [string]$Category = "All",
    [switch]$CleanupAfter = $true,
    [switch]$Verbose = $false,
    [switch]$DryRun = $false,
    [switch]$Help = $fa    try {
        # Test invalid path handling (with isolated script)
        $scriptPath = Get-IsolatedScriptPath
        $result = & PowerShell -ExecutionPolicy Bypass -File $scriptPath -RepositoryPath "C:\NonExistent\Path" 2>&1
        $invalidPathHandled = ($result -join " ") -like "*does not exist*"
        
        Write-TestResult "Invalid Repository Path" $invalidPathHandled "Error handled: $invalidPathHandled"
        
    } catch {
        Write-TestResult "Invalid Repository Path" $false $_.Exception.Message
    }
    
    try {
        # Test invalid repository structure (with isolated script)
        $tempInvalidRepo = "$TestBase\InvalidRepo-$(Get-Random)"
        New-Item -Path $tempInvalidRepo -ItemType Directory -Force | Out-Null
        Set-Content -Path "$tempInvalidRepo\README.md" -Value "# Not a terraform repo" -Force
        
        $scriptPath = Get-IsolatedScriptPath
        $result = & PowerShell -ExecutionPolicy Bypass -File $scriptPath -RepositoryPath $tempInvalidRepo 2>&1p) {
    Write-Host "AI Local Install Test Suite" -ForegroundColor Cyan
    Write-Host "===========================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "USAGE:" -ForegroundColor Yellow
    Write-Host "  .\run-all-tests.ps1 [-Category <string>] [-CleanupAfter] [-Verbose] [-DryRun] [-Help]" -ForegroundColor White
    Write-Host ""
    Write-Host "PARAMETERS:" -ForegroundColor Yellow
    Write-Host "  -Category        Test category to run (All, Install, Cleanup, EdgeCases)" -ForegroundColor White
    Write-Host "  -CleanupAfter    Clean up test directories after completion (default: true)" -ForegroundColor White
    Write-Host "  -Verbose         Show detailed output during tests" -ForegroundColor White
    Write-Host "  -DryRun          Show what tests would run without executing them" -ForegroundColor White  
    Write-Host "  -Help            Show this help message" -ForegroundColor White
    Write-Host ""
    Write-Host "EXAMPLES:" -ForegroundColor Yellow
    Write-Host "  .\run-all-tests.ps1                    # Run all tests with cleanup" -ForegroundColor Gray
    Write-Host "  .\run-all-tests.ps1 -DryRun            # Show test plan without execution" -ForegroundColor Gray
    Write-Host "  .\run-all-tests.ps1 -Category Install  # Run only installation tests" -ForegroundColor Gray
    return
}

if ($DryRun) {
    Write-Host "DRY RUN MODE - Tests will be listed but not executed" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "PLANNED TESTS:" -ForegroundColor Cyan
    Write-Host "1. Fresh Installation - Install on clean VS Code" -ForegroundColor White
    Write-Host "2. Merge Installation - Install with existing settings" -ForegroundColor White  
    Write-Host "3. Complex Regex Cleanup - Test surgical removal" -ForegroundColor White
    Write-Host "4. Edge Cases - Empty files, corruption, missing backups" -ForegroundColor White
    Write-Host "5. Error Handling - Invalid paths and repository structures" -ForegroundColor White
    Write-Host ""
    Write-Host "Use '.\run-all-tests.ps1' to execute these tests." -ForegroundColor Yellow
    return
}

$ErrorActionPreference = "Stop"

$ErrorActionPreference = "Stop"

# Test configuration with isolation
$TestBase = "C:\Users\$env:USERNAME\AppData\Local\Temp"
$SourceInstallDir = Resolve-Path "..\..\AILocalInstall"
$script:TestResults = @()

function Copy-InstallerToTestEnvironment {
    param([string]$TestDir)
    
    $installerDir = Join-Path $TestDir "installer"
    Copy-Item $SourceInstallDir -Destination $installerDir -Recurse -Force
    return Join-Path $installerDir "install-copilot-setup.ps1"
}

function Get-IsolatedScriptPath {
    # For tests that don't need full environment isolation but still need script isolation
    $tempDir = "$TestBase\AzureRM-Script-$(Get-Random)"
    return Copy-InstallerToTestEnvironment $tempDir
}

function Write-TestHeader {
    param([string]$Title, [string]$Description)
    
    Write-Host ""
    Write-Host "=" * 80 -ForegroundColor Cyan
    Write-Host "TEST: $Title" -ForegroundColor Yellow
    Write-Host "DESC: $Description" -ForegroundColor Gray
    Write-Host "=" * 80 -ForegroundColor Cyan
}

function Write-TestResult {
    param([string]$TestName, [bool]$Passed, [string]$Details = "")
    
    $Status = if ($Passed) { "PASS" } else { "FAIL" }
    $Color = if ($Passed) { "Green" } else { "Red" }
    
    Write-Host "RESULT: $TestName - $Status" -ForegroundColor $Color
    if ($Details) { Write-Host "DETAIL: $Details" -ForegroundColor Gray }
    
    $script:TestResults += @{
        Name = $TestName
        Passed = $Passed
        Details = $Details
    }
}

function Test-FreshInstallation {
    Write-TestHeader "Fresh Installation" "Install on clean VS Code with no existing settings"
    
    try {
        $testDir = "$TestBase\AzureRM-Test-Fresh-$(Get-Random)"
        $codeUserDir = "$testDir\Code\User"
        New-Item -Path $codeUserDir -ItemType Directory -Force | Out-Null
        
        # Copy installer to isolated test environment
        $scriptPath = Copy-InstallerToTestEnvironment $testDir
        
        # Run installation in isolated environment
        $originalAppData = $env:APPDATA
        $env:APPDATA = $testDir
        
        & PowerShell -ExecutionPolicy Bypass -File $scriptPath -Force
        
        # Verify installation
        $settingsExists = Test-Path "$codeUserDir\settings.json"
        $instructionsExist = Test-Path "$codeUserDir\instructions\terraform-azurerm"
        $copilotInstructionsExist = Test-Path "$codeUserDir\copilot-instructions.md"
        
        $env:APPDATA = $originalAppData
        
        $success = $settingsExists -and $instructionsExist -and $copilotInstructionsExist
        Write-TestResult "Fresh Installation" $success "Settings: $settingsExists, Instructions: $instructionsExist, Copilot: $copilotInstructionsExist"
        
        if ($CleanupAfter) { Remove-Item $testDir -Recurse -Force -ErrorAction SilentlyContinue }
        
    } catch {
        Write-TestResult "Fresh Installation" $false $_.Exception.Message
    }
}

function Test-MergeInstallation {
    Write-TestHeader "Merge Installation" "Install with existing user settings preserved"
    
    try {
        $testDir = "$TestBase\AzureRM-Test-Merge-$(Get-Random)"
        $codeUserDir = "$testDir\Code\User"
        New-Item -Path $codeUserDir -ItemType Directory -Force | Out-Null
        
        # Create existing settings
        $existingSettings = @{
            "editor.fontSize" = 14
            "workbench.colorTheme" = "Dark+ (default dark)"
            "user.custom.setting" = "preserve me"
        }
        $existingSettings | ConvertTo-Json | Set-Content "$codeUserDir\settings.json"
        
        # Copy installer to isolated test environment
        $scriptPath = Copy-InstallerToTestEnvironment $testDir
        
        # Run installation  
        $originalAppData = $env:APPDATA
        $env:APPDATA = $testDir
        
        & PowerShell -ExecutionPolicy Bypass -File $scriptPath -Force
        
        # Verify merge
        $newSettings = Get-Content "$codeUserDir\settings.json" | ConvertFrom-Json
        $userSettingsPreserved = $newSettings."editor.fontSize" -eq 14 -and $newSettings."user.custom.setting" -eq "preserve me"
        $azureSettingsAdded = $newSettings."github.copilot.chat.reviewSelection.instructions" -ne $null
        
        $env:APPDATA = $originalAppData
        
        $success = $userSettingsPreserved -and $azureSettingsAdded
        Write-TestResult "Merge Installation" $success "User preserved: $userSettingsPreserved, Azure added: $azureSettingsAdded"
        
        if ($CleanupAfter) { Remove-Item $testDir -Recurse -Force -ErrorAction SilentlyContinue }
        
    } catch {
        Write-TestResult "Merge Installation" $false $_.Exception.Message
    }
}

function Test-ComplexRegexCleanup {
    Write-TestHeader "Complex Regex Cleanup" "Test surgical removal of AzureRM settings while preserving user settings"
    
    try {
        $testDir = "$TestBase\AzureRM-Test-Regex-$(Get-Random)"
        $codeUserDir = "$testDir\Code\User"
        New-Item -Path "$codeUserDir\instructions\terraform-azurerm" -ItemType Directory -Force | Out-Null
        Set-Content -Path "$codeUserDir\instructions\terraform-azurerm\test.instructions.md" -Value "# Test" -Force
        Set-Content -Path "$codeUserDir\copilot-instructions.md" -Value "# Test" -Force
        
        # Create complex settings with mixed user and AzureRM content
        $complexSettings = @{
            "editor.fontSize" = 14
            "workbench.colorTheme" = "Dark+ (default dark)"
            "python.defaultInterpreter" = "/usr/bin/python3"
            "github.copilot.enable" = @{ "terminal" = $true; "*" = $true }
            "github.copilot.chat.reviewSelection.instructions" = @(
                @{ "file" = "copilot-instructions.md" }
                @{ "file" = "instructions/terraform-azurerm/test.instructions.md" }
            )
            "files.associations" = @{
                "*.azurerm.md" = "markdown"
                "*.instructions.md" = "markdown"
            }
            "// AZURERM_BACKUP_LENGTH" = "120"
            "user.custom.setting" = "should be preserved"
        }
        $complexSettings | ConvertTo-Json -Depth 10 | Set-Content "$codeUserDir\settings.json"
        
        # Create backup
        $userBackup = @{
            "editor.fontSize" = 14
            "workbench.colorTheme" = "Dark+ (default dark)"
            "python.defaultInterpreter" = "/usr/bin/python3"
            "user.custom.setting" = "should be preserved"
        }
        $userBackup | ConvertTo-Json | Set-Content "$codeUserDir\settings.json.backup.$(Get-Date -Format 'yyyyMMdd-HHmmss')"
        
        # Copy installer to isolated test environment
        $scriptPath = Copy-InstallerToTestEnvironment $testDir
        
        # Run cleanup
        $originalAppData = $env:APPDATA
        $env:APPDATA = $testDir
        
        & PowerShell -ExecutionPolicy Bypass -File $scriptPath -Clean -Force
        
        # Verify surgical cleanup
        $restoredSettings = Get-Content "$codeUserDir\settings.json" | ConvertFrom-Json
        $userSettingsPreserved = $restoredSettings."editor.fontSize" -eq 14 -and $restoredSettings."user.custom.setting" -eq "should be preserved"
        $azureSettingsRemoved = $restoredSettings."github.copilot.enable" -eq $null -and $restoredSettings."// AZURERM_BACKUP_LENGTH" -eq $null
        
        $env:APPDATA = $originalAppData
        
        $success = $userSettingsPreserved -and $azureSettingsRemoved
        Write-TestResult "Complex Regex Cleanup" $success "User preserved: $userSettingsPreserved, Azure removed: $azureSettingsRemoved"
        
        if ($CleanupAfter) { Remove-Item $testDir -Recurse -Force -ErrorAction SilentlyContinue }
        
    } catch {
        Write-TestResult "Complex Regex Cleanup" $false $_.Exception.Message
    }
}

function Test-EdgeCases {
    Write-TestHeader "Edge Cases" "Test empty files, corruption, missing backups"
    
    # Test empty settings.json
    try {
        $testDir = "$TestBase\AzureRM-Test-Empty-$(Get-Random)"
        $codeUserDir = "$testDir\Code\User"
        New-Item -Path "$codeUserDir\instructions\terraform-azurerm" -ItemType Directory -Force | Out-Null
        Set-Content -Path "$codeUserDir\instructions\terraform-azurerm\test.instructions.md" -Value "# Test" -Force
        Set-Content -Path "$codeUserDir\copilot-instructions.md" -Value "# Test" -Force
        Set-Content -Path "$codeUserDir\settings.json" -Value "" -Force
        Set-Content -Path "$codeUserDir\settings.json.backup.$(Get-Date -Format 'yyyyMMdd-HHmmss')" -Value '{"user.setting": "restored"}' -Force
        
        # Copy installer to isolated test environment
        $scriptPath = Copy-InstallerToTestEnvironment $testDir
        
        $originalAppData = $env:APPDATA
        $env:APPDATA = $testDir
        
        & PowerShell -ExecutionPolicy Bypass -File $scriptPath -Clean -Force
        
        $restored = Get-Content "$codeUserDir\settings.json" | ConvertFrom-Json
        $success = $restored."user.setting" -eq "restored"
        
        $env:APPDATA = $originalAppData
        
        Write-TestResult "Empty Settings File" $success "Restoration: $success"
        
        if ($CleanupAfter) { Remove-Item $testDir -Recurse -Force -ErrorAction SilentlyContinue }
        
    } catch {
        Write-TestResult "Empty Settings File" $false $_.Exception.Message
    }
}

function Test-ErrorHandling {
    Write-TestHeader "Error Handling" "Test invalid paths and repository structures"
    
    try {
        # Test invalid repository path
        $result = & PowerShell -ExecutionPolicy Bypass -File $ScriptPath -RepositoryPath "C:\NonExistent\Path" 2>&1
        $invalidPathHandled = ($result -join " ") -like "*does not exist*"
        
        Write-TestResult "Invalid Repository Path" $invalidPathHandled "Error handled: $invalidPathHandled"
        
    } catch {
        Write-TestResult "Invalid Repository Path" $false $_.Exception.Message
    }
    
    try {
        # Test invalid repository structure
        $tempInvalidRepo = "$TestBase\InvalidRepo-$(Get-Random)"
        New-Item -Path $tempInvalidRepo -ItemType Directory -Force | Out-Null
        Set-Content -Path "$tempInvalidRepo\README.md" -Value "# Not a terraform repo" -Force
        
        $result = & PowerShell -ExecutionPolicy Bypass -File $ScriptPath -RepositoryPath $tempInvalidRepo 2>&1
        $invalidStructureHandled = ($result -join " ") -like "*go.mod*not found*"
        
        Write-TestResult "Invalid Repository Structure" $invalidStructureHandled "Error handled: $invalidStructureHandled"
        
        if ($CleanupAfter) { Remove-Item $tempInvalidRepo -Recurse -Force -ErrorAction SilentlyContinue }
        
    } catch {
        Write-TestResult "Invalid Repository Structure" $false $_.Exception.Message
    }
}

# Main execution
Write-Host "AI Local Install - Comprehensive Test Suite" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan

$startTime = Get-Date

switch ($Category) {
    "All" {
        Test-FreshInstallation
        Test-MergeInstallation  
        Test-ComplexRegexCleanup
        Test-EdgeCases
        Test-ErrorHandling
    }
    "Installation" {
        Test-FreshInstallation
        Test-MergeInstallation
    }
    "EdgeCases" {
        Test-ComplexRegexCleanup
        Test-EdgeCases
    }
    "ErrorHandling" {
        Test-ErrorHandling
    }
    default {
        Write-Host "Invalid category. Use: All, Installation, EdgeCases, ErrorHandling" -ForegroundColor Red
        exit 1
    }
}

# Test summary
$endTime = Get-Date
$duration = $endTime - $startTime

Write-Host ""
Write-Host "=" * 80 -ForegroundColor Cyan
Write-Host "TEST SUMMARY" -ForegroundColor Yellow
Write-Host "=" * 80 -ForegroundColor Cyan

$totalTests = $TestResults.Count
$passedTests = ($TestResults | Where-Object { $_.Passed }).Count
$failedTests = $totalTests - $passedTests

Write-Host "Total Tests:  $totalTests" -ForegroundColor White
Write-Host "Passed:       $passedTests" -ForegroundColor Green
Write-Host "Failed:       $failedTests" -ForegroundColor Red
Write-Host "Duration:     $($duration.TotalSeconds) seconds" -ForegroundColor Gray

if ($failedTests -eq 0) {
    Write-Host ""
    Write-Host "ALL TESTS PASSED! Script is production ready." -ForegroundColor Green
} else {
    Write-Host ""
    Write-Host "Some tests failed. Review the results above." -ForegroundColor Red
    
    $TestResults | Where-Object { -not $_.Passed } | ForEach-Object {
        Write-Host "FAILED: $($_.Name) - $($_.Details)" -ForegroundColor Red
    }
}

Write-Host "=" * 80 -ForegroundColor Cyan
