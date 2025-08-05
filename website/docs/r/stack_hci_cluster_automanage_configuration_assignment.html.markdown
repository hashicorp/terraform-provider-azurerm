---
subcategory: "Automanage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stack_hci_cluster_automanage_configuration_assignment"
description: |-
  Manages a Azure Stack HCI Cluster Automanage Configuration Profile Assignment.
---

# azurerm_stack_hci_cluster_automanage_configuration_assignment

Manages an Azure Stack HCI Cluster Automanage Configuration Profile Assignment.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_automanage_configuration" "example" {
  name                = "example-configuration"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_stack_hci_cluster" "example" {
  name                = "example-stack"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  client_id           = data.azurerm_client_config.current.client_id
  tenant_id           = data.azurerm_client_config.current.tenant_id
}

resource "azurerm_stack_hci_cluster_automanage_configuration_assignment" "example" {
  stack_hci_cluster_id = azurerm_stack_hci_cluster.example.id
  configuration_id     = azurerm_automanage_configuration.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `configuration_id` - (Required) The ID of the Automanage Configuration to assign to the Azure Stack HCI Cluster. Changing this forces a new Azure Stack HCI Cluster Automanage Configuration Profile Assignment to be created.

* `stack_hci_cluster_id` - (Required) The ID of the Azure Stack HCI Cluster. Changing this forces a new Azure Stack HCI Cluster Automanage Configuration Profile Assignment to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Azure Stack HCI Cluster Automanage Configuration Profile Assignment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Stack HCI Cluster Automanage Configuration Profile Assignment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Stack HCI Cluster Automanage Configuration Profile Assignment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Stack HCI Cluster Automanage Configuration Profile Assignment.

## Import

Azure Stack HCI Cluster Automanage Configuration Profile Assignments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stack_hci_cluster_automanage_configuration_assignment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AzureStackHCI/clusters/cluster1/providers/Microsoft.AutoManage/configurationProfileAssignments/default
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.AzureStackHCI` - 2022-05-04
