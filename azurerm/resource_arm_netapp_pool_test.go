package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMNetAppPool_basic(t *testing.T) {
	resourceName := "azurerm_netapp_pool.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetAppPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppPool_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppPoolExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMNetAppPool_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_netapp_pool.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetAppPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppPool_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppPoolExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMNetAppPool_requiresImport(ri, acceptance.Location()),
				ExpectError: acceptance.RequiresImportError("azurerm_netapp_pool"),
			},
		},
	})
}

func TestAccAzureRMNetAppPool_complete(t *testing.T) {
	resourceName := "azurerm_netapp_pool.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetAppPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppPool_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "service_level", "Standard"),
					resource.TestCheckResourceAttr(resourceName, "size_in_tb", "15"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMNetAppPool_update(t *testing.T) {
	resourceName := "azurerm_netapp_pool.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetAppPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppPool_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "service_level", "Premium"),
					resource.TestCheckResourceAttr(resourceName, "size_in_tb", "4"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMNetAppPool_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "service_level", "Standard"),
					resource.TestCheckResourceAttr(resourceName, "size_in_tb", "15"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMNetAppPoolExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NetApp Pool not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["account_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).NetApp.PoolClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testAccAzureRMNetAppPool_basic(rInt int, location string) string {
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
`, rInt, location, rInt, rInt)
}

func testAccAzureRMNetAppPool_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s
resource "azurerm_netapp_pool" "import" {
  name                = "${azurerm_netapp_pool.test.name}"
  location            = "${azurerm_netapp_pool.test.location}"
  resource_group_name = "${azurerm_netapp_pool.test.name}"
}
}
`, testAccAzureRMNetAppPool_basic(rInt, location))
}

func testAccAzureRMNetAppPool_complete(rInt int, location string) string {
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
`, rInt, location, rInt, rInt)
}
