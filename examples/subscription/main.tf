# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

resource "azurerm_subscription" "subscription" {
  alias             = "${var.prefix}-sub"
  subscription_name = "${var.prefix}-sub"
  billing_scope_id  = var.subscription_billing_scope_id
  workload          = "Production"
  tags = {
    "Environment" = "Production"
  }
}
