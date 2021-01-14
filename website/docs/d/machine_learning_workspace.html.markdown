---
subcategory: "Machine Learning"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_machine_learning_workspace"
description: |-
  Gets information about an existing Machine Learning Workspace
---

# Data Source: azurerm_machine_learning_workspace

Use this data source to access information about an existing Machine Learning Workspace.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_machine_learning_workspace" "existing" {
  name                = "example-workspace"
  resource_group_name = "example-resources"
}

output "id" {
  value = azurerm_machine_learning_workspace.existing.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Machine Learning Workspace exists.

* `resource_group_name` - (Required) The name of the Resource Group where the Machine Learning Workspace exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Machine Learning Workspace.

* `location` - The location where the Machine Learning Workspace exists.

* `tags` - A mapping of tags assigned to the Machine Learning Workspace.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Machine Learning Workspace.
