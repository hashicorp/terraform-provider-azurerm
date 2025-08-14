# UI Module for Terraform AzureRM Provider AI Setup
# Handles all user interface, output formatting, and user interaction

#region Public Functions

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
    Write-Host $("=" * 60) -ForegroundColor Cyan
    Write-Host " $Title" -ForegroundColor Cyan
    Write-Host " Version: $Version" -ForegroundColor Cyan
    Write-Host $("=" * 60) -ForegroundColor Cyan
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

function Write-Warning {
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

function Write-Error {
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

function Write-Info {
    <#
    .SYNOPSIS
    Display informational message in cyan
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Message,
        
        [string]$Prefix = "[INFO]"
    )
    
    Write-Host "$Prefix $Message" -ForegroundColor Cyan
}

function Write-Progress {
    <#
    .SYNOPSIS
    Display progress information
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Activity,
        
        [string]$Status = "",
        
        [int]$PercentComplete = -1
    )
    
    if ($PercentComplete -ge 0) {
        Write-Progress -Activity $Activity -Status $Status -PercentComplete $PercentComplete
    } else {
        Write-Host "[PROGRESS] $Activity" -ForegroundColor Blue
        if ($Status) {
            Write-Host "   $Status" -ForegroundColor Gray
        }
    }
}

function Write-FileOperation {
    <#
    .SYNOPSIS
    Display file operation status (copied, downloaded, etc.)
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Operation,
        
        [Parameter(Mandatory)]
        [string]$FileName,
        
        [Parameter(Mandatory)]
        [ValidateSet("Success", "Failed", "Skipped")]
        [string]$Status
    )
    
    $symbol = switch ($Status) {
        "Success" { "  +" }
        "Failed"  { "  -" }
        "Skipped" { "  ~" }
    }
    
    $color = switch ($Status) {
        "Success" { "Green" }
        "Failed"  { "Red" }
        "Skipped" { "Yellow" }
    }
    
    Write-Host "$symbol $FileName" -ForegroundColor $color
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

