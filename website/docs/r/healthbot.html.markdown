---
subcategory: "Bot"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_healthbot"
description: |-
  Manages a Healthbot Service.
---

# azurerm_healthbot

Manages a Healthbot Service.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-healthbot"
  location = "West Europe"
}

resource "azurerm_healthbot" "example" {
  name                = "example-bot"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku_name            = "F0"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies The name of the Healthbot Service resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies The name of the Resource Group in which to create the Healthbot Service. changing this forces a new resource to be created.

* `location` - (Required) Specifies The Azure Region where the resource exists. Changing this force a new resource to be created.

* `sku_name` - (Required) The name which should be used for the SKU of the service. Possible values are `C0`, `F0` and `S1`.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the service.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the resource.

* `bot_management_portal_url` - The management portal url.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Healthbot Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Healthbot Service.
* `update` - (Defaults to 30 minutes) Used when updating the Healthbot Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the Healthbot Service.

## Import

Healthbot Service can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_healthbot.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.HealthBot/healthBots/bot1
```
