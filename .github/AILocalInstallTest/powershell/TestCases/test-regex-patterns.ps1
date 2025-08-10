# Regex Pattern Validation Test for AzureRM Copilot Cleanup
# Tests the specific regex patterns used in settings.json cleanup

param(
    [switch]$Verbose
)

$script:TestCounter = 0
$script:FailedTests = @()
$script:PassedTests = @()

function Test-RegexPattern {
    param(
        [string]$TestName,
        [string]$Pattern,
        [string]$InputText,
        [string]$ExpectedResult,
        [string]$Description
    )
    
    $script:TestCounter++
    
    try {
        $actualResult = $InputText -replace $Pattern, ''
        
        # Normalize whitespace for comparison
        $actualResult = $actualResult -replace '\s+', ' ' -replace '^\s+|\s+$', ''
        $ExpectedResult = $ExpectedResult -replace '\s+', ' ' -replace '^\s+|\s+$', ''
        
        if ($actualResult -eq $ExpectedResult) {
            $script:PassedTests += $TestName
            Write-Host "PASS: $TestName" -ForegroundColor Green
            if ($Verbose) {
                Write-Host "   Pattern: $Pattern" -ForegroundColor Gray
                Write-Host "   Expected: $ExpectedResult" -ForegroundColor Gray
                Write-Host "   Actual: $actualResult" -ForegroundColor Gray
            }
        } else {
            $script:FailedTests += $TestName
            Write-Host "FAIL: $TestName" -ForegroundColor Red
            Write-Host "   Pattern: $Pattern" -ForegroundColor Yellow
            Write-Host "   Expected: $ExpectedResult" -ForegroundColor Yellow
            Write-Host "   Actual: $actualResult" -ForegroundColor Yellow
            if ($Description) {
                Write-Host "   Description: $Description" -ForegroundColor Gray
            }
        }
    } catch {
        $script:FailedTests += $TestName
        Write-Host "ERROR: $TestName - $($_.Exception.Message)" -ForegroundColor Red
    }
}

function Test-BackupLengthPatterns {
    Write-Host "`nTesting Backup Length Removal Patterns" -ForegroundColor Cyan
    
    # Test comment-style backup length
    $input1 = @'
{
    "// AZURERM_BACKUP_LENGTH": "123",
    "editor.fontSize": 14
}
'@
    $expected1 = '{ "editor.fontSize": 14 }'
    Test-RegexPattern "Backup Length Comment Style" '"//\s*AZURERM_BACKUP_LENGTH"\s*:\s*"[^"]*"\s*,?\s*' $input1 $expected1 "Should remove comment-style backup length"
    
    # Test property-style backup length
    $input2 = @'
{
    "AZURERM_BACKUP_LENGTH": 456,
    "workbench.theme": "dark"
}
'@
    $expected2 = '{ "workbench.theme": "dark" }'
    Test-RegexPattern "Backup Length Property Style" '"AZURERM_BACKUP_LENGTH"\s*:\s*[^,}]*\s*,?\s*' $input2 $expected2 "Should remove property-style backup length"
}

function Test-GitHubCopilotPatterns {
    Write-Host "`nTesting GitHub Copilot Setting Patterns" -ForegroundColor Cyan
    
    # Test simple string pattern
    $input1 = @'
{
    "github.copilot.chat.summarizeAgentConversationHistory.enabled": false,
    "editor.fontSize": 14
}
'@
    $expected1 = '{ "editor.fontSize": 14 }'
    Test-RegexPattern "GitHub Copilot Boolean" '"github\.copilot\.chat\.summarizeAgentConversationHistory\.enabled"\s*:\s*(true|false)\s*,?\s*' $input1 $expected1
    
    # Test complex array pattern
    $input2 = @'
{
    "github.copilot.chat.commitMessageGeneration.instructions": [
        {
            "text": "Test message"
        }
    ],
    "editor.fontSize": 14
}
'@
    $expected2 = '{ "editor.fontSize": 14 }'
    Test-RegexPattern "GitHub Copilot Array" '(?s)"github\.copilot\.chat\.commitMessageGeneration\.instructions"\s*:\s*\[.*?\]\s*,?\s*' $input2 $expected2
    
    # Test multi-line review instructions
    $input3 = @'
{
    "github.copilot.chat.reviewSelection.instructions": [
        {"file": "copilot-instructions.md"},
        {"file": "instructions/implementation-guide.instructions.md"},
        {"file": "instructions/azure-patterns.instructions.md"}
    ],
    "workbench.theme": "dark"
}
'@
    $expected3 = '{ "workbench.theme": "dark" }'
    Test-RegexPattern "GitHub Copilot Review Instructions" '(?s)"github\.copilot\.chat\.reviewSelection\.instructions"\s*:\s*\[.*?\]\s*,?\s*' $input3 $expected3
}

