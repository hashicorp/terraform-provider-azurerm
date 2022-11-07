---
subcategory: "Automanage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automanage_configuration_profile_hci_assignment"
description: |-
  Manages a automanage ConfigurationProfileHCIAssignment.
---

# azurerm_automanage_configuration_profile_hci_assignment

Manages a automanage ConfigurationProfileHCIAssignment.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "example-automanage"
  location = "West Europe"
}

data "azurerm_client_config" "current" {}

resource "azurerm_stack_hci_cluster" "test" {
  name                = "example-azshci"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  client_id           = data.azurerm_client_config.current.client_id
  tenant_id           = data.azurerm_client_config.current.tenant_id
}

resource "azurerm_automanage_configuration_profile_hci_assignment" "test" {
  name = "default"
  resource_group_name = azurerm_resource_group.test.name
  cluster_name = azurerm_stack_hci_cluster.test.name
  configuration_profile = "/providers/Microsoft.Automanage/bestPractices/AzureBestPracticesProduction"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this automanage ConfigurationProfileHCIAssignment. Changing this forces a new automanage ConfigurationProfileHCIAssignment to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the automanage ConfigurationProfileHCIAssignment should exist. Changing this forces a new automanage ConfigurationProfileHCIAssignment to be created.

* `cluster_name` - (Required) The name of the Arc machine. Changing this forces a new automanage ConfigurationProfileHCIAssignment to be created.

* `configuration_profile` - (Required) The Automanage configurationProfile ARM Resource URI.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the automanage ConfigurationProfileHCIAssignment.

* `managed_by` - Azure resource id. Indicates if this resource is managed by another Azure resource.

* `target_id` - The ID of the target.

* `type` - The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts".

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the automanage ConfigurationProfileHCIAssignment.
* `read` - (Defaults to 5 minutes) Used when retrieving the automanage ConfigurationProfileHCIAssignment.
* `delete` - (Defaults to 30 minutes) Used when deleting the automanage ConfigurationProfileHCIAssignment.

## Import

automanage ConfigurationProfileHCIAssignments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automanage_configuration_profile_hci_assignment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AzureStackHci/clusters/cluster1/providers/Microsoft.Automanage/configurationProfileAssignments/configurationProfileAssignment1
```