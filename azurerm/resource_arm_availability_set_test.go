package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAvailabilitySet_basic(t *testing.T) {
	resourceName := "azurerm_availability_set.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAvailabilitySet_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAvailabilitySetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAvailabilitySetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "platform_update_domain_count", "5"),
					resource.TestCheckResourceAttr(resourceName, "platform_fault_domain_count", "3"),
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

func TestAccAzureRMAvailabilitySet_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_availability_set.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAvailabilitySetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAvailabilitySet_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAvailabilitySetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "platform_update_domain_count", "5"),
					resource.TestCheckResourceAttr(resourceName, "platform_fault_domain_count", "3"),
				),
			},
			{
				Config:      testAccAzureRMAvailabilitySet_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_availability_set"),
			},
		},
	})
}

func TestAccAzureRMAvailabilitySet_disappears(t *testing.T) {
	resourceName := "azurerm_availability_set.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAvailabilitySet_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAvailabilitySetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAvailabilitySetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "platform_update_domain_count", "5"),
					resource.TestCheckResourceAttr(resourceName, "platform_fault_domain_count", "3"),
					testCheckAzureRMAvailabilitySetDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMAvailabilitySet_withTags(t *testing.T) {
	resourceName := "azurerm_availability_set.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	preConfig := testAccAzureRMAvailabilitySet_withTags(ri, location)
	postConfig := testAccAzureRMAvailabilitySet_withUpdatedTags(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAvailabilitySetDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAvailabilitySetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(resourceName, "tags.cost_center", "MSFT"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAvailabilitySetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "staging"),
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

func TestAccAzureRMAvailabilitySet_withDomainCounts(t *testing.T) {
	resourceName := "azurerm_availability_set.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAvailabilitySet_withDomainCounts(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAvailabilitySetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAvailabilitySetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "platform_update_domain_count", "3"),
					resource.TestCheckResourceAttr(resourceName, "platform_fault_domain_count", "3"),
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

func TestAccAzureRMAvailabilitySet_managed(t *testing.T) {
	resourceName := "azurerm_availability_set.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAvailabilitySet_managed(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAvailabilitySetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAvailabilitySetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "managed", "true"),
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

func testCheckAzureRMAvailabilitySetExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		availSetName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for availability set: %s", availSetName)
		}

		client := testAccProvider.Meta().(*ArmClient).compute.AvailabilitySetsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, availSetName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Availability Set %q (resource group: %q) does not exist", availSetName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on availSetClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMAvailabilitySetDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		availSetName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for availability set: %s", availSetName)
		}

		client := testAccProvider.Meta().(*ArmClient).compute.AvailabilitySetsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Delete(ctx, resourceGroup, availSetName)
		if err != nil {
			if !response.WasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Delete on availSetClient: %+v", err)
			}
		}

		return nil
	}
}

func testCheckAzureRMAvailabilitySetDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_availability_set" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).compute.AvailabilitySetsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return fmt.Errorf("Availability Set still exists:\n%#v", resp.AvailabilitySetProperties)
	}

	return nil
}

func testAccAzureRMAvailabilitySet_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_availability_set" "test" {
  name                = "acctestavset-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt)
}

func testAccAzureRMAvailabilitySet_requiresImport(rInt int, location string) string {
	template := testAccAzureRMAvailabilitySet_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_availability_set" "import" {
  name                = "${azurerm_availability_set.test.name}"
  location            = "${azurerm_availability_set.test.location}"
  resource_group_name = "${azurerm_availability_set.test.resource_group_name}"
}
`, template)
}

func testAccAzureRMAvailabilitySet_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_availability_set" "test" {
  name                = "acctestavset-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMAvailabilitySet_withUpdatedTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_availability_set" "test" {
  name                = "acctestavset-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags = {
    environment = "staging"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMAvailabilitySet_withDomainCounts(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_availability_set" "test" {
  name                         = "acctestavset-%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  platform_update_domain_count = 3
  platform_fault_domain_count  = 3
}
`, rInt, location, rInt)
}

func testAccAzureRMAvailabilitySet_managed(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_availability_set" "test" {
  name                         = "acctestavset-%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  platform_update_domain_count = 3
  platform_fault_domain_count  = 3
  managed                      = true
}
`, rInt, location, rInt)
}
