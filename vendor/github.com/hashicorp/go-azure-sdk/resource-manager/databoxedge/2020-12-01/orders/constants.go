package orders

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OrderState string

const (
	OrderStateArriving               OrderState = "Arriving"
	OrderStateAwaitingDrop           OrderState = "AwaitingDrop"
	OrderStateAwaitingFulfilment     OrderState = "AwaitingFulfilment"
	OrderStateAwaitingPickup         OrderState = "AwaitingPickup"
	OrderStateAwaitingPreparation    OrderState = "AwaitingPreparation"
	OrderStateAwaitingReturnShipment OrderState = "AwaitingReturnShipment"
	OrderStateAwaitingShipment       OrderState = "AwaitingShipment"
	OrderStateCollectedAtMicrosoft   OrderState = "CollectedAtMicrosoft"
	OrderStateDeclined               OrderState = "Declined"
	OrderStateDelivered              OrderState = "Delivered"
	OrderStateLostDevice             OrderState = "LostDevice"
	OrderStatePickupCompleted        OrderState = "PickupCompleted"
	OrderStateReplacementRequested   OrderState = "ReplacementRequested"
	OrderStateReturnInitiated        OrderState = "ReturnInitiated"
	OrderStateShipped                OrderState = "Shipped"
	OrderStateShippedBack            OrderState = "ShippedBack"
	OrderStateUntracked              OrderState = "Untracked"
)

func PossibleValuesForOrderState() []string {
	return []string{
		string(OrderStateArriving),
		string(OrderStateAwaitingDrop),
		string(OrderStateAwaitingFulfilment),
		string(OrderStateAwaitingPickup),
		string(OrderStateAwaitingPreparation),
		string(OrderStateAwaitingReturnShipment),
		string(OrderStateAwaitingShipment),
		string(OrderStateCollectedAtMicrosoft),
		string(OrderStateDeclined),
		string(OrderStateDelivered),
		string(OrderStateLostDevice),
		string(OrderStatePickupCompleted),
		string(OrderStateReplacementRequested),
		string(OrderStateReturnInitiated),
		string(OrderStateShipped),
		string(OrderStateShippedBack),
		string(OrderStateUntracked),
	}
}

func parseOrderState(input string) (*OrderState, error) {
	vals := map[string]OrderState{
		"arriving":               OrderStateArriving,
		"awaitingdrop":           OrderStateAwaitingDrop,
		"awaitingfulfilment":     OrderStateAwaitingFulfilment,
		"awaitingpickup":         OrderStateAwaitingPickup,
		"awaitingpreparation":    OrderStateAwaitingPreparation,
		"awaitingreturnshipment": OrderStateAwaitingReturnShipment,
		"awaitingshipment":       OrderStateAwaitingShipment,
		"collectedatmicrosoft":   OrderStateCollectedAtMicrosoft,
		"declined":               OrderStateDeclined,
		"delivered":              OrderStateDelivered,
		"lostdevice":             OrderStateLostDevice,
		"pickupcompleted":        OrderStatePickupCompleted,
		"replacementrequested":   OrderStateReplacementRequested,
		"returninitiated":        OrderStateReturnInitiated,
		"shipped":                OrderStateShipped,
		"shippedback":            OrderStateShippedBack,
		"untracked":              OrderStateUntracked,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OrderState(input)
	return &out, nil
}

type ShipmentType string

const (
	ShipmentTypeNotApplicable     ShipmentType = "NotApplicable"
	ShipmentTypeSelfPickup        ShipmentType = "SelfPickup"
	ShipmentTypeShippedToCustomer ShipmentType = "ShippedToCustomer"
)

func PossibleValuesForShipmentType() []string {
	return []string{
		string(ShipmentTypeNotApplicable),
		string(ShipmentTypeSelfPickup),
		string(ShipmentTypeShippedToCustomer),
	}
}

func parseShipmentType(input string) (*ShipmentType, error) {
	vals := map[string]ShipmentType{
		"notapplicable":     ShipmentTypeNotApplicable,
		"selfpickup":        ShipmentTypeSelfPickup,
		"shippedtocustomer": ShipmentTypeShippedToCustomer,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ShipmentType(input)
	return &out, nil
}
