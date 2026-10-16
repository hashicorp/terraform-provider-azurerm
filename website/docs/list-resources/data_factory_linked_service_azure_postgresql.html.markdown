---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_linked_service_azure_postgresql"
description: |-
  Lists Data Factory Linked Service Azure PostgreSQL resources.
---

# List resource: azurerm_data_factory_linked_service_azure_postgresql

Lists Data Factory Linked Service Azure PostgreSQL resources.

## Example Usage

### List all Data Factory Linked Service Azure PostgreSQL in a specific Data Factory

```hcl
list "azurerm_data_factory_linked_service_azure_postgresql" "example" {
  provider = azurerm
  config {
    data_factory_id = azurerm_data_factory.example.id
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `data_factory_id` - (Required) The ID of the data factory to query.