function Write-Verbose {
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

function Show-Help {
    <#
    .SYNOPSIS
    Display help menu and usage instructions
    #>
    param(
        [string]$ScriptName = "install-copilot-setup.ps1"
    )
    
    # Get current branch for context-aware help
    $currentBranch = try { git branch --show-current 2>$null } catch { "unknown" }
    $isSourceBranch = $currentBranch -eq "exp/terraform_copilot"
    
    Clear-Host
    Show-Banner -Title "Terraform AzureRM Provider - AI Infrastructure Installer" -Version "1.0.0"
    
    Write-Host "DESCRIPTION:" -ForegroundColor Cyan
    Write-Host "  Installs GitHub Copilot AI instructions and VS Code settings for the"
    Write-Host "  Terraform AzureRM Provider development environment."
    Write-Host ""
    
    # Show current context
    Write-Host "CURRENT SITUATION:" -ForegroundColor Cyan
    Write-Host "  Branch: $currentBranch" -ForegroundColor Gray
    if ($isSourceBranch) {
        Write-Host "  Status: SOURCE BRANCH - Only bootstrap operations allowed" -ForegroundColor Yellow
        Write-Host "  Next:   Run with -Bootstrap to prepare installer for feature branches" -ForegroundColor Yellow
    } else {
        Write-Host "  Status: FEATURE BRANCH - Ready for AI infrastructure installation" -ForegroundColor Green
        Write-Host "  Next:   Run installer to set up AI infrastructure" -ForegroundColor Green
    }
    Write-Host ""
    
    if ($isSourceBranch) {
        Write-Host "WORKFLOW (You are on source branch):" -ForegroundColor Cyan
        Write-Host "  Step 1: .\$ScriptName -Bootstrap" -ForegroundColor White
        Write-Host "          (Copies installer to your user profile)" -ForegroundColor Gray
        Write-Host "  Step 2: git switch -c your-feature-branch" -ForegroundColor White
        Write-Host "          (Switch to your working branch)" -ForegroundColor Gray
        Write-Host "  Step 3: %USERPROFILE%\.terraform-ai-installer\install-copilot-setup.ps1" -ForegroundColor White
        Write-Host "          (Run installer from user profile)" -ForegroundColor Gray
    } else {
        Write-Host "WORKFLOW (You are on feature branch):" -ForegroundColor Cyan
        Write-Host "  If installer is not available locally:" -ForegroundColor Yellow
        Write-Host "    1. git switch exp/terraform_copilot" -ForegroundColor White
        Write-Host "    2. .\$ScriptName -Bootstrap" -ForegroundColor White
        Write-Host "    3. git switch $currentBranch" -ForegroundColor White
        Write-Host "    4. %USERPROFILE%\.terraform-ai-installer\install-copilot-setup.ps1" -ForegroundColor White
        Write-Host ""
        Write-Host "  If installer is available locally:" -ForegroundColor Green
        Write-Host "    .\$ScriptName                    # Install AI infrastructure" -ForegroundColor White
    }
    Write-Host ""
    
    Write-Host "OPTIONS:" -ForegroundColor Cyan
    Write-Host "  -Bootstrap         Copy installer to user profile (source branch only)" -ForegroundColor $(if ($isSourceBranch) { "Green" } else { "Gray" })
    Write-Host "  -Help              Show this help message" -ForegroundColor Green
    Write-Host "  -Verify            Check current AI infrastructure status" -ForegroundColor $(if ($isSourceBranch) { "Green" } else { "White" })
    Write-Host "  -Auto-Approve      Install without prompting (feature branch only)" -ForegroundColor $(if ($isSourceBranch) { "Gray" } else { "White" })
    Write-Host "  -Dry-Run           Preview installation (feature branch only)" -ForegroundColor $(if ($isSourceBranch) { "Gray" } else { "White" })
    Write-Host "  -Clean             Remove AI infrastructure (feature branch only)" -ForegroundColor $(if ($isSourceBranch) { "Gray" } else { "White" })
    Write-Host ""
    
    Write-Host "INSTALLATION BEHAVIOR:" -ForegroundColor Cyan
    Write-Host "  - Default behavior installs/updates all current AI infrastructure files"
    Write-Host "  - Automatically removes deprecated files (no longer in manifest)"
    Write-Host "  - Creates backup of modified settings.json files"
    Write-Host "  - Preserves existing VS Code user settings"
    Write-Host ""
    
    Write-Host "FILES INSTALLED:" -ForegroundColor Cyan
    Write-Host "  .github/copilot-instructions.md      # Main Copilot instructions"
    Write-Host "  .github/instructions/*.md            # Detailed instruction files"
    Write-Host "  .github/prompts/*.md                 # Copilot prompt templates"
    Write-Host "  .vscode/settings.json                # VS Code configuration"
    Write-Host ""
    
    Write-Host "SAFETY NOTES:" -ForegroundColor Yellow
    Write-Host "  - Installation operations are blocked on source branch (exp/terraform_copilot)"
    Write-Host "  - Bootstrap operation safely copies installer without modifying source branch"
    Write-Host "  - Created files are automatically excluded from git commits"
    Write-Host "  - Requires internet connection to download files during bootstrap"
    Write-Host ""
}

function Confirm-UserAction {
    <#
    .SYNOPSIS
    Prompt user for confirmation
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Message,
        
        [string]$Title = "Confirmation Required",
        
        [bool]$DefaultYes = $false
    )
    
    Write-Host ""
    Write-Host "[CONFIRMATION] $Title" -ForegroundColor Yellow
    Write-Host "   $Message" -ForegroundColor White
    Write-Host ""
    
    $prompt = if ($DefaultYes) { "[Y/n]" } else { "[y/N]" }
    $response = Read-Host "   Continue? $prompt"
    
    if ([string]::IsNullOrWhiteSpace($response)) {
        return $DefaultYes
    }
    
    return $response -match "^[Yy]"
}

