---
subcategory: "Healthcare"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_healthcare_medtech_service"
description: |-
  Get information about an existing Healthcare Med Tech Service
---

# Data Source: azurerm_healthcare_medtech_service

Use this data source to access information about an existing Healthcare Med Tech Service

## Example Usage

```hcl
data "azurerm_healthcare_medtech_service" "example" {
  name         = "tfexmedtech"
  workspace_id = "tfexwks"
}

output "azurerm_healthcare_medtech_service_id" {
  value = data.azurerm_healthcare_medtech_service.example.id
}
```

## Argument Reference

* `name` - The name of the Healthcare Med Tech Service.

* `workspace_id` - The id of the Healthcare Workspace in which the Healthcare Med Tech Service exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Healthcare Med Tech Service.

* `identity` - The `identity` block as defined below.

* `eventhub_namespace_name` - The namespace name of the Event Hub of the Healthcare Med Tech Service.

* `eventhub_name` - The name of the Event Hub of the Healthcare Med Tech Service.

* `eventhub_consumer_group_name` - The Consumer Group of the Event Hub of the Healthcare Med Tech Service.

* `device_mapping_json` - The Device Mappings of the Med Tech Service.

---
An `identity` block supports the following:

* `type` The type of identity used for the Healthcare Med Tech Service. Possible values are `SystemAssigned`.

* `principal_id` - The Principal ID of the System Assigned Managed Service Identity that is configured on this Healthcare Med Tech Service.

* `tenant_id` - The Tenant ID of the System Assigned Managed Service Identity that is configured on this Healthcare Med Tech Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Healthcare Med Tech Service.
