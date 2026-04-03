// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containerapps_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/httprouteconfig"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ContainerAppEnvironmentHttpRouteConfigResource struct{}

func TestAccContainerAppEnvironmentHttpRouteConfig_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment_http_route_config", "test")
	r := ContainerAppEnvironmentHttpRouteConfigResource{}

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

func TestAccContainerAppEnvironmentHttpRouteConfig_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment_http_route_config", "test")
	r := ContainerAppEnvironmentHttpRouteConfigResource{}

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

func TestAccContainerAppEnvironmentHttpRouteConfig_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment_http_route_config", "test")
	r := ContainerAppEnvironmentHttpRouteConfigResource{}

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

func TestAccContainerAppEnvironmentHttpRouteConfig_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment_http_route_config", "test")
	r := ContainerAppEnvironmentHttpRouteConfigResource{}

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
	})
}

func (r ContainerAppEnvironmentHttpRouteConfigResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := httprouteconfig.ParseHTTPRouteConfigID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.ContainerApps.HttpRouteConfigClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ContainerAppEnvironmentHttpRouteConfigResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_container_app_environment_http_route_config" "test" {
  name                         = "testroute%[2]d"
  container_app_environment_id = azurerm_container_app_environment.test.id

  rules {
    targets {
      container_app = azurerm_container_app.test.name
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ContainerAppEnvironmentHttpRouteConfigResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_container_app_environment_http_route_config" "import" {
  name                         = azurerm_container_app_environment_http_route_config.test.name
  container_app_environment_id = azurerm_container_app_environment_http_route_config.test.container_app_environment_id

  rules {
    targets {
      container_app = azurerm_container_app.test.name
    }
  }
}
`, r.basic(data))
}

func (r ContainerAppEnvironmentHttpRouteConfigResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_container_app_environment_http_route_config" "test" {
  name                         = "testroute%[2]d"
  container_app_environment_id = azurerm_container_app_environment.test.id

  rules {
    description = "Main routing rule"

    routes {
      match {
        prefix         = "/api"
        case_sensitive = true
      }

      action {
        prefix_rewrite = "/v1"
      }
    }

    targets {
      container_app = azurerm_container_app.test.name
    }
  }

  rules {
    description = "Fallback rule"

    routes {
      match {
        prefix = "/"
      }
    }

    targets {
      container_app = azurerm_container_app.test.name
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ContainerAppEnvironmentHttpRouteConfigResource) template(data acceptance.TestData) string {
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

resource "azurerm_container_app_environment" "test" {
  name                       = "acctest-CAEnv%[1]d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
}

resource "azurerm_container_app" "test" {
  name                         = "testcapp%[1]d"
  container_app_environment_id = azurerm_container_app_environment.test.id
  resource_group_name          = azurerm_resource_group.test.name
  revision_mode                = "Single"

  template {
    container {
      name   = "testcontainer"
      image  = "mcr.microsoft.com/k8se/quickstart:latest"
      cpu    = 0.25
      memory = "0.5Gi"
    }
  }

  ingress {
    external_enabled = true
    target_port      = 80
    transport        = "http"

    traffic_weight {
      latest_revision = true
      percentage      = 100
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
