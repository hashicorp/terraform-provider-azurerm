// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containerapps_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ContainerAppEnvironmentResource struct{}

func TestAccContainerAppEnvironment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment", "test")
	r := ContainerAppEnvironmentResource{}

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

func TestAccContainerAppEnvironment_consumptionWorkloadProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment", "test")
	r := ContainerAppEnvironmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.consumptionWorkloadProfile(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("log_analytics_workspace_id", "workload_profile"),
	})
}

func TestAccContainerAppEnvironment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment", "test")
	r := ContainerAppEnvironmentResource{}

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

func TestAccContainerAppEnvironment_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment", "test")
	r := ContainerAppEnvironmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppEnvironment_logsDestinationWithoutWorkspaceShouldFail(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment", "test")
	r := ContainerAppEnvironmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.logsDestinationWithoutWorkspaceShouldFail(data),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile("`log_analytics_workspace_id` must be set when `logs_destination` is set to `log-analytics`"),
		},
	})
}

func TestAccContainerAppEnvironment_logsAzureMonitorWithWorkspaceShouldFail(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment", "test")
	r := ContainerAppEnvironmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.logsAzureMonitorWithWorkspaceShouldFail(data),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile("`log_analytics_workspace_id` can only be set when `logs_destination` is set to `log-analytics`"),
		},
	})
}

func TestAccContainerAppEnvironment_updateLogsDestination(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment", "test")
	r := ContainerAppEnvironmentResource{}

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
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeNoLoggingDestination(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppEnvironment_updateWorkloadProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment", "test")
	r := ContainerAppEnvironmentResource{}

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
		{
			Config: r.completeMultipleWorkloadProfiles(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppEnvironment_daprApplicationInsightsConnectionString(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment", "test")
	r := ContainerAppEnvironmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.daprApplicationInsightsConnectionString(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("dapr_application_insights_connection_string"),
	})
}

func TestAccContainerAppEnvironment_zoneRedundant(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment", "test")
	r := ContainerAppEnvironmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeZoneRedundant(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppEnvironment_infraResourceGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment", "test")
	r := ContainerAppEnvironmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.infraResourceGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppEnvironment_infraResourceGroupWithoutName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment", "test")
	r := ContainerAppEnvironmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.infraResourceGroupWithoutName(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ContainerAppEnvironmentResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := managedenvironments.ParseManagedEnvironmentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.ContainerApps.ManagedEnvironmentClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ContainerAppEnvironmentResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_container_app_environment" "test" {
  name                = "acctest-CAEnv%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, r.template(data), data.RandomInteger)
}

func (r ContainerAppEnvironmentResource) basicNoProvider(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_container_app_environment" "test" {
  name                = "acctest-CAEnv%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, r.template(data), data.RandomInteger)
}

func (r ContainerAppEnvironmentResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`


%[1]s

resource "azurerm_container_app_environment" "import" {
  name                = azurerm_container_app_environment.test.name
  resource_group_name = azurerm_container_app_environment.test.resource_group_name
  location            = azurerm_container_app_environment.test.location
}
`, r.basic(data), data.RandomInteger)
}

func (r ContainerAppEnvironmentResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_container_app_environment" "test" {
  name                       = "acctest-CAEnv%[2]d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  logs_destination           = "log-analytics"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
  infrastructure_subnet_id   = azurerm_subnet.control.id

  internal_load_balancer_enabled = true
  zone_redundancy_enabled        = true
  mutual_tls_enabled             = true

  workload_profile {
    maximum_count         = 3
    minimum_count         = 0
    name                  = "D4-01"
    workload_profile_type = "D4"
  }

  tags = {
    Foo    = "Bar"
    secret = "sauce"
  }
}

resource "azurerm_monitor_diagnostic_setting" "test" {
  name                       = "diagnostics"
  target_resource_id         = azurerm_container_app_environment.test.id
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id

  enabled_log {
    category_group = "allLogs"
  }

  metric {
    category = "AllMetrics"
    enabled  = true
  }
}
`, r.templateVNet(data), data.RandomInteger)
}

func (r ContainerAppEnvironmentResource) updateTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_container_app_environment" "test" {
  name                       = "acctest-CAEnv%[2]d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  logs_destination           = "log-analytics"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
  infrastructure_subnet_id   = azurerm_subnet.control.id

  internal_load_balancer_enabled = true
  zone_redundancy_enabled        = true
  mutual_tls_enabled             = true

  workload_profile {
    maximum_count         = 3
    minimum_count         = 0
    name                  = "D4-01"
    workload_profile_type = "D4"
  }

  tags = {
    Foo    = "Bar"
    secret = "sauce"
    Hello  = "World"
  }
}

resource "azurerm_monitor_diagnostic_setting" "test" {
  name                       = "diagnostics"
  target_resource_id         = azurerm_container_app_environment.test.id
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id

  enabled_log {
    category_group = "allLogs"
  }

  metric {
    category = "AllMetrics"
    enabled  = true
  }
}
`, r.templateVNet(data), data.RandomInteger)
}

func (r ContainerAppEnvironmentResource) completeNoLoggingDestination(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_container_app_environment" "test" {
  name                     = "acctest-CAEnv%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  infrastructure_subnet_id = azurerm_subnet.control.id

  internal_load_balancer_enabled = true
  zone_redundancy_enabled        = true
  mutual_tls_enabled             = true

  workload_profile {
    maximum_count         = 3
    minimum_count         = 0
    name                  = "D4-01"
    workload_profile_type = "D4"
  }

  tags = {
    Foo    = "Bar"
    secret = "sauce"
  }
}

resource "azurerm_monitor_diagnostic_setting" "test" {
  name                       = "diagnostics"
  target_resource_id         = azurerm_container_app_environment.test.id
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id

  enabled_log {
    category_group = "allLogs"
  }

  metric {
    category = "AllMetrics"
    enabled  = true
  }
}
`, r.templateVNet(data), data.RandomInteger)
}

func (r ContainerAppEnvironmentResource) logsDestinationWithoutWorkspaceShouldFail(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_container_app_environment" "test" {
  name                     = "acctest-CAEnv%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  logs_destination         = "log-analytics"
  infrastructure_subnet_id = azurerm_subnet.control.id

  internal_load_balancer_enabled = true
  zone_redundancy_enabled        = true
  mutual_tls_enabled             = true

  workload_profile {
    maximum_count         = 3
    minimum_count         = 0
    name                  = "D4-01"
    workload_profile_type = "D4"
  }

  tags = {
    Foo    = "Bar"
    secret = "sauce"
  }
}

resource "azurerm_monitor_diagnostic_setting" "test" {
  name                       = "diagnostics"
  target_resource_id         = azurerm_container_app_environment.test.id
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id

  enabled_log {
    category_group = "allLogs"
  }

  metric {
    category = "AllMetrics"
    enabled  = true
  }
}
`, r.templateVNet(data), data.RandomInteger)
}

