---
subcategory: "Mobile Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mobile_network_sim_group"
description: |-
  Manages a Mobile Network Sim Group.
---

# azurerm_mobile_network_sim_group

Manages a Mobile Network Sim Group.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_mobile_network" "example" {
  name                = "example-mn"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  mobile_country_code = "001"
  mobile_network_code = "01"
}

data "azurerm_user_assigned_identity" "example" {
  name                = "name_of_user_assigned_identity"
  resource_group_name = "name_of_resource_group"
}

data "azurerm_key_vault" "example" {
  name                = "example-kv"
  resource_group_name = "some-resource-group"
}

data "azurerm_key_vault_key" "example" {
  name         = "example-key"
  key_vault_id = data.azurerm_key_vault.example.id
}

resource "azurerm_mobile_network_sim_group" "example" {
  name               = "example-mnsg"
  location           = azurerm_resource_group.example.location
  mobile_network_id  = azurerm_mobile_network.example.id
  encryption_key_url = data.azurerm_key_vault_key.example.id

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [data.azurerm_user_assigned_identity.example.id]
  }

  tags = {
    key = "value"
  }

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Mobile Network Sim Groups. Changing this forces a new Mobile Network Sim Group to be created.

* `location` - (Required) Specifies the Azure Region where the Mobile Network Sim Groups should exist. Changing this forces a new Mobile Network Sim Group to be created.

* `mobile_network_id` - (Required) The ID of Mobile Network which the Mobile Network Sim Group belongs to. Changing this forces a new Mobile Network Slice to be created.

* `encryption_key_url` - (Optional) A key to encrypt the SIM data that belongs to this SIM group.

* `identity` - (Optional) An `identity` block as defined below.

-> **Note:** A `UserAssigned` identity must be specified when `encryption_key_url` is specified.

* `tags` - (Optional) A mapping of tags which should be assigned to the Mobile Network Sim Groups.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity. Possible value is `UserAssigned`.

* `identity_ids` - (Required) A list of IDs for User Assigned Managed Identity resources to be assigned.


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mobile Network Sim Groups.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the Mobile Network Sim Groups.
* `read` - (Defaults to 5 minutes) Used when retrieving the Mobile Network Sim Groups.
* `update` - (Defaults to 3 hours) Used when updating the Mobile Network Sim Groups.
* `delete` - (Defaults to 3 hours) Used when deleting the Mobile Network Sim Groups.

## Import

Mobile Network Sim Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mobile_network_sim_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.MobileNetwork/simGroups/simGroup1
```
