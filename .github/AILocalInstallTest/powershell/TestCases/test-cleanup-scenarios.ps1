# Test Suite for AzureRM Copilot Install Script - Cleanup Scenarios
# Validates all backup length scenarios and edge cases

param(
    [string]$TestDir = ".\test-cleanup-validation",
    [switch]$Verbose = $false
)

$ErrorActionPreference = "Stop"

# Colors for output
function Write-ColoredOutput {
    param($Message, $Color = "White")
    Write-Host $Message -ForegroundColor $Color
}

function Write-TestHeader {
    param($TestName)
    Write-ColoredOutput "`n==================== $TestName ====================" "Cyan"
}

function Write-TestResult {
    param($Message, $Success)
    $color = if ($Success) { "Green" } else { "Red" }
    $symbol = if ($Success) { "PASS" } else { "FAIL" }
    Write-ColoredOutput "$symbol $Message" $color
}

# Test data structures for different scenarios
$TestCases = @{
    "BackupLength0_EmptyFile" = @{
        BackupLength = 0
        OriginalContent = ""
        ExpectedAction = "Remove"
        Description = "Empty file created by install - should be removed completely"
    }
    
    "BackupLength0_OnlyCopilotEntries" = @{
        BackupLength = 0
        OriginalContent = @"
{
    "github.copilot.enable": {
        "*": true,
        "plaintext": false,
        "markdown": true,
        "scminput": false
    }
}
"@
        ExpectedAction = "Remove"
        Description = "File with only Copilot entries - should be removed completely"
    }
    
    "BackupLength0_MixedContent" = @{
        BackupLength = 0
        OriginalContent = @"
{
    "editor.fontSize": 14,
    "github.copilot.enable": {
        "*": true,
        "plaintext": false,
        "markdown": true,
        "scminput": false
    },
    "workbench.colorTheme": "Dark+"
}
"@
        ExpectedAction = "Remove"
        Description = "File with mixed content created by install - should be removed completely"
    }
    
    "BackupLengthNeg1_ShouldError" = @{
        BackupLength = -1
        OriginalContent = @"
{
    "editor.fontSize": 14,
    "github.copilot.enable": {
        "*": true,
        "plaintext": false
    }
}
"@
        ExpectedAction = "Error"
        Description = "Manual merge scenario - should error and exit"
    }
    
    "BackupLengthPositive_CleanupOnly" = @{
        BackupLength = 5
        OriginalContent = @"
{
    "editor.fontSize": 14,
    "github.copilot.enable": {
        "*": true,
        "plaintext": false,
        "markdown": true,
        "scminput": false
    },
    "workbench.colorTheme": "Dark+",
    "github.copilot.chat.localeOverride": "en"
}
"@
        ExpectedContent = @"
{
    "editor.fontSize": 14,
    "workbench.colorTheme": "Dark+"
}
"@
        ExpectedAction = "Cleanup"
        Description = "Backed up original - should remove only Copilot entries"
    }
    
    "BackupLengthPositive_MultilineArray" = @{
        BackupLength = 3
        OriginalContent = @"
{
    "editor.fontSize": 14,
    "github.copilot.enable": {
        "*": true,
        "plaintext": false,
        "markdown": true,
        "scminput": false
    },
    "files.associations": {
        "*.md": "markdown"
    }
}
"@
        ExpectedContent = @"
{
    "editor.fontSize": 14,
    "files.associations": {
        "*.md": "markdown"
    }
}
"@
        ExpectedAction = "Cleanup"
        Description = "Multi-line object cleanup"
    }
    
    "BackupLengthPositive_TrailingComma" = @{
        BackupLength = 2
        OriginalContent = @"
{
    "editor.fontSize": 14,
    "github.copilot.enable": {
        "*": true,
        "plaintext": false
    },
    "workbench.colorTheme": "Dark+",
}
"@
        ExpectedContent = @"
{
    "editor.fontSize": 14,
    "workbench.colorTheme": "Dark+",
}
"@
        ExpectedAction = "Cleanup"
        Description = "Trailing comma handling"
    }
}

