# ValidationEngine Module for Terraform AzureRM Provider AI Setup
# Handles comprehensive validation, dependency checking, and system requirements

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

function Test-WorkspaceValid {
    <#
    .SYNOPSIS
    Test if current directory is a valid Terraform AzureRM workspace
    #>
    
    # Smart workspace detection - find the actual workspace root
    $currentPath = Get-Location
    $workspaceRoot = Find-WorkspaceRoot -StartPath $currentPath.Path
    
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
    }
    
    # Check for go.mod file in workspace root
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

function Test-GitRepository {
    <#
    .SYNOPSIS
    Test if current directory is a valid git repository with branch safety checks
    #>
    param(
        [bool]$AllowBootstrapOnSource = $false
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
        # Test if we're in a git repository
        $null = git status --porcelain 2>$null
        $results.IsGitRepo = $LASTEXITCODE -eq 0
        
        if ($results.IsGitRepo) {
            # Get current branch
            try {
                $results.CurrentBranch = git branch --show-current 2>$null
            }
            catch {
                $results.CurrentBranch = "unknown"
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
    catch {
        $results.Reason = "Error checking git repository: $($_.Exception.Message)"
    }
    
    return $results
}

#endregion

#region Public Functions

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
    $results.Git = Test-GitRepository -AllowBootstrapOnSource $AllowBootstrapOnSource
    
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
    $results.Workspace = Test-WorkspaceValid
    $results.SystemRequirements = Test-SystemRequirements
    
    # Check overall validity - Git validation (including branch safety) is critical
    $results.OverallValid = $results.Git.Valid -and
                           $results.Workspace.Valid -and
                           $results.SystemRequirements.OverallValid
    
    return $results
}

function Test-PostInstallation {
    <#
    .SYNOPSIS
    Run comprehensive post-installation validation
    #>
    param(
        [string]$Branch = "exp/terraform_copilot",
        [hashtable]$ManifestConfig = $null
    )
    
    # Use provided manifest configuration or get it directly
    if ($ManifestConfig) {
        $manifestConfig = $ManifestConfig
    } else {
        # Fallback: require ConfigParser to be loaded in parent scope
        if (-not (Get-Command Get-ManifestConfig -ErrorAction SilentlyContinue)) {
            throw "ManifestConfig parameter required or Get-ManifestConfig must be available"
        }
        $manifestConfig = Get-ManifestConfig -Branch $Branch
    }
    $allFiles = @()
    foreach ($section in $manifestConfig.Sections.Keys) {
        $allFiles += $manifestConfig.Sections[$section]
    }
    
    $workspaceStatus = Get-WorkspaceStatus
    
    $results = @{
        OverallValid = $true
        InstallationComplete = $true
        FilesInstalled = $workspaceStatus.InstalledCount
        TotalFiles = $workspaceStatus.TotalFiles
        MissingFiles = $workspaceStatus.MissingFiles
        FileDetails = @{}
        Timestamp = Get-Date
    }
    
    # Validate each installed file
    foreach ($filePath in $allFiles) {
        $fileValidation = @{
            Exists = Test-Path $filePath
            Required = $true  # All files in manifest are considered required
            Integrity = $null
            UpToDate = $null
        }
        
        if ($fileValidation.Exists) {
            # Test file integrity
            $fileValidation.Integrity = Test-FileIntegrity -FilePath $filePath
            
            # Test if up to date
            $fileValidation.UpToDate = Test-FileUpToDate -FilePath $filePath -Branch $Branch
        } else {
            if ($fileInfo.Required) {
                $results.InstallationComplete = $false
                $results.OverallValid = $false
            }
        }
        
        $results.FileDetails[$filePath] = $fileValidation
    }
    
    return $results
}

function Test-AIInfrastructure {
    <#
    .SYNOPSIS
    Comprehensive test of AI infrastructure status
    #>
    param(
        [string]$Branch = "exp/terraform_copilot",
        [hashtable]$ManifestConfig = $null
    )
    
    $results = @{
        OverallHealth = "Unknown"
        InstallationStatus = "Unknown"
        FileHealth = @{}
        DirectoryHealth = @{}
        GitIgnoreStatus = @{}
        Recommendations = @()
        Timestamp = Get-Date
    }
    
    # Get installation status
    $workspaceStatus = Get-WorkspaceStatus
    
    # Use provided manifest configuration or get it directly
    if ($ManifestConfig) {
        $manifestConfig = $ManifestConfig
    } else {
        # Fallback: require ConfigParser to be loaded in parent scope
        if (-not (Get-Command Get-ManifestConfig -ErrorAction SilentlyContinue)) {
            throw "ManifestConfig parameter required or Get-ManifestConfig must be available"
        }
        $manifestConfig = Get-ManifestConfig -Branch $Branch
    }
    $allFiles = @()
    foreach ($section in $manifestConfig.Sections.Keys) {
        $allFiles += $manifestConfig.Sections[$section]
    }
    
    # Determine installation status
    if ($workspaceStatus.InstalledCount -eq 0) {
        $results.InstallationStatus = "Not Installed"
        $results.OverallHealth = "Missing"
        $results.Recommendations += "Run installation to set up AI infrastructure"
    }
    elseif ($workspaceStatus.InstalledCount -eq $workspaceStatus.TotalFiles) {
        $results.InstallationStatus = "Complete"
        $results.OverallHealth = "Healthy"
    }
    else {
        $results.InstallationStatus = "Partial"
        $results.OverallHealth = "Degraded"
        $results.Recommendations += "Some files are missing. Consider reinstalling."
    }
    
    # Check file health
    foreach ($filePath in $allFiles) {
        $fileHealth = @{
            Status = "Unknown"
            Issues = @()
            Size = 0
            LastModified = $null
        }
        
        if (Test-Path $filePath) {
            $fileData = Get-Item $filePath
            $fileHealth.Size = $fileData.Length
            $fileHealth.LastModified = $fileData.LastWriteTime
            
            # Check integrity
            $integrity = Test-FileIntegrity -FilePath $filePath
            if ($integrity.Valid) {
                $fileHealth.Status = "Healthy"
            } else {
                $fileHealth.Status = "Corrupted"
                $fileHealth.Issues += $integrity.Message
                $results.OverallHealth = "Degraded"
            }
            
            # Check if up to date
            $upToDate = Test-FileUpToDate -FilePath $filePath -Branch $Branch
            if (-not $upToDate.UpToDate) {
                $fileHealth.Issues += "File may be outdated"
                if ($results.OverallHealth -eq "Healthy") {
                    $results.OverallHealth = "Outdated"
                }
            }
        } else {
            $fileHealth.Status = "Missing"
            # All files in manifest are considered required
            $results.OverallHealth = "Degraded"
            $fileHealth.Issues += "Required file is missing"
        }
        
        $results.FileHealth[$filePath] = $fileHealth
    }
    
    # Check directory health
    $requiredDirectories = @(
        ".github",
        ".github/instructions",
        ".github/prompts",
        ".vscode"
    )
    
    foreach ($directory in $requiredDirectories) {
        $results.DirectoryHealth[$directory] = @{
            Exists = Test-Path $directory -PathType Container
            Status = if (Test-Path $directory -PathType Container) { "Present" } else { "Missing" }
        }
    }
    
    # Check .gitignore status
    $gitIgnorePath = ".gitignore"
    if (Test-Path $gitIgnorePath) {
        $gitIgnoreContent = Get-Content $gitIgnorePath -Raw
        $hasAIEntries = $gitIgnoreContent -match "AI Infrastructure"
        
        $results.GitIgnoreStatus = @{
            Exists = $true
            HasAIEntries = $hasAIEntries
            Status = if ($hasAIEntries) { "Configured" } else { "Not Configured" }
        }
        
        if (-not $hasAIEntries -and $results.InstallationStatus -eq "Complete") {
            $results.Recommendations += "Consider adding AI files to .gitignore"
        }
    } else {
        $results.GitIgnoreStatus = @{
            Exists = $false
            HasAIEntries = $false
            Status = "Missing"
        }
    }
    
    # Add specific recommendations based on health
    switch ($results.OverallHealth) {
        "Missing" {
            $results.Recommendations += "Install AI infrastructure to enable enhanced development experience"
        }
        "Degraded" {
            $results.Recommendations += "Reinstall missing or corrupted files"
            $results.Recommendations += "Verify internet connectivity and repository access"
        }
        "Outdated" {
            $results.Recommendations += "Update AI infrastructure files to latest versions"
        }
        "Healthy" {
            $results.Recommendations += "AI infrastructure is functioning properly"
        }
    }
    
    return $results
}

function Test-InstallationDependencies {
    <#
    .SYNOPSIS
    Test all dependencies required for installation
    #>
    param(
        [string]$Branch = "exp/terraform_copilot"
    )
    
    $results = @{
        OverallValid = $true
        Dependencies = @{}
        MissingDependencies = @()
        Recommendations = @()
    }
    
    # Test workspace
    $workspaceTest = Test-WorkspaceValid
    $results.Dependencies["Workspace"] = $workspaceTest
    if (-not $workspaceTest.Valid) {
        $results.OverallValid = $false
        $results.MissingDependencies += "Valid Terraform AzureRM workspace"
        $results.Recommendations += $workspaceTest.Reason
    }
    
    # Test git repository
    $gitTest = Test-GitRepository
    $results.Dependencies["Git"] = $gitTest
    if (-not $gitTest.Valid) {
        $results.OverallValid = $false
        $results.MissingDependencies += "Valid git repository"
        $results.Recommendations += $gitTest.Reason
    }
    
    # Test internet connectivity
    $internetTest = Test-InternetConnectivity
    $results.Dependencies["Internet"] = $internetTest
    if (-not $internetTest.Connected) {
        $results.OverallValid = $false
        $results.MissingDependencies += "Internet connectivity"
        $results.Recommendations += $internetTest.Reason
    }
    
    # Test system requirements
    $systemTest = Test-SystemRequirements
    $results.Dependencies["System"] = $systemTest
    if (-not $systemTest.OverallValid) {
        $results.OverallValid = $false
        
        if (-not $systemTest.PowerShell.Valid) {
            $results.MissingDependencies += "PowerShell $($systemTest.PowerShell.MinimumVersion)+"
            $results.Recommendations += $systemTest.PowerShell.Reason
        }
        
        if (-not $systemTest.ExecutionPolicy.Valid) {
            $results.MissingDependencies += "Proper PowerShell execution policy"
            $results.Recommendations += $systemTest.ExecutionPolicy.Reason
        }
        
        if (-not $systemTest.Commands.Valid) {
            foreach ($command in $systemTest.Commands.Commands.Keys) {
                if (-not $systemTest.Commands.Commands[$command].Available) {
                    $results.MissingDependencies += "Command: $command"
                    $results.Recommendations += "Install $command and ensure it's in PATH"
                }
            }
        }
    }
    
    # Test branch access
    if ($results.OverallValid) {
        try {
            $apiUrl = "https://api.github.com/repos/hashicorp/terraform-provider-azurerm/branches/$Branch"
            # Disable progress bar to prevent console flashing
            $ProgressPreference = 'SilentlyContinue'
            $branchData = Invoke-RestMethod -Uri $apiUrl -Method GET -TimeoutSec 10
            
            $results.Dependencies["BranchAccess"] = @{
                Valid = $true
                Branch = $Branch
                BranchName = $branchData.name
                CommitSha = $branchData.commit.sha
                LastModified = $branchData.commit.commit.committer.date
                Reason = "Branch is accessible and up-to-date"
            }
        }
        catch {
            $results.Dependencies["BranchAccess"] = @{
                Valid = $false
                Branch = $Branch
                Reason = "Cannot access branch '$Branch': $($_.Exception.Message)"
            }
            $results.OverallValid = $false
            $results.MissingDependencies += "Access to branch '$Branch'"
            $results.Recommendations += "Verify branch name and GitHub connectivity"
        }
        finally {
            # Restore progress preference
            $ProgressPreference = 'Continue'
        }
    }
    
    return $results
}

function Get-ValidationReport {
    <#
    .SYNOPSIS
    Generate a comprehensive validation report
    #>
    param(
        [string]$Branch = "exp/terraform_copilot",
        [bool]$IncludeRecommendations = $true
    )
    
    $report = @{
        Timestamp = Get-Date
        Summary = @{}
        Details = @{}
        Recommendations = @()
    }
    
    # Run all validations
    $report.Details.SystemRequirements = Test-SystemRequirements
    $report.Details.Workspace = Test-WorkspaceValid
    $report.Details.Git = Test-GitRepository
    $report.Details.Internet = Test-InternetConnectivity
    $report.Details.AIInfrastructure = Test-AIInfrastructure -Branch $Branch
    $report.Details.Dependencies = Test-InstallationDependencies -Branch $Branch
    
    # Generate summary
    $report.Summary.SystemReady = $report.Details.SystemRequirements.OverallValid
    $report.Summary.WorkspaceValid = $report.Details.Workspace.Valid
    $report.Summary.GitReady = $report.Details.Git.Valid
    $report.Summary.InternetConnected = $report.Details.Internet.Connected
    $report.Summary.AIInfrastructureHealth = $report.Details.AIInfrastructure.OverallHealth
    $report.Summary.DependenciesValid = $report.Details.Dependencies.OverallValid
    
    $report.Summary.OverallStatus = if (
        $report.Summary.SystemReady -and 
        $report.Summary.WorkspaceValid -and 
        $report.Summary.GitReady -and 
        $report.Summary.InternetConnected -and 
        $report.Summary.DependenciesValid
    ) { "Ready" } else { "Not Ready" }
    
    # Collect recommendations
    if ($IncludeRecommendations) {
        # Check if we have a branch safety violation
        $hasBranchSafetyViolation = -not $report.Summary.GitReady -and 
                                   $report.Details.Git.Reason -like "*SAFETY VIOLATION*"
        
        if (-not $report.Summary.SystemReady) {
            $report.Recommendations += "Address system requirement issues"
        }
        if (-not $report.Summary.WorkspaceValid -and -not $hasBranchSafetyViolation) {
            $report.Recommendations += $report.Details.Workspace.Reason
        }
        if (-not $report.Summary.GitReady) {
            $report.Recommendations += $report.Details.Git.Reason
        }
        if (-not $report.Summary.InternetConnected -and -not $hasBranchSafetyViolation) {
            $report.Recommendations += $report.Details.Internet.Reason
        }
        
        # Only add AI infrastructure recommendations if no branch safety violation
        if (-not $hasBranchSafetyViolation) {
            $report.Recommendations += $report.Details.AIInfrastructure.Recommendations
            $report.Recommendations += $report.Details.Dependencies.Recommendations
        }
        
        # Remove duplicates
        $report.Recommendations = $report.Recommendations | Select-Object -Unique
    }
    
    return $report
}

function Get-CurrentBranch {
    <#
    .SYNOPSIS
    Get the current git branch name
    .DESCRIPTION
    Returns the current git branch name using git rev-parse
    .PARAMETER WorkspaceRoot
    Optional workspace root directory. If not provided, uses current directory.
    .EXAMPLE
    $branch = Get-CurrentBranch
    $branch = Get-CurrentBranch -WorkspaceRoot "C:\path\to\repo"
    #>
    param(
        [string]$WorkspaceRoot
    )
    
    try {
        if ($WorkspaceRoot -and (Test-Path $WorkspaceRoot)) {
            Push-Location $WorkspaceRoot
            try {
                $branch = git rev-parse --abbrev-ref HEAD 2>$null
                if ($LASTEXITCODE -eq 0 -and $branch) {
                    return $branch.Trim()
                }
            }
            finally {
                Pop-Location
            }
        } else {
            $branch = git rev-parse --abbrev-ref HEAD 2>$null
            if ($LASTEXITCODE -eq 0 -and $branch) {
                return $branch.Trim()
            }
        }
        return "unknown"
    }
    catch {
        return "unknown"
    }
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

function Get-DynamicSpacing {
    <#
    .SYNOPSIS
    Calculate dynamic spacing for aligned output based on branch type
    #>
    param(
        [ValidateSet("source", "feature", "unknown")]
        [string]$BranchType = "feature"
    )
    
    # Calculate the length of the branch detection line
    $branchLabels = @{
        "source"  = "SOURCE BRANCH DETECTED"
        "feature" = "FEATURE BRANCH DETECTED"
        "unknown" = "UNKNOWN BRANCH"
    }
    
    $workspaceLabel = "WORKSPACE"
    $branchLabel = $branchLabels[$BranchType]
    
    # Calculate spacing needed to align colons
    $branchLineLength = $branchLabel.Length + 2  # +2 for ": "
    $workspaceLineLength = $workspaceLabel.Length + 2  # +2 for ": "
    
    $maxLength = [Math]::Max($branchLineLength, $workspaceLineLength)
    $spacingNeeded = $maxLength - $workspaceLineLength
    
    return " " * $spacingNeeded
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
    The type of branch (source, feature, unknown) for dynamic spacing calculation
    
    .OUTPUTS
    Returns verification results and displays status to console
    
    .NOTES
    This function maintains source repository awareness and provides different
    behavior for source vs target repositories.
    #>
    param(
        [ValidateSet("source", "feature", "unknown")]
        [string]$BranchType = "feature"
    )
    
    # Use the dynamically determined workspace root
    $workspaceRoot = $Global:WorkspaceRoot
    Push-Location $workspaceRoot
    
    try {
        # CRITICAL: Check if we're on the source branch/repository
        $isSourceRepo = Test-SourceRepository
        
        $results = @{
            Success = $true
            Files = @()
            Issues = @()
            IsSourceRepo = $isSourceRepo
        }
        
        Write-Host $("=" * 60) -ForegroundColor Cyan
        Write-Host "Verifying AI infrastructure files" -ForegroundColor Cyan
        Write-Host $("=" * 60) -ForegroundColor Cyan
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
        
        # Check for deprecated files (files that exist but are no longer in manifest)
        if (-not $isSourceRepo) {
            $deprecatedFiles = Remove-DeprecatedFiles -ManifestConfig $ManifestConfig -WorkspaceRoot $workspaceRoot -DryRun $true -Quiet $true
            if ($deprecatedFiles.Count -gt 0) {
                $results.Issues += "Found $($deprecatedFiles.Count) deprecated files"
                Write-Host ""
                Write-Host "TIP: Deprecated files will be automatically removed during installation" -ForegroundColor Cyan
            }
        }
        
        if ($results.Issues.Count -gt 0) {
            $results.Success = $false
            Write-Host "Issues found:" -ForegroundColor Yellow
            foreach ($issue in $results.Issues) {
                Write-Host "  - $issue" -ForegroundColor Red
            }
            
            if (-not $isSourceRepo) {
                Write-Host ""
                if ($RepoDirectory) {
                    Write-Host "TIP: To install missing files, run: .\install-copilot-setup.ps1 -RepoDirectory `"$RepoDirectory`"" -ForegroundColor Cyan
                } else {
                    Write-Host "TIP: To install missing files, run: .\install-copilot-setup.ps1 -RepoDirectory `"<path-to-your-repo>`"" -ForegroundColor Cyan
                }
            }
        } else {
            if ($isSourceRepo) {
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

function Get-WorkspaceVerificationData {
    <#
    .SYNOPSIS
    Returns AI infrastructure verification data without any UI output
    
    .DESCRIPTION
    Checks for the presence of AI infrastructure files and returns structured data.
    This is a pure data function - no UI output is performed.
    
    .PARAMETER ManifestConfig
    The manifest configuration to use for verification
    
    .OUTPUTS
    Hashtable with verification results including Success, Files, Issues, and summary information
    #>
    param(
        [hashtable]$ManifestConfig = $null
    )
    
    $workspaceRoot = $Global:WorkspaceRoot
    if (-not $workspaceRoot) {
        throw "Global workspace root not set"
    }
    
    # Use provided manifest configuration or get it directly  
    if ($ManifestConfig) {
        $manifestConfig = $ManifestConfig
    } else {
        # Fallback: require ConfigParser to be loaded in parent scope
        if (-not (Get-Command Get-ManifestConfig -ErrorAction SilentlyContinue)) {
            throw "ManifestConfig parameter required or Get-ManifestConfig must be available"
        }
        $manifestConfig = Get-ManifestConfig
    }
    
    Push-Location $workspaceRoot
    
    try {
        # CRITICAL: Check if we're on the source branch/repository
        $isSourceRepo = Test-SourceRepository
        
        $results = @{
            Success = $true
            Files = @()
            Issues = @()
            IsSourceRepo = $isSourceRepo
            WorkspaceRoot = $workspaceRoot
            Summary = @{
                FilesFound = 0
                FilesMissing = 0
                DeprecatedFiles = 0
            }
        }
        
        # Check main instructions file
        $instructionsFile = $Global:InstallerConfig.Files.Instructions.Target
        if (Test-Path $instructionsFile) {
            $results.Files += @{
                Path = $instructionsFile
                Status = "Present"
                Description = "Main Copilot instructions"
                RelativePath = (Resolve-Path $instructionsFile -Relative)
            }
            $results.Summary.FilesFound++
        } else {
            $results.Files += @{
                Path = $instructionsFile
                Status = "Missing"
                Description = "Main Copilot instructions"
                RelativePath = $instructionsFile
            }
            $results.Issues += "Missing: $instructionsFile"
            $results.Summary.FilesMissing++
        }
        
        # Check instructions directory
        $instructionsDir = $Global:InstallerConfig.Files.InstructionFiles.Target
        if (Test-Path $instructionsDir -PathType Container) {
            $childResults = @()
            
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
                    $childResults += @{
                        File = $file
                        Status = "Present"
                    }
                    $results.Summary.FilesFound++
                } else {
                    $childResults += @{
                        File = $file
                        Status = "Missing"
                    }
                    $results.Issues += "Missing: $filePath"
                    $results.Summary.FilesMissing++
                }
            }
            
            $results.Files += @{
                Path = $instructionsDir
                Status = "Present"
                Description = "Instructions directory"
                RelativePath = (Resolve-Path $instructionsDir -Relative) + "/"
                Children = $childResults
            }
        } else {
            $results.Files += @{
                Path = $instructionsDir
                Status = "Missing"
                Description = "Instructions directory"
                RelativePath = $instructionsDir
            }
            $results.Issues += "Missing: $instructionsDir"
            $results.Summary.FilesMissing++
        }
        
        # Check prompts directory
        $promptsDir = $Global:InstallerConfig.Files.PromptFiles.Target
        if (Test-Path $promptsDir -PathType Container) {
            $childResults = @()
            
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
                    $childResults += @{
                        File = $file
                        Status = "Present"
                    }
                    $results.Summary.FilesFound++
                } else {
                    $childResults += @{
                        File = $file
                        Status = "Missing"
                    }
                    $results.Issues += "Missing: $filePath"
                    $results.Summary.FilesMissing++
                }
            }
            
            $results.Files += @{
                Path = $promptsDir
                Status = "Present"
                Description = "Prompts directory"
                RelativePath = (Resolve-Path $promptsDir -Relative) + "/"
                Children = $childResults
            }
        } else {
            $results.Files += @{
                Path = $promptsDir
                Status = "Missing"
                Description = "Prompts directory"
                RelativePath = $promptsDir
            }
            $results.Issues += "Missing: $promptsDir"
            $results.Summary.FilesMissing++
        }
        
        # Check .vscode directory and settings
        $vscodeDir = Join-Path $workspaceRoot ".vscode"
        if (Test-Path $vscodeDir -PathType Container) {
            $settingsFile = Join-Path $vscodeDir "settings.json"
            $vscodeChildren = @()
            
            if (Test-Path $settingsFile) {
                $vscodeChildren += @{
                    File = "settings.json"
                    Status = "Present"
                }
                $results.Summary.FilesFound++
            } else {
                $vscodeChildren += @{
                    File = "settings.json"
                    Status = "Missing"
                }
                $results.Issues += "Missing: $settingsFile"
                $results.Summary.FilesMissing++
            }
            
            $results.Files += @{
                Path = $vscodeDir
                Status = "Present"
                Description = ".vscode directory"
                RelativePath = ".vscode/"
                Children = $vscodeChildren
            }
        } else {
            $results.Files += @{
                Path = $vscodeDir
                Status = "Missing"
                Description = ".vscode directory"
                RelativePath = ".vscode/"
            }
            $results.Issues += "Missing: $vscodeDir"
            $results.Summary.FilesMissing++
        }
        
        # Check for deprecated files (files that exist but are no longer in manifest)
        if (-not $isSourceRepo) {
            $deprecatedFiles = Remove-DeprecatedFiles -ManifestConfig $manifestConfig -WorkspaceRoot $workspaceRoot -DryRun $true -Quiet $true
            $results.Summary.DeprecatedFiles = $deprecatedFiles.Count
            if ($deprecatedFiles.Count -gt 0) {
                $results.Issues += "Found $($deprecatedFiles.Count) deprecated files"
            }
        }
        
        # Final success determination
        $results.Success = ($results.Issues.Count -eq 0)
        
        return $results
    }
    finally {
        Pop-Location
    }
}

function Get-WorkspaceVerificationData {
    <#
    .SYNOPSIS
    Returns workspace verification data without UI output
    
    .DESCRIPTION
    Performs the same verification as Invoke-VerifyWorkspace but returns
    structured data instead of outputting to console. This allows the main
    script to control all UI formatting.
    
    .OUTPUTS
    Hashtable with Success, Issues, and Details properties
    #>
    
    $verificationData = @{
        Success = $true
        Issues = @()
        Details = @{}
    }
    
    try {
        # Initialize configuration
        Initialize-Configuration
        
        # Get basic workspace info
        $currentBranch = Get-CurrentBranch -WorkspaceRoot $Global:WorkspaceRoot
        $isSourceRepo = Test-SourceRepository
        $workspaceType = if ($isSourceRepo) { "Source Repository" } else { "Feature Branch" }
        
        $verificationData.Details.WorkspaceType = $workspaceType
        $verificationData.Details.Branch = $currentBranch
        
        # Verify workspace structure
        $requiredDirs = @("internal", ".github")
        $checkedDirs = 0
        
        foreach ($dir in $requiredDirs) {
            $dirPath = Join-Path $Global:WorkspaceRoot $dir
            if (-not (Test-Path $dirPath)) {
                $verificationData.Issues += "Required directory missing: $dir"
                $verificationData.Success = $false
            } else {
                $checkedDirs++
            }
        }
        
        $verificationData.Details.DirectoriesChecked = $checkedDirs
        
        # Verify key files
        $requiredFiles = @("go.mod", "main.go")
        $checkedFiles = 0
        
        foreach ($file in $requiredFiles) {
            $filePath = Join-Path $Global:WorkspaceRoot $file
            if (-not (Test-Path $filePath)) {
                $verificationData.Issues += "Required file missing: $file"
                $verificationData.Success = $false
            } else {
                $checkedFiles++
            }
        }
        
        $verificationData.Details.FilesChecked = $checkedFiles
        
        # Check AI infrastructure status
        $aiFilesStatus = Get-AIInfrastructureStatus -WorkspaceRoot $Global:WorkspaceRoot -ManifestConfig $Global:ManifestConfig
        $verificationData.Details.AIInfrastructure = $aiFilesStatus
        
        if ($aiFilesStatus.MissingFiles -gt 0) {
            $verificationData.Issues += "AI infrastructure incomplete: $($aiFilesStatus.MissingFiles) missing files"
        }
        
    }
    catch {
        $verificationData.Success = $false
        $verificationData.Issues += "Verification failed: $($_.Exception.Message)"
    }
    
    return $verificationData
}

#endregion

#region Export Module Members

Export-ModuleMember -Function @(
    'Test-SystemRequirements',
    'Test-PreInstallation',
    'Test-PostInstallation',
    'Test-AIInfrastructure',
    'Test-InstallationDependencies',
    'Get-ValidationReport',
    'Test-InternetConnectivity',
    'Test-WorkspaceValid',
    'Test-GitRepository',
    'Get-CurrentBranch',
    'Test-SourceRepository',
    'Get-DynamicSpacing',
    'Invoke-VerifyWorkspace',
    'Get-WorkspaceVerificationData'
)

#endregion
