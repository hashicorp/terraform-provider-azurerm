package parsers

import (
	"testing"
)

func TestHPCCacheTargetID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *StorageTargetId
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
			Name:  "Resource Group ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/",
			Error: true,
		},
		{
			Name:  "Missing Cache Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.StorageCache/caches/",
			Error: true,
		},
		{
			Name:  "With Cache Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.StorageCache/caches/cache1",
			Error: true,
		},
		{
			Name:  "Missing Storage Target Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.StorageCache/caches/cache1/storageTargets",
			Error: true,
		},
		{
			Name:  "With Storage Target Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.StorageCache/caches/cache1/storageTargets/target1",
			Expect: &StorageTargetId{
				ResourceGroup: "resGroup1",
				Cache:         "cache1",
				Name:          "target1",
			},
		},
		{
			Name:  "Wrong Casing",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.StorageCache/caches/cache1/StorageTargets/target1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := StorageTargetID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.ResourceGroup != v.Expect.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expect.ResourceGroup, actual.ResourceGroup)
		}

		if actual.Cache != v.Expect.Cache {
			t.Fatalf("Expected %q but got %q for Cache", v.Expect.Cache, actual.Cache)
		}

		if actual.Name != v.Expect.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.Name, actual.Name)
		}
	}
}
