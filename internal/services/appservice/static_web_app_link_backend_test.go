// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/staticsites"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StaticWebAppLinkBackendResource struct{}

func TestStaticWebAppLinkBackendResource_basicApiManagement(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_static_web_app_link_backend", "test")
	r := StaticWebAppLinkBackendResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicApiManagement(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestStaticWebAppLinkBackendResource_basicAppService(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_static_web_app_link_backend", "test")
	r := StaticWebAppLinkBackendResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicAppService(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestStaticWebAppLinkBackendResource_basicContainerApp(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_static_web_app_link_backend", "test")
	r := StaticWebAppLinkBackendResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicContainerApp(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestStaticWebAppLinkBackendResource_multipleExpectError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_static_web_app_link_backend", "test")
	r := StaticWebAppLinkBackendResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicApiManagement(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config:      r.multipleApiManagement(data),
			ExpectError: regexp.MustCompile("already has a backend and cannot have another"),
		},
	})
}

func TestStaticWebAppLinkBackendResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_static_web_app_link_backend", "test")
	r := StaticWebAppLinkBackendResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicApiManagement(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (s StaticWebAppLinkBackendResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := staticsites.ParseLinkedBackendID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.AppService.StaticSitesClient.GetLinkedBackend(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %q: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r StaticWebAppLinkBackendResource) basicApiManagement(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_static_web_app_link_backend" "test" {
  static_web_app_id   = azurerm_static_web_app.test.id
  backend_resource_id = azurerm_api_management.test.id
}
`, r.apiManagement(data))
}

func (r StaticWebAppLinkBackendResource) basicAppService(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_static_web_app_link_backend" "test" {
  static_web_app_id   = azurerm_static_web_app.test.id
  backend_resource_id = azurerm_linux_web_app.test.id
}
`, r.appService(data))
}

func (r StaticWebAppLinkBackendResource) basicContainerApp(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_static_web_app_link_backend" "test" {
  static_web_app_id   = azurerm_static_web_app.test.id
  backend_resource_id = azurerm_container_app.test.id
}
`, r.containerApps(data))
}

func (r StaticWebAppLinkBackendResource) multipleApiManagement(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management" "test" {
  name                = "acctestAPIM-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Consumption_0"
}

resource "azurerm_api_management" "test2" {
  name                = "acctestAPIM-%[2]d-2"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Consumption_0"
}

resource "azurerm_static_web_app_link_backend" "test" {
  static_web_app_id   = azurerm_static_web_app.test.id
  backend_resource_id = azurerm_api_management.test.id
}

resource "azurerm_static_web_app_link_backend" "test2" {
  static_web_app_id   = azurerm_static_web_app.test.id
  backend_resource_id = azurerm_api_management.test2.id
}
`, r.template(data), data.RandomInteger)
}

func (r StaticWebAppLinkBackendResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_static_web_app_link_backend" "import" {
  static_web_app_id   = azurerm_static_web_app_link_backend.test.static_web_app_id
  backend_resource_id = azurerm_static_web_app_link_backend.test.backend_resource_id
}
`, r.basicApiManagement(data))
}

func (r StaticWebAppLinkBackendResource) apiManagement(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

resource "azurerm_api_management" "test" {
  name                = "acctestAPIM-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Consumption_0"
}

`, r.template(data), data.RandomInteger)
}

func (r StaticWebAppLinkBackendResource) appService(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

resource "azurerm_service_plan" "test" {
  name                = "acctestASP-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"
  sku_name            = "F1"
}

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWP-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_service_plan.test.location
  service_plan_id     = azurerm_service_plan.test.id

  site_config {
    always_on = false
  }

  lifecycle {
    ignore_changes = [auth_settings_v2]
  }
}

`, r.template(data), data.RandomInteger)
}

func (r StaticWebAppLinkBackendResource) containerApps(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

resource "azurerm_container_app_environment" "test" {
  name                = "acctestCAEnv-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  workload_profile {
    name                  = "Consumption"
    workload_profile_type = "Consumption"
    maximum_count         = 0
    minimum_count         = 0
  }
}

resource "azurerm_container_app" "test" {
  name                         = "acctest-capp-%[2]d"
  container_app_environment_id = azurerm_container_app_environment.test.id
  resource_group_name          = azurerm_resource_group.test.name
  revision_mode                = "Single"
  workload_profile_name        = "Consumption"

  template {
    container {
      name   = "examplecontainerapp"
      image  = "mcr.microsoft.com/azuredocs/containerapps-helloworld:latest"
      cpu    = 0.25
      memory = "0.5Gi"
    }
  }

  ingress {
    allow_insecure_connections = false
    external_enabled           = false
    target_port                = 443
    transport                  = "auto"

    traffic_weight {
      latest_revision = true
      percentage      = 100
    }
  }
}
	`, r.template(data), data.RandomInteger)
}

func (r StaticWebAppLinkBackendResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`


resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_static_web_app" "test" {
  name                = "acctestSWA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_size            = "Standard"
  sku_tier            = "Standard"
}

`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
