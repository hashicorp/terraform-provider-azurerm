package network_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	network2 "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type NetworkInterfaceBackendAddressPoolResource struct {
}

func TestAccNetworkInterfaceBackendAddressPoolAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface_backend_address_pool_association", "test")
	r := NetworkInterfaceBackendAddressPoolResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		// intentional as this is a Virtual Resource
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkInterfaceBackendAddressPoolAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface_backend_address_pool_association", "test")
	r := NetworkInterfaceBackendAddressPoolResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		// intentional as this is a Virtual Resource
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_network_interface_backend_address_pool_association"),
		},
	})
}

func TestAccNetworkInterfaceBackendAddressPoolAssociation_deleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface_backend_address_pool_association", "test")
	r := NetworkInterfaceBackendAddressPoolResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		// intentionally not using a DisppearsStep as this is a Virtual Resource
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClient(r.destroy),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccNetworkInterfaceBackendAddressPoolAssociation_updateNIC(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface_backend_address_pool_association", "test")
	r := NetworkInterfaceBackendAddressPoolResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		// intentional as this is a Virtual Resource
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateNIC(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t NetworkInterfaceBackendAddressPoolResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	splitId := strings.Split(state.ID, "|")
	if len(splitId) != 2 {
		return nil, fmt.Errorf("expected ID to be in the format {networkInterfaceId}/ipConfigurations/{ipConfigurationName}|{backendAddressPoolId} but got %q", state.ID)
	}

	id, err := azure.ParseAzureResourceID(splitId[0])
	if err != nil {
		return nil, err
	}

	ipConfigurationName := id.Path["ipConfigurations"]
	networkInterfaceName := id.Path["networkInterfaces"]
	resourceGroup := id.ResourceGroup
	backendAddressPoolId := splitId[1]

	read, err := clients.Network.InterfacesClient.Get(ctx, resourceGroup, networkInterfaceName, "")
	if err != nil {
		return nil, fmt.Errorf("reading NetworkInterfaceApplicationGatewayBackendAddressPoolAssociation (%s): %+v", id, err)
	}

	nicProps := read.InterfacePropertiesFormat
	if nicProps == nil {
		return nil, fmt.Errorf("`properties` was nil for Network Interface (%s): %+v", id, err)
	}

	c := network2.FindNetworkInterfaceIPConfiguration(read.InterfacePropertiesFormat.IPConfigurations, ipConfigurationName)
	if c == nil {
		return nil, fmt.Errorf("IP Configuration %q wasn't found for Network Interface %q", ipConfigurationName, id)
	}
	config := *c

	found := false
	if config.InterfaceIPConfigurationPropertiesFormat.LoadBalancerBackendAddressPools != nil {
		for _, pool := range *config.InterfaceIPConfigurationPropertiesFormat.LoadBalancerBackendAddressPools {
			if *pool.ID == backendAddressPoolId {
				found = true
				break
			}
		}
	}

	return utils.Bool(found), nil
}

func (NetworkInterfaceBackendAddressPoolResource) destroy(ctx context.Context, client *clients.Client, state *terraform.InstanceState) error {
	nicID, err := parse.NetworkInterfaceID(state.Attributes["network_interface_id"])
	if err != nil {
		return err
	}

	nicName := nicID.Name
	resourceGroup := nicID.ResourceGroup
	backendAddressPoolId := state.Attributes["backend_address_pool_id"]
	ipConfigurationName := state.Attributes["ip_configuration_name"]

	read, err := client.Network.InterfacesClient.Get(ctx, resourceGroup, nicName, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Network Interface %q (Resource Group %q): %+v", nicName, resourceGroup, err)
	}

	c := network2.FindNetworkInterfaceIPConfiguration(read.InterfacePropertiesFormat.IPConfigurations, ipConfigurationName)
	if c == nil {
		return fmt.Errorf("IP Configuration %q wasn't found for Network Interface %q (Resource Group %q)", ipConfigurationName, nicName, resourceGroup)
	}
	config := *c

	updatedPools := make([]network.BackendAddressPool, 0)
	if config.InterfaceIPConfigurationPropertiesFormat.LoadBalancerBackendAddressPools != nil {
		for _, pool := range *config.InterfaceIPConfigurationPropertiesFormat.LoadBalancerBackendAddressPools {
			if *pool.ID != backendAddressPoolId {
				updatedPools = append(updatedPools, pool)
			}
		}
	}
	config.InterfaceIPConfigurationPropertiesFormat.LoadBalancerBackendAddressPools = &updatedPools

	future, err := client.Network.InterfacesClient.CreateOrUpdate(ctx, resourceGroup, nicName, read)
	if err != nil {
		return fmt.Errorf("Error removing Backend Address Pool Association for Network Interface %q (Resource Group %q): %+v", nicName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Network.InterfacesClient.Client); err != nil {
		return fmt.Errorf("Error waiting for removal of Backend Address Pool Association for NIC %q (Resource Group %q): %+v", nicName, resourceGroup, err)
	}

	return nil
}

func (r NetworkInterfaceBackendAddressPoolResource) basic(data acceptance.TestData) string {
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

resource "azurerm_network_interface_backend_address_pool_association" "test" {
  network_interface_id    = azurerm_network_interface.test.id
  ip_configuration_name   = "testconfiguration1"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkInterfaceBackendAddressPoolResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface_backend_address_pool_association" "import" {
  network_interface_id    = azurerm_network_interface_backend_address_pool_association.test.network_interface_id
  ip_configuration_name   = azurerm_network_interface_backend_address_pool_association.test.ip_configuration_name
  backend_address_pool_id = azurerm_network_interface_backend_address_pool_association.test.backend_address_pool_id
}
`, r.basic(data))
}

func (r NetworkInterfaceBackendAddressPoolResource) updateNIC(data acceptance.TestData) string {
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

resource "azurerm_network_interface_backend_address_pool_association" "test" {
  network_interface_id    = azurerm_network_interface.test.id
  ip_configuration_name   = "testconfiguration1"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id
}
`, r.template(data), data.RandomInteger)
}

func (NetworkInterfaceBackendAddressPoolResource) template(data acceptance.TestData) string {
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
  name                 = "testsubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  frontend_ip_configuration {
    name                 = "primary"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "acctestpool"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
