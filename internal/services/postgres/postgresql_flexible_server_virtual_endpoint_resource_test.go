// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
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
			),
		},
		data.ImportStep(),
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccPostgresqlFlexibleServerVirtualEndpoint_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_virtual_endpoint", "test")
	r := PostgresqlFlexibleServerVirtualEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.update(data, "azurerm_postgresql_flexible_server.test_replica_0.id"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data, "azurerm_postgresql_flexible_server.test_replica_1.id"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r PostgresqlFlexibleServerVirtualEndpointResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := virtualendpoints.ParseVirtualEndpointID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Postgres.VirtualEndpointClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.Members != nil), nil
}

func (r PostgresqlFlexibleServerVirtualEndpointResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := virtualendpoints.ParseVirtualEndpointID(state.ID)
	if err != nil {
		return nil, err
	}

	if err := postgres.DeletePostgresFlexibileServerVirtualEndpoint(ctx, client.Postgres.VirtualEndpointClient, id); err != nil {
		return nil, fmt.Errorf("deleting %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (PostgresqlFlexibleServerVirtualEndpointResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-ve-rg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_postgresql_flexible_server" "test" {
  name                          = "acctest-ve-primary-%[1]d"
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
  name                          = "acctest-ve-replica-%[1]d"
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
  name              = "acctest-ve-%[1]d"
  source_server_id  = azurerm_postgresql_flexible_server.test.id
  replica_server_id = azurerm_postgresql_flexible_server.test_replica.id
  type              = "ReadWrite"
}
`, data.RandomInteger, "eastus") // force region due to SKU constraints
}

func (PostgresqlFlexibleServerVirtualEndpointResource) update(data acceptance.TestData, replicaId string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-ve-rg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_postgresql_flexible_server" "test" {
  name                          = "acctest-ve-primary-%[1]d"
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

resource "azurerm_postgresql_flexible_server" "test_replica_0" {
  name                          = "${azurerm_postgresql_flexible_server.test.name}-replica-0"
  resource_group_name           = azurerm_postgresql_flexible_server.test.resource_group_name
  location                      = azurerm_postgresql_flexible_server.test.location
  create_mode                   = "Replica"
  source_server_id              = azurerm_postgresql_flexible_server.test.id
  version                       = azurerm_postgresql_flexible_server.test.version
  public_network_access_enabled = azurerm_postgresql_flexible_server.test.public_network_access_enabled
  zone                          = azurerm_postgresql_flexible_server.test.zone
  storage_mb                    = azurerm_postgresql_flexible_server.test.storage_mb
  storage_tier                  = azurerm_postgresql_flexible_server.test.storage_tier
}

resource "azurerm_postgresql_flexible_server" "test_replica_1" {
  name                          = "${azurerm_postgresql_flexible_server.test.name}-replica-1"
  resource_group_name           = azurerm_postgresql_flexible_server.test.resource_group_name
  location                      = azurerm_postgresql_flexible_server.test.location
  create_mode                   = "Replica"
  source_server_id              = azurerm_postgresql_flexible_server.test.id
  version                       = azurerm_postgresql_flexible_server.test.version
  public_network_access_enabled = azurerm_postgresql_flexible_server.test.public_network_access_enabled
  zone                          = azurerm_postgresql_flexible_server.test.zone
  storage_mb                    = azurerm_postgresql_flexible_server.test.storage_mb
  storage_tier                  = azurerm_postgresql_flexible_server.test.storage_tier

  ## this prevents a race condition that can occur when 2 replicas are created simultaneously
  depends_on = [azurerm_postgresql_flexible_server.test_replica_0]
}

resource "azurerm_postgresql_flexible_server_virtual_endpoint" "test" {
  name              = "acctest-ve-%[1]d"
  source_server_id  = azurerm_postgresql_flexible_server.test.id
  replica_server_id = %[3]s
  type              = "ReadWrite"

  ## this prevents a race condition that can occur if the virtual endpoint is created while a replica is still initializing
  depends_on = [azurerm_postgresql_flexible_server.test_replica_0, azurerm_postgresql_flexible_server.test_replica_1]
}
`, data.RandomInteger, "eastus", replicaId) // force region due to SKU constraints
}
