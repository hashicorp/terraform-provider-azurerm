---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_server_security_alert_policy"
description: |-
  Manages a Security Alert Policy for a MS SQL Server.

---

# azurerm_mssql_server_security_alert_policy

Manages a Security Alert Policy for a MSSQL Server.

-> **Note:** Security Alert Policy is currently only available for MS SQL databases.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
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
  retention_days             = 20

  disabled_alerts = [
    "Sql_Injection",
    "Data_Exfiltration"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group that contains the MS SQL Server. Changing this forces a new resource to be created.

* `server_name` - (Required) Specifies the name of the MS SQL Server. Changing this forces a new resource to be created.

* `state` - (Required) Specifies the state of the policy. Possible values are `Disabled` or `Enabled`.

* `disabled_alerts` - (Optional) Specifies an array of alerts that are disabled. Allowed values are: `Sql_Injection`, `Sql_Injection_Vulnerability`, `Access_Anomaly`, `Data_Exfiltration`, `Unsafe_Action`.

* `email_account_admins` - (Optional) Are the alerts sent to the account administrators? Possible values are `true` or `false`. Defaults to `false`.

* `email_addresses` - (Optional) Specifies an array of email addresses to which the alert is sent.

* `retention_days` - (Optional) Specifies the number of days to keep the Threat Detection audit logs. Defaults to `0`.

* `storage_endpoint` - (Optional) Specifies the blob storage endpoint that will hold all Threat Detection audit logs (e.g., `https://example.blob.core.windows.net`).

-> **Note:** The `storage_account_access_key` field is required when the `storage_endpoint` field has been set.

-> **Note:** Storage accounts configured with `shared_access_key_enabled = false` cannot be used for the `storage_endpoint` field.

* `storage_account_access_key` - (Optional) Specifies the primary access key of the Threat Detection audit logs blob storage endpoint.

-> **Note:** The `storage_account_access_key` only applies if the storage account is not behind a virtual network or a firewall.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the MS SQL Server Security Alert Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the MSSQL Server Security Alert Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the MSSQL Server Security Alert Policy.
* `update` - (Defaults to 30 minutes) Used when updating the MSSQL Server Security Alert Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the MSSQL Server Security Alert Policy.

## Import

MS SQL Server Security Alert Policy can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_server_security_alert_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acceptanceTestResourceGroup1/providers/Microsoft.Sql/servers/mssqlserver/securityAlertPolicies/Default
```
