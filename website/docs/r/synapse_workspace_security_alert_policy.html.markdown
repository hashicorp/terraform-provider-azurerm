---
subcategory: "Synapse"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_synapse_workspace_security_alert_policy"
description: |-
  Manages a Security Alert Policy for a Synapse Workspace.

---

# azurerm_synapse_workspace_security_alert_policy

Manages a Security Alert Policy for a Synapse Workspace.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestorageacc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "StorageV2"
  is_hns_enabled           = "true"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "example" {
  name               = "example"
  storage_account_id = azurerm_storage_account.example.id
}

resource "azurerm_synapse_workspace" "example" {
  name                                 = "example"
  resource_group_name                  = azurerm_resource_group.example.name
  location                             = azurerm_resource_group.example.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.example.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"

  aad_admin {
    login     = "AzureAD Admin"
    object_id = "00000000-0000-0000-0000-000000000000"
    tenant_id = "00000000-0000-0000-0000-000000000000"
  }

  identity {
    type = "SystemAssigned"
  }

  tags = {
    Env = "production"
  }
}

resource "azurerm_storage_account" "audit_logs" {
  name                     = "examplesa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_synapse_workspace_security_alert_policy" "example" {
  synapse_workspace_id       = azurerm_synapse_workspace.example.id
  policy_state               = "Enabled"
  storage_endpoint           = azurerm_storage_account.audit_logs.primary_blob_endpoint
  storage_account_access_key = azurerm_storage_account.audit_logs.primary_access_key
  disabled_alerts = [
    "Sql_Injection",
    "Data_Exfiltration"
  ]
  retention_days = 20
}
```

## Argument Reference

The following arguments are supported:

* `synapse_workspace_id` - (Required) Specifies the ID of the Synapse Workspace. Changing this forces a new resource to be created.

* `policy_state` - (Required) Specifies the state of the policy, whether it is enabled or disabled or a policy has not been applied yet on the specific workspace. Possible values are `Disabled`, `Enabled` and `New`.

* `disabled_alerts` - (Optional) Specifies an array of alerts that are disabled. Allowed values are: `Sql_Injection`, `Sql_Injection_Vulnerability`, `Access_Anomaly`, `Data_Exfiltration`, `Unsafe_Action`.

* `email_account_admins_enabled` - (Optional) Boolean flag which specifies if the alert is sent to the account administrators or not. Defaults to `false`.

* `email_addresses` - (Optional) Specifies an array of email addresses to which the alert is sent.

* `retention_days` - (Optional) Specifies the number of days to keep in the Threat Detection audit logs. Defaults to `0`.

* `storage_account_access_key` - (Optional) Specifies the identifier key of the Threat Detection audit storage account.

* `storage_endpoint` - (Optional) Specifies the blob storage endpoint (e.g. <https://example.blob.core.windows.net>). This blob storage will hold all Threat Detection audit logs.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Synapse Workspace Security Alert Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Synapse Workspace Security Alert Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Synapse Workspace Security Alert Policy.
* `update` - (Defaults to 30 minutes) Used when updating the Synapse Workspace Security Alert Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Synapse Workspace Security Alert Policy.

## Import

Synapse Workspace Security Alert Policies can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_synapse_workspace_security_alert_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Synapse/workspaces/workspace1/securityAlertPolicies/Default
```
