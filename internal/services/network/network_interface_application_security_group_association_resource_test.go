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
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/applicationsecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetworkInterfaceApplicationSecurityGroupAssociationResource struct{}

func TestAccNetworkInterfaceApplicationSecurityGroupAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface_application_security_group_association", "test")
	r := NetworkInterfaceApplicationSecurityGroupAssociationResource{}
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

func TestAccNetworkInterfaceApplicationSecurityGroupAssociation_multipleIPConfigurations(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface_application_security_group_association", "test")
	r := NetworkInterfaceApplicationSecurityGroupAssociationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleIPConfigurations(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkInterfaceApplicationSecurityGroupAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface_application_security_group_association", "test")
	r := NetworkInterfaceApplicationSecurityGroupAssociationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_network_interface_application_security_group_association"),
		},
	})
}

func TestAccNetworkInterfaceApplicationSecurityGroupAssociation_deleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface_application_security_group_association", "test")
	r := NetworkInterfaceApplicationSecurityGroupAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
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

func TestAccNetworkInterfaceApplicationSecurityGroupAssociation_updateNIC(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface_application_security_group_association", "test")
	r := NetworkInterfaceApplicationSecurityGroupAssociationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
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

func (t NetworkInterfaceApplicationSecurityGroupAssociationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseCompositeResourceID(state.ID, &commonids.NetworkInterfaceId{}, &applicationsecuritygroups.ApplicationSecurityGroupId{})
	if err != nil {
		return nil, err
	}

	read, err := clients.Network.Client.NetworkInterfaces.Get(ctx, *id.First, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id.First, err)
	}

	found := false
	if model := read.Model; model != nil {
		if props := model.Properties; props != nil {
			if props.IPConfigurations != nil {
				for _, config := range *props.IPConfigurations {
					if ipConfigProps := config.Properties; ipConfigProps != nil {
						if ipConfigProps.ApplicationSecurityGroups != nil {
							for _, group := range *ipConfigProps.ApplicationSecurityGroups {
								if *group.Id == id.Second.ID() {
									found = true
									break
								}
							}
						}
					}
				}
			}
		}
	}

	return pointer.To(found), nil
}

func (NetworkInterfaceApplicationSecurityGroupAssociationResource) destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
	id, err := commonids.ParseCompositeResourceID(state.ID, &commonids.NetworkInterfaceId{}, &applicationsecuritygroups.ApplicationSecurityGroupId{})
	if err != nil {
		return err
	}

	ctx2, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()
	read, err := client.Network.Client.NetworkInterfaces.Get(ctx2, *id.First, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id.First, err)
	}
	if read.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id.First)
	}
	if read.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id.First)
	}
	if read.Model.Properties.IPConfigurations == nil {
		return fmt.Errorf("retrieving %s: `properties.ipConfigurations` was nil", id.First)
	}

	configs := *read.Model.Properties.IPConfigurations
	for _, config := range configs {
		if props := config.Properties; props != nil {
			groups := make([]networkinterfaces.ApplicationSecurityGroup, 0)
			for _, group := range *props.ApplicationSecurityGroups {
				if *group.Id != id.Second.ID() {
					groups = append(groups, group)
				}
			}
			props.ApplicationSecurityGroups = &groups
		}
	}

	read.Model.Properties.IPConfigurations = &configs

	if err := client.Network.Client.NetworkInterfaces.CreateOrUpdateThenPoll(ctx2, *id.First, *read.Model); err != nil {
		return fmt.Errorf("removing Application Security Group Association for %s: %+v", id.First, err)
	}

	return nil
}

func (r NetworkInterfaceApplicationSecurityGroupAssociationResource) basic(data acceptance.TestData) string {
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

resource "azurerm_network_interface_application_security_group_association" "test" {
  network_interface_id          = azurerm_network_interface.test.id
  application_security_group_id = azurerm_application_security_group.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkInterfaceApplicationSecurityGroupAssociationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface_application_security_group_association" "import" {
  network_interface_id          = azurerm_network_interface_application_security_group_association.test.network_interface_id
  application_security_group_id = azurerm_network_interface_application_security_group_association.test.application_security_group_id
}
`, r.basic(data))
}

func (r NetworkInterfaceApplicationSecurityGroupAssociationResource) multipleIPConfigurations(data acceptance.TestData) string {
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
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_network_interface_application_security_group_association" "test" {
  network_interface_id          = azurerm_network_interface.test.id
  application_security_group_id = azurerm_application_security_group.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkInterfaceApplicationSecurityGroupAssociationResource) updateNIC(data acceptance.TestData) string {
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

resource "azurerm_network_interface_application_security_group_association" "test" {
  network_interface_id          = azurerm_network_interface.test.id
  application_security_group_id = azurerm_application_security_group.test.id
}
`, r.template(data), data.RandomInteger)
}

func (NetworkInterfaceApplicationSecurityGroupAssociationResource) template(data acceptance.TestData) string {
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

resource "azurerm_application_security_group" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
