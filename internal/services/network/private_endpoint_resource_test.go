// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/privateendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type PrivateEndpointResource struct{}

func TestAccPrivateEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint", "test")
	r := PrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subnet_id").Exists(),
				check.That(data.ResourceName).Key("network_interface.0.id").Exists(),
				check.That(data.ResourceName).Key("network_interface.0.name").Exists(),
				check.That(data.ResourceName).Key("private_service_connection.0.private_ip_address").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPrivateEndpoint_updateTag(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint", "test")
	r := PrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withTag(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPrivateEndpoint_withCustomNicName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint", "test")
	r := PrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withCustomNicName(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("custom_network_interface_name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPrivateEndpoint_requestMessage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint", "test")
	r := PrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.requestMessage(data, "CATS: ALL YOUR BASE ARE BELONG TO US."),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subnet_id").Exists(),
				check.That(data.ResourceName).Key("network_interface.0.id").Exists(),
				check.That(data.ResourceName).Key("network_interface.0.name").Exists(),
				check.That(data.ResourceName).Key("private_service_connection.0.request_message").HasValue("CATS: ALL YOUR BASE ARE BELONG TO US."),
			),
		},
		data.ImportStep(),
		{
			Config: r.requestMessage(data, "CAPTAIN: WHAT YOU SAY!!"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subnet_id").Exists(),
				check.That(data.ResourceName).Key("network_interface.0.id").Exists(),
				check.That(data.ResourceName).Key("network_interface.0.name").Exists(),
				check.That(data.ResourceName).Key("private_service_connection.0.request_message").HasValue("CAPTAIN: WHAT YOU SAY!!"),
			),
		},
		data.ImportStep(),
	})
}

// The update and complete test cases had to be totally removed since there is a bug with tags and the support for
// tags has been removed, all other attributes are ForceNew.
// API Issue "Unable to remove Tags from Private Endpoint": https://github.com/Azure/azure-sdk-for-go/issues/6467

func TestAccPrivateEndpoint_privateDnsZoneGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint", "test")
	r := PrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateDnsZoneGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_dns_zone_group.0.private_dns_zone_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("private_dns_zone_configs.#").HasValue("1"),
				check.That(data.ResourceName).Key("private_dns_zone_group.#").HasValue("1"),
			),
		},
		data.ImportStep("private_dns_zone_configs", "private_dns_zone_group"),
	})
}

func TestAccPrivateEndpoint_privateDnsZoneRename(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint", "test")
	r := PrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateDnsZoneGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_dns_zone_group.0.private_dns_zone_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("private_dns_zone_configs.#").HasValue("1"),
				check.That(data.ResourceName).Key("private_dns_zone_group.#").HasValue("1"),
			),
		},
		data.ImportStep("private_dns_zone_configs", "private_dns_zone_group"),
		{
			Config: r.privateDnsZoneGroupRename(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_dns_zone_group.0.private_dns_zone_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("private_dns_zone_configs.#").HasValue("1"),
				check.That(data.ResourceName).Key("private_dns_zone_group.#").HasValue("1"),
			),
		},
		data.ImportStep("private_dns_zone_configs", "private_dns_zone_group"),
	})
}

func TestAccPrivateEndpoint_privateDnsZoneUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint", "test")
	r := PrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateDnsZoneGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_dns_zone_group.0.private_dns_zone_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("private_dns_zone_configs.#").HasValue("1"),
				check.That(data.ResourceName).Key("private_dns_zone_group.#").HasValue("1"),
			),
		},
		data.ImportStep("private_dns_zone_configs", "private_dns_zone_group"),
		{
			Config: r.privateDnsZoneGroupUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_dns_zone_group.0.private_dns_zone_ids.#").HasValue("2"),
				check.That(data.ResourceName).Key("private_dns_zone_configs.#").HasValue("2"),
				check.That(data.ResourceName).Key("private_dns_zone_group.#").HasValue("1"),
			),
		},
		data.ImportStep("private_dns_zone_configs", "private_dns_zone_group"),
	})
}

