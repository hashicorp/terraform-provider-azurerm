# AI Infrastructure Installer

> **Intelligent setup tool for AI-powered Terraform AzureRM Provider development**

This installer provides GitHub Copilot instructions, VS Code configurations, and AI-powered development workflows for the Terraform AzureRM Provider repository.

## ✨ User Experience

The installer provides a **clean, professional output** focused on what matters:

- **📋 Clear progress indicators** - Section headers and completion status
- **🎯 Focused messaging** - Only essential information, no technical noise
- **✅ Success confirmations** - Clear indication when operations complete
- **📊 File operation tracking** - Detailed file copy/install status

**Example output:**
```
[SECTION] Bootstrap - Copying Installer to User Profile
-------------------------------------------------------
  Using existing directory: C:\Users\username\.terraform-ai-installer
  Copying installer files from local source repository...

    Copying: file-manifest.config [OK]
    Copying: install-copilot-setup.ps1 [OK]
    Copying: ConfigParser.psm1 [OK]

[SUCCESS] Bootstrap completed successfully!
```

## 🚀 Getting Started

### 📋 Prerequisites

Before using the AI-powered development features, ensure you have the following VS Code extensions installed:

#### Required Extensions
- **[GitHub Copilot](https://marketplace.visualstudio.com/items?itemName=GitHub.copilot)** - Core AI assistance and code generation
- **[GitHub Copilot Chat](https://marketplace.visualstudio.com/items?itemName=GitHub.copilot-chat)** - Interactive AI chat and slash commands

#### Recommended Configuration
For optimal performance and the best AI assistance experience:

- **LLM Model**: **Claude Sonnet 4** (optimal for Terraform and Go development)
- **Copilot Chat Model**: Set to Claude Sonnet 4 in VS Code settings
- **Alternative**: GPT-4o (good alternative if Claude is unavailable)

#### Quick Setup
1. Install the required extensions from VS Code marketplace
2. Sign in to GitHub Copilot with your account
3. Configure Claude Sonnet 4 as your preferred model:
   - Open VS Code Settings (`Ctrl+,`)
   - Search for "copilot chat model"
   - Select "Claude Sonnet 4" from the dropdown

**Note**: Claude Sonnet 4 provides superior understanding of Terraform patterns, Azure API specifics, and Go code generation compared to other models.

### First Time Setup

If you don't have the repository yet, clone it and switch to the AI installer branch:

```powershell
# Clone the repository
git clone https://github.com/hashicorp/terraform-provider-azurerm.git
cd terraform-provider-azurerm

# Switch to the AI installer source branch
git checkout exp/terraform_copilot

# Navigate to the installer
cd .github\AIinstaller

# Bootstrap the installer for use on feature branches
.\install-copilot-setup.ps1 -Bootstrap
```

### ⚠️ PowerShell Execution Policy

If you encounter execution policy errors, you have several options:

#### Option 1: Bypass for single execution (Recommended)
```powershell
# Run with execution policy bypass (safest for one-time use)
powershell -ExecutionPolicy Bypass -File .\install-copilot-setup.ps1 -Bootstrap

# Or for the user profile installer
powershell -ExecutionPolicy Bypass -File "$env:USERPROFILE\.terraform-ai-installer\install-copilot-setup.ps1"
```

#### Option 2: Unblock the downloaded files
```powershell
# Unblock all installer files
Get-ChildItem -Path . -Recurse | Unblock-File
.\install-copilot-setup.ps1 -Bootstrap
```

#### Option 3: Set execution policy for current user (Permanent)
```powershell
# Allow local scripts for current user
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
.\install-copilot-setup.ps1 -Bootstrap
```

## 🚀 Quick Start

### If you're on the source branch (`exp/terraform_copilot`):
```powershell
# 1. Bootstrap the installer
.\install-copilot-setup.ps1 -Bootstrap

# 2. Switch to your feature branch
git checkout feature/your-branch-name

# 3. Run installer from user profile (MUST specify RepoDirectory)
& "$env:USERPROFILE\.terraform-ai-installer\install-copilot-setup.ps1" -RepoDirectory "C:\github.com\hashicorp\terraform-provider-azurerm"
```

### If you're on a feature branch and have local installer:
```powershell
# Run installer directly (RepoDirectory auto-detected)
.\install-copilot-setup.ps1
```

### If you're on a feature branch and using user profile installer:
```powershell
# MUST specify RepoDirectory when running from user profile
& "$env:USERPROFILE\.terraform-ai-installer\install-copilot-setup.ps1" -RepoDirectory "C:\path\to\your\terraform-provider-azurerm"
```

## 📋 What Gets Installed

The installer sets up a complete AI development environment:

### 📄 Core AI Instructions
- **`.github/copilot-instructions.md`** - Main Copilot configuration
- **`.github/instructions/`** - 14 specialized instruction files:
  - `implementation-guide.instructions.md` - Complete coding standards
  - `azure-patterns.instructions.md` - Azure-specific patterns
  - `testing-guidelines.instructions.md` - Testing requirements
  - `documentation-guidelines.instructions.md` - Doc standards
  - `error-patterns.instructions.md` - Error handling
  - `migration-guide.instructions.md` - Version migration
  - `provider-guidelines.instructions.md` - Provider patterns
  - `schema-patterns.instructions.md` - Schema design
  - `code-clarity-enforcement.instructions.md` - Code quality
  - `performance-optimization.instructions.md` - Performance
  - `security-compliance.instructions.md` - Security patterns
  - `troubleshooting-decision-trees.instructions.md` - Debugging
  - `api-evolution-patterns.instructions.md` - API versioning

### 🎨 Development Templates & AI Prompts
- **`.github/prompts/`** - AI prompt templates for common development tasks:

#### Available Prompts

| Prompt File | Purpose | Usage |
|-------------|---------|-------|
| `add-unit-tests.prompt.md` | **Generate unit tests** for Go files with proper testing patterns | Use when you need comprehensive unit tests for utility functions, parsers, or validators |
| `code-review-local-changes.prompt.md` | **Review uncommitted changes** with Terraform provider best practices | Use before committing to get expert feedback on your local changes |
| `code-review-committed-changes.prompt.md` | **Review committed changes** for pull request feedback | Use to review git commits with detailed technical analysis |
| `setup-go-dev-environment.prompt.md` | **Set up Go development environment** for the Terraform provider | Use when setting up a new development machine or environment |
| `summarize-repo.prompt.md` | **Generate repository overview** with key components and structure | Use to understand the codebase structure and organization |
| `summarize-repo-deep-dive.prompt.md` | **Deep technical analysis** of repository architecture and patterns | Use for comprehensive codebase analysis and architecture review |

#### How to Use Prompts

**In GitHub Copilot Chat:**
**Prerequisites**: Requires [GitHub Copilot](https://marketplace.visualstudio.com/items?itemName=GitHub.copilot) and [GitHub Copilot Chat](https://marketplace.visualstudio.com/items?itemName=GitHub.copilot-chat) extensions. Use **Claude Sonnet 4** for optimal results.

Simply use slash commands to invoke the prompts directly:

| Slash Command | Prompt File | Description |
|---------------|-------------|-------------|
| `/code-review-local-changes` | `code-review-local-changes.prompt.md` | Review your uncommitted changes |
| `/code-review-committed-changes` | `code-review-committed-changes.prompt.md` | Review committed changes |
| `/add-unit-tests` | `add-unit-tests.prompt.md` | Generate unit tests for Go files |
| `/setup-go-dev-environment` | `setup-go-dev-environment.prompt.md` | Set up Go development environment |
| `/summarize-repo` | `summarize-repo.prompt.md` | Generate repository overview |
| `/summarize-repo-deep-dive` | `summarize-repo-deep-dive.prompt.md` | Deep technical analysis |

**Example Usage:**
```
/code-review-local-changes
/code-review-local-committed-changes
```

**For File-Specific Actions:**
```
/code-review-local-committed-changes for focusing on Azure patterns that have already been committed to your branch
/code-review-local-changes focusing on Azure patterns
```

### ⚙️ VS Code Configuration
- **`.vscode/settings.json`** - Optimized VS Code settings for Terraform development


## 🎛️ Command Reference

### Core Commands

| Command | Description | Available On |
|---------|-------------|--------------|
| `.\install-copilot-setup.ps1` | **Install AI infrastructure** | Feature branches |
| `.\install-copilot-setup.ps1 -Bootstrap` | **Copy installer to user profile** | Source branch only |
| `.\install-copilot-setup.ps1 -Verify` | **Check installation status** | Any branch |
| `.\install-copilot-setup.ps1 -Clean` | **Remove AI infrastructure** | Feature branches |
| `.\install-copilot-setup.ps1 -Help` | **Show detailed help** | Any branch |

### Parameters

| Parameter | Description | Required When | Example |
|-----------|-------------|---------------|---------|
| `-RepoDirectory` | **Specify repository path** | Running from user profile | `-RepoDirectory "C:\path\to\terraform-provider-azurerm"` |
| `-Auto-Approve` | Skip confirmation prompts | Optional | `-Auto-Approve` |
| `-Dry-Run` | Preview changes without applying | Optional | `-Dry-Run` |

### 🚨 Important: `-RepoDirectory` Parameter

When running the installer **from your user profile** (after bootstrap), you **MUST** specify the `-RepoDirectory` parameter:

```powershell
# ✅ CORRECT: Running from user profile with RepoDirectory
& "$env:USERPROFILE\.terraform-ai-installer\install-copilot-setup.ps1" -RepoDirectory "C:\github.com\hashicorp\terraform-provider-azurerm"

# ❌ INCORRECT: Running from user profile without RepoDirectory
& "$env:USERPROFILE\.terraform-ai-installer\install-copilot-setup.ps1"
```

**Why is this required?**
- The installer needs to know where your git repository is located
- Enables proper branch detection and workspace validation  
- Ensures files are installed in the correct repository directory

## 🌊 Workflow Overview

### Branch-Aware Architecture

The installer adapts its behavior based on your current Git branch:

#### 🔹 Source Branch (`exp/terraform_copilot`)
- Contains the master AI infrastructure files
- **Bootstrap mode**: Copies installer to user profile for feature branch use
- **Verification**: Checks source files integrity
- **Protection**: Prevents accidental deletion of development files

#### 🔹 Feature Branches (any other branch)
- Target for AI infrastructure installation
- **Install mode**: Sets up complete AI development environment
- **Clean mode**: Removes AI infrastructure when needed
- **Verification**: Checks installed files status

### Typical Development Workflow

```mermaid
graph LR
    A[Source Branch<br/>exp/terraform_copilot] --> B[Bootstrap<br/>-Bootstrap]
    B --> C[Switch Branch<br/>git checkout feature/xyz]
    C --> D[Install AI<br/>Run from user profile]
    D --> E[Develop<br/>With AI assistance]
    E --> F[Clean<br/>-Clean when done]
```

## 🏗️ Architecture

### Module Structure

```
AIinstaller/
├── install-copilot-setup.ps1      # Main entry point
├── README.md                      # This file
└── modules/powershell/            # PowerShell modules
    ├── ConfigParser.psm1          # Configuration management
    ├── FileOperations.psm1        # File installation/removal
    ├── ValidationEngine.psm1      # System validation
    └── UI.psm1                    # User interface functions
```

### Module Responsibilities

#### 📋 ConfigParser.psm1
- `Get-ManifestConfig` - File lists and target paths
- `Get-InstallationConfig` - Installation configuration
- `Get-FileDownloadUrl` - GitHub download URLs
- `ConvertTo-RelativePath` - Path utilities

#### 📁 FileOperations.psm1
- `Install-AllAIFiles` / `Remove-AllAIFiles` - Bulk operations
- `Install-AIFile` / `Remove-AIFile` - Individual file operations
- `Get-FileFromGitHub` - Remote file downloads
- `Update-GitIgnore` - Git ignore management
- `Test-FileIntegrity` - File validation

#### ✅ ValidationEngine.psm1
- `Test-SystemRequirements` - Validate prerequisites
- `Test-PreInstallation` / `Test-PostInstallation` - Installation checks
- `Test-WorkspaceValid` - Workspace validation
- `Get-ValidationReport` - Comprehensive status reports

#### 🎨 UI.psm1
- `Write-Success` / `Write-Warning` / `Write-Error` - Status messages
- `Show-Help` - Interactive help system
- `Confirm-UserAction` - User confirmations
- `Show-CompletionSummary` - Installation summaries

## 🎯 AI Prompt Usage Patterns

### Quick Start with Prompts

Once the AI infrastructure is installed, you can leverage the powerful prompt templates for common development tasks:

#### 🔍 Code Review Workflow
```
# Review your current changes before committing
/code-review-local-changes

# Review specific committed changes
/code-review-committed-changes for commit abc123
```

#### 🧪 Testing Workflow
```
# Add comprehensive unit tests to a file
/add-unit-tests for internal/services/cdn/frontdoor_resource.go

# Generate tests following Terraform provider patterns
/add-unit-tests with both typed and untyped resource patterns
```

#### 🏗️ Environment Setup
```
# Set up development environment
/setup-go-dev-environment

# Configure Go toolchain for Terraform development
/setup-go-dev-environment with Go 1.22+ and required tools
```

#### 📊 Repository Analysis
```
# Get high-level repository overview
/summarize-repo

# Deep dive into architecture and patterns
/summarize-repo-deep-dive
```

### Advanced Prompt Techniques

#### 🎯 Context-Specific Usage
Combine slash commands with specific context for better results:

```
# Review specific Azure service implementation
/code-review-local-changes focusing on Azure CDN Front Door patterns

# Generate tests for Azure API integration
/add-unit-tests for Azure SDK client code with proper mocking patterns

# Set up environment for specific Azure services
/setup-go-dev-environment and include Azure CLI tools for testing
```

#### 🔄 Iterative Development
Use slash commands in sequence for complete development workflows:

```
1. /summarize-repo to understand the codebase
2. # Make your code changes
3. /add-unit-tests for the changed files
4. /code-review-local-changes before committing
```

#### 🎨 Custom Prompt Combinations
Combine multiple commands for complex tasks:

```
# Comprehensive development review
/code-review-local-changes AND ensure the code follows Azure patterns from .github/instructions/

# Full testing coverage
/add-unit-tests and ensure tests follow .github/instructions/testing-guidelines.instructions.md
```

## 🔍 Usage Examples

### Check Current Status
```powershell
# From local installer (auto-detects repository)
.\install-copilot-setup.ps1 -Verify

# From user profile installer (must specify repository)
& "$env:USERPROFILE\.terraform-ai-installer\install-copilot-setup.ps1" -Verify -RepoDirectory "C:\github.com\hashicorp\terraform-provider-azurerm"
```

### Clean Installation
```powershell
# From local installer (auto-detects repository)
.\install-copilot-setup.ps1 -Clean

# From user profile installer (must specify repository)
& "$env:USERPROFILE\.terraform-ai-installer\install-copilot-setup.ps1" -Clean -RepoDirectory "C:\github.com\hashicorp\terraform-provider-azurerm"

# Remove without prompts
.\install-copilot-setup.ps1 -Clean -Auto-Approve

# Preview what would be removed (with RepoDirectory)
& "$env:USERPROFILE\.terraform-ai-installer\install-copilot-setup.ps1" -Clean -Dry-Run -RepoDirectory "C:\github.com\hashicorp\terraform-provider-azurerm"
```

### Bootstrap for Multiple Feature Branches
```powershell
# One-time setup from source branch
.\install-copilot-setup.ps1 -Bootstrap

# Then use from any feature branch (MUST specify RepoDirectory)
& "$env:USERPROFILE\.terraform-ai-installer\install-copilot-setup.ps1" -RepoDirectory "C:\github.com\hashicorp\terraform-provider-azurerm"

# For different repository locations
& "$env:USERPROFILE\.terraform-ai-installer\install-copilot-setup.ps1" -RepoDirectory "C:\your-projects\another-azurerm-fork"
```

## 🛠️ Troubleshooting

### Common Issues

#### ❌ "Bootstrap can only be run from source branch"
**Solution**: Switch to `exp/terraform_copilot` branch before running bootstrap.

#### ❌ "Clean operation not available on source branch"
**Solution**: Switch to a feature branch before running clean operations.

#### ❌ "DIRECTORY NOT FOUND: The specified RepoDirectory does not exist"
**Solution**: 
- Check the path spelling and ensure it exists
- Use an absolute path (e.g., `C:\path\to\repo`)
- Ensure you have permissions to access the directory

#### ❌ "INVALID REPOSITORY: The specified directory does not appear to be a terraform-provider-azurerm repository"
**Solution**:
- Ensure you're pointing to the repository ROOT directory
- Verify the directory contains `go.mod`, `main.go`, and `internal/` folder
- Example: `-RepoDirectory 'C:\github.com\hashicorp\terraform-provider-azurerm'`

#### ❌ Module import errors
**Solution**: Ensure you're running from the correct directory with all PowerShell modules present.

### Manual Recovery

If the installer becomes corrupted:

```powershell
# Re-download from source branch
git switch exp/terraform_copilot
git pull origin exp/terraform_copilot
.\install-copilot-setup.ps1 -Bootstrap

# Then return to your branch
git switch your-feature-branch
```

## 🎯 Advanced Usage

### Automation Scripts

```powershell
# Automated setup for CI/CD
$params = @{
    'Auto-Approve' = $true
    'Verify' = $true
}
& .\install-copilot-setup.ps1 @params
```

### Custom Installation Paths

The installer respects environment variables:
- `$env:USERPROFILE\.terraform-ai-installer` - Bootstrap location
- Target paths are defined in the configuration manifest

### Integration with Other Tools

```powershell
# Check if AI infrastructure is installed (returns boolean)
$result = .\install-copilot-setup.ps1 -Verify
$isInstalled = $result.Success -and ($result.Issues.Count -eq 0)
```

## 📚 Additional Resources

- **[Implementation Guide](.github/instructions/implementation-guide.instructions.md)** - Complete coding standards
- **[Testing Guidelines](.github/instructions/testing-guidelines.instructions.md)** - Testing requirements
- **[Azure Patterns](.github/instructions/azure-patterns.instructions.md)** - Azure-specific development patterns
- **[Troubleshooting Guide](.github/instructions/troubleshooting-decision-trees.instructions.md)** - Debugging workflows

## 🤝 Contributing

### Development Setup

1. Fork the repository
2. Switch to `exp/terraform_copilot` branch
3. Make changes to installer files
4. Test with `-Bootstrap` and feature branch installation
5. Submit pull request

### Testing Changes

```powershell
# Test bootstrap functionality
.\install-copilot-setup.ps1 -Bootstrap

# Test installation on clean feature branch
git checkout -b test/installer-changes
& "$env:USERPROFILE\.terraform-ai-installer\install-copilot-setup.ps1" -Dry-Run

# Test cleanup
& "$env:USERPROFILE\.terraform-ai-installer\install-copilot-setup.ps1" -Clean -Dry-Run
```

## 📄 License

This installer is part of the Terraform AzureRM Provider and follows the same licensing terms.

---

**Need help?** Run `.\install-copilot-setup.ps1 -Help` for interactive assistance.
