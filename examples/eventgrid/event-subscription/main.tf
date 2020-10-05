provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location
}

resource "azurerm_storage_account" "example" {
  name                     = "${var.prefix}sa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_queue" "example" {
  name                 = "${var.prefix}-sq"
  storage_account_name = azurerm_storage_account.example.name
}

resource "azurerm_storage_container" "example" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

resource "azurerm_storage_blob" "example" {
  name = "herpderp1.vhd"

  storage_account_name   = azurerm_storage_account.example.name
  storage_container_name = azurerm_storage_container.example.name

  type = "Page"
  size = 5120
}

resource "azurerm_eventgrid_event_subscription" "example" {
  name  = "${var.prefix}-eventsubs"
  scope = azurerm_resource_group.example.id

  storage_queue_endpoint {
    storage_account_id = azurerm_storage_account.example.id
    queue_name         = azurerm_storage_queue.example.name
  }

  storage_blob_dead_letter_destination {
    storage_account_id          = azurerm_storage_account.example.id
    storage_blob_container_name = azurerm_storage_container.example.name
  }

  retry_policy {
    event_time_to_live    = 11
    max_delivery_attempts = 11
  }

  labels = ["test", "test1", "test2"]
}
