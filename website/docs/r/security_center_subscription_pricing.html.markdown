---
subcategory: "Security Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_security_center_subscription_pricing"
description: |-
    Manages the Pricing Tier for Azure Security Center in the current subscription.
---

# azurerm_security_center_subscription_pricing

Manages the Pricing Tier for Azure Security Center in the current subscription.

~> **Note:** Deletion of this resource will reset the pricing tier to `Free`

## Example Usage

### Basic usage

```hcl
resource "azurerm_security_center_subscription_pricing" "example" {
  tier          = "Standard"
  resource_type = "VirtualMachines"
}
```

### Using Extensions with Defender CSPM

```hcl
resource "azurerm_security_center_subscription_pricing" "example1" {
  tier          = "Standard"
  resource_type = "CloudPosture"

  extension {
    name = "ContainerRegistriesVulnerabilityAssessments"
  }

  extension {
    name = "AgentlessVmScanning"
    additional_extension_properties = {
      ExclusionTags = "[]"
    }
  }

  extension {
    name = "AgentlessDiscoveryForKubernetes"
  }

  extension {
    name = "SensitiveDataDiscovery"
  }
}
```

## Argument Reference

The following arguments are supported:

* `tier` - (Required) The pricing tier to use. Possible values are `Free` and `Standard`.

* `resource_type` - (Optional) The resource type this setting affects. Possible values are `Api`, `AppServices`, `ContainerRegistry`, `KeyVaults`, `KubernetesService`, `SqlServers`, `SqlServerVirtualMachines`, `StorageAccounts`, `VirtualMachines`, `Arm`, `Dns`, `OpenSourceRelationalDatabases`, `Containers`, `CosmosDbs` and `CloudPosture`. Defaults to `VirtualMachines`

* `subplan` - (Optional) Resource type pricing subplan. Contact your MSFT representative for possible values. Changing this forces a new resource to be created.

* `extension` - (Optional) One or more `extension` blocks as defined below.

---

A `extension` block supports the following:

* `name` - (Required) The name of extension.

* `additional_extension_properties` - (Optional) Key/Value pairs that are required for some extensions.

~> **Note:** If an extension is not defined, it will not be enabled.

~> **Note:** Changing the pricing tier to `Standard` affects all resources of the given type in the subscription and could be quite costly.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The subscription pricing ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Security Center Subscription Pricing.
* `read` - (Defaults to 5 minutes) Used when retrieving the Security Center Subscription Pricing.
* `update` - (Defaults to 1 hour) Used when updating the Security Center Subscription Pricing.
* `delete` - (Defaults to 1 hour) Used when deleting the Security Center Subscription Pricing.

## Import

The pricing tier can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_security_center_subscription_pricing.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Security/pricings/<resource_type>
```
