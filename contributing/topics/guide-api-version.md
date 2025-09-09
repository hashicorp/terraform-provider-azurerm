# Guide: ARM API Version

Provider logic should be implemented using stable Azure Resource Manager (ARM) API version. Preview versions are prone to breaking changes which results in very poor azurerm user experience (eg: removed property). There are [automated checks on azure-rest-api-specs that prevents breaking changes against stable version](https://github.com/Azure/azure-rest-api-specs/blob/main/documentation/ci-fix.md#sdk-breaking-change-review). Such checks are not applicable to preview versions.

Breaking API changes often materialise into [azurerm breaking change](guide-breaking-changes.md) which involves non-trivial upgrade steps and long major version release wait time. v3.0.0 was released in March 2022 and v4.0.0 in August 2024.

Day-0 Terraform support for preview API should be done via [azapi provider](https://registry.terraform.io/providers/Azure/azapi/latest/docs) where users have full control of API version choice.

In September 2025 we implemented an API version check on PRs that prevents the use of preview versions. All historical usages of preview versions have been allow-listed as exceptions. See `internal/tools/api-version-lint` for the implementation details.

## Rerunning checks locally

If you came to this page through a build failure, once you have removed the preview API dependency, rerun this check locally using the command:

```
go run internal/tools/api-version-lint/main.go
```

## Obtaining exception to use preview API

> [!WARNING]
> preview API version usage is risky, prone to human error and can result in very poor azurerm user experience. Add an exception as a last resort only when all the consequences are fully understood.

To add an exception to use preview API version, following criteria must be met:

1. There is a clear business reason for not using stable API version, for example: reputational damage to Microsoft / Azure because azurerm Terraform support for the feature cannot be provided in time to meet customers' expectation and azapi support is insufficient
1. Guarantee that no breaking change will be made to the relevant preview API that could negatively impact azurerm users
1. There is a commitment to release a stable API version in the near future, a specific target date has to be set
1. There is a responsible individual with deep knowledge of the API that can be contacted in the future if required
1. There is an agreement between Microsoft and Hashicorp that the exception is appropriate

> [!NOTE]
> Feature being in preview phase is not a sufficient reason to add this exception. The concept of preview should be decoupled between feature and ARM API. It is okay to leave the feature in preview phase while having the API promoted to stable. This will safeguard the API against breaking changes and ensure azurerm support for the feature can be shipped sooner to customers. Otherwise Terraform support for the feature should be provided via azapi.

To add an exception, insert an entry to `internal/toos/api-version-lint/exceptions.yml` as per below example:

```yml
- module: github.com/hashicorp/go-azure-sdk/resource-manager
  service: compute
  version: 2021-06-01-preview
  stableVersionTargetDate: 2026-01-01
  responsibleIndividual: github.com/gerrytan
```

- `module`: go module name as per go.mod, see internal/tools/api-version-lint/sdk/sdk_types.go for supported modules
- `service`: service name as per the vendor path, for `vendor/github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-01-06-preview` the service name is `compute`
- `version`: preview version as per the vendor path
- `stableVersionTargetDate`: estimated stable API version release date, does not have to be the actual stable version string
- `responsibleIndividual`: individual with deep expertise of the API that can be contacted in the future, has to be a `github.com/myuser` GitHub handle or an email address

Entries have to be sorted alphabetically by `module`, `service` and `version`.

Once added, check the linter is passing by running `go run internal/tools/api-version-lint/main.go`.