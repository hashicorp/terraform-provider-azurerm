# ValidationEngine Module for Terraform AzureRM Provider AI Setup
# Handles comprehensive validation, dependency checking, and system requirements
# STREAMLINED VERSION - Contains only functions actually used by main script and dependencies

#region Private Functions

function Find-WorkspaceRoot {
    <#
    .SYNOPSIS
    Find the workspace root by looking for go.mod file in current or parent directories
    #>
    param(
        [Parameter(Mandatory)]
        [string]$StartPath
    )
    
    $currentPath = $StartPath
    $maxDepth = 10  # Prevent infinite loops
    $depth = 0
    
    while ($depth -lt $maxDepth -and $currentPath) {
        $goModPath = Join-Path $currentPath "go.mod"
        if (Test-Path $goModPath) {
            return $currentPath
        }
        
        # Move to parent directory
        $parentPath = Split-Path $currentPath -Parent
        if ($parentPath -eq $currentPath) {
            # Reached root directory
            break
        }
        $currentPath = $parentPath
        $depth++
    }
    
    return $null
}

function Test-PowerShellVersion {
    <#
    .SYNOPSIS
    Test if PowerShell version meets requirements
    #>
    
    $minimumVersion = [Version]"5.1"
    $currentVersion = $PSVersionTable.PSVersion
    
    return @{
        Valid = $currentVersion -ge $minimumVersion
        CurrentVersion = $currentVersion.ToString()
        MinimumVersion = $minimumVersion.ToString()
        Reason = if ($currentVersion -ge $minimumVersion) { 
            "PowerShell version is supported" 
        } else { 
            "PowerShell $minimumVersion or later is required" 
        }
    }
}

function Test-ExecutionPolicy {
    <#
    .SYNOPSIS
    Test if execution policy allows script execution
    #>
    
    $policy = Get-ExecutionPolicy -Scope CurrentUser
    $allowedPolicies = @("RemoteSigned", "Unrestricted", "Bypass")
    
    return @{
        Valid = $policy -in $allowedPolicies
        CurrentPolicy = $policy.ToString()
        AllowedPolicies = $allowedPolicies
        Reason = if ($policy -in $allowedPolicies) {
            "Execution policy allows script execution"
        } else {
            "Execution policy '$policy' prevents script execution. Use: Set-ExecutionPolicy RemoteSigned -Scope CurrentUser"
        }
    }
}

function Test-RequiredCommands {
    <#
    .SYNOPSIS
    Test if required external commands are available
    #>
    
    $requiredCommands = @("git")
    $results = @{}
    $allValid = $true
    
    foreach ($command in $requiredCommands) {
        try {
            $commandInfo = Get-Command $command -ErrorAction Stop
            $results[$command] = @{
                Available = $true
                Version = ""
                Path = $commandInfo.Source
                Reason = "Command is available"
            }
            
            # Try to get version for git
            if ($command -eq "git") {
                try {
                    $version = git --version 2>$null
                    if ($version -match "git version (.+)") {
                        $results[$command].Version = $matches[1]
                    }
                }
                catch {
                    # Version detection failed, but command exists
                }
            }
        }
        catch {
            $results[$command] = @{
                Available = $false
                Version = ""
                Path = ""
                Reason = "Command not found in PATH"
            }
            $allValid = $false
        }
    }
    
    return @{
        Valid = $allValid
        Commands = $results
        Reason = if ($allValid) { "All required commands are available" } else { "Some required commands are missing" }
    }
}

function Test-InternetConnectivity {
    <#
    .SYNOPSIS
    Test internet connectivity to required endpoints
    #>
    
    $testUrls = @(
        "https://api.github.com",
        "https://raw.githubusercontent.com"
    )
    
    $results = @{
        Connected = $false
        TestedEndpoints = @{}
        Reason = ""
    }
    
    $successCount = 0
    
    foreach ($url in $testUrls) {
        try {
            # Disable progress bar to prevent console flashing
            $ProgressPreference = 'SilentlyContinue'
            $response = Invoke-WebRequest -Uri $url -Method Head -TimeoutSec 10 -UseBasicParsing
            $results.TestedEndpoints[$url] = @{
                Success = $true
                StatusCode = $response.StatusCode
                ResponseTime = 0
            }
            $successCount++
        }
        catch {
            $results.TestedEndpoints[$url] = @{
                Success = $false
                StatusCode = 0
                Error = $_.Exception.Message
            }
        }
        finally {
            # Restore progress preference
            $ProgressPreference = 'Continue'
        }
    }
    
    $results.Connected = $successCount -gt 0
    $results.Reason = if ($results.Connected) {
        "Internet connectivity verified ($successCount/$($testUrls.Count) endpoints reachable)"
    } else {
        "No internet connectivity detected. Check network connection and firewall settings."
    }
    
    return $results
}

