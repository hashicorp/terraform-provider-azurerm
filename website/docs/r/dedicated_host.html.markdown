---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dedicated_host"
description: |-
  Manage an Dedicated Host.
---

# azurerm_dedicated_host

Manage an Dedicated Host.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West US"
}

resource "azurerm_dedicated_host_group" "example" {
  name                        = "example-dhg"
  resource_group_name         = azurerm_resource_group.example.name
  location                    = azurerm_resource_group.example.location
  platform_fault_domain_count = 2
}


resource "azurerm_dedicated_host" "example" {
  name                  = "example-dh"
  location              = azurerm_resource_group.example.location
  resource_group_name   = azurerm_resource_group.example.name
  host_group_name       = azurerm_dedicated_host_group.example.name
  sku_name              = "DSv3-Type1"
  platform_fault_domain = 1
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specify the name of the Dedicated Host. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Dedicated Host. Changing this forces a new resource to be created.

* `location` - (Required) Specify the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `host_group_name` - (Required) Specify the name of the Dedicated Host Group in which to create the Dedicated Host. Changing this forces a new resource to be created.

* `sku_name` - (Required) Specify the sku name of the Dedicated Host. Possible values are `DSv3-Type1`, `ESv3-Type1`, `FSv2-Type2`. Changing this forces a new resource to be created.

* `platform_fault_domain` - (Required) Specify the fault domain of the Dedicated Host Group in which to create the Dedicated Host. Changing this forces a new resource to be created.

* `auto_replace_on_failure` - (Optional) Specifies whether the Dedicated Host should be replaced automatically in case of a failure. The value is defaulted to `true` when not provided.

* `license_type` - (Optional) Specifies the software license type that will be applied to the VMs deployed on the Dedicated Host. Possible values are: `None`, `Windows_Server_Hybrid`, `Windows_Server_Perpetual`. The value is defaulted to `None` when not provided.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Dedicated Host.

## Import

Dedicated Host can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_dedicated_host.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/hostGroups/group1/hosts/host1
```
