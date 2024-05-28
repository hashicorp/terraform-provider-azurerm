// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/publicipprefixes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/natgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NatGatewayPublicIpPrefixAssociationResource struct{}

func TestAccNatGatewayPublicIpPrefixAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway_public_ip_prefix_association", "test")
	r := NatGatewayPublicIpPrefixAssociationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		// intentional as this is a Virtual Resource
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNatGatewayPublicIpPrefixAssociation_updateNatGateway(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway_public_ip_prefix_association", "test")
	r := NatGatewayPublicIpPrefixAssociationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		// intentional as this is a Virtual Resource
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateNatGateway(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNatGatewayPublicIpPrefixAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway_public_ip_prefix_association", "test")
	r := NatGatewayPublicIpPrefixAssociationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		// intentional as this is a Virtual Resource
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccNatGatewayPublicIpPrefixAssociation_deleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway_public_ip_prefix_association", "test")
	r := NatGatewayPublicIpPrefixAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		// intentional as this is a Virtual Resource
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func (t NatGatewayPublicIpPrefixAssociationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseCompositeResourceID(state.ID, &natgateways.NatGatewayId{}, &publicipprefixes.PublicIPPrefixId{})
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.Client.NatGateways.Get(ctx, *id.First, natgateways.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id.First, err)
	}

	found := false
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if props.PublicIPPrefixes != nil {
				for _, pip := range *props.PublicIPPrefixes {
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

func (NatGatewayPublicIpPrefixAssociationResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseCompositeResourceID(state.ID, &natgateways.NatGatewayId{}, &publicipprefixes.PublicIPPrefixId{})
	if err != nil {
		return nil, err
	}

	ctx2, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()
	resp, err := client.Network.Client.NatGateways.Get(ctx2, *id.First, natgateways.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id.First, err)
	}

	if resp.Model == nil {
		return nil, fmt.Errorf("retrieving %s: `model` was nil", id.First)
	}
	if resp.Model.Properties == nil {
		return nil, fmt.Errorf("retrieving %s: `properties` was nil", id.First)
	}

	updatedPrefixes := make([]natgateways.SubResource, 0)
	if publicIpPrefixes := resp.Model.Properties.PublicIPPrefixes; publicIpPrefixes != nil {
		for _, publicIpPrefix := range *publicIpPrefixes {
			if !strings.EqualFold(*publicIpPrefix.Id, id.Second.ID()) {
				updatedPrefixes = append(updatedPrefixes, publicIpPrefix)
			}
		}
	}
	resp.Model.Properties.PublicIPPrefixes = &updatedPrefixes

	if err := client.Network.Client.NatGateways.CreateOrUpdateThenPoll(ctx2, *id.First, *resp.Model); err != nil {
		return nil, fmt.Errorf("deleting %s: %+v", id, err)
	}

	return pointer.To(true), nil
}

func (r NatGatewayPublicIpPrefixAssociationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_nat_gateway" "test" {
  name                = "acctest-NatGateway-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard"
}

resource "azurerm_nat_gateway_public_ip_prefix_association" "test" {
  nat_gateway_id      = azurerm_nat_gateway.test.id
  public_ip_prefix_id = azurerm_public_ip_prefix.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r NatGatewayPublicIpPrefixAssociationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_nat_gateway_public_ip_prefix_association" "import" {
  nat_gateway_id      = azurerm_nat_gateway_public_ip_prefix_association.test.nat_gateway_id
  public_ip_prefix_id = azurerm_nat_gateway_public_ip_prefix_association.test.public_ip_prefix_id
}
`, r.basic(data))
}

func (r NatGatewayPublicIpPrefixAssociationResource) updateNatGateway(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_nat_gateway" "test" {
  name                = "acctest-NatGateway-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard"
  tags = {
    Hello = "World"
  }
}

resource "azurerm_nat_gateway_public_ip_prefix_association" "test" {
  nat_gateway_id      = azurerm_nat_gateway.test.id
  public_ip_prefix_id = azurerm_public_ip_prefix.test.id
}
`, r.template(data), data.RandomInteger)
}

func (NatGatewayPublicIpPrefixAssociationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ngpi-%d"
  location = "%s"
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "acctestpublicIPPrefix-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  prefix_length       = 30
  zones               = ["1"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
