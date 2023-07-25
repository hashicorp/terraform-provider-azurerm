// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package synapse_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SynapseSqlPoolSecurityAlertPolicyResource struct{}

func TestAccSynapseSqlPoolSecurityAlertPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool_security_alert_policy", "test")
	r := SynapseSqlPoolSecurityAlertPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_access_key"),
	})
}

func TestAccSynapseSqlPoolSecurityAlertPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool_security_alert_policy", "test")
	r := SynapseSqlPoolSecurityAlertPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_access_key"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_access_key"),
	})
}

func (SynapseSqlPoolSecurityAlertPolicyResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SqlPoolSecurityAlertPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Synapse.SqlPoolSecurityAlertPolicyClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.SqlPoolName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r SynapseSqlPoolSecurityAlertPolicyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_synapse_sql_pool_security_alert_policy" "test" {
  sql_pool_id                = azurerm_synapse_sql_pool.test.id
  policy_state               = "Enabled"
  storage_endpoint           = azurerm_storage_account.test.primary_blob_endpoint
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  retention_days             = 20

  disabled_alerts = [
    "Sql_Injection",
    "Data_Exfiltration"
  ]
}
`, r.template(data))
}

func (r SynapseSqlPoolSecurityAlertPolicyResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_synapse_sql_pool_security_alert_policy" "test" {
  sql_pool_id                  = azurerm_synapse_sql_pool.test.id
  policy_state                 = "Enabled"
  email_account_admins_enabled = true
  retention_days               = 30
}
`, r.template(data))
}

func (SynapseSqlPoolSecurityAlertPolicyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-synapse-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctest-%[1]d"
  storage_account_id = azurerm_storage_account.test.id
}

resource "azurerm_synapse_workspace" "test" {
  name                                 = "acctestsw%[1]d"
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_synapse_sql_pool" "test" {
  name                 = "acctestsp%[3]s"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  sku_name             = "DW100c"
  create_mode          = "Default"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
