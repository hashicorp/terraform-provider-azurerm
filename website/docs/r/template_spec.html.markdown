---
subcategory: "Template"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_template_spec"
description: |-
  Manages a Template Spec.
---

# azurerm_template_spec

Manages a Template Spec.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_template_spec" "example" {
  name                = "example-templatespec"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Template Spec. Changing this forces a new Template Spec to be created.

* `location` - (Required) The Azure Region where the Template Spec should exist. Changing this forces a new Template Spec to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Template Spec should exist. Changing this forces a new Template Spec to be created.

---

* `description` - (Optional) The description of the Template Spec.

* `display_name` - (Optional) The display name of the Template Spec.

* `tags` - (Optional) A mapping of tags which should be assigned to the Template Spec.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Template Spec.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the Template Spec.
* `read` - (Defaults to 5 minutes) Used when retrieving the Template Spec.
* `update` - (Defaults to 3 hours) Used when updating the Template Spec.
* `delete` - (Defaults to 3 hours) Used when deleting the Template Spec.

## Import

Template Specs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_template_spec.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Resources/templateSpecs/spec1
```
