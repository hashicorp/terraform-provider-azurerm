// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
)

type PtuReservationOrderId struct {
	OrderId string
}

func NewPtuReservationOrderId(orderId string) PtuReservationOrderId {
	return PtuReservationOrderId{OrderId: orderId}
}

func (id PtuReservationOrderId) ID() string {
	return fmt.Sprintf("/providers/Microsoft.Capacity/reservationOrders/%s", id.OrderId)
}

func (id PtuReservationOrderId) String() string {
	return fmt.Sprintf("PTU Reservation Order %q", id.OrderId)
}

func ValidatePtuReservationOrderID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := PtuReservationOrderID(v); err != nil {
		errors = append(errors, fmt.Errorf("cannot parse %q as a PTU Reservation Order ID: %+v", k, err))
	}

	return
}

func PtuReservationOrderID(input string) (*PtuReservationOrderId, error) {
	if input == "" {
		return nil, fmt.Errorf("input was empty")
	}

	prefix := "/providers/Microsoft.Capacity/reservationOrders/"
	if !strings.HasPrefix(strings.ToLower(input), strings.ToLower(prefix)) {
		return nil, fmt.Errorf("expected ID to start with %q but got %q", prefix, input)
	}

	orderId := input[len(prefix):]
	if orderId == "" {
		return nil, fmt.Errorf("order ID segment was empty in %q", input)
	}

	if strings.Contains(orderId, "/") {
		return nil, fmt.Errorf("expected order ID to be a single UUID segment but got %q in %q", orderId, input)
	}

	if _, err := uuid.ParseUUID(orderId); err != nil {
		return nil, fmt.Errorf("expected order ID to be a valid UUID but got %q in %q", orderId, input)
	}

	return &PtuReservationOrderId{OrderId: orderId}, nil
}
