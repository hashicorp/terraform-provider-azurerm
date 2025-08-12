# Installation Module
# Handles the core installation logic using real patterns from working system

function Read-FileManifest {
    <#
    .SYNOPSIS
        Reads the file-manifest.config to get the list of files to install
    #>
    param(
        [string]$RepositoryPath
    )
    
    $manifestPath = Join-Path $RepositoryPath ".github\AILocalInstall\file-manifest.config"
    
    if (-not (Test-Path $manifestPath)) {
        Write-StatusMessage "File manifest not found: $manifestPath" "Error"
        return @{
            InstructionFiles = @()
            PromptFiles = @()
            MainFiles = @()
            Success = $false
        }
    }
    
    try {
        $content = Get-Content $manifestPath -ErrorAction Stop
        $manifest = @{
            InstructionFiles = @()
            PromptFiles = @()
            MainFiles = @()
            Success = $true
        }
        
        $currentSection = ""
        
        foreach ($line in $content) {
            $line = $line.Trim()
            
            # Skip empty lines and comments
            if (-not $line -or $line.StartsWith("#")) {
                continue
            }
            
            # Check for section headers
            if ($line -match '^\[(.+)\]$') {
                $currentSection = $matches[1]
                continue
            }
            
            # Add files to appropriate section
            switch ($currentSection) {
                "INSTRUCTION_FILES" { $manifest.InstructionFiles += $line }
                "PROMPT_FILES" { $manifest.PromptFiles += $line }
                "MAIN_FILES" { $manifest.MainFiles += $line }
            }
        }
        
        Write-StatusMessage "Read manifest: $($manifest.InstructionFiles.Count) instruction files, $($manifest.PromptFiles.Count) prompt files, $($manifest.MainFiles.Count) main files" "Info"
        return $manifest
        
    } catch {
        Write-StatusMessage "Failed to read file manifest: $($_.Exception.Message)" "Error"
        return @{
            InstructionFiles = @()
            PromptFiles = @()
            MainFiles = @()
            Success = $false
        }
    }
}

# Installation Module
# Handles the core installation logic using real patterns from working system

function Get-FileManifest {
    <#
    .SYNOPSIS
        Reads the file-manifest.config to get the list of files to install
    #>
    param(
        [string]$RepositoryPath
    )
    
    $manifestPath = Join-Path $RepositoryPath ".github\AILocalInstall\file-manifest.config"
    
    if (-not (Test-Path $manifestPath)) {
        Write-StatusMessage "File manifest not found: $manifestPath" "Error"
        return $null
    }
    
    try {
        $content = Get-Content $manifestPath -ErrorAction Stop
        $manifest = @{
            InstructionFiles = @()
            PromptFiles = @()
            MainFiles = @()
        }
        
        $currentSection = $null
        
        foreach ($line in $content) {
            $line = $line.Trim()
            
            # Skip empty lines and comments
            if ([string]::IsNullOrEmpty($line) -or $line.StartsWith('#')) {
                continue
            }
            
            # Check for section headers
            if ($line -match '^\[(.+)\]$') {
                $currentSection = $matches[1]
                continue
            }
            
            # Add files to appropriate section
            switch ($currentSection) {
                'INSTRUCTION_FILES' { $manifest.InstructionFiles += $line }
                'PROMPT_FILES' { $manifest.PromptFiles += $line }
                'MAIN_FILES' { $manifest.MainFiles += $line }
            }
        }
        
        $instructionText = if ($manifest.InstructionFiles.Count -eq 1) { "instruction file" } else { "instruction files" }
        $promptText = if ($manifest.PromptFiles.Count -eq 1) { "prompt file" } else { "prompt files" }
        $mainText = if ($manifest.MainFiles.Count -eq 1) { "main file" } else { "main files" }
        
        Write-StatusMessage "Loaded manifest: $($manifest.InstructionFiles.Count) $instructionText, $($manifest.PromptFiles.Count) $promptText, $($manifest.MainFiles.Count) $mainText" "Info"
        return $manifest
        
    } catch {
        Write-StatusMessage "Failed to read file manifest: $($_.Exception.Message)" "Error"
        return $null
    }
}

