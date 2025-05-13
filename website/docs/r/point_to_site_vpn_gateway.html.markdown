---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_point_to_site_vpn_gateway"
description: |-
  Manages a Point-to-Site VPN Gateway.

---

# azurerm_point_to_site_vpn_gateway

Manages a Point-to-Site VPN Gateway.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_wan" "example" {
  name                = "example-virtualwan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_virtual_hub" "example" {
  name                = "example-virtualhub"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  virtual_wan_id      = azurerm_virtual_wan.example.id
  address_prefix      = "10.0.0.0/23"
}

resource "azurerm_vpn_server_configuration" "example" {
  name                     = "example-config"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  vpn_authentication_types = ["Certificate"]

  client_root_certificate {
    name             = "DigiCert-Federated-ID-Root-CA"
    public_cert_data = <<EOF
MIIDuzCCAqOgAwIBAgIQCHTZWCM+IlfFIRXIvyKSrjANBgkqhkiG9w0BAQsFADBn
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
d3cuZGlnaWNlcnQuY29tMSYwJAYDVQQDEx1EaWdpQ2VydCBGZWRlcmF0ZWQgSUQg
Um9vdCBDQTAeFw0xMzAxMTUxMjAwMDBaFw0zMzAxMTUxMjAwMDBaMGcxCzAJBgNV
BAYTAlVTMRUwEwYDVQQKEwxEaWdpQ2VydCBJbmMxGTAXBgNVBAsTEHd3dy5kaWdp
Y2VydC5jb20xJjAkBgNVBAMTHURpZ2lDZXJ0IEZlZGVyYXRlZCBJRCBSb290IENB
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvAEB4pcCqnNNOWE6Ur5j
QPUH+1y1F9KdHTRSza6k5iDlXq1kGS1qAkuKtw9JsiNRrjltmFnzMZRBbX8Tlfl8
zAhBmb6dDduDGED01kBsTkgywYPxXVTKec0WxYEEF0oMn4wSYNl0lt2eJAKHXjNf
GTwiibdP8CUR2ghSM2sUTI8Nt1Omfc4SMHhGhYD64uJMbX98THQ/4LMGuYegou+d
GTiahfHtjn7AboSEknwAMJHCh5RlYZZ6B1O4QbKJ+34Q0eKgnI3X6Vc9u0zf6DH8
Dk+4zQDYRRTqTnVO3VT8jzqDlCRuNtq6YvryOWN74/dq8LQhUnXHvFyrsdMaE1X2
DwIDAQABo2MwYTAPBgNVHRMBAf8EBTADAQH/MA4GA1UdDwEB/wQEAwIBhjAdBgNV
HQ4EFgQUGRdkFnbGt1EWjKwbUne+5OaZvRYwHwYDVR0jBBgwFoAUGRdkFnbGt1EW
jKwbUne+5OaZvRYwDQYJKoZIhvcNAQELBQADggEBAHcqsHkrjpESqfuVTRiptJfP
9JbdtWqRTmOf6uJi2c8YVqI6XlKXsD8C1dUUaaHKLUJzvKiazibVuBwMIT84AyqR
QELn3e0BtgEymEygMU569b01ZPxoFSnNXc7qDZBDef8WfqAV/sxkTi8L9BkmFYfL
uGLOhRJOFprPdoDIUBB+tmCl3oDcBy3vnUeOEioz8zAkprcb3GHwHAK+vHmmfgcn
WsfMLH4JCLa/tRYL+Rw/N3ybCkDp00s0WUZ+AoDywSl0Q/ZEnNY0MsFiw6LyIdbq
M/s/1JRtO3bDSzD9TazRVzn2oBqzSa8VgIo5C1nOnoAKJTlsClJKvIhnRlaLQqk=
EOF
  }
}

resource "azurerm_point_to_site_vpn_gateway" "example" {
  name                        = "example-vpn-gateway"
  location                    = azurerm_resource_group.example.location
  resource_group_name         = azurerm_resource_group.example.name
  virtual_hub_id              = azurerm_virtual_hub.example.id
  vpn_server_configuration_id = azurerm_vpn_server_configuration.example.id
  scale_unit                  = 1
  connection_configuration {
    name = "example-gateway-config"

    vpn_client_address_pool {
      address_prefixes = [
        "10.0.2.0/24"
      ]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Point-to-Site VPN Gateway. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Point-to-Site VPN Gateway. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `connection_configuration` - (Required) A `connection_configuration` block as defined below.

* `scale_unit` - (Required) The [Scale Unit](https://docs.microsoft.com/azure/virtual-wan/virtual-wan-faq#what-is-a-virtual-wan-gateway-scale-unit) for this Point-to-Site VPN Gateway.

* `virtual_hub_id` - (Required) The ID of the Virtual Hub where this Point-to-Site VPN Gateway should exist. Changing this forces a new resource to be created.

* `vpn_server_configuration_id` - (Required) The ID of the VPN Server Configuration which this Point-to-Site VPN Gateway should use. Changing this forces a new resource to be created.

* `dns_servers` - (Optional) A list of IP Addresses of DNS Servers for the Point-to-Site VPN Gateway.

* `routing_preference_internet_enabled` - (Optional) Is the Routing Preference for the Public IP Interface of the VPN Gateway enabled? Defaults to `false`. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the Point-to-Site VPN Gateway.

---

A `connection_configuration` block supports the following:

* `name` - (Required) The Name which should be used for this Connection Configuration.

* `vpn_client_address_pool` - (Required) A `vpn_client_address_pool` block as defined below.

* `route` - (Optional) A `route` block as defined below.

* `internet_security_enabled` - (Optional) Should Internet Security be enabled to secure internet traffic? Changing this forces a new resource to be created. Defaults to `false`.

---

A `vpn_client_address_pool` block supports the following:

* `address_prefixes` - (Required) A list of CIDR Ranges which should be used as Address Prefixes.

---

A `route` block supports the following:

* `associated_route_table_id` - (Required) The Virtual Hub Route Table resource id associated with this Routing Configuration.

* `inbound_route_map_id` - (Optional) The resource ID of the Route Map associated with this Routing Configuration for inbound learned routes.

* `outbound_route_map_id` - (Optional) The resource ID of the Route Map associated with this Routing Configuration for outbound advertised routes.

* `propagated_route_table` - (Optional) A `propagated_route_table` block as defined below.

---

A `propagated_route_table` block supports the following:

* `ids` - (Required) The list of Virtual Hub Route Table resource id which the routes will be propagated to.

* `labels` - (Optional) The list of labels to logically group Virtual Hub Route Tables which the routes will be propagated to.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Point-to-Site VPN Gateway.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the Point-to-Site VPN Gateway.
* `read` - (Defaults to 5 minutes) Used when retrieving the Point-to-Site VPN Gateway.
* `update` - (Defaults to 90 minutes) Used when updating the Point-to-Site VPN Gateway.
* `delete` - (Defaults to 90 minutes) Used when deleting the Point-to-Site VPN Gateway.

## Import

Point-to-Site VPN Gateway's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_point_to_site_vpn_gateway.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/p2sVpnGateways/gateway1
```
