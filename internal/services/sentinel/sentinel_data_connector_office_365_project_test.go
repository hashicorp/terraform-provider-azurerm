// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SentinelDataConnectorOffice365ProjectResource struct{}

func TestAccSentinelDataConnectorOffice365Project_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_data_connector_office_365_project", "test")
	r := SentinelDataConnectorOffice365ProjectResource{}

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

func TestAccSentinelDataConnectorOffice365Project_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_data_connector_office_365_project", "test")
	r := SentinelDataConnectorOffice365ProjectResource{}

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

func TestAccSentinelDataConnectorOffice365Project_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_data_connector_office_365_project", "test")
	r := SentinelDataConnectorOffice365ProjectResource{}

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

func (r SentinelDataConnectorOffice365ProjectResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	client := clients.Sentinel.DataConnectorsClient

	id, err := parse.DataConnectorID(state.ID)
	if err != nil {
		return nil, err
	}

	if resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name); err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r SentinelDataConnectorOffice365ProjectResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_data_connector_office_365_project" "test" {
  name                       = "accTestDC-%d"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
}
`, template, data.RandomInteger)
}

func (r SentinelDataConnectorOffice365ProjectResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "test" {}

resource "azurerm_sentinel_data_connector_office_365_project" "test" {
  name                       = "accTestDC-%d"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  tenant_id                  = data.azurerm_client_config.test.tenant_id
}
`, template, data.RandomInteger)
}

func (r SentinelDataConnectorOffice365ProjectResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_data_connector_office_365_project" "import" {
  name                       = azurerm_sentinel_data_connector_office_365_project.test.name
  log_analytics_workspace_id = azurerm_sentinel_data_connector_office_365_project.test.log_analytics_workspace_id
}
`, template)
}

func (r SentinelDataConnectorOffice365ProjectResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sentinel-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_sentinel_log_analytics_workspace_onboarding" "test" {
  workspace_id = azurerm_log_analytics_workspace.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
