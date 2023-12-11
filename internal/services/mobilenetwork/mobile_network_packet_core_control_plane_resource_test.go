// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mobilenetwork_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/packetcorecontrolplane"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MobileNetworkPacketCoreControlPlaneResource struct{}

func TestAccMobileNetworkPacketCoreControlPlane_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_packet_core_control_plane", "test")
	r := MobileNetworkPacketCoreControlPlaneResource{}
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

func TestAccMobileNetworkPacketCoreControlPlane_withAccessInterface(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_packet_core_control_plane", "test")
	r := MobileNetworkPacketCoreControlPlaneResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withAccessInterface(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMobileNetworkPacketCoreControlPlane_withInteropJSON(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_packet_core_control_plane", "test")
	r := MobileNetworkPacketCoreControlPlaneResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withInteropJson(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMobileNetworkPacketCoreControlPlane_withUeMTU(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_packet_core_control_plane", "test")
	r := MobileNetworkPacketCoreControlPlaneResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withUeMTU(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMobileNetworkPacketCoreControlPlane_withCertificateUserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_packet_core_control_plane", "test")
	r := MobileNetworkPacketCoreControlPlaneResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withCertificateUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMobileNetworkPacketCoreControlPlane_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_packet_core_control_plane", "test")
	r := MobileNetworkPacketCoreControlPlaneResource{}
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

func TestAccMobileNetworkPacketCoreControlPlane_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_packet_core_control_plane", "test")
	r := MobileNetworkPacketCoreControlPlaneResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
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
	})
}

func TestAccMobileNetworkPacketCoreControlPlane_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_packet_core_control_plane", "test")
	r := MobileNetworkPacketCoreControlPlaneResource{}
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

func (r MobileNetworkPacketCoreControlPlaneResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := packetcorecontrolplane.ParsePacketCoreControlPlaneID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.MobileNetwork.PacketCoreControlPlaneClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r MobileNetworkPacketCoreControlPlaneResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_databox_edge_device" "test" {
  name                = "acct%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku_name = "EdgeP_Base-Standard"
}

`, MobileNetworkSiteResource{}.basic(data), data.RandomInteger)
}

func (r MobileNetworkPacketCoreControlPlaneResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
				%s

resource "azurerm_mobile_network_packet_core_control_plane" "test" {
  name                = "acctest-mnpccp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  sku                 = "G0"
  site_ids            = [azurerm_mobile_network_site.test.id]

  local_diagnostics_access {
    authentication_type = "AAD"
  }

  platform {
    type           = "AKS-HCI"
    edge_device_id = azurerm_databox_edge_device.test.id
  }

  depends_on = [azurerm_mobile_network.test]
}
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkPacketCoreControlPlaneResource) withAccessInterface(data acceptance.TestData) string {
	return fmt.Sprintf(`
				%s

resource "azurerm_mobile_network_packet_core_control_plane" "test" {
  name                              = "acctest-mnpccp-%d"
  resource_group_name               = azurerm_resource_group.test.name
  location                          = "%s"
  sku                               = "G0"
  site_ids                          = [azurerm_mobile_network_site.test.id]
  control_plane_access_name         = "default-interface"
  control_plane_access_ipv4_address = "192.168.1.199"
  control_plane_access_ipv4_gateway = "192.168.1.1"
  control_plane_access_ipv4_subnet  = "192.168.1.0/25"

  local_diagnostics_access {
    authentication_type = "AAD"
  }

  platform {
    type           = "AKS-HCI"
    edge_device_id = azurerm_databox_edge_device.test.id
  }

  depends_on = [azurerm_mobile_network.test]
}
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkPacketCoreControlPlaneResource) withCertificateUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

data "azurerm_client_config" "test" {}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-mn-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_key_vault" "test" {
  name                = "acct-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.test.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id               = data.azurerm_client_config.test.tenant_id
    object_id               = data.azurerm_client_config.test.object_id
    secret_permissions      = ["Delete", "Get", "Set", "Purge"]
    certificate_permissions = ["Create", "Delete", "Get", "Import", "Purge"]
  }

  access_policy {
    tenant_id               = data.azurerm_client_config.test.tenant_id
    object_id               = azurerm_user_assigned_identity.test.principal_id
    secret_permissions      = ["Delete", "Get", "Set"]
    certificate_permissions = ["Create", "Delete", "Get", "Import"]
  }

}

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctest-mn-%[2]d"
  key_vault_id = azurerm_key_vault.test.id

  certificate {
    contents = filebase64("testdata/rsa_bundle.pfx")
    password = ""
  }
}

resource "azurerm_mobile_network_packet_core_control_plane" "test" {
  name                              = "acctest-mnpccp-%[2]d"
  resource_group_name               = azurerm_resource_group.test.name
  location                          = "%[3]s"
  sku                               = "G0"
  core_network_technology           = "5GC"
  site_ids                          = [azurerm_mobile_network_site.test.id]
  control_plane_access_name         = "default-interface"
  control_plane_access_ipv4_address = "192.168.1.199"
  control_plane_access_ipv4_gateway = "192.168.1.1"
  control_plane_access_ipv4_subnet  = "192.168.1.0/25"

  local_diagnostics_access {
    authentication_type          = "AAD"
    https_server_certificate_url = azurerm_key_vault_certificate.test.versionless_secret_id
  }

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }

  platform {
    type           = "AKS-HCI"
    edge_device_id = azurerm_databox_edge_device.test.id
  }

  depends_on = [azurerm_mobile_network.test]
}




`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkPacketCoreControlPlaneResource) withInteropJson(data acceptance.TestData) string {
	return fmt.Sprintf(`
				%s

