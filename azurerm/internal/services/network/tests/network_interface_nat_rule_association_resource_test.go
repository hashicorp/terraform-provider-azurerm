package tests

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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type NetworkInterfaceNATRuleAssociationResource struct {
}

func TestAccNetworkInterfaceNATRuleAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface_nat_rule_association", "test")
	r := NetworkInterfaceNATRuleAssociationResource{}
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

func TestAccNetworkInterfaceNATRuleAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface_nat_rule_association", "test")
	r := NetworkInterfaceNATRuleAssociationResource{}
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
			ExpectError: acceptance.RequiresImportError("azurerm_network_interface_nat_rule_association"),
		},
	})
}

func TestAccNetworkInterfaceNATRuleAssociation_deleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface_nat_rule_association", "test")
	r := NetworkInterfaceNATRuleAssociationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		// intentional as this is a Virtual Resource
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				testCheckNetworkInterfaceNATRuleAssociationDisappears(data.ResourceName),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccNetworkInterfaceNATRuleAssociation_updateNIC(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface_nat_rule_association", "test")
	r := NetworkInterfaceNATRuleAssociationResource{}
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

func (t NetworkInterfaceNATRuleAssociationResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	splitId := strings.Split(state.ID, "|")
	if len(splitId) != 2 {
		return nil, fmt.Errorf("expected ID to be in the format {networkInterfaceId}|{networkSecurityGroupId} but got %q", state.ID)
	}

	nicID, err := azure.ParseAzureResourceID(splitId[0])
	if err != nil {
		return nil, err
	}

	ipConfigurationName := nicID.Path["ipConfigurations"]
	networkInterfaceName := nicID.Path["networkInterfaces"]
	resourceGroup := nicID.ResourceGroup
	natRuleId := splitId[1]

	read, err := clients.Network.InterfacesClient.Get(ctx, resourceGroup, networkInterfaceName, "")
	if err != nil {
		return nil, fmt.Errorf("retrieving Network Interface %q: %+v", nicID, err)
	}

	c := azure.FindNetworkInterfaceIPConfiguration(read.InterfacePropertiesFormat.IPConfigurations, ipConfigurationName)
	if c == nil {
		return nil, fmt.Errorf("IP Configuration %q wasn't found for Network Interface %q (Resource Group %q)", ipConfigurationName, networkInterfaceName, resourceGroup)
	}
	config := *c

	found := false
	if config.InterfaceIPConfigurationPropertiesFormat.LoadBalancerInboundNatRules != nil {
		for _, rule := range *config.InterfaceIPConfigurationPropertiesFormat.LoadBalancerInboundNatRules {
			if *rule.ID == natRuleId {
				found = true
				break
			}
		}
	}

	return utils.Bool(found), nil
}

func testCheckNetworkInterfaceNATRuleAssociationDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.InterfacesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		nicID, err := azure.ParseAzureResourceID(rs.Primary.Attributes["network_interface_id"])
		if err != nil {
			return err
		}

		nicName := nicID.Path["networkInterfaces"]
		resourceGroup := nicID.ResourceGroup
		ipConfigurationName := rs.Primary.Attributes["ip_configuration_name"]
		natRuleId := rs.Primary.Attributes["nat_rule_id"]

		read, err := client.Get(ctx, resourceGroup, nicName, "")
		if err != nil {
			return fmt.Errorf("Error retrieving Network Interface %q (Resource Group %q): %+v", nicName, resourceGroup, err)
		}

		c := azure.FindNetworkInterfaceIPConfiguration(read.InterfacePropertiesFormat.IPConfigurations, ipConfigurationName)
		if c == nil {
			return fmt.Errorf("IP Configuration %q wasn't found for Network Interface %q (Resource Group %q)", ipConfigurationName, nicName, resourceGroup)
		}
		config := *c

		updatedRules := make([]network.InboundNatRule, 0)
		if config.InterfaceIPConfigurationPropertiesFormat.LoadBalancerInboundNatRules != nil {
			for _, rule := range *config.InterfaceIPConfigurationPropertiesFormat.LoadBalancerInboundNatRules {
				if *rule.ID != natRuleId {
					updatedRules = append(updatedRules, rule)
				}
			}
		}
		config.InterfaceIPConfigurationPropertiesFormat.LoadBalancerInboundNatRules = &updatedRules

		future, err := client.CreateOrUpdate(ctx, resourceGroup, nicName, read)
		if err != nil {
			return fmt.Errorf("Error removing NAT Rule Association for Network Interface %q (Resource Group %q): %+v", nicName, resourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for removal of NAT Rule Association for NIC %q (Resource Group %q): %+v", nicName, resourceGroup, err)
		}

		return nil
	}
}

func (r NetworkInterfaceNATRuleAssociationResource) basic(data acceptance.TestData) string {
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

resource "azurerm_network_interface_nat_rule_association" "test" {
  network_interface_id  = azurerm_network_interface.test.id
  ip_configuration_name = "testconfiguration1"
  nat_rule_id           = azurerm_lb_nat_rule.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkInterfaceNATRuleAssociationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface_nat_rule_association" "import" {
  network_interface_id  = azurerm_network_interface_nat_rule_association.test.network_interface_id
  ip_configuration_name = azurerm_network_interface_nat_rule_association.test.ip_configuration_name
  nat_rule_id           = azurerm_network_interface_nat_rule_association.test.nat_rule_id
}
`, r.basic(data))
}

func (r NetworkInterfaceNATRuleAssociationResource) updateNIC(data acceptance.TestData) string {
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

resource "azurerm_network_interface_nat_rule_association" "test" {
  network_interface_id  = azurerm_network_interface.test.id
  ip_configuration_name = "testconfiguration1"
  nat_rule_id           = azurerm_lb_nat_rule.test.id
}
`, r.template(data), data.RandomInteger)
}

func (NetworkInterfaceNATRuleAssociationResource) template(data acceptance.TestData) string {
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

resource "azurerm_lb_nat_rule" "test" {
  resource_group_name            = azurerm_resource_group.test.name
  loadbalancer_id                = azurerm_lb.test.id
  name                           = "RDPAccess"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = "primary"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
