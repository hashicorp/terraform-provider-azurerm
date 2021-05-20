---
subcategory: "Communication"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_communication_service"
description: |-
  Manages a Communication Service.
---

# azurerm_communication_service

Manages a Communication Service.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_communication_service" "example" {
  name                = "example-communicationservice"
  resource_group_name = azurerm_resource_group.example.name
  data_location       = "United States"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Communication Service resource. Changing this forces a new Communication Service to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Communication Service should exist. Changing this forces a new Communication Service to be created.

---

* `data_location` - (Optional) The location where the Communication service stores its data at rest. Possible values are `Asia Pacific`, `Australia`, `Europe`, `UK` and `United States`. Defaults to `United States`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Communication Service.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Communication Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Communication Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Communication Service.
* `update` - (Defaults to 30 minutes) Used when updating the Communication Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the Communication Service.

## Import

Communication Services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_communication_service.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Communication/CommunicationServices/communicationService1
```
