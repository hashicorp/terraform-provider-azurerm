---
subcategory: "Fluid Relay Service"
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

* `tags` - (Optional) A mapping of tags which should be assigned to the Fuild Relay Server.

* `identity` - (Optional) An `identity` block as defined below.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this API Management Service. Possible values are `SystemAssigned`, `UserAssigned`.

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this API Management Service.

~> **NOTE:** This is required when `type` is set to `UserAssigned`

~> **NOTE:** When `type` is set to `SystemAssigned`, the assigned `principal_id` and `tenant_id` can be retrieved after the Fluid Relay Server has been created. More details are available below.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Fuild Relay Server.

* `frs_tenant_id` - The Fluid tenantId for this server.

* `orderer_endpoints` - An array of the Fluid Relay Orderer endpoints.

* `principal_id` - The principal ID of the Fluid Relay Server.

* `storage_endpoints` - An array of the Fluid Relay storage endpoints.

* `tenant_id` - The tenant ID of the Fluid Relay Server.

---

`identity` exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Identity of this SQL Server.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Identity of this SQL Server.

-> You can access the Principal ID via `azurerm_fluid_relay_server.example.identity.0.principal_id` and the Tenant ID via `azurerm_fluid_relay_server.example.identity.0.tenant_id`

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Fuild Relay Server.
* `read` - (Defaults to 5 minutes) Used when retrieving the Fuild Relay Server.
* `update` - (Defaults to 10 minutes) Used when updating the Fuild Relay Server.
* `delete` - (Defaults to 10 minutes) Used when deleting the Fuild Relay Server.

## Import

Fuild Relay Servers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_fluid_relay_server.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.FluidRelay/fluidRelayServers/server1
```
