package appconfiguration_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type AppConfigurationKeysDataSource struct{}

func TestAccAppConfigurationKeysDataSource_allkeys(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_configuration_keys", "test")
	d := AppConfigurationKeysDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.allKeys(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("items.#").HasValue("4"),
			),
		},
	})
}

func TestAccAppConfigurationKeysDataSource_key(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_configuration_keys", "test")
	d := AppConfigurationKeysDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.key(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("items.#").HasValue("2"),
			),
		},
	})
}

func TestAccAppConfigurationKeysDataSource_label(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_configuration_keys", "test")
	d := AppConfigurationKeysDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.label(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("items.#").HasValue("2"),
			),
		},
	})
}

func (t AppConfigurationKeysDataSource) keys() string {
	return fmt.Sprintf(`
resource "azurerm_app_configuration_key" "test" {
  configuration_store_id = azurerm_app_configuration.test.id
  key                    = "key1"
  content_type           = "test"
  label                  = "label1"
  value                  = "a test"
}

resource "azurerm_app_configuration_key" "test2" {
  configuration_store_id = azurerm_app_configuration.test.id
  key                    = "key1"
  content_type           = "test"
  label                  = "label2"
  value                  = "a test"
}

resource "azurerm_app_configuration_key" "test3" {
  configuration_store_id = azurerm_app_configuration.test.id
  key                    = "key2"
  content_type           = "test"
  label                  = "testlabel"
  value                  = "a test"
}

resource "azurerm_app_configuration_key" "test4" {
  configuration_store_id = azurerm_app_configuration.test.id
  key                    = "key3"
  content_type           = "test"
  label                  = "testlabel"
  value                  = "a test"
}
`)
}

func (t AppConfigurationKeysDataSource) allKeys(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

%s

data "azurerm_app_configuration_keys" "test" {
  configuration_store_id = azurerm_app_configuration.test.id

  depends_on = [
    azurerm_app_configuration_key.test,
    azurerm_app_configuration_key.test2,
    azurerm_app_configuration_key.test3,
    azurerm_app_configuration_key.test4
  ]
}
`, AppConfigurationKeyResource{}.base(data), t.keys())
}

func (t AppConfigurationKeysDataSource) key(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

%s

data "azurerm_app_configuration_keys" "test" {
  configuration_store_id = azurerm_app_configuration.test.id
  key                    = "key1"

  depends_on = [
    azurerm_app_configuration_key.test,
    azurerm_app_configuration_key.test2,
    azurerm_app_configuration_key.test3,
    azurerm_app_configuration_key.test4
  ]
}
`, AppConfigurationKeyResource{}.base(data), t.keys())
}

func (t AppConfigurationKeysDataSource) label(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

%s

data "azurerm_app_configuration_keys" "test" {
  configuration_store_id = azurerm_app_configuration.test.id
  label                  = "testlabel"

  depends_on = [
    azurerm_app_configuration_key.test,
    azurerm_app_configuration_key.test2,
    azurerm_app_configuration_key.test3,
    azurerm_app_configuration_key.test4
  ]
}
`, AppConfigurationKeyResource{}.base(data), t.keys())
}
