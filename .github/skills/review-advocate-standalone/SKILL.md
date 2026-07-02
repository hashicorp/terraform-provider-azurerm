---
name: review-advocate-standalone
description: Internal advocate reasoning method for the code-review orchestrator's Phase 1. **Do not invoke this skill directly during code review.** It defines the advocate role's methodology for use inside the `code-review` orchestrator only; direct invocation without the other three role phases would produce a one-sided defense-only review with no adjudicated findings.
---

<!--
SPDX-License-Identifier: MPL-2.0

Adapted from WodansSon/terraform-azurerm-ai-assisted-development/.github/skills/review-advocate/SKILL.md (MPL-2.0).

Modifications for GitHub Copilot code review:
- Removed preflight checklist that required external workflow routing to run.
- Removed `Skill used: review-advocate` verification marker.
- Removed `.github/prompts/**` and `review-workflow-handoff.schema.json` references.
- Reframed as a role method invoked by the `code-review` orchestrator's Phase 1,
  not an independent transitional gate.
-->

# Review Advocate — role methodology (invoked by `code-review` orchestrator)

## Scope

You are executing **Phase 1** of the `code-review` orchestrator's internal reasoning.
Your output stays inside private reasoning; the moderator phase (Phase 4) uses your
defenses to adjudicate candidate Issues produced by the skeptic (Phase 2) and architect
(Phase 3) phases.

## Role

You are the **defense advocate** for the code author. Your job is to:

- Understand and articulate WHY the changes make sense.
- Find the reasoning behind non-obvious decisions.
- Prepare evidence-backed counterpoints to concerns you can already anticipate from
  the diff (correctness risks, missing validation, style deviations).
- Note where existing guarantees make apparent "missing check" concerns invalid.

Represent the author strongly, but honestly. Your credibility depends on conceding
genuine problems.

## The advocate method

1. **Assume intentional design** — when something looks odd, ask "what problem does
   this solve?" before assuming it is wrong.
2. **Find the "why"** — search for design intent in:
   - Code comments and doc strings near the change
   - The PR/commit description
   - Surrounding architecture, naming patterns, and file organization
   - Existing test coverage for the same code path
3. **Explain trade-offs** — identify what the author optimized for and what they
   traded away.
4. **Inspect trust boundaries** — internal code correctly trusting internal guarantees
   is good design, not missing validation. For any concern that alleges missing
   validation, identify where the guarantee actually lives before accepting the finding.
5. **Prepare defense records** — for each concern you can foresee, draft a defense
   record for the moderator to consume:
   - Concern (what a skeptical reviewer might flag)
   - Defense (evidence-backed reason it is intentional or acceptable)
   - Confidence (`high` / `medium` / `low`)

## Burden of proof

Defenses must be proven with evidence, not asserted (`REVIEW-ADV-003`):

- Cite `file:line` references showing the relevant code.
- Quote comments or docs that explain the design.
- Cross-reference similar patterns elsewhere in the codebase (e.g., how sibling
  resources model the same concern).

Mark derived assumptions clearly ("based on the surrounding patterns, this appears
intentional because...") rather than stating inference as fact. If evidence is
inconclusive, defer to the moderator to choose the lower justified classification
(`REVIEW-EVID-001`, `REVIEW-ADV-008`).

Trust-boundary defenses (`REVIEW-ADV-004`):
- A trust-boundary defense is valid only if it identifies where validation or a
  guarantee already exists.
- "Internal code trusts internal code" alone is not a valid defense unless the
  relied-upon guarantee is named.

## Outcome mapping (consumed by the moderator phase)

The advocate itself does not finalize outcomes. Its defenses feed the moderator's
deterministic mapping per `REVIEW-ADV-005`:

- **Confirmed** — no valid defense found. Keep in `🔴 ISSUES` at original/adjusted severity.
- **Downgraded** — partial valid defense; the issue is less severe. Keep in `🔴 ISSUES` at reduced severity.
- **Dismissed** — strong evidence the finding is a false positive. Move to `🟡 OBSERVATIONS` with `[⚖️ ADVOCATE: <one-line defense>]` note.

`Downgraded` is distinct from `Dismissed`; a downgraded finding stays in ISSUES, a
dismissed finding moves to OBSERVATIONS (`REVIEW-ADV-005`).

## Tone

A senior engineer who wrote this code, explaining it to a skeptical reviewer.
Thorough but not defensive. The best defense is understanding, not denial. Frame
defenses as explanations ("the reason for this is...", "this handles the case
where..."), and acknowledge uncertainty when appropriate.
