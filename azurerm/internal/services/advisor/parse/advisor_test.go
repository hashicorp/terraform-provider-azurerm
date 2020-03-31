package parse

import (
	"testing"
)

func TestAdvisorRecommendationID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *AdvisorRecommendationId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "Missing Recommendations",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000",
			Expected: nil,
		},
		{
			Name:     "Missing Recommendations 2",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/microsoft.web/sites/site1",
			Expected: nil,
		},
		{
			Name:     "Invalid Resource Uri",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/providers/microsoft.web/sites/site1/providers/Microsoft.Advisor/recommendations/recommendation1",
			Expected: nil,
		},
		{
			Name:     "Missing Recommendations 3",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/microsoft.web/sites/site1/providers/Microsoft.Advisor/",
			Expected: nil,
		},
		{
			Name:     "Missing Recommendations Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/microsoft.web/sites/site1/providers/Microsoft.Advisor/recommendations/",
			Expected: nil,
		},
		{
			Name:  "Advisor Recommendation 1",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/microsoft.web/sites/site1/providers/Microsoft.Advisor/recommendations/recommendation1",
			Expected: &AdvisorRecommendationId{
				Name:        "recommendation1",
				ResourceUri: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/microsoft.web/sites/site1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/microsoft.web/sites/site1/providers/Microsoft.Advisor/Recommendations/recommendation1",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := AdvisorRecommendationID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.ResourceUri != v.Expected.ResourceUri {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}

func TestAdvisorSuppressionID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *AdvisorSuppressionId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Expected: nil,
		},
		{
			Name:     "Resource Group ID",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/",
			Expected: nil,
		},
		{
			Name:     "Missing Recommendations",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/microsoft.web/sites/site1",
			Expected: nil,
		},
		{
			Name:     "Missing Recommendation Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/microsoft.web/sites/site1/providers/Microsoft.Advisor/recommendations/",
			Expected: nil,
		},
		{
			Name:     "Missing Suppression",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/microsoft.web/sites/site1/providers/Microsoft.Advisor/recommendations/recommendation1",
			Expected: nil,
		},
		{
			Name:     "Missing Suppression Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/microsoft.web/sites/site1/providers/Microsoft.Advisor/recommendations/recommendation1/suppressions",
			Expected: nil,
		},
		{
			Name:  "Advisor Suppression ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/microsoft.web/sites/site1/providers/Microsoft.Advisor/recommendations/recommendation1/suppressions/suppression1",
			Expected: &AdvisorSuppressionId{
				Name:               "suppression1",
				RecommendationName: "recommendation1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/microsoft.web/sites/site1/providers/Microsoft.Advisor/recommendations/recommendation1/Suppressions/suppression1",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := AdvisorSuppressionID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.RecommendationName != v.Expected.RecommendationName {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.RecommendationName, actual.RecommendationName)
		}
	}
}
