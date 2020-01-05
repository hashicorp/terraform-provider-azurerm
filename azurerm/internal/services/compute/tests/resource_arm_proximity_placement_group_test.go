package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccProximityPlacementGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_proximity_placement_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMProximityPlacementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProximityPlacementGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMProximityPlacementGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccProximityPlacementGroup_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_proximity_placement_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMProximityPlacementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProximityPlacementGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMProximityPlacementGroupExists(data.ResourceName),
				),
			},
			{
				Config:      testAccProximityPlacementGroup_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_proximity_placement_group"),
			},
		},
	})
}

func TestAccProximityPlacementGroup_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_proximity_placement_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMProximityPlacementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProximityPlacementGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMProximityPlacementGroupExists(data.ResourceName),
					testCheckAzureRMProximityPlacementGroupDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccProximityPlacementGroup_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_proximity_placement_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMProximityPlacementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProximityPlacementGroup_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMProximityPlacementGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.cost_center", "MSFT"),
				),
			},
			{
				Config: testAccProximityPlacementGroup_withUpdatedTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMProximityPlacementGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "staging"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMProximityPlacementGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		ppgName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for proximity placement group: %s", ppgName)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.ProximityPlacementGroupsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := client.Get(ctx, resourceGroup, ppgName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Availability Set %q (resource group: %q) does not exist", ppgName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on Proximity Placement Groups Client: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMProximityPlacementGroupDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		ppgName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for proximity placement group: %s", ppgName)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.ProximityPlacementGroupsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := client.Delete(ctx, resourceGroup, ppgName)
		if err != nil {
			if !response.WasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Delete on Proximity Placement Groups Client: %+v", err)
			}
		}

		return nil
	}
}

func testCheckAzureRMProximityPlacementGroupDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_proximity_placement_group" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.ProximityPlacementGroupsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return fmt.Errorf("Proximity placement group still exists:\n%#v", resp.ProximityPlacementGroupProperties)
	}

	return nil
}

func testAccProximityPlacementGroup_basic(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccProximityPlacementGroup_requiresImport(data acceptance.TestData) string {
	template := testAccProximityPlacementGroup_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_proximity_placement_group" "import" {
  name                = "${azurerm_proximity_placement_group.test.name}"
  location            = "${azurerm_proximity_placement_group.test.location}"
  resource_group_name = "${azurerm_proximity_placement_group.test.resource_group_name}"
}
`, template)
}

func testAccProximityPlacementGroup_withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_proximity_placement_group" "test" {
  name                = "acctestPPG-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccProximityPlacementGroup_withUpdatedTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_proximity_placement_group" "test" {
  name                = "acctestPPG-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
