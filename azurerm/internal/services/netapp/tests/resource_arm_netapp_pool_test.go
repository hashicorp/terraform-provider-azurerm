package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMNetAppPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetAppPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppPool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppPoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetAppPool_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_netapp_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetAppPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppPool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppPoolExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMNetAppPool_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_netapp_pool"),
			},
		},
	})
}

func TestAccAzureRMNetAppPool_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetAppPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppPool_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "service_level", "Standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "size_in_tb", "15"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetAppPool_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetAppPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppPool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "service_level", "Premium"),
					resource.TestCheckResourceAttr(data.ResourceName, "size_in_tb", "4"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMNetAppPool_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "service_level", "Standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "size_in_tb", "15"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMNetAppPoolExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).NetApp.PoolClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NetApp Pool not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["account_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, accountName, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: NetApp Pool %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on netapp.PoolClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMNetAppPoolDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).NetApp.PoolClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_netapp_pool" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["account_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, accountName, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on netapp.PoolClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMNetAppPool_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-netapp-%d"
  location = "%s"
}

resource "azurerm_netapp_account" "test" {
  name                = "acctest-NetAppAccount-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_netapp_pool" "test" {
  name                = "acctest-NetAppPool-%d"
  account_name        = "${azurerm_netapp_account.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  service_level       = "Premium"
  size_in_tb          = 4
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMNetAppPool_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_netapp_pool" "import" {
  name                = "${azurerm_netapp_pool.test.name}"
  location            = "${azurerm_netapp_pool.test.location}"
  resource_group_name = "${azurerm_netapp_pool.test.resource_group_name}"
}
`, testAccAzureRMNetAppPool_basic(data))
}

func testAccAzureRMNetAppPool_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-netapp-%d"
  location = "%s"
}

resource "azurerm_netapp_account" "test" {
  name                = "acctest-NetAppAccount-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_netapp_pool" "test" {
  name                = "acctest-NetAppPool-%d"
  account_name        = "${azurerm_netapp_account.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  service_level       = "Standard"
  size_in_tb          = 15
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
