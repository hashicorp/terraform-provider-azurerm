# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "example" {
  name                       = "${var.prefix}-mg"
  display_name               = "${var.prefix}-mg"
  parent_management_group_id = var.parent_management_group_id
  subscription_ids           = var.subscription_ids
}
