# Maintainer Specific: Updating the Changelog

> **Note:** When sending a Pull Request you should not include a changelog entry as a part of the Pull Request - this is to avoid conflicts. Contributors should not be concerned with updating the changelog as that is something only maintainers will do during merge.

When a PR is merged it may or may not be included in the changelog. While most PRs deserve a changelog entry not every change should be included in the changelog as some have no user-facing impact. Some examples of PRs that should **not** be included are:

- Unit and acceptance test fixes
- Refactoring
- Documentation changes

Otherwise, every PR that affects users should be added to the appropriate section:

* `FEATURES` - new resources and data sources
* `ENHANCEMENTS` - new properties, functionality, and features (including SDK/API upgrades)
* `BUG FIXES` - bug fixes

When adding a changelog entry, the following rules should be followed:

* Be consistent! Follow the formatting and language of the surrounding entries.
* Entries should start with a lower case, not end in a period, and always use the [serial (oxford) comma](https://en.wikipedia.org/wiki/Serial_comma).
* Each resource affected should be listed in full, i.e. do not use something like `azurerm_cosmosdb_*`.
* Each entry should link to the pull request with the placeholder `[GH-{number}]` (e.g. `[GH-1234]`), this will be replaced with a link during the release process.
* Entries should read as complete sentences such as ``add support for the property `new_feature` `` or ``improve validation of the property `old_feature` `` not ``support `new_feature` ``.

And finally, when making the edit commit, the PR number should be included in the commit message so the edit is linked to the PR, and the entry from the pr. For example `CHANGELOG.md for #1234`.

Here is a list of common changelog entries and how they should be formatted:

```
# X.YY.0 (Unreleased)

FEATURES:

* **New Data Source**: `azurerm_data_source` [GH-12345]
* **New Resource**: `azurerm_resource` [GH-12345]

ENHANCEMENTS:

* dependencies: `go-azure-sdk` - update to `v0.20250101.1123456` [GH-12345]
* dependencies: `service` - update API version to `2021-12-01` [GH-12345]
* Data Source: `azurerm_data_source` - export the `value` attribute [GH-12345]
* `azurerm_resource` - the `sku` property can now be updated to `Basic` or `Standard` without recreating the resource [GH-12345]
* `azurerm_resource` - add support for the `thing1` property [GH-12345]
* `azurerm_resource` - add support for the `thing2`, `thing3`, and `thing4` properties [GH-12345]
* `azurerm_resource` - improve validation for the `timeout` property within the `termination_notification` block [GH-12345]

BUG FIXES:

* Data Source: `azurerm_data_source` - prevent a possible crash by setting `queue_name` correctly [GH-12345]
* Data Source: `azurerm_data_source` - correctly populate the `kind` and `os_type` attributes [GH-12345]
* `azurerm_data_factory_dataset_delimited_text` - set defaults properly for `column_delimiter`, `quote_character`, `escape_character`, `first_row_as_header`, and `null_value` [GH-12345]
* `azurerm_linux_function_app` - correctly deduplicate user `app_settings` [GH-12345]
* `azurerm_windows_function_app_slot` - correctly deduplicate user `app_settings` [GH-12345]
```
