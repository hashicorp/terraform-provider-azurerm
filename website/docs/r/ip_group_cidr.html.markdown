---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_ip_group_cidr"
description: |-
  Manages a IP Group CIDR.
---

# azurerm_ip_group_cidr

Manages IP Group CIDR records.

~> **Note:** Warning Do not use this resource at the same time as the `cidrs` property of the
`azurerm_ip_group` resource for the same IP Group. Doing so will cause a conflict and
CIDRS will be removed.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "test-rg"
  location = "West Europe"
}

resource "azurerm_ip_group" "example" {
  name                = "test-ipgroup"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_ip_group_cidr" "example" {
  ip_group_id = azurerm_ip_group.example.id
  cidr        = "10.10.10.0/24"
}
```

## Arguments Reference

The following arguments are supported:

* `ip_group_id` - (Required) The ID of the destination IP Group.
Changing this forces a new IP Group CIDR to be created.

* `cidr` - (Required) The `CIDR` that should be added to the IP Group.
Changing this forces a new IP Group CIDR to be created.

~> **Note:** The AzureRM Terraform provider provides cidr support via this standalone resource and in-line within [azurerm_ip_group](ip_group.html) using the `cidrs` property. You cannot use both methods simultaneously. If cidrs are set via this resource then `ignore_changes` should be used in the resource `azurerm_ip_group_cidr` configuration.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the IP Group CIDR.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IP Group CIDR.
* `read` - (Defaults to 5 minutes) Used when retrieving the IP Group CIDR.
* `delete` - (Defaults to 30 minutes) Used when deleting the IP Group CIDR.

## Import

IP Group CIDRs can be imported using the `resource id` of the IP Group and
the CIDR value (`/` characters have to be replaced by `_`), e.g.

```shell
terraform import azurerm_ip_group_cidr.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/ipGroups/test-ipgroup/cidrs/10.1.0.0_24
```
