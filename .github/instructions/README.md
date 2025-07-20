# Terraform AzureRM Provider - AI Development Instructions

This directory contains comprehensive coding guidelines and instruction files designed to work with GitHub Copilot and other AI development tools for the Terraform AzureRM provider.

## 📁 File Organization

### Core Instruction Files

| File | Purpose | Scope |
|------|---------|-------|
| [`coding-style.instructions.md`](./coding-style.instructions.md) | Basic Go formatting and style rules | Fundamental code formatting, import organization |
| [`coding-standards.instructions.md`](./coding-standards.instructions.md) | Comprehensive coding standards | Naming conventions, error handling, file organization |
| [`coding-patterns.instructions.md`](./coding-patterns.instructions.md) | Implementation patterns and examples | Resource implementation, client management, schema design |
| [`provider-guidelines.instructions.md`](./provider-guidelines.instructions.md) | Azure-specific provider patterns | ARM integration, CustomizeDiff, Azure API patterns |
| [`testing-guidelines.instructions.md`](./testing-guidelines.instructions.md) | Testing patterns and best practices | Unit tests, acceptance tests, Azure-specific testing |
| [`documentation-guidelines.instructions.md`](./documentation-guidelines.instructions.md) | Documentation standards | Resource vs data source docs, examples, formatting |

## 🏗️ Implementation Approaches

The Terraform AzureRM provider supports two implementation approaches:

### 🎯 **Typed Resource Implementation (Preferred)**
- **Framework**: Uses `internal/sdk` with type-safe models
- **Pattern**: Receiver methods on resource structs
- **State Management**: `tfschema` tags with `metadata.Decode()`/`metadata.Encode()`
- **Recommended For**: All new resources and data sources

### 🔧 **UnTyped Resource Implementation (Maintenance)**
- **Framework**: Traditional Plugin SDK patterns
- **Pattern**: Function-based CRUD operations
- **State Management**: Direct `d.Set()`/`d.Get()` manipulation
- **Used For**: Existing resource maintenance only

Both approaches are fully documented with examples in the instruction files.

## 📚 Learning Path

For new contributors, we recommend reviewing the files in this order:

1. **Start Here**: [`coding-style.instructions.md`](./coding-style.instructions.md) - Basic formatting rules
2. **Core Standards**: [`coding-standards.instructions.md`](./coding-standards.instructions.md) - Naming and error patterns
3. **Implementation**: [`coding-patterns.instructions.md`](./coding-patterns.instructions.md) - Resource implementation patterns
4. **Azure-Specific**: [`provider-guidelines.instructions.md`](./provider-guidelines.instructions.md) - Azure integration patterns
5. **Testing**: [`testing-guidelines.instructions.md`](./testing-guidelines.instructions.md) - Comprehensive testing guide
6. **Documentation**: [`documentation-guidelines.instructions.md`](./documentation-guidelines.instructions.md) - Writing great docs

## 🔗 Cross-References

### Common Topics Across Files

| Topic | Primary File | Related Files |
|-------|-------------|---------------|
| **Error Handling** | `coding-standards.instructions.md` | `coding-patterns.instructions.md`, `testing-guidelines.instructions.md` |
| **Azure SDK Usage** | `provider-guidelines.instructions.md` | `coding-patterns.instructions.md`, `coding-standards.instructions.md` |
| **CustomizeDiff** | `provider-guidelines.instructions.md` | `coding-patterns.instructions.md`, `testing-guidelines.instructions.md` |
| **Import Management** | `coding-patterns.instructions.md` | `coding-style.instructions.md`, `coding-standards.instructions.md` |
| **Testing Patterns** | `testing-guidelines.instructions.md` | All implementation files |
| **State Management** | `coding-patterns.instructions.md` | `coding-standards.instructions.md`, `provider-guidelines.instructions.md` |

### Implementation-Specific Guidance

#### For Typed Resource Development
1. [`coding-patterns.instructions.md`](./coding-patterns.instructions.md) - Typed implementation patterns
2. [`coding-standards.instructions.md`](./coding-standards.instructions.md) - Typed error handling
3. [`testing-guidelines.instructions.md`](./testing-guidelines.instructions.md) - Testing typed resources

#### For UnTyped Resource Maintenance
1. [`coding-patterns.instructions.md`](./coding-patterns.instructions.md) - UnTyped maintenance patterns
2. [`coding-standards.instructions.md`](./coding-standards.instructions.md) - Traditional error handling
3. [`testing-guidelines.instructions.md`](./testing-guidelines.instructions.md) - Testing untyped resources

## 🎯 AI Development Integration

These instruction files are designed to work seamlessly with:

- **GitHub Copilot**: Referenced in `.vscode/settings.json` for context-aware assistance
- **AI Prompts**: Used by prompt files in `../.github/prompts/` for structured development tasks
- **Code Reviews**: Provide standards for AI-assisted code review processes
- **Documentation**: Guide AI in generating consistent, high-quality documentation

## 🔄 Maintenance

### Keeping Instructions Current

- **Regular Reviews**: Instructions should be reviewed quarterly for accuracy
- **Code Example Validation**: Automated testing of code examples (future enhancement)
- **Azure SDK Updates**: Update patterns when Azure SDK versions change
- **Provider Evolution**: Adjust guidelines as provider architecture evolves

### Contributing to Instructions

When updating instruction files:

1. **Maintain Consistency**: Ensure examples work across all referenced files
2. **Update Cross-References**: Check related files for impact
3. **Test Examples**: Verify all code examples compile and follow current patterns
4. **Azure Alignment**: Ensure Azure-specific guidance matches current SDK and API patterns

## 📋 Quick Reference

### Essential Commands
```bash
# Run tests following guidelines
make test
make testacc TEST=./internal/services/cdn TESTARGS='-run=TestAccCdnFrontDoorProfile'

# Format code according to style guide
gofmt -w .
goimports -w .
```

### Key Patterns
- **Resource Implementation**: See `coding-patterns.instructions.md` for CRUD patterns
- **Error Handling**: Use `fmt.Errorf("action `%s`: %+v", resource, err)` format
- **Azure Client Usage**: Follow patterns in `provider-guidelines.instructions.md`
- **Testing**: Always include import tests with `data.ImportStep()`

---

For questions about these guidelines or suggestions for improvements, please refer to the individual instruction files or contribute updates following the patterns established in this documentation system.
