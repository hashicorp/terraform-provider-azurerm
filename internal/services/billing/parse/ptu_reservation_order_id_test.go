// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package parse

import "testing"

func TestPtuReservationOrderIDFormatter(t *testing.T) {
	actual := NewPtuReservationOrderId("10146e93-a8f3-458a-83ce-7e88f2300632").ID()
	expected := "/providers/Microsoft.Capacity/reservationOrders/10146e93-a8f3-458a-83ce-7e88f2300632"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestPtuReservationOrderID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *PtuReservationOrderId
	}{
		{
			// empty
			Input: "",
			Error: true,
		},
		{
			// wrong provider namespace
			Input: "/providers/Microsoft.Billing/reservationOrders/10146e93-a8f3-458a-83ce-7e88f2300632",
			Error: true,
		},
		{
			// missing order ID value
			Input: "/providers/Microsoft.Capacity/reservationOrders/",
			Error: true,
		},
		{
			// extra path segment after UUID
			Input: "/providers/Microsoft.Capacity/reservationOrders/10146e93-a8f3-458a-83ce-7e88f2300632/reservations/abc",
			Error: true,
		},
		{
			// not a UUID
			Input: "/providers/Microsoft.Capacity/reservationOrders/not-a-uuid",
			Error: true,
		},
		{
			// valid
			Input: "/providers/Microsoft.Capacity/reservationOrders/10146e93-a8f3-458a-83ce-7e88f2300632",
			Expected: &PtuReservationOrderId{
				OrderId: "10146e93-a8f3-458a-83ce-7e88f2300632",
			},
		},
		{
			// valid — lowercase prefix (Azure API sometimes returns lowercase IDs)
			Input: "/providers/microsoft.capacity/reservationorders/10146e93-a8f3-458a-83ce-7e88f2300632",
			Expected: &PtuReservationOrderId{
				OrderId: "10146e93-a8f3-458a-83ce-7e88f2300632",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := PtuReservationOrderID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}
			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.OrderId != v.Expected.OrderId {
			t.Fatalf("Expected OrderId %q but got %q", v.Expected.OrderId, actual.OrderId)
		}
	}
}
