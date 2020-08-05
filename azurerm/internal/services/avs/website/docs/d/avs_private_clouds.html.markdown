---
subcategory: "Avs"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_avs_private_cloud"
description: |-
  Gets information about an existing avs PrivateCloud.
---

# Data Source: azurerm_avs_private_cloud

Use this data source to access information about an existing avs PrivateCloud.

## Example Usage

```hcl
data "azurerm_avs_private_cloud" "example" {
  name = "example-privatecloud"
  resource_group_name = "example-resource-group"
}

output "id" {
  value = data.azurerm_avs_private_cloud.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this avs PrivateCloud.

* `resource_group_name` - (Required) The name of the Resource Group where the avs PrivateCloud exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the avs PrivateCloud.

* `location` - The Azure Region where the avs PrivateCloud exists.

* `tags` - A mapping of tags assigned to the avs PrivateCloud.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the avs PrivateCloud.