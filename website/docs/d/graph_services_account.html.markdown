---
subcategory: "Graph Services"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_graph_services_account"
description: |-
  Gets information about an existing Graph Services Account.
---

# Data Source: azurerm_graph_services_account

Use this data source to access information about an existing Graph Services Account.

## Example Usage

```hcl
data "azurerm_graph_services_account" "example" {
  name                = "example-graph-services-account"
  resource_group_name = "example-resources"
}

output "application_id" {
  value = data.azurerm_graph_services_account.example.application_id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Graph Services Account.

* `resource_group_name` - (Required) The name of the Resource Group where the Graph Services Account exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Graph Services Account.

* `application_id` - The application ID associated with the Graph Services Account.

* `billing_plan_id` - The billing plan ID for the Graph Services Account.

* `tags` - A mapping of tags assigned to the Graph Services Account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Graph Services Account.