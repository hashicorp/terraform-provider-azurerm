# UI Module for Terraform AzureRM Provider AI Setup
# STREAMLINED VERSION - Contains only functions actually used by main script and dependencies

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

function Format-AlignedLabel {
    <#
    .SYNOPSIS
    Format a label with dynamic spacing to align with other labels
    .DESCRIPTION
    Returns a formatted string with appropriate spacing to align labels in a list.
    Calculates the required padding based on the longest label provided to ensure
    consistent vertical alignment when displaying multiple label-value pairs.
    
    .PARAMETER Label
    The label text to format (without decorative characters like colons)
    
    .PARAMETER LongestLabel
    The longest label in the set (without decorative characters like colons or separators)
    Used as the baseline for calculating alignment spacing
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Label,
        
        [Parameter(Mandatory)]
        [string]$LongestLabel
    )
    
    # Calculate required spacing for alignment based on the actual longest label
    $requiredWidth = $LongestLabel.Length - $Label.Length
    if ($requiredWidth -lt 1) { $requiredWidth = 1 }
    
    return "$Label$(' ' * $requiredWidth)"
}

function Show-BranchDetection {
    <#
    .SYNOPSIS
    Display current branch detection with type-based formatting
    #>
    param(
        [string]$BranchName = "Unknown",
        
        [ValidateSet("source", "feature", "Unknown")]
        [string]$BranchType = "Unknown"
    )
    
    switch ($BranchType) {
        "source" {
            $branchLabel = "SOURCE BRANCH DETECTED"
            Write-Host "${branchLabel}: " -NoNewline -ForegroundColor Cyan
            Write-Host "$BranchName" -ForegroundColor Yellow
        }
        "feature" {
            $branchLabel = "FEATURE BRANCH DETECTED"
            Write-Host "${branchLabel}: " -NoNewline -ForegroundColor Cyan
            Write-Host "$BranchName" -ForegroundColor Yellow
        }
        default {
            $branchLabel = "BRANCH DETECTED"
            Write-Host "${branchLabel}: " -NoNewline -ForegroundColor Cyan
            Write-Host "$BranchName" -ForegroundColor Yellow
        }
    }
    
    # Dynamic workspace label with proper alignment and colors
    if ($Global:WorkspaceRoot) {
        $formattedWorkspaceLabel = Format-AlignedLabel -Label "WORKSPACE" -LongestLabel $branchLabel
        Write-Host "${formattedWorkspaceLabel}: " -NoNewline -ForegroundColor Cyan
        Write-Host "$Global:WorkspaceRoot" -ForegroundColor Green
    }
    
    Write-Host ""
    Write-Separator
}

function Show-Help {
    <#
    .SYNOPSIS
    Display contextual help information based on branch type
    #>
    param(
        [string]$BranchType = "Unknown",
        [bool]$WorkspaceValid = $true,
        [string]$WorkspaceIssue = ""
    )
    
    Write-Host ""
    Write-Host "DESCRIPTION:" -ForegroundColor Cyan
    Write-Host "  Interactive installer for AI-assisted development infrastructure that enhances"
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
            Show-UnknownBranchHelp -WorkspaceValid $WorkspaceValid -WorkspaceIssue $WorkspaceIssue
        }
    }
    
    Write-Host "For more information, visit: https://github.com/hashicorp/terraform-provider-azurerm" -ForegroundColor Cyan
    Write-Host ""
}

