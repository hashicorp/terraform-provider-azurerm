provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-HDInsight-Hadoop"
  location = var.location
}

resource "azurerm_user_assigned_identity" "example" {
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  name = "${var.prefix}-identity"
}

resource "azurerm_role_assignment" "example" {
  scope                = data.azurerm_subscription.primary.id
  role_definition_name = "Storage Blob Data Owner"
  principal_id         = azurerm_user_assigned_identity.example.principal_id

  skip_service_principal_aad_check = true
}

resource "azurerm_virtual_network" "example" {
  name                = "${var.prefix}-vnet"
  address_space       = ["172.16.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "${var.prefix}-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["172.16.11.0/26"]

  enforce_private_link_service_network_policies = true
}

resource "azurerm_public_ip" "example" {
  name                = "${var.prefix}-pip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1"]
}

resource "azurerm_nat_gateway" "example" {
  name                    = "${var.prefix}-natGateway"
  location                = azurerm_resource_group.example.location
  resource_group_name     = azurerm_resource_group.example.name
  sku_name                = "Standard"
  idle_timeout_in_minutes = 10
  zones                   = ["1"]
}

resource "azurerm_nat_gateway_public_ip_association" "example" {
  nat_gateway_id       = azurerm_nat_gateway.example.id
  public_ip_address_id = azurerm_public_ip.example.id
}

resource "azurerm_subnet_nat_gateway_association" "example" {
  subnet_id      = azurerm_subnet.example.id
  nat_gateway_id = azurerm_nat_gateway.example.id
}

resource "azurerm_subnet_network_security_group_association" "example" {
  subnet_id                 = azurerm_subnet.example.id
  network_security_group_id = azurerm_network_security_group.example.id
}

