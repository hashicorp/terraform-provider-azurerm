package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
)

func TestAccAzureRMNatGatewayPublicIpAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway_public_ip_association", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		// intentional as this is a Virtual Resource
		CheckDestroy: testCheckAzureRMNatGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNatGatewayPublicIpAssociation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNatGatewayPublicIpAssociationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNatGatewayPublicIpAssociation_updateNatGateway(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway_public_ip_association", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		// intentional as this is a Virtual Resource
		CheckDestroy: testCheckAzureRMNatGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNatGatewayPublicIpAssociation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNatGatewayPublicIpAssociationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMNatGatewayPublicIpAssociation_updateNatGateway(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNatGatewayPublicIpAssociationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNatGatewayPublicIpAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway_public_ip_association", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		// intentional as this is a Virtual Resource
		CheckDestroy: testCheckAzureRMNatGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNatGatewayPublicIpAssociation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNatGatewayPublicIpAssociationExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMNatGatewayPublicIpAssociation_requiresImport),
		},
	})
}

func TestAccAzureRMNatGatewayPublicIpAssociation_deleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway_public_ip_association", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		// intentional as this is a Virtual Resource
		CheckDestroy: testCheckAzureRMNatGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNatGatewayPublicIpAssociation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNatGatewayPublicIpAssociationExists(data.ResourceName),
					testCheckAzureRMNatGatewayPublicIpAssociationDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMNatGatewayPublicIpAssociationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.NatGatewayClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.NatGatewayPublicIPAddressAssociationID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.NatGateway.ResourceGroup, id.NatGateway.Name, "")
		if err != nil {
			return fmt.Errorf("failed to retrieve Nat Gateway %q (Resource Group %q): %+v", id.NatGateway.Name, id.NatGateway.ResourceGroup, err)
		}

		if publicIpAddresses := resp.PublicIPAddresses; publicIpAddresses != nil {
			for _, publicIpAddress := range *publicIpAddresses {
				if strings.EqualFold(*publicIpAddress.ID, id.PublicIPAddressID) {
					return nil
				}
			}
		}

		return fmt.Errorf("Association between Nat Gateway %q and Public IP %q was not found.", id.NatGateway.Name, id.PublicIPAddressID)
	}
}

func testCheckAzureRMNatGatewayPublicIpAssociationDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.NatGatewayClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.NatGatewayPublicIPAddressAssociationID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.NatGateway.ResourceGroup, id.NatGateway.Name, "")
		if err != nil {
			return fmt.Errorf("failed to retrieve Nat Gateway %q (Resource Group %q): %+v", id.NatGateway.Name, id.NatGateway.ResourceGroup, err)
		}

		updatedAddresses := make([]network.SubResource, 0)
		if publicIpAddresses := resp.PublicIPAddresses; publicIpAddresses != nil {
			for _, publicIpAddress := range *publicIpAddresses {
				if !strings.EqualFold(*publicIpAddress.ID, id.PublicIPAddressID) {
					updatedAddresses = append(updatedAddresses, publicIpAddress)
				}
			}
		}
		resp.PublicIPAddresses = &updatedAddresses

		future, err := client.CreateOrUpdate(ctx, id.NatGateway.ResourceGroup, id.NatGateway.Name, resp)
		if err != nil {
			return fmt.Errorf("failed to remove Nat Gateway Public Ip Association for Nat Gateway %q (Resource Group %q): %+v", id.NatGateway.Name, id.NatGateway.ResourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("failed to wait for removal of Nat Gateway Public Ip Association for Nat Gateway %q (Resource Group %q): %+v", id.NatGateway.Name, id.NatGateway.ResourceGroup, err)
		}

		return nil
	}
}

func testAccAzureRMNatGatewayPublicIpAssociation_basic(data acceptance.TestData) string {
	template := testAccAzureRMNatGatewayPublicIpAssociation_template(data)
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
`, template, data.RandomInteger)
}

func testAccAzureRMNatGatewayPublicIpAssociation_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMNatGatewayPublicIpAssociation_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_nat_gateway_public_ip_association" "import" {
  nat_gateway_id       = azurerm_nat_gateway_public_ip_association.test.nat_gateway_id
  public_ip_address_id = azurerm_nat_gateway_public_ip_association.test.public_ip_address_id
}
`, template)
}

func testAccAzureRMNatGatewayPublicIpAssociation_updateNatGateway(data acceptance.TestData) string {
	template := testAccAzureRMNatGatewayPublicIpAssociation_template(data)
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
`, template, data.RandomInteger)
}

func testAccAzureRMNatGatewayPublicIpAssociation_template(data acceptance.TestData) string {
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
