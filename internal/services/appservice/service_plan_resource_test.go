// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ServicePlanResource struct{}

func TestAccServicePlan_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_plan", "test")
	r := ServicePlanResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServicePlan_linuxConsumption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_plan", "test")
	r := ServicePlanResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxFlexConsumption(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServicePlan_linuxFlexConsumption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_plan", "test")
	r := ServicePlanResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxConsumption(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServicePlan_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_plan", "test")
	r := ServicePlanResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccServicePlan_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_plan", "test")
	r := ServicePlanResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServicePlan_completeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_plan", "test")
	r := ServicePlanResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServicePlan_maxElasticWorkerCountForAllSupportedSku(t *testing.T) {
	for _, sku := range []string{"WS1", "EP1"} {
		t.Run(sku, func(t *testing.T) {
			data := acceptance.BuildTestData(t, "azurerm_service_plan", "test")
			r := ServicePlanResource{}

			data.ResourceTest(t, r, []acceptance.TestStep{
				{
					Config: r.maxElasticWorkerCountWithSku(data, 5, sku),
					Check: acceptance.ComposeTestCheckFunc(
						check.That(data.ResourceName).ExistsInAzure(r),
					),
				},
				data.ImportStep(),
				{
					Config: r.maxElasticWorkerCountWithSku(data, 10, sku),
					Check: acceptance.ComposeTestCheckFunc(
						check.That(data.ResourceName).ExistsInAzure(r),
					),
				},
				data.ImportStep(),
			})
		})
	}
}

func TestAccServicePlan_maxElasticWorkerCount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_plan", "test")
	r := ServicePlanResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.maxElasticWorkerCount(data, 5),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.maxElasticWorkerCount(data, 10),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServicePlan_memoryOptimized(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_plan", "test")
	r := ServicePlanResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.memoryOptimized(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

// ASE tests given longer prefix to allow them to be more easily filtered out due to exceptionally long running time
func TestAccServicePlanIsolated_appServiceEnvironmentV2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_plan", "test")
	r := ServicePlanResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.aseV2(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

// ASE tests given longer prefix to allow them to be more easily filtered out due to exceptionally long running time
func TestAccServicePlanIsolated_appServiceEnvironmentV3(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_plan", "test")
	r := ServicePlanResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.aseV3(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ServicePlanResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseAppServicePlanID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.AppService.ServicePlanClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retreiving %s: %v", id, err)
	}
	if response.WasNotFound(resp.HttpResponse) {
		return utils.Bool(false), nil
	}
	return utils.Bool(true), nil
}

// Configs

func (r ServicePlanResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appserviceplan-%[1]d"
  location = "%s"
}

resource "azurerm_service_plan" "test" {
  name                = "acctest-SP-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "B1"
  os_type             = "Windows"

  tags = {
    environment = "AccTest"
    Foo         = "bar"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ServicePlanResource) linuxConsumption(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appserviceplan-%[1]d"
  location = "%s"
}

resource "azurerm_service_plan" "test" {
  name                = "acctest-SP-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Y1"
  os_type             = "Linux"

  tags = {
    environment = "AccTest"
    Foo         = "bar"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ServicePlanResource) linuxFlexConsumption(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appserviceplan-%[1]d"
  location = "%s"
}

resource "azurerm_service_plan" "test" {
  name                = "acctest-SP-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "FC1"
  os_type             = "Linux"

  tags = {
    environment = "AccTest"
    Foo         = "bar"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

// (@jackofallops) - `complete` deliberately omits ASE testing for the moment and will be tested separately later
func (r ServicePlanResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appserviceplan-%[1]d"
  location = "%s"
}

resource "azurerm_service_plan" "test" {
  name                     = "acctest-SP-%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  sku_name                 = "P1v3"
  os_type                  = "Linux"
  per_site_scaling_enabled = true
  worker_count             = 3

  zone_balancing_enabled = true

  tags = {
    environment = "AccTest"
    Foo         = "bar"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ServicePlanResource) completeUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appserviceplan-%[1]d"
  location = "%s"
}

resource "azurerm_service_plan" "test" {
  name                     = "acctest-SP-%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  sku_name                 = "P1v2"
  os_type                  = "Linux"
  per_site_scaling_enabled = true
  worker_count             = 3

  zone_balancing_enabled = true

  tags = {
    Foo = "bar"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ServicePlanResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

resource "azurerm_service_plan" "import" {
  name                = azurerm_service_plan.test.name
  resource_group_name = azurerm_service_plan.test.resource_group_name
  location            = azurerm_service_plan.test.location
  sku_name            = azurerm_service_plan.test.sku_name
  os_type             = azurerm_service_plan.test.os_type
}
`, r.basic(data))
}

func (r ServicePlanResource) maxElasticWorkerCount(data acceptance.TestData, count int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appserviceplan-%[1]d"
  location = "%s"
}

resource "azurerm_service_plan" "test" {
  name                         = "acctest-SP-%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  sku_name                     = "EP1"
  os_type                      = "Linux"
  maximum_elastic_worker_count = %[3]d

  tags = {
    environment = "AccTest"
    Foo         = "bar"
  }
}
`, data.RandomInteger, data.Locations.Primary, count)
}

func (r ServicePlanResource) maxElasticWorkerCountWithSku(data acceptance.TestData, count int, sku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appserviceplan-%[1]d"
  location = "%s"
}

resource "azurerm_service_plan" "test" {
  name                         = "acctest-SP-%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  sku_name                     = "%[3]s"
  os_type                      = "Linux"
  maximum_elastic_worker_count = %[4]d

  tags = {
    environment = "AccTest"
    Foo         = "bar"
  }
}
`, data.RandomInteger, data.Locations.Primary, sku, count)
}

func (r ServicePlanResource) memoryOptimized(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appserviceplan-%[1]d"
  location = "%s"
}

resource "azurerm_service_plan" "test" {
  name                     = "acctest-SP-%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  sku_name                 = "P1mv3"
  os_type                  = "Linux"
  per_site_scaling_enabled = true
  worker_count             = 3

  tags = {
    environment = "AccTest"
    Foo         = "bar"
  }
}
`, data.RandomInteger, data.Locations.Secondary)
}

func (r ServicePlanResource) aseV2(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appserviceplan-%[1]d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "ase" {
  name                 = "asesubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_subnet" "gateway" {
  name                 = "gatewaysubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_app_service_environment" "test" {
  name                = "acctest-ase-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.ase.id
}

resource "azurerm_service_plan" "test" {
  name                = "acctest-SP-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Windows"
  sku_name            = "I1"

  app_service_environment_id = azurerm_app_service_environment.test.id

  tags = {
    environment = "AccTest"
    Foo         = "bar"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ServicePlanResource) aseV3(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appserviceplan-%[1]d"
  location = "%[2]s"
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG2-ase-%[1]d"
  location = "%[2]s"
}


resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-subnet-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
  delegation {
    name = "asedelegation"
    service_delegation {
      name    = "Microsoft.Web/hostingEnvironments"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_app_service_environment_v3" "test" {
  name                = "acctest-ase-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.test.id
}

resource "azurerm_service_plan" "test" {
  name                = "acctest-SP-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Windows"
  sku_name            = "I1v2"

  app_service_environment_id = azurerm_app_service_environment_v3.test.id

  tags = {
    environment = "AccTest"
    Foo         = "bar"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ServicePlanResource) aseV3Linux(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appserviceplan-%[1]d"
  location = "%[2]s"
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG2-ase-%[1]d"
  location = "%[2]s"
}


resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-subnet-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
  delegation {
    name = "asedelegation"
    service_delegation {
      name    = "Microsoft.Web/hostingEnvironments"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_app_service_environment_v3" "test" {
  name                = "acctest-ase-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.test.id
}

resource "azurerm_service_plan" "test" {
  name                = "acctest-SP-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"
  sku_name            = "I1v2"

  app_service_environment_id = azurerm_app_service_environment_v3.test.id

  tags = {
    environment = "AccTest"
    Foo         = "bar"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
