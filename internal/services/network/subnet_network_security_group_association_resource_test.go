// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/subnets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SubnetNetworkSecurityGroupAssociationResource struct{}

func TestAccSubnetNetworkSecurityGroupAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_network_security_group_association", "test")
	r := SubnetNetworkSecurityGroupAssociationResource{}

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

func TestAccSubnetNetworkSecurityGroupAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_network_security_group_association", "test")
	r := SubnetNetworkSecurityGroupAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		// intentional as this is a Virtual Resource
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportAssociationErrorStep(r.requiresImport),
	})
}

func TestAccSubnetNetworkSecurityGroupAssociation_updateSubnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_network_security_group_association", "test")
	r := SubnetNetworkSecurityGroupAssociationResource{}

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
			Config: r.updateSubnet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSubnetNetworkSecurityGroupAssociation_deleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_network_security_group_association", "test")
	r := SubnetNetworkSecurityGroupAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		// intentional not using a DisappearsStep as this is a Virtual Resource
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClient(r.destroy),
				data.CheckWithClientForResource(SubnetResource{}.hasNoNetworkSecurityGroup, "azurerm_subnet.test"),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func (SubnetNetworkSecurityGroupAssociationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseSubnetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.Client.Subnets.Get(ctx, *id, subnets.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	model := resp.Model
	if model == nil {
		return nil, fmt.Errorf("model was nil for %s", *id)
	}

	props := model.Properties
	if props == nil || props.NetworkSecurityGroup == nil {
		return nil, fmt.Errorf("properties was nil for %s", *id)
	}

	return utils.Bool(props.NetworkSecurityGroup.Id != nil), nil
}

func (SubnetNetworkSecurityGroupAssociationResource) destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(15*time.Minute))
	defer cancel()

	subnetId := state.Attributes["subnet_id"]
	id, err := commonids.ParseSubnetID(subnetId)
	if err != nil {
		return err
	}

	read, err := client.Network.Client.Subnets.Get(ctx, *id, subnets.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(read.HttpResponse) {
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}
	}

	read.Model.Properties.NetworkSecurityGroup = nil

	if err := client.Network.Client.Subnets.CreateOrUpdateThenPoll(ctx, *id, *read.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return nil
}

func (r SubnetNetworkSecurityGroupAssociationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}
`, r.template(data))
}

func (r SubnetNetworkSecurityGroupAssociationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet_network_security_group_association" "internal" {
  subnet_id                 = azurerm_subnet_network_security_group_association.test.subnet_id
  network_security_group_id = azurerm_subnet_network_security_group_association.test.network_security_group_id
}
`, r.basic(data))
}

func (r SubnetNetworkSecurityGroupAssociationResource) updateSubnet(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  private_endpoint_network_policies = "Disabled"
}

resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}
`, r.template(data))
}

func (SubnetNetworkSecurityGroupAssociationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  security_rule {
    name                       = "test123"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "*"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
