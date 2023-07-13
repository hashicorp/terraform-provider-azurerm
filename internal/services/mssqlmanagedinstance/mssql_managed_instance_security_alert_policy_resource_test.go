// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssqlmanagedinstance_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MsSqlManagedInstanceSecurityAlertPolicyResource struct{}

func TestAccMssqlManagedInstanceSecurityAlertPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance_security_alert_policy", "test")
	r := MsSqlManagedInstanceSecurityAlertPolicyResource{}

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

func TestAccMssqlManagedInstanceSecurityAlertPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance_security_alert_policy", "test")
	r := MsSqlManagedInstanceSecurityAlertPolicyResource{}

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

func (MsSqlManagedInstanceSecurityAlertPolicyResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ManagedInstancesSecurityAlertPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.MSSQLManagedInstance.ManagedInstanceServerSecurityAlertPoliciesClient.Get(ctx, id.ResourceGroup, id.ManagedInstanceName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil, fmt.Errorf("SQL Managed Instance Security Alert Policy for server %q (Resource Group %q) does not exist", id.ManagedInstanceName, id.ResourceGroup)
		}
		return nil, fmt.Errorf("reading SQL Managed Instance Security Alert Policy for server %q (Resource Group %q): %v", id.ManagedInstanceName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r MsSqlManagedInstanceSecurityAlertPolicyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_managed_instance_security_alert_policy" "test" {
  resource_group_name          = azurerm_resource_group.test.name
  managed_instance_name        = azurerm_mssql_managed_instance.test.name
  enabled                      = true
  storage_endpoint             = azurerm_storage_account.test.primary_blob_endpoint
  storage_account_access_key   = azurerm_storage_account.test.primary_access_key
  retention_days               = 20
  email_account_admins_enabled = true
  email_addresses              = ["pearcec@example.com"]

  disabled_alerts = [
    "Sql_Injection",
    "Data_Exfiltration"
  ]
}
`, MsSqlManagedInstanceResource{}.basic(data), data.RandomInteger, data.Locations.Primary)
}

func (r MsSqlManagedInstanceSecurityAlertPolicyResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_managed_instance_security_alert_policy" "test" {
  resource_group_name        = azurerm_resource_group.test.name
  managed_instance_name      = azurerm_mssql_managed_instance.test.name
  enabled                    = true
  storage_endpoint           = azurerm_storage_account.test.primary_blob_endpoint
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  retention_days             = 30
}
`, MsSqlManagedInstanceResource{}.basic(data), data.RandomInteger, data.Locations.Primary)
}
