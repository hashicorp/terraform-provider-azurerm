resource "azurerm_resource_group" "rg" {
  name     = "${var.resource_group_name}"
  location = "${var.resource_group_location}"
}

resource "azurerm_scheduler_job_collection" "jobs" {
  name                = "example_job_collection"
  location            = "${azurerm_resource_group.rg.location}"
  resource_group_name = "${azurerm_resource_group.rg.name}"
  sku                 = "free"
  state               = "enabled"

  quota {
    max_job_count            = 5
    max_recurrence_interval  = 24
    max_recurrence_frequency = "hour"
  }
}
