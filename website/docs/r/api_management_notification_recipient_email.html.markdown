---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_notification_recipient_email"
description: |-
  Manages a API Management Notification Recipient Email.
---

# azurerm_api_management_notification_recipient_email

Manages a API Management Notification Recipient Email.

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

resource "azurerm_api_management_notification_recipient_email" "example" {
  api_management_id = azurerm_api_management.example.id
  notification_type = "AccountClosedPublisher"
  email             = "foo@bar.com"
}
```

## Arguments Reference

The following arguments are supported:

* `api_management_id` - (Required) The ID of the API Management Service from which to create this Notification Recipient Email. Changing this forces a new API Management Notification Recipient Email to be created.

* `email` - (Required) The recipient email address. Changing this forces a new API Management Notification Recipient Email to be created.

* `notification_type` - (Required) The Notification Name to be received. Changing this forces a new API Management Notification Recipient Email to be created. Possible values are `AccountClosedPublisher`, `BCC`, `NewApplicationNotificationMessage`, `NewIssuePublisherNotificationMessage`, `PurchasePublisherNotificationMessage`, `QuotaLimitApproachingPublisherNotificationMessage`, and `RequestPublisherNotificationMessage`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management Notification Recipient Email.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Notification Recipient Email.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Notification Recipient Email.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Notification Recipient Email.

## Import

API Management Notification Recipient Emails can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_notification_recipient_email.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/notifications/notificationName1/recipientEmails/email1
```
