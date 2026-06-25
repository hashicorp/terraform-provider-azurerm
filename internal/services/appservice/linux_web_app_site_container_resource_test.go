// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LinuxWebAppSiteContainerResource struct{}

func TestAccLinuxWebAppSiteContainer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_site_container", "test")
	r := LinuxWebAppSiteContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password_secret"),
	})
}

func TestAccLinuxWebAppSiteContainer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_site_container", "test")
	r := LinuxWebAppSiteContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password_secret"),
	})
}

func TestAccLinuxWebAppSiteContainer_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_site_container", "test")
	r := LinuxWebAppSiteContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password_secret"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password_secret"),
	})
}

func TestAccLinuxWebAppSiteContainer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_site_container", "test")
	r := LinuxWebAppSiteContainerResource{}

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

func (r LinuxWebAppSiteContainerResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := webapps.ParseSitecontainerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.AppService.WebAppsClient.GetSiteContainer(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(true), nil
}

func (r LinuxWebAppSiteContainerResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_service_plan" "test" {
  name                = "acctestSP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Linux"
  sku_name            = "P1v2"
}

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {
    application_stack {
      site_containers_enabled = true
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r LinuxWebAppSiteContainerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_web_app_site_container" "test" {
  name             = "app"
  linux_web_app_id = azurerm_linux_web_app.test.id
  image            = "mcr.microsoft.com/appsvc/sample-hello-world:latest"
  target_port      = 80
  primary          = true
}
`, r.template(data))
}

func (r LinuxWebAppSiteContainerResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_web_app_site_container" "test" {
  name             = "app"
  linux_web_app_id = azurerm_linux_web_app.test.id
  image            = "mcr.microsoft.com/appsvc/sample-hello-world:latest"
  target_port      = 80
  primary          = true
  startup_command  = "sleep 5"

  environment_variable {
    name             = "EXAMPLE"
    app_setting_name = "EXAMPLE_SETTING"
  }

  volume_mount {
    container_mount_path = "/mnt/data"
    volume_sub_path      = "data"
    read_only            = true
  }
}

resource "azurerm_linux_web_app_site_container" "sidecar" {
  name             = "sidecar"
  linux_web_app_id = azurerm_linux_web_app.test.id
  image            = "mcr.microsoft.com/appsvc/sample-hello-world:latest"
  target_port      = 8080
}
`, r.template(data))
}

func (r LinuxWebAppSiteContainerResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_web_app_site_container" "test" {
  name             = "app"
  linux_web_app_id = azurerm_linux_web_app.test.id
  image            = "mcr.microsoft.com/appsvc/sample-hello-world:latest"
  target_port      = 8080
  primary          = true

  environment_variable {
    name             = "EXAMPLE"
    app_setting_name = "EXAMPLE_UPDATED"
  }
}
`, r.template(data))
}

func (r LinuxWebAppSiteContainerResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_web_app_site_container" "import" {
  name             = azurerm_linux_web_app_site_container.test.name
  linux_web_app_id = azurerm_linux_web_app_site_container.test.linux_web_app_id
  image            = azurerm_linux_web_app_site_container.test.image
  target_port      = azurerm_linux_web_app_site_container.test.target_port
  primary          = azurerm_linux_web_app_site_container.test.primary
}
`, r.basic(data))
}
