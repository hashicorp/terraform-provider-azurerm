package parse

import (
	"testing"
)

func TestEventGridSystemTopicId(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *SystemTopicId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No System Topic",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/providers/Microsoft.EventGrid",
			Expected: nil,
		},
		{
			Name:  "EventGrid System Topic ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.EventGrid/systemTopics/systemTopic1",
			Expected: &SystemTopicId{
				Name:          "systemTopic1",
				ResourceGroup: "resGroup1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := SystemTopicID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
