---
subcategory: ""
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

* `name` - (Required) The name of the private link endpoint.
* `resource_group_name` - (Required) The name of the resource group in which the private link endpoint resides.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Prviate Link Endpoint.
* `location` - The supported Azure location where the resource exists.
* `network_interface_ids` - A list of network interface resource IDs that are being used by the endpoint.
* `subnet_id` - The resource ID of the subnet to be used by the endpoint.
* `tags` - A mapping of tags assigned to the resource.
