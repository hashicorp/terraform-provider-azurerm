---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_fleet"
description: |-
  Manages a Cosmos DB Fleet.
---

# azurerm_cosmosdb_fleet

Manages a Cosmos DB Fleet.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cosmosdb_fleet" "example" {
  name                = "fleet-test"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Cosmos DB Fleet. Changing this forces new Cosmos DB Fleet to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Cosmos DB Fleet should exist. Changing this forces new Cosmos DB Fleet to be created.

* `location` - (Required) The Azure Region where the Cosmos DB Fleet should exist. Changing this forces new Cosmos DB Fleet to be created.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Cosmos DB Fleet. Changing this forces new Cosmos DB Fleet to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Cosmos DB Fleet.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cosmos DB Fleet.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cosmos DB Fleet.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cosmos DB Fleet.

## Import

Cosmos DB Fleet can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmosdb_fleet.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DocumentDB/fleets/fleet1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.DocumentDB` - 2025-10-15
