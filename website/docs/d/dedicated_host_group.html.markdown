---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dedicated_host_group"
sidebar_current: "docs-azurerm-datasource-dedicated-host-group"
description: |-
  Gets information about an existing Dedicated Host Group
---

# Data Source: azurerm_dedicated_host_group

Use this data source to access information about an existing Dedicated Host Group.


## Example Usage

```hcl
data "azurerm_dedicated_host_group" "example" {
  resource_group_name = "example-rg"
  name                = "example-dedicated-host-group" 
}

output "id" {
  value = data.azurerm_dedicated_host_group.example.id
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Dedicated Host Group.

* `resource_group_name` - (Required) Specifies the name of the resource group the Dedicated Host Group is located in.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Dedicated Host Group.

* `location` - The Azure location where the Dedicated Host Group exists.

* `host_ids` - A list of ID of all dedicated hosts in the Dedicated Host Group.

* `platform_fault_domain_count` - Number of fault domains that the Dedicated Host Group can span.

* `zones` - Availability Zone the Dedicated Host Group is scoped. Only single zone is supported. If it is absent, it means the group is using all zones in the region.

* `tags` - A mapping of tags assigned to the resource.
