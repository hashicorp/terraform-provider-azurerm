// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package web_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AppServicePlanResource struct{}

func TestAccAppServicePlan_basicWindows(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_plan", "test")
	r := AppServicePlanResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWindows(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("per_site_scaling").HasValue("false"),
				check.That(data.ResourceName).Key("reserved").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppServicePlan_basicLinux(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_plan", "test")
	r := AppServicePlanResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicLinux(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.basicLinuxNew(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("per_site_scaling").HasValue("false"),
				check.That(data.ResourceName).Key("reserved").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppServicePlan_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_plan", "test")
	r := AppServicePlanResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicLinux(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAppServicePlan_standardWindows(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_plan", "test")
	r := AppServicePlanResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standardWindows(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppServicePlan_premiumWindows(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_plan", "test")
	r := AppServicePlanResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.premiumWindows(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppServicePlan_premiumWindowsUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_plan", "test")
	r := AppServicePlanResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.premiumWindows(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("1"),
			),
		},
		{
			Config: r.premiumWindowsUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppServicePlan_completeWindows(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_plan", "test")
	r := AppServicePlanResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeWindows(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("per_site_scaling").HasValue("true"),
				check.That(data.ResourceName).Key("reserved").HasValue("false"),
			),
		},
		{
			Config: r.completeWindowsNew(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("per_site_scaling").HasValue("true"),
				check.That(data.ResourceName).Key("reserved").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppServicePlan_consumptionPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_plan", "test")
	r := AppServicePlanResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.consumptionPlan(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.tier").HasValue("Dynamic"),
				check.That(data.ResourceName).Key("sku.0.size").HasValue("Y1"),
			),
		},
	})
}

func TestAccAppServicePlan_linuxConsumptionPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_plan", "test")
	r := AppServicePlanResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxConsumptionPlan(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppServicePlan_premiumConsumptionPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_plan", "test")
	r := AppServicePlanResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.premiumConsumptionPlan(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.tier").HasValue("ElasticPremium"),
				check.That(data.ResourceName).Key("sku.0.size").HasValue("EP1"),
				check.That(data.ResourceName).Key("maximum_elastic_worker_count").HasValue("20"),
			),
		},
	})
}

func TestAccAppServicePlan_basicWindowsContainer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_plan", "test")
	r := AppServicePlanResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWindowsContainer(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("xenon"),
				check.That(data.ResourceName).Key("is_xenon").HasValue("true"),
				check.That(data.ResourceName).Key("sku.0.tier").HasValue("PremiumV3"),
				check.That(data.ResourceName).Key("sku.0.size").HasValue("P1v3"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppServicePlan_zoneRedundant(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_plan", "test")
	r := AppServicePlanResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.zoneRedundant(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone_redundant").HasValue("true"),
				check.That(data.ResourceName).Key("sku.0.tier").HasValue("PremiumV2"),
				check.That(data.ResourceName).Key("sku.0.size").HasValue("P1v2"),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("3"),
			),
		},
		data.ImportStep(),
	})
}

func (r AppServicePlanResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.AppServicePlanID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Web.AppServicePlansClient.Get(ctx, id.ResourceGroup, id.ServerFarmName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving App Service Plan %q (Resource Group %q): %+v", id.ServerFarmName, id.ResourceGroup, err)
	}

	// The SDK defines 404 as an "ok" status code..
	if utils.ResponseWasNotFound(resp.Response) {
		return utils.Bool(false), nil
	}

	return utils.Bool(true), nil
}

func (r AppServicePlanResource) basicWindows(data acceptance.TestData) string {
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

func (r AppServicePlanResource) basicLinux(data acceptance.TestData) string {
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

func (r AppServicePlanResource) requiresImport(data acceptance.TestData) string {
	template := r.basicLinux(data)
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

func (r AppServicePlanResource) basicLinuxNew(data acceptance.TestData) string {
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

func (r AppServicePlanResource) standardWindows(data acceptance.TestData) string {
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

func (r AppServicePlanResource) premiumWindows(data acceptance.TestData) string {
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

func (r AppServicePlanResource) premiumWindowsUpdated(data acceptance.TestData) string {
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

func (r AppServicePlanResource) completeWindows(data acceptance.TestData) string {
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

func (r AppServicePlanResource) completeWindowsNew(data acceptance.TestData) string {
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

func (r AppServicePlanResource) consumptionPlan(data acceptance.TestData) string {
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

func (r AppServicePlanResource) linuxConsumptionPlan(data acceptance.TestData) string {
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

func (r AppServicePlanResource) premiumConsumptionPlan(data acceptance.TestData) string {
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

func (r AppServicePlanResource) basicWindowsContainer(data acceptance.TestData) string {
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
    tier = "PremiumV3"
    size = "P1v3"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r AppServicePlanResource) zoneRedundant(data acceptance.TestData) string {
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

  zone_redundant = true

  sku {
    tier     = "PremiumV2"
    size     = "P1v2"
    capacity = 3
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
