---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dedicated_host"
description: |-
  Gets information about an existing Dedicated Host
---

# Data Source: azurerm_dedicated_host

Use this data source to access information about an existing Dedicated Host.

## Example Usage

```hcl
data "azurerm_dedicated_host" "example" {
  name                = "example-dh"
  resource_group_name = "example-rg"
  host_group_name     = "example-dhg"
}

output "dedicated_host_id" {
  value = "${data.azurerm_dedicated_host.example.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Dedicated Host.

* `resource_group_name` - (Required) Specifies the name of the resource group the Dedicated Host is located in.

* `host_group_name` - (Required) Specifies the name of the Dedicated Host Group the Dedicated Host is located in.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of Dedicated Host.

* `location` - The location where the Dedicated Host exists.

* `tags` - A mapping of tags assigned to the Dedicated Host.

