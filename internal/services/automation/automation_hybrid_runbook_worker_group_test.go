// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/hybridrunbookworkergroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type HybridRunbookWorkerGroupResource struct{}

func (a HybridRunbookWorkerGroupResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := hybridrunbookworkergroup.ParseHybridRunbookWorkerGroupID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Automation.HybridRunbookWorkerGroup.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving HybridRunbookWorkerGroup %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (a HybridRunbookWorkerGroupResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%[1]d"
  location = "%[2]s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestAA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_credential" "test" {
  name                    = "acctest-%[1]d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  username                = "test_user"
  password                = "test_pwd"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (a HybridRunbookWorkerGroupResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_automation_hybrid_runbook_worker_group" "test" {
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  name                    = "acctest-%[2]d"
  credential_name         = azurerm_automation_credential.test.name
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a HybridRunbookWorkerGroupResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_automation_credential" "test2" {
  name                    = "acctest2-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  username                = "test_user"
  password                = "test_pwd"
}

resource "azurerm_automation_hybrid_runbook_worker_group" "test" {
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  name                    = "acctest-%[2]d"
  credential_name         = azurerm_automation_credential.test2.name
}
`, a.template(data), data.RandomInteger)
}

func TestAccHybridRunbookWorkerGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, automation.HybridRunbookWorkerGroupResource{}.ResourceType(), "test")
	r := HybridRunbookWorkerGroupResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccHybridRunbookWorkerGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, automation.HybridRunbookWorkerGroupResource{}.ResourceType(), "test")
	r := HybridRunbookWorkerGroupResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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
