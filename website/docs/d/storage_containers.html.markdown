---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_storage_containers"
description: |-
  Gets information about an existing Storage Containers.
---

# Data Source: azurerm_storage_containers

Use this data source to access information about the existing Storage Containers within a Storage Account.

## Example Usage

```hcl
data "azurerm_storage_containers" "example" {
  storage_account_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Storage/storageAccounts/sa1"
}

output "container_id" {
  value = data.azurerm_storage_containers.example.containers[0].resource_manager_id
}
```

## Arguments Reference

The following arguments are supported:

* `storage_account_id` - (Required) The ID of the Storage Account that the Storage Containers reside in.

---

* `name_prefix` - (Optional) A prefix match used for the Storage Container `name` field.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Storage Containers.

* `containers` - A `containers` block as defined below.

---

A `containers` block exports the following:

* `data_plane_id` - The data plane ID of the Storage Container.

* `name` - The name of this Storage Container.

* `resource_manager_id` - The resource manager ID of the Storage Container.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Containers.
