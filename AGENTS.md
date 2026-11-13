# AGENTS.md

This file provides guidance to AI Agents working in this repository with basic context and guidelines needed for `terraform-provider-azurerm` development.

### Contributing

When implementing a new feature, fixing a bug, or modifying the provider in any other way, first reference any relevant topic documentation for best practices and requirements. See [AzureRM Provider Contributor Guides - Topics](contributing/README.md#topics).

## Code Review

When performing a code review, apply **all** coding practices defined under the [/contributing](/contributing) folder.

## Tests

By default, most AI models generate extensive unit tests for even the most trivial change, these tests are frequently meaningless when it comes to testing the provider's logic.
These tests should not be included as they add an unnecessary maintenance and review burden.

Examples of unit tests that the maintainers of this repository will not accept:

- Tests that test CRUD methods using mocked clients and hardcoded API request/response bodies.
- Tests that validate internal SDK functions (e.g. tests validating `ValidateFunc`'s such as `validation.StringInSlice` return the correct error message).
- Tests that assert schema fields are as defined in the schema.

When adding tests, first reference our [acceptance testing guide](contributing/topics/reference-acceptance-testing.md), the majority of the provider's logic is tested using real Terraform operations (plan, apply, import, etc.).

When adding a unit test, the reasoning for including it **must** be described in the PR description or in a comment in the unit test file itself if the agent is not opening the PR.

## AI Usage

AI usage must be disclosed in pull requests, the [PR template](.github/pull_request_template.md) contains a checkbox that must be ticked.

For the complete AI usage policy, see [topics/ai-usage](contributing/topics/ai-usage.md).
