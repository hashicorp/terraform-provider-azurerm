---
subcategory: "Base"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_resource_tag"
description: |-
    Manages a single tag on a resource.
---

# azurerm_resource_tag

Manages a single tag on a resource.

~> **NOTE:** Terraform currently
provides both a standalone tag resource, and allows tags to be defined in-line within resources.
Using `azurerm_resource_tag` may remove the tag on the next run, unless a `lifecycle { ignore_changes = [tags["key"]]}` rule is defined on the resource.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"

  # Since the tag is managed outside of the resource, it needs to be defined separately
  lifecycle {
    ignore_changes = [tags["owner"]]
  }
}

resource "azurerm_resource_tag" "example" {
  resource_id = azurerm_resource_group.example.id

  key   = "owner"
  value = "Terraform"
}
```

## Arguments Reference

The following arguments are supported:

* `resource_id` - (Required) The resource id where the tag should be applied. Changing this forces a new Resource Group to be created.

* `key` - (Required) The key of the tag. Changing this forces a new Resource Group to be created.

* `value` - (Required) The value of the tag. Changing this forces a new Resource Group to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The internal ID of the Tag.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 15 minutes) Used when creating the Resource Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Resource Group.
* `update` - (Defaults to 15 minutes) Used when updating the Resource Group.
* `delete` - (Defaults to 15 minutes) Used when deleting the Resource Group.

## Import

Resource Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_resource_tag.example '/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example|tagKey'
```
