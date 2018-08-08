resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_storage_account" "example" {
  name                     = "${var.prefix}stor"
  resource_group_name      = "${azurerm_resource_group.example.name}"
  location                 = "${azurerm_resource_group.example.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_queue" "example" {
  name                 = "${var.prefix}-queue"
  resource_group_name  = "${azurerm_resource_group.example.name}"
  storage_account_name = "${azurerm_storage_account.example.name}"
}

resource "azurerm_scheduler_job" "example" {
  name                = "${var.prefix}-job"
  resource_group_name = "${azurerm_resource_group.example.name}"
  job_collection_name = "${azurerm_scheduler_job_collection.example.name}"

  action_storage_queue = {
    storage_account_name = "${azurerm_storage_account.example.name}"
    storage_queue_name   = "${azurerm_storage_queue.example.name}"
    sas_token            = "${azurerm_storage_account.example.primary_access_key}"
    message              = "storage message"
  }
}
