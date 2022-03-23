---
subcategory: "Healthcare"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_healthcare_iot_connector"
description: |-
  Get information about an existing Healthcare Iot Connector
---

# Data Source: azurerm_healthcare_iot_connector

Use this data source to access information about an existing Healthcare Iot Connector

## Example Usage

```hcl
data "azurerm_healthcare_iot_connector" "example" {
  name = "tfexiot"
  workspace_id = "tfexwks"
}

output "azurerm_healthcare_iot_connector_id" {
  value = data.azurerm_healthcare_iot_connector.example.id
}
```
## Argument Reference

* `name` - The name of the Healthcare Iot Connector.

* `workspace_id` - The id of the Healthcare Workspace in which the Healthcare Iot Connector exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Healthcare Iot Connector.

* `identity` - The `identity` block as defined below.

* `eventhub_namespace_name` - The namespace name of the Event Hub of the Healthcare Iot Connector.

* `eventhub_name` - The name of the Event Hub of the Healthcare Iot Connector.

* `eventhub_consumer_group_name` - The Consumer Group of the Event Hub of the Healthcare Iot Connector.

* `device_mapping` - The Device Mappings of the Iot Connector.

---
An `identity` block supports the following:

* `type` The type of identity used for the Healthcare Iot Connector. Possible values are `SystemAssigned`.

* `principal_id` - The Principal ID of the System Assigned Managed Service Identity that is configured on this Healthcare Iot Connector.

* `tenant_id` - The Tenant ID of the System Assigned Managed Service Identity that is configured on this Healthcare Iot Connector.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Healthcare Iot Connector.
