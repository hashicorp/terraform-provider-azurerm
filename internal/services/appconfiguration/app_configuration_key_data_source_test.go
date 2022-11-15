package appconfiguration_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration"
)

type AppConfigurationKeyDataSource struct{}

func TestAccAppConfigurationKeyDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_configuration_key", "test")
	d := AppConfigurationKeyDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("content_type").HasValue("test"),
				check.That(data.ResourceName).Key("value").HasValue("a test"),
				check.That(data.ResourceName).Key("locked").HasValue("false"),
				check.That(data.ResourceName).Key("type").HasValue(appconfiguration.KeyTypeKV),
				check.That(data.ResourceName).Key("vault_key_reference").HasValue(""),
				check.That(data.ResourceName).Key("etag").IsSet(),
				check.That(data.ResourceName).Key("label").IsSet(),
			),
		},
	})
}

func TestAccAppConfigurationKeyDataSource_basicNoLabel(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_configuration_key", "test")
	d := AppConfigurationKeyDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basicNoLabel(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("content_type").HasValue("test"),
				check.That(data.ResourceName).Key("value").HasValue("a test"),
				check.That(data.ResourceName).Key("locked").HasValue("false"),
				check.That(data.ResourceName).Key("type").HasValue(appconfiguration.KeyTypeKV),
				check.That(data.ResourceName).Key("vault_key_reference").HasValue(""),
				check.That(data.ResourceName).Key("etag").IsSet(),
				check.That(data.ResourceName).Key("label").HasValue(""),
			),
		},
	})
}

func TestAccAppConfigurationKeyDataSource_basicVault(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_configuration_key", "test")
	d := AppConfigurationKeyDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.vaultKeyBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("content_type").HasValue(appconfiguration.VaultKeyContentType),
				check.That(data.ResourceName).Key("value").MatchesOtherKey(check.That(data.ResourceName).Key("vault_key_reference")),
				check.That(data.ResourceName).Key("locked").HasValue("false"),
				check.That(data.ResourceName).Key("type").HasValue(appconfiguration.KeyTypeVault),
				check.That(data.ResourceName).Key("vault_key_reference").MatchesOtherKey(check.That(data.ResourceName).Key("value")),
				check.That(data.ResourceName).Key("etag").IsSet(),
				check.That(data.ResourceName).Key("label").IsSet(),
			),
		},
	})
}

func (AppConfigurationKeyDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_app_configuration_key" "test" {
  key                    = azurerm_app_configuration_key.test.key
  label                  = azurerm_app_configuration_key.test.label
  configuration_store_id = azurerm_app_configuration.test.id
}
`, AppConfigurationKeyResource{}.basic(data))
}

func (AppConfigurationKeyDataSource) basicNoLabel(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_app_configuration_key" "test" {
  key                    = azurerm_app_configuration_key.test.key
  configuration_store_id = azurerm_app_configuration.test.id
}
`, AppConfigurationKeyResource{}.basicNoLabel(data))
}

func (AppConfigurationKeyDataSource) vaultKeyBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_app_configuration_key" "test" {
  key                    = azurerm_app_configuration_key.test.key
  label                  = azurerm_app_configuration_key.test.label
  configuration_store_id = azurerm_app_configuration.test.id
}
`, AppConfigurationKeyResource{}.vaultKeyBasic(data))
}
