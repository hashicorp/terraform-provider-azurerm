---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_manager"
description: |-
  Get information about an existing Network Manager.
---

# azurerm_network_manager

Use this data source to access information about a Network Manager.

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
}

data "azurerm_network_manager" "example" {
  name                = azurerm_network_manager.example.name
  resource_group_name = azurerm_network_manager.example.resource_group_name
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Network Manager.

* `resource_group_name` - (Required) The Name of the Resource Group where the Network Manager exists.


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Manager.

* `cross_tenant_scopes` - One or more `cross_tenant_scopes` blocks as defined below.

* `description` - A description of the Network Manager.

* `location` - The Azure Region where the Network Manager exists.

* `scope` - A `scope` block as defined below.

* `scope_accesses` - A list of configuration deployment type configured on the Network Manager.

* `tags` - A mapping of tags assigned to the Network Manager.

---

A `scope` block exports the following:

* `management_group_ids` - A list of management group IDs used a scope for the Network Manager.

* `subscription_ids` - A list of subscription IDs used as the scope for the Network Manager.

---

A `cross_tenant_scopes` block exports the following:

* `management_groups` - A list of management groups used as cross tenant scope for the Network Manager.

* `subscriptions` - A list of subscriptions used as cross tenant scope for the Network Manager.

* `tenant_id` - The tenant ID of the cross tenant scope.

 
## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Network Manager.
