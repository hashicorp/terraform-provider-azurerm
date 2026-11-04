// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package discovery_test

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

func TestAccDiscoveryBookshelf_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_bookshelf", "testlist")
	r := DiscoveryBookshelfResource{}
	listResourceAddress := "azurerm_discovery_bookshelf.list"

	expectedIdentity := map[string]knownvalue.Check{
		"name":                knownvalue.StringExact(fmt.Sprintf("acctest-%s", data.RandomString)),
		"resource_group_name": knownvalue.StringExact(fmt.Sprintf("acctestrg-%d", data.RandomInteger)),
		"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
	}

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
				Config: r.subscriptionListQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 1),
					querycheck.ExpectIdentity(listResourceAddress, expectedIdentity),
				},
			},
			{
				Query:  true,
				Config: r.resourceGroupListQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 1),
					querycheck.ExpectIdentity(listResourceAddress, expectedIdentity),
				},
			},
		},
	})
}

func (r DiscoveryBookshelfResource) subscriptionListQuery() string {
	return `
list "azurerm_discovery_bookshelf" "list" {
  provider = azurerm
  config {}
}
`
}

func (r DiscoveryBookshelfResource) resourceGroupListQuery() string {
	return `
list "azurerm_discovery_bookshelf" "list" {
  provider = azurerm
  config {
    resource_group_name = azurerm_resource_group.test.name
  }
}
`
}
