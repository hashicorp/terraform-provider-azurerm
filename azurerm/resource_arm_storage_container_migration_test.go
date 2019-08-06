package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/terraform"
)

// NOTE: this is intentionally an acceptance test (and we're not explicitly setting the env)
// as we want to run this depending on the cloud we're in.
func TestAccAzureRMStorageContainerMigrateState(t *testing.T) {
	config := testGetAzureConfig(t)
	if config == nil {
		t.SkipNow()
		return
	}

	client, err := getArmClient(config, false, "", true)
	if err != nil {
		t.Fatal(fmt.Errorf("Error building ARM Client: %+v", err))
		return
	}

	client.StopContext = testAccProvider.StopContext()

	suffix := client.environment.StorageEndpointSuffix

	cases := map[string]struct {
		StateVersion       int
		ID                 string
		InputAttributes    map[string]string
		ExpectedAttributes map[string]string
	}{
		"v0_1_without_value": {
			StateVersion: 0,
			ID:           "some_id",
			InputAttributes: map[string]string{
				"name":                 "container",
				"storage_account_name": "example",
			},
			ExpectedAttributes: map[string]string{
				"id": fmt.Sprintf("https://example.blob.%s/container", suffix),
			},
		},
	}

	for tn, tc := range cases {
		is := &terraform.InstanceState{
			ID:         tc.ID,
			Attributes: tc.InputAttributes,
		}
		is, err := resourceStorageContainerMigrateState(tc.StateVersion, is, client)

		if err != nil {
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		for k, v := range tc.ExpectedAttributes {
			actual := is.Attributes[k]
			if actual != v {
				t.Fatalf("Bad Storage Container Migrate for %q: %q\n\n expected: %q", k, actual, v)
			}
		}
	}
}
