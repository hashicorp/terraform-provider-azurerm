---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_manager"
description: |-
  Manages a Network Managers.

---

# azurerm_network_manager

Manages a Network Managers.

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

* `name` - (Required) Specifies the name which should be used for this Network Managers. Changing this forces a new Network Managers to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Network Managers should exist. Changing this forces a new Network Managers to be created.

* `location` - (Required) Specifies the Azure Region where the Network Managers should exist. Changing this forces a new resource to be created.

* `scope` - (Required) A `scope` block as defined below.

* `scope_accesses` - (Required) A list of configuration deployment type. Possible values are `Connectivity` and `SecurityAdmin`, corresponds to if Connectivity Configuration and Security Admin Configuration is allowed for the Network Manager.

* `description` - (Optional) A description of the network manager.

* `tags` - (Optional) A mapping of tags which should be assigned to the Network Managers.

---

A `scope` block supports the following:

* `management_group_ids` - (Optional) A list of management group IDs.

* `subscription_ids` - (Optional) A list of subscription IDs.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Managers.

* `cross_tenant_scopes` - A `cross_tenant_scopes` block as defined below.

---

A `cross_tenant_scopes` block exports the following:

* `management_groups` - List of management groups.

* `subscriptions` - List of subscriptions.

* `tenant_id` - Tenant ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Managers.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Managers.
* `update` - (Defaults to 30 minutes) Used when updating the Network Managers.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Managers.

## Import

Network Managers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_manager.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/networkManagers/networkManager1
```
