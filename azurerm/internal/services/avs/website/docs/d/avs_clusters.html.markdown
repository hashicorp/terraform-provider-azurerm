---
subcategory: "Avs"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_avs_cluster"
description: |-
  Gets information about an existing avs Cluster.
---

# Data Source: azurerm_avs_cluster

Use this data source to access information about an existing avs Cluster.

## Example Usage

```hcl
data "azurerm_avs_cluster" "example" {
  name = "example-cluster"
  resource_group_name = "example-resource-group"
  private_cloud_name = "existing"
}

output "id" {
  value = data.azurerm_avs_cluster.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this avs Cluster.

* `resource_group_name` - (Required) The name of the Resource Group where the avs Cluster exists.

* `private_cloud_name` - (Required) Name of the private cloud.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the avs Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the avs Cluster.