function Show-SourceBranchHelp {
    <#
    .SYNOPSIS
    Display help specific to source branch operations
    #>
    
    Write-Host "USAGE:" -ForegroundColor Cyan
    Write-Host "  .\install-copilot-setup.ps1 [OPTIONS]"
    Write-Host ""
    Write-Host "AVAILABLE OPTIONS:" -ForegroundColor Cyan
    Write-Host "  -Bootstrap        Copy installer to user profile (~\.terraform-ai-installer\)"
    Write-Host "                    Run this from the source branch to set up for feature branch use"
    Write-Host "  -Verify           Check current workspace status and validate setup"
    Write-Host "  -Help             Show this help information"
    Write-Host ""
    Write-Host "EXAMPLES:" -ForegroundColor Cyan
    Write-Host "  Bootstrap installer (run from source branch):"
    Write-Host "    .\install-copilot-setup.ps1 -Bootstrap"
    Write-Host ""
    Write-Host "  Verify setup:"
    Write-Host "    .\install-copilot-setup.ps1 -Verify"
    Write-Host ""
    Write-Host "BOOTSTRAP WORKFLOW:" -ForegroundColor Cyan
    Write-Host "  1. Run -Bootstrap from source branch (exp/terraform_copilot) to copy installer to user profile"
    Write-Host "  2. Switch to your feature branch: git checkout feature/your-branch-name"
    Write-Host "  3. Navigate to user profile: cd ~\.terraform-ai-installer\"
    Write-Host "  4. Run installer: .\install-copilot-setup.ps1 -RepoDirectory `"C:\path\to\your\feature\branch`""
    Write-Host ""
}

function Show-FeatureBranchHelp {
    <#
    .SYNOPSIS
    Display help specific to feature branch operations
    #>
    
    Write-Host "USAGE:" -ForegroundColor Cyan
    Write-Host "  .\install-copilot-setup.ps1 [OPTIONS]"
    Write-Host ""
    
    Write-Host "AVAILABLE OPTIONS:" -ForegroundColor Cyan
    Write-Host "  -RepoDirectory    Repository path (path to your feature branch directory)"
    Write-Host "  -Auto-Approve     Overwrite existing files without prompting"
    Write-Host "  -Dry-Run          Show what would be done without making changes"
    Write-Host "  -Verify           Check current workspace status and validate setup"
    Write-Host "  -Clean            Remove AI infrastructure from workspace"
    Write-Host "  -Help             Show this help information"
    Write-Host ""
    
    Write-Host "EXAMPLES:" -ForegroundColor Cyan
    Write-Host "  Install AI infrastructure:"
    Write-Host "    cd ~\.terraform-ai-installer\"
    Write-Host "    .\install-copilot-setup.ps1 -RepoDirectory `"C:\path\to\your\feature\branch`""
    Write-Host ""
    Write-Host "  Dry-Run (preview changes):"
    Write-Host "    cd ~\.terraform-ai-installer\"
    Write-Host "    .\install-copilot-setup.ps1 -RepoDirectory `"C:\path\to\your\feature\branch`" -Dry-Run"
    Write-Host ""
    Write-Host "  Auto-Approve installation:"
    Write-Host "    cd ~\.terraform-ai-installer\"
    Write-Host "    .\install-copilot-setup.ps1 -RepoDirectory `"C:\path\to\your\feature\branch`" -Auto-Approve"
    Write-Host ""
    Write-Host "  Clean removal:"
    Write-Host "    cd ~\.terraform-ai-installer\"
    Write-Host "    .\install-copilot-setup.ps1 -RepoDirectory `"C:\path\to\your\feature\branch`" -Clean"
    Write-Host ""
    
    Write-Host "WORKFLOW:" -ForegroundColor Cyan
    Write-Host "  1. Navigate to user profile installer directory: cd ~\.terraform-ai-installer\"
    Write-Host "  2. Run installer with path to your feature branch"
    Write-Host "  3. Start developing with enhanced GitHub Copilot AI features"
    Write-Host "  4. Use -Clean to remove AI infrastructure when done"
    Write-Host ""
}

function Show-UnknownBranchHelp {
    <#
    .SYNOPSIS
    Display generic help when branch type cannot be determined
    #>
    param(
        [bool]$WorkspaceValid = $true,
        [string]$WorkspaceIssue = ""
    )
    
    # Show workspace issue if detected
    if (-not $WorkspaceValid -and $WorkspaceIssue) {
        Write-Host "WORKSPACE ISSUE DETECTED:" -ForegroundColor Cyan
        Write-Host "  $WorkspaceIssue" -ForegroundColor Yellow
        Write-Host ""
        Write-Host "SOLUTION:" -ForegroundColor Cyan
        Write-Host "  Navigate to a terraform-provider-azurerm repository, or use the -RepoDirectory parameter:"
        Write-Host "  .\install-copilot-setup.ps1 -RepoDirectory `"C:\path\to\terraform-provider-azurerm`" -Help"
        Write-Host ""
        Write-Separator
        Write-Host ""
    }
    
    Write-Host "USAGE:" -ForegroundColor Cyan
    Write-Host "  .\install-copilot-setup.ps1 [OPTIONS]"
    Write-Host ""
    
    Write-Host "ALL OPTIONS:" -ForegroundColor Cyan
    Write-Host "  -Bootstrap        Copy installer to user profile (~\.terraform-ai-installer\)"
    Write-Host "  -RepoDirectory    Repository path for git operations (when running from user profile)"
    Write-Host "  -Auto-Approve     Overwrite existing files without prompting"
    Write-Host "  -Dry-Run          Show what would be done without making changes"
    Write-Host "  -Verify           Check current workspace status and validate setup"
    Write-Host "  -Clean            Remove AI infrastructure from workspace"
    Write-Host "  -Help             Show this help information"
    Write-Host ""
    
    Write-Host "EXAMPLES:" -ForegroundColor Cyan
    Write-Host "  Source Branch Operations:" -ForegroundColor DarkCyan
    Write-Host "    .\install-copilot-setup.ps1 -Bootstrap"
    Write-Host "    .\install-copilot-setup.ps1 -Verify"
    Write-Host ""
    Write-Host "  Feature Branch Operations:" -ForegroundColor DarkCyan
    Write-Host "    cd ~\.terraform-ai-installer\"
    Write-Host "    .\install-copilot-setup.ps1 -RepoDirectory `"C:\path\to\your\feature\branch`""
    Write-Host "    .\install-copilot-setup.ps1 -RepoDirectory `"C:\path\to\your\feature\branch`" -Clean"
    Write-Host ""
    
    Write-Host "BRANCH DETECTION:" -ForegroundColor Cyan
    Write-Host "  The installer automatically detects your branch type and shows appropriate options."
    Write-Host "  If branch detection fails, use the examples above as guidance."
    Write-Host ""
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
        Write-Host "[SUCCESS] Successfully installed $($Results.Successful) files" -ForegroundColor Green
        if ($Results.Skipped -gt 0) {
            Write-Host "[WARNING] Skipped $($Results.Skipped) existing files (use -Auto-Approve to overwrite)" -ForegroundColor Yellow
        }
    } else {
        Write-Host "[WARNING] Installation completed with some failures:" -ForegroundColor Yellow
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
    
    Write-Host "WELCOME TO AI-ASSISTED AZURERM TERRAFORM DEVELOPMENT" -ForegroundColor Green
    Write-Host ""
    Write-Host "Use the contextual help system above to get started." -ForegroundColor Cyan
    Write-Host ""
}

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
        $stepNumber = 1
        foreach ($step in $NextSteps) {
            Write-Host "  $stepNumber. $step" -ForegroundColor White
            $stepNumber++
        }
        Write-Host ""
    }
    
    # Show contextual completion message
    switch ($BranchType) {
        "source" {
            Write-Host "AI infrastructure is now available in the source repository!" -ForegroundColor Green
            Write-Host "Switch to a feature branch to start developing with AI assistance." -ForegroundColor Cyan
        }
        "feature" {
            Write-Host "AI infrastructure is now installed in your workspace!" -ForegroundColor Green
            Write-Host "Start coding and experience enhanced GitHub Copilot with Terraform expertise." -ForegroundColor Cyan
        }
        default {
            Write-Host "AI infrastructure installation completed!" -ForegroundColor Green
        }
    }
    
    Write-Host ""
}

