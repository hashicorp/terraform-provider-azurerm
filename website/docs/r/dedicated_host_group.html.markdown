---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dedicated_host_group"
sidebar_current: "docs-azurerm-resource-dedicated-host-group"
description: |-
  Manage Azure DedicatedHostGroup instance.
---

# azurerm_dedicated_host_group

Manage a Azure Dedicated Host Group instance.


## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg-compute"
  location = "West Europe"
}

resource "azurerm_dedicated_host_group" "example" {
  resource_group_name         = "${azurerm_resource_group.example.name}"
  location                    = "${azurerm_resource_group.example.location}"
  name                        = "example-dedicated-host-group-compute"
  platform_fault_domain_count = 1
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Dedicated Host Group. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the resource group the Dedicated Host Group is located in. Changing this forces a new resource to be created.

* `location` - (Required) The Azure location where the Dedicated Host Group exists. Changing this forces a new resource to be created.

* `platform_fault_domain_count` - (Required) Number of fault domains that the host group can span. Changing this forces a new resource to be created.

* `zones` - (Optional) Availability Zone to use for this host group. Only single zone is supported. The zone can be assigned only during creation. If not provided, the group supports all zones in the region. If provided, enforces each host in the group to be in the same zone. Changing this forces a new resource to be created. 

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Dedicated Host Group.

## Import

Dedicated Host Group can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_dedicated_host_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Compute/hostGroups/group1
```
