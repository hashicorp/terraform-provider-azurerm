# Guide: ARM API Versions

The provider should be implemented using stable Azure Resource Manager (ARM) API/SDK version. Preview versions are prone to sudden breaking changes which can result in a less than ideal user experience (eg: removed property or behavioural change). There are [automated checks on azure-rest-api-specs that prevents breaking changes against stable version](https://github.com/Azure/azure-rest-api-specs/blob/main/documentation/ci-fix.md#sdk-breaking-change-review) but they do not catch everything and are not applicable to preview versions.

These breaking API changes often materialise into [breaking changes](guide-breaking-changes.md) which can involve non-trivial upgrade steps and/or require waiting until a major version release to make the breaking change. v3.0.0 was released in March 2022 and v4.0.0 in August 2024.

In November 2025 we implemented an API version check on PRs that prevents the use of preview versions. All historical usages of preview versions have been allow-listed as exceptions. See `internal/tools/preview-api-version-linter` for the implementation details.

## Rerunning checks locally

If you came to this page through a build failure, once you have removed the preview API dependency, rerun this check locally using the command:

```
go run internal/tools/preview-api-version-linter/main.go
```

## Obtaining exception to use preview API

> [!WARNING]
> Using a preview API version can be risky, prone to human error, and can result in a substandard user experience. An exception is a last resort only when all the consequences are fully understood and there is no alternative.

To add an exception to use preview API version, the following criteria must be met:

1. There is a clear and compelling reason for not using a stable API version. For example: the fix for a critical security vulnerability or impactful bug is only available in the preview version.
1. There is a commitment from the service that no breaking changes will be made to the relevant preview API that could negatively impact azurerm users.
1. There is a commitment from the service team to release a stable API version in the near future, a specific target date has to be set.
1. There is a responsible individual with deep knowledge of the API that can be contacted in the future if required.
1. There is an agreement between Microsoft and Hashicorp that the exception is appropriate.

> [!NOTE]
> A feature being in preview phase is not a sufficient reason to add this exception. The concept of preview should be decoupled between feature and ARM API. It is okay to leave the feature in preview phase while having the API promoted to stable. This will safeguard the API against breaking changes and ensure azurerm support for the feature can be shipped sooner to customers.

To add an exception, insert an entry to `internal/tools/preview-api-version-linter/exceptions.yml` as per below example:

```yml
- module: github.com/hashicorp/go-azure-sdk/resource-manager
  service: compute
  version: 2021-06-01-preview
  stableVersionTargetDate: 2026-01-01
  responsibleIndividual: github.com/gerrytan
```

- `module`: go module name as per go.mod, see internal/tools/preview-api-version-linter/sdk/sdk_types.go for supported modules
- `service`: service name as per the vendor path, for `vendor/github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-01-06-preview` the service name is `compute`
- `version`: preview version as per the vendor path
- `stableVersionTargetDate`: estimated stable API version release date, does not have to be the actual stable version string
- `responsibleIndividual`: individual with deep expertise of the API that can be contacted in the future, has to be a `github.com/myuser` GitHub handle or an email address

Entries have to be sorted alphabetically by `module`, `service` and `version`.

Once added, check the linter is passing by running `go run internal/tools/preview-api-version-linter/main.go`.