// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_test

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
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle"
)

func TestAccNetworkAnchorResource_list(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.NetworkAnchorResource{}.ResourceType(), "test")
	r := NetworkAnchorResource{}
	listResourceAddress := "azurerm_oracle_network_anchor.list"
	identity := map[string]knownvalue.Check{
		"name":                knownvalue.StringExact("acctestNA" + data.RandomString),
		"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(`^acctestRG[0-9]+$`)),
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
				Config: r.basicListQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 1),
					querycheck.ExpectIdentity(listResourceAddress, identity),
				},
			},
			{
				Query:  true,
				Config: r.basicListQueryByResourceGroupName(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 1),
					querycheck.ExpectIdentity(listResourceAddress, identity),
				},
			},
			{
				Query:  true,
				Config: r.basicListQueryBySubscription(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 1),
					querycheck.ExpectIdentity(listResourceAddress, identity),
				},
			},
		},
	})
}

func (NetworkAnchorResource) basicListQuery() string {
	return `
list "azurerm_oracle_network_anchor" "list" {
  provider = azurerm
  config {}
}
`
}

func (NetworkAnchorResource) basicListQueryByResourceGroupName() string {
	return `
list "azurerm_oracle_network_anchor" "list" {
  provider = azurerm
  config {
    resource_group_name = azurerm_resource_group.test.name
  }
}
`
}

func (NetworkAnchorResource) basicListQueryBySubscription(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_oracle_network_anchor" "list" {
  provider = azurerm
  config {
    subscription_id = "%s"
  }
}
`, data.Subscriptions.Primary)
}
