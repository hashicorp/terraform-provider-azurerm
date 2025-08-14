# Main AI Infrastructure Installer for Terraform AzureRM Provider
# Version: 1.0.0
# Description: Interactive installer for AI-powered development infrastructure

#requires -version 5.1

[CmdletBinding()]
param(
    [Parameter(HelpMessage = "Overwrite existing files without prompting")]
    [switch]${Auto-Approve},
    
    [Parameter(HelpMessage = "Show what would be done without making changes")]
    [switch]${Dry-Run},
    
    [Parameter(HelpMessage = "Check the current state of the workspace")]
    [switch]$Verify,
    
    [Parameter(HelpMessage = "Remove AI infrastructure from the workspace")]
    [switch]$Clean,
    
    [Parameter(HelpMessage = "Show detailed help information")]
    [switch]$Help
)

# Import required modules
$ModulesPath = Join-Path $PSScriptRoot "modules\powershell"
$RequiredModules = @(
    "ConfigParser",
    "FileOperations",
    "ValidationEngine", 
    "UI"
)

foreach ($module in $RequiredModules) {
    $modulePath = Join-Path $ModulesPath "$module.psm1"
    if (Test-Path $modulePath) {
        Import-Module $modulePath -Force
    } else {
        Write-Error "Required module '$module' not found at: $modulePath"
        exit 1
    }
}

# Global variables
$Global:InstallerConfig = @{
    Version = "1.0.0"
    Branch = "exp/terraform_copilot"
    StartTime = Get-Date
    AutoApprove = ${Auto-Approve}.IsPresent
    DryRun = ${Dry-Run}.IsPresent
    Verify = $Verify.IsPresent
    Clean = $Clean.IsPresent
}

function Initialize-Installer {
    <#
    .SYNOPSIS
    Initialize the installer environment and validate prerequisites
    #>
    
    Show-Banner -Title "Terraform AzureRM Provider - AI Infrastructure Installer" -Version $Global:InstallerConfig.Version
    
    Write-Host "> Initializing AI Infrastructure Installer..." -ForegroundColor "Cyan"
    Write-Host "   Branch: $($Global:InstallerConfig.Branch)" -ForegroundColor "Gray"
    Write-Host "   Mode: $(if ($Global:InstallerConfig.Verify) { 'Verification Only' } elseif ($Global:InstallerConfig.DryRun) { 'Dry Run' } elseif ($Global:InstallerConfig.Clean) { 'Clean Mode' } else { 'Full Installation' })" -ForegroundColor "Gray"
    Write-Host ""
    
    # Run pre-installation validation
    $validation = Test-PreInstallation
    
    Show-ValidationResults -Results $validation -ShowDetails $true
    
    if (-not $validation.OverallValid) {
        Write-Error "Pre-installation validation failed. Please resolve the issues above before continuing."
        
        if (-not $Global:InstallerConfig.AutoApprove) {
            $showDetails = Confirm-UserAction -Message "Would you like to see detailed recommendations?" -DefaultYes $true
            if ($showDetails) {
                $report = Get-ValidationReport -Branch $Global:InstallerConfig.Branch
                
                if ($report.Recommendations.Count -gt 0) {
                    Write-Section "Recommendations"
                    foreach ($recommendation in $report.Recommendations) {
                        Write-Host "  - $recommendation" -ForegroundColor "Yellow"
                    }
                    Write-Host ""
                }
            }
        }
        
        exit 1
    }
    
    Show-Success "Pre-installation validation completed successfully!"
    
    return $validation
}

