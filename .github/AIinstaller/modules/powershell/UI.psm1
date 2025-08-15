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
        Write-Host "    4. %USERPROFILE%\.terraform-ai-installer\install-copilot-setup.ps1 -RepoDirectory `"C:\path\to\repo`"" -ForegroundColor White
        Write-Host ""
        Write-Host "  If installer is available locally:" -ForegroundColor Green
        Write-Host "    .\$ScriptName                    # Install AI infrastructure" -ForegroundColor White
        Write-Host "    .\$ScriptName -RepoDirectory `"C:\path\to\repo`"  # When running from user profile" -ForegroundColor White
    }
    Write-Host ""
    
    Write-Host "OPTIONS:" -ForegroundColor Cyan
    Write-Host "  -Bootstrap         Copy installer to user profile (source branch only)" -ForegroundColor $(if ($isSourceBranch) { "Green" } else { "Gray" })
    Write-Host "  -RepoDirectory     Path to repository for git operations (when running from user profile)" -ForegroundColor $(if ($isSourceBranch) { "Gray" } else { "White" })
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
    Write-Host "  - Use -RepoDirectory when running from user profile to specify repository location"
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

function Show-ErrorBlock {
    <#
    .SYNOPSIS
    Display a structured error block with issue description and solutions
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Issue,
        
        [string[]]$Solutions = @(),
        
        [string]$ExampleUsage,
        
        [string]$AdditionalInfo,
        
        [string[]]$AdditionalCommands = @()
    )
    
    Write-Host ""
    Write-Host "ISSUE:" -ForegroundColor Red
    Write-Host "  $Issue" -ForegroundColor Yellow
    Write-Host ""
    
    if ($Solutions.Count -gt 0) {
        Write-Host "HOW TO FIX:" -ForegroundColor Cyan
        for ($i = 0; $i -lt $Solutions.Count; $i++) {
            Write-Host "  $($i + 1). $($Solutions[$i])" -ForegroundColor White
        }
        Write-Host ""
    }
    
    if ($AdditionalCommands.Count -gt 0) {
        Write-Host "TROUBLESHOOTING COMMANDS:" -ForegroundColor Cyan
        foreach ($command in $AdditionalCommands) {
            Write-Host "  $command" -ForegroundColor White
        }
        Write-Host ""
    }
    
    if ($ExampleUsage) {
        Write-Host "EXAMPLE USAGE:" -ForegroundColor Cyan
        Write-Host "  $ExampleUsage" -ForegroundColor Gray
        Write-Host ""
    }
    
    if ($AdditionalInfo) {
        Write-Host "NOTE:" -ForegroundColor Yellow
        Write-Host "  $AdditionalInfo" -ForegroundColor Yellow
        Write-Host ""
    }
}

function Show-BranchDetection {
    <#
    .SYNOPSIS
    Display branch detection result with appropriate styling
    #>
    param(
        [Parameter(Mandatory)]
        [string]$BranchName,
        
        [ValidateSet("source", "feature", "unknown")]
        [string]$BranchType = "feature"
    )
    
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
}

function Show-FileOperation {
    <#
    .SYNOPSIS
    Display file operation with inline status
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Operation,
        
        [Parameter(Mandatory)]
        [string]$FileName,
        
        [ValidateSet("OK", "FAILED", "SKIPPED")]
        [string]$Status,
        
        [switch]$NoNewLine
    )
    
    Write-Host "    $Operation`: $FileName" -ForegroundColor "Gray" -NoNewline
    
    if (-not $NoNewLine -and $Status) {
        switch ($Status) {
            "OK" { Write-Host " [OK]" -ForegroundColor "Green" }
            "FAILED" { Write-Host " [FAILED]" -ForegroundColor "Red" }
            "SKIPPED" { Write-Host " [SKIPPED]" -ForegroundColor "Yellow" }
        }
    }
}

function Show-DirectoryOperation {
    <#
    .SYNOPSIS
    Display directory operation status
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
            Write-Warning "Skipped $($Results.Skipped) existing files (use -Auto-Approve to overwrite)"
        }
    } else {
        Write-Warning "Installation completed with some failures:"
        Write-Host "  Successful: $($Results.Successful)" -ForegroundColor Green
        Write-Host "  Failed: $($Results.Failed)" -ForegroundColor Red
        Write-Host "  Skipped: $($Results.Skipped)" -ForegroundColor Yellow
    }
}

function Show-RepositoryInfo {
    <#
    .SYNOPSIS
    Display repository directory information
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Directory
    )
    
    Write-Host "Using repository directory: $Directory" -ForegroundColor Green
}

function Write-OperationStatus {
    <#
    .SYNOPSIS
    Write operation status with color coding
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Message,
        
        [ValidateSet("Info", "Success", "Warning", "Error", "Progress")]
        [string]$Type = "Info"
    )
    
    switch ($Type) {
        "Info" { Write-Host $Message -ForegroundColor "White" }
        "Success" { Write-Host $Message -ForegroundColor "Green" }
        "Warning" { Write-Host $Message -ForegroundColor "Yellow" }
        "Error" { Write-Host $Message -ForegroundColor "Red" }
        "Progress" { Write-Host $Message -ForegroundColor "Cyan" }
    }
}

