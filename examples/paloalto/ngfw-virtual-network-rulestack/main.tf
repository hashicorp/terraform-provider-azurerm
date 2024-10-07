# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location
}

resource "azurerm_public_ip" "example" {
  name                = "${var.prefix}-public-ip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_public_ip" "egress" {
  name                = "${var.prefix}-public-ip-egress"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_network_security_group" "example" {
  name                = "${var.prefix}SecurityGroup1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_virtual_network" "example" {
  name                = "${var.prefix}-virtual-network"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  tags = {
    environment = "ManualTest"
  }
}

resource "azurerm_subnet" "trust" {
  name                 = "${var.prefix}-trust-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "trusted"

    service_delegation {
      name = "PaloAltoNetworks.Cloudngfw/firewalls"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_subnet_network_security_group_association" "trust" {
  subnet_id                 = azurerm_subnet.trust.id
  network_security_group_id = azurerm_network_security_group.example.id
}

resource "azurerm_subnet" "untrust" {
  name                 = "${var.prefix}-untrust-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "untrusted"

    service_delegation {
      name = "PaloAltoNetworks.Cloudngfw/firewalls"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_subnet_network_security_group_association" "untrust" {
  subnet_id                 = azurerm_subnet.untrust.id
  network_security_group_id = azurerm_network_security_group.example.id
}

resource "azurerm_palo_alto_local_rulestack" "example" {
  name                = "${var.prefix}-rulestack"
  resource_group_name = azurerm_resource_group.example.name
  location            = "westeurope"
}

resource "azurerm_palo_alto_local_rulestack_rule" "example" {
  name         = "${var.prefix}-rulestack-rule"
  rulestack_id = azurerm_palo_alto_local_rulestack.example.id
  priority     = 9999
  action       = "DenySilent"

  applications = ["any"]

  destination {
    cidrs = ["any"]
  }

  source {
    cidrs = ["any"]
  }
}

resource "azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack" "example" {
  name                = "${var.prefix}-ngfw-vnet-lrs"
  resource_group_name = azurerm_resource_group.example.name
  rulestack_id        = azurerm_palo_alto_local_rulestack.example.id

  network_profile {
    public_ip_address_ids     = [azurerm_public_ip.example.id]
    egress_nat_ip_address_ids = [azurerm_public_ip.egress.id]

    vnet_configuration {
      virtual_network_id  = azurerm_virtual_network.example.id
      trusted_subnet_id   = azurerm_subnet.trust.id
      untrusted_subnet_id = azurerm_subnet.untrust.id
    }
  }


  dns_settings {
    use_azure_dns = true
  }

  destination_nat {
    name     = "${var.prefix}DNAT-1"
    protocol = "TCP"
    frontend_config {
      public_ip_address_id = azurerm_public_ip.example.id
      port                 = 8081
    }
    backend_config {
      public_ip_address = "10.0.1.101"
      port              = 18081
    }
  }

  destination_nat {
    name     = "${var.prefix}DNAT-2"
    protocol = "UDP"
    frontend_config {
      public_ip_address_id = azurerm_public_ip.example.id
      port                 = 8082
    }
    backend_config {
      public_ip_address = "10.0.1.102"
      port              = 18082
    }
  }
}
