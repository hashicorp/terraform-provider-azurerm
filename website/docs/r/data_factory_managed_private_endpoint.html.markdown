---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_managed_private_endpoint"
description: |-
  Manages a Data Factory Managed Private Endpoint.
---

# azurerm_data_factory_managed_private_endpoint

Manages a Data Factory Managed Private Endpoint.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_data_factory" "example" {
  name                            = "example"
  location                        = azurerm_resource_group.example.location
  resource_group_name             = azurerm_resource_group.example.name
  managed_virtual_network_enabled = true
}

resource "azurerm_storage_account" "example" {
  name                     = "example"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_data_factory_managed_private_endpoint" "example" {
  name               = "example"
  data_factory_id    = azurerm_data_factory.example.id
  target_resource_id = azurerm_storage_account.example.id
  subresource_name   = "blob"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Managed Private Endpoint. Changing this forces a new resource to be created.

* `data_factory_id` - (Required) The ID of the Data Factory on which to create the Managed Private Endpoint. Changing this forces a new resource to be created.

* `target_resource_id` - (Required) The ID of the Private Link Enabled Remote Resource which this Data Factory Private Endpoint should be connected to. Changing this forces a new resource to be created.

* `subresource_name` - (Required) Specifies the sub resource name which the Data Factory Private Endpoint is able to connect to. Changing this forces a new resource to be created.

-> **NOTE:** Possible values are listed in [documentation](https://docs.microsoft.com/en-us/azure/private-link/private-endpoint-overview#dns-configuration).

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Data Factory Managed Private Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory Managed Private Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Managed Private Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory Managed Private Endpoint.

## Import

Data Factory Managed Private Endpoint can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_managed_private_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/managedVirtualNetworks/default/managedPrivateEndpoints/endpoint1
```