function Test-GitRepository {
    <#
    .SYNOPSIS
    Test if current directory is a valid git repository with branch safety checks
    #>
    param(
        [bool]$AllowBootstrapOnSource = $false,
        [string]$WorkspacePath = ""
    )
    
    $results = @{
        Valid = $false
        IsGitRepo = $false
        HasRemote = $false
        CurrentBranch = ""
        RemoteUrl = ""
        Reason = ""
    }
    
    try {
        # Save current location and switch to workspace if provided
        $originalLocation = Get-Location
        if ($WorkspacePath -and (Test-Path $WorkspacePath)) {
            Set-Location $WorkspacePath
        }
        
        try {
            # Test if we're in a git repository
            $null = git status --porcelain 2>$null
            $results.IsGitRepo = $LASTEXITCODE -eq 0
            
            if ($results.IsGitRepo) {
                # Get current branch
                try {
                    $results.CurrentBranch = git branch --show-current 2>$null
                    if (-not $results.CurrentBranch -or $results.CurrentBranch.Trim() -eq "") {
                        $results.CurrentBranch = "Unknown"
                    }
                }
                catch {
                    $results.CurrentBranch = "Unknown"
                }
                
                # Get remote URL
                try {
                    $results.RemoteUrl = git remote get-url origin 2>$null
                    $results.HasRemote = $LASTEXITCODE -eq 0 -and $results.RemoteUrl
                }
                catch {
                    $results.HasRemote = $false
                }
                
                # CRITICAL SAFETY CHECK: Prevent running on source branch (unless bootstrap)
                $isSourceBranch = $results.CurrentBranch -eq "exp/terraform_copilot"
                
                $results.Valid = $results.IsGitRepo -and $results.HasRemote -and (-not $isSourceBranch -or $AllowBootstrapOnSource)
                
                if ($isSourceBranch -and -not $AllowBootstrapOnSource) {
                    $results.Reason = "SAFETY VIOLATION: Cannot run installer on source branch 'exp/terraform_copilot'. Switch to a different branch to install AI infrastructure."
                }
                elseif ($results.Valid) {
                    if ($isSourceBranch -and $AllowBootstrapOnSource) {
                        $results.Reason = "Source branch - bootstrap operations allowed"
                    } else {
                        $results.Reason = "Valid git repository with remote origin"
                    }
                }
                elseif (-not $results.HasRemote) {
                    $results.Reason = "Git repository has no remote origin configured"
                }
            }
            else {
                $results.Reason = "Not a git repository"
            }
        }
        finally {
            # Always restore original location
            Set-Location $originalLocation
        }
    }
    catch {
        $results.Reason = "Error checking git repository: $($_.Exception.Message)"
    }
    
    return $results
}

#endregion

#region Public Functions

