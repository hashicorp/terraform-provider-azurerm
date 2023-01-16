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
resource "azurerm_healthcare_medtech_service_fhir_destination" "test" {
  name                                 = "tfexmtdes"
  location                             = "east us"
  medtech_service_id                   = "mt_service_id"
  destination_fhir_service_id          = "fhir_service_id"
  destination_identity_resolution_type = "Create"

  destination_fhir_mapping_json = <<JSON
  {
            "templateType": "CollectionFhirTemplate",
            "template": [
              {
                "templateType": "CodeValueFhir",
                "template": {
                  "codes": [
                    {
                      "code": "8867-4",
                      "system": "http://loinc.org",
                      "display": "Heart rate"
                    }
                  ],
                  "periodInterval": 60,
                  "typeName": "heartrate",
                  "value": {
                    "defaultPeriod": 5000,
                    "unit": "count/min",
                    "valueName": "hr",
                    "valueType": "SampledData"
                  }
                }
              }
            ]
  }
  JSON
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Healthcare Med Tech Service Fhir Destination. Changing this forces a new Healthcare Med Tech Service Fhir Destination to be created.

* `medtech_service_id` - (Required) Specifies the name of the Healthcare Med Tech Service where the Healthcare Med Tech Service Fhir Destination should exist. Changing this forces a new Healthcare Med Tech Service Fhir Destination to be created.

* `location` - (Required) Specifies the Azure Region where the Healthcare Med Tech Service Fhir Destination should be created. Changing this forces a new Healthcare Med Tech Service Fhir Destination to be created.

* `destination_identity_resolution_type` - (Required) Specifies the destination identity resolution type where the Healthcare Med Tech Service Fhir Destination should be created. Possible values are `Create`, `Lookup`.

* `destination_fhir_mapping_json` - (Required) Specifies the destination Fhir mappings of the Med Tech Service Fhir Destination.

## Attributes Reference

The following arguments are supported:

* `id` - The ID of the Healthcare Med Tech Service Fhir Destination.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the Healthcare Med Tech Service Fhir Destination.
* `update` - (Defaults to 30 minut es) Used when updating the Healthcare Med Tech Service Fhir Destination.
* `read` - (Defaults to 5 minutes) Used when retrieving the Healthcare Med Tech Service Fhir Destination.
* `delete` - (Defaults to 90 minutes) Used when deleting the Healthcare Med Tech Service Fhir Destination.
* `update` - (Defaults to 90 minutes) Used when updating the Healthcare Medtech Service Fhir Destination.

## Import

Healthcare Med Tech Service Fhir Destination can be imported using the resource`id`, e.g.

```shell
terraform import azurerm_healthcare_medtech_service_fhir_destination.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/iotConnectors/iotconnector1/fhirDestinations/destination1
```
