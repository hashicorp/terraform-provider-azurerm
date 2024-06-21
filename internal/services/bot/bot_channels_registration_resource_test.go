// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package bot_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type BotChannelsRegistrationResource struct{}

func TestAccBotChannelsRegistration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channels_registration", "test")
	r := BotChannelsRegistrationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("developer_app_insights_api_key"),
	})
}

func TestAccBotChannelsRegistration_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channels_registration", "test")
	r := BotChannelsRegistrationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("developer_app_insights_api_key"),
		{
			Config: r.completeConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("developer_app_insights_api_key"),
		{
			Config: r.updateConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("developer_app_insights_api_key"),
		{
			Config: r.withoutCMKKeyVaultURL(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("developer_app_insights_api_key"),
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("developer_app_insights_api_key"),
	})
}

func TestAccBotChannelsRegistration_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channels_registration", "test")
	r := BotChannelsRegistrationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("developer_app_insights_api_key"),
	})
}

func TestAccBotChannelsRegistration_streamingEndpointEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channels_registration", "test")
	r := BotChannelsRegistrationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.streamingEndpointEnabled(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("developer_app_insights_api_key"),
		{
			Config: r.streamingEndpointEnabled(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("developer_app_insights_api_key"),
	})
}

func (t BotChannelsRegistrationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.BotServiceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Bot.BotClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (BotChannelsRegistrationResource) basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azuread_application_registration" "test" {
  display_name = "acctestReg-%d"
}

resource "azurerm_bot_channels_registration" "test" {
  name                = "acctestdf%d"
  location            = "global"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "F0"
  microsoft_app_id    = azuread_application_registration.test.client_id

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (BotChannelsRegistrationResource) updateConfig(data acceptance.TestData) string {
	if !features.FourPointOhBeta() {
		return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {
}

data "azuread_service_principal" "test" {
  display_name = "Bot Service CMEK Prod"
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
  read_permissions        = ["aggregate", "api", "draft", "extendqueries", "search"]
}

resource "azurerm_application_insights" "test2" {
  name                = "acctestappinsights2-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_application_insights_api_key" "test2" {
  name                    = "acctestappinsightsapikey2-%d"
  application_insights_id = azurerm_application_insights.test2.id
  read_permissions        = ["aggregate", "api", "draft", "extendqueries", "search"]
}

resource "azurerm_key_vault" "test" {
  name                     = "acctestkv-%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions    = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"]
  secret_permissions = ["Get"]
}

resource "azurerm_key_vault_access_policy" "test2" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azuread_service_principal.test.id

  key_permissions    = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"]
  secret_permissions = ["Get"]
}

resource "azurerm_key_vault_key" "test" {
  name         = "acctestkvkey-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.test,
    azurerm_key_vault_access_policy.test2,
  ]
}

resource "azurerm_key_vault_key" "test2" {
  name         = "acctestkvkey2-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.test,
    azurerm_key_vault_access_policy.test2,
  ]
}

resource "azuread_application_registration" "test" {
  display_name = "acctestReg-%d"
}

resource "azurerm_bot_channels_registration" "test" {
  name                = "acctestdf%d"
  location            = "global"
  resource_group_name = azurerm_resource_group.test.name
  microsoft_app_id    = azuread_application_registration.test.client_id
  sku                 = "F0"

  endpoint                              = "https://example2.com"
  developer_app_insights_api_key        = azurerm_application_insights_api_key.test2.api_key
  developer_app_insights_application_id = azurerm_application_insights.test2.app_id
  developer_app_insights_key            = azurerm_application_insights.test2.instrumentation_key

  description              = "TestDescription2"
  isolated_network_enabled = false
  icon_url                 = "http://myprofile/myicon2.png"
  cmk_key_vault_url        = azurerm_key_vault_key.test2.id

  tags = {
    environment = "production2"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomString, data.RandomString, data.RandomInteger, data.RandomInteger)
	}
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {
}

data "azuread_service_principal" "test" {
  display_name = "Bot Service CMEK Prod"
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
  read_permissions        = ["aggregate", "api", "draft", "extendqueries", "search"]
}

resource "azurerm_application_insights" "test2" {
  name                = "acctestappinsights2-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_application_insights_api_key" "test2" {
  name                    = "acctestappinsightsapikey2-%d"
  application_insights_id = azurerm_application_insights.test2.id
  read_permissions        = ["aggregate", "api", "draft", "extendqueries", "search"]
}

resource "azurerm_key_vault" "test" {
  name                     = "acctestkv-%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions    = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"]
  secret_permissions = ["Get"]
}

