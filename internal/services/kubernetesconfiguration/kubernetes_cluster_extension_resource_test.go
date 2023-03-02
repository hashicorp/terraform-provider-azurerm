package kubernetesconfiguration_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2022-11-01/extensions"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KubernetesClusterExtensionResource struct{}

func TestAccKubernetesClusterExtension_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_extension", "test")
	r := KubernetesClusterExtensionResource{}
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

func TestAccKubernetesClusterExtension_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_extension", "test")
	r := KubernetesClusterExtensionResource{}
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

func TestAccKubernetesClusterExtension_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_extension", "test")
	r := KubernetesClusterExtensionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("configuration_protected_settings"),
	})
}

func TestAccKubernetesClusterExtension_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_extension", "test")
	r := KubernetesClusterExtensionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("configuration_protected_settings"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("configuration_protected_settings"),
	})
}

func TestAccKubernetesClusterExtension_plan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_extension", "test")
	r := KubernetesClusterExtensionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.plan(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterExtension_arc(t *testing.T) {
	// follow https://learn.microsoft.com/en-us/azure/azure-arc/kubernetes/quickstart-connect-cluster?tabs=azure-cli%2Cazure-cloud to create an arc Kubernetes cluster
	if os.Getenv("ARM_TEST_ARC_K8S_RG") == "" || os.Getenv("ARM_TEST_ARC_K8S_CLUSTER") == "" {
		t.Skip("Skipping as ARM_TEST_ARC_K8S_RG and ARM_TEST_ARC_K8S_CLUSTER are not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_extension", "test")
	r := KubernetesClusterExtensionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.arc(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r KubernetesClusterExtensionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := extensions.ParseExtensionID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.KubernetesConfiguration.ExtensionsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r KubernetesClusterExtensionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestAKC-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[1]d"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r KubernetesClusterExtensionResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_kubernetes_cluster_extension" "test" {
  name                  = "acctest-kce-%d"
  resource_group_name   = azurerm_resource_group.test.name
  cluster_name          = azurerm_kubernetes_cluster.test.name
  cluster_resource_name = "managedClusters"
  extension_type        = "microsoft.flux"
}
`, template, data.RandomInteger)
}

func (r KubernetesClusterExtensionResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_kubernetes_cluster_extension" "import" {
  name                  = azurerm_kubernetes_cluster_extension.test.name
  resource_group_name   = azurerm_kubernetes_cluster_extension.test.resource_group_name
  cluster_name          = azurerm_kubernetes_cluster_extension.test.cluster_name
  cluster_resource_name = azurerm_kubernetes_cluster_extension.test.cluster_resource_name
  extension_type        = azurerm_kubernetes_cluster_extension.test.extension_type
}
`, config)
}

func (r KubernetesClusterExtensionResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_kubernetes_cluster_extension" "test" {
  name                  = "acctest-kce-%d"
  resource_group_name   = azurerm_resource_group.test.name
  cluster_name          = azurerm_kubernetes_cluster.test.name
  cluster_resource_name = "managedClusters"
  extension_type        = "microsoft.flux"
  version               = "1.6.3"
  release_namespace     = "flux-system"

  configuration_protected_settings = {
    "omsagent.secret.key" = "secretKeyValue1"
  }

  configuration_settings = {
    "omsagent.env.clusterName" = "clusterName1"
  }
}
`, template, data.RandomInteger)
}

func (r KubernetesClusterExtensionResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_kubernetes_cluster_extension" "test" {
  name                  = "acctest-kce-%d"
  resource_group_name   = azurerm_resource_group.test.name
  cluster_name          = azurerm_kubernetes_cluster.test.name
  cluster_resource_name = "managedClusters"
  extension_type        = "microsoft.flux"
  version               = "1.6.3"
  release_namespace     = "flux-system"

  configuration_protected_settings = {
    "omsagent.secret.key" = "secretKeyValue2"
  }

  configuration_settings = {
    "omsagent.env.clusterName" = "clusterName2"
  }
}
`, template, data.RandomInteger)
}

