# .tfproviderlint — AST-based schema linter

A fast, AST-based linter for the AzureRM provider's **resource** and **data
source** schemas. It reimplements the schema-lint rules (`SL001`–`SL013`) as
small checks over a syntactic schema tree built with the standard `go/ast`.

Unlike the JSON-schema linter in `internal/tools/schema-lint`, this tool never
compiles the provider or renders its schema. It parses Go source directly (no
type checking), which makes it fast enough to lint a single file and to run on
every pull request. Linting the whole `internal/services` tree takes well under
a second.

This is a standalone Go module (it lives in the repository but is not part of the
`terraform-provider-azurerm` module), so it does not affect the provider build or
its vendored dependencies. It has **no external dependencies** — only the Go
standard library — so there is no `go.sum` to maintain and nothing to fetch when
building it.

## Usage

```bash
cd .tfproviderlint/schemalint

# list the available rules
go run . list

# lint a single file (paths resolve against -C, default the current directory)
go run . check -C ../.. ../../internal/services/foo/foo_resource.go

# lint a directory tree
go run . check -C ../.. ../../internal/services/foo

# include suggested fixes for fixable findings
go run . check -C ../.. -fix ../../internal/services/foo

# run or disable specific rules
go run . check -C ../.. -rules SL002,SL007 ../../internal/services/foo
go run . check -C ../.. -disable SL001 ../../internal/services/foo

# machine-readable output
go run . check -C ../.. -format json ../../internal/services/foo

# lint only the schema properties added since a base branch (see Diff mode)
go run . check -diff-base origin/main
```

Or build a binary and run it from the repository root:

```bash
( cd .tfproviderlint/schemalint && go build -o /tmp/schemalint . )
/tmp/schemalint check internal/services/foo/foo_resource.go
```

### Flags (`check`)

| Flag | Default | Description |
|------|---------|-------------|
| `-C` | `.` | repository root; targets and `git` run relative to it |
| `-format` | `text` | output format: `text` or `json` |
| `-fix` | `false` | include suggested fixes for fixable findings |
| `-diff-base` | | only report findings on lines added since this git ref (see [Diff mode](#diff-mode)) |
| `-rules` | | comma-separated rule IDs to run (default: all) |
| `-disable` | | comma-separated rule IDs to disable (takes precedence) |
| `-fail-on-error` | `true` | exit non-zero when any `error`-severity finding is present |

## Rules

| ID | Name | Severity | Checks |
|----|------|----------|--------|
| SL001 | property-description-required | warning | every property sets a non-empty `Description` |
| SL002 | single-property-block | warning | a `MaxItems: 1` block with a single nested property should be flattened |
| SL003 | limits-on-non-collection | error | `MinItems`/`MaxItems` are only set on `TypeList`/`TypeSet` |
| SL004 | avoid-none-value | warning | a user-settable enum should not accept `None`/`Off`/`Default`/`Disabled`; detects `StringInSlice` literals, `string(pkg.…None)` enum constants, and validators nested in `validation.Any`/`All` |
| SL005 | validation-required | warning | user-settable string/numeric arguments set validation |
| SL006 | block-needs-constraint | warning | a block with no required fields sets `AtLeastOneOf`/`ExactlyOneOf` |
| SL007 | array-limits | warning | a scalar array declares `MaxItems` |
| SL008 | sku-field-naming | warning | group multiple `sku_*` fields into a single `sku` block |
| SL009 | unit-in-naming | warning | unit-of-measure suffixes use the `_in_<unit>` form (e.g. `size_in_mb`) |
| SL010 | no-abbreviations | warning | property names use full words, not abbreviations |
| SL011 | redundant-is-prefix | warning | boolean names do not start with a redundant `is_` |
| SL012 | redundant-suffix | warning | names drop a redundant grouping-word suffix (`_properties`/`_config`/`_profile`) |
| SL013 | id-reference-validation | warning | `*_id` references use a resource-specific ID validator or `IsUUID`, not just `StringIsNotEmpty` |

Only `SL003` is `error` severity; the rest are `warning`. With
`-fail-on-error` (the default), the process exits non-zero only when an `error`
finding is present.

### Coverage note

Because the linter reads source rather than the compiled schema, properties
defined by opaque helper calls (for example `commonschema.Location()` or
`commonschema.ResourceIDReferenceOptional(...)`) have no visible schema body.
Rules that need the schema (type, validation, limits) skip these; name-only
rules (SL009, SL010, SL012) still apply. This is by design: newly added
properties in a pull request are written as inline `*pluginsdk.Schema` literals,
which are fully analysed.

Two consequences follow for whole-resource runs (they do not affect diff mode,
which only lints newly added inline properties):

- **External composition is invisible.** Enum values sourced from an external
  SDK function (`validation.StringInSlice(pkg.PossibleValuesForX(), ...)`) and
  schemas produced by `commonschema.*` helpers live in other modules, so their
  properties cannot be inspected. The compiled-schema linter can, because it
  evaluates that code at runtime.
- **Local schema helpers are resolved.** When a property value or a block's
  `Elem` schema is a call to a package-level `schemaXxx()` helper that returns a
  literal (directly, via `s := ...; return s`, or via an inline `func() {...}()`),
  the linter follows it, so block-level rules run and children are reported with
  full dotted paths. Helpers that build their map imperatively (appending to it
  rather than returning a literal), or that live in another package, are not
  followed; their leaf properties are then linted as their own root instead.

## Diff mode

In a pull request, existing properties usually cannot be changed to satisfy a
rule — renaming a released property is a breaking change. Diff mode reports only
findings on the lines a change set **adds**.

Pass `-diff-base <ref>`. The tool computes the merge base of `<ref>` and `HEAD`,
diffs the working tree against it (restricted to resource and data source
files), and reports a finding only when the property's map key is on an added
line. A brand new resource is linted in full; a new child added to a pre-existing
block is linted while its siblings are not.

```bash
GITHUB_BASE_REF=main bash ./scripts/run-schema-lint-diff.sh
```

## Adding a rule

1. Add a `slNNN.go` file in `rules/` defining a `*Rule` with a `Check` function
   that iterates the property nodes:

   ```go
   var slNNN = &Rule{
       ID: "SLNNN", Name: "my-rule", Severity: Warning,
       Check: checkSLNNN,
   }

   func checkSLNNN(res *schematree.Result, report ReportFunc) {
       for _, n := range res.All {
           // inspect n.Name / n.Path / n.Schema and call report(n, msg, fix)
       }
   }
   ```

2. Register it by appending to `Registry` in [rules/rules.go](rules/rules.go).
3. Add a table-driven test alongside it (see the existing `rules/slNNN_test.go`).

To make a rule fixable, set `Fixable: true` and pass a non-empty suggestion as
the last argument to `report`; it is surfaced by `-fix`.
