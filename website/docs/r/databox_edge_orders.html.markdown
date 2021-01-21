---
subcategory: "Databox Edge"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_databox_edge_order"
description: |-
  Manages a Databox Edge Order.
---

# azurerm_databox_edge_order

Manages a Databox Edge Order.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-databoxedge"
  location = "West Europe"
}

resource "azurerm_databox_edge_order" "example" {
  resource_group_name = azurerm_resource_group.example.name
  device_name = "somename"
}
```

## Arguments Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the Resource Group where the Databox Edge Order should exist. Changing this forces a new Databox Edge Order to be created.

* `device_name` - (Required) The order details of a device. Changing this forces a new Databox Edge Order to be created.

---

* `contact_information` - (Optional)  A `contact_information` block as defined below.

* `current_status` - (Optional)  A `current_status` block as defined below.

* `shipping_address` - (Optional)  A `shipping_address` block as defined below.

---

An `contact_information` block exports the following:

* `company_name` - (Required) The name of the company.

* `contact_person` - (Required) The contact person name.

* `emails` - (Required) The email list.

* `phone_number` - (Required) The phone number.

---

An `current_status` block exports the following:

* `comments` - (Optional) Comments related to this status change.

---

An `shipping_address` block exports the following:

* `address_line1` - (Required) The address line1.

* `city` - (Required) The city name.

* `country` - (Required) The country name.

* `postal_code` - (Required) The postal code.

* `state` - (Required) The state name.

---

* `address_line2` - (Optional) The address line2.

* `address_line3` - (Optional) The address line3.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Databox Edge Order.

* `name` - The Name of this Databox Edge Order.

* `delivery_tracking_info` - Tracking information for the package delivered to the customer whether it has an original or a replacement device. A `delivery_tracking_info` block as defined below.

* `order_history` - List of status changes in the order. A `order_history` block as defined below.

* `return_tracking_info` - Tracking information for the package returned from the customer whether it has an original or a replacement device. A `return_tracking_info` block as defined below.

* `serial_number` - Serial number of the device.

* `type` - The hierarchical type of the object.

---

An `delivery_tracking_info` block exports the following:

* `carrier_name` - Name of the carrier used in the delivery.

* `serial_number` - Serial number of the device being tracked.

* `tracking_id` - The ID of the tracking.

* `tracking_url` - Tracking URL of the shipment.

---

An `order_history` block exports the following:

* `additional_order_details` - Dictionary to hold generic information which is not stored
by the already existing properties.

* `comments` - Comments related to this status change.

* `update_date_time` - Time of status update.

---

An `return_tracking_info` block exports the following:

* `carrier_name` - Name of the carrier used in the delivery.

* `serial_number` - Serial number of the device being tracked.

* `tracking_id` - The ID of the tracking.

* `tracking_url` - Tracking URL of the shipment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Databox Edge Order.
* `read` - (Defaults to 5 minutes) Used when retrieving the Databox Edge Order.
* `delete` - (Defaults to 30 minutes) Used when deleting the Databox Edge Order.

## Import

Databox Edge Orders can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_databoxedge_order.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataBoxEdge/dataBoxEdgeDevices/device1/orders/default
```