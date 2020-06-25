---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_web_application_firewall_policy"
description: |-
  Gets information about an existing Web Application Firewall Policy.
---

# Data Source: azurerm_web_application_firewall_policy

Use this data source to access information about an existing Web Application Firewall Policy.

## Example Usage

```hcl
data "azurerm_web_application_firewall_policy" "example" {
  resource_group_name = "existing"
  name                = "existing"
}

output "id" {
  value = data.azurerm_web_application_firewall_policy.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Web Application Firewall Policy

* `resource_group_name` - (Required) The name of the Resource Group where the Web Application Firewall Policy exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Web Application Firewall Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Web Application Firewall Policy.
