package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/reachabilityanalysisintents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ManagerVerifierWorkspaceReachabilityAnalysisIntentResource struct{}

func testAccNetworkManagerVerifierWorkspaceReachabilityAnalysisIntent_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_verifier_workspace_reachability_analysis_intent", "test")
	r := ManagerVerifierWorkspaceReachabilityAnalysisIntentResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkManagerVerifierWorkspaceReachabilityAnalysisIntent_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_verifier_workspace_reachability_analysis_intent", "test")
	r := ManagerVerifierWorkspaceReachabilityAnalysisIntentResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func testAccNetworkManagerVerifierWorkspaceReachabilityAnalysisIntent_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_verifier_workspace_reachability_analysis_intent", "test")
	r := ManagerVerifierWorkspaceReachabilityAnalysisIntentResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ManagerVerifierWorkspaceReachabilityAnalysisIntentResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := reachabilityanalysisintents.ParseReachabilityAnalysisIntentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Network.ReachabilityAnalysisIntents.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ManagerVerifierWorkspaceReachabilityAnalysisIntentResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {}
}

resource "azurerm_network_manager_verifier_workspace_reachability_analysis_intent" "test" {
  name                    = "acctest-vw-%[2]d"
  verifier_workspace_id   = azurerm_network_manager_verifier_workspace.test.id
  source_resource_id      = azurerm_linux_virtual_machine.test.id
  destination_resource_id = azurerm_linux_virtual_machine.test.id
  ip_traffic {
    source_ips        = ["10.0.2.0"]
    source_ports      = ["80"]
    destination_ips   = ["10.0.3.0"]
    destination_ports = ["*"]
    protocols         = ["Any"]
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ManagerVerifierWorkspaceReachabilityAnalysisIntentResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_verifier_workspace_reachability_analysis_intent" "import" {
  name                    = azurerm_network_manager_verifier_workspace_reachability_analysis_intent.test.name
  verifier_workspace_id   = azurerm_network_manager_verifier_workspace.test.id
  source_resource_id      = azurerm_linux_virtual_machine.test.id
  destination_resource_id = azurerm_linux_virtual_machine.test.id
  ip_traffic {
    source_ips        = azurerm_network_manager_verifier_workspace_reachability_analysis_intent.test.ip_traffic[0].source_ips
    source_ports      = azurerm_network_manager_verifier_workspace_reachability_analysis_intent.test.ip_traffic[0].source_ports
    destination_ips   = azurerm_network_manager_verifier_workspace_reachability_analysis_intent.test.ip_traffic[0].destination_ips
    destination_ports = azurerm_network_manager_verifier_workspace_reachability_analysis_intent.test.ip_traffic[0].destination_ports
    protocols         = azurerm_network_manager_verifier_workspace_reachability_analysis_intent.test.ip_traffic[0].protocols
  }
}
`, r.basic(data))
}

func (r ManagerVerifierWorkspaceReachabilityAnalysisIntentResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {}
}

resource "azurerm_network_manager_verifier_workspace_reachability_analysis_intent" "test" {
  name                    = "acctest-vw-%[2]d"
  verifier_workspace_id   = azurerm_network_manager_verifier_workspace.test.id
  source_resource_id      = azurerm_linux_virtual_machine.test.id
  destination_resource_id = azurerm_linux_virtual_machine.test.id
  description             = "test"
  ip_traffic {
    source_ips        = ["10.0.2.1", "10.0.2.3"]
    source_ports      = ["80", "88", "100-120"]
    destination_ips   = ["10.0.2.2", "10.0.2.5"]
    destination_ports = ["60", "89"]
    protocols         = ["UDP", "ICMP"]
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ManagerVerifierWorkspaceReachabilityAnalysisIntentResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-manager-vw-%[1]d"
  location = "%[2]s"
}

data "azurerm_subscription" "current" {}

resource "azurerm_network_manager" "test" {
  name                = "acctest-nm-vw-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  scope {
    subscription_ids = [data.azurerm_subscription.current.id]
  }
  scope_accesses = ["Connectivity"]
}

resource "azurerm_network_manager_verifier_workspace" "test" {
  name               = "acctest-vw-%[1]d"
  network_manager_id = azurerm_network_manager.test.id
  location           = azurerm_resource_group.test.location
}

resource "azurerm_virtual_network" "test" {
  name                = "test-network"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "test" {
  name                = "test-nic"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "test" {
  name                            = "test-machine"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_B1ls"
  admin_username                  = "adminuser"
  admin_password                  = "P@ssw0rd1234!"
  disable_password_authentication = false
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
