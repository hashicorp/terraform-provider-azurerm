---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_notification_recipient_user"
description: |-
  Manages a API Management Notification Recipient User.
---

# azurerm_api_management_notification_recipient_user

Manages a API Management Notification Recipient User.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_api_management" "example" {
  name                = "example-apim"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  publisher_name      = "My Company"
  publisher_email     = "company@terraform.io"

  sku_name = "Developer_1"
}

resource "azurerm_api_management_user" "example" {
  user_id             = "123"
  api_management_name = azurerm_api_management.example.name
  resource_group_name = azurerm_resource_group.example.name
  first_name          = "Example"
  last_name           = "User"
  email               = "foo@bar.com"
  state               = "active"
}

resource "azurerm_api_management_notification_recipient_user" "example" {
  api_management_id = azurerm_api_management.example.id
  notification_type = "AccountClosedPublisher"
  user_id           = azurerm_api_management_user.example.user_id
}
```

## Arguments Reference

The following arguments are supported:

* `api_management_id` - (Required) The ID of the API Management Service from which to create this Notification Recipient User. Changing this forces a new API Management Notification Recipient User to be created.

* `user_id` - (Required) The recipient user ID. Changing this forces a new API Management Notification Recipient User to be created.

* `notification_type` - (Required) The Notification Name to be received. Changing this forces a new API Management Notification Recipient User to be created. Possible values are `AccountClosedPublisher`, `BCC`, `NewApplicationNotificationMessage`, `NewIssuePublisherNotificationMessage`, `PurchasePublisherNotificationMessage`, `QuotaLimitApproachingPublisherNotificationMessage`, and `RequestPublisherNotificationMessage`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management Notification Recipient User.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Notification Recipient User.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Notification Recipient User.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Notification Recipient User.

## Import

API Management Notification Recipient Users can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_notification_recipient_user.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/notifications/notificationName1/recipientUsers/userid1
```
