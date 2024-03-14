// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/subnets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SubnetNatGatewayAssociationResource struct{}

func TestAccAzureRMSubnetNatGatewayAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_nat_gateway_association", "test")
	r := SubnetNatGatewayAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		// intentional since this is a virtual resource
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMSubnetNatGatewayAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_nat_gateway_association", "test")
	r := SubnetNatGatewayAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		// intentional since this is a virtual resource
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError(data.ResourceType),
		},
	})
}

func TestAccAzureRMSubnetNatGatewayAssociation_deleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_nat_gateway_association", "test")
	r := SubnetNatGatewayAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		// intentionally not using a DisappearsStep since this is virtual resource
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClient(r.destroy),
				data.CheckWithClientForResource(SubnetResource{}.hasNoNatGateway, "azurerm_subnet.test"),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccAzureRMSubnetNatGatewayAssociation_updateSubnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_nat_gateway_association", "test")
	r := SubnetNatGatewayAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		// intentional since this is a virtual resource
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateSubnet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t SubnetNatGatewayAssociationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(5*time.Minute))
	defer cancel()

	id, err := commonids.ParseSubnetID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroupName
	virtualNetworkName := id.VirtualNetworkName
	subnetName := id.SubnetName

	resp, err := clients.Network.SubnetsClient.Get(ctx, *id, subnets.GetOperationOptions{})
	if err != nil {
		return nil, fmt.Errorf("reading Subnet Nat Gateway Association (%s): %+v", id, err)
	}

	props := resp.Model.Properties
	if props == nil || props.NatGateway == nil {
		return nil, fmt.Errorf("properties was nil for Subnet %q (Virtual Network %q / Resource Group: %q)", subnetName, virtualNetworkName, resourceGroup)
	}

	return utils.Bool(props.NatGateway.Id != nil), nil
}

func (SubnetNatGatewayAssociationResource) destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(5*time.Minute))
	defer cancel()

	parsedSubnetId, err := commonids.ParseSubnetID(state.Attributes["subnet_id"])
	if err != nil {
		return err
	}

	resourceGroup := parsedSubnetId.ResourceGroupName
	virtualNetworkName := parsedSubnetId.VirtualNetworkName
	subnetName := parsedSubnetId.SubnetName

	subnet, err := client.Network.SubnetsClient.Get(ctx, *parsedSubnetId, subnets.GetOperationOptions{})
	if err != nil {
		if subnet.HttpResponse.StatusCode != http.StatusNotFound {
			return fmt.Errorf("retrieving Subnet %q (Network %q / Resource Group %q): %+v", subnetName, virtualNetworkName, resourceGroup, err)
		}
		return fmt.Errorf("Bad: Get on subnetClient: %+v", err)
	}

	props := subnet.Model.Properties
	if props == nil {
		return fmt.Errorf("Properties was nil for Subnet %q (Virtual Network %q / Resource Group: %q)", subnetName, virtualNetworkName, resourceGroup)
	}
	props.NatGateway = nil

	err = client.Network.SubnetsClient.CreateOrUpdateThenPoll(ctx, *parsedSubnetId, *subnet.Model)
	if err != nil {
		return fmt.Errorf("updating Subnet %q (Network %q / Resource Group %q): %+v", subnetName, virtualNetworkName, resourceGroup, err)
	}

	return nil
}

func (r SubnetNatGatewayAssociationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_subnet_nat_gateway_association" "test" {
  subnet_id      = azurerm_subnet.test.id
  nat_gateway_id = azurerm_nat_gateway.test.id
}
`, r.template(data))
}

func (r SubnetNatGatewayAssociationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet_nat_gateway_association" "import" {
  subnet_id      = azurerm_subnet_nat_gateway_association.test.subnet_id
  nat_gateway_id = azurerm_subnet_nat_gateway_association.test.nat_gateway_id
}
`, r.basic(data))
}

func (r SubnetNatGatewayAssociationResource) updateSubnet(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  enforce_private_link_endpoint_network_policies = true
}

resource "azurerm_subnet_nat_gateway_association" "test" {
  subnet_id      = azurerm_subnet.test.id
  nat_gateway_id = azurerm_nat_gateway.test.id
}
`, r.template(data))
}

func (SubnetNatGatewayAssociationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_nat_gateway" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
