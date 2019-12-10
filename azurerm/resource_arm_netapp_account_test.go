package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
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
	resourceName := "azurerm_netapp_account.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetAppAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppAccount_basicConfig(ri, location),
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

func testAccAzureRMNetAppAccount_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_netapp_account.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetAppAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppAccount_basicConfig(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppAccountExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMNetAppAccount_requiresImportConfig(ri, testLocation()),
				ExpectError: testRequiresImportError("azurerm_netapp_account"),
			},
		},
	})
}

func testAccAzureRMNetAppAccount_complete(t *testing.T) {
	resourceName := "azurerm_netapp_account.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetAppAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppAccount_completeConfig(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "active_directory.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"active_directory",
				},
			},
		},
	})
}

func testAccAzureRMNetAppAccount_update(t *testing.T) {
	resourceName := "azurerm_netapp_account.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetAppAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppAccount_basicConfig(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "active_directory.#", "0"),
				),
			},
			{
				Config: testAccAzureRMNetAppAccount_completeConfig(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "active_directory.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"active_directory",
				},
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

		client := testAccProvider.Meta().(*ArmClient).NetApp.AccountClient
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
	client := testAccProvider.Meta().(*ArmClient).NetApp.AccountClient
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

func testAccAzureRMNetAppAccount_basicConfig(rInt int, location string) string {
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

func testAccAzureRMNetAppAccount_requiresImportConfig(rInt int, location string) string {
	return fmt.Sprintf(`
%s
resource "azurerm_netapp_account" "import" {
  name                = "${azurerm_netapp_account.test.name}"
  location            = "${azurerm_netapp_account.test.location}"
  resource_group_name = "${azurerm_netapp_account.test.name}"
}
}
`, testAccAzureRMNetAppAccount_basicConfig(rInt, location))
}

func testAccAzureRMNetAppAccount_completeConfig(rInt int, location string) string {
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
}
`, rInt, location, rInt)
}
