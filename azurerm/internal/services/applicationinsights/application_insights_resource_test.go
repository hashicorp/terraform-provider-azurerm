package applicationinsights_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/applicationinsights/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type AppInsightsResource struct {
}

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

func (t AppInsightsResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ComponentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.AppInsights.ComponentsClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Application Insights %q (resource group: %q) does not exist", id.Name, id.ResourceGroup)
	}

	return utils.Bool(resp.ApplicationInsightsComponentProperties != nil), nil
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
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.Hello").HasValue("World"),
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

  tags = {
    Hello = "World"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, applicationType)
}
