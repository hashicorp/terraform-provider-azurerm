---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_custom_ip_prefix"
description: |-
  Manages a Custom Ip Prefix.

---

# azurerm_custom_ip_prefix

Manages a Custom Ip Prefix.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_custom_ip_prefix" "example" {
  name                  = "example-CustomIpPrefix"
  location              = azurerm_resource_group.example.location
  resource_group_name   = azurerm_resource_group.example.name
  cidr                  = "0.0.0.0/24"
  action                = "Provision"
  authorization_message = "00000000-0000-0000-0000-000000000000|0.0.0.0/24|20991212"
  signed_message        = "singed message for WAN validation"
  zones                 = ["1","2","3"]

  tags = {
    env = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Custom Ip Prefix. Changing this forces a new resource to be created.

* `location` - (Required) The location where the Custom Ip Prefix should exist. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which to create the Custom Ip Prefix. Changing this forces a new resource to be created.

* `cidr` - (Required) The `cidr` of the Custom Ip Prefix. Changing this forces a new resource to be created.

~> **NOTE:** Currently, the `cidr` only supports IPv4.

* `zones` - (Required) A list of availability zones which the Custom Ip Prefix should be allocated. Possible values are `1`, `2`, `3`. Changing this forces a new resource to be created.

* `authorization_message` - (Optional) The authorization message for WAN validation. Changing this forces a new resource to be created.

-> **NOTE:** The `authorization_message` should be formatted as "<subscriptionId>|<cidr>|<yyyMMdd>", such as "00000000-0000-0000-0000-00000000|0.0.0.0/24|20221231".

* `signed_message` - (Optional) The signed message for WAN validation. Changing this forces a new resource to be created.

* `action` - (Optional) The commission action of the Custom Ip Prefix. Possible values are `Provision`,`Commission`,`Decommission` or `Deprovision`. The default is `Provision`.

* `tags` - (Optional) A mapping of tags to assign to the Custom Ip Prefix.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Custom Ip Prefix.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Custom Ip Prefix.
* `update` - (Defaults to 30 minutes) Used when updating the Custom Ip Prefix.
* `read` - (Defaults to 5 minutes) Used when retrieving the Custom Ip Prefix.
* `delete` - (Defaults to 30 minutes) Used when deleting the Custom Ip Prefix.

## Import

Custom Ip Prefix can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_custom_ip_prefix.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/customIPPrefixes/customIpPrefix1
```
