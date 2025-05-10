---
subcategory: "Active Directory Domain Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_active_directory_domain_service"
description: |-
  Manages an Active Directory Domain Service.
---

# azurerm_active_directory_domain_service

Manages an Active Directory Domain Service.

~> **Note:** Before using this resource, there must exist in your tenant a service principal for the Domain Services published application. This service principal cannot be easily managed by Terraform and it's recommended to create this manually, as it does not exist by default. See [official documentation](https://docs.microsoft.com/azure/active-directory-domain-services/powershell-create-instance#create-required-azure-ad-resources) for details.

-> **Note:** At present this resource only supports **User Forest** mode and _not_ **Resource Forest** mode. [Read more](https://docs.microsoft.com/azure/active-directory-domain-services/concepts-resource-forest) about the different operation modes for this service.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "deploy" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "deploy" {
  name                = "deploy-vnet"
  location            = azurerm_resource_group.deploy.location
  resource_group_name = azurerm_resource_group.deploy.name
  address_space       = ["10.0.1.0/16"]
}

resource "azurerm_subnet" "deploy" {
  name                 = "deploy-subnet"
  resource_group_name  = azurerm_resource_group.deploy.name
  virtual_network_name = azurerm_virtual_network.deploy.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_network_security_group" "deploy" {
  name                = "deploy-nsg"
  location            = azurerm_resource_group.deploy.location
  resource_group_name = azurerm_resource_group.deploy.name

  security_rule {
    name                       = "AllowSyncWithAzureAD"
    priority                   = 101
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "443"
    source_address_prefix      = "AzureActiveDirectoryDomainServices"
    destination_address_prefix = "*"
  }

  security_rule {
    name                       = "AllowRD"
    priority                   = 201
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "3389"
    source_address_prefix      = "CorpNetSaw"
    destination_address_prefix = "*"
  }

  security_rule {
    name                       = "AllowPSRemoting"
    priority                   = 301
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "5986"
    source_address_prefix      = "AzureActiveDirectoryDomainServices"
    destination_address_prefix = "*"
  }

  security_rule {
    name                       = "AllowLDAPS"
    priority                   = 401
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "636"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}

resource "azurerm_subnet_network_security_group_association" "deploy" {
  subnet_id                 = azurerm_subnet.deploy.id
  network_security_group_id = azurerm_network_security_group.deploy.id
}

resource "azuread_group" "dc_admins" {
  display_name     = "AAD DC Administrators"
  security_enabled = true
}

resource "azuread_user" "admin" {
  user_principal_name = "dc-admin@hashicorp-example.com"
  display_name        = "DC Administrator"
  password            = "Pa55w0Rd!!1"
}

resource "azuread_group_member" "admin" {
  group_object_id  = azuread_group.dc_admins.object_id
  member_object_id = azuread_user.admin.object_id
}

resource "azuread_service_principal" "example" {
  application_id = "2565bd9d-da50-47d4-8b85-4c97f669dc36" // published app for domain services
}

resource "azurerm_resource_group" "aadds" {
  name     = "aadds-rg"
  location = "westeurope"
}

resource "azurerm_active_directory_domain_service" "example" {
  name                = "example-aadds"
  location            = azurerm_resource_group.aadds.location
  resource_group_name = azurerm_resource_group.aadds.name

  domain_name           = "widgetslogin.net"
  sku                   = "Enterprise"
  filtered_sync_enabled = false

  initial_replica_set {
    subnet_id = azurerm_subnet.deploy.id
  }

  notifications {
    additional_recipients = ["notifyA@example.net", "notifyB@example.org"]
    notify_dc_admins      = true
    notify_global_admins  = true
  }

  security {
    sync_kerberos_passwords = true
    sync_ntlm_passwords     = true
    sync_on_prem_passwords  = true
  }

  tags = {
    Environment = "prod"
  }

  depends_on = [
    azuread_service_principal.example,
    azurerm_subnet_network_security_group_association.deploy,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required) The Active Directory domain to use. See [official documentation](https://docs.microsoft.com/azure/active-directory-domain-services/tutorial-create-instance#create-a-managed-domain) for constraints and recommendations. Changing this forces a new resource to be created.

* `domain_configuration_type` - (Optional) The configuration type of this Active Directory Domain. Possible values are `FullySynced` and `ResourceTrusting`. Changing this forces a new resource to be created.

* `filtered_sync_enabled` - (Optional) Whether to enable group-based filtered sync (also called scoped synchronisation). Defaults to `false`.

* `secure_ldap` - (Optional) A `secure_ldap` block as defined below.

* `location` - (Required) The Azure location where the Domain Service exists. Changing this forces a new resource to be created.

* `name` - (Required) The display name for your managed Active Directory Domain Service resource. Changing this forces a new resource to be created.

* `notifications` - (Optional) A `notifications` block as defined below.

* `initial_replica_set` - (Required) An `initial_replica_set` block as defined below. The initial replica set inherits the same location as the Domain Service resource.

* `resource_group_name` - (Required) The name of the Resource Group in which the Domain Service should exist. Changing this forces a new resource to be created.

* `security` - (Optional) A `security` block as defined below.

* `sku` - (Required) The SKU to use when provisioning the Domain Service resource. One of `Standard`, `Enterprise` or `Premium`.

* `tags` - (Optional) A mapping of tags assigned to the resource.

---

A `secure_ldap` block supports the following:

* `enabled` - (Required) Whether to enable secure LDAP for the managed domain. For more information, please see [official documentation on enabling LDAPS](https://docs.microsoft.com/azure/active-directory-domain-services/tutorial-configure-ldaps), paying particular attention to the section on network security to avoid unnecessarily exposing your service to Internet-borne bruteforce attacks.

* `external_access_enabled` - (Optional) Whether to enable external access to LDAPS over the Internet. Defaults to `false`.

* `pfx_certificate` - (Required) The certificate/private key to use for LDAPS, as a base64-encoded TripleDES-SHA1 encrypted PKCS#12 bundle (PFX file).

* `pfx_certificate_password` - (Required) The password to use for decrypting the PKCS#12 bundle (PFX file).

---

A `notifications` block supports the following:

* `additional_recipients` - (Optional) A list of additional email addresses to notify when there are alerts in the managed domain.

* `notify_dc_admins` - (Optional) Whether to notify members of the _AAD DC Administrators_ group when there are alerts in the managed domain.

* `notify_global_admins` - (Optional) Whether to notify all Global Administrators when there are alerts in the managed domain.

---

An `initial_replica_set` block supports the following:

* `subnet_id` - (Required) The ID of the subnet in which to place the initial replica set. Changing this forces a new resource to be created.

---

A `security` block supports the following:

* `kerberos_armoring_enabled` - (Optional) Whether to enable Kerberos Armoring. Defaults to `false`.

* `kerberos_rc4_encryption_enabled` - (Optional) Whether to enable Kerberos RC4 Encryption. Defaults to `false`.

* `ntlm_v1_enabled` - (Optional) Whether to enable legacy NTLM v1 support. Defaults to `false`.

* `sync_kerberos_passwords` - (Optional) Whether to synchronize Kerberos password hashes to the managed domain. Defaults to `false`.

* `sync_ntlm_passwords` - (Optional) Whether to synchronize NTLM password hashes to the managed domain. Defaults to `false`.

* `sync_on_prem_passwords` - (Optional) Whether to synchronize on-premises password hashes to the managed domain. Defaults to `false`.

* `tls_v1_enabled` - (Optional) Whether to enable legacy TLS v1 support. Defaults to `false`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Domain Service.
  
* `deployment_id` - A unique ID for the managed domain deployment.

* `resource_id` - The Azure resource ID for the domain service.

---

A `secure_ldap` block exports the following:

* `certificate_expiry` - The expiry time of the certificate.

* `certificate_thumbprint` - The thumbprint of the certificate.

* `public_certificate` - The public certificate.

---

An `initial_replica_set` block exports the following:

* `domain_controller_ip_addresses` - A list of subnet IP addresses for the domain controllers in the initial replica set, typically two.

* `external_access_ip_address` - The publicly routable IP address for the domain controllers in the initial replica set.

* `location` - The Azure location in which the initialreplica set resides.

* `id` - A unique ID for the replica set.

* `service_status` - The current service status for the initial replica set.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the Domain Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Domain Service.
* `update` - (Defaults to 2 hours) Used when updating the Domain Service.
* `delete` - (Defaults to 1 hour) Used when deleting the Domain Service.

## Import

Domain Services can be imported using the resource ID, together with the Replica Set ID that you wish to designate as the initial replica set, e.g.

```shell
terraform import azurerm_active_directory_domain_service.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.AAD/domainServices/instance1/initialReplicaSetId/00000000-0000-0000-0000-000000000000
```
