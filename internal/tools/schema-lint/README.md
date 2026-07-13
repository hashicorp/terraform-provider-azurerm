# schema-lint

A modular, pluggable linter for the AzureRM provider's **resource** and **data source** schemas.

In the spirit of a markdown linter (e.g. markdownlint), each check is an independent *rule* with a
stable ID. Rules can be enabled, disabled and re-configured through a `.schema-lint.json` config file
or via CLI flags. The linter walks every property of every resource and data source — including nested
blocks — and reports findings. Rules that can recommend a remediation implement a `Fixer` hint; pass
`-fix` to include their suggested fixes in the output. The linter does not rewrite source itself:
because it inspects the compiled provider schema, a fix is surfaced as an actionable suggestion.

## Usage

```bash
# list the available rules
go run ./internal/tools/schema-lint list

# lint the whole provider schema with defaults
go run ./internal/tools/schema-lint check

# run only specific rules
go run ./internal/tools/schema-lint check -rules=SL002,SL005

# disable a rule
go run ./internal/tools/schema-lint check -disable=SL001

# include suggested fixes for fixable findings
go run ./internal/tools/schema-lint check -fix

# scope to specific resources / data sources
go run ./internal/tools/schema-lint check -resource=azurerm_kubernetes_cluster

# lint only properties added since a base schema (see Diff mode)
go run ./internal/tools/schema-lint check -diff .release/provider-schema.json

# machine readable output
go run ./internal/tools/schema-lint check -format=json
```

There is also a Make target:

```bash
make schema-lint
```

### Flags (`check`)

