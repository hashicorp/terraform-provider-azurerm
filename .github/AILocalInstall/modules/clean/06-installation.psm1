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
    
    # Create backup if settings exist (real pattern)
    if ((Test-Path $settingsPath) -and -not $Force) {
        try {
            $backupResult = New-SafeBackup -SourcePath $settingsPath -BackupReason "AI installation"
            if (-not $backupResult) {
                Write-StatusMessage "Failed to create backup of VS Code settings" "Error"
                return $false
            }
        } catch {
            Write-StatusMessage "Backup creation failed: $_" "Error"
            return $false
        }
    }
    
    # Read existing settings or create new (handle JSONC with comments)
    $settings = @{}
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
            Write-StatusMessage "CRITICAL ERROR: VS Code settings.json exists but contains invalid JSON" "Error"
            Write-StatusMessage "This indicates VS Code environment corruption" "Error"
            Write-StatusMessage "Cannot proceed safely with corrupted JSON files" "Error"
            Write-StatusMessage "Corrupted VS Code settings.json detected: $settingsPath. Cannot proceed with installation." "Error"
            throw "Installation halted due to corrupted VS Code settings.json"
        }
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

        # Enable code review with instruction files (LOCAL PATHS)
        "github.copilot.chat.reviewSelection.enabled" = $true
        "github.copilot.chat.reviewSelection.instructions" = @(
            @{"file" = "copilot-instructions.md"}
            @{"file" = "instructions/terraform-azurerm/implementation-guide.instructions.md"}
            @{"file" = "instructions/terraform-azurerm/azure-patterns.instructions.md"}
            @{"file" = "instructions/terraform-azurerm/testing-guidelines.instructions.md"}
            @{"file" = "instructions/terraform-azurerm/documentation-guidelines.instructions.md"}
            @{"file" = "instructions/terraform-azurerm/provider-guidelines.instructions.md"}
            @{"file" = "instructions/terraform-azurerm/code-clarity-enforcement.instructions.md"}
            @{"file" = "instructions/terraform-azurerm/error-patterns.instructions.md"}
            @{"file" = "instructions/terraform-azurerm/migration-guide.instructions.md"}
            @{"file" = "instructions/terraform-azurerm/schema-patterns.instructions.md"}
            @{"file" = "instructions/terraform-azurerm/performance-optimization.instructions.md"}
            @{"file" = "instructions/terraform-azurerm/security-compliance.instructions.md"}
            @{"file" = "instructions/terraform-azurerm/troubleshooting-decision-trees.instructions.md"}
            @{"file" = "instructions/terraform-azurerm/api-evolution-patterns.instructions.md"}
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
        
        # Check overall success (real pattern)
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
        } else {
            Write-StatusMessage "Installation completed with issues" "Warning"
            if (-not $results.InstructionFiles.Success) { $results.Errors += $results.InstructionFiles.Errors }
            if (-not $results.PromptFiles.Success) { $results.Errors += $results.PromptFiles.Errors }
            if (-not $results.MainFiles.Success) { $results.Errors += $results.MainFiles.Errors }
            if (-not $results.VSCodeSettings) { $results.Errors += "VS Code settings configuration failed" }
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
        
        # Restore VS Code settings from backup (original logic)
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
                # No backup found, remove our settings manually (original logic)
                if (Test-Path $paths.VSCodeSettings) {
                    if (Remove-TerraformSettingsFromVSCode -SettingsPath $paths.VSCodeSettings) {
                        Write-StatusMessage "Terraform settings removed from VS Code" "Success"
                    } else {
                        Write-StatusMessage "Failed to clean Terraform settings from VS Code" "Warning"
                        $results.Errors += "Failed to clean Terraform settings from VS Code"
                    }
                }
            }
            
            # NEVER delete backup directory - users need it to restore original settings!
            # Remove-BackupFiles -BackupDirectory $paths.VSCodeBackupDir  # DANGEROUS - commented out
            Write-StatusMessage "Backup directory preserved at: $($paths.VSCodeBackupDir)" "Info"
            Write-StatusMessage "Backups contain your original VS Code settings - keep them safe!" "Info"
        } else {
            Write-StatusMessage "No backup directory found - checking for manual cleanup" "Info"
            
            # Still try manual cleanup if no backup directory exists
            if (Test-Path $paths.VSCodeSettings) {
                if (Remove-TerraformSettingsFromVSCode -SettingsPath $paths.VSCodeSettings) {
                    Write-StatusMessage "Terraform settings removed from VS Code" "Success"
                } else {
                    Write-StatusMessage "Failed to clean Terraform settings from VS Code" "Warning"
                    $results.Errors += "Failed to clean Terraform settings from VS Code"
                }
            }
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
        
        Copy-Item -Path $BackupPath -Destination $TargetPath -Force -ErrorAction Stop
        Write-StatusMessage "Restored from backup: $TargetPath" "Success"
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
                    # File has .json extension but isn't valid JSON - this indicates corruption
                    Write-StatusMessage "CRITICAL ERROR: File $SourcePath has .json extension but contains invalid JSON" "Error"
                    Write-StatusMessage "This indicates VS Code environment corruption" "Error"
                    Write-StatusMessage "Cannot proceed safely with corrupted JSON files" "Error"
                    throw "Corrupted JSON file detected: $SourcePath. VS Code environment may be damaged."
                }
            } else {
                # For now, only JSON files are expected in backup operations
                # If we're backing up non-JSON files, that's unexpected
                Write-StatusMessage "CRITICAL ERROR: Unexpected file type for backup: $SourcePath" "Error"
                Write-StatusMessage "Backup system only handles JSON files (like settings.json)" "Error"
                Write-StatusMessage "Non-JSON file backup indicates environment corruption" "Error"
                throw "Unexpected file type for backup: $SourcePath. Expected JSON files only."
            }
            
            Write-StatusMessage "Created backup: $backupFileName (original length: $originalLength)" "Success"
            
            return @{
                Success = $true
                BackupPath = $backupPath
                BackupFileName = $backupFileName
                OriginalFileLength = $originalLength
            }
            
        } catch {
            # File exists but couldn't read it - this is an ERROR condition, not a backup scenario
            Write-StatusMessage "CRITICAL ERROR: Cannot read existing file $SourcePath`: $($_.Exception.Message)" "Error"
            Write-StatusMessage "This file may be corrupted, locked, or have permission issues" "Error"
            Write-StatusMessage "Backup operation cannot proceed safely" "Error"
            throw "Cannot backup unreadable file: $SourcePath. Error: $($_.Exception.Message)"
        }
        
    } catch {
        # Any other error during backup process - this should also halt
        Write-StatusMessage "CRITICAL ERROR: Backup operation failed for $SourcePath`: $($_.Exception.Message)" "Error"
        Write-StatusMessage "Cannot proceed with installation without proper backup" "Error"
        throw "Backup operation failed: $($_.Exception.Message)"
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
        
        # Remove Terraform AzureRM specific settings
        $terraformKeys = $settings.Keys | Where-Object { $_ -like "terraform_azurerm_*" }
        foreach ($key in $terraformKeys) {
            $settings.Remove($key)
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
