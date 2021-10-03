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

type AppConfigurationFeatureResource struct {
}

func TestAccAppConfigurationFeature_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration_feature", "test")
	r := AppConfigurationFeatureResource{}
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

func TestAccAppConfigurationFeature_basicNoLabel(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration_feature", "test")
	r := AppConfigurationFeatureResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicNoLabel(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppConfigurationFeature_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration_feature", "test")
	r := AppConfigurationFeatureResource{}
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

func TestAccAppConfigurationFeature_lockUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration_feature", "test")
	r := AppConfigurationFeatureResource{}
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

func TestAccAppConfigurationFeature_enabledUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration_feature", "test")
	r := AppConfigurationFeatureResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.enabledUpdate(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("false"),
			),
		},
		{
			Config: r.enabledUpdate(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
			),
		},
	})
}

func (t AppConfigurationFeatureResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {

	resourceID, err := parse.FeatureId(state.ID)
	if err != nil {
		return nil, fmt.Errorf("while parsing resource ID: %+v", err)
	}

	client, err := clients.AppConfiguration.DataPlaneClient(ctx, resourceID.ConfigurationStoreId)
	if err != nil {
		return nil, err
	}

	res, err := client.GetKeyValues(ctx, resourceID.Name, resourceID.Label, "", "", []string{})
	if err != nil {
		return nil, fmt.Errorf("while checking for key's %q existence: %+v", resourceID.Name, err)
	}

	return utils.Bool(res.Response().StatusCode == 200), nil
}

func (t AppConfigurationFeatureResource) basic(data acceptance.TestData) string {
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

resource "azurerm_app_configuration_feature" "test" {
  configuration_store_id = azurerm_app_configuration.test.id
  description            = "test description"
  name                   = "acctest-ackey-%d"
  label                  = "acctest-ackeylabel-%d"
  enabled                = true

  percentage_filter_value = 10

  timewindow_filter {
    start = "2019-11-12T07:20:50.52Z"
    end   = "2019-11-13T07:20:50.52Z"
  }

  targeting_filter {
    default_rollout_percentage = 39
    users                      = ["random", "user"]

    groups {
      name               = "testgroup"
      rollout_percentage = 50
    }

    groups {
      name               = "testgroup2"
      rollout_percentage = 30
    }
  }
}


`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (t AppConfigurationFeatureResource) basicNoLabel(data acceptance.TestData) string {
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

resource "azurerm_app_configuration_feature" "test" {
  configuration_store_id = azurerm_app_configuration.test.id
  description            = "test description"
  name                   = "acctest-ackey-%d"
  enabled                = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (t AppConfigurationFeatureResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_configuration_feature" "import" {
  configuration_store_id = azurerm_app_configuration_feature.test.configuration_store_id
  description            = azurerm_app_configuration_feature.test.description
  name                   = azurerm_app_configuration_feature.test.name
  label                  = azurerm_app_configuration_feature.test.label
  enabled                = azurerm_app_configuration_feature.test.enabled

  percentage_filter_value = 10

  timewindow_filter {
    start = "2019-11-12T07:20:50.52Z"
    end   = "2019-11-12T07:20:50.52Z"
  }

  targeting_filter {
    default_rollout_percentage = 39
    users                      = ["random", "user"]

    groups {
      name               = "testgroup"
      rollout_percentage = 50
    }

    groups {
      name               = "testgroup2"
      rollout_percentage = 30
    }
  }
}
`, t.basic(data))
}

func (t AppConfigurationFeatureResource) lockUpdate(data acceptance.TestData, lockStatus bool) string {
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

resource "azurerm_app_configuration_feature" "test" {
  configuration_store_id = azurerm_app_configuration.test.id
  description            = "test description"
  name                   = "acctest-ackey-%d"
  label                  = "acctest-ackeylabel-%d"
  enabled                = true
  locked                 = %t
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, lockStatus)
}

func (t AppConfigurationFeatureResource) enabledUpdate(data acceptance.TestData, enabledStatus bool) string {
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

resource "azurerm_app_configuration_feature" "test" {
  configuration_store_id = azurerm_app_configuration.test.id
  description            = "test description"
  name                   = "acctest-ackey-%d"
  label                  = "acctest-ackeylabel-%d"
  enabled                = %t
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, enabledStatus)
}
