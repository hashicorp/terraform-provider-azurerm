package network_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-05-01/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetworkInterfaceApplicationSecurityGroupAssociationResource struct {
}

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
	splitId := strings.Split(state.ID, "|")
	if len(splitId) != 2 {
		return nil, fmt.Errorf("expected ID to be in the format {networkInterfaceId}/ipConfigurations/{ipConfigurationName}|{backendAddressPoolId} but got %q", state.ID)
	}

	id, err := parse.NetworkInterfaceID(splitId[0])
	if err != nil {
		return nil, err
	}
	applicationSecurityGroupId := splitId[1]

	read, err := clients.Network.InterfacesClient.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading (%s): %+v", *id, err)
	}

	found := false
	for _, config := range *read.InterfacePropertiesFormat.IPConfigurations {
		if config.ApplicationSecurityGroups != nil {
			for _, group := range *config.ApplicationSecurityGroups {
				if *group.ID == applicationSecurityGroupId {
					found = true
					break
				}
			}
		}
	}

	return utils.Bool(found), nil
}

func (NetworkInterfaceApplicationSecurityGroupAssociationResource) destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
	nicID, err := parse.NetworkInterfaceID(state.Attributes["network_interface_id"])
	if err != nil {
		return err
	}

	applicationSecurityGroupId := state.Attributes["application_security_group_id"]

	read, err := client.Network.InterfacesClient.Get(ctx, nicID.ResourceGroup, nicID.Name, "")
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *nicID, err)
	}

	configs := *read.InterfacePropertiesFormat.IPConfigurations
	for _, config := range configs {
		if config.ApplicationSecurityGroups != nil {
			groups := make([]network.ApplicationSecurityGroup, 0)
			for _, group := range *config.ApplicationSecurityGroups {
				if *group.ID != applicationSecurityGroupId {
					groups = append(groups, group)
				}
			}
			config.ApplicationSecurityGroups = &groups
		}
	}

	read.InterfacePropertiesFormat.IPConfigurations = &configs

	future, err := client.Network.InterfacesClient.CreateOrUpdate(ctx, nicID.ResourceGroup, nicID.Name, read)
	if err != nil {
		return fmt.Errorf("removing Application Security Group Association for %s: %+v", *nicID, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Network.InterfacesClient.Client); err != nil {
		return fmt.Errorf("waiting for removal of Application Security Group Association for %s: %+v", *nicID, err)
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
    private_ip_address_allocation = "dynamic"
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
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_application_security_group" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
