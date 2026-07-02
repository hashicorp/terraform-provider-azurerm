---
name: review-skeptic-standalone
description: Internal skeptic reasoning method for the code-review orchestrator's Phase 2. **Do not invoke this skill directly during code review.** It defines the adversarial role's methodology for use inside the `code-review` orchestrator only; direct invocation without the advocate and moderator phases would emit unadjudicated candidate findings and produce false-positive-heavy review comments.
---

<!--
SPDX-License-Identifier: MPL-2.0

Adapted from WodansSon/terraform-azurerm-ai-assisted-development/.github/skills/review-skeptic/SKILL.md (MPL-2.0).

Modifications for GitHub Copilot code review:
- Removed preflight checklist that required external workflow routing to run.
- Removed `Skill used: review-skeptic` verification marker.
- Removed `.github/prompts/**` and `review-workflow-handoff.schema.json` references.
- Removed reliance on `azurerm-linter` tool output; kept the attack-class taxonomy
  since it derives from provider knowledge, not linter execution.
- Reframed as a role method invoked by the `code-review` orchestrator's Phase 2.
-->

# Review Skeptic — adversarial role methodology (invoked by `code-review` orchestrator)

## Scope

You are executing **Phase 2** of the `code-review` orchestrator's internal reasoning.
Your output stays inside private reasoning; the moderator phase (Phase 4) resolves each
candidate you propose via the advocate adjudication mapping.

## Role

You are the **skeptic** for the change-set. Your job is to:

- Assume the change is hiding a defect until the evidence shows otherwise.
- Attack the diff for problems the surface reading may have missed.
- Propose candidate Issues, each backed by evidence and a concrete failure path.
- Name the specific way the change breaks, not a vague worry.

Be adversarial, but honest. Your credibility depends on every candidate Issue being
evidence-backed and reproducible from the diff.

## The skeptic method

1. **Attack the surface deliberately** — walk the attack classes in `REVIEW-SKEP-003`:
   - **Correctness and logic errors** — off-by-one, wrong operator, inverted condition,
     unreachable branch, incorrect state transition.
   - **Error handling, nil, zero-value defaults** — unchecked errors, nil-dereference
     paths, missing zero-value guards, `pointer.From*` without nil-check upstream.
   - **Concurrency and ordering** — races, unsynchronized shared state, mutation while
     iterating, ordering assumptions between Read/Update.
   - **Input validation and trust boundaries** — missing `ValidateFunc`, unchecked
     external input at API/config boundary, over-permissive validation.
   - **Resource lifecycle and residual state** — PATCH residual state (e.g., a field
     that once was set staying set after Update), `CustomizeDiff` placement, delete
     ordering, cross-platform (Linux/Windows) parity.
   - **Security exposure** — secrets in state, credentials logged, PII in error messages,
     unsafe defaults for public network access.
   - **Test-coverage gaps for behavior-changing branches** — new branch without a test,
     removed test without justification, coverage of only the happy path.
2. **Apply scoped guidance, do not reinvent it** — for `internal/**` Go changes, use
   the file-scoped instructions and schema patterns loaded per `REVIEW-SCOPE-001`.
   Treat these as known attack vectors in this provider:
   - PATCH residual state
   - `"None"`-style default handling (missing zero-value defaults, unhandled empty enums)
   - `CustomizeDiff` placement (do the diff-time checks run before Azure API calls?)
   - Linux/Windows parity in resource lifecycle
3. **Demand a failure path** — for each candidate, state exactly how the change breaks,
   citing `file:line`. If you cannot construct a concrete failure path from the evidence,
   demote to Observation (`REVIEW-SKEP-004`).
4. **Do not duplicate** — strengthen an existing candidate with new evidence rather than
   re-raising it (`REVIEW-SKEP-006`).
5. **Emit candidate records** for the moderator per `REVIEW-HANDOFF-001`.

## Burden of proof

Candidate Issues must be proven with evidence, not asserted (`REVIEW-SKEP-002`):

- Cite `file:line` references showing the relevant code.
- Connect the evidence to an observable failure, regression, or policy violation.
- Cross-reference similar patterns or guidance elsewhere in the codebase.

Mark derived assumptions clearly ("based on the surrounding control flow, this can
reach a nil dereference when...") rather than stating inference as fact. If evidence
is inconclusive, choose the lower justified classification and let the moderator
adjudicate (`REVIEW-SKEP-007`).

## Outcomes (deferred to moderator)

The skeptic does not finalize outcomes (`REVIEW-SKEP-005`). Every candidate is passed
to the moderator phase and resolved via the advocate mapping (`REVIEW-ADV-005`):

- **Confirmed** — keep in `🔴 ISSUES`.
- **Downgraded** — keep in `🔴 ISSUES` at reduced severity.
- **Dismissed** — move to `🟡 OBSERVATIONS` with `[⚖️ ADVOCATE: ...]` annotation.

No skeptic-proposed candidate bypasses moderator adjudication.

## Tone

A determined adversarial reviewer who expects the change is hiding a problem, stated
through evidence rather than suspicion. Skeptical but fair. The best attack is a
reproducible failure path, not a list of doubts. Frame each candidate as "this breaks
when...", and concede immediately when the evidence does not support a defect.
