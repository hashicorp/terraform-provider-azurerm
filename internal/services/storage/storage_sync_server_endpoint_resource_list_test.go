// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccStorageSyncServerEndpoint_list_basic(t *testing.T) {
	t.Skip("@mbfrahry: temporarily skipping as the server must be registered manually. Will come back to this when the server can be registered programmatically")

	r := StorageSyncServerEndpointResource{}
	listResourceAddress := "azurerm_storage_sync_server_endpoint.list"
	data := acceptance.BuildTestData(t, "azurerm_storage_sync_server_endpoint", "test")

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
			},
			{
				Query:  true,
				Config: r.basicListQuery(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 1),
				},
			},
		},
	})
}

func (r StorageSyncServerEndpointResource) basicListQuery(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_storage_sync_server_endpoint" "list" {
  provider = azurerm
  config {
    storage_sync_group_id = "/subscriptions/%s/resourceGroups/acctestRG-StorageSync-%d/providers/Microsoft.StorageSync/storageSyncServices/acctest-StorageSync-%d/syncGroups/acctest-StorageSyncGroup-%d"
  }
}
`, data.Subscriptions.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
