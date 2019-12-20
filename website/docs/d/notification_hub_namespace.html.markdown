---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_notification_hub_namespace"
sidebar_current: "docs-azurerm-datasource-notification-hub-namespace"
description: |-
  Gets information about an existing Notification Hub Namespace.
---

# Data Source: azurerm_notification_hub_namespace

Use this data source to access information about an existing Notification Hub Namespace.

## Example Usage

```hcl
data "azurerm_notification_hub_namespace" "example" {
  name                = "my-namespace"
  resource_group_name = "my-resource-group"
}

output "servicebus_endpoint" {
  value = "${data.azurerm_notification_hub_namespace.example.servicebus_endpoint}"
}
```

## Argument Reference

* `name` - (Required) Specifies the Name of the Notification Hub Namespace.

* `resource_group_name` - (Required) Specifies the Name of the Resource Group within which the Notification Hub exists.

## Attributes Reference

* `id` - The ID of the Notification Hub Namespace.

* `location` - The Azure Region in which this Notification Hub Namespace exists.

* `namespace_type` - The Type of Namespace, such as `Messaging` or `NotificationHub`.

* `sku` - A `sku` block as defined below.

* `enabled` - Is this Notification Hub Namespace enabled?

---

A `sku` block exports the following:

* `name` - (Required) The name of the SKU to use for this Notification Hub Namespace. Possible values are `Free`, `Basic` or `Standard.`