function Show-Summary {
    <#
    .SYNOPSIS
    Display operation summary with details
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Title,
        
        [hashtable]$Details = @{},
        
        [string[]]$Errors = @(),
        
        [string[]]$Warnings = @(),
        
        [string[]]$NextSteps = @()
    )
    
    Write-Host ""
    Write-Host $Title.ToUpper() -ForegroundColor Cyan
    Write-Separator
    Write-Host ""
    
    if ($Details.Count -gt 0) {
        Write-Host "DETAILS:" -ForegroundColor Cyan
        
        # Calculate the maximum key length for proper alignment
        $maxKeyLength = 0
        foreach ($key in $Details.Keys) {
            if ($key.Length -gt $maxKeyLength) {
                $maxKeyLength = $key.Length
            }
        }
        
        # Display each detail with aligned colons and proper colors
        foreach ($key in $Details.Keys) {
            $value = $Details[$key]
            $padding = " " * ($maxKeyLength - $key.Length)
            
            # Write key in cyan
            Write-Host "  $key$padding" -ForegroundColor Cyan -NoNewline
            Write-Host ": " -ForegroundColor Cyan -NoNewline
            
            # Determine value color based on content
            if ($value -match '^\d+$') {
                # Numbers in green
                Write-Host $value -ForegroundColor Green
            } else {
                # Text values in yellow
                Write-Host $value -ForegroundColor Yellow
            }
        }
        Write-Host ""
    }
    
    if ($Warnings.Count -gt 0) {
        Write-Host "WARNINGS:" -ForegroundColor Yellow
        foreach ($warningMsg in $Warnings) {
            Write-Host "  ! $warningMsg" -ForegroundColor Yellow
        }
        Write-Host ""
    }
    
    if ($Errors.Count -gt 0) {
        Write-Host "ERRORS:" -ForegroundColor Red
        foreach ($errorMsg in $Errors) {
            Write-Host "  X $errorMsg" -ForegroundColor Red
        }
    }
    
    Write-Host ""
}

