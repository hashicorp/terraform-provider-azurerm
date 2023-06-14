package extendedlocation_test

import (
	"context"
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"math/rand"
	"os"
	"testing"
)

type CustomLocationResource struct{}

func (r CustomLocationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := customlocations.ParseCustomLocationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.ExtendedLocation.CustomLocationsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func TestAccExtendedLocationCustomLocations_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_extended_custom_locations", "test")
	r := CustomLocationResource{}
	credential := fmt.Sprintf("P@$$w0rd%d!", rand.Intn(10000))
	privateKey, publicKey, err := CustomLocationResource{}.generateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %+v", err)
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, credential, privateKey, publicKey),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r CustomLocationResource) basic(data acceptance.TestData, credential string, privateKey string, publicKey string) string {
	template := r.template(data, credential, privateKey, publicKey)
	return fmt.Sprintf(`
%s

//resource "azurerm_role_assignment" "admin" {
//  scope                = azurerm_kubernetes_cluster.test.id
//  role_definition_name = "Azure Kubernetes Service RBAC Cluster Admin"
//  principal_id         = "51dfe1e8-70c6-4de5-a08e-e18aff23d815"
//}

resource "azurerm_extended_custom_locations" "test" {
  name = "acctestcustomlocation%d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
  cluster_extension_ids = [
	"${azurerm_arc_kubernetes_cluster_extension.test.id}"
  ]
  display_name = "customlocation%[2]d"
  namespace = "namespace%[2]d"
  host_resource_id = azurerm_arc_kubernetes_cluster.test.id
}
`, template, data.RandomInteger)
}

func (r CustomLocationResource) template(data acceptance.TestData, credential string, privateKey string, publicKey string) string {
	data.Locations.Primary = "eastus"
	return fmt.Sprintf(`
provider "azurerm" {
  features {
	resource_group {
	  prevent_deletion_if_contains_resources = false
	}
  }
}

resource "azurerm_resource_group" "test" {
  name = "acctestRG-%[1]d"
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
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }
}

data "azurerm_extended_locations" "test" {
  location = azurerm_resource_group.test.location
}

resource "azurerm_arc_kubernetes_cluster" "test" {
  name                         = "acctest-akcc-%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  agent_public_key_certificate = "%[4]s"
  identity {
    type = "SystemAssigned"
  }

  %[5]s

  depends_on = [
    azurerm_linux_virtual_machine.test
  ]
}

resource "azurerm_arc_kubernetes_cluster_extension" "test" {
  name				= "foo"
  cluster_id 	  = azurerm_arc_kubernetes_cluster.test.id
  extension_type 	= "microsoft.flux"

  identity {
	type = "SystemAssigned"
  }

  depends_on = [
  	azurerm_linux_virtual_machine.test
  ]
}
`, data.RandomInteger, data.Locations.Primary, credential, publicKey, r.provisionTemplate(data, credential, privateKey))
}

func (r CustomLocationResource) generateKey() (string, string, error) {
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

func (r CustomLocationResource) provisionTemplate(data acceptance.TestData, credential string, privateKey string) string {
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
