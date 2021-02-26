---
subcategory: "Security Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_security_center_subscription_pricing"
description: |-
    Manages the Pricing Tier for Azure Security Center in the current subscription.
---

# azurerm_security_center_subscription_pricing

Manages the Pricing Tier for Azure Security Center in the current subscription.

~> **NOTE:** This resource requires the `Owner` permission on the Subscription.

~> **NOTE:** Deletion of this resource does not change or reset the pricing tier to `Free`

## Example Usage

```hcl
resource "azurerm_security_center_subscription_pricing" "example" {
  tier          = "Standard"
  resource_type = "VirtualMachines"
}
```

## Argument Reference

The following arguments are supported:

* `tier` - (Required) The pricing tier to use. Possible values are `Free` and `Standard`.
* `resource_type` - (Required) The resource type this setting affects. Possible values are `AppServices`, `ContainerRegistry`, `KeyVaults`, `KubernetesService`, `SqlServers`, `SqlServerVirtualMachines`, `StorageAccounts`, `VirtualMachines`, `Arm` and `Dns`. 

~> **NOTE:** Changing the pricing tier to `Standard` affects all resources of the given type in the subscription and could be quite costly.

## Attributes Reference

The following attributes are exported:

* `id` - The subscription pricing ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Security Center Subscription Pricing.
* `update` - (Defaults to 60 minutes) Used when updating the Security Center Subscription Pricing.
* `read` - (Defaults to 5 minutes) Used when retrieving the Security Center Subscription Pricing.
* `delete` - (Defaults to 60 minutes) Used when deleting the Security Center Subscription Pricing.

## Import

The pricing tier can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_security_center_subscription_pricing.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Security/pricings/<resource_type>
```
