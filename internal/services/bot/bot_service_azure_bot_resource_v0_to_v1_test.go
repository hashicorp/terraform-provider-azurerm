// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

// TestAccBotServiceAzureBot_V0ToV1_501 tests the state migration path from an `id` with lowercased
// static segments to their canonicalized format. It uses v5.0.1 as the setup version because it is
// the last release where the `id` could have been stored in state with lowercased static segments via an import using a non-canonical ID.
func TestAccBotServiceAzureBot_V0ToV1_501(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_service_azure_bot", "test")
	r := BotServiceAzureBotResource{}

	importedResourceName := data.ResourceName + "-import"

	data.ResourceRegressionAdditionalStepsTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acctestRG-%[2]d/providers/Microsoft.BotService/botServices/acctestdf%[2]d", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
		{
			Config: r.basicV0(data),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					// Ensure import is a no-op to prevent the resource from normalizing the ID due to a combined CreateUpdate func
					plancheck.ExpectResourceAction(importedResourceName, plancheck.ResourceActionNoop),
				},
			},
			Check: acceptance.ComposeTestCheckFunc(
				check.That(importedResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourcegroups/acctestRG-%[2]d/providers/microsoft.botservice/botServices/acctestdf%[2]d", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
		{
			Config: r.basicImported(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(importedResourceName).ExistsInAzure(r),
				check.That(importedResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acctestRG-%[2]d/providers/Microsoft.BotService/botServices/acctestdf%[2]d", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
	}, "5.0.1")
}

func (BotServiceAzureBotResource) v0MigrationTemplate(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary)
}

func (r BotServiceAzureBotResource) basicV0(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_bot_service_azure_bot" "test-import" {
  name                    = "acctestdf%[2]d"
  resource_group_name     = azurerm_resource_group.test.name
  location                = "global"
  sku                     = "F0"
  microsoft_app_id        = azuread_application_registration.test.client_id
  microsoft_app_type      = "SingleTenant"
  microsoft_app_tenant_id = data.azurerm_client_config.current.tenant_id

  tags = {
    environment = "test"
  }

  lifecycle {
    # ignore developer_app_insights_api_key, it's not returned on imports
    ignore_changes = [developer_app_insights_api_key]
  }
}

removed {
  from = azurerm_bot_service_azure_bot.test
  lifecycle {
    destroy = false
  }
}

import {
  to = azurerm_bot_service_azure_bot.test-import
  id = "/subscriptions/%[3]s/resourcegroups/acctestRG-%[2]d/providers/microsoft.botservice/botServices/acctestdf%[2]d"
}
`, r.v0MigrationTemplate(data), data.RandomInteger, data.Subscriptions.Primary)
}

func (r BotServiceAzureBotResource) basicImported(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_bot_service_azure_bot" "test-import" {
  name                    = "acctestdf%[2]d"
  resource_group_name     = azurerm_resource_group.test.name
  location                = "global"
  sku                     = "F0"
  microsoft_app_id        = azuread_application_registration.test.client_id
  microsoft_app_type      = "SingleTenant"
  microsoft_app_tenant_id = data.azurerm_client_config.current.tenant_id

  tags = {
    environment = "test"
  }
}
`, r.v0MigrationTemplate(data), data.RandomInteger)
}
