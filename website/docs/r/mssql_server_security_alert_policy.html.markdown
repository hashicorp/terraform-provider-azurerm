---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_server_security_alert_policy"
sidebar_current: "docs-azurerm-resource-database-mssql-server-security-alert-policy-x"
description: |-
  Manages a Security Alert Policy for a MS SQL Server.

---

# azurerm_mssql_server_security_alert_policy

Manages a Security Alert Policy for a MSSQL Server.

-> **NOTE** Security Alert Policy is currently only available for MS SQL databases.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "acceptanceTestResourceGroup1"
  location = "West US"
}

resource "azurerm_sql_server" "example" {
  name                         = "mysqlserver"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_storage_account" "example" {
  name                     = "accteststorageaccount"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_mssql_server_security_alert_policy" "example" {
  resource_group_name        = azurerm_resource_group.example.name
  server_name                = azurerm_sql_server.example.name
  state                      = "Enabled"
  storage_endpoint           = azurerm_storage_account.example.primary_blob_endpoint
  storage_account_access_key = azurerm_storage_account.example.primary_access_key
  disabled_alerts = [
    "Sql_Injection",
    "Data_Exfiltration"
  ]
  retention_days = 20
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group that contains the MS SQL Server. Changing this forces a new resource to be created.

* `server_name` - (Required) Specifies the name of the MS SQL Server. Changing this forces a new resource to be created.

* `state` - (Required) Specifies the state of the policy, whether it is enabled or disabled or a policy has not been applied yet on the specific database server. Allowed values are: `Disabled`, `Enabled`, `New`.

* `disabled_alerts` - (Optional) Specifies an array of alerts that are disabled. Allowed values are: `Sql_Injection`, `Sql_Injection_Vulnerability`, `Access_Anomaly`, `Data_Exfiltration`, `Unsafe_Action`.

* `email_account_admins` - (Optional) Boolean flag which specifies if the alert is sent to the account administrators or not. Defaults to `false`.
    
* `email_addresses` - (Optional) Specifies an array of e-mail addresses to which the alert is sent.

* `retention_days` - (Optional) Specifies the number of days to keep in the Threat Detection audit logs. Defaults to `0`.

* `storage_account_access_key` - (Optional) Specifies the identifier key of the Threat Detection audit storage account.

* `storage_endpoint` - (Optional) Specifies the blob storage endpoint (e.g. https://MyAccount.blob.core.windows.net). This blob storage will hold all Threat Detection audit logs.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the MS SQL Server Security Alert Policy.

## Import

MS SQL Server Security Alert Policy can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_server_security_alert_policy.example  /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acceptanceTestResourceGroup1/providers/Microsoft.Sql/servers/mssqlserver/securityAlertPolicies/Default 
```
