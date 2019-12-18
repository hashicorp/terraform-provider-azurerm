package azurerm

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

const ArmAccCapacityAdminEmail = "ARM_ACCTEST_CAPACITY_ADMIN_EMAIL"
const ArmAccCapacityAdminServicePrincipal = "ARM_ACCTEST_ADMIN_SERVICE_PRINCIPAL"

func TestAccAzureRMPowerBIDedicatedCapacity_basic(t *testing.T) {
	resourceName := "azurerm_powerbi_dedicated_capacity.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	capacityAdminEmail := os.Getenv(ArmAccCapacityAdminEmail)
	if capacityAdminEmail == "" {
		t.Skip(fmt.Sprintf("Acceptance test skipped unless env '%s' set", ArmAccCapacityAdminEmail))
		return
	}
	capacityAdmins := []string{capacityAdminEmail}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPowerBIDedicatedCapacityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPowerBIDedicatedCapacity_basic(ri, location, capacityAdmins),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPowerBIDedicatedCapacityExists(resourceName),
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

func TestAccAzureRMPowerBIDedicatedCapacity_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_powerbi_dedicated_capacity.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	capacityAdminEmail := os.Getenv(ArmAccCapacityAdminEmail)
	if capacityAdminEmail == "" {
		t.Skip(fmt.Sprintf("Acceptance test skipped unless env '%s' set", ArmAccCapacityAdminEmail))
		return
	}
	capacityAdmins := []string{capacityAdminEmail}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPowerBIDedicatedCapacityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPowerBIDedicatedCapacity_basic(ri, location, capacityAdmins),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPowerBIDedicatedCapacityExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMPowerBIDedicatedCapacity_requiresImport(ri, location, capacityAdmins),
				ExpectError: testRequiresImportError("azurerm_powerbi_dedicated_capacity"),
			},
		},
	})
}

func TestAccAzureRMPowerBIDedicatedCapacity_complete(t *testing.T) {
	resourceName := "azurerm_powerbi_dedicated_capacity.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	capacityAdminServicePrincipal := os.Getenv(ArmAccCapacityAdminServicePrincipal)
	if capacityAdminServicePrincipal == "" {
		t.Skip(fmt.Sprintf("Acceptance test skipped unless env '%s' set", ArmAccCapacityAdminServicePrincipal))
		return
	}
	capacityAdmins := []string{capacityAdminServicePrincipal}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPowerBIDedicatedCapacityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPowerBIDedicatedCapacity_complete(ri, location, capacityAdmins),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPowerBIDedicatedCapacityExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku_name", "A2"),
					resource.TestCheckResourceAttr(resourceName, "administrators.#", "1"),
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

func TestAccAzureRMPowerBIDedicatedCapacity_update(t *testing.T) {
	resourceName := "azurerm_powerbi_dedicated_capacity.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	capacityAdminEmail := os.Getenv(ArmAccCapacityAdminEmail)
	capacityAdminServicePrincipal := os.Getenv(ArmAccCapacityAdminServicePrincipal)
	if capacityAdminEmail == "" || capacityAdminServicePrincipal == "" {
		t.Skip(fmt.Sprintf("Acceptance test skipped unless env '%s' and '%s' set", ArmAccCapacityAdminEmail, ArmAccCapacityAdminServicePrincipal))
		return
	}
	preCapacityAdmins := []string{capacityAdminEmail, capacityAdminServicePrincipal}
	postCapacityAdmins := []string{capacityAdminEmail}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPowerBIDedicatedCapacityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPowerBIDedicatedCapacity_basic(ri, location, preCapacityAdmins),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPowerBIDedicatedCapacityExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku_name", "A1"),
					resource.TestCheckResourceAttr(resourceName, "administrators.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMPowerBIDedicatedCapacity_complete(ri, location, postCapacityAdmins),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPowerBIDedicatedCapacityExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku_name", "A2"),
					resource.TestCheckResourceAttr(resourceName, "administrators.#", "1"),
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

func testCheckAzureRMPowerBIDedicatedCapacityExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Capacity not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).PowerBIDedicated.CapacityClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		if resp, err := client.GetDetails(ctx, resourceGroup, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Capacity (Capacity Name %q / Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on PowerBI Dedicated.CapacityClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMPowerBIDedicatedCapacityDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).PowerBIDedicated.CapacityClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_powerbi_dedicated_capacity" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.GetDetails(ctx, resourceGroup, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on CapacityClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMPowerBIDedicatedCapacity_basic(rInt int, location string, capacityAdmins []string) string {
	template := testAccAzureRMPowerBIDedicatedCapacity_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_powerbi_dedicated_capacity" "test" {
  name                = "acctestpowerbidedicatedcapacity%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku_name            = "A1"
  administrators      = ["%s"]
}
`, template, rInt, strings.Join(capacityAdmins, "\", \""))
}

func testAccAzureRMPowerBIDedicatedCapacity_requiresImport(rInt int, location string, capacityAdmins []string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_powerbi_dedicated_capacity" "import" {
  name                = "${azurerm_powerbi_dedicated_capacity.test.name}"
  location            = "${azurerm_powerbi_dedicated_capacity.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, testAccAzureRMPowerBIDedicatedCapacity_basic(rInt, location, capacityAdmins))
}

func testAccAzureRMPowerBIDedicatedCapacity_complete(rInt int, location string, capacityAdmins []string) string {
	template := testAccAzureRMPowerBIDedicatedCapacity_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_powerbi_dedicated_capacity" "test" {
  name                = "acctestpowerbidedicatedcapacity%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku_name            = "A2"
  administrators      = ["%s"]

  tags = {
    ENV = "Test"
  }
}
`, template, rInt, strings.Join(capacityAdmins, "\", \""))
}

func testAccAzureRMPowerBIDedicatedCapacity_template(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-powerbidedicated-%d"
  location = "%s"
}
`, rInt, location)
}
