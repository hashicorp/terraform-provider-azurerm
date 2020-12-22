---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_disk_access"
description: |-
  Gets information about an existing Disk Access.
---

# Data Source: azurerm_disk_access

Use this data source to access information about an existing Disk Access.

## Example Usage

```hcl
data "azurerm_disk_access" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_disk_access.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Disk Access.

* `resource_group_name` - (Required) The name of the Resource Group where the Disk Access exists.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Disk Access.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Disk Access.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Disk.