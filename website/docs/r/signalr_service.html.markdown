---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_signalr_service"
sidebar_current: "docs-azurerm-resource-messaging-signalr"
description: |-
  Manages an Azure SignalR service.
---

# azurerm_signalr_service

Manages an Azure SignalR service.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "terraform-signalr"
  location = "West US"
}

resource "azurerm_signalr_service" "example" {
  name                = "tfex-signalr"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  sku {
    name     = "Free"
    capacity = 1
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the SignalR service. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the SignalR service. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the SignalR service exists. Changing this forces a new resource to be created.

* `sku` - A `sku` block as documented below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `sku` block supports the following:

* `name` - (Required) Specifies which tier to use. Valid values are `Free` and `Standard`.

* `capacity` - (Required) Specifies the number of units associated with this SignalR service. Valid values are `1`, `2`, `5`, `10`, `20`, `50` and `100`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the SignalR service.

* `hostname` - The FQDN of the SignalR service.

* `ip_address` - The publicly accessible IP of the SignalR service.

* `public_port` - The publicly accessible port of the SignalR service which is designed for browser/client use.

* `server_port` - The publicly accessible port of the SignalR service which is designed for customer server side use.

## Import

SignalR services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_signalr_service.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/terraform-signalr/providers/Microsoft.SignalRService/SignalR/tfex-signalr
```