func (r ContainerAppEnvironmentResource) logsAzureMonitorWithWorkspaceShouldFail(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_container_app_environment" "test" {
  name                       = "acctest-CAEnv%[2]d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  logs_destination           = "azure-monitor"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
  infrastructure_subnet_id   = azurerm_subnet.control.id

  internal_load_balancer_enabled = true
  zone_redundancy_enabled        = true
  mutual_tls_enabled             = true

  workload_profile {
    maximum_count         = 3
    minimum_count         = 0
    name                  = "D4-01"
    workload_profile_type = "D4"
  }

  tags = {
    Foo    = "Bar"
    secret = "sauce"
  }
}

resource "azurerm_monitor_diagnostic_setting" "test" {
  name                       = "diagnostics"
  target_resource_id         = azurerm_container_app_environment.test.id
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id

  enabled_log {
    category_group = "allLogs"
  }

  metric {
    category = "AllMetrics"
    enabled  = true
  }
}
`, r.templateVNet(data), data.RandomInteger)
}

func (r ContainerAppEnvironmentResource) completeWithoutWorkloadProfile(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_container_app_environment" "test" {
  name                       = "acctest-CAEnv%[2]d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
  infrastructure_subnet_id   = azurerm_subnet.control.id

  internal_load_balancer_enabled = true
  zone_redundancy_enabled        = true

  tags = {
    Foo    = "Bar"
    secret = "sauce"
  }
}
`, r.templateVnetSubnetNotDelegated(data), data.RandomInteger)
}

