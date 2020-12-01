package devtestlabs_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/devtestlabs"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type DevTestVirtualNetworkResource struct {
}

func TestValidateDevTestVirtualNetworkName(t *testing.T) {
	validNames := []string{
		"valid-name",
		"valid02-name",
		"validName1",
		"-validname1",
		"valid_name",
		"double-hyphen--valid",
	}
	for _, v := range validNames {
		_, errors := devtestlabs.ValidateDevTestVirtualNetworkName()(v, "example")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Dev Test Virtual Network Name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"invalid!",
		"!@£",
	}
	for _, v := range invalidNames {
		_, errors := devtestlabs.ValidateDevTestVirtualNetworkName()(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid Dev Test Virtual Network Name", v)
		}
	}
}

func TestAccDevTestVirtualNetwork_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_virtual_network", "test")
	r := DevTestVirtualNetworkResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDevTestVirtualNetwork_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_virtual_network", "test")
	r := DevTestVirtualNetworkResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_dev_test_virtual_network"),
		},
	})
}

func TestAccDevTestVirtualNetwork_subnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_virtual_network", "test")
	r := DevTestVirtualNetworkResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.subnets(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subnet.#").HasValue("1"),
				check.That(data.ResourceName).Key("subnet.0.use_public_ip_address").HasValue("Deny"),
				check.That(data.ResourceName).Key("subnet.0.use_in_virtual_machine_creation").HasValue("Allow"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func (DevTestVirtualNetworkResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	labName := id.Path["labs"]
	name := id.Path["virtualnetworks"]

	resp, err := clients.DevTestLabs.VirtualNetworksClient.Get(ctx, id.ResourceGroup, labName, name, "")
	if err != nil {
		return nil, fmt.Errorf("retrieving DevTest Virtual Network %q (Lab %q / Resource Group: %q) does not exist", name, labName, id.ResourceGroup)
	}

	return utils.Bool(resp.VirtualNetworkProperties != nil), nil
}

func (DevTestVirtualNetworkResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctestdtl%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_dev_test_virtual_network" "test" {
  name                = "acctestdtvn%d"
  lab_name            = azurerm_dev_test_lab.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r DevTestVirtualNetworkResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_virtual_network" "import" {
  name                = azurerm_dev_test_virtual_network.test.name
  lab_name            = azurerm_dev_test_virtual_network.test.lab_name
  resource_group_name = azurerm_dev_test_virtual_network.test.resource_group_name
}
`, r.basic(data))
}

func (DevTestVirtualNetworkResource) subnets(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctestdtl%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_dev_test_virtual_network" "test" {
  name                = "acctestdtvn%d"
  lab_name            = azurerm_dev_test_lab.test.name
  resource_group_name = azurerm_resource_group.test.name

  subnet {
    use_public_ip_address           = "Deny"
    use_in_virtual_machine_creation = "Allow"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
