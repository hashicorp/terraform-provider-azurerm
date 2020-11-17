package tests

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMVpnGatewayConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVpnGatewayConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVpnGatewayConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVpnGatewayConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVpnGatewayConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVpnGatewayConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVpnGatewayConnection_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVpnGatewayConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVpnGatewayConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVpnGatewayConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVpnGatewayConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVpnGatewayConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMVpnGatewayConnection_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVpnGatewayConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMVpnGatewayConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVpnGatewayConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVpnGatewayConnection_customRouteTable(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVpnGatewayConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVpnGatewayConnection_customRouteTable(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVpnGatewayConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMVpnGatewayConnection_customRouteTableUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVpnGatewayConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVpnGatewayConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVpnGatewayConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVpnGatewayConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVpnGatewayConnectionExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMVpnGatewayConnection_requiresImport),
		},
	})
}

func testCheckAzureRMVpnGatewayConnectionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VpnConnectionsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Vpn Gateway Connection not found: %s", resourceName)
		}

		id, err := parse.VPNGatewayConnectionID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Gateway, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Vpn Gateway Connection %q (Resource Group %q / VPN Gateway %q) does not exist", id.Name, id.ResourceGroup, id.Gateway)
			}
			return fmt.Errorf("Getting on Network.VpnConnetions: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMVpnGatewayConnectionDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VpnConnectionsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_vpn_gateway_connection" {
			continue
		}

		id, err := parse.VPNGatewayConnectionID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Gateway, id.Name)
		if err == nil {
			return fmt.Errorf("Network.VpnConnetions still exists")
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Getting on Network.VpnConnetions: %+v", err)
		}
		return nil
	}

	return nil
}

func testAccAzureRMVpnGatewayConnection_basic(data acceptance.TestData) string {
	template := testAccAzureRMVpnGatewayConnection_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway_connection" "test" {
  name               = "acctest-VpnGwConn-%[2]d"
  vpn_gateway_id     = azurerm_vpn_gateway.test.id
  remote_vpn_site_id = azurerm_vpn_site.test.id
  vpn_link {
    name             = "link1"
    vpn_site_link_id = azurerm_vpn_site.test.link[0].id
  }
  vpn_link {
    name             = "link2"
    vpn_site_link_id = azurerm_vpn_site.test.link[1].id
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMVpnGatewayConnection_complete(data acceptance.TestData) string {
	template := testAccAzureRMVpnGatewayConnection_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway_connection" "test" {
  name               = "acctest-VpnGwConn-%[2]d"
  vpn_gateway_id     = azurerm_vpn_gateway.test.id
  remote_vpn_site_id = azurerm_vpn_site.test.id
  vpn_link {
    name             = "link1"
    vpn_site_link_id = azurerm_vpn_site.test.link[0].id
  }
  vpn_link {
    name             = "link2"
    vpn_site_link_id = azurerm_vpn_site.test.link[1].id
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMVpnGatewayConnection_customRouteTable(data acceptance.TestData) string {
	template := testAccAzureRMVpnGatewayConnection_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_route_table" "test" {
  name           = "acctest-RouteTable-%[2]d"
  virtual_hub_id = azurerm_virtual_hub.test.id
}

resource "azurerm_vpn_gateway_connection" "test" {
  name               = "acctest-VpnGwConn-%[2]d"
  vpn_gateway_id     = azurerm_vpn_gateway.test.id
  remote_vpn_site_id = azurerm_vpn_site.test.id
  routing {
    associated_route_table  = azurerm_virtual_hub_route_table.test.id
    propagated_route_tables = [azurerm_virtual_hub_route_table.test.id]
  }
  vpn_link {
    name             = "link1"
    vpn_site_link_id = azurerm_vpn_site.test.link[0].id
    ipsec_policy {
      sa_lifetime_sec          = 300
      sa_data_size_kb          = 1024
      encryption_algorithm     = "AES256"
      integrity_algorithm      = "SHA256"
      ike_encryption_algorithm = "AES128"
      ike_integrity_algorithm  = "SHA256"
      dh_group                 = "DHGroup14"
      pfs_group                = "PFS14"
    }
    bandwidth_mbps                        = 30
    protocol                              = "IKEv2"
    ratelimit_enabled                     = true
    route_weight                          = 2
    shared_key                            = "secret"
    local_azure_ip_address_enabled        = true
    policy_based_traffic_selector_enabled = true
  }

  vpn_link {
    name             = "link3"
    vpn_site_link_id = azurerm_vpn_site.test.link[1].id
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMVpnGatewayConnection_customRouteTableUpdate(data acceptance.TestData) string {
	template := testAccAzureRMVpnGatewayConnection_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_route_table" "test" {
  name           = "acctest-RouteTable-%[2]d"
  virtual_hub_id = azurerm_virtual_hub.test.id
}

resource "azurerm_virtual_hub_route_table" "test2" {
  name           = "acctest-RouteTable-%[2]d-2"
  virtual_hub_id = azurerm_virtual_hub.test.id
}

resource "azurerm_vpn_gateway_connection" "test" {
  name               = "acctest-VpnGwConn-%[2]d"
  vpn_gateway_id     = azurerm_vpn_gateway.test.id
  remote_vpn_site_id = azurerm_vpn_site.test.id
  routing {
    associated_route_table  = azurerm_virtual_hub_route_table.test2.id
    propagated_route_tables = [azurerm_virtual_hub_route_table.test2.id]
  }
  vpn_link {
    name             = "link1"
    vpn_site_link_id = azurerm_vpn_site.test.link[0].id
    ipsec_policy {
      sa_lifetime_sec          = 300
      sa_data_size_kb          = 1024
      encryption_algorithm     = "AES256"
      integrity_algorithm      = "SHA256"
      ike_encryption_algorithm = "AES128"
      ike_integrity_algorithm  = "SHA256"
      dh_group                 = "DHGroup14"
      pfs_group                = "PFS14"
    }
    bandwidth_mbps                        = 30
    protocol                              = "IKEv2"
    ratelimit_enabled                     = true
    route_weight                          = 2
    shared_key                            = "secret"
    local_azure_ip_address_enabled        = true
    policy_based_traffic_selector_enabled = true
  }

  vpn_link {
    name             = "link3"
    vpn_site_link_id = azurerm_vpn_site.test.link[1].id
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMVpnGatewayConnection_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMVpnGatewayConnection_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway_connection" "import" {
  name               = azurerm_vpn_gateway_connection.test.name
  vpn_gateway_id     = azurerm_vpn_gateway_connection.test.vpn_gateway_id
  remote_vpn_site_id = azurerm_vpn_gateway_connection.test.remote_vpn_site_id
  dynamic "vpn_link" {
    for_each = azurerm_vpn_gateway_connection.test.vpn_link
    iterator = v
    content {
      name             = v.value["name"]
      vpn_site_link_id = v.value["vpn_site_link_id"]
    }
  }
}
`, template)
}

func testAccAzureRMVpnGatewayConnection_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vpn-%[1]d"
  location = "%[2]s"
}


resource "azurerm_virtual_wan" "test" {
  name                = "acctest-vwan-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_hub" "test" {
  name                = "acctest-vhub-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.0.0/24"
}

resource "azurerm_vpn_gateway" "test" {
  name                = "acctest-vpngw-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_hub_id      = azurerm_virtual_hub.test.id
}

resource "azurerm_vpn_site" "test" {
  name                = "acctest-vpnsite-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_wan_id      = azurerm_virtual_wan.test.id
  link {
    name       = "link1"
    ip_address = "10.0.0.1"
  }
  link {
    name       = "link2"
    ip_address = "10.0.0.2"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
