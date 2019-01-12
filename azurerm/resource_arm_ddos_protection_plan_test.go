package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMDDosProtectionPlan_basic(t *testing.T) {
	resourceName := "azurerm_ddos_protection_plan.test"
	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDDosProtectionPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDDosProtectionPlan_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDDosProtectionPlanExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_network_ids.#"),
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

func TestAccAzureRMDDosProtectionPlan_requiresImport(t *testing.T) {
	if !requireResourcesToBeImported {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_ddos_protection_plan.test"
	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDDosProtectionPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDDosProtectionPlan_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDDosProtectionPlanExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMDDosProtectionPlan_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_ddos_protection_plan"),
			},
		},
	})
}

func TestAccAzureRMDDosProtectionPlan_withTags(t *testing.T) {
	resourceName := "azurerm_ddos_protection_plan.test"
	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDDosProtectionPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDDosProtectionPlan_withTags(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDDosProtectionPlanExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(resourceName, "tags.cost_center", "MSFT"),
				),
			},
			{
				Config: testAccAzureRMDDosProtectionPlan_withUpdatedTags(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDDosProtectionPlanExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "Staging"),
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

func TestAccAzureRMDDosProtectionPlan_disappears(t *testing.T) {
	resourceName := "azurerm_ddos_protection_plan.test"
	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDDosProtectionPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDDosProtectionPlan_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDDosProtectionPlanExists(resourceName),
					testCheckAzureRMDDosProtectionPlanDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMDDosProtectionPlanExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for DDos Protection Plan: %q", name)
		}

		client := testAccProvider.Meta().(*ArmClient).ddosProtectionPlanClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: DDos Protection Plan %q (Resource Group: %q) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on ddosProtectionPlanClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMDDosProtectionPlanDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for DDos Protection Plan: %q", name)
		}

		client := testAccProvider.Meta().(*ArmClient).ddosProtectionPlanClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		future, err := client.Delete(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Delete on ddosProtectionPlanClient: %+v", err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Bad: waiting for Deletion on ddosProtectionPlanClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMDDosProtectionPlanDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).ddosProtectionPlanClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_ddos_protection_plan" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("DDos Protection Plan still exists:\n%#v", resp.DdosProtectionPlanPropertiesFormat)
	}

	return nil
}

func testAccAzureRMDDosProtectionPlan_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_ddos_protection_plan" "test" {
  name                = "acctestddospplan-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt)
}

func testAccAzureRMDDosProtectionPlan_requiresImport(rInt int, location string) string {
	basicConfig := testAccAzureRMDDosProtectionPlan_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_ddos_protection_plan" "import" {
  name                = "${azurerm_ddos_protection_plan.test.name}"
  location            = "${azurerm_ddos_protection_plan.test.location}"
  resource_group_name = "${azurerm_ddos_protection_plan.test.resource_group_name}"
}
`, basicConfig)
}

func testAccAzureRMDDosProtectionPlan_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_ddos_protection_plan" "test" {
  name                = "acctestddospplan-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMDDosProtectionPlan_withUpdatedTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_ddos_protection_plan" "test" {
  name                = "acctestddospplan-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags {
    environment = "Staging"
  }
}
`, rInt, location, rInt)
}
