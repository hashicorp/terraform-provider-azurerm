// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/serverdnsaliases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ServerDNSAliasResource struct{}

func TestAccServerDNSAlias_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server_dns_alias", "test")
	r := ServerDNSAliasResource{}

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

func TestAccServerDNSAlias_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server_dns_alias", "test")
	r := ServerDNSAliasResource{}

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

func (r ServerDNSAliasResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := serverdnsaliases.ParseDnsAliasID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.MSSQL.ServerDNSAliasClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retreiving %s: %v", id, err)
	}
	if response.WasNotFound(resp.HttpResponse) {
		return pointer.To(false), nil
	}
	return utils.Bool(true), nil
}

// Configs

func (r ServerDNSAliasResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appServerDNSAlias-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "sql" {
  administrator_login          = "umtacc"
  administrator_login_password = "random81jdpwd_$#fs"
  location                     = azurerm_resource_group.test.location
  name                         = "acctestrg-sql-sever-%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  version                      = "12.0"
}

resource "azurerm_mssql_server_dns_alias" "test" {
  mssql_server_id = azurerm_mssql_server.sql.id
  name            = "acctest-dns-alias-%[1]d"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ServerDNSAliasResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_mssql_server_dns_alias" "import" {
  name            = azurerm_mssql_server_dns_alias.test.name
  mssql_server_id = azurerm_mssql_server_dns_alias.test.mssql_server_id
}
`, r.basic(data))
}
