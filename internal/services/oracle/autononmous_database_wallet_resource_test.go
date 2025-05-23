// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AutonomousDatabaseWalletResource struct{}

func (r AutonomousDatabaseWalletResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := autonomousdatabases.ParseAutonomousDatabaseID(state.ID)
	if err != nil {
		return nil, err
	}

	// Since the wallet is generated on-demand and doesn't have a persistent state to check,
	// we just verify the autonomous database exists
	resp, err := client.Oracle.OracleClient.AutonomousDatabases.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func TestAccOracleAutonomousDatabaseWallet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_oracle_autonomous_database_wallet", "test")
	r := AutonomousDatabaseWalletResource{}
	adbsResource := AdbsRegularResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, adbsResource),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("wallet_files").Exists(),
				check.That(data.ResourceName).Key("generate_type").HasValue("SINGLE"),
				check.That(data.ResourceName).Key("is_regional").HasValue("false"),
			),
		},
	})
}

func TestAccOracleAutonomousDatabaseWallet_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_oracle_autonomous_database_wallet", "test")
	r := AutonomousDatabaseWalletResource{}
	adbsResource := AdbsRegularResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, adbsResource),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("wallet_files").Exists(),
				check.That(data.ResourceName).Key("generate_type").HasValue("ALL"),
				check.That(data.ResourceName).Key("is_regional").HasValue("true"),
			),
		},
	})
}

func (r AutonomousDatabaseWalletResource) basic(data acceptance.TestData, adbsResource AdbsRegularResource) string {
	return fmt.Sprintf(`
%s
resource "azurerm_oracle_autonomous_database_wallet" "test" {
  autonomous_database_id = azurerm_oracle_autonomous_database.test.id
  password               = "TestPass#2024#"
}
`, adbsResource.basic(data))
}

func (r AutonomousDatabaseWalletResource) complete(data acceptance.TestData, adbsResource AdbsRegularResource) string {
	return fmt.Sprintf(`
%s
resource "azurerm_oracle_autonomous_database_wallet" "test" {
  autonomous_database_id = azurerm_oracle_autonomous_database.test.id
  password               = "TestPass#2024#"
  generate_type          = "ALL"
  is_regional            = true
}
`, adbsResource.basic(data))
}
