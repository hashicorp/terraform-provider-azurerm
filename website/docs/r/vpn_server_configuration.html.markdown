---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_vpn_server_configuration"
description: |-
    Manages a VPN Server Configuration.
---

# azurerm_vpn_server_configuration

Manages a VPN Server Configuration.

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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Name which should be used for this VPN Server Configuration. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The Name of the Resource Group in which this VPN Server Configuration should be created. Changing this forces a new resource to be created.

* `location` - (Required) The Azure location where this VPN Server Configuration should be created. Changing this forces a new resource to be created.

* `vpn_authentication_types` - (Required) A list of one of more Authentication Types applicable for this VPN Server Configuration. Possible values are `AAD` (Azure Active Directory), `Certificate` and `Radius`.

-> **NOTE:** At this time a maximum of one VPN Authentication Types can be specified.

---

* `ipsec_policy` - (Optional) A `ipsec_policy` block as defined below.

* `vpn_protocols` - (Optional) A list of VPN Protocols to use for this Server Configuration. Possible values are `IkeV2` and `OpenVPN`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

When `vpn_authentication_types` contains `AAD` the following arguments are supported:

* `azure_active_directory_authentication` - (Required) A `azure_active_directory_authentication` block as defined below.

---

When `vpn_authentication_types` contains `Certificate` the following arguments are supported:

* `client_root_certificate` - (Required) One or more `client_root_certificate` blocks as defined below.

* `client_revoked_certificate` - (Optional) One or more `client_revoked_certificate` blocks as defined below.

---

When `vpn_authentication_types` contains `Radius` the following arguments are supported:

* `radius_server` - (Optional / **Deprecated**) A `radius_server` block as defined below.
* `radius` - (Optional) A `radius` block as defined below.

---

A `azure_active_directory_authentication` block supports the following:

* `audience` - (Required) The Audience which should be used for authentication.

* `issuer` - (Required) The Issuer which should be used for authentication.

* `tenant` - (Required) The Tenant which should be used for authentication.

---

A `client_revoked_certificate` block supports the following:

* `name` - (Required) A name used to uniquely identify this certificate.

* `thumbprint` - (Required) The Thumbprint of the Certificate.

---

A `client_root_certificate` block at the root of the resource supports the following:

* `name` - (Required) A name used to uniquely identify this certificate.

* `public_cert_data` - (Required) The Public Key Data associated with the Certificate.

---

A `client_root_certificate` block nested within the `radius_server` block supports the following:

* `name` - (Required) A name used to uniquely identify this certificate.

* `thumbprint` - (Required) The Thumbprint of the Certificate.

---

A `ipsec_policy` block supports the following:

* `dh_group` - (Required) The DH Group, used in IKE Phase 1. Possible values include `DHGroup1`, `DHGroup2`, `DHGroup14`, `DHGroup24`, `DHGroup2048`, `ECP256`, `ECP384` and `None`.

* `ike_encryption` - (Required) The IKE encryption algorithm, used for IKE Phase 2. Possible values include `AES128`, `AES192`, `AES256`, `DES`, `DES3`, `GCMAES128` and `GCMAES256`.

* `ike_integrity` - (Required) The IKE encryption integrity algorithm, used for IKE Phase 2. Possible values include `GCMAES128`, `GCMAES256`, `MD5`, `SHA1`, `SHA256` and `SHA384`.

* `ipsec_encryption` - (Required) The IPSec encryption algorithm, used for IKE phase 1. Possible values include `AES128`, `AES192`, `AES256`, `DES`, `DES3`, `GCMAES128`, `GCMAES192`, `GCMAES256` and `None`.

* `ipsec_integrity` - (Required) The IPSec integrity algorithm, used for IKE phase 1. Possible values include `GCMAES128`, `GCMAES192`, `GCMAES256`, `MD5`, `SHA1` and `SHA256`.

* `pfs_group` - (Required) The Pfs Group, used in IKE Phase 2. Possible values include `ECP256`, `ECP384`, `PFS1`, `PFS2`, `PFS14`, `PFS24`, `PFS2048`, `PFSMM` and `None`.

* `sa_lifetime_seconds` - (Required) The IPSec Security Association lifetime in seconds for a Site-to-Site VPN tunnel.

* `sa_data_size_kilobytes` - (Required) The IPSec Security Association payload size in KB for a Site-to-Site VPN tunnel.

---

A `radius_server` (**Deprecated) Use it to configure single Radius Server. The block supports the following:

* `address` - (Required) The Address of the Radius Server.

* `secret` - (Required) The Secret used to communicate with the Radius Server. Us

* `client_root_certificate` - (Optional) One or more `client_root_certificate` blocks as defined above.

* `server_root_certificate` - (Required) One or more `server_root_certificate` blocks as defined below.

---

A `radius` The block is used to configure single Radius Server. The block supports the following:

* `server` - (Required) One or more `server` blocks as defined below.

* `client_root_certificate` - (Optional) One or more `client_root_certificate` blocks as defined above.

* `server_root_certificate` - (Required) One or more `server_root_certificate` blocks as defined below.

---

A `server` block supports the following:

* `address` - (Required) The Address of the Radius Server.

* `secret` - (Required) The Secret used to communicate with the Radius Server.

* `score` - (Required) The score of the Radius Server determines the priority of the server. Ranges from 1 to 30.

---

A `server_root_certificate` block supports the following:

* `name` - (Required) A name used to uniquely identify this certificate.

* `public_cert_data` - (Required) The Public Key Data associated with the Certificate.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the VPN Server Configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the VPN Server Configuration.
* `update` - (Defaults to 90 minutes) Used when updating the VPN Server Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the VPN Server Configuration.
* `delete` - (Defaults to 90 minutes) Used when deleting the VPN Server Configuration.

## Import

VPN Server Configurations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_vpn_server_configuration.config1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/vpnServerConfigurations/config1
```
