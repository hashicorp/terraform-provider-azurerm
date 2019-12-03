---
subcategory: ""
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dedicated_host_group"
sidebar_current: "docs-azurerm-resource-dedicated-host-group"
description: |-
  Manage Azure DedicatedHostGroup instance.
---

# azurerm_dedicated_host_group

Manage Azure DedicatedHostGroup instance.


## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West US"
}

resource "azurerm_dedicated_host_group" "example" {
  resource_group_name         = "${azurerm_resource_group.example.name}"
  location                    = "${azurerm_resource_group.example.location}"
  name                        = "example-dedicated-host-group"
  platform_fault_domain_count = 1
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the dedicated host group. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group. Changing this forces a new resource to be created.

* `location` - (Required) Resource location Changing this forces a new resource to be created.

* `platform_fault_domain_count` - (Required) Number of fault domains that the host group can span.

* `zones` - (Optional) Availability Zone to use for this host group. Only single zone is supported. The zone can be assigned only during creation. If not provided, the group supports all zones in the region. If provided, enforces each host in the group to be in the same zone. Changing this forces a new resource to be created.

* `tags` - (Optional) Resource tags Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `hosts` - A list of references to all dedicated `host`s in the dedicated host group. The definition of `host` is defined below.

* `id` - Resource Id

* `type` - Resource type


---

The `host` block contains the following:

* `id` - Resource Id

## Import

Dedicated Host Group can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_dedicated_host_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Compute/hostGroups/group1
```
