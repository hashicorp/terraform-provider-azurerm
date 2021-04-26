---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_email_template"
description: |-
  Manages a API Management Email Template.
---

# azurerm_api_management_email_template

Manages a API Management Email Template.

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

resource "azurerm_api_management_email_template" "example" {
  template_name       = "ConfirmSignUpIdentityDefault"
  resource_group_name = azurerm_resource_group.example.resource_group_name
  api_management_name = azurerm_api_management.example.name
  subject             = "Customized confirmation email for your new $OrganizationName API account"
  body                = <<EOF
<!DOCTYPE html >
<html>
<head>
  <meta charset="UTF-8" />
  <title>Customized Letter Title</title>
</head>
<body>
  <p style="font-size:12pt;font-family:'Segoe UI'">Dear $DevFirstName $DevLastName,</p>
</body>
</html>
EOF
}
```

## Arguments Reference

The following arguments are supported:

* `template_name` - (Required) The name of the Email Template. Possible values are `AccountClosedDeveloper`, `ApplicationApprovedNotificationMessage`, `ConfirmSignUpIdentityDefault`, `EmailChangeIdentityDefault`, `InviteUserNotificationMessage`, `NewCommentNotificationMessage`, `NewDeveloperNotificationMessage`, `NewIssueNotificationMessage`, `PasswordResetByAdminNotificationMessage`, `PasswordResetIdentityDefault`, `PurchaseDeveloperNotificationMessage`, `QuotaLimitApproachingDeveloperNotificationMessage`, `RejectDeveloperNotificationMessage`, `RequestDeveloperNotificationMessage`. Changing this forces a new API Management Email Template to be created.

* `api_management_name` - (Required) The name of the API Management Service in which the Email Template should exist. Changing this forces a new API Management Email Template to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the API Management Email Template should exist. Changing this forces a new API Management Email Template to be created.

* `subject` - (Required) The subject of the Email.

* `body` - (Required) The body of the Email. Its format has to be a well-formed HTML document.

-> **NOTE:** In `subject` and `body` predefined parameters can be used. The available parameters depend on the template. Schema to use a parameter: `$` followed by the `parameter.name` - `$<parameter.name>`. The available parameters can be seen in the Notification templates section of the API-Management Service instance within the Azure Portal.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management Email Template.

* `title` - The title of the Email Template.

* `description` - The description of the Email Template.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Email Template.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Email Template.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Email Template.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Email Template.

## Import

API Management Email Templates can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_email_template.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/instance1/templates/template1
```
