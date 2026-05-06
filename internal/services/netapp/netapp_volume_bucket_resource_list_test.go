// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccNetAppVolumeBucket_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_bucket", "test")
	r := NetAppVolumeBucketResource{}
	bucketName := fmt.Sprintf("acctest-bucket-%d", data.RandomInteger)
	resourceGroupName := fmt.Sprintf("acctestRG-netapp-%d", data.RandomInteger)
	accountName := fmt.Sprintf("acctest-NetAppAccount-%d", data.RandomInteger)
	poolName := fmt.Sprintf("acctest-NetAppPool-%d", data.RandomInteger)
	volumeName := fmt.Sprintf("acctest-NetAppVolume-%d", data.RandomInteger)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{Config: r.basic(data)},
			{
				Query:  true,
				Config: r.listQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_netapp_volume_bucket.list", 1),
					querycheck.ExpectIdentity("azurerm_netapp_volume_bucket.list", map[string]knownvalue.Check{
						"bucket_name":          knownvalue.StringExact(bucketName),
						"volume_name":          knownvalue.StringExact(volumeName),
						"capacity_pool_name":   knownvalue.StringExact(poolName),
						"net_app_account_name": knownvalue.StringExact(accountName),
						"resource_group_name":  knownvalue.StringExact(resourceGroupName),
						"subscription_id":      knownvalue.StringExact(data.Subscriptions.Primary),
					}),
				},
			},
		},
	})
}

func (NetAppVolumeBucketResource) listQuery() string {
	return `
list "azurerm_netapp_volume_bucket" "list" {
  provider = azurerm
  config {
    volume_id = azurerm_netapp_volume.test.id
  }
}
`
}
