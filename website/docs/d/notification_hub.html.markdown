---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_notification_hub"
description: |-
  Gets information about an existing Notification Hub within a Notification Hub Namespace.
---

# Data Source: azurerm_notification_hub

Use this data source to access information about an existing Notification Hub within a Notification Hub Namespace.

## Example Usage

```hcl
data "azurerm_notification_hub" "example" {
  name                = "notification-hub"
  namespace_name      = "namespace-name"
  resource_group_name = "resource-group-name"
}

output "id" {
  value = data.azurerm_notification_hub.example.id
}
```

## Argument Reference

* `name` - Specifies the Name of the Notification Hub.

* `namespace_name` -  Specifies the Name of the Notification Hub Namespace which contains the Notification Hub.

* `resource_group_name` - Specifies the Name of the Resource Group within which the Notification Hub exists.

## Attributes Reference

* `id` - The ID of the Notification Hub.

* `location` - The Azure Region in which this Notification Hub exists.

* `apns_credential` - A `apns_credential` block as defined below.

* `gcm_credential` - A `gcm_credential` block as defined below.

* `tags` - A mapping of tags to assign to the resource.

---

A `apns_credential` block exports:

* `application_mode` - The Application Mode which defines which server the APNS Messages should be sent to. Possible values are `Production` and `Sandbox`.

* `bundle_id` - The Bundle ID of the iOS/macOS application to send push notifications for, such as `com.hashicorp.example`.

* `key_id` - The Apple Push Notifications Service (APNS) Key.

* `team_id` - The ID of the team the Token.

* `token` - The Push Token associated with the Apple Developer Account.

---

A `gcm_credential` block exports:

* `api_key` - The API Key associated with the Google Cloud Messaging service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Notification Hub within a Notification Hub Namespace.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.NotificationHubs`: 2023-09-01
