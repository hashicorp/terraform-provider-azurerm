---
subcategory: "Discovery"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_discovery_bookshelf_private_endpoint_connection"
description: |-
  Lists Discovery Bookshelf Private Endpoint Connection resources.
---

# List resource: azurerm_discovery_bookshelf_private_endpoint_connection

Lists Discovery Bookshelf Private Endpoint Connection resources.

## Example Usage

### List Private Endpoint Connections in a Discovery Bookshelf

```hcl
list "azurerm_discovery_bookshelf_private_endpoint_connection" "example" {
  provider = azurerm
  config {
    bookshelf_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Discovery/bookshelves/example-bookshelf"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `bookshelf_id` - (Required) The ID of the Discovery Bookshelf to query.
