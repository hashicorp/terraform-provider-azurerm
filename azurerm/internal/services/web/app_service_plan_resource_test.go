package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAppServicePlan_basicWindows(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_plan", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServicePlan_basicWindows(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "per_site_scaling", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "reserved", "false"),
				),
			},

			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServicePlan_basicLinux(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_plan", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServicePlan_basicLinux(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMAppServicePlan_basicLinuxNew(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "per_site_scaling", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "reserved", "true"),
				),
			},

			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServicePlan_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_plan", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServicePlan_basicLinux(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMAppServicePlan_requiresImport),
		},
	})
}

func TestAccAzureRMAppServicePlan_standardWindows(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_plan", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServicePlan_standardWindows(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(data.ResourceName),
				),
			},

			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServicePlan_premiumWindows(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_plan", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServicePlan_premiumWindows(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(data.ResourceName),
				),
			},

			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServicePlan_premiumWindowsUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_plan", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServicePlan_premiumWindows(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "1"),
				),
			},
			{
				Config: testAccAzureRMAppServicePlan_premiumWindowsUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "2"),
				),
			},

			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServicePlan_completeWindows(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_plan", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServicePlan_completeWindows(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "per_site_scaling", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "reserved", "false"),
				),
			},
			{
				Config: testAccAzureRMAppServicePlan_completeWindowsNew(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "per_site_scaling", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "reserved", "false"),
				),
			},

			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServicePlan_consumptionPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_plan", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServicePlan_consumptionPlan(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.tier", "Dynamic"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.size", "Y1"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServicePlan_linuxConsumptionPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_plan", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServicePlan_linuxConsumptionPlan(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServicePlan_premiumConsumptionPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_plan", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServicePlan_premiumConsumptionPlan(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.tier", "ElasticPremium"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.size", "EP1"),
					resource.TestCheckResourceAttr(data.ResourceName, "maximum_elastic_worker_count", "20"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServicePlan_basicWindowsContainer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_plan", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServicePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServicePlan_basicWindowsContainer(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServicePlanExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "xenon"),
					resource.TestCheckResourceAttr(data.ResourceName, "is_xenon", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.tier", "PremiumContainer"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.size", "PC2"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMAppServicePlanDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Web.AppServicePlansClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_app_service_plan" {
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

		return nil
	}

	return nil
}

func testCheckAzureRMAppServicePlanExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Web.AppServicePlansClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		appServicePlanName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for App Service Plan: %s", appServicePlanName)
		}

		resp, err := conn.Get(ctx, resourceGroup, appServicePlanName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: App Service Plan %q (resource group: %q) does not exist", appServicePlanName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on appServicePlansClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMAppServicePlan_basicWindows(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Basic"
    size = "B1"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAppServicePlan_basicLinux(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "Linux"

  sku {
    tier = "Basic"
    size = "B1"
  }

  reserved = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAppServicePlan_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMAppServicePlan_basicLinux(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_plan" "import" {
  name                = azurerm_app_service_plan.test.name
  location            = azurerm_app_service_plan.test.location
  resource_group_name = azurerm_app_service_plan.test.resource_group_name
  kind                = azurerm_app_service_plan.test.kind

  sku {
    tier = "Basic"
    size = "B1"
  }

  reserved = true
}
`, template)
}

func testAccAzureRMAppServicePlan_basicLinuxNew(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "Linux"

  sku {
    tier = "Basic"
    size = "B1"
  }

  reserved = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAppServicePlan_standardWindows(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAppServicePlan_premiumWindows(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Premium"
    size = "P1"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAppServicePlan_premiumWindowsUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier     = "Premium"
    size     = "P1"
    capacity = 2
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAppServicePlan_completeWindows(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "Windows"

  sku {
    tier = "Standard"
    size = "S1"
  }

  per_site_scaling = true
  reserved         = false

  tags = {
    environment = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAppServicePlan_completeWindowsNew(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "Windows"

  sku {
    tier = "Standard"
    size = "S1"
  }

  per_site_scaling = true
  reserved         = false

  tags = {
    environment = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAppServicePlan_consumptionPlan(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  kind                = "FunctionApp"

  sku {
    tier = "Dynamic"
    size = "Y1"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAppServicePlan_linuxConsumptionPlan(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  kind                = "FunctionApp"
  reserved            = true

  sku {
    tier = "Dynamic"
    size = "Y1"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAppServicePlan_premiumConsumptionPlan(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  kind                = "elastic"

  maximum_elastic_worker_count = 20

  sku {
    tier = "ElasticPremium"
    size = "EP1"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAppServicePlan_basicWindowsContainer(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "xenon"
  is_xenon            = true

  sku {
    tier = "PremiumContainer"
    size = "PC2"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
