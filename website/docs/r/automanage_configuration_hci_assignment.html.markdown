---
subcategory: "Automanage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automanage_configuration_hci_assignment"
description: |-
  Manages an Automanage Configuration Azure Stack HCI Assignment.
---

# azurerm_automanage_configuration_hci_assignment

Manages an Automanage Configuration Azure Stack HCI Assignment.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "west europe"
}

resource "azurerm_automanage_configuration" "example" {
  name                = "example-configuration"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_stack_hci_cluster" "example" {
  name                = "example-cluster"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  client_id           = data.azurerm_client_config.current.client_id
  tenant_id           = data.azurerm_client_config.current.tenant_id
}

resource "azurerm_automanage_configuration_hci_assignment" "example" {
  // Currently default is the only possible value for name, or it can be omitted
  name                = "default"
  resource_group_name = azurerm_resource_group.example.name
  configuration_id    = azurerm_automanage_configuration.example.id
  cluster_name        = azurerm_stack_hci_cluster.example.name
}
```

## Arguments Reference

The following arguments are supported:

* `cluster_name` - (Required) Azure Stack HCI cluster name. Note that cluster should be in the same subscription as the automanage configuration. Changing this forces a new assignment to be created.

* `configuration_id` - (Required) The ID of the automanage configuration. Changing this forces a new assignment to be created.

* `resource_group_name` - (Required) The name of the resource group where the assignment should exist. Changing this forces a new assignment to be created.

---

* `name` - (Optional) The name which should be used for this assignment. The only possible value is `default`. Defaults to `default`. Changing this forces a new assignment to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Automanage Configuration Azure Stack HCI Assignment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Automanage Configuration Azure Stack HCI Assignment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automanage Configuration Azure Stack HCI Assignment.
* `delete` - (Defaults to 5 minutes) Used when deleting the Automanage Configuration Azure Stack HCI Assignment.

## Import

Automanage Configuration Azure Stack HCI Assignments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automanage_configuration_hci_assignment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AzureStackHci/clusters/clusterName1/providers/Microsoft.Automanage/configurationProfileAssignments/default
```