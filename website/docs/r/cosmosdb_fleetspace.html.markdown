---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_fleetspace"
description: |-
  Manages a Cosmos DB Fleetspace.
---

# azurerm_cosmosdb_fleetspace

Manages a Cosmos DB Fleetspace.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "southeastasia"
}

resource "azurerm_cosmosdb_fleet" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_cosmosdb_fleetspace" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  fleet_name          = azurerm_cosmosdb_fleet.example.name
  service_tier        = "GeneralPurpose"
  data_regions = [
    azurerm_resource_group.example.location
  ]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Cosmos DB Fleetspace. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Cosmos DB Fleetspace should exist. Changing this forces a new resource to be created.

* `data_regions` - (Required) The list of data regions assigned to the Cosmos DB Fleetspace. Changing this forces a new resource to be created.

* `fleet_name` - (Required) The Cosmos DB Fleet name in which the Cosmos DB Fleetspace is created. Changing this forces a new resource to be created.

* `service_tier` - (Required) The service tier for the Cosmos DB Fleetspace. Possible values include `GeneralPurpose` and `BusinessCritical`. `GeneralPurpose` types refers to single write region accounts that can be added to the Cosmos DB Fleetspace, whereas `BusinessCritical` refers to multi write region. Changing this forces a new resource to be created.

---

* `maximum_throughput` - (Optional) The maximum throughput for the throughput pool of the Cosmos DB Fleetspace. Must be divisible by 1000, more than or equal to `minimum_throughput`, and less than or equal to 10 times of `minimum_throughput`.

* `minimum_throughput` - (Optional) The minimum throughput for the throughput pool of the Cosmos DB Fleetspace. Must be divisible by 1000, more than or equal to 100000, and less than or equal to 10000000.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Cosmos DB Fleetspace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cosmos DB Fleetspace.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cosmos DB Fleetspace.
* `update` - (Defaults to 30 minutes) Used when updating the Cosmos DB Fleetspace.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cosmos DB Fleetspace.

## Import

Cosmos DB Fleetspaces can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmosdb_fleetspace.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DocumentDB/fleets/fleet1/fleetspaces/fleetspace1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.DocumentDB` - 2025-10-15
