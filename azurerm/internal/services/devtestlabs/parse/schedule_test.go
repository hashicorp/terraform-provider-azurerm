package parse

import "testing"

func TestGlobalSchedule(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *ScheduleId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No DevTest Lab segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/",
			Expected: nil,
		},
		{
			Name:     "No schedule name",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/Microsoft.DevTestLab/schedules/",
			Expected: nil,
		},
		{
			Name:     "Case incorrect in path element",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/Microsoft.DevTestLab/Schedules/schedule1",
			Expected: nil,
		},
		{
			Name:  "Valid",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/Microsoft.DevTestLab/schedules/schedule1",
			Expected: &ScheduleId{
				ResourceGroup: "myGroup1",
				Name:          "schedule1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := ScheduleID(v.Input)
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
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
