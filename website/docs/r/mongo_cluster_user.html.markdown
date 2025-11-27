---
subcategory: "Mongo Cluster"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mongo_cluster_user"
description: |-
  Manages a Mongo Cluster User.
---

# azurerm_mongo_cluster_user

Manages a Mongo Cluster User.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_mongo_cluster" "example" {
  name                   = "example-mc"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  administrator_username = "adminTerraform"
  administrator_password = "QAZwsx123"
  shard_count            = "1"
  compute_tier           = "M30"
  high_availability_mode = "Disabled"
  storage_size_in_gb     = "32"
  version                = "8.0"
}

resource "azurerm_mongo_cluster_user" "example" {
  object_id              = data.azurerm_client_config.current.object_id
  mongo_cluster_id       = azurerm_mongo_cluster.example.id
  identity_provider_type = "MicrosoftEntraID"
  principal_type         = "servicePrincipal"

  role {
    database = "admin"
    role     = "root"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `mongo_cluster_id` - (Required) The ID of the Mongo Cluster where the User should exist. Changing this forces a new resource to be created.

* `object_id` - (Required) The object ID of a user, service principal for the Mongo Cluster User. Changing this forces a new resource to be created.

* `principal_type` - (Required) The principal type for the Mongo Cluster User. Possible values are `user` and `servicePrincipal`. Changing this forces a new resource to be created.

---

* `role` - (Optional) One or more `role` blocks as defined below. Changing this forces a new resource to be created.


---

A `role` block supports the following:

* `database` - (Required) The database name for the role.

* `role` - (Required) The role name. Possible values are `root`. 

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mongo Cluster User.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Mongo Cluster User.
* `read` - (Defaults to 5 minutes) Used when retrieving the Mongo Cluster User.
* `delete` - (Defaults to 30 minutes) Used when deleting the Mongo Cluster User.

## Import

Mongo Cluster Users can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mongo_cluster_user.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DocumentDB/mongoClusters/cluster1/users/user1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.DocumentDB` - 2025-09-01
