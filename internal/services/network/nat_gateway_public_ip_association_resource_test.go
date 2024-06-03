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
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/natgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NatGatewayPublicAssociationResource struct{}

func TestAccNatGatewayPublicIpAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway_public_ip_association", "test")
	r := NatGatewayPublicAssociationResource{}
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

func TestAccNatGatewayPublicIpAssociation_updateNatGateway(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway_public_ip_association", "test")
	r := NatGatewayPublicAssociationResource{}
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

func TestAccNatGatewayPublicIpAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway_public_ip_association", "test")
	r := NatGatewayPublicAssociationResource{}
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

func TestAccNatGatewayPublicIpAssociation_deleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway_public_ip_association", "test")
	r := NatGatewayPublicAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		// intentional as this is a Virtual Resource
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func (t NatGatewayPublicAssociationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseCompositeResourceID(state.ID, &natgateways.NatGatewayId{}, &commonids.PublicIPAddressId{})
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
			if props.PublicIPAddresses != nil {
				for _, pip := range *props.PublicIPAddresses {
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

func (NatGatewayPublicAssociationResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseCompositeResourceID(state.ID, &natgateways.NatGatewayId{}, &commonids.PublicIPAddressId{})
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

	updatedAddresses := make([]natgateways.SubResource, 0)
	if publicIpAddresses := resp.Model.Properties.PublicIPAddresses; publicIpAddresses != nil {
		for _, publicIpAddress := range *publicIpAddresses {
			if !strings.EqualFold(*publicIpAddress.Id, id.Second.ID()) {
				updatedAddresses = append(updatedAddresses, publicIpAddress)
			}
		}
	}
	resp.Model.Properties.PublicIPAddresses = &updatedAddresses

	if err := client.Network.Client.NatGateways.CreateOrUpdateThenPoll(ctx2, *id.First, *resp.Model); err != nil {
		return nil, fmt.Errorf("removing Association between %s and %s: %+v", id.First, id.Second, err)
	}

	return pointer.To(true), nil
}

func (r NatGatewayPublicAssociationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_nat_gateway" "test" {
  name                = "acctest-NatGateway-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard"
}

resource "azurerm_nat_gateway_public_ip_association" "test" {
  nat_gateway_id       = azurerm_nat_gateway.test.id
  public_ip_address_id = azurerm_public_ip.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r NatGatewayPublicAssociationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_nat_gateway_public_ip_association" "import" {
  nat_gateway_id       = azurerm_nat_gateway_public_ip_association.test.nat_gateway_id
  public_ip_address_id = azurerm_nat_gateway_public_ip_association.test.public_ip_address_id
}
`, r.basic(data))
}

func (r NatGatewayPublicAssociationResource) updateNatGateway(data acceptance.TestData) string {
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

resource "azurerm_nat_gateway_public_ip_association" "test" {
  nat_gateway_id       = azurerm_nat_gateway.test.id
  public_ip_address_id = azurerm_public_ip.test.id
}
`, r.template(data), data.RandomInteger)
}

func (NatGatewayPublicAssociationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ngpi-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-PIP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
