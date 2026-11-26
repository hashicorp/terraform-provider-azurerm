// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package bot_channel_web_chat_test

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type BotChannelsRegistrationResource struct{}

func (r BotChannelsRegistrationResource) basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_bot_channels_registration" "test" {
  name                    = "acctestdf%[2]d"
  location                = "global"
  resource_group_name     = azurerm_resource_group.test.name
  sku                     = "F0"
  microsoft_app_id        = azuread_application_registration.test.client_id
  microsoft_app_type      = "SingleTenant"
  microsoft_app_tenant_id = data.azurerm_client_config.current.tenant_id
}
`, r.template(data), data.RandomInteger)
}

func (BotChannelsRegistrationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azuread_application_registration" "test" {
  display_name = "acctestReg-%[1]d"
}
`, data.RandomInteger, data.Locations.Primary)
}
