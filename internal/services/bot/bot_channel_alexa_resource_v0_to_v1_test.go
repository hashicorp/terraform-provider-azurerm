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

// TestAccBotChannelAlexa_V0ToV1_501 tests the state migration path from an `id` with lowercased
// static segments to their canonicalized format. It uses v5.0.1 as the setup version because it is
// the last release where the `id` could have been stored in state with lowercased static segments via an import using a non-canonical ID.
func TestAccBotChannelAlexa_V0ToV1_501(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_alexa", "test")
	r := BotChannelAlexaResource{}

	importedResourceName := data.ResourceName + "-import"

	data.ResourceRegressionAdditionalStepsTest(t, r, []acceptance.TestStep{
		{
			// the existing basic config generates a new random skill_id per invocation, so a variant
			// with a fixed skill_id is used to keep the config stable across the test steps
			Config: r.basicV0Setup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acctestRG-%[2]d/providers/Microsoft.BotService/botServices/acctestdf%[2]d/channels/AlexaChannel", data.Subscriptions.Primary, data.RandomInteger)),
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
				check.That(importedResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourcegroups/acctestRG-%[2]d/providers/microsoft.botservice/botServices/acctestdf%[2]d/channels/AlexaChannel", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
		{
			Config: r.basicImported(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(importedResourceName).ExistsInAzure(r),
				check.That(importedResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acctestRG-%[2]d/providers/Microsoft.BotService/botServices/acctestdf%[2]d/channels/AlexaChannel", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
	}, "5.0.1")
}

func (BotChannelAlexaResource) basicV0Setup(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_alexa" "test" {
  bot_name            = azurerm_bot_channels_registration.test.name
  location            = azurerm_bot_channels_registration.test.location
  resource_group_name = azurerm_resource_group.test.name
  skill_id            = "amzn1.ask.skill.34b7d17b-90d6-4b0a-8d69-b3d21b8dd3e8"
}
`, BotChannelsRegistrationResource{}.basicConfig(data))
}

func (BotChannelAlexaResource) basicV0(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_bot_channel_alexa" "test-import" {
  bot_name            = azurerm_bot_channels_registration.test.name
  location            = azurerm_bot_channels_registration.test.location
  resource_group_name = azurerm_resource_group.test.name
  skill_id            = "amzn1.ask.skill.34b7d17b-90d6-4b0a-8d69-b3d21b8dd3e8"
}

removed {
  from = azurerm_bot_channel_alexa.test
  lifecycle {
    destroy = false
  }
}

import {
  to = azurerm_bot_channel_alexa.test-import
  id = "/subscriptions/%[3]s/resourcegroups/acctestRG-%[2]d/providers/microsoft.botservice/botServices/acctestdf%[2]d/channels/AlexaChannel"
}
`, BotChannelsRegistrationResource{}.basicConfig(data), data.RandomInteger, data.Subscriptions.Primary)
}

func (BotChannelAlexaResource) basicImported(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_alexa" "test-import" {
  bot_name            = azurerm_bot_channels_registration.test.name
  location            = azurerm_bot_channels_registration.test.location
  resource_group_name = azurerm_resource_group.test.name
  skill_id            = "amzn1.ask.skill.34b7d17b-90d6-4b0a-8d69-b3d21b8dd3e8"
}
`, BotChannelsRegistrationResource{}.basicConfig(data))
}
