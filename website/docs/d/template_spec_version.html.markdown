---
subcategory: "Template"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_template_spec_version"
description: |-
  Gets information about an existing Template Spec Version.
---

# Data Source: azurerm_template_spec_version

Use this data source to access information about an existing Template Spec Version.

## Example Usage

```hcl
data "azurerm_template_spec_version" "example" {
  name                = "exampleTemplateSpec"
  resource_group_name = "MyResourceGroup"
  version             = "v1.0.4"
}

output "id" {
  value = data.azurerm_template_spec_version.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Template Spec.

* `resource_group_name` - (Required) The name of the Resource Group where the Template Spec exists.

* `version` - (Required) The Version Name of the Template Spec.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Template Spec version.

* `template_body` - The ARM Template body of the Template Spec Version.

* `tags` - A mapping of tags assigned to the Template.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Template.
