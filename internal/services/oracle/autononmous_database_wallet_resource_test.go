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
				check.That(data.ResourceName).Key("content").Exists(),
				check.That(data.ResourceName).Key("generate_type").HasValue("SINGLE"),
				check.That(data.ResourceName).Key("base64_encode").HasValue("false"),
			),
		},
	})
}

func TestAccOracleAutonomousDatabaseWallet_withBase64Encode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_oracle_autonomous_database_wallet", "test")
	r := AutonomousDatabaseWalletResource{}
	adbsResource := AdbsRegularResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withBase64Encode(data, adbsResource),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("content").Exists(),
				check.That(data.ResourceName).Key("base64_encode").HasValue("true"),
			),
		},
	})
}

func TestAccOracleAutonomousDatabaseWallet_withRegionalGenerateType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_oracle_autonomous_database_wallet", "test")
	r := AutonomousDatabaseWalletResource{}
	adbsResource := AdbsRegularResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withRegionalGenerateType(data, adbsResource),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("content").Exists(),
				check.That(data.ResourceName).Key("generate_type").HasValue("REGIONAL"),
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

func (r AutonomousDatabaseWalletResource) withBase64Encode(data acceptance.TestData, adbsResource AdbsRegularResource) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_autonomous_database_wallet" "test" {
  autonomous_database_id = azurerm_oracle_autonomous_database.test.id
  password               = "TestPass#2024#"
  base64_encode          = true
}
`, adbsResource.basic(data))
}

func (r AutonomousDatabaseWalletResource) withRegionalGenerateType(data acceptance.TestData, adbsResource AdbsRegularResource) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_autonomous_database_wallet" "test" {
  autonomous_database_id = azurerm_oracle_autonomous_database.test.id
  password               = "TestPass#2024#"
  generate_type          = "REGIONAL"
}
`, adbsResource.basic(data))
}
