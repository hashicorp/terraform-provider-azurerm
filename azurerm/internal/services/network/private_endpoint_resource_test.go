package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type PrivateEndpointResource struct {
}

func TestAccPrivateEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint", "test")
	r := PrivateEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subnet_id").Exists(),
				check.That(data.ResourceName).Key("private_service_connection.0.private_ip_address").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPrivateEndpoint_updateTag(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint", "test")
	r := PrivateEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withTag(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPrivateEndpoint_requestMessage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint", "test")
	r := PrivateEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.requestMessage(data, "CATS: ALL YOUR BASE ARE BELONG TO US."),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subnet_id").Exists(),
				check.That(data.ResourceName).Key("private_service_connection.0.request_message").HasValue("CATS: ALL YOUR BASE ARE BELONG TO US."),
			),
		},
		data.ImportStep(),
		{
			Config: r.requestMessage(data, "CAPTAIN: WHAT YOU SAY!!"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subnet_id").Exists(),
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.privateDnsZoneGroup(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.privateDnsZoneGroup(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_dns_zone_group.0.private_dns_zone_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("private_dns_zone_configs.#").HasValue("1"),
				check.That(data.ResourceName).Key("private_dns_zone_group.#").HasValue("1"),
			),
		},
		data.ImportStep("private_dns_zone_configs", "private_dns_zone_group"),
		{
			Config: r.privateDnsZoneGroupRename(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.privateDnsZoneGroup(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_dns_zone_group.0.private_dns_zone_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("private_dns_zone_configs.#").HasValue("1"),
				check.That(data.ResourceName).Key("private_dns_zone_group.#").HasValue("1"),
			),
		},
		data.ImportStep("private_dns_zone_configs", "private_dns_zone_group"),
		{
			Config: r.privateDnsZoneGroupUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_dns_zone_group.0.private_dns_zone_ids.#").HasValue("2"),
				check.That(data.ResourceName).Key("private_dns_zone_configs.#").HasValue("2"),
				check.That(data.ResourceName).Key("private_dns_zone_group.#").HasValue("1"),
			),
		},
		data.ImportStep("private_dns_zone_configs", "private_dns_zone_group"),
	})
}

func TestAccPrivateEndpoint_privateDnsZoneRemove(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint", "test")
	r := PrivateEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.privateDnsZoneGroup(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_dns_zone_group.0.private_dns_zone_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("private_dns_zone_configs.#").HasValue("1"),
				check.That(data.ResourceName).Key("private_dns_zone_group.#").HasValue("1"),
			),
		},
		data.ImportStep("private_dns_zone_configs", "private_dns_zone_group"),
		{
			Config: r.privateDnsZoneGroupRemove(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("private_dns_zone_configs", "private_dns_zone_group"),
		{
			Config: r.privateDnsZoneGroup(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_dns_zone_group.0.private_dns_zone_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("private_dns_zone_configs.#").HasValue("1"),
				check.That(data.ResourceName).Key("private_dns_zone_group.#").HasValue("1"),
			),
		},
		data.ImportStep("private_dns_zone_configs", "private_dns_zone_group"),
	})
}

func (t PrivateEndpointResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.PrivateEndpointID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.PrivateEndpointClient.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading Private Endpoint (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
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

  enforce_private_link_service_network_policies = true
}

resource "azurerm_subnet" "endpoint" {
  name                 = "acctestsnetendpoint-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.2.0/24"]

  enforce_private_link_endpoint_network_policies = true
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

  enforce_private_link_service_network_policies = true
}

resource "azurerm_subnet" "endpoint" {
  name                 = "acctestsnetendpoint-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.2.0/24"]

  enforce_private_link_endpoint_network_policies = true
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

  administrator_login          = "psqladminun"
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

  enforce_private_link_service_network_policies = true
}

resource "azurerm_subnet" "endpoint" {
  name                 = "acctestsnetendpoint-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.2.0/24"]

  enforce_private_link_endpoint_network_policies = true
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

  administrator_login          = "psqladminun"
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

  enforce_private_link_service_network_policies = true
}

resource "azurerm_subnet" "endpoint" {
  name                 = "acctestsnetendpoint-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.2.0/24"]

  enforce_private_link_endpoint_network_policies = true
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

  administrator_login          = "psqladminun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement_enabled      = true
}

resource "azurerm_private_dns_zone" "finance" {
  name                = "privatelink.postgres.database.azure.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_zone" "sales" {
  name                = "acctest.pdz.%d"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_endpoint" "test" {
  name                = "acctest.privatelink.%d"
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

  enforce_private_link_service_network_policies = true
}

resource "azurerm_subnet" "endpoint" {
  name                 = "acctestsnetendpoint-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.2.0/24"]

  enforce_private_link_endpoint_network_policies = true
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

  administrator_login          = "psqladminun"
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
