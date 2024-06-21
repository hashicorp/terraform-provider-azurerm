// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package notificationhub_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/notificationhubs/2017-04-01/notificationhubs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NotificationHubAuthorizationRuleResource struct{}

func TestAccNotificationHubAuthorizationRule_listen(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_notification_hub_authorization_rule", "test")
	r := NotificationHubAuthorizationRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.listen(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("manage").HasValue("false"),
				check.That(data.ResourceName).Key("send").HasValue("false"),
				check.That(data.ResourceName).Key("listen").HasValue("true"),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNotificationHubAuthorizationRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_notification_hub_authorization_rule", "test")
	r := NotificationHubAuthorizationRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.listen(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("manage").HasValue("false"),
				check.That(data.ResourceName).Key("send").HasValue("false"),
				check.That(data.ResourceName).Key("listen").HasValue("true"),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccNotificationHubAuthorizationRule_manage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_notification_hub_authorization_rule", "test")
	r := NotificationHubAuthorizationRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.manage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("manage").HasValue("true"),
				check.That(data.ResourceName).Key("send").HasValue("true"),
				check.That(data.ResourceName).Key("listen").HasValue("true"),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNotificationHubAuthorizationRule_send(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_notification_hub_authorization_rule", "test")
	r := NotificationHubAuthorizationRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.send(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("manage").HasValue("false"),
				check.That(data.ResourceName).Key("send").HasValue("true"),
				check.That(data.ResourceName).Key("listen").HasValue("true"),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNotificationHubAuthorizationRule_multi(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_notification_hub_authorization_rule", "test1")
	r := NotificationHubAuthorizationRuleResource{}
	resourceTwoName := "azurerm_notification_hub_authorization_rule.test2"
	resourceThreeName := "azurerm_notification_hub_authorization_rule.test3"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multi(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("manage").HasValue("false"),
				check.That(data.ResourceName).Key("send").HasValue("true"),
				check.That(data.ResourceName).Key("listen").HasValue("true"),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
				check.That("azurerm_notification_hub_authorization_rule.test2").ExistsInAzure(r),
				check.That(resourceTwoName).Key("manage").HasValue("false"),
				check.That(resourceTwoName).Key("send").HasValue("true"),
				check.That(resourceTwoName).Key("listen").HasValue("true"),
				check.That(resourceTwoName).Key("primary_access_key").Exists(),
				check.That(resourceTwoName).Key("secondary_access_key").Exists(),
				check.That(resourceTwoName).Key("primary_connection_string").Exists(),
				check.That(resourceTwoName).Key("secondary_connection_string").Exists(),
				check.That("azurerm_notification_hub_authorization_rule.test3").ExistsInAzure(r),
				check.That(resourceThreeName).Key("manage").HasValue("false"),
				check.That(resourceThreeName).Key("send").HasValue("true"),
				check.That(resourceThreeName).Key("listen").HasValue("true"),
				check.That(resourceThreeName).Key("primary_access_key").Exists(),
				check.That(resourceThreeName).Key("secondary_access_key").Exists(),
				check.That(resourceThreeName).Key("primary_connection_string").Exists(),
				check.That(resourceThreeName).Key("secondary_connection_string").Exists(),
			),
		},
		data.ImportStep(),
		data.ImportStepFor(resourceTwoName),
		data.ImportStepFor(resourceThreeName),
	})
}

func TestAccNotificationHubAuthorizationRule_updated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_notification_hub_authorization_rule", "test")
	r := NotificationHubAuthorizationRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.listen(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("manage").HasValue("false"),
				check.That(data.ResourceName).Key("send").HasValue("false"),
				check.That(data.ResourceName).Key("listen").HasValue("true"),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
			),
		},
		{
			Config: r.manage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("manage").HasValue("true"),
				check.That(data.ResourceName).Key("send").HasValue("true"),
				check.That(data.ResourceName).Key("listen").HasValue("true"),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
			),
		},
	})
}

func (NotificationHubAuthorizationRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := notificationhubs.ParseNotificationHubAuthorizationRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.NotificationHubs.HubsClient.GetAuthorizationRule(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (NotificationHubAuthorizationRuleResource) listen(data acceptance.TestData) string {
	template := NotificationHubAuthorizationRuleResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_notification_hub_authorization_rule" "test" {
  name                  = "acctestrule-%d"
  notification_hub_name = azurerm_notification_hub.test.name
  namespace_name        = azurerm_notification_hub_namespace.test.name
  resource_group_name   = azurerm_resource_group.test.name
  listen                = true
}
`, template, data.RandomInteger)
}

func (NotificationHubAuthorizationRuleResource) requiresImport(data acceptance.TestData) string {
	template := NotificationHubAuthorizationRuleResource{}.listen(data)
	return fmt.Sprintf(`
%s

resource "azurerm_notification_hub_authorization_rule" "import" {
  name                  = azurerm_notification_hub_authorization_rule.test.name
  notification_hub_name = azurerm_notification_hub_authorization_rule.test.notification_hub_name
  namespace_name        = azurerm_notification_hub_authorization_rule.test.namespace_name
  resource_group_name   = azurerm_notification_hub_authorization_rule.test.resource_group_name
  listen                = azurerm_notification_hub_authorization_rule.test.listen
}
`, template)
}

func (NotificationHubAuthorizationRuleResource) send(data acceptance.TestData) string {
	template := NotificationHubAuthorizationRuleResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_notification_hub_authorization_rule" "test" {
  name                  = "acctestrule-%d"
  notification_hub_name = azurerm_notification_hub.test.name
  namespace_name        = azurerm_notification_hub_namespace.test.name
  resource_group_name   = azurerm_resource_group.test.name
  send                  = true
  listen                = true
}
`, template, data.RandomInteger)
}

func (NotificationHubAuthorizationRuleResource) multi(data acceptance.TestData) string {
	template := NotificationHubAuthorizationRuleResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_notification_hub_authorization_rule" "test1" {
  name                  = "acctestruleone-%d"
  notification_hub_name = azurerm_notification_hub.test.name
  namespace_name        = azurerm_notification_hub_namespace.test.name
  resource_group_name   = azurerm_resource_group.test.name
  send                  = true
  listen                = true
}

resource "azurerm_notification_hub_authorization_rule" "test2" {
  name                  = "acctestruletwo-%d"
  notification_hub_name = azurerm_notification_hub.test.name
  namespace_name        = azurerm_notification_hub_namespace.test.name
  resource_group_name   = azurerm_resource_group.test.name
  send                  = true
  listen                = true
}

resource "azurerm_notification_hub_authorization_rule" "test3" {
  name                  = "acctestrulethree-%d"
  notification_hub_name = azurerm_notification_hub.test.name
  namespace_name        = azurerm_notification_hub_namespace.test.name
  resource_group_name   = azurerm_resource_group.test.name
  send                  = true
  listen                = true
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (NotificationHubAuthorizationRuleResource) manage(data acceptance.TestData) string {
	template := NotificationHubAuthorizationRuleResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_notification_hub_authorization_rule" "test" {
  name                  = "acctestrule-%d"
  notification_hub_name = azurerm_notification_hub.test.name
  namespace_name        = azurerm_notification_hub_namespace.test.name
  resource_group_name   = azurerm_resource_group.test.name
  manage                = true
  send                  = true
  listen                = true
}
`, template, data.RandomInteger)
}

func (NotificationHubAuthorizationRuleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
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
