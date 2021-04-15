---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_signalr_service"
description: |-
  Manages an Azure SignalR service.
---

# azurerm_signalr_service

Manages an Azure SignalR service.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "terraform-signalr"
  location = "West Europe"
}

resource "azurerm_signalr_service" "example" {
  name                = "tfex-signalr"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku {
    name     = "Free_F1"
    capacity = 1
  }

  cors {
    allowed_origins = ["http://www.example.com"]
  }

  features {
    flag  = "ServiceMode"
    value = "Default"
  }

  upstream_endpoint {
    category_pattern = ["connections", "messages"]
    event_pattern    = ["*"]
    hub_pattern      = ["hub1"]
    url_template     = "http://foo.com"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the SignalR service. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the SignalR service. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the SignalR service exists. Changing this forces a new resource to be created.

* `sku` - A `sku` block as documented below.

* `cors` - (Optional) A `cors` block as documented below.

* `features` - (Optional) A `features` block as documented below.

* `upstream_endpoint` - (Optional) An `upstream_endpoint` block as documented below. Using this block requires the SignalR service to be Serverless. When creating multiple blocks they will be processed in the order they are defined in.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `cors` block supports the following:

* `allowed_origins` - (Required) A list of origins which should be able to make cross-origin calls. `*` can be used to allow all calls.

---

A `features` block supports the following:

* `flag` - (Required) The kind of Feature. Possible values are `EnableConnectivityLogs`, `EnableMessagingLogs`, and `ServiceMode`.

* `value` - (Required) A value of a feature flag. Possible values are `Classic`, `Default` and `Serverless`.

---

An `upstream_endpoint` block supports the following:

* `url_template` - (Required) The upstream URL Template. This can be a url or a template such as `http://host.com/{hub}/api/{category}/{event}`.

* `category_pattern` - (Optional) The categories to match on, or `*` for all.

* `event_pattern` - (Optional) The events to match on, or `*` for all.

* `hub_pattern` - (Optional) The hubs to match on, or `*` for all.

---

A `sku` block supports the following:

* `name` - (Required) Specifies which tier to use. Valid values are `Free_F1` and `Standard_S1`.

* `capacity` - (Required) Specifies the number of units associated with this SignalR service. Valid values are `1`, `2`, `5`, `10`, `20`, `50` and `100`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the SignalR service.

* `hostname` - The FQDN of the SignalR service.

* `ip_address` - The publicly accessible IP of the SignalR service.

* `public_port` - The publicly accessible port of the SignalR service which is designed for browser/client use.

* `server_port` - The publicly accessible port of the SignalR service which is designed for customer server side use.

* `primary_access_key` - The primary access key for the SignalR service.

* `primary_connection_string` - The primary connection string for the SignalR service.

* `secondary_access_key` - The secondary access key for the SignalR service.

* `secondary_connection_string` - The secondary connection string for the SignalR service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the SignalR Service.
* `update` - (Defaults to 30 minutes) Used when updating the SignalR Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the SignalR Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the SignalR Service.

## Import

SignalR services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_signalr_service.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/terraform-signalr/providers/Microsoft.SignalRService/SignalR/tfex-signalr
```
