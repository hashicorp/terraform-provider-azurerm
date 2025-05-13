---
subcategory: "Dev Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dev_center_network_connection"
description: |-
  Manages a Dev Center Network Connection.
---

# azurerm_dev_center_network_connection

Manages a Dev Center Network Connection.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_dev_center_network_connection" "example" {
  name                = "example-dcnc"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  domain_join_type    = "AzureADJoin"
  subnet_id           = azurerm_subnet.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of this Dev Center Network Connection. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group within which this Dev Center Network Connection should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Dev Center Network Connection should exist. Changing this forces a new resource to be created.

* `domain_join_type` - (Required) The Azure Active Directory Join type. Possible values are `AzureADJoin` and `HybridAzureADJoin`. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the Subnet that is used to attach Virtual Machines.

* `domain_name` - (Optional) The name of the Azure Active Directory domain.

* `domain_password` - (Optional) The password for the account used to join domain.

* `domain_username` - (Optional) The username of the Azure Active Directory account (user or service account) that has permissions to create computer objects in Active Directory.

* `organization_unit` - (Optional) The Azure Active Directory domain Organization Unit (OU).

* `tags` - (Optional) A mapping of tags which should be assigned to the Dev Center Network Connection.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dev Center Network Connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Dev Center Network Connection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Dev Center Network Connection.
* `update` - (Defaults to 30 minutes) Used when updating the Dev Center Network Connection.
* `delete` - (Defaults to 30 minutes) Used when deleting the Dev Center Network Connection.

## Import

An existing Dev Center Network Connection can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_dev_center_network_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DevCenter/networkConnections/networkConnection1
```
