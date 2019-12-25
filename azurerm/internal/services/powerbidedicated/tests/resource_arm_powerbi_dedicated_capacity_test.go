package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMPowerBIDedicatedCapacity_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_powerbi_dedicated_capacity", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPowerBIDedicatedCapacityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPowerBIDedicatedCapacity_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPowerBIDedicatedCapacityExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPowerBIDedicatedCapacity_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_powerbi_dedicated_capacity", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPowerBIDedicatedCapacityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPowerBIDedicatedCapacity_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPowerBIDedicatedCapacityExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMPowerBIDedicatedCapacity_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_powerbi_dedicated_capacity"),
			},
		},
	})
}

func TestAccAzureRMPowerBIDedicatedCapacity_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_powerbi_dedicated_capacity", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPowerBIDedicatedCapacityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPowerBIDedicatedCapacity_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPowerBIDedicatedCapacityExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "A2"),
					resource.TestCheckResourceAttr(data.ResourceName, "administrators.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.ENV", "Test"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPowerBIDedicatedCapacity_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_powerbi_dedicated_capacity", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPowerBIDedicatedCapacityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPowerBIDedicatedCapacity_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPowerBIDedicatedCapacityExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "A1"),
					resource.TestCheckResourceAttr(data.ResourceName, "administrators.#", "1"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMPowerBIDedicatedCapacity_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPowerBIDedicatedCapacityExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "A2"),
					resource.TestCheckResourceAttr(data.ResourceName, "administrators.#", "2"),
				),
			},
			data.ImportStep(),
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

		client := acceptance.AzureProvider.Meta().(*clients.Client).PowerBIDedicated.CapacityClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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
	client := acceptance.AzureProvider.Meta().(*clients.Client).PowerBIDedicated.CapacityClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testAccAzureRMPowerBIDedicatedCapacity_basic(data acceptance.TestData) string {
	template := testAccAzureRMPowerBIDedicatedCapacity_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_powerbi_dedicated_capacity" "test" {
  name                = "acctestpowerbidedicatedcapacity%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku_name            = "A1"
  administrators      = ["test2@microsoft.onmicrosoft.com"]
}
`, template, data.RandomInteger)
}

func testAccAzureRMPowerBIDedicatedCapacity_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_powerbi_dedicated_capacity" "import" {
  name                = "${azurerm_powerbi_dedicated_capacity.test.name}"
  location            = "${azurerm_powerbi_dedicated_capacity.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, testAccAzureRMPowerBIDedicatedCapacity_basic(data))
}

func testAccAzureRMPowerBIDedicatedCapacity_complete(data acceptance.TestData) string {
	template := testAccAzureRMPowerBIDedicatedCapacity_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_powerbi_dedicated_capacity" "test" {
  name                = "acctestpowerbidedicatedcapacity%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku_name            = "A2"
  administrators      = ["test2@microsoft.onmicrosoft.com", "b1b1f3bc-050d-401c-857a-b872ce501819"]

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMPowerBIDedicatedCapacity_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-powerbidedicated-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
