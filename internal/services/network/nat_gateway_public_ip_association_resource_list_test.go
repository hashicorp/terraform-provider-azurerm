// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
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

func TestAccNatGatewayPublicIpAssociation_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway_public_ip_association", "list")
	r := NatGatewayPublicIpAssociationResource{}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.multipleAssociations(data),
			},
			{
				Query:  true,
				Config: r.basicListQueryByNatGateway(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_nat_gateway_public_ip_association.list", 3),
					querycheck.ExpectIdentity(
						"azurerm_nat_gateway_public_ip_association.list",
						map[string]knownvalue.Check{
							"resource_id1": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_id2": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
						},
					),
				},
			},
			{
				Query:  true,
				Config: r.basicListQueryByResourceGroup(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_nat_gateway_public_ip_association.list", 3),
					querycheck.ExpectIdentity(
						"azurerm_nat_gateway_public_ip_association.list",
						map[string]knownvalue.Check{
							"resource_id1": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_id2": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
						},
					),
				},
			},
			{
				Query:  true,
				Config: r.basicListQueryBySubscription(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("azurerm_nat_gateway_public_ip_association.list", 3),
					querycheck.ExpectIdentity(
						"azurerm_nat_gateway_public_ip_association.list",
						map[string]knownvalue.Check{
							"resource_id1": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_id2": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
						},
					),
				},
			},
		},
	})
}

func (NatGatewayPublicIpAssociationResource) basicListQueryByNatGateway() string {
	return `
list "azurerm_nat_gateway_public_ip_association" "list" {
  provider = azurerm
  config {
    nat_gateway_id = azurerm_nat_gateway.test.id
  }
}
`
}

func (NatGatewayPublicIpAssociationResource) basicListQueryByResourceGroup() string {
	return `
list "azurerm_nat_gateway_public_ip_association" "list" {
  provider = azurerm
  config {
    resource_group_id = azurerm_resource_group.test.id
  }
}
`
}

func (NatGatewayPublicIpAssociationResource) basicListQueryBySubscription() string {
	return `
list "azurerm_nat_gateway_public_ip_association" "list" {
  provider = azurerm
  config {
  }
}
`
}
