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
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type BotServiceAzureBotResource struct{}

func TestAccBotServiceAzureBot_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_service_azure_bot", "test")
	r := BotServiceAzureBotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku").HasValue("F0"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("test"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBotServiceAzureBot_completeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_service_azure_bot", "test")
	r := BotServiceAzureBotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("developer_app_insights_api_key"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("developer_app_insights_api_key"),
	})
}

func TestAccBotServiceAzureBot_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_service_azure_bot", "test")
	r := BotServiceAzureBotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_bot_service_azure_bot"),
		},
	})
}

func TestAccBotServiceAzureBot_msaAppType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_service_azure_bot", "test")
	r := BotServiceAzureBotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.msaAppType(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBotServiceAzureBot_streamingEndpointEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_service_azure_bot", "test")
	r := BotServiceAzureBotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.steamingEndpointEnabled(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.steamingEndpointEnabled(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBotServiceAzureBot_cmekEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_service_azure_bot", "test")
	r := BotServiceAzureBotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cmekEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cmk_key_vault_key_url").IsNotEmpty(),
			),
		},
		data.ImportStep(),
		{
			Config: r.cmekEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cmk_key_vault_key_url").IsNotEmpty(),
			),
		},
		data.ImportStep(),
	})
}

func (t BotServiceAzureBotResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.BotServiceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Bot.BotClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", *id, err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (BotServiceAzureBotResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azuread_application_registration" "test" {
  display_name = "acctestReg-%[1]d"
}

resource "azurerm_bot_service_azure_bot" "test" {
  name                = "acctestdf%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "global"
  sku                 = "F0"
  microsoft_app_id    = azuread_application_registration.test.client_id

  tags = {
    environment = "test"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (BotServiceAzureBotResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_application_insights_api_key" "test" {
  name                    = "acctestappinsightsapikey-%[1]d"
  application_insights_id = azurerm_application_insights.test.id
  read_permissions        = ["aggregate", "api", "draft", "extendqueries", "search"]
}

resource "azuread_application_registration" "test" {
  display_name = "acctestReg-%[1]d"
}

resource "azurerm_bot_service_azure_bot" "test" {
  name                                  = "acctestdf%[1]d"
  resource_group_name                   = azurerm_resource_group.test.name
  location                              = "global"
  microsoft_app_id                      = azuread_application_registration.test.client_id
  sku                                   = "F0"
  local_authentication_enabled          = false
  public_network_access_enabled         = false
  icon_url                              = "https://registry.terraform.io/images/providers/azure.png"
  endpoint                              = "https://example.com"
  developer_app_insights_api_key        = azurerm_application_insights_api_key.test.api_key
  developer_app_insights_application_id = azurerm_application_insights.test.app_id

  tags = {
    environment = "test2"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (BotServiceAzureBotResource) requiresImport(data acceptance.TestData) string {
	template := BotServiceAzureBotResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_bot_service_azure_bot" "import" {
  name                = azurerm_bot_service_azure_bot.test.name
  resource_group_name = azurerm_bot_service_azure_bot.test.resource_group_name
  location            = azurerm_bot_service_azure_bot.test.location
  sku                 = azurerm_bot_service_azure_bot.test.sku
  microsoft_app_id    = azurerm_bot_service_azure_bot.test.microsoft_app_id
}
`, template)
}

func (BotServiceAzureBotResource) msaAppType(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azuread_application_registration" "test" {
  display_name = "acctestReg-%[1]d"
}

resource "azurerm_bot_service_azure_bot" "test" {
  name                = "acctestdf%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "global"
  sku                 = "F0"
  microsoft_app_id    = azuread_application_registration.test.client_id

  microsoft_app_type      = "UserAssignedMSI"
  microsoft_app_tenant_id = data.azurerm_client_config.current.tenant_id
  microsoft_app_msi_id    = azurerm_user_assigned_identity.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (BotServiceAzureBotResource) steamingEndpointEnabled(data acceptance.TestData, streamingEndpointEnabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azuread_application_registration" "test" {
  display_name = "acctestReg-%[1]d"
}

resource "azurerm_bot_service_azure_bot" "test" {
  name                       = "acctestdf%[1]d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = "global"
  sku                        = "F0"
  microsoft_app_id           = azuread_application_registration.test.client_id
  streaming_endpoint_enabled = %[3]t
}
`, data.RandomInteger, data.Locations.Primary, streamingEndpointEnabled)
}

func (BotServiceAzureBotResource) cmekEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy    = true
      recover_soft_deleted_key_vaults = true
    }
  }
}

data "azurerm_client_config" "current" {}

data "azuread_service_principal" "test" {
  display_name = "Bot Service CMEK Prod"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_key_vault" "test" {
  name                        = "acctestKV-%[1]d"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  enabled_for_disk_encryption = true
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  soft_delete_retention_days  = 7
  purge_protection_enabled    = true
  enable_rbac_authorization   = true
  sku_name                    = "standard"
}

resource "azurerm_role_assignment" "test_deployer" {
  scope                = azurerm_key_vault.test.id
  role_definition_name = "Key Vault Administrator"
  principal_id         = data.azurerm_client_config.current.object_id
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_key_vault.test.id
  role_definition_name = "Key Vault Crypto Service Encryption User"
  principal_id         = data.azuread_service_principal.test.object_id
}

resource "azurerm_key_vault_key" "test" {
  name         = "acctestKey-%[1]d"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  depends_on = [azurerm_role_assignment.test_deployer]
}

resource "azuread_application_registration" "test" {
  display_name = "acctestReg-%[1]d"
}

resource "azurerm_bot_service_azure_bot" "test" {
  name                  = "acctestdf%[1]d"
  resource_group_name   = azurerm_resource_group.test.name
  location              = "global"
  sku                   = "F0"
  microsoft_app_id      = azuread_application_registration.test.client_id
  cmk_key_vault_key_url = azurerm_key_vault_key.test.id
  endpoint              = "https://example2.com"

  depends_on = [azurerm_role_assignment.test]
}
`, data.RandomIntOfLength(8), data.Locations.Primary)
}
