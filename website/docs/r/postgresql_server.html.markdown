---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_postgresql_server"
description: |-
  Manages a PostgreSQL Server.
---

# azurerm_postgresql_server

Manages a PostgreSQL Server.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_postgresql_server" "example" {
  name                = "example-psqlserver"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  administrator_login          = "psqladmin"
  administrator_login_password = "H@Sh1CoR3!"

  sku_name   = "GP_Gen5_4"
  version    = "11"
  storage_mb = 640000

  backup_retention_days        = 7
  geo_redundant_backup_enabled = true
  auto_grow_enabled            = true

  public_network_access_enabled    = false
  ssl_enforcement_enabled          = true
  ssl_minimal_tls_version_enforced = "TLS1_2"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the PostgreSQL Server. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the PostgreSQL Server. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku_name` - (Required) Specifies the SKU Name for this PostgreSQL Server. The name of the SKU, follows the `tier` + `family` + `cores` pattern (e.g. `B_Gen4_1`, `GP_Gen5_8`). For more information see the [product documentation](https://docs.microsoft.com/rest/api/postgresql/singleserver/servers/create#sku). Possible values are `B_Gen4_1`, `B_Gen4_2`, `B_Gen5_1`, `B_Gen5_2`, `GP_Gen4_2`, `GP_Gen4_4`, `GP_Gen4_8`, `GP_Gen4_16`, `GP_Gen4_32`, `GP_Gen5_2`, `GP_Gen5_4`, `GP_Gen5_8`, `GP_Gen5_16`, `GP_Gen5_32`, `GP_Gen5_64`, `MO_Gen5_2`, `MO_Gen5_4`, `MO_Gen5_8`, `MO_Gen5_16` and `MO_Gen5_32`.

~> **Note:** When replication is set up and `sku_name` is changed to a higher tier or more capacity for the primary, all replicas are scaled up to the same tier/capacity. This is an Azure requirement, for more information see the [replica scaling documentation](https://docs.microsoft.com/azure/postgresql/concepts-read-replicas#scaling)

* `version` - (Required) Specifies the version of PostgreSQL to use. Valid values are `9.5`, `9.6`, `10`, `10.0`, `10.2` and `11`. Changing this forces a new resource to be created.

* `administrator_login` - (Optional) The Administrator login for the PostgreSQL Server. Required when `create_mode` is `Default`. Changing this forces a new resource to be created.

* `administrator_login_password` - (Optional) The Password associated with the `administrator_login` for the PostgreSQL Server.

* `administrator_login_password_wo` - (Optional, Write-Only) The Password associated with the `administrator_login` for the PostgreSQL Server.

* ~> **Note:** Either `administrator_login_password` or `administrator_login_password_wo` is required when `create_mode` is `Default`.

* `administrator_login_password_wo_version` - (Optional) An integer value used to trigger an update for `administrator_login_password_wo`. This property should be incremented when updating `administrator_login_password_wo`.

* `auto_grow_enabled` - (Optional) Enable/Disable auto-growing of the storage. Storage auto-grow prevents your server from running out of storage and becoming read-only. If storage auto grow is enabled, the storage automatically grows without impacting the workload. Defaults to `true`.

* `backup_retention_days` - (Optional) Backup retention days for the server, supported values are between `7` and `35` days.

* `create_mode` - (Optional) The creation mode. Can be used to restore or replicate existing servers. Possible values are `Default`, `Replica`, `GeoRestore`, and `PointInTimeRestore`. Defaults to `Default`.

* `creation_source_server_id` - (Optional) For creation modes other than `Default`, the source server ID to use.

* `geo_redundant_backup_enabled` - (Optional) Turn Geo-redundant server backups on/off. This allows you to choose between locally redundant or geo-redundant backup storage in the General Purpose and Memory Optimized tiers. When the backups are stored in geo-redundant backup storage, they are not only stored within the region in which your server is hosted, but are also replicated to a paired data center. This provides better protection and ability to restore your server in a different region in the event of a disaster. This is not support for the Basic tier. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `infrastructure_encryption_enabled` - (Optional) Whether or not infrastructure is encrypted for this server. Changing this forces a new resource to be created.

~> **Note:** This property is currently still in development and not supported by Microsoft. If the `infrastructure_encryption_enabled` attribute is set to `true` the PostgreSQL instance will incur a substantial performance degradation due to a second encryption pass on top of the existing default encryption that is already provided by Azure Storage. It is strongly suggested to leave this value `false` as not doing so can lead to unclear error messages.

* `public_network_access_enabled` - (Optional) Whether or not public network access is allowed for this server. Defaults to `true`.

* `restore_point_in_time` - (Optional) When `create_mode` is `PointInTimeRestore` the point in time to restore from `creation_source_server_id`. It should be provided in [RFC3339](https://www.rfc-editor.org/rfc/rfc3339) format, e.g. `2013-11-08T22:00:40Z`.

* `ssl_enforcement_enabled` - (Required) Specifies if SSL should be enforced on connections. Possible values are `true` and `false`.

-> **Note:** `ssl_minimal_tls_version_enforced` must be set to `TLSEnforcementDisabled` when `ssl_enforcement_enabled` is set to `false`.

* `ssl_minimal_tls_version_enforced` - (Optional) The minimum TLS version to support on the sever. Possible values are `TLSEnforcementDisabled`, `TLS1_0`, `TLS1_1`, and `TLS1_2`. Defaults to `TLS1_2`.

* `storage_mb` - (Optional) Max storage allowed for a server. Possible values are between `5120` MB(5GB) and `1048576` MB(1TB) for the Basic SKU and between `5120` MB(5GB) and `16777216` MB(16TB) for General Purpose/Memory Optimized SKUs. For more information see the [product documentation](https://docs.microsoft.com/azure/postgresql/concepts-pricing-tiers#storage).

* `threat_detection_policy` - (Optional) Threat detection policy configuration, known in the API as Server Security Alerts Policy. The `threat_detection_policy` block supports fields documented below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this PostgreSQL Server. The only possible value is `SystemAssigned`.

---

a `threat_detection_policy` block supports the following:

* `enabled` - (Optional) Is the policy enabled?

* `disabled_alerts` - (Optional) Specifies a list of alerts which should be disabled. Possible values are `Sql_Injection`, `Sql_Injection_Vulnerability`, `Access_Anomaly`, `Data_Exfiltration` and `Unsafe_Action`.

* `email_account_admins` - (Optional) Should the account administrators be emailed when this alert is triggered?

* `email_addresses` - (Optional) A list of email addresses which alerts should be sent to.

* `retention_days` - (Optional) Specifies the number of days to keep in the Threat Detection audit logs.

* `storage_account_access_key` - (Optional) Specifies the identifier key of the Threat Detection audit storage account.

* `storage_endpoint` - (Optional) Specifies the blob storage endpoint (e.g. <https://example.blob.core.windows.net>). This blob storage will hold all Threat Detection audit logs.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the PostgreSQL Server.

* `fqdn` - The FQDN of the PostgreSQL Server.

* `identity` - An `identity` block as documented below.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the PostgreSQL Server.
* `read` - (Defaults to 5 minutes) Used when retrieving the PostgreSQL Server.
* `update` - (Defaults to 1 hour) Used when updating the PostgreSQL Server.
* `delete` - (Defaults to 1 hour) Used when deleting the PostgreSQL Server.

## Import

PostgreSQL Server's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_postgresql_server.server1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DBforPostgreSQL/servers/server1
```
