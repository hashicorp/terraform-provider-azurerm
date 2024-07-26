// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/networksecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetworkInterfaceNetworkSecurityGroupAssociationResource struct{}

func TestAccNetworkInterfaceSecurityGroupAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface_security_group_association", "test")
	r := NetworkInterfaceNetworkSecurityGroupAssociationResource{}
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

func TestAccNetworkInterfaceSecurityGroupAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface_security_group_association", "test")
	r := NetworkInterfaceNetworkSecurityGroupAssociationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		// intentional as this is a Virtual Resource
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_network_interface_security_group_association"),
		},
	})
}

func TestAccNetworkInterfaceSecurityGroupAssociation_deleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface_security_group_association", "test")
	r := NetworkInterfaceNetworkSecurityGroupAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		// intentionally not using a DisappearsStep since this is a Virtual Resource
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClient(r.destroy),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccNetworkInterfaceSecurityGroupAssociation_updateNIC(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface_security_group_association", "test")
	r := NetworkInterfaceNetworkSecurityGroupAssociationResource{}
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
			Config: r.updateNIC(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t NetworkInterfaceNetworkSecurityGroupAssociationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseCompositeResourceID(state.ID, &commonids.NetworkInterfaceId{}, &networksecuritygroups.NetworkSecurityGroupId{})
	if err != nil {
		return nil, err
	}

	read, err := clients.Network.NetworkInterfaces.Get(ctx, *id.First, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id.First, err)
	}

	found := false
	if model := read.Model; model != nil {
		if props := model.Properties; props != nil {
			if props.NetworkSecurityGroup != nil && props.NetworkSecurityGroup.Id != nil {
				found = *props.NetworkSecurityGroup.Id == state.Attributes["network_security_group_id"]
			}
		}
	}

	return pointer.To(found), nil
}

func (NetworkInterfaceNetworkSecurityGroupAssociationResource) destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
	id, err := commonids.ParseCompositeResourceID(state.ID, &commonids.NetworkInterfaceId{}, &networksecuritygroups.NetworkSecurityGroupId{})
	if err != nil {
		return err
	}

	ctx2, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()
	read, err := client.Network.NetworkInterfaces.Get(ctx2, *id.First, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id.First, err)
	}

	read.Model.Properties.NetworkSecurityGroup = nil

	if err := client.Network.NetworkInterfaces.CreateOrUpdateThenPoll(ctx2, *id.First, *read.Model); err != nil {
		return fmt.Errorf("removing Network Security Group Association for %s: %+v", id.First, err)

	}

	return nil
}

func (r NetworkInterfaceNetworkSecurityGroupAssociationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_network_interface_security_group_association" "test" {
  network_interface_id      = azurerm_network_interface.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkInterfaceNetworkSecurityGroupAssociationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface_security_group_association" "import" {
  network_interface_id      = azurerm_network_interface_security_group_association.test.network_interface_id
  network_security_group_id = azurerm_network_interface_security_group_association.test.network_security_group_id
}
`, r.basic(data))
}

func (r NetworkInterfaceNetworkSecurityGroupAssociationResource) updateNIC(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    primary                       = true
  }

  ip_configuration {
    name                          = "testconfiguration2"
    private_ip_address_version    = "IPv6"
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_network_interface_security_group_association" "test" {
  network_interface_id      = azurerm_network_interface.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}
`, r.template(data), data.RandomInteger)
}

func (NetworkInterfaceNetworkSecurityGroupAssociationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