function Show-ParameterError {
    <#
    .SYNOPSIS
    Display a friendly error message for invalid parameters
    
    .DESCRIPTION
    Shows a clean error message when invalid parameters are detected,
    followed by the help information.
    #>
    param(
        [string]$ParameterName
    )
    
    Write-Host ""
    Write-Host "PARAMETER ERROR" -ForegroundColor Red
    Write-Separator -Character "-" -Color Red
    Write-Host ""
    Write-Host "'$ParameterName' is not a recognized parameter." -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Here are the available options:" -ForegroundColor Cyan
    Write-Host ""
    
    # Show the standard help (minimal version)
    Show-Help -BranchName "Unknown" -BranchType "Unknown" -SkipHeader:$true -WorkspaceValid $true
}

function Show-SafetyViolation {
    <#
    .SYNOPSIS
    Display safety violation message for source branch operations
    
    .DESCRIPTION
    Shows a standardized safety violation message when operations are attempted
    on the source branch that should only be performed on feature branches.
    #>
    param(
        [string]$BranchName = "exp/terraform_copilot",
        [string]$Operation = "operation",
        [switch]$FromUserProfile
    )
    
    Write-Host "SAFETY VIOLATION: Cannot perform operations on source branch" -ForegroundColor Red
    Write-Separator
    Write-Host ""
    Write-Host "The -RepoDirectory points to the source branch '$BranchName'." -ForegroundColor Yellow
    Write-Host "Operations other than -Verify, -Help, and -Bootstrap are not allowed on the source branch." -ForegroundColor Yellow
    Write-Host ""
    Write-Host "SOLUTION:" -ForegroundColor Cyan
    Write-Host "  Switch to a feature branch in your target repository:" -ForegroundColor DarkCyan
    
    if ($FromUserProfile) {
        Write-Host "    cd `"<path-to-your-terraform-provider-azurerm>`"" -ForegroundColor Gray
    } else {
        Write-Host "    cd `"$Global:WorkspaceRoot`"" -ForegroundColor Gray
    }
    
    Write-Host "    git checkout -b feature/your-branch-name" -ForegroundColor Gray
    Write-Host ""
    
    if ($FromUserProfile) {
        Write-Host "  Then run the installer from your user profile:" -ForegroundColor DarkCyan
        Write-Host "    cd `"$env:USERPROFILE\.terraform-ai-installer`"" -ForegroundColor Gray
        Write-Host "    .\install-copilot-setup.ps1 -RepoDirectory `"<path-to-your-terraform-provider-azurerm>`"" -ForegroundColor Gray
        Write-Host ""
    }
}

#region Operation Summary Functions

