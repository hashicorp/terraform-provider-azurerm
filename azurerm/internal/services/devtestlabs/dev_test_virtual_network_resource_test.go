package devtestlabs_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/devtestlabs"
)

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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDevTestVirtualNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDevTestVirtualNetwork_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDevTestVirtualNetworkExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccDevTestVirtualNetwork_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_virtual_network", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDevTestVirtualNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDevTestVirtualNetwork_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDevTestVirtualNetworkExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			{
				Config:      testAccDevTestVirtualNetwork_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_dev_test_virtual_network"),
			},
		},
	})
}

func TestAccDevTestVirtualNetwork_subnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_virtual_network", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDevTestVirtualNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDevTestVirtualNetwork_subnets(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDevTestVirtualNetworkExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "subnet.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "subnet.0.use_public_ip_address", "Deny"),
					resource.TestCheckResourceAttr(data.ResourceName, "subnet.0.use_in_virtual_machine_creation", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckDevTestVirtualNetworkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).DevTestLabs.VirtualNetworksClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		virtualNetworkName := rs.Primary.Attributes["name"]
		labName := rs.Primary.Attributes["lab_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, labName, virtualNetworkName, "")
		if err != nil {
			return fmt.Errorf("Bad: Get devTestVirtualNetworksClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: DevTest Virtual Network %q (Lab %q / Resource Group: %q) does not exist", virtualNetworkName, labName, resourceGroup)
		}

		return nil
	}
}

func testCheckDevTestVirtualNetworkDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).DevTestLabs.VirtualNetworksClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dev_test_virtual_network" {
			continue
		}

		virtualNetworkName := rs.Primary.Attributes["name"]
		labName := rs.Primary.Attributes["lab_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, labName, virtualNetworkName, "")
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("DevTest Virtual Network still exists:\n%#v", resp)
	}

	return nil
}

func testAccDevTestVirtualNetwork_basic(data acceptance.TestData) string {
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

func testAccDevTestVirtualNetwork_requiresImport(data acceptance.TestData) string {
	template := testAccDevTestVirtualNetwork_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_virtual_network" "import" {
  name                = azurerm_dev_test_virtual_network.test.name
  lab_name            = azurerm_dev_test_virtual_network.test.lab_name
  resource_group_name = azurerm_dev_test_virtual_network.test.resource_group_name
}
`, template)
}

func testAccDevTestVirtualNetwork_subnets(data acceptance.TestData) string {
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
