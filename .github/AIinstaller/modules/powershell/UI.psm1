# UI Module for Terraform AzureRM Provider AI Setup
# Cleaned version with only used functions

#region Public Functions

function Write-Separator {
    <#
    .SYNOPSIS
    Display a separator line with consistent formatting
    
    .DESCRIPTION
    Displays a colored separator line for visual separation in UI output.
    Matches the bash script's print_separator() function behavior.
    
    .PARAMETER Length
    The length of the separator line. Defaults to 60 characters.
    
    .PARAMETER Color
    The color of the separator line. Defaults to Cyan.
    
    .PARAMETER Character
    The character to use for the separator. Defaults to "=".
    #>
    param(
        [int]$Length = 60,
        [string]$Color = "Cyan",
        [string]$Character = "="
    )
    
    Write-Host $($Character * $Length) -ForegroundColor $Color
}

function Write-Header {
    <#
    .SYNOPSIS
    Display the main application header
    #>
    param(
        [string]$Title = "Terraform AzureRM Provider - AI Infrastructure Installer",
        [string]$Version = "1.0.0"
    )
    
    Write-Host ""
    Write-Separator
    Write-Host " $Title" -ForegroundColor Cyan
    Write-Host " Version: $Version" -ForegroundColor Cyan
    Write-Separator
    Write-Host ""
}

function Write-Success {
    <#
    .SYNOPSIS
    Display success message in green
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Message,
        
        [string]$Prefix = "[SUCCESS]"
    )
    
    Write-Host "$Prefix $Message" -ForegroundColor Green
}

function Format-AlignedLabel {
    <#
    .SYNOPSIS
    Format a label with dynamic spacing to align with branch detection labels
    .DESCRIPTION
    Returns a complete formatted string "LABEL    : " that aligns perfectly with 
    branch detection output like "SOURCE BRANCH DETECTED: "
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Label,
        
        [string]$CurrentBranchType = "source"
    )
    
    # Define branch labels exactly as they appear in Show-BranchDetection (including ": ")
    $branchLabels = @{
        "source"  = "SOURCE BRANCH DETECTED: "
        "feature" = "FEATURE BRANCH DETECTED: " 
        "unknown" = "UNKNOWN BRANCH: "
    }
    
    # Get the current branch label length (including ": ")
    $currentBranchLabelLength = $branchLabels[$CurrentBranchType].Length
    
    # Calculate spacing needed to align the colons perfectly
    $spacingNeeded = $currentBranchLabelLength - $Label.Length - 2
    $spacing = " " * [Math]::Max(0, $spacingNeeded)
    
    # Return the complete formatted string with colon - centralized formatting
    return "${Label}${spacing}: "
}

function Show-BranchDetection {
    <#
    .SYNOPSIS
    Display branch detection result with workspace path and appropriate styling
    #>
    param(
        [Parameter(Mandatory)]
        [string]$BranchName,
        
        [ValidateSet("source", "feature", "unknown")]
        [string]$BranchType = "feature"
    )
    
    # Display branch detection
    switch ($BranchType) {
        "source" {
            Write-Host "SOURCE BRANCH DETECTED: " -ForegroundColor "Cyan" -NoNewline
            Write-Host "$BranchName" -ForegroundColor "Green"
        }
        "feature" {
            Write-Host "FEATURE BRANCH DETECTED: " -ForegroundColor "Cyan" -NoNewline
            Write-Host "$BranchName" -ForegroundColor "Yellow"
        }
        "unknown" {
            Write-Host "UNKNOWN BRANCH: " -ForegroundColor "Red" -NoNewline
            Write-Host "$BranchName" -ForegroundColor "Gray"
        }
    }
    
    # Display workspace path with consistent formatting
    if ($Global:WorkspaceRoot) {
        $formattedWorkspaceLabel = Format-AlignedLabel -Label "WORKSPACE" -CurrentBranchType $BranchType
        Write-Host $formattedWorkspaceLabel -ForegroundColor Cyan -NoNewline
        Write-Host $Global:WorkspaceRoot -ForegroundColor Green
        Write-Host ""
    }
}

function Write-WarningMessage {
    <#
    .SYNOPSIS
    Display warning message in yellow
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Message,
        
        [string]$Prefix = "[WARNING]"
    )
    
    Write-Host "$Prefix $Message" -ForegroundColor Yellow
}

