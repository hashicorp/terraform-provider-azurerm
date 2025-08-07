# Edge Case Testing for AzureRM Copilot Install/Cleanup
# Tests unusual scenarios and edge cases

param(
    [switch]$Interactive,
    [switch]$KeepTestFiles
)

$script:TestResults = @()
$testDir = Join-Path $env:TEMP "azurerm-edge-test-$(Get-Date -Format 'yyyyMMdd-HHmmss')"
New-Item -ItemType Directory -Path $testDir -Force | Out-Null

function Add-TestResult {
    param(
        [string]$TestName,
        [string]$Status,
        [string]$Message,
        [string]$Details = ""
    )
    
    $script:TestResults += [PSCustomObject]@{
        Test = $TestName
        Status = $Status
        Message = $Message
        Details = $Details
    }
    
    $color = switch ($Status) {
        "PASS" { "Green" }
        "FAIL" { "Red" }
        "WARN" { "Yellow" }
        "INFO" { "Cyan" }
        default { "White" }
    }
    
    Write-Host "[$Status] $TestName - $Message" -ForegroundColor $color
    if ($Details) {
        Write-Host "  Details: $Details" -ForegroundColor Gray
    }
}

function Test-PermissionIssues {
    Write-Host "`nTesting Permission and Access Issues" -ForegroundColor Cyan
    
    $userDir = Join-Path $testDir "permission-test"
    New-Item -ItemType Directory -Path $userDir -Force | Out-Null
    $settingsPath = Join-Path $userDir "settings.json"
    
    # Test read-only file
    try {
        "{ `"test`": `"value`" }" | Out-File -FilePath $settingsPath -Force
        Set-ItemProperty -Path $settingsPath -Name IsReadOnly -Value $true
        
        $output = & "$PSScriptRoot\install-copilot-setup.ps1" -Clean -UserDirectory $userDir -Force 2>&1
        
        if ($output -match "permission|access|readonly") {
            Add-TestResult "Permission - Read-Only File" "PASS" "Correctly handled read-only file"
        } else {
            Add-TestResult "Permission - Read-Only File" "FAIL" "Did not properly handle read-only file"
        }
        
        # Cleanup
        Set-ItemProperty -Path $settingsPath -Name IsReadOnly -Value $false -ErrorAction SilentlyContinue
        Remove-Item -Path $settingsPath -Force -ErrorAction SilentlyContinue
        
    } catch {
        Add-TestResult "Permission - Read-Only File" "WARN" "Could not test read-only scenario" $_.Exception.Message
    }
}

function Test-LargeFileHandling {
    Write-Host "`nTesting Large File Handling" -ForegroundColor Cyan
    
    $userDir = Join-Path $testDir "large-file-test"
    New-Item -ItemType Directory -Path $userDir -Force | Out-Null
    $settingsPath = Join-Path $userDir "settings.json"
    
    # Create a large settings file with many entries
    $largeSettings = @{
        "// AZURERM_BACKUP_LENGTH" = "500"
        "github.copilot.chat.reviewSelection.enabled" = $true
    }
    
    # Add many dummy settings to make file large
    for ($i = 1; $i -le 100; $i++) {
        $largeSettings["dummy.setting.$i"] = "value_$i"
    }
    
    $largeSettingsJson = $largeSettings | ConvertTo-Json -Depth 10
    Set-Content -Path $settingsPath -Value $largeSettingsJson -Force
    
    $fileSizeBefore = (Get-Item $settingsPath).Length
    
    try {
        $startTime = Get-Date
        & "$PSScriptRoot\install-copilot-setup.ps1" -Clean -UserDirectory $userDir -Force | Out-Null
        $endTime = Get-Date
        $duration = ($endTime - $startTime).TotalSeconds
        
        if (Test-Path $settingsPath) {
            $content = Get-Content $settingsPath -Raw
            $hasAzureRM = $content -match "github\.copilot|AZURERM_BACKUP"
            
            if (-not $hasAzureRM -and $duration -lt 30) {
                Add-TestResult "Large File Handling" "PASS" "Successfully processed large file in $duration seconds"
            } else {
                Add-TestResult "Large File Handling" "FAIL" "Issues with large file processing" "Duration: $duration, AzureRM still present: $hasAzureRM"
            }
        } else {
            Add-TestResult "Large File Handling" "WARN" "Large file was removed entirely"
        }
        
    } catch {
        Add-TestResult "Large File Handling" "FAIL" "Exception during large file processing" $_.Exception.Message
    }
}

