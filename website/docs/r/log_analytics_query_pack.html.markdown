---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_query_pack"
description: |-
  Manages a Log Analytics Query Pack.
---

# azurerm_log_analytics_query_pack

Manages a Log Analytics Query Pack.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_log_analytics_query_pack" "example" {
  name                = "example-laqp"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Log Analytics Query Pack. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Log Analytics Query Pack should exist. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Log Analytics Query Pack should exist. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Log Analytics Query Pack.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Log Analytics Query Pack.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Log Analytics Query Pack.
* `read` - (Defaults to 5 minutes) Used when retrieving the Log Analytics Query Pack.
* `update` - (Defaults to 30 minutes) Used when updating the Log Analytics Query Pack.
* `delete` - (Defaults to 30 minutes) Used when deleting the Log Analytics Query Pack.

## Import

Log Analytics Query Packs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_log_analytics_query_pack.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.OperationalInsights/queryPacks/queryPack1
```
