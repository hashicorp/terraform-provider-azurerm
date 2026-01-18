# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 2.0"
    }

    azuread = {
      source  = "hashicorp/azuread"
      version = "~> 2.0"
    }
  }
}
