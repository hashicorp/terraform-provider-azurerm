# pr-guide-suggestions — AST-based schema linter

A fast, AST-based linter for the AzureRM provider's **resource** and **data
source** schemas. It implements the pr-guide-suggestions rules (`SL001`–`SL013`) as
small checks over a syntactic schema tree built with the standard `go/ast`.

It is a normal package of the provider module (at `internal/tools/pr-guide-suggestions`)
with **no external dependencies** — only the Go standard library — so it adds
nothing to the provider's `go.mod` or vendor tree and needs no build of the
provider itself.

## Usage

Run it from the repository root:

```bash
# list the available rules
go run ./internal/tools/pr-guide-suggestions list

# lint a single file
go run ./internal/tools/pr-guide-suggestions check internal/services/foo/foo_resource.go

# lint a directory tree
go run ./internal/tools/pr-guide-suggestions check internal/services/foo

# apply auto-fixable fixes (property renames) in place
go run ./internal/tools/pr-guide-suggestions check -fix internal/services/foo

# run or disable specific rules
go run ./internal/tools/pr-guide-suggestions check -rules SL002,SL007 internal/services/foo
go run ./internal/tools/pr-guide-suggestions check -disable SL001 internal/services/foo

# machine-readable output
go run ./internal/tools/pr-guide-suggestions check -format json internal/services/foo

# lint only the schema properties added since a base branch (see Diff mode)
go run ./internal/tools/pr-guide-suggestions check -diff-base origin/main
```

Or build a binary once and reuse it:

```bash
go build -o /tmp/pr-guide-suggestions ./internal/tools/pr-guide-suggestions
/tmp/pr-guide-suggestions check internal/services/foo/foo_resource.go
```

### Flags (`check`)

| Flag | Default | Description |
|------|---------|-------------|
| `-C` | `.` | repository root; targets and `git` run relative to it |
| `-format` | `text` | output format: `text` or `json` |
| `-fix` | `false` | apply auto-fixable fixes (property renames) to files in place (see [Applying fixes](#applying-fixes)) |
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
GITHUB_BASE_REF=main bash ./scripts/run-pr-guide-suggestions-diff.sh
```

## Applying fixes

By default the linter is read-only: it reports findings (with a suggested
remediation for each fixable one) but does not touch your files. Pass `-fix` to
*apply* the auto-fixable findings in place. Today the naming rules
(`SL009`–`SL012`) are auto-applicable: each is a **rename**, and applying it
replaces every quoted occurrence of the property name across the file — not just
the schema map key, but also its `d.Get`/`d.Set` calls, `tfschema` struct tags
and `AtLeastOneOf`/`ExactlyOneOf`/`ConflictsWith` cross-references:

```bash
# rename "vm_count" -> "virtual_machine_count" everywhere in the file
go run ./internal/tools/pr-guide-suggestions check -rules SL010 -fix internal/services/foo/foo_resource.go
```

```text
fixed internal/services/foo/foo_resource.go: renamed "vm_count" → "virtual_machine_count"

applied 1 fix(es)
```

Details and caveats:

- Replacement matches the name as a whole **path segment** inside string
  literals — bounded by a double quote or a dot — so both the standalone key
  (`"vm_count"`) and dotted cross-field references (`"vm_count.0.child"` in
  `AtLeastOneOf`/`ConflictsWith`/`d.Get`) are updated, while superstrings like
  `"vm_count_total"`, bare words in comments and the Go field identifier
  (`VmCount`) are left untouched. Rename the Go field yourself if you want it to
  match.
- Every matching occurrence in the file is replaced. If two blocks happen to
  share a child name, both are renamed; review the diff.
- Non-rename fixable findings (`SL002` flatten, `SL003` remove limits) are
  structural and are **not** auto-applied — `-fix` reports them as remaining
  work.
- `-fix` writes files in place; run it on a clean working tree so you can review
  the change with `git diff`. Re-run the linter afterwards to converge on any
  follow-on findings and to reformat with `gofmt`/`goimports` if needed.

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

To make a rule fixable, set `Fixable: true` and pass a `*Fix` as the last
argument to `report` (pass `nil` when there is no fix). A `Fix` with only a
`Suggestion` is printed as the remediation hint; a `Fix` that also sets `Rename`
to the new property name is auto-applied by `-fix`:

```go
report(n, msg, &Fix{
    Suggestion: fmt.Sprintf("rename %q to %q", n.Name, preferred),
    Rename:     preferred,
})
```
