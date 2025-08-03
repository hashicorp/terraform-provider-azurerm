# 🚀 Terraform AzureRM Provider - Development Guidelines

**Quick navigation:** [🎯 Core Guides](#🎯-core-development-guides-start-here) | [🔧 Specialized Guides](#🔧-specialized-development-guides) | [🎓 How to Use](#🎓-how-to-use-this-system) | [🚀 Next Steps](#🚀-next-steps) | [⚠️ Critical Policy](#🎯-critical-policy-reminder)

## Quick Access to Development Guidelines

Welcome to the instruction file system for the Terraform AzureRM provider. These guides provide comprehensive development guidance for contributors and AI assistance.

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
5. **Schema design**: [� Schema Patterns](./schema-patterns.instructions.md) - Understand schema design patterns and validation

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
