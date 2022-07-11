provider "azurerm" {
  features {}
}

variable "prefix" {
  description = "The prefix which should be used for all resources in this example"
}

variable "location" {
  description = "The Azure Region in which all primary resources in this example should be created."
}

variable "location_secondary" {
  description = "The Azure Region in which all secondary resources in this example should be created."
}

resource "azurerm_resource_group" "primary" {
  name     = "${var.prefix}-primary-resources"
  location = var.location
}

module "primary" {
  source = "./modules/sql_managed_instance"

  prefix              = "${var.prefix}primary"
  location            = azurerm_resource_group.primary.location
  resource_group_name = azurerm_resource_group.primary.name
}

resource "azurerm_resource_group" "secondary" {
  name     = "${var.prefix}-secondary-resources"
  location = var.location_secondary
}

module "secondary" {
  source = "./modules/sql_managed_instance"

  prefix              = "${var.prefix}secondary"
  location            = azurerm_resource_group.secondary.location
  resource_group_name = azurerm_resource_group.secondary.name
  dns_zone_partner_id = module.primary.managed_instance_id

  vnet_range   = "10.1.0.0/16"
  subnet_range = "10.1.0.0/24"
}

module "vnet_gateway" {
  source = "./modules/vnet_gateway"

  shared_key = "s3cre7_key!"

  prefix                 = var.prefix
  location_1             = azurerm_resource_group.primary.location
  resource_group_name_1  = azurerm_resource_group.primary.name
  vnet_name_1            = module.primary.vnet_name
  gateway_subnet_range_1 = "10.0.1.0/24"

  location_2             = azurerm_resource_group.secondary.location
  resource_group_name_2  = azurerm_resource_group.secondary.name
  vnet_name_2            = module.secondary.vnet_name
  gateway_subnet_range_2 = "10.0.2.0/24"
}

resource "azurerm_sql_managed_instance_failover_group" "example" {
  name                        = "${var.prefix}-fog"
  resource_group_name         = azurerm_resource_group.primary.name
  location                    = azurerm_resource_group.primary.location
  managed_instance_name       = module.primary.managed_instance_name
  partner_managed_instance_id = module.secondary.managed_instance_id

  read_write_endpoint_failover_policy {
    mode = "Manual"
  }

  depends_on = [
    module.vnet_gateway
  ]
}


