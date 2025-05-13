---
subcategory: "Graph Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_graph_services_account"
description: |-
  Manages a Microsoft Graph Services Account.
---

# azurerm_graph_services_account

Manages a Microsoft Graph Services Account.

## Example Usage

```hcl
resource "azuread_application" "example" {
  display_name = "example-app"
}
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}
resource "azurerm_graph_services_account" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  application_id      = azuread_application.example.application_id
  tags = {
    environment = "Production"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of this Account. Changing this forces a new Account to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group within which this Account should exist. Changing this forces a new Account to be created.

* `application_id` - (Required) Customer owned application ID. Changing this forces a new Account to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Account.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Account.

* `billing_plan_id` - Billing Plan Id.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the Account.
* `update` - (Defaults to 30 minutes) Used when updating the Account.
* `delete` - (Defaults to 30 minutes) Used when deleting the Account.

## Import

An existing Account can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_graph_services_account.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.GraphServices/accounts/account1
```
