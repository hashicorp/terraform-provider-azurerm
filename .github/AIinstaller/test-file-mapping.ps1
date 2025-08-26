# Test script to show file source→destination mappings
# Bypasses safety checks to show actual download logic

# Import modules
$modulePath = Join-Path $PSScriptRoot "modules\powershell"
Import-Module (Join-Path $modulePath "ConfigParser.psm1") -Force
Import-Module (Join-Path $modulePath "FileOperations.psm1") -Force

Write-Host "=== TESTING FILE MAPPING LOGIC ===" -ForegroundColor Cyan
Write-Host ""

# Test with your actual workspace
$workspaceRoot = "C:\github.com\hashicorp\terraform-provider-azurerm"

# Get manifest config
Write-Host "Loading manifest configuration..." -ForegroundColor Yellow
$manifestConfig = Get-ManifestConfig

Write-Host "Found $($manifestConfig.Sections.Keys.Count) sections in manifest" -ForegroundColor Green
Write-Host "Base URL: $($manifestConfig.BaseUrl)" -ForegroundColor Green
Write-Host ""

# Show first few files from each section
foreach ($sectionName in $manifestConfig.Sections.Keys) {
    $files = $manifestConfig.Sections[$sectionName]
    Write-Host "Section: $sectionName ($($files.Count) files)" -ForegroundColor Cyan
    
    # Show first 3 files from this section
    $sampleFiles = $files | Select-Object -First 3
    foreach ($file in $sampleFiles) {
        $downloadUrl = "$($manifestConfig.BaseUrl)/$file"
        $targetPath = Join-Path $workspaceRoot $file
        
        Write-Host "  SOURCE: $downloadUrl" -ForegroundColor Yellow
        Write-Host "  TARGET: $targetPath" -ForegroundColor Green
        Write-Host ""
    }
    
    if ($files.Count -gt 3) {
        Write-Host "  ... and $($files.Count - 3) more files" -ForegroundColor Gray
        Write-Host ""
    }
}

Write-Host "=== MAPPING TEST COMPLETE ===" -ForegroundColor Cyan
