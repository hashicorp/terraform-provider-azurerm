---
subcategory: "Authorization"
layout: "azurerm"
page_title: "Azure Resource Manager: azure_user_assigned_identity"
description: |-
  Manages a new user assigned identity.
---

# azurerm_user_assigned_identity

Manages a user assigned identity.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "acceptanceTestResourceGroup1"
  location = "eastus"
}

resource "azurerm_user_assigned_identity" "example" {
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  name = "search-api"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the user assigned identity. Changing this forces a
    new identity to be created.

* `resource_group_name` - (Required) The name of the resource group in which to
    create the user assigned identity.

* `location` - (Required) The location/region where the user assigned identity is
    created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The user assigned identity ID.

* `principal_id` - Service Principal ID associated with the user assigned identity.

* `client_id` - Client ID associated with the user assigned identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the User Assigned Identity.
* `update` - (Defaults to 30 minutes) Used when updating the User Assigned Identity.
* `read` - (Defaults to 5 minutes) Used when retrieving the User Assigned Identity.
* `delete` - (Defaults to 30 minutes) Used when deleting the User Assigned Identity.

## Import

User Assigned Identities can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_user_assigned_identity.exampleIdentity /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/acceptanceTestResourceGroup1/providers/Microsoft.ManagedIdentity/userAssignedIdentities/testIdentity
```