function Write-ErrorMessage {
    <#
    .SYNOPSIS
    Display error message in red
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Message,
        
        [string]$Prefix = "[ERROR]"
    )
    
    Write-Host "$Prefix $Message" -ForegroundColor Red
}

function Write-Section {
    <#
    .SYNOPSIS
    Display section header
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Title
    )
    
    Write-Host ""
    Write-Host "[SECTION] $Title" -ForegroundColor Cyan
    Write-Host ("-" * ($Title.Length + 10)) -ForegroundColor Cyan
}

function Write-VerboseMessage {
    <#
    .SYNOPSIS
    Display verbose information
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Message
    )
    
    Write-Host "[VERBOSE] $Message" -ForegroundColor DarkGray
}

function Write-OperationStatus {
    <#
    .SYNOPSIS
    Write operation status message with color coding
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Message,
        
        [ValidateSet('Info', 'Success', 'Warning', 'Error', 'Progress')]
        [string]$Type = 'Info'
    )
    
    switch ($Type) {
        'Info'     { Write-Host "[INFO] $Message" -ForegroundColor Cyan }
        'Success'  { Write-Host "[SUCCESS] $Message" -ForegroundColor Green }
        'Warning'  { Write-Host "[WARNING] $Message" -ForegroundColor Yellow }
        'Error'    { Write-Host "[ERROR] $Message" -ForegroundColor Red }
        'Progress' { Write-Host "[PROGRESS] $Message" -ForegroundColor Blue }
    }
}

function Show-Help {
    <#
    .SYNOPSIS
    Display contextual help information based on branch type
    #>
    param(
        [string]$BranchName = "",
        [string]$BranchType = "unknown",
        [switch]$SkipHeader
    )
    
    # Only show header if not already shown by caller
    if (-not $SkipHeader) {
        Write-Header -Title "Terraform AzureRM Provider - AI Infrastructure Installer" -Version $Global:InstallerConfig.Version
    }
    
    # Show current branch context if available
    if ($BranchName) {
        Show-BranchDetection -BranchName $BranchName -BranchType $BranchType
        Write-Separator
        Write-Host ""
    }
    
    Write-Host "DESCRIPTION:" -ForegroundColor Cyan
    Write-Host "  Interactive installer for AI-powered development infrastructure that enhances"
    Write-Host "  GitHub Copilot with Terraform-specific knowledge, patterns, and best practices."
    Write-Host ""
    
    # Dynamic options and examples based on branch type
    switch ($BranchType) {
        "source" {
            Show-SourceBranchHelp
        }
        "feature" {
            Show-FeatureBranchHelp
        }
        default {
            Show-UnknownBranchHelp
        }
    }
    
    Write-Host ""
    Write-Host "For more information, visit: https://github.com/hashicorp/terraform-provider-azurerm" -ForegroundColor Cyan
    Write-Host ""
}

function Show-SourceBranchHelp {
    <#
    .SYNOPSIS
    Display help specific to source branch operations
    #>
    param(
        [string]$BranchName = "",
        [string]$WorkspacePath = ""
    )
    
    # Show branch and workspace context if provided
    if ($BranchName) {
        Show-BranchDetection -BranchName $BranchName -BranchType "source"
        Write-Separator
        Write-Host ""
    }
    
    Write-Host "USAGE:" -ForegroundColor Cyan
    Write-Host "  .\install-copilot-setup.ps1 [OPTIONS]"
    Write-Host ""
    Write-Host "AVAILABLE OPTIONS:" -ForegroundColor Cyan
    Write-Host "  -Bootstrap        Copy installer to user profile (~\.terraform-ai-installer\)"
    Write-Host "  -Verify           Check current workspace status and validate setup"
    Write-Host "  -Help             Show this help information"
    Write-Host ""
    Write-Host "EXAMPLES:" -ForegroundColor Cyan
    Write-Host "  Bootstrap installer:"
    Write-Host "    .\install-copilot-setup.ps1 -Bootstrap"
    Write-Host ""
    Write-Host "  Verify setup:"
    Write-Host "    .\install-copilot-setup.ps1 -Verify"
    Write-Host ""
    Write-Host "NEXT STEPS:" -ForegroundColor Cyan
    Write-Host "  1. Run bootstrap to set up the installer in your user profile"
    Write-Host "  2. Switch to your feature branch: git checkout feature/your-branch-name"
    Write-Host "  3. Run installer from user profile to install AI infrastructure"
}

