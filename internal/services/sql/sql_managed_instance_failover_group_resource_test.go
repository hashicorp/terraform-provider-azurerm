package sql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SqlInstanceFailoverGroupResource struct{}

func TestAccAzureRMSqlInstanceFailoverGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_managed_instance_failover_group", "test")
	r := SqlInstanceFailoverGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: SqlManagedInstanceResource{}.dnsZonePartner(data),
		},
		{
			// It speeds up deletion to remove the explicit dependency between the instances
			Config: SqlManagedInstanceResource{}.emptyDnsZonePartner(data),
		},
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: resource.ComposeTestCheckFunc(
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

func TestAccAzureRMSqlInstanceFailoverGroup_change(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_managed_instance_failover_group", "test")
	r := SqlInstanceFailoverGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: SqlManagedInstanceResource{}.dnsZonePartners(data),
		},
		{
			// It speeds up deletion to remove the explicit dependency between the instances
			Config: SqlManagedInstanceResource{}.emptyDnsZonePartners(data),
		},
		{
			Config: r.connectSecondary(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.changeSecondary(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// disconnect
			Config: SqlManagedInstanceResource{}.emptyDnsZonePartners(data),
		},
	})
}

func (r SqlInstanceFailoverGroupResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
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
  name                  = "acctest-fog-%d"
  resource_group_name   = azurerm_resource_group.test_1.name
  location              = azurerm_sql_managed_instance.test_1.location
  managed_instance_name = azurerm_sql_managed_instance.test_1.name

  partner_managed_instances {
    id = azurerm_sql_managed_instance.test_2.id
  }

  read_write_endpoint_failover_policy {
    mode = "Manual"
  }

  depends_on = [
    azurerm_virtual_network_gateway_connection.test_1,
    azurerm_virtual_network_gateway_connection.test_2,
  ]
}
`, SqlManagedInstanceResource{}.emptyDnsZonePartner(data), r.vnetToVnetGateway(data), data.RandomInteger)
}

func (r SqlInstanceFailoverGroupResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

%s

resource "azurerm_sql_managed_instance_failover_group" "test" {
  name                  = "acctest-fog-%d"
  resource_group_name   = azurerm_resource_group.test_1.name
  location              = azurerm_sql_managed_instance.test_1.location
  managed_instance_name = azurerm_sql_managed_instance.test_1.name

  partner_managed_instances {
    id = azurerm_sql_managed_instance.test_2.id
  }

  readonly_endpoint_failover_policy {
    mode = "Enabled"
  }

  read_write_endpoint_failover_policy {
    mode          = "Automatic"
    grace_minutes = 60
  }

  depends_on = [
    azurerm_virtual_network_gateway_connection.test_1,
    azurerm_virtual_network_gateway_connection.test_2,
  ]
}
`, SqlManagedInstanceResource{}.emptyDnsZonePartner(data), r.vnetToVnetGateway(data), data.RandomInteger)
}

func (r SqlInstanceFailoverGroupResource) connectSecondary(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

%s

resource "azurerm_sql_managed_instance_failover_group" "test" {
  name                  = "acctest-fog-%d"
  resource_group_name   = azurerm_resource_group.test_1.name
  location              = azurerm_sql_managed_instance.test_1.location
  managed_instance_name = azurerm_sql_managed_instance.test_1.name

  partner_managed_instances {
    id = azurerm_sql_managed_instance.test_2.id
  }

  read_write_endpoint_failover_policy {
    mode = "Manual"
  }

  depends_on = [
    azurerm_virtual_network_gateway_connection.test_1,
    azurerm_virtual_network_gateway_connection.test_2,
  ]
}
`, SqlManagedInstanceResource{}.emptyDnsZonePartners(data), r.vnetToVnetsGateway(data), data.RandomInteger)
}

func (r SqlInstanceFailoverGroupResource) changeSecondary(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

%s

resource "azurerm_sql_managed_instance_failover_group" "test" {
  name                  = "acctest-fog-%d"
  resource_group_name   = azurerm_resource_group.test_1.name
  location              = azurerm_sql_managed_instance.test_1.location
  managed_instance_name = azurerm_sql_managed_instance.test_1.name

  partner_managed_instances {
    id = azurerm_sql_managed_instance.test_3.id
  }

  read_write_endpoint_failover_policy {
    mode = "Manual"
  }

  depends_on = [
    azurerm_virtual_network_gateway_connection.test_3,
    azurerm_virtual_network_gateway_connection.test_4,
  ]
}
`, SqlManagedInstanceResource{}.emptyDnsZonePartners(data), r.vnetToVnetsGateway(data), data.RandomInteger)
}

