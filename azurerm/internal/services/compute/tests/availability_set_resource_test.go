package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAvailabilitySet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_availability_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAvailabilitySetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAvailabilitySet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAvailabilitySetExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "platform_update_domain_count", "5"),
					resource.TestCheckResourceAttr(data.ResourceName, "platform_fault_domain_count", "3"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAvailabilitySet_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_availability_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAvailabilitySetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAvailabilitySet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAvailabilitySetExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "platform_update_domain_count", "5"),
					resource.TestCheckResourceAttr(data.ResourceName, "platform_fault_domain_count", "3"),
				),
			},
			{
				Config:      testAccAzureRMAvailabilitySet_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_availability_set"),
			},
		},
	})
}

func TestAccAzureRMAvailabilitySet_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_availability_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAvailabilitySetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAvailabilitySet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAvailabilitySetExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "platform_update_domain_count", "5"),
					resource.TestCheckResourceAttr(data.ResourceName, "platform_fault_domain_count", "3"),
					testCheckAzureRMAvailabilitySetDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMAvailabilitySet_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_availability_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAvailabilitySetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAvailabilitySet_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAvailabilitySetExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.cost_center", "MSFT"),
				),
			},
			{
				Config: testAccAzureRMAvailabilitySet_withUpdatedTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAvailabilitySetExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "staging"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAvailabilitySet_withPPG(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_availability_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAvailabilitySetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAvailabilitySet_withPPG(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAvailabilitySetExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "proximity_placement_group_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAvailabilitySet_withDomainCounts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_availability_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAvailabilitySetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAvailabilitySet_withDomainCounts(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAvailabilitySetExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "platform_update_domain_count", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "platform_fault_domain_count", "3"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAvailabilitySet_unmanaged(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_availability_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAvailabilitySetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAvailabilitySet_unmanaged(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAvailabilitySetExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "managed", "false"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMAvailabilitySetExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.AvailabilitySetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.AvailabilitySetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Availability Set %q (Resource Group %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("Bad: Get on vmScaleSetClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMAvailabilitySetDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.AvailabilitySetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.AvailabilitySetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Delete(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !response.WasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Delete on availSetClient: %+v", err)
			}
		}

		return nil
	}
}

func testCheckAzureRMAvailabilitySetDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.AvailabilitySetsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_availability_set" {
			continue
		}

		id, err := parse.AvailabilitySetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return fmt.Errorf("Bad: Availability Set still exists:\n%#v", resp.AvailabilitySetProperties)
	}

	return nil
}

func testAccAzureRMAvailabilitySet_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_availability_set" "test" {
  name                = "acctestavset-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAvailabilitySet_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMAvailabilitySet_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_availability_set" "import" {
  name                = azurerm_availability_set.test.name
  location            = azurerm_availability_set.test.location
  resource_group_name = azurerm_availability_set.test.resource_group_name
}
`, template)
}

func testAccAzureRMAvailabilitySet_withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_availability_set" "test" {
  name                = "acctestavset-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAvailabilitySet_withUpdatedTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_availability_set" "test" {
  name                = "acctestavset-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAvailabilitySet_withPPG(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_proximity_placement_group" "test" {
  name                = "acctestPPG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_availability_set" "test" {
  name                = "acctestavset-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  proximity_placement_group_id = azurerm_proximity_placement_group.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAvailabilitySet_withDomainCounts(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_availability_set" "test" {
  name                         = "acctestavset-%d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  platform_update_domain_count = 3
  platform_fault_domain_count  = 3
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAvailabilitySet_unmanaged(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_availability_set" "test" {
  name                         = "acctestavset-%d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  platform_update_domain_count = 3
  platform_fault_domain_count  = 3
  managed                      = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
