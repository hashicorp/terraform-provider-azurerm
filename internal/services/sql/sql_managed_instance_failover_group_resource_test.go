// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SqlInstanceFailoverGroupResource struct{}

func TestAccAzureRMSqlInstanceFailoverGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_managed_instance_failover_group", "test")
	r := SqlInstanceFailoverGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: SqlManagedInstanceResource{}.dnsZonePartner(data),
		},
		{
			// It speeds up deletion to remove the explicit dependency between the instances
			Config: SqlManagedInstanceResource{}.emptyDnsZonePartner(data),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// disconnect
			Config: SqlManagedInstanceResource{}.emptyDnsZonePartner(data),
		},
	})
}

func (r SqlInstanceFailoverGroupResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.InstanceFailoverGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Sql.InstanceFailoverGroupsClient.Get(ctx, id.ResourceGroup, id.LocationName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving SQL Instance Failover Group %q: %+v", id.ID(), err)
	}
	return utils.Bool(true), nil
}

func (r SqlInstanceFailoverGroupResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

%s

resource "azurerm_sql_managed_instance_failover_group" "test" {
  name                        = "acctest-fog-%d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_sql_managed_instance.test.location
  managed_instance_name       = azurerm_sql_managed_instance.test.name
  partner_managed_instance_id = azurerm_sql_managed_instance.secondary.id

  read_write_endpoint_failover_policy {
    mode = "Manual"
  }

  depends_on = [
    azurerm_virtual_network_gateway_connection.test,
    azurerm_virtual_network_gateway_connection.secondary,
  ]
}
`, SqlManagedInstanceResource{}.emptyDnsZonePartner(data), r.vnetToVnetGateway(data), data.RandomInteger)
}

func (r SqlInstanceFailoverGroupResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

%s

resource "azurerm_sql_managed_instance_failover_group" "test" {
  name                        = "acctest-fog-%d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_sql_managed_instance.test.location
  managed_instance_name       = azurerm_sql_managed_instance.test.name
  partner_managed_instance_id = azurerm_sql_managed_instance.secondary.id

  readonly_endpoint_failover_policy_enabled = true

  read_write_endpoint_failover_policy {
    mode          = "Automatic"
    grace_minutes = 60
  }

  depends_on = [
    azurerm_virtual_network_gateway_connection.test,
    azurerm_virtual_network_gateway_connection.secondary,
  ]
}
`, SqlManagedInstanceResource{}.emptyDnsZonePartner(data), r.vnetToVnetGateway(data), data.RandomInteger)
}

func (r SqlInstanceFailoverGroupResource) vnetToVnetGateway(data acceptance.TestData) string {
	return fmt.Sprintf(`
variable "shared_key" {
  default = "s3cr37"
}

resource "azurerm_subnet" "gateway_snet_test" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-pip-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctest-vnetgway-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "Basic"

  ip_configuration {
    name                          = "vnetGatewayConfig"
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.gateway_snet_test.id
  }
}

resource "azurerm_virtual_network_gateway_connection" "test" {
  name                = "acctest-gwc-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type                            = "Vnet2Vnet"
  virtual_network_gateway_id      = azurerm_virtual_network_gateway.test.id
  peer_virtual_network_gateway_id = azurerm_virtual_network_gateway.secondary.id

  shared_key = var.shared_key
}

resource "azurerm_subnet" "gateway_snet_secondary" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.secondary.name
  virtual_network_name = azurerm_virtual_network.secondary.name
  address_prefixes     = ["10.1.1.0/24"]
}

resource "azurerm_public_ip" "secondary" {
  name                = "acctest-pip2-%[1]d"
  location            = azurerm_resource_group.secondary.location
  resource_group_name = azurerm_resource_group.secondary.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "secondary" {
  name                = "acctest-vnetgway2-%[1]d"
  location            = azurerm_resource_group.secondary.location
  resource_group_name = azurerm_resource_group.secondary.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "Basic"

  ip_configuration {
    name                          = "vnetGatewayConfig"
    public_ip_address_id          = azurerm_public_ip.secondary.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.gateway_snet_secondary.id
  }
}

resource "azurerm_virtual_network_gateway_connection" "secondary" {
  name                = "acctest-gwc2-%[1]d"
  location            = azurerm_resource_group.secondary.location
  resource_group_name = azurerm_resource_group.secondary.name

  type                            = "Vnet2Vnet"
  virtual_network_gateway_id      = azurerm_virtual_network_gateway.secondary.id
  peer_virtual_network_gateway_id = azurerm_virtual_network_gateway.test.id

  shared_key = var.shared_key
}
`, data.RandomInteger)
}