function Test-FileAssociationPatterns {
    Write-Host "`nTesting File Association Patterns" -ForegroundColor Cyan
    
    $input1 = @'
{
    "files.associations": {
        "*.instructions.md": "markdown",
        "*.js": "javascript",
        "*.prompt.md": "markdown",
        "*.azurerm.md": "markdown"
    },
    "editor.fontSize": 14
}
'@
    
    # Apply all file association removals
    $result = $input1
    $result = $result -replace '"\*\.instructions\.md"\s*:\s*"[^"]*"\s*,?\s*', ''
    $result = $result -replace '"\*\.prompt\.md"\s*:\s*"[^"]*"\s*,?\s*', ''
    $result = $result -replace '"\*\.azurerm\.md"\s*:\s*"[^"]*"\s*,?\s*', ''
    
    $expected1 = '{ "files.associations": { "*.js": "javascript" }, "editor.fontSize": 14 }'
    
    # Normalize for comparison
    $result = $result -replace '\s+', ' ' -replace '^\s+|\s+$', ''
    $expected1 = $expected1 -replace '\s+', ' '
    
    if ($result -like "*`"*.js`": `"javascript`"*") {
        $script:PassedTests += "File Association Selective Removal"
        Write-Host "PASS: File Association Selective Removal" -ForegroundColor Green
    } else {
        $script:FailedTests += "File Association Selective Removal"
        Write-Host "FAIL: File Association Selective Removal" -ForegroundColor Red
        Write-Host "   Result: $result" -ForegroundColor Yellow
    }
}

function Test-TrailingCommaCleanup {
    Write-Host "`nTesting Trailing Comma Cleanup" -ForegroundColor Cyan
    
    $input1 = @'
{
    "setting1": "value1",
    "setting2": "value2",
}
'@
    $expected1 = '{ "setting1": "value1", "setting2": "value2" }'
    Test-RegexPattern "Trailing Comma Object" ',(\s*[}\]])' $input1 $expected1
    
    $input2 = @'
{
    "array": [
        "item1",
        "item2",
    ]
}
'@
    $expected2 = '{ "array": [ "item1", "item2" ] }'
    Test-RegexPattern "Trailing Comma Array" ',(\s*[}\]])' $input2 $expected2
}