| Flag | Default | Description |
|------|---------|-------------|
| `-config` | `.schema-lint.json` if present | path to a config file |
| `-rules` | `all` | comma-separated rule IDs to run, or `all` |
| `-disable` | | comma-separated rule IDs to disable (takes precedence) |
| `-resource` | | comma-separated resource/data source type names to include |
| `-skip-resource` | | comma-separated resource/data source type names to skip |
| `-format` | `text` | output format: `text` or `json` |
| `-fail-on-error` | `true` | exit non-zero when any `error`-severity finding is present |
| `-fix` | `false` | include suggested fixes for fixable findings (see [Suggested fixes](#suggested-fixes)) |
| `-diff` | | path to a base schema dump; report only findings on properties added since the base (see [Diff mode](#diff-mode)) |

Precedence for enabling/disabling and severity is **CLI flags > config file > rule defaults**.

## Configuration

An optional `.schema-lint.json` file (markdownlint-style) configures the run:

```json
{
  "default": true,
  "rules": {
    "SL001": { "enabled": false },
    "SL004": { "severity": "error" }
  },
  "skipResources": ["azurerm_some_resource"]
}
```

- `default` — whether rules are enabled unless individually overridden (default `true`).
- `rules.<ID>.enabled` — enable/disable a specific rule.
- `rules.<ID>.severity` — override a rule's severity (`error` or `warning`).
- `rules.<ID>.options` — arbitrary per-rule options (for rules that support them).
- `includeResources` / `skipResources` — scope which resource/data source types are linted.

## Rules

| ID | Name | Severity | Fixable | Checks |
|----|------|----------|---------|--------|
| SL001 | property-description-required | warning | | every property sets a non-empty `Description` |
| SL002 | optional-and-required | error | | a property is not both `Optional` and `Required` |
| SL003 | required-and-computed | error | | a property is not both `Required` and `Computed` |
| SL004 | computed-only-forcenew | warning | ✓ | a Computed-only property does not set `ForceNew` |
| SL005 | collection-elem-required | error | | a `TypeList`/`TypeSet` property defines an `Elem` |
| SL006 | single-property-block | warning | ✓ | a `MaxItems: 1` block with a single nested property should be flattened |
| SL007 | limits-on-non-collection | error | ✓ | `MinItems`/`MaxItems` are only set on `TypeList`/`TypeSet` |
| SL008 | avoid-none-value | warning | | a user-settable enum should not accept `None`/`Off`/`Default`/`Disabled` |
| SL009 | validation-required | warning | | user-settable string/numeric arguments set validation |
| SL010 | block-needs-constraint | warning | | a block with no required fields sets `AtLeastOneOf`/`ExactlyOneOf` |
| SL011 | array-limits | warning | | a scalar array declares `MinItems`/`MaxItems` |
| SL012 | sku-field-naming | warning | | group multiple `sku_*` fields into a single `sku` block |
| SL013 | unit-in-naming | warning | ✓ | unit-of-measure suffixes use the `_in_<unit>` form (e.g. `size_in_mb`) |
| SL014 | no-abbreviations | warning | ✓ | property names use full words, not abbreviations (e.g. `virtual_machine` not `vm`, `percentage` not `pct`) |
| SL015 | redundant-is-prefix | warning | ✓ | boolean names do not start with a redundant `is_` |
| SL016 | redundant-suffix | warning | ✓ | names drop a redundant grouping-word suffix (`_properties`/`_config`/`_profile`) |
| SL017 | id-reference-validation | warning | | `*_id` references use a resource-specific ID validator or `IsUUID`, not just `StringIsNotEmpty` |

## Suggested fixes

Fixable rules attach a concrete remediation to each finding. Pass `-fix` to include these in the
output (a `→ fix:` line in text mode, or a `fixSuggestion` field in JSON mode):

```text
azurerm_api_management (resource)
  warning  SL006   sign_in: block "sign_in" has a single nested property "enabled" (MaxItems 1); consider flattening it
           → fix: replace the block "sign_in" with a single top-level boolean "sign_in_enabled"
```

Because the linter inspects the compiled provider schema rather than its Go source, `-fix` reports the
recommended change; it does not edit source files.

## Diff mode

In a pull request, existing properties usually cannot be changed to satisfy a rule — renaming a released
property is a breaking change. Diff mode holds only the properties a change set **adds** to the rules,
leaving pre-existing ones untouched.

Pass `-diff <base>`, where `<base>` is a provider schema dump (as produced by `schema-api -export`). A
finding is reported only when the property (by its dotted path) — or the whole resource / data source —
is absent from the base schema:

```bash
# export a base schema, then lint only what changed since it
go run internal/tools/schema-api/main.go -export /tmp/base.json
go run ./internal/tools/schema-lint check -diff /tmp/base.json
```

A brand new resource is linted in full; a new child added to a pre-existing block is linted while its
siblings are not.

### Continuous integration

The [schema-lint workflow](../../../.github/workflows/schema-lint.yaml) runs two jobs on pull requests
that touch `internal/services/**` schemas:

- **check** — lints entirely new resources / data sources
  ([run-schema-lint.sh](../../../scripts/run-schema-lint.sh)).
- **diff** — exports the base branch's schema into a temporary `git worktree` and runs `-diff`, so only
  newly added properties are checked ([run-schema-lint-diff.sh](../../../scripts/run-schema-lint-diff.sh)).

Both are skipped when a maintainer applies the `schema-accepted` label. Reproduce the diff check locally
against the branch you are merging into:

```bash
GITHUB_BASE_REF=main bash ./scripts/run-schema-lint-diff.sh
```

## Design considerations coverage

These rules are derived from the provider's schema design considerations. The mechanically checkable
considerations are implemented as rules; the remainder require semantic judgement a linter cannot make.
Because the linter is intended for new resources, the checkable rules are enabled by default and noise
on the existing provider is expected.

| Consideration | Rule / status |
|---------------|---------------|
| Features toggled by `Enabled` / flattening nested properties | SL006 |
| MinItems / MaxItems placement | SL007 |
| `None`/`Off`/`Default` value removal | SL008 |
| Validation required (strings, numeric ranges) | SL009 |
| `TypeList` with no required fields needs `AtLeastOneOf`/`ExactlyOneOf` | SL010 |
| Array fields should declare MinItems / MaxItems | SL011 |
| SKU field naming (`sku_name`, `sku_family`) | SL012 |
| Unit-of-measure naming (`_in_<unit>`, from reference-naming.md) | SL013 |
| No abbreviations, incl. percentage (`vm` -> `virtual_machine`, `pct` -> `percentage`, from reference-naming.md) | SL014 |
| Redundant `is_` boolean prefix (from reference-naming.md) | SL015 |
| Redundant grouping-word suffix (`_properties`/`_config`/`_profile`, from reference-naming.md) | SL016 |
| ID reference validation (`commonids`/`IsUUID` for `*_id`, from schema-design-considerations.md) | SL017 |
| Portal terminology, argument grouping, collection ambiguity, `type` discriminator split, preview fields, portal-required | not lintable (semantic / architectural) |

SL008 only applies to enum validators. The special sentinel values and the enum heuristic live in the rule
itself; `providerjson.SchemaJSON` exposes an `AcceptsValue(string) bool` probe (backed by the property's
`ValidateDiagFunc`/`ValidateFunc`), and the rule reports a value only when the validator accepts
`None`/`Off`/`Default`/`Disabled` but rejects same-shaped decoy strings — the signature of a finite
`StringInSlice`-style enum rather than a length or format validator.

## Adding a rule

1. Create a file in `rules/` implementing `PropertyRule` (per-property) or `ResourceRule` (per-resource):

   ```go
   type myRule struct{}

   func (myRule) ID() string                { return "SL018" }
   func (myRule) Name() string              { return "my-rule" }
   func (myRule) Description() string       { return "what it checks" }
   func (myRule) DefaultSeverity() rules.Severity { return rules.SeverityWarning }

   func (r myRule) CheckProperty(ctx PropertyContext) []Finding {
       // inspect ctx.Schema / ctx.Path and return Findings
   }
   ```

2. Register it by appending to `AllRules` in [rules/rule.go](rules/rule.go).
3. Add a table-driven test alongside it (see the existing `*_test.go` files).

To make a rule fixable, also implement `FixHint() string` (a short label shown by `list`) and build
findings with `propertyFindingFix(r, ctx, message, suggestion)` so each finding carries a concrete
remediation surfaced by `-fix`.

## Architecture

- `rules/` — rule contracts (`Rule`, `PropertyRule`, `ResourceRule`, `Fixer`), the `Finding`/context
  types, the `AllRules` registry, and the concrete rules.
- `engine/` — loads the schema via `providerschema`, walks resources/data sources and nested
  blocks, applies the active rules and resolves severity.
- `config/` — the `.schema-lint.json` model and loader.
- `report/` — text and JSON reporters.
- `main.go` — the CLI.

The schema model is reused from `internal/tools/providerschema` (shared with `schema-api`), so the
linter does not duplicate schema extraction.
