# Maintainer Specific: Merging Pull Requests

> **Note:** All pull requests must be reviewed and approved before they are merged.

## Commit Type

All pull requests must be merged using the "Squash and merge" option. This ensures a clean commit history and simplifies either reverting changes or cherry-picking commits in the future.

## Commit Message Format

When merging a PR, the **commit message** should clearly describe the change being introduced. If the PR is correctly named (as described in [this guide](guide-opening-a-pr.md)), then the title can be used as-is. Otherwise, update the title to reflect the purpose of the PR in a way that will be meaningful in the Git history and use that as the message.

## Adding a changelog entry

When a PR is opened, it may require a changelog entry. While most PRs deserve a changelog entry, not every change should be included in the changelog as some have no user-facing impact. Some examples of PRs that should **not** be included are:

- Unit and acceptance test fixes
- Refactoring
- Documentation changes
- Deprecations (these must have an entry in the `{major}.0-upgrade-guide.html.markdown` file instead)

Otherwise, every PR that affects users should include a changelog entry file. This file should be located in the `.changelog` directory, and the file's name must be `<PR number>.hcl`.

To ensure consistency, and to group entries in the appropriate section, this repository uses a CLI tool to validate and add entries. To view the types of changes, run `make changelog-types`.

When adding a changelog entry, the following rules should be followed:

* Run `make changelog-types` to print a list of change types, use the most specific one. Each type provides an example format, this is validated using regex.
* Run `make changelog TYPE=<type> BODY='<body>' [PR=<num>]` to add a new entry, commit the resulting file to the PR. If `PR` is omitted, e.g. when adding an entry before opening a PR, a placeholder filename of `0.hcl` will be used. This will need to be renamed once the PR has been opened.
* Be consistent! Follow the formatting and language of the surrounding entries.
* Entries should start with a lower case, not end in a period, and always use the [serial (oxford) comma](https://en.wikipedia.org/wiki/Serial_comma).
* Each resource affected should be listed in full, i.e. do not use something like `azurerm_cosmosdb_*`.
* Entries should read as complete sentences such as ``add support for the `new_feature` property `` or ``improve validation of the `existing_feature` property ``, not ``support `new_feature` ``.

Here is a list of common changelog entries and how they should be formatted:

```
# X.YY.0 (Unreleased)

FEATURES:

* **New Action**: `azurerm_action` [GH-12345]
* **New Data Source**: `azurerm_data_source` [GH-12345]
* **New List Resource**: `azurerm_resource` [GH-12345]
* **New Resource**: `azurerm_resource` [GH-12345]

ENHANCEMENTS:

* dependencies: `go-azure-sdk` - update to `v0.20250101.1123456` [GH-12345]
* dependencies: `service` - update API version to `2021-12-01` [GH-12345]
* Data Source: `azurerm_data_source` - export the `value` attribute [GH-12345]
* `azurerm_resource` - the `sku` property can now be updated to `Basic` or `Standard` without recreating the resource [GH-12345]
* `azurerm_resource` - add support for the `thing1` property [GH-12345]
* Action: `azurerm_action` - add support for the `thing1` property [GH-12345]
* List Resource: `azurerm_resource` - add support for the `thing1` property [GH-12345]
* `azurerm_resource` - add support for the `block1` block [GH-12345]
* `azurerm_resource` - add support for the `thing2`, `thing3`, and `thing4` properties [GH-12345]
* `azurerm_resource` - add support for the `block2`, `block3`, and `block4` blocks [GH-12345]
* `azurerm_resource` - improve validation for the `termination_nofication.timeout` property [GH-12345]

BUG FIXES:

* provider: `provider_feature` now correctly reads the default value from an environment variable [GH-12345]
* Data Source: `azurerm_data_source` - prevent a possible crash by setting `queue_name` correctly [GH-12345]
* Data Source: `azurerm_data_source` - correctly populate the `kind` and `os_type` attributes [GH-12345]
* `azurerm_data_factory_dataset_delimited_text` - set defaults properly for `column_delimiter`, `quote_character`, `escape_character`, `first_row_as_header`, and `null_value` [GH-12345]
* `azurerm_linux_function_app` - correctly deduplicate user `app_settings` [GH-12345]
* `azurerm_windows_function_app_slot` - correctly deduplicate user `app_settings` [GH-12345]
```
