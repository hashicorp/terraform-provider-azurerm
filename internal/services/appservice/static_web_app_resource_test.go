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

type StaticWebAppResource struct{}

func TestAccAzureStaticWebApp_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_static_web_app", "test")
	r := StaticWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_host_name").Exists(),
				check.That(data.ResourceName).Key("api_key").Exists(),
				check.That(data.ResourceName).Key("tags.environment").HasValue("acceptance"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureStaticWebApp_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_static_web_app", "test")
	r := StaticWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_host_name").Exists(),
				check.That(data.ResourceName).Key("api_key").Exists(),
			),
		},
		data.ImportStep("basic_auth.0.password"),
	})
}

func TestAccAzureStaticWebApp_completeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_static_web_app", "test")
	r := StaticWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_host_name").Exists(),
				check.That(data.ResourceName).Key("api_key").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_host_name").Exists(),
				check.That(data.ResourceName).Key("api_key").Exists(),
			),
		},
		data.ImportStep("basic_auth.0.password"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_host_name").Exists(),
				check.That(data.ResourceName).Key("api_key").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureStaticWebApp_withSystemAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_static_web_app", "test")
	r := StaticWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withSystemAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureStaticWebApp_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_static_web_app", "test")
	r := StaticWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withSystemAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withSystemAssignedUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withUserAssignedIdentity(data),
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

func TestAccAzureStaticWebApp_withSystemAssignedUserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_static_web_app", "test")
	r := StaticWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withSystemAssignedUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned, UserAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureStaticWebApp_withUserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_static_web_app", "test")
	r := StaticWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("UserAssigned"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureStaticWebApp_basicUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_static_web_app", "test")
	r := StaticWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_host_name").Exists(),
				check.That(data.ResourceName).Key("api_key").Exists(),
				check.That(data.ResourceName).Key("tags.environment").HasValue("acceptance"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_host_name").Exists(),
				check.That(data.ResourceName).Key("api_key").Exists(),
				check.That(data.ResourceName).Key("tags.environment").HasValue("acceptance"),
				check.That(data.ResourceName).Key("tags.updated").HasValue("true"),
			),
		},
	})
}

func TestAccAzureStaticWebApp_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_static_web_app", "test")
	r := StaticWebAppResource{}

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

func TestAccStaticWebApp_appSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_static_web_app", "test")
	r := StaticWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_settings.foo").HasValue("bar"),
			),
		},
		data.ImportStep(),
		{
			Config: r.appSettingsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_settings.foo").HasValue("bar"),
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

func TestAccStaticWebApp_basicAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_static_web_app", "test")
	r := StaticWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withBasicAuth(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("basic_auth.0.password"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureStaticWebApp_basicWithConfigShouldFail(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_static_web_app", "test")
	r := StaticWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.freeWithAuth(data),
			ExpectError: regexp.MustCompile("basic_auth cannot be used with the Free tier of Static Web Apps"),
		},
		{
			Config:      r.freeWithIdentity(data),
			ExpectError: regexp.MustCompile("identities cannot be used with the Free tier of Static Web Apps"),
		},
	})
}

func TestAccAzureStaticWebApp_publicNetworkAccessDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_static_web_app", "test")
	r := StaticWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.publicNetworkAccessDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureStaticWebApp_publicNetworkAccessUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_static_web_app", "test")
	r := StaticWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.publicNetworkAccessDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func (r StaticWebAppResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := staticsites.ParseStaticSiteID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.AppService.StaticSitesClient.GetStaticSite(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %q: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r StaticWebAppResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_static_web_app" "test" {
  name                = "acctestSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "acceptance"
  }
}
`, data.RandomInteger, data.Locations.Primary) // TODO - Put back to primary when support ticket is resolved
}

func (r StaticWebAppResource) withSystemAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_static_web_app" "test" {
  name                = "acctestSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_size            = "Standard"
  sku_tier            = "Standard"

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r StaticWebAppResource) freeWithIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_static_web_app" "test" {
  name                = "acctestSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_size            = "Free"
  sku_tier            = "Free"

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r StaticWebAppResource) withSystemAssignedUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  name = "acctest-%[1]d"
}

resource "azurerm_static_web_app" "test" {
  name                = "acctestSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_size            = "Standard"
  sku_tier            = "Standard"

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, data.RandomInteger, data.Locations.Primary) // TODO - Put back to primary when support ticket is resolved
}

func (r StaticWebAppResource) withUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  name = "acctest-%[1]d"
}

resource "azurerm_static_web_app" "test" {
  name                = "acctestSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_size            = "Standard"
  sku_tier            = "Standard"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, data.RandomInteger, data.Locations.Primary) // TODO - Put back to primary when support ticket is resolved
}

func (r StaticWebAppResource) basicUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_static_web_app" "test" {
  name                = "acctestSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_size            = "Standard"
  sku_tier            = "Standard"

  tags = {
    environment = "acceptance"
    updated     = "true"
  }
}
`, data.RandomInteger, data.Locations.Primary) // TODO - Put back to primary when support ticket is resolved
}

func (r StaticWebAppResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_static_web_app" "import" {
  name                = azurerm_static_web_app.test.name
  location            = azurerm_static_web_app.test.location
  resource_group_name = azurerm_static_web_app.test.resource_group_name
  sku_size            = azurerm_static_web_app.test.sku_size
  sku_tier            = azurerm_static_web_app.test.sku_tier
}
`, template)
}

func (r StaticWebAppResource) appSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_static_web_app" "test" {
  name                = "acctestSS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_size            = "Standard"
  sku_tier            = "Standard"

  app_settings = {
    "foo" = "bar"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r StaticWebAppResource) appSettingsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_static_web_app" "test" {
  name                = "acctestSS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  app_settings = {
    "foo" = "bar"
    "baz" = "foo"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r StaticWebAppResource) withBasicAuth(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_static_web_app" "test" {
  name                = "acctestSS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_size            = "Standard"
  sku_tier            = "Standard"

  basic_auth {
    password     = "Super$3cretPassW0rd"
    environments = "AllEnvironments"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r StaticWebAppResource) freeWithAuth(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_static_web_app" "test" {
  name                = "acctestSS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_size            = "Free"
  sku_tier            = "Free"

  basic_auth {
    password     = "Super$3cretPassW0rd"
    environments = "AllEnvironments"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r StaticWebAppResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_static_web_app" "test" {
  name                = "acctestSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_size            = "Standard"
  sku_tier            = "Standard"

  configuration_file_changes_enabled = false
  preview_environments_enabled       = false

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  app_settings = {
    "foo" = "bar"
  }

  basic_auth {
    password     = "Super$3cretPassW0rd"
    environments = "AllEnvironments"
  }

  tags = {
    environment = "acceptance"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r StaticWebAppResource) publicNetworkAccessDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_static_web_app" "test" {
  name                = "acctestSS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_size            = "Free"
  sku_tier            = "Free"

  public_network_access_enabled = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
