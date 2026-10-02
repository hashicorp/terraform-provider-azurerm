---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_ddos_custom_policy"
description: |-
  Manages a Network DDoS Custom Policy.

---

# azurerm_network_ddos_custom_policy

Manages a Network DDoS Custom Policy.

~> **Note:** Azure DDoS Protection custom policy is currently in preview. Please see [the product documentation](https://learn.microsoft.com/azure/ddos-protection/ddos-custom-policy-overview) for the regions in which this feature is available - creating this resource in an unsupported region fails with a `RegionNotEnabledForFeature` error.

-> **Note:** A DDoS Custom Policy is applied to a Standard Load Balancer Frontend IP Configuration from the Load Balancer side, rather than from this resource. A Frontend IP Configuration can only be associated with one DDoS Custom Policy at a time. Support for configuring this association is not yet available in the Provider.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_network_ddos_custom_policy" "example" {
  name                = "example-ddos-custom-policy"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  detection_rule {
    name               = "detectionRuleTcp"
    packets_per_second = 1000000
    traffic_type       = "Tcp"
  }

  detection_rule {
    name               = "detectionRuleUdp"
    packets_per_second = 100000
    traffic_type       = "Udp"
  }

  tags = {
    environment = "production"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Network DDoS Custom Policy. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Network DDoS Custom Policy should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Network DDoS Custom Policy should exist. Changing this forces a new resource to be created.

* `detection_rule` - (Required) One or more `detection_rule` blocks as defined below. At least one and at most three can be specified.

* `tags` - (Optional) A mapping of tags which should be assigned to the Network DDoS Custom Policy.

---

A `detection_rule` block supports the following:

* `name` - (Required) The name of this detection rule.

* `packets_per_second` - (Required) The number of inbound packets per second at which DDoS mitigation is triggered for this `traffic_type`. The supported range depends on `traffic_type`:

    | `traffic_type` | Minimum | Maximum   |
    |----------------|---------|-----------|
    | `Tcp`          | 50000   | 2000000   |
    | `Udp`          | 20000   | 200000    |
    | `TcpSyn`       | 1000    | 100000    |

* `traffic_type` - (Required) The type of traffic this detection rule applies to. Possible values are `Tcp`, `Udp`, and `TcpSyn`. Each `traffic_type` can only be used by one `detection_rule` block.

-> **Note:** Configuring a `detection_rule` for a `traffic_type` disables Azure's adaptive auto-tuning for that traffic type. Traffic types without a `detection_rule` continue to use adaptive auto-tuning.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network DDoS Custom Policy.

* `public_ip_address_ids` - A list of Public IP Address IDs which are associated with the Network DDoS Custom Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network DDoS Custom Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network DDoS Custom Policy.
* `update` - (Defaults to 30 minutes) Used when updating the Network DDoS Custom Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network DDoS Custom Policy.

## Import

A Network DDoS Custom Policy can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_ddos_custom_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/ddosCustomPolicies/ddosCustomPolicy1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Network` - 2025-07-01
