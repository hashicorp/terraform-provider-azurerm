# Test Case: Fresh Installation Scenarios
# Tests installation on clean VS Code environments

param(
    [switch]$CleanupAfter = $true
)

$ErrorActionPreference = "Stop"
$TestBase = "C:\Users\$env:USERNAME\AppData\Local\Temp"
$ScriptPath = "..\..\AILocalInstall\install-copilot-setup.ps1"

function Test-FreshInstallation {
    Write-Host "=== FRESH INSTALLATION TEST ===" -ForegroundColor Cyan
    
    $testDir = "$TestBase\AzureRM-Fresh-$(Get-Random)"
    $codeUserDir = "$testDir\Code\User"
    
    try {
        # Create clean environment
        New-Item -Path $codeUserDir -ItemType Directory -Force | Out-Null
        Write-Host "Created clean test environment: $testDir" -ForegroundColor Gray
        
        # Run installation
        $originalAppData = $env:APPDATA
        $env:APPDATA = $testDir
        
        Write-Host "Running fresh installation..." -ForegroundColor Yellow
        & PowerShell -ExecutionPolicy Bypass -File $ScriptPath -Force
        
        # Verify installation
        $settingsExists = Test-Path "$codeUserDir\settings.json"
        $instructionsExist = Test-Path "$codeUserDir\instructions\terraform-azurerm"
        $copilotInstructionsExist = Test-Path "$codeUserDir\copilot-instructions.md"
        
        Write-Host "Verification Results:" -ForegroundColor White
        Write-Host "  Settings.json created: $settingsExists" -ForegroundColor $(if ($settingsExists) { "Green" } else { "Red" })
        Write-Host "  Instructions directory: $instructionsExist" -ForegroundColor $(if ($instructionsExist) { "Green" } else { "Red" })
        Write-Host "  Copilot instructions: $copilotInstructionsExist" -ForegroundColor $(if ($copilotInstructionsExist) { "Green" } else { "Red" })
        
        if ($settingsExists) {
            $settings = Get-Content "$codeUserDir\settings.json" | ConvertFrom-Json
            $hasAzureRM = $settings."github.copilot.chat.reviewSelection.instructions" -ne $null
            Write-Host "  AzureRM configuration: $hasAzureRM" -ForegroundColor $(if ($hasAzureRM) { "Green" } else { "Red" })
        }
        
        $env:APPDATA = $originalAppData
        
        $success = $settingsExists -and $instructionsExist -and $copilotInstructionsExist
        
        if ($success) {
            Write-Host "FRESH INSTALLATION TEST PASSED" -ForegroundColor Green
        } else {
            Write-Host "FRESH INSTALLATION TEST FAILED" -ForegroundColor Red
        }
        
        return $success
        
    } catch {
        Write-Host "FRESH INSTALLATION TEST FAILED: $($_.Exception.Message)" -ForegroundColor Red
        return $false
    } finally {
        if ($CleanupAfter -and (Test-Path $testDir)) {
            Remove-Item $testDir -Recurse -Force -ErrorAction SilentlyContinue
            Write-Host "Cleaned up test environment" -ForegroundColor Gray
        }
    }
}

# Run the test
$result = Test-FreshInstallation

if ($result) {
    exit 0
} else {
    exit 1
}
