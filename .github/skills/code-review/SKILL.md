---
name: code-review
description: Multi-role code review for pull requests in this Terraform provider. Executes an internal four-phase reasoning pass (advocate → skeptic → architect → moderator) within a single review, then emits a synthesized final review with 🔴 ISSUES, 🟡 OBSERVATIONS, and optionally 🟢 STRENGTHS. Use for any pull request review of Go code, acceptance tests, or reference documentation in this repository.
---

<!--
SPDX-License-Identifier: MPL-2.0

Adapted orchestration pattern; the four role sub-skills are adapted from
WodansSon/terraform-azurerm-ai-assisted-development (MPL-2.0), where they run under
VS Code prompt orchestration. This master skill provides equivalent orchestration
inside a single GitHub Copilot code review invocation.
-->

# Multi-role code review orchestrator

You are performing a code review as the top-level `code-review` skill. Execute the
review in four internal phases. **Only Phase 4 output becomes visible review comments.**
Do not narrate the phases; the reader sees only the final synthesized review body.

## Governing rules

Apply the shared code review rules for this repository throughout all four phases:

- `.github/instructions/code-review-shared-rules.instructions.md` — `REVIEW-*` rule taxonomy
- `.github/copilot-instructions.md` — repo-wide expectations and scope

Provider-specific companion guidance for Go files under `internal/**` (load if present):
- `implementation-guide.instructions.md` — typed vs untyped model identification, CRUD patterns
- `schema-patterns.instructions.md` — field types, validation patterns
- `azure-patterns.instructions.md` — Azure SDK integration
- `testing-guidelines.instructions.md` — acceptance test discipline

## Phase 1 — Advocate (private reasoning)

Apply the method defined in `.github/skills/review-advocate-standalone/SKILL.md`.

- Assume intentional design; search for the "why" in comments, PR description,
  surrounding architecture, naming patterns, and test coverage.
- Inspect trust boundaries: internal code correctly trusting internal guarantees is
  good design, not missing validation. Identify where the guarantee lives before
  accepting a "missing check" finding.
- Produce draft defenses for likely concerns you can already foresee.

**Do not emit Phase 1 output.** Store internally as advocate reasoning notes.

## Phase 2 — Skeptic (private reasoning)

Apply the method defined in `.github/skills/review-skeptic-standalone/SKILL.md`.

Walk the six attack classes per `REVIEW-SKEP-003`:
1. Correctness and logic errors
2. Error handling, nil, zero-value defaults
3. Concurrency and ordering
4. Input validation and trust boundaries
5. Resource lifecycle and residual state (PATCH residual state, `CustomizeDiff` placement, Linux/Windows parity)
6. Security exposure
7. Test-coverage gaps for behavior-changing branches

For each candidate finding, construct a reproducible failure path citing `file:line`.
If you cannot construct one, demote it to an Observation.

**Do not emit Phase 2 output.** Store internally as candidate finding records per
`REVIEW-HANDOFF-001`.

## Phase 3 — Architect (private reasoning)

Apply the method defined in `.github/skills/review-architect-standalone/SKILL.md`.

Walk the direction areas per `REVIEW-ARCH-003`:
- Schema shape and field naming
- Argument grouping, singular vs plural naming
- Resource decomposition and singleton modeling
- Typed vs untyped implementation approach (`internal/sdk` framework vs Plugin SDK)
- Cross-resource and cross-platform consistency
- Required companion artifacts (Resource Identity, list resources, ephemeral resources,
  provider-defined functions) per `REVIEW-SCOPE-002` / `REVIEW-SCOPE-003`
- Overall maintainability and diff readability

Escalate a design concern to a candidate Issue **only** when a mandatory source
supports the rule (`REVIEW-ARCH-004`). Otherwise record as Observation
(`REVIEW-OBS-001` / `REVIEW-ARCH-002`).

**Do not emit Phase 3 output.** Store internally.

## Phase 4 — Moderator synthesis (published output)

Apply the method defined in `.github/skills/review-moderator-standalone/SKILL.md`.

For each candidate finding produced by Phase 2 and Phase 3, apply the advocate
adjudication mapping (`REVIEW-ADV-005`) using the Phase 1 defenses:

- **Confirmed** — no valid defense found → keep in `🔴 ISSUES` at original or adjusted severity.
- **Downgraded** — partial valid defense → keep in `🔴 ISSUES` at reduced severity.
- **Dismissed** — strong evidence the finding is a false positive or intentional design →
  move to `🟡 OBSERVATIONS` with `[⚖️ ADVOCATE: <one-line defense>]` annotation.

Then merge duplicates and normalize wording per `REVIEW-MOD-003` / `REVIEW-MOD-004`.

### Emit exactly this output shape

```markdown
### 🔴 **ISSUES**

<!-- Each entry: severity + file:line + concrete failure or violation + fix direction. -->

### 🟡 **OBSERVATIONS**

<!-- Non-blocking design concerns, dismissed candidates with advocate annotation,
     follow-up ideas outside change-set scope. -->

### 🟢 **STRENGTHS** (optional, only when concretely evidenced)

<!-- Factual, evidence-backed positives. Skip if none apply. -->
```

## Guardrails for the moderator phase

- **No candidate silently dropped.** Every Phase 2/3 candidate must resolve to
  `Confirmed`, `Downgraded`, or `Dismissed` per `REVIEW-ADV-006`.
- **One concern, one classification.** No item appears in both `🔴 ISSUES` and
  `🟡 OBSERVATIONS` (`REVIEW-CLASS-004`).
- **Every Issue cites `file:line`** (`REVIEW-OUT-003`).
- **No phase narration.** Do not say "in the advocate phase I found..." or
  "the skeptic proposed...". The reader sees findings, not workflow.
- **No bookkeeping markers.** Do not emit `Skill used: xxx` or similar workflow traces.
- **Inconclusive evidence chooses the lower classification** (`REVIEW-EVID-001`,
  `REVIEW-ADV-008`): prefer Observation over Issue, Downgraded over Confirmed, Dismissed
  over Downgraded, only when evidence supports it.
