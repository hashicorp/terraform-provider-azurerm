---
subcategory: "Chaos Studio"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_chaos_studio_target"
description: |-
  Manages a Chaos Studio Target.
---

<!-- Note: This documentation is generated. Any manual changes will be overwritten -->

# azurerm_chaos_studio_target

Manages a Chaos Studio Target.

## Example Usage

```hcl
resource "azurerm_kubernetes_cluster" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  dns_prefix          = "acctestaksexample"
  default_node_pool {
    name       = "example-value"
    node_count = "example-value"
    vm_size    = "example-value"
    upgrade_settings {
      max_surge = "example-value"
    }
  }
  identity {
    type = "example-value"
  }
}
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}
resource "azurerm_chaos_studio_target" "example" {
  location           = azurerm_resource_group.example.location
  target_resource_id = azurerm_kubernetes_cluster.example.id
  target_type        = "example-value"
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Chaos Studio Target should exist. Changing this forces a new Chaos Studio Target to be created.

* `target_resource_id` - (Required) Specifies the Target Resource Id within which this Chaos Studio Target should exist. Changing this forces a new Chaos Studio Target to be created.

* `target_type` - (Required) The name of the Chaos Studio Target. This has the format of [publisher]-[targetType] e.g. `Microsoft-StorageAccount`. For supported values please see this Target Type column in [this table](https://learn.microsoft.com/azure/chaos-studio/chaos-studio-fault-providers). Changing this forces a new Chaos Studio Target to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Chaos Studio Target.

---



## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating this Chaos Studio Target.
* `delete` - (Defaults to 30 minutes) Used when deleting this Chaos Studio Target.
* `read` - (Defaults to 5 minutes) Used when retrieving this Chaos Studio Target.
* `update` - (Defaults to 30 minutes) Used when updating this Chaos Studio Target.

## Import

An existing Chaos Studio Target can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_chaos_studio_target.example /{scope}/providers/Microsoft.Chaos/targets/{targetName}
```

* Where `{scope}` is the ID of the Azure Resource under which the Chaos Studio Target exists. For example `/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group`.
* Where `{targetName}` is the name of the Target. For example `targetValue`.
