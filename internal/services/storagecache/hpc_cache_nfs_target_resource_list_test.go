// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storagecache_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func testAccHPCCacheNFSTarget_list_basic(t *testing.T) {
	r := HpcCacheNfsTargetResource{}
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_nfs_target", "test")

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
				Config: r.basicListQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_hpc_cache_nfs_target.list", 1),
				},
			},
		},
	})
}

func (r HpcCacheNfsTargetResource) basicListQuery() string {
	return `
provider "azurerm" {
  features {}
}

list "azurerm_hpc_cache_nfs_target" "list" {
  provider = azurerm
  config {
    cache_id = azurerm_hpc_cache.test.id
  }
}
`
}
