---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_application_gateway"
description: |-
  Gets information about an existing Application Gateway.
---

# Data Source: azurerm_application_gateway

Use this data source to access information about an existing Application Gateway.

## Example Usage

```hcl
data "azurerm_application_gateway" "example" {
  name                = "existing-app-gateway"
  resource_group_name = "existing-resources"
}

output "id" {
  value = data.azurerm_application_gateway.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Application Gateway.

* `resource_group_name` - (Required) The name of the Resource Group where the Application Gateway exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Application Gateway.

* `backend_address_pool` - A `backend_address_pool` block as defined below.

* `identity` - A `identity` block as defined below.

* `location` - The Azure Region where the Application Gateway exists.

* `tags` - A mapping of tags assigned to the Application Gateway.

---

A `backend_address_pool` block exports the following:

* `id` - The ID of the Backend Address Pool.

* `name` - The name of the Backend Address Pool.

* `fqdns` - A list of FQDN's that are included in the Backend Address Pool.

* `ip_addresses` - A list of IP Addresses that are included in the Backend Address Pool.

---

A `identity` block exports the following:

* `identity_ids` - A list of Managed Identity IDs assigned to this Application Gateway.

* `type` - The type of Managed Identity assigned to this Application Gateway.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Application Gateway.
