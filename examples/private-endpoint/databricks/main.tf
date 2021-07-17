provider "azurerm" {
  features {}
}

data "azurerm_databricks_workspace_private_endpoint_connection" "test" {
  workspace_id        = azurerm_databricks_workspace.test.id
  private_endpoint_id = azurerm_private_endpoint.databricks.id
}

resource "azurerm_resource_group" "test" {
  name     = "${var.prefix}-databricks-private-endpoint"
  location = "eastus2"
}

resource "azurerm_virtual_network" "test" {
  name                = "${var.prefix}-vnet-databricks"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "public" {
  name                 = "${var.prefix}-sn-public"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "databricks-del-pub-${var.prefix}"

    service_delegation {
      actions = [
          "Microsoft.Network/virtualNetworks/subnets/join/action",
          "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
          "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
        ]
      name = "Microsoft.Databricks/workspaces"
    }
  }
}

resource "azurerm_subnet" "private" {
  name                 = "${var.prefix}-sn-private"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "databricks-del-pri-${var.prefix}"

    service_delegation {
      actions = [
          "Microsoft.Network/virtualNetworks/subnets/join/action",
          "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
          "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
        ]
      name = "Microsoft.Databricks/workspaces"
    }
  }
}

resource "azurerm_subnet" "private_endpoint" {
  name                 = "${var.prefix}-sn-private-endpoint"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.3.0/24"]

  enforce_private_link_endpoint_network_policies = true
}

resource "azurerm_subnet_network_security_group_association" "private" {
  subnet_id                 = azurerm_subnet.private.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_subnet_network_security_group_association" "public" {
  subnet_id                 = azurerm_subnet.public.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_network_security_group" "test" {
  name                = "${var.prefix}-nsg-databricks"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_databricks_workspace" "test" {
  name                        = "acctestDBW-${var.prefix}"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku                         = "premium"
  managed_resource_group_name = "acctestRG-DBW-${var.prefix}-managed"

  public_network_access_enabled        = false
  require_network_security_group_rules = "NoAzureDatabricksRules"

  custom_parameters {
    no_public_ip        = true
    public_subnet_name  = azurerm_subnet.public.name
    private_subnet_name = azurerm_subnet.private.name
    virtual_network_id  = azurerm_virtual_network.test.id

    public_subnet_network_security_group_association_id  = azurerm_subnet_network_security_group_association.public.id
    private_subnet_network_security_group_association_id = azurerm_subnet_network_security_group_association.private.id
  }

  tags = {
    Environment = "Production"
    Pricing     = "Premium"
  }
}

resource "azurerm_private_endpoint" "databricks" {
  name                = "${var.prefix}-pe-databricks"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.private_endpoint.id

  private_service_connection {
    name                           = "${var.prefix}-psc"
    is_manual_connection           = false
    private_connection_resource_id = azurerm_databricks_workspace.test.id
    subresource_names              = ["databricks_ui_api"]
  }
}

resource "azurerm_private_dns_zone" "test" {
  depends_on = [azurerm_private_endpoint.databricks]

  name                = "privatelink.azuredatabricks.net"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_cname_record" "test" {
  name                = azurerm_databricks_workspace.test.workspace_url
  zone_name           = azurerm_private_dns_zone.test.name
  resource_group_name = azurerm_resource_group.test.name
  ttl                 = 300
  record              = "eastus2-c2.azuredatabricks.net"
}

output "databricks_workspace_private_endpoint_connection_workspace_id" {
  value = data.azurerm_databricks_workspace_private_endpoint_connection.test.workspace_id
}

output "databricks_workspace_private_endpoint_connection_private_endpoint_id" {
  value = data.azurerm_databricks_workspace_private_endpoint_connection.test.private_endpoint_id
}

output "databricks_workspace_private_endpoint_connection_name" {
  value = data.azurerm_databricks_workspace_private_endpoint_connection.test.connections.0.name
}

output "databricks_workspace_private_endpoint_connection_workspace_private_endpoint_id" {
  value = data.azurerm_databricks_workspace_private_endpoint_connection.test.connections.0.workspace_private_endpoint_id
}

output "databricks_workspace_private_endpoint_connection_status" {
  value = data.azurerm_databricks_workspace_private_endpoint_connection.test.connections.0.status
}

output "databricks_workspace_private_endpoint_connection_description" {
  value = data.azurerm_databricks_workspace_private_endpoint_connection.test.connections.0.description
}

output "databricks_workspace_private_endpoint_connection_action_required" {
  value = data.azurerm_databricks_workspace_private_endpoint_connection.test.connections.0.action_required
}