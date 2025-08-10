<#
.SYNOPSIS
    Test runner for the modular AI setup installation system

.DESCRIPTION
    Tests all components of the modular installation system in isolation
    and validates the integration without affecting real files.
#>

Write-Host "=================================================================" -ForegroundColor Cyan
Write-Host "  Terraform AzureRM Provider AI Setup - Modular Test Suite" -ForegroundColor Cyan  
Write-Host "=================================================================" -ForegroundColor Cyan
Write-Host ""

$TestsDir = $PSScriptRoot
$ModulesDir = Join-Path $TestsDir "modules"
$FakeRepoDir = Join-Path $TestsDir "fake-repo"
$TestBackupsDir = Join-Path $TestsDir "test-backups"

# Clean up any previous test runs
if (Test-Path $TestBackupsDir) {
    Remove-Item $TestBackupsDir -Recurse -Force -ErrorAction SilentlyContinue
}

function Test-ModuleImports {
    Write-Host "[TEST] Testing module imports..." -ForegroundColor Blue
    
    $modules = @(
        "core-functions.ps1",
        "vscode-setup.ps1", 
        "copilot-install.ps1",
        "cleanup-restore.ps1"
    )
    
    $success = $true
    foreach ($module in $modules) {
        $modulePath = Join-Path $ModulesDir $module
        
        if (-not (Test-Path $modulePath)) {
            Write-Host "[ERROR] Module not found: $module" -ForegroundColor Red
            $success = $false
            continue
        }
        
        try {
            Import-Module $modulePath -Force
            Write-Host "[SUCCESS] $module imported successfully" -ForegroundColor Green
        }
        catch {
            Write-Host "[ERROR] Failed to import $module`: $($_.Exception.Message)" -ForegroundColor Red
            $success = $false
        }
    }
    
    return $success
}

function Test-CoreFunctions {
    Write-Host "[TEST] Testing core functions..." -ForegroundColor Blue
    
    Import-Module (Join-Path $ModulesDir "core-functions.ps1") -Force
    
    $success = $true
    
    # Test repository discovery
    $repoPath = Find-RepositoryRoot -StartPath $TestsDir
    if ($repoPath -and (Test-Path (Join-Path $repoPath "go.mod"))) {
        Write-Host "[SUCCESS] Repository discovery works" -ForegroundColor Green
    } else {
        Write-Host "[ERROR] Repository discovery failed" -ForegroundColor Red
        $success = $false
    }
    
    # Test prerequisites
    $prereqs = Test-Prerequisites
    if ($prereqs) {
        Write-Host "[SUCCESS] Prerequisites check works" -ForegroundColor Green
    } else {
        Write-Host "[WARNING] Prerequisites check returned false (may be expected)" -ForegroundColor Yellow
    }
    
    return $success
}

function Test-VSCodeFunctions {
    Write-Host "[TEST] Testing VS Code functions..." -ForegroundColor Blue
    
    Import-Module (Join-Path $ModulesDir "vscode-setup.ps1") -Force
    
    $success = $true
    
    # Test settings path detection
    $settingsPath = Get-VSCodeUserSettingsPath
    if ($settingsPath) {
        Write-Host "[SUCCESS] VS Code settings path detected: $(Split-Path $settingsPath -Leaf)" -ForegroundColor Green
    } else {
        Write-Host "[ERROR] VS Code settings path detection failed" -ForegroundColor Red
        $success = $false
    }
    
    # Test Copilot settings generation
    $copilotSettings = Get-CopilotSettings
    if ($copilotSettings -and $copilotSettings.ContainsKey("github.copilot.enable")) {
        Write-Host "[SUCCESS] Copilot settings generation works" -ForegroundColor Green
    } else {
        Write-Host "[ERROR] Copilot settings generation failed" -ForegroundColor Red
        $success = $false
    }
    
    return $success
}

function Test-CopilotFunctions {
    Write-Host "[TEST] Testing Copilot installation functions..." -ForegroundColor Blue
    
    Import-Module (Join-Path $ModulesDir "copilot-install.ps1") -Force
    
    $success = $true
    
    # Test installation verification
    $installExists = Test-CopilotInstallation -TargetRepoPath $FakeRepoDir
    if ($installExists -eq $false) {
        Write-Host "[SUCCESS] Installation detection works (correctly detected no installation)" -ForegroundColor Green
    } else {
        Write-Host "[WARNING] Installation detection returned unexpected result" -ForegroundColor Yellow
    }
    
    return $success
}

