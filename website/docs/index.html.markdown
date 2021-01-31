---
layout: "azurerm"
page_title: "Provider: Azure"
description: |-
  The Azure Provider is used to interact with the many resources supported by Azure Resource Manager (also known as AzureRM) through its APIs.

---

# Azure Provider

The Azure Provider can be used to configure infrastructure in [Microsoft Azure](https://azure.microsoft.com/en-us/) using the Azure Resource Manager API's. Documentation regarding the [Data Sources](/docs/configuration/data-sources.html) and [Resources](/docs/configuration/resources.html) supported by the Azure Provider can be found in the navigation to the left.

Interested in the provider's latest features, or want to make sure you're up to date? Check out the [changelog](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG.md) for version information and release notes.

## Authenticating to Azure

Terraform supports a number of different methods for authenticating to Azure:

* [Authenticating to Azure using the Azure CLI](guides/azure_cli.html)
* [Authenticating to Azure using Managed Service Identity](guides/managed_service_identity.html)
* [Authenticating to Azure using a Service Principal and a Client Certificate](guides/service_principal_client_certificate.html)
* [Authenticating to Azure using a Service Principal and a Client Secret](guides/service_principal_client_secret.html)

---

We recommend using either a Service Principal or Managed Service Identity when running Terraform non-interactively (such as when running Terraform in a CI server) - and authenticating using the Azure CLI when running Terraform locally.

## Example Usage

```hcl
# We strongly recommend using the required_providers block to set the
# Azure Provider source and version being used
terraform {
  required_providers {
    azurerm = {
      source = "hashicorp/azurerm"
      version = "=2.40.1"
    }
  }
}

# Configure the Microsoft Azure Provider
provider "azurerm" {
  features {}
}

# Create a resource group
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

# Create a virtual network within the resource group
resource "azurerm_virtual_network" "example" {
  name                = "example-network"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
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

* An issue assigned to an individual maintainer indicates that the maintainer is working
on the issue

* If you're interested in working on an issue please leave a comment on that issue

---

If you have configuration questions, or general questions about using the provider, try checking out:

* [Terraform's community resources](https://www.terraform.io/docs/extend/community/index.html)
* [HashiCorp support](https://support.hashicorp.com) for Terraform Enterprise customers

## Argument Reference

The following arguments are supported:

* `features` - (Required) A `features` block as defined below which can be used to customize the behaviour of certain Azure Provider resources.

* `client_id` - (Optional) The Client ID which should be used. This can also be sourced from the `ARM_CLIENT_ID` Environment Variable.

* `environment` - (Optional) The Cloud Environment which should be used. Possible values are `public`, `usgovernment`, `german`, and `china`. Defaults to `public`. This can also be sourced from the `ARM_ENVIRONMENT` environment variable.

* `subscription_id` - (Optional) The Subscription ID which should be used. This can also be sourced from the `ARM_SUBSCRIPTION_ID` Environment Variable.

* `tenant_id` - (Optional) The Tenant ID should be used. This can also be sourced from the `ARM_TENANT_ID` Environment Variable.

---

When authenticating as a Service Principal using a Client Certificate, the following fields can be set:

* `client_certificate_password` - (Optional) The password associated with the Client Certificate. This can also be sourced from the `ARM_CLIENT_CERTIFICATE_PASSWORD` Environment Variable.

* `client_certificate_path` - (Optional) The path to the Client Certificate associated with the Service Principal which should be used. This can also be sourced from the `ARM_CLIENT_CERTIFICATE_PATH` Environment Variable.

More information on [how to configure a Service Principal using a Client Certificate can be found in this guide](guides/service_principal_client_certificate.html).

---

When authenticating as a Service Principal using a Client Secret, the following fields can be set:

* `client_secret` - (Optional) The Client Secret which should be used. This can also be sourced from the `ARM_CLIENT_SECRET` Environment Variable.

More information on [how to configure a Service Principal using a Client Secret can be found in this guide](guides/service_principal_client_secret.html).

---

When authenticating using Managed Service Identity, the following fields can be set:

* `msi_endpoint` - (Optional) The path to a custom endpoint for Managed Service Identity - in most circumstances, this should be detected automatically. This can also, be sourced from the `ARM_MSI_ENDPOINT` Environment Variable.

* `use_msi` - (Optional) Should Managed Service Identity be used for Authentication? This can also be sourced from the `ARM_USE_MSI` Environment Variable. Defaults to `false`.

More information on [how to configure a Service Principal using Managed Service Identity can be found in this guide](guides/managed_service_identity.html).

---

For some advanced scenarios, such as where more granular permissions are necessary - the following properties can be set:

* `disable_terraform_partner_id` - (Optional) Disable sending the Terraform Partner ID if a custom `partner_id` isn't specified, which allows Microsoft to better understand the usage of Terraform. The Partner ID does not give HashiCorp any direct access to usage information. This can also be sourced from the `ARM_DISABLE_TERRAFORM_PARTNER_ID` environment variable. Defaults to `false`.

* `metadata_host` - (Optional) The Hostname of the Azure Metadata Service (for example `management.azure.com`), used to obtain the Cloud Environment when using a Custom Azure Environment. This can also be sourced from the `ARM_METADATA_HOST` Environment Variable.

~> **Note:** `environment` must be set to the requested environment name in the list of available environments held in the `metadata_host`.

* `partner_id` - (Optional) A GUID/UUID that is [registered](https://docs.microsoft.com/azure/marketplace/azure-partner-customer-usage-attribution#register-guids-and-offers) with Microsoft to facilitate partner resource usage attribution. This can also be sourced from the `ARM_PARTNER_ID` Environment Variable.

* `skip_credentials_validation` - (Optional) Should the AzureRM Provider skip verifying the credentials being used are valid? This can also be sourced from the `ARM_SKIP_CREDENTIALS_VALIDATION` Environment Variable. Defaults to `false`.

~> **Note:** if `skip_credentials_validation` is false, AzureRM Provider will send a request to list all available providers.

* `skip_provider_registration` - (Optional) Should the AzureRM Provider skip registering the Resource Providers it supports? This can also be sourced from the `ARM_SKIP_PROVIDER_REGISTRATION` Environment Variable. Defaults to `false`.

-> By default, Terraform will attempt to register any Resource Providers that it supports, even if they're not used in your configurations to be able to display more helpful error messages. If you're running in an environment with restricted permissions, or wish to manage Resource Provider Registration outside of Terraform you may wish to disable this flag; however, please note that the error messages returned from Azure may be confusing as a result (example: `API version 2019-01-01 was not found for Microsoft.Foo`).

* `storage_use_azuread` - (Optional) Should the AzureRM Provider use AzureAD to connect to the Storage Blob & Queue API's, rather than the SharedKey from the Storage Account? This can also be sourced from the `ARM_STORAGE_USE_AZUREAD` Environment Variable. Defaults to `false`.

~> **Note:** This requires that the User/Service Principal being used has the associated `Storage` roles - which are added to new Contributor/Owner role-assignments, but **have not** been backported by Azure to existing role-assignments.

~> **Note:** The Files & Table Storage API's do not support authenticating via AzureAD and will continue to use a SharedKey to access the API's.

It's also possible to use multiple Provider blocks within a single Terraform configuration, for example, to work with resources across multiple Subscriptions - more information can be found [in the documentation for Providers](https://www.terraform.io/docs/configuration/providers.html#multiple-provider-instances).

## Features

It's possible to configure the behaviour of certain resources using the `features` block - more details can be found below.

The `features` block supports the following:

* `key_vault` - (Optional) A `key_vault` block as defined below.

* `template_deployment` - (Optional) A `template_deployment` block as defined below.

* `virtual_machine` - (Optional) A `virtual_machine` block as defined below.

* `virtual_machine_scale_set` - (Optional) A `virtual_machine_scale_set` block as defined below.

---

The `key_vault` block supports the following:

* `recover_soft_deleted_key_vaults` - (Optional) Should the `azurerm_key_vault`, `azurerm_key_vault_certificate`, `azurerm_key_vault_key` and `azurerm_key_vault_secret` resources recover a Soft-Deleted Key Vault/Item? Defaults to `true`.

~> **Note:** When recovering soft-deleted Key Vault items (Keys, Certificates, and Secrets) the Principal used by Terraform needs the `"recover"` permission.

* `purge_soft_delete_on_destroy` - (Optional) Should the `azurerm_key_vault`, `azurerm_key_vault_certificate`, `azurerm_key_vault_key` and `azurerm_key_vault_secret` resources be permanently deleted (e.g. purged) when destroyed? Defaults to `true`.

~> **Note:** When purge protection is enabled, a key vault or an object in the deleted state cannot be purged until the retention period (7-90 days) has passed.

---

The `template_deployment` block supports the following:

* `delete_nested_items_during_deletion` - (Optional) Should the `azurerm_resource_group_template_deployment` resource attempt to delete resources that have been provisioned by the ARM Template, when the Resource Group Template Deployment is deleted? Defaults to `true`.

---

The `virtual_machine` block supports the following:

* `delete_os_disk_on_deletion` - (Optional) Should the `azurerm_linux_virtual_machine` and `azurerm_windows_virtual_machine` resources delete the OS Disk attached to the Virtual Machine when the Virtual Machine is destroyed? Defaults to `true`.

~> **Note:** This does not affect the older `azurerm_virtual_machine` resource, which has its own flags for managing this within the resource.

* `graceful_shutdown` - (Optional) Should the `azurerm_linux_virtual_machine` and `azurerm_windows_virtual_machine` request a graceful shutdown when the Virtual Machine is destroyed? Defaults to `false`.

~> **Note:** When using a graceful shutdown, Azure gives the Virtual Machine a 5 minutes window in which to complete the shutdown process, at which point the machine will be force powered off - [more information can be found in this blog post](https://azure.microsoft.com/en-us/blog/linux-and-graceful-shutdowns-2/).

---

The `virtual_machine_scale_set` block supports the following:

* `roll_instances_when_required` - (Optional) Should the `azurerm_linux_virtual_machine_scale_set` and `azurerm_windows_virtual_machine_scale_set` resources automatically roll the instances in the Scale Set when Required (for example when updating the Sku/Image). Defaults to `true`.