resource "azurerm_key_vault_access_policy" "test2" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azuread_service_principal.test.id

  key_permissions    = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"]
  secret_permissions = ["Get"]
}

resource "azurerm_key_vault_key" "test" {
  name         = "acctestkvkey-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.test,
    azurerm_key_vault_access_policy.test2,
  ]
}

resource "azurerm_key_vault_key" "test2" {
  name         = "acctestkvkey2-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.test,
    azurerm_key_vault_access_policy.test2,
  ]
}

resource "azuread_application_registration" "test" {
  display_name = "acctestReg-%d"
}

resource "azurerm_bot_channels_registration" "test" {
  name                = "acctestdf%d"
  location            = "global"
  resource_group_name = azurerm_resource_group.test.name
  microsoft_app_id    = azuread_application_registration.test.client_id
  sku                 = "F0"

  endpoint                              = "https://example2.com"
  developer_app_insights_api_key        = azurerm_application_insights_api_key.test2.api_key
  developer_app_insights_application_id = azurerm_application_insights.test2.app_id
  developer_app_insights_key            = azurerm_application_insights.test2.instrumentation_key

  description                   = "TestDescription2"
  public_network_access_enabled = false
  icon_url                      = "http://myprofile/myicon2.png"
  cmk_key_vault_url             = azurerm_key_vault_key.test2.id

  tags = {
    environment = "production2"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomString, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (BotChannelsRegistrationResource) completeConfig(data acceptance.TestData) string {
	if !features.FourPointOhBeta() {
		return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {
}

data "azuread_service_principal" "test" {
  display_name = "Bot Service CMEK Prod"
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
  read_permissions        = ["aggregate", "api", "draft", "extendqueries", "search"]
}

resource "azurerm_key_vault" "test" {
  name                     = "acctestkv-%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions    = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"]
  secret_permissions = ["Get"]
}

resource "azurerm_key_vault_access_policy" "test2" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azuread_service_principal.test.id

  key_permissions    = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"]
  secret_permissions = ["Get"]
}

resource "azurerm_key_vault_key" "test" {
  name         = "acctestkvkey-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.test,
    azurerm_key_vault_access_policy.test2,
  ]
}

resource "azuread_application_registration" "test" {
  display_name = "acctestReg-%d"
}

resource "azurerm_bot_channels_registration" "test" {
  name                = "acctestdf%d"
  location            = "global"
  resource_group_name = azurerm_resource_group.test.name
  microsoft_app_id    = azuread_application_registration.test.client_id
  sku                 = "F0"

  endpoint                              = "https://example.com"
  developer_app_insights_api_key        = azurerm_application_insights_api_key.test.api_key
  developer_app_insights_application_id = azurerm_application_insights.test.app_id
  developer_app_insights_key            = azurerm_application_insights.test.instrumentation_key

  description              = "TestDescription"
  isolated_network_enabled = true
  icon_url                 = "http://myprofile/myicon.png"
  cmk_key_vault_url        = azurerm_key_vault_key.test.id

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomString, data.RandomInteger, data.RandomInteger)
	}
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {
}

data "azuread_service_principal" "test" {
  display_name = "Bot Service CMEK Prod"
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
  read_permissions        = ["aggregate", "api", "draft", "extendqueries", "search"]
}

resource "azurerm_key_vault" "test" {
  name                     = "acctestkv-%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions    = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"]
  secret_permissions = ["Get"]
}

resource "azurerm_key_vault_access_policy" "test2" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azuread_service_principal.test.id

  key_permissions    = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"]
  secret_permissions = ["Get"]
}

resource "azurerm_key_vault_key" "test" {
  name         = "acctestkvkey-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.test,
    azurerm_key_vault_access_policy.test2,
  ]
}

resource "azuread_application_registration" "test" {
  display_name = "acctestReg-%d"
}

