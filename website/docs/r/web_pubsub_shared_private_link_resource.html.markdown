---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_web_pubsub_shared_private_link_resource"
description: |-
  Manages the Shared Private Link Resource for a Web Pubsub service.
---

# azurerm_web_pubsub_shared_private_link_resource

Manages the Shared Private Link Resource for a Web Pubsub service.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "terraform-webpubsub"
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
      "managecontacts",
    ]

    key_permissions = [
      "create",
    ]

    secret_permissions = [
      "set",
    ]
  }
}

resource "azurerm_web_pubsub" "example" {
  name                = "tfex-webpubsub"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard_S1"
  capacity            = 1
}

resource "azurerm_web_pubsub_shared_private_link_resource" "example" {
  name               = "tfex-webpubsub-splr"
  web_pubsub_id      = azurerm_web_pubsub.example.id
  subresource_name   = "vault"
  target_resource_id = azurerm_key_vault.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specify the name of the Web Pubsub Shared Private Link Resource. Changing this forces a new resource to be created.

* `web_pubsub_id` - (Required) Specify the id of the Web Pubsub. Changing this forces a new resource to be created.

* `subresource_name` - (Required) Specify the sub resource name which the Web Pubsub Private Endpoint is able to connect to. Changing this forces a new resource to be created.

-> **Note:** The available sub resource can be retrieved by using `azurerm_web_pubsub_private_link_resource` data source.

* `target_resource_id` - (Required) Specify the ID of the Shared Private Link Enabled Remote Resource which this Web Pubsub Private Endpoint should be connected to. Changing this forces a new resource to be created.

-> **Note:** The sub resource name should match with the type of the target resource id that's being specified.

* `request_message` - (Optional) Specify the request message for requesting approval of the Shared Private Link Enabled Remote Resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Web Pubsub Shared Private Link resource.

* `status` - The status of a private endpoint connection. Possible values are Pending, Approved, Rejected or Disconnected.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Web Pubsub Shared Private Link Resource.
* `read` - (Defaults to 5 minutes) Used when retrieving the Web Pubsub Shared Private Link Resource.
* `update` - (Defaults to 30 minutes) Used when updating the Web Pubsub Shared Private Link Resource.
* `delete` - (Defaults to 30 minutes) Used when deleting the Web Pubsub Shared Private Link Resource.

## Import

Web Pubsub Shared Private Link Resource can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_web_pubsub_shared_private_link_resource.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.SignalRService/webPubSub/webPubSub1/sharedPrivateLinkResources/resource1
```