function Test-AIInstallation {
    <#
    .SYNOPSIS
        Verifies that all AI files from manifest are installed in user's VS Code directory
    #>
    param(
        [string]$RepositoryPath
    )
    
    Write-StatusMessage "Verifying AI installation in VS Code..." "Info"
    
    # 1. Read the manifest file
    $manifest = Get-FileManifest -RepositoryPath $RepositoryPath
    if ($null -eq $manifest) {
        Write-StatusMessage "Cannot verify: manifest file not found" "Error"
        return @{
            Success = $false
            Errors = @("Manifest file not found")
        }
    }
    
    $errors = @()
    $userCodePath = Join-Path $env:APPDATA "Code\User"
    
    # Track what's actually found vs expected
    $foundInstructionFiles = @()
    $foundPromptFiles = @()
    $foundMainFiles = @()
    $settingsConfigured = $false
    
    # 2. Check instruction files exist in user's VS Code
    $instructionsPath = Join-Path $userCodePath "instructions\terraform-azurerm"
    foreach ($file in $manifest.InstructionFiles) {
        $filePath = Join-Path $instructionsPath $file
        if (Test-Path $filePath) {
            $foundInstructionFiles += $file
        } else {
            $errors += "Missing instruction file: $file"
            Write-StatusMessage "Missing instruction file: $file" "Error"
        }
    }
    
    # 3. Check prompt files exist in user's VS Code (should be errors if missing after install)
    $promptsPath = Join-Path $userCodePath "prompts"
    foreach ($file in $manifest.PromptFiles) {
        $filePath = Join-Path $promptsPath $file
        if (Test-Path $filePath) {
            $foundPromptFiles += $file
        } else {
            $errors += "Missing prompt file: $file"
            Write-StatusMessage "Missing prompt file: $file" "Error"
        }
    }
    
    # 4. Check main files exist in user's VS Code
    foreach ($file in $manifest.MainFiles) {
        $filePath = Join-Path $userCodePath $file
        if (Test-Path $filePath) {
            $foundMainFiles += $file
        } else {
            $errors += "Missing main file: $file"
            Write-StatusMessage "Missing main file: $file" "Error"
        }
    }
    
    # 5. Check settings.json points to local instruction files
    $settingsPath = Join-Path $userCodePath "settings.json"
    if (Test-Path $settingsPath) {
        try {
            $settingsContent = Get-Content $settingsPath -Raw | ConvertFrom-Json
            
            # Convert to hashtable if needed for PowerShell 5.1 compatibility
            if ($settingsContent -is [PSCustomObject]) {
                $settings = @{}
                $settingsContent.PSObject.Properties | ForEach-Object {
                    $settings[$_.Name] = $_.Value
                }
            } else {
                $settings = $settingsContent
            }
            
            # Check for the specific github.copilot.chat.reviewSelection.instructions setting
            $settingKey = "github.copilot.chat.reviewSelection.instructions"
            if ($settings.ContainsKey($settingKey)) {
                $instructions = $settings[$settingKey]
                $hasCorrectPaths = $false  # Start with false, set to true if ANY correct path found
                $hasIncorrectPaths = $false
                
                # Check each instruction file path
                foreach ($instruction in $instructions) {
                    if ($instruction -is [PSCustomObject] -and $instruction.file) {
                        $filePath = $instruction.file
                        # Check if path points to repository vs local installation
                        if ($filePath -like "*/.github/instructions/*" -or $filePath -like "*\.github\instructions\*") {
                            $hasIncorrectPaths = $true
                            Write-StatusMessage "Found repository path in settings.json: $filePath (will break local AI installation)" "Error"
                        } elseif ($filePath -like "*/instructions/terraform-azurerm/*" -or $filePath -like "*\instructions\terraform-azurerm\*" -or 
                                  $filePath -like "instructions/terraform-azurerm/*" -or $filePath -like "instructions\terraform-azurerm\*") {
                            $hasCorrectPaths = $true  # Local instruction files directory
                            Write-StatusMessage "Found correct local path: $filePath" "Success"
                        } elseif ($filePath -eq "copilot-instructions.md") {
                            $hasCorrectPaths = $true  # Local root file (correct)
                            Write-StatusMessage "Found correct local path: $filePath" "Success"
                        } elseif ($filePath -like "*/copilot-instructions.md") {
                            $hasIncorrectPaths = $true
                            Write-StatusMessage "Found repository path in settings.json: $filePath (will break local AI installation - should be just 'copilot-instructions.md')" "Error"
                        } else {
                            Write-StatusMessage "Found unexpected path format: $filePath" "Warning"
                        }
                    }
                }
                
                # Settings are considered correct if we have correct local paths AND no incorrect repository paths
                if ($hasCorrectPaths -and -not $hasIncorrectPaths) {
                    Write-StatusMessage "VS Code settings.json correctly configured for local AI installation" "Success"
                    $settingsConfigured = $true  # Mark as successfully configured
                } elseif ($hasCorrectPaths -and $hasIncorrectPaths) {
                    $errors += "VS Code settings.json: contains repository paths that will break local AI installation"
                    Write-StatusMessage "Settings contain both local and repository paths - repository paths will cause failures" "Error"
                } else {
                    $errors += "VS Code settings.json: contains no local instruction paths (should point to instructions/terraform-azurerm/)"
                }
            } else {
                $errors += "VS Code settings.json: missing github.copilot.chat.reviewSelection.instructions setting"
                Write-StatusMessage "VS Code settings.json: missing github.copilot.chat.reviewSelection.instructions setting" "Error"
            }
            
        } catch {
            $errors += "VS Code settings.json: failed to read - $($_.Exception.Message)"
            Write-StatusMessage "Failed to read VS Code settings.json: $($_.Exception.Message)" "Error"
        }
    } else {
        $errors += "VS Code settings.json: not found"
        Write-StatusMessage "VS Code settings.json not found" "Error"
    }
    
    # 6. Report results
    if ($errors.Count -eq 0) {
        Write-StatusMessage "AI installation verification: SUCCESS" "Success"
        Write-StatusMessage "All $($manifest.InstructionFiles.Count + $manifest.PromptFiles.Count + $manifest.MainFiles.Count) files found in VS Code" "Success"
        return @{
            Success = $true
            Errors = @()
            InstructionFiles = $foundInstructionFiles  # Actually found files
            PromptFiles = $foundPromptFiles            # Actually found files
            MainFiles = $foundMainFiles                # Actually found files
            SettingsConfigured = $settingsConfigured   # Whether settings.json is correctly configured
            ExpectedInstructionFiles = $manifest.InstructionFiles  # For detailed reporting if needed
            ExpectedPromptFiles = $manifest.PromptFiles
            ExpectedMainFiles = $manifest.MainFiles
        }
    } else {
        Write-StatusMessage "AI installation verification: FAILED ($($errors.Count) issues)" "Error"
        return @{
            Success = $false
            Errors = $errors
            InstructionFiles = $foundInstructionFiles  # Actually found files (partial)
            PromptFiles = $foundPromptFiles            # Actually found files (partial)
            MainFiles = $foundMainFiles                # Actually found files (partial)
            SettingsConfigured = $settingsConfigured   # Whether settings.json is correctly configured
            ExpectedInstructionFiles = $manifest.InstructionFiles  # Expected files for comparison
            ExpectedPromptFiles = $manifest.PromptFiles
            ExpectedMainFiles = $manifest.MainFiles
        }
    }
}

function Install-InstructionFiles {
    <#
    .SYNOPSIS
        Copies instruction files from repository to user's VS Code instructions directory
    #>
    param(
        [string]$RepositoryPath,
        [array]$FileList
    )
    
    $sourcePath = Join-Path $RepositoryPath ".github\instructions"
    $targetPath = Join-Path $env:APPDATA "Code\User\instructions\terraform-azurerm"
    
    Write-StatusMessage "Installing instruction files to VS Code..." "Info"
    
    # Validate source directory exists in repository
    if (-not (Test-Path $sourcePath)) {
        Write-StatusMessage "Instructions directory not found in repository: $sourcePath" "Error"
        return @{
            Success = $false
            Installed = @()
            Errors = @("Instructions directory not found in repository: $sourcePath")
        }
    }
    
    # Validate specific files from manifest exist in repository
    $missingFiles = @()
    foreach ($fileName in $FileList) {
        $sourceFile = Join-Path $sourcePath $fileName
        if (-not (Test-Path $sourceFile)) {
            $missingFiles += $fileName
        }
    }
    
    if ($missingFiles.Count -gt 0) {
        Write-StatusMessage "Missing instruction files in repository: $($missingFiles -join ', ')" "Error"
        return @{
            Success = $false
            Installed = @()
            Errors = @("Missing instruction files: $($missingFiles -join ', ')")
        }
    }
    
    # Create target directory in user's VS Code
    if (-not (Test-Path $targetPath)) {
        try {
            New-Item -ItemType Directory -Path $targetPath -Force | Out-Null
            Write-StatusMessage "Created instructions directory: $targetPath" "Info"
        } catch {
            Write-StatusMessage "Failed to create target directory: $targetPath" "Error"
            return @{
                Success = $false
                Installed = @()
                Errors = @("Failed to create target directory: $($_.Exception.Message)")
            }
        }
    }
    
    # Copy specific instruction files from manifest
    $copiedFiles = @()
    $errors = @()
    
    foreach ($fileName in $FileList) {
        try {
            $sourceFile = Join-Path $sourcePath $fileName
            $targetFile = Join-Path $targetPath $fileName
            Copy-Item -Path $sourceFile -Destination $targetFile -Force -ErrorAction Stop
            $copiedFiles += $fileName
            Write-StatusMessage "Copied instruction file: $fileName" "Success"
        } catch {
            $errors += "Failed to copy $fileName`: $($_.Exception.Message)"
            Write-StatusMessage "Failed to copy $fileName`: $($_.Exception.Message)" "Error"
        }
    }
    
    Write-StatusMessage "Copied $($copiedFiles.Count) instruction files to VS Code" "Success"
    
    return @{
        Success = ($errors.Count -eq 0)
        Installed = $copiedFiles
        Errors = $errors
        Count = $copiedFiles.Count
        TargetPath = $targetPath
    }
}