function Test-UnicodeAndSpecialCharacters {
    Write-Host "`nTesting Unicode and Special Characters" -ForegroundColor Cyan
    
    $userDir = Join-Path $testDir "unicode-test"
    New-Item -ItemType Directory -Path $userDir -Force | Out-Null
    $settingsPath = Join-Path $userDir "settings.json"
    
    # Create settings with Unicode characters
    $unicodeSettings = @'
{
    "// AZURERM_BACKUP_LENGTH": "100",
    "editor.fontSize": 14,
    "github.copilot.chat.commitMessageGeneration.instructions": [
        {
            "text": "测试中文字符 and émojis 🚀 and spëcial châractërs"
        }
    ],
    "files.associations": {
        "*.instructions.md": "markdown"
    },
    "workbench.theme": "Тёмная тема"
}
'@
    
    Set-Content -Path $settingsPath -Value $unicodeSettings -Encoding UTF8 -Force
    
    try {
        & "$PSScriptRoot\install-copilot-setup.ps1" -Clean -UserDirectory $userDir -Force | Out-Null
        
        if (Test-Path $settingsPath) {
            $content = Get-Content $settingsPath -Raw -Encoding UTF8
            $hasUnicode = $content -match "Тёмная"
            $hasAzureRM = $content -match "github\.copilot"
            
            if ($hasUnicode -and -not $hasAzureRM) {
                Add-TestResult "Unicode Characters" "PASS" "Preserved Unicode while removing AzureRM settings"
            } else {
                Add-TestResult "Unicode Characters" "FAIL" "Issues with Unicode handling" "Unicode preserved: $hasUnicode, AzureRM removed: $(-not $hasAzureRM)"
            }
        } else {
            Add-TestResult "Unicode Characters" "WARN" "File was removed entirely"
        }
        
    } catch {
        Add-TestResult "Unicode Characters" "FAIL" "Exception during Unicode processing" $_.Exception.Message
    }
}

function Test-ConcurrentAccess {
    Write-Host "`nTesting Concurrent Access Scenarios" -ForegroundColor Cyan
    
    $userDir = Join-Path $testDir "concurrent-test"
    New-Item -ItemType Directory -Path $userDir -Force | Out-Null
    $settingsPath = Join-Path $userDir "settings.json"
    
    # Create initial settings
    $settings = @'
{
    "// AZURERM_BACKUP_LENGTH": "150",
    "editor.fontSize": 14,
    "github.copilot.chat.reviewSelection.enabled": true
}
'@
    Set-Content -Path $settingsPath -Value $settings -Force
    
    # Simulate concurrent access by running two cleanup operations
    try {
        $job1 = Start-Job -ScriptBlock {
            param($scriptPath, $userDir)
            & $scriptPath -Clean -UserDirectory $userDir -Force
        } -ArgumentList "$PSScriptRoot\install-copilot-setup.ps1", $userDir
        
        $job2 = Start-Job -ScriptBlock {
            param($scriptPath, $userDir)
            Start-Sleep -Milliseconds 100  # Slight delay
            & $scriptPath -Clean -UserDirectory $userDir -Force
        } -ArgumentList "$PSScriptRoot\install-copilot-setup.ps1", $userDir
        
        $results1 = Wait-Job $job1 | Receive-Job
        $results2 = Wait-Job $job2 | Receive-Job
        
        Remove-Job $job1, $job2
        
        # Check if either operation succeeded without corruption
        if (Test-Path $settingsPath) {
            $content = Get-Content $settingsPath -Raw
            try {
                $json = $content | ConvertFrom-Json
                Add-TestResult "Concurrent Access" "PASS" "No corruption during concurrent operations"
            } catch {
                Add-TestResult "Concurrent Access" "FAIL" "File corruption during concurrent access" $_.Exception.Message
            }
        } else {
            Add-TestResult "Concurrent Access" "PASS" "File was safely removed during concurrent operations"
        }
        
    } catch {
        Add-TestResult "Concurrent Access" "FAIL" "Exception during concurrent access test" $_.Exception.Message
    }
}

