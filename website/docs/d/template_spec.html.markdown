---
subcategory: "Template"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_template_spec"
description: |-
  Gets information about an existing Template Spec
---

# Data Source: azurerm_template_spec

Use this data source to access information about an existing Template Spec

## Example Usage

```hcl
data "azurerm_template_spec" "example" {
  name                = "example-templatespec"
  resource_group_name = "example-resources"
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the Template Spec resource.

* `resource_group_name` - The name of the Resource Group where the Template Spec exists.

## Attributes Reference

The following attributes are exported:

* `id` - The Template Spec ID.

* `location` - The Azure location where the Template Spec exists.

* `description` - The description of the Template Spec.

* `display_name` - The display name of the Template Spec.

* `tags` - A mapping of tags assigned to the Template Spec.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Template Spec.
