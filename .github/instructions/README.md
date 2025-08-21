# 🚀 Terraform AzureRM Provider - AI-Powered Development Guidelines

> **💡 TL;DR - Best Practice**: Use instruction files directly from this repository! Add them to your `.vscode/settings.json` instead of copying to user directories. This keeps instructions in sync with your code and eliminates branch management overhead.

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
> **ℹ️ NOTE**: While Method 1 (Repository-Based) is typically recommended, **Method 3 (Complete Local Installation)** is currently the preferred option until the AI instructions are merged into the `main` branch.

#### **✅ Method 1: Repository-Based (RECOMMENDED)**
---
**Use instruction files directly from the repository** - No copying needed!

**Advantages:**
- Always up-to-date with latest changes
- No sync issues between branches
- Instructions evolve with codebase
- Easier collaboration

**Setup:**
1. Simply use the instruction, prompt, and settings files from the repository
2. Start coding with AI assistance

#### **✅ Method 2: Dedicated AI Repository (Multi-Project)**
---
**For teams working across multiple Terraform providers:**

```bash
# Create dedicated instruction repository
git clone https://github.com/yourorg/terraform-ai-instructions.git

# Structure for multiple providers:
terraform-ai-instructions/
├── azurerm/           # Azure RM provider instructions
├── aws/               # AWS provider instructions
├── google/            # Google provider instructions
└── shared/            # Common Terraform patterns
```

#### **✅ Method 3: Complete Local Installation**
---
> **🌟 RECOMMENDED FOR NOW**: This is the preferred method until the AI instructions are merged into the main branch. It provides the most complete and immediate setup experience.

**Full local installation with automated scripts and comprehensive AI setup:**

**🎯 What It Does:**
- Installs all AI instructions and prompts directly to your VS Code user directory
- Configures local VS Code settings with intelligent merge of existing configuration
- Creates automatic backups and provides clean uninstall capability
- Works across all your projects (not just this repository)

**📂 Installation Scripts:**
- **Windows**: `.\AILocalInstall\install-copilot-setup.ps1`
- **Linux/macOS/WSL**: `./AILocalInstall/install-copilot-setup.sh`

**📋 Key Benefits:**
- **🔍 Auto-discovery**: Automatically finds repository location
- **🛡️ Safe setup**: Smart backups and clean uninstall
- **🔀 Intelligent merge**: Preserves your existing VS Code settings
- **🧹 Easy cleanup**: Complete removal with original settings restoration

**📖 For detailed usage, features, and troubleshooting:**
➡️ **See [`../AILocalInstall/README.md`](../AILocalInstall/README.md)** for complete installation guide

---
### ⚙️ **VS Code Configuration**
---

✅ **Automatic Setup**: The installation scripts automatically copy the repository's VS Code settings which include:

- **Auto-loading instruction files** during code review
- **Optimized AI response settings** (temperature, length)
- **Copilot enabled** for Go, Terraform, Markdown files
- **Smart commit message generation** following provider conventions

The settings are copied directly from `.vscode/settings.json` in the repository to ensure they're always up-to-date.

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
| `/summarize-repo` | Repository overview | Understanding structure |

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

1. **First Time Setup**: Run the installation script for your platform
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

# Get repository overview
/summarize-repo

# Environment setup help
/setup-go-dev-environment
```

## 🤖 **AI Prompts & Workflows**

### **Available AI Prompts**
---

| Prompt | Purpose | Usage |
|--------|---------|-------|
| `/code-review-local-changes` | Review uncommitted changes for compliance | Type in Copilot Chat |
| `/code-review-committed-changes` | Review commits and PRs for standards | Type in Copilot Chat |
| `/setup-go-dev-environment` | Complete development environment setup | Type in Copilot Chat |
| `/summarize-repo` | High-level repository overview | Type in Copilot Chat |
| `/summarize-repo-deep-dive` | Detailed technical analysis | Type in Copilot Chat |

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

#### Deep dive analysis:
* `/summarize-repo-deep-dive`

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
        "text": "Provide a concise and clear commit message that summarizes the changes made in the code. For complex changes, include the following details: 1) Specify if the change introduces a breaking change and describe its impact. 2) Highlight any new resources or features added. 3) Mention updates to Azure services or APIs. Aim to keep the message under 72 characters per line for readability."
    }
]
```

