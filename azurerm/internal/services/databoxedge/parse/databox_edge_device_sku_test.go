package parse

import (
	"testing"
)

func TestDataboxEdgeDeviceSkuName(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *DataboxEdgeDeviceSku
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "Delimiter only",
			Input:    "-",
			Expected: nil,
		},
		{
			Name:     "Name Segment is a space with delimiter and Tier",
			Input:    " -Standard",
			Expected: nil,
		},
		{
			Name:     "Name Segment is a space with delimiter only",
			Input:    " -",
			Expected: nil,
		},
		{
			Name:     "Name Segment is a space without delimiter",
			Input:    " ",
			Expected: nil,
		},
		{
			Name:     "No Name Segment with delimiter",
			Input:    "-Standard",
			Expected: nil,
		},
		{
			Name:  "No Tier Segment without delimiter",
			Input: "Edge",
			Expected: &DataboxEdgeDeviceSku{
				Name: "Edge",
				Tier: "Standard",
			},
		},
		{
			Name:  "No Tier Segment with delimiter",
			Input: "Edge-",
			Expected: &DataboxEdgeDeviceSku{
				Name: "Edge",
				Tier: "Standard",
			},
		},
		{
			Name:  "Databox Edge Device Sku",
			Input: "Edge-Standard",
			Expected: &DataboxEdgeDeviceSku{
				Name: "Edge",
				Tier: "Standard",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.Name)

		actual, err := DataboxEdgeDeviceSkuName(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}
			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.Tier != v.Expected.Tier {
			t.Fatalf("Expected %q but got %q for Tier", v.Expected.Tier, actual.Tier)
		}
	}
}
