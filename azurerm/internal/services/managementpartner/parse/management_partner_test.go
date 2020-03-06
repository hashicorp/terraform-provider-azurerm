package parse

import (
	"testing"
)

func TestParseVirtualHubConnection(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *ManagementPartnerId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No PartnerId Segment",
			Input:    "/providers/Microsoft.ManagementPartner/",
			Expected: nil,
		},
		{
			Name:     "No PartnerId Value",
			Input:    "/providers/Microsoft.ManagementPartner/partners/",
			Expected: nil,
		},
		{
			Name:  "Completed",
			Input: "/providers/Microsoft.ManagementPartner/partners/5127255",
			Expected: &ManagementPartnerId{
				PartnerId: "5127255",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := ParseManagementPartnerID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.PartnerId != v.Expected.PartnerId {
			t.Fatalf("Expected %q but got %q for PartnerId", v.Expected.PartnerId, actual.PartnerId)
		}
	}
}