function Show-FeatureBranchHelp {
    <#
    .SYNOPSIS
    Display help specific to feature branch operations
    #>
    Write-Host "USAGE:" -ForegroundColor Cyan
    Write-Host "  & `"~\.terraform-ai-installer\install-copilot-setup.ps1`" [OPTIONS]"
    Write-Host ""
    
    Write-Host "AVAILABLE OPTIONS:" -ForegroundColor Cyan
    Write-Host "  -RepoDirectory    Repository path for git operations (when running from user profile)"
    Write-Host "  -Auto-Approve     Overwrite existing files without prompting"
    Write-Host "  -Dry-Run          Show what would be done without making changes"
    Write-Host "  -Verify           Check current workspace status and validate setup"
    Write-Host "  -Clean            Remove AI infrastructure from workspace"
    Write-Host "  -Help             Show this help information"
    Write-Host ""
    
    Write-Host "EXAMPLES:" -ForegroundColor Cyan
    Write-Host "  Install AI infrastructure:"
    Write-Host "    & `"~\.terraform-ai-installer\install-copilot-setup.ps1`" -RepoDirectory `"C:\path\to\repo`""
    Write-Host ""
    Write-Host "  Dry run (preview changes):"
    Write-Host "    & `"~\.terraform-ai-installer\install-copilot-setup.ps1`" -RepoDirectory `"C:\path\to\repo`" -Dry-Run"
    Write-Host ""
    Write-Host "  Auto-approve installation:"
    Write-Host "    & `"~\.terraform-ai-installer\install-copilot-setup.ps1`" -RepoDirectory `"C:\path\to\repo`" -Auto-Approve"
    Write-Host ""
    Write-Host "  Clean removal:"
    Write-Host "    & `"~\.terraform-ai-installer\install-copilot-setup.ps1`" -RepoDirectory `"C:\path\to\repo`" -Clean"
    Write-Host ""
    
    Write-Host "WORKFLOW:" -ForegroundColor Cyan
    Write-Host "  1. Ensure you're on your feature branch"
    Write-Host "  2. Run installer from user profile (installed via bootstrap)"
    Write-Host "  3. Start developing with enhanced GitHub Copilot AI features"
    Write-Host "  4. Use -Clean to remove AI infrastructure when done"
    Write-Host ""
}

function Show-UnknownBranchHelp {
    <#
    .SYNOPSIS
    Display generic help when branch type cannot be determined
    #>
    Write-Host "USAGE:" -ForegroundColor Cyan
    Write-Host "  .\install-copilot-setup.ps1 [OPTIONS]"
    Write-Host ""
    
    Write-Host "ALL OPTIONS:" -ForegroundColor Cyan
    Write-Host "  -Bootstrap        Copy installer to user profile (source branch only)"
    Write-Host "  -RepoDirectory    Repository path for git operations"
    Write-Host "  -Auto-Approve     Overwrite existing files without prompting"
    Write-Host "  -Dry-Run          Show what would be done without making changes"
    Write-Host "  -Verify           Check current workspace status and validate setup"
    Write-Host "  -Clean            Remove AI infrastructure from workspace"
    Write-Host "  -Help             Show this help information"
    Write-Host ""
    
    Write-Host "BRANCH-SPECIFIC WORKFLOW:" -ForegroundColor Cyan
    Write-Host "  Source Branch: Use -Bootstrap and -Verify only"
    Write-Host "  Feature Branch: Use all options except -Bootstrap"
    Write-Host ""
    
    Write-Host "EXAMPLES:" -ForegroundColor Cyan
    Write-Host "  Check current state:"
    Write-Host "    .\install-copilot-setup.ps1 -Verify"
    Write-Host ""
    Write-Host "  Get contextual help:"
    Write-Host "    .\install-copilot-setup.ps1 -Help"
    Write-Host ""
}

