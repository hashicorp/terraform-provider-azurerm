package parse

import (
	"testing"
)

func TestSecurityCenterSettingID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *SecurityCenterSettingId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "No Settings Segment",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000",
			Error: true,
		},
		{
			Name:  "No Settings Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/settings/",
			Error: true,
		},
		{
			Name:  "Security Center Setting ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/settings/MCAS",
			Expect: &SecurityCenterSettingId{
				SettingName: "MCAS",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := SecurityCenterSettingID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.SettingName != v.Expect.SettingName {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.SettingName, actual.SettingName)
		}
	}
}
