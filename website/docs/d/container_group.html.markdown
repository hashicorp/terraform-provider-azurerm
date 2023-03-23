---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_group"
description: |-
  Gets information about an existing Container Group instance.
---

# Data Source: azurerm_container_group

Use this data source to access information about an existing Container Group instance.

## Example Usage

```hcl
data "azurerm_container_group" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_container_group.example.id
}

output "ip_address" {
  value = data.azurerm_container_group.example.ip_address
}

output "fqdn" {
  value = data.azurerm_container_group.example.fqdn
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Container Group instance.

* `resource_group_name` - (Required) The name of the Resource Group where the Container Group instance exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container Group instance.

* `ip_address` - The IP address allocated to the Container Group instance.

* `fqdn` - The FQDN of the Container Group instance derived from `dns_name_label`.

* `location` - The Azure Region where the Container Group instance exists.

* `identity` - A `identity` block as defined below.

* `subnet_ids` - The subnet resource IDs for a container group.

* `zones` - A list of Availability Zones in which this Container Group is located.

* `tags` - A mapping of tags assigned to the Container Group instance.

---

A `identity` block exports the following:

* `type` - Type of Managed Service Identity configured on this Container Group.

* `principal_id` - The Principal ID of the System Assigned Managed Service Identity that is configured on this Container Group.

* `tenant_id` - The Tenant ID of the System Assigned Managed Service Identity that is configured on this Container Group.

* `identity_ids` - The list of User Assigned Managed Identity IDs assigned to this Container Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Container Group instance.
