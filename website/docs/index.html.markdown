---
layout: "azurerm"
page_title: "Provider: Azure"
sidebar_current: "docs-azurerm-index"
description: |-
  The Azure Provider is used to interact with the many resources supported by Azure Resource Manager (also known as AzureRM) through its APIs.

---

# Azure Provider

The Azure Provider is used to interact with the many resources supported by Azure Resource Manager (AzureRM) through its APIs.

~> **Note:** This supercedes the [legacy Azure provider](/docs/providers/azure/index.html), which interacts with Azure using the Service Management API.

Use the navigation to the left to read about the available resources.

# Creating Credentials

Terraform supports authenticating to Azure through a Service Principal or the Azure CLI.

We recommend [using a Service Principal when running in a shared environment](authenticating_via_service_principal.html) (such as within a CI server/automation) - and [authenticating via the Azure CLI](authenticating_via_azure_cli.html) when you're running Terraform locally.

~> **NOTE:** Authenticating via the Azure CLI is only supported when using a User Account. If you're using a Service Principal (e.g. via `az login --service-principal`) you should instead [authenticate via the Service Principal directly](authenticating_via_service_principal.html).

## Example Usage

```hcl
# Configure the Azure Provider
provider "azurerm" { }

# Create a resource group
resource "azurerm_resource_group" "network" {
  name     = "production"
  location = "West US"
}

# Create a virtual network within the resource group
resource "azurerm_virtual_network" "network" {
  name                = "production-network"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.network.location}"
  resource_group_name = "${azurerm_resource_group.network.name}"

  subnet {
    name           = "subnet1"
    address_prefix = "10.0.1.0/24"
  }

  subnet {
    name           = "subnet2"
    address_prefix = "10.0.2.0/24"
  }

  subnet {
    name           = "subnet3"
    address_prefix = "10.0.3.0/24"
  }
}
```

## Argument Reference

The following arguments are supported:

* `subscription_id` - (Optional) The subscription ID to use. It can also
  be sourced from the `ARM_SUBSCRIPTION_ID` environment variable.

* `client_id` - (Optional) The client ID to use. It can also be sourced from
  the `ARM_CLIENT_ID` environment variable.

* `client_secret` - (Optional) The client secret to use. It can also be sourced from
  the `ARM_CLIENT_SECRET` environment variable.

* `tenant_id` - (Optional) The tenant ID to use. It can also be sourced from the
  `ARM_TENANT_ID` environment variable.

* `use_msi` - (Optional) Set to true to authenticate using managed service identity.
  It can also be sourced from the `ARM_USE_MSI` environment variable.

* `msi_endpoint` - (Optional) The REST endpoint to retrieve an MSI token from. Terraform
  will attempt to discover this automatically but it can be specified manually here.
  It can also be sourced from the `ARM_MSI_ENDPOINT` environment variable.

* `environment` - (Optional) The cloud environment to use. It can also be sourced
  from the `ARM_ENVIRONMENT` environment variable. Supported values are:
  * `public` (default)
  * `usgovernment`
  * `german`
  * `china`

* `skip_credentials_validation` - (Optional) Prevents the provider from validating
  the given credentials. When set to `true`, `skip_provider_registration` is assumed.
  It can also be sourced from the `ARM_SKIP_CREDENTIALS_VALIDATION` environment
  variable; defaults to `false`.

* `skip_provider_registration` - (Optional) Prevents the provider from registering
  the ARM provider namespaces, this can be used if you don't wish to give the Active
  Directory Application permission to register resource providers. It can also be
  sourced from the `ARM_SKIP_PROVIDER_REGISTRATION` environment variable; defaults
  to `false`.

## Testing

The following Environment Variables must be set to run the acceptance tests:

~> **NOTE:** The Acceptance Tests require the use of a Service Principal - authenticating via either the Azure CLI or MSI is not supported.

* `ARM_SUBSCRIPTION_ID` - The ID of the Azure Subscription in which to run the Acceptance Tests.
* `ARM_CLIENT_ID` - The Client ID of the Service Principal.
* `ARM_CLIENT_SECRET` - The Client Secret associated with the Service Principal.
* `ARM_TENANT_ID` - The Tenant ID to use.
* `ARM_ENVIRONMENT` - The Azure Cloud Environment to use, such as `public`, `german` etc. Defaults to `public`.
* `ARM_TEST_LOCATION` - The primary Azure Region to provision resources in for the Acceptance Tests.
* `ARM_TEST_LOCATION_ALT` - The secondary Azure Region to provision resources in for the Acceptance Tests. This needs to be a different region to `ARM_TEST_LOCATION`.
