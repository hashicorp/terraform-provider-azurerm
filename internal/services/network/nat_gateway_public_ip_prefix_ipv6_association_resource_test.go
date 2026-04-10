// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/natgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/publicipprefixes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NatGatewayPublicIpPrefixIPv6AssociationResource struct{}

func TestAccNatGatewayPublicIpPrefixIPv6Association_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway_public_ip_prefix_ipv6_association", "test")
	r := NatGatewayPublicIpPrefixIPv6AssociationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNatGatewayPublicIpPrefixIPv6Association_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway_public_ip_prefix_ipv6_association", "test")
	r := NatGatewayPublicIpPrefixIPv6AssociationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccNatGatewayPublicIpPrefixIPv6Association_multipleAssociations(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway_public_ip_prefix_ipv6_association", "test")
	r := NatGatewayPublicIpPrefixIPv6AssociationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleAssociations(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNatGatewayPublicIpPrefixIPv6Association_publicIPPrefixMustBeIPv6(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway_public_ip_prefix_ipv6_association", "test")
	r := NatGatewayPublicIpPrefixIPv6AssociationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.prerequisites(data, "StandardV2", "StandardV2", "IPv4"),
		},
		{
			Config:      r.publicIPPrefixMustBeIPv6(data),
			ExpectError: regexp.MustCompile("`public_ip_prefix_id` must reference an `IPv6` Public IP Prefix with SKU `StandardV2`"),
		},
	})
}

func TestAccNatGatewayPublicIpPrefixIPv6Association_standardSkuNatGatewayCannotUseIPv6PublicIPAddressesOrPrefixes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway_public_ip_prefix_ipv6_association", "test")
	r := NatGatewayPublicIpPrefixIPv6AssociationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.prerequisites(data, "Standard", "StandardV2", "IPv6"),
		},
		{
			Config:      r.standardSkuNatGatewayCannotUseIPv6PublicIPAddressesOrPrefixes(data),
			ExpectError: regexp.MustCompile("`nat_gateway_id` must reference a NAT Gateway with SKU `StandardV2`"),
		},
	})
}

func TestAccNatGatewayPublicIpPrefixIPv6Association_standardV2SkuNatGatewayRequiresPublicIPPrefixWithStandardV2Sku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway_public_ip_prefix_ipv6_association", "test")
	r := NatGatewayPublicIpPrefixIPv6AssociationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.prerequisites(data, "StandardV2", "Standard", "IPv6"),
		},
		{
			Config:      r.standardV2SkuNatGatewayRequiresPublicIPPrefixWithStandardV2Sku(data),
			ExpectError: regexp.MustCompile("`public_ip_prefix_id` must reference an `IPv6` Public IP Prefix with SKU `StandardV2`"),
		},
	})
}

func (t NatGatewayPublicIpPrefixIPv6AssociationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseCompositeResourceID(state.ID, &natgateways.NatGatewayId{}, &publicipprefixes.PublicIPPrefixId{})
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.NatGateways.Get(ctx, *id.First, natgateways.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id.First, err)
	}

	found := false
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if props.PublicIPPrefixesV6 != nil {
				for _, pip := range *props.PublicIPPrefixesV6 {
					if pip.Id == nil {
						continue
					}

					if strings.EqualFold(*pip.Id, id.Second.ID()) {
						found = true
						break
					}
				}
			}
		}
	}

	return pointer.To(found), nil
}

func (r NatGatewayPublicIpPrefixIPv6AssociationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_nat_gateway" "test" {
  name                = "acctest-NatGateway-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "StandardV2"
}

resource "azurerm_nat_gateway_public_ip_prefix_ipv6_association" "test" {
  nat_gateway_id      = azurerm_nat_gateway.test.id
  public_ip_prefix_id = azurerm_public_ip_prefix.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r NatGatewayPublicIpPrefixIPv6AssociationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_nat_gateway_public_ip_prefix_ipv6_association" "import" {
  nat_gateway_id      = azurerm_nat_gateway_public_ip_prefix_ipv6_association.test.nat_gateway_id
  public_ip_prefix_id = azurerm_nat_gateway_public_ip_prefix_ipv6_association.test.public_ip_prefix_id
}
`, r.basic(data))
}

func (r NatGatewayPublicIpPrefixIPv6AssociationResource) multipleAssociations(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_nat_gateway" "test" {
  name                = "acctest-NatGateway-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "StandardV2"
}

resource "azurerm_nat_gateway_public_ip_prefix_ipv6_association" "test" {
  nat_gateway_id      = azurerm_nat_gateway.test.id
  public_ip_prefix_id = azurerm_public_ip_prefix.test.id
}

resource "azurerm_public_ip_prefix" "test2" {
  name                = "acctest-pipPrefix2V6-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_version          = "IPv6"
  prefix_length       = 127
  sku                 = "StandardV2"
}

resource "azurerm_nat_gateway_public_ip_prefix_ipv6_association" "test2" {
  nat_gateway_id      = azurerm_nat_gateway.test.id
  public_ip_prefix_id = azurerm_public_ip_prefix.test2.id
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r NatGatewayPublicIpPrefixIPv6AssociationResource) publicIPPrefixMustBeIPv6(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_nat_gateway_public_ip_prefix_ipv6_association" "test" {
  nat_gateway_id      = azurerm_nat_gateway.test.id
  public_ip_prefix_id = azurerm_public_ip_prefix.test.id
}
`, r.prerequisites(data, "StandardV2", "StandardV2", "IPv4"))
}

func (r NatGatewayPublicIpPrefixIPv6AssociationResource) standardSkuNatGatewayCannotUseIPv6PublicIPAddressesOrPrefixes(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_nat_gateway_public_ip_prefix_ipv6_association" "test" {
  nat_gateway_id      = azurerm_nat_gateway.test.id
  public_ip_prefix_id = azurerm_public_ip_prefix.test.id
}
`, r.prerequisites(data, "Standard", "StandardV2", "IPv6"))
}

func (r NatGatewayPublicIpPrefixIPv6AssociationResource) standardV2SkuNatGatewayRequiresPublicIPPrefixWithStandardV2Sku(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_nat_gateway_public_ip_prefix_ipv6_association" "test" {
  nat_gateway_id      = azurerm_nat_gateway.test.id
  public_ip_prefix_id = azurerm_public_ip_prefix.test.id
}
`, r.prerequisites(data, "StandardV2", "Standard", "IPv6"))
}

func (NatGatewayPublicIpPrefixIPv6AssociationResource) prerequisites(data acceptance.TestData, natGatewaySkuName, publicIPPrefixSku, publicIPPrefixVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ngpipv6-%d"
  location = "%s"
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "acctest-pipPrefixV6-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_version          = "%s"
  prefix_length       = 127
  sku                 = "%s"
}

resource "azurerm_nat_gateway" "test" {
  name                = "acctest-NatGateway-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, publicIPPrefixVersion, publicIPPrefixSku, data.RandomInteger, natGatewaySkuName)
}

func (NatGatewayPublicIpPrefixIPv6AssociationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ngpipv6-%d"
  location = "%s"
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "acctest-pipPrefixV6-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_version          = "IPv6"
  prefix_length       = 127
  sku                 = "StandardV2"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
