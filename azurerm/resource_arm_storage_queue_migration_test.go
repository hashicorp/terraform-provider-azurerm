package azurerm

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// NOTE: this is intentionally an acceptance test (and we're not explicitly setting the env)
// as we want to run this depending on the cloud we're in.
func TestAccAzureRMStorageQueueMigrateState(t *testing.T) {
	config := testGetAzureConfig(t)
	if config == nil {
		t.SkipNow()
		return
	}

	builder := armClientBuilder{
		authConfig:                  config,
		terraformVersion:            "0.0.0",
		partnerId:                   "",
		disableCorrelationRequestID: true,
		disableTerraformPartnerID:   false,
		skipProviderRegistration:    false,
	}
	client, err := getArmClient(context.Background(), builder)
	if err != nil {
		t.Fatal(fmt.Errorf("Error building ARM Client: %+v", err))
		return
	}

	client.StopContext = testAccProvider.StopContext()

	suffix := client.Account.Environment.StorageEndpointSuffix

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
				"name":                 "queue",
				"storage_account_name": "example",
			},
			ExpectedAttributes: map[string]string{
				"id": fmt.Sprintf("https://example.queue.%s/queue", suffix),
			},
		},
	}

	for tn, tc := range cases {
		is := &terraform.InstanceState{
			ID:         tc.ID,
			Attributes: tc.InputAttributes,
		}
		is, err := resourceStorageQueueMigrateState(tc.StateVersion, is, client)

		if err != nil {
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		for k, v := range tc.ExpectedAttributes {
			actual := is.Attributes[k]
			if actual != v {
				t.Fatalf("Bad Storage Queue Migrate for %q: %q\n\n expected: %q", k, actual, v)
			}
		}
	}
}
