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

type FWStaticWebAppResource struct{}

func TestAccAzureFWStaticWebApp_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_fw_static_web_app", "test")
	r := FWStaticWebAppResource{}

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

func TestAccAzureFWStaticWebApp_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_fw_static_web_app", "test")
	r := FWStaticWebAppResource{}

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

func TestAccAzureFWStaticWebApp_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_fw_static_web_app", "test")
	r := FWStaticWebAppResource{}

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

func TestAccAzureFWStaticWebApp_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_fw_static_web_app", "test")
	r := FWStaticWebAppResource{}

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

func TestAccAzureFWStaticWebApp_shouldFailBasicAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_fw_static_web_app", "test")
	r := FWStaticWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.shouldFailFreeBasicAuth(data),
			ExpectError: regexp.MustCompile("basic_auth cannot be used with the Free tier of Static Web Apps"),
		},
	})
}

func TestAccAzureFWStaticWebApp_shouldFailIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_fw_static_web_app", "test")
	r := FWStaticWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.shouldFailFreeIdentity(data),
			ExpectError: regexp.MustCompile("identities cannot be used with the Free tier of Static Web Apps"),
		},
	})
}

func (r FWStaticWebAppResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r FWStaticWebAppResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
 features {}
}

resource "azurerm_fw_resource_group" "test" {
 name     = "acctestRG-%[1]d"
 location = "%[2]s"
}

resource "azurerm_fw_static_web_app" "test" {
  name                = "acctestSS-%[1]d"
  location            = azurerm_fw_resource_group.test.location
  resource_group_name = azurerm_fw_resource_group.test.name
  sku_size            = "Standard"
  sku_tier            = "Standard"

}
`, data.RandomInteger, data.Locations.Secondary)
}

func (r FWStaticWebAppResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

resource "azurerm_fw_static_web_app" "import" {
  name                = azurerm_fw_static_web_app.test.name
  location            = azurerm_fw_static_web_app.test.location
  resource_group_name = azurerm_fw_static_web_app.test.resource_group_name
  sku_size            = azurerm_fw_static_web_app.test.sku_size
  sku_tier            = azurerm_fw_static_web_app.test.sku_tier

}
`, r.basic(data))
}

func (r FWStaticWebAppResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_fw_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-%[1]d"
  resource_group_name = azurerm_fw_resource_group.test.name
  location            = azurerm_fw_resource_group.test.location
}

resource "azurerm_fw_static_web_app" "test" {
  name                = "acctestSS-%[1]d"
  location            = azurerm_fw_resource_group.test.location
  resource_group_name = azurerm_fw_resource_group.test.name
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

func (r FWStaticWebAppResource) shouldFailFreeIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_fw_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-%[1]d"
  resource_group_name = azurerm_fw_resource_group.test.name
  location            = azurerm_fw_resource_group.test.location
}

resource "azurerm_fw_static_web_app" "test" {
  name                = "acctestSS-%[1]d"
  location            = azurerm_fw_resource_group.test.location
  resource_group_name = azurerm_fw_resource_group.test.name
  sku_size            = "Free"
  sku_tier            = "Free"

  configuration_file_changes_enabled = false
  preview_environments_enabled       = false

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  app_settings = {
    "foo" = "bar"
  }

  tags = {
    environment = "acceptance"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r FWStaticWebAppResource) shouldFailFreeBasicAuth(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_fw_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_fw_static_web_app" "test" {
  name                = "acctestSS-%[1]d"
  location            = azurerm_fw_resource_group.test.location
  resource_group_name = azurerm_fw_resource_group.test.name
  sku_size            = "Free"
  sku_tier            = "Free"

  configuration_file_changes_enabled = false
  preview_environments_enabled       = false

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
