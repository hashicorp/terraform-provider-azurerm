// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package applicationinsights_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	components "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AppInsightsResource struct{}

func TestAccApplicationInsights_basicWeb(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights", "test")
	r := AppInsightsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "web"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("application_type").HasValue("web"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationInsights_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights", "test")
	r := AppInsightsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "web"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("application_type").HasValue("web"),
			),
		},
		{
			Config:      r.requiresImport(data, "web"),
			ExpectError: acceptance.RequiresImportError("azurerm_application_insights"),
		},
	})
}

func TestAccApplicationInsights_basicJava(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights", "test")
	r := AppInsightsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "java"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("application_type").HasValue("java"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationInsights_basicMobileCenter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights", "test")
	r := AppInsightsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "MobileCenter"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("application_type").HasValue("MobileCenter"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationInsights_basicOther(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights", "test")
	r := AppInsightsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "other"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("application_type").HasValue("other"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationInsights_basicPhone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights", "test")
	r := AppInsightsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "phone"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("application_type").HasValue("phone"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationInsights_basicStore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights", "test")
	r := AppInsightsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "store"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("application_type").HasValue("store"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationInsights_basiciOS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights", "test")
	r := AppInsightsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "ios"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("application_type").HasValue("ios"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationInsights_basicWorkspaceMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights", "test")
	r := AppInsightsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWorkspaceMode(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWorkspaceModeUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t AppInsightsResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := components.ParseComponentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.AppInsights.ComponentsClient.ComponentsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func TestAccApplicationInsights_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights", "test")
	r := AppInsightsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, "web"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("application_type").HasValue("web"),
				check.That(data.ResourceName).Key("retention_in_days").HasValue("120"),
				check.That(data.ResourceName).Key("sampling_percentage").HasValue("50"),
				check.That(data.ResourceName).Key("daily_data_cap_in_gb").HasValue("50"),
				check.That(data.ResourceName).Key("daily_data_cap_notifications_disabled").HasValue("true"),
				check.That(data.ResourceName).Key("local_authentication_disabled").HasValue("true"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.Hello").HasValue("World"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationInsights_withInternetQueryEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights", "test")
	r := AppInsightsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withInternetQueryEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withInternetQueryEnabledUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationInsights_withInternetIngestionEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights", "test")
	r := AppInsightsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withInternetIngestionEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withInternetIngestionEnabledUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationInsights_disableGeneratedRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights", "test")
	r := AppInsightsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disableGeneratedRule(data, "web"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("application_type").HasValue("web"),
			),
		},
		data.ImportStep(),
	})
}

func (AppInsightsResource) basic(data acceptance.TestData, applicationType string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appinsights-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, applicationType)
}

func (AppInsightsResource) basicWorkspaceMode(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appinsights-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_log_analytics_workspace.test.id
  application_type    = "web"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (AppInsightsResource) basicWorkspaceModeUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appinsights-%[2]d"
  location = "%[1]s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_log_analytics_workspace" "test2" {
  name                = "acctest2-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_log_analytics_workspace.test2.id
  application_type    = "web"
}
`, data.Locations.Primary, data.RandomInteger)
}

func (AppInsightsResource) requiresImport(data acceptance.TestData, applicationType string) string {
	template := AppInsightsResource{}.basic(data, applicationType)
	return fmt.Sprintf(`
%s

resource "azurerm_application_insights" "import" {
  name                = azurerm_application_insights.test.name
  location            = azurerm_application_insights.test.location
  resource_group_name = azurerm_application_insights.test.resource_group_name
  application_type    = azurerm_application_insights.test.application_type
}
`, template)
}

func (AppInsightsResource) complete(data acceptance.TestData, applicationType string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appinsights-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                                  = "acctestappinsights-%d"
  location                              = azurerm_resource_group.test.location
  resource_group_name                   = azurerm_resource_group.test.name
  application_type                      = "%s"
  retention_in_days                     = 120
  sampling_percentage                   = 50
  daily_data_cap_in_gb                  = 50
  daily_data_cap_notifications_disabled = true
  disable_ip_masking                    = true
  force_customer_storage_for_profiler   = true
  local_authentication_disabled         = true

  tags = {
    Hello = "World"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, applicationType)
}

func (AppInsightsResource) withInternetQueryEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                   = "acctestappinsights-%d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  application_type       = "web"
  internet_query_enabled = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (AppInsightsResource) withInternetQueryEnabledUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                   = "acctestappinsights-%d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  application_type       = "web"
  internet_query_enabled = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (AppInsightsResource) withInternetIngestionEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                       = "acctestappinsights-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  application_type           = "web"
  internet_ingestion_enabled = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (AppInsightsResource) withInternetIngestionEnabledUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                       = "acctestappinsights-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  application_type           = "web"
  internet_ingestion_enabled = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (AppInsightsResource) disableGeneratedRule(data acceptance.TestData, applicationType string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    application_insights {
      disable_generated_rule = true
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appinsights-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, applicationType)
}
