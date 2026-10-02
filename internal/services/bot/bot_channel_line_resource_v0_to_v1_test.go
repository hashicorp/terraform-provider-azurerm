// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

// TestAccBotChannelLine_V0ToV1_501 tests the state migration path from an `id` with lowercased
// static segments to their canonicalized format. It uses v5.0.1 as the setup version because it is
// the last release where the `id` could have been stored in state with lowercased static segments via an import using a non-canonical ID.
func TestAccBotChannelLine_V0ToV1_501(t *testing.T) {
	skipLineChannel(t)

	data := acceptance.BuildTestData(t, "azurerm_bot_channel_line", "test")
	r := BotChannelLineResource{}

	importedResourceName := data.ResourceName + "-import"

	data.ResourceRegressionAdditionalStepsTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acctestRG-%[2]d/providers/Microsoft.BotService/botServices/acctestdf%[2]d/channels/LineChannel", data.Subscriptions.Primary, data.RandomInteger)),
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
				check.That(importedResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourcegroups/acctestRG-%[2]d/providers/microsoft.botservice/botServices/acctestdf%[2]d/channels/LineChannel", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
		{
			Config: r.basicImported(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(importedResourceName).ExistsInAzure(r),
				check.That(importedResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acctestRG-%[2]d/providers/Microsoft.BotService/botServices/acctestdf%[2]d/channels/LineChannel", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
	}, "5.0.1")
}

func (BotChannelLineResource) basicV0(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_bot_channel_line" "test-import" {
  bot_name            = azurerm_bot_channels_registration.test.name
  location            = azurerm_bot_channels_registration.test.location
  resource_group_name = azurerm_resource_group.test.name

  line_channel {
    access_token = "%[4]s"
    secret       = "%[5]s"
  }
}

removed {
  from = azurerm_bot_channel_line.test
  lifecycle {
    destroy = false
  }
}

import {
  to = azurerm_bot_channel_line.test-import
  id = "/subscriptions/%[3]s/resourcegroups/acctestRG-%[2]d/providers/microsoft.botservice/botServices/acctestdf%[2]d/channels/LineChannel"
}
`, BotChannelsRegistrationResource{}.basicConfig(data), data.RandomInteger, data.Subscriptions.Primary, os.Getenv("ARM_TEST_CHANNEL_ACCESS_TOKEN"), os.Getenv("ARM_TEST_CHANNEL_SECRET"))
}

func (BotChannelLineResource) basicImported(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_bot_channel_line" "test-import" {
  bot_name            = azurerm_bot_channels_registration.test.name
  location            = azurerm_bot_channels_registration.test.location
  resource_group_name = azurerm_resource_group.test.name

  line_channel {
    access_token = "%[2]s"
    secret       = "%[3]s"
  }
}
`, BotChannelsRegistrationResource{}.basicConfig(data), os.Getenv("ARM_TEST_CHANNEL_ACCESS_TOKEN"), os.Getenv("ARM_TEST_CHANNEL_SECRET"))
}
