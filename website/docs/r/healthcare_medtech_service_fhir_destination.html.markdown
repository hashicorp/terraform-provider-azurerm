---
subcategory: "Healthcare"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_healthcare_medtech_service_fhir_destination"
description: |-
  Manages a Healthcare Med Tech (Internet of Medical Things) Service Fhir Destination.
---

# azurerm_healthcare_medtech_service_fhir_Destination

Manages a Healthcare Med Tech Service Fhir Destination.

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}

data "azurerm_client_config" "current" {
}

resource "azurerm_healthcare_workspace" "example" {
  name                = "exampleworkspace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_eventhub_namespace" "example" {
  name                = "example-ehn"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"
}

resource "azurerm_eventhub" "example" {
  name                = "example-eh"
  namespace_name      = azurerm_eventhub_namespace.example.name
  resource_group_name = azurerm_resource_group.example.name
  partition_count     = 1
  message_retention   = 1
}

resource "azurerm_eventhub_consumer_group" "example" {
  name                = "$default"
  namespace_name      = azurerm_eventhub_namespace.example.name
  eventhub_name       = azurerm_eventhub.example.name
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_healthcare_fhir_service" "example" {
  name                = "examplefhir"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  workspace_id        = azurerm_healthcare_workspace.example.id
  kind                = "fhir-R4"

  authentication {
    authority = "https://login.microsoftonline.com/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    audience  = "https://examplefhir.fhir.azurehealthcareapis.com"
  }
}

resource "azurerm_healthcare_medtech_service" "example" {
  name         = "examplemt"
  workspace_id = azurerm_healthcare_workspace.example.id
  location     = azurerm_resource_group.example.location

  eventhub_namespace_name      = azurerm_eventhub_namespace.example.name
  eventhub_name                = azurerm_eventhub.example.name
  eventhub_consumer_group_name = azurerm_eventhub_consumer_group.example.name

  device_mapping_json = jsonencode({
    "templateType" : "CollectionContent",
    "template" : []
  })
}

resource "azurerm_healthcare_medtech_service_fhir_destination" "example" {
  name                                 = "examplemtdes"
  location                             = "east us"
  medtech_service_id                   = azurerm_healthcare_medtech_service.example.id
  destination_fhir_service_id          = azurerm_healthcare_fhir_service.example.id
  destination_identity_resolution_type = "Create"

  destination_fhir_mapping_json = jsonencode({
    "templateType" : "CollectionFhirTemplate",
    "template" : [
      {
        "templateType" : "CodeValueFhir",
        "template" : {
          "codes" : [
            {
              "code" : "8867-4",
              "system" : "http://loinc.org",
              "display" : "Heart rate"
            }
          ],
          "periodInterval" : 60,
          "typeName" : "heartrate",
          "value" : {
            "defaultPeriod" : 5000,
            "unit" : "count/min",
            "valueName" : "hr",
            "valueType" : "SampledData"
          }
        }
      }
    ]
  })
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Healthcare Med Tech Service Fhir Destination. Changing this forces a new Healthcare Med Tech Service Fhir Destination to be created.

* `medtech_service_id` - (Required) Specifies the name of the Healthcare Med Tech Service where the Healthcare Med Tech Service Fhir Destination should exist. Changing this forces a new Healthcare Med Tech Service Fhir Destination to be created.

* `location` - (Required) Specifies the Azure Region where the Healthcare Med Tech Service Fhir Destination should be created. Changing this forces a new Healthcare Med Tech Service Fhir Destination to be created.

* `destination_fhir_service_id` - (Required) Specifies the destination fhir service id of the Med Tech Service Fhir Destination.

* `destination_identity_resolution_type` - (Required) Specifies the destination identity resolution type where the Healthcare Med Tech Service Fhir Destination should be created. Possible values are `Create`, `Lookup`.

* `destination_fhir_mapping_json` - (Required) Specifies the destination Fhir mappings of the Med Tech Service Fhir Destination.

## Attributes Reference

The following arguments are supported:

* `id` - The ID of the Healthcare Med Tech Service Fhir Destination.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the Healthcare Med Tech Service Fhir Destination.
* `read` - (Defaults to 5 minutes) Used when retrieving the Healthcare Med Tech Service Fhir Destination.
* `update` - (Defaults to 90 minutes) Used when updating the Healthcare Med Tech Service Fhir Destination.
* `delete` - (Defaults to 90 minutes) Used when deleting the Healthcare Med Tech Service Fhir Destination.

## Import

Healthcare Med Tech Service Fhir Destination can be imported using the resource`id`, e.g.

```shell
terraform import azurerm_healthcare_medtech_service_fhir_destination.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/iotConnectors/iotconnector1/fhirDestinations/destination1
```
