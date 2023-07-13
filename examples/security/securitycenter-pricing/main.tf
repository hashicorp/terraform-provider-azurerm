# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

resource "azurerm_security_center_subscription_pricing" "example" {
  tier          = "Standard"
  resource_type = "VirtualMachines"
}
