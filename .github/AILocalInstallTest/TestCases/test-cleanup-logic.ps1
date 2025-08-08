# Test script for AzureRM Copilot cleanup logic
# This script tests various edge cases and scenarios

param(
    [switch]$Verbose,
    [switch]$KeepTestFiles
)

# Test configuration
$script:TestCounter = 0
$script:FailedTests = @()
$script:PassedTests = @()

# Create temporary test directory
$testDir = Join-Path $env:TEMP "azurerm-copilot-test-$(Get-Date -Format 'yyyyMMdd-HHmmss')"
New-Item -ItemType Directory -Path $testDir -Force | Out-Null
Write-Host "Test directory: $testDir" -ForegroundColor Cyan

function Test-ScenarioResult {
    param(
        [string]$TestName,
        [string]$Expected,
        [string]$Actual,
        [string]$Description
    )
    
    $script:TestCounter++
    
    if ($Expected -eq $Actual) {
        $script:PassedTests += $TestName
        Write-Host "PASS: $TestName" -ForegroundColor Green
        if ($Verbose) {
            Write-Host "   Expected: $Expected" -ForegroundColor Gray
            Write-Host "   Actual: $Actual" -ForegroundColor Gray
        }
    } else {
        $script:FailedTests += $TestName
        Write-Host "FAIL: $TestName" -ForegroundColor Red
        Write-Host "   Expected: $Expected" -ForegroundColor Yellow
        Write-Host "   Actual: $Actual" -ForegroundColor Yellow
        if ($Description) {
            Write-Host "   Description: $Description" -ForegroundColor Gray
        }
    }
}

function Create-TestVSCodeDir {
    param([string]$TestName)
    
    $vscodeDir = Join-Path $testDir "vscode-$TestName"
    $userDir = Join-Path $vscodeDir "User"
    New-Item -ItemType Directory -Path $userDir -Force | Out-Null
    return $userDir
}

function Test-BackupLengthZero {
    Write-Host "`nTesting Backup Length = 0 (File created by install)" -ForegroundColor Cyan
    
    $userDir = Create-TestVSCodeDir "backup-zero"
    $settingsPath = Join-Path $userDir "settings.json"
    
    # Create settings.json with backup length 0
    $settingsContent = @'
{
    "// AZURERM_BACKUP_LENGTH": "0",
    "github.copilot.chat.commitMessageGeneration.instructions": [
        {
            "text": "Test instructions"
        }
    ],
    "github.copilot.chat.reviewSelection.enabled": true,
    "files.associations": {
        "*.instructions.md": "markdown"
    }
}
'@
    Set-Content -Path $settingsPath -Value $settingsContent -Force
    
    # Run cleanup
    & "$PSScriptRoot\install-copilot-setup.ps1" -Clean -UserDirectory $userDir -Force
    
    # Verify file was removed completely
    $fileExists = Test-Path $settingsPath
    Test-ScenarioResult "Backup Length 0 - File Removal" "False" $fileExists.ToString() "File should be completely removed when backup length is 0"
}

function Test-BackupLengthMinusOne {
    Write-Host "`nTesting Backup Length = -1 (Manual merge)" -ForegroundColor Cyan
    
    $userDir = Create-TestVSCodeDir "backup-minus-one"
    $settingsPath = Join-Path $userDir "settings.json"
    
    # Create settings.json with backup length -1
    $settingsContent = @'
{
    "// AZURERM_BACKUP_LENGTH": "-1",
    "github.copilot.chat.commitMessageGeneration.instructions": [
        {
            "text": "Test instructions"
        }
    ],
    "github.copilot.chat.reviewSelection.enabled": true,
    "editor.fontSize": 14
}
'@
    Set-Content -Path $settingsPath -Value $settingsContent -Force
    
    # Capture output to verify error message
    $output = & "$PSScriptRoot\install-copilot-setup.ps1" -Clean -UserDirectory $userDir -Force 2>&1
    
    # Verify file still exists (not auto-cleaned)
    $fileExists = Test-Path $settingsPath
    Test-ScenarioResult "Backup Length -1 - File Preserved" "True" $fileExists.ToString() "File should be preserved when backup length is -1"
    
    # Verify error message about manual cleanup
    $hasManualMessage = $output -match "Manual merge detected"
    Test-ScenarioResult "Backup Length -1 - Manual Message" "True" $hasManualMessage.ToString() "Should show manual cleanup message"
}

