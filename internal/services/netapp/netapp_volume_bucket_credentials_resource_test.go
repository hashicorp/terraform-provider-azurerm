// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-01-01/buckets"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
)

type NetAppVolumeBucketCredentialsResource struct{}

func TestAccNetAppVolumeBucketCredentials_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_bucket_credentials", "test")
	r := NetAppVolumeBucketCredentialsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("access_key").Exists(),
				check.That(data.ResourceName).Key("secret_key").Exists(),
			),
		},
		data.ImportStep("access_key", "secret_key", "key_pair_expiry"),
	})
}

func (t NetAppVolumeBucketCredentialsResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := buckets.ParseBucketID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.NetApp.BucketsClient.Get(ctx, *id)
	if err != nil {
		if resp.HttpResponse != nil && resp.HttpResponse.StatusCode == http.StatusNotFound {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(true), nil
}

func (NetAppVolumeBucketCredentialsResource) basic(data acceptance.TestData) string {
	template := NetAppVolumeBucketResource{}.basic(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_bucket_credentials" "test" {
  bucket_id            = azurerm_netapp_volume_bucket.test.id
  key_pair_expiry_days = 30
}
`, template)
}
