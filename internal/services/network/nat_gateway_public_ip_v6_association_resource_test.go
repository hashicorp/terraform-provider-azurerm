// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/natgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NatGatewayPublicIpV6AssociationResource struct{}

func TestAccNatGatewayPublicIpV6Association_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway_public_ip_v6_association", "test")
	r := NatGatewayPublicIpV6AssociationResource{}
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

func TestAccNatGatewayPublicIpV6Association_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway_public_ip_v6_association", "test")
	r := NatGatewayPublicIpV6AssociationResource{}
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

func TestAccNatGatewayPublicIpV6Association_multipleAssociations(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway_public_ip_v6_association", "test")
	r := NatGatewayPublicIpV6AssociationResource{}
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

func (t NatGatewayPublicIpV6AssociationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseCompositeResourceID(state.ID, &natgateways.NatGatewayId{}, &commonids.PublicIPAddressId{})
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
			if props.PublicIPAddressesV6 != nil {
				for _, pip := range *props.PublicIPAddressesV6 {
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

func (r NatGatewayPublicIpV6AssociationResource) basic(data acceptance.TestData) string {
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
  zones               = ["1", "2", "3"]
}

resource "azurerm_nat_gateway_public_ip_v6_association" "test" {
  nat_gateway_id       = azurerm_nat_gateway.test.id
  public_ip_address_id = azurerm_public_ip.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r NatGatewayPublicIpV6AssociationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_nat_gateway_public_ip_v6_association" "import" {
  nat_gateway_id       = azurerm_nat_gateway_public_ip_v6_association.test.nat_gateway_id
  public_ip_address_id = azurerm_nat_gateway_public_ip_v6_association.test.public_ip_address_id
}
`, r.basic(data))
}

func (r NatGatewayPublicIpV6AssociationResource) multipleAssociations(data acceptance.TestData) string {
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
  zones               = ["1", "2", "3"]
}

resource "azurerm_nat_gateway_public_ip_v6_association" "test" {
  nat_gateway_id       = azurerm_nat_gateway.test.id
  public_ip_address_id = azurerm_public_ip.test.id
}

resource "azurerm_public_ip" "test2" {
  name                = "acctest-PIP2v6-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "StandardV2"
  ip_version          = "IPv6"
}

resource "azurerm_nat_gateway_public_ip_v6_association" "test2" {
  nat_gateway_id       = azurerm_nat_gateway.test.id
  public_ip_address_id = azurerm_public_ip.test2.id
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (NatGatewayPublicIpV6AssociationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ngpi-v6-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-PIPv6-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "StandardV2"
  ip_version          = "IPv6"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
