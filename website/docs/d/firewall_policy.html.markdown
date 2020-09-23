---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_firewall_policy"
description: |-
  Gets information about an existing Firewall Policy.
---

# Data Source: azurerm_firewall_policy

Use this data source to access information about an existing Firewall Policy.

## Example Usage

```hcl
data "azurerm_firewall_policy" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_firewall_policy.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Firewall Policy.

* `resource_group_name` - (Required) The name of the Resource Group where the Firewall Policy exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Firewall Policy.

* `tags` - A mapping of tags assigned to the Firewall Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Firewall Policy.
