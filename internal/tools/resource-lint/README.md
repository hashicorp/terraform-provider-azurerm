# resource-lint

A code linting tool for the AzureRM Terraform Provider, designed to check code consistency with rules defined in `/contributing`.

## Lint Checks

For additional information about each check, see the documentation in the passes directory (e.g., `passes/doc.go`).

### Azure Best Practice Checks

| Check | Description |
|-------|-------------|
| AZBP001 | check for all String arguments have `ValidateFunc` |
| AZBP002 | check for `Optional+Computed` fields follow conventions |
| AZBP003 | check for `pointer.ToEnum` to convert Enum type instead of explicitly type conversion |
| AZBP004 | check for zero-value initialization followed by nil check and pointer dereference that should use `pointer.From` |
| AZBP005 | check that Go source files have the correct licensing header |

### Azure New Resource Checks

| Check | Description | Comments |
|-------|-------------|----------|
| AZNR001 | check for Schema field ordering | When git filter is on, this analyzer only runs on newly created resources/data sources |
| AZNR002 | check for top-level updatable arguments are included in Update func | This analyzer currently only runs on typed resources |

### Azure Naming Rule Checks

| Check | Description |
|-------|-------------|
| AZRN001 | check for percentage properties use `_percentage` suffix instead of `_in_percent` |

### Azure Reference Error Checks

| Check | Description |
|-------|-------------|
| AZRE001 | check for fixed error strings using `fmt.Errorf` instead of `errors.New` |

### Azure Schema Design Checks

| Check | Description |
|-------|-------------|
| AZSD001 | check for `MaxItems:1` blocks with single property should be flattened |
| AZSD002 | check for `AtLeastOneOf` or `ExactlyOneOf` validation on TypeList fields with all optional nested fields |

## Usage

### Quick Start

```bash
# Run from the terraform-provider-azurerm root directory

# Check your local branch changes (auto-detect changed lines and packages)
go run ./internal/tools/resource-lint

# Check from a diff file
go run ./internal/tools/resource-lint --diff=changes.txt

# Check specific packages
go run ./internal/tools/resource-lint ./internal/services/compute/...

# Check all issues in packages (no filtering)
go run ./internal/tools/resource-lint --all ./internal/services/...
```

### Options

```bash
--diff=<file>      # Read diff from file instead of git
--all              # Check all issues in packages (not just changes)
--remote=<name>    # Specify git remote (default: auto-detect origin/upstream)
--base=<branch>    # Specify base branch (default: auto-detect from git config or 'main')
--list             # List all available checks
--help             # Show help
```

**Note**: By default, all code is analyzed but only issues on changed lines are reported. Use `--all` to report issues on all lines.

### Output

Results are printed to **standard output**:

| Scenario | Output | Exit Code |
|----------|--------|-----------|
| Issues found | File path, line number, check ID for each issue + summary | 1 |
| No issues | `✓ Analysis completed successfully with no issues found` | 0 |
| Errors (build failures, etc.) | Error message with details | 2 |

#### Example Output

```
go run ./internal/tools/resource-lint
2026/01/05 10:39:01 Using local git diff mode
2026/01/05 10:39:01 Current branch: my-feature-branch
2026/01/05 10:39:02 Merge-base with origin/main: 0aac888
2026/01/05 10:39:03 ✓ Found 9 changed files with 1553 changed lines
2026/01/05 10:39:03 Auto-detected 1 changed packages:
2026/01/05 10:39:03   ./internal/services/policy
2026/01/05 10:39:03 Loading packages...
2026/01/05 10:40:36 Running analysis...

internal/services/policy/policy_resource.go:55:19: AZBP001: string argument "display_name" must have ValidateFunc
internal/services/policy/policy_resource.go:94:18: AZBP002: field "policy_rule" is Optional+Computed but missing required comment
internal/services/policy/policy_resource.go:162:19: AZBP003: use `pointer.ToEnum` to convert Enum type

2026/01/05 10:40:40 Found 3 issue(s)
```

## Limitations

