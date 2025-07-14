// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2024-08-01/virtualendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
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

func TestAccPostgresqlFlexibleServerVirtualEndpoint_crossRegion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_virtual_endpoint", "test")
	r := PostgresqlFlexibleServerVirtualEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.crossRegion(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.crossRegion,
			TestResource: r,
		}),
	})
}

func TestAccPostgresqlFlexibleServerVirtualEndpoint_identicalSourceAndReplica(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_virtual_endpoint", "test")
	r := PostgresqlFlexibleServerVirtualEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identicalSourceAndReplica(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r PostgresqlFlexibleServerVirtualEndpointResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseCompositeResourceID(state.ID, &virtualendpoints.VirtualEndpointId{}, &virtualendpoints.VirtualEndpointId{})
	if err != nil {
		return nil, err
	}

	virtualEndpointId := virtualendpoints.NewVirtualEndpointID(id.First.SubscriptionId, id.First.ResourceGroupName, id.First.FlexibleServerName, id.First.VirtualEndpointName)

	resp, err := clients.Postgres.VirtualEndpointClient.Get(ctx, virtualEndpointId)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.Members != nil), nil
}

func (r PostgresqlFlexibleServerVirtualEndpointResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseCompositeResourceID(state.ID, &virtualendpoints.VirtualEndpointId{}, &virtualendpoints.VirtualEndpointId{})
	if err != nil {
		return nil, err
	}

	virtualEndpointId := virtualendpoints.NewVirtualEndpointID(id.First.SubscriptionId, id.First.ResourceGroupName, id.First.FlexibleServerName, id.First.VirtualEndpointName)

	if err := client.Postgres.VirtualEndpointClient.DeleteThenPoll(ctx, virtualEndpointId); err != nil {
		return nil, fmt.Errorf("deleting %s: %+v", *id, err)
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

/** Complex test cases across regions and resource groups */
func (PostgresqlFlexibleServerVirtualEndpointResource) crossRegion(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

###### EAST RG ######

resource "azurerm_resource_group" "east" {
  name     = "acctest%[1]d-east"
  location = "eastus"
}

resource "azurerm_virtual_network" "east" {
  name                = "acctest%[1]d-east-vn"
  location            = azurerm_resource_group.east.location
  resource_group_name = azurerm_resource_group.east.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_network_security_group" "east" {
  name                = "acctest%[1]d-east-nsg"
  location            = azurerm_resource_group.east.location
  resource_group_name = azurerm_resource_group.east.name

  security_rule {
    name                       = "allow_all"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "*"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}

resource "azurerm_subnet" "east" {
  name                 = "acctest%[1]d-east-sn"
  resource_group_name  = azurerm_resource_group.east.name
  virtual_network_name = azurerm_virtual_network.east.name
  address_prefixes     = ["10.0.1.0/24"]
  service_endpoints    = ["Microsoft.Storage"]

  delegation {
    name = "fs"
    service_delegation {
      name = "Microsoft.DBforPostgreSQL/flexibleServers"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_subnet_network_security_group_association" "east" {
  subnet_id                 = azurerm_subnet.east.id
  network_security_group_id = azurerm_network_security_group.east.id
}

resource "azurerm_private_dns_zone" "east" {
  name                = "acctest%[1]d-pdz.postgres.database.azure.com"
  resource_group_name = azurerm_resource_group.east.name

  depends_on = [azurerm_subnet_network_security_group_association.east]
}

resource "azurerm_virtual_network_peering" "east" {
  name                         = "east-to-west"
  resource_group_name          = azurerm_resource_group.east.name
  virtual_network_name         = azurerm_virtual_network.east.name
  remote_virtual_network_id    = azurerm_virtual_network.west.id
  allow_virtual_network_access = true
  allow_forwarded_traffic      = true
}

resource "azurerm_private_dns_zone_virtual_network_link" "east" {
  name                  = "acctest%[1]d-east-pdzvnetlink.com"
  private_dns_zone_name = azurerm_private_dns_zone.east.name
  virtual_network_id    = azurerm_virtual_network.east.id
  resource_group_name   = azurerm_resource_group.east.name

  depends_on = [azurerm_virtual_network_peering.east]
}

resource "azurerm_postgresql_flexible_server" "east" {
  name                          = "acctest%[1]d-east-pg"
  resource_group_name           = azurerm_resource_group.east.name
  location                      = azurerm_resource_group.east.location
  version                       = "13"
  public_network_access_enabled = false
  administrator_login           = "adminTerraform"
  administrator_password        = "maLTq5SnDBrWfyban7Wz"
  sku_name                      = "GP_Standard_D2s_v3"

  delegated_subnet_id = azurerm_subnet.east.id
  private_dns_zone_id = azurerm_private_dns_zone.east.id

  depends_on = [azurerm_private_dns_zone_virtual_network_link.east]

  lifecycle {
    ignore_changes = [zone]
  }

  timeouts {
    create = "120m"
  }
}

resource "azurerm_postgresql_flexible_server_virtual_endpoint" "test" {
  name              = "acctest%[1]d-endpoint"
  source_server_id  = azurerm_postgresql_flexible_server.east.id
  replica_server_id = azurerm_postgresql_flexible_server.west.id
  type              = "ReadWrite"
}

###### WEST RG ######

resource "azurerm_resource_group" "west" {
  name     = "acctest%[1]d-west"
  location = "westus"
}

resource "azurerm_virtual_network" "west" {
  name                = "acctest%[1]d-west-vn"
  location            = azurerm_resource_group.west.location
  resource_group_name = azurerm_resource_group.west.name
  address_space       = ["11.0.0.0/16"]
}

resource "azurerm_network_security_group" "west" {
  name                = "acctest%[1]d-west-nsg"
  location            = azurerm_resource_group.west.location
  resource_group_name = azurerm_resource_group.west.name

  security_rule {
    name                       = "allow_all"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "*"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}

resource "azurerm_subnet" "west" {
  name                 = "acctest%[1]d-west-sn"
  resource_group_name  = azurerm_resource_group.west.name
  virtual_network_name = azurerm_virtual_network.west.name
  address_prefixes     = ["11.0.1.0/24"]
  service_endpoints    = ["Microsoft.Storage"]

  delegation {
    name = "fs"
    service_delegation {
      name = "Microsoft.DBforPostgreSQL/flexibleServers"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_subnet_network_security_group_association" "west" {
  subnet_id                 = azurerm_subnet.west.id
  network_security_group_id = azurerm_network_security_group.west.id
}

resource "azurerm_private_dns_zone" "west" {
  name                = "acctest%[1]d-pdz.postgres.database.azure.com"
  resource_group_name = azurerm_resource_group.west.name

  depends_on = [azurerm_subnet_network_security_group_association.west]
}

resource "azurerm_virtual_network_peering" "west" {
  name                         = "acctest-pfs%[1]d-west-to-east"
  resource_group_name          = azurerm_resource_group.west.name
  virtual_network_name         = azurerm_virtual_network.west.name
  remote_virtual_network_id    = azurerm_virtual_network.east.id
  allow_virtual_network_access = true
  allow_forwarded_traffic      = true
}

resource "azurerm_private_dns_zone_virtual_network_link" "west" {
  name                  = "acctest%[1]d-west-pdzvnetlink.com"
  private_dns_zone_name = azurerm_private_dns_zone.west.name
  virtual_network_id    = azurerm_virtual_network.west.id
  resource_group_name   = azurerm_resource_group.west.name

  depends_on = [azurerm_virtual_network_peering.west]
}

resource "azurerm_postgresql_flexible_server" "west" {
  name                          = "acctest%[1]d-west-pg"
  resource_group_name           = azurerm_resource_group.west.name
  location                      = azurerm_resource_group.west.location
  create_mode                   = "Replica"
  source_server_id              = azurerm_postgresql_flexible_server.east.id
  version                       = azurerm_postgresql_flexible_server.east.version
  public_network_access_enabled = azurerm_postgresql_flexible_server.east.public_network_access_enabled
  sku_name                      = azurerm_postgresql_flexible_server.east.sku_name

  delegated_subnet_id = azurerm_subnet.west.id
  private_dns_zone_id = azurerm_private_dns_zone.west.id

  depends_on = [azurerm_private_dns_zone_virtual_network_link.west]

  lifecycle {
    ignore_changes = [zone]
  }

  timeouts {
    create = "120m"
  }
}
`, data.RandomInteger)
}

func (PostgresqlFlexibleServerVirtualEndpointResource) identicalSourceAndReplica(data acceptance.TestData) string {
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

resource "azurerm_postgresql_flexible_server_virtual_endpoint" "test" {
  name              = "acctest-ve-%[1]d"
  source_server_id  = azurerm_postgresql_flexible_server.test.id
  replica_server_id = azurerm_postgresql_flexible_server.test.id
  type              = "ReadWrite"
}
`, data.RandomInteger, "eastus") // force region due to SKU constraints
}