function Get-UserInput {
    <#
    .SYNOPSIS
    Get user input with validation
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Prompt,
        
        [string]$DefaultValue = "",
        
        [string[]]$ValidValues = @(),
        
        [switch]$Secure
    )
    
    do {
        $promptText = $Prompt
        if ($DefaultValue) {
            $promptText += " (default: $DefaultValue)"
        }
        $promptText += ": "
        
        if ($Secure) {
            $userInput = Read-Host -Prompt $promptText -AsSecureString
            $userInput = [System.Runtime.InteropServices.Marshal]::PtrToStringAuto([System.Runtime.InteropServices.Marshal]::SecureStringToBSTR($userInput))
        } else {
            $userInput = Read-Host -Prompt $promptText
        }
        
        if ([string]::IsNullOrWhiteSpace($userInput) -and $DefaultValue) {
            $userInput = $DefaultValue
        }
        
        if ($ValidValues.Count -gt 0 -and $userInput -notin $ValidValues) {
            Write-Host "Invalid input. Valid values are: $($ValidValues -join ', ')" -ForegroundColor Red
            continue
        }
        
        return $userInput
        
    } while ($true)
}

function Show-ErrorBlock {
    <#
    .SYNOPSIS
    Display detailed error information with solutions
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Issue,
        
        [string[]]$Solutions = @(),
        
        [string]$ExampleUsage = "",
        
        [string]$AdditionalInfo = ""
    )
    
    Write-Host ""
    Write-Host "ISSUE:" -ForegroundColor Red
    Write-Host "  $Issue" -ForegroundColor White
    Write-Host ""
    
    if ($Solutions.Count -gt 0) {
        Write-Host "SOLUTIONS:" -ForegroundColor Yellow
        foreach ($solution in $Solutions) {
            Write-Host "  - $solution" -ForegroundColor White
        }
        Write-Host ""
    }
    
    if ($ExampleUsage) {
        Write-Host "EXAMPLE:" -ForegroundColor Green
        Write-Host "  $ExampleUsage" -ForegroundColor White
        Write-Host ""
    }
    
    if ($AdditionalInfo) {
        Write-Host "ADDITIONAL INFO:" -ForegroundColor Cyan
        Write-Host "  $AdditionalInfo" -ForegroundColor Gray
        Write-Host ""
    }
}

function Show-InstallationResults {
    <#
    .SYNOPSIS
    Display installation results summary
    #>
    param(
        [Parameter(Mandatory)]
        [hashtable]$Results
    )
    
    if ($Results.OverallSuccess) {
        Write-Success "Successfully installed $($Results.Successful) files"
        if ($Results.Skipped -gt 0) {
            Write-WarningMessage "Skipped $($Results.Skipped) existing files (use -Auto-Approve to overwrite)"
        }
    } else {
        Write-WarningMessage "Installation completed with some failures:"
        Write-Host "  Successful: $($Results.Successful)" -ForegroundColor Green
        Write-Host "  Failed    : $($Results.Failed)" -ForegroundColor Red
        Write-Host "  Skipped   : $($Results.Skipped)" -ForegroundColor Yellow
    }
}

function Show-SourceBranchWelcome {
    <#
    .SYNOPSIS
    Display streamlined welcome message for source branch users
    #>
    param(
        [Parameter(Mandatory)]
        [string]$BranchName
    )
    
    Write-Host "WELCOME TO AI-POWERED TERRAFORM DEVELOPMENT!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Use the contextual help system above to get started." -ForegroundColor Cyan
    Write-Host ""
}

function Get-BootstrapConfirmation {
    <#
    .SYNOPSIS
    Get user confirmation for bootstrap operation
    #>
    
    $response = Read-Host "Would you like to run bootstrap now? [Y/n]"
    if ([string]::IsNullOrWhiteSpace($response) -or $response -match "^[Yy]") {
        Write-Host "Running bootstrap automatically..." -ForegroundColor Green
        Write-Host ""
        return $true
    }
    
    return $false
}

function Show-ContextualError {
    <#
    .SYNOPSIS
    Display contextual error message with appropriate solutions based on branch type
    #>
    param(
        [Parameter(Mandatory)]
        [string]$ErrorMessage,
        
        [string]$BranchType = "unknown",
        
        [string[]]$AdditionalSolutions = @()
    )
    
    Write-ErrorMessage $ErrorMessage
    Write-Host ""
    
    # Provide contextual solutions based on branch type
    $solutions = switch ($BranchType) {
        "source" {
            @(
                "Use -Bootstrap to set up the installer in your user profile",
                "Use -Verify to check the current workspace status",
                "Use -Help to see source branch specific options"
            )
        }
        "feature" {
            @(
                "Use -Verify to check the current workspace status", 
                "Run installer from user profile: ~\.terraform-ai-installer\install-copilot-setup.ps1",
                "Use -Help to see feature branch specific options"
            )
        }
        default {
            @(
                "Use -Verify to check the current workspace status",
                "Use -Help to see all available options",
                "Ensure you're in the correct repository directory"
            )
        }
    }
    
    # Combine with any additional solutions
    $allSolutions = $solutions + $AdditionalSolutions
    
    if ($allSolutions.Count -gt 0) {
        Write-Host "SUGGESTED ACTIONS:" -ForegroundColor Yellow
        foreach ($solution in $allSolutions) {
            Write-Host "  - $solution" -ForegroundColor White
        }
        Write-Host ""
    }
}

