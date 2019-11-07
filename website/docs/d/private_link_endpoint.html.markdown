---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_link_endpoint"
sidebar_current: "docs-azurerm-datasource-private-endpoint"
description: |-
  Gets information about an existing Private Link Endpoint
---

# Data Source: azurerm_private_link_endpoint

Use this data source to access information about an existing Private Link Endpoint.

## Example Usage

```hcl
data "azurerm_private_link_endpoint" "example" {
  resource_group_name = "example-rg"
  name                = "example-private-endpoint"
}

output "subnet_id" {
  value = "${data.azurerm_private_link_endpoint.example.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the Name of the `Private Link Endpoint`.
* `resource_group_name` - (Required) Specifies the Name of the Resource Group within which the `Private Link Endpoint` exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Prviate Link Endpoint.
* `location` - The supported Azure location where the resource exists.
* `network_interface_ids` - A list of network interfaces IDs.
* `subnet_id` - The subnet ID.
* `tags` - A mapping of tags assigned to the resource.
