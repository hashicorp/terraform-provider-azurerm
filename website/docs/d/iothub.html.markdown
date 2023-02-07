---
subcategory: "IoT Hub"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_iothub"
description: |-
  Gets information about an existing IoTHub.
---

# Data Source: azurerm_iothub

Use this data source to access information about an existing IoTHub.

## Example Usage

```hcl
data "azurerm_iothub" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_iothub.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this IoTHub.

* `resource_group_name` - (Required) The name of the Resource Group where the IoTHub exists.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the IoTHub.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the IoTHub.
* `hostname` - The Hostname of the IoTHub.
* `identity` - A `identity` block as defined below.

---

An `identity` block exports the following:

* `identity_ids` - The list of User Managed Identity IDs which are assigned to the Virtual Machine.

* `principal_id` - The ID of the System Managed Service Principal assigned to the Virtual Machine.

* `tenant_id` - The ID of the Tenant of the System Managed Service Principal assigned to the Virtual Machine.

* `type` - The identity type of the Managed Identity assigned to the Virtual Machine.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the IoTHub.
