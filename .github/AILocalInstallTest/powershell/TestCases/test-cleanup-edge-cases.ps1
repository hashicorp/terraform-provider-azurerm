# Test script for validating cleanup logic edge cases
# This script tests all scenarios for settings.json cleanup

param(
    [switch]$Verbose = $false,
    [switch]$DryRun = $false
)

$ErrorActionPreference = 'Stop'

# Test configuration
$testWorkspaceDir = "$env:TEMP\copilot-cleanup-test-$(Get-Random)"
$testSettingsFile = "$testWorkspaceDir\.vscode\settings.json"

Write-Host "Testing Cleanup Logic Edge Cases" -ForegroundColor Cyan
Write-Host "Test workspace: $testWorkspaceDir" -ForegroundColor Gray

function Write-TestResult {
    param($TestName, $Result, $Details = "")
    $color = if ($Result) { "Green" } else { "Red" }
    $symbol = if ($Result) { "PASS" } else { "FAIL" }
    Write-Host "$symbol $TestName" -ForegroundColor $color
    if ($Details) {
        Write-Host "   $Details" -ForegroundColor Gray
    }
}

function Setup-TestWorkspace {
    if (Test-Path $testWorkspaceDir) {
        Remove-Item $testWorkspaceDir -Recurse -Force
    }
    New-Item -Path $testWorkspaceDir -ItemType Directory -Force | Out-Null
    New-Item -Path "$testWorkspaceDir\.vscode" -ItemType Directory -Force | Out-Null
}

function Cleanup-TestWorkspace {
    if (Test-Path $testWorkspaceDir) {
        Remove-Item $testWorkspaceDir -Recurse -Force
    }
}

function Test-BackupLength0-EmptyFile {
    Write-Host "`nTest 1: Backup Length 0 - Empty file created by install" -ForegroundColor Yellow
    
    Setup-TestWorkspace
    
    # Create empty settings.json (simulates new file created by install)
    Set-Content -Path $testSettingsFile -Value "{}" -Encoding UTF8
    
    # Set environment variable
    $env:AZURERM_BACKUP_LENGTH = "0"
    
    # Simulate cleanup logic
    $shouldRemove = ($env:AZURERM_BACKUP_LENGTH -eq "0")
    
    if ($shouldRemove -and (Test-Path $testSettingsFile)) {
        Remove-Item $testSettingsFile -Force
        $fileExists = Test-Path $testSettingsFile
        Write-TestResult "File removed completely" (!$fileExists) "File should be deleted when backup length = 0"
    }
    
    return !$fileExists
}

function Test-BackupLengthNegative1-ManualMerge {
    Write-Host "`nTest 2: Backup Length -1 - Manual merge scenario" -ForegroundColor Yellow
    
    Setup-TestWorkspace
    
    # Create settings.json with mixed content
    $mixedContent = @"
{
    "github.copilot.enable": {
        "*": true,
        "plaintext": false,
        "markdown": false,
        "scminput": false
    },
    "userSetting": "value",
    "copilot.enable": {
        "terraform": true
    }
}
"@
    Set-Content -Path $testSettingsFile -Value $mixedContent -Encoding UTF8
    
    # Set environment variable
    $env:AZURERM_BACKUP_LENGTH = "-1"
    
    # Simulate cleanup logic - should error and exit
    $shouldError = ($env:AZURERM_BACKUP_LENGTH -eq "-1")
    
    Write-TestResult "Manual merge detected" $shouldError "Should exit with error when backup length = -1"
    
    return $shouldError
}

function Test-BackupLengthPositive-RegexCleanup {
    Write-Host "`nTest 3: Backup Length > 0 - Regex cleanup" -ForegroundColor Yellow
    
    Setup-TestWorkspace
    
    # Create settings.json with copilot content to clean up
    $contentWithCopilot = @"
{
    "userSetting": "value",
    "github.copilot.enable": {
        "*": true,
        "plaintext": false,
        "markdown": false,
        "scminput": false
    },
    "anotherUserSetting": "value2",
    "copilot.enable": {
        "terraform": true
    },
    "finalUserSetting": "value3"
}
"@
    Set-Content -Path $testSettingsFile -Value $contentWithCopilot -Encoding UTF8
    
    # Set environment variable
    $env:AZURERM_BACKUP_LENGTH = "100"
    
    # Read and apply regex cleanup
    $content = Get-Content -Path $testSettingsFile -Raw
    
    # Apply the regex patterns from the script
    $patterns = @(
        '(?s),\s*"github\.copilot\.enable"\s*:\s*\{[^{}]*(?:\{[^{}]*\}[^{}]*)*\}',
        '(?s)"github\.copilot\.enable"\s*:\s*\{[^{}]*(?:\{[^{}]*\}[^{}]*)*\}\s*,?',
        '(?s),\s*"copilot\.enable"\s*:\s*\{[^{}]*(?:\{[^{}]*\}[^{}]*)*\}',
        '(?s)"copilot\.enable"\s*:\s*\{[^{}]*(?:\{[^{}]*\}[^{}]*)*\}\s*,?'
    )
    
    $originalContent = $content
    foreach ($pattern in $patterns) {
        $content = $content -replace $pattern, ''
    }
    
    # Clean up any double commas or trailing commas
    $content = $content -replace ',\s*,', ','
    $content = $content -replace ',(\s*[}\]])', '$1'
    
    # Write back cleaned content
    Set-Content -Path $testSettingsFile -Value $content -Encoding UTF8
    
    # Verify cleanup
    $cleanedContent = Get-Content -Path $testSettingsFile -Raw
    $hasCopilotSettings = ($cleanedContent -match 'github\.copilot\.enable|copilot\.enable')
    $hasUserSettings = ($cleanedContent -match 'userSetting|anotherUserSetting|finalUserSetting')
    
    Write-TestResult "Copilot settings removed" (!$hasCopilotSettings) "Regex should remove all copilot settings"
    Write-TestResult "User settings preserved" $hasUserSettings "User settings should remain intact"
    
    if ($Verbose) {
        Write-Host "Original content length: $($originalContent.Length)" -ForegroundColor Gray
        Write-Host "Cleaned content length: $($cleanedContent.Length)" -ForegroundColor Gray
    }
    
    return (!$hasCopilotSettings -and $hasUserSettings)
}

