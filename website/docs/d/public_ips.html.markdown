---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_public_ips"
description: |-
  Gets information about a set of existing Public IP Addresses.
---

# Data Source: azurerm_public_ips

Use this data source to access information about a set of existing Public IP Addresses.

## Example Usage

```hcl
data "azurerm_public_ips" "example" {
  resource_group_name = "pip-test"
  attachment_status   = "Attached"
}
```

## Argument Reference

* `resource_group_name` - Specifies the name of the resource group.
* `attachment_status` - (Optional) Filter to include IP Addresses which are attached to a device, such as a VM/LB (`Attached`) or unattached (`Unattached`).
* `name_prefix` - (Optional) A prefix match used for the IP Addresses `name` field, case sensitive.
* `allocation_type` - (Optional) The Allocation Type for the Public IP Address. Possible values include `Static` or `Dynamic`.

## Attributes Reference

* `public_ips` - A List of `public_ips` blocks as defined below filtered by the criteria above.

A `public_ips` block contains:

* `id` - The ID of the Public IP Address
* `domain_name_label` - The Domain Name Label of the Public IP Address
* `fqdn` - The FQDN of the Public IP Address
* `name` - The Name of the Public IP Address
* `ip_address` - The IP address of the Public IP Address

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Public IP Addresses.
