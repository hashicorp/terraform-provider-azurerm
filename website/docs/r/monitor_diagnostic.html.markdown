---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_diagnostics"
sidebar_current: "docs-azurerm-resource-monitor-diagnostics"
description: |-
  Manages Diagnostic Logging settings.

---

# azurerm_metric_alertrule

Use this resource to manage [Monitor Diagnostic](https://docs.microsoft.com/en-us/azure/monitoring-and-diagnostics/monitoring-overview-of-diagnostic-logs) settings.

## Example Usage (Enable logging on a Key Vault)

```
resource "azurerm_key_vault" "vault" {
  name                = "vault"
  location            = "eastus2"
  resource_group_name = "resource_group"
  tenant_id           = "xxx"

  sku {
    name = "premium"
  }
}

resource "azurerm_storage_account" "logging_storage" {
  name = "loggingstorage"
  resource_group_name = "resource_group"
  location = "${azurerm_resource_group.test_rg.location}"
  account_replication_type = "LRS"
  account_tier = "Standard"
}
  
resource "azurerm_monitor_diagnostics" "diagnostic_logging" {
  name = "diagnostic_logging"
  resource_id = "${azurerm_key_vault.vault.id}"
  storage_account_id = "${azurerm_storage_account.logging_storage.id}"
  disabled_settings = ["AuditEvent"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Diagnostic Settings. Changing this will force a new resource.

* `resource_id` - (Required) Provide the ID of the resource you want to configure Diagnostic Settings for. Must be an existing Azure resource. Changing this will force a new resource.

* `storage_account_id` - (Optional) With this parameter you can specify a storage account which should be used to send the logs to. Parameter must be a valid Azure Resource ID.

* `event_hub_name` - (Optional) Specify the name of the event hub used to log to.

* `event_hub_authorization_rule_id` - (Optional)  Specify the id of the rule used for authorization with event hub.

* `workspace_id` - (Optional) Specify the Azure resource ID of a Log Analytics workspace to send the logs to.

* `disabled_settings` - (Optional) Specify type of logging you do not want. If you omit this setting, all Metrics and Log Settings will be enabled.

* `retention_days` - (Optional) With this parameter you can specify the number of days to keep these Metrics/Logs. If you leave this setting out or provide a `0` logs will not be deleted. This setting only applies to storage accounts and is ignored for other log targets!


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the diagnostic setting.

## Import

Diagnostic Settings can be imported using the `resource id`, e.g.

```
terraform import azurerm_monitor_diagnostics.diagnostic_logging /subscriptions/XXX/resourcegroups/resource_group/providers/microsoft.keyvault/vaults/vault/providers/microsoft.insights/diagnosticSettings/diagnostic_logging
```
