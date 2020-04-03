package parse

import (
	"testing"
)

func TestPeerAsnID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *PeerAsnId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "No Resource Groups Segment",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000",
			Error: true,
		},
		{
			Name:  "No Provider Name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers",
			Error: true,
		},
		{
			Name:  "No PeerAsn component",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Peering",
			Error: true,
		},
		{
			Name:  "No PeerAsn name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Peering/peerAsns",
			Error: true,
		},
		{
			Name:  "Unexpected casing of peerasn segment",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Peering/peerasns/peerasn1",
			Error: true,
		},
		{
			Name:  "Complete id",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Peering/peerAsns/peerasn1",
			Expect: &PeerAsnId{
				Name: "peerasn1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := PeerAsnID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expect.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.Name, actual.Name)
		}
	}
}