function Test-WorkspaceValid {
    <#
    .SYNOPSIS
    Test if current directory is a valid Terraform AzureRM workspace
    #>
    param(
        [string]$WorkspacePath = ""
    )
    
    # Smart workspace detection - use provided path or find from current location
    if ($WorkspacePath) {
        # When WorkspacePath is provided, check if it's already the workspace root
        $goModInProvidedPath = Join-Path $WorkspacePath "go.mod"
        if (Test-Path $goModInProvidedPath) {
            $workspaceRoot = $WorkspacePath
        } else {
            # If not, search from the provided path
            $workspaceRoot = Find-WorkspaceRoot -StartPath $WorkspacePath
        }
        # Use FullName property for DirectoryInfo objects from Get-Item
        $itemResult = Get-Item $WorkspacePath
        $currentPath = @{ Path = $itemResult.FullName }
    } else {
        $currentPath = Get-Location
        $workspaceRoot = Find-WorkspaceRoot -StartPath $currentPath.Path
    }
    
    $results = @{
        Valid = $false
        Path = $workspaceRoot
        CurrentPath = $currentPath.Path
        IsAzureRMProvider = $false
        HasGoMod = $false
        HasMainGo = $false
        Reason = ""
    }

    if (-not $workspaceRoot) {
        $results.Reason = "Could not locate workspace root (no go.mod found in current path or parent directories)"
        return $results
    }    # Check for go.mod file in workspace root
    $goModPath = Join-Path $workspaceRoot "go.mod"
    $results.HasGoMod = Test-Path $goModPath
    
    # Check for main.go file in workspace root  
    $mainGoPath = Join-Path $workspaceRoot "main.go"
    $results.HasMainGo = Test-Path $mainGoPath
    
    # Check if this is the terraform-provider-azurerm repository
    if ($results.HasGoMod) {
        try {
            $goModContent = Get-Content $goModPath -Raw
            $results.IsAzureRMProvider = $goModContent -match "terraform-provider-azurerm"
        }
        catch {
            $results.IsAzureRMProvider = $false
        }
    }
    
    # Determine validity
    $results.Valid = $results.HasGoMod -and $results.HasMainGo -and $results.IsAzureRMProvider
    
    if ($results.Valid) {
        $results.Reason = "Valid Terraform AzureRM provider workspace"
    }
    elseif (-not $results.HasGoMod) {
        $results.Reason = "Not a Go module (missing go.mod file)"
    }
    elseif (-not $results.HasMainGo) {
        $results.Reason = "Not a Go application (missing main.go file)"
    }
    elseif (-not $results.IsAzureRMProvider) {
        $results.Reason = "Not the Terraform AzureRM provider repository"
    }
    else {
        $results.Reason = "Workspace validation failed"
    }
    
    return $results
}

function Test-SystemRequirements {
    <#
    .SYNOPSIS
    Test all system requirements for the AI installer
    #>
    
    $results = @{
        OverallValid = $true
        PowerShell = Test-PowerShellVersion
        ExecutionPolicy = Test-ExecutionPolicy
        Commands = Test-RequiredCommands
        Internet = Test-InternetConnectivity
    }
    
    # Check if any requirement failed
    $results.OverallValid = $results.PowerShell.Valid -and 
                           $results.ExecutionPolicy.Valid -and 
                           $results.Commands.Valid -and 
                           $results.Internet.Connected
    
    return $results
}

function Test-PreInstallation {
    <#
    .SYNOPSIS
    Run comprehensive pre-installation validation
    #>
    param(
        [bool]$AllowBootstrapOnSource = $false
    )
    
    $results = @{
        OverallValid = $true
        Git = $null
        Workspace = $null
        SystemRequirements = $null
        Timestamp = Get-Date
    }
    
    # CRITICAL: Check Git first for branch safety
    # Use the workspace root for git operations if available
    $gitPath = if ($Global:WorkspaceRoot) { $Global:WorkspaceRoot } else { (Get-Location).Path }
    $results.Git = Test-GitRepository -AllowBootstrapOnSource $AllowBootstrapOnSource -WorkspacePath $gitPath
    
    # If Git validation fails due to branch safety, short-circuit other validations
    # This prevents running unnecessary tests when we know we can't proceed
    if (-not $results.Git.Valid -and $results.Git.Reason -like "*SAFETY VIOLATION*") {
        $results.OverallValid = $false
        
        # Still run system requirements (these are always safe to check)
        $results.SystemRequirements = Test-SystemRequirements
        
        # Skip workspace and detailed checks due to safety violation
        $results.Workspace = @{
            Valid = $false
            Reason = "Skipped due to Git branch safety violation"
            Skipped = $true
        }
        
        return $results
    }
    
    # Continue with full validation if Git is safe
    $results.Workspace = Test-WorkspaceValid -WorkspacePath $Global:WorkspaceRoot
    $results.SystemRequirements = Test-SystemRequirements
    
    # Check overall validity - Git validation (including branch safety) is critical
    $results.OverallValid = $results.Git.Valid -and
                           $results.Workspace.Valid -and
                           $results.SystemRequirements.OverallValid
    
    return $results
}