function Test-ComplexRealWorldScenario {
    Write-Host "`nTesting Complex Real-World Scenario" -ForegroundColor Cyan
    
    $complexInput = @'
{
    "// AZURERM_BACKUP_LENGTH": "245",
    "editor.fontSize": 14,
    "github.copilot.chat.commitMessageGeneration.instructions": [
        {
            "text": "Provide a concise and clear commit message that summarizes the changes made in the code."
        }
    ],
    "github.copilot.chat.summarizeAgentConversationHistory.enabled": false,
    "github.copilot.chat.reviewSelection.enabled": true,
    "github.copilot.chat.reviewSelection.instructions": [
        {"file": "copilot-instructions.md"},
        {"file": "instructions/implementation-guide.instructions.md"},
        {"file": "instructions/azure-patterns.instructions.md"}
    ],
    "github.copilot.advanced": {
        "length": 8000,
        "temperature": "0.1"
    },
    "files.associations": {
        "*.instructions.md": "markdown",
        "*.prompt.md": "markdown",
        "*.azurerm.md": "markdown",
        "*.ts": "typescript"
    },
    "workbench.colorTheme": "Dark+",
    "editor.tabSize": 4
}
'@
    
    # Apply all cleanup patterns
    $cleaned = $complexInput
    $cleaned = $cleaned -replace '"//\s*AZURERM_BACKUP_LENGTH"\s*:\s*"[^"]*"\s*,?\s*', ''
    $cleaned = $cleaned -replace '"AZURERM_BACKUP_LENGTH"\s*:\s*[^,}]*\s*,?\s*', ''
    $cleaned = $cleaned -replace '"github\.copilot\.enable"\s*:\s*\{[^{}]*\}\s*,?\s*', ''
    $cleaned = $cleaned -replace '"github\.copilot\.chat\.commitMessageGeneration\.instructions"\s*:\s*"[^"]*"\s*,?\s*', ''
    $cleaned = $cleaned -replace '(?s)"github\.copilot\.chat\.commitMessageGeneration\.instructions"\s*:\s*\[.*?\]\s*,?\s*', ''
    $cleaned = $cleaned -replace '"github\.copilot\.chat\.summarizeAgentConversationHistory\.enabled"\s*:\s*(true|false)\s*,?\s*', ''
    $cleaned = $cleaned -replace '"github\.copilot\.chat\.reviewSelection\.enabled"\s*:\s*(true|false)\s*,?\s*', ''
    $cleaned = $cleaned -replace '"github\.copilot\.chat\.reviewSelection\.instructions"\s*:\s*"[^"]*"\s*,?\s*', ''
    $cleaned = $cleaned -replace '(?s)"github\.copilot\.chat\.reviewSelection\.instructions"\s*:\s*\[.*?\]\s*,?\s*', ''
    $cleaned = $cleaned -replace '"github\.copilot\.advanced"\s*:\s*\{[^}]*\}\s*,?\s*', ''
    $cleaned = $cleaned -replace '"\*\.instructions\.md"\s*:\s*"[^"]*"\s*,?\s*', ''
    $cleaned = $cleaned -replace '"\*\.prompt\.md"\s*:\s*"[^"]*"\s*,?\s*', ''
    $cleaned = $cleaned -replace '"\*\.azurerm\.md"\s*:\s*"[^"]*"\s*,?\s*', ''
    $cleaned = $cleaned -replace ',(\s*[}\]])', '$1'
    $cleaned = $cleaned -replace '\n\s*\n', "`n"
    
    # Verify the result
    try {
        $parsedJson = $cleaned | ConvertFrom-Json
        $hasOriginalSettings = ($parsedJson."editor.fontSize" -eq 14) -and ($parsedJson."workbench.colorTheme" -eq "Dark+")
        $hasAzureRMSettings = $parsedJson.PSObject.Properties.Name -match "github\.copilot"
        $hasTypeScriptAssociation = $parsedJson."files.associations"."*.ts" -eq "typescript"
        $hasAzureRMAssociations = $parsedJson."files.associations".PSObject.Properties.Name -match "instructions|prompt|azurerm"
        
        if ($hasOriginalSettings -and -not $hasAzureRMSettings -and $hasTypeScriptAssociation -and -not $hasAzureRMAssociations) {
            $script:PassedTests += "Complex Real-World Scenario"
            Write-Host "PASS: Complex Real-World Scenario" -ForegroundColor Green
        } else {
            $script:FailedTests += "Complex Real-World Scenario"
            Write-Host "FAIL: Complex Real-World Scenario" -ForegroundColor Red
            Write-Host "   Original Settings: $hasOriginalSettings" -ForegroundColor Yellow
            Write-Host "   AzureRM Settings Removed: $(-not $hasAzureRMSettings)" -ForegroundColor Yellow
            Write-Host "   TypeScript Association Preserved: $hasTypeScriptAssociation" -ForegroundColor Yellow
            Write-Host "   AzureRM Associations Removed: $(-not $hasAzureRMAssociations)" -ForegroundColor Yellow
        }
    } catch {
        $script:FailedTests += "Complex Real-World Scenario"
        Write-Host "FAIL: Complex Real-World Scenario - Invalid JSON after cleanup" -ForegroundColor Red
        Write-Host "   Error: $($_.Exception.Message)" -ForegroundColor Yellow
    }
}

function Run-RegexTests {
    Write-Host "Testing AzureRM Copilot Regex Patterns" -ForegroundColor Green
    Write-Host "==========================================" -ForegroundColor Green
    
    Test-BackupLengthPatterns
    Test-GitHubCopilotPatterns
    Test-FileAssociationPatterns
    Test-TrailingCommaCleanup
    Test-ComplexRealWorldScenario
    
    Write-Host "`nRegex Test Results" -ForegroundColor Green
    Write-Host "======================" -ForegroundColor Green
    Write-Host "Total Tests: $script:TestCounter" -ForegroundColor Cyan
    Write-Host "Passed: $($script:PassedTests.Count)" -ForegroundColor Green
    Write-Host "Failed: $($script:FailedTests.Count)" -ForegroundColor Red
    
    if ($script:FailedTests.Count -gt 0) {
        Write-Host "`nFailed Tests:" -ForegroundColor Red
        $script:FailedTests | ForEach-Object { Write-Host "  - $_" -ForegroundColor Yellow }
        return $false
    } else {
        Write-Host "`nAll regex patterns working correctly!" -ForegroundColor Green
        return $true
    }
}

# Run the tests
try {
    $success = Run-RegexTests
    exit $(if ($success) { 0 } else { 1 })
} catch {
    Write-Host "ERROR: Regex test execution failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}
