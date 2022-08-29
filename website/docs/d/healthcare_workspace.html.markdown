---
subcategory: "Healthcare"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_healthcare_workspace"
description: |-
  Get information about an existing Healthcare Workspace
---

# Data Source: azurerm_healthcare_workspace

Use this data source to access information about an existing Healthcare Workspace

## Example Usage

```hcl
data "azurerm_healthcare_workspace" "example" {
  name                = "example-healthcare_service"
  resource_group_name = "example-resources"
}

output "healthcare_workspace_id" {
  value = data.azurerm_healthcare_workspace.example.id
}
```

## Argument Reference

* `name` - The name of the Healthcare Workspace.

* `resource_group_name` - The name of the Resource Group in which the Healthcare Workspace exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Healthcare Workspace.

* `location` - The Azure Region where the Healthcare Workspace is located.

* `tags` - A map of tags assigned to the Healthcare Workspace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Healthcare Workspace.
