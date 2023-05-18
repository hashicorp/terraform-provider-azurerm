provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "iottestRG-healthcareapi"
  location = "east us"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "iottestenh"
  location            = "east us"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_eventhub" "test" {
  name                = "iotTestEh"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 1
  message_retention   = 1
}

resource "azurerm_eventhub_consumer_group" "test" {
  name                = "iotTestcg"
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_healthcare_workspace" "test" {
  name                = "iotwks"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_healthcare_medtech_service" "test" {
  name                         = "tftest"
  workspace_id                 = azurerm_healthcare_workspace.test.id
  location                     = "east us"
  identity {
    type = "SystemAssigned"
  }
  eventhub_namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name                = azurerm_eventhub.test.name
  eventhub_consumer_group_name = azurerm_eventhub_consumer_group.test.name
  device_mapping               = <<JSON
{
            "templateType": "CollectionContent",
            "template": []
}
JSON
}
resource "azurerm_healthcare_fhir_service" "test" {
  name                = "fhirtftest"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_healthcare_workspace.test.id
  kind                = "fhir-R4"
  authentication {
    authority = "https://login.microsoftonline.com/72f988bf-86f1-41af-91ab-2d7cd011db47"
    audience  = "https://acctestfhir.fhir.azurehealthcareapis.com"
  }
}

resource "azurerm_healthcare_medtech_service_fhir_destination" "test" {
  name = "xiaxintestdes"
  location = "east us"
  medtech_service_id = azurerm_healthcare_medtech_service.test.id
  destination_fhir_service_id = azurerm_healthcare_fhir_service.test.id
  destination_identity_resolution_type = "Create"
  destination_fhir_mapping = <<JSON
{
  "templateType": "CollectionFhirTemplate",
            "template": []
}
JSON
}