// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2015-10-31/webhook"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AutomationWebhookResource struct{}

func TestAccAutomationWebhook_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_webhook", "test")
	r := AutomationWebhookResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("uri"),
	})
}

func TestAccAutomationWebhook_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_webhook", "test")
	r := AutomationWebhookResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("uri"),
	})
}

func TestAccAutomationWebhook_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_webhook", "test")
	r := AutomationWebhookResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("uri"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("uri"),
	})
}

func TestAccAutomationWebhook_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_webhook", "test")
	r := AutomationWebhookResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_automation_webhook"),
		},
	})
}

func (t AutomationWebhookResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := webhook.ParseWebHookID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Automation.WebhookClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (AutomationWebhookResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_runbook" "test" {
  name                    = "Get-AzureVMTutorial"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name

  log_verbose  = "true"
  log_progress = "true"
  description  = "This is a test runbook for terraform acceptance test"
  runbook_type = "PowerShell"

  content = <<CONTENT
# Some test content
# for Terraform acceptance test
CONTENT
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (AutomationWebhookResource) basic(data acceptance.TestData) string {
	template := AutomationWebhookResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_webhook" "test" {
  name                    = "TestRunbook_webhook"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  expiry_time             = "%s"
  runbook_name            = azurerm_automation_runbook.test.name
}
`, template, time.Now().UTC().Add(time.Hour).Format(time.RFC3339))
}

func (AutomationWebhookResource) complete(data acceptance.TestData) string {
	template := AutomationWebhookResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_credential" "test" {
  name                    = "acctest-%[2]d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  username                = "test_user"
  password                = "test_pwd"
}

resource "azurerm_automation_hybrid_runbook_worker_group" "test" {
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  name                    = "acctest-%[2]d"
  credential_name         = azurerm_automation_credential.test.name
}

resource "azurerm_automation_webhook" "test" {
  name                    = "TestRunbook_webhook"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  expiry_time             = "%[3]s"
  enabled                 = false
  runbook_name            = azurerm_automation_runbook.test.name
  run_on_worker_group     = azurerm_automation_hybrid_runbook_worker_group.test.name
  parameters = {
    input = "parameter"
  }
}
`, template, data.RandomInteger, time.Now().UTC().Add(time.Hour).Format(time.RFC3339))
}

func (AutomationWebhookResource) requiresImport(data acceptance.TestData) string {
	template := AutomationWebhookResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_webhook" "import" {
  name                    = azurerm_automation_webhook.test.name
  resource_group_name     = azurerm_automation_webhook.test.resource_group_name
  automation_account_name = azurerm_automation_webhook.test.automation_account_name
  expiry_time             = azurerm_automation_webhook.test.expiry_time
  runbook_name            = azurerm_automation_webhook.test.runbook_name
}
`, template)
}
