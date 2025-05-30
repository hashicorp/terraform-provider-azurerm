---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_account_queue_properties"
description: |-
  Manages the Queue Properties of an Azure Storage Account.
---

# azurerm_storage_account_queue_properties

Manages the Queue Properties of an Azure Storage Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "storageaccountname"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_account_queue_properties" "example" {
  storage_account_id = azurerm_storage_account.example.id
  cors_rule {
    allowed_origins    = ["http://www.example.com"]
    exposed_headers    = ["x-tempo-*"]
    allowed_headers    = ["x-tempo-*"]
    allowed_methods    = ["GET", "PUT"]
    max_age_in_seconds = "500"
  }

  logging {
    version               = "1.0"
    delete                = true
    read                  = true
    write                 = true
    retention_policy_days = 7
  }

  hour_metrics {
    version               = "1.0"
    retention_policy_days = 7
  }

  minute_metrics {
    version               = "1.0"
    retention_policy_days = 7
  }
}
```

## Argument Reference

The following arguments are supported:

* `storage_account_id` - (Required) The ID of the Storage Account to set Queue Properties on. Changing this forces a new resource to be created.

* `cors_rule` - (Optional) A `cors_rule` block as defined above.

* `logging` - (Optional) A `logging` block as defined below.

* `minute_metrics` - (Optional) A `minute_metrics` block as defined below.

* `hour_metrics` - (Optional) A `hour_metrics` block as defined below.

~> **Note:** At least one of `cors_rule`, `logging`, `minute_metrics`, or `hour_metrics` must be specified.

---

A `cors_rule` block supports the following:

* `allowed_headers` - (Required) A list of headers that are allowed to be a part of the cross-origin request.

* `allowed_methods` - (Required) A list of HTTP methods that are allowed to be executed by the origin. Valid options are
  `DELETE`, `GET`, `HEAD`, `MERGE`, `POST`, `OPTIONS`, `PUT` or `PATCH`.

* `allowed_origins` - (Required) A list of origin domains that will be allowed by CORS.

* `exposed_headers` - (Required) A list of response headers that are exposed to CORS clients.

* `max_age_in_seconds` - (Required) The number of seconds the client should cache a preflight response.

---

An `hour_metrics` block supports the following:

* `version` - (Required) The version of storage analytics to configure.

* `include_apis` - (Optional) Indicates whether metrics should generate summary statistics for called API operations.

* `retention_policy_days` - (Optional) Specifies the number of days that logs will be retained.

---

A `logging` block supports the following:

* `delete` - (Required) Indicates whether all delete requests should be logged.

* `read` - (Required) Indicates whether all read requests should be logged.

* `version` - (Required) The version of storage analytics to configure.

* `write` - (Required) Indicates whether all write requests should be logged.

* `retention_policy_days` - (Optional) Specifies the number of days that logs will be retained.

---

A `minute_metrics` block supports the following:

* `version` - (Required) The version of storage analytics to configure.

* `include_apis` - (Optional) Indicates whether metrics should generate summary statistics for called API operations.

* `retention_policy_days` - (Optional) Specifies the number of days that logs will be retained.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Account Queue Properties.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Account Queue Properties.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Account Queue Properties.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Account Queue Properties.

## Import

Storage Account Queue Properties can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_account_queue_properties.queueprops /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Storage/storageAccounts/myaccount
```