func (r KubernetesClusterExtensionResource) plan(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_kubernetes_cluster_extension" "test" {
  name                  = "acctest-kce-%d"
  resource_group_name   = azurerm_resource_group.test.name
  cluster_name          = azurerm_kubernetes_cluster.test.name
  cluster_resource_name = "managedClusters"
  extension_type        = "cognosys.nodejs-on-alpine"

  configuration_settings = {
    "title" = "Title",
  }

  plan {
    name           = "nodejs-18-alpine-container"
    product        = "nodejs18-alpine-container"
    publisher      = "cognosys"
  }
}
`, template, data.RandomInteger)
}

func (r KubernetesClusterExtensionResource) arc(data acceptance.TestData) string {
	credential := fmt.Sprintf("P@$$w0rd%d!", rand.Intn(10000))
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}
resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
resource "azurerm_role_assignment" "test" {
  scope                = azurerm_resource_group.test.id
  role_definition_name = "Owner"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}
resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%[1]d"
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
resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}
resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test.id
  }
}
resource "azurerm_network_security_group" "my_terraform_nsg" {
  name                = "myNetworkSecurityGroup"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  security_rule {
    name                       = "SSH"
    priority                   = 1001
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "22"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}
resource "azurerm_network_interface_security_group_association" "example" {
  network_interface_id      = azurerm_network_interface.test.id
  network_security_group_id = azurerm_network_security_group.my_terraform_nsg.id
}
resource "azurerm_linux_virtual_machine" "test" {
  name                            = "acctestVM-%[1]d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_F2"
  admin_username                  = "adminuser"
  admin_password                  = "%[3]s"
  provision_vm_agent              = false
  allow_extension_operations      = false
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
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }
  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
  tags = {
    ip_address = azurerm_public_ip.test.ip_address
  }
  connection {
    type     = "ssh"
    host     = self.tags.ip_address
    user     = "adminuser"
    password = "%[3]s"
  }
  provisioner "file" {
    content = templatefile("testdata/%[5]s.sh.tftpl", {
      subscription_id     = "%[7]s"
      user_id             = azurerm_user_assigned_identity.test.id
      resource_group_name = azurerm_resource_group.test.name
      cluster_name        = "acctest-akcc-%[1]d"
      working_dir         = "%[4]s"
    })
    destination = "%[4]s/%[5]s.sh"
  }
  provisioner "file" {
    content = templatefile("testdata/%[6]s.sh.tftpl", {
      subscription_id     = "%[7]s"
      user_id             = azurerm_user_assigned_identity.test.id
      resource_group_name = azurerm_resource_group.test.name
      cluster_name        = "acctest-akcc-%[1]d"
      working_dir         = "%[4]s"
    })
    destination = "%[4]s/%[6]s.sh"
  }
  provisioner "file" {
    source      = "testdata/kind.yaml"
    destination = "%[4]s/kind.yaml"
  }
  provisioner "remote-exec" {
    inline = [
      "sudo sed -i 's/\r$//' %[4]s/%[5]s.sh",
      "sudo chmod +x %[4]s/%[5]s.sh",
      "bash %[4]s/%[5]s.sh > %[4]s/%[5]s_log",
    ]
  }
  provisioner "remote-exec" {
    when = destroy
    inline = [
      "sudo sed -i 's/\r$//' %[4]s/%[6]s.sh",
      "sudo chmod +x %[4]s/%[6]s.sh",
      "bash %[4]s/%[6]s.sh > %[4]s/%[6]s_log",
    ]
  }

  depends_on = [
    azurerm_role_assignment.test
  ]

}

resource "azurerm_kubernetes_cluster_extension" "test" {
  name                  = "acctest-kce-%[1]d"
  resource_group_name   = azurerm_resource_group.test.name
  cluster_name          = "acctest-akcc-%[1]d"
  cluster_resource_name = "connectedClusters"
  extension_type        = "microsoft.flux"

  identity {
    type = "SystemAssigned"
  }

  depends_on = [
    azurerm_linux_virtual_machine.test
  ]
}
`, data.RandomInteger, data.Locations.Primary, credential, "/home/adminuser", "create_arc_kubernetes", "delete_arc_kubernetes", os.Getenv("ARM_SUBSCRIPTION_ID"))
}
