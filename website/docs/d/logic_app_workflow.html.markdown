---
subcategory: "Logic App"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_logic_app_workflow"
description: |-
  Gets information about an existing Logic App Workflow.
---

# Data Source: azurerm_logic_app_workflow

Use this data source to access information about an existing Logic App Workflow.

## Example Usage

```hcl
data "azurerm_logic_app_workflow" "example" {
  name                = "workflow1"
  resource_group_name = "my-resource-group"
}

output "access_endpoint" {
  value = "${data.azurerm_logic_app_workflow.example.access_endpoint}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Logic App Workflow.

* `resource_group_name` - (Required) The name of the Resource Group in which the Logic App Workflow exists.

## Attributes Reference

The following attributes are exported:

* `id` - The Logic App Workflow ID.

* `location` - The Azure location where the Logic App Workflow exists.

* `workflow_schema` - The Schema used for this Logic App Workflow.

* `workflow_version` - The version of the Schema used for this Logic App Workflow. Defaults to `1.0.0.0`.

* `parameters` - A map of Key-Value pairs.

* `tags` - A mapping of tags assigned to the resource.

* `access_endpoint` - The Access Endpoint for the Logic App Workflow

### Timeouts

~> **Note:** Custom Timeouts are available [as an opt-in Beta in version 1.43 of the Azure Provider](/docs/providers/azurerm/guides/2.0-beta.html) and will be enabled by default in version 2.0 of the Azure Provider.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Logic App Workflow.
