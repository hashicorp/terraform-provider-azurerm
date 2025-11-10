---
subcategory: "Mongo Cluster"
layout: "azurerm"
page_title: "Azure Resource Manager: `azurerm_mongo_cluster"
description: |-
  Manages a MongoDB Cluster using vCore Architecture.
---

# azurerm_mongo_cluster

Manages a MongoDB Cluster using vCore Architecture.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "East US"
}

resource "azurerm_mongo_cluster" "example" {
  name                   = "example-mc"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  administrator_username = "adminTerraform"
  administrator_password = "QAZwsx123"
  shard_count            = "1"
  compute_tier           = "Free"
  high_availability_mode = "Disabled"
  storage_size_in_gb     = "32"
}

```

## Example Usage (with Microsoft Entra ID Authentication)

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "East US"
}

resource "azurerm_mongo_cluster" "example" {
  name                      = "example-mc"
  resource_group_name       = azurerm_resource_group.example.name
  location                  = azurerm_resource_group.example.location
  administrator_username    = "adminTerraform"
  administrator_password    = "QAZwsx123"
  shard_count               = "1"
  compute_tier              = "M30"
  high_availability_mode    = "ZoneRedundantPreferred"
  storage_size_in_gb        = "64"
  auth_config_allowed_modes = ["NativeAuth", "MicrosoftEntraID"]
}
```

## Example Usage (Preview feature GeoReplicas)

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "East US"
}

resource "azurerm_mongo_cluster" "example" {
  name                   = "example-mc"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  administrator_username = "adminTerraform"
  administrator_password = "QAZwsx123"
  shard_count            = "1"
  compute_tier           = "M30"
  high_availability_mode = "ZoneRedundantPreferred"
  storage_size_in_gb     = "64"
  preview_features       = ["GeoReplicas"]
}

resource "azurerm_mongo_cluster" "example_geo_replica" {
  name                = "example-mc-geo"
  resource_group_name = azurerm_resource_group.example.name
  location            = "Central US"
  source_server_id    = azurerm_mongo_cluster.example.id
  source_location     = azurerm_mongo_cluster.example.location
  create_mode         = "GeoReplica"

  lifecycle {
    ignore_changes = ["administrator_username", "high_availability_mode", "preview_features", "shard_count", "storage_size_in_gb", "compute_tier", "version"]
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for the MongoDB Cluster. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the MongoDB Cluster. Changing this forces a new resource to be created.

* `location` - (Required) The supported Azure location where the MongoDB Cluster exists. Changing this forces a new resource to be created.

* `administrator_username` - (Optional) The administrator username of the MongoDB Cluster. Changing this forces a new resource to be created.

* `create_mode` - (Optional) The creation mode for the MongoDB Cluster. Possibles values are `Default`, `GeoReplica` and `PointInTimeRestore`. Defaults to `Default`. Changing this forces a new resource to be created.

* `customer_managed_key` - (Optional) A `customer_managed_key` block as defined below. Changing this forces a new resource to be created.

* `data_api_mode_enabled` - (Optional) Is the Data API for the MongoDB Cluster enabled? Defaults to `false`. Once enabled, changing this back to `false` will force a new resource to be created.

* `identity` - (Optional) An `identity` block as detailed below. Changing this forces a new resource to be created.

* `preview_features` - (Optional) The preview features that can be enabled on the MongoDB Cluster. Changing this forces a new resource to be created.

* `restore` - (Optional) A `restore` block as defined below. Required when `create_mode` is set to `PointInTimeRestore`. Changing this forces a new resource to be created.

* `shard_count` -  (Optional) The Number of shards to provision on the MongoDB Cluster. Changing this forces a new resource to be created.

* `source_location` - (Optional) The location of the source MongoDB Cluster. Changing this forces a new resource to be created.

* `source_server_id` - (Optional) The ID of the replication source MongoDB Cluster. Changing this forces a new resource to be created.

* `administrator_password` - (Optional) The Password associated with the `administrator_username` for the MongoDB Cluster.

* `auth_config_allowed_modes` - (Optional) A list of allowed authentication modes for the MongoDB Cluster. Possible values are `NativeAuth` and `MicrosoftEntraID`.

* `compute_tier` - (Optional) The compute tier to assign to the MongoDB Cluster. Possible values are `Free`, `M10`, `M20`, `M25`, `M30`, `M40`, `M50`, `M60`, `M80`, and `M200`.

* `high_availability_mode` - (Optional) The high availability mode for the MongoDB Cluster. Possibles values are `Disabled` and `ZoneRedundantPreferred`.

* `public_network_access` - (Optional) The Public Network Access setting for the MongoDB Cluster. Possibles values are `Disabled` and `Enabled`. Defaults to `Enabled`.

* `storage_size_in_gb` - (Optional) The size of the data disk space for the MongoDB Cluster.

* `tags` - (Optional) A mapping of tags to assign to the MongoDB Cluster.

* `version` - (Optional) The version for the MongoDB Cluster. Possibles values are `5.0`, `6.0`, `7.0` and `8.0`.

---

A `customer_managed_key` block supports the following:

* `key_vault_key_id` - (Required) The ID of the key vault key used for encryption. For example: `https://example-vault-name.vault.azure.net/keys/example-key-name/a1b2c3d4`.

* `user_assigned_identity_id` - (Required) The ID of the User Assigned Identity that has access to the Key Vault Key.

---

An `identity` block supports the following:

* `type` - (Required) The type of managed identity to assign. Possible value is `UserAssigned`. Changing this forces a new resource to be created.

* `identity_ids` - (Optional) - A list of one or more Resource IDs for User Assigned Managed identities to assign. Required when `type` is set to `UserAssigned`. Changing this forces a new resource to be created.

---

A `restore` block supports the following:

* `point_in_time_utc` - (Required) The point in time (in UTC) to restore from, in ISO 8601 format (e.g., `2024-01-01T00:00:00Z`). Changing this forces a new resource to be created.

* `source_id` - (Required) The ID of the source MongoDB Cluster to restore from. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the MongoDB Cluster.

* `connection_strings` - The list of `connection_strings` blocks as defined below.

---

A `connection_strings` exports the following:

* `name` - The name of the connection string.

* `description` - The description of the connection string.

* `value` - The value of the Mongo Cluster connection string. The `<user>:<password>` placeholder returned from API will be replaced by the real `administrator_username` and `administrator_password` if available in the state.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the MongoDB Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the MongoDB Cluster.
* `update` - (Defaults to 1 hour) Used when updating the MongoDB Cluster.
* `delete` - (Defaults to 1 hour) Used when deleting the MongoDB Cluster.

## Import

MongoDB Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mongo_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DocumentDB/mongoClusters/myMongoCluster
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.DocumentDB` - 2025-09-01
