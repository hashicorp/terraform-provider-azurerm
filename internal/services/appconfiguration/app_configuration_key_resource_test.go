package appconfiguration_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AppConfigurationKeyResource struct {
}

func TestAccAppConfigurationKey_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration_key", "test")
	r := AppConfigurationKeyResource{}
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

func TestAccAppConfigurationKey_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration_key", "test")
	r := AppConfigurationKeyResource{}
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
func TestAccAppConfigurationKey_lockUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration_key", "test")
	r := AppConfigurationKeyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.lockUpdate(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("locked").HasValue("false"),
			),
		},
		{
			Config: r.lockUpdate(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("locked").HasValue("true"),
			),
		},
	})
}
func (t AppConfigurationKeyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {

	resourceID, err := parse.AppConfigurationKeyID(state.ID)
	if err != nil {
		return nil, fmt.Errorf("while parsing resource ID: %+v", err)
	}

	client, err := clients.AppConfiguration.DataPlaneClient(ctx, resourceID.ConfigurationStoreId)
	if err != nil {
		return nil, err
	}

	res, err := client.GetKeyValues(ctx, resourceID.Key, resourceID.Label, "", "", []string{})
	if err != nil {
		return nil, fmt.Errorf("while checking for key's %q existence: %+v", resourceID.Key, err)
	}

	return utils.Bool(res.Response().StatusCode == 200), nil
}

func (t AppConfigurationKeyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appconfig-%d"
  location = "%s"
}

resource "azurerm_app_configuration" "test" {
  name                = "testacc-appconf%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "standard"
}

resource "azurerm_app_configuration_key" "test" {
  configuration_store_id = azurerm_app_configuration.test.id
  key                    = "acctest-ackey-%d"
  content_type           = "test"
  label                  = "acctest-ackeylabel-%d"
  value                  = "a test"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (t AppConfigurationKeyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_configuration_key" "import" {
  configuration_store_id = azurerm_app_configuration_key.test.configuration_store_id
  key                    = azurerm_app_configuration_key.test.key
  content_type           = azurerm_app_configuration_key.test.content_type
  label                  = azurerm_app_configuration_key.test.label
  value                  = azurerm_app_configuration_key.test.value
}
`, t.basic(data))
}

func (t AppConfigurationKeyResource) lockUpdate(data acceptance.TestData, lockStatus bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appconfig-%d"
  location = "%s"
}

resource "azurerm_app_configuration" "test" {
  name                = "testacc-appconf%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "standard"
}

resource "azurerm_app_configuration_key" "test" {
  configuration_store_id = azurerm_app_configuration.test.id
  key                    = "acctest-ackey-%d"
  content_type           = "test"
  label                  = "acctest-ackeylabel-%d"
  value                  = "a test"
  locked                 = %t
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, lockStatus)
}
