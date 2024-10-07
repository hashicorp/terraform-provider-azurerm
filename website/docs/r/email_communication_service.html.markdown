---
subcategory: "Communication"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_email_communication_service"
description: |-
  Manages an Email Communication Service.
---

# azurerm_email_communication_service

Manages an Email Communication Service.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_email_communication_service" "example" {
  name                = "example-emailcommunicationservice"
  resource_group_name = azurerm_resource_group.example.name
  data_location       = "United States"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Email Communication Service resource. Changing this forces a new Email Communication Service to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Email Communication Service should exist. Changing this forces a new Email Communication Service to be created.

---

* `data_location` - (Required) The location where the Email Communication service stores its data at rest. Possible values are `Africa`, `Asia Pacific`, `Australia`, `Brazil`, `Canada`, `Europe`, `France`, `Germany`, `India`, `Japan`, `Korea`, `Norway`, `Switzerland`, `UAE`, `UK` and `United States`. Changing this forces a new Email Communication Service to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Email Communication Service.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Email Communication Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Email Communication Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Email Communication Service.
* `update` - (Defaults to 30 minutes) Used when updating the Email Communication Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the Email Communication Service.

## Import

Communication Services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_email_communication_service.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Communication/emailServices/emailCommunicationService1
```
