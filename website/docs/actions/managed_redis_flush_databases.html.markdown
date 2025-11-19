---
subcategory: "Managed Redis"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_redis_databases_flush"
description: |-
  Flushes all the keys in the Managed Redis Database and also from its linked databases.
---

# Action: azurerm_managed_redis_databases_flush

~> **Note:** `azurerm_managed_redis_databases_flush` is in beta. Its interface and behaviour may change as the feature evolves, and breaking changes are possible. It is offered as a technical preview without compatibility guarantees until Terraform 1.14 is generally available.

Flushes all the keys in the Managed Redis Database and also from its linked databases.


## Example Usage

```terraform
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_managed_redis" "example" {
  name                = "example-managed-redis"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku_name            = "Balanced_B3"

  default_database {
    geo_replication_group_name = "myGeoGroup"
  }
}

action "azurerm_managed_redis_databases_flush" "example" {
  config {
    managed_redis_database_id = azurerm_managed_redis.example.default_database[0].id
  }
}
```

## Argument Reference

This action supports the following arguments:

* `managed_redis_database_id` - (Required) The ID of the Managed Redis Database where the keys will be flushed.

* `linked_database_ids` - (Optional) The IDs of any Linked Databases where the keys will be flushed.

---

* `timeout` - (Optional) Timeout duration for the action to complete. Defaults to `15m`.
