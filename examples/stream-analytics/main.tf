provider "azurerm" {}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-example-resources"
  location = "${var.location}"
}

resource "azurerm_storage_account" "example" {
  name                     = "${var.prefix}-examplestoracc"
  resource_group_name      = "${azurerm_resource_group.example.name}"
  location                 = "${azurerm_resource_group.example.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                  = "${var.prefix}example"
  resource_group_name   = "${azurerm_resource_group.example.name}"
  storage_account_name  = "${azurerm_storage_account.example.name}"
  container_access_type = "private"
}

resource "azurerm_stream_analytics_job" "example" {
  name                                     = "${var.prefix}-example-job"
  resource_group_name                      = "${azurerm_resource_group.example.name}"
  location                                 = "${azurerm_resource_group.example.location}"
  compatibility_level                      = "1.1"
  data_locale                              = "en-US"
  events_late_arrival_max_delay_in_seconds = 60
  events_out_of_order_max_delay_in_seconds = 50
  events_out_of_order_policy               = "Adjust"
  output_error_policy                      = "Drop"
  streaming_units                          = 3

  tags = {
    environment = "Example"
  }

  transformation_query = <<QUERY
    SELECT *
    INTO [YourOutputAlias]
    FROM [YourInputAlias]
QUERY
}

resource "azurerm_stream_analytics_reference_input_blob" "test" {
  name                         = "${var.prefix}-blob-reference-input"
  stream_analytics_job_name    = "${azurerm_stream_analytics_job.example.name}"
  resource_group_name          = "${azurerm_stream_analytics_job.example.resource_group_name}"
  storage_account_name         = "${azurerm_storage_account.example.name}"
  storage_account_key          = "${azurerm_storage_account.example.primary_access_key}"
  storage_container_name       = "${azurerm_storage_container.example.name}"
  path_pattern                 = "some-random-pattern"
  date_format                  = "yyyy/MM/dd"
  time_format                  = "HH"

  serialization {
    type     = "Json"
    encoding = "UTF8"
  }
}
