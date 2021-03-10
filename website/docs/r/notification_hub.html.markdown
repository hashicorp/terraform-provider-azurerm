---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_notification_hub"
description: |-
  Manages a Notification Hub within a Notification Hub Namespace.

---

# azurerm_notification_hub

Manages a Notification Hub within a Notification Hub Namespace.

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

  sku_name = "Free"
}

resource "azurerm_notification_hub" "example" {
  name                = "mynotificationhub"
  namespace_name      = azurerm_notification_hub_namespace.example.name
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name to use for this Notification Hub. Changing this forces a new resource to be created.

* `namespace_name` - (Required) The name of the Notification Hub Namespace in which to create this Notification Hub. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Notification Hub Namespace exists. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region in which this Notification Hub Namespace exists. Changing this forces a new resource to be created.

* `apns_credential` - (Optional) A `apns_credential` block as defined below.

~> **NOTE:** Removing the `apns_credential` block will currently force a recreation of this resource [due to this bug in the Azure SDK for Go](https://github.com/Azure/azure-sdk-for-go/issues/2246) - we'll remove this limitation when the SDK bug is fixed.

* `gcm_credential` - (Optional) A `gcm_credential` block as defined below.

~> **NOTE:** Removing the `gcm_credential` block will currently force a recreation of this resource [due to this bug in the Azure SDK for Go](https://github.com/Azure/azure-sdk-for-go/issues/2246) - we'll remove this limitation when the SDK bug is fixed.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `apns_credential` block contains:

* `application_mode` - (Required) The Application Mode which defines which server the APNS Messages should be sent to. Possible values are `Production` and `Sandbox`.

* `bundle_id` - (Required) The Bundle ID of the iOS/macOS application to send push notifications for, such as `com.hashicorp.example`.

* `key_id` - (Required) The Apple Push Notifications Service (APNS) Key.

* `team_id` - (Required) The ID of the team the Token.

* `token` - (Required) The Push Token associated with the Apple Developer Account. This is the contents of the `key` downloaded from [the Apple Developer Portal](https://developer.apple.com/account/ios/authkey/) between the `-----BEGIN PRIVATE KEY-----` and `-----END PRIVATE KEY-----` blocks.

---

A `gcm_credential` block contains:

* `api_key` - (Required) The API Key associated with the Google Cloud Messaging service.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Notification Hub.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Notification Hub.
* `update` - (Defaults to 30 minutes) Used when updating the Notification Hub.
* `read` - (Defaults to 5 minutes) Used when retrieving the Notification Hub.
* `delete` - (Defaults to 30 minutes) Used when deleting the Notification Hub.

## Import

Notification Hubs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_notification_hub.hub1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.NotificationHubs/namespaces/{namespaceName}/notificationHubs/hub1
```
