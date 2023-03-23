---
subcategory: "PostgreSQL HyperScale"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_postgresql_hyperscale_cluster"
description: |-
  Manages a PostgreSQL HyperScale Cluster.
---

# azurerm_postgresql_hyperscale_cluster

Manages a PostgreSQL HyperScale Cluster.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_postgresql_hyperscale_cluster" "example" {
  name                            = "example-cluster"
  resource_group_name             = azurerm_resource_group.example.name
  location                        = azurerm_resource_group.example.location
  administrator_login_password    = "H@Sh1CoR3!"
  coordinator_storage_quota_in_mb = 131072
  coordinator_vcores              = 2
  node_count                      = 0
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this PostgreSQL HyperScale Cluster. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the PostgreSQL HyperScale Cluster should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the PostgreSQL HyperScale Cluster should exist. Changing this forces a new resource to be created.

* `administrator_login_password` - (Required) The password of the administrator login.

* `coordinator_storage_quota_in_mb` - (Required) The coordinator storage allowed for the PostgreSQL HyperScale Cluster. Possible values are `32768`, `65536`, `131072`, `262144`, `524288`, `1048576`, `2097152`, `4194304`, `8388608` and `16777216`.

* `coordinator_vcores` - (Required) The coordinator vCores count for the PostgreSQL HyperScale Cluster. Possible values are `1`, `2`, `4`, `8`, `16`, `32`, `64` and `96`.

* `node_count` - (Required) The worker node count of the PostgreSQL HyperScale Cluster. Possible value is between `0` and `20` except `1`.

* `citus_version` - (Optional) The citus extension version on the PostgreSQL HyperScale Cluster. Possible values are `8.3`, `9.0`, `9.1`, `9.2`, `9.3`, `9.4`, `9.5`, `10.0`, `10.1`, `10.2`, `11.0`, `11.1` and `11.2`.

* `coordinator_public_ip_access_enabled` - (Optional) Is public access enabled on coordinator? Defaults to `true`.

* `coordinator_server_edition` - (Optional) The edition of the coordinator server. Defaults to `GeneralPurpose`.

* `ha_enabled` - (Optional) Is high availability enabled for the PostgreSQL HyperScale cluster? Defaults to `false`.

* `maintenance_window` - (Optional) A `maintenance_window` block as defined below.

* `node_public_ip_access_enabled` - (Optional) Is public access enabled on worker nodes. Defaults to `false`.

* `node_server_edition` - (Optional) The edition of the node server. Defaults to `MemoryOptimized`.

* `node_storage_quota_in_mb` - (Optional) The storage quota in MB on each worker node. Possible values are `32768`, `65536`, `131072`, `262144`, `524288`, `1048576`, `2097152`, `4194304`, `8388608` and `16777216`.

* `node_vcores` - (Optional) The vCores count on each worker node. Possible values are `1`, `2`, `4`, `8`, `16`, `32`, `64`, `96` and `104`.

* `point_in_time_in_utc` - (Optional) The date and time in UTC (ISO8601 format) for the PostgreSQL HyperScale cluster restore. Changing this forces a new resource to be created.

* `preferred_primary_zone` - (Optional) The preferred primary availability zone for the PostgreSQL HyperScale cluster.

* `shards_on_coordinator_enabled` - (Optional) Is shards on coordinator enabled for the PostgreSQL HyperScale cluster.

* `source_location` - (Optional) The Azure region of the source PostgreSQL HyperScale cluster for read replica clusters. Changing this forces a new resource to be created.

* `source_resource_id` - (Optional) The resource ID of the source PostgreSQL HyperScale cluster for read replica clusters. Changing this forces a new resource to be created.

* `sql_version` - (Optional) The major PostgreSQL version on the PostgreSQL HyperScale cluster. Possible values are `11`, `12`, `13`, `14` and `15`.

* `tags` - (Optional) A mapping of tags which should be assigned to the PostgreSQL HyperScale Cluster.

---

A `maintenance_window` block supports the following:

* `day_of_week` - (Optional) The day of week for maintenance window, where the week starts on a Sunday, i.e. Sunday = `0`, Monday = `1`. Defaults to `0`.

* `start_hour` - (Optional) The start hour for maintenance window. Defaults to `0`.

* `start_minute` - (Optional) The start minute for maintenance window. Defaults to `0`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the PostgreSQL HyperScale Cluster.

* `earliest_restore_time` - The earliest restore point time (ISO8601 format) for the PostgreSQL HyperScale Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the PostgreSQL HyperScale Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the PostgreSQL HyperScale Cluster.
* `update` - (Defaults to 60 minutes) Used when updating the PostgreSQL HyperScale Cluster.
* `delete` - (Defaults to 60 minutes) Used when deleting the PostgreSQL HyperScale Cluster.

## Import

PostgreSQL HyperScale Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_postgresql_hyperscale_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.DBforPostgreSQL/serverGroupsv2/cluster1
```
