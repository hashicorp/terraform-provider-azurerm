---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_ip_group"
description: |-
  Manages an IP group which contains a list of CIDRs and/or IP addresses.

---

# azurerm_ip_group

Manages an IP group that contains a list of CIDRs and/or IP addresses.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "westus"
}

resource "azurerm_ip_group" "example" {
  name                = "example1-ipgroup"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  cidrs = ["192.168.0.1", "172.16.240.0/20", "10.48.0.0/12"]

  tags = {
    environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the IP group. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the IP group. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `cidrs` - (Optional) A list of CIDRs or IP addresses.

* `tags` - (Optional) A mapping of tags to assign to the resource.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the IP Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IP Group.
* `update` - (Defaults to 30 minutes) Used when updating the IP Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the IP Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the IP Group.

## Import

IP Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_ip_group.ipgroup1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/ipGroups/myIpGroup
```
