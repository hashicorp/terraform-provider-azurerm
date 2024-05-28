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

type BotWebAppResource struct{}

func TestAccBotWebApp_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_web_app", "test")
	r := BotWebAppResource{}

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

func TestAccBotWebApp_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_web_app", "test")
	r := BotWebAppResource{}

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
	})
}

func TestAccBotWebApp_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_web_app", "test")
	r := BotWebAppResource{}

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

func (t BotWebAppResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (BotWebAppResource) basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
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

resource "azurerm_bot_web_app" "test" {
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

func (BotWebAppResource) updateConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
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

resource "azurerm_bot_web_app" "test" {
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

func (BotWebAppResource) completeConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
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

resource "azuread_application_registration" "test" {
  display_name = "acctestReg-%d"
}

resource "azurerm_bot_web_app" "test" {
  name                = "acctestdf%d"
  location            = "global"
  resource_group_name = azurerm_resource_group.test.name
  microsoft_app_id    = azuread_application_registration.test.client_id
  sku                 = "F0"

  endpoint                              = "https://example.com"
  developer_app_insights_api_key        = azurerm_application_insights_api_key.test.api_key
  developer_app_insights_application_id = azurerm_application_insights.test.app_id
  developer_app_insights_key            = azurerm_application_insights.test.instrumentation_key

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
