---
subcategory: "Healthcare"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_healthcare_iot_fhir_destination"
description: |-
  Manages a Healthcare Iot Connector Fhir Destination.
---

# azurerm_healthcare_iot_fhir_Destination

Manages a Healthcare Iot Connector Fhir Destination.

```hcl
resource "azurerm_healthcare_iot_fhir_destination" "test" {
  name                                 = "tfexiotdes"
  location                             = "east us"
  iot_connector_id                     = "iotconnector_id"
  destination_fhir_service_id          = "fhir_service_id"
  destination_identity_resolution_type = "Create"
  destination_fhir_mapping             = <<JSON
  {
    "content": {
              "templateType": "CollectionFhirTemplate",
              "template": []
  }
  JSON
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Healthcare Iot Connector Fhir Destination. Changing this forces a new Healthcare Iot Connector Fhir Destination to be created.

* `iot_connector_id`  - (Required) Specifies the name of the Healthcare Iot Connector where the Healthcare Iot Connector Fhir Destination should exist. Changing this forces a new Healthcare Iot Connector Fhir Destination to be created.

* `location` - (Required) Specifies the Azure Region where the Healthcare Iot Connector Fhir Destination should be created. Changing this forces a new Healthcare Iot Connector Fhir Destination to be created.

* `destination_identity_resolution_type` - (Required) Specifies the Destination Identity Resolution Type where the Healthcare Iot Connector Fhir Destination should be created. Possible values are `Create`, `Lookup`. Defaults to `Create`.

* `destination_fhir_mapping` - (Required) Specifies the Destination Fhir Mappings of the Iot Connector Fhir Destination.

## Attributes Reference

The following arguments are supported:

* `id` - The ID of the Healthcare Iot Connector Fhir Destination.

## Timeouts
The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Healthcare Iot Connector Fhir Destination.
* `update` - (Defaults to 30 minut es) Used when updating the Healthcare Iot Connector Fhir Destination.
* `read` - (Defaults to 5 minutes) Used when retrieving the Healthcare Iot Connector Fhir Destination.
* `delete` - (Defaults to 30 minutes) Used when deleting the Healthcare Iot Connector Fhir Destination.

## Import

Healthcare Iot Connector Fhir Destination can be imported using the resource`id`, e.g.

```shell
terraform import azurerm_healthcare_iot_fhir_destination.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/iotconnectors/iotconnector1/fhirdestinations/destination1
```
