---
subcategory: "Fluid Relay"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_fluid_relay_server"
description: |-
  Manages a Fluid Relay Server.
---

# azurerm_fluid_relay_server

Manages a Fluid Relay Server.

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
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Fluid Relay Server should exist. Changing this forces a new Fluid Relay Server to be created.

* `name` - (Required) The name which should be used for this Fluid Relay Server. Changing this forces a new Fluid Relay Server to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Fluid Relay Server should exist. Changing this forces a new Fluid Relay Server to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `storage_sku` - (Optional) Sku of the storage associated with the resource, Possible values are `standard` and `basic`. Changing this forces a new Fluid Relay Server to be created.

* `customer_managed_key` - (Optional) A `customer_managed_key` block as defined below. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Fluid Relay Server.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Fluid Relay Service. Possible values are `SystemAssigned`,`UserAssigned` and `SystemAssigned, UserAssigned`.

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Fluid Relay Service.

---

An `customer_managed_key` block supports the following:

* `key_vault_key_id` - (Required) The Key Vault Key Id that will be used to encrypt the Fluid Relay Server.

* `user_assigned_identity_id` - (Required) The User Assigned Managed Identity ID to be used for accessing the Customer Managed Key for encryption.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Fluid Relay Server.

* `frs_tenant_id` - The Fluid tenantId for this server.

* `primary_key` - The primary key for this server.

* `secondary_key` - The secondary key for this server.

* `orderer_endpoints` - An array of the Fluid Relay Orderer endpoints. This will be deprecated in future version of fluid relay server and will always be empty, [more details](https://learn.microsoft.com/en-us/azure/azure-fluid-relay/concepts/version-compatibility).

* `storage_endpoints` - An array of storage endpoints for this Fluid Relay Server. This will be deprecated in future version of fluid relay server and will always be empty, [more details](https://learn.microsoft.com/en-us/azure/azure-fluid-relay/concepts/version-compatibility).

* `service_endpoints` - An array of service endpoints for this Fluid Relay Server.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Identity of this Fluid Relay Server.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Identity of this Fluid Relay Server.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Fluid Relay Server.
* `read` - (Defaults to 5 minutes) Used when retrieving the Fluid Relay Server.
* `update` - (Defaults to 10 minutes) Used when updating the Fluid Relay Server.
* `delete` - (Defaults to 10 minutes) Used when deleting the Fluid Relay Server.

## Import

Fluid Relay Servers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_fluid_relay_server.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.FluidRelay/fluidRelayServers/server1
```