func (r SqlInstanceFailoverGroupResource) vnetToVnetGateway(data acceptance.TestData) string {
	return fmt.Sprintf(`
variable "shared_key" {
  default = "s3cr37"
}

resource "azurerm_subnet" "gateway_snet_test_1" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test_1.name
  virtual_network_name = azurerm_virtual_network.test_1.name
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_public_ip" "test_1" {
  name                = "acctest-pip-%[1]d"
  location            = azurerm_resource_group.test_1.location
  resource_group_name = azurerm_resource_group.test_1.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test_1" {
  name                = "acctest-vnetgway-%[1]d"
  location            = azurerm_resource_group.test_1.location
  resource_group_name = azurerm_resource_group.test_1.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "Basic"

  ip_configuration {
    name                          = "vnetGatewayConfig"
    public_ip_address_id          = azurerm_public_ip.test_1.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.gateway_snet_test_1.id
  }
}

resource "azurerm_virtual_network_gateway_connection" "test_1" {
  name                = "acctest-gwc-%[1]d"
  location            = azurerm_resource_group.test_1.location
  resource_group_name = azurerm_resource_group.test_1.name

  type                            = "Vnet2Vnet"
  virtual_network_gateway_id      = azurerm_virtual_network_gateway.test_1.id
  peer_virtual_network_gateway_id = azurerm_virtual_network_gateway.test_2.id

  shared_key = var.shared_key
}

resource "azurerm_subnet" "gateway_snet_test_2" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test_2.name
  virtual_network_name = azurerm_virtual_network.test_2.name
  address_prefix       = "10.1.1.0/24"
}

resource "azurerm_public_ip" "test_2" {
  name                = "acctest-pip2-%[1]d"
  location            = azurerm_resource_group.test_2.location
  resource_group_name = azurerm_resource_group.test_2.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test_2" {
  name                = "acctest-vnetgway2-%[1]d"
  location            = azurerm_resource_group.test_2.location
  resource_group_name = azurerm_resource_group.test_2.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "Basic"

  ip_configuration {
    name                          = "vnetGatewayConfig"
    public_ip_address_id          = azurerm_public_ip.test_2.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.gateway_snet_test_2.id
  }
}

resource "azurerm_virtual_network_gateway_connection" "test_2" {
  name                = "acctest-gwc2-%[1]d"
  location            = azurerm_resource_group.test_2.location
  resource_group_name = azurerm_resource_group.test_2.name

  type                            = "Vnet2Vnet"
  virtual_network_gateway_id      = azurerm_virtual_network_gateway.test_2.id
  peer_virtual_network_gateway_id = azurerm_virtual_network_gateway.test_1.id

  shared_key = var.shared_key
}
`, data.RandomInteger)
}

func (r SqlInstanceFailoverGroupResource) vnetToVnetsGateway(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_subnet" "gateway_snet_test_3" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test_3.name
  virtual_network_name = azurerm_virtual_network.test_3.name
  address_prefix       = "10.2.1.0/24"
}

resource "azurerm_public_ip" "test_3" {
  name                = "acctest-pip3-%[2]d"
  location            = azurerm_resource_group.test_3.location
  resource_group_name = azurerm_resource_group.test_3.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test_3" {
  name                = "acctest-vnetgway3-%[2]d"
  location            = azurerm_resource_group.test_3.location
  resource_group_name = azurerm_resource_group.test_3.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "Basic"

  ip_configuration {
    name                          = "vnetGatewayConfig"
    public_ip_address_id          = azurerm_public_ip.test_3.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.gateway_snet_test_3.id
  }
}

resource "azurerm_virtual_network_gateway_connection" "test_3" {
  name                = "acctest-gwc3-%[2]d"
  location            = azurerm_resource_group.test_3.location
  resource_group_name = azurerm_resource_group.test_3.name

  type                            = "Vnet2Vnet"
  virtual_network_gateway_id      = azurerm_virtual_network_gateway.test_3.id
  peer_virtual_network_gateway_id = azurerm_virtual_network_gateway.test_1.id

  shared_key = var.shared_key
}

resource "azurerm_virtual_network_gateway_connection" "test_4" {
  name                = "acctest-gwc4-%[2]d"
  location            = azurerm_resource_group.test_1.location
  resource_group_name = azurerm_resource_group.test_1.name

  type                            = "Vnet2Vnet"
  virtual_network_gateway_id      = azurerm_virtual_network_gateway.test_1.id
  peer_virtual_network_gateway_id = azurerm_virtual_network_gateway.test_3.id

  shared_key = var.shared_key
}
`, r.vnetToVnetGateway(data), data.RandomInteger)
}
