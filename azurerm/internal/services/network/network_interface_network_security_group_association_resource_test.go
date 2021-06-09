package network_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/azuresdkhacks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type NetworkInterfaceNetworkSecurityGroupAssociationResource struct {
}

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
	splitId := strings.Split(state.ID, "|")
	if len(splitId) != 2 {
		return nil, fmt.Errorf("expected ID to be in the format {networkInterfaceId}|{networkSecurityGroupId} but got %q", state.ID)
	}

	nicID, err := azure.ParseAzureResourceID(splitId[0])
	if err != nil {
		return nil, err
	}

	name := nicID.Path["networkInterfaces"]
	resourceGroup := nicID.ResourceGroup

	read, err := clients.Network.InterfacesClient.Get(ctx, resourceGroup, name, "")
	if err != nil {
		return nil, fmt.Errorf("retrieving Network Interface %q: %+v", nicID, err)
	}

	found := false
	if read.InterfacePropertiesFormat != nil {
		if read.InterfacePropertiesFormat.NetworkSecurityGroup != nil && read.InterfacePropertiesFormat.NetworkSecurityGroup.ID != nil {
			found = *read.InterfacePropertiesFormat.NetworkSecurityGroup.ID == state.Attributes["network_security_group_id"]
		}
	}

	return utils.Bool(found), nil
}

func (NetworkInterfaceNetworkSecurityGroupAssociationResource) destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
	nicID, err := parse.NetworkInterfaceID(state.Attributes["network_interface_id"])
	if err != nil {
		return err
	}

	resourceGroup := nicID.ResourceGroup
	nicName := nicID.Name

	read, err := client.Network.InterfacesClient.Get(ctx, resourceGroup, nicName, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Network Interface %q (Resource Group %q): %+v", nicName, resourceGroup, err)
	}

	read.InterfacePropertiesFormat.NetworkSecurityGroup = nil

	future, err := azuresdkhacks.UpdateNetworkInterfaceAllowingRemovalOfNSG(ctx, client.Network.InterfacesClient, resourceGroup, nicName, read)
	if err != nil {
		return fmt.Errorf("Error removing Network Security Group Association for Network Interface %q (Resource Group %q): %+v", nicName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Network.InterfacesClient.Client); err != nil {
		return fmt.Errorf("Error waiting for removal of Network Security Group Association for NIC %q (Resource Group %q): %+v", nicName, resourceGroup, err)
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
    private_ip_address_allocation = "dynamic"
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
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