function Install-PromptFiles {
    <#
    .SYNOPSIS
        Copies prompt files from repository to user's VS Code prompts directory
    #>
    param(
        [string]$RepositoryPath,
        [array]$FileList
    )
    
    $sourcePath = Join-Path $RepositoryPath ".github\prompts"
    $targetPath = Join-Path $env:APPDATA "Code\User\prompts"
    
    Write-StatusMessage "Installing prompt files to VS Code..." "Info"
    
    # Check if prompts exist in repository (optional, but if FileList has items, they should exist)
    if ($FileList.Count -eq 0) {
        Write-StatusMessage "No prompt files specified in manifest (optional)" "Info"
        return @{
            Success = $true
            Installed = @()
            Errors = @()
            Count = 0
            TargetPath = $null
        }
    }
    
    if (-not (Test-Path $sourcePath)) {
        Write-StatusMessage "Prompts directory not found in repository: $sourcePath" "Error"
        return @{
            Success = $false
            Installed = @()
            Errors = @("Prompts directory not found in repository: $sourcePath")
        }
    }
    
    # Validate specific files from manifest exist in repository
    $missingFiles = @()
    foreach ($fileName in $FileList) {
        $sourceFile = Join-Path $sourcePath $fileName
        if (-not (Test-Path $sourceFile)) {
            $missingFiles += $fileName
        }
    }
    
    if ($missingFiles.Count -gt 0) {
        Write-StatusMessage "Missing prompt files in repository: $($missingFiles -join ', ')" "Error"
        return @{
            Success = $false
            Installed = @()
            Errors = @("Missing prompt files: $($missingFiles -join ', ')")
        }
    }
    
    # Create target directory in user's VS Code
    if (-not (Test-Path $targetPath)) {
        try {
            New-Item -ItemType Directory -Path $targetPath -Force | Out-Null
            Write-StatusMessage "Created prompts directory: $targetPath" "Info"
        } catch {
            Write-StatusMessage "Failed to create prompts directory: $targetPath" "Error"
            return @{
                Success = $false
                Installed = @()
                Errors = @("Failed to create prompts directory: $($_.Exception.Message)")
            }
        }
    }
    
    # Copy specific prompt files from manifest
    $copiedFiles = @()
    $errors = @()
    
    foreach ($fileName in $FileList) {
        try {
            $sourceFile = Join-Path $sourcePath $fileName
            $targetFile = Join-Path $targetPath $fileName
            Copy-Item -Path $sourceFile -Destination $targetFile -Force -ErrorAction Stop
            $copiedFiles += $fileName
            Write-StatusMessage "Copied prompt file: $fileName" "Success"
        } catch {
            $errors += "Failed to copy $fileName`: $($_.Exception.Message)"
            Write-StatusMessage "Failed to copy $fileName`: $($_.Exception.Message)" "Error"
        }
    }
    
    Write-StatusMessage "Copied $($copiedFiles.Count) prompt files to VS Code" "Success"
    
    return @{
        Success = ($errors.Count -eq 0)
        Installed = $copiedFiles
        Errors = $errors
        Count = $copiedFiles.Count
        TargetPath = $targetPath
    }
}