function Invoke-Installation {
    <#
    .SYNOPSIS
    Execute the main installation process
    #>
    
    $startTime = Get-Date
    $statistics = @{
        "Files Downloaded" = 0
        "Files Skipped" = 0
        "Files Failed" = 0
        "Total Size" = 0
    }
    
    try {
        # Get installation configuration
        $config = Get-InstallationConfig -Branch $Global:InstallerConfig.Branch
        
        Write-Section "Installing AI Infrastructure"
        Write-Host "Installing $($config.Files.Count) files from branch '$($Global:InstallerConfig.Branch)'..." -ForegroundColor "White"
        Write-Host ""
        
        $fileIndex = 0
        foreach ($filePath in $config.Files.Keys) {
            $fileIndex++
            $fileInfo = $config.Files[$filePath]
            
            # Check if file exists and AutoApprove flag
            if ((Test-Path $filePath) -and -not $Global:InstallerConfig.AutoApprove) {
                Show-FileOperationStatus -Operation "Skipped" -FilePath $filePath -Status "Skipped"
                $statistics["Files Skipped"]++
                continue
            }
            
            # Show progress
            Show-InstallationProgress -CurrentFile $filePath -Current $fileIndex -Total $config.Files.Count -Status "Installing"
            
            # Download file
            try {
                $downloadResult = Get-FileFromGitHub -FilePath $fileInfo.GitHubPath -LocalPath $filePath -Branch $Global:InstallerConfig.Branch
                
                if ($downloadResult.Success) {
                    $statistics["Files Downloaded"]++
                    $statistics["Total Size"] += $downloadResult.Size
                    
                    Write-Host ""  # New line after progress
                    Show-FileOperationStatus -Operation "Downloaded" -FilePath $filePath -Status "Success" -Size $downloadResult.Size
                } else {
                    $statistics["Files Failed"]++
                    
                    Write-Host ""  # New line after progress
                    Show-FileOperationStatus -Operation "Download" -FilePath $filePath -Status "Failed" -ErrorMessage $downloadResult.Error
                    
                    if ($fileInfo.Required) {
                        throw "Failed to download required file: $filePath"
                    }
                }
            }
            catch {
                $statistics["Files Failed"]++
                
                Write-Host ""  # New line after progress
                Show-FileOperationStatus -Operation "Download" -FilePath $filePath -Status "Failed" -ErrorMessage $_.Exception.Message
                
                if ($fileInfo.Required) {
                    throw
                }
            }
        }
        
        # Clear progress line
        Write-Host ""
        
        # Update .gitignore if requested
        if (-not $Global:InstallerConfig.AutoApprove) {
            $currentGitIgnore = Get-GitIgnoreStatus
            if (-not $currentGitIgnore.HasAIEntries) {
                $addToGitIgnore = Confirm-UserAction -Message "Would you like to add AI infrastructure files to .gitignore?" -DefaultYes $true
                
                if ($addToGitIgnore) {
                    Update-GitIgnore
                    Show-Success "Updated .gitignore with AI infrastructure entries."
                }
            }
        }
        
        $duration = (Get-Date) - $startTime
        
        Show-Summary -Operation "Installation" -Statistics $statistics -Duration $duration -Success ($statistics["Files Failed"] -eq 0)
        
        return @{
            Success = ($statistics["Files Failed"] -eq 0)
            Statistics = $statistics
            Duration = $duration
        }
    }
    catch {
        $duration = (Get-Date) - $startTime
        
        Write-Error "Installation failed: $($_.Exception.Message)"
        Show-Summary -Operation "Installation" -Statistics $statistics -Duration $duration -Success $false
        
        return @{
            Success = $false
            Error = $_.Exception.Message
            Statistics = $statistics
            Duration = $duration
        }
    }
}

function Invoke-PostInstallationValidation {
    <#
    .SYNOPSIS
    Validate the installation and show results
    #>
    
    Write-Section "Post-Installation Validation"
    
    # Run post-installation validation
    $validation = Test-PostInstallation -Branch $Global:InstallerConfig.Branch
    
    if ($validation.OverallValid) {
        Show-Success "All files installed and validated successfully!"
        
        Write-Host "Installation Summary:" -ForegroundColor "Cyan"
        Write-Host "   Files Installed: $($validation.FilesInstalled)/$($validation.TotalFiles)" -ForegroundColor "White"
        Write-Host "   Missing Files: $($validation.MissingFiles.Count)" -ForegroundColor $(if ($validation.MissingFiles.Count -eq 0) { "Green" } else { "Red" })
            
            if ($validation.MissingFiles.Count -gt 0) {
                Write-Host "   Missing:" -ForegroundColor "Red"
            foreach ($missingFile in $validation.MissingFiles) {
                Write-Host "     - $missingFile" -ForegroundColor "Red"
            }
        }
    } else {
        Write-Error "Post-installation validation failed!"
        
        Write-Host "Issues found:" -ForegroundColor "Red"
        foreach ($filePath in $validation.FileDetails.Keys) {
            $fileValidation = $validation.FileDetails[$filePath]
            if (-not $fileValidation.Exists -and $fileValidation.Required) {
                Write-Host "  - Missing required file: $filePath" -ForegroundColor "Red"
            }
            elseif ($fileValidation.Exists -and $fileValidation.Integrity -and -not $fileValidation.Integrity.Valid) {
                Write-Host "  - Corrupted file: $filePath" -ForegroundColor "Red"
            }
        }
    }
    
    Write-Host ""
    
    return $validation
}

