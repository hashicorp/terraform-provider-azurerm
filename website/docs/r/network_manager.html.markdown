---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_manager"
description: |-
  Manages a Network Manager.

---

# azurerm_network_manager

Manages a Network Manager.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_subscription" "current" {
}

resource "azurerm_network_manager" "example" {
  name                = "example-network-manager"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  scope {
    subscription_ids = [data.azurerm_subscription.current.id]
  }
  scope_accesses = ["Connectivity", "SecurityAdmin"]
  description    = "example network manager"
  tags = {
    foo = "bar"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Network Manager. Changing this forces a new Network Manager to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Network Manager should exist. Changing this forces a new Network Manager to be created.

* `location` - (Required) Specifies the Azure Region where the Network Manager should exist. Changing this forces a new resource to be created.

* `scope` - (Required) A `scope` block as defined below.

* `description` - (Optional) A description of the Network Manager.

* `scope_accesses` - (Optional) A list of configuration deployment types. Possible values are `Connectivity`, `SecurityAdmin` and `Routing`, which specify whether Connectivity Configuration, Security Admin Configuration or Routing Configuration are allowed for the Network Manager.

* `tags` - (Optional) A mapping of tags which should be assigned to the Network Manager.

---

A `scope` block supports the following:

* `management_group_ids` - (Optional) A list of management group IDs.

~> **Note:** When specifying a scope at the management group level, you need to register the `Microsoft.Network` at the management group scope before deploying a Network Manager, more information can be found in the [Azure document](https://learn.microsoft.com/en-us/azure/virtual-network-manager/concept-network-manager-scope#scope).

* `subscription_ids` - (Optional) A list of subscription IDs.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Manager.

* `cross_tenant_scopes` - One or more `cross_tenant_scopes` blocks as defined below.

---

A `cross_tenant_scopes` block exports the following:

* `management_groups` - List of management groups.

* `subscriptions` - List of subscriptions.

* `tenant_id` - Tenant ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Manager.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Manager.
* `update` - (Defaults to 30 minutes) Used when updating the Network Manager.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Manager.

## Import

Network Manager can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_manager.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/networkManagers/networkManager1
```
