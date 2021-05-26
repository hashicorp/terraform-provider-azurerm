---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dedicated_host_group"
description: |-
  Gets information about an existing Dedicated Host Group
---

# Data Source: azurerm_dedicated_host_group

Use this data source to access information about an existing Dedicated Host Group.

## Example Usage

```hcl
data "azurerm_dedicated_host_group" "example" {
  name                = "example-dedicated-host-group"
  resource_group_name = "example-rg"
}

output "id" {
  value = data.azurerm_dedicated_host_group.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - Specifies the name of the Dedicated Host Group.

* `resource_group_name` - Specifies the name of the resource group the Dedicated Host Group is located in.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Dedicated Host Group.

* `location` - The Azure location where the Dedicated Host Group exists.

* `platform_fault_domain_count` - The number of fault domains that the Dedicated Host Group spans.

* `automatic_placement_enabled` - Whether virtual machines or virtual machine scale sets be placed automatically on this Dedicated Host Group.

* `zones` - The Availability Zones in which this Dedicated Host Group is located.

* `tags` - A mapping of tags assigned to the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Dedicated Host Group.
