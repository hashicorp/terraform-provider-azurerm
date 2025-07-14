---
subcategory: "Azure Stack HCI"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_stack_hci_storage_path"
description: |-
  Gets information about an existing Stack HCI Storage Path.
---

# Data Source: azurerm_stack_hci_storage_path

Use this data source to access information about an existing Stack HCI Storage Path.

## Example Usage

```hcl
data "azurerm_stack_hci_storage_path" "example" {
  name                = "example-hci-storage-path-name"
  resource_group_name = "example-rg"
}

output "id" {
  value = data.azurerm_stack_hci_storage_path.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Stack HCI Storage Path.

* `resource_group_name` - (Required) The name of the Resource Group where the Stack HCI Storage Path exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Stack HCI Storage Path.

* `custom_location_id` - The ID of the Custom Location where the Stack HCI Storage Path exists.

* `location` - The Azure Region where the Stack HCI Storage Path exists.

* `path` - The file path on the disk where the Stack HCI Storage Path was created.

* `tags` - A mapping of tags assigned to the Stack HCI Storage Path.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Stack HCI Storage Path.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.AzureStackHCI`: 2024-01-01
