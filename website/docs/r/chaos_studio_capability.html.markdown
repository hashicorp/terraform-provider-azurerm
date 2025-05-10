---
subcategory: "Chaos Studio"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_chaos_studio_capability"
description: |-
  Manages a Chaos Studio Capability.
---

# azurerm_chaos_studio_capability

Manages a Chaos Studio Capability.

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
resource "azurerm_chaos_studio_capability" "example" {
  capability_type        = "example-value"
  chaos_studio_target_id = azurerm_chaos_studio_target.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `capability_type` - (Required) The capability that should be applied to the Chaos Studio Target. For supported values please see this Chaos Studio [Fault Library](https://learn.microsoft.com/azure/chaos-studio/chaos-studio-fault-library). Changing this forces a new Chaos Studio Capability to be created.

* `chaos_studio_target_id` - (Required) The Chaos Studio Target that the capability should be applied to. Changing this forces a new Chaos Studio Capability to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Chaos Studio Capability.

* `urn` - The Unique Resource Name of the Capability.

---


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Chaos Studio Capability.
* `read` - (Defaults to 5 minutes) Used when retrieving the Chaos Studio Capability.
* `delete` - (Defaults to 30 minutes) Used when deleting the Chaos Studio Capability.

## Import

An existing Chaos Studio Target can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_chaos_studio_capability.example /{scope}/providers/Microsoft.Chaos/targets/{targetName}/capabilities/{capabilityName}
```

* Where `{scope}` is the ID of the Azure Resource under which the Chaos Studio Target exists. For example `/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group`.
* Where `{targetName}` is the name of the Target. For example `targetValue`.
* Where `{capabilityName}` is the name of the Capability. For example `capabilityName`.
