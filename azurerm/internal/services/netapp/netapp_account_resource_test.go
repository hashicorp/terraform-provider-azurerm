package netapp_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMNetAppAccount(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests since
	// Azure allows only one active directory can be joined to a single subscription at a time for NetApp Account.
	// The CI system runs all tests in parallel, so the tests need to be changed to run one at a time.
	testCases := map[string]map[string]func(t *testing.T){
		"Resource": {
			"basic":          testAccAzureRMNetAppAccount_basic,
			"requiresImport": testAccAzureRMNetAppAccount_requiresImport,
			"complete":       testAccAzureRMNetAppAccount_complete,
			"update":         testAccAzureRMNetAppAccount_update,
		},
		"DataSource": {
			"basic": testAccDataSourceAzureRMNetAppAccount_basic,
		},
	}

	for group, m := range testCases {
		for name, tc := range m {
			t.Run(group, func(t *testing.T) {
				t.Run(name, func(t *testing.T) {
					tc(t)
				})
			})
		}
	}
}

func testAccAzureRMNetAppAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetAppAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppAccount_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMNetAppAccount_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetAppAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppAccount_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppAccountExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMNetAppAccount_requiresImportConfig(data),
				ExpectError: acceptance.RequiresImportError("azurerm_netapp_account"),
			},
		},
	})
}

func testAccAzureRMNetAppAccount_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetAppAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppAccount_completeConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "active_directory.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.FoO", "BaR"),
				),
			},
			data.ImportStep("active_directory"),
		},
	})
}

func testAccAzureRMNetAppAccount_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetAppAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppAccount_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "active_directory.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			{
				Config: testAccAzureRMNetAppAccount_completeConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "active_directory.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.FoO", "BaR"),
				),
			},
			data.ImportStep("active_directory"),
		},
	})
}

func testCheckAzureRMNetAppAccountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).NetApp.AccountClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NetApp Account not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: NetApp Account %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on netapp.AccountClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMNetAppAccountDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).NetApp.AccountClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_netapp_account" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on netapp.AccountClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMNetAppAccount_basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-netapp-%d"
  location = "%s"
}

resource "azurerm_netapp_account" "test" {
  name                = "acctest-NetAppAccount-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMNetAppAccount_requiresImportConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_netapp_account" "import" {
  name                = azurerm_netapp_account.test.name
  location            = azurerm_netapp_account.test.location
  resource_group_name = azurerm_netapp_account.test.resource_group_name
}
`, testAccAzureRMNetAppAccount_basicConfig(data))
}

func testAccAzureRMNetAppAccount_completeConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-netapp-%d"
  location = "%s"
}

resource "azurerm_netapp_account" "test" {
  name                = "acctest-NetAppAccount-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  active_directory {
    username            = "aduser"
    password            = "aduserpwd"
    smb_server_name     = "SMBSERVER"
    dns_servers         = ["1.2.3.4"]
    domain              = "westcentralus.com"
    organizational_unit = "OU=FirstLevel"
  }

  tags = {
    "FoO" = "BaR"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
