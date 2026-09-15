<!-- markdownlint-disable-file first-line-heading -->
<!--  All Submissions -->

## Community Note
<!-- Please leave the community note as is. -->
* Please vote on this PR by adding a :thumbsup: [reaction](https://blog.github.com/2016-03-10-add-reactions-to-pull-requests-issues-and-comments/) to the original PR to help the community and maintainers prioritize for review
* Please do not leave comments along the lines of "+1", "me too" or "any updates", they generate extra noise for PR followers and do not help prioritize for review

## Description

<!-- Please include a description below with the reason for the PR, what it is doing, what it is trying to accomplish, and anything relevant for a reviewer to know. 

If this is a breaking change for users please detail how it cannot be avoided and why it should be made in a minor version of the provider -->

## PR Checklist

- [ ] I have followed the guidelines in our [Contributing Documentation](../blob/main/contributing/README.md).
- [ ] I have checked to ensure there aren't other open [Pull Requests](../pulls) for the same update/change.
- [ ] I have checked if my changes close any open issues. If so please include appropriate [closing keywords](https://docs.github.com/en/issues/tracking-your-work-with-issues/linking-a-pull-request-to-an-issue#linking-a-pull-request-to-an-issue-using-a-keyword) below.
- [ ] I have updated/added Documentation as required, written in a helpful and kind way to assist users that may be unfamiliar with the resource / data source.
- [ ] I have used a meaningful PR title to help maintainers and other users understand this change and help prevent duplicate work.
For example: “`resource_name_here` - description of change e.g. adding property `new_property_name_here`”

## Changes to existing Resource / Data Source

- [ ] I have added an explanation of what my changes do and why I'd like you to include them (This may be covered by linking to an issue above, but may benefit from additional explanation).
- [ ] I have written new tests for my resource or datasource changes & updated any relevant documentation.
- [ ] I have successfully run tests with my changes locally. If not, please provide details on testing challenges that prevented you running the tests.
- [ ] (For changes that include a **state migration only**). I have manually tested the migration path between relevant versions of the provider.

## Testing

<!-- You can erase this section if it is not applicable to your Pull Request. For example, a documentation fix requires no testing -->

- [ ] My submission includes Test coverage as described in the [Contribution Guide](../blob/main/contributing/topics/guide-new-resource.md) and the tests pass. (if this is not possible for any reason, please include details of why you did or could not add test coverage)

<!-- Please include testing logs or evidence here or an explanation on why no testing evidence can be provided. 

For state migrations, ensure a state migration test has been added. For further details on testing state migration changes please see our guide on [state migrations](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/contributing/topics/guide-state-migrations.md#testing) in the contributor documentation. -->

## Change Log

Changelog entries must be added to a file in the pull request's source branch. The file must be located under `./changelog/{PR Number}.hcl`.

If this pull request does not require a changelog entry, check the box below. Reference [maintainer-merging](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/contributing/topics/maintainer-merging.md) to determine whether this needs an entry.

- [ ] This change does not require a changelog entry, I have checked the guide linked above to determine this

Adding a changelog entry can be done manually or using the `make changelog` target. Each changelog entry must specify a type, for a list of types, run `make changelog-types`.

---

<!-- What type of PR is this? -->
This is a (please select all that apply):

- [ ] Bug Fix
- [ ] New Feature (ie adding a service, resource, or data source)
- [ ] Enhancement
- [ ] Breaking Change

## Related Issue(s)

<!-- You can erase this section if it is not applicable to your Pull Request. For example, if there is not an active issue to close. -->

Fixes #0000

## AI Assistance Disclosure

<!-- 
IMPORTANT!

If you are using any kind of AI/LLM assistance to contribute to the AzureRM provider, this must be disclosed in the pull request.

If this is the case, please check the box below, and include the extent to which AI was used. (e.g. documentation only, code generation, etc.)

If responses to this pull request are/will be generated using AI, disclose this as well.
-->

- [ ] AI Assisted - This contribution was made by, or with the assistance of, AI/LLMs

<!-- extent of AI usage must be described here if the box above is checked -->

<!-- heimdall_github_prtemplate:grc-pci_dss-2024-01-05 -->

## Rollback Plan

If a change needs to be reverted, we will publish an updated version of the provider.

## Changes to Security Controls

Are there any changes to security controls (access controls, encryption, logging) in this pull request? If so, explain.

> [!NOTE]
> If this PR changes meaningfully during the course of review please update the title and description as required.