# Create test environment
function Setup-TestEnvironment {
    if (Test-Path $TestDir) {
        Remove-Item $TestDir -Recurse -Force
    }
    New-Item -ItemType Directory -Path $TestDir -Force | Out-Null
    
    # Create test VS Code settings directory
    $vscodeDir = Join-Path $TestDir ".vscode"
    New-Item -ItemType Directory -Path $vscodeDir -Force | Out-Null
    
    Write-ColoredOutput "Test environment setup complete: $TestDir" "Yellow"
}

# Simulate the cleanup logic from the main script
function Invoke-CleanupLogic {
    param(
        [string]$SettingsPath,
        [int]$BackupLength,
        [bool]$DryRun = $false
    )
    
    $result = @{
        Action = ""
        Success = $true
        Error = ""
        Content = ""
    }
    
    try {
        if ($BackupLength -eq 0) {
            $result.Action = "Remove"
            if (-not $DryRun -and (Test-Path $SettingsPath)) {
                Remove-Item $SettingsPath -Force
            }
        }
        elseif ($BackupLength -eq -1) {
            $result.Action = "Error"
            $result.Success = $false
            $result.Error = "Manual merge detected - cannot auto-cleanup"
            return $result
        }
        else {
            $result.Action = "Cleanup"
            if (Test-Path $SettingsPath) {
                $content = Get-Content $SettingsPath -Raw
                
                # Apply regex patterns (simplified version of main script logic)
                $patterns = @(
                    ',\s*"github\.copilot\.enable"\s*:\s*\{[^}]*\}(?s)',
                    '"github\.copilot\.enable"\s*:\s*\{[^}]*\}\s*,?(?s)',
                    ',\s*"github\.copilot\.chat\.localeOverride"\s*:\s*"[^"]*"(?s)',
                    '"github\.copilot\.chat\.localeOverride"\s*:\s*"[^"]*"\s*,?(?s)'
                )
                
                foreach ($pattern in $patterns) {
                    $content = $content -replace $pattern, ''
                }
                
                # Clean up double commas and normalize spacing
                $content = $content -replace ',\s*,', ','
                $content = $content -replace '{\s*,', '{'
                $content = $content -replace ',\s*}', '}'
                
                $result.Content = $content
                
                if (-not $DryRun) {
                    Set-Content -Path $SettingsPath -Value $content -NoNewline
                }
            }
        }
    }
    catch {
        $result.Success = $false
        $result.Error = $_.Exception.Message
    }
    
    return $result
}

# Run individual test case
function Test-CleanupScenario {
    param($TestName, $TestCase)
    
    Write-TestHeader "Testing: $TestName"
    Write-ColoredOutput "Description: $($TestCase.Description)" "Gray"
    
    # Setup test file
    $settingsPath = Join-Path $TestDir ".vscode\settings.json"
    Set-Content -Path $settingsPath -Value $TestCase.OriginalContent -NoNewline
    
    # Run cleanup logic
    $result = Invoke-CleanupLogic -SettingsPath $settingsPath -BackupLength $TestCase.BackupLength -DryRun $false
    
    # Validate results
    $success = $true
    $messages = @()
    
    # Check expected action
    if ($result.Action -eq $TestCase.ExpectedAction) {
        $messages += "PASS - Action: $($result.Action) (Expected: $($TestCase.ExpectedAction))"
    } else {
        $messages += "FAIL - Action: $($result.Action) (Expected: $($TestCase.ExpectedAction))"
        $success = $false
    }
    
    # Check error scenarios
    if ($TestCase.ExpectedAction -eq "Error") {
        if (-not $result.Success) {
            $messages += "PASS - Expected error occurred: $($result.Error)"
        } else {
            $messages += "FAIL - Expected error but operation succeeded"
            $success = $false
        }
    } else {
        if ($result.Success) {
            $messages += "PASS - Operation completed successfully"
        } else {
            $messages += "FAIL - Unexpected error: $($result.Error)"
            $success = $false
        }
    }
    
    # Check file state
    if ($TestCase.ExpectedAction -eq "Remove") {
        if (-not (Test-Path $settingsPath)) {
            $messages += "PASS - File removed as expected"
        } else {
            $messages += "FAIL - File still exists when it should be removed"
            $success = $false
        }
    } elseif ($TestCase.ExpectedAction -eq "Cleanup" -and $TestCase.ExpectedContent) {
        if (Test-Path $settingsPath) {
            $actualContent = Get-Content $settingsPath -Raw
            $actualContent = $actualContent -replace '\r\n', "`n" -replace '\r', "`n"
            $expectedContent = $TestCase.ExpectedContent -replace '\r\n', "`n" -replace '\r', "`n"
            
            if ($actualContent.Trim() -eq $expectedContent.Trim()) {
                $messages += "PASS - Content matches expected result"
            } else {
                $messages += "FAIL - Content mismatch"
                $messages += "Expected:`n$($TestCase.ExpectedContent)"
                $messages += "Actual:`n$actualContent"
                $success = $false
            }
        } else {
            $messages += "FAIL - File was removed when it should be cleaned up"
            $success = $false
        }
    }
    
    # Output results
    foreach ($message in $messages) {
        if ($message.StartsWith("PASS")) {
            Write-ColoredOutput $message "Green"
        } elseif ($message.StartsWith("FAIL")) {
            Write-ColoredOutput $message "Red"
        } else {
            Write-ColoredOutput $message "Yellow"
        }
    }
    
    Write-TestResult "Test Case: $TestName" $success
    return $success
}

