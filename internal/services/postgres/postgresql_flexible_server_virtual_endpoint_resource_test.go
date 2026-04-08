// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package postgres_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/virtualendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
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

func TestAccPostgresqlFlexibleServerVirtualEndpoint_crossSubscription(t *testing.T) {
	t.Skip("Skipping: cross subscription replication is non-standard operation and need to add the subscriptions to a service allow list")
	altSubscription := getAltSubscription()

	if altSubscription == nil {
		t.Skip("Skipping: Test requires `ARM_SUBSCRIPTION_ID_ALT` and `ARM_TENANT_ID` environment variables to be specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_virtual_endpoint", "test")
	r := PostgresqlFlexibleServerVirtualEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.crossSubscription(data, altSubscription),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("replica_server_id"),
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

	return pointer.To(true), nil
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
`, data.RandomInteger, data.Locations.Primary) // force region due to SKU constraints
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
`, data.RandomInteger, data.Locations.Primary, replicaId) // force region due to SKU constraints
}

/** Complex test cases across regions and resource groups */
func (PostgresqlFlexibleServerVirtualEndpointResource) crossRegion(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

###### Primary RG ######

resource "azurerm_resource_group" "primary" {
  name     = "acctest%[1]d-primary"
  location = "%[3]s"
}

resource "azurerm_virtual_network" "primary" {
  name                = "acctest%[1]d-primary-vn"
  location            = azurerm_resource_group.primary.location
  resource_group_name = azurerm_resource_group.primary.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_network_security_group" "primary" {
  name                = "acctest%[1]d-primary-nsg"
  location            = azurerm_resource_group.primary.location
  resource_group_name = azurerm_resource_group.primary.name

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

resource "azurerm_subnet" "primary" {
  name                 = "acctest%[1]d-primary-sn"
  resource_group_name  = azurerm_resource_group.primary.name
  virtual_network_name = azurerm_virtual_network.primary.name
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

resource "azurerm_subnet_network_security_group_association" "primary" {
  subnet_id                 = azurerm_subnet.primary.id
  network_security_group_id = azurerm_network_security_group.primary.id
}

resource "azurerm_private_dns_zone" "primary" {
  name                = "acctest%[1]d-pdz.postgres.database.azure.com"
  resource_group_name = azurerm_resource_group.primary.name

  depends_on = [azurerm_subnet_network_security_group_association.primary]
}

resource "azurerm_virtual_network_peering" "primary" {
  name                         = "primary-to-secondary"
  resource_group_name          = azurerm_resource_group.primary.name
  virtual_network_name         = azurerm_virtual_network.primary.name
  remote_virtual_network_id    = azurerm_virtual_network.secondary.id
  allow_virtual_network_access = true
  allow_forwarded_traffic      = true
}

resource "azurerm_private_dns_zone_virtual_network_link" "primary" {
  name                  = "acctest%[1]d-primary-pdzvnetlink.com"
  private_dns_zone_name = azurerm_private_dns_zone.primary.name
  virtual_network_id    = azurerm_virtual_network.primary.id
  resource_group_name   = azurerm_resource_group.primary.name

  depends_on = [azurerm_virtual_network_peering.primary]
}

resource "azurerm_postgresql_flexible_server" "primary" {
  name                          = "acctest%[1]d-primary-pg"
  resource_group_name           = azurerm_resource_group.primary.name
  location                      = azurerm_resource_group.primary.location
  version                       = "13"
  public_network_access_enabled = false
  administrator_login           = "adminTerraform"
  administrator_password        = "maLTq5SnDBrWfyban7Wz"
  sku_name                      = "GP_Standard_D2s_v3"

  delegated_subnet_id = azurerm_subnet.primary.id
  private_dns_zone_id = azurerm_private_dns_zone.primary.id

  depends_on = [azurerm_private_dns_zone_virtual_network_link.primary]

  lifecycle {
    ignore_changes = [zone]
  }

  timeouts {
    create = "120m"
  }
}

resource "azurerm_postgresql_flexible_server_virtual_endpoint" "test" {
  name              = "acctest%[1]d-endpoint"
  source_server_id  = azurerm_postgresql_flexible_server.primary.id
  replica_server_id = azurerm_postgresql_flexible_server.secondary.id
  type              = "ReadWrite"
}

###### Secondary RG ######

resource "azurerm_resource_group" "secondary" {
  name     = "acctest%[1]d-secondary"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "secondary" {
  name                = "acctest%[1]d-secondary-vn"
  location            = azurerm_resource_group.secondary.location
  resource_group_name = azurerm_resource_group.secondary.name
  address_space       = ["11.0.0.0/16"]
}

resource "azurerm_network_security_group" "secondary" {
  name                = "acctest%[1]d-secondary-nsg"
  location            = azurerm_resource_group.secondary.location
  resource_group_name = azurerm_resource_group.secondary.name

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

resource "azurerm_subnet" "secondary" {
  name                 = "acctest%[1]d-secondary-sn"
  resource_group_name  = azurerm_resource_group.secondary.name
  virtual_network_name = azurerm_virtual_network.secondary.name
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

resource "azurerm_subnet_network_security_group_association" "secondary" {
  subnet_id                 = azurerm_subnet.secondary.id
  network_security_group_id = azurerm_network_security_group.secondary.id
}

resource "azurerm_private_dns_zone" "secondary" {
  name                = "acctest%[1]d-pdz.postgres.database.azure.com"
  resource_group_name = azurerm_resource_group.secondary.name

  depends_on = [azurerm_subnet_network_security_group_association.secondary]
}

resource "azurerm_virtual_network_peering" "secondary" {
  name                         = "acctest-pfs%[1]d-secondary-to-primary"
  resource_group_name          = azurerm_resource_group.secondary.name
  virtual_network_name         = azurerm_virtual_network.secondary.name
  remote_virtual_network_id    = azurerm_virtual_network.primary.id
  allow_virtual_network_access = true
  allow_forwarded_traffic      = true
}

resource "azurerm_private_dns_zone_virtual_network_link" "secondary" {
  name                  = "acctest%[1]d-secondary-pdzvnetlink.com"
  private_dns_zone_name = azurerm_private_dns_zone.secondary.name
  virtual_network_id    = azurerm_virtual_network.secondary.id
  resource_group_name   = azurerm_resource_group.secondary.name

  depends_on = [azurerm_virtual_network_peering.secondary]
}

resource "azurerm_postgresql_flexible_server" "secondary" {
  name                          = "acctest%[1]d-secondary-pg"
  resource_group_name           = azurerm_resource_group.secondary.name
  location                      = azurerm_resource_group.secondary.location
  create_mode                   = "Replica"
  source_server_id              = azurerm_postgresql_flexible_server.primary.id
  version                       = azurerm_postgresql_flexible_server.primary.version
  public_network_access_enabled = azurerm_postgresql_flexible_server.primary.public_network_access_enabled
  sku_name                      = azurerm_postgresql_flexible_server.primary.sku_name

  delegated_subnet_id = azurerm_subnet.secondary.id
  private_dns_zone_id = azurerm_private_dns_zone.secondary.id

  depends_on = [azurerm_private_dns_zone_virtual_network_link.secondary]

  lifecycle {
    ignore_changes = [zone]
  }

  timeouts {
    create = "120m"
  }
}
`, data.RandomInteger, data.Locations.Secondary, data.Locations.Ternary)
}

type alternateSubscription struct {
	tenant_id       string
	subscription_id string
}

func getAltSubscription() *alternateSubscription {
	altSubscriptonID := os.Getenv("ARM_SUBSCRIPTION_ID_ALT")
	altTenantID := os.Getenv("ARM_TENANT_ID")

	if altSubscriptonID == "" || altTenantID == "" {
		return nil
	}

	return &alternateSubscription{
		tenant_id:       altTenantID,
		subscription_id: altSubscriptonID,
	}
}

func (PostgresqlFlexibleServerVirtualEndpointResource) crossSubscription(data acceptance.TestData, altSub *alternateSubscription) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
  }
}

provider "azurerm-alt" {
  features {}

  tenant_id       = "%[2]s"
  subscription_id = "%[3]s"
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[4]s" // force region due to service allow list
}

resource "azurerm_resource_group" "alt" {
  provider = azurerm-alt

  name     = "acctestRG-alt-%[1]d"
  location = "eastus2" // force region due to service allow list
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
  sku_name                      = "GP_Standard_D2s_v3"
}

resource "azurerm_postgresql_flexible_server" "test_replica" {
  provider = azurerm-alt

  name                          = "acctest-ve-replica-%[1]d"
  resource_group_name           = azurerm_resource_group.alt.name
  location                      = azurerm_resource_group.alt.location
  create_mode                   = "Replica"
  source_server_id              = azurerm_postgresql_flexible_server.test.id
  version                       = azurerm_postgresql_flexible_server.test.version
  public_network_access_enabled = azurerm_postgresql_flexible_server.test.public_network_access_enabled
  storage_mb                    = azurerm_postgresql_flexible_server.test.storage_mb
  storage_tier                  = azurerm_postgresql_flexible_server.test.storage_tier
  zone                          = "1"
}

resource "azurerm_postgresql_flexible_server_virtual_endpoint" "test" {
  name              = "acctest-ve-%[1]d"
  source_server_id  = azurerm_postgresql_flexible_server.test.id
  replica_server_id = azurerm_postgresql_flexible_server.test_replica.id
  type              = "ReadWrite"
}
`, data.RandomInteger, altSub.tenant_id, altSub.subscription_id, data.Locations.Primary)
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
`, data.RandomInteger, data.Locations.Primary) // force region due to SKU constraints
}
