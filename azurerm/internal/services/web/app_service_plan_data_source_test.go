package web_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

type AppServicePlanDataSource struct{}

func TestAccAppServicePlanDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service_plan", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: AppServicePlanDataSource{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				acceptance.TestCheckResourceAttr(data.ResourceName, "kind", "Windows"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "sku.#", "1"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "sku.0.tier", "Basic"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "sku.0.size", "B1"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
			),
		},
	})
}

func TestAccAppServicePlanDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service_plan", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: AppServicePlanDataSource{}.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				acceptance.TestCheckResourceAttr(data.ResourceName, "kind", "Windows"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "sku.#", "1"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "sku.0.tier", "Standard"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "sku.0.size", "S1"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Test"),
			),
		},
	})
}

func TestAccAppServicePlanDataSource_premiumSKU(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service_plan", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: AppServicePlanDataSource{}.premiumSKU(data),
			Check: acceptance.ComposeTestCheckFunc(
				acceptance.TestCheckResourceAttr(data.ResourceName, "kind", "elastic"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "sku.#", "1"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "sku.0.tier", "ElasticPremium"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "sku.0.size", "EP1"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "maximum_elastic_worker_count", "20"),
			),
		},
	})
}

func TestAccAppServicePlanDataSource_basicWindowsContainer(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service_plan", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: AppServicePlanDataSource{}.basicWindowsContainer(data),
			Check: acceptance.ComposeTestCheckFunc(
				acceptance.TestCheckResourceAttr(data.ResourceName, "kind", "xenon"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "sku.#", "1"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "sku.0.tier", "PremiumV3"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "sku.0.size", "P1v3"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "is_xenon", "true"),
			),
		},
	})
}

func (d AppServicePlanDataSource) basic(data acceptance.TestData) string {
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

data "azurerm_app_service_plan" "test" {
  name                = azurerm_app_service_plan.test.name
  resource_group_name = azurerm_app_service_plan.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (d AppServicePlanDataSource) complete(data acceptance.TestData) string {
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

  tags = {
    environment = "Test"
  }
}

data "azurerm_app_service_plan" "test" {
  name                = azurerm_app_service_plan.test.name
  resource_group_name = azurerm_app_service_plan.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (d AppServicePlanDataSource) premiumSKU(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                         = "acctestASP-%d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  kind                         = "elastic"
  maximum_elastic_worker_count = 20

  sku {
    tier = "ElasticPremium"
    size = "EP1"
  }

  per_site_scaling = true

  tags = {
    environment = "Test"
  }
}

data "azurerm_app_service_plan" "test" {
  name                = azurerm_app_service_plan.test.name
  resource_group_name = azurerm_app_service_plan.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (d AppServicePlanDataSource) basicWindowsContainer(data acceptance.TestData) string {
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
  is_xenon            = true
  kind                = "xenon"

  sku {
    tier = "PremiumV3"
    size = "P1v3"
  }
}

data "azurerm_app_service_plan" "test" {
  name                = azurerm_app_service_plan.test.name
  resource_group_name = azurerm_app_service_plan.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
