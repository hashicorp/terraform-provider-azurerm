package tests

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type NatGatewayPublicAssociationResource struct {
}

func TestAccNatGatewayPublicIpAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway_public_ip_association", "test")
	r := NatGatewayPublicAssociationResource{}
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

func TestAccNatGatewayPublicIpAssociation_updateNatGateway(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway_public_ip_association", "test")
	r := NatGatewayPublicAssociationResource{}
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
			Config: r.updateNatGateway(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNatGatewayPublicIpAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway_public_ip_association", "test")
	r := NatGatewayPublicAssociationResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		// intentional as this is a Virtual Resource
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccNatGatewayPublicIpAssociation_deleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway_public_ip_association", "test")
	r := NatGatewayPublicAssociationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		// intentional as this is a Virtual Resource
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				testCheckNatGatewayPublicIpAssociationDisappears(data.ResourceName),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func (t NatGatewayPublicAssociationResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.NatGatewayPublicIPAddressAssociationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.NatGatewayClient.Get(ctx, id.NatGateway.ResourceGroup, id.NatGateway.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading Nat Gateway Public IP Association (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func testCheckNatGatewayPublicIpAssociationDisappears(resourceName string) resource.TestCheckFunc {
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
