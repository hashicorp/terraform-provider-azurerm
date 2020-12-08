package parse

import (
	"testing"
)

func TestDataBoxJobID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *DataBoxJobId
	}{
		{
			Input: "",
			Error: true,
		},
		{
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000",
			Error: true,
		},
		{
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups",
			Error: true,
		},
		{
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/hello",
			Error: true,
		},
		{
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/hello/providers",
			Error: true,
		},
		{
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/hello/providers/Microsoft.DataBox",
			Error: true,
		},
		{
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/hello/providers/Microsoft.DataBox/jobs",
			Error: true,
		},
		{
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/hello/providers/Microsoft.DataBox/jobs/job1",
			Expected: &DataBoxJobId{
				Name:          "job1",
				ResourceGroup: "hello",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := DataBoxJobID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Data Box Job Name", v.Expected.Name, actual.Name)
		}
	}
}
