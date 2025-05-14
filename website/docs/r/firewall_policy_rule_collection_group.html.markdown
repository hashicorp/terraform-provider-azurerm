---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_firewall_policy_rule_collection_group"
description: |-
  Manages a Firewall Policy Rule Collection Group.
---

# azurerm_firewall_policy_rule_collection_group

Manages a Firewall Policy Rule Collection Group.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_firewall_policy" "example" {
  name                = "example-fwpolicy"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_firewall_policy_rule_collection_group" "example" {
  name               = "example-fwpolicy-rcg"
  firewall_policy_id = azurerm_firewall_policy.example.id
  priority           = 500
  application_rule_collection {
    name     = "app_rule_collection1"
    priority = 500
    action   = "Deny"
    rule {
      name = "app_rule_collection1_rule1"
      protocols {
        type = "Http"
        port = 80
      }
      protocols {
        type = "Https"
        port = 443
      }
      source_addresses  = ["10.0.0.1"]
      destination_fqdns = ["*.microsoft.com"]
    }
  }

  network_rule_collection {
    name     = "network_rule_collection1"
    priority = 400
    action   = "Deny"
    rule {
      name                  = "network_rule_collection1_rule1"
      protocols             = ["TCP", "UDP"]
      source_addresses      = ["10.0.0.1"]
      destination_addresses = ["192.168.1.1", "192.168.1.2"]
      destination_ports     = ["80", "1000-2000"]
    }
  }

  nat_rule_collection {
    name     = "nat_rule_collection1"
    priority = 300
    action   = "Dnat"
    rule {
      name                = "nat_rule_collection1_rule1"
      protocols           = ["TCP", "UDP"]
      source_addresses    = ["10.0.0.1", "10.0.0.2"]
      destination_address = "192.168.1.1"
      destination_ports   = ["80"]
      translated_address  = "192.168.0.1"
      translated_port     = "8080"
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Firewall Policy Rule Collection Group. Changing this forces a new Firewall Policy Rule Collection Group to be created.

* `firewall_policy_id` - (Required) The ID of the Firewall Policy where the Firewall Policy Rule Collection Group should exist. Changing this forces a new Firewall Policy Rule Collection Group to be created.

* `priority` - (Required) The priority of the Firewall Policy Rule Collection Group. The range is 100-65000.

---

* `application_rule_collection` - (Optional) One or more `application_rule_collection` blocks as defined below.

* `nat_rule_collection` - (Optional) One or more `nat_rule_collection` blocks as defined below.

* `network_rule_collection` - (Optional) One or more `network_rule_collection` blocks as defined below.

---

A `application_rule_collection` block supports the following:

* `name` - (Required) The name which should be used for this application rule collection.

* `action` - (Required) The action to take for the application rules in this collection. Possible values are `Allow` and `Deny`.

* `priority` - (Required) The priority of the application rule collection. The range is `100` - `65000`.

* `rule` - (Required) One or more `application_rule` blocks as defined below.

---

A `network_rule_collection` block supports the following:

* `name` - (Required) The name which should be used for this network rule collection.

* `action` - (Required) The action to take for the network rules in this collection. Possible values are `Allow` and `Deny`.

* `priority` - (Required) The priority of the network rule collection. The range is `100` - `65000`.

* `rule` - (Required) One or more `network_rule` blocks as defined below.

---

A `nat_rule_collection` block supports the following:

* `name` - (Required) The name which should be used for this NAT rule collection.

* `action` - (Required) The action to take for the NAT rules in this collection. Currently, the only possible value is `Dnat`.

* `priority` - (Required) The priority of the NAT rule collection. The range is `100` - `65000`.

* `rule` - (Required) A `nat_rule` block as defined below.

---

A `application_rule` (application rule) block supports the following:

* `name` - (Required) The name which should be used for this rule.

* `description` - (Optional) The description which should be used for this rule.

* `protocols` - (Optional) One or more `protocols` blocks as defined below.

* `http_headers` - (Optional) Specifies a list of HTTP/HTTPS headers to insert. One or more `http_headers` blocks as defined below.

* `source_addresses` - (Optional) Specifies a list of source IP addresses (including CIDR, IP range and `*`).

* `source_ip_groups` - (Optional) Specifies a list of source IP groups.

* `destination_addresses` - (Optional) Specifies a list of destination IP addresses (including CIDR, IP range and `*`).

* `destination_urls` - (Optional) Specifies a list of destination URLs for which policy should hold. Needs Premium SKU for Firewall Policy. Conflicts with `destination_fqdns`.

* `destination_fqdns` - (Optional) Specifies a list of destination FQDNs. Conflicts with `destination_urls`.

* `destination_fqdn_tags` - (Optional) Specifies a list of destination FQDN tags.

* `terminate_tls` - (Optional) Boolean specifying if TLS shall be terminated (true) or not (false). Must be `true` when using `destination_urls`. Needs Premium SKU for Firewall Policy.

* `web_categories` - (Optional) Specifies a list of web categories to which access is denied or allowed depending on the value of `action` above. Needs Premium SKU for Firewall Policy.

---

A `network_rule` (network rule) block supports the following:

* `name` - (Required) The name which should be used for this rule.

* `description` - (Optional) The description which should be used for this rule.

* `protocols` - (Required) Specifies a list of network protocols this rule applies to. Possible values are `Any`, `TCP`, `UDP`, `ICMP`.

* `destination_ports` - (Required) Specifies a list of destination ports.

* `source_addresses` - (Optional) Specifies a list of source IP addresses (including CIDR, IP range and `*`).

* `source_ip_groups` - (Optional) Specifies a list of source IP groups.

* `destination_addresses` - (Optional) Specifies a list of destination IP addresses (including CIDR, IP range and `*`) or Service Tags.

* `destination_ip_groups` - (Optional) Specifies a list of destination IP groups.

* `destination_fqdns` - (Optional) Specifies a list of destination FQDNs.

---

A `nat_rule` (NAT rule) block supports the following:

* `name` - (Required) The name which should be used for this rule.

* `description` - (Optional) The description which should be used for this rule.

* `protocols` - (Required) Specifies a list of network protocols this rule applies to. Possible values are `TCP`, `UDP`.

* `source_addresses` - (Optional) Specifies a list of source IP addresses (including CIDR, IP range and `*`).

* `source_ip_groups` - (Optional) Specifies a list of source IP groups.

* `destination_address` - (Optional) The destination IP address (including CIDR).

* `destination_ports` - (Optional) Specifies a list of destination ports. Only one destination port is supported in a NAT rule.

* `translated_address` - (Optional) Specifies the translated address.

* `translated_fqdn` - (Optional) Specifies the translated FQDN.

~> **Note:** Exactly one of `translated_address` and `translated_fqdn` should be set.

* `translated_port` - (Required) Specifies the translated port.

---

A `protocols` block supports the following:

* `type` - (Required) Protocol type. Possible values are `Http` and `Https`.

* `port` - (Required) Port number of the protocol. Range is 0-64000.

---

A `http_headers` block supports the following:

* `name` - (Required) Specifies the name of the header.

* `value` - (Required) Specifies the value of the value.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Firewall Policy Rule Collection Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Firewall Policy Rule Collection Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Firewall Policy Rule Collection Group.
* `update` - (Defaults to 30 minutes) Used when updating the Firewall Policy Rule Collection Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Firewall Policy Rule Collection Group.

## Import

Firewall Policy Rule Collection Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_firewall_policy_rule_collection_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/firewallPolicies/policy1/ruleCollectionGroups/gruop1
```
