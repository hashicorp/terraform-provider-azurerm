# FileOperations Module for Terraform AzureRM Provider AI Setup
# Handles file downloading, installation, removal, and management

#region Private Functions

function Assert-DirectoryExists {
    <#
    .SYNOPSIS
    Ensure a directory exists, creating it if necessary
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Path
    )
    
    if (-not (Test-Path $Path -PathType Container)) {
        try {
            New-Item -Path $Path -ItemType Directory -Force | Out-Null
            return $true
        }
        catch {
            Write-Error "Failed to create directory '$Path': $($_.Exception.Message)"
            return $false
        }
    }
    
    return $true
}

function Get-GitHubFileContent {
    <#
    .SYNOPSIS
    Download file content from GitHub
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Url,
        
        [int]$TimeoutSeconds = 30
    )
    
    try {
        $response = Invoke-WebRequest -Uri $Url -UseBasicParsing -TimeoutSec $TimeoutSeconds
        return @{
            Success = $true
            Content = $response.Content
            Size = $response.Content.Length
        }
    }
    catch {
        return @{
            Success = $false
            Content = $null
            Size = 0
            ErrorMessage = $_.Exception.Message
        }
    }
}

function Get-FileHash {
    <#
    .SYNOPSIS
    Get file hash for comparison
    #>
    param(
        [Parameter(Mandatory)]
        [string]$FilePath
    )
    
    if (Test-Path $FilePath) {
        try {
            $hash = Get-FileHash -Path $FilePath -Algorithm SHA256
            return $hash.Hash
        }
        catch {
            return $null
        }
    }
    
    return $null
}

#endregion

#region Public Functions

function Install-AIFile {
    <#
    .SYNOPSIS
    Download and install a single AI infrastructure file
    #>
    param(
        [Parameter(Mandatory)]
        [string]$FilePath,
        
        [Parameter(Mandatory)]
        [string]$DownloadUrl,
        
        [bool]$Force = $false,
        
        [bool]$DryRun = $false
    )
    
    $result = @{
        FilePath = $FilePath
        Success = $false
        Action = "None"
        Message = ""
        Size = 0
    }
    
    try {
        # Check if file already exists
        $fileExists = Test-Path $FilePath
        
        if ($fileExists -and -not $Force) {
            $result.Action = "Skipped"
            $result.Success = $true
            $result.Message = "File already exists (use -Force to overwrite)"
            return $result
        }
        
        # Create directory if needed
        $directory = Split-Path $FilePath -Parent
        if ($directory -and -not (Assert-DirectoryExists $directory)) {
            $result.Message = "Failed to create directory: $directory"
            return $result
        }
        
        if ($DryRun) {
            $result.Action = if ($fileExists) { "Would Overwrite" } else { "Would Download" }
            $result.Success = $true
            $result.Message = "Dry run - no changes made"
            return $result
        }
        
        # Download file
        Write-Verbose "Downloading: $DownloadUrl"
        $response = Invoke-WebRequest -Uri $DownloadUrl -UseBasicParsing -ErrorAction Stop
        
        # Save file
        [System.IO.File]::WriteAllBytes($FilePath, $response.Content)
        
        # Verify file was created
        if (Test-Path $FilePath) {
            $fileInfo = Get-Item $FilePath
            $result.Size = $fileInfo.Length
            $result.Action = if ($fileExists) { "Overwritten" } else { "Downloaded" }
            $result.Success = $true
            $result.Message = "Successfully installed ($($result.Size) bytes)"
        } else {
            $result.Message = "File was not created"
        }
    }
    catch {
        $result.Message = "Download failed: $($_.Exception.Message)"
    }
    
    return $result
}

