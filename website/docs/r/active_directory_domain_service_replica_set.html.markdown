---
subcategory: "Active Directory Domain Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_active_directory_domain_service_replica_set"
description: |-
  Manages a Replica Set for an Active Directory Domain Service.
---

# azurerm_active_directory_domain_service_replica_set

Manages a Replica Set for an Active Directory Domain Service.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "primary" {
  name     = "aadds-primary-rg"
  location = "West Europe"
}

resource "azurerm_virtual_network" "primary" {
  name                = "aadds-primary-vnet"
  location            = azurerm_resource_group.primary.location
  resource_group_name = azurerm_resource_group.primary.name
  address_space       = ["10.0.1.0/16"]
}

resource "azurerm_subnet" "primary" {
  name                 = "aadds-primary-subnet"
  resource_group_name  = azurerm_resource_group.primary.name
  virtual_network_name = azurerm_virtual_network.primary.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_network_security_group" "primary" {
  name                = "aadds-primary-nsg"
  location            = azurerm_resource_group.primary.location
  resource_group_name = azurerm_resource_group.primary.name

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

resource "azurerm_subnet_network_security_group_association" "primary" {
  subnet_id                 = azurerm_subnet.primary.id
  network_security_group_id = azurerm_network_security_group.primary.id
}

resource "azuread_group" "dc_admins" {
  display_name     = "aad-dc-administrators"
  security_enabled = true
}

resource "azuread_user" "admin" {
  user_principal_name = "dc-admin@hashicorp-example.net"
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
    location  = azurerm_virtual_network.primary.location
    subnet_id = azurerm_subnet.primary.id
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
    azurerm_subnet_network_security_group_association.primary,
  ]
}

resource "azurerm_resource_group" "replica" {
  name     = "aadds-replica-rg"
  location = "North Europe"
}

resource "azurerm_virtual_network" "replica" {
  name                = "aadds-replica-vnet"
  location            = azurerm_resource_group.replica.location
  resource_group_name = azurerm_resource_group.replica.name
  address_space       = ["10.20.0.0/16"]
}

resource "azurerm_subnet" "aadds_replica" {
  name                 = "aadds-replica-subnet"
  resource_group_name  = azurerm_resource_group.replica.name
  virtual_network_name = azurerm_virtual_network.replica.name
  address_prefixes     = ["10.20.0.0/24"]
}

resource "azurerm_network_security_group" "aadds_replica" {
  name                = "aadds-replica-nsg"
  location            = azurerm_resource_group.replica.location
  resource_group_name = azurerm_resource_group.replica.name

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

resource "azurerm_subnet_network_security_group_association" "replica" {
  subnet_id                 = azurerm_subnet.aadds_replica.id
  network_security_group_id = azurerm_network_security_group.aadds_replica.id
}

resource "azurerm_virtual_network_peering" "primary_replica" {
  name                      = "aadds-primary-replica"
  resource_group_name       = azurerm_virtual_network.primary.resource_group_name
  virtual_network_name      = azurerm_virtual_network.primary.name
  remote_virtual_network_id = azurerm_virtual_network.replica.id

  allow_forwarded_traffic      = true
  allow_gateway_transit        = false
  allow_virtual_network_access = true
  use_remote_gateways          = false
}

resource "azurerm_virtual_network_peering" "replica_primary" {
  name                      = "aadds-replica-primary"
  resource_group_name       = azurerm_virtual_network.replica.resource_group_name
  virtual_network_name      = azurerm_virtual_network.replica.name
  remote_virtual_network_id = azurerm_virtual_network.primary.id

  allow_forwarded_traffic      = true
  allow_gateway_transit        = false
  allow_virtual_network_access = true
  use_remote_gateways          = false
}

resource "azurerm_virtual_network_dns_servers" "replica" {
  virtual_network_id = azurerm_virtual_network.replica.id
  dns_servers        = azurerm_active_directory_domain_service.example.initial_replica_set[0].domain_controller_ip_addresses
}

resource "azurerm_active_directory_domain_service_replica_set" "replica" {
  domain_service_id = azurerm_active_directory_domain_service.example.id
  location          = azurerm_resource_group.replica.location
  subnet_id         = azurerm_subnet.aadds_replica.id

  depends_on = [
    azurerm_subnet_network_security_group_association.replica,
    azurerm_virtual_network_peering.primary_replica,
    azurerm_virtual_network_peering.replica_primary,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `domain_service_id` - (Required) The ID of the Domain Service for which to create this Replica Set. Changing this forces a new resource to be created.
  
* `location` - (Required) The Azure location where this Replica Set should exist. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the subnet in which to place this Replica Set. Changing this forces a new resource to be created.
  
## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Domain Service Replica Set.
  
* `domain_controller_ip_addresses` - A list of subnet IP addresses for the domain controllers in this Replica Set, typically two.

* `external_access_ip_address` - The publicly routable IP address for the domain controllers in this Replica Set.

* `service_status` - The current service status for the replica set.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the Domain Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Domain Service.
* `delete` - (Defaults to 1 hour) Used when deleting the Domain Service.

## Import

Domain Service Replica Sets can be imported using the resource ID of the parent Domain Service and the Replica Set ID, e.g.

```shell
terraform import azurerm_active_directory_domain_service_replica_set.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.AAD/domainServices/instance1/replicaSets/00000000-0000-0000-0000-000000000000
```
