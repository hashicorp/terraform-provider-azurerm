// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package discovery_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccDiscoveryBookshelfPrivateEndpointConnection_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_bookshelf_private_endpoint_connection", "testlist")
	r := DiscoveryBookshelfPrivateEndpointConnectionResource{}
	listResourceAddress := "azurerm_discovery_bookshelf_private_endpoint_connection.list"

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
				Config: r.listQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 1),
					querycheck.ExpectIdentity(listResourceAddress, map[string]knownvalue.Check{
						// the service appends a generated GUID to the Private Endpoint name
						"name":                knownvalue.StringRegexp(regexp.MustCompile(fmt.Sprintf(`^acctestpe-%d\.`, data.RandomInteger))),
						"bookshelf_name":      knownvalue.StringExact(fmt.Sprintf("acctest-%s", data.RandomString)),
						"resource_group_name": knownvalue.StringExact(fmt.Sprintf("acctestrg-%d", data.RandomInteger)),
						"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
					}),
				},
			},
		},
	})
}

func (r DiscoveryBookshelfPrivateEndpointConnectionResource) listQuery() string {
	return `
list "azurerm_discovery_bookshelf_private_endpoint_connection" "list" {
  provider = azurerm
  config {
    bookshelf_id = azurerm_discovery_bookshelf.test.id
  }
}
`
}
