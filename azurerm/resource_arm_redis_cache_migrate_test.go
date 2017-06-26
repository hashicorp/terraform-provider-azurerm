package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/terraform"
)

func TestAzureRMRedisCacheMigrateState(t *testing.T) {
	cases := map[string]struct {
		StateVersion int
		ID           string
		Attributes   map[string]string
		Expected     string
		Meta         interface{}
	}{
		"v0_1_without_value": {
			StateVersion: 0,
			ID:           "azurerm_redis_cache.test",
			Attributes:   map[string]string{},
			Expected:     "false",
		},
	}

	for tn, tc := range cases {
		is := &terraform.InstanceState{
			ID:         tc.ID,
			Attributes: tc.Attributes,
		}
		is, err := resourceAzureRMRedisCacheMigrateState(tc.StateVersion, is, tc.Meta)

		if err != nil {
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		actual := is.Attributes["redis_configuration.0.rdb_backup_enabled"]

		if actual != tc.Expected {
			t.Fatalf("bad Redis Cache Migrate: %s\n\n expected: %s", actual, tc.Expected)
		}
	}
}
