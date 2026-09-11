// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccLinuxWebAppSiteContainer_list_basic(t *testing.T) {
	r := LinuxWebAppSiteContainerResource{}
	listResourceAddress := "azurerm_linux_web_app_site_container.list"

	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_site_container", "test")

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				Query:  true,
				Config: r.basicQueryByLinuxWebAppId(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 3),
				},
			},
		},
	})
}

func (r LinuxWebAppSiteContainerResource) basicList(data acceptance.TestData) string {
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

resource "azurerm_linux_web_app_site_container" "test" {
  count = 3

  name             = "container${count.index}"
  linux_web_app_id = azurerm_linux_web_app.test.id
  image            = "mcr.microsoft.com/appsvc/sample-hello-world:latest"
  target_port      = 80
  primary          = count.index == 0
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r LinuxWebAppSiteContainerResource) basicQueryByLinuxWebAppId(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_linux_web_app_site_container" "list" {
  provider = azurerm
  config {
    linux_web_app_id = "/subscriptions/%[1]s/resourceGroups/acctestRG-%[2]d/providers/Microsoft.Web/sites/acctestWA-%[2]d"
  }
}
`, data.Subscriptions.Primary, data.RandomInteger)
}
