# Import configuration management module
Import-Module (Join-Path $PSScriptRoot "config-management.psm1") -Force

function Test-HardcodedInstallationIntegrity {
    param(
        [Parameter(Mandatory = $true)]
        [string]$RepositoryPath
    )
    
    try {
        Write-Host "Verifying AI installation (hardcoded check)..." -ForegroundColor Cyan
        
        # Initialize counters
        $foundInstructionFiles = 0
        $foundPromptFiles = 0
        $hasCopilotFile = $false
        
        # Check VS Code settings
        $settingsPath = Join-Path $env:APPDATA "Code\User\settings.json"
        $hasVSCodeSettings = $false
        
        if (Test-Path $settingsPath) {
            try {
                $content = Get-Content $settingsPath -Raw
                if ($content -like "*AZURERM_BACKUP_LENGTH*" -or $content -like "*AZURERM_INSTALLATION_DATE*") {
                    $hasVSCodeSettings = $true
                    Write-Host "Found AI settings in VS Code" -ForegroundColor Green
                    
                    # Check backup length values
                    try {
                        $settingsObj = $content | ConvertFrom-Json
                        $backupLength = $settingsObj.'// AZURERM_BACKUP_LENGTH'
                        if ($backupLength -eq 0) {
                            Write-Host "  - Backup status: No original settings.json existed (fake backup created)" -ForegroundColor Yellow
                        } elseif ($backupLength -eq -1) {
                            Write-Host "  - Backup status: Manual merge scenario detected" -ForegroundColor Yellow
                        } elseif ($backupLength -gt 0) {
                            Write-Host "  - Backup status: Original settings.json backed up ($backupLength chars)" -ForegroundColor Yellow
                        }
                    } catch {
                        # Ignore JSON parsing errors for backup length - not critical
                    }
                } else {
                    Write-Host "No AI settings found in VS Code" -ForegroundColor Red
                }
            } catch {
                Write-Host "Could not parse VS Code settings: $_" -ForegroundColor Red
            }
        } else {
            Write-Host "VS Code settings.json not found" -ForegroundColor Red
        }
        
        # Check instruction files with specific file verification
        $userDir = Join-Path $env:APPDATA "Code\User"
        $instructionsDir = Join-Path $userDir "instructions\terraform-azurerm"
        $hasInstructions = $false
        
        if (Test-Path $instructionsDir) {
            # Check for specific instruction files from config
            $instructionFiles = Get-ExpectedInstructionFiles
            
            # Check for the root copilot instructions file (directly in User directory)
            $rootCopilotFile = Join-Path $userDir "copilot-instructions.md"
            
            $missingFiles = @()
            
            # Check root copilot file
            $hasCopilotFile = Test-Path $rootCopilotFile
            if (-not $hasCopilotFile) {
                $missingFiles += "copilot-instructions.md"
            }

            # Check instruction files
            foreach ($expectedFile in $instructionFiles) {
                $filePath = Join-Path $instructionsDir $expectedFile
                if (Test-Path $filePath) {
                    $foundInstructionFiles++
                } else {
                    $missingFiles += $expectedFile
                }
            }
            
            $totalExpected = $instructionFiles.Count  # 13 instruction files (without copilot file)
            
            if ($foundInstructionFiles -gt 0) {
                $hasInstructions = $true
                Write-Host "Found instruction files ($foundInstructionFiles of $totalExpected expected files)" -ForegroundColor Green
                
                if ($missingFiles.Count -gt 0) {
                    Write-Host "  - Missing files: $($missingFiles -join ', ')" -ForegroundColor Yellow
                }
                
                # Check for core instruction file content
                if ($hasCopilotFile) {
                    $coreContent = Get-Content $rootCopilotFile -Raw -ErrorAction SilentlyContinue
                    if ($coreContent -and $coreContent.Length -gt 1000) {
                        Write-Host "  - Core copilot instruction file verified ($(($coreContent.Length / 1024).ToString('F1')) KB)" -ForegroundColor Cyan
                    } else {
                        Write-Host "  - Core copilot instruction file too small or empty" -ForegroundColor Yellow
                    }
                }
            } else {
                Write-Host "Instruction directory exists but no expected files found" -ForegroundColor Yellow
            }
        } else {
            Write-Host "No instruction files found" -ForegroundColor Red
        }
        
        # Check prompts files with specific file verification
        $promptsDir = Join-Path $userDir "prompts"
        $hasPrompts = $false
        
        if (Test-Path $promptsDir) {
            # Check for specific prompt files from config
            $expectedPrompts = Get-ExpectedPromptFiles
            
            $missingPrompts = @()
            
            foreach ($expectedPrompt in $expectedPrompts) {
                $filePath = Join-Path $promptsDir $expectedPrompt
                if (Test-Path $filePath) {
                    $foundPromptFiles++
                } else {
                    $missingPrompts += $expectedPrompt
                }
            }
            
            if ($foundPromptFiles -gt 0) {
                $hasPrompts = $true
                Write-Host "Found prompt files ($foundPromptFiles of $($expectedPrompts.Count) expected files)" -ForegroundColor Green
                
                if ($missingPrompts.Count -gt 0) {
                    Write-Host "  - Missing prompts: $($missingPrompts -join ', ')" -ForegroundColor Yellow
                }
            } else {
                Write-Host "Prompts directory exists but no expected files found" -ForegroundColor Yellow
            }
        } else {
            Write-Host "No prompt files found" -ForegroundColor Red
        }
        
        # Overall assessment with detailed criteria
        $installationFound = $hasVSCodeSettings -or $hasInstructions -or $hasPrompts
        
        # Calculate installation completeness score (count individual files)
        $score = 0
        $maxScore = 20  # 13 instruction files + 6 prompt files + 1 main copilot file
        
        # Count instruction files found
        $score += $foundInstructionFiles
        
        # Add main copilot file if found
        if ($hasCopilotFile) {
            $score += 1
        }
        
        # Count prompt files found  
        $score += $foundPromptFiles
        
        $completeness = ($score / $maxScore) * 100
        
        Write-Host "`nInstallation Completeness: $completeness% ($score of $maxScore components)" -ForegroundColor Cyan
        
        if ($installationFound) {
            if ($score -eq $maxScore) {
                Write-Host "`nAI installation verification PASSED (Complete)" -ForegroundColor Green
                Write-Host "All components are properly installed and verified" -ForegroundColor Green
                return $true
            } elseif ($score -ge 18) {  # 90% of 20 files
                Write-Host "`nAI installation verification PASSED (Mostly Complete)" -ForegroundColor Green
                Write-Host "Most components found, installation appears functional" -ForegroundColor Green
                return $true
            } else {
                Write-Host "`nAI installation verification PARTIAL" -ForegroundColor Yellow
                Write-Host "Some components found but installation may be incomplete" -ForegroundColor Yellow
                return $false
            }
        } else {
            Write-Host "`nAI installation verification FAILED" -ForegroundColor Red
            Write-Host "No AI enhancement components found" -ForegroundColor Red
            return $false
        }
        
    } catch {
        Write-Host "Verification error: $($_.Exception.Message)" -ForegroundColor Red
        return $false
    }
}

function Test-InstallationIntegrity {
    param(
        [Parameter(Mandatory = $true)]
        [string]$RepositoryPath,
        
        [Parameter(Mandatory = $false)]
        [string]$ManifestPath
    )
    
    Write-Host "Using hardcoded verification (manifest parameter ignored)..." -ForegroundColor Yellow
    
    # Run hardcoded verification and convert result to expected object structure
    $isValid = Test-HardcodedInstallationIntegrity -RepositoryPath $RepositoryPath
    
    # Return object structure expected by main script
    return [PSCustomObject]@{
        IsValid = $isValid
        Issues = @()  # Hardcoded verification doesn't track specific issues
    }
}

Export-ModuleMember -Function Test-HardcodedInstallationIntegrity, Test-InstallationIntegrity
