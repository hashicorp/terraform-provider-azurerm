# FileOperations Module for Terraform AzureRM Provider AI Setup
# STREAMLINED VERSION - Contains only functions actually used by main script

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
        
        [bool]$DryRun = $false,
        
        [Parameter(Mandatory)]
        [string]$WorkspaceRoot
    )
    
    $result = @{
        FilePath = $FilePath
        Success = $false
        Action = "None"
        Message = ""
        Size = 0
        DebugInfo = @{}
    }
    
    try {
        # Resolve the full file path by joining workspace root with relative path
        $resolvedFilePath = Join-Path $WorkspaceRoot $FilePath
        $result.DebugInfo.WorkspaceRoot = $WorkspaceRoot
        $result.DebugInfo.OriginalPath = $FilePath
        $result.DebugInfo.ResolvedPath = $resolvedFilePath
        $result.DebugInfo.PathMethod = "Join-Path"
        
        # Update result with resolved path
        $result.FilePath = $resolvedFilePath
        
        # Record initial state
        $result.DebugInfo.StartTime = Get-Date
        $result.DebugInfo.DownloadUrl = $DownloadUrl
        $result.DebugInfo.TargetPath = $resolvedFilePath
        
        # Check if file already exists
        $fileExists = Test-Path $resolvedFilePath
        $result.DebugInfo.FileExisted = $fileExists
        
        if ($fileExists -and -not $Force) {
            $result.Action = "Skipped"
            $result.Success = $true
            $result.Message = "File already exists (use -Force to overwrite)"
            return $result
        }
        
        # Create directory if needed
        $directory = Split-Path $resolvedFilePath -Parent
        $result.DebugInfo.TargetDirectory = $directory
        
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
        
        $downloadStart = Get-Date
        
        # Hide progress bar during download
        $originalProgressPreference = $ProgressPreference
        $ProgressPreference = 'SilentlyContinue'
        
        try {
            $response = Invoke-WebRequest -Uri $DownloadUrl -UseBasicParsing -ErrorAction Stop
        }
        finally {
            $ProgressPreference = $originalProgressPreference
        }
        
        $downloadEnd = Get-Date
        
        if ($downloadEnd -and $downloadStart) {
            $result.DebugInfo.DownloadDuration = ($downloadEnd - $downloadStart).TotalMilliseconds
        } else {
            $result.DebugInfo.DownloadDuration = 0
        }
        $result.DebugInfo.ResponseSize = $response.Content.Length
        $result.DebugInfo.StatusCode = $response.StatusCode
        
        # Save file (all AI infrastructure files are text-based)
        try {
            # Use UTF8 encoding without BOM for all text files
            $utf8NoBom = New-Object System.Text.UTF8Encoding($false)
            [System.IO.File]::WriteAllText($resolvedFilePath, $response.Content, $utf8NoBom)
            $result.DebugInfo.SaveMethod = "WriteAllText (UTF8 no BOM)"
        } catch {
            $result.DebugInfo.SaveException = $_.Exception.Message
            throw
        }
        
        # Verify file was created
        if (Test-Path $resolvedFilePath) {
            $fileInfo = Get-Item $resolvedFilePath
            $result.Size = $fileInfo.Length
            $result.Action = if ($fileExists) { "Overwritten" } else { "Downloaded" }
            $result.Success = $true
            $result.Message = "Successfully installed ($($result.Size) bytes)"
            $result.DebugInfo.FinalSize = $result.Size
        } else {
            $result.Message = "File was not created"
        }
    }
    catch {
        $result.Message = "Download failed: $($_.Exception.Message)"
        $result.DebugInfo.Exception = $_.Exception.GetType().Name
        $result.DebugInfo.ExceptionMessage = $_.Exception.Message
        
        # Additional debug for specific error types
        if ($_.Exception -is [System.Net.WebException]) {
            $webEx = $_.Exception
            if ($webEx.Response) {
                $result.DebugInfo.HttpStatusCode = [int]$webEx.Response.StatusCode
                $result.DebugInfo.HttpStatusDescription = $webEx.Response.StatusDescription
            }
        }
    }
    
    $result.DebugInfo.EndTime = Get-Date
    if ($result.DebugInfo.EndTime -and $result.DebugInfo.StartTime) {
        $result.DebugInfo.TotalDuration = ($result.DebugInfo.EndTime - $result.DebugInfo.StartTime).TotalMilliseconds
    } else {
        $result.DebugInfo.TotalDuration = 0
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
        [string]$Branch = "exp/terraform_copilot",
        [string]$WorkspaceRoot = $null,
        [hashtable]$ManifestConfig = $null
    )
    
    # CRITICAL: Use centralized pre-installation validation (replaces scattered safety checks)
    Write-Host "Validating installation prerequisites..." -ForegroundColor Cyan
    $validation = Test-PreInstallation -AllowBootstrapOnSource:$false
    
    if (-not $validation.OverallValid) {
        Write-Host ""
        Write-Host "Pre-installation validation failed!" -ForegroundColor Red
        Write-Host ""
        
        # Show specific validation failures
        if (-not $validation.Git.Valid) {
            Write-Host "   Git Issue: $($validation.Git.Reason)" -ForegroundColor Yellow
        }
        if (-not $validation.Workspace.Valid -and -not $validation.Workspace.Skipped) {
            Write-Host "   Workspace Issue: $($validation.Workspace.Reason)" -ForegroundColor Yellow
        }
        if (-not $validation.SystemRequirements.OverallValid) {
            Write-Host "   System Issue: Missing requirements detected" -ForegroundColor Yellow
        }
        
        Write-Host ""
        Write-Host "Fix these issues and try again." -ForegroundColor Cyan
        
        return @{
            TotalFiles = 0
            Successful = 0
            Failed = 0
            Skipped = 0
            Files = @{}
            OverallSuccess = $false
            ValidationFailed = $true
            ValidationResults = $validation
            DebugInfo = @{
                StartTime = Get-Date
                Branch = $Branch
                FailureReason = "Pre-installation validation failed"
            }
        }
    }
    
    Write-Host "All prerequisites validated successfully!" -ForegroundColor Green
    Write-Host ""
    
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
    
    $results = @{
        TotalFiles = $allFiles.Count
        Successful = 0
        Failed = 0
        Skipped = 0
        Files = @{}
        OverallSuccess = $true
        DebugInfo = @{
            StartTime = Get-Date
            Branch = $Branch
            BaseUrl = $manifestConfig.BaseUrl
        }
    }
    
    Write-Host "Preparing to install $($allFiles.Count) files..." -ForegroundColor Cyan
    
    $fileIndex = 0
    foreach ($filePath in $allFiles) {
        $fileIndex++
        $downloadUrl = "$($manifestConfig.BaseUrl)/$filePath"
        
        # Debug: Show the constructed URL
        Write-Verbose "Constructed URL: $downloadUrl"
        Write-Verbose "Base URL: $($manifestConfig.BaseUrl)"
        Write-Verbose "File Path: $filePath"
        
        if (-not $downloadUrl -or $downloadUrl -eq "/$filePath" -or [string]::IsNullOrWhiteSpace($manifestConfig.BaseUrl)) {
            Write-Warning "Could not determine download URL for file: $filePath (BaseUrl: '$($manifestConfig.BaseUrl)')"
            $results.Files[$filePath] = @{
                FilePath = $filePath
                Success = $false
                Action = "Skipped"
                Message = "Could not determine download URL - BaseUrl is empty or invalid"
                Size = 0
                DebugInfo = @{
                    BaseUrl = $manifestConfig.BaseUrl
                    ConstructedUrl = $downloadUrl
                }
            }
            continue
        }
        
        $percentComplete = [math]::Round(($fileIndex / $allFiles.Count) * 100)
        # Dynamic padding to align closing brackets (1-digit=2 spaces, 2-digit=1 space, 3-digit=0 spaces)
        $progressPadding = if ($percentComplete -lt 10) { "  " } elseif ($percentComplete -lt 100) { " " } else { "" }
        $progressText = "[$percentComplete%$progressPadding]"
        Write-Host "  Downloading " -ForegroundColor Yellow -NoNewline
        Write-Host $progressText -ForegroundColor Green -NoNewline
        Write-Host ": " -ForegroundColor Yellow -NoNewline
        Write-Host $filePath -ForegroundColor White
        
        $fileResult = Install-AIFile -FilePath $filePath -DownloadUrl $downloadUrl -Force $Force -DryRun $DryRun -WorkspaceRoot $WorkspaceRoot
        $results.Files[$filePath] = $fileResult
        
        # Show error details if download failed
        if (-not $fileResult.Success) {
            Write-Host "   ERROR: $($fileResult.Message)" -ForegroundColor Red
            if ($fileResult.DebugInfo.ExceptionMessage) {
                Write-Host "   DETAILS: $($fileResult.DebugInfo.ExceptionMessage)" -ForegroundColor Red
            }
        }
        
        switch ($fileResult.Action) {
            { $_ -in @("Downloaded", "Overwritten") } { $results.Successful++ }
            "Skipped" { $results.Skipped++ }
            default { 
                $results.Failed++
                $results.OverallSuccess = $false
            }
        }
    }
    
    Write-Host "Installation completed!" -ForegroundColor Green
    
    # Show detailed debug summary
    $results.DebugInfo.EndTime = Get-Date
    if ($results.DebugInfo.StartTime -and $results.DebugInfo.EndTime) {
        $results.DebugInfo.TotalDuration = ($results.DebugInfo.EndTime - $results.DebugInfo.StartTime).TotalMilliseconds
    } else {
        $results.DebugInfo.TotalDuration = 0
    }
    
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
        
        [bool]$DryRun = $false,
        
        [string]$WorkspaceRoot = ""
    )
    
    # Resolve file path relative to workspace root if provided
    $resolvedFilePath = if ($WorkspaceRoot -and -not [System.IO.Path]::IsPathRooted($FilePath)) {
        Join-Path $WorkspaceRoot $FilePath
    } else {
        $FilePath
    }
    
    $result = @{
        FilePath = $resolvedFilePath
        Success = $false
        Action = "None"
        Message = ""
    }
    
    try {
        if (-not (Test-Path $resolvedFilePath)) {
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
        Remove-Item -Path $resolvedFilePath -Force -ErrorAction Stop
        
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
        [string]$Branch = "exp/terraform_copilot",
        [string]$WorkspaceRoot = "",
        [hashtable]$ManifestConfig = $null
    )
    
    # CRITICAL: Use centralized pre-installation validation for source repository protection
    Write-Host "Validating cleanup prerequisites..." -ForegroundColor Yellow
    $validation = Test-PreInstallation -AllowBootstrapOnSource:$false
    
    if (-not $validation.OverallValid) {
        # Build detailed error messages based on validation results
        $errorMessages = @()
        
        if ($validation.Git.Reason -like "*SAFETY VIOLATION*") {
            $errorMessages += "SAFETY VIOLATION: Cannot run clean operation on source branch '$($validation.Git.CurrentBranch)'. Switch to a feature branch to clean the AI infrastructure"
        } elseif (-not $validation.Git.Valid) {
            $errorMessages += $validation.Git.Reason
        }
        
        if (-not $validation.Workspace.Valid -and -not $validation.Workspace.Skipped) {
            $errorMessages += $validation.Workspace.Reason
        }
        
        if (-not $validation.SystemRequirements.OverallValid) {
            if (-not $validation.SystemRequirements.PowerShell.Valid) {
                $errorMessages += $validation.SystemRequirements.PowerShell.Reason
            }
            if (-not $validation.SystemRequirements.ExecutionPolicy.Valid) {
                $errorMessages += $validation.SystemRequirements.ExecutionPolicy.Reason
            }
            if (-not $validation.SystemRequirements.Commands.Valid) {
                $errorMessages += $validation.SystemRequirements.Commands.Reason
            }
            if (-not $validation.SystemRequirements.Internet.Connected) {
                $errorMessages += $validation.SystemRequirements.Internet.Reason
            }
        }
        
        # Fallback if no specific errors found
        if ($errorMessages.Count -eq 0) {
            $errorMessages += "Pre-cleanup validation failed"
        }

        Write-Host ""
        Write-Host "To Clean a Target Repository:" -ForegroundColor Cyan
        Write-Host "  1. Switch to a feature branch" -ForegroundColor White
        Write-Host "  2. Or run this command on a cloned repository" -ForegroundColor White
        Write-Host "  3. Or use -RepoDirectory to specify target directory" -ForegroundColor White
        
        return @{
            Success = $false
            Issues = $errorMessages
            FilesRemoved = 0
            DirectoriesCleaned = 0
            SourceRepoProtection = $true
            ValidationResults = $validation
        }
    }
    
    Write-Host "Cleanup validation passed!" -ForegroundColor Green
    Write-Host ""
    
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
    
    # Dynamically determine directories to check based on manifest file paths
    $directoriesToCheck = @()
    $uniqueDirectories = @{}
    
    # Directories that should NEVER be removed (important repository infrastructure)
    $protectedDirectories = @(
        ".github",           # Main GitHub directory contains workflows, issue templates, etc.
        ".vscode",           # VS Code workspace settings and configurations
        "internal",          # Core provider code
        "vendor",            # Go dependencies
        "scripts",           # Build and maintenance scripts
        "website",           # Documentation
        "examples",          # Example configurations
        "helpers"            # Helper utilities
    )
    
    foreach ($filePath in $allFiles) {
        $directory = Split-Path $filePath -Parent
        if ($directory -and -not $uniqueDirectories.ContainsKey($directory)) {
            # Skip protected directories - we only clean up specific AI subdirectories
            $isProtected = $false
            
            # Allow specific AI directories under .github
            $allowedAIDirectories = @(
                ".github/AIinstaller",
                ".github/AIinstaller/modules",
                ".github/AIinstaller/modules/powershell",
                ".github/AIinstaller/modules/bash",
                ".github/instructions", 
                ".github/prompts"
            )
            
            $isAllowedAI = $false
            foreach ($allowed in $allowedAIDirectories) {
                if ($directory -eq $allowed -or $directory.StartsWith("$allowed/")) {
                    $isAllowedAI = $true
                    break
                }
            }
            
            if (-not $isAllowedAI) {
                # Check if it's a hidden directory (starts with .) - these often contain user settings
                $dirName = Split-Path $directory -Leaf
                if ($dirName.StartsWith(".")) {
                    $isProtected = $true
                } else {
                    # Check against explicit protected directories list
                    foreach ($protected in $protectedDirectories) {
                        if ($directory -eq $protected) {
                            $isProtected = $true
                            break
                        }
                    }
                }
            }
            
            if (-not $isProtected) {
                $uniqueDirectories[$directory] = $true
                $directoriesToCheck += $directory
            }
        }
    }
    
    # Sort directories by depth (deepest first) for proper cleanup order
    # This ensures subdirectories are cleaned before parent directories
    $directoriesToCheck = $directoriesToCheck | Sort-Object { ($_ -split '[/\\]').Count } -Descending | Sort-Object
    
    # Calculate total work for accurate progress tracking
    $totalWork = $allFiles.Count + $directoriesToCheck.Count
    $workCompleted = 0
    
    # Calculate the longest filename for perfect status alignment
    $maxFileNameLength = 0
    $maxDirNameLength = 0
    
    foreach ($filePath in $allFiles) {
        $fileName = Split-Path $filePath -Leaf
        if ($fileName.Length -gt $maxFileNameLength) {
            $maxFileNameLength = $fileName.Length
        }
    }
    
    foreach ($dir in $directoriesToCheck) {
        $dirName = Split-Path $dir -Leaf
        if ($dirName.Length -gt $maxDirNameLength) {
            $maxDirNameLength = $dirName.Length
        }
    }
    
    # Use the longer of the two for universal alignment
    $maxNameLength = [math]::Max($maxFileNameLength, $maxDirNameLength)
    
    $results = @{
        TotalFiles = $allFiles.Count
        Removed = 0
        NotFound = 0
        Failed = 0
        Files = @{}
        Directories = @{}
        Success = $true
        FilesRemoved = 0
        DirectoriesCleaned = 0
        Issues = @()
    }
    
    # Pre-scan: Check if any AI files actually exist
    Write-Host "Scanning for AI infrastructure files..." -ForegroundColor Cyan
    $existingFiles = @()
    $existingDirectories = @()
    
    foreach ($filePath in $allFiles) {
        $fullPath = Join-Path $WorkspaceRoot $filePath
        if (Test-Path $fullPath) {
            $existingFiles += $filePath
        }
    }
    
    foreach ($dirPath in $directoriesToCheck) {
        $fullDirPath = Join-Path $WorkspaceRoot $dirPath
        if (Test-Path $fullDirPath) {
            # Skip child directories of AIinstaller since we'll remove it recursively
            $isChildOfAIinstaller = $dirPath.StartsWith(".github/AIinstaller/") -and $dirPath -ne ".github/AIinstaller"
            if (-not $isChildOfAIinstaller) {
                $existingDirectories += $dirPath
            }
        }
    }
    
    # If nothing exists, show clean message and exit early
    if ($existingFiles.Count -eq 0 -and $existingDirectories.Count -eq 0) {
        Write-Host ""
        Write-Host "No AI infrastructure files found to remove." -ForegroundColor Green
        Write-Host "Workspace is already clean!" -ForegroundColor Green
        
        return @{
            Success = $true
            Issues = @()
            FilesRemoved = 0
            DirectoriesCleaned = 0
            TotalFiles = $allFiles.Count
            Removed = 0
            NotFound = $allFiles.Count
            Failed = 0
            Files = @{}
            Directories = @{}
            CleanWorkspace = $true
        }
    }
    
    # Show what was found
    Write-Host "Found $($existingFiles.Count) AI files and $($existingDirectories.Count) directories to remove." -ForegroundColor Yellow
    Write-Host ""
    
    Write-Host "Removing AI Infrastructure Files" -ForegroundColor Cyan
    Write-Separator
    
    # Remove files (only process existing ones)
    $fileIndex = 0
    $totalWork = $existingFiles.Count + $existingDirectories.Count
    $workCompleted = 0
    
    # Calculate padding for clean display (only for existing files)
    $maxNameLength = 0
    foreach ($filePath in $existingFiles) {
        $fileName = Split-Path $filePath -Leaf
        if ($fileName.Length -gt $maxNameLength) {
            $maxNameLength = $fileName.Length
        }
    }
    
    # Also check directory names for padding
    foreach ($dirPath in $existingDirectories) {
        $dirName = Split-Path $dirPath -Leaf
        if ($dirName.Length -gt $maxNameLength) {
            $maxNameLength = $dirName.Length
        }
    }
    
    foreach ($filePath in $existingFiles) {
        $fileIndex++
        $workCompleted++
        $percentComplete = [math]::Round(($workCompleted / $totalWork) * 100)
        
        # Extract just the filename for cleaner display
        $fileName = Split-Path $filePath -Leaf
        
        # Calculate padding needed to align status indicators
        $fileNamePadding = " " * ($maxNameLength - $fileName.Length)
        
        # Pad "Removing File" to match "Removing Directory" length for perfect alignment
        # Dynamic padding to align closing brackets (1-digit=2 spaces, 2-digit=1 space, 3-digit=0 spaces)
        $progressPadding = if ($percentComplete -lt 10) { "  " } elseif ($percentComplete -lt 100) { " " } else { "" }
        $progressText = "[$percentComplete%$progressPadding]"
        Write-Host "  Removing File      " -ForegroundColor Cyan -NoNewline
        Write-Host $progressText -ForegroundColor Green -NoNewline
        Write-Host ": " -ForegroundColor Cyan -NoNewline
        Write-Host "$fileName$fileNamePadding " -ForegroundColor White -NoNewline
        
        $fileResult = Remove-AIFile -FilePath $filePath -DryRun $DryRun -WorkspaceRoot $WorkspaceRoot
        $results.Files[$filePath] = $fileResult
        
        switch ($fileResult.Action) {
            "Removed" { 
                $results.Removed++
                $results.FilesRemoved++
                Write-Host "[OK]" -ForegroundColor Green
            }
            "Would Remove" {
                $results.Removed++  # Count as success for dry run
                $results.FilesRemoved++
                Write-Host "[WOULD REMOVE]" -ForegroundColor Yellow
            }
            "Not Found" { 
                $results.NotFound++
                Write-Host "[NOT FOUND]" -ForegroundColor Yellow
            }
            default { 
                $results.Failed++
                $results.Success = $false
                Write-Host "[FAILED]" -ForegroundColor Red
                if ($fileResult.Message) {
                    $results.Issues += "Failed to remove ${filePath}: $($fileResult.Message)"
                }
            }
        }
    }
    
    # Remove empty directories (only process existing ones)
    $dirIndex = 0
    foreach ($dir in $existingDirectories) {
        $dirIndex++
        $workCompleted++
        $percentComplete = [math]::Round(($workCompleted / $totalWork) * 100)
        
        # Extract just the directory name for cleaner display
        $dirName = Split-Path $dir -Leaf
        
        # Calculate padding needed to align status indicators (same as files)
        $dirNamePadding = " " * ($maxNameLength - $dirName.Length)
        
        # "Removing Directory" is the longest operation name, so no padding needed
        # Dynamic padding to align closing brackets (1-digit=2 spaces, 2-digit=1 space, 3-digit=0 spaces)
        $progressPadding = if ($percentComplete -lt 10) { "  " } elseif ($percentComplete -lt 100) { " " } else { "" }
        $progressText = "[$percentComplete%$progressPadding]"
        Write-Host "  Removing Directory " -ForegroundColor Cyan -NoNewline
        Write-Host $progressText -ForegroundColor Green -NoNewline
        Write-Host ": " -ForegroundColor Cyan -NoNewline
        Write-Host "$dirName$dirNamePadding " -ForegroundColor White -NoNewline
        
        # Resolve directory path relative to workspace root if provided
        $resolvedDirPath = if ($WorkspaceRoot -and -not [System.IO.Path]::IsPathRooted($dir)) {
            Join-Path $WorkspaceRoot $dir
        } else {
            $dir
        }
        
        $dirResult = @{
            Path = $resolvedDirPath
            Action = "None"
            Success = $true
            Message = ""
        }
        
        if (Test-Path $resolvedDirPath -PathType Container) {
            # Special handling for AIinstaller directory - just nuke it recursively
            $relativePath = if ($WorkspaceRoot) {
                $resolvedDirPath.Replace($WorkspaceRoot, "").TrimStart('\', '/').Replace('\', '/')
            } else {
                $dir
            }
            
            if ($relativePath -eq ".github/AIinstaller" -or $relativePath.StartsWith(".github/AIinstaller/")) {
                # For AIinstaller, remove recursively with force
                if ($DryRun) {
                    $dirResult.Action = "Would Remove"
                    $dirResult.Message = "Directory would be removed recursively"
                    $results.DirectoriesCleaned++
                    Write-Host "[WOULD REMOVE]" -ForegroundColor Yellow
                } else {
                    try {
                        Remove-Item -Path $resolvedDirPath -Recurse -Force -ErrorAction Stop
                        $dirResult.Action = "Removed"
                        $dirResult.Message = "Directory removed recursively"
                        $results.DirectoriesCleaned++
                        Write-Host "[OK]" -ForegroundColor Green
                    }
                    catch {
                        $dirResult.Action = "Failed"
                        $dirResult.Success = $false
                        $dirResult.Message = "Failed to remove directory: $($_.Exception.Message)"
                        $results.Success = $false
                        $results.Issues += "Failed to remove directory ${resolvedDirPath}: $($_.Exception.Message)"
                        Write-Host "[FAILED]" -ForegroundColor Red
                    }
                }
            } else {
                # For other directories, check if empty first
                $dirContents = Get-ChildItem $resolvedDirPath -Force
                if ($dirContents.Count -eq 0) {
                    if ($DryRun) {
                        $dirResult.Action = "Would Remove"
                        $dirResult.Message = "Empty directory would be removed"
                        $results.DirectoriesCleaned++
                        Write-Host "[WOULD REMOVE]" -ForegroundColor Yellow
                    } else {
                        try {
                            Remove-Item -Path $resolvedDirPath -Force -ErrorAction Stop
                            $dirResult.Action = "Removed"
                            $dirResult.Message = "Empty directory removed"
                            $results.DirectoriesCleaned++
                            Write-Host "[OK]" -ForegroundColor Green
                        }
                        catch {
                            $dirResult.Action = "Failed"
                            $dirResult.Success = $false
                            $dirResult.Message = "Failed to remove directory: $($_.Exception.Message)"
                            $results.Success = $false
                            $results.Issues += "Failed to remove directory ${resolvedDirPath}: $($_.Exception.Message)"
                            Write-Host "[FAILED]" -ForegroundColor Red
                        }
                    }
                } else {
                    $dirResult.Action = "Not Empty"
                    $dirResult.Message = "Directory contains other files"
                    Write-Host "[NOT EMPTY]" -ForegroundColor Yellow
                }
            }
        } else {
            $dirResult.Action = "Not Found"
            $dirResult.Message = "Directory does not exist"
            Write-Host "[NOT FOUND]" -ForegroundColor Yellow
        }
        
        $results.Directories[$resolvedDirPath] = $dirResult
    }
    
    Write-Host ""
    Write-Host "Completed AI infrastructure removal." -ForegroundColor Green
    
    return $results
}

function Remove-EmptyParentDirectories {
    <#
    .SYNOPSIS
    Simple cleanup of empty directories after files are removed
    #>
    param(
        [string[]]$DirectoriesToCheck,
        [string]$WorkspaceRoot,
        [bool]$DryRun = $false
    )
    
    $cleanedCount = 0
    
    # Get all unique directories that might be empty now
    $allDirs = @()
    foreach ($dir in $DirectoriesToCheck) {
        $resolvedDir = if ($WorkspaceRoot -and -not [System.IO.Path]::IsPathRooted($dir)) {
            Join-Path $WorkspaceRoot $dir
        } else {
            $dir
        }
        
        # Add this directory and all its parents up to .github
        $currentDir = $resolvedDir
        while ($currentDir -and $currentDir -ne $WorkspaceRoot) {
            $relativePath = $currentDir.Replace($WorkspaceRoot, "").TrimStart('\', '/').Replace('\', '/')
            
            # Only process directories under .github that are AI-related
            if ($relativePath.StartsWith(".github/AIinstaller") -or 
                $relativePath.StartsWith(".github/instructions") -or 
                $relativePath.StartsWith(".github/prompts")) {
                $allDirs += $currentDir
            }
            
            $currentDir = Split-Path $currentDir -Parent
        }
    }
    
    # Remove duplicates and sort by depth (deepest first)
    $uniqueDirs = $allDirs | Sort-Object -Unique | Sort-Object { ($_ -split '[/\\]').Count } -Descending
    
    # Remove empty directories
    foreach ($dir in $uniqueDirs) {
        if ((Test-Path $dir -PathType Container) -and ((Get-ChildItem $dir -Force).Count -eq 0)) {
            $dirName = Split-Path $dir -Leaf
            Write-Host "  Cleaning Directory: $dirName" -NoNewline
            
            if ($DryRun) {
                Write-Host " [WOULD REMOVE]" -ForegroundColor Yellow
                $cleanedCount++
            } else {
                try {
                    Remove-Item -Path $dir -Force -ErrorAction Stop
                    Write-Host " [OK]" -ForegroundColor Green
                    $cleanedCount++
                }
                catch {
                    Write-Host " [FAILED]" -ForegroundColor Red
                }
            }
        }
    }
    
    return $cleanedCount
}

function Remove-DeprecatedFiles {
    <#
    .SYNOPSIS
    Removes files that were previously installed but are no longer in the manifest
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

function Invoke-Bootstrap {
    <#
    .SYNOPSIS
    Copy installer files to user profile for feature branch use
    #>
    
    try {
        # Show operation title (main header already displayed by caller)
        Write-Host " Bootstrap - Copying Installer to User Profile" -ForegroundColor Cyan
        Write-Separator
        
        # Create target directory
        $targetDirectory = Join-Path $env:USERPROFILE ".terraform-ai-installer"
        if (-not (Test-Path $targetDirectory)) {
            New-Item -ItemType Directory -Path $targetDirectory -Force | Out-Null
        }
        
        # Files to bootstrap from configuration
        $filesToBootstrap = $Global:InstallerConfig.Files.InstallerFiles.Files
        
        # CRITICAL: Always include the manifest file in bootstrap - it's required for user profile operations
        $manifestFile = "file-manifest.config"
        $manifestFileFullPath = ".github/AIinstaller/$manifestFile"
        if (($filesToBootstrap | Where-Object { $_ -like "*$manifestFile" }).Count -eq 0) {
            $filesToBootstrap += $manifestFileFullPath
        }
        
        # Statistics
        $statistics = @{
            "Files Copied" = 0
            "Files Downloaded" = 0
            "Files Failed" = 0
            "Total Size" = 0
        }
        
        # CRITICAL: Bootstrap should ONLY be allowed from the source branch (exp/terraform_copilot)
        # This ensures you're copying the correct, official installer files to your user profile
        # The validation in Test-PreInstallation should have already verified we're on the source branch
        $aiInstallerSourcePath = Join-Path $Global:WorkspaceRoot ".github/AIinstaller"
        
        if (Test-Path $aiInstallerSourcePath) {
            Write-Host ""
            Write-Host "Copying installer files from current repository..." -ForegroundColor Cyan
            Write-Host ""
            
            # Calculate maximum filename length for alignment
            $maxFileNameLength = 0
            foreach ($file in $filesToBootstrap) {
                $fileName = Split-Path $file -Leaf
                if ($fileName.Length -gt $maxFileNameLength) {
                    $maxFileNameLength = $fileName.Length
                }
            }
            
            # Copy files locally from source repository
            foreach ($file in $filesToBootstrap) {
                try {
                    # Handle full repository paths vs relative AIinstaller paths
                    if ($file.StartsWith('.github/AIinstaller/')) {
                        # This is a full repository path - use it directly from workspace root
                        $sourcePath = Join-Path $Global:WorkspaceRoot $file
                    } else {
                        # This is a relative path - join with AIinstaller directory
                        $sourcePath = Join-Path $aiInstallerSourcePath $file
                    }
                    
                    $fileName = Split-Path $file -Leaf
                    
                    # Determine target path based on file type and maintain directory structure
                    if ($fileName.EndsWith('.psm1')) {
                        # PowerShell modules go in modules/powershell/ subdirectory
                        $modulesDir = Join-Path $targetDirectory "modules\powershell"
                        if (-not (Test-Path $modulesDir)) {
                            New-Item -ItemType Directory -Path $modulesDir -Force | Out-Null
                        }
                        $targetPath = Join-Path $modulesDir $fileName
                    } elseif ($fileName.EndsWith('.sh')) {
                        # Bash modules and scripts go in modules/bash/ subdirectory or root for main scripts
                        if ($file -like "*modules/bash/*") {
                            $modulesDir = Join-Path $targetDirectory "modules\bash"
                            if (-not (Test-Path $modulesDir)) {
                                New-Item -ItemType Directory -Path $modulesDir -Force | Out-Null
                            }
                            $targetPath = Join-Path $modulesDir $fileName
                        } else {
                            # Main bash script goes in root directory
                            $targetPath = Join-Path $targetDirectory $fileName
                        }
                    } else {
                        # Other files (PowerShell script, config files like file-manifest.config) go directly in target directory
                        $targetPath = Join-Path $targetDirectory $fileName
                    }
                    
                    Write-Host "   Copying: " -ForegroundColor Cyan -NoNewline
                    Write-Host "$($fileName.PadRight($maxFileNameLength))" -ForegroundColor White -NoNewline
                    
                    if (Test-Path $sourcePath) {
                        Copy-Item $sourcePath $targetPath -Force
                        
                        if (Test-Path $targetPath) {
                            $fileSize = (Get-Item $targetPath).Length
                            $statistics["Files Copied"]++
                            $statistics["Total Size"] += $fileSize
                            
                            Write-Host " [OK]" -ForegroundColor "Green"
                        } else {
                            Write-Host " [FAILED]" -ForegroundColor "Red"
                            $statistics["Files Failed"]++
                        }
                    } else {
                        Write-Host " [SOURCE NOT FOUND]" -ForegroundColor "Red"
                        $statistics["Files Failed"]++
                    }
                }
                catch {
                    Write-Host " [ERROR] ($($_.Exception.Message))" -ForegroundColor "Red"
                    $statistics["Files Failed"]++
                }
            }
        } else {
            # AIinstaller directory not found in current repository
            Show-AIInstallerNotFoundError
            exit 1
        }
        
        # Prepare details for centralized summary
        $details = @()
        $totalSizeKB = [math]::Round($statistics["Total Size"] / 1KB, 1)
        
        if ($statistics["Files Copied"] -gt 0) {
            $details += "Files copied: $($statistics["Files Copied"])"
        }
        if ($statistics["Files Downloaded"] -gt 0) {
            $details += "Files downloaded: $($statistics["Files Downloaded"])"
        }
        $details += "Total size: $totalSizeKB KB"
        $details += "Location: $targetDirectory"
        
        if ($statistics["Files Failed"] -eq 0) {
            # Use centralized success reporting
            Show-OperationSummary -OperationName "Bootstrap" -Success $true -DryRun $false `
                -ItemsProcessed $statistics["Total Files"] `
                -ItemsSuccessful ($statistics["Files Copied"] + $statistics["Files Downloaded"]) `
                -ItemsFailed $statistics["Files Failed"] `
                -Details $details
            
            # Show next steps using UI module function
            Show-BootstrapNextSteps
            
            return @{
                Success = $true
                TargetDirectory = $targetDirectory
                Statistics = $statistics
            }
        } else {
            # Use centralized failure reporting
            Show-OperationSummary -OperationName "Bootstrap" -Success $false -DryRun $false `
                -ItemsProcessed $statistics["Total Files"] `
                -ItemsSuccessful ($statistics["Files Copied"] + $statistics["Files Downloaded"]) `
                -ItemsFailed $statistics["Files Failed"] `
                -Details $details
            
            return @{
                Success = $false
                Statistics = $statistics
            }
        }
    }
    catch {
        Show-OperationSummary -OperationName "Bootstrap" -Success $false -DryRun $false `
            -Details @{ "Error" = $_.Exception.Message }
        return @{
            Success = $false
            Error = $_.Exception.Message
        }
    }
}

function Invoke-CleanWorkspace {
    <#
    .SYNOPSIS
    High-level clean workspace operation with complete UI experience
    
    .PARAMETER AutoApprove
    Skip confirmation prompts
    
    .PARAMETER DryRun
    Show what would be done without making changes
    
    .PARAMETER WorkspaceRoot
    Root directory of the workspace
    #>
    param(
        [bool]$AutoApprove,
        [bool]$DryRun,
        [string]$WorkspaceRoot
    )
    
    Write-Host " Clean Workspace" -ForegroundColor Cyan
    Write-Separator
    Write-Host ""
    
    if ($DryRun) {
        Write-Host "DRY RUN - No files will be deleted" -ForegroundColor Yellow
        Write-Host ""
    }
    
    # Use the FileOperations module to properly remove all AI files
    try {
        $result = Remove-AllAIFiles -Force:$AutoApprove -DryRun:$DryRun -WorkspaceRoot $WorkspaceRoot
        
        if ($result.Success) {
            # Check if this was a clean workspace (no files found)
            if ($result.CleanWorkspace) {
                # Simple, friendly message for already clean workspace - matches UX pattern
                Write-Host ""
                Write-Host "DETAILS:" -ForegroundColor Cyan
                Write-Host "  Operation type     : " -ForegroundColor Cyan -NoNewline
                if ($DryRun) {
                    Write-Host "Dry run (simulation)" -ForegroundColor Yellow
                } else {
                    Write-Host "Live cleanup" -ForegroundColor Yellow
                }
                Write-Host "  Files removed      : " -ForegroundColor Cyan -NoNewline
                Write-Host "0" -ForegroundColor Green
                Write-Host "  Directories cleaned: " -ForegroundColor Cyan -NoNewline  
                Write-Host "0" -ForegroundColor Green
                Write-Host ""
            } else {
                # Full completion message for actual removals
                Write-Host ""
                Write-Host "All AI infrastructure files have been successfully cleaned!" -ForegroundColor Green
                Write-Separator
                Write-Host ""
                Write-Host "DETAILS:" -ForegroundColor Cyan
                
                # Find the longest key for proper alignment
                $operationType = if ($DryRun) { "Dry run (simulation)" } else { "Live cleanup" }
                $details = @{
                    "Operation type" = $operationType
                    "Files removed" = $result.FilesRemoved
                    "Directories cleaned" = $result.DirectoriesCleaned
                }
                
                $longestKey = ($details.Keys | Sort-Object Length -Descending | Select-Object -First 1)
                
                # Display each detail with consistent alignment
                foreach ($key in $details.Keys) {
                    $value = $details[$key]
                    $formattedLabel = Format-AlignedLabel -Label $key -LongestLabel $longestKey
                    
                    Write-Host "  ${formattedLabel}: " -ForegroundColor Cyan -NoNewline
                    
                    # Numbers in green, text in yellow
                    if ($value -match '^\d+$') {
                        Write-Host $value -ForegroundColor Green
                    } else {
                        Write-Host $value -ForegroundColor Yellow
                    }
                }
            }
        } else {
            Write-Host ""
            
            # Handle dry-run vs actual operation messaging differently
            if ($DryRun) {
                # For dry-run, show positive confirmation that files were verified
                $dryRunIssues = $result.Issues | Where-Object { $_ -match "Dry run - no changes made" }
                $actualIssues = $result.Issues | Where-Object { $_ -notmatch "Dry run - no changes made" }
                
                if ($dryRunIssues.Count -gt 0) {
                    Write-Host "Dry run completed successfully - all $($dryRunIssues.Count) files verified and ready for removal" -ForegroundColor Green
                    Write-Host ""
                    Write-Host "Files that would be removed:" -ForegroundColor Cyan
                    Write-Separator
                    foreach ($issue in $dryRunIssues) {
                        # Extract just the filename from the error message
                        $fileName = ($issue -split ": Dry run")[0] -replace "Failed to remove ", ""
                        Write-Host "  - $fileName" -ForegroundColor Gray
                    }
                }
                
                # Show any actual issues (non-dry-run related)
                if ($actualIssues.Count -gt 0) {
                    Write-Host "Actual Issues Encountered:" -ForegroundColor Cyan
                    Write-Separator
                    foreach ($issue in $actualIssues) {
                        Write-Host "  - $issue" -ForegroundColor Red
                    }
                    Write-Host ""
                }
            } else {
                # For actual operations, show the issues as errors
                Write-Host "Clean Operation Encountered Issues:" -ForegroundColor Cyan
                foreach ($issue in $result.Issues) {
                    Write-Host "  - $issue" -ForegroundColor Red
                }
                Write-Host ""
            }
        }
        
        return $result
    }
    catch {
        Write-Host "Failed to clean workspace: $($_.Exception.Message)" -ForegroundColor Red
        return @{ Success = $false; Issues = @($_.Exception.Message) }
    }
}

function Invoke-InstallInfrastructure {
    <#
    .SYNOPSIS
    High-level install infrastructure operation with complete UI experience
    
    .PARAMETER AutoApprove
    Skip confirmation prompts
    
    .PARAMETER DryRun
    Show what would be done without making changes
    
    .PARAMETER WorkspaceRoot
    Root directory of the workspace
    
    .PARAMETER ManifestConfig
    Manifest configuration object
    
    .PARAMETER TargetBranch
    Target repository branch name for summary display
    #>
    param(
        [bool]$AutoApprove,
        [bool]$DryRun,
        [string]$WorkspaceRoot,
        [hashtable]$ManifestConfig,
        [string]$TargetBranch = "Unknown"
    )
    
    Write-Host " Installing AI Infrastructure" -ForegroundColor Cyan
    Write-Separator
    Write-Host ""
    
    if ($DryRun) {
        Write-Host "DRY RUN - No files will be created or removed" -ForegroundColor Yellow
        Write-Host ""
    }
    
    # Step 1: Clean up deprecated files first (automatic part of installation)
    Write-Host "Checking for deprecated files..." -ForegroundColor Gray
    $deprecatedFiles = Remove-DeprecatedFiles -ManifestConfig $ManifestConfig -WorkspaceRoot $WorkspaceRoot -DryRun $DryRun -Quiet $true
    
    if ($deprecatedFiles.Count -gt 0) {
        Write-Host "  Removed $($deprecatedFiles.Count) deprecated files" -ForegroundColor Green
    } else {
        Write-Host "  No deprecated files found" -ForegroundColor Cyan
    }
    Write-Host ""
    
    # Step 2: Install/update current files
    Write-Host "Installing current AI infrastructure files..." -ForegroundColor Cyan
    
    # Use the FileOperations module to actually install files
    try {
        $result = Install-AllAIFiles -Force:$AutoApprove -DryRun:$DryRun -WorkspaceRoot $WorkspaceRoot -ManifestConfig $ManifestConfig
        
        if ($result.OverallSuccess) {
            # Use the superior completion summary function
            $nextSteps = @()
            if ($result.Skipped -gt 0) {
                $nextSteps += "Review skipped files and use -Auto-Approve if needed"
            }
            $nextSteps += "Start using GitHub Copilot with your new AI-assisted infrastructure"
            $nextSteps += "Check the .github/instructions/ folder for detailed guidelines"
            
            # Get branch information for completion summary - use target branch passed in
            $currentBranch = if ($TargetBranch -and $TargetBranch -ne "Unknown") { 
                $TargetBranch 
            } else { 
                "Unknown" 
            }
            
            $branchType = if (Test-SourceRepository) { "source" } else { 
                if ($currentBranch -eq "Unknown") { "Unknown" } else { "feature" }
            }
            
            Show-OperationSummary -OperationName "Installation" -Success $true -ItemsSuccessful $result.Successful -ItemsFailed $result.Failed -Details @{ "Branch" = $currentBranch; "Type" = $branchType }
        } else {
            Show-InstallationResults -Results $result
        }
        
        return @{ Success = $result.OverallSuccess; Details = $result }
    }
    catch {
        Write-Host "Installation failed: $($_.Exception.Message)" -ForegroundColor Red
        return @{ Success = $false; Error = $_.Exception.Message }
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
    'Invoke-Bootstrap',
    'Invoke-CleanWorkspace',
    'Invoke-InstallInfrastructure'
)

#endregion
