# AzureRM Provider AI-Powered Development Installation Scripts

This directory contains installation scripts for setting up AI-powered coding assistance specifically tuned for the Terraform AzureRM provider development.

## Quick Start

### Windows (PowerShell)
```powershell
# Auto-discover repository and install
.\AILocalInstall\install-copilot-setup.ps1

# Or specify repository path
.\AILocalInstall\install-copilot-setup.ps1 -RepositoryPath "D:\Projects\terraform-provider-azurerm"

# Clean up and restore original settings
.\AILocalInstall\install-copilot-setup.ps1 -Clean

# Show help
.\AILocalInstall\install-copilot-setup.ps1 -Help
```

### Linux/macOS/WSL (Bash)
```bash
# Auto-discover repository and install
./AILocalInstall/install-copilot-setup.sh

# Or specify repository path
./AILocalInstall/install-copilot-setup.sh -repository-path "/path/to/terraform-provider-azurerm"

# Clean up and restore original settings
./AILocalInstall/install-copilot-setup.sh -clean

# Show help
./AILocalInstall/install-copilot-setup.sh -help
```

## What These Scripts Do

1. **🔍 Auto-discover** the terraform-provider-azurerm repository location
2. **📚 Copy instruction files** to your VS Code user directory
3. **🤖 Copy AI prompt files** for enhanced coding assistance
4. **⚙️ Configure VS Code settings** with intelligent merge of existing configuration
5. **💾 Create backups** of your existing settings for safe restoration

## Features

- **Smart Repository Detection**: Automatically finds the repository from any subdirectory
- **Intelligent Settings Merge**: Preserves your existing VS Code configuration
- **Safe Backup System**: Creates timestamped backups before making changes
- **Clean Uninstall**: Complete removal with original settings restoration
- **Cross-Platform**: Works on Windows (PowerShell), Linux, macOS, and WSL

## Complete Local Installation Features

This installation method provides a full local setup with these comprehensive capabilities:

- **🔍 Auto-discovery**: Automatically finds your repository location
- **🛡️ Manual backup prompt**: Pauses to let you create your own backup copy
- **🕐 Smart backup**: Creates timestamped backups of existing settings
- **🔀 Intelligent merge**: Preserves your existing VS Code configuration
- **📂 User input fallback**: Prompts for repository path if auto-discovery fails
- **🧹 Clean mode**: Complete removal of AI setup with settings restoration
- **⚙️ Settings validation**: Verifies configuration integrity after installation
- **📋 Detailed logging**: Comprehensive feedback during installation process

## Detailed Usage Examples

### Windows PowerShell Examples
```powershell
# Auto-discover repository (run from repository root)
.\AILocalInstall\install-copilot-setup.ps1

# Specify repository path explicitly
.\AILocalInstall\install-copilot-setup.ps1 -RepositoryPath "D:\Projects\terraform-provider-azurerm"

# Clean up - remove all AI setup and restore original settings
.\AILocalInstall\install-copilot-setup.ps1 -Clean
.\AILocalInstall\install-copilot-setup.ps1 -RepositoryPath "D:\Projects\terraform-provider-azurerm" -Clean

# Show detailed help information
.\AILocalInstall\install-copilot-setup.ps1 -Help
```

### Linux/macOS/WSL Bash Examples
```bash
# Auto-discover repository (run from repository root)
./AILocalInstall/install-copilot-setup.sh

# Specify repository path explicitly  
./AILocalInstall/install-copilot-setup.sh --repository-path "/home/user/projects/terraform-provider-azurerm"

# Clean up - remove all AI setup and restore original settings
./AILocalInstall/install-copilot-setup.sh --clean
./AILocalInstall/install-copilot-setup.sh --repository-path "/home/user/projects/terraform-provider-azurerm" --clean

# Show detailed help information
./AILocalInstall/install-copilot-setup.sh --help
```

## Files Installed

### VS Code User Directory (`~/.vscode/`)
```
├── copilot-instructions.md           # Main AI coding instructions
├── instructions/terraform-azurerm/   # Detailed implementation guides
│   ├── implementation-guide.instructions.md
│   ├── azure-patterns.instructions.md
│   ├── testing-guidelines.instructions.md
│   ├── documentation-guidelines.instructions.md
│   └── ... (other specialized guides)
├── prompts/                          # AI conversation prompts
│   ├── *.prompt.md                   # Various coding scenario prompts
│   └── ...
└── settings.json                     # Enhanced with AI configuration
```