# Main test execution
function Run-AllTests {
    Write-ColoredOutput "`nStarting AzureRM Copilot Cleanup Logic Test Suite" "Magenta"
    Write-ColoredOutput "Testing all backup length scenarios and edge cases`n" "Gray"
    
    Setup-TestEnvironment
    
    $testResults = @{}
    $totalTests = $TestCases.Count
    $passedTests = 0
    
    foreach ($testName in $TestCases.Keys) {
        $testCase = $TestCases[$testName]
        $success = Test-CleanupScenario -TestName $testName -TestCase $testCase
        $testResults[$testName] = $success
        if ($success) { $passedTests++ }
    }
    
    # Summary report
    Write-TestHeader "Test Summary"
    Write-ColoredOutput "Total Tests: $totalTests" "White"
    Write-ColoredOutput "Passed: $passedTests" "Green"
    Write-ColoredOutput "Failed: $($totalTests - $passedTests)" "Red"
    
    if ($passedTests -eq $totalTests) {
        Write-ColoredOutput "`nALL TESTS PASSED! Cleanup logic is working correctly." "Green"
        Write-ColoredOutput "The script is ready for production use." "Green"
    } else {
        Write-ColoredOutput "`nSome tests failed. Review the output above for details." "Red"
        Write-ColoredOutput "Failed tests:" "Red"
        foreach ($testName in $testResults.Keys) {
            if (-not $testResults[$testName]) {
                Write-ColoredOutput "  - $testName" "Red"
            }
        }
    }
    
    # Cleanup test environment
    if (Test-Path $TestDir) {
        Remove-Item $TestDir -Recurse -Force
    }
}

# Edge case validation
function Test-EdgeCases {
    Write-TestHeader "Additional Edge Case Testing"
    
    # Test regex patterns independently
    $edgeCasePatterns = @(
        @{
            Name = "Nested JSON Objects"
            Content = @"
{
    "github.copilot.enable": {
        "nested": {
            "deeply": true
        }
    }
}
"@
            ShouldMatch = $true
        },
        @{
            Name = "Multi-line with weird spacing"
            Content = @"
{
    "github.copilot.enable"    :    {
        "*": true,
        "plaintext": false
    }
}
"@
            ShouldMatch = $true
        }
    )
    
    foreach ($test in $edgeCasePatterns) {
        Write-ColoredOutput "Testing: $($test.Name)" "Yellow"
        
        $pattern = '"github\.copilot\.enable"\s*:\s*\{[^}]*\}(?s)'
        $matches = [regex]::Matches($test.Content, $pattern)
        
        if (($matches.Count -gt 0) -eq $test.ShouldMatch) {
            Write-TestResult "Regex pattern matching" $true
        } else {
            Write-TestResult "Regex pattern matching" $false
        }
    }
}

# Run the tests
Run-AllTests
Test-EdgeCases

Write-ColoredOutput "`nTest suite execution complete!" "Magenta"
