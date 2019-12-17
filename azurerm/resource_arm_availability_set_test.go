package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAvailabilitySet_basic(t *testing.T) {
	resourceName := "azurerm_availability_set.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAvailabilitySet_basic(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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
				ExpectError: acceptance.RequiresImportError("azurerm_availability_set"),
			},
		},
	})
}

func TestAccAzureRMAvailabilitySet_disappears(t *testing.T) {
	resourceName := "azurerm_availability_set.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAvailabilitySet_basic(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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
	location := acceptance.Location()
	preConfig := testAccAzureRMAvailabilitySet_withTags(ri, location)
	postConfig := testAccAzureRMAvailabilitySet_withUpdatedTags(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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

func TestAccAzureRMAvailabilitySet_withPPG(t *testing.T) {
	resourceName := "azurerm_availability_set.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAvailabilitySet_withPPG(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAvailabilitySetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAvailabilitySetExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "proximity_placement_group_id"),
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
	config := testAccAzureRMAvailabilitySet_withDomainCounts(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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
	config := testAccAzureRMAvailabilitySet_managed(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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

		// Name of the actual scale set
		name := rs.Primary.Attributes["name"]

		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for availability set: %s", name)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.AvailabilitySetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		vmss, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on vmScaleSetClient: %+v", err)
		}

		if vmss.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: VirtualMachineScaleSet %q (resource group: %q) does not exist", name, resourceGroup)
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

		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.AvailabilitySetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
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

		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.AvailabilitySetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
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

func testAccAzureRMAvailabilitySet_withPPG(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_proximity_placement_group" "test" {
  name                = "acctestPPG-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_availability_set" "test" {
  name                = "acctestavset-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  proximity_placement_group_id = "${azurerm_proximity_placement_group.test.id}"
}
`, rInt, location, rInt, rInt)
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
