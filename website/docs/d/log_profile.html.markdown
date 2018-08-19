---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_profile"
sidebar_current: "docs-azurerm-datasource-log-profile"
description: |-
  Get information about the specified Log Profile.
---

# Data Source: azurerm_log_profile

Use this data source to access the properties of a Log Profile.

## Example Usage

```hcl
data "azurerm_log_profile" "test" {
  name = "test-logprofile"
}

output "log_profile_storage_account_id" {
  value = "${data.azurerm_log_profile.test.storage_account_id}"
}
```

## Argument Reference

* `name` - (Required) Specifies the Name of the Log Profile.


## Attributes Reference

* `id` - The ID of the Log Profile.
 
* `storage_account_id` - The resource id of the storage account in which the Activity Log is stored.

* `service_bus_rule_id` - The service bus (or event hub) rule ID of the service bus (or event hub) namespace in which the Activity Log is streamed to.

* `locations` - List of regions for which Activity Log events are stored or streamed.

* `categories` - List of categories of the logs.

* `retention_policy`- A retention policy for how long Activity Logs are retained in the storage account.

---

The `retention_policy` block supports:

* `enabled` - A boolean value indicating whether the retention policy is enabled.

* `days` - The number of days for the retention policy.
