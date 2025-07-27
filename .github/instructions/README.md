# 🚀 Terraform AzureRM Provider - Development Guidelines

## Quick Access to Development Guidelines

Welcome to the instruction file system for the Terraform AzureRM provider. These guides provide comprehensive development guidance for contributors and AI assistance.

### 🎯 **Core Development Guides** (Start Here)

| Guide | Purpose | Key Content |
|-------|---------|-------------|
| [🏗️ **Implementation Guide**](./implementation-guide.md) | **Primary reference for all coding standards, patterns, and style** | Unified coding standards, typed vs untyped patterns, naming conventions, file organization, Azure SDK integration |
| [🔷 **Azure Patterns**](./azure-patterns.md) | **Azure-specific implementation patterns and best practices** | PATCH operations, CustomizeDiff validation, schema flattening, security patterns, "None" value handling |
| [❌ **Error Patterns**](./error-patterns.md) | **Comprehensive error handling and debugging guidelines** | Error message formatting, debugging protocols, Azure API errors, console output interpretation |
| [📋 **Schema Patterns**](./schema-patterns.md) | **Schema design patterns and validation standards** | Field types, validation functions, Azure helpers, complex schemas, TypeSet vs TypeList |

### 🔧 **Specialized Development Guides**

| Guide | Purpose | Key Content |
|-------|---------|-------------|
| [🔄 **Migration Guide**](./migration-guide.md) | **Implementation approach transitions and upgrade procedures** | Typed resource migration, breaking changes, version compatibility, upgrade procedures |
| [🧪 **Testing Guidelines**](./testing-guidelines.instructions.md) | **Testing standards and patterns for acceptance and unit tests** | Test execution protocols, CustomizeDiff testing, Azure resource testing, cleanup patterns |
| [📚 **Documentation Guidelines**](./documentation-guidelines.instructions.md) | **Documentation standards for resources and data sources** | Resource vs data source patterns, example standards, field documentation |
| [🏢 **Provider Guidelines**](./provider-guidelines.instructions.md) | **Azure-specific provider patterns and integration guidelines** | ARM integration, client management, Azure service constraints |

## 🎓 **How to Use This System**

### For New Developers
1. **Start with**: [🏗️ Implementation Guide](./implementation-guide.md) - Get familiar with coding standards and patterns
2. **Learn Azure specifics**: [🔷 Azure Patterns](./azure-patterns.md) - Understand Azure-specific implementation requirements
3. **Master error handling**: [❌ Error Patterns](./error-patterns.md) - Learn proper error handling and debugging techniques
4. **Schema design**: [📋 Schema Patterns](./schema-patterns.md) - Understand schema design patterns and validation

### For Experienced Developers
- **Quick Reference**: Use the emoji navigation within each file to jump between related sections
- **Specific Tasks**: Use the purpose column above to find the most relevant guide
- **Migration Work**: Start with [🔄 Migration Guide](./migration-guide.md) for implementation approach changes

### For Code Reviews
- **Standards Check**: [🏗️ Implementation Guide](./implementation-guide.md) for coding standards compliance
- **Azure Compliance**: [🔷 Azure Patterns](./azure-patterns.md) for Azure-specific pattern verification
- **Error Handling**: [❌ Error Patterns](./error-patterns.md) for proper error handling review

## 🚀 **Next Steps**

1. **Bookmark this file** as your starting point for development guidance
2. **Use the emoji navigation** within each file to quickly find related information
3. **Contribute improvements** by following the patterns established in these guides
4. **Report issues** if you find gaps or inconsistencies in the guidance

---

*This instruction system provides operational guidance for Terraform AzureRM provider development. Each file includes cross-references to related content and clear navigation paths.*
