# 🚀 Terraform AzureRM Provider - AI-Powered Development Guidelines

> **💡 TL;DR - Best Practice**: Until this PR is merged to `main`, use **Method 2 (Automated Installer Setup)** for the most complete setup. Once merged, you'll be able to use instruction files directly from the repository without any copying.

**Quick navigation:** [🤖 AI Setup & Installation](#🤖-ai-setup--installation) | [🎯 Core Guides](#🎯-core-development-guides-start-here) | [🔧 Specialized Guides](#🔧-specialized-development-guides) | [🎓 How to Use](#🎓-how-to-use-this-system) | [🚀 Next Steps](#🚀-next-steps) | [⚠️ Critical Policy](#🎯-critical-policy-reminder)

## 🤖 **AI Setup & Installation**

This instruction system enables **AI-powered "Vibe Coding"** for generating Resources, Tests, and Documentation in the Terraform AzureRM Provider. Follow these steps to supercharge your development workflow with AI assistance.

---
### 📋 **Prerequisites**
---

**Required VS Code Extensions:**
- **GitHub Copilot** (GitHub.copilot) - Essential for AI code generation
- **GitHub Copilot Chat** (GitHub.copilot-chat) - Required for instruction-aware conversations
- **HashiCorp Terraform** (HashiCorp.terraform) - Terraform language support
- **Go** (golang.Go) - Go language support and debugging

**Optional but Recommended:**
- **Azure Account** (ms-vscode.azure-account) - Azure integration
- **Error Lens** (usernamehw.errorlens) - Inline error highlighting
- **GitLens** (eamodio.gitlens) - Enhanced Git capabilities

---
### 🎯 **Setup Methods (Choose One)**
---
> **🌟 RECOMMENDED FOR NOW**: **Method 2 (Automated Installer Setup)** is the preferred approach until this PR is merged to `main`. Once merged, `Method 1 (Repository-Based)` will become the recommended option.

#### **✅ Method 1: Repository-Based (FUTURE RECOMMENDED)**
---
**Use instruction files directly from the repository** - No copying needed!

> **⚠️ AVAILABILITY**: This method will be available once this PR is merged to the `main` branch.

**Advantages:**
- Always up-to-date with latest changes
- No sync issues between branches  
- Instructions evolve with codebase
- Easier collaboration
- Works directly from repository

**Setup:**
1. Simply use the instruction, prompt, and settings files from the repository
2. Start coding with AI assistance

#### **✅ Method 2: Automated Installer Setup (RECOMMENDED FOR NOW)**
---
> **🌟 CURRENT BEST OPTION**: This is the preferred method until the AI instructions are merged to the `main` branch. It provides automated setup using the installer script.

**Automated installation script that sets up AI instructions from this feature branch:**

**🎯 What It Does:**
- Downloads and installs AI instructions from this feature branch to your VS Code user directory
- Configures local VS Code settings with intelligent merge of existing configuration
- Creates automatic backups and provides clean uninstall capability
- Works across all your projects (pulls latest instructions from this branch)

**📂 Installation Scripts:**
- **Windows**: `.\.github\AIinstaller\install-copilot-setup.ps1`

**📋 Key Benefits:**
- **🔍 Branch-aware**: Automatically pulls instructions from the correct source branch
- **🛡️ Safe setup**: Smart file copying with verification  
- **🔀 Works everywhere**: Bootstrap once, works across all projects
- **🧹 Easy cleanup**: Complete removal and restoration
- **🚀 Ready now**: Available immediately without waiting for PR merge

**📖 For detailed usage, features, and troubleshooting:**
➡️ **See [`../AIinstaller/README.md`](../AIinstaller/README.md)** for complete installation guide

---
### ⚙️ **VS Code Configuration**
---

✅ **Automatic Setup**: The installation script automatically copies the repository's optimized VS Code settings which include:

- **Auto-loading instruction files** during code review
- **Optimized AI response settings** (temperature, length)  
- **Copilot enabled** for Go, Terraform, Markdown files
- **Smart commit message generation** following provider conventions

The settings are copied directly from `.vscode/settings.json` in this branch to ensure they're always up-to-date with the latest AI configuration.

---
### 🚀 **Quick Start Guide**
---

**Basic Workflow:**
1. **Ask AI for help**: "Generate a new Azure [ServiceName] resource using typed implementation"
2. **Review code**: Select generated code → use `/code-review-local-changes` prompt
3. **Generate docs**: "Create documentation following provider guidelines"
4. **Final review**: "/code-review-local-changes" before committing

---
### 🎯 **Essential AI Prompts**
---

The repository includes specialized prompt files for common development tasks:

| Prompt | Purpose | Usage |
|--------|---------|-------|
| `/code-review-local-changes` | Review uncommitted changes | Before committing |
| `/code-review-committed-changes` | Review committed changes | PR review |

**Example Usage:**
```
# In GitHub Copilot Chat
/code-review-local-changes

# Or ask directly
"How do I implement a new Azure CDN resource following provider patterns?"
```

### 💡 **Best Practices**

- **Be specific**: Include Azure service names for better context
- **Use code selection**: Select relevant code before asking for reviews
- **Reference guidelines**: Mention specific instruction files when needed
- **Test locally**: Always validate AI-generated code before committing

---

## 🎯 **Core AI Development Guides**

### **🤖 Essential AI Instruction Files**

| File | Purpose | AI Loading Behavior |
|------|---------|---------------------|
| **`copilot-instructions.md`** | Main AI behavior and provider guidelines | Always loaded - provides base context |
| **`implementation-guide.instructions.md`** | Complete Go implementation patterns | Auto-loaded for resource implementation tasks |
| **`azure-patterns.instructions.md`** | Azure-specific PATCH operations, CustomizeDiff patterns | Auto-loaded for Azure API integration work |
| **`testing-guidelines.instructions.md`** | Test execution protocols and patterns | Auto-loaded for testing and review tasks |
| **`documentation-guidelines.instructions.md`** | Resource and data source documentation standards | Auto-loaded for documentation tasks |
| **`provider-guidelines.instructions.md`** | Azure Resource Manager integration best practices | Auto-loaded for ARM-specific functionality |

### **🎯 Specialized AI Instruction Files**

| File | Purpose | AI Loading Behavior |
|------|---------|---------------------|
| **`error-patterns.instructions.md`** | Error handling patterns and debugging guidelines | Loaded on-demand for debugging tasks |
| **`schema-patterns.instructions.md`** | Schema design patterns and validation standards | Loaded on-demand for complex schema work |
| **`code-clarity-enforcement.instructions.md`** | Code clarity and quality enforcement guidelines | Loaded on-demand for code quality enforcement |
| **`migration-guide.instructions.md`** | Migration patterns and upgrade procedures | Loaded on-demand for migration tasks |
| **`api-evolution-patterns.instructions.md`** | API evolution and versioning patterns | Loaded on-demand for API compatibility work |
| **`performance-optimization.instructions.md`** | Performance optimization patterns and efficiency guidelines | Loaded on-demand for optimization tasks |
| **`security-compliance.instructions.md`** | Security and compliance patterns | Loaded on-demand for security implementation |
| **`troubleshooting-decision-trees.instructions.md`** | Troubleshooting decision trees and diagnostic patterns | Loaded on-demand for diagnostic workflows |

---
### **🚀 Quick Start for New Developers**
---

1. **First Time Setup**: Run the automated installer script (Method 2) for your platform
2. **Understand the Provider**: Ask AI "Summarize the terraform-provider-azurerm repository structure"  
3. **Learn Patterns**: Start with `/code-review-local-changes` on existing code to see standards
4. **Practice**: Try "Generate a simple Azure resource following typed implementation patterns"

---
### **🔧 Most Common AI Commands**
---
```bash
# Review any changes before committing
/code-review-local-changes

# Review committed code and PRs for standards
/code-review-committed-changes
```

## 🤖 **AI Prompts & Workflows**

### **Available AI Prompts**
---

| Prompt | Purpose | Usage |
|--------|---------|-------|
| `/code-review-local-changes` | Review uncommitted changes for compliance | Type in Copilot Chat |
| `/code-review-committed-changes` | Review commits and PRs for standards | Type in Copilot Chat |

---

### ⚡ **Resource Development Workflow**:
---

#### Step 1: Initialize new resource
* **Copilot Chat**: "Create a new Azure [`ServiceName`] resource using typed implementation"

#### Step 2: Review implementation
`/code-review-local-changes`

#### Step 3: Generate tests
* **Copilot Chat**: "Create comprehensive acceptance tests for azurerm_[`resource_name`]"

#### Step 4: Generate documentation
* **Copilot Chat**: "Create documentation following provider standards for azurerm_[`resource_name`]"

#### Step 5: Final review
* `/code-review-local-changes`

---

### 🐛 **Bug Fix Workflow**:
---

#### Step 1: Analyze issue
* **Copilot Chat**: "Help me understand this Azure API error and suggest a fix"

#### Step 2: Validate fix
* **Copilot Chat**: "Review my bug fix for provider standards compliance"

#### Step 3: Add regression tests
* **Copilot Chat**: "Create tests to prevent regression of this Azure resource issue"

#### Step 4: Use specific prompt
* `/code-review-local-changes`

---

### 🚀 **Feature Enhancement Workflow**:
---
#### Step 1: Plan enhancement
* **Copilot Chat**: "How should I add [`feature`] to azurerm_[`resource`] following Azure patterns?"

#### Step 2: Implementation guidance
* **Copilot Chat**: "Guide me through implementing this Azure feature with proper CustomizeDiff validation"

#### Step 3: Code review
* `/code-review-local-changes`

#### Step 4: Documentation
* **Copilot Chat**: "Update documentation for this new Azure feature following provider guidelines"

---

### 🎓 **Learning & Exploration**:

#### Repository overview:
* **Copilot Chat**: "Provide a summary of the Terraform AzureRM provider architecture"

#### Pattern understanding:
* **Copilot Chat**: "Explain the difference between typed and untyped resource implementations"

#### Azure service integration:
* **Copilot Chat**: "How do I integrate with Azure [`ServiceName`] APIs following provider patterns?"

---
## 🎯 **Advanced AI Integration Features**

### **Automatic Context Loading**:
The VS Code settings automatically include instruction files when reviewing code:

```jsonc
"github.copilot.chat.reviewSelection.instructions": [
    {"file": ".github/copilot-instructions.md"},
    {"file": ".github/instructions/implementation-guide.instructions.md"},
    {"file": ".github/instructions/azure-patterns.instructions.md"},
    {"file": ".github/instructions/testing-guidelines.instructions.md"},
    {"file": ".github/instructions/documentation-guidelines.instructions.md"},
    {"file": ".github/instructions/provider-guidelines.instructions.md"}
]
```

#### **How It Works**:
1. Select any code in VS Code
2. Right-click and choose "Copilot: Review Selection" **OR** use Copilot Chat directly with natural language, prompts (like `/code-review-local-changes`), or both
3. AI automatically loads all provider guidelines
4. Provides context-aware feedback following Azure provider standards

---
### **Smart Commit Message Generation**:
AI generates commit messages following Azure provider conventions:

```jsonc
"github.copilot.chat.commitMessageGeneration.instructions": [
    {
        "text": "Write clear, concise commit messages under 72 characters. Focus on what changed, not why. No need for conventional commit prefixes unless it's a breaking change."
    }
]
```

**Generated Commit Examples**:
```
add Azure CDN Front Door Profile log scrubbing support

update VMSS resiliency policy validation logic

fix authentication handling for managed identities

BREAKING CHANGE: remove deprecated log_scrubbing wrapper
```

---
### **File Type Optimization & Terminal AI Integration**
---
Optimized AI assistance for all relevant file types and contexts:

```jsonc
"github.copilot.enable": {
    "*": true,         // Enable for all file types (Go, Terraform, YAML, Markdown, etc.)
    "terminal": true   // Enable Copilot assistance directly in the terminal
}
```

**Use Cases:**
- Get Azure CLI command suggestions
- Debug Azure API errors
- Generate PowerShell scripts for Azure operations
- Troubleshoot authentication issues

---
### **Response Quality Tuning**
---
Optimized for detailed, accurate responses:

```jsonc
"github.copilot.advanced": {
    "length": 3000,    // Longer responses for comprehensive guidance
    "temperature": 0.1 // Lower temperature for focused, accurate suggestions
}
```
### **File Associations**
```jsonc
"files.associations": {
    "*.instructions.md": "markdown",  // Proper highlighting for instruction files
    ".github/*.md": "markdown"        // and GitHub documentation files
}
```

### **Privacy & Security Settings**
```jsonc
"github.copilot.chat.summarizeAgentConversationHistory.enabled": false
```
- Prevents conversation history storage
- Ensures sensitive Azure configurations remain private
- Complies with enterprise security requirements

### **Development Optimizations**
```jsonc
"github.copilot.chat.reviewSelection.enabled": true
```
- Enables automatic code review with instruction files
- Provides consistent feedback across team members
- Enforces provider coding standards automatically

---
## 🔧 **Advanced Troubleshooting & Best Practices**
---

### **AI Response Quality Issues**
---

**Problem**: AI not following provider guidelines

**Solutions:**
1. **Repository-based approach**: Verify instruction files are properly referenced in `.vscode/settings.json`
2. **User directory approach**: Check files are in `$env:APPDATA\Code\User\instructions\terraform-azurerm\`
3. **Symbolic link approach**: Ensure links are properly created and accessible
4. Restart VS Code completely after setup changes
5. Check GitHub Copilot authentication: `Ctrl+Shift+P` → `GitHub Copilot: Sign In`

---
**Problem**: Inconsistent or generic responses
---

**Solutions:**
1. Use `@workspace` prefix for workspace-aware responses
2. Include specific context: "following Azure provider patterns"
3. Reference specific instruction files: "using implementation-guide.instructions.md patterns"
4. Use selection-based prompts for targeted feedback

---
**Problem**: AI generates code that doesn't compile
---

**Solutions:**
1. Always specify implementation approach: "using typed resource implementation"
2. Include Azure service context: "for Azure CDN Front Door service"
3. Request complete implementations: "include all necessary imports and error handling"
4. Use code review prompts to validate generated code

---
**Problem**: GitHub Copilot not suggesting provider-specific patterns
---

**Solutions:**
1. Verify VS Code settings include provider-specific configurations
2. Ensure file associations are correctly configured for `.instructions.md` files
3. Check that `reviewSelection.instructions` includes all instruction files
4. Validate Copilot is enabled for Go and Terraform file types

---
**Problem**: Commit message generation not working
---

**Solutions:**
1. Verify `commitMessageGeneration.instructions` is configured in settings
2. Ensure you're using VS Code's built-in source control commit interface
3. Check that commit message generation is enabled in Copilot settings
4. Test with a simple file change to verify functionality

---
#### **Performance Optimization**
---

**Best Practices for Large Codebase:**
1. Use `.gitignore` patterns in VS Code search to exclude vendor directories
2. Configure file watching exclusions for better performance
3. Use targeted searches with specific file patterns
4. Limit AI context to relevant files when possible

**Memory Management:**
1. Close unused editor tabs when working with large instruction files
2. Use VS Code's workspace folders for focused development
3. Configure reasonable limits for AI response length
4. Clear terminal output regularly during long development sessions

---
#### **Team Collaboration Guidelines**
---

**Shared Configuration:**
1. Commit `.vscode/settings.json` to repository for team consistency
2. Document any local-only settings that shouldn't be shared
3. Use workspace-level settings for team standards
4. Keep user-level settings for personal preferences

**Prompt Sharing:**
1. Create custom prompts in `.github/prompts/` for team use
2. Document prompt usage patterns in team documentation
3. Share successful AI interaction patterns with team members
4. Maintain prompt files alongside code changes

---
### 🎓 **Training & Onboarding**
---

#### **New Developer Checklist**
- [ ] Install required VS Code extensions
- [ ] Run automated installer script (Method 2) to set up AI instructions
- [ ] Test AI integration with simple prompt
- [ ] Review Azure provider architecture using repository analysis prompts
- [ ] Practice resource development workflow with AI assistance

#### **Advanced Developer Workflow**
- [ ] Create custom prompts for specialized workflows
- [ ] Optimize VS Code settings for performance
- [ ] Contribute improvements to instruction files
- [ ] Mentor new developers on AI-assisted development
- [ ] Share best practices and successful patterns with team

---
### 🎯 **Maximizing AI Effectiveness**
---

#### **Key Success Factors**
1. **Proper Setup**: Install instruction files and configure VS Code settings correctly
2. **Correct Syntax**: Use `/prompt-name` format for prompt invocation
3. **Context Awareness**: Select relevant code and provide specific Azure service context
4. **Iterative Approach**: Use multiple prompts in sequence for complex tasks
5. **Manual Review**: Always validate AI-generated code against provider standards

#### **AI + Human Collaboration Model**
- **AI Strengths**: Code generation, pattern recognition, comprehensive documentation
- **Human Oversight**: Business logic validation, Azure service expertise, code review
- **Best Results**: Combine AI efficiency with human knowledge and judgment

#### **Continuous Improvement**
- Contribute feedback on instruction file effectiveness
- Share successful AI interaction patterns with the team
- Update prompts based on evolving provider patterns
- Document new AI workflows that prove effective

---

*This instruction system provides operational guidance for Terraform AzureRM provider development. Each file includes cross-references to related content and clear navigation paths.*

---
[⬆️ Back to top](#🚀-terraform-azurerm-provider---development-guidelines)
