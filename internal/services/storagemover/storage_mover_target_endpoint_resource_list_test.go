// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storagemover_test

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

func TestAccStorageMoverTargetEndpoint_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_target_endpoint", "testlist")
	r := StorageMoverTargetEndpointResource{}
	resourceName := fmt.Sprintf("acctest-smte-%d", data.RandomInteger)
	storageMoverName := fmt.Sprintf("acctest-ssm-%d", data.RandomInteger)
	resourceGroupName := fmt.Sprintf("acctest-rg-%d", data.RandomInteger)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{Config: r.listConfig(data)},
			{
				Query:  true,
				Config: r.listQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_storage_mover_target_endpoint.list", 1),
					querycheck.ExpectIdentity("azurerm_storage_mover_target_endpoint.list", map[string]knownvalue.Check{
						"name":                knownvalue.StringExact(resourceName),
						"resource_group_name": knownvalue.StringExact(resourceGroupName),
						"storage_mover_name":  knownvalue.StringExact(storageMoverName),
						"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
					}),
				},
			},
		},
	})
}

func (r StorageMoverTargetEndpointResource) listConfig(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_storage_mover_source_endpoint" "other" {
  name             = "acctest-smse-%d"
  storage_mover_id = azurerm_storage_mover.test.id
  host             = "192.168.0.1"
}

resource "azurerm_storage_mover_target_endpoint" "test" {
  name                   = "acctest-smte-%d"
  storage_mover_id       = azurerm_storage_mover.test.id
  storage_account_id     = azurerm_storage_account.test.id
  storage_container_name = azurerm_storage_container.test.name
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r StorageMoverTargetEndpointResource) listQuery() string {
	return `
list "azurerm_storage_mover_target_endpoint" "list" {
  provider = azurerm
  config {
    storage_mover_id = azurerm_storage_mover.test.id
  }
}
`
}
