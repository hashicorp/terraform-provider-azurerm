// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package arckubernetes_test

import (
	"context"
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	arckubernetes "github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2024-01-01/connectedclusters"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ArcKubernetesClusterResource struct{}

func TestAccArcKubernetesCluster_basic(t *testing.T) {
	credential, privateKey, publicKey := ArcKubernetesClusterResource{}.getCredentials(t)
	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_cluster", "test")
	r := ArcKubernetesClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, credential, privateKey, publicKey),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data, credential, privateKey, publicKey),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("distribution").HasValue("generic"),
				check.That(data.ResourceName).Key("infrastructure").HasValue("generic"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccArcKubernetesCluster_requiresImport(t *testing.T) {
	credential, privateKey, publicKey := ArcKubernetesClusterResource{}.getCredentials(t)
	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_cluster", "test")
	r := ArcKubernetesClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, credential, privateKey, publicKey),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data, credential, privateKey, publicKey),
			ExpectError: acceptance.RequiresImportError(data.ResourceType),
		},
	})
}

func TestAccArcKubernetesCluster_complete(t *testing.T) {
	credential, privateKey, publicKey := ArcKubernetesClusterResource{}.getCredentials(t)
	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_cluster", "test")
	r := ArcKubernetesClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, credential, privateKey, publicKey),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data, credential, privateKey, publicKey),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("distribution").HasValue("generic"),
				check.That(data.ResourceName).Key("infrastructure").HasValue("generic"),
			),
		},
		data.ImportStep(),
	})
}

func (r ArcKubernetesClusterResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := arckubernetes.ParseConnectedClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.ArcKubernetes.ArcKubernetesClient
	resp, err := client.ConnectedClusterGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r ArcKubernetesClusterResource) template(data acceptance.TestData, credential string) string {
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
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary, credential)
}

func (r ArcKubernetesClusterResource) provisionTemplate(data acceptance.TestData, credential string, privateKey string) string {
	return fmt.Sprintf(`
connection {
  type     = "ssh"
  host     = azurerm_public_ip.test.ip_address
  user     = "adminuser"
  password = "%[1]s"
}

provisioner "file" {
  content = templatefile("testdata/install_agent.sh.tftpl", {
    subscription_id     = "%[4]s"
    resource_group_name = azurerm_resource_group.test.name
    cluster_name        = "acctest-akcc-%[2]d"
    location            = azurerm_resource_group.test.location
    tenant_id           = "%[5]s"
    working_dir         = "%[3]s"
  })
  destination = "%[3]s/install_agent.sh"
}

provisioner "file" {
  source      = "testdata/install_agent.py"
  destination = "%[3]s/install_agent.py"
}

provisioner "file" {
  source      = "testdata/kind.yaml"
  destination = "%[3]s/kind.yaml"
}

provisioner "file" {
  content     = <<EOT
%[6]s
EOT
  destination = "%[3]s/private.pem"
}

provisioner "remote-exec" {
  inline = [
    "sudo sed -i 's/\r$//' %[3]s/install_agent.sh",
    "sudo chmod +x %[3]s/install_agent.sh",
    "bash %[3]s/install_agent.sh > %[3]s/agent_log",
  ]
}
`, credential, data.RandomInteger, "/home/adminuser", os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_TENANT_ID"), privateKey)
}

func (r ArcKubernetesClusterResource) basic(data acceptance.TestData, credential string, privateKey string, publicKey string) string {
	template := r.template(data, credential)
	provisionTemplate := r.provisionTemplate(data, credential, privateKey)
	return fmt.Sprintf(`
				%[1]s

resource "azurerm_arc_kubernetes_cluster" "test" {
  name                         = "acctest-akcc-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  agent_public_key_certificate = "%[4]s"
  identity {
    type = "SystemAssigned"
  }

  %[3]s

  depends_on = [
    azurerm_linux_virtual_machine.test
  ]
}
`, template, data.RandomInteger, provisionTemplate, publicKey)
}

func (r ArcKubernetesClusterResource) requiresImport(data acceptance.TestData, credential string, privateKey string, publicKey string) string {
	config := r.basic(data, credential, privateKey, publicKey)
	return fmt.Sprintf(`
			%s

resource "azurerm_arc_kubernetes_cluster" "import" {
  name                         = azurerm_arc_kubernetes_cluster.test.name
  resource_group_name          = azurerm_arc_kubernetes_cluster.test.resource_group_name
  location                     = azurerm_arc_kubernetes_cluster.test.location
  agent_public_key_certificate = azurerm_arc_kubernetes_cluster.test.agent_public_key_certificate

  identity {
    type = "SystemAssigned"
  }
}
`, config)
}

func (r ArcKubernetesClusterResource) complete(data acceptance.TestData, credential string, privateKey string, publicKey string) string {
	template := r.template(data, credential)
	provisionTemplate := r.provisionTemplate(data, credential, privateKey)
	return fmt.Sprintf(`
			%[1]s

resource "azurerm_arc_kubernetes_cluster" "test" {
  name                         = "acctest-akcc-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  agent_public_key_certificate = "%[4]s"

  identity {
    type = "SystemAssigned"
  }

  tags = {
    ENV = "Test"
  }

  %[3]s

  depends_on = [
    azurerm_linux_virtual_machine.test
  ]
}
`, template, data.RandomInteger, provisionTemplate, publicKey)
}

func (r ArcKubernetesClusterResource) update(data acceptance.TestData, credential string, privateKey string, publicKey string) string {
	template := r.template(data, credential)
	provisionTemplate := r.provisionTemplate(data, credential, privateKey)
	return fmt.Sprintf(`
			%[1]s

resource "azurerm_arc_kubernetes_cluster" "test" {
  name                         = "acctest-akcc-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  agent_public_key_certificate = "%[4]s"

  identity {
    type = "SystemAssigned"
  }

  tags = {
    ENV = "TestUpdate"
  }

  %[3]s

  depends_on = [
    azurerm_linux_virtual_machine.test
  ]
}
`, template, data.RandomInteger, provisionTemplate, publicKey)
}

func (r ArcKubernetesClusterResource) generateKey() (string, string, error) {
	privateKey, err := rsa.GenerateKey(cryptoRand.Reader, 4096)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate RSA key")
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	privatePem := pem.EncodeToMemory(privateKeyBlock)
	if privatePem == nil {
		return "", "", fmt.Errorf("failed to encode pem")
	}

	return string(privatePem), base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)), nil
}

func (r ArcKubernetesClusterResource) getCredentials(t *testing.T) (credential, privateKey, publicKey string) {
	// generateKey() is a time-consuming operation, we only run this test if an env var is set.
	if os.Getenv(resource.EnvTfAcc) == "" {
		t.Skipf("Acceptance tests skipped unless env '%s' set", resource.EnvTfAcc)
	}

	credential = fmt.Sprintf("P@$$w0rd%d!", rand.Intn(10000))
	privateKey, publicKey, err := r.generateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %+v", err)
	}

	return
}