Schema-related checks (e.g., AZNR002, AZSD001, AZSD002) analyze schemas defined as `map[string]*pluginsdk.Schema` or `map[string]*schema.Schema` composite literals returned from functions. This includes:

- **Direct returns**: `return &map[string]*pluginsdk.Schema{...}`
- **Variable returns**: `output := map[string]*pluginsdk.Schema{...}; return output` (captures initial `:=` definition only)
- **Inline schema definitions**: `return &pluginsdk.Schema{...}`
- **Cross-package function calls**: Only `commonschema` package is currently supported
- **Same-package helper functions** returning schemas

Schemas defined in other ways are excluded to reduce false positives from runtime modifications that cannot be determined through static analysis.

For detailed limitations of each analyzer, refer to the documentation in the respective analyzer files (e.g., `passes/AZNR002.go`).

## Contributing: Adding a New Rule

### Rule Naming Convention

Rules follow this naming pattern: `AZ{Category}{Number}`

| Category | Prefix | Description |
|----------|--------|-------------|
| Best Practice | `AZBP` | General coding best practices |
| New Resource | `AZNR` | Rules specific to new resources/data sources |
| Naming Rule | `AZRN` | Naming conventions |
| Reference Error | `AZRE` | Error handling patterns |
| Schema Design | `AZSD` | Schema structure and design |

### Step-by-Step Guide

#### 1. Create the Analyzer File

Create `passes/AZ{XX}{NNN}.go`:

```go
package passes

import (
    "go/ast"
    
    "github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/loader"
    "golang.org/x/tools/go/analysis"
    "golang.org/x/tools/go/analysis/passes/inspect"
    "golang.org/x/tools/go/ast/inspector"
)

// Documentation for the analyzer (shown by --list flag)
const AZXXNNNDoc = `short description of the check

Detailed explanation of what the analyzer checks.

Example violations:
  // bad code example

Valid usage:
  // good code example`

const azxxnnnName = "AZXXNNN"

