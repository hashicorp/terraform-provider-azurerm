---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_firewall_policy_application_rule_collection"
description: |-
  Manages a Firewall Policy Application Rule Collection.
---

# azurerm_firewall_policy_application_rule_collection

Manages a Firewall Policy Application Rule Collection.

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
}

resource "azurerm_firewall_policy_application_rule_collection" "example" {
  name                     = "example-fwpolicy-arc"
  rule_collection_group_id = azurerm_firewall_policy_rule_collection_group.example.id
  priority                 = 500
  action                   = "Deny"
  rule {
    name = "example-fwpolicy-arc-rule1"
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
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Firewall Policy Application Rule Collection. Changing this forces a new Firewall Policy Application Rule Collection to be created.

* `rule_collection_group_id` - (Required) The ID of the Firewall Policy Rule Collection Group where the Firewall Policy Application Rule Collection should exist. Changing this forces a new Firewall Policy Application Rule Collection to be created.

* `action` - (Required) The action to take for the application rules in this collection. Possible values are `Allow` and `Deny`.

* `priority` - (Required) The priority of the application rule collection. The range is `100` - `65000`.

* `rule` - (Required) One or more `rule` (application rule) blocks as defined below.

---

A `rule` (application rule) block supports the following:

* `name` - (Required) The name which should be used for this rule.

* `description` - (Optional) The description which should be used for this rule.

* `protocols` - (Optional) One or more `protocols` blocks as defined below. Not required when specifying `destination_fqdn_tags`, but required when specifying `destination_fqdns`.

* `source_addresses` - (Optional) Specifies a list of source IP addresses (including CIDR, IP range and `*`).

* `source_ip_groups` - (Optional) Specifies a list of source IP groups.

* `destination_addresses` - (Optional) Specifies a list of destination IP addresses (including CIDR, IP range and `*`).

* `destination_urls` - (Optional) Specifies a list of destination URLs for which policy should hold. Needs Premium SKU for Firewall Policy. Conflicts with `destination_fqdns`.

* `destination_fqdns` - (Optional) Specifies a list of destination FQDNs. Conflicts with `destination_urls`.

* `destination_fqdn_tags` - (Optional) Specifies a list of destination FQDN tags.

* `terminate_tls` - (Optional) Boolean specifying if TLS shall be terminated (true) or not (false). Must be `true` when using `destination_urls`. Needs Premium SKU for Firewall Policy.

* `web_categories` - (Optional) Specifies a list of web categories to which access is denied or allowed depending on the value of `action` above. Needs Premium SKU for Firewall Policy.

---

A `protocols` block supports the following:

* `type` - (Required) Protocol type. Possible values are `Http` and `Https`.

* `port` - (Required) Port number of the protocol. Range is 0-64000.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Firewall Policy Application Rule Collection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Firewall Policy Application Rule Collection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Firewall Policy Application Rule Collection.
* `update` - (Defaults to 30 minutes) Used when updating the Firewall Policy Application Rule Collection.
* `delete` - (Defaults to 30 minutes) Used when deleting the Firewall Policy Application Rule Collection.

## Import

Firewall Policy Application Rule Collections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_firewall_policy_application_rule_collection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/firewallPolicies/policy1/ruleCollectionGroups/gruop1/ruleCollections/collection1
```
