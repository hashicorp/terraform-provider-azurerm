---
subcategory: "AAD Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_aad_diagnostic_setting"
description: |-
  Manages Azure Active Directory diagnostic settings.
---

# azurerm_aad_diagnostic_setting

Manages Azure Active Directory diagnostic settings.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-law"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
}

resource "azurerm_storage_account" "example" {
  name                     = "aaddiagsettingsstgacc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_eventhub_namespace" "example" {
  name                = "aaddiagsettingseventhubns"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"
  capacity            = 1

  tags = {
    environment = "Production"
  }
}

resource "azurerm_eventhub" "example" {
  name                = "aaddiagsettingseventhub"
  namespace_name      = azurerm_eventhub_namespace.example.name
  resource_group_name = azurerm_resource_group.example.name
  partition_count     = 2
  message_retention   = 1
}


resource "azurerm_aad_diagnostic_setting" "example" {
  name                   = "aad-diag-test"
  storage_account_id     = azurerm_storage_account.example.id
  workspace_id           = azurerm_log_analytics_workspace.example.id
  event_hub_name         = azurerm_eventhub.example.name
  event_hub_auth_rule_id = "${azurerm_eventhub_namespace.example.id}/authorizationRules/RootManageSharedAccessKey"
  logs {
    category = "AuditLogs"
    enabled  = false
  }

  logs {
    enabled  = true
    category = "SignInLogs"
    retention_policy {
      retention_policy_days    = 20
      retention_policy_enabled = true
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of AAD Diagnostic settings. Changing this forces a new AAD Diagnostic settings to be created.

* `storage_account_id` - (Optional) The Storage account ID to which logs should be saved and archieved.

* `workspace_id` - (Optional) The Id of the Log Analytics Workspace to which the diagnostic logs should be sent to.

* `event_hub_name` - (Optional) The event hub name for sending logs to event hub namespace.

* `event_hub_auth_rule_id` - (Optional) The event hub authorization rule id to use for sending logs.

* `logs` - (Required) One or more `logs` block(s) as defined below.

---

A `logs` block supports the following:

* `category` - (Required) The log category to monitor. Possible values include 'SignInLogs', 'AuditLogs' 'ManagedIdentitySignInLogs', 'NonInteractiveUserSignInLogs', 'ProvisioningLogs' and 'ServicePrincipalSignInLogs'.

* `enabled` - (Optional) Whether to enable monitoring of the Log category. Defaults to true.

* `retention_policy` - (Optional) One block of Retention policy for arhieval. Retention policy block is defined as below.

---

A `retention_policy` block supports the following:

* `retention_policy_days` - (Optional) The number of days to archieve logs for. Possible values range from 0 to 365(inclusive).

* `retention_policy_enabled` - (Optional) Whether to enable retention policy settings for the log category.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 15 minutes) Used when creating the AAD diagnostic settings resource.
* `read` - (Defaults to 5 minutes) Used when retrieving the AAD diagnostic settings resource.
* `update` - (Defaults to 15 minutes) Used when updating the AAD diagnostic settings resource.
* `delete` - (Defaults to 5 minutes) Used when deleting the AAD diagnostic settings resource.

## Import

Azure AD Diagnostic settings resources can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_aad_diagnostic_setting.example /providers/microsoft.aadiam/diagnosticSettings/aad-diag-test
```
