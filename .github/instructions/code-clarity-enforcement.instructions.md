---
applyTo: "internal/**/*.go"
description: Code clarity and policy enforcement guidelines for Terraform AzureRM provider Go files. Includes detailed rules for comments, imports, implementation patterns, and quality standards.
---

# Code Clarity and Policy Enforcement Guidelines

This document provides detailed enforcement guidelines for maintaining code clarity and quality standards in the Terraform AzureRM provider.

**Quick Navigation:** [🚫 Comment Policy](#-zero-tolerance-for-unnecessary-comments) | [🔄 CustomizeDiff](#customizediff-import-requirements) | [☁️ Azure Standards](#azure-api-integration-standards) | [🎯 Enforcement Priority](#enforcement-priority)

**Related Guidelines:**
- 🏗️ **Core Implementation**: [implementation-guide.instructions.md](./implementation-guide.instructions.md) - Main coding standards and patterns
- ☁️ **Azure Patterns**: [azure-patterns.instructions.md](./azure-patterns.instructions.md) - PATCH operations, CustomizeDiff validation, Azure-specific behaviors
- 🧪 **Testing Standards**: [testing-guidelines.instructions.md](./testing-guidelines.instructions.md) - Comprehensive test requirements and patterns

## 🎯 Strategic Decision-Making Guidance

**Implementation Context Awareness**: When making coding decisions during pair programming, always consider:

**1. Comment Policy Enforcement Priority**
- **Zero tolerance for unnecessary comments** - This is the highest priority enforcement guideline
- **Before ANY comment**: Ask whether code structure, naming, or extraction can eliminate the need
- **Exception criteria**: Only Azure API quirks, complex business logic, SDK workarounds, or non-obvious state management patterns

**2. Implementation Pattern Context**
- **Typed vs Untyped resources**: Apply same comment standards regardless of implementation approach
- **Azure service constraints**: Comments acceptable for Azure-specific behaviors that cannot be expressed through code structure
- **CustomizeDiff patterns**: Complex validation logic may require explanation of Azure API constraints

**3. Performance-Critical Decisions**
- **Code clarity over comments**: Always prefer refactoring to commenting
- **Cross-pattern consistency**: Ensure comment policies apply uniformly across resource variants (Linux/Windows VMSS, etc.)
- **Maintainability impact**: Favor self-documenting code patterns that reduce long-term maintenance burden

**4. Quality Gate Integration**
- **Pre-submission validation**: Every comment must have explicit justification documented in review response
- **Cross-file consistency**: Validate related implementations maintain identical comment policies
- **Azure API alignment**: Comments must reflect actual Azure service behavior, not implementation assumptions

## CustomizeDiff Import Requirements

**Conditional Import Pattern**: Import requirements depend on the CustomizeDiff implementation:

**Dual Imports Required**: When using *schema.ResourceDiff directly in CustomizeDiff functions:
- `github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema` (for *schema.ResourceDiff)
- `github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk` (for helpers)

**Single Import Sufficient**: When using *pluginsdk.ResourceDiff (legacy resources):
- `github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema` (only this import needed)

**Function Signature Analysis**: Check the function signature to determine import requirements:
- *schema.ResourceDiff → Usually typed resources, may need dual imports
- *pluginsdk.ResourceDiff → Usually legacy resources, single import sufficient

**Helper Function Usage**: Verify proper use of pluginsdk.All(), pluginsdk.ForceNewIfChange() helpers

**For comprehensive CustomizeDiff patterns and Azure validation examples, see:** [azure-patterns.instructions.md](./azure-patterns.instructions.md#customizediff-validation)

---
[⬆️ Back to top](#code-clarity-and-policy-enforcement-guidelines)

## Resource Implementation Standards

**CRUD Operations**: Ensure Create, Read, Update, Delete functions handle all edge cases

**Schema Validation**: Verify all required fields, validation functions, and type definitions

**ForceNew Logic**: Check that properties requiring resource recreation are properly marked

**Timeouts**: Ensure appropriate timeout values for Azure operations (often long-running)

---
[⬆️ Back to top](#code-clarity-and-policy-enforcement-guidelines)

## Code Comments Policy Enforcement

### 🚫 ZERO TOLERANCE for Unnecessary Comments

**MANDATORY comment criteria - comments ONLY allowed for**:
- Azure API-specific quirks or behaviors not obvious from code
- Complex business logic that cannot be made clear through code structure alone
- Workarounds for Azure SDK limitations or API bugs
- Non-obvious state management patterns (PATCH operations, residual state handling)
- Azure service constraints requiring explanation (timeout ranges, SKU limitations)

**For detailed Azure API patterns requiring comments, see:** [azure-patterns.instructions.md](./azure-patterns.instructions.md#patch-operations)

### 🚫 FORBIDDEN Comments - Flag These Immediately

**Never comment**:
- Variable assignments, struct initialization, basic operations
- Standard Terraform patterns (CRUD operations, schema definitions)
- Self-explanatory function calls or routine Azure API calls
- Field mappings between Terraform and Azure API models
- Obvious conditional logic or loops
- Standard Go patterns (error handling, nil checks, etc.)

### Comment Review Process

**JUSTIFICATION REQUIREMENT**: If ANY comment exists, the developer MUST provide explicit justification:
- Which exception case this comment falls under
- Why the code cannot be self-explanatory through better naming/structure
- What specific Azure API behavior requires documentation (if applicable)

**SUGGESTED ACTION**: When flagging unnecessary comments, suggest how to make code self-explanatory instead:
- Better variable naming
- Function extraction
- Structure reorganization
- Pattern clarification

### Comment Validation Questions

Before allowing any comment, ask:
1. "Is this code unclear without a comment?" → Refactor the code instead
2. "Would a developer be confused by this logic?" → Only then consider a comment
3. "Is this documenting an Azure API quirk?" → Comment may be acceptable

---
[⬆️ Back to top](#code-clarity-and-policy-enforcement-guidelines)

## Azure API Integration Standards

**Error Handling**: Verify proper handling of Azure API errors, including 404s during Read operations

**Polling**: Check for proper implementation of long-running operation polling

**API Versions**: Ensure correct and consistent Azure API versions are used

**Authentication**: Verify proper use of Azure client authentication patterns

---
[⬆️ Back to top](#code-clarity-and-policy-enforcement-guidelines)

## State Management Requirements

**Drift Detection**: Ensure Read operations properly detect and handle resource drift

**Import Functionality**: Verify resource import works correctly and sets all required attributes

**Nested Resources**: Check proper handling of complex nested Azure resource structures

**Resource IDs**: Ensure consistent Azure resource ID parsing and formatting

---
[⬆️ Back to top](#code-clarity-and-policy-enforcement-guidelines)

## Testing Standards

**Acceptance Tests**: Verify comprehensive test coverage including error scenarios

**Test Cleanup**: Ensure tests properly clean up Azure resources

**Multiple Regions**: Check if tests account for regional Azure service availability

**Test Configuration**: Verify test fixtures use appropriate Azure resource configurations

---
[⬆️ Back to top](#code-clarity-and-policy-enforcement-guidelines)

## Documentation Quality

**Examples**: Ensure realistic and working Terraform configuration examples

**Attributes**: Verify all resource attributes are documented with correct types

**Import Documentation**: Check that import syntax and requirements are clearly documented

---
[⬆️ Back to top](#code-clarity-and-policy-enforcement-guidelines)

## Enforcement Priority

1. **Highest**: Code Comments Policy - Zero tolerance for unnecessary comments
2. **High**: Strategic Decision-Making - Performance-critical choices during pair programming
3. **High**: CustomizeDiff Import Requirements - Critical for compilation
4. **High**: Azure API Integration - Essential for functionality
5. **Medium**: Resource Implementation - Quality standards
6. **Medium**: State Management - Reliability standards
7. **Medium**: Testing and Documentation - Completeness standards

**Performance Decision Framework**: Use strategic guidance above to make rapid, correct decisions during active development work.

## ⚡ Quick Decision Trees

### **Comment Decision Tree (30-second evaluation)**
```text
Is this code being written/reviewed?
├─ YES → Apply comment evaluation
│  ├─ Azure API quirk that's non-obvious? → Comment MAY be acceptable
│  ├─ Complex business logic? → Can it be refactored instead? → Refactor FIRST
│  ├─ SDK workaround/limitation? → Comment MAY be acceptable
│  └─ Everything else → NO COMMENT (refactor instead)
└─ NO → Skip comment evaluation
```

### **Cross-Pattern Consistency Check (15-second scan)**
```text
Working on resource with variants (Linux/Windows VMSS, etc.)?
├─ YES → Quick consistency validation required
│  ├─ Check sibling implementation for identical patterns
│  ├─ Ensure validation logic matches
│  └─ Verify error messages use same format
└─ NO → Standard implementation check
```

### **Azure API Integration Priority (10-second assessment)**
```text
Azure API behavior involved?
├─ YES → High priority validation
│  ├─ PATCH operation? → Check residual state handling
│  ├─ Long-running operation? → Verify polling implementation
│  └─ Error handling? → Ensure 404 detection patterns
└─ NO → Standard coding patterns apply
```

## 📊 Performance Metrics & Success Indicators

### **Real-Time Decision Quality Checklist**
- ✅ **Comment Decision**: Made in <30 seconds using decision tree
- ✅ **Cross-Pattern Check**: Sibling resource validated in <15 seconds
- ✅ **Azure Integration**: Priority assessment completed in <10 seconds
- ✅ **Quality Gate**: Pre-submission validation criteria met
- ✅ **Consistency**: Related implementations checked for alignment

### **Session Performance Indicators**
- **High Performance**: 90%+ decisions made using decision trees
- **Optimal Consistency**: Zero cross-pattern validation misses
- **Enforcement Success**: Zero unnecessary comments accepted
- **Strategic Focus**: Primary effort on code clarity over commenting

### **Continuous Improvement Signals**
- **Decision Speed**: Decreasing time to reach enforcement decisions
- **Pattern Recognition**: Faster identification of Azure API quirks vs standard patterns
- **Refactoring Suggestions**: Increasing ratio of refactoring suggestions vs comment acceptance

## 🎯 Context-Aware AI Optimization

### **Session Context Indicators**
- **Active Development**: User actively coding → Apply real-time decision trees
- **Code Review**: User reviewing code → Focus on consistency validation
- **Architecture Discussion**: User planning → Emphasize strategic decision framework
- **Problem Solving**: User debugging → Prioritize Azure API integration patterns

### **Smart Pattern Recognition**
- **Resource Type Context**: Automatically apply VMSS/Storage/Network specific patterns
- **Implementation Approach**: Detect typed vs untyped resource patterns for appropriate guidance
- **Azure Service Context**: Recognize CDN/Compute/Database specific enforcement needs
- **Development Phase**: Adjust guidance intensity based on implementation vs maintenance mode

### **Adaptive Enforcement Intensity**
- **High Intensity**: New resource implementation, complex Azure services, cross-pattern validation
- **Medium Intensity**: Bug fixes, updates, standard patterns
- **Low Intensity**: Documentation updates, minor configuration changes

---
[⬆️ Back to top](#code-clarity-and-policy-enforcement-guidelines)
