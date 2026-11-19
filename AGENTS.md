# AGENTS.md

This file provides guidance to AI Agents working in this repository with basic context and guidelines needed for `terraform-provider-azurerm` development.

## Contributing

When contributing, apply **all** coding practices defined in the following guides:

- [Best Practices](contributing/topics/best-practices.md)
- [Naming Conventions](contributing/topics/reference-naming.md)
- [Error Handling](contributing/topics/reference-errors.md)

### Code Additions

When implementing a new feature, fixing a bug, or modifying the provider in any other way, first reference **only** relevant topic documentation for best practices and requirements. See [AzureRM Provider Contributor Guides - Topics](contributing/README.md#topics).

When writing code, favour readability. Separate statements, or distinct groupings of statements, with blank lines. Avoid flattening large structs, maps, or other nested types into a single line.

### Code Review

In addition to the guides linked above, see [AzureRM Provider Contributor Guides - Topics](contributing/README.md#topics) and reference **only** guides that apply to the change introduced in the PR.

CI already checks formatting, lint, copyright headers, spelling, generated code, and that a new resource has a list resource. Spend the review on what it cannot: is the resource registered, does every property get sent on create and update and read back, is this a breaking change, does it need a state migration.

## Tests

When adding tests, first reference our [acceptance testing guide](contributing/topics/reference-acceptance-testing.md), the majority of the provider's logic is tested using real Terraform operations (plan, apply, import, etc.).

Running acceptance tests costs real money. Ask to run them, if permitted run a narrow selection of tests covering the change instead of tests across whole service packages.

By default, most AI models generate extensive unit tests for even the most trivial change, these tests are frequently meaningless when it comes to testing the provider's logic.
These tests should not be included as they add an unnecessary maintenance and review burden.

Examples of unit tests that the maintainers of this repository **will not** accept:

- Tests that test CRUD methods or flatten/expand functions using mocked clients and hardcoded API request/response bodies.
- Tests that validate internal SDK functions (e.g. tests validating `ValidateFunc` functions such as `validation.StringInSlice` return the correct error message).
- Tests that assert schema fields are as defined in the schema.

When adding a unit test, the reasoning for including it **must** be described in a comment in the unit test file.

## AI Usage

Do not open pull requests or write pull request descriptions.

For the complete AI usage policy, see [topics/ai-usage](contributing/topics/ai-usage.md).