function Test-SourceRepository {
    <#
    .SYNOPSIS
    Determines if we're running on the source repository vs a target repository
    
    .DESCRIPTION
    Checks various indicators to determine if this is the source repository where
    AI infrastructure files are maintained vs a target repository where they
    would be installed.
    
    .OUTPUTS
    Boolean - True if this is the source repository, False if target
    
    .NOTES
    CRITICAL FUNCTION: This provides essential source repository protection.
    The logic here determines whether files should be copied locally or downloaded
    remotely, preventing accidental overwriting of source files.
    #>
    
    # Check if we're on the exp/terraform_copilot branch (source branch)
    try {
        Push-Location $Global:WorkspaceRoot
        $currentBranch = git rev-parse --abbrev-ref HEAD 2>$null
        if ($currentBranch -eq "exp/terraform_copilot") {
            return $true
        }
    } catch {
        # Git not available or not in a git repo
    } finally {
        Pop-Location
    }
    
    # Check if AIinstaller directory exists (only in source)
    $aiInstallerPath = Join-Path $Global:WorkspaceRoot ".github/AIinstaller"
    if (Test-Path $aiInstallerPath) {
        return $true
    }
    
    # Check if this directory structure looks like the source
    $copilotInstructionsPath = Join-Path $Global:WorkspaceRoot ".github/copilot-instructions.md"
    $instructionsPath = Join-Path $Global:WorkspaceRoot ".github/instructions"
    $promptsPath = Join-Path $Global:WorkspaceRoot ".github/prompts"
    
    if ((Test-Path $copilotInstructionsPath) -and 
        (Test-Path $instructionsPath) -and 
        (Test-Path $promptsPath) -and
        (Test-Path $aiInstallerPath)) {
        return $true
    }
    
    return $false
}