function Install-MainFiles {
    <#
    .SYNOPSIS
        Copies and modifies main files from repository to user's VS Code root directory
    #>
    param(
        [string]$RepositoryPath,
        [array]$FileList
    )
    
    $sourceDir = Join-Path $RepositoryPath ".github"
    $targetDir = Join-Path $env:APPDATA "Code\User"
    $installed = @()
    $errors = @()
    
    Write-StatusMessage "Installing main files to VS Code..." "Info"
    
    # Process each file specified in manifest
    foreach ($fileName in $FileList) {
        $sourcePath = Join-Path $sourceDir $fileName
        $targetPath = Join-Path $targetDir $fileName
        
        # Validate source file exists in repository
        if (-not (Test-Path $sourcePath)) {
            Write-StatusMessage "Main file not found in repository: $sourcePath" "Error"
            $errors += "Main file not found in repository: $fileName"
            continue
        }
        
        # Verify source file integrity
        if (-not (Test-FileIntegrity -FilePath $sourcePath)) {
            Write-StatusMessage "Main file in repository may be corrupted: $fileName" "Error"
            $errors += "Main file may be corrupted: $fileName"
            continue
        }
        
        try {
            # Read the source file content
            $content = Get-Content $sourcePath -Raw -ErrorAction Stop
            
            # For copilot-instructions.md, modify paths to point to local VS Code directories
            if ($fileName -eq "copilot-instructions.md") {
                $localInstructionsPath = Join-Path $env:APPDATA "Code\User\instructions\terraform-azurerm"
                
                # Update paths in the content to point to local VS Code directories
                $modifiedContent = $content -replace '\.github\\instructions', $localInstructionsPath.Replace('\', '\\')
                $modifiedContent = $modifiedContent -replace '\.github/instructions', $localInstructionsPath.Replace('\', '/')
                
                if ($modifiedContent -ne $content) {
                    Write-StatusMessage "Modified instruction file paths to point to local VS Code directories" "Info"
                }
                
                $content = $modifiedContent
            }
            
            # Write the (possibly modified) content to user's VS Code directory
            Set-Content -Path $targetPath -Value $content -Encoding UTF8 -ErrorAction Stop
            
            Write-StatusMessage "Installed main file: $fileName" "Success"
            $installed += $fileName
            
        } catch {
            $errors += "Failed to copy and modify $fileName`: $($_.Exception.Message)"
            Write-StatusMessage "Failed to copy and modify $fileName`: $($_.Exception.Message)" "Error"
        }
    }
    
    return @{
        Success = ($errors.Count -eq 0)
        Installed = $installed
        Errors = $errors
        Count = $installed.Count
        TargetPath = $targetDir
    }
}

function Update-VSCodeSettings {
    <#
    .SYNOPSIS
        Updates VS Code settings.json with Terraform AzureRM configuration (real implementation)
    #>
    param(
        [string]$RepositoryPath,
        [switch]$Force
    )
    
    # Real implementation pattern from working system
    $userDataPath = Join-Path $env:APPDATA "Code\User"
    $settingsPath = Join-Path $userDataPath "settings.json"
    
    Write-StatusMessage "Configuring VS Code settings..." "Info"
    
    # Check if settings already contain our installation markers - if so, don't create redundant backup
    $isReinstall = $false
    if (Test-Path $settingsPath) {
        try {
            $existingContent = Get-Content $settingsPath -Raw
            if ($existingContent -like "*AZURERM_BACKUP_LENGTH*" -or $existingContent -like "*AZURERM_INSTALLATION_DATE*") {
                $isReinstall = $true
                Write-StatusMessage "Detected existing AI installation - skipping backup creation" "Info"
            }
        } catch {
            Write-StatusMessage "Could not read existing settings to check for previous installation" "Warning"
        }
    }
    
    # Create backup if settings exist and this is NOT a reinstall
    if ((Test-Path $settingsPath) -and -not $Force -and -not $isReinstall) {
        try {
            $backupResult = New-SafeBackup -SourcePath $settingsPath -BackupReason "AI installation"
            if (-not $backupResult) {
                Write-StatusMessage "Could not create backup - proceeding with manual merge scenario" "Warning"
                # Don't return false here - continue to manual merge handling
            }
        } catch {
            Write-StatusMessage "Backup failed - proceeding with manual merge scenario" "Warning"
            # Don't return false here - continue to manual merge handling
        }
    }
    
    # Read existing settings or create new (handle JSONC with comments)
    $settings = @{}
    $manualMergeRequired = $false
    
    if (Test-Path $settingsPath) {
        try {
            $existingContent = Get-Content $settingsPath -Raw
            if ($existingContent.Trim()) {
                # Remove comments for JSON parsing (but preserve them for final output)
                $cleanJson = $existingContent -replace '//.*$', '' | ForEach-Object { $_ -replace '\s+$', '' }
                $cleanJson = ($cleanJson -split "`n" | Where-Object { $_.Trim() -ne '' }) -join "`n"
                
                if ($cleanJson.Trim()) {
                    $tempSettings = $cleanJson | ConvertFrom-Json
                    
                    # Convert to hashtable for PowerShell 5.1 compatibility (real pattern)
                    if ($tempSettings -is [PSCustomObject]) {
                        $settings = @{}
                        $tempSettings.PSObject.Properties | ForEach-Object {
                            $settings[$_.Name] = $_.Value
                        }
                    } else {
                        $settings = $tempSettings
                    }
                }
            }
        } catch {
            # Manual merge scenario - settings.json exists but is unreadable/corrupted
            $manualMergeRequired = $true
            
            Write-StatusMessage "MANUAL MERGE REQUIRED: Cannot read existing settings.json" "Warning"
            Write-StatusMessage "File exists but contains invalid JSON or is corrupted" "Warning"
            Write-StatusMessage "Creating minimal backup marker and proceeding with file-only installation" "Info"
            
            # Create minimal backup with just metadata to mark this as manual merge scenario
            $backupDir = Join-Path $env:APPDATA "Code\User\.terraform-azurerm-backups"
            if (-not (Test-Path $backupDir)) {
                New-Item -ItemType Directory -Path $backupDir -Force | Out-Null
            }
            
            $timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
            $backupPath = Join-Path $backupDir "settings_backup_$timestamp.json"
            
            # Create minimal backup file with only manual merge metadata
            $manualMergeBackup = @{
                "// AZURERM_BACKUP_LENGTH" = -1
                "// AZURERM_INSTALLATION_DATE" = (Get-Date -Format "yyyy-MM-dd HH:mm:ss")
                "// MANUAL_MERGE_NOTICE" = "This was a manual merge scenario - user must manually add AI settings to their settings.json"
            }
            
            try {
                $manualMergeBackup | ConvertTo-Json -Depth 10 | Set-Content -Path $backupPath -Encoding UTF8
                Write-StatusMessage "Created manual merge marker: $($backupPath | Split-Path -Leaf)" "Info"
            } catch {
                Write-StatusMessage "Failed to create manual merge marker, but continuing" "Warning"
            }
            
            # DO NOT modify settings.json - user must manually merge
            Write-StatusMessage "INSTALLATION INCOMPLETE: Settings.json requires manual configuration" "Warning"
            Write-StatusMessage "Instruction and prompt files have been installed successfully" "Success"
            Write-StatusMessage "" "Info"
            Write-StatusMessage "=== MANUAL MERGE INSTRUCTIONS ===" "Warning"
            Write-StatusMessage "Your VS Code settings.json file contains invalid JSON and could not be automatically updated." "Info"
            Write-StatusMessage "To complete the AI setup, you must manually add the following settings to your settings.json:" "Info"
            Write-StatusMessage "" "Info"
            Write-StatusMessage "1. Fix any JSON syntax errors in your current settings.json file" "Info"
            Write-StatusMessage "2. Add the AI settings from this repository's .vscode/settings.json file" "Info"
            Write-StatusMessage "3. The required settings include GitHub Copilot configuration and file associations" "Info"
            Write-StatusMessage "" "Info"
            Write-StatusMessage "Repository settings file location: $RepositoryPath\.vscode\settings.json" "Info"
            Write-StatusMessage "Your settings file location: $settingsPath" "Info"
            Write-StatusMessage "" "Info"
            Write-StatusMessage "After manually merging settings, the AI features will be fully functional." "Success"
            Write-StatusMessage "================================" "Warning"
        }
    }
    
    # Exit early if manual merge is required - DO NOT TOUCH settings.json
    if ($manualMergeRequired) {
        return $false
    }
    
    # Add AI system settings - exact format from repository's .vscode/settings.json
    $terraformSettings = @{
        # Commit message generation with Azure provider context
        "github.copilot.chat.commitMessageGeneration.instructions" = @(
            @{
                "text" = "Provide a concise and clear commit message that summarizes the changes made in the code. For complex changes, include the following details: 1) Specify if the change introduces a breaking change and describe its impact. 2) Highlight any new resources or features added. 3) Mention updates to Azure services or APIs. Aim to keep the message under 72 characters per line for readability."
            }
        )

        # Disable conversation history for privacy
        "github.copilot.chat.summarizeAgentConversationHistory.enabled" = $false

        # Enable code review with instruction files (CORE 6 FILES ONLY - matches repository .vscode/settings.json)
        "github.copilot.chat.reviewSelection.enabled" = $true
        "github.copilot.chat.reviewSelection.instructions" = @(
            @{"file" = "copilot-instructions.md"}
            @{"file" = "instructions/terraform-azurerm/implementation-guide.instructions.md"}
            @{"file" = "instructions/terraform-azurerm/azure-patterns.instructions.md"}
            @{"file" = "instructions/terraform-azurerm/testing-guidelines.instructions.md"}
            @{"file" = "instructions/terraform-azurerm/documentation-guidelines.instructions.md"}
            @{"file" = "instructions/terraform-azurerm/provider-guidelines.instructions.md"}
        )

        # File associations for proper syntax highlighting
        "files.associations" = @{
            "*.instructions.md" = "markdown"
            ".github/*.md" = "markdown"
        }

        # Additional Copilot optimization settings
        "github.copilot.advanced" = @{
            "length" = 3000
            "temperature" = 0.1
        }

        # Enable Copilot across all relevant contexts and file types
        "github.copilot.enable" = @{
            "*" = $true
            "terminal" = $true
        }
    }
    
    # Remove any old invalid properties from previous installations
    $invalidKeys = @(
        "terraform_azurerm_provider_mode",
        "terraform_azurerm_ai_enhanced", 
        "terraform_azurerm_installation_date",
        "terraform_azurerm_backup_length",
        "github.copilot.chat.localeOverride",
        "// AZURERM_BACKUP_LENGTH",
        "// AZURERM_INSTALLATION_DATE"
    )
    
    foreach ($invalidKey in $invalidKeys) {
        if ($settings.ContainsKey($invalidKey)) {
            $settings.Remove($invalidKey)
            Write-StatusMessage "Removed invalid setting: $invalidKey" "Info"
        }
    }

    # Merge settings (real pattern)
    foreach ($key in $terraformSettings.Keys) {
        if ($key -notlike "// AZURERM_*") {  # Skip metadata, we'll add it separately
            $settings[$key] = $terraformSettings[$key]
        }
    }
    
    try {
        # Create ordered settings with metadata at the top
        $orderedSettings = [ordered]@{}
        
        # Add metadata first
        $currentDate = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
        $backupLength = 0  # Default for installation
        $orderedSettings["// AZURERM_BACKUP_LENGTH"] = $backupLength
        $orderedSettings["// AZURERM_INSTALLATION_DATE"] = $currentDate
        
        # Add all other settings (sorted alphabetically for consistency)
        $sortedKeys = $settings.Keys | Sort-Object
        foreach ($key in $sortedKeys) {
            $orderedSettings[$key] = $settings[$key]
        }
        
        # Write ordered settings as JSON
        $json = $orderedSettings | ConvertTo-Json -Depth 10
        Set-Content -Path $settingsPath -Value $json -Encoding UTF8
        
        Write-StatusMessage "VS Code settings updated successfully" "Success"
        return $true
    } catch {
        Write-StatusMessage "Failed to update VS Code settings: $_" "Error"
        return $false
    }
}

function Start-CompleteInstallation {
    <#
    .SYNOPSIS
        Orchestrates the complete installation process using file manifest
    #>
    param(
        [string]$RepositoryPath,
        [switch]$Force
    )
    
    Write-StatusMessage "Starting Terraform AzureRM Provider AI Setup installation..." "Info"
    
    $results = @{
        Success = $false
        InstructionFiles = @{}
        PromptFiles = @{}
        MainFiles = @{}
        VSCodeSettings = $false
        Errors = @()
    }
    
    try {
        # Read file manifest to get exact list of files to install
        Write-StatusMessage "Reading file manifest..." "Info"
        $manifest = Read-FileManifest -RepositoryPath $RepositoryPath
        
        if (-not $manifest.Success) {
            $results.Errors += "Failed to read file manifest"
            Write-StatusMessage "Failed to read file manifest - cannot proceed with installation" "Error"
            return $results
        }
        
        # Install files based on manifest
        $results.InstructionFiles = Install-InstructionFiles -RepositoryPath $RepositoryPath -FileList $manifest.InstructionFiles
        $results.PromptFiles = Install-PromptFiles -RepositoryPath $RepositoryPath -FileList $manifest.PromptFiles
        $results.MainFiles = Install-MainFiles -RepositoryPath $RepositoryPath -FileList $manifest.MainFiles
        
        # Install VS Code settings (real pattern - only external modification)
        $results.VSCodeSettings = Update-VSCodeSettings -RepositoryPath $RepositoryPath -Force:$Force
        
        # Check overall success - modified to handle manual merge scenarios
        $fileValidationSuccess = $results.InstructionFiles.Success -and $results.PromptFiles.Success -and $results.MainFiles.Success
        $results.Success = $fileValidationSuccess -and $results.VSCodeSettings
        
        # Transform results for UI compatibility
        $results.ExpectedInstructionFiles = @{ Count = $manifest.InstructionFiles.Count }
        $results.ExpectedPromptFiles = @{ Count = $manifest.PromptFiles.Count }
        $results.ExpectedMainFiles = @{ Count = $manifest.MainFiles.Count }
        $results.SettingsConfigured = $results.VSCodeSettings
        
        if ($results.Success) {
            Write-StatusMessage "Installation completed successfully!" "Success"
            
            # Create installation marker (real pattern from working system)
            $backupDir = Join-Path $env:APPDATA "Code\User\.terraform-azurerm-backups"
            if (-not (Test-Path $backupDir)) {
                New-Item -ItemType Directory -Path $backupDir -Force | Out-Null
            }
        } elseif ($fileValidationSuccess -and -not $results.VSCodeSettings) {
            # Special case: Files installed successfully but settings configuration failed (manual merge scenario)
            Write-StatusMessage "Installation partially completed - files installed successfully" "Warning"
            Write-StatusMessage "MANUAL ACTION REQUIRED: You must manually configure VS Code settings.json" "Warning"
            Write-StatusMessage "Copy the AI settings from the repository's .vscode/settings.json to your settings.json" "Info"
            Write-StatusMessage "File installation was successful, only settings configuration requires manual intervention" "Info"
            
            # Add the settings failure to errors so it gets counted properly
            $results.Errors += "VS Code settings configuration failed - manual merge required"
        } else {
            Write-StatusMessage "Installation completed with issues" "Warning"
            if (-not $results.InstructionFiles.Success) { $results.Errors += $results.InstructionFiles.Errors }
            if (-not $results.PromptFiles.Success) { $results.Errors += $results.PromptFiles.Errors }
            if (-not $results.MainFiles.Success) { $results.Errors += $results.MainFiles.Errors }
            if (-not $results.VSCodeSettings) { $results.Errors += "VS Code settings configuration failed - see manual merge instructions above" }
        }
        
        return $results
        
    } catch {
        $results.Errors += "Installation failed: $($_.Exception.Message)"
        Write-StatusMessage "Installation failed: $($_.Exception.Message)" "Error"
        throw
    }
}

function Remove-AIInstallation {
    <#
    .SYNOPSIS
        Removes AI installation and restores backups - complete implementation
    #>
    param([string]$RepositoryPath)
    
    Write-StatusMessage "Starting AI installation removal..." "Info"
    
    $results = @{
        Success = $true
        Errors = @()
    }
    
    try {
        # Get proper installation paths (files copied to user's VS Code directories)
        $paths = @{
            VSCodeSettings = Join-Path $env:APPDATA "Code\User\settings.json"
            VSCodeBackupDir = Join-Path $env:APPDATA "Code\User\.terraform-azurerm-backups"
            CopilotInstructions = Join-Path $env:APPDATA "Code\User\copilot-instructions.md"
            InstructionsDir = Join-Path $env:APPDATA "Code\User\instructions\terraform-azurerm"
            PromptsDir = Join-Path $env:APPDATA "Code\User\prompts\terraform-azurerm"
        }
        
        # Smart VS Code settings restoration based on metadata
        if (Test-Path $paths.VSCodeSettings) {
            # Check if we created this file from scratch by reading our metadata
            $shouldDelete = $false
            $isManualMerge = $false
            try {
                $settingsContent = Get-Content $paths.VSCodeSettings -Raw
                $settings = $settingsContent | ConvertFrom-Json
                
                # Convert to hashtable if needed
                if ($settings -is [PSCustomObject]) {
                    $hashtable = @{}
                    $settings.PSObject.Properties | ForEach-Object {
                        $hashtable[$_.Name] = $_.Value
                    }
                    $settings = $hashtable
                }
                
                # Check our metadata to determine restoration strategy
                if ($settings.ContainsKey("// AZURERM_BACKUP_LENGTH")) {
                    $backupLength = $settings["// AZURERM_BACKUP_LENGTH"]
                    
                    if ($backupLength -eq 0) {
                        # We created this file from scratch - just delete it
                        $shouldDelete = $true
                        Write-StatusMessage "No original settings.json existed - will delete our created file" "Info"
                    } elseif ($backupLength -eq -1) {
                        # Manual merge scenario - user manually merged settings, just clean our entries
                        $isManualMerge = $true
                        Write-StatusMessage "Manual merge scenario detected - will clean AI settings but preserve user's merged content" "Info"
                    } else {
                        # Original file existed - try to restore from backup (backupLength > 0)
                        Write-StatusMessage "Original settings.json existed (length: $backupLength) - will restore from backup" "Info"
                    }
                }
            } catch {
                Write-StatusMessage "Could not read metadata from settings.json - will attempt manual cleanup" "Warning"
            }
            
            if ($shouldDelete) {
                # Delete the entire file since we created it from scratch
                try {
                    Remove-Item -Path $paths.VSCodeSettings -Force -ErrorAction Stop
                    Write-StatusMessage "Deleted settings.json (no original file existed)" "Success"
                } catch {
                    Write-StatusMessage "Failed to delete settings.json: $($_.Exception.Message)" "Error"
                    $results.Success = $false
                    $results.Errors += "Failed to delete settings.json"
                }
            } elseif ($isManualMerge) {
                # Manual merge scenario - DO NOT touch settings.json, user manually merged and must manually clean
                Write-StatusMessage "Manual merge scenario detected - user must manually clean AI settings from settings.json" "Warning"
                Write-StatusMessage "We will NOT modify your settings.json file (you manually merged, you manually clean)" "Info"
                Write-StatusMessage "Please remove Terraform AzureRM AI settings from your settings.json manually" "Warning"
                
                # Remove the manual merge backup marker file
                if (Test-Path $paths.VSCodeBackupDir) {
                    try {
                        $manualMergeBackups = Get-ChildItem -Path $paths.VSCodeBackupDir -Filter "settings_backup_*.json" | ForEach-Object {
                            try {
                                $content = Get-Content $_.FullName -Raw | ConvertFrom-Json
                                if ($content.'// AZURERM_BACKUP_LENGTH' -eq -1) {
                                    $_
                                }
                            } catch {
                                # Ignore malformed backup files
                            }
                        }
                        
                        foreach ($backup in $manualMergeBackups) {
                            Remove-Item -Path $backup.FullName -Force
                            Write-StatusMessage "Removed manual merge marker: $($backup.Name)" "Success"
                        }
                    } catch {
                        Write-StatusMessage "Failed to remove manual merge marker files" "Warning"
                    }
                }
                
                Write-StatusMessage "Manual merge cleanup completed - settings.json left untouched" "Success"
            } else {
                # Try to restore from backup or clean manually
                if (Test-Path $paths.VSCodeBackupDir) {
                    $latestBackup = Get-MostRecentBackup -BackupDirectory $paths.VSCodeBackupDir -FilePattern "settings_backup_*.json"
                    
                    if ($latestBackup) {
                        if (Restore-FromBackup -BackupPath $latestBackup -TargetPath $paths.VSCodeSettings) {
                            Write-StatusMessage "VS Code settings restored from backup" "Success"
                        } else {
                            Write-StatusMessage "Failed to restore VS Code settings" "Error"
                            $results.Success = $false
                            $results.Errors += "Failed to restore VS Code settings from backup"
                        }
                    } else {
                        # No backup found, remove our settings manually
                        if (Remove-TerraformSettingsFromVSCode -SettingsPath $paths.VSCodeSettings) {
                            Write-StatusMessage "Terraform settings removed from VS Code" "Success"
                        } else {
                            Write-StatusMessage "Failed to clean Terraform settings from VS Code" "Warning"
                            $results.Errors += "Failed to clean Terraform settings from VS Code"
                        }
                    }
                    
                    # Preserve backup directory
                    Write-StatusMessage "Backup directory preserved at: $($paths.VSCodeBackupDir)" "Info"
                    Write-StatusMessage "Backups contain your original VS Code settings - keep them safe!" "Info"
                } else {
                    # No backup directory, just clean manually
                    if (Remove-TerraformSettingsFromVSCode -SettingsPath $paths.VSCodeSettings) {
                        Write-StatusMessage "Terraform settings removed from VS Code" "Success"
                    } else {
                        Write-StatusMessage "Failed to clean Terraform settings from VS Code" "Warning"
                        $results.Errors += "Failed to clean Terraform settings from VS Code"
                    }
                }
            }
        } else {
            Write-StatusMessage "No settings.json found - nothing to clean" "Info"
        }
        
        # Remove copied instruction files from user's VS Code
        if (Test-Path $paths.InstructionsDir) {
            try {
                Remove-Item -Path $paths.InstructionsDir -Recurse -Force -ErrorAction Stop
                Write-StatusMessage "Removed copied instruction files from VS Code" "Success"
            } catch {
                Write-StatusMessage "Failed to remove instruction files: $($_.Exception.Message)" "Warning"
                $results.Errors += "Failed to remove instruction files"
            }
        }
        
        # Remove copied prompt files from user's VS Code
        if (Test-Path $paths.PromptsDir) {
            try {
                Remove-Item -Path $paths.PromptsDir -Recurse -Force -ErrorAction Stop
                Write-StatusMessage "Removed copied prompt files from VS Code" "Success"
            } catch {
                Write-StatusMessage "Failed to remove prompt files: $($_.Exception.Message)" "Warning"
                $results.Errors += "Failed to remove prompt files"
            }
        }
        
        # Remove copied copilot-instructions.md from user's VS Code
        if (Test-Path $paths.CopilotInstructions) {
            try {
                Remove-Item -Path $paths.CopilotInstructions -Force -ErrorAction Stop
                Write-StatusMessage "Removed copied copilot-instructions.md from VS Code" "Success"
            } catch {
                Write-StatusMessage "Failed to remove copilot-instructions.md: $($_.Exception.Message)" "Warning"
                $results.Errors += "Failed to remove copilot-instructions.md"
            }
        }
        
        # Note: We removed the copied AI files from VS Code directories
        # The original files remain in the repository untouched
        
        if ($results.Success) {
            Write-StatusMessage "AI installation removal completed successfully!" "Success"
            Write-StatusMessage "Your original VS Code settings backups are preserved for safety" "Info"
        } else {
            Write-StatusMessage "AI installation removal completed with some errors" "Warning"
            Write-StatusMessage "Your original VS Code settings backups are preserved for safety" "Info"
        }
        
    } catch {
        $results.Errors += "Removal failed: $($_.Exception.Message)"
        $results.Success = $false
        Write-StatusMessage "Removal failed: $($_.Exception.Message)" "Error"
    }
    
    return $results
}

# Helper functions (migrated from original modules)
function Get-MostRecentBackup {
    param([string]$BackupDirectory, [string]$FilePattern = "*")
    
    if (-not (Test-Path $BackupDirectory)) {
        return $null
    }
    
    try {
        $backups = Get-ChildItem -Path $BackupDirectory -Filter $FilePattern | 
                   Where-Object { -not $_.PSIsContainer } |
                   Sort-Object LastWriteTime -Descending
        
        if ($backups.Count -gt 0) {
            return $backups[0].FullName
        }
    } catch {
        Write-StatusMessage "Error finding backup files: $_" "Warning"
    }
    
    return $null
}

function Restore-FromBackup {
    param([string]$BackupPath, [string]$TargetPath)
    
    if (-not (Test-Path $BackupPath)) {
        Write-StatusMessage "Backup file not found: $BackupPath" "Error"
        return $false
    }
    
    try {
        # Verify backup integrity before restore
        if (-not (Test-FileIntegrity -FilePath $BackupPath)) {
            Write-StatusMessage "Backup file is corrupted: $BackupPath" "Error"
            return $false
        }
        
        # Create target directory if needed
        $targetDir = Split-Path $TargetPath -Parent
        if (-not (Test-Path $targetDir)) {
            New-Item -ItemType Directory -Path $targetDir -Force | Out-Null
        }
        
        # Read backup content and strip our metadata before restoring
        $backupContent = Get-Content $BackupPath -Raw
        $backupData = $backupContent | ConvertFrom-Json
        
        # Convert to hashtable if needed for PowerShell 5.1 compatibility
        if ($backupData -is [PSCustomObject]) {
            $hashtable = @{}
            $backupData.PSObject.Properties | ForEach-Object {
                $hashtable[$_.Name] = $_.Value
            }
            $backupData = $hashtable
        }
        
        # Remove our metadata from the restored content
        $metadataKeys = @("// AZURERM_BACKUP_LENGTH", "// AZURERM_INSTALLATION_DATE")
        foreach ($key in $metadataKeys) {
            if ($backupData.ContainsKey($key)) {
                $backupData.Remove($key)
            }
        }
        
        # Save cleaned content to target
        $backupData | ConvertTo-Json -Depth 10 | Set-Content -Path $TargetPath -Encoding UTF8
        Write-StatusMessage "Restored from backup (metadata cleaned): $TargetPath" "Success"
        return $true
    } catch {
        Write-StatusMessage "Failed to restore from backup: $_" "Error"
        return $false
    }
}

function New-SafeBackup {
    <#
    .SYNOPSIS
        Creates a safe backup with proper metadata
        Metadata: AZURERM_BACKUP_LENGTH (0 = didn't exist, -1 = couldn't read, length > 0 = valid backup)
                 AZURERM_INSTALLATION_DATE
    #>
    param(
        [string]$SourcePath,
        [string]$BackupReason = "backup"
    )
    
    # Create timestamped backup path
    $timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
    $sourceFileName = [System.IO.Path]::GetFileNameWithoutExtension($SourcePath)
    $sourceExtension = [System.IO.Path]::GetExtension($SourcePath)
    $backupFileName = "${sourceFileName}_backup_${timestamp}${sourceExtension}"
    
    # Use standard backup directory
    $backupDir = Join-Path $env:APPDATA "Code\User\.terraform-azurerm-backups"
    if (-not (Test-Path $backupDir)) {
        New-Item -ItemType Directory -Path $backupDir -Force | Out-Null
    }
    
    $backupPath = Join-Path $backupDir $backupFileName
    $installationDate = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    
    try {
        if (-not (Test-Path $SourcePath)) {
            # File doesn't exist - create backup indicating this (length = 0)
            $noFileBackup = @{
                "// AZURERM_BACKUP_LENGTH" = 0
                "// AZURERM_INSTALLATION_DATE" = $installationDate
            }
            
            $noFileBackup | ConvertTo-Json -Depth 3 | Set-Content -Path $backupPath -Encoding UTF8
            Write-StatusMessage "Created backup (no original file): $backupFileName" "Success"
            
            return @{
                Success = $true
                BackupPath = $backupPath
                BackupFileName = $backupFileName
                OriginalFileLength = 0
            }
        }
        
        # File exists - try to read it
        try {
            $originalContent = Get-Content -Path $SourcePath -Raw -ErrorAction Stop
            $originalLength = (Get-Item $SourcePath).Length
            
            # If original is JSON, merge metadata into it
            if ($SourcePath -like "*.json") {
                try {
                    $originalJson = $originalContent | ConvertFrom-Json -ErrorAction Stop
                    
                    # Convert to hashtable and add metadata
                    $mergedContent = @{}
                    if ($originalJson -is [PSCustomObject]) {
                        $originalJson.PSObject.Properties | ForEach-Object {
                            $mergedContent[$_.Name] = $_.Value
                        }
                    } else {
                        $mergedContent = $originalJson
                    }
                    
                    # Add the two required metadata fields
                    $mergedContent["// AZURERM_BACKUP_LENGTH"] = $originalLength
                    $mergedContent["// AZURERM_INSTALLATION_DATE"] = $installationDate
                    
                    $mergedContent | ConvertTo-Json -Depth 10 | Set-Content -Path $backupPath -Encoding UTF8
                    
                } catch {
                    # File has .json extension but isn't valid JSON - manual merge scenario
                    Write-StatusMessage "Settings.json contains invalid JSON syntax - manual merge required" "Warning"
                    throw "Manual merge required: Settings.json contains invalid JSON syntax"
                }
            } else {
                # For now, only JSON files are expected in backup operations
                Write-StatusMessage "Backup system only handles JSON files - manual merge required" "Warning"
                throw "Manual merge required: Backup system only handles JSON files"
            }
            
            Write-StatusMessage "Created backup: $backupFileName (original length: $originalLength)" "Success"
            
            return @{
                Success = $true
                BackupPath = $backupPath
                BackupFileName = $backupFileName
                OriginalFileLength = $originalLength
            }
            
        } catch {
            # File exists but couldn't read it - could be permissions, lock, or corruption
            Write-StatusMessage "Cannot read existing settings.json file - manual merge required" "Warning"
            throw "Manual merge required: Cannot read settings.json file"
        }
        
    } catch {
        # Any other error during backup process - manual merge scenario
        Write-StatusMessage "Backup creation failed - proceeding with manual merge" "Warning"
        throw "Manual merge required: $($_.Exception.Message)"
    }
}


function Remove-TerraformSettingsFromVSCode {
    param([string]$SettingsPath)
    
    if (-not (Test-Path $SettingsPath)) {
        return $true
    }
    
    try {
        $content = Get-Content $SettingsPath -Raw
        $settings = $content | ConvertFrom-Json
        
        # Convert to hashtable if needed for PowerShell 5.1 compatibility
        if ($settings -is [PSCustomObject]) {
            $hashtable = @{}
            $settings.PSObject.Properties | ForEach-Object {
                $hashtable[$_.Name] = $_.Value
            }
            $settings = $hashtable
        }
        
        # Remove Terraform AzureRM specific settings (both old format and current format)
        $terraformKeys = $settings.Keys | Where-Object { 
            $_ -like "terraform_azurerm_*" -or $_ -like "// AZURERM_*" 
        }
        foreach ($key in $terraformKeys) {
            $settings.Remove($key)
        }
        
        # Also remove our specific GitHub Copilot settings that were added
        $copilotKeys = @(
            "github.copilot.chat.commitMessageGeneration.instructions",
            "github.copilot.chat.summarizeAgentConversationHistory.enabled", 
            "github.copilot.chat.reviewSelection.enabled",
            "github.copilot.chat.reviewSelection.instructions",
            "github.copilot.advanced",
            "github.copilot.enable"
        )
        
        foreach ($key in $copilotKeys) {
            if ($settings.ContainsKey($key)) {
                $settings.Remove($key)
            }
        }
        
        # Remove file associations we added
        if ($settings.ContainsKey("files.associations")) {
            $fileAssoc = $settings["files.associations"]
            if ($fileAssoc -is [PSCustomObject]) {
                $assocHash = @{}
                $fileAssoc.PSObject.Properties | ForEach-Object {
                    $assocHash[$_.Name] = $_.Value
                }
                $fileAssoc = $assocHash
            }
            
            # Remove our specific associations
            if ($fileAssoc.ContainsKey("*.instructions.md")) {
                $fileAssoc.Remove("*.instructions.md")
            }
            if ($fileAssoc.ContainsKey(".github/*.md")) {
                $fileAssoc.Remove(".github/*.md")
            }
            
            # If file associations is now empty, remove the whole section
            if ($fileAssoc.Count -eq 0) {
                $settings.Remove("files.associations")
            } else {
                $settings["files.associations"] = $fileAssoc
            }
        }
        
        # Write back the cleaned settings
        $settings | ConvertTo-Json -Depth 10 | Set-Content -Path $SettingsPath -Encoding UTF8
        
        Write-StatusMessage "Cleaned Terraform settings from VS Code" "Success"
        return $true
    } catch {
        Write-StatusMessage "Failed to clean VS Code settings: $_" "Error"
        return $false
    }
}

function Test-FileIntegrity {
    param([string]$FilePath, [string]$ExpectedPattern = "")
    
    if (-not (Test-Path $FilePath)) {
        return $false
    }
    
    try {
        $content = Get-Content $FilePath -Raw -ErrorAction Stop
        
        # Basic length check (minimum viable content)
        if ($content.Length -lt 50) {
            return $false
        }
        
        # Pattern check if provided
        if ($ExpectedPattern -and $content -notmatch $ExpectedPattern) {
            return $false
        }
        
        return $true
    } catch {
        return $false
    }
}

Export-ModuleMember -Function Get-FileManifest, Test-AIInstallation, Install-InstructionFiles, Install-PromptFiles, Install-MainFiles, Update-VSCodeSettings, Start-CompleteInstallation, Remove-AIInstallation
