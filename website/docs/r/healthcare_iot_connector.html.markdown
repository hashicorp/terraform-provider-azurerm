---
subcategory: "Healthcare"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_healthcare_iot_connector"
description: |-
  Manages a Healthcare Iot Connector.
---

# azurerm_healthcare_iot_connector

Manages a Healthcare Iot Connector.

## Example Usage

```hcl
resource "azurerm_healthcare_iot_connector" "test" {
  name         = "tftest"
  workspace_id = "tfex-workspace_id"
  location     = "east us"
  identity {
    type = "SystemAssigned"
  }
  eventhub_namespace_name      = "tfex-eventhub-namespace.name"
  eventhub_name                = "tfex-eventhub.name"
  eventhub_consumer_group_name = "tfex-eventhub-consumer-group.name"
  device_mapping               = <<JSON
{
"templateType": "CollectionContent",
"template": []
}
JSON
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Healthcare Iot Connector. Changing this forces a new Healthcare Iot Connector to be created.

* `workspace_id`  - (Required) Specifies the id of the Healthcare Workspace where the Healthcare Iot Connector should exist. Changing this forces a new Healthcare Iot Connector to be created.

* `location` - (Required) Specifies the Azure Region where the Healthcare Iot Connector should be created. Changing this forces a new Healthcare Iot Connector to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `eventhub_namespace_name` - (Required) Specifies the namespace name of the Event Hub to connect to.

* `eventhub_name` - (Required) Specifies the name of the Event Hub to connect to.

* `eventhub_consumer_group_name` - (Required) Specifies the Consumer Group of the Event Hub to connect to.

* `device_mapping` - (Required) Specifies the Device Mappings of the Iot Connector.

---
A `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Healthcare Iot Connector. Possible values are `SystemAssigned`.

## Attributes Reference

The following arguments are supported:

* `id` - The ID of the Healthcare Iot Connector.

*`identity` - An `identity` block as defined below.

---
An `identity` block exports the following:

* `type` The type of identity used for the Healthcare Fhir service.

* `principal_id` - The Principal ID associated with this System Assigned Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this System Assigned Managed Service Identity.

## Timeouts
The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Healthcare Iot Connector.
* `update` - (Defaults to 30 minutes) Used when updating the Healthcare Iot Connector.
* `read` - (Defaults to 5 minutes) Used when retrieving the Healthcare Iot Connector.
* `delete` - (Defaults to 30 minutes) Used when deleting the Healthcare Iot Connector.

## Import

Healthcare Iot Connector can be imported using the resource`id`, e.g.

```shell
terraform import azurerm_healthcare_iot_connector.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/iotconnectors/iotconnector1
```
