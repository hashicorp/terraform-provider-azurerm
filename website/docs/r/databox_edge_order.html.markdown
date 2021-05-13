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

resource "azurerm_databox_edge_device" "example" {
  name                = "example-device"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku_name = "EdgeP_Base-Standard"
}

resource "azurerm_databox_edge_order" "example" {
  resource_group_name = azurerm_resource_group.example.name
  device_name         = azurerm_databox_edge_device.example.name

  contact {
    company_name = "Contoso Corporation"
    name         = "Bart"
    email_lists  = ["bart@example.com"]
    phone        = "(800) 867-5309"
  }

  shipment_address {
    address     = ["740 Evergreen Terrace"]
    city        = "Springfield"
    country     = "United States"
    postal_code = "97403"
    state       = "OR"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the Resource Group where the Databox Edge Order should exist. Changing this forces a new Databox Edge Order to be created.

* `device_name` - (Required) The name of the Databox Edge Device this order is for. Changing this forces a new Databox Edge Order to be created.

* `contact` - (Required)  A `contact` block as defined below.

* `shipment_address` - (Required)  A `shipment_address block as defined below.

---

An `contact` block includes the following:

* `company_name` - (Required) The name of the company. Changing this forces a new Databox Edge Order to be created.

* `name` - (Required) The contact person name. Changing this forces a new Databox Edge Order to be created.

* `emails` - (Required) A list of email address to send order notification to. Changing this forces a new Databox Edge Order to be created.

* `phone_number` - (Required) The phone number. Changing this forces a new Databox Edge Order to be created.

---

An `shipment_address` block includes the following:

* `address` - (Required) The list of upto 3 lines for address information. Changing this forces a new Databox Edge Order to be created.

* `city` - (Required) The city name. Changing this forces a new Databox Edge Order to be created.

* `country` - (Required) The name of the country to ship the Databox Edge Device to. Valid values are "Algeria", "Argentina", "Australia", "Austria", "Bahamas", "Bahrain", "Bangladesh", "Barbados", "Belgium", "Bermuda", "Bolivia", "Bosnia and Herzegovina", "Brazil", "Bulgaria", "Canada", "Cayman Islands", "Chile", "Colombia", "Costa Rica", "Croatia", "Cyprus", "Czechia", "CÃ´te D'ivoire", "Denmark", "Dominican Republic", "Ecuador", "Egypt", "El Salvador", "Estonia", "Ethiopia", "Finland", "France", "Georgia", "Germany", "Ghana", "Greece", "Guatemala", "Honduras", "Hong Kong SAR", "Hungary", "Iceland", "India", "Indonesia", "Ireland", "Israel", "Italy", "Jamaica", "Japan", "Jordan", "Kazakhstan", "Kenya", "Kuwait", "Kyrgyzstan", "Latvia", "Libya", "Liechtenstein", "Lithuania", "Luxembourg", "Macao SAR", "Malaysia", "Malta", "Mauritius", "Mexico", "Moldova", "Monaco", "Mongolia", "Montenegro", "Morocco", "Namibia", "Nepal", "Netherlands", "New Zealand", "Nicaragua", "Nigeria", "Norway", "Oman", "Pakistan", "Palestinian Authority", "Panama", "Paraguay", "Peru", "Philippines", "Poland", "Portugal", "Puerto Rico", "Qatar", "Republic of Korea", "Romania", "Russia", "Rwanda", "Saint Kitts And Nevis", "Saudi Arabia", "Senegal", "Serbia", "Singapore", "Slovakia", "Slovenia", "South Africa", "Spain", "Sri Lanka", "Sweden", "Switzerland", "Taiwan", "Tajikistan", "Tanzania", "Thailand", "Trinidad And Tobago", "Tunisia", "Turkey", "Turkmenistan", "U.S. Virgin Islands", "Uganda", "Ukraine", "United Arab Emirates", "United Kingdom", "United States", "Uruguay", "Uzbekistan", "Venezuela", "Vietnam", "Yemen", "Zambia" or "Zimbabwe". Changing this forces a new Databox Edge Order to be created.

* `postal_code` - (Required) The postal code. Changing this forces a new Databox Edge Order to be created.

* `state` - (Required) The name of the state to ship the Databox Edge Device to. Changing this forces a new Databox Edge Order to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Databox Edge Order.

* `name` - The Name of this Databox Edge Order.

* `shipment_tracking` - Tracking information for the package delivered to the customer whether it has an original or a replacement device. A `shipment_tracking` block as defined below.

* `status` - The current status of the order. A `status` block as defined below.

* `shipment_history` - List of status changes in the order. A `shipment_history` block as defined below.

* `return_tracking` - Tracking information for the package returned from the customer whether it has an original or a replacement device. A `return_tracking` block as defined below.

* `serial_number` - Serial number of the device.

---

A `shipment_tracking` block exports the following:

* `carrier_name` - Name of the carrier used in the delivery.

* `serial_number` - Serial number of the device being tracked.

* `tracking_id` - The ID of the tracking.

* `tracking_url` - Tracking URL of the shipment.

---

A `status` block exports the following:

 * `info` - The current status of the order. Possible values include `Untracked`, `AwaitingFulfilment`, `AwaitingPreparation`, `AwaitingShipment`, `Shipped`, `Arriving`, `Delivered`, `ReplacementRequested`, `LostDevice`, `Declined`, `ReturnInitiated`, `AwaitingReturnShipment`, `ShippedBack` or `CollectedAtMicrosoft`.

* `additional_details` - Dictionary to hold generic information which is not stored by the already existing properties.

* `comments` - Comments related to this status change.

* `last_update` - Time of status update.

---

A `shipment_history` block exports the following:

* `additional_details` - Dictionary to hold generic information which is not stored by the already existing properties.

* `comments` - Comments related to this status change.

* `last_update` - Time of status update.

---

A `return_tracking` block exports the following:

* `carrier_name` - Name of the carrier used in the delivery.

* `serial_number` - Serial number of the device being tracked.

* `tracking_id` - The ID of the tracking.

* `tracking_url` - Tracking URL of the shipment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating a Databox Edge Order.
* `read` - (Defaults to 5 minutes) Used when retrieving a Databox Edge Order.
* `update` - (Defaults to 30 minutes) Used when updating a Databox Edge Order.
* `delete` - (Defaults to 30 minutes) Used when deleting a Databox Edge Order.

## Import

Databox Edge Orders can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_databoxedge_order.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataBoxEdge/dataBoxEdgeDevices/device1/orders/default
```
