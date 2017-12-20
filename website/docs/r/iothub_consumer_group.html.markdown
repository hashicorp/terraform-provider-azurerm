---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iothub_consumer_group"
sidebar_current: "docs-azurerm-resource-iothub-consumer-group"
description: |-
    Adds a consumer group to an Event Hub-compatible endpoint in an IoT hub.
---

# azurerm\_iothub\_consumer\_group

Creates a new consumer group

## Example Usage

```hcl
resource "azurerm_resource_group" "foo" {
	name = "acctestIot-%d"
	location = "%s"
}

resource "azurerm_iothub" "bar" {
	name = "acctestiothub-%d"
	location = "${azurerm_resource_group.foo.location}"
	resource_group_name = "${azurerm_resource_group.foo.name}"
	sku {
		name = "S1"
		tier = "Standard"
		capacity = "1"
	}

	tags {
		"purpose" = "testing"
	}
}

resource "azurerm_iothub_consumer_group" "foo" {
	name = "acctestiothubgroup-%d"
	resource_group_name = "${azurerm_resource_group.foo.location}"
	iothub_name = "${azurerm_iothub.bar.name}"
	event_hub_endpoint = "test"
}


```

## Argument Reference

The following arguments are supported:

* `name` (Required) Specifies the name of the Consumer Group resource. Changing this forces a new resource to be created.

* `resource_group_name` (Required) The name of the resource group that contains the IoT hub.

* `iothub_name` (Required) The name of the IoT hub.

* `event_hub_endpoint` (Required) The name of the Event Hub-compatible endpoint in the IoT hub.

## Attributes Reference

The following attributes are exported:

WIP