function Show-UnknownBranchError {
    <#
    .SYNOPSIS
    Display error for unknown branch scenarios
    #>
    param(
        [bool]$HasRepoDirectory,
        [string]$RepoDirectory = "",
        [string]$ScriptPath = ""
    )
    
    if ($HasRepoDirectory) {
        Show-ErrorBlock -Issue "Cannot determine git branch from the specified repository directory: $RepoDirectory" -Solutions @(
            "Verify the repository directory is correct",
            "Ensure git is available in PATH",
            "Check that the directory is a valid git repository"
        ) -ExampleUsage "`"$ScriptPath`" -RepoDirectory `"C:\correct\path\to\terraform-provider-azurerm`""
    } else {
        Show-ErrorBlock -Issue "Cannot determine git branch from current location" -Solutions @(
            "Run from within the terraform-provider-azurerm repository",
            "Use -RepoDirectory to specify the repository path",
            "Ensure git is available in PATH"
        ) -ExampleUsage "`"$ScriptPath`" -RepoDirectory `"C:\path\to\terraform-provider-azurerm`""
    }
}

# Superior functions that should be used instead of simple ones

function Show-CompletionSummary {
    <#
    .SYNOPSIS
    Display installation completion summary with next steps
    #>
    param(
        [int]$FilesInstalled = 0,
        [int]$FilesSkipped = 0,
        [int]$FilesFailed = 0,
        [string[]]$NextSteps = @(),
        [string]$BranchName = "",
        [string]$BranchType = "feature"
    )
    
    Write-Host ""
    Write-Host "INSTALLATION COMPLETE" -ForegroundColor Green
    Write-Separator -Length 40 -Color Green
    Write-Host ""
    
    # Show branch information if provided
    if ($BranchName) {
        Show-BranchDetection -BranchName $BranchName -BranchType $BranchType
        Write-Host ""
    }
    
    Write-Host "SUMMARY:" -ForegroundColor Cyan
    Write-Host "  Files installed: $FilesInstalled" -ForegroundColor Green
    if ($FilesSkipped -gt 0) {
        Write-Host "  Files skipped:   $FilesSkipped" -ForegroundColor Yellow
    }
    if ($FilesFailed -gt 0) {
        Write-Host "  Files failed:    $FilesFailed" -ForegroundColor Red
    }
    Write-Host ""
    
    if ($NextSteps.Count -gt 0) {
        Write-Host "NEXT STEPS:" -ForegroundColor Cyan
        foreach ($step in $NextSteps) {
            Write-Host "  - $step" -ForegroundColor White
        }
        Write-Host ""
    }
    
    Write-Success "GitHub Copilot AI setup is now active for this workspace!"
}

function Show-Summary {
    <#
    .SYNOPSIS
    Display operation summary with details, warnings, and errors
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Title,
        
        [hashtable]$Details = @{},
        
        [string[]]$Warnings = @(),
        
        [string[]]$Errors = @()
    )
    
    Write-Host ""
    Write-Host "SUMMARY: $Title" -ForegroundColor Cyan
    Write-Separator
    
    # Calculate perfect alignment like our main operation display
    if ($Details.Keys.Count -gt 0) {
        # Find the longest key name for perfect alignment
        $maxKeyLength = ($Details.Keys | Measure-Object -Property Length -Maximum).Maximum
        
        # Display each detail with perfect alignment and beautiful colors
        foreach ($key in $Details.Keys) {
            $paddedKey = $key.PadRight($maxKeyLength)
            Write-Host "  $paddedKey`: " -ForegroundColor Cyan -NoNewline
            
            # Special formatting for operation type - always highlight in yellow
            if ($key -eq "Operation type") {
                $value = $Details[$key]
                if ($value -match '^(.+?)\s+(\(.+\))$') {
                    # Format: "Operation mode (details)" - mode in yellow, details in cyan
                    Write-Host "$($matches[1]) " -ForegroundColor Yellow -NoNewline
                    Write-Host "$($matches[2])" -ForegroundColor Cyan
                } else {
                    # Single operation mode - all in yellow
                    Write-Host "$value" -ForegroundColor Yellow
                }
            } else {
                Write-Host "$($Details[$key])" -ForegroundColor Green
            }
        }
    }
    
    if ($Warnings.Count -gt 0) {
        Write-Host ""
        Write-Host "WARNINGS:" -ForegroundColor Yellow
        foreach ($warning in $Warnings) {
            Write-Host "  ! $warning" -ForegroundColor Yellow
        }
    }
    
    if ($Errors.Count -gt 0) {
        Write-Host ""
        Write-Host "ERRORS:" -ForegroundColor Red
        foreach ($errorMsg in $Errors) {
            Write-Host "  X $errorMsg" -ForegroundColor Red
        }
    }
    
    Write-Host ""
}

