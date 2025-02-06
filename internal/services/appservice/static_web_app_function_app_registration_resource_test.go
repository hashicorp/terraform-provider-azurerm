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

type StaticWebAppFunctionAppRegistrationResource struct{}

func TestStaticWebAppFunctionAppRegistrationResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_static_web_app_function_app_registration", "test")
	r := StaticWebAppFunctionAppRegistrationResource{}

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

func TestStaticWebAppFunctionAppRegistrationResource_multipleExpectError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_static_web_app_function_app_registration", "test")
	r := StaticWebAppFunctionAppRegistrationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config:      r.multipleFunctions(data),
			ExpectError: regexp.MustCompile("already has a backend and cannot have another"),
		},
	})
}

func TestStaticWebAppFunctionAppRegistrationResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_static_web_app_function_app_registration", "test")
	r := StaticWebAppFunctionAppRegistrationResource{}

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

func (s StaticWebAppFunctionAppRegistrationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := staticsites.ParseUserProvidedFunctionAppID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.AppService.StaticSitesClient.GetUserProvidedFunctionAppForStaticSite(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %q: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r StaticWebAppFunctionAppRegistrationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_static_web_app_function_app_registration" "test" {
  static_web_app_id = azurerm_static_web_app.test.id
  function_app_id   = azurerm_linux_function_app.test.id
}
`, r.template(data))
}

func (r StaticWebAppFunctionAppRegistrationResource) multipleFunctions(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_function_app" "test2" {
  name                = "acctest-LFA-%[2]d-2"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {}

  lifecycle {
    ignore_changes = [auth_settings_v2]
  }
}

resource "azurerm_static_web_app_function_app_registration" "test" {
  static_web_app_id = azurerm_static_web_app.test.id
  function_app_id   = azurerm_linux_function_app.test.id
}

resource "azurerm_static_web_app_function_app_registration" "test2" {
  static_web_app_id = azurerm_static_web_app.test.id
  function_app_id   = azurerm_linux_function_app.test2.id
}
`, r.template(data), data.RandomInteger)
}

func (r StaticWebAppFunctionAppRegistrationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_static_web_app_function_app_registration" "import" {
  static_web_app_id = azurerm_static_web_app_function_app_registration.test.static_web_app_id
  function_app_id   = azurerm_static_web_app_function_app_registration.test.function_app_id
}
`, r.basic(data))
}

func (r StaticWebAppFunctionAppRegistrationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`


resource "azurerm_resource_group" "test" {
  name     = "acctestRG-WFA-%[1]d"
  location = "%[2]s"
}

resource "azurerm_static_web_app" "test" {
  name                = "acctestSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_size            = "Standard"
  sku_tier            = "Standard"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Linux"
  sku_name            = "S1"
}

resource "azurerm_linux_function_app" "test" {
  name                = "acctest-LFA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {}

  lifecycle {
    ignore_changes = [auth_settings_v2]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