var AZXXNNNAnalyzer = &analysis.Analyzer{
    Name:     azxxnnnName,
    Doc:      AZXXNNNDoc,
    Run:      runAZXXNNN,
    Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func runAZXXNNN(pass *analysis.Pass) (interface{}, error) {
    inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
    
    nodeFilter := []ast.Node{
        (*ast.CallExpr)(nil),  // filter for specific AST node types
    }
    
    inspector.Preorder(nodeFilter, func(n ast.Node) {
        // Get position info
        pos := pass.Fset.Position(n.Pos())
        filename := pos.Filename
        
        // Skip if file/line not in changed set (for incremental analysis)
        if !loader.IsFileChanged(filename) {
            return
        }
        
        // Your analysis logic here...
        
        // Report issue (only if line is changed)
        if loader.ShouldReport(filename, pos.Line) {
            pass.Reportf(n.Pos(), "%s: description of the issue", azxxnnnName)
        }
    })
    
    return nil, nil
}
```

#### 2. Register the Analyzer

Add your analyzer to `passes/checks.go`:

```go
var AllChecks = []*analysis.Analyzer{
    // ... existing analyzers
    AZXXNNNAnalyzer,  // Add your new analyzer here
}
```

#### 3. Create Test File

Create `passes/AZ{XX}{NNN}_test.go`:

```go
package passes_test

import (
    "testing"
    
    "github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/passes"
    "golang.org/x/tools/go/analysis/analysistest"
)

func TestAZXXNNN(t *testing.T) {
    testdata := analysistest.TestData()
    analysistest.Run(t, testdata, passes.AZXXNNNAnalyzer, "testdata/src/azxxnnn")
}
```

#### 4. Create Test Data

Create `passes/testdata/src/azxxnnn/a.go`:

```go
package azxxnnn

// Valid cases - should NOT trigger warnings
func validCases() {
    // good code examples
}

// Invalid cases - SHOULD trigger warnings
func invalidCases() {
    // bad code example // want `AZXXNNN`
    // The `// want` comment tells the test framework to expect a diagnostic
}
```

#### 5. Run Tests

```bash
# Run your specific test
go test -v ./internal/tools/resource-lint/passes -run TestAZXXNNN

# Run all tests
go test -v ./internal/tools/resource-lint/...
```

### Key APIs

#### `loader` Package - Change Filtering

| Function | Description |
|----------|-------------|
| `loader.IsFileChanged(filename)` | Check if file has changes (for filtering) |
| `loader.ShouldReport(filename, line)` | Check if specific line should be reported |
| `loader.IsNewFile(filename)` | Check if file is newly created |

#### `helper` Package - Common Utilities

| Function/Type | Description |
|---------------|-------------|
| **Output Formatting** | |
| `helper.FixedCode(s)` | Format string as suggested fix (green) |
| `helper.IssueLine(s)` | Format string as problematic code (yellow) |
| `helper.Bold(s)` | Format string as bold |
| `helper.FormatCode(s)` | Format string as code (magenta) |
| **Schema Detection** | |
| `helper.IsSchemaMap(cl, info)` | Check if composite literal is `map[string]*pluginsdk.Schema` |
| `helper.IsSchemaSchema(info, cl)` | Check if composite literal is `pluginsdk.Schema` type |
| `helper.IsNestedSchemaMap(file, cl)` | Check if schema map is nested within an `Elem` field |
| **Type Detection** | |
| `helper.IsResourceData(info, sel)` | Check if selector is `*pluginsdk.ResourceData` |
| `helper.GetReceiverTypeName(expr)` | Get receiver type name from method declaration |
| **Typed Resource** | |
| `helper.TypedResourceInfo` | Struct containing typed resource metadata (model, schema, CRUD funcs) |
| `helper.NewTypedResourceInfo(name, file, info)` | Parse typed resource from file |
| `helper.GetTypedServices()` | Get all registered typed services from provider |

#### `passes/schema` Package - Schema Analysis (Fact Exporters)

These analyzers export facts that can be reused by other analyzers via `pass.ResultOf`:

| Analyzer | Result Type | Description |
|----------|-------------|-------------|
| `schema.CompleteSchemaAnalyzer` | `*CompleteSchemaInfo` | Resolves all schema fields including cross-package calls (e.g., `commonschema.ResourceGroupName()`) |
| `schema.TypedResourceInfoAnalyzer` | `[]*helper.TypedResourceInfo` | Extracts typed resource metadata (Arguments, Attributes, CRUD methods) |
| `schema.CommonAnalyzer` | `*CommonSchemaInfo` | Caches `commonschema` package schema definitions |
| `schema.InlineSchemaAnalyzer` | `*InlineSchemaInfo` | Parses inline `&pluginsdk.Schema{}` literals |

**Usage Example:**

```go
var MyAnalyzer = &analysis.Analyzer{
    Name:     "myanalyzer",
    Requires: []*analysis.Analyzer{schema.CompleteSchemaAnalyzer},
    Run:      runMyAnalyzer,
}

func runMyAnalyzer(pass *analysis.Pass) (interface{}, error) {
    // Get pre-computed schema info
    schemaInfo := pass.ResultOf[schema.CompleteSchemaAnalyzer].(*schema.CompleteSchemaInfo)
    
    // Use it to get resolved schema fields for a schema map
    fields := schemaInfo.SchemaFields[schemaMapLit.Pos()]
    // ...
}
```

### Tips

1. **Always use change filtering** - Call `loader.IsFileChanged()` and `loader.ShouldReport()` to support incremental analysis
2. **Pre-filter AST nodes** - Use `inspector.Preorder` with specific node types for better performance
3. **Skip test files** - Usually add `strings.HasSuffix(filename, "_test.go")` check
4. **Skip migration packages** - Add `strings.Contains(pass.Pkg.Path(), "/migration")` check if needed
5. **Reuse schema analyzers** - Add `schema.CompleteSchemaAnalyzer` to `Requires` instead of re-parsing schemas
6. **Write comprehensive tests** - Include both valid and invalid cases in test data