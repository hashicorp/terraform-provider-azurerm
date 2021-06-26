---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_security_group"
description: |-
  Gets information about an existing Network Security Group.
---

# Data Source: azurerm_network_security_group

Use this data source to access information about an existing Network Security Group.

## Example Usage

```hcl
data "azurerm_network_security_group" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
}

output "location" {
  value = data.azurerm_network_security_group.example.location
}
```

## Argument Reference

* `name` - Specifies the Name of the Network Security Group.
* `resource_group_name` - Specifies the Name of the Resource Group within which the Network Security Group exists


## Attributes Reference

* `id` - The ID of the Network Security Group.

* `location` - The supported Azure location where the resource exists.

* `security_rule` - One or more `security_rule` blocks as defined below.

* `tags` - A mapping of tags assigned to the resource.


The `security_rule` block supports:

* `name` - The name of the security rule.

* `description` - The description for this rule.

* `protocol` - The network protocol this rule applies to.

* `source_port_range` - The Source Port or Range.

* `destination_port_range` - The Destination Port or Range.

* `source_address_prefix` - CIDR or source IP range or * to match any IP.

* `source_address_prefixes` - A list of CIDRs or source IP ranges.

* `destination_address_prefix` - CIDR or destination IP range or * to match any IP.

* `destination_address_prefixes` - A list of CIDRs or destination IP ranges.

* `source_application_security_group_ids` - A List of source Application Security Group ID's

* `destination_application_security_group_ids` - A List of destination Application Security Group ID's

* `access` - Is network traffic is allowed or denied?

* `priority` - The priority of the rule

* `direction` - The direction specifies if rule will be evaluated on incoming or outgoing traffic.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Network Security Group.
