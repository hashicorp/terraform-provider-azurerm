# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-databricks-esc"
  location = "West Europe"
}

resource "azurerm_databricks_workspace" "example" {
  name                        = "${var.prefix}-DBW"
  resource_group_name         = azurerm_resource_group.example.name
  location                    = azurerm_resource_group.example.location
  sku                         = "premium"
  managed_resource_group_name = "${var.prefix}-DBW-managed-esc"

  enhanced_security_compliance {
    automatic_cluster_update_enabled      = true
    compliance_security_profile_enabled   = true
    compliance_security_profile_standards = ["HIPAA", "PCI_DSS", "GERMANY_C5", "HITRUST"]
    enhanced_security_monitoring_enabled  = true
  }
}
