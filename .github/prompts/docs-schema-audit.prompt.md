---
description: "Docs + Schema Audit Prompt for Terraform AzureRM Provider"
---

# 📋 Docs + Schema Audit (AzureRM)

# 🚫 EXECUTION GUARDRAILS (READ FIRST)

## Renderer artifact
Some chat UIs may display a leading `-` before this prompt's content. Treat that as a rendering artifact and **do not** comment on it. Proceed with the audit.

## Required active file
This prompt audits the **currently-open documentation page** under `website/docs/**`.

If the active editor is not a file under `website/docs/**` (for example if the active editor is this prompt file, or a README), do **not** attempt the audit.

Instead, respond with:

"Cannot run docs/schema audit: active file is not under `website/docs/**`. Open the target docs page and re-run this prompt."

Audit the **currently-open** documentation page under `website/docs/**` for:
- AzureRM documentation standards, and
- parity with the provider schema under `internal/**`.

When reviewing documentation standards, treat these as authoritative:
- `contributing/topics/reference-documentation-standards.md`
- `.github/instructions/documentation-guidelines.instructions.md`

This audit is **optional** and **user-invoked** (no CI enforcement).

## ⚡ Mandatory procedure

### 1) Identify the Terraform object from the doc path
- Resource docs: `website/docs/r/<name>.html.markdown` → `azurerm_<name>`
- Data source docs: `website/docs/d/<name>.html.markdown` → `azurerm_<name>`

Also record the **doc type** from the path:
- `website/docs/r/**` => **Resource** documentation rules
- `website/docs/d/**` => **Data Source** documentation rules

### 2) Locate the schema in `internal/**`
- Search under `internal/**` for the Terraform name (e.g. `azurerm_<name>`).
- Open the relevant registration/implementation files until you find the schema definition.
- Record the schema file path(s) used.

If you cannot find the schema, say so explicitly and continue with a docs-only standards review.

### 3) Extract schema facts (from the schema definition)
From the schema, extract:
- required arguments
- optional arguments
- computed attributes
- ForceNew fields (`ForceNew: true`)
- constraints that affect docs (e.g. `ConflictsWith`, `ExactlyOneOf`, `AtLeastOneOf`, validations), if clearly visible

### 4) Audit the documentation for standards + parity

#### A) Formatting and structure
Validate:
- Frontmatter includes `subcategory`, `layout`, `page_title`, `description`
- H1 matches `azurerm_<name>`

**Resource vs Data Source hard rules (must enforce):**

- **Resources** (`website/docs/r/**`)
  - Must include: Example Usage, Arguments Reference, Attributes Reference, Import
  - Timeouts: required **only if** the resource schema defines timeouts (look for `Timeouts:` in the resource implementation)

- **Data Sources** (`website/docs/d/**`)
  - Must include: Example Usage, Arguments Reference, Attributes Reference
  - Must **not** include: Import
  - Timeouts: required **only if** the data source schema defines timeouts (look for `Timeouts:` in the data source implementation)

  **Timeouts link standard (new vs existing docs):**
- When a Timeouts section is present, validate the link uses the current format for **new** documentation pages:
  - `https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts`
- If the page uses the legacy Terraform.io link (for example `https://www.terraform.io/language/resources/syntax#operation-timeouts`):
  - If the docs file appears to be **newly added** in git (e.g. `git status` shows it as untracked/added), mark this as an **Issue** and fail the relevant standards check.
  - If the docs file already existed (modified but not newly added), record this as an **Observation** (existing pages may keep the older link for consistency).
  - If you cannot determine whether the file is new vs existing from the available context, default to **Observation**.

#### B) Arguments Reference parity and ordering
- All schema required args must be documented
- Documented args must exist in schema
- **Schema shape parity (block vs inline):** docs must match the schema's structural shape.
  - If schema defines an argument as a **nested block** (typically `TypeList`/`TypeSet` with `Elem: &Resource{Schema: ...}` and `MaxItems: 1` for single blocks), docs must describe it as a `... block` and include a section like: "A `${block}` block supports the following:" listing the nested fields.
  - If schema defines an argument as a **scalar/inline field** (`TypeString`/`TypeBool`/`TypeInt`/etc.), docs must not describe it as a block and must not document nested subfields under it.
  - If schema defines an argument as a **collection of primitives** (`TypeList`/`TypeSet` with `Elem: &Schema{Type: ...}`), docs should describe it as a list/set of values (not as a block with named subfields).
  - If schema defines an argument as a **map** (`TypeMap`), docs must describe it as a map and not as a block.
  - If docs describe `${arg}` as a block but schema indicates `${arg}` is an inline field (common when blocks have been flattened), mark as a parity failure and suggest updating the docs to reflect the flattened field shape.
- Argument ordering must follow `contributing/topics/reference-documentation-standards.md`:
  1. ID arguments first, with the last user-specified segment (usually `name`) first
  2. `location` (if present)
  3. remaining required arguments (alphabetical)
  4. optional arguments (alphabetical), with `tags` last (if present)
- **Resources only:** for every ForceNew field in schema, the argument description must end with a sentence of the form: "Changing this forces a new … to be created."
- **Data sources:** do not use "Changing this forces a new … to be created" wording (data sources do not create resources)
- If schema validations constrain values (e.g. `validation.StringInSlice`, `validation.IntBetween`), docs must include "Possible values …" using the standard phrasing.
- If schema defines a default value, docs must include "Defaults to `...`."

#### C) Attributes Reference parity
- All schema computed attributes must be present in Attributes Reference
- Ordering must follow `contributing/topics/reference-documentation-standards.md`: `id` first, then remaining attributes alphabetical
- Attribute descriptions must be concise and must not include possible/default values