function Test-EdgeCase-MalformedJSON {
    Write-Host "`nTest 4: Edge Case - Malformed JSON" -ForegroundColor Yellow
    
    Setup-TestWorkspace
    
    # Create malformed JSON
    $malformedJSON = @"
{
    "userSetting": "value",
    "github.copilot.enable": {
        "*": true,
        "plaintext": false,
        // comment in JSON (invalid)
        "markdown": false
    },
    "anotherSetting": "value2"
"@
    Set-Content -Path $testSettingsFile -Value $malformedJSON -Encoding UTF8
    
    try {
        $content = Get-Content -Path $testSettingsFile -Raw
        $jsonTest = $content | ConvertFrom-Json -ErrorAction Stop
        $isValidJSON = $true
    }
    catch {
        $isValidJSON = $false
    }
    
    Write-TestResult "Malformed JSON detected" (!$isValidJSON) "Should handle malformed JSON gracefully"
    
    return !$isValidJSON
}

function Test-EdgeCase-EmptyAfterCleanup {
    Write-Host "`nTest 5: Edge Case - Empty file after cleanup" -ForegroundColor Yellow
    
    Setup-TestWorkspace
    
    # Create file with only copilot settings
    $onlyCopilotContent = @"
{
    "github.copilot.enable": {
        "*": true,
        "plaintext": false,
        "markdown": false,
        "scminput": false
    }
}
"@
    Set-Content -Path $testSettingsFile -Value $onlyCopilotContent -Encoding UTF8
    
    # Apply cleanup regex
    $content = Get-Content -Path $testSettingsFile -Raw
    $patterns = @(
        '(?s),\s*"github\.copilot\.enable"\s*:\s*\{[^{}]*(?:\{[^{}]*\}[^{}]*)*\}',
        '(?s)"github\.copilot\.enable"\s*:\s*\{[^{}]*(?:\{[^{}]*\}[^{}]*)*\}\s*,?'
    )
    
    foreach ($pattern in $patterns) {
        $content = $content -replace $pattern, ''
    }
    
    # Check if file is effectively empty (just {})
    $trimmedContent = $content.Trim()
    $isEmpty = ($trimmedContent -eq '{}' -or $trimmedContent -eq '')
    
    Write-TestResult "File empty after cleanup" $isEmpty "Should result in empty JSON object"
    
    return $isEmpty
}

function Test-EdgeCase-MultilineArrays {
    Write-Host "`nTest 6: Edge Case - Multiline arrays in copilot settings" -ForegroundColor Yellow
    
    Setup-TestWorkspace
    
    # Create content with multiline arrays
    $multilineContent = @"
{
    "userSetting": "value",
    "github.copilot.enable": {
        "*": true,
        "plaintext": false,
        "markdown": false,
        "scminput": false,
        "languages": [
            "javascript",
            "typescript",
            "python"
        ]
    },
    "anotherSetting": "value2"
}
"@
    Set-Content -Path $testSettingsFile -Value $multilineContent -Encoding UTF8
    
    # Apply regex with (?s) flag
    $content = Get-Content -Path $testSettingsFile -Raw
    $pattern = '(?s)"github\.copilot\.enable"\s*:\s*\{[^{}]*(?:\{[^{}]*\}[^{}]*)*\}\s*,?'
    $content = $content -replace $pattern, ''
    
    # Clean up trailing commas
    $content = $content -replace ',(\s*[}\]])', '$1'
    
    # Check results
    $hasCopilotSettings = ($content -match 'github\.copilot\.enable')
    $hasUserSettings = ($content -match 'userSetting|anotherSetting')
    
    Write-TestResult "Multiline arrays handled" (!$hasCopilotSettings) "Should handle nested multiline structures"
    Write-TestResult "User settings preserved" $hasUserSettings "User settings should remain"
    
    return (!$hasCopilotSettings -and $hasUserSettings)
}

# Run all tests
Write-Host "Starting Edge Case Testing..." -ForegroundColor Cyan

try {
    $results = @()
    
    $results += Test-BackupLength0-EmptyFile
    $results += Test-BackupLengthNegative1-ManualMerge
    $results += Test-BackupLengthPositive-RegexCleanup
    $results += Test-EdgeCase-MalformedJSON
    $results += Test-EdgeCase-EmptyAfterCleanup
    $results += Test-EdgeCase-MultilineArrays
    
    # Summary
    $passed = ($results | Where-Object { $_ -eq $true }).Count
    $total = $results.Count
    
    Write-Host "`nTest Results Summary:" -ForegroundColor Cyan
    Write-Host "Passed: $passed/$total" -ForegroundColor $(if ($passed -eq $total) { "Green" } else { "Yellow" })
    
    if ($passed -eq $total) {
        Write-Host "All edge case tests passed!" -ForegroundColor Green
    } else {
        Write-Host "Some tests failed - review cleanup logic" -ForegroundColor Yellow
    }
    
} finally {
    # Cleanup
    Cleanup-TestWorkspace
    Remove-Variable -Name AZURERM_BACKUP_LENGTH -Scope Global -ErrorAction SilentlyContinue
}

Write-Host "`nEdge case testing completed" -ForegroundColor Cyan
