---
layout: "azurerm"
page_title: "Provider: Azure"
sidebar_current: "docs-azurerm-index"
description: |-
  The Azure Provider is used to interact with the many resources supported by Azure Resource Manager (also known as AzureRM) through its APIs.

---

# Azure Provider

The Azure Provider can be used to configure infrastructure in [Microsoft Azure](https://azure.microsoft.com/en-us/) using the Azure Resource Manager API's. Documentation regarding the [Data Sources](/docs/configuration/data-sources.html) and [Resources](/docs/configuration/resources.html) supported by the Azure Provider can be found in the navigation to the left.

Interested in the provider's latest features, or want to make sure you're up to date? Check out the [changelog](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG.md) for version information and release notes.

## Authenticating to Azure

Terraform supports a number of different methods for authenticating to Azure:

* [Authenticating to Azure using the Azure CLI](auth/azure_cli.html)
* [Authenticating to Azure using Managed Service Identity](auth/managed_service_identity.html)
* [Authenticating to Azure using a Service Principal and a Client Certificate](auth/service_principal_client_certificate.html)
* [Authenticating to Azure using a Service Principal and a Client Secret](auth/service_principal_client_secret.html)

---

We recommend using either a Service Principal or Managed Service Identity when running Terraform non-interactively (such as when running Terraform in a CI server) - and authenticating using the Azure CLI when running Terraform locally.

## Example Usage

```hcl
# Configure the Azure Provider
provider "azurerm" {
  # whilst the `version` attribute is optional, we recommend pinning to a given version of the Provider
  version = "=1.34.0"
}

# Create a resource group
resource "azurerm_resource_group" "test" {
  name     = "production"
  location = "West US"
}

# Create a virtual network within the resource group
resource "azurerm_virtual_network" "test" {
  name                = "production-network"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  address_space       = ["10.0.0.0/16"]
}
```

## Features and Bug Requests

The Azure provider's bugs and feature requests can be found in the [GitHub repo issues](https://github.com/terraform-providers/terraform-provider-azurerm/issues).
Please avoid "me too" or "+1" comments. Instead, use a thumbs up [reaction](https://blog.github.com/2016-03-10-add-reactions-to-pull-requests-issues-and-comments/)
on enhancement requests. Provider maintainers will often prioritize work based on the number of thumbs on an issue.

Community input is appreciated on outstanding issues! We love to hear what use
cases you have for new features, and want to provide the best possible
experience for you using the Azure provider.

If you have a bug or feature request without an existing issue

* if an existing resource or field is working in an unexpected way, [file a bug](https://github.com/terraform-providers/terraform-provider-azurerm/issues/new?template=bug.md).

* if you'd like the provider to support a new resource or field, [file an enhancement/feature request](https://github.com/terraform-providers/terraform-provider-azurerm/issues/new?template=enhancement.md).

The provider maintainers will often use the assignee field on an issue to mark
who is working on it.

* An issue assigned to an individual maintainer indicates that maintainer is working
on the issue

* If you're interested in working on an issue please leave a comment in that issue

---

If you have configuration questions, or general questions about using the provider, try checking out:

* [Terraform's community resources](https://www.terraform.io/docs/extend/community/index.html)
* [HashiCorp support](https://support.hashicorp.com) for Terraform Enterprise customers


## Argument Reference

The following arguments are supported:

* `client_id` - (Optional) The Client ID which should be used. This can also be sourced from the `ARM_CLIENT_ID` Environment Variable.

* `environment` - (Optional) The Cloud Environment which should be used. Possible values are `public`, `usgovernment`, `german`, and `china`. Defaults to `public`. This can also be sourced from the `ARM_ENVIRONMENT` environment variable.

* `subscription_id` - (Optional) The Subscription ID which should be used. This can also be sourced from the `ARM_SUBSCRIPTION_ID` Environment Variable.

* `tenant_id` - (Optional) The Tenant ID which should be used. This can also be sourced from the `ARM_TENANT_ID` Environment Variable.

---

When authenticating as a Service Principal using a Client Certificate, the following fields can be set:

* `client_certificate_password` - (Optional) The password associated with the Client Certificate. This can also be sourced from the `ARM_CLIENT_CERTIFICATE_PASSWORD` Environment Variable.

* `client_certificate_path` - (Optional) The path to the Client Certificate associated with the Service Principal which should be used. This can also be sourced from the `ARM_CLIENT_CERTIFICATE_PATH` Environment Variable.

More information on [how to configure a Service Principal using a Client Certificate can be found in this guide](auth/service_principal_client_certificate.html).

---

When authenticating as a Service Principal using a Client Secret, the following fields can be set:

* `client_secret` - (Optional) The Client Secret which should be used. This can also be sourced from the `ARM_CLIENT_SECRET` Environment Variable.

More information on [how to configure a Service Principal using a Client Secret can be found in this guide](auth/service_principal_client_secret.html).

---

When authenticating using Managed Service Identity, the following fields can be set:

* `msi_endpoint` - (Optional) The path to a custom endpoint for Managed Service Identity - in most circumstances this should be detected automatically. This can also be sourced from the `ARM_MSI_ENDPOINT` Environment Variable.

* `use_msi` - (Optional) Should Managed Service Identity be used for Authentication? This can also be sourced from the `ARM_USE_MSI` Environment Variable. Defaults to `false`.

More information on [how to configure a Service Principal using Managed Service Identity can be found in this guide](auth/managed_service_identity.html).

---

For some advanced scenarios, such as where more granular permissions are necessary - the following properties can be set:

* `partner_id` - (Optional) A GUID/UUID that is [registered](https://docs.microsoft.com/azure/marketplace/azure-partner-customer-usage-attribution#register-guids-and-offers) with Microsoft to facilitate partner resource usage attribution. This can also be sourced from the `ARM_PARTNER_ID` Environment Variable.

* `skip_credentials_validation` - (Optional) Should the AzureRM Provider skip verifying the credentials being used are valid? This can also be sourced from the `ARM_SKIP_CREDENTIALS_VALIDATION` Environment Variable. Defaults to `false`.

* `skip_provider_registration` - (Optional) Should the AzureRM Provider skip registering the Resource Providers it supports? This can also be sourced from the `ARM_SKIP_PROVIDER_REGISTRATION` Environment Variable. Defaults to `false`.

-> By default, Terraform will attempt to register any Resource Providers that it supports, even if they're not used in your configurations to be able to display more helpful error messages. If you're running in an environment with restricted permissions, or wish to manage Resource Provider Registration outside of Terraform you may wish to disable this flag; however please note that the error messages returned from Azure may be confusing as a result (example: `API version 2019-01-01 was not found for Microsoft.Foo`).

It's also possible to use multiple Provider blocks within a single Terraform configuration, for example to work with resources across multiple Subscriptions - more information can be found [in the documentation for Providers](https://www.terraform.io/docs/configuration/providers.html#multiple-provider-instances).
