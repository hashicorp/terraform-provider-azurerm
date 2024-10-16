---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_vpn_server_configuration"
description: |-
  Gets information about an existing VPN Server Configuration.
---

# Data Source: azurerm_vpn_server_configuration

Use this data source to access information about an existing VPN Server Configuration.

## Example Usage

```hcl
data "azurerm_vpn_server_configuration" "example" {
  name                = "existing-local-vpn-server-configuration"
  resource_group_name = "existing-resource-group"
}

output "azurerm_vpn_server_configuration" {
  value = data.azurerm_vpn_server_configuration.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Name of the VPN Server Configuration.

* `resource_group_name` - (Required) The name of the Resource Group where the VPN Server Configuration exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the VPN Server Configuration.

* `location` - The Azure Region where the VPN Server Configuration exists.

* `vpn_authentication_types` -  The list of Authentication Types applicable for the VPN Server Configuration.

* `ipsec_policy` - The `bgp_settings` block as defined below.

* `vpn_protocols` -  The list of VPN Protocols to use for the VPN Server Configuration.

* `tags` - A mapping of tags to assign to the VPN Server Configuration.

---

When `vpn_authentication_types` contains `AAD` the following arguments are exported:

* `azure_active_directory_authentication` - A `azure_active_directory_authentication` block as defined below.

---

When `vpn_authentication_types` contains `Certificate` the following arguments are supported:

* `client_root_certificate` - One or more `client_root_certificate` blocks as defined below.

* `client_revoked_certificate` - One or more `client_revoked_certificate` blocks as defined below.

---

When `vpn_authentication_types` contains `Radius` the following arguments are exported:

* `radius` - A `radius` block as defined below.

A `azure_active_directory_authentication` block exports the following:

* `audience` - The Audience which should be used for authentication.

* `issuer` - The Issuer which should be used for authentication.

* `tenant` - The Tenant which should be used for authentication.

---

A `client_revoked_certificate` block exports the following:

* `name` - The name used to uniquely identify this certificate.

* `thumbprint` - The Thumbprint of the Certificate.

---

A `client_root_certificate` block at the root of the resource exports the following:

* `name` - The name used to uniquely identify this certificate.

* `public_cert_data` - The Public Key Data associated with the Certificate.

---

A `ipsec_policy` block exports the following:

* `dh_group` - The DH Group, used in IKE Phase 1.

* `ike_encryption` - The IKE encryption algorithm, used for IKE Phase 2.

* `ike_integrity` - The IKE encryption integrity algorithm, used for IKE Phase 2.

* `ipsec_encryption` - The IPSec encryption algorithm, used for IKE phase 1.

* `ipsec_integrity` - The IPSec integrity algorithm, used for IKE phase 1.

* `pfs_group` - The Pfs Group, used in IKE Phase 2.

* `sa_lifetime_seconds` - The IPSec Security Association lifetime in seconds for a Site-to-Site VPN tunnel.

* `sa_data_size_kilobytes` - The IPSec Security Association payload size in KB for a Site-to-Site VPN tunnel.

---

A `radius` block exports the following:

* `server` - One or more `server` blocks as defined below.

* `client_root_certificate` - One or more `client_root_certificate` blocks as defined below.

* `server_root_certificate` - One or more `server_root_certificate` blocks as defined below.

---

A `server` nested within the `radius` block exports the following::

* `address` - The Address of the Radius Server.

* `secret` - The Secret used to communicate with the Radius Server.

* `score` - The Score of the Radius Server determines the priority of the server.

---

A `client_root_certificate` block nested within the `radius` block exports the following:

* `name` - The name used to uniquely identify this certificate.

* `thumbprint` - The Thumbprint of the Certificate.

---

A `server_root_certificate` block exports the following:

* `name` - The name used to uniquely identify this certificate.

* `public_cert_data` - The Public Key Data associated with the Certificate.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the VPN Server Configuration.
