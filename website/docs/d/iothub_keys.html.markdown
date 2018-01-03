---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iothub_keys"
sidebar_current: "docs-azurerm-datasource-iothub-keys"
description: |-
  Get shared access policy
---

# azurerm\_iothub_keys

Get a shared access policy by name from an IoT hub.

## Example Usage

```hcl
data "azurerm_iothub_keys" "test" {
	name = "Shared Access Policy"
	resource_group_name = "ResourceGroupTest"
	iot_hub_name = "IoTHub"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the shared access policy which contains access keys and permissions.

* `resource_group_name` - (Required) The name of the resource group that contains the IoT hub.

* `iot_hub_name` - (Required) The name of the IoT hub.

## Attributes Reference

The following attributes are exported:

* `primary_key` - The primary key.

* `secondary_key` - The secondary key.

* `permissions` - The permissions assigned to the shared access policy.
