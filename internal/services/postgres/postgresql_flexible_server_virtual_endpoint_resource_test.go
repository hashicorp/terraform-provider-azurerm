// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2023-06-01-preview/virtualendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type PostgresqlFlexibleServerVirtualEndpointResource struct{}

func TestAccPostgresqlFlexibleServerVirtualEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_virtual_endpoint", "test")
	r := PostgresqlFlexibleServerVirtualEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("source_server_id").Exists(),
				check.That(data.ResourceName).Key("replica_server_id").Exists(),
				check.That(data.ResourceName).Key("type").Exists(),
			),
		},
		data.ImportStep("type"),
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func (r PostgresqlFlexibleServerVirtualEndpointResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := virtualendpoints.ParseVirtualEndpointID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Postgres.VirtualEndpointClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Postgresql Virtual Endpoint (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.Members != nil), nil
}

func (r PostgresqlFlexibleServerVirtualEndpointResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := virtualendpoints.ParseVirtualEndpointID(state.ID)
	if err != nil {
		return nil, err
	}

	if err := postgres.DeletePostgresFlexibileServerVirtualEndpoint(ctx, client.Postgres.VirtualEndpointClient, id); err != nil {
		return nil, fmt.Errorf("deleting Postgresql Virtual Endpoint (%s): %+v", id.String(), err)
	}

	return utils.Bool(true), nil
}

func (PostgresqlFlexibleServerVirtualEndpointResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-psql-virtualendpoint-rg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_postgresql_flexible_server" "test" {
  name                          = "acctest-psql-virtualendpoint-primary-%[1]d"
  resource_group_name           = azurerm_resource_group.test.name
  location                      = azurerm_resource_group.test.location
  version                       = "16"
  public_network_access_enabled = false
  administrator_login           = "psqladmin"
  administrator_password        = "H@Sh1CoR3!"
  zone                          = "1"
  storage_mb                    = 32768
  storage_tier                  = "P30"
  sku_name                      = "GP_Standard_D2ads_v5"
}

resource "azurerm_postgresql_flexible_server" "test_replica" {
  name                          = "acctest-psql-virtualendpoint-replica-%[1]d"
  resource_group_name           = azurerm_postgresql_flexible_server.test.resource_group_name
  location                      = azurerm_postgresql_flexible_server.test.location
  create_mode                   = "Replica"
  source_server_id              = azurerm_postgresql_flexible_server.test.id
  version                       = "16"
  public_network_access_enabled = false
  zone                          = "1"
  storage_mb                    = 32768
  storage_tier                  = "P30"
}

resource "azurerm_postgresql_flexible_server_virtual_endpoint" "test" {
  name              = "acctest-psqlvirtualendpoint-endpoint-%[1]d"
  source_server_id  = azurerm_postgresql_flexible_server.test.id
  replica_server_id = azurerm_postgresql_flexible_server.test_replica.id
  type              = "ReadWrite"
}
`, data.RandomInteger, "eastus") // force region due to SKU constraints
}
