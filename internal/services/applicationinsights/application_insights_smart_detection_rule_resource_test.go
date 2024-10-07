// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package applicationinsights_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	smartdetection "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/componentproactivedetectionapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AppInsightsSmartDetectionRule struct{}

func TestAccApplicationInsightsSmartDetectionRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_smart_detection_rule", "test")
	r := AppInsightsSmartDetectionRule{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationInsightsSmartDetectionRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_smart_detection_rule", "test")
	r := AppInsightsSmartDetectionRule{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccApplicationInsightsSmartDetectionRule_multiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_smart_detection_rule", "test")
	r := AppInsightsSmartDetectionRule{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That("azurerm_application_insights_smart_detection_rule.test2").ExistsInAzure(r),
				check.That("azurerm_application_insights_smart_detection_rule.test3").ExistsInAzure(r),
				check.That("azurerm_application_insights_smart_detection_rule.test4").ExistsInAzure(r),
				check.That("azurerm_application_insights_smart_detection_rule.test5").ExistsInAzure(r),
				check.That("azurerm_application_insights_smart_detection_rule.test6").ExistsInAzure(r),
				check.That("azurerm_application_insights_smart_detection_rule.test7").ExistsInAzure(r),
				check.That("azurerm_application_insights_smart_detection_rule.test8").ExistsInAzure(r),
				check.That("azurerm_application_insights_smart_detection_rule.test9").ExistsInAzure(r),
				check.That("azurerm_application_insights_smart_detection_rule.test10").ExistsInAzure(r),
			),
		},
	})
}

func TestAccApplicationInsightsSmartDetectionRule_longDependencyDuration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_smart_detection_rule", "test")
	r := AppInsightsSmartDetectionRule{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.longDependencyDuration(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

// A requires import test isn't possible here due to the behaviour of app insights. When a new app insights instance is
// created all the smart detection rules are created with it and are set to enabled. They cannot be deleted, only disabled -
// but this still causes issues when the resource performs a d.IsNewResource() check, where the tf.ImportAsExistsError is thrown.

func (t AppInsightsSmartDetectionRule) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := smartdetection.ParseProactiveDetectionConfigID(state.Attributes["id"])
	if err != nil {
		return nil, err
	}

	resp, err := clients.AppInsights.SmartDetectionRuleClient.ProactiveDetectionConfigurationsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (AppInsightsSmartDetectionRule) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
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

resource "azurerm_application_insights_smart_detection_rule" "test" {
  name                    = "Slow page load time"
  application_insights_id = azurerm_application_insights.test.id
  enabled                 = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (AppInsightsSmartDetectionRule) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
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

resource "azurerm_application_insights_smart_detection_rule" "test" {
  name                    = "Slow page load time"
  application_insights_id = azurerm_application_insights.test.id
  enabled                 = false

  send_emails_to_subscription_owners = false
  additional_email_recipients        = ["test@example.com", "test2@example.com"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (AppInsightsSmartDetectionRule) multiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
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

resource "azurerm_application_insights_smart_detection_rule" "test" {
  name                    = "Slow page load time"
  application_insights_id = azurerm_application_insights.test.id
  enabled                 = false
}

resource "azurerm_application_insights_smart_detection_rule" "test2" {
  name                    = "Slow server response time"
  application_insights_id = azurerm_application_insights.test.id
  enabled                 = false
}

resource "azurerm_application_insights_smart_detection_rule" "test3" {
  name                    = "Long dependency duration"
  application_insights_id = azurerm_application_insights.test.id
  enabled                 = false
}

resource "azurerm_application_insights_smart_detection_rule" "test4" {
  name                    = "Degradation in server response time"
  application_insights_id = azurerm_application_insights.test.id
  enabled                 = false
}

resource "azurerm_application_insights_smart_detection_rule" "test5" {
  name                    = "Degradation in dependency duration"
  application_insights_id = azurerm_application_insights.test.id
  enabled                 = false
}

resource "azurerm_application_insights_smart_detection_rule" "test6" {
  name                    = "Degradation in trace severity ratio"
  application_insights_id = azurerm_application_insights.test.id
  enabled                 = false
}

resource "azurerm_application_insights_smart_detection_rule" "test7" {
  name                    = "Abnormal rise in exception volume"
  application_insights_id = azurerm_application_insights.test.id
  enabled                 = false
}

resource "azurerm_application_insights_smart_detection_rule" "test8" {
  name                    = "Potential memory leak detected"
  application_insights_id = azurerm_application_insights.test.id
  enabled                 = false
}

resource "azurerm_application_insights_smart_detection_rule" "test9" {
  name                    = "Potential security issue detected"
  application_insights_id = azurerm_application_insights.test.id
  enabled                 = false
}

resource "azurerm_application_insights_smart_detection_rule" "test10" {
  name                    = "Abnormal rise in daily data volume"
  application_insights_id = azurerm_application_insights.test.id
  enabled                 = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (AppInsightsSmartDetectionRule) longDependencyDuration(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
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

resource "azurerm_application_insights_smart_detection_rule" "test" {
  name                    = "Long dependency duration"
  application_insights_id = azurerm_application_insights.test.id
  enabled                 = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
