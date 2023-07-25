package orders

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OrderProperties struct {
	ContactInformation   ContactDetails  `json:"contactInformation"`
	CurrentStatus        *OrderStatus    `json:"currentStatus,omitempty"`
	DeliveryTrackingInfo *[]TrackingInfo `json:"deliveryTrackingInfo,omitempty"`
	OrderHistory         *[]OrderStatus  `json:"orderHistory,omitempty"`
	OrderId              *string         `json:"orderId,omitempty"`
	ReturnTrackingInfo   *[]TrackingInfo `json:"returnTrackingInfo,omitempty"`
	SerialNumber         *string         `json:"serialNumber,omitempty"`
	ShipmentType         *ShipmentType   `json:"shipmentType,omitempty"`
	ShippingAddress      *Address        `json:"shippingAddress,omitempty"`
}