func (r ContainerAppEnvironmentResource) consumptionWorkloadProfile(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_container_app_environment" "test" {
  name                = "acctest-CAEnv%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  workload_profile {
    name                  = "Consumption"
    workload_profile_type = "Consumption"
  }
}
`, r.templateVNet(data), data.RandomInteger)
}

func (r ContainerAppEnvironmentResource) completeUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_log_analytics_workspace" "second" {
  name                = "acctestLAW-second-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_container_app_environment" "test" {
  name                     = "acctest-CAEnv%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  logs_destination         = "azure-monitor"
  infrastructure_subnet_id = azurerm_subnet.control.id

  internal_load_balancer_enabled = true
  zone_redundancy_enabled        = true

  workload_profile {
    maximum_count         = 2
    minimum_count         = 0
    name                  = "E4-01"
    workload_profile_type = "E4"
  }

  tags = {
    Foo    = "Bar"
    secret = "sauce"
  }
}
`, r.templateVNet(data), data.RandomInteger)
}

func (r ContainerAppEnvironmentResource) completeMultipleWorkloadProfiles(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_container_app_environment" "test" {
  name                     = "acctest-CAEnv%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  infrastructure_subnet_id = azurerm_subnet.control.id

  internal_load_balancer_enabled = true
  zone_redundancy_enabled        = true

  workload_profile {
    maximum_count         = 2
    minimum_count         = 0
    name                  = "E4-01"
    workload_profile_type = "E4"
  }

  workload_profile {
    maximum_count         = 2
    minimum_count         = 0
    name                  = "D4-02"
    workload_profile_type = "E4"
  }

  workload_profile {
    maximum_count         = 2
    minimum_count         = 0
    name                  = "D4-01"
    workload_profile_type = "D4"
  }

  tags = {
    Foo    = "Bar"
    secret = "sauce"
  }
}
`, r.templateVNet(data), data.RandomInteger)
}

func (r ContainerAppEnvironmentResource) completeZoneRedundant(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_container_app_environment" "test" {
  name                           = "acctest-CAEnv%[2]d"
  resource_group_name            = azurerm_resource_group.test.name
  location                       = azurerm_resource_group.test.location
  infrastructure_subnet_id       = azurerm_subnet.control.id
  zone_redundancy_enabled        = true
  internal_load_balancer_enabled = true

  tags = {
    Foo    = "Bar"
    secret = "sauce"
  }
}
`, r.templateVNet(data), data.RandomInteger)
}

func (r ContainerAppEnvironmentResource) daprApplicationInsightsConnectionString(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_container_app_environment" "test" {
  name                = "acctest-CAEnv%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  dapr_application_insights_connection_string = azurerm_application_insights.test.connection_string
}
`, r.template(data), data.RandomInteger)
}

func (r ContainerAppEnvironmentResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-CAE-%[1]d"
  location = "%[2]s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ContainerAppEnvironmentResource) templateVNet(data acceptance.TestData) string {
	return fmt.Sprintf(`


%[1]s

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[2]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "control" {
  name                 = "control-plane"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.0.0/23"]

  delegation {
    name = "acctestdelegation%[2]d"
    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.App/environments"
    }
  }
}


`, r.template(data), data.RandomInteger)
}

func (r ContainerAppEnvironmentResource) templateVnetSubnetNotDelegated(data acceptance.TestData) string {
	return fmt.Sprintf(`

%[1]s

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[2]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "control" {
  name                 = "control-plane"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.0.0/23"]
}
`, r.template(data), data.RandomInteger)
}

func (r ContainerAppEnvironmentResource) infraResourceGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_container_app_environment" "test" {
  name                     = "acctest-CAEnv%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  infrastructure_subnet_id = azurerm_subnet.control.id

  infrastructure_resource_group_name = "rg-acctest-CAEnv%[2]d"

  workload_profile {
    maximum_count         = 2
    minimum_count         = 0
    name                  = "D4-01"
    workload_profile_type = "D4"
  }

  zone_redundancy_enabled = true

  tags = {
    Foo    = "Bar"
    secret = "sauce"
  }
}
`, r.templateVNet(data), data.RandomInteger)
}

func (r ContainerAppEnvironmentResource) infraResourceGroupWithoutName(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_container_app_environment" "test" {
  name                     = "acctest-CAEnv%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  infrastructure_subnet_id = azurerm_subnet.control.id

  workload_profile {
    maximum_count         = 2
    minimum_count         = 0
    name                  = "D4-01"
    workload_profile_type = "D4"
  }

  zone_redundancy_enabled = true

  tags = {
    Foo    = "Bar"
    secret = "sauce"
  }
}
`, r.templateVNet(data), data.RandomInteger)
}