### Configuration Added to settings.json
```json
{
  // Commit message generation with Azure provider context
  "github.copilot.chat.commitMessageGeneration.instructions": [
    {
      "text": "Provide a concise and clear commit message that summarizes the changes made in the code. For complex changes, include the following details: 1) Specify if the change introduces a breaking change and describe its impact. 2) Highlight any new resources or features added. 3) Mention updates to Azure services or APIs. Aim to keep the message under 72 characters per line for readability."
    }
  ],

  // Disable conversation history for privacy
  "github.copilot.chat.summarizeAgentConversationHistory.enabled": false,

  // Enable code review with instruction files
  "github.copilot.chat.reviewSelection.enabled": true,
  "github.copilot.chat.reviewSelection.instructions": [
    {"file": ".github/copilot-instructions.md"},
    {"file": ".github/instructions/implementation-guide.instructions.md"},
    {"file": ".github/instructions/azure-patterns.instructions.md"},
    {"file": ".github/instructions/testing-guidelines.instructions.md"},
    {"file": ".github/instructions/documentation-guidelines.instructions.md"},
    {"file": ".github/instructions/provider-guidelines.instructions.md"}
  ],

  // File associations for proper syntax highlighting
  "files.associations": {
    "*.instructions.md": "markdown",
    ".github/*.md": "markdown"
  },

  // Additional Copilot optimization settings
  "github.copilot.advanced": {
    "length": 3000,
    "temperature": 0.1
  },

  // Enable Copilot across all relevant contexts and file types
  "github.copilot.enable": {
    "*": true,
    "terminal": true
  }
}
```

## Safety Features

- **Automatic Backups**: Creates timestamped backups before making changes
- **Conflict Detection**: Prevents installing over existing installations
- **Integrity Verification**: Tracks original file state for perfect restoration
- **Manual Backup Prompts**: Encourages additional manual backups for extra safety
- **Rollback Support**: Complete uninstall with original settings restoration

## Troubleshooting

### Installation Issues
- **Repository not found**: Ensure you're in the terraform-provider-azurerm directory or provide the correct path
- **Permission errors**: Run with appropriate permissions for your system
- **JSON merge errors**: The script will provide manual merge instructions

### Clean Up Issues
- **Backup not found**: Check for timestamped backup files in `~/.vscode/`
- **Manual merge detected**: Follow the provided manual cleanup instructions
- **Settings corruption**: Restore from your own backup if available

## Support

For more information and detailed documentation:
- **Repository**: https://github.com/hashicorp/terraform-provider-azurerm
- **Instructions**: See the `.github/instructions/` directory
- **Help**: Run the scripts with `-Help` (PowerShell) or `-help` (Bash) for detailed usage information

## Advanced Usage

### Development Environment Setup
The scripts are designed for developers working on the Terraform AzureRM provider. They configure VS Code with:
- Azure-specific coding patterns and best practices
- Terraform resource implementation guidance
- Testing and documentation standards
- Error handling and debugging assistance

### Customization
After installation, you can customize the AI behavior by:
- Editing the instruction files in `~/.vscode/instructions/terraform-azurerm/`
- Modifying prompts in `~/.vscode/prompts/`
- Adjusting VS Code settings in `settings.json`

All changes are preserved during clean operations, and you can safely reinstall to get updates.


## 🛡️ Settings Safety & Backup Strategy

### Settings Precedence & Conflicts

**🎯 Important: Understanding VS Code Settings Hierarchy**