func TestAccPrivateEndpoint_privateDnsZoneIdsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint", "test")
	r := PrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateDnsZoneGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_dns_zone_group.0.private_dns_zone_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("private_dns_zone_configs.#").HasValue("1"),
				check.That(data.ResourceName).Key("private_dns_zone_group.#").HasValue("1"),
			),
		},
		data.ImportStep("private_dns_zone_configs", "private_dns_zone_group"),
		{
			Config: r.privateDnsZoneGroupIdsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_dns_zone_group.0.private_dns_zone_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("private_dns_zone_configs.#").HasValue("1"),
				check.That(data.ResourceName).Key("private_dns_zone_group.#").HasValue("1"),
			),
		},
		data.ImportStep("private_dns_zone_configs", "private_dns_zone_group"),
	})
}

func TestAccPrivateEndpoint_staticIpAddress(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint", "test")
	r := PrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.staticIpAddress(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subnet_id").Exists(),
				check.That(data.ResourceName).Key("network_interface.0.id").Exists(),
				check.That(data.ResourceName).Key("network_interface.0.name").Exists(),
				check.That(data.ResourceName).Key("private_service_connection.0.private_ip_address").Exists(),
				check.That(data.ResourceName).Key("private_service_connection.0.private_ip_address").HasValue("10.5.2.47"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPrivateEndpoint_privateDnsZoneRemove(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint", "test")
	r := PrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateDnsZoneGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_dns_zone_group.0.private_dns_zone_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("private_dns_zone_configs.#").HasValue("1"),
				check.That(data.ResourceName).Key("private_dns_zone_group.#").HasValue("1"),
			),
		},
		data.ImportStep("private_dns_zone_configs", "private_dns_zone_group"),
		{
			Config: r.privateDnsZoneGroupRemove(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("private_dns_zone_configs", "private_dns_zone_group"),
		{
			Config: r.privateDnsZoneGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_dns_zone_group.0.private_dns_zone_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("private_dns_zone_configs.#").HasValue("1"),
				check.That(data.ResourceName).Key("private_dns_zone_group.#").HasValue("1"),
			),
		},
		data.ImportStep("private_dns_zone_configs", "private_dns_zone_group"),
	})
}

func TestAccPrivateEndpoint_privateConnectionAlias(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint", "test")
	r := PrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateConnectionAlias(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subnet_id").Exists(),
				check.That(data.ResourceName).Key("network_interface.0.id").Exists(),
				check.That(data.ResourceName).Key("network_interface.0.name").Exists(),
				check.That(data.ResourceName).Key("private_service_connection.0.private_connection_resource_alias").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPrivateEndpoint_updateToPrivateConnectionAlias(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint", "test")
	r := PrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateConnectionAlias(data, false),
		},
		data.ImportStep(),
		{
			Config: r.privateConnectionAlias(data, true),
		},
		data.ImportStep(),
	})
}

func (t PrivateEndpointResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := privateendpoints.ParsePrivateEndpointID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.PrivateEndpoints.Get(ctx, *id, privateendpoints.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("reading Private Endpoint (%s): %+v", id.String(), err)
	}

	return pointer.To(resp.Model != nil), nil
}

func TestAccPrivateEndpoint_multipleInstances(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint", "test")
	r := PrivateEndpointResource{}

	instanceCount := 5
	var checks []pluginsdk.TestCheckFunc
	for i := 0; i < instanceCount; i++ {
		checks = append(checks, check.That(fmt.Sprintf("%s.%d", data.ResourceName, i)).ExistsInAzure(r))
	}

	config := r.multipleInstances(data, instanceCount, false)
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: config,
			Check:  acceptance.ComposeTestCheckFunc(checks...),
		},
	})
}

func TestAccPrivateEndpoint_multipleInstancesWithLinkAlias(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint", "test")
	r := PrivateEndpointResource{}

	instanceCount := 5
	var checks []pluginsdk.TestCheckFunc
	for i := 0; i < instanceCount; i++ {
		checks = append(checks, check.That(fmt.Sprintf("%s.%d", data.ResourceName, i)).ExistsInAzure(r))
	}

	config := r.multipleInstances(data, instanceCount, true)
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: config,
			Check:  acceptance.ComposeTestCheckFunc(checks...),
		},
	})
}