resource "azurerm_storage_account" "example" {
  name                     = "${var.prefix}sa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "StorageV2"
  is_hns_enabled           = "true"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "example" {
  name               = "${var.prefix}-fs"
  storage_account_id = azurerm_storage_account.example.id
}

resource "azurerm_hdinsight_hadoop_cluster" "example" {
  depends_on = [azurerm_role_assignment.example, azurerm_nat_gateway.example, azurerm_subnet_network_security_group_association.example]

  name                = "${var.prefix}-hadoop"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  cluster_version     = "4.0"
  tier                = "Standard"

  component_version {
    hadoop = "3.1"
  }

  network {
    connection_direction = "Outbound"
    private_link_enabled = true
  }

  gateway {
    username = "accexampleusrgw"
    password = "TerrAform123!"
  }

  storage_account_gen2 {
    storage_resource_id          = azurerm_storage_account.example.id
    filesystem_id                = azurerm_storage_data_lake_gen2_filesystem.example.id
    managed_identity_resource_id = azurerm_user_assigned_identity.example.id
    is_default                   = true
  }

  roles {
    head_node {
      vm_size  = "Standard_D3_V2"
      username = "accexampleusrvm"
      password = "AccTestvdSC4daf986!"

     subnet_id          = azurerm_subnet.example.id
     virtual_network_id = azurerm_virtual_network.example.id
    }

    worker_node {
      vm_size               = "Standard_D4_V2"
      username              = "accexampleusrvm"
      password              = "AccTestvdSC4daf986!"
      target_instance_count = 3

     subnet_id          = azurerm_subnet.example.id
     virtual_network_id = azurerm_virtual_network.example.id
    }

    zookeeper_node {
      vm_size  = "Standard_D3_V2"
      username = "accexampleusrvm"
      password = "AccTestvdSC4daf986!"

     subnet_id          = azurerm_subnet.example.id
     virtual_network_id = azurerm_virtual_network.example.id
    }
  }
}

resource "azurerm_network_security_group" "example" {
  name                = "${var.prefix}-NSG"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  security_rule = [
    {
      access                                     = "Allow"
      description                                = "Rule can be deleted but do not change source ips."
      destination_address_prefix                 = "*"
      destination_address_prefixes               = []
      destination_application_security_group_ids = []
      destination_port_range                     = "443"
      destination_port_ranges                    = []
      direction                                  = "Inbound"
      name                                       = "Rule-101"
      priority                                   = 101
      protocol                                   = "Tcp"
      source_address_prefix                      = "VirtualNetwork"
      source_address_prefixes                    = []
      source_application_security_group_ids      = []
      source_port_range                          = "*"
      source_port_ranges                         = []
    },
    {
      access                                     = "Allow"
      description                                = "Rule can be deleted but do not change source ips."
      destination_address_prefix                 = "*"
      destination_address_prefixes               = []
      destination_application_security_group_ids = []
      destination_port_range                     = "*"
      destination_port_ranges                    = []
      direction                                  = "Inbound"
      name                                       = "Rule-103"
      priority                                   = 103
      protocol                                   = "*"
      source_address_prefix                      = "CorpNetPublic"
      source_address_prefixes                    = []
      source_application_security_group_ids      = []
      source_port_range                          = "*"
      source_port_ranges                         = []
    },
    {
      access                                     = "Allow"
      description                                = "Rule can be deleted but do not change source ips."
      destination_address_prefix                 = "*"
      destination_address_prefixes               = []
      destination_application_security_group_ids = []
      destination_port_range                     = "*"
      destination_port_ranges                    = []
      direction                                  = "Inbound"
      name                                       = "Rule-104"
      priority                                   = 104
      protocol                                   = "*"
      source_address_prefix                      = "CorpNetSaw"
      source_address_prefixes                    = []
      source_application_security_group_ids      = []
      source_port_range                          = "*"
      source_port_ranges                         = []
    },
    {
      access                                     = "Deny"
      description                                = "DO NOT DELETE"
      destination_address_prefix                 = "*"
      destination_address_prefixes               = []
      destination_application_security_group_ids = []
      destination_port_range                     = ""
      destination_port_ranges = [
        "111",
        "11211",
        "123",
        "13",
        "17",
        "19",
        "1900",
        "512",
        "514",
        "53",
        "5353",
        "593",
        "69",
        "873",
      ]
      direction                             = "Inbound"
      name                                  = "Rule-108"
      priority                              = 108
      protocol                              = "*"
      source_address_prefix                 = "Internet"
      source_address_prefixes               = []
      source_application_security_group_ids = []
      source_port_range                     = "*"
      source_port_ranges                    = []
    },
    {
      access                                     = "Deny"
      description                                = "DO NOT DELETE"
      destination_address_prefix                 = "*"
      destination_address_prefixes               = []
      destination_application_security_group_ids = []
      destination_port_range                     = ""
      destination_port_ranges = [
        "119",
        "137",
        "138",
        "139",
        "161",
        "162",
        "2049",
        "2301",
        "2381",
        "3268",
        "389",
        "5800",
        "5900",
        "636",
      ]
      direction                             = "Inbound"
      name                                  = "Rule-109"
      priority                              = 109
      protocol                              = "*"
      source_address_prefix                 = "Internet"
      source_address_prefixes               = []
      source_application_security_group_ids = []
      source_port_range                     = "*"
      source_port_ranges                    = []
    },
    {
      access                                     = "Deny"
      description                                = "DO NOT DELETE"
      destination_address_prefix                 = "*"
      destination_address_prefixes               = []
      destination_application_security_group_ids = []
      destination_port_range                     = ""
      destination_port_ranges = [
        "135",
        "23",
        "445",
        "5985",
        "5986",
      ]
      direction                             = "Inbound"
      name                                  = "Rule-107"
      priority                              = 107
      protocol                              = "Tcp"
      source_address_prefix                 = "Internet"
      source_address_prefixes               = []
      source_application_security_group_ids = []
      source_port_range                     = "*"
      source_port_ranges                    = []
    },
    {
      access                                     = "Deny"
      description                                = "DO NOT DELETE"
      destination_address_prefix                 = "*"
      destination_address_prefixes               = []
      destination_application_security_group_ids = []
      destination_port_range                     = ""
      destination_port_ranges = [
        "1433",
        "1434",
        "16379",
        "26379",
        "27017",
        "3306",
        "4333",
        "5432",
        "6379",
        "7000",
        "7001",
        "7199",
        "9042",
        "9160",
        "9300",
      ]
      direction                             = "Inbound"
      name                                  = "Rule-105"
      priority                              = 105
      protocol                              = "*"
      source_address_prefix                 = "Internet"
      source_address_prefixes               = []
      source_application_security_group_ids = []
      source_port_range                     = "*"
      source_port_ranges                    = []
    },
    {
      access                                     = "Deny"
      description                                = "DO NOT DELETE"
      destination_address_prefix                 = "*"
      destination_address_prefixes               = []
      destination_application_security_group_ids = []
      destination_port_range                     = ""
      destination_port_ranges = [
        "22",
        "3389",
      ]
      direction                             = "Inbound"
      name                                  = "Rule-106"
      priority                              = 106
      protocol                              = "Tcp"
      source_address_prefix                 = "Internet"
      source_address_prefixes               = []
      source_application_security_group_ids = []
      source_port_range                     = "*"
      source_port_ranges                    = []
    },
  ]
}
