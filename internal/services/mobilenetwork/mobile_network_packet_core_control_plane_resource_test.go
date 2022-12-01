package mobilenetwork_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-04-01-preview/packetcorecontrolplane"
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

func TestAccMobileNetworkPacketCoreControlPlane_withAKSHCI(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_packet_core_control_plane", "test")
	r := MobileNetworkPacketCoreControlPlaneResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withAKSHCI(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMobileNetworkPacketCoreControlPlane_withIneropSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_packet_core_control_plane", "test")
	r := MobileNetworkPacketCoreControlPlaneResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withInteropSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMobileNetworkPacketCoreControlPlane_withCertificateUserAssignedIdentity(t *testing.T) {
	t.Skip("Skipping as the service is still in progress.")
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
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-mn-rg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mobile_network" "test" {
  name                = "acctest-mn-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[2]s"
  mobile_country_code = "001"
  mobile_network_code = "01"
}

`, data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkPacketCoreControlPlaneResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_mobile_network_packet_core_control_plane" "test" {
  name                = "acctest-mnpccp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  sku                 = "EvaluationPackage"
  mobile_network_id   = azurerm_mobile_network.test.id

  control_plane_access_interface {
    name         = "default-interface"
    ipv4_address = "192.168.1.199"
    ipv4_gateway = "192.168.1.1"
    ipv4_subnet  = "192.168.1.0/25"
  }

  platform {
    type = "BaseVM"
  }

}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkPacketCoreControlPlaneResource) withAKSHCI(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s
resource "azurerm_databox_edge_device" "test" {
  name                = "acct%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku_name = "EdgeP_Base-Standard"
}

resource "azurerm_mobile_network_packet_core_control_plane" "test" {
  name                = "acctest-mnpccp-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
  sku                 = "EvaluationPackage"
  mobile_network_id   = azurerm_mobile_network.test.id

  control_plane_access_interface {
    name         = "default-interface"
    ipv4_address = "192.168.1.199"
    ipv4_gateway = "192.168.1.1"
    ipv4_subnet  = "192.168.1.0/25"
  }

  platform {
    type           = "AKS-HCI"
    edge_device_id = azurerm_databox_edge_device.test.id
  }

}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkPacketCoreControlPlaneResource) withCertificateUserAssignedIdentity(data acceptance.TestData) string {
	template := r.template(data)
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
  name                                     = "acctest-mnpccp-%[2]d"
  resource_group_name                      = azurerm_resource_group.test.name
  location                                 = "%[3]s"
  sku                                      = "EvaluationPackage"
  core_network_technology                  = "5GC"
  mobile_network_id                        = azurerm_mobile_network.test.id
  local_diagnostics_access_certificate_url = azurerm_key_vault_certificate.test.versionless_secret_id

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }

  control_plane_access_interface {
    name         = "default-interface"
    ipv4_address = "192.168.1.199"
    ipv4_gateway = "192.168.1.1"
    ipv4_subnet  = "192.168.1.0/25"
  }

  platform {
    type = "BaseVM"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkPacketCoreControlPlaneResource) withInteropSettings(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_mobile_network_packet_core_control_plane" "test" {
  name                = "acctest-mnpccp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  sku                 = "EvaluationPackage"
  mobile_network_id   = azurerm_mobile_network.test.id

  control_plane_access_interface {
    name         = "default-interface"
    ipv4_address = "192.168.1.199"
    ipv4_gateway = "192.168.1.1"
    ipv4_subnet  = "192.168.1.0/25"
  }

  platform {
    type = "BaseVM"
  }

  interop_settings = jsonencode({
    "mtu" = 1440
  })

}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkPacketCoreControlPlaneResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_mobile_network_packet_core_control_plane" "import" {
  name                = azurerm_mobile_network_packet_core_control_plane.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  sku                 = "EvaluationPackage"
  mobile_network_id   = azurerm_mobile_network.test.id

  control_plane_access_interface {
    name         = "default-interface"
    ipv4_address = "192.168.1.199"
    ipv4_gateway = "192.168.1.1"
    ipv4_subnet  = "192.168.1.0/25"
  }

  platform {
    type = "BaseVM"
  }

}
`, config, data.Locations.Primary)
}

func (r MobileNetworkPacketCoreControlPlaneResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_mobile_network_packet_core_control_plane" "test" {
  name                = "acctest-mnpccp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  sku                 = "EdgeSite4GBPS"
  mobile_network_id   = azurerm_mobile_network.test.id

  control_plane_access_interface {
    name         = "default-interface"
    ipv4_address = "192.168.1.199"
    ipv4_gateway = "192.168.1.1"
    ipv4_subnet  = "192.168.1.0/25"
  }

  platform {
    type = "BaseVM"
  }

  tags = {
    "Environment" = "non-prod"
  }

}
`, template, data.RandomInteger, data.Locations.Primary)
}
