# document-fmt

A validation and auto-fixing tool for Terraform provider documentation in the terraform-provider-azurerm repository. It enforces documentation standards by applying a set of customizable rules to ensure consistency and completeness across all resource and data source documentation files.

## Quickstart

### Validate Documentation
```bash
# From the repository root
go run ./internal/tools/document-fmt/main.go validate

# Or use the make target
make document-validate

# Or use the shell script directly
./scripts/documentfmt-validate.sh
```

### Fix Documentation Issues
```bash
# Automatically fix issues where possible
go run ./internal/tools/document-fmt/main.go fix

# Or use the make target
make document-fix

# Or use the shell script directly
./scripts/documentfmt-fix.sh
```

### Command-Line Options
```bash
# Validate specific service
go run ./internal/tools/document-fmt/main.go validate --service=storage

# Validate specific resource
go run ./internal/tools/document-fmt/main.go validate --resource=azurerm_storage_account

# Enable debug logging
go run ./internal/tools/document-fmt/main.go validate --debug

# Specify provider directory (defaults to current directory)
go run ./internal/tools/document-fmt/main.go validate --provider-dir=/path/to/provider
```

## Structure & Design

The document-fmt tool follows a modular architecture with clear separation of concerns:

### Architecture Overview

```
document-fmt/
├── main.go                    # Entry point
├── cmd/                       # Command-line interface
│   ├── cmds.go               # Command implementations (validate, fix, scaffold)
│   └── flags.go              # CLI flag definitions
├── rule/                      # Validation rules
│   ├── rule.go               # Rule interface and registration
│   ├── rule_g001_*.go        # Global rules (G prefix)
│   ├── rule_s001_*.go        # Section-specific rules (S prefix)
│   └── rule_exceptions.go    # Exception handling
├── data/                      # Data models
├── markdown/                  # Markdown parsing & manipulation
├── template/                  # Template rendering
├── validator/                 # Validation engine
├── differror/                 # Error handling
└── util/                      # Utilities
    ├── util.go               # General helpers
    ├── log.go                # Logging setup
```

### Major Components

#### 1. Command Layer (`cmd/`)
- **cmds.go**: Implements the three main commands using the Cobra CLI framework:
  - `validate`: Scans documentation and reports issues
  - `fix`: Attempts to automatically correct found issues
  - `scaffold`: (Not yet implemented) Will generate documentation templates
- **flags.go**: Defines CLI parameters for filtering resources/services and controlling behavior

#### 2. Rule System (`rule/`)
Rules implement the `Rule` interface with four methods:
- `ID()`: Returns rule identifier (e.g., "G001", "S001")
- `Name()`: Returns human-readable name
- `Description()`: Explains what the rule validates
- `Run()`: Executes validation logic, optionally fixes issues

**Rule Categories:**
- **Global Rules (G prefix)**: Apply to entire documents
- **Section Rules (S prefix)**: Validate specific documentation sections

Rules are registered in the `rule/rule.go` file and can be found in the `rule/` directory with the naming pattern `rule_<category><number>_<name>.go`.

#### 3. Data Models (`data/`)
- **TerraformNodeData**: Central data structure containing:
  - Resource/data source metadata (name, type, service)
  - Document reference
  - API version information
  - Timeout configurations
  - Validation errors
- **Document**: Represents parsed markdown with sections
- **API**: Tracks Azure API provider names and versions
- **Service**: Represents service package metadata

#### 4. Markdown Parser (`markdown/`)
- Parses documentation files into structured sections
- Each section type (Frontmatter, Title, Example, Arguments, Attributes, Timeouts, Import, API, etc.) has its own implementation
- Supports reading, modifying, and writing sections back to files
- Preserves formatting and structure during modifications

#### 5. Validator (`validator/`)
- Orchestrates rule execution across resources
- Maintains state and accumulates errors
- Coordinates between validation mode and fix mode

#### 6. Template Engine (`template/`)
- Uses Go's `text/template` package to generate section content
- Renders API provider sections, timeout blocks, etc.
- Ensures consistent formatting across generated content

### Data Flow

