// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AutomationConnectionResource struct{}

func TestAccAutomationConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_connection", "test")
	r := AutomationConnectionResource{}

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

func TestAccAutomationConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_connection", "test")
	r := AutomationConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAutomationConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_connection", "test")
	r := AutomationConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutomationConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_connection", "test")
	r := AutomationConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t AutomationConnectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := connection.ParseConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Automation.Connection.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (AutomationConnectionResource) basic(data acceptance.TestData) string {
	template := AutomationConnectionResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_connection" "test" {
  name                    = "acctestAAC-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  type                    = "AzureServicePrincipal"

  values = {
    "ApplicationId" : "00000000-0000-0000-0000-000000000000"
    "TenantId" : data.azurerm_client_config.test.tenant_id
    "SubscriptionId" : data.azurerm_client_config.test.subscription_id
    "CertificateThumbprint" : file("testdata/automation_certificate_test.thumb")
  }
}
`, template, data.RandomInteger)
}

func (AutomationConnectionResource) requiresImport(data acceptance.TestData) string {
	template := AutomationConnectionResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_connection" "import" {
  name                    = azurerm_automation_connection.test.name
  resource_group_name     = azurerm_automation_connection.test.resource_group_name
  automation_account_name = azurerm_automation_connection.test.automation_account_name
  type                    = azurerm_automation_connection.test.type
  values                  = azurerm_automation_connection.test.values
}
`, template)
}

func (AutomationConnectionResource) complete(data acceptance.TestData) string {
	template := AutomationConnectionResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_connection" "test" {
  name                    = "acctestAAC-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  type                    = "AzureServicePrincipal"
  description             = "acceptance test for automation connection"

  values = {
    "ApplicationId" : "00000000-0000-0000-0000-000000000000"
    "TenantId" : data.azurerm_client_config.test.tenant_id
    "SubscriptionId" : data.azurerm_client_config.test.subscription_id
    "CertificateThumbprint" : file("testdata/automation_certificate_test.thumb")
  }
}
`, template, data.RandomInteger)
}

func (AutomationConnectionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%d"
  location = "%s"
}

data "azurerm_client_config" "test" {}

resource "azurerm_automation_account" "test" {
  name                = "acctestAA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
