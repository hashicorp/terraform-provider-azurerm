package applicationinsights_test

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type AppInsightsAPIKey struct {
}

func TestAccApplicationInsightsAPIKey_no_permission(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_api_key", "test")
	r := AppInsightsAPIKey{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:      r.basic(data, "[]", "[]"),
			ExpectError: regexp.MustCompile("The API Key needs to have a Role"),
		},
	})
}

func TestAccApplicationInsightsAPIKey_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_api_key", "test")
	r := AppInsightsAPIKey{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, "[]", `["annotations"]`),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("read_permissions.#").HasValue("0"),
				check.That(data.ResourceName).Key("write_permissions.#").HasValue("1"),
			),
		},
		{
			Config:      r.requiresImport(data, "[]", `["annotations"]`),
			ExpectError: acceptance.RequiresImportError("azurerm_application_insights_api_key"),
		},
	})
}

func TestAccApplicationInsightsAPIKey_read_telemetry_permissions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_api_key", "test")
	r := AppInsightsAPIKey{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, `["aggregate", "api", "draft", "extendqueries", "search"]`, "[]"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("read_permissions.#").HasValue("5"),
				check.That(data.ResourceName).Key("write_permissions.#").HasValue("0"),
			),
		},
		data.ImportStep("api_key"),
	})
}

func TestAccApplicationInsightsAPIKey_write_annotations_permission(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_api_key", "test")
	r := AppInsightsAPIKey{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, "[]", `["annotations"]`),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("read_permissions.#").HasValue("0"),
				check.That(data.ResourceName).Key("write_permissions.#").HasValue("1"),
			),
		},
		data.ImportStep("api_key"),
	})
}

func TestAccApplicationInsightsAPIKey_authenticate_permission(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_api_key", "test")
	r := AppInsightsAPIKey{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, `["agentconfig"]`, "[]"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("read_permissions.#").HasValue("1"),
				check.That(data.ResourceName).Key("write_permissions.#").HasValue("0"),
			),
		},
		data.ImportStep("api_key"),
	})
}

func TestAccApplicationInsightsAPIKey_full_permissions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_api_key", "test")
	r := AppInsightsAPIKey{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, `["agentconfig", "aggregate", "api", "draft", "extendqueries", "search"]`, `["annotations"]`),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("read_permissions.#").HasValue("6"),
				check.That(data.ResourceName).Key("write_permissions.#").HasValue("1"),
			),
		},
		data.ImportStep("api_key"),
	})
}

func (t AppInsightsAPIKey) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.Attributes["id"])
	if err != nil {
		return nil, err
	}
	keyID := id.Path["APIKeys"]
	resGroup := id.ResourceGroup
	appInsightsName := id.Path["components"]

	resp, err := clients.AppInsights.APIKeysClient.Get(ctx, resGroup, appInsightsName, keyID)
	if err != nil {
		return nil, fmt.Errorf("retrieving Application Insights API Key '%q' (resource group: '%q') does not exist", keyID, resGroup)
	}

	return utils.Bool(resp.StatusCode != http.StatusNotFound), nil
}

func (AppInsightsAPIKey) basic(data acceptance.TestData, readPerms, writePerms string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_application_insights_api_key" "test" {
  name                    = "acctestappinsightsapikey-%d"
  application_insights_id = azurerm_application_insights.test.id
  read_permissions        = %s
  write_permissions       = %s
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, readPerms, writePerms)
}

func (AppInsightsAPIKey) requiresImport(data acceptance.TestData, readPerms, writePerms string) string {
	template := AppInsightsAPIKey{}.basic(data, readPerms, writePerms)
	return fmt.Sprintf(`
%s

resource "azurerm_application_insights_api_key" "import" {
  name                    = azurerm_application_insights_api_key.test.name
  application_insights_id = azurerm_application_insights_api_key.test.application_insights_id
  read_permissions        = azurerm_application_insights_api_key.test.read_permissions
  write_permissions       = azurerm_application_insights_api_key.test.write_permissions
}
`, template)
}
