---
subcategory: "Mobile Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mobile_network_sim_group"
description: |-
  Get information about a Mobile Network Sim Group.
---

# azurerm_mobile_network_sim_group

Get information about a Mobile Network Sim Group.

## Example Usage

```hcl
data "azurerm_mobile_network_sim_group" "example" {
  name                = "example-mnsg"
  resource_group_name = "example-rg"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Mobile Network Sim Groups. Changing this forces a new Mobile Network Sim Group to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Mobile Network Sim Groups should exist. Changing this forces a new Mobile Network Sim Group to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mobile Network Sim Groups.

* `location` - The Azure Region where the Mobile Network Sim Groups should exist.

* `mobile_network_id` - The ID of Mobile Network which the Mobile Network Sim Group belongs to.

* `encryption_key_url` - A key to encrypt the SIM data that belongs to this SIM group.

* `identity` - An `identity` block as defined below.

-> **NOTE:** A `UserAssigned` identity must be specified when `encryption_key_url` is specified.

* `tags` - A mapping of tags which should be assigned to the Mobile Network Sim Groups.

---

An `identity` block supports the following:

* `type` - The type of Managed Service Identity. Possible values are `SystemAssigned`, `UserAssigned`.

* `identity_ids` - A list of IDs for User Assigned Managed Identity resources to be assigned.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Mobile Network Sim Groups.

