---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_mysql_server"
description: |-
  Gets information about an existing MySQL Server.

---

# azurerm_mysql_server

Use this data source to access information about an existing MySQL Server.

~> **Note:** Azure Database for MySQL Single Server and its sub resources are scheduled for retirement by 2024-09-16 and will migrate to using Azure Database for MySQL Flexible Server: https://go.microsoft.com/fwlink/?linkid=2216041. The `azurerm_mysql_server` data source is deprecated and will be removed in v4.0 of the AzureRM Provider. Please use the `azurerm_mysql_flexible_server` data source instead.

## Example Usage

```hcl
data "azurerm_mysql_server" "example" {
  name                = "existingMySqlServer"
  resource_group_name = "existingResGroup"
}

output "id" {
  value = data.azurerm_mysql_server.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the MySQL Server.

* `resource_group_name` - (Required) The name of the resource group for the MySQL Server.

## Attributes Reference

* `id` - The ID of the MySQL Server.

* `fqdn` - The FQDN of the MySQL Server.

* `location` - The Azure location where the resource exists.

* `sku_name` - The SKU Name for this MySQL Server.

* `version` - The version of this MySQL Server.

* `administrator_login` - The Administrator login for the MySQL Server.

* `auto_grow_enabled` - The auto grow setting for this MySQL Server.

* `backup_retention_days` - The backup retention days for this MySQL server.

* `geo_redundant_backup_enabled` - The geo redundant backup setting for this MySQL Server.

* `identity` - An `identity` block as defined below.

* `infrastructure_encryption_enabled` - Whether or not infrastructure is encrypted for this MySQL Server.

* `public_network_access_enabled` - Whether or not public network access is allowed for this MySQL Server.

* `ssl_enforcement_enabled` -  Specifies if SSL should be enforced on connections for this MySQL Server.

* `ssl_minimal_tls_version_enforced` - The minimum TLS version to support for this MySQL Server.

* `storage_mb` -  Max storage allowed for this MySQL Server.

* `threat_detection_policy` - Threat detection policy configuration, known in the API as Server Security Alerts Policy. The `threat_detection_policy` block exports fields documented below.

* `tags` - A mapping of tags to assign to the resource.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

* `type` - The identity type of this Managed Service Identity.

---

A `threat_detection_policy` block exports the following:

* `enabled` -  Is the policy enabled?

* `disabled_alerts` - Specifies a list of alerts which should be disabled. Possible values include `Access_Anomaly`, `Sql_Injection` and `Sql_Injection_Vulnerability`.

* `email_account_admins` - Should the account administrators be emailed when this alert is triggered?

* `email_addresses` - A list of email addresses which alerts should be sent to.

* `retention_days` - Specifies the number of days to keep in the Threat Detection audit logs.

* `storage_account_access_key` - Specifies the identifier key of the Threat Detection audit storage account.

* `storage_endpoint` - Specifies the blob storage endpoint (e.g. <https://example.blob.core.windows.net>). This blob storage will hold all Threat Detection audit logs.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the MySQL Server.