function Install-AllAIFiles {
    <#
    .SYNOPSIS
    Install all AI infrastructure files
    #>
    param(
        [bool]$Force = $false,
        [bool]$DryRun = $false,
        [string]$Branch = "exp/terraform_copilot"
    )
    
    $config = Get-InstallationConfig -Branch $Branch
    $results = @{
        TotalFiles = $config.Files.Count
        Successful = 0
        Failed = 0
        Skipped = 0
        Files = @{}
        OverallSuccess = $true
    }
    
    Write-Progress -Activity "Installing AI Infrastructure" -Status "Preparing..." -PercentComplete 0
    
    $fileIndex = 0
    foreach ($filePath in $config.Files.Keys) {
        $fileIndex++
        $fileInfo = $config.Files[$filePath]
        $downloadUrl = $config.BaseUrl + $fileInfo.Url
        
        $percentComplete = [math]::Round(($fileIndex / $config.Files.Count) * 100)
        Write-Progress -Activity "Installing AI Infrastructure" -Status "Processing: $filePath" -PercentComplete $percentComplete
        
        $fileResult = Install-AIFile -FilePath $filePath -DownloadUrl $downloadUrl -Force $Force -DryRun $DryRun
        $results.Files[$filePath] = $fileResult
        
        switch ($fileResult.Action) {
            { $_ -in @("Downloaded", "Overwritten") } { $results.Successful++ }
            "Skipped" { $results.Skipped++ }
            default { 
                $results.Failed++
                $results.OverallSuccess = $false
            }
        }
    }
    
    Write-Progress -Activity "Installing AI Infrastructure" -Completed
    
    return $results
}

function Remove-AIFile {
    <#
    .SYNOPSIS
    Remove a single AI infrastructure file
    #>
    param(
        [Parameter(Mandatory)]
        [string]$FilePath,
        
        [bool]$DryRun = $false
    )
    
    $result = @{
        FilePath = $FilePath
        Success = $false
        Action = "None"
        Message = ""
    }
    
    try {
        if (-not (Test-Path $FilePath)) {
            $result.Action = "Not Found"
            $result.Success = $true
            $result.Message = "File does not exist"
            return $result
        }
        
        if ($DryRun) {
            $result.Action = "Would Remove"
            $result.Success = $true
            $result.Message = "Dry run - no changes made"
            return $result
        }
        
        # Remove file
        Remove-Item -Path $FilePath -Force -ErrorAction Stop
        
        $result.Action = "Removed"
        $result.Success = $true
        $result.Message = "File successfully removed"
    }
    catch {
        $result.Message = "Failed to remove file: $($_.Exception.Message)"
    }
    
    return $result
}

function Remove-AllAIFiles {
    <#
    .SYNOPSIS
    Remove all AI infrastructure files and clean up
    #>
    param(
        [bool]$DryRun = $false,
        [string]$Branch = "exp/terraform_copilot"
    )
    
    $config = Get-InstallationConfig -Branch $Branch
    $results = @{
        TotalFiles = $config.Files.Count
        Removed = 0
        NotFound = 0
        Failed = 0
        Files = @{}
        Directories = @{}
        OverallSuccess = $true
    }
    
    Write-Progress -Activity "Removing AI Infrastructure" -Status "Preparing..." -PercentComplete 0
    
    # Remove files
    $fileIndex = 0
    foreach ($filePath in $config.Files.Keys) {
        $fileIndex++
        $percentComplete = [math]::Round(($fileIndex / $config.Files.Count) * 50)
        Write-Progress -Activity "Removing AI Infrastructure" -Status "Removing: $filePath" -PercentComplete $percentComplete
        
        $fileResult = Remove-AIFile -FilePath $filePath -DryRun $DryRun
        $results.Files[$filePath] = $fileResult
        
        switch ($fileResult.Action) {
            "Removed" { $results.Removed++ }
            "Not Found" { $results.NotFound++ }
            default { 
                $results.Failed++
                $results.OverallSuccess = $false
            }
        }
    }
    
    # Remove empty directories
    Write-Progress -Activity "Removing AI Infrastructure" -Status "Cleaning up directories..." -PercentComplete 75
    
    $directoriesToCheck = @(
        ".github/instructions",
        ".github/prompts"
    )
    
    foreach ($dir in $directoriesToCheck) {
        $dirResult = @{
            Path = $dir
            Action = "None"
            Success = $true
            Message = ""
        }
        
        if (Test-Path $dir -PathType Container) {
            $dirContents = Get-ChildItem $dir -Force
            if ($dirContents.Count -eq 0) {
                if ($DryRun) {
                    $dirResult.Action = "Would Remove"
                    $dirResult.Message = "Empty directory would be removed"
                } else {
                    try {
                        Remove-Item -Path $dir -Force -ErrorAction Stop
                        $dirResult.Action = "Removed"
                        $dirResult.Message = "Empty directory removed"
                    }
                    catch {
                        $dirResult.Action = "Failed"
                        $dirResult.Success = $false
                        $dirResult.Message = "Failed to remove directory: $($_.Exception.Message)"
                        $results.OverallSuccess = $false
                    }
                }
            } else {
                $dirResult.Action = "Not Empty"
                $dirResult.Message = "Directory contains other files"
            }
        } else {
            $dirResult.Action = "Not Found"
            $dirResult.Message = "Directory does not exist"
        }
        
        $results.Directories[$dir] = $dirResult
    }
    
    Write-Progress -Activity "Removing AI Infrastructure" -Completed
    
    return $results
}

