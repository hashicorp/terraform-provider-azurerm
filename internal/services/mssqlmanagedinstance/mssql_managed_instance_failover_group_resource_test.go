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

type MsSqlManagedInstanceFailoverGroupResource struct{}

func TestAccMsSqlManagedInstanceFailoverGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance_failover_group", "test")
	r := MsSqlManagedInstanceFailoverGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: MsSqlManagedInstanceResource{}.dnsZonePartner(data),
		},
		{
			// It speeds up deletion to remove the explicit dependency between the instances
			Config: MsSqlManagedInstanceResource{}.emptyDnsZonePartner(data),
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
			Config: MsSqlManagedInstanceResource{}.emptyDnsZonePartner(data),
		},
	})
}

func (r MsSqlManagedInstanceFailoverGroupResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ManagedInstanceFailoverGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.MSSQLManagedInstance.ManagedInstanceFailoverGroupsClient.Get(ctx, id.ResourceGroup, id.LocationName, id.InstanceFailoverGroupName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r MsSqlManagedInstanceFailoverGroupResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_managed_instance_failover_group" "test" {
  name                        = "acctest-fog-%[2]d"
  location                    = "%[3]s"
  managed_instance_id         = azurerm_mssql_managed_instance.test.id
  partner_managed_instance_id = azurerm_mssql_managed_instance.secondary.id

  read_write_endpoint_failover_policy {
    mode = "Manual"
  }

  depends_on = [
    azurerm_virtual_network_gateway_connection.test,
    azurerm_virtual_network_gateway_connection.secondary,
  ]
}
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r MsSqlManagedInstanceFailoverGroupResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_managed_instance_failover_group" "test" {
  name                        = "acctest-fog-%[2]d"
  location                    = "%[3]s"
  managed_instance_id         = azurerm_mssql_managed_instance.test.id
  partner_managed_instance_id = azurerm_mssql_managed_instance.secondary.id

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
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r MsSqlManagedInstanceFailoverGroupResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

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
  name                = "acctest-pip-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctest-vnetgway-%[2]d"
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
  name                = "acctest-gwc-%[2]d"
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
  name                = "acctest-pip2-%[2]d"
  location            = azurerm_resource_group.secondary.location
  resource_group_name = azurerm_resource_group.secondary.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "secondary" {
  name                = "acctest-vnetgway2-%[2]d"
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
  name                = "acctest-gwc2-%[2]d"
  location            = azurerm_resource_group.secondary.location
  resource_group_name = azurerm_resource_group.secondary.name

  type                            = "Vnet2Vnet"
  virtual_network_gateway_id      = azurerm_virtual_network_gateway.secondary.id
  peer_virtual_network_gateway_id = azurerm_virtual_network_gateway.test.id

  shared_key = var.shared_key
}
`, MsSqlManagedInstanceResource{}.emptyDnsZonePartner(data), data.RandomInteger)
}
