---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_vpn_server_configuration_policy_group"
description: |-
  Manages a VPN Server Configuration Policy Group.
---

# azurerm_vpn_server_configuration_policy_group

Manages a VPN Server Configuration Policy Group.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_vpn_server_configuration" "test" {
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

resource "azurerm_vpn_server_configuration_policy_group" "example" {
  name                        = "example-pg"
  resource_group_name         = azurerm_resource_group.example.name
  vpn_server_configuration_id = azurerm_network_vpn_server_configuration.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Name which should be used for this VPN Server Configuration Policy Group. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The Name of the Resource Group in which this VPN Server Configuration Policy Group should be created. Changing this forces a new resource to be created.

* `vpn_server_configuration_id` - (Required) The ID of the VPN Server Configuration that this VPN Server Configuration Policy Group belongs to. Changing this forces a new resource to be created.

* `is_default` - (Optional) Is this a default VPN Server Configuration Policy Group? 

* `policy_member` - (Optional) One or more `policy_member` blocks as documented below.

* `priority` - (Optional) The priority for this VPN Server Configuration Policy Group.

---

A `policy_member` block supports the following:

* `name` - (Optional) The name of the VPN Server Configuration Policy Group member.

* `attribute_type` - (Optional) The attribute type of the VPN Server Configuration Policy Group member. Possible values are `AADGroupId`, `CertificateGroupId` and `RadiusAzureGroupId`.

* `attribute_value` - (Optional) The value of the attribute that is used for the VPN Server Configuration Policy Group member.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the VPN Server Configuration Policy Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the VPN Server Configuration Policy Group.
* `update` - (Defaults to 30 minutes) Used when updating the VPN Server Configuration Policy Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the VPN Server Configuration Policy Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the VPN Server Configuration Policy Group.

## Import

VPN Server Configuration Policy Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_vpn_server_configuration_policy_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Network/vpnServerConfigurations/serverConfiguration1/configurationPolicyGroups/configurationPolicyGroup1
```
