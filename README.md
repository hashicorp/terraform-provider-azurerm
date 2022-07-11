<a href="https://terraform.io">
    <img src=".github/tf.png" alt="Terraform logo" title="Terraform" align="left" height="50" />
</a>

# Terraform Provider for Azure (Resource Manager)

The AzureRM Terraform Provider allows managing resources within Azure Resource Manager.

When using version 3.0 of the AzureRM Provider we recommend using Terraform 1.x ([the latest version can be found here](https://www.terraform.io/downloads)). Whilst older versions of Terraform Core (0.12.x and later) remain compatible with v3.0 of the AzureRM Provider - support for versions prior to 1.0 will be removed in the next major release of the AzureRM Provider (v4.0).

* [Terraform Website](https://www.terraform.io)
* [AzureRM Provider Documentation](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs)
* [AzureRM Provider Usage Examples](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples)
* [Slack Workspace for Contributors](https://terraform-azure.slack.com) ([Request Invite](https://join.slack.com/t/terraform-azure/shared_invite/enQtNDMzNjQ5NzcxMDc3LWNiY2ZhNThhNDgzNmY0MTM0N2MwZjE4ZGU0MjcxYjUyMzRmN2E5NjZhZmQ0ZTA1OTExMGNjYzA4ZDkwZDYxNDE))

## Usage Example

```hcl
# 1. Specify the version of the AzureRM Provider to use
terraform {
  required_providers {
    azurerm = {
      source = "hashicorp/azurerm"
      version = "=3.0.1"
    }
  }
}

# 2. Configure the AzureRM Provider
provider "azurerm" {
  # The AzureRM Provider supports authenticating using via the Azure CLI, a Managed Identity
  # and a Service Principal. More information on the authentication methods supported by
  # the AzureRM Provider can be found here:
  # https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs#authenticating-to-azure

  # The features block allows changing the behaviour of the Azure Provider, more
  # information can be found here:
  # https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/features-block
  features {}
}

# 3. Create a resource group
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

# 4. Create a virtual network within the resource group
resource "azurerm_virtual_network" "example" {
  name                = "example-network"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  address_space       = ["10.0.0.0/16"]
}
```

* [Usage documentation for the AzureRM Provider can be found in the Terraform Registry](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs).
* [Learn more about Terraform and the AzureRM Provider on HashiCorp Learn](https://learn.hashicorp.com/collections/terraform/azure-get-started).
* [Additional examples can be found in the `./examples` folder within this repository](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples).

## Developing & Contributing to the Provider

The [DEVELOPER.md](DEVELOPER.md) file is a basic outline on how to build and develop the provider while more detailed guides geared towards contributors can be found in the [`/contributing`] (https://github.com/hashicorp/terraform-provider-azurerm/tree/main/contributing) directory of this repository.
