<!--
SPDX-License-Identifier: MPL-2.0

Portions of the review methodology encoded here are adapted from
WodansSon/terraform-azurerm-ai-assisted-development (MPL-2.0).
Runtime adaptations for GitHub Copilot code review by the terraform-provider-azurerm
fork maintainer. See .github/skills/*/SKILL.md for the multi-role review orchestration.
-->

# GitHub Copilot instructions

This repository is the Azure Resource Manager provider for Terraform, written in Go.
The code under review here is production infrastructure code used by real users; treat
correctness, backwards compatibility, and clarity as first-class concerns.

## Preferred review approach

When reviewing pull requests, Copilot **should invoke the `code-review` skill**
(see `.github/skills/code-review/SKILL.md`). The `code-review` skill orchestrates a
multi-role internal reasoning pass — advocate, skeptic, architect, moderator — and
emits a synthesized final review. This is the intended review flow for this repository.

- The `code-review` skill is the top-level entry point for PR review.
- The four standalone role skills (`review-advocate-standalone`, `review-skeptic-standalone`,
  `review-architect-standalone`, `review-moderator-standalone`) are invoked by the
  `code-review` skill. They are not intended to be discovered and run independently
  during PR review; their descriptions include an explicit self-limitation.

## What counts as review scope

- **Primary scope**: files under `internal/**` (Go implementation and acceptance tests).
- **Docs**: files under `website/docs/**/*.html.markdown` (reference documentation).
- **Non-actionable scope**: files under `vendor/**` (vendored third-party code — do not
  raise Issues that ask contributors to hand-edit vendored content; instead, point to the
  dependency update or generation input that introduces the change).

## Provider-specific expectations

- New resources typically need companion artifacts: Resource Identity, list resource
  (unless the change uses the maintainer-reviewed exception path such as
  `allow-without-list` / `list-not-supported`), documentation under `website/docs/`,
  and acceptance tests.
- Ephemeral resources (`*_ephemeral.go`) need docs under `website/docs/ephemeral-resources/`
  and Terraform 1.10-gated tests.
- Provider-defined functions under `internal/provider/function/` need docs under
  `website/docs/functions/` and Terraform 1.8-gated unit tests.
- Prefer the typed `internal/sdk` framework for new resources. Untyped Plugin SDK
  patterns are maintenance-only.
- Do not flag "internal code trusts internal code" as missing validation unless the
  concern identifies where the relied-upon guarantee actually lives.

## Evidence discipline

- Every Issue must cite `file:line` from the actual diff.
- If evidence is inconclusive, choose the lower justified classification
  (Observation over Issue, Downgraded over Confirmed).
- Do not narrate intermediate reasoning steps in the visible review; present only
  final evidence-backed conclusions.
- Do not invent policies that no workspace document, instruction file, or skill supports.

## What review MUST NOT do

- Do not require reviewers to run local binaries (e.g., `azurerm-linter`, `git`, `gh`) —
  the GitHub Copilot code review environment cannot execute them. If a rule in an
  instructions file references such tools, treat that rule as informational context, not
  as a required verification step.
- Do not emit workflow bookkeeping markers such as `Skill used: xxx` in review output.
- Do not add review sections beyond `🔴 ISSUES`, `🟡 OBSERVATIONS`, and (optionally)
  `🟢 STRENGTHS` unless the change explicitly warrants it.
