// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package costmanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2023-08-01/scheduledactions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AnomalyAlertResource struct{}

func TestAccResourceAnomalyAlert_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cost_anomaly_alert", "test")
	r := AnomalyAlertResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceAnomalyAlert_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cost_anomaly_alert", "test")
	r := AnomalyAlertResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceAnomalyAlert_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cost_anomaly_alert", "test")
	r := AnomalyAlertResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicForImport(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(),
			ExpectError: acceptance.RequiresImportError("azurerm_cost_anomaly_alert"),
		},
	})
}

func TestAccResourceAnomalyAlert_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cost_anomaly_alert", "test")
	r := AnomalyAlertResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.update(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update2(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
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

	return pointer.To(resp.Model.Properties != nil), nil
}

func (AnomalyAlertResource) basic() string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_cost_anomaly_alert" "test" {
  name            = "acctest-basic"
  display_name    = "Finance budget"
  email_subject   = "Hi"
  email_addresses = ["test@test.com", "test@hashicorp.developer"]
  message         = "Oops, cost anomaly"
}
`)
}

func (AnomalyAlertResource) basicForImport() string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_cost_anomaly_alert" "test" {
  name            = "acctest-import"
  display_name    = "Finance budget"
  email_subject   = "Hi"
  email_addresses = ["test@test.com", "test@hashicorp.developer"]
  message         = "Oops, cost anomaly"
}
`)
}

func (AnomalyAlertResource) complete() string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "test" {}

resource "azurerm_cost_anomaly_alert" "test" {
  name               = "acctest-complete"
  display_name       = "Finance budget"
  subscription_id    = data.azurerm_subscription.test.id
  email_subject      = "Hi"
  email_addresses    = ["test@test.com", "test@hashicorp.developer"]
  notification_email = "othertest@hashicorp.developer"
  message            = "Cost anomaly complete test"
}
`)
}

func (r AnomalyAlertResource) requiresImport() string {
	return fmt.Sprintf(`
%s

resource "azurerm_cost_anomaly_alert" "import" {
  name            = azurerm_cost_anomaly_alert.test.name
  display_name    = azurerm_cost_anomaly_alert.test.display_name
  email_subject   = azurerm_cost_anomaly_alert.test.email_subject
  email_addresses = azurerm_cost_anomaly_alert.test.email_addresses
  message         = azurerm_cost_anomaly_alert.test.message
}
`, r.basicForImport())
}

func (AnomalyAlertResource) update() string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_cost_anomaly_alert" "test" {
  name            = "acctest-update"
  display_name    = "Budget"
  email_subject   = "Hi you!"
  email_addresses = ["tester@test.com", "test2@hashicorp.developer"]
  message         = "A cost anomaly for you"
}
`)
}

func (AnomalyAlertResource) update2() string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "test" {}

resource "azurerm_cost_anomaly_alert" "test" {
  name            = "acctest-update"
  display_name    = "Finance budget"
  email_subject   = "Hello you!"
  email_addresses = ["tester@test.com", "test2@hashicorp.developer", "test@test.com"]
  message         = "An updated cost anomaly for you"
  subscription_id = data.azurerm_subscription.test.id
}
`)
}
