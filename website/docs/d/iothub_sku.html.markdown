---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iothub_sku"
sidebar_current: "docs-azurerm-datasource-iothub-sku"
description: |-
  Get sku
---

# azurerm\_iothub_sku

Get the list of valid SKUs for an IoT hub.

## Example Usage

```hcl
data "azurerm_iothub_sku" "test" {
	resource_group_name = "ResourceGroupTest"
	iot_hub_name = "IoTHub"
}

```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group that contains the IoT hub.

* `iot_hub_name` - (Required) The name of the IoT hub.

## Attributes Reference

The following attributes are exported:

* `tier` - The billing tier for the IoT hub.

* `name` - The name of the SKU.
