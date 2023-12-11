# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-example-resources"
  location = var.location
}

### Stream Analytics Job ###
resource "azurerm_stream_analytics_job" "example" {
  name                                     = "${var.prefix}-example-job"
  resource_group_name                      = azurerm_resource_group.example.name
  location                                 = azurerm_resource_group.example.location
  compatibility_level                      = "1.1"
  data_locale                              = "en-US"
  events_late_arrival_max_delay_in_seconds = 60
  events_out_of_order_max_delay_in_seconds = 50
  events_out_of_order_policy               = "Adjust"
  output_error_policy                      = "Drop"
  streaming_units                          = 3

  identity {
    type = "SystemAssigned"
  }

  tags = {
    environment = "Example"
  }

  transformation_query = <<QUERY
    SELECT *
    INTO [YourOutputAlias]
    FROM [YourInputAlias]
QUERY
}

### REFERENCE INPUT: Blob storage/ADLS Gen2 ###
resource "azurerm_storage_account" "example" {
  name                     = "${var.prefix}examplestoracc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_role_assignment" "example" {
  scope                = azurerm_storage_account.example.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = azurerm_stream_analytics_job.example.identity[0].principal_id
}

resource "azurerm_storage_container" "example" {
  name                  = "${var.prefix}example"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

resource "azurerm_stream_analytics_reference_input_blob" "test" {
  name                      = "${var.prefix}-blob-reference-input"
  stream_analytics_job_name = azurerm_stream_analytics_job.example.name
  resource_group_name       = azurerm_stream_analytics_job.example.resource_group_name
  storage_account_name      = azurerm_storage_account.example.name
  storage_container_name    = azurerm_storage_container.example.name
  authentication_mode       = "Msi"
  path_pattern              = "some-random-pattern"
  date_format               = "yyyy/MM/dd"
  time_format               = "HH"

  serialization {
    type     = "Json"
    encoding = "UTF8"
  }
}

### STREAM INPUT: Event Hub ###
resource "azurerm_eventhub_namespace" "example" {
  name                = "${var.prefix}-example-namespace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"
  capacity            = 1
}

resource "azurerm_eventhub" "example" {
  name                = "${var.prefix}-example-eventhub"
  namespace_name      = azurerm_eventhub_namespace.example.name
  resource_group_name = azurerm_resource_group.example.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_role_assignment" "example_eventhub" {
  scope                = azurerm_eventhub.example.id
  role_definition_name = "Azure Event Hubs Data receiver"
  principal_id         = azurerm_stream_analytics_job.example.identity[0].principal_id
}

resource "azurerm_stream_analytics_stream_input_eventhub_v2" "example" {
  name                    = "${var.prefix}-eventhub-stream-input"
  stream_analytics_job_id = azurerm_stream_analytics_job.example.id
  eventhub_name           = azurerm_eventhub.example.name
  servicebus_namespace    = azurerm_eventhub_namespace.example.name
  authentication_mode     = "Msi"

  serialization {
    type     = "Json"
    encoding = "UTF8"
  }
}