1. **Discovery**: Scans provider directory to find all resources/data sources
2. **Parsing**: Reads source code to extract API versions and timeout configurations
3. **Document Loading**: Parses markdown documentation files into structured sections
4. **Validation**: Executes each rule against each resource's data
5. **Reporting/Fixing**: Either reports errors or applies fixes
6. **Writing**: Modified documents are written back to disk (fix mode only)

## How to Extend: Adding Rules

### Rule Naming Convention
- **G###** - Global rules (document-wide validations)
- **S###** - Section-specific rules (validate individual sections)
- Use sequential numbering (G003, S004, etc.)

### Creating a New Rule

1. **Create the rule file** in `internal/tools/document-fmt/rule/`:

```go
// rule_s004_example_section.go
package rule

import (
    "fmt"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/markdown"
)

type S004 struct{}

var _ Rule = S004{}

func (r S004) ID() string {
    return "S004"
}

func (r S004) Name() string {
    return "Example Section"
}

func (r S004) Description() string {
    return "validates the Example Usage section exists and follows conventions"
}

func (r S004) Run(d *data.TerraformNodeData, fix bool) []error {
    if !d.Document.Exists {
        return nil
    }

    errs := make([]error, 0)
    
    // Find the section in the document
    var section *markdown.ExampleSection
    for _, sec := range d.Document.Sections {
        if s, ok := sec.(*markdown.ExampleSection); ok {
            section = s
            break
        }
    }
    
    if section == nil {
        errs = append(errs, fmt.Errorf("%s: missing Example Usage section", IdAndName(r)))
        
        if fix {
            // Create and add the missing section
            section = &markdown.ExampleSection{}
            // Set default content...
            d.Document.Sections = append(d.Document.Sections, section)
            d.Document.HasChange = true
        }
        return errs
    }
    
    // Perform validation on the section
    // Compare current content with expected content
    // If fixing, update the section content
    
    return errs
}
```

2. **Register the rule** in `rule/rule.go`:

```go
var Registration = map[string]Rule{
    // ... existing rules ...
    S004{}.ID(): S004{}, // Example Section
}
```

3. **Test the rule**:

```bash
# Run validation with your new rule
go run ./internal/tools/document-fmt/main.go validate

# Test fixing capability
go run ./internal/tools/document-fmt/main.go fix
```

### Rule Best Practices

- **Graceful Handling**: Return early if the document doesn't exist or if rule preconditions aren't met
- **Clear Error Messages**: Use `IdAndName(r)` prefix for consistency
- **Idempotent Fixes**: Ensure running fix multiple times produces the same result
- **Section Ordering**: Maintain proper section order when adding new sections
- **Template Usage**: Use the template engine for generating consistent section content
- **Logging**: Add debug logging for troubleshooting (use `log.WithFields()`)
- **Testing Edge Cases**: Handle resources with unusual configurations gracefully

### Adding Exception Handling

If certain resources need to skip a rule, add exceptions in `rule/rule_exceptions.go`:

```go
func shouldSkipRule(resourceName string, ruleID string) bool {
    exceptions := map[string][]string{
        "azurerm_special_resource": {"S001", "S002"},
    }
    
    if rules, ok := exceptions[resourceName]; ok {
        for _, r := range rules {
            if r == ruleID {
                return true
            }
        }
    }
    return false
}
```

### Adding New Section Types

If your rule needs to work with a new section type:

1. Create section implementation in `markdown/section_*.go`
2. Implement the `Section` interface
3. Register the section type in `markdown/registration.go`
4. Add template if needed in `template/template.go`

## Integration with CI/CD

The tool integrates with the repository's CI/CD pipeline:

- **Pre-commit**: Can be run locally before committing
- **CI Validation**: `scripts/documentfmt-validate.sh` runs in CI to catch issues
- **Automated Fixing**: Developers can run `make document-fix` to auto-correct issues
- **Exit Codes**: Non-zero exit code on validation failure stops the build

## Additional Notes

- The tool uses the `afero` filesystem abstraction for testability
- All rules can operate in both validate-only and fix modes
- The validator maintains state across all resources being processed
- Template rendering uses Go's `text/template` package for flexibility
- Color-coded output helps distinguish errors, warnings, and success messages
