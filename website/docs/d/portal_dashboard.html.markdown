---
subcategory: "Portal"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_portal_dashboard"
description: |-
  Gets information about an existing shared dashboard in the Azure Portal.
---

# Data Source: azurerm_portal_dashboard

Use this data source to access information about an existing shared dashboard in the Azure Portal. This is the data source of the `azurerm_dashboard` resource.

## Example Usage

```hcl

data "azurerm_portal_dashboard" "example" {
  name                = "existing-dashboard"
  resource_group_name = "dashboard-rg"
}

output "id" {
  value = data.azurerm_dashboard.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `resource_group_name` - (Required) Specifies the name of the resource group the shared Azure Portal Dashboard is located in.

* `name` - (Optional) Specifies the name of the shared Azure Portal Dashboard.

* `display_name` - (Optional) Specifies the display name of the shared Azure Portal Dashboard.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the shared Azure Portal dashboard.

* `location` - The Azure Region where the shared Azure Portal dashboard exists.

* `dashboard_properties` - JSON data representing dashboard body.

* `tags` - A mapping of tags assigned to the shared Azure Portal dashboard.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the shared Azure Dashboard.
