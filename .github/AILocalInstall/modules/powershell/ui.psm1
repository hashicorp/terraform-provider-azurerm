# User Interface Module  
# Handles user interaction, progress reporting, and output formatting using real patterns

function Show-WelcomeBanner {
    <#
    .SYNOPSIS
        Displays the installation welcome banner
    #>
    param([string]$Version = "2.0.0")
    
    Write-Host ""
    Write-Host "============================================================================" -ForegroundColor Cyan
    Write-Host " Terraform AzureRM Provider - AI Setup Installer v$Version" -ForegroundColor White
    Write-Host "============================================================================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host " This installer will configure VS Code with:" -ForegroundColor Gray
    Write-Host "   - AI-powered coding instructions for AzureRM provider development" -ForegroundColor Gray
    Write-Host "   - GitHub Copilot prompt templates and examples" -ForegroundColor Gray
    Write-Host "   - Optimized VS Code settings for Terraform development" -ForegroundColor Gray
    Write-Host ""
}

function Show-InstallationSummary {
    <#
    .SYNOPSIS
        Displays a comprehensive installation summary
    #>
    param([hashtable]$Results)
    
    Write-Host ""
    Write-Host "============================================================================" -ForegroundColor Cyan
    Write-Host " Installation Summary" -ForegroundColor White
    Write-Host "============================================================================" -ForegroundColor Cyan
    Write-Host ""
    
    if ($Results.Success) {
        Write-Host " Status: INSTALLATION SUCCESSFUL" -ForegroundColor Green
        Write-Host ""
        Write-Host " Components Installed:" -ForegroundColor White
        
        # Show found vs expected counts
        $instructionStatus = if ($Results.InstructionFiles.Count -eq $Results.ExpectedInstructionFiles.Count) { "Green" } else { "Yellow" }
        $promptStatus = if ($Results.PromptFiles.Count -eq $Results.ExpectedPromptFiles.Count) { "Green" } else { "Yellow" }
        $mainStatus = if ($Results.MainFiles.Count -eq $Results.ExpectedMainFiles.Count) { "Green" } else { "Yellow" }
        $settingsStatus = if ($Results.SettingsConfigured) { "Green" } else { "Red" }
        
        Write-Host "   - Instruction Files: $($Results.InstructionFiles.Count) of $($Results.ExpectedInstructionFiles.Count) files found" -ForegroundColor $instructionStatus
        Write-Host "   - Prompt Files: $($Results.PromptFiles.Count) of $($Results.ExpectedPromptFiles.Count) files found" -ForegroundColor $promptStatus  
        Write-Host "   - Main Files: $($Results.MainFiles.Count) of $($Results.ExpectedMainFiles.Count) files found" -ForegroundColor $mainStatus
        Write-Host "   - VS Code Settings: $(if ($Results.SettingsConfigured) { 'Correctly configured' } else { 'Not configured' })" -ForegroundColor $settingsStatus
    } else {
        # Determine if this is no installation or partial installation
        $totalFound = $Results.InstructionFiles.Count + $Results.PromptFiles.Count + $Results.MainFiles.Count
        $totalExpected = $Results.ExpectedInstructionFiles.Count + $Results.ExpectedPromptFiles.Count + $Results.ExpectedMainFiles.Count
        
        if ($totalFound -eq 0 -and -not $Results.SettingsConfigured) {
            Write-Host " Status: NO INSTALLATION FOUND" -ForegroundColor Red
        } else {
            Write-Host " Status: INSTALLATION COMPLETED WITH ISSUES" -ForegroundColor Yellow
        }
        Write-Host ""
        
        # Show partial installation status
        Write-Host " Components Found:" -ForegroundColor White
        
        # Show each component with appropriate color based on actual status
        $instructionStatus = if ($Results.InstructionFiles.Count -eq $Results.ExpectedInstructionFiles.Count) { "Green" } else { "Red" }
        $promptStatus = if ($Results.PromptFiles.Count -eq $Results.ExpectedPromptFiles.Count) { "Green" } else { "Red" }
        $mainStatus = if ($Results.MainFiles.Count -eq $Results.ExpectedMainFiles.Count) { "Green" } else { "Red" }
        $settingsStatus = if ($Results.SettingsConfigured) { "Green" } else { "Red" }
        
        Write-Host "   - Instruction Files: $($Results.InstructionFiles.Count) of $($Results.ExpectedInstructionFiles.Count) files found" -ForegroundColor $instructionStatus
        Write-Host "   - Prompt Files: $($Results.PromptFiles.Count) of $($Results.ExpectedPromptFiles.Count) files found" -ForegroundColor $promptStatus  
        Write-Host "   - Main Files: $($Results.MainFiles.Count) of $($Results.ExpectedMainFiles.Count) files found" -ForegroundColor $mainStatus
        Write-Host "   - VS Code Settings: $(if ($Results.SettingsConfigured) { 'Correctly configured' } else { 'Not configured' })" -ForegroundColor $settingsStatus
        Write-Host ""
        
        Write-Host " Errors encountered: $($Results.Errors.Count)" -ForegroundColor Red
        foreach ($errorMsg in $Results.Errors) {
            Write-Host "   - $errorMsg" -ForegroundColor Red
        }
    }
    
    Write-Host ""
}

