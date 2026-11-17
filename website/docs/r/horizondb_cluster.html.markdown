---
subcategory: "HorizonDB"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_horizondb_cluster"
description: |-
  Manages a HorizonDB Cluster.
---

# azurerm_horizondb_cluster

Manages a HorizonDB Cluster.

-> **Note:** The required `read_write_pool` block provisions compute after the Cluster is created. This resource manages the Cluster's designated read/write Pool, not other Pools under the Cluster.

## Example Usage

```hcl
variable "horizondb_administrator_password" {
  type        = string
  description = "The PostgreSQL administrator password for the HorizonDB Cluster."
  sensitive   = true
}

resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "East US 2"
}

resource "azurerm_horizondb_cluster" "example" {
  name                = "example-horizondb-cluster"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  postgresql_version  = "18"

  administrator_login    = "exampleadmin"
  administrator_password = var.horizondb_administrator_password

  read_write_pool {
    name                   = "example-pool"
    node_count             = 1
    high_availability_mode = "Disabled"

    compute {
      type       = "Provisioned"
      processor  = "AMD"
      generation = 1
      vcores     = 4
    }
  }

  tags = {
    environment = "example"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the HorizonDB Cluster. Changing this forces a new resource to be created.

~> **Note:** The name must contain between 1 and 63 lowercase letters, numbers, or hyphens and must start and end with a letter or number.

* `resource_group_name` - (Required) The name of the Resource Group where the HorizonDB Cluster should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the HorizonDB Cluster should exist. Changing this forces a new resource to be created.

* `administrator_login` - (Required) The PostgreSQL administrator login name, which must be between 1 and 63 characters in length. Changing this forces a new resource to be created.

* `read_write_pool` - (Required) A `read_write_pool` block as defined below.

~> **Note:** Refresh and import resolve the read/write Pool using the Cluster's `readWritePool.designatedPoolResourceId`, even if the Pool was created outside Terraform. They fail if the designation is missing or the designated Pool cannot be retrieved. The block cannot be omitted or removed from configuration. Do not manage the same Pool's compute configuration through another resource or deployment.

-> **Note:** New Pools use the default parameter group for the Cluster's region and PostgreSQL version with `Immediate` application. Compute updates preserve the Pool's parameter-group association, and this resource does not automatically restart Pools.

* `administrator_password` - (Optional) The PostgreSQL administrator password, which must be between 8 and 128 characters in length.

~> **Note:** Exactly one of `administrator_password` and `administrator_password_wo` must be specified. The value of `administrator_password` is stored in Terraform state.

* `administrator_password_wo` - (Optional, Write-Only) The PostgreSQL administrator password, which must be between 8 and 128 characters in length. This value is not stored in Terraform state.

~> **Note:** Write-only arguments require Terraform 1.11 or later. `administrator_password_wo` and `administrator_password_wo_version` must be specified together.

* `administrator_password_wo_version` - (Optional) A positive integer used to trigger an update for `administrator_password_wo`. Increment this value when updating the write-only password.

* `backup_retention_days` - (Optional) The backup retention period in days. Possible values range between `1` and `35`. Defaults to `7`.

~> **Note:** Removing this argument from an existing configuration resets the retention period to `7` days.

* `create_mode` - (Optional) The creation mode of the HorizonDB Cluster. Possible values are `Create` and `PointInTimeRestore`. Defaults to `Create`. Changing this forces a new resource to be created.

* `geo_redundant_backup_enabled` - (Optional) Whether `geo redundant backup` is enabled. Defaults to `false`. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `postgresql_version` - (Optional) The PostgreSQL major version. Possible values are `17` and `18`. Changing this forces a new resource to be created.

~> **Note:** `postgresql_version` is required when `create_mode` is `Create` and must be omitted when `create_mode` is `PointInTimeRestore`. For restored Clusters, the version is obtained from the source backup and exported by this resource.

* `public_network_access_enabled` - (Optional) Whether `public network access` is enabled. Defaults to `false`.

* `restore` - (Optional) A `restore` block as defined below. Changing this forces a new resource to be created.

~> **Note:** This block is required when `create_mode` is `PointInTimeRestore` and must be omitted when `create_mode` is `Create`.

-> **Note:** HorizonDB can omit restore-source details and the restore timestamp from API responses. This resource preserves the existing `restore` values in Terraform state when those details are not returned.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `compute` block supports the following:

* `generation` - (Required) The compute generation.

~> **Note:** The generation must be between `1` and `2147483647`. Azure validates availability for the selected region and processor.

* `processor` - (Required) The processor type. Possible values are `AMD` and `Intel`.

* `vcores` - (Required) The provisioned compute capacity in vCores.

~> **Note:** The number of vCores must be between `2` and `192`. Azure validates the supported combinations of processor, generation, and capacity.

* `type` - (Optional) The compute model type. The only possible value is `Provisioned`. Defaults to `Provisioned`.

---

An `identity` block supports the following:

* `type` - (Required) The type of Managed Identity assigned to the Cluster. The only possible value is `UserAssigned`.

* `identity_ids` - (Required) A set of User Assigned Identity IDs assigned to the Cluster.

-> **Note:** Identity updates replace the complete set of assigned identities in place. Removing the `identity` block clears all assignments.

---

A `read_write_pool` block supports the following:

* `compute` - (Required) A `compute` block as defined above.

* `high_availability_mode` - (Required) The high-availability mode. Possible values are `AnyZone`, `Disabled`, and `ZoneRedundant`.

~> **Note:** Use `Disabled` with one node, or `AnyZone` or `ZoneRedundant` with two or more nodes. Change the mode and node count together when crossing these ranges. `ZoneRedundant` requires allocation across availability zones; changing the mode can replace nodes and disconnect their clients.

* `name` - (Required) The name of the Pool.

~> **Note:** The name must contain between 1 and 63 lowercase letters, numbers, or hyphens, must start and end with a letter or number, and must differ from the Cluster name. In-place renaming is rejected. Replace the Cluster to use a different Pool name.

* `node_count` - (Required) The number of nodes allocated to the Pool.

~> **Note:** The node count must be between `1` and `16`. Zero nodes is not a supported configuration; stopping compute is a separate Azure operation.

---

A `restore` block supports the following:

* `point_in_time_utc` - (Required) The restore timestamp in RFC3339 format.

~> **Note:** For a same-region restore, HorizonDB uses the latest available backup at or before this timestamp. For a cross-region restore, the timestamp is required but ignored by the service, which uses the latest available geo-replicated backup.

* `source_horizondb_cluster_id` - (Required) The ID of the source HorizonDB Cluster.

* `source_location` - (Optional) The Azure Region of the source HorizonDB Cluster.

~> **Note:** This is required for cross-region restores and optional for same-region restores. Cross-region restores require an eligible geo-replicated backup from the source Cluster.

-> **Note:** Restoring a Cluster does not clone the source Pool or its child resources. A new Pool is created from the required `read_write_pool` block after restoration using the restored PostgreSQL version to select its default parameter group.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the HorizonDB Cluster.

* `read_write_pool` - A `read_write_pool` block as defined below.

---

A `read_write_pool` block exports the following:

* `id` - The ID of the HorizonDB Pool.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the HorizonDB Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the HorizonDB Cluster.
* `update` - (Defaults to 1 hour) Used when updating the HorizonDB Cluster.
* `delete` - (Defaults to 1 hour) Used when deleting the HorizonDB Cluster.

## Import

A HorizonDB Cluster can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_horizondb_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.HorizonDb/clusters/cluster1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.HorizonDb` - 2026-10-01-preview
