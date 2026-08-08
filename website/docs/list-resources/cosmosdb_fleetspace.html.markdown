---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_fleetspace"
description: |-
  Lists Cosmos DB Fleetspace resources.
---

# List resource: azurerm_cosmosdb_fleetspace

Lists Cosmos DB Fleetspace resources.

## Example Usage

### List all Cosmos DB Fleetspaces in a specific Fleet

```hcl
list "azurerm_cosmosdb_fleetspace" "example" {
  provider = azurerm
  config {
    fleet_id = azurerm_cosmosdb_fleet.example.id
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `fleet_id` - (Required) The ID of the Cosmos DB Fleet to query.
