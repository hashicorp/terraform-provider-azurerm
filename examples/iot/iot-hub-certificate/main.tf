# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

module "naming" {
  source  = "Azure/naming/azurerm"
  version = "0.4.0"
  prefix  = [var.prefix]
}

resource "azurerm_resource_group" "example" {
  name     = module.naming.resource_group.name
  location = var.location
}

resource "azurerm_iothub_dps_certificate" "example" {
  name                = module.naming.iothub_dps_certificate.name
  resource_group_name = azurerm_resource_group.example.name
  iot_dps_name        = var.iot_hub_dps_name
  certificate_content = filebase64(var.dps_cert_content)
  is_verified         = false
}
