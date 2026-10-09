# High Level Overview

The AzureRM Provider is a Plugin which is invoked by Terraform (Core) and comprised of Data Sources and Resources.

Within the AzureRM Provider, these Data Sources and Resources are grouped into Service Packages - which are logical groupings of Data Sources/Resources based on the Azure Service they're related to.

Each of these Data Sources and Resources has both Acceptance Tests and Documentation associated with each Data Source/Resource - the Acceptance Tests are also located within this Service Package, however the Documentation exists within a dedicated folder.

## Project Structure

The Azure Provider is a large codebase which has evolved over time - but tends to follow consistent patterns for the most-part.

The Provider is split up into Service Packages (see [terminology](reference-glossary.md)) - with some other logic sprinkled across several packages.

At a high-level, the Provider structure is:

| Directory/Package | Description |
| ------------------- | ------------- |
| `./examples` | More complete example usages of Data Sources and Resources offered by this Provider. |
| `./helpers` | **This package is deprecated (and so intentionally not documented) - new functionality should instead be added to either the Service Package or [go-azure-helpers](https://github.com/hashicorp/go-azure-helpers)**. |
| `./internal/acceptance` | The Acceptance Test wrappers that we use in the Azure Provider, offering common patterns across the Provider to be reused. |
| `./internal/clients` | Refers to the Client from each Service Package, which is used in Data Sources and Resources to access the Azure APIs. |
| `./internal/common` | Helper functions for registering Clients (for example, setting the user agent, configuring credentials etc.). |
| `./internal/features` | Feature Toggles for Provider functionality and behaviour (for example, enabling Betas or changing a resource type's soft delete or purge protection). This also contains the struct and parsing of/default values for the `features` block (within the Provider block). |
| `./internal/locks` | Common locking across resources where necessary to workaround API consistency issues. |
| `./internal/provider` | The Provider implementation itself, the Provider schema and a reference to each Service Registration so that Data Sources and Resources can be surfaced within the Provider. |
| `./internal/resourceproviders` | The list of Resource Providers which should be auto-registered by the Provider. |
| `./internal/sdk` | The Typed Plugin SDK functionality used in this Provider. |
| `./internal/services` | Packages for each service that the provider supports (e.g. `appconfiguration`, `compute`) which contain the Data Sources and Resources supported by the service. |
| `./internal/tags` | Helpers for parsing Tags from the Terraform Configuration and setting Tags into the Terraform State. |
| `./internal/tf` | Helpers and abstractions on top of the Terraform Plugin SDK. |
| `./internal/timeouts` | Helpers for computing the Timeouts for a Data Source / Resource - used in Untyped Data Sources and Untyped Resources. |
| `./internal/tools` | Tooling used to generate functionality within the Provider, for example for Resource IDs and Website Documentation. |
| `./scripts` | Scripts used during testing, linting, and building the provider. |
| `./vendor` | Vendored copies of the go modules the provider uses. For more information please refer to the official [Go Documentation](https://go.dev/ref/mod#vendoring). |
| `./website` | Guides and documentation for each resource (in `./website/docs/r`) and data source (in `./website/docs/d`) that are published to the Terraform [registry](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs). |

> **Note:** Due to the size of the codebase and open Pull Requests - when functionality is moved we use aliasing to try and avoid breaking open Pull Requests / big-bang migrations. These aliases stick around for a few weeks to allow open PRs to be merged without extra out-of-scope changes - at which point these aliases are removed.

Each Service Package consists of (to take `appconfiguration` as an example):

| File/Directory | Description |
| ---------------- | ------------- |
| `./services/appconfiguration` | |
| `./client` | A Client struct, with a reference to any SDK Clients used to access the Azure APIs within this Service Package. |
| `./parse` | Resource ID Formatters and Parsers. |
| `./validate` | Validation functions for this Service Package, including Resource ID Validators. |
| `./app_configuration_data_source.go` | The Data Source `azurerm_app_configuration`. |
| `./app_configuration_data_source_test.go` | Acceptance tests for the Data Source `azurerm_app_configuration`. |
| `./app_configuration_key_resource.go` | The Resource `azurerm_app_configuration_key`. |
| `./app_configuration_key_resource_test.go` | Acceptance Tests for the Resource `azurerm_app_configuration_key`. |
| `./app_configuration_resource.go` | The Resource `azurerm_app_configuration`. |
| `./app_configuration_resource_test.go` | Acceptance tests for the Resource `azurerm_app_configuration`. |
| `./registration.go` | The Service Registration for this Service Package. |

Some Service Packages may also contain:

| File/Directory         | Description                                                                                  |
|------------------------|----------------------------------------------------------------------------------------------|
| `./migration`          | Any State Migrations used in Resources.                                                      |
| `./resourceids.go`     | Used to generate Resource ID Formatters, Parsers and Validators.                             |

---

* Data Sources use the filename format: `{name}_data_source.go`
* Acceptance Tests for Data Sources use the filename format: `{name}_data_source_test.go` (note: Golang requires that Tests are contained within a `test.go` file)
* Resources use the filename format: `{name}_resource.go`
* Acceptance Tests for Resources use the filename format: `{name}_resource_test.go` (note: Golang requires that Tests are contained within a `test.go` file)

> **Note:** there are a handful of exceptions to these to reduce stuttering (e.g. Resource Provider Registration Resource)

## Types of Data Sources/Resources within the Provider

The provider is built on [the Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk), with a Typed SDK of our own layered on top in [`internal/sdk`](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/internal/sdk). That gives four kinds of Data Source and Resource:

1. Untyped Data Sources ([example](https://github.com/hashicorp/terraform-provider-azurerm/blob/2ff15cca48adc7315f67d8b653409e621963ca64/internal/services/search/search_service_data_source.go#L16-L131)).
2. Untyped Resources ([example](https://github.com/hashicorp/terraform-provider-azurerm/blob/2ff15cca48adc7315f67d8b653409e621963ca64/internal/services/search/search_service_resource.go#L24-L289)).
3. Typed Data Sources ([example](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/internal/services/privatednsresolver/private_dns_resolver_data_source.go)).
4. Typed Resources ([example](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/internal/services/privatednsresolver/private_dns_resolver_resource.go)).

Most of the codebase is still untyped. New work is typed, and existing resources move across over time, which is also what will let the provider move to the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework) later. See [Typed vs Untyped](best-practices.md#typed-vs-untyped-resources) for the differences.

## Interaction with Azure

This Provider makes use of a number of SDKs to interact with both the Azure Resource Manager and a number of associated Data Plane APIs, these are:

* [go-azure-sdk](https://github.com/hashicorp/go-azure-sdk) - HashiCorp's generated SDK for Azure Resource Manager and a number of Data Plane APIs. New work uses this.
* [The Azure SDK for Go](https://github.com/Azure/azure-sdk-for-go) - the older Microsoft SDK, built on [Azure/go-autorest](https://github.com/Azure/go-autorest). A handful of services still use it and are being moved to go-azure-sdk.
* [Giovanni](https://github.com/jackofallops/giovanni) - for interaction with the Azure Storage Data Plane APIs.

## Testing the Provider

Since the behaviour of the Azure API can change over time, the Provider leans on Acceptance Tests over Unit Tests for asserting that the Data Sources and Resources within the Provider work as expected.

More details and guidance on how to test Data Sources/Resources can be found in [the Acceptance Testing reference](reference-acceptance-testing.md).
