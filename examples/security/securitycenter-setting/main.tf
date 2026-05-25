# Copyright IBM Corp. 2023, 2025
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

resource "azurerm_security_center_setting" "example" {
  setting_name = "MCAS"
  enabled      = true
}
