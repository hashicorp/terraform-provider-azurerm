# Opening a PR

Firstly all contributions are welcome!

There is no change too small for us to accept and minor formatting, consistency and documentation PRs are very welcome! However, before making any large or structural changes it is recommended to seek feedback (preferably by reaching out in our community slack) to prevent wasted time and effort. We may already be working on a solution, or have a different direction we would like to take.

If you are ever unsure please just reach out, we are more than happy to guide you in the right direction!

## Considerations

As a general rule, the smaller the PR the quicker it's merged - as such when upgrading an SDK and introducing new properties we'd ask that you split that into multiple smaller PRs, for example if you were planning on updating an SDK to add a new resource and update an existing one we would prefer `3` separate PRs:

1. Update the Cosmos DB SDK to use API Version `2022-02-02` from `2020-01-01`.
2. Add the new property `new_feature` to the `azurerm_cosmosdb_*` resources.
3. Introduce the New Resource `azurerm_cosmosdb_resource`.

In the same vein, please do not reorder existing properties in the same PR, as doing so degrades diff readability. Make any reordering changes in a separate follow-up PR instead.

We also recommend not opening a PR based on your `main` branch. By doing this any changed pushed to the PR may inadvertently be also pushed to your `main` branch without warning.

Due to the high volume of PRs on the project and to ensure maintainers are able to focus on changes which are ready to review, please do not open Draft PRs or work that is not yet ready to be reviewed.

## Process

Pull Requests generally go through a number of phases which vary slightly depending on what's being changed.

The following guides cover the more common scenarios we see:

* [Extending an existing Resource](guide-new-fields-to-resource.md)
* [Extending an existing Data Source](guide-new-fields-to-data-source.md)
* [Adding a new Resource](guide-new-resource.md)
* [Adding a new Data Source](guide-new-data-source.md)
* [Adding a new Service Package](guide-new-service-package.md)

In general, Pull Requests which add/change either code or SDK's go through the following steps:

1. Make / commit the changes.
2. Run GitHub Actions linting and checks locally with the make command `make pr-check`.
3. Run all relevant [Acceptance Tests](running-the-tests.md).
4. Add a changelog entry if the change is user-facing (see [Adding a changelog entry](maintainer-merging.md#adding-a-changelog-entry)).
5. Open a Pull Request (see below on `What makes a good PR?`).
6. GitHub actions will trigger and run all linters.
7. A Maintainer will review the PR and also run acceptance tests against our test subscription.
8. Once all comments have been addressed and tests pass the PR will be merged.

## What makes a good PR?

* Don't change your forked repo's `main` branch, instead, make a feature branch.
* The PR Title is obvious/clear about what it's changing (see `Title` below).
* The PR Body contains a summary of what/why is included (see `Body` below).
* any linked Issues (see `Body` below)
* Avoid force pushing - it removes commit history and forces reviewers to re-review everything
* Reply to every reviewer comment and only mark it as resolved if you are 100% sure it has been resolved, such as "remove this comment" and you have removed it.

### Title

The title of the PR should clearly state what the PR is doing, and ideally should match the entry that will end up in the changelog.

Examples of good PR titles:

- `azurerm_storage_management_policy - Mark rule.filters.blob_type as required`
- `azurerm_container_registry - support updating replications on demand`
- `azurerm_automation_account - support for the encryption, local_authentication_enabled, and tags properties`
- `Data Source: azurerm_automation_account - prevent panic (#15474) by adding a nil check`
- `Upgrade bot API version from 2021-03-01 to 2021-05-01-preview`
- `New Resource: azurerm_managed_disk_sas_token`
- `New Data Source: azurerm_managed_disk_sas_token`
- `Docs: Fix wrong command in 3.0-upgrade-guide`

Examples of poorly written PR titles:

- `fix sql bug`
- `fixes #1234`
- `new resource`
- `upgrade sdk`
- `upgrade compute api`
- `add cosmos property`
- `support encryption, local_authentication_enabled properties`

### Body

GitHub pre-fills the body from [the pull request template](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/.github/pull_request_template.md). Fill in the description, say which issues it closes using [closing keywords](https://docs.github.com/en/issues/tracking-your-work-with-issues/linking-a-pull-request-to-an-issue#linking-a-pull-request-to-an-issue-using-a-keyword), and tick the boxes that apply. Leave the community note as it is.

> **Note:** If this PR changes meaningfully during the course of review please update the title and description as required.