function Test-BackupLengthPositive {
    Write-Host "`nTesting Backup Length > 0 (Original backed up)" -ForegroundColor Cyan
    
    $userDir = Create-TestVSCodeDir "backup-positive"
    $settingsPath = Join-Path $userDir "settings.json"
    $backupPath = Join-Path $userDir "settings_backup_azurerm.json"
    
    # Create original backup
    $backupContent = @'
{
    "editor.fontSize": 16,
    "workbench.colorTheme": "Dark+"
}
'@
    Set-Content -Path $backupPath -Value $backupContent -Force
    
    # Create current settings with AzureRM additions
    $settingsContent = @'
{
    "// AZURERM_BACKUP_LENGTH": "78",
    "editor.fontSize": 16,
    "workbench.colorTheme": "Dark+",
    "github.copilot.chat.commitMessageGeneration.instructions": [
        {
            "text": "Test instructions"
        }
    ],
    "github.copilot.chat.reviewSelection.enabled": true,
    "files.associations": {
        "*.instructions.md": "markdown"
    }
}
'@
    Set-Content -Path $settingsPath -Value $settingsContent -Force
    
    # Run cleanup
    & "$PSScriptRoot\install-copilot-setup.ps1" -Clean -UserDirectory $userDir -Force
    
    # Verify backup was restored
    $currentContent = Get-Content $settingsPath -Raw
    $hasOriginalSettings = $currentContent -match "workbench.colorTheme"
    $hasAzureRMSettings = $currentContent -match "github.copilot"
    
    Test-ScenarioResult "Backup Length >0 - Original Restored" "True" $hasOriginalSettings.ToString() "Original settings should be restored"
    Test-ScenarioResult "Backup Length >0 - AzureRM Removed" "False" $hasAzureRMSettings.ToString() "AzureRM settings should be removed"
    
    # Verify backup file was cleaned up
    $backupExists = Test-Path $backupPath
    Test-ScenarioResult "Backup Length >0 - Backup Cleanup" "False" $backupExists.ToString() "Backup file should be removed after successful restore"
}

function Test-ComplexArrayCleanup {
    Write-Host "`nTesting Complex Array Structure Cleanup" -ForegroundColor Cyan
    
    $userDir = Create-TestVSCodeDir "complex-array"
    $settingsPath = Join-Path $userDir "settings.json"
    
    # Create settings with complex nested arrays
    $settingsContent = @'
{
    "editor.fontSize": 14,
    "github.copilot.chat.commitMessageGeneration.instructions": [
        {
            "text": "Provide a concise and clear commit message"
        }
    ],
    "github.copilot.chat.reviewSelection.instructions": [
        {"file": "copilot-instructions.md"},
        {"file": "instructions/implementation-guide.instructions.md"},
        {"file": "instructions/azure-patterns.instructions.md"}
    ],
    "files.associations": {
        "*.instructions.md": "markdown",
        "*.prompt.md": "markdown",
        "*.js": "javascript"
    },
    "workbench.colorTheme": "Dark+"
}
'@
    Set-Content -Path $settingsPath -Value $settingsContent -Force
    
    # Run cleanup (simulate no backup scenario)
    & "$PSScriptRoot\install-copilot-setup.ps1" -Clean -UserDirectory $userDir -Force
    
    # Verify selective cleanup
    $cleanedContent = Get-Content $settingsPath -Raw
    $hasOriginalSettings = $cleanedContent -match "editor.fontSize.*14"
    $hasJSAssociation = $cleanedContent -match '"\*\.js"'
    $hasAzureRMArrays = $cleanedContent -match "github\.copilot"
    $hasAzureRMAssociations = $cleanedContent -match "\*\.instructions\.md"
    
    Test-ScenarioResult "Complex Array - Original Preserved" "True" $hasOriginalSettings.ToString() "Original settings should be preserved"
    Test-ScenarioResult "Complex Array - Non-AzureRM Preserved" "True" $hasJSAssociation.ToString() "Non-AzureRM file associations should be preserved"
    Test-ScenarioResult "Complex Array - GitHub Copilot Removed" "False" $hasAzureRMArrays.ToString() "GitHub Copilot arrays should be removed"
    Test-ScenarioResult "Complex Array - AzureRM Associations Removed" "False" $hasAzureRMAssociations.ToString() "AzureRM file associations should be removed"
}

