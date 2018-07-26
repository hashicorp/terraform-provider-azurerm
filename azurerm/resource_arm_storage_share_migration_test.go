package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/terraform"
)

func TestAzureRMStorageShareMigrateState(t *testing.T) {
	cases := map[string]struct {
		StateVersion       int
		ID                 string
		InputAttributes    map[string]string
		ExpectedAttributes map[string]string
		Meta               interface{}
	}{
		"v0_1": {
			StateVersion: 0,
			ID:           "some_id",
			InputAttributes: map[string]string{
				"name":                 "some_id",
				"resource_group_name":  "some_rgn",
				"storage_account_name": "some_sgn",
			},
			ExpectedAttributes: map[string]string{
				"id": "some_id/some_rgn/some_sgn",
			},
		},
	}

	for tn, tc := range cases {
		is := &terraform.InstanceState{
			ID:         tc.ID,
			Attributes: tc.InputAttributes,
		}
		is, err := resourceStorageShareMigrateState(tc.StateVersion, is, tc.Meta)

		if err != nil {
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		for k, v := range tc.ExpectedAttributes {
			actual := is.Attributes[k]
			if actual != v {
				t.Fatalf("Bad Storage Share Migrate for %q: %q\n\n expected: %q", k, actual, v)
			}
		}
	}
}
