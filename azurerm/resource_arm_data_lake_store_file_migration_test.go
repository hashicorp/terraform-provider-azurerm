package azurerm

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// NOTE: this is intentionally an acceptance test (and we're not explicitly setting the env)
// as we want to run this depending on the cloud we're in.
func TestAccAzureRMDataLakeStoreFileMigrateState(t *testing.T) {
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

	filesClient := client.Datalake.StoreFilesClient

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
				"remote_file_path": "/test/blob.vhd",
				"account_name":     "example",
			},
			ExpectedAttributes: map[string]string{
				"id": fmt.Sprintf("example.%s/test/blob.vhd", filesClient.AdlsFileSystemDNSSuffix),
			},
		},
	}

	for tn, tc := range cases {
		is := &terraform.InstanceState{
			ID:         tc.ID,
			Attributes: tc.InputAttributes,
		}
		is, err := resourceDataLakeStoreFileMigrateState(tc.StateVersion, is, client)

		if err != nil {
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		for k, v := range tc.ExpectedAttributes {
			actual := is.Attributes[k]
			if actual != v {
				t.Fatalf("Bad Data Lake Store File Migrate for %q: %q\n\n expected: %q", k, actual, v)
			}
		}
	}
}
