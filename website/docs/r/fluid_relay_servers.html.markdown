---
subcategory: "Active Directory Domain Service"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_fluid_relay_server"
description: |-
  Manages a Fuild Relay Server.
---

# azurerm_fluid_relay_server

Manages a Fuild Relay Server.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_fluid_relay_server" "example" {
  name = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Fuild Relay Server should exist. Changing this forces a new Fuild Relay Server to be created.

* `name` - (Required) The name which should be used for this Fuild Relay Server. Changing this forces a new Fuild Relay Server to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Fuild Relay Server should exist. Changing this forces a new Fuild Relay Server to be created.

---

* `encryption` - (Optional) One or more `encryption` blocks as defined below.

* `identity_type` - (Optional) The identity type,  value can be `SystemAssigned`, `UserAssigned`,`SystemAssigned, UserAssigned`, `None`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Fuild Relay Server.

* `user_assigned_identity` - (Optional) One or more `user_assigned_identity` blocks as defined below.

---

A `encryption` block supports the following:

* `identity_resource_id` - (Optional) user assigned identity to use for accessing key encryption key Url. Ex: /subscriptions/fa5fc227-a624-475e-b696-cdd604c735bc/resourceGroups/<resource group>/providers/Microsoft.ManagedIdentity/userAssignedIdentities/myId. Mutually exclusive with identityType systemAssignedIdentity.

* `identity_type` - (Optional) Values can be `SystemAssigned` or `UserAssigned`.

* `key_encryption_key_url` - (Optional) key encryption key Url, with or without a version. Ex: https://contosovault.vault.azure.net/keys/contosokek/562a4bb76b524a1493a6afe8e536ee78 or https://contosovault.vault.azure.net/keys/contosokek. Key auto rotation is enabled by providing a key uri without version. Otherwise, customer is responsible for rotating the key. The keyEncryptionKeyIdentity(either SystemAssigned or UserAssigned) should have permission to access this key url.

---

A `user_assigned_identity` block supports the following:

* `client_id` - (Optional) The client id of user assigned identity.

* `identity_id` - (Optional) user assigned identity to use for accessing key encryption key Url. Ex: /subscriptions/fa5fc227-a624-475e-b696-cdd604c735bc/resourceGroups/<resource group>/providers/Microsoft.ManagedIdentity/userAssignedIdentities/myId. Mutually exclusive with identityType systemAssignedIdentity..

* `principal_id` - (Optional) The principal id of user assigned identity.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Fuild Relay Server.

* `frs_tenant_id` - The Fluid tenantId for this server.

* `orderer_endpoints` - A `orderer_endpoints` block as defined below.

* `principal_id` - The principal ID of resource identity.

* `provisioning_state` - Provision states for FluidRelay RP, value can be `Succeeded`, `Failed`, `Canceled`.

* `storage_endpoints` - A `storage_endpoints` block as defined below.

* `tenant_id` - The tenant ID of resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Fuild Relay Server.
* `read` - (Defaults to 5 minutes) Used when retrieving the Fuild Relay Server.
* `update` - (Defaults to 10 minutes) Used when updating the Fuild Relay Server.
* `delete` - (Defaults to 10 minutes) Used when deleting the Fuild Relay Server.

## Import

Fuild Relay Servers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_fluid_relay_server.example /subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/myrg/providers/Microsoft.FluidRelay/fluidRelayServers/myFluid
```