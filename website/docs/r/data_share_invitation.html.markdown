---
subcategory: "Data Share"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_share_invitation"
description: |-
  Manages a Data Share Invitation.
---

# azurerm_data_share_invitation

Manages a Data Share Invitation.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_data_share_account" "example" {
  name = "example-dsa"
resource_group_name = azurerm_resource_group.example.name
  location = azurerm_resource_group.example.location
  
  tags = {
    foo = "bar"
  }
}

resource "azurerm_data_share_share" "example" {
  name = "example_dss"
  account_id = azurerm_data_share_account.example.id
share_kind = "CopyBased"
}

resource "azurerm_data_share_invitation" "example" {
  name = "example"
  share_id = azurerm_data_share_share.example.id
target_email = "123456@microsoft.com"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Data Share Invitation. Changing this forces a new Data Share Invitation to be created.

* `share_id` - (Required) The ID of the TODO. Changing this forces a new Data Share Invitation to be created.

* `target_email` - (Required) TODO. Changing this forces a new Data Share Invitation to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Data Share Invitation.

* `invitation_id` - The ID of the TODO.

* `invitation_status` - TODO.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Share Invitation.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Share Invitation.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Share Invitation.

## Import

Data Share Invitations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_share_invitation.example C:/Program Files/Git/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataShare/accounts/account1/shares/share1/invitations/invitation1
```