resource "azurerm_bot_channels_registration" "test" {
  name                = "acctestdf%d"
  location            = "global"
  resource_group_name = azurerm_resource_group.test.name
  microsoft_app_id    = azuread_application_registration.test.client_id
  sku                 = "F0"

  endpoint                              = "https://example.com"
  developer_app_insights_api_key        = azurerm_application_insights_api_key.test.api_key
  developer_app_insights_application_id = azurerm_application_insights.test.app_id
  developer_app_insights_key            = azurerm_application_insights.test.instrumentation_key

  description                   = "TestDescription"
  public_network_access_enabled = true
  icon_url                      = "http://myprofile/myicon.png"
  cmk_key_vault_url             = azurerm_key_vault_key.test.id

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (BotChannelsRegistrationResource) withoutCMKKeyVaultURL(data acceptance.TestData) string {
	if !features.FourPointOhBeta() {
		return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {
}

data "azuread_service_principal" "test" {
  display_name = "Bot Service CMEK Prod"
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
  read_permissions        = ["aggregate", "api", "draft", "extendqueries", "search"]
}

resource "azurerm_application_insights" "test2" {
  name                = "acctestappinsights2-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_application_insights_api_key" "test2" {
  name                    = "acctestappinsightsapikey2-%d"
  application_insights_id = azurerm_application_insights.test2.id
  read_permissions        = ["aggregate", "api", "draft", "extendqueries", "search"]
}

resource "azurerm_key_vault" "test" {
  name                     = "acctestkv-%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions    = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"]
  secret_permissions = ["Get"]
}

resource "azurerm_key_vault_access_policy" "test2" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azuread_service_principal.test.id

  key_permissions    = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"]
  secret_permissions = ["Get"]
}

resource "azurerm_key_vault_key" "test" {
  name         = "acctestkvkey-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.test,
    azurerm_key_vault_access_policy.test2,
  ]
}

resource "azurerm_key_vault_key" "test2" {
  name         = "acctestkvkey2-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.test,
    azurerm_key_vault_access_policy.test2,
  ]
}

resource "azuread_application_registration" "test" {
  display_name = "acctestReg-%d"
}

resource "azurerm_bot_channels_registration" "test" {
  name                = "acctestdf%d"
  location            = "global"
  resource_group_name = azurerm_resource_group.test.name
  microsoft_app_id    = azuread_application_registration.test.client_id
  sku                 = "F0"

  endpoint                              = "https://example2.com"
  developer_app_insights_api_key        = azurerm_application_insights_api_key.test2.api_key
  developer_app_insights_application_id = azurerm_application_insights.test2.app_id
  developer_app_insights_key            = azurerm_application_insights.test2.instrumentation_key

  description              = "TestDescription2"
  isolated_network_enabled = false
  icon_url                 = "http://myprofile/myicon2.png"

  tags = {
    environment = "production2"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomString, data.RandomString, data.RandomInteger, data.RandomInteger)
	}
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {
}

data "azuread_service_principal" "test" {
  display_name = "Bot Service CMEK Prod"
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
  read_permissions        = ["aggregate", "api", "draft", "extendqueries", "search"]
}

resource "azurerm_application_insights" "test2" {
  name                = "acctestappinsights2-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_application_insights_api_key" "test2" {
  name                    = "acctestappinsightsapikey2-%d"
  application_insights_id = azurerm_application_insights.test2.id
  read_permissions        = ["aggregate", "api", "draft", "extendqueries", "search"]
}

resource "azurerm_key_vault" "test" {
  name                     = "acctestkv-%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions    = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"]
  secret_permissions = ["Get"]
}

resource "azurerm_key_vault_access_policy" "test2" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azuread_service_principal.test.id

  key_permissions    = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"]
  secret_permissions = ["Get"]
}

resource "azurerm_key_vault_key" "test" {
  name         = "acctestkvkey-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.test,
    azurerm_key_vault_access_policy.test2,
  ]
}

resource "azurerm_key_vault_key" "test2" {
  name         = "acctestkvkey2-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.test,
    azurerm_key_vault_access_policy.test2,
  ]
}

resource "azuread_application_registration" "test" {
  display_name = "acctestReg-%d"
}

resource "azurerm_bot_channels_registration" "test" {
  name                = "acctestdf%d"
  location            = "global"
  resource_group_name = azurerm_resource_group.test.name
  microsoft_app_id    = azuread_application_registration.test.client_id
  sku                 = "F0"

  endpoint                              = "https://example2.com"
  developer_app_insights_api_key        = azurerm_application_insights_api_key.test2.api_key
  developer_app_insights_application_id = azurerm_application_insights.test2.app_id
  developer_app_insights_key            = azurerm_application_insights.test2.instrumentation_key

  description                   = "TestDescription2"
  public_network_access_enabled = false
  icon_url                      = "http://myprofile/myicon2.png"

  tags = {
    environment = "production2"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomString, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (BotChannelsRegistrationResource) streamingEndpointEnabled(data acceptance.TestData, streamingEndpointEnabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azuread_application_registration" "test" {
  display_name = "acctestReg-%d"
}

resource "azurerm_bot_channels_registration" "test" {
  name                       = "acctestdf%d"
  location                   = "global"
  resource_group_name        = azurerm_resource_group.test.name
  sku                        = "F0"
  microsoft_app_id           = azuread_application_registration.test.client_id
  streaming_endpoint_enabled = %t
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, streamingEndpointEnabled)
}
