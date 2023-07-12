---
subcategory: "Communication"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_communication_service"
description: |-
  Gets information about an existing Communication Service.
---

# Data Source: azurerm_communication_service

Use this data source to access information about an existing Communication Service.

## Example Usage

```hcl
data "azurerm_communication_service" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_communication_service.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Communication Service.
* 
* `resource_group_name` - (Required) The name of the Resource Group where the Communication Service exists.
* 
## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Communication Service.

* `data_location` - The location where the Communication service stores its data at rest.

* `primary_connection_string` - The primary connection string of the Communication Service.

* `primary_key` - The primary key of the Communication Service.

* `secondary_connection_string` - The secondary connection string of the Communication Service.

* `secondary_key` - The secondary key of the Communication Service.

* `tags` - A mapping of tags assigned to the Communication Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Communication Service.
