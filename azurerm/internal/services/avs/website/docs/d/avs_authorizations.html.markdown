---
subcategory: "Avs"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_avs_authorization"
description: |-
  Gets information about an existing avs Authorization.
---

# Data Source: azurerm_avs_authorization

Use this data source to access information about an existing avs Authorization.

## Example Usage

```hcl
data "azurerm_avs_authorization" "example" {
  name = "example-authorization"
  resource_group_name = "example-resource-group"
  private_cloud_name = "existing"
}

output "id" {
  value = data.azurerm_avs_authorization.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this avs Authorization.

* `resource_group_name` - (Required) The name of the Resource Group where the avs Authorization exists.

* `private_cloud_name` - (Required) Name of the private cloud.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the avs Authorization.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the avs Authorization.