---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_relay_hybrid_connection_authorization_rule"
description: |-
  Manages an Azure Relay Hybrid Connection Authorization Rule.
---

# azurerm_relay_hybrid_connection_authorization_rule

Manages an Azure Relay Hybrid Connection Authorization Rule.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_relay_namespace" "example" {
  name                = "example-relay"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku_name = "Standard"

  tags = {
    source = "terraform"
  }
}

resource "azurerm_relay_hybrid_connection" "example" {
  name                          = "acctestrnhc-%d"
  resource_group_name           = azurerm_resource_group.example.name
  relay_namespace_name          = azurerm_relay_namespace.example.name
  requires_client_authorization = false
  user_metadata                 = "testmetadata"
}

resource "azurerm_relay_hybrid_connection_authorization_rule" "example" {
  name                   = "example"
  resource_group_name    = azurerm_resource_group.example.name
  hybrid_connection_name = azurerm_relay_hybrid_connection.example.name
  namespace_name         = azurerm_relay_namespace.example.name


  listen = true
  send   = true
  manage = false
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Azure Relay Hybrid Connection Authorization Rule. Changing this forces a new Azure Relay Hybrid Connection Authorization Rule to be created.

* `hybrid_connection_name` - (Required) Name of the Azure Relay Hybrid Connection for which this Azure Relay Hybrid Connection Authorization Rule will be created. Changing this forces a new Azure Relay Hybrid Connection Authorization Rule to be created.

* `namespace_name` - (Required) Name of the Azure Relay Namespace for which this Azure Relay Hybrid Connection Authorization Rule will be created. Changing this forces a new Azure Relay Hybrid Connection Authorization Rule to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Azure Relay Hybrid Connection Authorization Rule should exist. Changing this forces a new Azure Relay Hybrid Connection Authorization Rule to be created.

---

* `listen` - (Optional) Grants listen access to this Authorization Rule. Defaults to `false`.

* `send` - (Optional) Grants send access to this Authorization Rule. Defaults to `false`.

* `manage` - (Optional) Grants manage access to this Authorization Rule. When this property is `true` - both `listen` and `send` must be set to `true` too. Defaults to `false`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure Relay Hybrid Connection Authorization Rule.

* `primary_key` - The Primary Key for the Azure Relay Hybrid Connection Authorization Rule.

* `primary_connection_string` - The Primary Connection String for the Azure Relay Hybrid Connection Authorization Rule.

* `secondary_key` - The Secondary Key for the Azure Relay Hybrid Connection Authorization Rule.

* `secondary_connection_string` - The Secondary Connection String for the Azure Relay Hybrid Connection Authorization Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Relay Hybrid Connection Authorization Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Relay Hybrid Connection Authorization Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Azure Relay Hybrid Connection Authorization Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Relay Hybrid Connection Authorization Rule.

## Import

Azure Relay Hybrid Connection Authorization Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_relay_hybrid_connection_authorization_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Relay/namespaces/namespace1/hybridConnections/connection1/authorizationRules/rule1
```