function Show-RepositoryInfo {
    <#
    .SYNOPSIS
        Displays repository validation information
    #>
    param(
        [string]$RepositoryPath,
        [hashtable]$ValidationResults
    )
    
    Write-Host ""
    Write-Host "Repository Information:" -ForegroundColor White
    Write-Host "  Path: $RepositoryPath" -ForegroundColor Gray
    
    if ($ValidationResults.IsValidRepository) {
        Write-Host "  Status: Valid Terraform AzureRM Provider repository" -ForegroundColor Green
    } else {
        Write-Host "  Status: Invalid repository or missing files" -ForegroundColor Red
    }
    Write-Host ""
}

function Request-UserConfirmation {
    <#
    .SYNOPSIS
        Requests user confirmation for installation
    #>
    param(
        [string]$Message = "Proceed with installation?",
        [switch]$Force
    )
    
    if ($Force) {
        Write-Host "Force mode enabled - proceeding automatically" -ForegroundColor Yellow
        return $true
    }
    
    Write-Host ""
    Write-Host $Message -ForegroundColor White
    Write-Host "  [Y] Yes  [N] No  (Default: Y): " -NoNewline -ForegroundColor Yellow
    
    $response = Read-Host
    if ($response -eq "" -or $response.ToLower() -eq "y" -or $response.ToLower() -eq "yes") {
        return $true
    }
    
    return $false
}

function Show-NextSteps {
    <#
    .SYNOPSIS
        Shows next steps after installation
    #>
    Write-Host ""
    Write-Host "============================================================================" -ForegroundColor Cyan
    Write-Host " Next Steps" -ForegroundColor White
    Write-Host "============================================================================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host " 1. Restart VS Code to ensure all settings take effect" -ForegroundColor White
    Write-Host " 2. Open the Terraform AzureRM Provider repository in VS Code" -ForegroundColor White
    Write-Host " 3. Try using GitHub Copilot with the new AI instructions!" -ForegroundColor White
    Write-Host ""
    Write-Host " The AI is now configured to help you with:" -ForegroundColor Gray
    Write-Host "   - Writing resource implementations following provider patterns" -ForegroundColor Gray
    Write-Host "   - Creating comprehensive acceptance tests" -ForegroundColor Gray
    Write-Host "   - Generating proper documentation" -ForegroundColor Gray
    Write-Host "   - Following Azure SDK integration best practices" -ForegroundColor Gray
    Write-Host ""
    Write-Host " Happy coding!" -ForegroundColor Green
    Write-Host ""
}

function Write-InstallationProgress {
    <#
    .SYNOPSIS
        Shows progress during installation (simplified - no progress bar)
    #>
    param(
        [string]$Stage,
        [int]$StageNumber,
        [int]$TotalStages
    )
    
    # Calculate percentage and show status with percentage
    if ($TotalStages -eq 0) { $TotalStages = 1 }
    $percent = [math]::Round(($StageNumber / $TotalStages) * 100)
    
    Write-StatusMessage "$Stage - Step $StageNumber of $TotalStages ($percent%)" "Info"
}

function Show-Help {
    <#
    .SYNOPSIS
        Displays help information for the installer
    #>
    Write-Host "Terraform AzureRM Provider AI Setup" -ForegroundColor Cyan
    Write-Host "====================================================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "USAGE:" -ForegroundColor Yellow
    Write-Host "  .\install-copilot-setup.ps1 [OPTIONS]"
    Write-Host ""
    Write-Host "OPTIONS:" -ForegroundColor Yellow
    Write-Host "  -RepositoryPath <path>    Path to terraform-provider-azurerm repository"
    Write-Host "  -Clean                    Remove all installed files and restore backups"
    Write-Host "  -Auto-Approve             Skip interactive approval prompts"
    Write-Host "  -Verify                   Run verification only without installing"
    Write-Host "  -Help                     Show this help message"
    Write-Host ""
    Write-Host "EXAMPLES:" -ForegroundColor Yellow
    Write-Host "  .\install-copilot-setup.ps1                                  # Auto-discover repository"
    Write-Host "  .\install-copilot-setup.ps1 -RepositoryPath C:\path\to\repo  # Use specific path"
    Write-Host "  .\install-copilot-setup.ps1 -Auto-Approve                    # Non-interactive install"
    Write-Host "  .\install-copilot-setup.ps1 -Clean                           # Remove installation"
    Write-Host "  .\install-copilot-setup.ps1 -Verify                          # Verify current installation"
    Write-Host ""
}

function Write-StatusMessage {
    <#
    .SYNOPSIS
        Writes formatted status messages
    #>
    param(
        [string]$Message,
        [ValidateSet("Info", "Success", "Warning", "Error")]
        [string]$Type = "Info"
    )
    
    $colors = @{
        "Info" = "Cyan"
        "Success" = "Green"
        "Warning" = "Yellow"
        "Error" = "Red"
    }
    
    $icons = @{
        "Info" = "[INFO]"
        "Success" = "[SUCCESS]"
        "Warning" = "[WARNING]"
        "Error" = "[ERROR]"
    }
    
    Write-Host "$($icons[$Type]) $Message" -ForegroundColor $colors[$Type]
}

Export-ModuleMember -Function Show-WelcomeBanner, Show-InstallationSummary, Show-RepositoryInfo, Request-UserConfirmation, Show-NextSteps, Write-InstallationProgress, Show-Help, Write-StatusMessage
