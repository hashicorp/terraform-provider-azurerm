---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mysql_server"
description: |-
  Manages a MySQL Server.

---

# azurerm_mysql_server

Manages a MySQL Server.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_mysql_server" "example" {
  name                = "example-mysqlserver"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  administrator_login          = "mysqladminun"
  administrator_login_password = "H@Sh1CoR3!"

  sku_name   = "B_Gen5_2"
  storage_mb = 5120
  version    = "5.7"

  auto_grow_enabled                 = true
  backup_retention_days             = 7
  geo_redundant_backup_enabled      = false
  infrastructure_encryption_enabled = false
  public_network_access_enabled     = true
  ssl_enforcement_enabled           = true
  ssl_minimal_tls_version_enforced  = "TLS1_2"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the MySQL Server. Changing this forces a new resource to be created. This needs to be globally unique within Azure.

* `resource_group_name` - (Required) The name of the resource group in which to create the MySQL Server. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku_name` - (Required) Specifies the SKU Name for this MySQL Server. The name of the SKU, follows the `tier` + `family` + `cores` pattern (e.g. `B_Gen4_1`, `GP_Gen5_8`). For more information see the [product documentation](https://docs.microsoft.com/en-us/rest/api/mysql/servers/create#sku).

* `version` - (Required) Specifies the version of MySQL to use. Valid values are `5.6`, `5.7`, and `8.0`. Changing this forces a new resource to be created.

* `administrator_login` - (Optional) The Administrator Login for the MySQL Server. Required when `create_mode` is `Default`. Changing this forces a new resource to be created.

* `administrator_login_password` - (Optional) The Password associated with the `administrator_login` for the MySQL Server. Required when `create_mode` is `Default`.

* `auto_grow_enabled` - (Optional) Enable/Disable auto-growing of the storage. Storage auto-grow prevents your server from running out of storage and becoming read-only. If storage auto grow is enabled, the storage automatically grows without impacting the workload. The default value if not explicitly specified is `true`.

* `backup_retention_days` - (Optional) Backup retention days for the server, supported values are between `7` and `35` days.

* `create_mode` - (Optional) The creation mode. Can be used to restore or replicate existing servers. Possible values are `Default`, `Replica`, `GeoRestore`, and `PointInTimeRestore`. Defaults to `Default`.

* `creation_source_server_id` - (Optional) For creation modes other than `Default`, the source server ID to use.

* `geo_redundant_backup_enabled` - (Optional) Turn Geo-redundant server backups on/off. This allows you to choose between locally redundant or geo-redundant backup storage in the General Purpose and Memory Optimized tiers. When the backups are stored in geo-redundant backup storage, they are not only stored within the region in which your server is hosted, but are also replicated to a paired data center. This provides better protection and ability to restore your server in a different region in the event of a disaster. This is not supported for the Basic tier.

* `identity` - (Optional) An `identity` block as defined below. 

* `infrastructure_encryption_enabled` - (Optional) Whether or not infrastructure is encrypted for this server. Defaults to `false`. Changing this forces a new resource to be created.

* `public_network_access_enabled` - (Optional) Whether or not public network access is allowed for this server. Defaults to `true`.

* `restore_point_in_time` - (Optional) When `create_mode` is `PointInTimeRestore`, specifies the point in time to restore from `creation_source_server_id`.

* `ssl_enforcement_enabled` - (Required) Specifies if SSL should be enforced on connections. Possible values are `true` and `false`.

* `ssl_minimal_tls_version_enforced` - (Optional) The minimum TLS version to support on the sever. Possible values are `TLSEnforcementDisabled`, `TLS1_0`, `TLS1_1`, and `TLS1_2`. Defaults to `TLSEnforcementDisabled`.

* `storage_mb` - (Required) Max storage allowed for a server. Possible values are between `5120` MB(5GB) and `1048576` MB(1TB) for the Basic SKU and between `5120` MB(5GB) and `4194304` MB(4TB) for General Purpose/Memory Optimized SKUs. For more information see the [product documentation](https://docs.microsoft.com/en-us/rest/api/mysql/servers/create#StorageProfile).

* `threat_detection_policy` - (Optional) Threat detection policy configuration, known in the API as Server Security Alerts Policy. The `threat_detection_policy` block supports fields documented below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `identity` block supports the following:
   
* `type` - (Required) The Type of Identity which should be used for this MySQL Server. At this time the only possible value is `SystemAssigned`.

---

a `threat_detection_policy` block supports the following:

* `enabled` - (Required) Is the policy enabled?

* `disabled_alerts` - (Optional) Specifies a list of alerts which should be disabled. Possible values include `Access_Anomaly`, `Sql_Injection` and `Sql_Injection_Vulnerability`.

* `email_account_admins` - (Optional) Should the account administrators be emailed when this alert is triggered?

* `email_addresses` - (Optional) A list of email addresses which alerts should be sent to.

* `retention_days` - (Optional) Specifies the number of days to keep in the Threat Detection audit logs.

* `storage_account_access_key` - (Optional) Specifies the identifier key of the Threat Detection audit storage account.

* `storage_endpoint` - (Optional) Specifies the blob storage endpoint (e.g. https://MyAccount.blob.core.windows.net). This blob storage will hold all Threat Detection audit logs.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the MySQL Server.

* `fqdn` - The FQDN of the MySQL Server.

---

A `identity` block exports the following:

* `principal_id` - The Client ID of the Service Principal assigned to this MySQL Server.

* `tenant_id` - The ID of the Tenant the Service Principal is assigned in.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the MySQL Server.
* `update` - (Defaults to 60 minutes) Used when updating the MySQL Server.
* `read` - (Defaults to 5 minutes) Used when retrieving the MySQL Server.
* `delete` - (Defaults to 60 minutes) Used when deleting the MySQL Server.

## Import

MySQL Server's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mysql_server.server1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DBforMySQL/servers/server1
```
