package helper

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
)

func TestExpandEntityStatus(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected servicebus.EntityStatus
	}{
		{
			Name:     "Active",
			Input:    "Active",
			Expected: servicebus.Active,
		},
		{
			Name:     "ReceiveDisabled",
			Input:    "ReceiveDisabled",
			Expected: servicebus.ReceiveDisabled,
		},
		{
			Name:     "Disabled",
			Input:    "Disabled",
			Expected: servicebus.Disabled,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual := ExpandEntityStatus(v.Input)

		if actual != v.Expected {
			t.Fatalf("Expected %q but got %q", v.Expected, actual)
		}
	}
}

func TestFlattenEntityStatus(t *testing.T) {
	testData := []struct {
		Name     string
		Input    servicebus.EntityStatus
		Expected string
	}{
		{
			Name:     "Active",
			Input:    servicebus.Active,
			Expected: "Active",
		},
		{
			Name:     "Disabled",
			Input:    servicebus.Disabled,
			Expected: "Disabled",
		},
		{
			Name:     "ReceiveDisabled",
			Input:    servicebus.ReceiveDisabled,
			Expected: "ReceiveDisabled",
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual := FlattenEntityStatus(v.Input)

		if *actual != v.Expected {
			t.Fatalf("Expected %q but got %q", v.Expected, *actual)
		}
	}
}
