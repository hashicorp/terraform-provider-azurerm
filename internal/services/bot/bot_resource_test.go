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

type BotResource struct {
}

func TestAccBot_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot", "test")
	r := BotResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("azurebot"),
			),
		},
		data.ImportStep("developer_app_insights_api_key"),
	})
}

func TestAccBot_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot", "test")
	r := BotResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("azurebot"),
			),
		},
		data.ImportStep("developer_app_insights_api_key"),
		{
			Config: r.updateConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("sdk"),
			),
		},
		data.ImportStep("developer_app_insights_api_key"),
	})
}

func TestAccBot_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot", "test")
	r := BotResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("azurebot"),
			),
		},
		data.ImportStep("developer_app_insights_api_key"),
	})
}

func (t BotResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (BotResource) basicConfig(data acceptance.TestData) string {
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

resource "azurerm_bot" "test" {
  name                = "acctestdf%d"
  location            = "global"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "F0"
  microsoft_app_id    = data.azurerm_client_config.current.client_id

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (BotResource) updateConfig(data acceptance.TestData) string {
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

resource "azurerm_bot" "test" {
  name                = "acctestdf%d"
  location            = "global"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "F0"
  kind                = "sdk"
  microsoft_app_id    = data.azurerm_client_config.current.client_id

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (BotResource) completeConfig(data acceptance.TestData) string {
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

resource "azurerm_bot" "test" {
  name                = "acctestdf%d"
  location            = "global"
  resource_group_name = azurerm_resource_group.test.name
  microsoft_app_id    = data.azurerm_client_config.current.client_id
  sku                 = "F0"
  kind                = "azurebot"

  endpoint                              = "https://example.com"
  developer_app_insights_api_key        = azurerm_application_insights_api_key.test.api_key
  developer_app_insights_application_id = azurerm_application_insights.test.app_id

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
