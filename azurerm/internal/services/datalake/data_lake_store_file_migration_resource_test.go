package datalake_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

// NOTE: this is intentionally an acceptance test (and we're not explicitly setting the env)
// as we want to run this depending on the cloud we're in.
func TestAccAzureRMDataLakeStoreFileMigrateState(t *testing.T) {
	config := acceptance.GetAuthConfig(t)
	if config == nil {
		t.SkipNow()
		return
	}

	builder := clients.ClientBuilder{
		AuthConfig:                  config,
		TerraformVersion:            "0.0.0",
		PartnerId:                   "",
		DisableCorrelationRequestID: true,
		DisableTerraformPartnerID:   false,
		SkipProviderRegistration:    false,
	}
	client, err := clients.Build(context.Background(), builder)
	if err != nil {
		t.Fatal(fmt.Errorf("Error building ARM Client: %+v", err))
		return
	}

	client.StopContext = context.Background()

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
		is, err := ResourceDataLakeStoreFileMigrateState(tc.StateVersion, is, client)
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
