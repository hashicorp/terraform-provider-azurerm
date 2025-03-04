// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package costmanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2023-08-01/scheduledactions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AnomalyAlertResource struct{}

func TestAccResourceAnomalyAlert_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cost_anomaly_alert", "test")
	testResource := AnomalyAlertResource{}
	data.ResourceTest(t, testResource, []acceptance.TestStep{
		data.ApplyStep(testResource.basicConfig, testResource),
		data.ImportStep(),
		data.ApplyStep(testResource.updateConfig, testResource),
		data.ImportStep(),
		data.ApplyStep(testResource.basicConfig, testResource),
		data.ImportStep(),
	})
}

func TestAccResourceAnomalyAlert_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cost_anomaly_alert", "test")
	testResource := AnomalyAlertResource{}
	data.ResourceTest(t, testResource, []acceptance.TestStep{
		data.ApplyStep(testResource.completeConfig, testResource),
		data.ImportStep(),
		data.ApplyStep(testResource.updateConfig, testResource),
		data.ImportStep(),
	})
}

func TestAccResourceAnomalyAlert_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cost_anomaly_alert", "test")
	testResource := AnomalyAlertResource{}
	data.ResourceTest(t, testResource, []acceptance.TestStep{
		data.ApplyStep(testResource.basicConfig, testResource),
		data.RequiresImportErrorStep(testResource.requiresImportConfig),
	})
}

func TestAccResourceAnomalyAlert_emailAddressSender(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cost_anomaly_alert", "test")
	testResource := AnomalyAlertResource{}
	data.ResourceTest(t, testResource, []acceptance.TestStep{
		data.ApplyStep(testResource.notificationEmailConfig, testResource),
		data.ImportStep(),
		data.ApplyStep(testResource.updateConfig, testResource),
		data.ImportStep(),
	})
}

func (AnomalyAlertResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := scheduledactions.ParseScopedScheduledActionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.CostManagement.ScheduledActionsClient.GetByScope(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model.Properties != nil), nil
}

func (AnomalyAlertResource) basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_cost_anomaly_alert" "test" {
  name            = "-acctest-%d"
  display_name    = "acctest %d"
  email_subject   = "Hi"
  email_addresses = ["test@test.com", "test@hashicorp.developer"]
  message         = "Oops, cost anomaly"
}
`, data.RandomInteger, data.RandomInteger)
}

func (AnomalyAlertResource) completeConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "test" {}

resource "azurerm_cost_anomaly_alert" "test" {
  name            = "-acctest-%d"
  display_name    = "acctest %d"
  subscription_id = data.azurerm_subscription.test.id
  email_subject   = "Hi"
  email_addresses = ["test@test.com", "test@hashicorp.developer"]
  message         = "Cost anomaly complete test"
}
`, data.RandomInteger, data.RandomInteger)
}

func (r AnomalyAlertResource) requiresImportConfig(data acceptance.TestData) string {
	template := r.basicConfig(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cost_anomaly_alert" "import" {
  name            = azurerm_cost_anomaly_alert.test.name
  display_name    = azurerm_cost_anomaly_alert.test.display_name
  email_subject   = azurerm_cost_anomaly_alert.test.email_subject
  email_addresses = azurerm_cost_anomaly_alert.test.email_addresses
  message         = azurerm_cost_anomaly_alert.test.message
}
`, template)
}

func (AnomalyAlertResource) updateConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_cost_anomaly_alert" "test" {
  name            = "-acctest-%d"
  display_name    = "acctest name update %d"
  email_subject   = "Hi you!"
  email_addresses = ["tester@test.com", "test2@hashicorp.developer"]
  message         = "An updated cost anomaly for you"
}
`, data.RandomInteger, data.RandomInteger)
}

func (AnomalyAlertResource) notificationEmailConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_cost_anomaly_alert" "test" {
  name               = "-acctest-%d"
  display_name       = "acctest %d"
  email_subject      = "Hi"
  email_addresses    = ["test@test.com", "test@hashicorp.developer"]
  notification_email = "othertest@hashicorp.developer"
  message            = "Custom sender email configured"
}
`, data.RandomInteger, data.RandomInteger)
}
