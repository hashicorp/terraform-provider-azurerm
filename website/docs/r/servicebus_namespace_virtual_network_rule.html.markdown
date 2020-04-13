---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_namespace_virtual_network_rule"
description: |-
  Manages a ServiceBus Namespace Virtual Network Rule within a ServiceBus.
---

# azurerm_servicebus_namespace_virtual_network_rule

Manages a ServiceBus Namespace Virtual Network Rule within a ServiceBus.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "terraform-servicebus"
  location = "West US"
}

resource "azurerm_virtual_network" "vnet" {
  name                = "example-vnet"
  address_space       = ["10.7.29.0/29"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "subnet" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.vnet.name
  address_prefix       = "10.7.29.0/29"
  service_endpoints    = ["Microsoft.ServiceBus"]
}

resource "azurerm_servicebus_namespace" "servicebus" {
  name                = "tfex_sevicebus_namespace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Premium"

  tags = {
    source = "terraform"
  }
}

resource "azurerm_servicebus_namespace_virtual_network_rule" "example" {
  name                = "servicebus-namespace-vnet-rule"
  namespace_name      = azurerm_servicebus_namespace.servicebus.name
  resource_group_name = azurerm_resource_group.example.name
  subnet_id           = azurerm_subnet.subnet.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the ServiceBus Namespace Virtual Network Rule resource. Changing this forces a new resource to be created.

~> **NOTE:** `name` must be between 2-64 characters long and must satisfy all of the requirements below:
1. Contains only alphanumeric, underscores and hyphen characters
2. Cannot start with an underscore or hyphen
3. Cannot end with a hyphen

* `namespace_name` - (Required) Specifies the name of the ServiceBus Namespace. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the ServiceBus Namespace exists. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the subnet that the ServiceBus Namespace will be connected to.


## Attributes Reference

The following attributes are exported:

* `id` - The ServiceBus Namespace ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the ServiceBus Namespace Virtual Network Rule.
* `update` - (Defaults to 30 minutes) Used when updating the ServiceBus Namespace Virtual Network Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the ServiceBus Namespace Virtual Network Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the ServiceBus Namespace Virtual Network Rule.

## Import

ServiceBus Namespace virtual network rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_servicebus_namespace_virtual_network_rule.rule1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ServiceBus/namespaces/namespace1/virtualnetworkrules/rule1
```
