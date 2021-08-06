package migration

import (
	"context"
	"testing"
)

func TestDataFactoryMigrateState(t *testing.T) {
	cases := map[string]struct {
		StateVersion    int
		InputAttributes map[string]interface{}
		ExpectedNewID   string
	}{
		"name_upper_case": {
			StateVersion: 1,
			InputAttributes: map[string]interface{}{
				"id":                  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DataFactory/factories/acctest",
				"name":                "ACCTEST",
				"location":            "westeurope",
				"resource_group_name": "resGroup1",
			},
			ExpectedNewID: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DataFactory/factories/ACCTEST",
		},
	}

	for _, tc := range cases {
		newID, _ := DataFactoryV1ToV2{}.UpgradeFunc()(context.TODO(), tc.InputAttributes, nil)

		if newID["id"].(string) != tc.ExpectedNewID {
			t.Fatalf("ID migration failed, expected %q, got: %q", tc.ExpectedNewID, newID["id"].(string))
		}
	}
}
