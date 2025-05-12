---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_signalr_shared_private_link_resource"
description: |-
  Manages the Shared Private Link Resource for a Signalr service.
---

# azurerm_signalr_shared_private_link_resource

Manages the Shared Private Link Resource for a Signalr service.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "terraform-signalr"
  location = "east us"
}

resource "azurerm_key_vault" "example" {
  name                       = "examplekeyvault"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7
  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id
    certificate_permissions = [
      "ManageContacts",
    ]
    key_permissions = [
      "Create",
    ]
    secret_permissions = [
      "Set",
    ]
  }
}

resource "azurerm_signalr_service" "test" {
  name                = "tfex-signalr"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Standard_S1"
    capacity = 1
  }
}

resource "azurerm_signalr_shared_private_link_resource" "example" {
  name               = "tfex-signalr-splr"
  signalr_service_id = azurerm_signalr_service.example.id
  sub_resource_name  = "vault"
  target_resource_id = azurerm_key_vault.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Signalr Shared Private Link Resource. Changing this forces a new resource to be created.

* `signalr_service_id` - (Required) The id of the Signalr Service. Changing this forces a new resource to be created.

* `sub_resource_name` - (Required) The sub resource name which the Signalr Private Endpoint can connect to. Possible values are `sites`, `vault`. Changing this forces a new resource to be created.

* `target_resource_id` - (Required) The ID of the Shared Private Link Enabled Remote Resource which this Signalr Private Endpoint should be connected to. Changing this forces a new resource to be created.

-> **Note:** The `sub_resource_name` should match with the type of the `target_resource_id` that's being specified.

* `request_message` - (Optional) The request message for requesting approval of the Shared Private Link Enabled Remote Resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Signalr Shared Private Link resource.

* `status` - The status of a private endpoint connection. Possible values are `Pending`, `Approved`, `Rejected` or `Disconnected`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Signalr Shared Private Link Resource.
* `read` - (Defaults to 5 minutes) Used when retrieving the Signalr Shared Private Link Resource.
* `update` - (Defaults to 30 minutes) Used when updating the Signalr Shared Private Link Resource.
* `delete` - (Defaults to 30 minutes) Used when deleting the Signalr Shared Private Link Resource.

## Import

Signalr Shared Private Link Resource can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_signalr_shared_private_link_resource.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.SignalRService/signalR/signalr1/sharedPrivateLinkResources/resource1
```
