---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_app_service_environment_v3"
description: |-
  Gets information about an existing 3rd Generation (v3) App Service Environment.
---

# Data Source: azurerm_app_service_environment_v3

Use this data source to access information about an existing 3rd Generation (v3) App Service Environment.

## Example Usage

```hcl
data "azurerm_app_service_environment_v3" "example" {
  name                = "example-ASE"
  resource_group_name = "example-resource-group"
}

output "id" {
  value = data.azurerm_app_service_environment_v3.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this v3 App Service Environment.

* `resource_group_name` - (Required) The name of the Resource Group where the v3 App Service Environment exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the v3 App Service Environment.

* `cluster_setting` - A `cluster_setting` block as defined below.

* `subnet_id` - The ID of the v3 App Service Environment Subnet.

* `tags` - A mapping of tags assigned to the v3 App Service Environment.

---

A `cluster_setting` block exports the following:

* `name` - The name of the Cluster Setting.

* `value` - The value for the Cluster Setting.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the 3rd Generation (v3) App Service Environment.
