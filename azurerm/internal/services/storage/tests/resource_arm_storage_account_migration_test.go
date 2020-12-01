package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage"
)

func TestAzureRMStorageAccountMigrateState(t *testing.T) {
	cases := map[string]struct {
		StateVersion       int
		ID                 string
		InputAttributes    map[string]string
		ExpectedAttributes map[string]string
		Meta               interface{}
	}{
		"v0_1_with_standard": {
			StateVersion: 0,
			ID:           "some_id",
			InputAttributes: map[string]string{
				"account_type": "Standard_LRS",
			},
			ExpectedAttributes: map[string]string{
				"account_tier":             "Standard",
				"account_replication_type": "LRS",
			},
		},
		"v0_1_with_premium": {
			StateVersion: 0,
			ID:           "some_id",
			InputAttributes: map[string]string{
				"account_type": "Premium_GRS",
			},
			ExpectedAttributes: map[string]string{
				"account_tier":             "Premium",
				"account_replication_type": "GRS",
			},
		},
		"v1_2_empty": {
			StateVersion:    1,
			ID:              "some_id",
			InputAttributes: map[string]string{},
			ExpectedAttributes: map[string]string{
				"account_encryption_source": "Microsoft.Storage",
			},
		},
	}

	for tn, tc := range cases {
		is := &terraform.InstanceState{
			ID:         tc.ID,
			Attributes: tc.InputAttributes,
		}
		is, err := storage.ResourceStorageAccountMigrateState(tc.StateVersion, is, tc.Meta)
		if err != nil {
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		for k, v := range tc.ExpectedAttributes {
			actual := is.Attributes[k]
			if actual != v {
				t.Fatalf("Bad Storage Account Migrate for %q: %q\n\n expected: %q", k, actual, v)
			}
		}
	}
}