**Generated Commit Examples**:
```
feat: add Azure CDN Front Door Profile log scrubbing support

- Add log_scrubbing_rule schema with QueryStringArgNames, RequestIPAddress, RequestUri support
- Implement expand/flatten functions following Azure API patterns
- Add comprehensive acceptance tests with CustomizeDiff validation
- Update documentation with examples and proper note formatting

BREAKING CHANGE: log_scrubbing configuration moved from nested block to direct rules
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
- [ ] Copy instruction files to user directory
- [ ] Configure VS Code settings for AI integration
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

## Quick Access to Development Guidelines

Welcome to the AI-enhanced instruction file system for the Terraform AzureRM provider. These guides provide comprehensive development guidance optimized for both human contributors and AI assistance.

### 🎯 **Core Development Guides** (Start Here)

| Guide | Purpose | Key Content |
|-------|---------|-------------|
| [🏗️ **Implementation Guide**](./implementation-guide.instructions.md) | **Primary reference for all coding standards, patterns, and style** | Unified coding standards, typed vs untyped patterns, naming conventions, file organization, Azure SDK integration |
| [📋 **Code Clarity Enforcement**](./code-clarity-enforcement.instructions.md) | **Code clarity and policy enforcement guidelines** | Zero tolerance comment policy, strategic decision-making, CustomizeDiff requirements, quality standards |
| [☁️ **Azure Patterns**](./azure-patterns.instructions.md) | **Azure-specific implementation patterns and best practices** | PATCH operations, CustomizeDiff validation, schema flattening, security patterns, "None" value handling |
| [❌ **Error Patterns**](./error-patterns.instructions.md) | **Comprehensive error handling and debugging guidelines** | Error message formatting, debugging protocols, Azure API errors, console output interpretation |
| [📐 **Schema Patterns**](./schema-patterns.instructions.md) | **Schema design patterns and validation standards** | Field types, validation functions, Azure helpers, complex schemas, TypeSet vs TypeList |

### 🔧 **Specialized Development Guides**

| Guide | Purpose | Key Content |
|-------|---------|-------------|
| [🔄 **Migration Guide**](./migration-guide.instructions.md) | **Implementation approach transitions and upgrade procedures** | Typed resource migration, breaking changes, version compatibility, upgrade procedures |
| [🧪 **Testing Guidelines**](./testing-guidelines.instructions.md) | **Testing standards and patterns for acceptance and unit tests** | Test execution protocols, CustomizeDiff testing, Azure resource testing, cleanup patterns |
| [📚 **Documentation Guidelines**](./documentation-guidelines.instructions.md) | **Documentation standards for resources and data sources** | Resource vs data source patterns, example standards, field documentation |
| [🏢 **Provider Guidelines**](./provider-guidelines.instructions.md) | **Azure-specific provider patterns and integration guidelines** | ARM integration, client management, Azure service constraints |

### 🚀 **Enhanced Guidance Files**

| Guide | Purpose | Key Content |
|-------|---------|-------------|
| [🔄 **API Evolution**](./api-evolution-patterns.instructions.md) | **API evolution and versioning patterns** | Version management, backward compatibility, migration strategies, deprecation management |
| [⚡ **Performance**](./performance-optimization.instructions.md) | **Performance optimization and efficiency guidelines** | Azure API efficiency, resource management optimization, monitoring patterns, scalability |
| [🔐 **Security**](./security-compliance.instructions.md) | **Security and compliance patterns** | Input validation, credential management, security best practices, compliance requirements |
| [🔧 **Troubleshooting**](./troubleshooting-decision-trees.instructions.md) | **Troubleshooting decision trees and diagnostic patterns** | Common issues resolution, root cause analysis, error diagnostics, state management troubleshooting |

---
[⬆️ Back to top](#🚀-terraform-azurerm-provider---development-guidelines)

## 🎓 **How to Use This System**

### For New Developers
1. **Start with**: [🏗️ Implementation Guide](./implementation-guide.instructions.md) - Get familiar with coding standards and patterns
2. **Understand policy enforcement**: [📋 Code Clarity Enforcement](./code-clarity-enforcement.instructions.md) - Learn critical comment policies and quality standards
3. **Learn Azure specifics**: [☁️ Azure Patterns](./azure-patterns.instructions.md) - Understand Azure-specific implementation requirements
4. **Master error handling**: [❌ Error Patterns](./error-patterns.instructions.md) - Learn proper error handling and debugging techniques
5. **Schema design**: [📐 Schema Patterns](./schema-patterns.instructions.md) - Understand schema design patterns and validation

### For Experienced Developers
- **Quick Reference**: Use the emoji navigation within each file to jump between related sections
- **Specific Tasks**: Use the purpose column above to find the most relevant guide
- **Migration Work**: Start with [🔄 Migration Guide](./migration-guide.instructions.md) for implementation approach changes
- **Policy Enforcement**: Reference [📋 Code Clarity Enforcement](./code-clarity-enforcement.instructions.md) for code review standards
- **Performance Issues**: Use [⚡ Performance](./performance-optimization.instructions.md) for optimization guidance
- **Security Reviews**: Reference [🔐 Security](./security-compliance.instructions.md) for security pattern compliance
- **Troubleshooting**: Use [🔧 Troubleshooting](./troubleshooting-decision-trees.instructions.md) for systematic issue resolution

### For Code Reviews
- **Comment Policy**: [📋 Code Clarity Enforcement](./code-clarity-enforcement.instructions.md) - **CRITICAL**: Zero tolerance for unnecessary comments
- **Standards Check**: [🏗️ Implementation Guide](./implementation-guide.instructions.md) for coding standards compliance
- **Azure Compliance**: [☁️ Azure Patterns](./azure-patterns.instructions.md) for Azure-specific pattern verification
- **Error Handling**: [❌ Error Patterns](./error-patterns.instructions.md) for proper error handling review

---
[⬆️ Back to top](#🚀-terraform-azurerm-provider---development-guidelines)

## 🚀 **Next Steps**

1. **Bookmark this file** as your starting point for development guidance
2. **Review comment policy FIRST**: [📋 Code Clarity Enforcement](./code-clarity-enforcement.instructions.md) - Understanding the zero tolerance comment policy is critical
3. **Use the emoji navigation** within each file to quickly find related information
4. **Contribute improvements** by following the patterns established in these guides
5. **Report issues** if you find gaps or inconsistencies in the guidance

---
[⬆️ Back to top](#🚀-terraform-azurerm-provider---development-guidelines)

## 🎯 **Critical Policy Reminder**

**⚠️ ZERO TOLERANCE FOR UNNECESSARY COMMENTS** - Before writing ANY comment, review the [Code Clarity Enforcement Guidelines](./code-clarity-enforcement.instructions.md#🚫-zero-tolerance-for-unnecessary-comments-policy). Comments are only allowed for Azure API quirks, complex business logic, SDK workarounds, or non-obvious state patterns.

---

## Quick Reference Links

- ☁️ **Azure Patterns**: [azure-patterns.instructions.md](./azure-patterns.instructions.md)
- 📋 **Code Clarity Enforcement**: [code-clarity-enforcement.instructions.md](./code-clarity-enforcement.instructions.md)
- 📝 **Documentation Guide**: [documentation-guidelines.instructions.md](./documentation-guidelines.instructions.md)
- ❌ **Error Patterns**: [error-patterns.instructions.md](./error-patterns.instructions.md)
- 🏗️ **Implementation Guide**: [implementation-guide.instructions.md](./implementation-guide.instructions.md)
- 🔄 **Migration Guide**: [migration-guide.instructions.md](./migration-guide.instructions.md)
- 🏢 **Provider Guidelines**: [provider-guidelines.instructions.md](./provider-guidelines.instructions.md)
- 📐 **Schema Patterns**: [schema-patterns.instructions.md](./schema-patterns.instructions.md)
- 🧪 **Testing Guide**: [testing-guidelines.instructions.md](./testing-guidelines.instructions.md)

### 🚀 Enhanced Guidance Files

- 🔄 **API Evolution**: [api-evolution-patterns.instructions.md](./api-evolution-patterns.instructions.md)
- ⚡ **Performance**: [performance-optimization.instructions.md](./performance-optimization.instructions.md)
- 🔐 **Security**: [security-compliance.instructions.md](./security-compliance.instructions.md)
- 🔧 **Troubleshooting**: [troubleshooting-decision-trees.instructions.md](./troubleshooting-decision-trees.instructions.md)

---

*This instruction system provides operational guidance for Terraform AzureRM provider development. Each file includes cross-references to related content and clear navigation paths.*

---
[⬆️ Back to top](#🚀-terraform-azurerm-provider---development-guidelines)
