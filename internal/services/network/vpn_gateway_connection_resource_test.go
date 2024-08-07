// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type VPNGatewayConnectionResource struct{}

func TestAccVpnGatewayConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway_connection", "test")
	r := VPNGatewayConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVpnGatewayConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway_connection", "test")
	r := VPNGatewayConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVpnGatewayConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway_connection", "test")
	r := VPNGatewayConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVpnGatewayConnection_customRouteTable(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway_connection", "test")
	r := VPNGatewayConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customRouteTable(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.customRouteTableUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVpnGatewayConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway_connection", "test")
	r := VPNGatewayConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccVpnGatewayConnection_updateConnectionMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway_connection", "test")
	r := VPNGatewayConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateConnectionMode(data, "Default"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateConnectionMode(data, "InitiatorOnly"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVpnGatewayConnection_updateTrafficSelectorPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway_connection", "test")
	r := VPNGatewayConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateTrafficSelectorPolicy(data, "10.0.0.0/24", "10.0.1.0/24"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateTrafficSelectorPolicy(data, "10.0.2.0/24", "10.0.3.0/24"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVpnGatewayConnection_natRuleIds(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway_connection", "test")
	r := VPNGatewayConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.natRuleIds(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateNatRuleIds(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVpnGatewayConnection_customBgpAddress(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway_connection", "test")
	r := VPNGatewayConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customBgpAddress(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVpnGatewayConnection_routeMap(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway_connection", "test")
	r := VPNGatewayConnectionResource{}
	nameSuffix := randString()

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.routeMap(data, nameSuffix),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t VPNGatewayConnectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseVPNConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.VirtualWANs.VpnConnectionsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading VPN Gateway Connnection (%s): %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r VPNGatewayConnectionResource) basic(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r VPNGatewayConnectionResource) complete(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r VPNGatewayConnectionResource) customRouteTable(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway_connection" "test" {
  name               = "acctest-VpnGwConn-%[2]d"
  vpn_gateway_id     = azurerm_vpn_gateway.test.id
  remote_vpn_site_id = azurerm_vpn_site.test.id
  routing {
    associated_route_table = azurerm_virtual_hub.test.default_route_table_id

    propagated_route_table {
      route_table_ids = [azurerm_virtual_hub.test.default_route_table_id]
      labels          = ["label1"]
    }
  }
  vpn_link {
    name             = "link1"
    vpn_site_link_id = azurerm_vpn_site.test.link[0].id
    ipsec_policy {
      sa_lifetime_sec          = 300
      sa_data_size_kb          = 0
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
`, r.template(data), data.RandomInteger)
}

func (r VPNGatewayConnectionResource) customRouteTableUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway_connection" "test" {
  name               = "acctest-VpnGwConn-%[2]d"
  vpn_gateway_id     = azurerm_vpn_gateway.test.id
  remote_vpn_site_id = azurerm_vpn_site.test.id
  routing {
    associated_route_table = azurerm_virtual_hub.test.default_route_table_id

    propagated_route_table {
      route_table_ids = [azurerm_virtual_hub.test.default_route_table_id]
      labels          = ["label2"]
    }
  }
  vpn_link {
    name             = "link1"
    vpn_site_link_id = azurerm_vpn_site.test.link[0].id
    ipsec_policy {
      sa_lifetime_sec          = 300
      sa_data_size_kb          = 0
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
`, r.template(data), data.RandomInteger)
}

func (r VPNGatewayConnectionResource) requiresImport(data acceptance.TestData) string {
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
`, r.basic(data))
}

func (r VPNGatewayConnectionResource) updateConnectionMode(data acceptance.TestData, connectionMode string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway_connection" "test" {
  name               = "acctest-VpnGwConn-%[2]d"
  vpn_gateway_id     = azurerm_vpn_gateway.test.id
  remote_vpn_site_id = azurerm_vpn_site.test.id
  vpn_link {
    name             = "link1"
    vpn_site_link_id = azurerm_vpn_site.test.link[0].id
    connection_mode  = "%s"
  }
  vpn_link {
    name             = "link2"
    vpn_site_link_id = azurerm_vpn_site.test.link[1].id
  }
}
`, r.template(data), data.RandomInteger, connectionMode)
}

func (r VPNGatewayConnectionResource) updateTrafficSelectorPolicy(data acceptance.TestData, localAddressRange string, remoteAddressRange string) string {
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
  traffic_selector_policy {
    local_address_ranges  = ["%s"]
    remote_address_ranges = ["%s"]
  }
}
`, r.template(data), data.RandomInteger, localAddressRange, remoteAddressRange)
}

func (r VPNGatewayConnectionResource) natRuleIds(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway_nat_rule" "test" {
  name           = "acctest-vpngwnatrule-%[2]d"
  vpn_gateway_id = azurerm_vpn_gateway.test.id
  external_mapping {
    address_space = "192.168.21.0/26"
  }

  internal_mapping {
    address_space = "10.4.0.0/26"
  }

  mode = "EgressSnat"
  type = "Static"
}

resource "azurerm_vpn_gateway_nat_rule" "test2" {
  name           = "acctest-vpngwnatrule2-%[2]d"
  vpn_gateway_id = azurerm_vpn_gateway.test.id
  external_mapping {
    address_space = "192.168.22.0/26"
  }

  internal_mapping {
    address_space = "10.5.0.0/26"
  }

  mode = "IngressSnat"
  type = "Static"
}

resource "azurerm_vpn_gateway_nat_rule" "test3" {
  name           = "acctest-vpngwnatrule3-%[2]d"
  vpn_gateway_id = azurerm_vpn_gateway.test.id
  external_mapping {
    address_space = "192.168.23.0/26"
  }

  internal_mapping {
    address_space = "10.6.0.0/26"
  }

  mode = "EgressSnat"
  type = "Static"
}

resource "azurerm_vpn_gateway_nat_rule" "test4" {
  name           = "acctest-vpngwnatrule4-%[2]d"
  vpn_gateway_id = azurerm_vpn_gateway.test.id
  external_mapping {
    address_space = "192.168.24.0/26"
  }

  internal_mapping {
    address_space = "10.7.0.0/26"
  }

  mode = "IngressSnat"
  type = "Static"
}

resource "azurerm_vpn_gateway_connection" "test" {
  name               = "acctest-VpnGwConn-%[2]d"
  vpn_gateway_id     = azurerm_vpn_gateway.test.id
  remote_vpn_site_id = azurerm_vpn_site.test.id

  vpn_link {
    name                = "link1"
    vpn_site_link_id    = azurerm_vpn_site.test.link[0].id
    egress_nat_rule_ids = [azurerm_vpn_gateway_nat_rule.test.id, azurerm_vpn_gateway_nat_rule.test3.id]
  }

  vpn_link {
    name                 = "link2"
    vpn_site_link_id     = azurerm_vpn_site.test.link[1].id
    ingress_nat_rule_ids = [azurerm_vpn_gateway_nat_rule.test2.id, azurerm_vpn_gateway_nat_rule.test4.id]
  }

  lifecycle {
    create_before_destroy = true
  }
}
`, r.template(data), data.RandomInteger)
}

func (r VPNGatewayConnectionResource) updateNatRuleIds(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway_nat_rule" "test" {
  name           = "acctest-vpngwnatrule-%[2]d"
  vpn_gateway_id = azurerm_vpn_gateway.test.id
  external_mapping {
    address_space = "192.168.21.0/26"
  }

  internal_mapping {
    address_space = "10.4.0.0/26"
  }

  mode = "EgressSnat"
  type = "Static"
}

resource "azurerm_vpn_gateway_nat_rule" "test2" {
  name           = "acctest-vpngwnatrule2-%[2]d"
  vpn_gateway_id = azurerm_vpn_gateway.test.id
  external_mapping {
    address_space = "192.168.22.0/26"
  }

  internal_mapping {
    address_space = "10.5.0.0/26"
  }

  mode = "IngressSnat"
  type = "Static"
}

resource "azurerm_vpn_gateway_connection" "test" {
  name               = "acctest-VpnGwConn-%[2]d"
  vpn_gateway_id     = azurerm_vpn_gateway.test.id
  remote_vpn_site_id = azurerm_vpn_site.test.id

  vpn_link {
    name                = "link1"
    vpn_site_link_id    = azurerm_vpn_site.test.link[0].id
    egress_nat_rule_ids = [azurerm_vpn_gateway_nat_rule.test.id]
  }

  vpn_link {
    name                 = "link2"
    vpn_site_link_id     = azurerm_vpn_site.test.link[1].id
    ingress_nat_rule_ids = [azurerm_vpn_gateway_nat_rule.test2.id]
  }

  lifecycle {
    create_before_destroy = true
  }
}
`, r.template(data), data.RandomInteger)
}

func (r VPNGatewayConnectionResource) customBgpAddress(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctestVWAN-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_hub" "test" {
  name                = "acctestVHUB-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.0.0/24"
}

resource "azurerm_vpn_gateway" "test" {
  name                = "acctestVPNGW-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_hub_id      = azurerm_virtual_hub.test.id

  bgp_settings {
    asn         = 65515
    peer_weight = 0

    instance_0_bgp_peering_address {
      custom_ips = ["169.254.21.5"]
    }

    instance_1_bgp_peering_address {
      custom_ips = ["169.254.21.10"]
    }
  }
}

resource "azurerm_vpn_site" "test" {
  name                = "acctestVPNSite-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_wan_id      = azurerm_virtual_wan.test.id

  link {
    name       = "link1"
    ip_address = "169.254.21.5"

    bgp {
      asn             = 1234
      peering_address = "169.254.21.5"
    }
  }
}

resource "azurerm_vpn_gateway_connection" "test" {
  name               = "acctestVPNGWConn-%[1]d"
  vpn_gateway_id     = azurerm_vpn_gateway.test.id
  remote_vpn_site_id = azurerm_vpn_site.test.id

  vpn_link {
    name             = "link1"
    vpn_site_link_id = azurerm_vpn_site.test.link[0].id
    bgp_enabled      = true

    custom_bgp_address {
      ip_address          = "169.254.21.5"
      ip_configuration_id = azurerm_vpn_gateway.test.bgp_settings.0.instance_0_bgp_peering_address.0.ip_configuration_id
    }

    custom_bgp_address {
      ip_address          = "169.254.21.10"
      ip_configuration_id = azurerm_vpn_gateway.test.bgp_settings.0.instance_1_bgp_peering_address.0.ip_configuration_id
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r VPNGatewayConnectionResource) routeMap(data acceptance.TestData, nameSuffix string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_route_map" "test" {
  name           = "acctestrm-%[2]s"
  virtual_hub_id = azurerm_virtual_hub.test.id

  rule {
    name                 = "rule1"
    next_step_if_matched = "Continue"

    action {
      type = "Add"

      parameter {
        as_path = ["22334"]
      }
    }

    match_criterion {
      match_condition = "Contains"
      route_prefix    = ["10.0.0.0/8"]
    }
  }
}

resource "azurerm_route_map" "test2" {
  name           = "acctestrmn-%[2]s"
  virtual_hub_id = azurerm_virtual_hub.test.id

  rule {
    name                 = "rule1"
    next_step_if_matched = "Continue"

    action {
      type = "Add"

      parameter {
        as_path = ["22334"]
      }
    }

    match_criterion {
      match_condition = "Contains"
      route_prefix    = ["10.0.0.0/8"]
    }
  }
}

resource "azurerm_vpn_gateway_connection" "test" {
  name               = "acctest-VpnGwConn-%[3]d"
  vpn_gateway_id     = azurerm_vpn_gateway.test.id
  remote_vpn_site_id = azurerm_vpn_site.test.id

  routing {
    associated_route_table = azurerm_virtual_hub.test.default_route_table_id
    inbound_route_map_id   = azurerm_route_map.test.id
    outbound_route_map_id  = azurerm_route_map.test2.id

    propagated_route_table {
      route_table_ids = [azurerm_virtual_hub.test.default_route_table_id]
      labels          = ["label1"]
    }
  }

  vpn_link {
    name             = "link1"
    vpn_site_link_id = azurerm_vpn_site.test.link[0].id

    ipsec_policy {
      sa_lifetime_sec          = 300
      sa_data_size_kb          = 0
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
`, r.template(data), nameSuffix, data.RandomInteger)
}

func (VPNGatewayConnectionResource) template(data acceptance.TestData) string {
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
  address_cidrs       = ["10.0.1.0/24"]

  link {
    name       = "link1"
    ip_address = "10.0.1.1"
  }

  link {
    name       = "link2"
    ip_address = "10.0.1.2"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
