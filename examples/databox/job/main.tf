provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location
}

resource "azurerm_storage_account" "example" {
  name                     = "${var.prefix}storageaccount"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location

  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "RAGRS"
}

resource "azurerm_data_box_job" "example" {
  name                = "${var.prefix}-databoxjob"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  contact_details {
    name         = "DataBoxJobTester"
    emails       = ["some.user@example.com"]
    phone_number = "+112345678912"
  }

  destination_storage_account {
    storage_account_id = azurerm_storage_account.example.id
  }

  preferred_shipment_type = "MicrosoftManaged"

  shipping_address {
    city              = "San Francisco"
    country           = "US"
    postal_code       = "94107"
    state_or_province = "CA"
    street_address_1  = "16 TOWNSEND ST"
  }

  sku_name = "DataBox"
}