function Show-CompletionSummary {
    <#
    .SYNOPSIS
    Display installation completion summary
    #>
    param(
        [int]$FilesInstalled = 0,
        [int]$FilesSkipped = 0,
        [int]$FilesFailed = 0,
        [string[]]$NextSteps = @()
    )
    
    Write-Host ""
    Write-Host "INSTALLATION COMPLETE" -ForegroundColor Green
    Write-Host $("=" * 40) -ForegroundColor Green
    Write-Host ""
    
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

function Write-Separator {
    <#
    .SYNOPSIS
    Display a visual separator line
    #>
    param(
        [int]$Length = 60,
        [string]$Character = "-",
        [string]$Color = "Gray"
    )
    
    Write-Host ($Character * $Length) -ForegroundColor $Color
}

function Write-StatusMessage {
    <#
    .SYNOPSIS
    Write a status message with formatting
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Message,
        
        [ValidateSet('Info', 'Warning', 'Error', 'Success')]
        [string]$Type = 'Info',
        
        [bool]$NewLine = $true
    )
    
    switch ($Type) {
        'Info'    { Write-Host "[INFO] $Message" -ForegroundColor Cyan -NoNewline:(-not $NewLine) }
        'Warning' { Write-Host "[WARNING] $Message" -ForegroundColor Yellow -NoNewline:(-not $NewLine) }
        'Error'   { Write-Host "[ERROR] $Message" -ForegroundColor Red -NoNewline:(-not $NewLine) }
        'Success' { Write-Host "[SUCCESS] $Message" -ForegroundColor Green -NoNewline:(-not $NewLine) }
    }
}

function Show-Banner {
    <#
    .SYNOPSIS
    Display formatted application banner
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Title,
        
        [string]$Version = "",
        
        [int]$Width = 80
    )
    
    Write-Host ("=" * $Width) -ForegroundColor Cyan
    $titleLine = " $Title"
    if ($Version) {
        $titleLine += " - Version: $Version"
    }
    Write-Host $titleLine -ForegroundColor White
    Write-Host ("=" * $Width) -ForegroundColor Cyan
    Write-Host ""
}

