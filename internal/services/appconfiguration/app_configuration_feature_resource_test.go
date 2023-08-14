// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	"github.com/tombuildsstuff/kermit/sdk/appconfiguration/1.0/appconfiguration"
)

type AppConfigurationFeatureResource struct{}

func TestAccAppConfigurationFeature_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration_feature", "test")
	r := AppConfigurationFeatureResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("percentage_filter_value").HasValue("10"),

				check.That(data.ResourceName).Key("timewindow_filter.#").HasValue("1"),
				check.That(data.ResourceName).Key("timewindow_filter.0.start").HasValue("2019-11-12T07:20:50.52Z"),
				check.That(data.ResourceName).Key("timewindow_filter.0.end").HasValue("2019-11-13T07:20:50.52Z"),

				check.That(data.ResourceName).Key("targeting_filter.#").HasValue("1"),
				check.That(data.ResourceName).Key("targeting_filter.0.default_rollout_percentage").HasValue("39"),
				check.That(data.ResourceName).Key("targeting_filter.0.users.#").HasValue("2"),
				check.That(data.ResourceName).Key("targeting_filter.0.users.0").HasValue("random"),
				check.That(data.ResourceName).Key("targeting_filter.0.users.1").HasValue("user"),

				check.That(data.ResourceName).Key("targeting_filter.0.groups.#").HasValue("2"),
				check.That(data.ResourceName).Key("targeting_filter.0.groups.0.name").HasValue("testgroup"),
				check.That(data.ResourceName).Key("targeting_filter.0.groups.0.rollout_percentage").HasValue("50"),
				check.That(data.ResourceName).Key("targeting_filter.0.groups.1.name").HasValue("testgroup2"),
				check.That(data.ResourceName).Key("targeting_filter.0.groups.1.rollout_percentage").HasValue("30"),
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

func TestAccAppConfigurationFeature_customKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration_feature", "test")
	r := AppConfigurationFeatureResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppConfigurationFeature_basicWithSlash(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration_feature", "test")
	r := AppConfigurationFeatureResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithSlash(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppConfigurationFeature_basicNoFilters(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration_feature", "test")
	r := AppConfigurationFeatureResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicNoFilters(data),
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

func TestAccAppConfigurationFeature_noLabelUpdate(t *testing.T) {
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
		{
			Config: r.updateNoLabel(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
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

func TestAccAppConfigurationFeature_complicatedKeyLabel(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration_feature", "test")
	r := AppConfigurationFeatureResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complicatedKeyLabel(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
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
	nestedItemId, err := parse.ParseNestedItemID(state.ID)
	if err != nil {
		return nil, fmt.Errorf("while parsing resource ID: %+v", err)
	}

	client, err := clients.AppConfiguration.DataPlaneClientWithEndpoint(nestedItemId.ConfigurationStoreEndpoint)
	if err != nil {
		return nil, err
	}

	res, err := client.GetKeyValue(ctx, nestedItemId.Key, nestedItemId.Label, "", "", "", []appconfiguration.KeyValueFields{})
	if err != nil {
		return nil, fmt.Errorf("while checking for key's %q existence: %+v", nestedItemId.Key, err)
	}

	return utils.Bool(res.Response.StatusCode == 200), nil
}

func (t AppConfigurationFeatureResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_configuration_feature" "test" {
  configuration_store_id = azurerm_app_configuration.test.id
  description            = "test description"
  name                   = "acctest-ackey-%d"
  label                  = "acctest-ackeylabel-%[2]d"
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

`, t.template(data), data.RandomInteger)
}

func (t AppConfigurationFeatureResource) basicNoLabel(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_configuration_feature" "test" {
  configuration_store_id = azurerm_app_configuration.test.id
  description            = "test description"
  name                   = "acctest-ackey-%d"
  enabled                = true
}
`, t.template(data), data.RandomInteger)
}

func (t AppConfigurationFeatureResource) customKey(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_configuration_feature" "test" {
  configuration_store_id = azurerm_app_configuration.test.id
  description            = "test description"
  name                   = "acctest-ackey-%d"
  key                    = "custom/:-key-%[2]d"
  enabled                = true
}
`, t.template(data), data.RandomInteger)
}

func (t AppConfigurationFeatureResource) updateNoLabel(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_configuration_feature" "test" {
  configuration_store_id = azurerm_app_configuration.test.id
  description            = "test description"
  name                   = "acctest-ackey-%d"
  enabled                = false
}
`, t.template(data), data.RandomInteger)
}

func (t AppConfigurationFeatureResource) basicWithSlash(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_configuration_feature" "test" {
  configuration_store_id = azurerm_app_configuration.test.id
  description            = "test description"
  name                   = "acctest/ackey/%d"
  label                  = "acctest/label"
  enabled                = true
}
`, t.template(data), data.RandomInteger)
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

func (t AppConfigurationFeatureResource) complicatedKeyLabel(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_configuration_feature" "test" {
  configuration_store_id = azurerm_app_configuration.test.id
  name                   = "acctest-ackey-%d/Label/AppConfigurationKey/Label/"
  label                  = "/Key/AppConfigurationKey/Label/acctest-ackeylabel-%[2]d"
  enabled                = true
}
`, t.template(data), data.RandomInteger)
}

func (t AppConfigurationFeatureResource) lockUpdate(data acceptance.TestData, lockStatus bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_configuration_feature" "test" {
  configuration_store_id = azurerm_app_configuration.test.id
  description            = "test description"
  name                   = "acctest-ackey-%d"
  label                  = "acctest-ackeylabel-%[2]d"
  enabled                = true
  locked                 = %[3]t
}
`, t.template(data), data.RandomInteger, lockStatus)
}

func (t AppConfigurationFeatureResource) enabledUpdate(data acceptance.TestData, enabledStatus bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_configuration_feature" "test" {
  configuration_store_id = azurerm_app_configuration.test.id
  description            = "test description"
  name                   = "acctest-ackey-%d"
  label                  = "acctest-ackeylabel-%[2]d"
  enabled                = %[3]t
}
`, t.template(data), data.RandomInteger, enabledStatus)
}

func (t AppConfigurationFeatureResource) basicNoFilters(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_configuration_feature" "test" {
  configuration_store_id = azurerm_app_configuration.test.id
  description            = "test description"
  name                   = "acctest-ackey-%[2]d"
  label                  = "acctest-ackeylabel-%[2]d"
  enabled                = true
}
`, t.template(data), data.RandomInteger)
}

func (t AppConfigurationFeatureResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appconfig-%d"
  location = "%s"
}

data "azurerm_client_config" "test" {
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_resource_group.test.id
  role_definition_name = "App Configuration Data Owner"
  principal_id         = data.azurerm_client_config.test.object_id
}

resource "azurerm_app_configuration" "test" {
  name                = "testacc-appconf%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "standard"

  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, data.RandomInteger, data.Locations.Primary)
}
