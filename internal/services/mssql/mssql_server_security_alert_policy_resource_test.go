// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MsSqlServerSecurityAlertPolicyResource struct{}

func TestAccMsSqlServerSecurityAlertPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server_security_alert_policy", "test")
	data = setTestDataLocationBySubscription(data)

	r := MsSqlServerSecurityAlertPolicyResource{}

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

func TestAccMsSqlServerSecurityAlertPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server_security_alert_policy", "test")
	data = setTestDataLocationBySubscription(data)

	r := MsSqlServerSecurityAlertPolicyResource{}

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

func (MsSqlServerSecurityAlertPolicyResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ServerSecurityAlertPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	serverId := commonids.NewSqlServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)

	result, err := client.MSSQL.ServerSecurityAlertPoliciesClient.Get(ctx, serverId)
	if err != nil {
		if response.WasNotFound(result.HttpResponse) {
			return nil, fmt.Errorf("SQL Security Alert Policy for server %q (Resource Group %q) does not exist", id.ServerName, id.ResourceGroup)
		}
		return nil, fmt.Errorf("reading SQL Security Alert Policy for server %q (Resource Group %q): %v", id.ServerName, id.ResourceGroup, err)
	}

	model := result.Model
	if model == nil {
		return nil, fmt.Errorf("reading SQL Security Alert Policy for server %q (Resource Group %q): Model was nil", id.ServerName, id.ResourceGroup)
	}

	return pointer.To(model.Id != nil), nil
}

// Sets the tests 'Primary' and 'Secondary' locations based on the 'ARM_MICROSOFT_TEST' environment variable due to subscription quota policies
func setTestDataLocationBySubscription(data acceptance.TestData) acceptance.TestData {
	if isMicrosoftTest := os.Getenv("ARM_MICROSOFT_TEST"); isMicrosoftTest != "" {
		data.Locations.Primary = "eastus"
		data.Locations.Secondary = "swedencentral"
	}

	return data
}

func (r MsSqlServerSecurityAlertPolicyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_server_security_alert_policy" "test" {
  resource_group_name        = azurerm_resource_group.test.name
  server_name                = azurerm_mssql_server.test.name
  state                      = "Enabled"
  storage_endpoint           = azurerm_storage_account.test.primary_blob_endpoint
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  retention_days             = 20

  disabled_alerts = [
    "Sql_Injection",
    "Data_Exfiltration"
  ]

}
`, r.server(data))
}

func (r MsSqlServerSecurityAlertPolicyResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_server_security_alert_policy" "test" {
  resource_group_name  = azurerm_resource_group.test.name
  server_name          = azurerm_mssql_server.test.name
  state                = "Enabled"
  email_account_admins = true
  retention_days       = 30
}
`, r.server(data))
}

func (MsSqlServerSecurityAlertPolicyResource) server(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsqlserver%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}
`, data.RandomInteger, data.Locations.Primary)
}
