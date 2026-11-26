// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccSubnetRouteTableAssociation_listByRouteTableID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_route_table_association", "list")
	r := SubnetRouteTableAssociationResource{}

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
				Config: r.basicQuery(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_subnet_route_table_association.list", 3),
					querycheck.ExpectIdentity(
						"azurerm_subnet_route_table_association.list",
						map[string]knownvalue.Check{
							"name":                 knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name":  knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"virtual_network_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":      knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
		},
	})
}

func (r SubnetRouteTableAssociationResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  count                = 3
  name                 = "acctest%d${count.index}"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.${count.index}.0/24"]
}


resource "azurerm_subnet_route_table_association" "test" {
  count          = 3
  subnet_id      = azurerm_subnet.test[count.index].id
  route_table_id = azurerm_route_table.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r SubnetRouteTableAssociationResource) basicQuery(data acceptance.TestData) string {
	return `
list "azurerm_subnet_route_table_association" "list" {
  provider = azurerm
  config {
    virtual_network_id = azurerm_virtual_network.test.id
  }
}
`
}
