---
subcategory: "Web Pubsub"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_web_pubsub"
description: |-
  Manages an Azure Web Pubsub service.
---

# azurerm_web_pubsub

Manages an Azure Web Pubsub Service.

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

  sku {
    name     = "Standard_S1"
    capactiy = 1
  }

  public_network_access = "Disabled"

  live_trace_configuration {
    enabled = "false"
    categories {
      name    = "MessagingLogs"
      enabled = "true"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Web Pubsub service. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Web Pubsub service. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the Web Pubsub service exists. Changing this forces a new resource to be created.

* `sku` - (Required) A `sku` block as documented below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `live_trace_configuration` - (Optional) A `live_trace_configuration` block as documented below

* `disable_local_auth` - (Optional) A `boolean` type. Specifies whether to disable local auth, defaults to `false`

* `disable_aad_auth` - (Optional) A `boolean` type. Specifies whether to disable AAD auth, defaults to `false`

* `tls_client_cert_enabled` - (Optional)  Specifies whether to request client certificate during TLS handshake if enabled

---

A `sku` block supports the following:

* `name` - (Required) Specifies which tier to use. Valid values are `Free_F1` and `Standard_S1`.

* `capacity` - (Required) Specifies the number of units associated with this Web Pubsub resource. Valid values are `1`, `2`, `5`, `10`, `20`, `50` and `100`.

---

A `live_trace_configuration` block supports the following:

* `enabled` - (Optional) Is this live trace enabled? Defaults to `true`. Possible values are `true`, `false`.

-> **NOTE:** Indicates whether or not enable live trace. When it's set to true, live trace client can connect to the service. Otherwise, live trace client can't connect to the service, so that you are unable to receive any log, no matter what you configure in `categories`.

* `categories` - (Optional) block as documented below.

---

A `categories` block supports the following:

* `name` - (Required) The name of the Log Category for this Resource. Possible values are `ConnectivityLogs`, `MessagingLogs` and `HttpRequestLogs`.

* `enabled` - (Optional) Is this log category enabled? Defaults to `true`. Possible values are `true`, `false`.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Web Pubsub service.

* `hostname` - The FQDN of the Web Pubsub service.

* `ip_address` - The publicly accessible IP of the Web Pubsub service.

* `public_port` - The publicly accessible port of the Web Pubsub service which is designed for browser/client use.

* `server_port` - The publicly accessible port of the Web Pubsub service which is designed for customer server side use.

* `primary_access_key` - The primary access key for the Web Pubsub service.

* `primary_connection_string` - The primary connection string for the Web Pubsub service.

* `secondary_access_key` - The secondary access key for the Web Pubsub service.

* `secondary_connection_string` - The secondary connection string for the Web Pubsub service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Web Pubsub Service.
* `update` - (Defaults to 30 minutes) Used when updating the Web Pubsub Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Web Pubsub Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the Web Pubsub Service.

## Import

Web Pubsub services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_web_pubsub.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/terraform-wps/providers/Microsoft.SignalRService/webPubSub/tfex-webpubsub
```