function Remove-DeprecatedFiles {
    <#
    .SYNOPSIS
    Removes files that were previously installed but are no longer in the manifest
    
    .DESCRIPTION
    Scans the target directories for files that exist but are not listed in the 
    current manifest configuration, indicating they were deprecated/removed
    
    .PARAMETER ManifestConfig
    The manifest configuration containing current file lists
    
    .PARAMETER WorkspaceRoot
    The root directory of the workspace
    
    .PARAMETER DryRun
    If true, only reports what would be removed without actually removing files
    
    .PARAMETER Quiet
    If true, suppresses output (useful for verification checks)
    
    .OUTPUTS
    Array of deprecated files found
    #>
    param(
        [Parameter(Mandatory)]
        [hashtable]$ManifestConfig,
        
        [Parameter(Mandatory)]
        [string]$WorkspaceRoot,
        
        [bool]$DryRun = $false,
        [bool]$Quiet = $false
    )
    
    $deprecatedFiles = @()
    
    # Check for deprecated instruction files
    $instructionsDir = Join-Path $WorkspaceRoot ".github\instructions"
    if (Test-Path $instructionsDir -PathType Container) {
        $currentFiles = $ManifestConfig.Sections.INSTRUCTION_FILES
        $existingFiles = Get-ChildItem $instructionsDir -File | Where-Object { $_.Name -like "*.instructions.md" }
        
        foreach ($existingFile in $existingFiles) {
            if ($existingFile.Name -notin $currentFiles) {
                $deprecatedFiles += @{
                    Path = $existingFile.FullName
                    Type = "Instruction"
                    Name = $existingFile.Name
                    RelativePath = $existingFile.FullName.Replace($WorkspaceRoot, "").TrimStart('\').TrimStart('/')
                }
            }
        }
    }
    
    # Check for deprecated prompt files
    $promptsDir = Join-Path $WorkspaceRoot ".github\prompts"
    if (Test-Path $promptsDir -PathType Container) {
        $currentPrompts = $ManifestConfig.Sections.PROMPT_FILES
        $existingPrompts = Get-ChildItem $promptsDir -File | Where-Object { $_.Name -like "*.prompt.md" }
        
        foreach ($existingPrompt in $existingPrompts) {
            if ($existingPrompt.Name -notin $currentPrompts) {
                $deprecatedFiles += @{
                    Path = $existingPrompt.FullName
                    Type = "Prompt"
                    Name = $existingPrompt.Name
                    RelativePath = $existingPrompt.FullName.Replace($WorkspaceRoot, "").TrimStart('\').TrimStart('/')
                }
            }
        }
    }
    
    if ($deprecatedFiles.Count -gt 0) {
        if (-not $Quiet) {
            Write-Host ""
            Write-Host "Found deprecated files (no longer in manifest):" -ForegroundColor Yellow
            foreach ($file in $deprecatedFiles) {
                Write-Host "  [$($file.Type)] $($file.RelativePath)" -ForegroundColor Gray
            }
        }
        
        if (-not $DryRun) {
            if (-not $Quiet) {
                Write-Host ""
                $confirm = Read-Host "Remove deprecated files? (y/N)"
            } else {
                # Auto-approve in quiet mode (typically during installation)
                $confirm = 'y'
            }
            
            if ($confirm -eq 'y' -or $confirm -eq 'Y') {
                $removedCount = 0
                foreach ($file in $deprecatedFiles) {
                    try {
                        Remove-Item -Path $file.Path -Force
                        if (-not $Quiet) {
                            Write-Host "  Removed: $($file.RelativePath)" -ForegroundColor Green
                        }
                        $removedCount++
                    }
                    catch {
                        if (-not $Quiet) {
                            Write-Host "  Failed to remove: $($file.RelativePath) - $($_.Exception.Message)" -ForegroundColor Red
                        }
                    }
                }
                if (-not $Quiet) {
                    Write-Host ""
                    Write-Host "Removed $removedCount deprecated files." -ForegroundColor Green
                }
            } else {
                if (-not $Quiet) {
                    Write-Host "Deprecated files kept." -ForegroundColor Yellow
                }
            }
        }
    } elseif (-not $DryRun -and -not $Quiet) {
        Write-Host "No deprecated files found." -ForegroundColor Green
    }
    
    return $deprecatedFiles
}

function Update-GitIgnore {
    <#
    .SYNOPSIS
    Add or remove AI infrastructure entries from .gitignore
    #>
    param(
        [Parameter(Mandatory)]
        [ValidateSet("Add", "Remove")]
        [string]$Action,
        
        [bool]$DryRun = $false,
        [string]$Branch = "exp/terraform_copilot"
    )
    
    $config = Get-InstallationConfig -Branch $Branch
    $gitIgnorePath = ".gitignore"
    
    $result = @{
        Action = $Action
        Success = $false
        Message = ""
        EntriesProcessed = 0
        GitIgnoreExists = Test-Path $gitIgnorePath
    }
    
    try {
        # Read current .gitignore
        $gitIgnoreContent = @()
        if ($result.GitIgnoreExists) {
            $gitIgnoreContent = Get-Content $gitIgnorePath -ErrorAction Stop
        }
        
        $entriesSection = $config.GitIgnoreEntries
        $sectionStart = $entriesSection[0] # Comment line
        
        if ($Action -eq "Add") {
            # Check if section already exists
            $sectionExists = $gitIgnoreContent -contains $sectionStart
            
            if ($sectionExists -and -not $DryRun) {
                $result.Success = $true
                $result.Message = "AI infrastructure entries already exist in .gitignore"
                return $result
            }
            
            if ($DryRun) {
                $result.Success = $true
                $result.Message = "Would add $($entriesSection.Count) entries to .gitignore"
                $result.EntriesProcessed = $entriesSection.Count
                return $result
            }
            
            # Add entries
            $newContent = $gitIgnoreContent + "" + $entriesSection
            Set-Content -Path $gitIgnorePath -Value $newContent -ErrorAction Stop
            $result.EntriesProcessed = $entriesSection.Count
            $result.Message = "Added $($entriesSection.Count) entries to .gitignore"
        }
        elseif ($Action -eq "Remove") {
            # Find and remove the section
            $newContent = @()
            $inAISection = $false
            $removedCount = 0
            
            foreach ($line in $gitIgnoreContent) {
                if ($line -eq $sectionStart) {
                    $inAISection = $true
                    $removedCount++
                    continue
                }
                
                if ($inAISection) {
                    # Check if this line is part of our section
                    if ($line -in $entriesSection[1..($entriesSection.Count - 1)]) {
                        $removedCount++
                        continue
                    } elseif ([string]::IsNullOrWhiteSpace($line)) {
                        # Skip empty lines in our section
                        continue
                    } else {
                        # End of our section
                        $inAISection = $false
                    }
                }
                
                $newContent += $line
            }
            
            if ($removedCount -eq 0) {
                $result.Success = $true
                $result.Message = "No AI infrastructure entries found in .gitignore"
                return $result
            }
            
            if ($DryRun) {
                $result.Success = $true
                $result.Message = "Would remove $removedCount entries from .gitignore"
                $result.EntriesProcessed = $removedCount
                return $result
            }
            
            # Save updated content
            Set-Content -Path $gitIgnorePath -Value $newContent -ErrorAction Stop
            $result.EntriesProcessed = $removedCount
            $result.Message = "Removed $removedCount entries from .gitignore"
        }
        
        $result.Success = $true
    }
    catch {
        $result.Message = "Failed to update .gitignore: $($_.Exception.Message)"
    }
    
    return $result
}

function Backup-ExistingFile {
    <#
    .SYNOPSIS
    Create a backup of an existing file before overwriting
    #>
    param(
        [Parameter(Mandatory)]
        [string]$FilePath
    )
    
    if (-not (Test-Path $FilePath)) {
        return @{
            Success = $true
            Message = "File does not exist, no backup needed"
            BackupPath = $null
        }
    }
    
    try {
        $timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
        $backupPath = "$FilePath.backup.$timestamp"
        
        Copy-Item -Path $FilePath -Destination $backupPath -ErrorAction Stop
        
        return @{
            Success = $true
            Message = "Backup created successfully"
            BackupPath = $backupPath
        }
    }
    catch {
        return @{
            Success = $false
            Message = "Failed to create backup: $($_.Exception.Message)"
            BackupPath = $null
        }
    }
}

function Test-FileIntegrity {
    <#
    .SYNOPSIS
    Verify file integrity after download
    #>
    param(
        [Parameter(Mandatory)]
        [string]$FilePath,
        
        [int]$MinimumSize = 100
    )
    
    $result = @{
        Valid = $false
        Message = ""
        Size = 0
        Readable = $false
    }
    
    try {
        if (-not (Test-Path $FilePath)) {
            $result.Message = "File does not exist"
            return $result
        }
        
        $fileInfo = Get-Item $FilePath
        $result.Size = $fileInfo.Length
        
        if ($result.Size -lt $MinimumSize) {
            $result.Message = "File size ($($result.Size) bytes) is smaller than expected minimum ($MinimumSize bytes)"
            return $result
        }
        
        # Try to read file to verify it's not corrupted
        $content = Get-Content $FilePath -Raw -ErrorAction Stop
        $result.Readable = $true
        
        if ([string]::IsNullOrWhiteSpace($content)) {
            $result.Message = "File appears to be empty"
            return $result
        }
        
        $result.Valid = $true
        $result.Message = "File integrity verified"
    }
    catch {
        $result.Message = "File integrity check failed: $($_.Exception.Message)"
    }
    
    return $result
}

function Test-InternetConnectivity {
    <#
    .SYNOPSIS
    Test internet connectivity and GitHub access
    #>
    
    $connectivityResults = @{
        Connected = $false
        GitHubAccess = $false
        TestResults = @{}
        Reason = ""
    }
    
    # Test basic internet connectivity
    try {
        $githubUrl = "https://github.com"
        $response = Invoke-WebRequest -Uri $githubUrl -UseBasicParsing -TimeoutSec 10 -ErrorAction Stop
        
        $connectivityResults.TestResults["GitHub"] = @{
            Success = $true
            StatusCode = $response.StatusCode
            ResponseTime = "< 10s"
        }
        $connectivityResults.Connected = $true
        $connectivityResults.GitHubAccess = $true
        $connectivityResults.Reason = "Internet connectivity and GitHub access confirmed"
    }
    catch {
        $connectivityResults.TestResults["GitHub"] = @{
            Success = $false
            ErrorDetails = $_.Exception.Message
            ResponseTime = "Timeout"
        }
        $connectivityResults.Reason = "Cannot connect to GitHub: $($_.Exception.Message)"
    }
    
    # Test GitHub API access
    if ($connectivityResults.Connected) {
        try {
            $apiUrl = "https://api.github.com/repos/hashicorp/terraform-provider-azurerm"
            $apiResponse = Invoke-RestMethod -Uri $apiUrl -Method GET -TimeoutSec 10
            
            $connectivityResults.TestResults["GitHubAPI"] = @{
                Success = $true
                RepoName = $apiResponse.name
                LastUpdated = $apiResponse.updated_at
            }
        }
        catch {
            $connectivityResults.TestResults["GitHubAPI"] = @{
                Success = $false
                ErrorDetails = $_.Exception.Message
            }
            $connectivityResults.GitHubAccess = $false
            if ($connectivityResults.Reason -eq "Internet connectivity and GitHub access confirmed") {
                $connectivityResults.Reason = "GitHub API access failed: $($_.Exception.Message)"
            }
        }
    }
    
    return $connectivityResults
}

function Test-WorkspaceValid {
    <#
    .SYNOPSIS
    Validate that the current directory is a valid Terraform AzureRM Provider workspace
    #>
    
    $requiredFiles = @(
        "go.mod",
        "main.go",
        "internal",
        "website"
    )
    
    $requiredContent = @{
        "go.mod" = "terraform-provider-azurerm"
        "main.go" = "github.com/hashicorp/terraform-provider-azurerm"
    }
    
    # Check required files exist
    foreach ($file in $requiredFiles) {
        if (-not (Test-Path $file)) {
            return @{
                Valid = $false
                Reason = "Missing required file or directory: $file"
            }
        }
    }
    
    # Check required content
    foreach ($file in $requiredContent.Keys) {
        if (Test-Path $file) {
            $content = Get-Content $file -Raw -ErrorAction SilentlyContinue
            if ($content -notmatch [regex]::Escape($requiredContent[$file])) {
                return @{
                    Valid = $false
                    Reason = "File $file does not contain expected content: $($requiredContent[$file])"
                }
            }
        }
    }
    
    return @{
        Valid = $true
        Reason = "Valid Terraform AzureRM Provider workspace"
    }
}

function Test-GitRepository {
    <#
    .SYNOPSIS
    Validate git repository state and branch
    #>
    
    # Check if git repository
    if (-not (Test-Path ".git")) {
        return @{
            Valid = $false
            Reason = "Not a git repository"
        }
    }
    
    # Get current branch
    try {
        $currentBranch = git rev-parse --abbrev-ref HEAD 2>$null
        if ($LASTEXITCODE -ne 0) {
            return @{
                Valid = $false
                Reason = "Unable to determine current git branch"
            }
        }
        
        # Check if on source branch (would cause conflicts)
        if ($currentBranch -eq "exp/terraform_copilot") {
            return @{
                Valid = $false
                Reason = "Cannot run installer from source branch 'exp/terraform_copilot'. Switch to a different branch first."
            }
        }
        
        return @{
            Valid = $true
            CurrentBranch = $currentBranch
            Reason = "Valid git repository on branch '$currentBranch'"
        }
    }
    catch {
        return @{
            Valid = $false
            Reason = "Error checking git status: $($_.Exception.Message)"
        }
    }
}

function Get-WorkspaceStatus {
    <#
    .SYNOPSIS
    Get comprehensive workspace status including installed files
    #>
    
    $config = Get-InstallationConfig
    $workspaceStatus = @{
        Workspace = Test-WorkspaceValid
        Git = Test-GitRepository
        InstalledFiles = @{}
        MissingFiles = @()
        OutdatedFiles = @()
        TotalFiles = $config.Files.Count
        InstalledCount = 0
    }
    
    # Check each file
    foreach ($filePath in $config.Files.Keys) {
        $fileDefinition = $config.Files[$filePath]
        $fileStatus = @{
            Path = $filePath
            Exists = Test-Path $filePath
            Required = $fileDefinition.Required
            Type = $fileDefinition.Type
            Description = $fileDefinition.Description
            Size = 0
            LastModified = $null
        }
        
        if ($fileStatus.Exists) {
            $fileDetails = Get-Item $filePath -ErrorAction SilentlyContinue
            if ($fileDetails) {
                $fileStatus.Size = $fileDetails.Length
                $fileStatus.LastModified = $fileDetails.LastWriteTime
            }
            $workspaceStatus.InstalledCount++
        } else {
            if ($fileDefinition.Required) {
                $workspaceStatus.MissingFiles += $filePath
            }
        }
        
        $workspaceStatus.InstalledFiles[$filePath] = $fileStatus
    }
    
    return $workspaceStatus
}

function Get-FileFromGitHub {
    <#
    .SYNOPSIS
    Download a file from GitHub and save to local path
    #>
    param(
        [Parameter(Mandatory)]
        [string]$GitHubPath,
        
        [Parameter(Mandatory)]
        [string]$LocalPath,
        
        [string]$Branch = "exp/terraform_copilot"
    )
    
    $downloadResult = @{
        Success = $false
        Size = 0
        ErrorMessage = ""
        LocalPath = $LocalPath
    }
    
    try {
        # Ensure directory exists
        $directory = Split-Path $LocalPath -Parent
        if ($directory -and -not (Assert-DirectoryExists $directory)) {
            $downloadResult.ErrorMessage = "Failed to create directory: $directory"
            return $downloadResult
        }
        
        # Construct download URL
        $baseUrl = "https://raw.githubusercontent.com/hashicorp/terraform-provider-azurerm/$Branch"
        $downloadUrl = $baseUrl + $GitHubPath
        
        # Download file
        $downloadData = Get-GitHubFileContent -Url $downloadUrl
        
        if ($downloadData.Success) {
            # Save file
            [System.IO.File]::WriteAllBytes($LocalPath, $downloadData.Content)
            
            # Verify file was created
            if (Test-Path $LocalPath) {
                $fileData = Get-Item $LocalPath
                $downloadResult.Size = $fileData.Length
                $downloadResult.Success = $true
            } else {
                $downloadResult.ErrorMessage = "File was not created after download"
            }
        } else {
            $downloadResult.ErrorMessage = $downloadData.ErrorMessage
        }
    }
    catch {
        $downloadResult.ErrorMessage = "Download failed: $($_.Exception.Message)"
    }
    
    return $downloadResult
}

function Test-FileUpToDate {
    <#
    .SYNOPSIS
    Check if a local file is up to date with the remote version
    #>
    param(
        [Parameter(Mandatory)]
        [string]$FilePath,
        
        [string]$Branch = "exp/terraform_copilot"
    )
    
    if (-not (Test-Path $FilePath)) {
        return @{
            UpToDate = $false
            Reason = "File does not exist locally"
        }
    }
    
    try {
        $config = Get-InstallationConfig -Branch $Branch
        $fileConfig = $config.Files[$FilePath]
        
        if (-not $fileConfig) {
            return @{
                UpToDate = $false
                Reason = "File not found in configuration"
            }
        }
        
        $downloadUrl = $config.BaseUrl + $fileConfig.Url
        
        # Get remote file info
        $remoteResponse = Invoke-WebRequest -Uri $downloadUrl -Method HEAD -UseBasicParsing -ErrorAction Stop
        $remoteLastModified = $remoteResponse.Headers['Last-Modified']
        
        # Get local file info
        $localFile = Get-Item $FilePath
        
        if ($remoteLastModified) {
            $remoteDate = [DateTime]::Parse($remoteLastModified)
            if ($localFile.LastWriteTime -lt $remoteDate) {
                return @{
                    UpToDate = $false
                    Reason = "Local file is older than remote version"
                    LocalDate = $localFile.LastWriteTime
                    RemoteDate = $remoteDate
                }
            }
        }
        
        return @{
            UpToDate = $true
            Reason = "File appears to be up to date"
        }
    }
    catch {
        return @{
            UpToDate = $false
            Reason = "Unable to check remote file: $($_.Exception.Message)"
        }
    }
}

function Get-GitIgnoreStatus {
    <#
    .SYNOPSIS
    Check the status of .gitignore AI entries
    #>
    
    $gitIgnorePath = ".gitignore"
    
    if (-not (Test-Path $gitIgnorePath)) {
        return @{
            Exists = $false
            HasAIEntries = $false
            Status = "Missing"
        }
    }
    
    $gitIgnoreContent = Get-Content $gitIgnorePath -Raw
    $hasAIEntries = $gitIgnoreContent -match "AI Infrastructure"
    
    return @{
        Exists = $true
        HasAIEntries = $hasAIEntries
        Status = if ($hasAIEntries) { "Configured" } else { "Not Configured" }
    }
}

#endregion

#region Export Module Members

Export-ModuleMember -Function @(
    'Install-AIFile',
    'Install-AllAIFiles',
    'Remove-AIFile',
    'Remove-AllAIFiles',
    'Remove-DeprecatedFiles',
    'Update-GitIgnore',
    'Backup-ExistingFile',
    'Test-FileIntegrity',
    'Test-InternetConnectivity',
    'Test-WorkspaceValid',
    'Test-GitRepository',
    'Get-WorkspaceStatus',
    'Get-FileFromGitHub',
    'Test-FileUpToDate',
    'Get-GitIgnoreStatus'
)

#endregion
