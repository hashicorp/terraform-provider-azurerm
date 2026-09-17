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

// TestAccBotChannelFacebook_V0ToV1_501 tests the state migration path from an `id` with lowercased
// static segments to their canonicalized format. It uses v5.0.1 as the setup version because it is
// the last release where the `id` could have been stored in state with lowercased static segments via an import using a non-canonical ID.
func TestAccBotChannelFacebook_V0ToV1_501(t *testing.T) {
	skipFacebookChannel(t)

	data := acceptance.BuildTestData(t, "azurerm_bot_channel_facebook", "test")
	r := BotChannelFacebookResource{}

	importedResourceName := data.ResourceName + "-import"

	data.ResourceRegressionAdditionalStepsTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acctestRG-%[2]d/providers/Microsoft.BotService/botServices/acctestdf%[2]d/channels/FacebookChannel", data.Subscriptions.Primary, data.RandomInteger)),
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
				check.That(importedResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourcegroups/acctestRG-%[2]d/providers/microsoft.botservice/botServices/acctestdf%[2]d/channels/FacebookChannel", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
		{
			Config: r.basicImported(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(importedResourceName).ExistsInAzure(r),
				check.That(importedResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acctestRG-%[2]d/providers/Microsoft.BotService/botServices/acctestdf%[2]d/channels/FacebookChannel", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
	}, "5.0.1")
}

func (BotChannelFacebookResource) basicV0(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_bot_channel_facebook" "test-import" {
  bot_name                    = azurerm_bot_channels_registration.test.name
  location                    = azurerm_bot_channels_registration.test.location
  resource_group_name         = azurerm_resource_group.test.name
  facebook_application_id     = "%[4]s"
  facebook_application_secret = "%[5]s"

  page {
    id           = "%[6]s"
    access_token = "%[7]s"
  }
}

removed {
  from = azurerm_bot_channel_facebook.test
  lifecycle {
    destroy = false
  }
}

import {
  to = azurerm_bot_channel_facebook.test-import
  id = "/subscriptions/%[3]s/resourcegroups/acctestRG-%[2]d/providers/microsoft.botservice/botServices/acctestdf%[2]d/channels/FacebookChannel"
}
`, BotChannelsRegistrationResource{}.basicConfig(data), data.RandomInteger, data.Subscriptions.Primary, os.Getenv("ARM_TEST_FACEBOOK_APPLICATION_ID"), os.Getenv("ARM_TEST_FACEBOOK_APPLICATION_SECRET"), os.Getenv("ARM_TEST_PAGE_ID"), os.Getenv("ARM_TEST_PAGE_ACCESS_TOKEN"))
}

func (BotChannelFacebookResource) basicImported(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_bot_channel_facebook" "test-import" {
  bot_name                    = azurerm_bot_channels_registration.test.name
  location                    = azurerm_bot_channels_registration.test.location
  resource_group_name         = azurerm_resource_group.test.name
  facebook_application_id     = "%[2]s"
  facebook_application_secret = "%[3]s"

  page {
    id           = "%[4]s"
    access_token = "%[5]s"
  }
}
`, BotChannelsRegistrationResource{}.basicConfig(data), os.Getenv("ARM_TEST_FACEBOOK_APPLICATION_ID"), os.Getenv("ARM_TEST_FACEBOOK_APPLICATION_SECRET"), os.Getenv("ARM_TEST_PAGE_ID"), os.Getenv("ARM_TEST_PAGE_ACCESS_TOKEN"))
}
