// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containerapps_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/containerappssessionpools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ContainerAppSessionPoolResource struct{}

func TestAccContainerAppSessionPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_session_pool", "test")
	r := ContainerAppSessionPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("container_type").HasValue("PythonLTS"),
				check.That(data.ResourceName).Key("max_concurrent_sessions").HasValue("5"),
				check.That(data.ResourceName).Key("lifecycle_configuration.0.lifecycle_type").HasValue("Timed"),
				check.That(data.ResourceName).Key("lifecycle_configuration.0.cooldown_period_in_seconds").HasValue("300"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppSessionPool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_session_pool", "test")
	r := ContainerAppSessionPoolResource{}

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

func TestAccContainerAppSessionPool_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_session_pool", "test")
	r := ContainerAppSessionPoolResource{}

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

func TestAccContainerAppSessionPool_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_session_pool", "test")
	r := ContainerAppSessionPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppSessionPool_customContainer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_session_pool", "test")
	r := ContainerAppSessionPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customContainer(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("secret"),
	})
}

func TestAccContainerAppSessionPool_onContainerExit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_session_pool", "test")
	r := ContainerAppSessionPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.onContainerExit(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("lifecycle_configuration.0.lifecycle_type").HasValue("OnContainerExit"),
				check.That(data.ResourceName).Key("lifecycle_configuration.0.max_alive_period_in_seconds").HasValue("600"),
			),
		},
		data.ImportStep(),
	})
}

func (r ContainerAppSessionPoolResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := containerappssessionpools.ParseSessionPoolID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.ContainerApps.SessionPoolClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ContainerAppSessionPoolResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_container_app_session_pool" "test" {
  name                = "acctestcasp%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  lifecycle_configuration {
    cooldown_period_in_seconds = 300
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ContainerAppSessionPoolResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_container_app_session_pool" "import" {
  name                = azurerm_container_app_session_pool.test.name
  resource_group_name = azurerm_container_app_session_pool.test.resource_group_name
  location            = azurerm_container_app_session_pool.test.location

  lifecycle_configuration {
    lifecycle_type             = azurerm_container_app_session_pool.test.lifecycle_configuration[0].lifecycle_type
    cooldown_period_in_seconds = azurerm_container_app_session_pool.test.lifecycle_configuration[0].cooldown_period_in_seconds
  }
}
`, r.basic(data))
}

func (r ContainerAppSessionPoolResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_container_app_session_pool" "test" {
  name                    = "acctestcasp%[2]d"
  resource_group_name     = azurerm_resource_group.test.name
  location                = azurerm_resource_group.test.location
  container_type          = "PythonLTS"
  max_concurrent_sessions = 10

  lifecycle_configuration {
    lifecycle_type             = "Timed"
    cooldown_period_in_seconds = 600
  }

  network_egress_enabled = true

  tags = {
    environment = "acctest"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ContainerAppSessionPoolResource) customContainer(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_container_app_session_pool" "test" {
  name                         = "acctestcasp%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  container_app_environment_id = azurerm_container_app_environment.test.id
  container_type               = "CustomContainer"
  max_concurrent_sessions      = 10

  lifecycle_configuration {
    lifecycle_type             = "Timed"
    cooldown_period_in_seconds = 600
  }

  ready_session_instances = 2

  identity {
    type = "SystemAssigned"
  }

  session_managed_identities = ["System"]

  custom_container_template {
    ingress_target_port = 80

    registry {
      server               = azurerm_container_registry.test.login_server
      username             = azurerm_container_registry.test.admin_username
      password_secret_name = "registry-password"
    }

    container {
      name   = "acctestcontainer"
      image  = "mcr.microsoft.com/k8se/quickstart:latest"
      cpu    = 0.25
      memory = "0.5Gi"

      args    = ["-c", "while true; do echo hello; sleep 10;done"]
      command = ["/bin/sh"]

      env {
        name  = "ACC_TEST"
        value = "acctestvalue"
      }

      env {
        name        = "ACC_TEST_SECRET"
        secret_name = "acctestsecret"
      }
    }
  }

  secret {
    name  = "acctestsecret"
    value = "acctestsecretvalue"
  }

  secret {
    name  = "registry-password"
    value = azurerm_container_registry.test.admin_password
  }
}
`, r.templateWithEnvironment(data), data.RandomInteger)
}

func (r ContainerAppSessionPoolResource) onContainerExit(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_container_app_session_pool" "test" {
  name                         = "acctestcasp%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  container_app_environment_id = azurerm_container_app_environment.test.id
  container_type               = "CustomContainer"
  max_concurrent_sessions      = 5

  lifecycle_configuration {
    lifecycle_type              = "OnContainerExit"
    max_alive_period_in_seconds = 600
  }

  ready_session_instances = 1

  custom_container_template {
    ingress_target_port = 80

    container {
      name   = "acctestcontainer"
      image  = "mcr.microsoft.com/k8se/quickstart:latest"
      cpu    = 0.25
      memory = "0.5Gi"
    }
  }
}
`, r.templateWithEnvironment(data), data.RandomInteger)
}

func (r ContainerAppSessionPoolResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-CASP-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ContainerAppSessionPoolResource) templateWithEnvironment(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_container_app_environment" "test" {
  name                = "acctest-CAEnv%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  # Session Pools are only supported on Workload profile environments.
  workload_profile {
    name                  = "Consumption"
    workload_profile_type = "Consumption"
  }
}

resource "azurerm_container_registry" "test" {
  name                = "acctestcr%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Basic"
  admin_enabled       = true
}
`, r.template(data), data.RandomInteger)
}
