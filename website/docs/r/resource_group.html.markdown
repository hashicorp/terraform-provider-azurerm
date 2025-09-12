---
subcategory: "Base"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_resource_group"
description: |-
  Manages a Resource Group.
---

# azurerm_resource_group

Manages a Resource Group.

-> **Note:** Azure automatically deletes any Resources nested within the Resource Group when a Resource Group is deleted.

-> **Note:** Version 2.72 and later of the Azure Provider include a Feature Toggle which can error if there are any Resources left within the Resource Group at deletion time. This Feature Toggle is disabled in 2.x but enabled by default from 3.0 onwards, and is intended to avoid the unintentional destruction of resources managed outside of Terraform (for example, provisioned by an ARM Template). See [the Features block documentation](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/features-block) for more information on Feature Toggles within Terraform.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Resource Group should exist. Changing this forces a new Resource Group to be created.

* `name` - (Required) The Name which should be used for this Resource Group. Changing this forces a new Resource Group to be created.

---

* `managed_by` - (Optional) The ID of the resource or application that manages this Resource Group.

* `tags` - (Optional) A mapping of tags which should be assigned to the Resource Group.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Resource Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the Resource Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Resource Group.
* `update` - (Defaults to 90 minutes) Used when updating the Resource Group.
* `delete` - (Defaults to 90 minutes) Used when deleting the Resource Group.

## Import

Resource Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_resource_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1
```
