---
subcategory: "Appservice"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_web_app_active_slot"
description: |-
	Manages a Web App Active Slot.
---

# azurerm_web_app_active_slot

Manages a Web App Active Slot.

## Example Usage

```hcl
resource "azurerm_web_app_active_slot" "example" {
  slot_id = "example"

}
```

## Arguments Reference

The following arguments are supported:

* `slot_id` - (Required) The ID of the Slot to swap with `Production`.

---

* `overwrite_network_config` - (Optional) The swap action should overwrite the Production slot's network configuration with the configuration from this slot. Defaults to `true`. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Web App Active Slot

* `last_successful_swap` - The timestamp of the last successful swap with `Production`


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Web App Active Slot.
* `update` - (Defaults to 30 minutes) Used when updating the Web App Active Slot.
* `read` - (Defaults to 5 minutes) Used when retrieving the Web App Active Slot.
* `delete` - (Defaults to 5 minutes) Used when deleting the Web App Active Slot.

## Import

a Web App Active Slot can be imported using the `resource id`, e.g.

```shell
terraform import WebAppActiveSlot.example "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1"
```