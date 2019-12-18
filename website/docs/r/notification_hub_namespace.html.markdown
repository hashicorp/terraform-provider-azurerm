---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_notification_hub_namespace"
sidebar_current: "docs-azurerm-resource-messaging-notification-hub-namespace"
description: |-
  Manages a Notification Hub Namespace.

---

# azurerm_notification_hub_namespace

Manages a Notification Hub Namespace.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "notificationhub-resources"
  location = "Australia East"
}

resource "azurerm_notification_hub_namespace" "example" {
  name                = "myappnamespace"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
  namespace_type      = "NotificationHub"

  sku_name = "Free"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name to use for this Notification Hub Namespace. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Notification Hub Namespace should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region in which this Notification Hub Namespace should be created.

* `namespace_type` - (Required) The Type of Namespace - possible values are `Messaging` or `NotificationHub`. Changing this forces a new resource to be created.

* `sku` - (Optional **Deprecated**)) A `sku` block as described below.

* `sku_name` - (Optional) The name of the SKU to use for this Notification Hub Namespace. Possible values are `Free`, `Basic` or `Standard`. Changing this forces a new resource to be created.

* `enabled` - (Optional) Is this Notification Hub Namespace enabled? Defaults to `true`.

----

A `sku` block contains:

* `name` - (Required) The name of the SKU to use for this Notification Hub Namespace. Possible values are `Free`, `Basic` or `Standard`. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Notification Hub Namespace.

* `servicebus_endpoint` - The ServiceBus Endpoint for this Notification Hub Namespace.

## Import

Notification Hub Namespaces can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_notification_hub_namespace.namespace1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.NotificationHubs/namespaces/{namespaceName}
```
