---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_log_profile"
description: |-
  Manages a Log Profile.
---

# azurerm_monitor_log_profile

Manages a [Log Profile](https://docs.microsoft.com/en-us/azure/monitoring-and-diagnostics/monitoring-overview-activity-logs#export-the-activity-log-with-a-log-profile). A Log Profile configures how Activity Logs are exported.

-> **NOTE:** It's only possible to configure one Log Profile per Subscription. If you are trying to create more than one Log Profile, an error with `StatusCode=409` will occur.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "logprofiletest-rg"
  location = "eastus"
}

resource "azurerm_storage_account" "example" {
  name                     = "afscsdfytw"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_eventhub_namespace" "example" {
  name                = "logprofileeventhub"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"
  capacity            = 2
}

resource "azurerm_monitor_log_profile" "example" {
  name     = "default"
  location = azurerm_resource_group.example.location

  categories = [
    "Action",
    "Delete",
    "Write",
  ]

  locations = [
    "westus",
    "global",
  ]

  # RootManageSharedAccessKey is created by default with listen, send, manage permissions
  servicebus_rule_id = "${azurerm_eventhub_namespace.example.id}/authorizationrules/RootManageSharedAccessKey"
  storage_account_id = azurerm_storage_account.example.id

  retention_policy {
    enabled = true
    days    = 7
  }

  tags = {
    "foo" = "bar"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Log Profile. Changing this forces a
    new resource to be created.

* `categories` - (Required) List of categories of the logs.

* `locations` - (Required) List of regions for which Activity Log events are stored or streamed.

* `storage_account_id` - (Optional) The resource ID of the storage account in which the Activity Log is stored. At least one of `storage_account_id` or `servicebus_rule_id` must be set.

* `servicebus_rule_id` - (Optional) The service bus (or event hub) rule ID of the service bus (or event hub) namespace in which the Activity Log is streamed to. At least one of `storage_account_id` or `servicebus_rule_id` must be set.

* `retention_policy` - (Required) A `retention_policy` block as documented below. A retention policy for how long Activity Logs are retained in the storage account.

* `location` - (Optional) The location of the Log Profile.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `retention_policy` block supports:

* `enabled` - (Required) A boolean value to indicate whether the retention policy is enabled.

* `days` - (Optional) The number of days for the retention policy. Defaults to 0.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Log Profile.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Log Profile.
* `update` - (Defaults to 30 minutes) Used when updating the Log Profile.
* `read` - (Defaults to 5 minutes) Used when retrieving the Log Profile.
* `delete` - (Defaults to 30 minutes) Used when deleting the Log Profile.

## Import

A Log Profile can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_log_profile.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/microsoft.insights/logprofiles/test
```
