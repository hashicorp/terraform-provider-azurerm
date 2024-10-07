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

type AutomationConnectionClassicCertificateResource struct{}

func TestAccAutomationConnectionClassicCertificate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_connection_classic_certificate", "test")
	r := AutomationConnectionClassicCertificateResource{}

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

func TestAccAutomationConnectionClassicCertificate_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_connection_classic_certificate", "test")
	r := AutomationConnectionClassicCertificateResource{}

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

func TestAccAutomationConnectionClassicCertificate_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_connection_classic_certificate", "test")
	r := AutomationConnectionClassicCertificateResource{}

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

func TestAccAutomationConnectionClassicCertificate_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_connection_classic_certificate", "test")
	r := AutomationConnectionClassicCertificateResource{}

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

func (t AutomationConnectionClassicCertificateResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (AutomationConnectionClassicCertificateResource) basic(data acceptance.TestData) string {
	template := AutomationConnectionClassicCertificateResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_connection_classic_certificate" "test" {
  name                    = "acctestACCC-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  certificate_asset_name  = "cert1"
  subscription_name       = "subs1"
  subscription_id         = data.azurerm_client_config.test.subscription_id
}
`, template, data.RandomInteger)
}

func (AutomationConnectionClassicCertificateResource) requiresImport(data acceptance.TestData) string {
	template := AutomationConnectionClassicCertificateResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_connection_classic_certificate" "import" {
  name                    = azurerm_automation_connection_classic_certificate.test.name
  resource_group_name     = azurerm_automation_connection_classic_certificate.test.resource_group_name
  automation_account_name = azurerm_automation_connection_classic_certificate.test.automation_account_name
  certificate_asset_name  = azurerm_automation_connection_classic_certificate.test.certificate_asset_name
  subscription_name       = azurerm_automation_connection_classic_certificate.test.subscription_name
  subscription_id         = azurerm_automation_connection_classic_certificate.test.subscription_id
}
`, template)
}

func (AutomationConnectionClassicCertificateResource) complete(data acceptance.TestData) string {
	template := AutomationConnectionClassicCertificateResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_connection_classic_certificate" "test" {
  name                    = "acctestACCC-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  certificate_asset_name  = "cert1"
  subscription_name       = "subs1"
  subscription_id         = data.azurerm_client_config.test.subscription_id
  description             = "acceptance test for automation connection"
}
`, template, data.RandomInteger)
}

func (AutomationConnectionClassicCertificateResource) template(data acceptance.TestData) string {
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
