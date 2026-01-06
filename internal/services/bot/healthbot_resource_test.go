// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthbot/2025-05-25/healthbots"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type HealthbotResource struct{}

func TestAccBotHealthbot_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthbot", "test")
	r := HealthbotResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, string(healthbots.SkuNameFZero)),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBotHealthbot_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthbot", "test")
	r := HealthbotResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, string(healthbots.SkuNameFZero)),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccBotHealthbot_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthbot", "test")
	r := HealthbotResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBotHealthbot_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthbot", "test")
	r := HealthbotResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, string(healthbots.SkuNameFZero)),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, string(healthbots.SkuNameCOne)),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r HealthbotResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := healthbots.ParseHealthBotID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Bot.HealthBotClient.HealthBots.BotsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r HealthbotResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-healthbot-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r HealthbotResource) basic(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_healthbot" "test" {
  name                = "acctest-hb-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "%[3]s"
}
`, r.template(data), data.RandomInteger, sku)
}

func (r HealthbotResource) basicForResourceIdentity(data acceptance.TestData) string {
	return r.basic(data, string(healthbots.SkuNameFZero))
}

func (r HealthbotResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_healthbot" "import" {
  name                = azurerm_healthbot.test.name
  resource_group_name = azurerm_healthbot.test.resource_group_name
  location            = azurerm_healthbot.test.location
  sku_name            = azurerm_healthbot.test.sku_name
}
`, r.basic(data, string(healthbots.SkuNameFZero)))
}

func (r HealthbotResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_healthbot" "test" {
  name                = "acctest-hb-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "C1"

  tags = {
    ENV = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}
