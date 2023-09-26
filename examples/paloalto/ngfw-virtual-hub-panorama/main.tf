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

resource "azurerm_virtual_wan" "example" {
  name                = "${var.prefix}-virtual-wan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_virtual_hub" "example" {
  name                = "${var.prefix}-virtual-hub"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  virtual_wan_id      = azurerm_virtual_wan.example.id
  address_prefix      = "10.0.1.0/24"

  tags = {
    hubSaaSPreview = "true"
  }
}

resource "azurerm_palo_alto_virtual_network_appliance" "example" {
  name           = "${var.prefix}-nva"
  virtual_hub_id = azurerm_virtual_hub.example.id
}

resource "azurerm_palo_alto_next_generation_firewall_virtual_hub_panorama" "example" {
  name                   = "${var.prefix}-ngfw-vhub"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  panorama_base64_config = var.panorama-config

  network_profile {
    virtual_hub_id               = azurerm_virtual_hub.example.id
    network_virtual_appliance_id = azurerm_palo_alto_virtual_network_appliance.example.id
    public_ip_address_ids        = [azurerm_public_ip.example.id]
    egress_nat_ip_address_ids    = [azurerm_public_ip.egress.id]
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
