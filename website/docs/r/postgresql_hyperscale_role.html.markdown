---
subcategory: "PostgreSQL HyperScale"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_postgresql_hyperscale_role"
description: |-
  Manages a PostgreSQL HyperScale Role.
---

# azurerm_postgresql_hyperscale_role

Manages a PostgreSQL HyperScale Role.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_postgresql_hyperscale_cluster" "example" {
  name                = "example-postgresqlhscsg"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_postgresql_hyperscale_role" "example" {
  name            = "example-postgresqlhscrole"
  server_group_id = azurerm_postgresql_hyperscale_cluster.example.id
  password        = "H@Sh1CoR3!"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this PostgreSQL HyperScale Role. Changing this forces a new resource to be created.

* `server_group_id` - (Required) The ID of the PostgreSQL HyperScale Cluster. Changing this forces a new resource to be created.

* `password` - (Required) The password of the PostgreSQL HyperScale Role.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the PostgreSQL HyperScale Role.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the PostgreSQL HyperScale Role.
* `read` - (Defaults to 5 minutes) Used when retrieving the PostgreSQL HyperScale Role.
* `delete` - (Defaults to 30 minutes) Used when deleting the PostgreSQL HyperScale Role.

## Import

PostgreSQL HyperScale Roles can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_postgresql_hyperscale_role.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.DBforPostgreSQL/serverGroupsv2/cluster1/roles/role1
```
