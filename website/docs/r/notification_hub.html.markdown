---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_notification_hub"
sidebar_current: "docs-azurerm-resource-messaging-notification-hub-x"
description: |-
  Manages a Notification Hub within a Notification Hub Namespace.

---

# azurerm_notification_hub

Manages a Notification Hub within a Notification Hub Namespace.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name = "notificationhub-resources"
  location = "Australia East"
}

resource "azurerm_notification_hub_namespace" "test" {
  name                = "myappnamespace"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  namespace_type      = "NotificationHub"

  sku {
    name = "Free"
  }
}

resource "azurerm_notification_hub" "test" {
  name                = "mynotificationhub"
  namespace_name      = "${azurerm_notification_hub_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name to use for this Notification Hub. Changing this forces a new resource to be created.

* `namespace_name` - (Required) The name of the Notification Hub Namespace in which to create this Notification Hub. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Notification Hub Namespace exists. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region in which this Notification Hub Namespace exists. Changing this forces a new resource to be created.

* `apns_credential` - (Optional) A `apns_credential` block as defined below.

* `gcm_credential` - (Optional) A `gcm_credential` block as defined below.

---

A `apns_credential` block contains:

* `application_id` - (Required) The reverse-domain identifier for the application, such as `com.hashicorp.example`.

* `application_mode` - (Required) The Application Mode which defines which server the APNS Messages should be sent to. Possible values are `Production` and `Sandbox`.

* `application_name` - (Required) The name of the Application.

* `key_id` - (Required) The Apple Push Notifications Service (APNS) Key.

* `token` - (Required) The Token associated with the Apple.

---

A `gcm_credential` block contains:

* `api_key` - (Required) The API Key associated with the Google Cloud Messaging service.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Notification Hub.

## Import

Notification Hubs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_notification_hub.hub1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.NotificationHubs/namespaces/{namespaceName}/notificationHubs/hub1
```
