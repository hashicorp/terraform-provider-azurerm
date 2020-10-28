provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-k8s-rg"
  location = var.location
}

resource "azurerm_virtual_network" "example" {
  name                = "${var.prefix}-vnet"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.10.0.0/16"]
}

resource "azurerm_subnet" "example-nodepool" {
  name                 = "default"
  virtual_network_name = azurerm_virtual_network.example.name
  resource_group_name  = azurerm_resource_group.example.name
  address_prefixes     = ["10.10.1.0/24"]
}

resource "azurerm_subnet" "example-aci" {
  name                 = "aci"
  virtual_network_name = azurerm_virtual_network.example.name
  resource_group_name  = azurerm_resource_group.example.name
  address_prefixes     = ["10.10.3.0/24"]

  delegation {
    name = "aciDelegation"
    service_delegation {
      name    = "Microsoft.ContainerInstance/containerGroups"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_kubernetes_cluster" "example" {
  name                = "${var.prefix}-k8s"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  dns_prefix          = "${var.prefix}-k8s"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
    vnet_subnet_id = azurerm_subnet.example-nodepool.id
  }

  network_profile {
    network_plugin = "azure"
    network_policy = "azure"
    load_balancer_sku = "standard"
  }

  service_principal {
    client_id     = azuread_service_principal.example.application_id
    client_secret = var.password
  }

  addon_profile {
    aci_connector_linux {
      enabled = true
      subnet_name = azurerm_subnet.example-aci.name
    }

    azure_policy {
      enabled = false
    }

    http_application_routing {
      enabled = false
    }

    kube_dashboard {
      enabled = true
    }

    oms_agent {
      enabled = false
    }
  }
}

resource "azurerm_role_assignment" "example" {
  scope                = azurerm_subnet.example-aci.id
  role_definition_name = "Network Contributor"
  principal_id         = azuread_service_principal.example.object_id
}

resource "azuread_application" "example" {
  name                       = "example"
  homepage                   = "http://homepage"
  identifier_uris            = ["http://uri"]
  reply_urls                 = ["http://replyurl"]
  available_to_other_tenants = false
  oauth2_allow_implicit_flow = true
}

resource "azuread_service_principal" "example" {
  application_id               = azuread_application.example.application_id
  app_role_assignment_required = false
}

resource "azuread_service_principal_password" "example" {
  service_principal_id = azuread_service_principal.example.id
  description          = "My managed password"
  value                = var.password
  end_date             = "2099-01-01T01:02:03Z"
}