If you work with both **local settings** (installed by our scripts) and **repository settings** (in the repo's `.vscode/` folder), here's how VS Code handles them:

**Precedence Order (Highest to Lowest):**
1. **Local User Settings** (`~/.vscode/settings.json`) - **ALWAYS WINS**
2. **Repository Workspace Settings** (`.vscode/settings.json`) - **FALLBACK ONLY**  
3. **Default VS Code Settings** - **BASELINE**

**Key Points:**
- **Local settings override repo settings** - Property by property, not file by file
- **No conflicts** - VS Code merges intelligently at the individual setting level
- **Best of both worlds** - You get repo settings for unconfigured properties + your local preferences for configured ones

**Example Scenario:**
```bash
# You have local settings with custom copilot length:
"github.copilot.advanced": { "length": 5000 }

# Repo has different length but additional settings:
"github.copilot.advanced": { "length": 3000, "temperature": 0.1 }

# Result: You get YOUR length (5000) + repo's temperature (0.1)
```

**Recommendation:**
- **Keep both** - Local settings for personal preferences, repo settings for project-specific needs
- **No cleanup needed** - They work together harmoniously
- **Local wins** - Your personal preferences are always respected

### 🔍 **Troubleshooting Settings Conflicts**

**Q: The script says "AzureRM Copilot settings already detected" - what does this mean?**

**A: Protection System Activation:**
- This means you've **already run the installation script before**
- Running it again would corrupt your backup by backing up **merged settings instead of original settings**
- **Solution**: Run clean mode first, then reinstall:
  ```bash
  # Clean first, then reinstall
  ./.github/AILocalInstall/install-copilot-setup.sh --clean
  ./.github/AILocalInstall/install-copilot-setup.sh
  ```

**Q: I get a warning about "original settings backup is missing" - what happened?**

**A: Missing Backup Recovery:**
- Your backup file (`settings.json.azurerm-backup`) was deleted or moved
- **Backup exists**: Follow normal clean → reinstall process
- **Backup missing**: You have two options:
  
  **Option 1 (Safe)**: Manual backup first
  ```bash
  # 1. Copy current settings to safe location
  cp ~/.vscode/settings.json ~/Desktop/my-settings-backup.json
  
  # 2. Run clean (will warn about missing backup)  
  ./.github/AILocalInstall/install-copilot-setup.sh --clean
  
  # 3. Manually restore if needed
  cp ~/Desktop/my-settings-backup.json ~/.vscode/settings.json
  
  # 4. Fresh install
  ./.github/AILocalInstall/install-copilot-setup.sh
  ```
  
  **Option 2 (Risky)**: Force proceed
  - Your current settings contain merged AI + original settings
  - Re-installing may create duplicate entries
  - Only do this if you're confident about your current settings

**Q: I see different copilot instructions than expected - which ones are being used?**

**A: Check VS Code's settings precedence:**

1. **Open VS Code Settings UI** (`Ctrl+,` or `Cmd+,`)
2. **Search for** `github.copilot.chat.reviewSelection.instructions`
3. **Look at the source indicator**:
   - `User` = Your local settings (takes precedence)
   - `Workspace` = Repository settings (fallback)

**Q: How do I see what settings are actually active?**

**A: Use VS Code's settings inspection:**
```bash
# View effective settings in VS Code
1. Press F1 → "Developer: Reload Window"
2. Press F1 → "Preferences: Open Settings (JSON)"
3. Check both User and Workspace tabs
```

**Q: I want to use ONLY repository settings, not local ones**

**A: Temporarily rename your local settings:**
```bash
# Windows PowerShell
Rename-Item "$env:APPDATA\Code\User\settings.json" "settings.json.backup"

# macOS/Linux Bash  
mv ~/.vscode/settings.json ~/.vscode/settings.json.backup
```

**Q: I want to completely remove the local AI setup**

**A: Use the clean mode:**
```bash
# This removes local settings and restores backups
./.github/AILocalInstall/install-copilot-setup.sh --clean
```

**✅ Pro Tip: Best of Both Worlds**
- Keep **project-specific** settings in repository `.vscode/settings.json`
- Keep **personal preferences** in local `~/.vscode/settings.json`  
- They merge intelligently - no conflicts!

---

### Smart Merge Technology

**🔐 Length-Based Installation Protection System**

The scripts use simple content length verification to ensure installation integrity and prevent data loss:

**How Length Protection Works:**
1. **Before installation**: Calculate byte length of your original settings.json (or "0" if no file exists)
2. **During installation**: Inject length comment into merged settings: `"// AZURERM_BACKUP_LENGTH": "1234"`
3. **On re-install attempt**: Verify length to determine safe recovery path
4. **During clean**: Use length to determine correct restoration behavior

**Length-Based Decision Matrix:**
```json
// Example settings.json after installation:
{
  "// AZURERM_BACKUP_LENGTH": "1234",
  "github.copilot.chat.reviewSelection.instructions": "...",
  "your": "existing settings preserved"
}
```

**Protection Scenarios:**
- **Length = "0"**: Originally empty file → Clean mode removes settings.json entirely
- **Length = "1234"**: Had original file → Clean mode restores from backup
- **Length = "-1"**: Manual merge installation → Clean mode provides manual removal guidance
- **Length missing**: Corruption detected → Manual intervention required
- **Backup missing + Length exists**: Backup corruption → Guided recovery options

**Benefits:**
- **Universal compatibility** - Content length available on all systems (no MD5 dependency)
- **Prevents false positives** - Won't trigger on fresh VS Code installations
- **Intelligent restoration** - Knows whether to restore or remove entirely  
- **Corruption detection** - Catches manual tampering or partial installations
- **User education** - Clear error messages explain exactly what happened

**Detection Examples:**
```bash
# Normal re-installation attempt:
AzureRM Copilot installation detected with length verification.
[OK] Installation verified: Backup integrity confirmed

To prevent losing your original configuration, please:
1. Run the clean command first: ./install-copilot-setup.sh --clean
2. Then run the installation again: ./install-copilot-setup.sh

# Manual merge clean attempt:
Manual merge detected - cannot auto-restore

MANUAL CLEANUP REQUIRED:
   This installation was done manually without automatic backup.
   You must manually remove AzureRM Copilot settings from your settings.json:

   Remove these entries from settings.json:
   - "// AZURERM_BACKUP_LENGTH": "-1"
   - "github.copilot.chat.commitMessageGeneration.instructions"
   - "github.copilot.chat.reviewSelection.instructions"
   - [Additional settings listed...]

   Be careful to preserve your other VS Code settings!

# Corruption detected:
INSTALLATION INTEGRITY ERROR
================================

We detected AzureRM Copilot settings in your settings.json,
but cannot verify the installation integrity.

This usually means:
  - Settings were manually modified
  - Length comment was removed or changed
  - Installation was partially corrupted

MANUAL RECOVERY REQUIRED:
  1. Back up your current settings.json to a safe location
  2. Manually remove all AzureRM Copilot entries from settings.json
  3. Run this script again for fresh installation
```

**Why This Protection Matters:**
- **Original → Merged → Backup** = ❌ You lose your original settings forever
- **Original → Backup → Merged** = ✅ You can always get back to your original settings

**Enhanced Detection Criteria:**
- Looks for `terraform-azurerm` patterns in your settings
- Checks for `copilot-instructions.md` references  
- **Verifies length comment** for installation integrity
- **Validates backup file** matches stored length
- Only triggers if patterns found AND length verification succeeds

**🚨 Missing Backup Protection:**
```bash
# If settings are installed but backup is missing, you'll see:
WARNING: Your original settings backup is missing!
Expected location: settings.json.azurerm-backup

This means clean mode cannot restore your original settings.
You have two options:

OPTION 1 (Recommended): Manual backup first
  1. Copy your current settings.json to a safe location
  2. Run: ./install-copilot-setup.sh --clean (will warn about missing backup)
  3. Manually restore your backed-up settings if needed
  4. Run: ./install-copilot-setup.sh

OPTION 2 (Advanced): Force proceed (RISKY)
  - Your current settings.json contains merged AI + original settings
  - Re-running install will merge AI settings again (may cause duplicates)
```

**Why This Protection Matters:**
- **User deleted backup**: Protection prevents permanent settings loss
- **Backup corruption**: Safeguards against installing over corrupted backups  
- **Manual intervention**: Provides clear recovery paths when automation fails
- **Data preservation**: Prioritizes keeping user settings over convenience

### Intelligent Settings Merging

The installation scripts use **intelligent merging** to protect your existing VS Code configuration:

**✅ What the scripts DO:**
- **Pause for manual backup opportunity** with clear file location and instructions
- **Detect existing settings.json** and create automatic backups
- **Preserve all your custom settings** (themes, extensions, keybindings, etc.)
- **Add only AzureRM provider-specific settings** without conflicts
- **Create timestamped backups** for easy restoration if needed

**❌ What the scripts DON'T do:**
- Overwrite your existing VS Code configuration
- Delete or modify your custom settings
- Change non-AzureRM related configurations

### Backup Strategy

**🛡️ Hash-Based Protection System:**

1. **Length Generation**: Content length of original settings.json before any changes
2. **Length Storage**: Embedded as comment in merged settings.json: `"// AZURERM_BACKUP_LENGTH": "length_value"`
3. **Length Verification**: Validates installation integrity and guides restoration
4. **Smart Restoration**: Uses length to determine correct clean behavior

**Length Values & Clean Behavior:**
- **Length "0"**: Originally no file existed → Clean removes settings.json entirely
- **Length "1234"**: Original file existed → Clean restores from backup
- **Length "-1"**: Manual merge installation → Clean provides manual removal instructions
- **Length missing**: Corruption detected → Clean requires manual intervention
- **Backup missing**: Guided recovery with multiple options

**Manual Backup Process:**
- Scripts display the exact path to your settings file
- You can copy it to Desktop, Documents, cloud storage, or anywhere you prefer
- This gives you 100% control over your backup location and method
- Choose **[C]** to continue or **[X]** to exit without making any changes

### Clean Mode Restoration

**🧹 Automatic Clean Mode (Recommended):**
```powershell
# PowerShell - Completely remove AI setup and restore original settings
.\install-copilot-setup.ps1 -Clean
```

```bash
# Bash - Completely remove AI setup and restore original settings  
./install-copilot-setup.sh --clean
```

**What Clean Mode Does:**
- ✅ **Removes all AI instruction files** from VS Code user directory
- ✅ **Removes all AI prompt files** (AzureRM-specific only)
- ✅ **Removes main copilot instructions** 
- ✅ **Restores original settings.json** from backup
- ✅ **Cleans up empty directories**

### Manual Restore Process

If you need to restore your original settings manually:

**PowerShell:**
```powershell
# Navigate to VS Code user directory
cd "$env:APPDATA\Code\User"

# List available backups
Get-ChildItem settings.json.backup.*

# Restore from specific backup
Copy-Item "settings.json.backup.20250806-143022" "settings.json" -Force
```

**Bash:**
```bash
# Navigate to VS Code user directory
cd ~/.vscode

# List available backups
ls settings.json.backup.*

# Restore from specific backup
cp settings.json.backup.20250806-143022 settings.json
```

### Merge Behavior Details

**Settings Precedence:**
1. **AzureRM provider settings** take precedence for Copilot chat instructions
2. **Your existing settings** are preserved for everything else
3. **File associations** are merged (both yours and AzureRM patterns work)

**Specific Merge Rules:**
- `github.copilot.chat.reviewSelection.instructions` → **Replaced** with AzureRM instructions
- `files.associations` → **Merged** (both yours and `*.instructions.md` patterns)
- All other settings → **Your settings preserved**

### Smart Merge Example

**BEFORE Installation (Your existing settings.json):**
```json
{
    "workbench.colorTheme": "Dark+ (default dark)",
    "editor.fontSize": 14,
    "editor.tabSize": 4,
    "files.associations": {
        "*.tf": "terraform",
        "*.tfvars": "terraform-vars"
    },
    "extensions.autoUpdate": true,
    "github.copilot.enable": {
        "*": true
    },
    "myCustomSetting": "important_value"
}
```

**AFTER Installation (Merged result):**
```json
{
    "workbench.colorTheme": "Dark+ (default dark)",
    "editor.fontSize": 14,
    "editor.tabSize": 4,
    "files.associations": {
        "*.tf": "terraform",
        "*.tfvars": "terraform-vars",
        "*.instructions.md": "markdown",
        "*.prompt.md": "markdown"
    },
    "extensions.autoUpdate": true,
    "github.copilot.enable": {
        "*": true,
        "terminal": true
    },
    "myCustomSetting": "important_value",
    "github.copilot.chat.commitMessageGeneration.instructions": [...],
    "github.copilot.chat.summarizeAgentConversationHistory.enabled": false,
    "github.copilot.chat.reviewSelection.enabled": true,
    "github.copilot.chat.reviewSelection.instructions": [
        {"file": "copilot-instructions.md"},
        {"file": "instructions/terraform-azurerm/implementation-guide.instructions.md"}
    ],
    "github.copilot.advanced": {
        "length": 3000,
        "temperature": 0.1
    }
}
```

**🎯 Key Benefits:**
- ✅ **All your settings preserved** (theme, font size, extensions, etc.)
- ✅ **File associations merged** (both terraform and markdown patterns work)
- ✅ **GitHub Copilot enhanced** with AzureRM provider intelligence
- ✅ **Custom settings untouched** (`myCustomSetting` remains intact)
- ✅ **Automatic backup created** for easy restoration

**⚠️ IMPORTANT: Known Limitation**

**Relative Path References**: Some instruction files contain relative path references to the contributing documentation (e.g., `../../../contributing/topics/`). These links will not work in the local installation. This is a known limitation that will be addressed in a future update.

**Current Workaround**: For complete documentation, refer to the original repository at:
- 📚 **Contributing Documentation**: [github.com/hashicorp/terraform-provider-azurerm/tree/main/contributing](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/contributing)
- 🔧 **Repository Instructions**: [github.com/hashicorp/terraform-provider-azurerm/tree/main/.github/instructions](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/.github/instructions)
