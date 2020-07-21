package parse

import (
	"testing"
)

func TestApplicationInsightsID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *ApplicationInsightsId
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
			Name:  "No Resource Groups Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Error: true,
		},
		{
			Name:  "No Provider",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers",
			Error: true,
		},
		{
			Name:  "No component",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/microsoft.insights/components",
			Error: true,
		},
		{
			Name:  "Correct",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/microsoft.insights/components/appinsights1",
			Expect: &ApplicationInsightsId{
				ResourceGroup: "group1",
				Name:          "appinsights1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := ApplicationInsightsID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get")
		}

		if actual.ResourceGroup != v.Expect.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expect.ResourceGroup, actual.ResourceGroup)
		}

		if actual.Name != v.Expect.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.Name, actual.Name)
		}
	}
}

func TestApplicationInsightsWebTestID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *ApplicationInsightsWebTestId
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
			Name:  "No Resource Groups Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Error: true,
		},
		{
			Name:  "No Provider",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers",
			Error: true,
		},
		{
			Name:  "No webtest",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/microsoft.insights/webtests",
			Error: true,
		},
		{
			Name:  "Correct",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/microsoft.insights/webtests/test1",
			Expect: &ApplicationInsightsWebTestId{
				ResourceGroup: "group1",
				Name:          "test1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := ApplicationInsightsWebTestID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get")
		}

		if actual.ResourceGroup != v.Expect.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expect.ResourceGroup, actual.ResourceGroup)
		}

		if actual.Name != v.Expect.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.Name, actual.Name)
		}
	}
}
