---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_notification_hub_namespace"
description: |-
  Manages a Notification Hub Namespace.

---

# azurerm_notification_hub_namespace

Manages a Notification Hub Namespace.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "notificationhub-resources"
  location = "West Europe"
}

resource "azurerm_notification_hub_namespace" "example" {
  name                = "myappnamespace"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  namespace_type      = "NotificationHub"
  sku_name            = "Free"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name to use for this Notification Hub Namespace. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Notification Hub Namespace should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region in which this Notification Hub Namespace should be created. Changing this forces a new resource to be created.

* `namespace_type` - (Required) The Type of Namespace - possible values are `Messaging` or `NotificationHub`. 

* `sku_name` - (Required) The name of the SKU to use for this Notification Hub Namespace. Possible values are `Free`, `Basic` or `Standard`. 

* `enabled` - (Optional) Is this Notification Hub Namespace enabled? Defaults to `true`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Notification Hub Namespace.

* `servicebus_endpoint` - The ServiceBus Endpoint for this Notification Hub Namespace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Notification Hub Namespace.
* `read` - (Defaults to 5 minutes) Used when retrieving the Notification Hub Namespace.
* `update` - (Defaults to 30 minutes) Used when updating the Notification Hub Namespace.
* `delete` - (Defaults to 30 minutes) Used when deleting the Notification Hub Namespace.

## Import

Notification Hub Namespaces can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_notification_hub_namespace.namespace1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.NotificationHubs/namespaces/namespace1
```
