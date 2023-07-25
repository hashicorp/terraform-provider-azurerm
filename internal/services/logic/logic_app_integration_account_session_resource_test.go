// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountsessions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LogicAppIntegrationAccountSessionResource struct{}

func TestAccLogicAppIntegrationAccountSession_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_session", "test")
	r := LogicAppIntegrationAccountSessionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "1234"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppIntegrationAccountSession_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_session", "test")
	r := LogicAppIntegrationAccountSessionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "1234"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data, "1234"),
			ExpectError: acceptance.RequiresImportError(data.ResourceType),
		},
	})
}

func TestAccLogicAppIntegrationAccountSession_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_session", "test")
	r := LogicAppIntegrationAccountSessionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "1234"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, "5678"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (LogicAppIntegrationAccountSessionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := integrationaccountsessions.ParseSessionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Logic.IntegrationAccountSessionClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %q: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r LogicAppIntegrationAccountSessionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-logic-%d"
  location = "%s"
}

resource "azurerm_logic_app_integration_account" "test" {
  name                = "acctest-ia-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r LogicAppIntegrationAccountSessionResource) basic(data acceptance.TestData, controlNumber string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account_session" "test" {
  name                     = "acctest-ias-%d"
  resource_group_name      = azurerm_resource_group.test.name
  integration_account_name = azurerm_logic_app_integration_account.test.name

  content = <<CONTENT
	{
       "controlNumber": "%s"
    }
  CONTENT
}
`, r.template(data), data.RandomInteger, controlNumber)
}

func (r LogicAppIntegrationAccountSessionResource) requiresImport(data acceptance.TestData, controlNumber string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account_session" "import" {
  name                     = azurerm_logic_app_integration_account_session.test.name
  resource_group_name      = azurerm_logic_app_integration_account_session.test.resource_group_name
  integration_account_name = azurerm_logic_app_integration_account_session.test.integration_account_name
  content                  = azurerm_logic_app_integration_account_session.test.content
}
`, r.basic(data, controlNumber))
}