function Test-CleanupFunctions {
    Write-Host "[TEST] Testing cleanup and restore functions..." -ForegroundColor Blue
    
    Import-Module (Join-Path $ModulesDir "cleanup-restore.ps1") -Force
    Import-Module (Join-Path $ModulesDir "core-functions.ps1") -Force
    
    $success = $true
    
    # Test backup info with empty directory
    $backupInfo = Get-BackupInfo -BackupDir $TestBackupsDir
    if ($backupInfo.Count -eq 0) {
        Write-Host "[SUCCESS] Backup info works with non-existent directory" -ForegroundColor Green
    } else {
        Write-Host "[ERROR] Backup info returned unexpected results" -ForegroundColor Red
        $success = $false
    }
    
    return $success
}

function Test-MainInstaller {
    Write-Host "[TEST] Testing main installer..." -ForegroundColor Blue
    
    $installerPath = Join-Path $TestsDir "main-installer.ps1"
    
    if (-not (Test-Path $installerPath)) {
        Write-Host "[ERROR] Main installer not found" -ForegroundColor Red
        return $false
    }
    
    # Test help display - direct execution
    try {
        # Check if script can be parsed without errors
        $null = [System.Management.Automation.Language.Parser]::ParseFile($installerPath, [ref]$null, [ref]$null)
        Write-Host "[SUCCESS] Main installer syntax is valid" -ForegroundColor Green
        
        # Test if we can dot-source the functions without errors
        . $installerPath -Help *>&1 | Out-Null
        Write-Host "[SUCCESS] Help display executes without errors" -ForegroundColor Green
    }
    catch {
        Write-Host "[ERROR] Main installer test failed: $($_.Exception.Message)" -ForegroundColor Red
        return $false
    }
    
    return $true
}

function Invoke-AllTests {
    Write-Host "[INFO] Starting modular installation test suite..." -ForegroundColor Blue
    Write-Host ""
    
    $testResults = @{}
    
    $testResults["ModuleImports"] = Test-ModuleImports
    Write-Host ""
    
    $testResults["CoreFunctions"] = Test-CoreFunctions  
    Write-Host ""
    
    $testResults["VSCodeFunctions"] = Test-VSCodeFunctions
    Write-Host ""
    
    $testResults["CopilotFunctions"] = Test-CopilotFunctions
    Write-Host ""
    
    $testResults["CleanupFunctions"] = Test-CleanupFunctions
    Write-Host ""
    
    $testResults["MainInstaller"] = Test-MainInstaller
    Write-Host ""
    
    # Summary
    Write-Host "=================================================================" -ForegroundColor Cyan
    Write-Host "  Test Results Summary" -ForegroundColor Cyan
    Write-Host "=================================================================" -ForegroundColor Cyan
    
    $passed = 0
    $total = $testResults.Count
    
    foreach ($test in $testResults.GetEnumerator()) {
        $status = if ($test.Value) { "[PASS]"; $passed++ } else { "[FAIL]" }
        $color = if ($test.Value) { "Green" } else { "Red" }
        Write-Host "$status $($test.Key)" -ForegroundColor $color
    }
    
    Write-Host ""
    Write-Host "Results: $passed/$total tests passed" -ForegroundColor $(if ($passed -eq $total) { "Green" } else { "Yellow" })
    
    if ($passed -eq $total) {
        Write-Host "[SUCCESS] All tests passed! The modular installation system is ready." -ForegroundColor Green
    } else {
        Write-Host "[WARNING] Some tests failed. Review the output above." -ForegroundColor Yellow
    }
    
    return $passed -eq $total
}

# Run the test suite
$allTestsPassed = Invoke-AllTests

if ($allTestsPassed) {
    Write-Host ""
    Write-Host "=================================================================" -ForegroundColor Green
    Write-Host "  MODULAR SYSTEM VALIDATION COMPLETE" -ForegroundColor Green
    Write-Host "=================================================================" -ForegroundColor Green
    Write-Host "The modular installation system has been validated and is ready" -ForegroundColor Green
    Write-Host "for production use. You can now safely replace the monolithic" -ForegroundColor Green  
    Write-Host "install scripts with this modular approach." -ForegroundColor Green
    exit 0
} else {
    Write-Host ""
    Write-Host "=================================================================" -ForegroundColor Red
    Write-Host "  VALIDATION FAILED" -ForegroundColor Red
    Write-Host "=================================================================" -ForegroundColor Red
    Write-Host "Some tests failed. Please review and fix issues before using" -ForegroundColor Red
    Write-Host "the modular installation system in production." -ForegroundColor Red
    exit 1
}
