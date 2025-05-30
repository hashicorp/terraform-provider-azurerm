---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_kubernetes_fleet_manager"
description: |-
  Gets information about an existing Kubernetes Fleet Manager.
---

# Data Source: azurerm_kubernetes_fleet_manager

Use this data source to access information about an existing Kubernetes Fleet Manager.

## Example Usage

```hcl
data "azurerm_kubernetes_fleet_manager" "example" {
  name                = "example"
  resource_group_name = "example-resource-group"
}

output "id" {
  value = data.azurerm_kubernetes_fleet_manager.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Kubernetes Fleet Manager.

* `resource_group_name` - (Required) The name of the Resource Group where the Kubernetes Fleet Manager exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Kubernetes Fleet Manager.

* `location` - The Azure Region where the Kubernetes Fleet Manager exists.

* `tags` - A mapping of tags assigned to the Kubernetes Fleet Manager.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Kubernetes Fleet Manager.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.ContainerService`: 2024-04-01