resource "azurerm_mobile_network_packet_core_control_plane" "test" {
  name                              = "acctest-mnpccp-%d"
  resource_group_name               = azurerm_resource_group.test.name
  location                          = "%s"
  sku                               = "G0"
  site_ids                          = [azurerm_mobile_network_site.test.id]
  control_plane_access_name         = "default-interface"
  control_plane_access_ipv4_address = "192.168.1.199"
  control_plane_access_ipv4_gateway = "192.168.1.1"
  control_plane_access_ipv4_subnet  = "192.168.1.0/25"

  local_diagnostics_access {
    authentication_type = "AAD"
  }

  platform {
    type           = "AKS-HCI"
    edge_device_id = azurerm_databox_edge_device.test.id
  }

  interoperability_settings_json = jsonencode({
    "unknownuser-causecode" = "eps-and-non-eps-service-not-allowed-8"
  })

  depends_on = [azurerm_mobile_network.test]
}
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkPacketCoreControlPlaneResource) withUeMTU(data acceptance.TestData) string {
	return fmt.Sprintf(`
				%s

resource "azurerm_mobile_network_packet_core_control_plane" "test" {
  name                              = "acctest-mnpccp-%d"
  resource_group_name               = azurerm_resource_group.test.name
  location                          = "%s"
  sku                               = "G0"
  user_equipment_mtu_in_bytes       = 1600
  site_ids                          = [azurerm_mobile_network_site.test.id]
  control_plane_access_name         = "default-interface"
  control_plane_access_ipv4_address = "192.168.1.199"
  control_plane_access_ipv4_gateway = "192.168.1.1"
  control_plane_access_ipv4_subnet  = "192.168.1.0/25"

  local_diagnostics_access {
    authentication_type = "AAD"
  }

  platform {
    type           = "AKS-HCI"
    edge_device_id = azurerm_databox_edge_device.test.id
  }

  depends_on = [azurerm_mobile_network.test]
}
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkPacketCoreControlPlaneResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
			%s

resource "azurerm_mobile_network_packet_core_control_plane" "import" {
  name                              = azurerm_mobile_network_packet_core_control_plane.test.name
  resource_group_name               = azurerm_resource_group.test.name
  location                          = "%s"
  sku                               = "G0"
  site_ids                          = [azurerm_mobile_network_site.test.id]
  control_plane_access_name         = "default-interface"
  control_plane_access_ipv4_address = "192.168.1.199"
  control_plane_access_ipv4_gateway = "192.168.1.1"
  control_plane_access_ipv4_subnet  = "192.168.1.0/25"

  platform {
    type           = "AKS-HCI"
    edge_device_id = azurerm_databox_edge_device.test.id
  }

  local_diagnostics_access {
    authentication_type = "AAD"
  }

  depends_on = [azurerm_mobile_network.test]
}
`, r.basic(data), data.Locations.Primary)
}

func (r MobileNetworkPacketCoreControlPlaneResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
			%s

resource "azurerm_mobile_network_packet_core_control_plane" "test" {
  name                              = "acctest-mnpccp-%d"
  resource_group_name               = azurerm_resource_group.test.name
  location                          = "%s"
  sku                               = "G0"
  site_ids                          = [azurerm_mobile_network_site.test.id]
  control_plane_access_name         = "default-interface"
  control_plane_access_ipv4_address = "192.168.1.199"
  control_plane_access_ipv4_gateway = "192.168.1.1"
  control_plane_access_ipv4_subnet  = "192.168.1.0/25"

  local_diagnostics_access {
    authentication_type = "AAD"
  }

  platform {
    type           = "AKS-HCI"
    edge_device_id = azurerm_databox_edge_device.test.id
  }

  tags = {
    "update" = "true"
  }

  depends_on = [azurerm_mobile_network.test]
}
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkPacketCoreControlPlaneResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
				%s
data "azurerm_client_config" "test" {}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-mn-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_key_vault" "test" {
  name                = "acct-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.test.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id               = data.azurerm_client_config.test.tenant_id
    object_id               = data.azurerm_client_config.test.object_id
    secret_permissions      = ["Delete", "Get", "Set", "Purge"]
    certificate_permissions = ["Create", "Delete", "Get", "Import", "Purge"]
  }

  access_policy {
    tenant_id               = data.azurerm_client_config.test.tenant_id
    object_id               = azurerm_user_assigned_identity.test.principal_id
    secret_permissions      = ["Delete", "Get", "Set"]
    certificate_permissions = ["Create", "Delete", "Get", "Import"]
  }

}

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctest-mn-%[2]d"
  key_vault_id = azurerm_key_vault.test.id

  certificate {
    contents = filebase64("testdata/rsa_bundle.pfx")
    password = ""
  }
}

resource "azurerm_mobile_network_packet_core_control_plane" "test" {
  name                              = "acctest-mnpccp-%[2]d"
  resource_group_name               = azurerm_resource_group.test.name
  location                          = "%[3]s"
  sku                               = "G0"
  user_equipment_mtu_in_bytes       = 1600
  site_ids                          = [azurerm_mobile_network_site.test.id]
  control_plane_access_name         = "default-interface"
  control_plane_access_ipv4_address = "192.168.1.199"
  control_plane_access_ipv4_gateway = "192.168.1.1"
  control_plane_access_ipv4_subnet  = "192.168.1.0/25"

  interoperability_settings_json = jsonencode({
    "unknownuser-causecode" = "eps-and-non-eps-service-not-allowed-8"
  })

  local_diagnostics_access {
    authentication_type          = "AAD"
    https_server_certificate_url = azurerm_key_vault_certificate.test.versionless_secret_id
  }

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }

  platform {
    type           = "AKS-HCI"
    edge_device_id = azurerm_databox_edge_device.test.id
  }

  depends_on = [azurerm_mobile_network.test]
}
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}