function Show-SourceBranchWelcome {
    <#
    .SYNOPSIS
    Shows the welcome message and options for source branch users
    
    .PARAMETER BranchName
    The name of the source branch
    
    .PARAMETER BootstrapCommand
    The command to run bootstrap
    
    .EXAMPLE
    Show-SourceBranchWelcome -BranchName "main" -BootstrapCommand ".\install-copilot-setup.ps1 -Bootstrap"
    #>
    param(
        [Parameter(Mandatory)]
        [string]$BranchName,
        
        [Parameter(Mandatory)]
        [string]$BootstrapCommand
    )
    
    Write-Host "WELCOME TO AI-POWERED TERRAFORM DEVELOPMENT!" -ForegroundColor "Cyan"
    Write-Host ""
    Write-Host "You're running from the source branch ($BranchName)." -ForegroundColor "White"
    Write-Host "To get started with AI infrastructure on your feature branch:" -ForegroundColor "White"
    Write-Host ""
    
    Write-Host "QUICK START (Recommended):" -ForegroundColor "Green"
    Write-Host "  Run the bootstrap command to set up the installer:" -ForegroundColor "White"
    Write-Host "  $BootstrapCommand" -ForegroundColor "Yellow"
    Write-Host ""
    
    Write-Host "MANUAL WORKFLOW:" -ForegroundColor "Cyan"
    Write-Host "  1. Bootstrap: $BootstrapCommand" -ForegroundColor "Gray"
    Write-Host "  2. Switch branch: git checkout feature/your-branch-name" -ForegroundColor "Gray"
    Write-Host "  3. Install: Run installer from user profile" -ForegroundColor "Gray"
    Write-Host ""
    
    Write-Host "OTHER OPTIONS:" -ForegroundColor "White"
    Write-Host "  -Verify       Check current workspace status" -ForegroundColor "Gray"
    Write-Host "  -Help         Show detailed help information" -ForegroundColor "Gray"
    Write-Host ""
}

function Get-BootstrapConfirmation {
    <#
    .SYNOPSIS
    Prompts user for bootstrap confirmation and returns their response
    
    .OUTPUTS
    Returns $true if user wants to bootstrap, $false otherwise
    
    .EXAMPLE
    $shouldBootstrap = Get-BootstrapConfirmation
    #>
    
    Write-Host "Would you like to run bootstrap now? [Y/n]: " -NoNewline -ForegroundColor "Yellow"
    $response = Read-Host
    
    if ($response -eq "" -or $response -eq "y" -or $response -eq "Y") {
        Write-Host ""
        Write-OperationStatus -Message "Running bootstrap automatically..." -Type "Success"
        return $true
    } else {
        Write-Host ""
        Write-Host "No problem! Run with -Bootstrap when you're ready." -ForegroundColor "Cyan"
        Write-Host "Or use -Help for more information." -ForegroundColor "Cyan"
        return $false
    }
}

function Show-UnknownBranchError {
    <#
    .SYNOPSIS
    Shows appropriate error message for unknown branch scenarios
    
    .PARAMETER HasRepoDirectory
    Whether RepoDirectory parameter was provided
    
    .PARAMETER RepoDirectory
    The repository directory path (if provided)
    
    .PARAMETER ScriptPath
    The script path for examples
    
    .EXAMPLE
    Show-UnknownBranchError -HasRepoDirectory $false -ScriptPath $PSCommandPath
    Show-UnknownBranchError -HasRepoDirectory $true -RepoDirectory "C:\invalid" -ScriptPath $PSCommandPath
    #>
    param(
        [Parameter(Mandatory)]
        [bool]$HasRepoDirectory,
        
        [string]$RepoDirectory,
        
        [Parameter(Mandatory)]
        [string]$ScriptPath
    )
    
    if (-not $HasRepoDirectory) {
        # Case 1: Running from user profile without RepoDirectory
        Write-Host ""
        Write-Error "Repository location not specified"
        
        Show-ErrorBlock -Issue "When running from user profile, you must specify the repository location" -Solutions @(
            "Use the -RepoDirectory parameter to specify the terraform-provider-azurerm repository path"
        ) -ExampleUsage "$ScriptPath -RepoDirectory `"C:\github.com\hashicorp\terraform-provider-azurerm`"" -AdditionalInfo "The -RepoDirectory parameter tells the installer where to find your git repository for branch detection and workspace identification."
    } else {
        # Case 2: RepoDirectory was provided but git operations failed
        Write-Host ""
        Write-Error "Invalid repository directory or git operations failed"
        
        Show-ErrorBlock -Issue "The specified repository directory has issues:" -Solutions @(
            "Path does not exist or is not accessible",
            "Directory is not a git repository", 
            "Git is not installed or not in PATH",
            "Repository is in a corrupted state"
        ) -ExampleUsage "$ScriptPath -RepoDirectory `"C:\github.com\hashicorp\terraform-provider-azurerm`"" -AdditionalInfo "Path: $RepoDirectory" -AdditionalCommands @(
            "Verify the path exists: Test-Path `"$RepoDirectory`"",
            "Check if it's a git repo: git status (from that directory)",
            "Verify git is available: git --version"
        )
    }
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
    'Show-ValidationResults',
    'Show-ErrorBlock',
    'Show-BranchDetection',
    'Show-FileOperation',
    'Show-DirectoryOperation',
    'Write-OperationStatus',
    'Show-InstallationResults',
    'Show-RepositoryInfo',
    'Show-SourceBranchWelcome',
    'Get-BootstrapConfirmation',
    'Show-UnknownBranchError'
)

#endregion