func TestAccPrivateEndpoint_multipleIpConfigurations(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint", "test")
	r := PrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.recoveryServiceVaultWithMultiIpConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (PrivateEndpointResource) template(data acceptance.TestData, seviceCfg string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-privatelink-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "service" {
  name                 = "acctestsnetservice-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.1.0/24"]

  private_link_service_network_policies_enabled = false
}

resource "azurerm_subnet" "endpoint" {
  name                 = "acctestsnetendpoint-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.2.0/24"]

  private_endpoint_network_policies = "Disabled"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%d"
  sku                 = "Standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%d"
  sku                 = "Standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  frontend_ip_configuration {
    name                 = azurerm_public_ip.test.name
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

%s
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, seviceCfg)
}

func (PrivateEndpointResource) serviceAutoApprove(data acceptance.TestData) string {
	return fmt.Sprintf(`

resource "azurerm_private_link_service" "test" {
  name                           = "acctestPLS-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  auto_approval_subscription_ids = [data.azurerm_subscription.current.subscription_id]
  visibility_subscription_ids    = [data.azurerm_subscription.current.subscription_id]

  nat_ip_configuration {
    name      = "primaryIpConfiguration-%d"
    primary   = true
    subnet_id = azurerm_subnet.service.id
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]
}
`, data.RandomInteger, data.RandomInteger)
}

func (PrivateEndpointResource) serviceManualApprove(data acceptance.TestData) string {
	return fmt.Sprintf(`

resource "azurerm_private_link_service" "test" {
  name                = "acctestPLS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  nat_ip_configuration {
    name      = "primaryIpConfiguration-%d"
    primary   = true
    subnet_id = azurerm_subnet.service.id
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]
}
`, data.RandomInteger, data.RandomInteger)
}

func (r PrivateEndpointResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-privatelink-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = azurerm_subnet.endpoint.id

  private_service_connection {
    name                           = azurerm_private_link_service.test.name
    is_manual_connection           = false
    private_connection_resource_id = azurerm_private_link_service.test.id
  }
}
`, r.template(data, r.serviceAutoApprove(data)), data.RandomInteger)
}

func (r PrivateEndpointResource) withTag(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-privatelink-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = azurerm_subnet.endpoint.id

  private_service_connection {
    name                           = azurerm_private_link_service.test.name
    is_manual_connection           = false
    private_connection_resource_id = azurerm_private_link_service.test.id
  }

  tags = {
    env = "TEST"
  }
}
`, r.template(data, r.serviceAutoApprove(data)), data.RandomInteger)
}

func (r PrivateEndpointResource) withCustomNicName(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_endpoint" "test" {
  name                          = "acctest-privatelink-%d"
  resource_group_name           = azurerm_resource_group.test.name
  location                      = azurerm_resource_group.test.location
  subnet_id                     = azurerm_subnet.endpoint.id
  custom_network_interface_name = "acctest-privatelink-%d-nic"

  private_service_connection {
    name                           = azurerm_private_link_service.test.name
    is_manual_connection           = false
    private_connection_resource_id = azurerm_private_link_service.test.id
  }
}
`, r.template(data, r.serviceAutoApprove(data)), data.RandomInteger, data.RandomInteger)
}

func (r PrivateEndpointResource) requestMessage(data acceptance.TestData, msg string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-privatelink-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = azurerm_subnet.endpoint.id

  private_service_connection {
    name                           = azurerm_private_link_service.test.name
    is_manual_connection           = true
    private_connection_resource_id = azurerm_private_link_service.test.id
    request_message                = %q
  }
}
`, r.template(data, r.serviceManualApprove(data)), data.RandomInteger, msg)
}

func (PrivateEndpointResource) privateDnsZoneGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-privatelink-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "service" {
  name                 = "acctestsnetservice-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.1.0/24"]

  private_link_service_network_policies_enabled = false
}

resource "azurerm_subnet" "endpoint" {
  name                 = "acctestsnetendpoint-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.2.0/24"]

  private_endpoint_network_policies = "Disabled"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-pe-server-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "GP_Gen5_4"

  storage_mb                   = 5120
  backup_retention_days        = 7
  geo_redundant_backup_enabled = false
  auto_grow_enabled            = true

  administrator_login          = "psqladmin"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement_enabled      = true
}