function Test-DiskSpaceHandling {
    Write-Host "`nTesting Low Disk Space Scenarios" -ForegroundColor Cyan
    
    # This is a simplified test - real disk space testing would be complex
    $userDir = Join-Path $testDir "diskspace-test"
    New-Item -ItemType Directory -Path $userDir -Force | Out-Null
    $settingsPath = Join-Path $userDir "settings.json"
    
    # Create a settings file
    $settings = '{"// AZURERM_BACKUP_LENGTH": "200", "test": "value"}'
    Set-Content -Path $settingsPath -Value $settings -Force
    
    # Check available disk space
    $drive = (Get-Item $testDir).Root.Name
    $diskSpace = Get-WmiObject -Class Win32_LogicalDisk -Filter "DeviceID='$drive'" | Select-Object -ExpandProperty FreeSpace
    
    if ($diskSpace -gt 100MB) {
        try {
            & "$PSScriptRoot\install-copilot-setup.ps1" -Clean -UserDirectory $userDir -Force | Out-Null
            Add-TestResult "Disk Space" "PASS" "Normal operation with sufficient disk space"
        } catch {
            Add-TestResult "Disk Space" "FAIL" "Failed even with sufficient space" $_.Exception.Message
        }
    } else {
        Add-TestResult "Disk Space" "WARN" "Cannot test - insufficient disk space for test"
    }
}

function Test-NetworkPathScenarios {
    Write-Host "`nTesting Network Path Scenarios" -ForegroundColor Cyan
    
    # Test with UNC-like paths (simulated)
    $networkLikePath = "\\localhost\$testDir"
    
    try {
        # This is a simulation - real network testing would require network setup
        Add-TestResult "Network Paths" "INFO" "Network path testing requires manual setup" "Test path: $networkLikePath"
    } catch {
        Add-TestResult "Network Paths" "WARN" "Could not test network scenarios" $_.Exception.Message
    }
}

function Test-BackupIntegrityCorruption {
    Write-Host "`nTesting Backup Integrity and Corruption" -ForegroundColor Cyan
    
    $userDir = Join-Path $testDir "integrity-test"
    New-Item -ItemType Directory -Path $userDir -Force | Out-Null
    $settingsPath = Join-Path $userDir "settings.json"
    $backupPath = Join-Path $userDir "settings_backup_azurerm.json"
    
    # Create corrupted backup file
    Set-Content -Path $backupPath -Value "{ corrupted json" -Force
    
    # Create current settings
    $settings = @'
{
    "// AZURERM_BACKUP_LENGTH": "50",
    "github.copilot.chat.reviewSelection.enabled": true,
    "editor.fontSize": 14
}
'@
    Set-Content -Path $settingsPath -Value $settings -Force
    
    try {
        $output = & "$PSScriptRoot\install-copilot-setup.ps1" -Clean -UserDirectory $userDir -Force 2>&1
        
        if ($output -match "backup.*corrupt|invalid.*backup") {
            Add-TestResult "Backup Integrity" "PASS" "Correctly detected corrupted backup"
        } else {
            Add-TestResult "Backup Integrity" "WARN" "May not have detected backup corruption"
        }
        
    } catch {
        Add-TestResult "Backup Integrity" "FAIL" "Exception during integrity test" $_.Exception.Message
    }
}

