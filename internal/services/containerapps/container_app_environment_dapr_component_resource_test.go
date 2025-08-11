// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containerapps_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/daprcomponents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ContainerAppEnvironmentDaprComponentResource struct{}

func TestAccContainerAppEnvironmentDaprComponent_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment_dapr_component", "test")
	r := ContainerAppEnvironmentDaprComponentResource{}

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

func TestAccContainerAppEnvironmentDaprComponent_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment_dapr_component", "test")
	r := ContainerAppEnvironmentDaprComponentResource{}

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

func TestAccContainerAppEnvironmentDaprComponent_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment_dapr_component", "test")
	r := ContainerAppEnvironmentDaprComponentResource{}

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

func TestAccContainerAppEnvironmentDaprComponent_completeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment_dapr_component", "test")
	r := ContainerAppEnvironmentDaprComponentResource{}

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

func (r ContainerAppEnvironmentDaprComponentResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := daprcomponents.ParseDaprComponentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.ContainerApps.DaprComponentsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ContainerAppEnvironmentDaprComponentResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_container_app_environment_dapr_component" "test" {
  name                         = "acctest-dapr-%[2]d"
  container_app_environment_id = azurerm_container_app_environment.test.id
  component_type               = "state.azure.blobstorage"
  version                      = "v1"
}

`, r.template(data), data.RandomInteger)
}

func (r ContainerAppEnvironmentDaprComponentResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`

%[1]s

resource "azurerm_container_app_environment_dapr_component" "import" {
  name                         = azurerm_container_app_environment_dapr_component.test.name
  container_app_environment_id = azurerm_container_app_environment_dapr_component.test.container_app_environment_id
  component_type               = azurerm_container_app_environment_dapr_component.test.component_type
  version                      = azurerm_container_app_environment_dapr_component.test.version
}

`, r.basic(data), data.RandomInteger)
}

func (r ContainerAppEnvironmentDaprComponentResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%[3]s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "container-app-storage"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_container_app_environment_dapr_component" "test" {
  name                         = "acctest-dapr-%[2]d"
  container_app_environment_id = azurerm_container_app_environment.test.id
  component_type               = "state.azure.blobstorage"
  version                      = "v1"

  init_timeout  = "10s"
  ignore_errors = true

  secret {
    name  = "secret"
    value = "sauce"
  }

  secret {
    name  = "storage-account-access-key"
    value = azurerm_storage_account.test.primary_access_key
  }

  metadata {
    name        = "storage-account-key"
    secret_name = "storage-account-access-key"
  }

  metadata {
    name  = "storage-container-name"
    value = azurerm_storage_container.test.name
  }

  metadata {
    name  = "SOME_APP_SETTING"
    value = "scwiffy"
  }

  scopes = ["testapp"]
}


`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r ContainerAppEnvironmentDaprComponentResource) completeUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%[3]s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "container-app-storage"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_container_app_environment_dapr_component" "test" {
  name                         = "acctest-dapr-%[2]d"
  container_app_environment_id = azurerm_container_app_environment.test.id
  component_type               = "state.azure.blobstorage"
  version                      = "v2"

  init_timeout  = "5s"
  ignore_errors = false

  secret {
    name  = "storage-account-access-key"
    value = azurerm_storage_account.test.secondary_access_key
  }

  metadata {
    name        = "storage-account-key"
    secret_name = "storage-account-access-key"
  }

  metadata {
    name  = "storage-container-name"
    value = azurerm_storage_container.test.name
  }

  metadata {
    name  = "SOME_APP_SETTING"
    value = "plumbus"
  }

  scopes = ["testapp", "updatedapp"]
}
`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r ContainerAppEnvironmentDaprComponentResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-CAEnv-%[1]d"
  location = "%[2]s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestCAEnv-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_container_app_environment" "test" {
  name                       = "acctest-CAEnv%[1]d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}
