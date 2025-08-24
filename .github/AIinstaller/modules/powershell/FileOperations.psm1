# FileOperations Module for Terraform AzureRM Provider AI Setup
# Handles file downloading, installation, removal, and management

# Import required modules
$ModulePath = Split-Path $PSScriptRoot -Parent
Import-Module (Join-Path $ModulePath "powershell\ConfigParser.psm1") -Force

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
        
        [bool]$DryRun = $false,
        
        [string]$WorkspaceRoot = $null,
        
        [switch]$UseManifestPath
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
        # Resolve the full file path
        if ($UseManifestPath -and $WorkspaceRoot) {
            # Use manifest-driven path mapping
            $resolvedFilePath = Get-FileLocalPath -FilePath $FilePath -WorkspaceRoot $WorkspaceRoot
            $result.DebugInfo.WorkspaceRoot = $WorkspaceRoot
            $result.DebugInfo.OriginalPath = $FilePath
            $result.DebugInfo.ResolvedPath = $resolvedFilePath
            $result.DebugInfo.PathMethod = "Manifest-driven"
        } elseif ($WorkspaceRoot) {
            # Legacy path resolution
            $resolvedFilePath = Join-Path $WorkspaceRoot $FilePath
            $result.DebugInfo.WorkspaceRoot = $WorkspaceRoot
            $result.DebugInfo.OriginalPath = $FilePath
            $result.DebugInfo.ResolvedPath = $resolvedFilePath
            $result.DebugInfo.PathMethod = "Legacy"
        } else {
            $resolvedFilePath = $FilePath
            $result.DebugInfo.WorkspaceRoot = "Not provided"
            $result.DebugInfo.ResolvedPath = $resolvedFilePath
            $result.DebugInfo.PathMethod = "Direct"
        }
        
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
        $response = Invoke-WebRequest -Uri $DownloadUrl -UseBasicParsing -ErrorAction Stop
        $downloadEnd = Get-Date
        
        $result.DebugInfo.DownloadDuration = ($downloadEnd - $downloadStart).TotalMilliseconds
        $result.DebugInfo.ResponseSize = $response.Content.Length
        $result.DebugInfo.StatusCode = $response.StatusCode
        
        # Save file (handle both text and binary content properly)
        try {
            # For text files (like .md, .instructions.md), use UTF8 encoding
            if ($resolvedFilePath -match '\.(md|txt|instructions\.md|yml|yaml|json)$') {
                # Use UTF8 encoding without BOM for text files
                $utf8NoBom = New-Object System.Text.UTF8Encoding($false)
                [System.IO.File]::WriteAllText($resolvedFilePath, $response.Content, $utf8NoBom)
                $result.DebugInfo.SaveMethod = "WriteAllText (UTF8 no BOM)"
            } else {
                # For binary files, use the raw content bytes
                [System.IO.File]::WriteAllBytes($resolvedFilePath, $response.Content)
                $result.DebugInfo.SaveMethod = "WriteAllBytes"
            }
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
    $result.DebugInfo.TotalDuration = ($result.DebugInfo.EndTime - $result.DebugInfo.StartTime).TotalMilliseconds
    
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
        [string]$WorkspaceRoot = $null
    )
    
    # Use manifest-driven configuration
    $manifestConfig = Get-ManifestConfig -Branch $Branch
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
            BaseUrl = $config.BaseUrl
        }
    }
    
    Write-ProgressMessage -Activity "Installing AI Infrastructure" -Status "Preparing..." -PercentComplete 0
    
    $fileIndex = 0
    foreach ($filePath in $allFiles) {
        $fileIndex++
        $downloadUrl = Get-FileDownloadUrl -FilePath $filePath -Branch $Branch
        
        if (-not $downloadUrl) {
            Write-Warning "Could not determine download URL for file: $filePath"
            $results.Files[$filePath] = @{
                FilePath = $filePath
                Success = $false
                Action = "Skipped"
                Message = "Could not determine download URL"
                Size = 0
            }
            continue
        }
        
        $percentComplete = [math]::Round(($fileIndex / $allFiles.Count) * 100)
        Write-ProgressMessage -Activity "Installing AI Infrastructure" -Status "Processing: $filePath" -PercentComplete $percentComplete
        
        $fileResult = Install-AIFile -FilePath $filePath -DownloadUrl $downloadUrl -Force $Force -DryRun $DryRun -WorkspaceRoot $WorkspaceRoot -UseManifestPath
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
    
    Write-ProgressMessage -Activity "Installing AI Infrastructure" -Status "Completed"
    
    # Show detailed debug summary
    $results.DebugInfo.EndTime = Get-Date
    $results.DebugInfo.TotalDuration = ($results.DebugInfo.EndTime - $results.DebugInfo.StartTime).TotalMilliseconds
    
    Write-ProgressMessage -Activity "Installing AI Infrastructure" -Status "Complete" -PercentComplete 100
    
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
        [string]$WorkspaceRoot = ""
    )
    
    # Use manifest-driven configuration
    $manifestConfig = Get-ManifestConfig -Branch $Branch
    $allFiles = @()
    foreach ($section in $manifestConfig.Sections.Keys) {
        $allFiles += $manifestConfig.Sections[$section]
    }
    
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
    
    Write-ProgressMessage -Activity "Removing AI Infrastructure" -Status "Preparing..." -PercentComplete 0
    
    # Remove files
    $fileIndex = 0
    foreach ($filePath in $allFiles) {
        $fileIndex++
        $percentComplete = [math]::Round(($fileIndex / $allFiles.Count) * 50)
        Write-ProgressMessage -Activity "Removing AI Infrastructure" -Status "Removing: $filePath" -PercentComplete $percentComplete
        
        $fileResult = Remove-AIFile -FilePath $filePath -DryRun $DryRun -WorkspaceRoot $WorkspaceRoot
        $results.Files[$filePath] = $fileResult
        
        switch ($fileResult.Action) {
            "Removed" { 
                $results.Removed++
                $results.FilesRemoved++
            }
            "Not Found" { $results.NotFound++ }
            default { 
                $results.Failed++
                $results.Success = $false
                if ($fileResult.Message) {
                    $results.Issues += "Failed to remove ${filePath}: $($fileResult.Message)"
                }
            }
        }
    }
    
    # Remove empty directories
    Write-ProgressMessage -Activity "Removing AI Infrastructure" -Status "Cleaning up directories..." -PercentComplete 75
    
    $directoriesToCheck = @(
        ".github/instructions",
        ".github/prompts",
        ".vscode"
    )
    
    foreach ($dir in $directoriesToCheck) {
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
            $dirContents = Get-ChildItem $resolvedDirPath -Force
            if ($dirContents.Count -eq 0) {
                if ($DryRun) {
                    $dirResult.Action = "Would Remove"
                    $dirResult.Message = "Empty directory would be removed"
                } else {
                    try {
                        Remove-Item -Path $resolvedDirPath -Force -ErrorAction Stop
                        $dirResult.Action = "Removed"
                        $dirResult.Message = "Empty directory removed"
                        $results.DirectoriesCleaned++
                    }
                    catch {
                        $dirResult.Action = "Failed"
                        $dirResult.Success = $false
                        $dirResult.Message = "Failed to remove directory: $($_.Exception.Message)"
                        $results.Success = $false
                        $results.Issues += "Failed to remove directory ${resolvedDirPath}: $($_.Exception.Message)"
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
        
        $results.Directories[$resolvedDirPath] = $dirResult
    }
    
    Write-ProgressMessage -Activity "Removing AI Infrastructure" -Status "Completed"
    
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

function Invoke-Bootstrap {
    <#
    .SYNOPSIS
    Copy installer files to user profile for feature branch use
    
    .DESCRIPTION
    Handles bootstrapping the installer files to the user profile directory.
    Preserves source repository protections and handles both local copying
    and remote downloading based on repository type.
    
    .NOTES
    This function maintains critical source repository protections:
    - Prevents downloading on source repositories
    - Uses Test-SourceRepository for proper detection
    - Maintains workspace root handling
    #>
    
    try {
        # Modern header style to match main UI
        Write-Host $("=" * 60) -ForegroundColor Cyan
        Write-Host " Bootstrap - Copying Installer to User Profile" -ForegroundColor Cyan
        Write-Host $("=" * 60) -ForegroundColor Cyan
        
        # Create target directory
        $targetDirectory = Join-Path $env:USERPROFILE ".terraform-ai-installer"
        if (-not (Test-Path $targetDirectory)) {
            New-Item -ItemType Directory -Path $targetDirectory -Force | Out-Null
            Write-Host ""
            Write-Host "PATH: " -ForegroundColor Cyan -NoNewline
            Write-Host "$targetDirectory" -ForegroundColor Yellow
        } else {
            Write-Host ""
            Write-Host "PATH: " -ForegroundColor Cyan -NoNewline
            Write-Host "$targetDirectory" -ForegroundColor Yellow
        }
        
        # Files to bootstrap from configuration
        $filesToBootstrap = $Global:InstallerConfig.Files.InstallerFiles.Files
        
        # Statistics
        $statistics = @{
            "Files Copied" = 0
            "Files Downloaded" = 0
            "Files Failed" = 0
            "Total Size" = 0
        }
        
        # CRITICAL: Determine if we should copy locally or download from remote
        # This protection prevents overwriting source files
        $isSourceRepo = Test-SourceRepository
        $aiInstallerSourcePath = Join-Path $Global:WorkspaceRoot ".github/AIinstaller"
        
        if ($isSourceRepo -and (Test-Path $aiInstallerSourcePath)) {
            Write-Host ""
            Write-Host "Copying installer files from local source repository..." -ForegroundColor Cyan
            Write-Host ""
            
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
                        # Other files (PowerShell script, config) go directly in target directory
                        $targetPath = Join-Path $targetDirectory $fileName
                    }
                    
                    Write-Host "   Copying: " -ForegroundColor Cyan -NoNewline
                    Write-Host "$fileName" -ForegroundColor White -NoNewline
                    
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
            # CRITICAL SOURCE REPO PROTECTION: Prevent downloading on source repository
            if ($isSourceRepo) {
                Write-Host "ERROR: Cannot download files on source repository!" -ForegroundColor Red
                Write-Host "This would overwrite source files with downloaded versions, potentially losing local changes." -ForegroundColor Red
                Write-Host "Source repository detected, but local AI installer files are missing." -ForegroundColor Yellow
                Write-Host "This suggests the repository may be in an inconsistent state." -ForegroundColor Yellow
                Write-Host "" -ForegroundColor Red
                Write-Host "To resolve this issue:" -ForegroundColor Yellow
                Write-Host "  1. Check if you're on the correct branch (should contain .github/AIinstaller/)" -ForegroundColor Yellow
                Write-Host "  2. If files are missing, restore them from the main branch" -ForegroundColor Yellow
                Write-Host "  3. Use -Verify flag to check repository state without making changes" -ForegroundColor Yellow
                exit 1
            }
            
            Write-Host "  Downloading installer files from remote source branch..." -ForegroundColor "Cyan"
            Write-Host ""
            
            # Download files from remote repository
            $baseUri = "$($Global:InstallerConfig.SourceRepository)/$($Global:InstallerConfig.Branch)/.github/AIinstaller"
            
            foreach ($file in $filesToBootstrap) {
                try {
                    $uri = "$baseUri/$file"
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
                        # Other files (PowerShell script, config) go directly in target directory
                        $targetPath = Join-Path $targetDirectory $fileName
                    }
                    
                    Show-FileOperation -Operation "Downloading" -FileName $fileName -NoNewLine
                    
                    # Download with progress
                    Invoke-WebRequest -Uri $uri -OutFile $targetPath -UseBasicParsing | Out-Null
                    
                    if (Test-Path $targetPath) {
                        $fileSize = (Get-Item $targetPath).Length
                        $statistics["Files Downloaded"]++
                        $statistics["Total Size"] += $fileSize
                        
                        Write-Host " [OK]" -ForegroundColor "Green"
                    } else {
                        Write-Host " [FAILED]" -ForegroundColor "Red"
                        $statistics["Files Failed"]++
                    }
                }
                catch {
                    Write-Host " [ERROR] ($($_.Exception.Message))" -ForegroundColor "Red"
                    $statistics["Files Failed"]++
                }
            }
        }
        
        Write-Host ""
        
        if ($statistics["Files Failed"] -eq 0) {
            $totalSizeKB = [math]::Round($statistics["Total Size"] / 1KB, 1)
            
            Write-Host "Bootstrap completed successfully!" -ForegroundColor Green
            Write-Host ""
            
            if ($statistics["Files Copied"] -gt 0) {
                $label = "Files copied".PadRight(13)
                Write-Host "  ${label}: " -ForegroundColor Cyan -NoNewline
                Write-Host "$($statistics["Files Copied"])" -ForegroundColor Green
            }
            if ($statistics["Files Downloaded"] -gt 0) {
                $label = "Files downloaded".PadRight(13)
                Write-Host "  ${label}: " -ForegroundColor Cyan -NoNewline
                Write-Host "$($statistics["Files Downloaded"])" -ForegroundColor Green
            }
            
            $label = "Total size".PadRight(13)
            Write-Host "  ${label}: " -ForegroundColor Cyan -NoNewline
            Write-Host "$totalSizeKB KB" -ForegroundColor Green
            $label = "Location".PadRight(13)
            Write-Host "  ${label}: " -ForegroundColor Cyan -NoNewline
            Write-Host "$targetDirectory" -ForegroundColor Yellow
            Write-Host ""
            Write-Host "NEXT STEPS:" -ForegroundColor "Cyan"
            Write-Host ""
            Write-Host "  1. Switch to your feature branch:" -ForegroundColor "Cyan"
            Write-Host "     git checkout feature/your-branch-name" -ForegroundColor "White"
            Write-Host ""
            Write-Host "  2. Run the installer from your user profile:" -ForegroundColor "Cyan"
            Write-Host ""
            Write-Host "     Windows (PowerShell):" -ForegroundColor "Green"
            Write-Host "     & `"`$env:USERPROFILE\.terraform-ai-installer\install-copilot-setup.ps1`" -RepoDirectory `"$($Global:WorkspaceRoot)`"" -ForegroundColor "White"
            Write-Host ""
            Write-Host "     macOS/Linux (Bash):" -ForegroundColor "Green"
            Write-Host "     ~/.terraform-ai-installer/install-copilot-setup.sh -repo-directory `"$($Global:WorkspaceRoot)`"" -ForegroundColor "White"
            Write-Host ""
            Write-Host "  " -NoNewline
            Write-Host "Note:" -ForegroundColor "Cyan" -NoNewline
            Write-Host " Choose the installer for your platform. Both PowerShell and Bash components" -ForegroundColor "Yellow"
            Write-Host "        are available for maximum flexibility (e.g., VMs, containers, WSL)." -ForegroundColor "Yellow"
            Write-Host ""
            
            return @{
                Success = $true
                TargetDirectory = $targetDirectory
                Statistics = $statistics
            }
        } else {
            Write-ErrorMessage "Bootstrap failed: $($statistics["Files Failed"]) files could not be processed"
            return @{
                Success = $false
                Statistics = $statistics
            }
        }
    }
    catch {
        Write-ErrorMessage "Bootstrap failed: $($_.Exception.Message)"
        return @{
            Success = $false
            Error = $_.Exception.Message
        }
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
    'Get-FileFromGitHub',
    'Invoke-Bootstrap'
)

#endregion
