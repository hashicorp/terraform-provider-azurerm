---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_fleet"
description: |-
  Manages a Cosmos DB Fleets.
---

# azurerm_cosmosdb_fleet

Manages a Cosmos DB Fleets.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cosmosdb_fleet" "example" {
  name                = "fleet-test"
  resource_group_name = "azurerm_resource_group.example.name"
  location            = "azurerm_resource_group.example.location"
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Cosmos DB Fleets should exist. Changing this forces new Cosmos DB Fleets to be created.

* `name` - (Required) The name which should be used for this Cosmos DB Fleets. Changing this forces new Cosmos DB Fleets to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Cosmos DB Fleets should exist. Changing this forces new Cosmos DB Fleets to be created.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Cosmos DB Fleets. Changing this forces new Cosmos DB Fleets to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Cosmos DB Fleets.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cosmos DB Fleets.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cosmos DB Fleets.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cosmos DB Fleets.

## Import

Cosmos DB Fleets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmosdb_fleet.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DocumentDB/fleets/fleets1
```
