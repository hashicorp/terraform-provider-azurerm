package arckubernetes_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2024-01-01/connectedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/arckubernetes/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// https://learn.microsoft.com/en-us/azure/aks/hybrid/aks-create-clusters-cli?toc=%2Fazure-stack%2Fhci%2Ftoc.json&bc=%2Fazure-stack%2Fbreadcrumb%2Ftoc.json#before-you-begin
// The resource can only be created on the customlocation generated after HCI deployment
const (
	customLocationIdEnv = "ARM_TEST_STACK_HCI_CUSTOM_LOCATION_ID"
)

type ArcKubernetesProvisionedClusterInstanceResource struct{}

func TestAccArcKubernetesProvisionedClusterInstance_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_provisioned_cluster_instance", "test")
	r := ArcKubernetesProvisionedClusterInstanceResource{}
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

func TestAccArcKubernetesProvisionedClusterInstance_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_provisioned_cluster_instance", "test")
	r := ArcKubernetesProvisionedClusterInstanceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError(data.ResourceType),
		},
	})
}

func (r ArcKubernetesProvisionedClusterInstanceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	client := clients.ArcKubernetes.ProvisionedClusterInstancesClient

	id, err := parse.ArcKubernetesProvisionedClusterInstanceID(state.ID)
	if err != nil {
		return nil, err
	}

	connectedClusterId := connectedclusters.NewConnectedClusterID(id.SubscriptionId, id.ResourceGroup, id.ConnectedClusterName)
	scopeId := commonids.NewScopeID(connectedClusterId.ID())

	resp, err := client.ProvisionedClusterInstancesGet(ctx, scopeId)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ArcKubernetesProvisionedClusterInstanceResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_arc_kubernetes_provisioned_cluster_instance" "test" {
  cluster_id = azurerm_arc_kubernetes_cluster.test.id
  custom_location_id = "%[3]s"
  kubernetes_version = "1.28.5"

  agent_pool_profile {
    count = 1
    name = "nodepool1"
    os_sku = "CBLMariner"
    os_type = "Linux"
    vm_size = "Standard_A4_v2"
  }

  cloud_provider_profile {
    infra_network_profile {
      vnet_subnet_ids = [azurerm_stack_hci_logical_network.test.id]
    }
  }

  control_plane_profile {
    count = 1
    host_ip = "192.168.1.190"
    vm_size = "Standard_A4_v2"
  }

  license_profile {
    azure_hybrid_benefit = "False"
  }

  linux_profile {
    ssh_key = [tls_private_key.rsaKey.public_key_openssh]
  }

  network_profile {
    network_policy = "calico"
    pod_cidr = "10.244.0.0/16"
  }

  storage_profile {
    smb_csi_driver_enabled = false
    nfs_csi_driver_enabled = false
  }
}
`, template, data.RandomInteger, os.Getenv(customLocationIdEnv))
}

func (r ArcKubernetesProvisionedClusterInstanceResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_arc_kubernetes_provisioned_cluster_instance" "import" {
  name       = azurerm_arc_kubernetes_provisioned_cluster_instance.test.name
  cluster_id = azurerm_arc_kubernetes_provisioned_cluster_instance.test.cluster_id
  namespace  = azurerm_arc_kubernetes_provisioned_cluster_instance.test.namespace

  git_repository {
    url             = "https://github.com/Azure/arc-k8s-demo"
    reference_type  = "branch"
    reference_value = "main"
  }

  kustomizations {
    name = "kustomization-1"
  }

  depends_on = [
    azurerm_arc_kubernetes_cluster_extension.test
  ]
}
`, config)
}

func (r ArcKubernetesProvisionedClusterInstanceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "tls_private_key" "rsaKey" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "azurerm_resource_group" "test" {
  name = "acctest-akpci-%[1]d"
  location = "%[2]s"
}

resource "azurerm_stack_hci_logical_network" "test" {
  name                = "acctest-ln-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = "%[3]s"
  virtual_switch_name = "ConvergedSwitch(managementcompute)"
  dns_servers         = ["192.168.1.254"]

  subnet {
    ip_allocation_method = "Static"
    address_prefix       = "192.168.1.0/24"
    ip_pool {
      start = "192.168.1.171"
      end   = "192.168.1.190"
    }
    route {
      name                = "test-route"
      address_prefix      = "10.0.0.0/0"
      next_hop_ip_address = "192.168.1.1"
    }
  }
}

# from https://github.com/Azure/Edge-infrastructure-quickstart-template/blob/main/modules/aks-arc/readiness.ps1
resource "terraform_data" "waitAksVhdReady" {
  provisioner "local-exec" {
    command     = "powershell.exe -ExecutionPolicy Bypass -NoProfile -File ./testdata/readiness.ps1 -customLocationResourceId '%[2]s' -kubernetesVersion '1.28.5' -osSku CBLMariner"
    interpreter = ["PowerShell", "-Command"]
  }
}

resource "azurerm_arc_kubernetes_cluster" "test" {
  name                = "acctest-akcc-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  kind                = "ProvisionedCluster"
  identity {
    type = "SystemAssigned"
  }
  depends_on = [terraform_data.waitAksVhdReady]
}

`, data.RandomInteger, data.Locations.Primary, os.Getenv(customLocationIdEnv))
}