function Invoke-VerifyWorkspace {
    <#
    .SYNOPSIS
    Verifies the presence of AI infrastructure files in the workspace
    
    .DESCRIPTION
    Checks for all required AI infrastructure files including:
    - Main copilot instructions
    - Detailed instruction files
    - Prompts directory
    - VS Code settings
    
    .PARAMETER BranchType
    The type of branch (source, feature, Unknown) for dynamic spacing calculation
    
    .OUTPUTS
    Returns verification results and displays status to console
    
    .NOTES
    This function maintains source repository awareness and provides different
    behavior for source vs target repositories.
    #>
    param(
        [ValidateSet("source", "feature", "Unknown")]
        [string]$BranchType = "feature"
    )
    
    # Use the dynamically determined workspace root
    $workspaceRoot = $Global:WorkspaceRoot
    Push-Location $workspaceRoot
    
    try {
        # CRITICAL: Use centralized validation (replaces Test-SourceRepository)
        $validation = Test-PreInstallation -AllowBootstrapOnSource:$true  # Allow verification on source
        
        $results = @{
            Success = $validation.OverallValid
            Files = @()
            Issues = @()
            IsSourceRepo = ($validation.Git.CurrentBranch -eq "exp/terraform_copilot")
            ValidationResults = $validation
        }
        
        # If basic validation failed, show that first
        if (-not $validation.OverallValid) {
            $results.Issues += "Workspace validation failed: $($validation.Git.Reason)"
            Write-Host "❌ Workspace validation failed!" -ForegroundColor Red
            Write-Host "   $($validation.Git.Reason)" -ForegroundColor Yellow
            return $results
        }
        
        Write-Host " Workspace Verification" -ForegroundColor Cyan
        Write-Separator
        Write-Host ""
        
        # Check main instructions file
        $instructionsFile = $Global:InstallerConfig.Files.Instructions.Target
        if (Test-Path $instructionsFile) {
            $results.Files += @{
                Path = $instructionsFile
                Status = "Present"
                Description = "Main Copilot instructions"
            }
            Write-Host "  [FOUND] $(Resolve-Path $instructionsFile -Relative)" -ForegroundColor Green
        } else {
            $results.Files += @{
                Path = $instructionsFile
                Status = "Missing"
                Description = "Main Copilot instructions"
            }
            $results.Issues += "Missing: $instructionsFile"
            Write-Host "  [MISSING] $(Resolve-Path $instructionsFile -Relative -ErrorAction SilentlyContinue)" -ForegroundColor Red
        }
        
        # Check instructions directory
        $instructionsDir = $Global:InstallerConfig.Files.InstructionFiles.Target
        if (Test-Path $instructionsDir -PathType Container) {
            $results.Files += @{
                Path = $instructionsDir
                Status = "Present"
                Description = "Instructions directory"
            }
            Write-Host "  [FOUND] $(Resolve-Path $instructionsDir -Relative)/" -ForegroundColor Green
            
            # Check specific instruction files
            $requiredFiles = $Global:InstallerConfig.Files.InstructionFiles.Files
            
            foreach ($file in $requiredFiles) {
                # Handle full repository paths vs relative paths
                if ($file.StartsWith('.github/')) {
                    # This is a full repository path - use it directly from workspace root
                    $filePath = Join-Path $Global:WorkspaceRoot $file
                } else {
                    # This is a relative path - join with target directory
                    $filePath = Join-Path $instructionsDir $file
                }
                
                if (Test-Path $filePath) {
                    Write-Host "    [FOUND] $file" -ForegroundColor Green
                } else {
                    Write-Host "    [MISSING] $file" -ForegroundColor Red
                    $results.Issues += "Missing: $filePath"
                }
            }
        } else {
            $results.Files += @{
                Path = $instructionsDir
                Status = "Missing"
                Description = "Instructions directory"
            }
            $results.Issues += "Missing: $instructionsDir"
            Write-Host "  [MISSING] $(Resolve-Path $instructionsDir -Relative -ErrorAction SilentlyContinue)/" -ForegroundColor Red
        }
        
        # Check prompts directory
        $promptsDir = $Global:InstallerConfig.Files.PromptFiles.Target
        if (Test-Path $promptsDir -PathType Container) {
            $results.Files += @{
                Path = $promptsDir
                Status = "Present"
                Description = "Prompts directory"
            }
            Write-Host "  [FOUND] $(Resolve-Path $promptsDir -Relative)/" -ForegroundColor Green
            
            # Check specific prompt files
            $requiredPrompts = $Global:InstallerConfig.Files.PromptFiles.Files
            
            foreach ($file in $requiredPrompts) {
                # Handle full repository paths vs relative paths
                if ($file.StartsWith('.github/')) {
                    # This is a full repository path - use it directly from workspace root
                    $filePath = Join-Path $Global:WorkspaceRoot $file
                } else {
                    # This is a relative path - join with target directory
                    $filePath = Join-Path $promptsDir $file
                }
                
                if (Test-Path $filePath) {
                    Write-Host "    [FOUND] $file" -ForegroundColor Green
                } else {
                    Write-Host "    [MISSING] $file" -ForegroundColor Red
                    $results.Issues += "Missing: $filePath"
                }
            }
        } else {
            $results.Files += @{
                Path = $promptsDir
                Status = "Missing"
                Description = "Prompts directory"
            }
            $results.Issues += "Missing: $promptsDir"
            Write-Host "  [MISSING] $(Resolve-Path $promptsDir -Relative -ErrorAction SilentlyContinue)/" -ForegroundColor Red
        }
        
        # Check .vscode directory and settings
        $vscodeDir = Join-Path $workspaceRoot ".vscode"
        if (Test-Path $vscodeDir -PathType Container) {
            Write-Host "  [FOUND] .vscode/" -ForegroundColor Green
            
            $settingsFile = Join-Path $vscodeDir "settings.json"
            if (Test-Path $settingsFile) {
                Write-Host "    [FOUND] settings.json" -ForegroundColor Green
            } else {
                Write-Host "    [MISSING] settings.json" -ForegroundColor Red
                $results.Issues += "Missing: $settingsFile"
            }
        } else {
            Write-Host "  [MISSING] .vscode/" -ForegroundColor Red
            $results.Issues += "Missing: $vscodeDir"
        }
        
        Write-Host ""
        
        # Show results summary
        if ($results.Issues.Count -gt 0) {
            $results.Success = $false
            Write-Host "Issues found:" -ForegroundColor Yellow
            foreach ($issue in $results.Issues) {
                Write-Host "  - $issue" -ForegroundColor Red
            }
            
            if (-not $results.IsSourceRepo) {
                Write-Host ""
                Write-Host "TIP: To install missing files, run the installer from user profile" -ForegroundColor Cyan
            }
        } else {
            if ($results.IsSourceRepo) {
                Write-Host "All AI infrastructure files are present in the source repository!" -ForegroundColor Green
            } else {
                Write-Host "All AI infrastructure files have been successfully installed!" -ForegroundColor Green
            }
        }
        
        Write-Host ""
        return $results
    }
    finally {
        Pop-Location
    }
}

#endregion

# Export only the functions actually used by the main script and inter-module dependencies
Export-ModuleMember -Function @(
    'Test-WorkspaceValid',
    'Test-PreInstallation',
    'Invoke-VerifyWorkspace',
    'Test-SourceRepository'
)
