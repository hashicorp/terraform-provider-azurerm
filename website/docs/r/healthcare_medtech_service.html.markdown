---
subcategory: "Healthcare"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_healthcare_medtech_service"
description: |-
  Manages a Healthcare MedTech (Internet of Medical Things) devices Service.
---

# azurerm_healthcare_medtech_service

Manages a Healthcare Med Tech Service.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "east us"
}

resource "azurerm_healthcare_workspace" "example" {
  name                = "examplewkspace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_healthcare_medtech_service" "example" {
  name         = "examplemed"
  workspace_id = azurerm_healthcare_workspace.example.id
  location     = "east us"

  identity {
    type = "SystemAssigned"
  }

  eventhub_namespace_name      = "example-eventhub-namespace"
  eventhub_name                = "example-eventhub"
  eventhub_consumer_group_name = "$Default"

  device_mapping_json = jsonencode({
    "templateType" : "CollectionContent",
    "template" : [
      {
        "templateType" : "JsonPathContent",
        "template" : {
          "typeName" : "heartrate",
          "typeMatchExpression" : "$..[?(@heartrate)]",
          "deviceIdExpression" : "$.deviceid",
          "timestampExpression" : "$.measurementdatetime",
          "values" : [
            {
              "required" : "true",
              "valueExpression" : "$.heartrate",
              "valueName" : "hr"
            }
          ]
        }
      }
    ]
  })
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Healthcare Med Tech Service. Changing this forces a new Healthcare Med Tech Service to be created.

* `workspace_id` - (Required) Specifies the id of the Healthcare Workspace where the Healthcare Med Tech Service should exist. Changing this forces a new Healthcare Med Tech Service to be created.

* `location` - (Required) Specifies the Azure Region where the Healthcare Med Tech Service should be created. Changing this forces a new Healthcare Med Tech Service to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `eventhub_namespace_name` - (Required) Specifies the namespace name of the Event Hub to connect to.

* `eventhub_name` - (Required) Specifies the name of the Event Hub to connect to.

* `eventhub_consumer_group_name` - (Required) Specifies the Consumer Group of the Event Hub to connect to.

* `device_mapping_json` - (Required) Specifies the Device Mappings of the Med Tech Service.

* `tags` - (Optional) A mapping of tags to assign to the Healthcare Med Tech Service.

---
A `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Healthcare Med Tech Service. Possible values are `SystemAssigned`.

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Healthcare Med Tech Service.

## Attributes Reference

The following arguments are supported:

* `id` - The ID of the Healthcare Med Tech Service.

*`identity` - An `identity` block as defined below.

---
An `identity` block exports the following:

* `type` - (Required) The type of identity used for the Healthcare Med Tech service.

* `principal_id` - The Principal ID associated with this System Assigned Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this System Assigned Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the Healthcare Med Tech Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Healthcare Med Tech Service.
* `update` - (Defaults to 90 minutes) Used when updating the Healthcare Med Tech Service.
* `delete` - (Defaults to 90 minutes) Used when deleting the Healthcare Med Tech Service.

## Import

Healthcare Med Tech Service can be imported using the resource`id`, e.g.

```shell
terraform import azurerm_healthcare_medtech_service.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/iotConnectors/iotconnector1
```