function Show-ValidationResults {
    <#
    .SYNOPSIS
    Display validation results with detailed breakdown
    #>
    param(
        [Parameter(Mandatory)]
        [hashtable]$Results,
        
        [bool]$ShowDetails = $false,
        [string]$BranchName = "",
        [string]$BranchType = "feature"
    )
    
    Write-Section "Validation Results"
    
    # Show branch context if available
    if ($BranchName) {
        Show-BranchDetection -BranchName $BranchName -BranchType $BranchType
        Write-Host ""
    }
    
    # Show overall status first
    $overallStatus = if ($Results.OverallValid) { "PASSED" } else { "FAILED" }
    $overallColor = if ($Results.OverallValid) { "Green" } else { "Red" }
    Write-Host "Overall Status: $overallStatus" -ForegroundColor $overallColor
    Write-Host ""
    
    # Show detailed results if requested
    if ($ShowDetails -and $Results.Details) {
        foreach ($category in $Results.Details.Keys) {
            $result = $Results.Details[$category]
            $status = if ($result.Valid) { "PASS" } else { "FAIL" }
            $color = if ($result.Valid) { "Green" } else { "Red" }
            Write-Host "$category`: $status" -ForegroundColor $color
            
            if (-not $result.Valid -and $result.Issues) {
                foreach ($issue in $result.Issues) {
                    Write-Host "  - $issue" -ForegroundColor Red
                }
            }
        }
    }
}

function Show-DirectoryOperation {
    <#
    .SYNOPSIS
    Display directory operation status with color coding
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Directory,
        
        [ValidateSet("Created", "Existing", "Failed")]
        [string]$Status = "Created"
    )
    
    switch ($Status) {
        "Created" {
            Write-Host "  Created directory: $Directory" -ForegroundColor "Green"
        }
        "Existing" {
            Write-Host "  Using existing directory: $Directory" -ForegroundColor "Yellow"
        }
        "Failed" {
            Write-Host "  Failed to create directory: $Directory" -ForegroundColor "Red"
        }
    }
}

function Show-FileOperation {
    <#
    .SYNOPSIS
    Display file operation status
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Operation,
        
        [Parameter(Mandatory)]
        [string]$FileName,
        
        [switch]$NoNewLine
    )
    
    Write-Host "   $Operation" -ForegroundColor Cyan -NoNewline
    Write-Host ": " -ForegroundColor Cyan -NoNewline
    Write-Host "$FileName" -ForegroundColor DarkCyan -NoNewline
    
    if (-not $NoNewLine) {
        Write-Host ""
    }
}

function Wait-ForUser {
    <#
    .SYNOPSIS
    Wait for user to press a key to continue
    #>
    param(
        [string]$Message = "Press any key to continue..."
    )
    
    Write-Host $Message -ForegroundColor Yellow -NoNewline
    $null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
    Write-Host ""
}

function Confirm-UserAction {
    <#
    .SYNOPSIS
    Get user confirmation for an action
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Message,
        
        [string]$DefaultChoice = "N"
    )
    
    $choices = if ($DefaultChoice -eq "Y") { "[Y/n]" } else { "[y/N]" }
    $response = Read-Host "$Message $choices"
    
    if ([string]::IsNullOrWhiteSpace($response)) {
        return $DefaultChoice -eq "Y"
    }
    
    return $response -match "^[Yy]"
}

