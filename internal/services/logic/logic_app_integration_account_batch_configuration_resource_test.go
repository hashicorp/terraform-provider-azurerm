// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountbatchconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LogicAppIntegrationAccountBatchConfigurationResource struct{}

func TestAccLogicAppIntegrationAccountBatchConfiguration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_batch_configuration", "test")
	r := LogicAppIntegrationAccountBatchConfigurationResource{}

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

func TestAccLogicAppIntegrationAccountBatchConfiguration_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_batch_configuration", "test")
	r := LogicAppIntegrationAccountBatchConfigurationResource{}

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

func TestAccLogicAppIntegrationAccountBatchConfiguration_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_batch_configuration", "test")
	r := LogicAppIntegrationAccountBatchConfigurationResource{}

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

func TestAccLogicAppIntegrationAccountBatchConfiguration_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_batch_configuration", "test")
	r := LogicAppIntegrationAccountBatchConfigurationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r LogicAppIntegrationAccountBatchConfigurationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := integrationaccountbatchconfigurations.ParseBatchConfigurationID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Logic.IntegrationAccountBatchConfigurationClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %q %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r LogicAppIntegrationAccountBatchConfigurationResource) template(data acceptance.TestData) string {
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
  sku_name            = "Standard"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r LogicAppIntegrationAccountBatchConfigurationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account_batch_configuration" "test" {
  name                     = "acctestiabc%s"
  resource_group_name      = azurerm_resource_group.test.name
  integration_account_name = azurerm_logic_app_integration_account.test.name
  batch_group_name         = "TestBatchGroup"

  release_criteria {
    message_count = 80
  }
}
`, r.template(data), data.RandomString)
}

func (r LogicAppIntegrationAccountBatchConfigurationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account_batch_configuration" "import" {
  name                     = azurerm_logic_app_integration_account_batch_configuration.test.name
  resource_group_name      = azurerm_logic_app_integration_account_batch_configuration.test.resource_group_name
  integration_account_name = azurerm_logic_app_integration_account_batch_configuration.test.integration_account_name
  batch_group_name         = azurerm_logic_app_integration_account_batch_configuration.test.batch_group_name

  release_criteria {
    message_count = azurerm_logic_app_integration_account_batch_configuration.test.release_criteria.0.message_count
  }
}
`, r.basic(data))
}

func (r LogicAppIntegrationAccountBatchConfigurationResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account_batch_configuration" "test" {
  name                     = "acctestiabc%s"
  resource_group_name      = azurerm_resource_group.test.name
  integration_account_name = azurerm_logic_app_integration_account.test.name
  batch_group_name         = "TestBatchGroup"

  release_criteria {
    batch_size    = 100
    message_count = 100

    recurrence {
      frequency  = "Month"
      interval   = 1
      start_time = "2021-09-01T01:00:00Z"
      end_time   = "2021-09-02T01:00:00Z"
      time_zone  = "Pacific Standard Time"

      schedule {
        hours      = [2, 3]
        minutes    = [4, 5]
        month_days = [6, 7]

        monthly {
          weekday = "Monday"
          week    = 1
        }
      }
    }
  }

  metadata = {
    foo = "bar"
  }
}
`, r.template(data), data.RandomString)
}

func (r LogicAppIntegrationAccountBatchConfigurationResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account_batch_configuration" "test" {
  name                     = "acctestiabc%s"
  resource_group_name      = azurerm_resource_group.test.name
  integration_account_name = azurerm_logic_app_integration_account.test.name
  batch_group_name         = "TestBatchGroup"

  release_criteria {
    batch_size    = 110
    message_count = 110

    recurrence {
      frequency  = "Week"
      interval   = 1
      start_time = "2021-09-02T01:00:00Z"
      end_time   = "2021-09-03T01:00:00Z"
      time_zone  = "Pacific SA Standard Time"

      schedule {
        hours     = [3, 4]
        minutes   = [5, 6]
        week_days = ["Monday", "Tuesday"]
      }
    }
  }

  metadata = {
    foo = "bar2"
  }
}
`, r.template(data), data.RandomString)
}
