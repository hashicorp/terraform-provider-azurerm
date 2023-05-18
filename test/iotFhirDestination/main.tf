provider "azurerm" {
  features {}
}

data "azurerm_resource_group" "test" {
  name = "iottestRG-healthcareapi"
}

data "azurerm_healthcare_workspace" "test" {
  name                = "iotwks"
  resource_group_name = data.azurerm_resource_group.test.name
}

resource "azurerm_healthcare_iot_connector" "test" {
  name = "tftest"
  workspace_id = data.azurerm_healthcare_workspace.test.id
  location = "east us"
  identity {
    type = "SystemAssigned"
  }
  eventhub_namespace_name = "iottestenh"
  eventhub_name = "iottesteh"
  eventhub_consumer_group_name = "iottestcg"
  device_mapping               = <<JSON
{
"templateType": "CollectionContent",
"template": []
}
JSON
}

resource "azurerm_healthcare_fhir_service" "test" {
  name                = "fhirtftest"
  location            = data.azurerm_resource_group.test.location
  resource_group_name = data.azurerm_resource_group.test.name
  workspace_id        = data.azurerm_healthcare_workspace.test.id
  kind                = "fhir-R4"
  authentication {
    authority = "https://login.microsoftonline.com/72f988bf-86f1-41af-91ab-2d7cd011db47"
    audience  = "https://acctestfhir.fhir.azurehealthcareapis.com"
  }
}

resource "azurerm_healthcare_iot_fhir_destination" "test" {
  name = "xiaxintestdes"
  location = "east us"
  iot_connector_id = azurerm_healthcare_iot_connector.test.id
  destination_fhir_service_id = azurerm_healthcare_fhir_service.test.id
  destination_identity_resolution_type = "Create"
  destination_fhir_mapping = <<JSON
{
  "templateType": "CollectionFhirTemplate",
            "template": []
}
JSON
}