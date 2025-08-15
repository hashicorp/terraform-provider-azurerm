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
        [string]$Branch = "exp/terraform_copilot"
    )
    
    $config = Get-InstallationConfig -Branch $Branch
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
    foreach ($filePath in $config.Files.Keys) {
        $fileInfo = $config.Files[$filePath]
        
        $fileValidation = @{
            Exists = Test-Path $filePath
            Required = $fileInfo.Required
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
        [string]$Branch = "exp/terraform_copilot"
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
    $config = Get-InstallationConfig -Branch $Branch
    
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
    foreach ($filePath in $config.Files.Keys) {
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
            if ($config.Files[$filePath].Required) {
                $results.OverallHealth = "Degraded"
                $fileHealth.Issues += "Required file is missing"
            }
        }
        
        $results.FileHealth[$filePath] = $fileHealth
    }
    
    # Check directory health
    foreach ($directory in $config.RequiredDirectories) {
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
    'Test-GitRepository'
)

#endregion
