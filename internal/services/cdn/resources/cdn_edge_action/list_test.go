// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_edge_action_test

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

func TestAccCdnEdgeAction_list_basic(t *testing.T) {
	r := CdnEdgeActionResource{}
	listResourceAddress := "azurerm_cdn_edge_action.list"

	data := acceptance.BuildTestData(t, "azurerm_cdn_edge_action", "test")

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				Query:  true,
				Config: r.basicQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 2),
				},
			},
			{
				Query:  true,
				Config: r.basicQueryByResourceGroupName(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 2),
				},
			},
		},
	})
}

func (r CdnEdgeActionResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_edge_action" "test" {
  count = 2

  name                = "acctestea${count.index}%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard"
}
`, r.template(data), data.RandomInteger)
}

func (CdnEdgeActionResource) basicQuery() string {
	return `
list "azurerm_cdn_edge_action" "list" {
  provider = azurerm
  config {}
}
`
}

func (CdnEdgeActionResource) basicQueryByResourceGroupName() string {
	return `
list "azurerm_cdn_edge_action" "list" {
  provider = azurerm
  config {
    resource_group_name = azurerm_resource_group.test.name
  }
}
`
}
