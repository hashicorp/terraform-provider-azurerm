package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMNetAppAccount_basic(t *testing.T) {
	resourceName := "azurerm_netapp_account.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetAppAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppAccount_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppAccountExists(resourceName),
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

func TestAccAzureRMNetAppAccount_complete(t *testing.T) {
	resourceName := "azurerm_netapp_account.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetAppAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppAccount_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "active_directory.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.ENV", "Test"),
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

func TestAccAzureRMNetAppAccount_update(t *testing.T) {
	resourceName := "azurerm_netapp_account.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetAppAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppAccount_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "active_directory.#", "0"),
				),
			},
			{
				Config: testAccAzureRMNetAppAccount_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "active_directory.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.ENV", "Test"),
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

func testCheckAzureRMNetAppAccountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NetApp Account not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).Netapp.AccountClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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
	client := testAccProvider.Meta().(*ArmClient).Netapp.AccountClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

func testAccAzureRMNetAppAccount_basic(rInt int, location string) string {
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
`, rInt, location, rInt)
}

func testAccAzureRMNetAppAccount_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-netapp-%d"
  location = "%s"
}

resource "azurerm_netapp_account" "test" {
  name                = "acctest-NetAppAccount-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  active_directory {
    username            = "aduser"
    password            = "aduserpwd"
    smb_server_name     = "SMBSERVER"
    dns_servers         = ["1.2.3.4"]
	domain              = "westcentralus.com"
	organizational_unit = "OU=FirstLevel"
  }

  tags = {
    ENV = "Test"
  }
}
`, rInt, location, rInt)
}
