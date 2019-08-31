package azurerm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAutomationAccount_basic(t *testing.T) {
	resourceName := "azurerm_automation_account.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationAccount_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku_name", "Basic"),
					resource.TestCheckResourceAttrSet(resourceName, "dsc_server_endpoint"),
					resource.TestCheckResourceAttrSet(resourceName, "dsc_primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "dsc_secondary_access_key"),
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

// Remove in 2.0
func TestAccAzureRMAutomationAccount_basicClassic(t *testing.T) {
	resourceName := "azurerm_automation_account.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationAccount_basicClassic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Basic"),
					resource.TestCheckResourceAttrSet(resourceName, "dsc_server_endpoint"),
					resource.TestCheckResourceAttrSet(resourceName, "dsc_primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "dsc_secondary_access_key"),
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

// Remove in 2.0
func TestAccAzureRMAutomationAccount_basicNotDefined(t *testing.T) {
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMAutomationAccount_basicNotDefined(ri, testLocation()),
				ExpectError: regexp.MustCompile("either 'sku_name' or 'sku' must be defined in the configuration file"),
			},
		},
	})
}

func TestAccAzureRMAutomationAccount_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_automation_account.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationAccount_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationAccountExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMAutomationAccount_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_automation_account"),
			},
		},
	})
}

func TestAccAzureRMAutomationAccount_complete(t *testing.T) {
	resourceName := "azurerm_automation_account.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationAccount_complete(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku_name", "Basic"),
					resource.TestCheckResourceAttr(resourceName, "tags.hello", "world"),
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

func testCheckAzureRMAutomationAccountDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).automation.AccountClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_automation_account" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Automation Account still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMAutomationAccountExists(resourceName string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Automation Account: '%s'", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).automation.AccountClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Automation Account '%s' (resource group: '%s') was not found: %+v", name, resourceGroup, err)
			}

			return fmt.Errorf("Bad: Get on automationClient: %s", err)
		}

		return nil
	}
}

func testAccAzureRMAutomationAccount_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku_name = "Basic"
}
`, rInt, location, rInt)
}

// Remove in 2.0
func testAccAzureRMAutomationAccount_basicClassic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name = "Basic"
  } 
}
`, rInt, location, rInt)
}

// Remove in 2.0
func testAccAzureRMAutomationAccount_basicNotDefined(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt)
}

func testAccAzureRMAutomationAccount_requiresImport(rInt int, location string) string {
	template := testAccAzureRMAutomationAccount_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_account" "import" {
  name                = "${azurerm_automation_account.test.name}"
  location            = "${azurerm_automation_account.test.location}"
  resource_group_name = "${azurerm_automation_account.test.resource_group_name}"

  sku_name = "Basic"
}
`, template)
}

func testAccAzureRMAutomationAccount_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku_name = "Basic"
  
  tags = {
    "hello" = "world"
  }
}
`, rInt, location, rInt)
}