function Show-Summary {
    <#
    .SYNOPSIS
    Display operation summary
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
    Write-Host ("-" * 50) -ForegroundColor Gray
    
    foreach ($key in $Details.Keys) {
        Write-Host "  $key`: $($Details[$key])" -ForegroundColor White
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

function Wait-ForUser {
    <#
    .SYNOPSIS
    Wait for user to press a key
    #>
    param(
        [string]$Message = "Press any key to continue..."
    )
    
    Write-Host $Message -ForegroundColor Yellow -NoNewline
    $null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
    Write-Host ""
}

function Show-ValidationResults {
    <#
    .SYNOPSIS
    Display validation results with appropriate formatting
    #>
    param(
        [Parameter(Mandatory)]
        [hashtable]$Results,
        
        [bool]$ShowDetails = $false
    )
    
    Write-Section "Validation Results"
    
    # Show overall status first
    $overallStatus = if ($Results.OverallValid) { "PASSED" } else { "FAILED" }
    $overallColor = if ($Results.OverallValid) { "Green" } else { "Red" }
    Write-Host "Overall Status: $overallStatus" -ForegroundColor $overallColor
    Write-Host ""
    
    # Check Git validation first (including branch safety)
    if ($Results.Git) {
        $gitStatus = if ($Results.Git.Valid) { "PASS" } else { "FAIL" }
        $gitColor = if ($Results.Git.Valid) { "Green" } else { "Red" }
        Write-Host "  [$gitStatus] Git Repository" -ForegroundColor $gitColor
        
        # Show branch safety error prominently if it exists
        if (-not $Results.Git.Valid -and $Results.Git.Reason -like "*SAFETY VIOLATION*") {
            Write-Host ""
            Write-Host "CRITICAL SAFETY ERROR:" -ForegroundColor Red -BackgroundColor Yellow
            Write-Host $Results.Git.Reason -ForegroundColor Red
            Write-Host ""
            Write-Host "TO FIX: Create a new branch for this work:" -ForegroundColor Yellow
            Write-Host "  git switch -c feature/my-branch-name" -ForegroundColor Cyan
            Write-Host "  # Or switch to an existing feature branch:" -ForegroundColor Gray
            Write-Host "  git switch existing-feature-branch" -ForegroundColor Cyan
            Write-Host ""
        }
    }
    
    # Check workspace validation
    if ($Results.Workspace) {
        if ($Results.Workspace.Skipped) {
            Write-Host "  [SKIP] Workspace" -ForegroundColor Yellow
        } else {
            $workspaceStatus = if ($Results.Workspace.Valid) { "PASS" } else { "FAIL" }
            $workspaceColor = if ($Results.Workspace.Valid) { "Green" } else { "Red" }
            Write-Host "  [$workspaceStatus] Workspace" -ForegroundColor $workspaceColor
            
            if ($ShowDetails -and -not $Results.Workspace.Valid -and -not $Results.Workspace.Skipped) {
                Write-Host "    - $($Results.Workspace.Reason)" -ForegroundColor Yellow
                if ($Results.Workspace.CurrentPath -and $Results.Workspace.Path -and 
                    $Results.Workspace.CurrentPath -ne $Results.Workspace.Path) {
                    Write-Host "    - Current directory: $($Results.Workspace.CurrentPath)" -ForegroundColor Gray
                    Write-Host "    - Workspace root: $($Results.Workspace.Path)" -ForegroundColor Gray
                }
            }
        }
    }
    
    # Check system requirements
    if ($Results.SystemRequirements) {
        $sysStatus = if ($Results.SystemRequirements.OverallValid) { "PASS" } else { "FAIL" }
        $sysColor = if ($Results.SystemRequirements.OverallValid) { "Green" } else { "Red" }
        Write-Host "  [$sysStatus] System Requirements" -ForegroundColor $sysColor
        
        if ($ShowDetails -and -not $Results.SystemRequirements.OverallValid) {
            if ($Results.SystemRequirements.PowerShell -and -not $Results.SystemRequirements.PowerShell.Valid) {
                Write-Host "    - PowerShell Version: $($Results.SystemRequirements.PowerShell.Reason)" -ForegroundColor Yellow
            }
            if ($Results.SystemRequirements.ExecutionPolicy -and -not $Results.SystemRequirements.ExecutionPolicy.Valid) {
                Write-Host "    - Execution Policy: $($Results.SystemRequirements.ExecutionPolicy.Reason)" -ForegroundColor Yellow
            }
            if ($Results.SystemRequirements.Commands -and -not $Results.SystemRequirements.Commands.Valid) {
                Write-Host "    - Required Commands: $($Results.SystemRequirements.Commands.Reason)" -ForegroundColor Yellow
            }
            if ($Results.SystemRequirements.Internet -and -not $Results.SystemRequirements.Internet.Connected) {
                Write-Host "    - Internet Connectivity: $($Results.SystemRequirements.Internet.Reason)" -ForegroundColor Yellow
            }
        }
    }
    
    Write-Host ""
}

#endregion

#region Export Module Members

Export-ModuleMember -Function @(
    'Write-Header',
    'Write-Success',
    'Write-Warning', 
    'Write-Error',
    'Write-Info',
    'Write-Progress',
    'Write-FileOperation',
    'Write-Section',
    'Write-Verbose',
    'Show-Help',
    'Confirm-UserAction',
    'Show-CompletionSummary',
    'Write-Separator',
    'Write-StatusMessage',
    'Show-Banner',
    'Show-Summary',
    'Get-UserInput',
    'Wait-ForUser',
    'Show-ValidationResults'
)

#endregion
