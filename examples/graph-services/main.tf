# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-graph-services-rg"
  location = var.location
}

resource "azurerm_graph_services_account" "example" {
  name                = "${var.prefix}-graph-services-account"
  resource_group_name = azurerm_resource_group.example.name
  application_id      = data.azurerm_client_config.current.client_id

  tags = {
    Environment = "Example"
  }
}
