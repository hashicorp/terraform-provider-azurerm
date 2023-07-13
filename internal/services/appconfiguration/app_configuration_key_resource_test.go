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

type AppConfigurationKeyResource struct{}

func TestAccAppConfigurationKey_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration_key", "test")
	r := AppConfigurationKeyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("etag").IsSet(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppConfigurationKey_basicNoLabel(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration_key", "test")
	r := AppConfigurationKeyResource{}
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

func TestAccAppConfigurationKey_complicatedKeyLabel(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration_key", "test")
	r := AppConfigurationKeyResource{}
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

func TestAccAppConfigurationKey_basicNoLabel_afterLabel(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration_key", "test")
	r := AppConfigurationKeyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicNoLabelAfterLabel(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppConfigurationKey_basicVault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration_key", "test")
	r := AppConfigurationKeyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.vaultKeyBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppConfigurationKey_KVToVault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration_key", "test")
	r := AppConfigurationKeyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("type").HasValue("kv"),
			),
		},
		{
			Config: r.vaultKeyBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("type").HasValue("vault"),
			),
		},
	})
}

func TestAccAppConfigurationKey_slash(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration_key", "test")
	r := AppConfigurationKeyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.slash(data),
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

func (t AppConfigurationKeyResource) base(data acceptance.TestData) string {
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
  name                = "testacc-appconf%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "standard"
  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (t AppConfigurationKeyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_configuration_key" "test" {
  configuration_store_id = azurerm_app_configuration.test.id
  key                    = "acctest-ackey-%d"
  content_type           = "test"
  label                  = "acctest-ackeylabel-%d"
  value                  = "a test"
}
`, t.base(data), data.RandomInteger, data.RandomInteger)
}

func (t AppConfigurationKeyResource) complicatedKeyLabel(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_configuration_key" "test" {
  configuration_store_id = azurerm_app_configuration.test.id
  key                    = "acctest-ackey-%d/Label/AppConfigurationKey/Label/"
  content_type           = "test"
  label                  = "/AppConfigurationKey/acctest-ackeylabel-%d"
  value                  = "a test"
}
`, t.base(data), data.RandomInteger, data.RandomInteger)
}

func (t AppConfigurationKeyResource) slash(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_configuration_key" "test" {
  configuration_store_id = azurerm_app_configuration.test.id
  key                    = "/acctest/-ackey/-%d"
  content_type           = "test"
  label                  = "/acctest/-ackeylabel/-%d"
  value                  = "a test"
}
`, t.base(data), data.RandomInteger, data.RandomInteger)
}

func (t AppConfigurationKeyResource) basicNoLabel(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_configuration_key" "test" {
  configuration_store_id = azurerm_app_configuration.test.id
  key                    = "acctest-ackey-%d"
  content_type           = "test"
  value                  = "a test"
}
`, t.base(data), data.RandomInteger)
}

func (t AppConfigurationKeyResource) basicNoLabelAfterLabel(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_configuration_key" "test1" {
  configuration_store_id = azurerm_app_configuration.test.id
  key                    = "acctest-ackey-%d"
  content_type           = "test"
  value                  = "a test"
}
`, t.basic(data), data.RandomInteger)
}

func (t AppConfigurationKeyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_configuration_key" "import" {
  configuration_store_id = azurerm_app_configuration_key.test.configuration_store_id
  key                    = azurerm_app_configuration_key.test.key
  content_type           = azurerm_app_configuration_key.test.content_type
  label                  = azurerm_app_configuration_key.test.label
  value                  = "${azurerm_app_configuration_key.test.value}-another"
}
`, t.basic(data))
}

func (t AppConfigurationKeyResource) lockUpdate(data acceptance.TestData, lockStatus bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_configuration_key" "test" {
  configuration_store_id = azurerm_app_configuration.test.id
  key                    = "acctest-ackey-%d"
  content_type           = "test"
  label                  = "acctest-ackeylabel-%d"
  value                  = "a test"
  locked                 = %t
}
`, t.base(data), data.RandomInteger, data.RandomInteger, lockStatus)
}

func (t AppConfigurationKeyResource) vaultKeyBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault" "example" {
  name                       = "a-v-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.test.tenant_id
  sku_name                   = "premium"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.test.tenant_id
    object_id = data.azurerm_client_config.test.object_id

    key_permissions = [
      "Create",
      "Get",
    ]

    secret_permissions = [
      "Set",
      "Get",
      "Delete",
      "Purge",
      "Recover"
    ]
  }
}

resource "azurerm_key_vault_secret" "example" {
  name         = "acctest-secret-%d"
  value        = "szechuan"
  key_vault_id = azurerm_key_vault.example.id
}

resource "azurerm_app_configuration_key" "test" {
  configuration_store_id = azurerm_app_configuration.test.id
  key                    = "acctest-ackey-%d"
  type                   = "vault"
  label                  = "acctest-ackeylabel-%d"
  vault_key_reference    = azurerm_key_vault_secret.example.id
}
`, t.base(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
