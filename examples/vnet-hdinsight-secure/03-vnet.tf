# Create a virtual network
resource "azurerm_virtual_network" "vnet" {
  name                = "${var.resource_prefix}-vnet"
  location            = "${azurerm_resource_group.rg.location}"
  address_space       = ["${var.address_space}"]
  resource_group_name = "${azurerm_resource_group.rg.name}"
  dns_servers         = "${var.dns_servers}"
  tags                = "${var.tags}"
}

# Create a subnet
resource "azurerm_subnet" "subnet" {
  name                 = "${var.subnet_name}"
  virtual_network_name = "${azurerm_virtual_network.vnet.name}"
  resource_group_name  = "${azurerm_resource_group.rg.name}"
  address_prefix       = "${var.subnet_prefix}"
  service_endpoints    = ["${var.service_endpoints}"]
}

# Create a network security group
resource "azurerm_network_security_group" "nsg" {
  name                = "${var.resource_prefix}-nsg"
  location            = "${azurerm_resource_group.rg.location}"
  resource_group_name = "${azurerm_resource_group.rg.name}"
  tags                = "${var.tags}"
}

# Create network security group rules to secure HDInsight management traffic
resource "azurerm_network_security_rule" "nsg_rule_allow_hdi_mgmt_traffic" {
  name                        = "allow_hdi_mgmt_traffic"
  priority                    = 300
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "*"
  source_port_range           = "*"
  source_address_prefixes     = ["${var.source_address_prefixes_mgmt}"]
  destination_port_range      = "443"
  destination_address_prefix  = "VirtualNetwork"
  resource_group_name         = "${azurerm_resource_group.rg.name}"
  network_security_group_name = "${azurerm_network_security_group.nsg.name}"
}

resource "azurerm_network_security_rule" "nsg_rule_allow_azure_resolver_traffic" {
  name                        = "allow_azure_resolver_traffic"
  priority                    = 301
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "*"
  source_port_range           = "*"
  source_address_prefix       = "${var.source_address_prefix_resolver}"
  destination_port_range      = "443"
  destination_address_prefix  = "VirtualNetwork"
  resource_group_name         = "${azurerm_resource_group.rg.name}"
  network_security_group_name = "${azurerm_network_security_group.nsg.name}"
}

resource "azurerm_network_security_rule" "nsg_rule_allow_hdi_mgmt_traffic_regional" {
  name                        = "allow_hdi_mgmt_traffic_regional"
  priority                    = 302
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "*"
  source_port_range           = "*"
  source_address_prefixes     = ["${var.source_address_prefixes_mgmt_region}"]
  destination_port_range      = "443"
  destination_address_prefix  = "VirtualNetwork"
  resource_group_name         = "${azurerm_resource_group.rg.name}"
  network_security_group_name = "${azurerm_network_security_group.nsg.name}"
}

resource "azurerm_subnet_network_security_group_association" "subnet_to_nsg" {
  subnet_id                 = "${azurerm_subnet.subnet.id}"
  network_security_group_id = "${azurerm_network_security_group.nsg.id}"
}
