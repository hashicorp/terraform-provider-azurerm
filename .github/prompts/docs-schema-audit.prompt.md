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

#### B) Arguments Reference parity and ordering
- All schema required args must be documented
- Documented args must exist in schema
- Required args listed first, then optional
- Alphabetical within each group
- `tags` last (if present)
- **Resources only:** for every ForceNew field in schema, docs must include: "Changing this forces a new resource."
- **Data sources:** do not use "Changing this forces a new resource" wording (data sources do not create resources)

#### C) Attributes Reference parity
- All schema computed attributes must be present in Attributes Reference
- Alphabetical unless a clear provider convention requires otherwise

#### D) Example Usage correctness
- Example must include all schema required args
- No hard-coded secrets (passwords/tokens/keys). Use `variable` with `sensitive = true` or a generator pattern.
- Example references must be internally consistent

#### E) Language
- Fix obvious grammar/spelling and consistency issues

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
- **Argument Ordering**: pass/fail (required first, alphabetical within groups, `tags` last)
- **Attributes Coverage**: pass/fail (computed attrs present)
- **ForceNew Notes**: pass/fail (missing “Changing this forces…” notes)
- **Examples**: pass/fail (required args present, no hard-coded secrets)

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
