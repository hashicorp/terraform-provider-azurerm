---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_aad_diagnostic_setting"
description: |-
  Manages an Azure Active Directory Diagnostic Setting for Azure Monitor.
---

# azurerm_monitor_aad_diagnostic_setting

Manages an Azure Active Directory Diagnostic Setting for Azure Monitor.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "west europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestorageaccount"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_kind             = "StorageV2"
  account_replication_type = "LRS"
}

resource "azurerm_monitor_aad_diagnostic_setting" "example" {
  name               = "setting1"
  storage_account_id = azurerm_storage_account.example.id
  log {
    category = "SignInLogs"
    enabled  = true
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  log {
    category = "AuditLogs"
    enabled  = true
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  log {
    category = "NonInteractiveUserSignInLogs"
    enabled  = true
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  log {
    category = "ServicePrincipalSignInLogs"
    enabled  = true
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  log {
    category = "ManagedIdentitySignInLogs"
    enabled  = false
    retention_policy {}
  }
  log {
    category = "ProvisioningLogs"
    enabled  = false
    retention_policy {}
  }
  log {
    category = "ADFSSignInLogs"
    enabled  = false
    retention_policy {}
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Monitor Azure Active Directory Diagnostic Setting. Changing this forces a new Monitor Azure Active Directory Diagnostic Setting to be created.
  
* `log` - (Required) One or more `log` blocks as defined below.

~> **Note:** At least one of the `log` blocks must have the `enabled` property set to `true`.

---

* `eventhub_authorization_rule_id` - (Optional) Specifies the ID of an Event Hub Namespace Authorization Rule used to send Diagnostics Data. Changing this forces a new resource to be created.

-> **NOTE:** This can be sourced from [the `azurerm_eventhub_namespace_authorization_rule` resource](eventhub_namespace_authorization_rule.html) and is different from [a `azurerm_eventhub_authorization_rule` resource](eventhub_authorization_rule.html).

* `eventhub_name` - (Optional) Specifies the name of the Event Hub where Diagnostics Data should be sent. If not specified, the default Event Hub will be used. Changing this forces a new resource to be created.

* `log_analytics_workspace_id` - (Optional) Specifies the ID of a Log Analytics Workspace where Diagnostics Data should be sent.

* `storage_account_id` - (Optional) The ID of the Storage Account where logs should be sent. Changing this forces a new resource to be created.

-> **NOTE:** One of `eventhub_authorization_rule_id`, `log_analytics_workspace_id` and `storage_account_id` must be specified.

---

A `log` block supports the following:

* `category` - (Required) The log category for the Azure Active Directory Diagnostic. Possible values are `AuditLogs`, `SignInLogs`, `ADFSSignInLogs`, `ManagedIdentitySignInLogs`, `NonInteractiveUserSignInLogs`, `ProvisioningLogs`, `ServicePrincipalSignInLogs`.

* `retention_policy` - (Required) A `retention_policy` block as defined below.

* `enabled` - (Optional) Is this Diagnostic Log enabled? Defaults to `true`.

---

A `retention_policy` block supports the following:

* `enabled` - (Optional) Is this Retention Policy enabled? Defaults to `false`.

* `days` - (Optional) The number of days for which this Retention Policy should apply. Defaults to `0`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Monitor Azure Active Directory Diagnostic Setting.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Monitor Azure Active Directory Diagnostic Setting.
* `read` - (Defaults to 5 minutes) Used when retrieving the Monitor Azure Active Directory Diagnostic Setting.
* `update` - (Defaults to 5 minutes) Used when updating the Monitor Azure Active Directory Diagnostic Setting.
* `delete` - (Defaults to 5 minutes) Used when deleting the Monitor Azure Active Directory Diagnostic Setting.

## Import

Monitor Azure Active Directory Diagnostic Settings can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_aad_diagnostic_setting.example /providers/Microsoft.AADIAM/diagnosticSettings/setting1
```