function Show-NextSteps {
    <#
    .SYNOPSIS
    Display next steps after successful installation
    #>
    
    Write-Section "Next Steps"
    
    Write-Host "Your AI infrastructure is now installed! Here's what you can do next:" -ForegroundColor "White"
    Write-Host ""
    
    Write-Host "1. Restart VS Code" -ForegroundColor "Cyan"
    Write-Host "   - Close and reopen VS Code to load the new configuration" -ForegroundColor "Gray"
    Write-Host "   - The AI instructions will be automatically loaded" -ForegroundColor "Gray"
    Write-Host ""
    
    Write-Host "2. Verify GitHub Copilot" -ForegroundColor "Cyan"
    Write-Host "   - Ensure GitHub Copilot extension is installed and authenticated" -ForegroundColor "Gray"
    Write-Host "   - Check that Copilot is active in VS Code status bar" -ForegroundColor "Gray"
    Write-Host ""
    
    Write-Host "3.  Test AI Features" -ForegroundColor "Cyan"
    Write-Host "   - Open a Go file in the internal/services directory" -ForegroundColor "Gray"
    Write-Host "   - Try AI-powered code completion and suggestions" -ForegroundColor "Gray"
    Write-Host "   - Use GitHub Copilot Chat for code explanations" -ForegroundColor "Gray"
    Write-Host ""
    
    Write-Host "4. Explore AI Instructions" -ForegroundColor "Cyan"
    Write-Host "   - Check .github/instructions/ for comprehensive coding guidelines" -ForegroundColor "Gray"
    Write-Host "   - Review .github/copilot-instructions.md for AI-specific guidance" -ForegroundColor "Gray"
    Write-Host ""
    
    Write-Host "5. Validate Installation" -ForegroundColor "Cyan"
    Write-Host "   - Run: .\Install-AIInfrastructure.ps1 -ValidateOnly" -ForegroundColor "Gray"
    Write-Host "   - Check AI infrastructure health anytime" -ForegroundColor "Gray"
    Write-Host ""
    
    Show-Info "For troubleshooting and advanced configuration, see the documentation in .github/AIinstaller/README.md"
}

function Invoke-CleanInstallation {
    <#
    .SYNOPSIS
    Remove all AI infrastructure files from the workspace
    #>
    
    $result = @{ Success = $false; Error = $null; RemovedFiles = @() }
    
    try {
        $config = Get-InstallationConfig -Branch $Global:InstallerConfig.Branch
        $removedCount = 0
        
        foreach ($filePath in $config.Files.Keys) {
            if (Test-Path $filePath) {
                try {
                    Remove-Item -Path $filePath -Force -ErrorAction Stop
                    $result.RemovedFiles += $filePath
                    $removedCount++
                    
                    Show-FileOperationStatus -Operation "Removed" -FilePath $filePath -Status "Success"
                } catch {
                    Show-FileOperationStatus -Operation "Failed to remove" -FilePath $filePath -Status "Error"
                }
            }
        }
        
        # Remove empty directories
        $dirsToCheck = @(".github\instructions", ".github\prompts", ".vscode")
        foreach ($dir in $dirsToCheck) {
            if ((Test-Path $dir) -and ((Get-ChildItem $dir -Force | Measure-Object).Count -eq 0)) {
                Remove-Item -Path $dir -Force -ErrorAction SilentlyContinue
                Show-FileOperationStatus -Operation "Removed" -FilePath $dir -Status "Success"
            }
        }
        
        Write-Host ""
        Show-Success "Removed $removedCount AI infrastructure files"
        
        $result.Success = $true
    } catch {
        $result.Error = $_.Exception.Message
    }
    
    return $result
}

