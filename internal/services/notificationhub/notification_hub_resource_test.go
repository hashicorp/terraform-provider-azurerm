// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package notificationhub_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/notificationhubs/2023-09-01/hubs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NotificationHubResource struct{}

func TestAccNotificationHub_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_notification_hub", "test")
	r := NotificationHubResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("apns_credential.#").HasValue("0"),
				check.That(data.ResourceName).Key("gcm_credential.#").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNotificationHub_browserCredential(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_notification_hub", "test")
	r := NotificationHubResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.browserCredential(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNotificationHub_updateTag(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_notification_hub", "test")
	r := NotificationHubResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withoutTag(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNotificationHub_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_notification_hub", "test")
	r := NotificationHubResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("apns_credential.#").HasValue("0"),
				check.That(data.ResourceName).Key("gcm_credential.#").HasValue("0"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (NotificationHubResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := hubs.ParseNotificationHubID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.NotificationHubs.HubsClient.NotificationHubsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (NotificationHubResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRGpol-%d"
  location = "%s"
}

resource "azurerm_notification_hub_namespace" "test" {
  name                = "acctestnhn-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  namespace_type      = "NotificationHub"
  sku_name            = "Free"
}

resource "azurerm_notification_hub" "test" {
  name                = "acctestnh-%d"
  namespace_name      = azurerm_notification_hub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  tags = {
    env = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (NotificationHubResource) browserCredential(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRGpol-%d"
  location = "%s"
}

resource "azurerm_notification_hub_namespace" "test" {
  name                = "acctestnhn-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  namespace_type      = "NotificationHub"
  sku_name            = "Free"
}

resource "azurerm_notification_hub" "test" {
  name                = "acctestnh-%d"
  namespace_name      = azurerm_notification_hub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  browser_credential {
    subject           = "testSubject"
    vapid_private_key = "X4X_Awjb4HyD70adCrw6FmFgA4wiu_TTWSZFcayBN6U"
    vapid_public_key  = "BC1XlIUxB6kQ2a214VqTMT4hnX44LRnhWDaiNxEi5bRtkdE5bFkRClX6gunX4_YWIn0UY8TD20gBGqvOg6T-go4"
  }

  tags = {
    env = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (NotificationHubResource) withoutTag(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRGpol-%d"
  location = "%s"
}

resource "azurerm_notification_hub_namespace" "test" {
  name                = "acctestnhn-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  namespace_type      = "NotificationHub"
  sku_name            = "Free"
}

resource "azurerm_notification_hub" "test" {
  name                = "acctestnh-%d"
  namespace_name      = azurerm_notification_hub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (NotificationHubResource) requiresImport(data acceptance.TestData) string {
	template := NotificationHubResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_notification_hub" "import" {
  name                = azurerm_notification_hub.test.name
  namespace_name      = azurerm_notification_hub.test.namespace_name
  resource_group_name = azurerm_notification_hub.test.resource_group_name
  location            = azurerm_notification_hub.test.location
}
`, template)
}