function Write-Status {
    <#
    .SYNOPSIS
    Centralized status message formatting for consistent output
    #>
    param(
        [Parameter(Mandatory = $true)]
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

function Show-OperationSummary {
    <#
    .SYNOPSIS
    Centralized operation summary display for consistent success/failure reporting
    
    .DESCRIPTION
    Provides a standardized way to report operation outcomes across all installer functions.
    This ensures consistent formatting and messaging for long-term maintenance.
    
    .PARAMETER OperationName
    The name of the operation being summarized (e.g., "Installation", "Cleanup", "Bootstrap")
    
    .PARAMETER Success
    Whether the operation was successful
    
    .PARAMETER ItemsProcessed
    Total number of items processed
    
    .PARAMETER ItemsSuccessful
    Number of items processed successfully
    
    .PARAMETER ItemsFailed
    Number of items that failed processing
    
    .PARAMETER DryRun
    Whether this was a dry run operation
    
    .PARAMETER Details
    Additional details to include in the summary
    
    .PARAMETER FailureDetails
    Array of failure messages for specific items
    #>
    param(
        [Parameter(Mandatory = $true)]
        [string]$OperationName,
        
        [Parameter(Mandatory = $true)]
        [bool]$Success,
        
        [int]$ItemsProcessed = 0,
        [int]$ItemsSuccessful = 0,
        [int]$ItemsFailed = 0,
        
        [switch]$DryRun,
        
        [string[]]$Details = @(),
        [string[]]$FailureDetails = @()
    )
    
    Write-Separator
    
    # Operation header with status
    if ($DryRun) {
        $operationTitle = "$OperationName Dry Run"
    } else {
        $operationTitle = $OperationName
    }
    
    if ($Success) {
        Write-Host "$operationTitle completed successfully!" -ForegroundColor Green
    } else {
        Write-Host "$operationTitle failed!" -ForegroundColor Red
    }
    
    Write-Host ""
    
    # Statistics (if provided)
    if ($ItemsProcessed -gt 0) {
        Write-Host "SUMMARY:" -ForegroundColor White
        
        # Collect all labels to calculate proper alignment
        $labels = @("Items processed")
        if ($ItemsSuccessful -gt 0) { $labels += "Successful" }
        if ($ItemsFailed -gt 0) { $labels += "Failed" }
        $labels += "Success Rate"
        
        # Find the longest label for alignment
        $longestLabel = ($labels | Sort-Object Length -Descending)[0]
        
        # Display aligned statistics
        $padding = " " * ($longestLabel.Length - "Items processed".Length)
        Write-Host "  Items processed$padding" -ForegroundColor Cyan -NoNewline
        Write-Host ": " -ForegroundColor Cyan -NoNewline
        Write-Host $ItemsProcessed -ForegroundColor Cyan
        
        if ($ItemsSuccessful -gt 0) {
            $padding = " " * ($longestLabel.Length - "Successful".Length)
            Write-Host "  Successful$padding" -ForegroundColor Cyan -NoNewline
            Write-Host ": " -ForegroundColor Cyan -NoNewline
            Write-Host $ItemsSuccessful -ForegroundColor Green
        }
        
        if ($ItemsFailed -gt 0) {
            $padding = " " * ($longestLabel.Length - "Failed".Length)
            Write-Host "  Failed$padding" -ForegroundColor Cyan -NoNewline
            Write-Host ": " -ForegroundColor Cyan -NoNewline
            Write-Host $ItemsFailed -ForegroundColor Red
        }
        
        # Success rate
        if ($ItemsProcessed -gt 0) {
            $successRate = [math]::Round(($ItemsSuccessful / $ItemsProcessed) * 100, 1)
            $color = if ($successRate -eq 100) { "Green" } elseif ($successRate -ge 80) { "Yellow" } else { "Red" }
            $padding = " " * ($longestLabel.Length - "Success Rate".Length)
            Write-Host "  Success Rate$padding" -ForegroundColor Cyan -NoNewline
            Write-Host ": " -ForegroundColor Cyan -NoNewline
            Write-Host "$successRate%" -ForegroundColor $color
        }
        
        Write-Host ""
    }
    
    # Additional details
    if ($Details.Count -gt 0) {
        Write-Host "DETAILS:" -ForegroundColor White
        foreach ($detail in $Details) {
            Write-Host "  - $detail" -ForegroundColor Gray
        }
        Write-Host ""
    }
    
    # Failure details (if any)
    if ($FailureDetails.Count -gt 0) {
        Write-Host "FAILURES:" -ForegroundColor Red
        foreach ($failure in $FailureDetails) {
            Write-Host "  - $failure" -ForegroundColor Red
        }
        Write-Host ""
    }
    
    # Dry run note
    if ($DryRun) {
        Write-Host "NOTE: This was a dry run - no actual changes were made." -ForegroundColor Yellow
        Write-Host ""
    }
    
    Write-Separator
}

function Show-QuickOperationResult {
    <#
    .SYNOPSIS
    Simplified operation result for quick feedback
    
    .DESCRIPTION
    For simple operations that just need basic success/failure feedback
    #>
    param(
        [Parameter(Mandatory = $true)]
        [string]$OperationName,
        
        [Parameter(Mandatory = $true)]
        [bool]$Success,
        
        [string]$Message = "",
        
        [switch]$DryRun
    )
    
    $prefix = if ($DryRun) { "[DRY RUN] " } else { "" }
    
    if ($Success) {
        $statusMessage = if ($Message) { "$OperationName completed: $Message" } else { "$OperationName completed successfully!" }
        Write-Status -Message "$prefix$statusMessage" -Type Success
    } else {
        $statusMessage = if ($Message) { "$OperationName failed: $Message" } else { "$OperationName failed!" }
        Write-Status -Message "$prefix$statusMessage" -Type Error
    }
}

function Show-OperationSuccess {
    <#
    .SYNOPSIS
    Shows standardized success message for operations
    
    .PARAMETER Operation
    Name of the operation that succeeded
    
    .PARAMETER Summary
    Brief summary of what was accomplished
    
    .PARAMETER NextSteps
    Array of next steps for the user
    #>
    param(
        [string]$Operation,
        [string]$Summary,
        [string[]]$NextSteps = @()
    )
    
    Write-Separator
    Write-Host "$($Operation.ToUpper()) COMPLETED SUCCESSFULLY" -ForegroundColor Green
    Write-Host ""
    Write-Host $Summary -ForegroundColor Green
    
    if ($NextSteps.Count -gt 0) {
        Write-Host ""
        Write-Host "Next Steps:" -ForegroundColor Cyan
        for ($i = 0; $i -lt $NextSteps.Count; $i++) {
            Write-Host "$($i + 1). $($NextSteps[$i])" -ForegroundColor White
        }
    }
    
    Write-Separator
}

function Show-OperationFailure {
    <#
    .SYNOPSIS
    Shows standardized failure message for operations
    
    .PARAMETER Operation
    Name of the operation that failed
    
    .PARAMETER ErrorMessage
    Error message or exception details
    
    .PARAMETER Suggestions
    Array of suggested next steps or troubleshooting tips
    #>
    param(
        [string]$Operation,
        [string]$ErrorMessage,
        [string[]]$Suggestions = @()
    )
    
    Write-Separator
    Write-Host "$($Operation.ToUpper()) FAILED" -ForegroundColor Red
    Write-Host ""
    Write-Host "Error: $ErrorMessage" -ForegroundColor Red
    
    if ($Suggestions.Count -gt 0) {
        Write-Host ""
        Write-Host "Suggestions:" -ForegroundColor Yellow
        for ($i = 0; $i -lt $Suggestions.Count; $i++) {
            Write-Host "$($i + 1). $($Suggestions[$i])" -ForegroundColor White
        }
    }
    
    Write-Separator
}

#endregion

#endregion

# Export only the functions actually used by the main script
Export-ModuleMember -Function @(
    'Write-Separator',
    'Write-Header',
    'Format-AlignedLabel',
    'Show-BranchDetection',
    'Write-OperationStatus',
    'Write-Status',
    'Show-OperationSummary',
    'Show-QuickOperationResult',
    'Show-Help',
    'Show-SourceBranchHelp',
    'Show-FeatureBranchHelp',
    'Show-UnknownBranchHelp',
    'Show-InstallationResults',
    'Show-SourceBranchWelcome',
    'Show-CompletionSummary',
    'Show-Summary',
    'Show-ParameterError',
    'Show-SafetyViolation',
    'Show-OperationSuccess',
    'Show-OperationFailure'
)
