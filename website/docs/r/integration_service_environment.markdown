---
subcategory: "Logic App"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_integration_service_environment"
description: |-
  Manages private and isolated Logic App instances within an Azure virtual network.
---

# azurerm_integration_service_environment

Manages private and isolated Logic App instances within an Azure virtual network.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "exampleRG1"
  location = "westeurope"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/22"]
}

resource "azurerm_subnet" "isesubnet1" {
  name                 = "isesubnet1"
  virtual_network_id = azurerm_virtual_network.example.id
  address_prefixes     = ["10.0.1.0/26"]

  delegation {
    name = "integrationServiceEnvironments"
    service_delegation {
      name = "Microsoft.Logic/integrationServiceEnvironments"
    }
  }
}

resource "azurerm_subnet" "isesubnet2" {
  name                 = "isesubnet2"
  virtual_network_id = azurerm_virtual_network.example.id
  address_prefixes     = ["10.0.1.64/26"]
}

resource "azurerm_subnet" "isesubnet3" {
  name                 = "isesubnet3"
  virtual_network_id = azurerm_virtual_network.example.id
  address_prefixes     = ["10.0.1.128/26"]
}

resource "azurerm_subnet" "isesubnet4" {
  name                 = "isesubnet4"
  virtual_network_id = azurerm_virtual_network.example.id
  address_prefixes     = ["10.0.1.192/26"]
}

resource "azurerm_integration_service_environment" "example" {

  name                 = "example-ise"
  location             = azurerm_resource_group.example.location
  resource_group_name  = azurerm_resource_group.example.name
  sku_name             = "Developer_0"
  access_endpoint_type = "Internal"
  virtual_network_subnet_ids = [
    azurerm_subnet.isesubnet1.id,
    azurerm_subnet.isesubnet2.id,
    azurerm_subnet.isesubnet3.id,
    azurerm_subnet.isesubnet4.id
  ]
  tags = {
    environment = "development"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Integration Service Environment. Changing this forces a new Integration Service Environment to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Integration Service Environment should exist. Changing this forces a new Integration Service Environment to be created.

* `location` - (Required) The Azure Region where the Integration Service Environment should exist. Changing this forces a new Integration Service Environment to be created.

* `sku_name` - (Required) The sku name and capacity of the Integration Service Environment. Possible Values for `sku` element are `Developer` and `Premium` and possible values for the `capacity` element are from `0` to `10`.  Defaults to `sku` of `Developer` with a `Capacity` of `0` (e.g. `Developer_0`). Changing this forces a new Integration Service Environment to be created when `sku` element is not the same with existing one.

~> **NOTE** For a `sku_name` using the `Developer` `sku` the `capacity` element must be always `0`. For a `sku_name` using the `sku` of `Premium` the `capacity` element can be between `0` and `10`.

* `access_endpoint_type` - (Required) The type of access endpoint to use for the Integration Service Environment. Possible Values are `Internal` and `External`. Changing this forces a new Integration Service Environment to be created.

* `virtual_network_subnet_ids` - (Required) A list of virtual network subnet ids to be used by Integration Service Environment. Exactly four distinct ids to subnets must be provided. Changing this forces a new Integration Service Environment to be created.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Integration Service Environment.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Integration Service Environment.

* `connector_endpoint_ip_addresses` - The list of access endpoint ip addresses of connector.

* `connector_outbound_ip_addresses` - The list of outgoing ip addresses of connector.

* `workflow_endpoint_ip_addresses` - The list of access endpoint ip addresses of workflow.

* `workflow_outbound_ip_addresses` - The list of outgoing ip addresses of workflow.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 hours) Used when creating the Integration Service Environment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Integration Service Environment.
* `update` - (Defaults to 5 hours) Used when updating the Integration Service Environment.
* `delete` - (Defaults to 5 hours) Used when deleting the Integration Service Environment.

## Import

Integration Service Environments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_integration_service_environment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Logic/integrationServiceEnvironments/ise1
```
