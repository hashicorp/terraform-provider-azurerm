---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_firewall_policy"
description: |-
  Manages a Firewall Policy.
---

# azurerm_firewall_policy

Manages a Firewall Policy.

## Example Usage

```hcl
resource "azurerm_firewall_policy" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "West Europe"
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Firewall Policy should exist. Changing this forces a new Firewall Policy to be created.

* `name` - (Required) The name which should be used for this Firewall Policy. Changing this forces a new Firewall Policy to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Firewall Policy should exist. Changing this forces a new Firewall Policy to be created.

---

* `base_policy_id` - (Optional) The ID of the base Firewall Policy.

* `threat_intelligence_mode` - (Optional) The operation mode for Threat Intelligence. Possible values are `Alert`, `Deny` and `Off`. Defaults to `Alert`.

* `threat_intelligence_allowlist` - (Optional) A `threat_intelligence_allowlist` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Firewall Policy.

---

A `threat_intelligence_allowlist` block supports the following:

* `ip_addresses` - (Optional) A list of IP addresses or IP address ranges that will be skipped for threat detection.

* `fqdns` - (Optional) A list of FQDNs that will be skipped for threat detection.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Firewall Policy.

* `child_policies` - A list of reference to child Firewall Policies of this Firewall Policy.

* `firewalls` - A list of references to Azure Firewalls that this Firewall Policy is associated with.

* `rule_collection_groups` - A list of references to Firewall Policy Rule Collection Groups that belongs to this Firewall Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Firewall Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Firewall Policy.
* `update` - (Defaults to 30 minutes) Used when updating the Firewall Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Firewall Policy.

## Import

networks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_firewall_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/firewallPolicies/policy1
```
