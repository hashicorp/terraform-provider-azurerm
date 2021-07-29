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
  name = "example"
  resource_group_name = "example"
  location = "West Europe"
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Firewall Policy should exist. Changing this forces a new Firewall Policy to be created.

* `name` - (Required) The name which should be used for this Firewall Policy. Changing this forces a new Firewall Policy to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Firewall Policy should exist. Changing this forces a new Firewall Policy to be created.

---

* `base_policy_id` - (Optional) The ID of the TODO.

* `dns` - (Optional) A `dns` block as defined below.

* `intrusion_detection` - (Optional) A `intrusion_detection` block as defined below.

* `sku` - (Optional) TODO. Changing this forces a new Firewall Policy to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Firewall Policy.

* `threat_intelligence_allowlist` - (Optional) A `threat_intelligence_allowlist` block as defined below.

* `threat_intelligence_mode` - (Optional) TODO.

* `tls_certificate` - (Optional) A `tls_certificate` block as defined below.

---

A `dns` block supports the following:

* `network_rule_fqdn_enabled` - (Optional) Should the TODO be enabled?

* `proxy_enabled` - (Optional) Should the TODO be enabled?

* `servers` - (Optional) Specifies a list of TODO.

---

A `intrusion_detection` block supports the following:

* `mode` - (Optional) TODO.

* `signature_overrides` - (Optional) A `signature_overrides` block as defined below.

* `traffic_bypass` - (Optional) One or more `traffic_bypass` blocks as defined below.

---

A `signature_overrides` block supports the following:

* `id` - (Optional) TODO.

* `state` - (Optional) TODO.

---

A `threat_intelligence_allowlist` block supports the following:

* `fqdns` - (Optional) Specifies a list of TODO.

* `ip_addresses` - (Optional) Specifies a list of TODO.

---

A `tls_certificate` block supports the following:

* `key_vault_secret_id` - (Required) The ID of the TODO.

* `name` - (Required) The name which should be used for this TODO.

---

A `traffic_bypass` block supports the following:

* `name` - (Required) The name which should be used for this TODO.

* `protocol` - (Required) TODO.

* `description` - (Optional) TODO.

* `destination_addresses` - (Optional) Specifies a list of TODO.

* `destination_ip_groups` - (Optional) Specifies a list of TODO.

* `destination_ports` - (Optional) Specifies a list of TODO.

* `source_addresses` - (Optional) Specifies a list of TODO.

* `source_ip_groups` - (Optional) Specifies a list of TODO.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Firewall Policy.

* `child_policies` - A `child_policies` block as defined below.

* `firewalls` - A `firewalls` block as defined below.

* `rule_collection_groups` - A `rule_collection_groups` block as defined below.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Firewall Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Firewall Policy.
* `update` - (Defaults to 30 minutes) Used when updating the Firewall Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Firewall Policy.

## Import

Firewall Policys can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_firewall_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/firewallPolicies/policy1
```