function Test-CorruptedSettingsHandling {
    Write-Host "`nTesting Corrupted Settings Handling" -ForegroundColor Cyan
    
    $userDir = Create-TestVSCodeDir "corrupted"
    $settingsPath = Join-Path $userDir "settings.json"
    
    # Create corrupted JSON with AzureRM content
    $corruptedContent = @'
{
    "editor.fontSize": 14,
    "github.copilot.chat.reviewSelection.enabled": true,
    "missing.closing.brace": "test"
    // Missing closing brace
'@
    Set-Content -Path $settingsPath -Value $corruptedContent -Force
    
    # Run cleanup
    $output = & "$PSScriptRoot\install-copilot-setup.ps1" -Clean -UserDirectory $userDir -Force 2>&1
    
    # Verify file was removed due to corruption
    $fileExists = Test-Path $settingsPath
    $hasCorruptionMessage = $output -match "corrupted|invalid JSON"
    
    Test-ScenarioResult "Corrupted Settings - File Removed" "False" $fileExists.ToString() "Corrupted file should be removed"
    Test-ScenarioResult "Corrupted Settings - Warning Message" "True" $hasCorruptionMessage.ToString() "Should show corruption warning"
}

function Test-EmptyResultHandling {
    Write-Host "`nTesting Empty Result After Cleanup" -ForegroundColor Cyan
    
    $userDir = Create-TestVSCodeDir "empty-result"
    $settingsPath = Join-Path $userDir "settings.json"
    
    # Create settings with only AzureRM content
    $onlyAzureRMContent = @'
{
    "github.copilot.chat.reviewSelection.enabled": true,
    "github.copilot.chat.commitMessageGeneration.instructions": [
        {
            "text": "Test"
        }
    ],
    "files.associations": {
        "*.instructions.md": "markdown"
    }
}
'@
    Set-Content -Path $settingsPath -Value $onlyAzureRMContent -Force
    
    # Run cleanup
    & "$PSScriptRoot\install-copilot-setup.ps1" -Clean -UserDirectory $userDir -Force
    
    # Verify file was removed (empty after cleanup)
    $fileExists = Test-Path $settingsPath
    Test-ScenarioResult "Empty Result - File Removed" "False" $fileExists.ToString() "Empty file should be removed after cleanup"
}

function Run-AllTests {
    Write-Host "Starting AzureRM Copilot Cleanup Logic Tests" -ForegroundColor Green
    Write-Host "=================================================" -ForegroundColor Green
    
    Test-BackupLengthZero
    Test-BackupLengthMinusOne
    Test-BackupLengthPositive
    Test-ComplexArrayCleanup
    Test-CorruptedSettingsHandling
    Test-EmptyResultHandling
    
    Write-Host "`nTest Results Summary" -ForegroundColor Green
    Write-Host "======================" -ForegroundColor Green
    Write-Host "Total Tests: $script:TestCounter" -ForegroundColor Cyan
    Write-Host "Passed: $($script:PassedTests.Count)" -ForegroundColor Green
    Write-Host "Failed: $($script:FailedTests.Count)" -ForegroundColor Red
    
    if ($script:FailedTests.Count -gt 0) {
        Write-Host "`nFailed Tests:" -ForegroundColor Red
        $script:FailedTests | ForEach-Object { Write-Host "  - $_" -ForegroundColor Yellow }
    } else {
        Write-Host "`nAll tests passed!" -ForegroundColor Green
    }
    
    if ($script:PassedTests.Count -gt 0) {
        Write-Host "`nPassed Tests:" -ForegroundColor Green
        $script:PassedTests | ForEach-Object { Write-Host "  - $_" -ForegroundColor Gray }
    }
    
    # Cleanup
    if (-not $KeepTestFiles) {
        Remove-Item -Path $testDir -Recurse -Force -ErrorAction SilentlyContinue
        Write-Host "`nTest files cleaned up" -ForegroundColor Gray
    } else {
        Write-Host "`nTest files preserved at: $testDir" -ForegroundColor Cyan
    }
    
    return $script:FailedTests.Count -eq 0
}

# Run all tests
try {
    $success = Run-AllTests
    exit $(if ($success) { 0 } else { 1 })
} catch {
    Write-Host "ERROR: Test execution failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}