resource "azurerm_private_dns_zone" "finance" {
  name                = "privatelink.postgres.database.azure.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-privatelink-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.endpoint.id

  private_dns_zone_group {
    name                 = "acctest-dzg-%d"
    private_dns_zone_ids = [azurerm_private_dns_zone.finance.id]
  }

  private_service_connection {
    name                           = "acctest-privatelink-psc-%d"
    private_connection_resource_id = azurerm_postgresql_server.test.id
    subresource_names              = ["postgresqlServer"]
    is_manual_connection           = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (PrivateEndpointResource) privateDnsZoneGroupRemove(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-privatelink-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "service" {
  name                 = "acctestsnetservice-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.1.0/24"]

  private_link_service_network_policies_enabled = false
}

resource "azurerm_subnet" "endpoint" {
  name                 = "acctestsnetendpoint-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.2.0/24"]

  private_endpoint_network_policies = "Disabled"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-pe-server-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "GP_Gen5_4"

  storage_mb                   = 5120
  backup_retention_days        = 7
  geo_redundant_backup_enabled = false
  auto_grow_enabled            = true

  administrator_login          = "psqladmin"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement_enabled      = true
}

resource "azurerm_private_dns_zone" "finance" {
  name                = "privatelink.postgres.database.azure.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-privatelink-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.endpoint.id

  private_service_connection {
    name                           = "acctest-privatelink-psc-%d"
    private_connection_resource_id = azurerm_postgresql_server.test.id
    subresource_names              = ["postgresqlServer"]
    is_manual_connection           = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (PrivateEndpointResource) privateDnsZoneGroupUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-privatelink-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "service" {
  name                 = "acctestsnetservice-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.1.0/24"]

  private_link_service_network_policies_enabled = false
}

resource "azurerm_subnet" "endpoint" {
  name                 = "acctestsnetendpoint-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.2.0/24"]

  private_endpoint_network_policies = "Disabled"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-pe-server-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "GP_Gen5_4"

  storage_mb                   = 5120
  backup_retention_days        = 7
  geo_redundant_backup_enabled = false
  auto_grow_enabled            = true

  administrator_login          = "psqladmin"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement_enabled      = true
}

resource "azurerm_private_dns_zone" "finance" {
  name                = "privatelink.postgres.database.azure.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_zone" "sales" {
  name                = "acceptance.pdz.%d"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-privatelink-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.endpoint.id

  private_dns_zone_group {
    name                 = "acctest-dzg-%d"
    private_dns_zone_ids = [azurerm_private_dns_zone.sales.id, azurerm_private_dns_zone.finance.id]
  }

  private_service_connection {
    name                           = "acctest-privatelink-psc-%d"
    private_connection_resource_id = azurerm_postgresql_server.test.id
    subresource_names              = ["postgresqlServer"]
    is_manual_connection           = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (PrivateEndpointResource) privateDnsZoneGroupIdsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-privatelink-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "service" {
  name                 = "acctestsnetservice-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.1.0/24"]

  private_link_service_network_policies_enabled = false
}

resource "azurerm_subnet" "endpoint" {
  name                 = "acctestsnetendpoint-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.2.0/24"]

  private_endpoint_network_policies = "Disabled"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-pe-server-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "GP_Gen5_4"

  storage_mb                   = 5120
  backup_retention_days        = 7
  geo_redundant_backup_enabled = false
  auto_grow_enabled            = true

  administrator_login          = "psqladmin"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement_enabled      = true
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG-privatelink2-%d"
  location = azurerm_resource_group.test.location
}

resource "azurerm_private_dns_zone" "finance2" {
  name                = "privatelink.postgres.database.azure.com"
  resource_group_name = azurerm_resource_group.test2.name
}

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-privatelink-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.endpoint.id

  private_dns_zone_group {
    name                 = "acctest-dzg-%d"
    private_dns_zone_ids = [azurerm_private_dns_zone.finance2.id]
  }

  private_service_connection {
    name                           = "acctest-privatelink-psc-%d"
    private_connection_resource_id = azurerm_postgresql_server.test.id
    subresource_names              = ["postgresqlServer"]
    is_manual_connection           = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r PrivateEndpointResource) staticIpAddress(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-privatelink-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.endpoint.id

  private_service_connection {
    name                           = azurerm_private_link_service.test.name
    is_manual_connection           = false
    private_connection_resource_id = azurerm_private_link_service.test.id
  }

  ip_configuration {
    name               = "acctest-ip-privatelink-%[2]d"
    private_ip_address = "10.5.2.47"
  }
}
`, r.template(data, r.serviceAutoApprove(data)), data.RandomInteger)
}

func (PrivateEndpointResource) privateDnsZoneGroupRename(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-privatelink-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "service" {
  name                 = "acctestsnetservice-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.1.0/24"]

  private_link_service_network_policies_enabled = false
}

resource "azurerm_subnet" "endpoint" {
  name                 = "acctestsnetendpoint-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.2.0/24"]

  private_endpoint_network_policies = "Disabled"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-pe-server-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "GP_Gen5_4"

  storage_mb                   = 5120
  backup_retention_days        = 7
  geo_redundant_backup_enabled = false
  auto_grow_enabled            = true

  administrator_login          = "psqladmin"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement_enabled      = true
}

resource "azurerm_private_dns_zone" "finance" {
  name                = "privatelink.postgres.database.azure.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-privatelink-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.endpoint.id

  private_dns_zone_group {
    name                 = "acctest-dzg-rn-%d"
    private_dns_zone_ids = [azurerm_private_dns_zone.finance.id]
  }

  private_service_connection {
    name                           = "acctest-privatelink-psc-%d"
    private_connection_resource_id = azurerm_postgresql_server.test.id
    subresource_names              = ["postgresqlServer"]
    is_manual_connection           = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r PrivateEndpointResource) privateConnectionAlias(data acceptance.TestData, withTags bool) string {
	tags := `
  tags = {
    env = "TEST"
  }
`
	if !withTags {
		tags = ""
	}
	return fmt.Sprintf(`
%s

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-privatelink-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = azurerm_subnet.endpoint.id

  private_service_connection {
    name                              = azurerm_private_link_service.test.name
    is_manual_connection              = true
    private_connection_resource_alias = azurerm_private_link_service.test.alias
    request_message                   = "test"
  }
%s
}
`, r.template(data, r.serviceAutoApprove(data)), data.RandomInteger, tags)
}

func (r PrivateEndpointResource) multipleInstances(data acceptance.TestData, count int, useAlias bool) string {
	privateConnectionAssignment := "private_connection_resource_id = azurerm_private_link_service.test.id"
	if useAlias {
		privateConnectionAssignment = `private_connection_resource_alias = azurerm_private_link_service.test.alias
                                       request_message                   = "test"`
	}

	return fmt.Sprintf(`
%s

resource "azurerm_private_endpoint" "test" {
  count               = %d
  name                = "acctest-privatelink-%d-${count.index}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = azurerm_subnet.endpoint.id

  private_service_connection {
    name                 = azurerm_private_link_service.test.name
    is_manual_connection = %t
    %s
  }
}
`, r.template(data, r.serviceAutoApprove(data)), count, data.RandomInteger, useAlias, privateConnectionAssignment)
}

func (r PrivateEndpointResource) recoveryServiceVaultWithMultiIpConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

locals {
  ip_configs = {
    "SiteRecovery-prot2" = "10.5.2.24"
    "SiteRecovery-srs1"  = "10.5.2.25"
    "SiteRecovery-id1"   = "10.5.2.26"
    "SiteRecovery-tel1"  = "10.5.2.27"
    "SiteRecovery-rcm1"  = "10.5.2.28"
  }
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-vault-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  soft_delete_enabled = false

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-privatelink-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = azurerm_subnet.endpoint.id

  private_service_connection {
    name                           = "acctest-privatelink-%[2]d"
    is_manual_connection           = false
    subresource_names              = ["AzureSiteRecovery"]
    private_connection_resource_id = azurerm_recovery_services_vault.test.id
  }

  dynamic "ip_configuration" {
    for_each = local.ip_configs

    content {
      name               = ip_configuration.key
      private_ip_address = ip_configuration.value
      subresource_name   = "AzureSiteRecovery"
      member_name        = ip_configuration.key
    }
  }
}
`, r.template(data, r.serviceAutoApprove(data)), data.RandomInteger)
}
