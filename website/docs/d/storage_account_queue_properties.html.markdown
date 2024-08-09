---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_account_queue_properties"
description: |-
  Gets information about an existing Storage Account Queue Properties.

---

# Data Source: azurerm_storage_account_queue_properties

Use this data source to access information about an existing Storage Account Queue Properties.

## Example Usage

```hcl
data "azurerm_storage_account_queue_properties" "example" {
  storage_account_id = azurerm_storage_account.example.id
}

output "storage_account_primary_queue_endpoint" {
  value = data.azurerm_storage_account_queue_properties.example.primary_queue_endpoint
}
```

## Argument Reference

* `storage_account_id` - Specifies the name of the Storage Account

## Attributes Reference

* `id` - The ID of the Storage Account Queue Properties resource.

* `storage_account_id` - The ID of the Storage Account.

* `primary_queue_endpoint` - The endpoint URL for queue storage in the primary location.

* `primary_queue_host` - The hostname with port if applicable for queue storage in the primary location.

* `primary_queue_microsoft_endpoint` - The microsoft routing endpoint URL for queue storage in the primary location.

* `primary_queue_microsoft_host` - The microsoft routing hostname with port if applicable for queue storage in the primary location.

* `secondary_queue_endpoint` - The endpoint URL for queue storage in the secondary location.

* `secondary_queue_host` - The hostname with port if applicable for queue storage in the secondary location.

* `secondary_queue_microsoft_endpoint` - The microsoft routing endpoint URL for queue storage in the secondary location.

* `secondary_queue_microsoft_host` - The microsoft routing hostname with port if applicable for queue storage in the secondary location.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Account Queue Properties.
