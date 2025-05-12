---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_web_pubsub"
description: |-
  Manages an Azure Web PubSub service.
---

# azurerm_web_pubsub

Manages an Azure Web PubSub Service.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "terraform-webpubsub"
  location = "east us"
}

resource "azurerm_web_pubsub" "example" {
  name                = "tfex-webpubsub"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku      = "Standard_S1"
  capacity = 1

  public_network_access_enabled = false

  live_trace {
    enabled                   = true
    messaging_logs_enabled    = true
    connectivity_logs_enabled = false
  }

  identity {
    type = "SystemAssigned"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Web PubSub service. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Web PubSub service. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the Web PubSub service exists. Changing this forces a new resource to be created.

* `sku` - (Required) Specifies which SKU to use. Possible values are `Free_F1`, `Standard_S1`, `Premium_P1` and `Premium_P2`.

* `capacity` - (Optional) Specifies the number of units associated with this Web PubSub resource. Valid values are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `20`, `30`, `40`, `50`, `60`, `70`, `80`, `90`, `100`, `200`, `300`, `400`, `500`, `600`, `700`, `800`, `900` and `1000`.

~> **Note:** The valid capacity range for sku `Free_F1` is `1`, for sku `Premium_P2` is from `100` to `1000`, and from `1` to `100` for sku `Standard_S1` and `Premium_P1`.

* `public_network_access_enabled` - (Optional) Whether to enable public network access? Defaults to `true`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `live_trace` - (Optional) A `live_trace` block as defined below.

* `identity` - (Optional) An `identity` block as defined below.

* `local_auth_enabled` - (Optional) Whether to enable local auth? Defaults to `true`.

* `aad_auth_enabled` - (Optional) Whether to enable AAD auth? Defaults to `true`.

* `tls_client_cert_enabled` - (Optional) Whether to request client certificate during TLS handshake? Defaults to `false`.

---

A `live_trace` block supports the following:

* `enabled` - (Optional) Whether the live trace is enabled? Defaults to `true`.

* `messaging_logs_enabled` - (Optional) Whether the log category `MessagingLogs` is enabled? Defaults to `true`

* `connectivity_logs_enabled` - (Optional) Whether the log category `ConnectivityLogs` is enabled? Defaults to `true`

* `http_request_logs_enabled` - (Optional) Whether the log category `HttpRequestLogs` is enabled? Defaults to `true`

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Web PubSub. Possible values are `SystemAssigned`, `UserAssigned`.

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Web PubSub.

~> **Note:** This is required when `type` is set to `UserAssigned`

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Web PubSub service.

* `hostname` - The FQDN of the Web PubSub service.

* `identity` - An `identity` block as defined below.

* `external_ip` - The publicly accessible IP of the Web PubSub service.

* `public_port` - The publicly accessible port of the Web PubSub service which is designed for browser/client use.

* `server_port` - The publicly accessible port of the Web PubSub service which is designed for customer server side use.

* `primary_access_key` - The primary access key for the Web PubSub service.

* `primary_connection_string` - The primary connection string for the Web PubSub service.

* `secondary_access_key` - The secondary access key for the Web PubSub service.

* `secondary_connection_string` - The secondary connection string for the Web PubSub service.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to
specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Web PubSub Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Web PubSub Service.
* `update` - (Defaults to 30 minutes) Used when updating the Web PubSub Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the Web PubSub Service.

## Import

Web PubSub services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_web_pubsub.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.SignalRService/webPubSub/pubsub1
```
