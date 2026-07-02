---
applyTo: "internal/**/*.go"
description: "Shared code review rules for Go code in this Terraform provider. Defines the REVIEW-* rule taxonomy consumed by the code-review skill and its role sub-skills. Loaded automatically for internal/**/*.go review scope."
---

<!--
SPDX-License-Identifier: MPL-2.0

Adapted from WodansSon/terraform-azurerm-ai-assisted-development
(.github/instructions/code-review-compliance-contract.instructions.md), MPL-2.0.

Modifications for GitHub Copilot code review:
- Removed rules that require local tool execution (`azurerm-linter`, `git`, `gh`).
- Removed prompt-file references (`.github/prompts/**`) — GitHub PR review does not
  load prompt files.
- Simplified REVIEW-HANDOFF-* to describe an internal reasoning shape rather than a
  runtime JSON schema (`review-workflow-handoff.schema.json` is not shipped here).
- Retained the terraform-provider-azurerm-specific REVIEW-SCOPE-* rules verbatim in
  intent, since they encode current maintainer expectations.
-->

# Code review shared rules

This file is the single source of truth for the `REVIEW-*` rule taxonomy consumed by
`.github/skills/code-review/SKILL.md` and its four role sub-skills.

## Rule ID format

`REVIEW-<AREA>-<NNN>`

Areas:
- `EVID` — evidence and verification guardrails
- `CLASS` — finding classification
- `SCOPE` — file-type-specific coverage
- `HANDOFF` — internal role-to-role reasoning shape
- `OBS` — observation-only design guidance
- `OUT` — required review output semantics

## Evidence hierarchy

When a claim affects correctness, severity, or merge readiness, weigh evidence in this order:

1. The changed files and the actual diff under review
2. Repo-level contributor guidance (`CONTRIBUTING.md`, `contributing/README.md`, `contributing/topics/**`)
3. File-scoped instructions and skills (`.github/instructions/**`, `.github/skills/**`)
4. Surrounding workspace code, tests, and patterns
5. PR/commit description and code comments that state design intent
6. External references for semantics only, when workspace evidence is insufficient

If evidence is missing for a claim that would change severity or requested action, do not guess.

# Contract Rules

## Evidence and verification

### REVIEW-EVID-001: Do not guess when evidence is required
- Rule: If a compliance-relevant or correctness-relevant claim cannot be backed by available evidence, do not invent it.
- Behavior: Downgrade to an Observation, ask for clarification, or explicitly state that evidence could not be proven.

### REVIEW-EVID-002: Attribute policies to real sources
- Rule: Do not claim a style or implementation rule is mandatory unless it is supported by a current contributor document, instruction file, skill, or this file.
- Behavior: Avoid invented policy language such as "must" or "required" when the source only supports a preference.

### REVIEW-EVID-003: Discover contributor-guidance paths before claiming absence
- Rule: Do not assume repo-level contributor guidance always lives at `CONTRIBUTING.md`. Check common locations such as `CONTRIBUTING.md`, `contributing/README.md`, and `contributing/topics/**` before claiming guidance is absent.

### REVIEW-EVID-004: Perform verification silently
- Rule: Do not narrate intermediate verification steps such as reading files or comparing content. The visible review should present only final evidence-backed conclusions.

### REVIEW-EVID-005: Every review invocation is a fresh audit
- Rule: Do not reuse prior review conclusions from earlier turns. Base findings on the diff under review in this invocation.

## Finding classification

### REVIEW-CLASS-001: Issues are for actual problems only
- Rule: An Issue must be a real defect, regression, policy violation, missing requirement, or correctness risk with evidence.
- Rule: Do not place stylistic preferences or speculative concerns in Issues.

### REVIEW-CLASS-002: Observations are non-blocking
- Rule: Observations capture design concerns, preferences, uncertainty, or follow-up ideas that are not clearly blocking.
- Rule: If the current implementation is acceptable under the available evidence, keep it out of Issues even if another design might be preferable.

### REVIEW-CLASS-003: Strengths must be factual
- Rule: Strengths call out concrete, evidenced positives. Do not use Strengths to pad the review with generic praise.

### REVIEW-CLASS-004: One finding, one classification
- Rule: The same underlying concern must not appear in both Observations and Issues.
- Rule: If severity is uncertain, choose the lower justified classification and explain why.

### REVIEW-CLASS-005: Fixes must be deterministic
- Rule: Each Issue should point to a single, concrete correction path. Do not present multiple alternative fixes unless explicitly asked.

## Internal reasoning handoff (for the code-review skill's four-phase orchestration)

### REVIEW-HANDOFF-001: Internal reasoning uses a stable finding shape
- Rule: Within the `code-review` skill's four-phase orchestration (advocate → skeptic → architect → moderator), each internal candidate finding produced by any phase should carry:
  - `id` — a stable short identifier the moderator can reference (e.g., `SKEP-001`, `ARCH-002`).
  - `title` — a one-line summary of the concern.
  - `scope` — the `file:line` or `file:line-range` under review.
  - `severity` — `high` / `medium` / `low`, with the classification's justification.
  - `evidence` — the diff excerpt, file quote, or cross-reference supporting the finding.
  - `reasoning` — why this evidence supports the finding.
  - `status` — `candidate` (from skeptic/architect) or `observation` (from architect defaults).
  - `roles` — which internal role(s) proposed or enriched this record.
