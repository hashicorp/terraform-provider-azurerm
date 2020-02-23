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

func TestAccAzureRMAppServiceEnvironment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_environment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceEnvironment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceEnvironmentExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "pricing_tier", "I1"),
					resource.TestCheckResourceAttr(data.ResourceName, "front_end_scale_factor", "15"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServiceEnvironment_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_app_service_environment", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceEnvironment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceEnvironmentExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMAppServiceEnvironment_requiresImport),
		},
	})
}

func TestAccAzureRMAppServiceEnvironment_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_environment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceEnvironment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceEnvironmentExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "pricing_tier", "I1"),
					resource.TestCheckResourceAttr(data.ResourceName, "front_end_scale_factor", "15"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAppServiceEnvironment_tierAndScaleFactor(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "pricing_tier", "I2"),
					resource.TestCheckResourceAttr(data.ResourceName, "front_end_scale_factor", "10"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServiceEnvironment_tierAndScaleFactor(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_environment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceEnvironment_tierAndScaleFactor(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceEnvironmentExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "pricing_tier", "I2"),
					resource.TestCheckResourceAttr(data.ResourceName, "front_end_scale_factor", "10"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServiceEnvironment_withAppServicePlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_environment", "test")
	aspData := acceptance.BuildTestData(t, "azurerm_app_service_plan", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceEnvironment_withAppServicePlan(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceEnvironmentExists(data.ResourceName),
					resource.TestCheckResourceAttrPair(data.ResourceName, "id", aspData.ResourceName, "app_service_environment_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMAppServiceEnvironmentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Web.AppServiceEnvironmentsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		appServiceEnvironmentName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for App Service Environment: %s", appServiceEnvironmentName)
		}

		resp, err := client.Get(ctx, resourceGroup, appServiceEnvironmentName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: App Service Environment %q (resource group %q) does not exist", appServiceEnvironmentName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on appServiceEnvironmentClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMAppServiceEnvironmentDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Web.AppServiceEnvironmentsClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_app_service_environment" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return nil
	}

	return nil
}

func testAccAzureRMAppServiceEnvironment_basic(data acceptance.TestData) string {
	template := testAccAzureRMAppServiceEnvironment_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_environment" "test" {
  name      = "acctest-ase-%d"
  subnet_id = azurerm_subnet.ase.id
}
`, template, data.RandomInteger)
}

func testAccAzureRMAppServiceEnvironment_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMAppServiceEnvironment_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_environment" "import" {
  name      = azurerm_app_service_environment.test.name
  subnet_id = azurerm_app_service_environment.test.subnet_id
}
`, template)
}

func testAccAzureRMAppServiceEnvironment_tierAndScaleFactor(data acceptance.TestData) string {
	template := testAccAzureRMAppServiceEnvironment_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_environment" "test" {
  name                   = "acctest-ase-%d"
  subnet_id              = azurerm_subnet.ase.id
  pricing_tier           = "I2"
  front_end_scale_factor = 10
}
`, template, data.RandomInteger)
}

func testAccAzureRMAppServiceEnvironment_withAppServicePlan(data acceptance.TestData) string {
	template := testAccAzureRMAppServiceEnvironment_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_plan" "test" {
  name                       = "acctest-ASP-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_environment_id = azurerm_app_service_environment.test.id

  sku {
    tier     = "Isolated"
    size     = "I1"
    capacity = 1
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMAppServiceEnvironment_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "ase" {
  name                 = "asesubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_subnet" "gateway" {
  name                 = "gatewaysubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
