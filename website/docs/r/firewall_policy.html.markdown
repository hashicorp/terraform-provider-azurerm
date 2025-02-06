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
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_firewall_policy" "example" {
  name                = "example-policy"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Firewall Policy should exist. Changing this forces a new Firewall Policy to be created.

* `name` - (Required) The name which should be used for this Firewall Policy. Changing this forces a new Firewall Policy to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Firewall Policy should exist. Changing this forces a new Firewall Policy to be created.

---

* `base_policy_id` - (Optional) The ID of the base Firewall Policy.

* `dns` - (Optional) A `dns` block as defined below.

* `identity` - (Optional) An `identity` block as defined below.

* `insights` - (Optional) An `insights` block as defined below.

* `intrusion_detection` - (Optional) A `intrusion_detection` block as defined below.

* `private_ip_ranges` - (Optional) A list of private IP ranges to which traffic will not be SNAT.

* `auto_learn_private_ranges_enabled` - (Optional) Whether enable auto learn private ip range.

* `sku` - (Optional) The SKU Tier of the Firewall Policy. Possible values are `Standard`, `Premium` and `Basic`. Defaults to `Standard`. Changing this forces a new Firewall Policy to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Firewall Policy.

* `threat_intelligence_allowlist` - (Optional) A `threat_intelligence_allowlist` block as defined below.

* `threat_intelligence_mode` - (Optional) The operation mode for Threat Intelligence. Possible values are `Alert`, `Deny` and `Off`. Defaults to `Alert`.

* `tls_certificate` - (Optional) A `tls_certificate` block as defined below.

* `sql_redirect_allowed` - (Optional) Whether SQL Redirect traffic filtering is allowed. Enabling this flag requires no rule using ports between `11000`-`11999`.

* `explicit_proxy` - (Optional) A `explicit_proxy` block as defined below.

---

A `dns` block supports the following:

* `proxy_enabled` - (Optional) Whether to enable DNS proxy on Firewalls attached to this Firewall Policy? Defaults to `false`.

* `servers` - (Optional) A list of custom DNS servers' IP addresses.

---

A `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Firewall Policy. Only possible value is `UserAssigned`.

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Firewall Policy.

---

An `insights` block supports the following:

* `enabled` - (Required) Whether the insights functionality is enabled for this Firewall Policy.

* `default_log_analytics_workspace_id` - (Required) The ID of the default Log Analytics Workspace that the Firewalls associated with this Firewall Policy will send their logs to, when there is no location matches in the `log_analytics_workspace`.

* `retention_in_days` - (Optional) The log retention period in days.

* `log_analytics_workspace` - (Optional) A list of `log_analytics_workspace` block as defined below.

---

A `intrusion_detection` block supports the following:

* `mode` - (Optional) In which mode you want to run intrusion detection: `Off`, `Alert` or `Deny`.

* `signature_overrides` - (Optional) One or more `signature_overrides` blocks as defined below.

* `traffic_bypass` - (Optional) One or more `traffic_bypass` blocks as defined below.

* `private_ranges` - (Optional) A list of Private IP address ranges to identify traffic direction. By default, only ranges defined by IANA RFC 1918 are considered private IP addresses.

---

A `log_analytics_workspace` block supports the following:

* `id` - (Required) The ID of the Log Analytics Workspace that the Firewalls associated with this Firewall Policy will send their logs to when their locations match the `firewall_location`.

* `firewall_location` - (Required) The location of the Firewalls, that when matches this Log Analytics Workspace will be used to consume their logs.

---

A `signature_overrides` block supports the following:

* `id` - (Optional) 12-digit number (id) which identifies your signature.

* `state` - (Optional) state can be any of `Off`, `Alert` or `Deny`.

---

A `threat_intelligence_allowlist` block supports the following:

* `fqdns` - (Optional) A list of FQDNs that will be skipped for threat detection.

* `ip_addresses` - (Optional) A list of IP addresses or CIDR ranges that will be skipped for threat detection.

---

A `tls_certificate` block supports the following:

* `key_vault_secret_id` - (Required) The ID of the Key Vault, where the secret or certificate is stored.

* `name` - (Required) The name of the certificate.

---

A `traffic_bypass` block supports the following:

* `name` - (Required) The name which should be used for this bypass traffic setting.

* `protocol` - (Required) The protocols any of `ANY`, `TCP`, `ICMP`, `UDP` that shall be bypassed by intrusion detection.

* `description` - (Optional) The description for this bypass traffic setting.

* `destination_addresses` - (Optional) Specifies a list of destination IP addresses that shall be bypassed by intrusion detection.

* `destination_ip_groups` - (Optional) Specifies a list of destination IP groups that shall be bypassed by intrusion detection.

* `destination_ports` - (Optional) Specifies a list of destination IP ports that shall be bypassed by intrusion detection.

* `source_addresses` - (Optional) Specifies a list of source addresses that shall be bypassed by intrusion detection.

* `source_ip_groups` - (Optional) Specifies a list of source IP groups that shall be bypassed by intrusion detection.

---

A `explicit_proxy` block supports the following:

* `enabled` - (Optional) Whether the explicit proxy is enabled for this Firewall Policy.

* `http_port` - (Optional) The port number for explicit http protocol.

* `https_port` - (Optional) The port number for explicit proxy https protocol.

* `enable_pac_file` - (Optional) Whether the pac file port and url need to be provided.

* `pac_file_port` - (Optional) Specifies a port number for firewall to serve PAC file.

* `pac_file` - (Optional) Specifies a SAS URL for PAC file.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Firewall Policy.

* `child_policies` - A list of reference to child Firewall Policies of this Firewall Policy.

* `firewalls` - A list of references to Azure Firewalls that this Firewall Policy is associated with.

* `rule_collection_groups` - A list of references to Firewall Policy Rule Collection Groups that belongs to this Firewall Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Firewall Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Firewall Policy.
* `update` - (Defaults to 30 minutes) Used when updating the Firewall Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Firewall Policy.

## Import

Firewall Policies can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_firewall_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/firewallPolicies/policy1
```