function Show-BootstrapLocationError {
    <#
    .SYNOPSIS
    Shows a detailed error message when bootstrap is run from wrong location
    
    .DESCRIPTION
    Displays a formatted error message with color-coded paths showing current location
    vs expected location for bootstrap operation
    
    .PARAMETER CurrentLocation
    The current directory where bootstrap is being run from
    
    .PARAMETER ExpectedLocation
    The expected location where bootstrap should be run from
    #>
    param(
        [Parameter(Mandatory)]
        [string]$CurrentLocation,
        
        [Parameter(Mandatory)]
        [string]$ExpectedLocation
    )
    
    Write-Separator
    Write-Host ""
    Write-ErrorMessage "Bootstrap must be run from the source repository, not from user profile directory."
    Write-Host ""
    Write-Host "CORRECT USAGE:" -ForegroundColor Cyan
    Write-Host "  cd C:\path\to\terraform-provider-azurerm" -ForegroundColor Gray
    Write-Host "  & .\.github\AIinstaller\install-copilot-setup.ps1 -Bootstrap" -ForegroundColor Gray
    Write-Host ""
    Write-Host "CURRENT LOCATION: " -ForegroundColor Cyan -NoNewline
    Write-Host "$CurrentLocation" -ForegroundColor Yellow
    Write-Host "EXPECTED LOCATION: " -ForegroundColor Cyan -NoNewline
    Write-Host "$ExpectedLocation" -ForegroundColor Green
    Write-Host ""
}

function Show-VerificationResults {
    <#
    .SYNOPSIS
    Displays workspace verification results in a consistent format
    
    .DESCRIPTION
    Takes verification data and displays it using the standard UI format.
    This ensures UI consistency across all operations.
    
    .PARAMETER VerificationData
    Hashtable containing verification results with Success, Issues, and Details
    #>
    param(
        [Parameter(Mandatory = $true)]
        [hashtable]$VerificationData
    )
    
    Write-Section "Workspace Verification Results"
    
    if ($VerificationData.Success) {
        Write-OperationStatus -Message "Workspace verification completed successfully" -Type "Success"
        
        # Show summary details if available
        if ($VerificationData.Details) {
            $details = @{}
            if ($VerificationData.Details.ContainsKey("WorkspaceType")) {
                $details["Workspace Type"] = $VerificationData.Details.WorkspaceType
            }
            if ($VerificationData.Details.ContainsKey("FilesChecked")) {
                $details["Files Checked"] = $VerificationData.Details.FilesChecked
            }
            if ($VerificationData.Details.ContainsKey("DirectoriesChecked")) {
                $details["Directories Checked"] = $VerificationData.Details.DirectoriesChecked
            }
            
            if ($details.Count -gt 0) {
                Write-Host ""
                Show-Summary -Title "Verification Summary" -Details $details
            }
        }
    } else {
        Write-OperationStatus -Message "Workspace verification encountered issues" -Type "Error"
        
        if ($VerificationData.Issues -and $VerificationData.Issues.Count -gt 0) {
            Write-Host ""
            Write-Host "Issues found:" -ForegroundColor Yellow
            foreach ($issue in $VerificationData.Issues) {
                Write-Host "  - $issue" -ForegroundColor Red
            }
        }
    }
    
    Write-Host ""
}

#endregion

#region Export Module Members

Export-ModuleMember -Function @(
    'Write-Header',
    'Write-Separator',
    'Write-Success',
    'Write-WarningMessage', 
    'Write-ErrorMessage',
    'Write-Section',
    'Write-VerboseMessage',
    'Write-OperationStatus',
    'Show-Help',
    'Show-SourceBranchHelp',
    'Show-FeatureBranchHelp',
    'Show-UnknownBranchHelp',
    'Get-UserInput',
    'Show-ErrorBlock',
    'Show-ContextualError',
    'Show-BranchDetection',
    'Show-InstallationResults',
    'Show-SourceBranchWelcome',
    'Get-BootstrapConfirmation',
    'Show-UnknownBranchError',
    'Show-CompletionSummary',
    'Show-Summary',
    'Show-ValidationResults',
    'Show-VerificationResults',
    'Show-DirectoryOperation',
    'Show-FileOperation',
    'Wait-ForUser',
    'Confirm-UserAction',
    'Format-AlignedLabel',
    'Show-BootstrapLocationError'
)

#endregion