- Rule: This is an internal reasoning shape only; the moderator translates it into the final visible ISSUES/OBSERVATIONS structure. Do not emit these fields to the reviewer.

### REVIEW-HANDOFF-002: Status transitions are stage-aware
- Rule: Before the moderator phase, allowed statuses are `candidate` and `observation`.
- Rule: The moderator phase resolves each `candidate` to `confirmed`, `downgraded`, or `dismissed`, then translates the final set into the visible ISSUES/OBSERVATIONS/STRENGTHS output.

### REVIEW-HANDOFF-003: Enrich, do not duplicate
- Rule: When multiple phases touch the same concern, enrich one record rather than clone duplicates.
- Rule: The skeptic must not restate advocate observations; it must add net-new evidence or net-new candidates.

## Observation-only design guidance

### REVIEW-OBS-001: Design preference is Observation-only by default
- Rule: A design preference, stylistic improvement, or "a different shape might be better" concern without a mandatory source stays an Observation.
- Rule: Do not use invented "must" or "required" language when the source only supports a preference.

## File-type-specific review coverage (terraform-provider-azurerm)

### REVIEW-SCOPE-001: Go implementation and test scope
- Rule: Files under `internal/**/*.go` and `internal/**/*_test.go` are primary review scope for correctness, provider conventions, and test discipline.
- Rule: Apply file-scoped companion guidance: `implementation-guide.instructions.md`, `schema-patterns.instructions.md`, `azure-patterns.instructions.md`, `testing-guidelines.instructions.md`, and any implementation-compliance / testing-compliance contracts when they are present in this repository. If they are absent, apply the maintainer-visible patterns already in place in surrounding code.

### REVIEW-SCOPE-002: New resources need required companion artifacts
- Rule: When the review scope adds a brand-new resource under `internal/**/*.go`, check whether the required companion artifacts are present or explicitly justified:
  - Resource Identity support
  - List resource (unless the maintainer-reviewed exception path is used, e.g., `allow-without-list`, `list-not-supported`)
  - List-resource query tests when a list resource is required
  - Documentation under `website/docs/` (and `website/docs/list-resources/` when a list resource exists)
- Rule: Missing companion artifacts without explicit exception justification is a reviewable Issue.

### REVIEW-SCOPE-003: Ephemeral resources and provider-defined functions need companions
- Rule: New `*_ephemeral.go` implementations need service registration, docs under `website/docs/ephemeral-resources/`, and Terraform 1.10-gated tests under `*_ephemeral_test.go`.
- Rule: New provider-defined functions under `internal/provider/function/` need docs under `website/docs/functions/` and Terraform 1.8-gated unit tests.
- Rule: Missing companions is a reviewable Issue.

### REVIEW-SCOPE-004: Singleton or get-only resources — exception-aware list review
- Rule: Do not raise a generic "missing list resource" Issue when the implementation evidence shows singleton child modeling (fixed child path, synthetic singleton ID type, CRUD operating on parent+fixed segment).
- Rule: If implementation evidence supports singleton behavior but the change context does not explicitly justify omitting the list resource, keep the finding but frame it as: "document the reason and use the maintainer-reviewed exception path (`allow-without-list` / `list-not-supported`)" rather than blocking on missing list.

### REVIEW-SCOPE-005: Generic lifecycle logging adds no unique value
- Rule: Additions of generic lifecycle logging such as `Import check`, `Creating`, `Reading`, `Updating`, `Deleting` are reviewable Issues when they only duplicate Terraform-core or provider-native logging.
- Rule: Do not require removal of targeted `not-found` / `removing-from-state` diagnostics when they add distinct debugging value.

### REVIEW-SCOPE-006: Vendored files are non-actionable
- Rule: Files under `vendor/**` are non-actionable for normal review. Identify them in the scope summary if present, but do not raise Issues asking contributors to hand-edit vendored content.
- Rule: When a correctness concern originates from vendored content, point to the first actionable non-vendored source (dependency update, generation input, service client wiring).

## Output shape

### REVIEW-OUT-001: Emit only the three canonical sections
- Rule: The visible review body uses only `🔴 ISSUES`, `🟡 OBSERVATIONS`, and optionally `🟢 STRENGTHS`. No other sections.
- Rule: Do not narrate the internal advocate/skeptic/architect/moderator phases in the visible output.

### REVIEW-OUT-002: Advocate-dismissed findings surface in OBSERVATIONS with rationale
- Rule: When the moderator dismisses a candidate Issue based on advocate defense, the finding appears in `🟡 OBSERVATIONS` with a brief `[⚖️ ADVOCATE: <one-line rationale>]` annotation.

### REVIEW-OUT-003: Every Issue cites `file:line`
- Rule: Every entry in `🔴 ISSUES` must include a `file:line` (or `file:line-range`) evidence pointer.

<!-- REVIEW-SHARED-RULES-EOF -->
