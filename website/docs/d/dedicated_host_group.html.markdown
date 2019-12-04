---
subcategory: ""
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
  value = "${data.azurerm_dedicated_host_group.example.id}"
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the dedicated host group.

* `resource_group_name` - (Required) The name of the resource group.


## Attributes Reference

The following attributes are exported:

* `id` - Resource Id

* `location` - Resource location

* `hosts` - A list of references to all dedicated `host`s in the dedicated host group. The definition of `host` is defined below.

* `platform_fault_domain_count` - Number of fault domains that the host group can span.

* `type` - Resource type

* `zones` - Availability Zone to use for this host group. Only single zone is supported. The zone can be assigned only during creation. If not provided, the group supports all zones in the region. If provided, enforces each host in the group to be in the same zone.

* `tags` - Resource tags


---

The `host` block contains the following:

* `id` - Resource Id
