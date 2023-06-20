---
subcategory: "Authorization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_user_assigned_identity"
description: |-
  Manages a User Assigned Identity.
---

<!-- Note: This documentation is generated. Any manual changes will be overwritten -->

# azurerm_user_assigned_identity

Manages a User Assigned Identity.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}
resource "azurerm_user_assigned_identity" "example" {
  location            = azurerm_resource_group.example.location
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the User Assigned Identity should exist. Changing this forces a new User Assigned Identity to be created.

* `name` - (Required) Specifies the name of this User Assigned Identity. Changing this forces a new User Assigned Identity to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group within which this User Assigned Identity should exist. Changing this forces a new User Assigned Identity to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the User Assigned Identity.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the User Assigned Identity.

* `client_id` - The ID of the app associated with the Identity.

* `principal_id` - The ID of the Service Principal object associated with the created Identity.

* `tenant_id` - The ID of the Tenant which the Identity belongs to.

---



## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating this User Assigned Identity.
* `delete` - (Defaults to 30 minutes) Used when deleting this User Assigned Identity.
* `read` - (Defaults to 5 minutes) Used when retrieving this User Assigned Identity.
* `update` - (Defaults to 30 minutes) Used when updating this User Assigned Identity.

## Import

An existing User Assigned Identity can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_user_assigned_identity.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/{userAssignedIdentityName}
```

* Where `{subscriptionId}` is the ID of the Azure Subscription where the User Assigned Identity exists. For example `12345678-1234-9876-4563-123456789012`.
* Where `{resourceGroupName}` is the name of Resource Group where this User Assigned Identity exists. For example `example-resource-group`.
* Where `{userAssignedIdentityName}` is the name of the User Assigned Identity. For example `userAssignedIdentityValue`.