#### D) Notes / note notation
- All note blocks must use the exact standard format: `(->|~>|!>) **Note:** ...`
- Flag invalid/legacy note styles (e.g. `Important:`, `NOTE:`, missing marker, wrong casing)
- **Semantic validation (marker must match meaning):** validate that the chosen marker is appropriate for what the note says.
  - `->` (informational): tips, extra context, recommendations, external links, clarifications that do not prevent errors or warn about irreversible impact.
  - `~>` (warning): guidance to avoid configuration errors or surprising behavior that is *reversible* (e.g. conditional requirements, conflicts, exactly-one-of, ForceNew behavior, API limitations that block create/update, deprecation/retirement where a configuration will error).
  - `!>` (caution): irreversible or high-impact guidance (e.g. data loss, permanent deletion, cannot be undone/disabled, security exposure with serious consequences).
  - **ForceNew-related guidance** should generally be `~> **Note:**` (do not use `->` for ForceNew warnings).
  - If a note’s content indicates one marker but another is used, mark **Note Notation** as fail and add an Issue suggesting the correct marker.
- Breaking changes should not be documented as notes (they belong in the changelog/upgrade guide)

#### E) Example Usage correctness
- Example must include all schema required args
- Example must not include a `terraform` or `provider` block
- Example should be functional and self-contained (no undefined references)
- Resource/data source instance name should generally be `example`
- Names in the example should generally be prefixed with `example-` (subject to service naming constraints)
- No hard-coded secrets (passwords/tokens/keys). Use `variable` with `sensitive = true` or a generator pattern.
- Example references must be internally consistent

#### F) Language
- Fix obvious grammar/spelling and consistency issues

#### G) Link hygiene
- Documentation links should be locale-neutral.
- Flag links containing locale path segments such as `/en-us/`, `/en-gb/`, `/de-de/`, etc.
- Suggested fix is to remove the locale segment (e.g. prefer `https://learn.microsoft.com/azure/...` over `https://learn.microsoft.com/en-us/azure/...`) unless there is a strong reason the localized link is required.

## ✅ Review output format (use this exact structure)

Output must be **rendered Markdown**.

- Do **not** wrap the review output in triple-backtick code fences.
- Use real headings, bullets, and bold text so it renders in chat.
- Use the section headings **exactly as written below** (including the emoji). Do not rename headings or remove emoji.


# 📋 **Docs Review**: ${terraform_name}

## 📌 **COMPLIANCE RESULT**
- **Status**: Valid / Invalid
- **Doc File**: ${docs_file_path}
- **Doc Type**: Resource / Data Source

## 🧾 **SCHEMA SNAPSHOT**
- **Schema File(s)**: ${schema_file_paths}
- **Required Args**: ${required_args}
- **Optional Args**: ${optional_args}
- **Computed Attributes**: ${computed_attrs}
- **ForceNew Fields**: ${force_new_fields}

## 📊 **DOC STANDARDS CHECK**
- **Frontmatter**: pass/fail + missing keys (if any)
- **Section Order**: pass/fail + missing sections (if any)
- **Argument Ordering**: pass/fail (ID args first, `location` next, then required alpha, then optional alpha, `tags` last)
- **Schema Shape**: pass/fail (docs describe blocks vs inline fields consistently with schema)
- **Attributes Coverage**: pass/fail (`id` first, computed attrs present, alphabetical)
- **ForceNew Wording**: pass/fail (resources only, missing “Changing this forces…” sentence)
- **Note Notation**: pass/fail (->/~>/!> exact format + marker meaning matches note content)
- **Link Locales**: pass/fail (no locale segments like `/en-us/` in URLs)
- **Examples**: pass/fail (functional/self-contained, no hard-coded secrets)

## 🟢 **STRENGTHS**
- ...

## 🟡 **OBSERVATIONS**
- ...

## 🔴 **ISSUES** (only actual problems)

### ${🔧/⛏️/❓} ${summary}
* **Priority**: 🔥 Critical / 🔴 High / 🟡 Medium / 🔵 Low / ✅ Good
* **Location**: ${doc_section_or_argument_name}
* **Schema Evidence**: ${what_in_schema_proves_this}
* **Problem**: clear description
* **Suggested Fix**: minimal edit/snippet

## 🛠️ **MINIMAL FIXES (PATCH-READY)**
Provide a minimal set of edits/snippets that fix all 🔴 Issues. Keep changes small and targeted.

## 🏆 **OVERALL ASSESSMENT**
One paragraph summary of what to change to become compliant.

### Notes
- Always cite the schema file path(s) you used.
- Prefer referencing doc section headings / argument names over line numbers.
- Do not invent schema fields; if schema cannot be located, explicitly say so and run a docs-only standards check.

### Individual Suggestions Format (legend)

**Priority System:** 🔥 Critical → 🔴 High → 🟡 Medium → 🔵 Low → ⭐ Notable → ✅ Good

**Review Type Icons:**
* 🔧 Change request - Standards/parity issues requiring fixes
* ❓ Question - Clarification needed about schema intent or doc meaning
* ⛏️ Nitpick - Minor style/consistency issues (typos, wording, formatting)
* ♻️ Refactor suggestion - Structural doc improvements (only when necessary)
* 🤔 Thought/concern - Potential mismatch or ambiguous behavior requiring discussion
* 🚀 Positive feedback - Excellent documentation patterns worth highlighting
* ℹ️ Explanatory note - Context about schema behavior or provider conventions
* 📌 Future consideration - Larger scope items for follow-up