function Test-NestedJsonComplexity {
    Write-Host "`nTesting Complex Nested JSON Structures" -ForegroundColor Cyan
    
    $userDir = Join-Path $testDir "nested-test"
    New-Item -ItemType Directory -Path $userDir -Force | Out-Null
    $settingsPath = Join-Path $userDir "settings.json"
    
    # Create deeply nested complex structure
    $complexSettings = @'
{
    "// AZURERM_BACKUP_LENGTH": "300",
    "editor.fontSize": 14,
    "github.copilot.chat.reviewSelection.instructions": [
        {
            "file": "instructions/implementation-guide.instructions.md",
            "nested": {
                "level1": {
                    "level2": ["item1", "item2"],
                    "level2b": {"key": "value"}
                }
            }
        },
        {"file": "instructions/azure-patterns.instructions.md"}
    ],
    "complex.nested.object": {
        "level1": {
            "level2": {
                "level3": {
                    "github.copilot.advanced": {"length": 8000},
                    "normalSetting": "keep this"
                }
            }
        }
    }
}
'@
    
    Set-Content -Path $settingsPath -Value $complexSettings -Force
    
    try {
        & "$PSScriptRoot\install-copilot-setup.ps1" -Clean -UserDirectory $userDir -Force | Out-Null
        
        if (Test-Path $settingsPath) {
            $content = Get-Content $settingsPath -Raw
            $json = $content | ConvertFrom-Json
            
            $hasOriginal = $json."editor.fontSize" -eq 14
            $hasNested = $json."complex.nested.object".level1.level2.level3.normalSetting -eq "keep this"
            $hasAzureRM = $content -match "github\.copilot"
            
            if ($hasOriginal -and $hasNested -and -not $hasAzureRM) {
                Add-TestResult "Nested JSON" "PASS" "Correctly handled complex nested structures"
            } else {
                Add-TestResult "Nested JSON" "FAIL" "Issues with nested structure handling" "Original: $hasOriginal, Nested: $hasNested, AzureRM removed: $(-not $hasAzureRM)"
            }
        } else {
            Add-TestResult "Nested JSON" "WARN" "Complex file was removed entirely"
        }
        
    } catch {
        Add-TestResult "Nested JSON" "FAIL" "Exception during nested JSON processing" $_.Exception.Message
    }
}

function Run-EdgeCaseTests {
    Write-Host "Running Edge Case Tests for AzureRM Copilot" -ForegroundColor Green
    Write-Host "=============================================" -ForegroundColor Green
    
    Test-PermissionIssues
    Test-LargeFileHandling  
    Test-UnicodeAndSpecialCharacters
    Test-ConcurrentAccess
    Test-DiskSpaceHandling
    Test-NetworkPathScenarios
    Test-BackupIntegrityCorruption
    Test-NestedJsonComplexity
    
    Write-Host "`nEdge Case Test Summary" -ForegroundColor Green
    Write-Host "=========================" -ForegroundColor Green
    
    $passCount = ($script:TestResults | Where-Object { $_.Status -eq "PASS" }).Count
    $failCount = ($script:TestResults | Where-Object { $_.Status -eq "FAIL" }).Count
    $warnCount = ($script:TestResults | Where-Object { $_.Status -eq "WARN" }).Count
    $infoCount = ($script:TestResults | Where-Object { $_.Status -eq "INFO" }).Count
    
    Write-Host "Total Tests: $($script:TestResults.Count)" -ForegroundColor Cyan
    Write-Host "Passed: $passCount" -ForegroundColor Green
    Write-Host "Failed: $failCount" -ForegroundColor Red
    Write-Host "Warnings: $warnCount" -ForegroundColor Yellow
    Write-Host "Info: $infoCount" -ForegroundColor Gray
    
    if ($failCount -gt 0) {
        Write-Host "`nFailed Tests:" -ForegroundColor Red
        $script:TestResults | Where-Object { $_.Status -eq "FAIL" } | ForEach-Object {
            Write-Host "  FAIL - $($_.Test): $($_.Message)" -ForegroundColor Yellow
        }
    }
    
    if ($warnCount -gt 0) {
        Write-Host "`nWarnings:" -ForegroundColor Yellow
        $script:TestResults | Where-Object { $_.Status -eq "WARN" } | ForEach-Object {
            Write-Host "  WARNING - $($_.Test): $($_.Message)" -ForegroundColor Gray
        }
    }
    
    # Cleanup
    if (-not $KeepTestFiles) {
        Remove-Item -Path $testDir -Recurse -Force -ErrorAction SilentlyContinue
        Write-Host "`nTest files cleaned up" -ForegroundColor Gray
    } else {
        Write-Host "`nTest files preserved at: $testDir" -ForegroundColor Cyan
    }
    
    return $failCount -eq 0
}

# Run the edge case tests
try {
    $success = Run-EdgeCaseTests
    
    if ($Interactive) {
        Write-Host "`nPress any key to continue..." -ForegroundColor Gray
        $null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
    }
    
    exit $(if ($success) { 0 } else { 1 })
} catch {
    Write-Host "ERROR: Edge case test execution failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}
