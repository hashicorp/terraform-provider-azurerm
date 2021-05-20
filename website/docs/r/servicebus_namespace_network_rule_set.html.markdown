---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_namespace_network_rule_set"
description: |-
  Manages a ServiceBus Namespace Network Rule Set Set.
---

# azurerm_servicebus_namespace_network_rule_set

Manages a ServiceBus Namespace Network Rule Set Set.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_servicebus_namespace" "example" {
  name                = "example-sb-namespace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Premium"

  capacity = 1
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["172.17.0.0/16"]
  dns_servers         = ["10.0.0.4", "10.0.0.5"]
}

resource "azurerm_subnet" "example" {
  name                 = "default"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["172.17.0.0/24"]

  service_endpoints = ["Microsoft.ServiceBus"]
}

resource "azurerm_servicebus_namespace_network_rule_set" "example" {
  namespace_name      = azurerm_servicebus_namespace.example.name
  resource_group_name = azurerm_resource_group.example.name

  default_action = "Deny"

  network_rules {
    subnet_id                            = azurerm_subnet.example.id
    ignore_missing_vnet_service_endpoint = false
  }

  ip_rules = ["1.1.1.1"]
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the ServiceBus Namespace Network Rule Set should exist. Changing this forces a new resource to be created.

* `namespace_name` - (Required) Specifies the ServiceBus Namespace name to which to attach the ServiceBus Namespace Network Rule Set. Changing this forces a new resource to be created.

~> **NOTE:** The ServiceBus Namespace must be `Premium` in order to attach a ServiceBus Namespace Network Rule Set.

* `default_action` - (Optional) Specifies the default action for the ServiceBus Namespace Network Rule Set. Possible values are `Allow` and `Deny`. Defaults to `Deny`.

* `ip_rules` - (Optional) One or more IP Addresses, or CIDR Blocks which should be able to access the ServiceBus Namespace.

* `network_rules` - (Optional) One or more `network_rules` blocks as defined below.

---

A `network_rules` block supports the following:

* `subnet_id` - (Required) The Subnet ID which should be able to access this ServiceBus Namespace.

* `ignore_missing_vnet_service_endpoint` - (Optional) Should the ServiceBus Namespace Network Rule Set ignore missing Virtual Network Service Endpoint option in the Subnet? Defaults to `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the ServiceBus Namespace Network Rule Set.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the ServiceBus Namespace Network Rule Set.
* `update` - (Defaults to 30 minutes) Used when updating the ServiceBus Namespace Network Rule Set.
* `read` - (Defaults to 5 minutes) Used when retrieving the ServiceBus Namespace Network Rule Set.
* `delete` - (Defaults to 30 minutes) Used when deleting the ServiceBus Namespace Network Rule Set.

## Import

Service Bus Namespace can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_servicebus_namespace_network_rule_set.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Servicebus/namespaces/sbns1/networkrulesets/default
```