function Show-DryRunPreview {
    <#
    .SYNOPSIS
    Show what would be installed without making changes
    #>
    
    $config = Get-InstallationConfig -Branch $Global:InstallerConfig.Branch
    
    Write-Host "The following files would be installed:" -ForegroundColor "White"
    Write-Host ""
    
    foreach ($filePath in $config.Files.Keys) {
        $fileInfo = $config.Files[$filePath]
        $status = if (Test-Path $filePath) { "Overwrite" } else { "Create" }
        $color = if ($status -eq "Create") { "Green" } else { "Yellow" }
        
        Write-Host "  [$status] " -ForegroundColor $color -NoNewline
        Write-Host $filePath -ForegroundColor "White"
        Write-Host "           $($fileInfo.Description)" -ForegroundColor "Gray"
    }
    
    Write-Host ""
    Write-Host "Total files: $($config.Files.Count)" -ForegroundColor "Cyan"
    Write-Host "Branch: $($Global:InstallerConfig.Branch)" -ForegroundColor "Gray"
    Write-Host ""
    Write-Host "To perform the actual installation, run without -DryRun" -ForegroundColor "Yellow"
}

function Main {
    <#
    .SYNOPSIS
    Main installer entry point
    #>
    
    try {
        # Show help if requested
        if ($Help) {
            Show-Help
            return
        }
        
        # Initialize installer
        Initialize-Installer
        
        # Handle Clean mode
        if ($Global:InstallerConfig.Clean) {
            Write-Host "Cleaning AI infrastructure..." -ForegroundColor "Yellow"
            $cleanResult = Invoke-CleanInstallation
            if ($cleanResult.Success) {
                Show-Success "AI infrastructure successfully removed."
            } else {
                Write-Error "Clean operation failed: $($cleanResult.Error)"
                exit 1
            }
            return
        }
        
        # Handle Verify mode
        if ($Global:InstallerConfig.Verify) {
            $infraStatus = Test-AIInfrastructure -Branch $Global:InstallerConfig.Branch
            Show-AIInfrastructureStatus -Status $infraStatus
            
            Show-Success "Verification completed."
            return
        }
        
        # Handle Dry Run mode
        if ($Global:InstallerConfig.DryRun) {
            Write-Host "DRY RUN MODE - No changes will be made" -ForegroundColor "Yellow"
            Write-Host ""
            Show-DryRunPreview
            return
        }
        
        # Confirm installation in interactive mode (unless AutoApprove is set)
        if (-not $Global:InstallerConfig.AutoApprove) {
            $proceed = Confirm-UserAction -Message "Ready to install AI infrastructure. Continue?" -DefaultYes $true
            if (-not $proceed) {
                Write-Host "Installation cancelled by user." -ForegroundColor "Yellow"
                return
            }
        }
        
        # Execute installation
        $installResult = Invoke-Installation
        
        if ($installResult.Success) {
            # Run post-installation validation
            $postValidation = Invoke-PostInstallationValidation
            
            if ($postValidation.OverallValid) {
                Show-NextSteps
                
                $totalDuration = (Get-Date) - $Global:InstallerConfig.StartTime
                $durationFormatted = "{0:N1} seconds" -f $totalDuration.TotalSeconds
                
                Show-Success "AI infrastructure installation completed successfully in $durationFormatted!"
                
                if (-not $Global:InstallerConfig.AutoApprove) {
                    Wait-ForUser "Press any key to exit..."
                }
                exit 0
            } else {
                Write-Error "Installation completed but post-validation failed. Some files may be missing or corrupted."
                exit 1
            }
        } else {
            Write-Error "Installation failed: $($installResult.Error)"
            exit 1
        }
    }
    catch {
        $errorMessage = $_.Exception.Message
        Write-Error "Installer failed with error: $errorMessage"
        
        Write-Host "Stack trace:" -ForegroundColor "Red"
        Write-Host $_.ScriptStackTrace -ForegroundColor "Gray"
        
        exit 1
    }
}

# Execute main function
Main
