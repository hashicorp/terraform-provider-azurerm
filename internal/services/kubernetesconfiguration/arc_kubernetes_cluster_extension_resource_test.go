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

type ArcKubernetesClusterExtensionResource struct{}

func TestAccArcKubernetesClusterExtension_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_cluster_extension", "test")
	r := ArcKubernetesClusterExtensionResource{}
	credential := fmt.Sprintf("P@$$w0rd%d!", rand.Intn(10000))
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, credential),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccArcKubernetesClusterExtension_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_cluster_extension", "test")
	r := ArcKubernetesClusterExtensionResource{}
	credential := fmt.Sprintf("P@$$w0rd%d!", rand.Intn(10000))
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, credential),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data, credential),
			ExpectError: acceptance.RequiresImportError(data.ResourceType),
		},
	})
}

func TestAccArcKubernetesClusterExtension_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_cluster_extension", "test")
	r := ArcKubernetesClusterExtensionResource{}
	credential := fmt.Sprintf("P@$$w0rd%d!", rand.Intn(10000))
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, credential),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("configuration_protected_settings"),
	})
}

func TestAccArcKubernetesClusterExtension_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_cluster_extension", "test")
	r := ArcKubernetesClusterExtensionResource{}
	credential := fmt.Sprintf("P@$$w0rd%d!", rand.Intn(10000))
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, credential),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("configuration_protected_settings"),
		{
			Config: r.update(data, credential),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("configuration_protected_settings"),
	})
}

func (r ArcKubernetesClusterExtensionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r ArcKubernetesClusterExtensionResource) template(data acceptance.TestData, credential string) string {
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
`, data.RandomInteger, data.Locations.Primary, credential, "/home/adminuser", "create_arc_kubernetes", "delete_arc_kubernetes", os.Getenv("ARM_SUBSCRIPTION_ID"))
}

func (r ArcKubernetesClusterExtensionResource) basic(data acceptance.TestData, credential string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_arc_kubernetes_cluster_extension" "test" {
  name                  = "acctest-kce-%[2]d"
  resource_group_name   = azurerm_resource_group.test.name
  cluster_name          = "acctest-akcc-%[2]d"
  extension_type        = "microsoft.flux"

  identity {
    type = "SystemAssigned"
  }

  depends_on = [
    azurerm_linux_virtual_machine.test
  ]
}
`, r.template(data, credential), data.RandomInteger)
}

func (r ArcKubernetesClusterExtensionResource) requiresImport(data acceptance.TestData, credential string) string {
	config := r.basic(data, credential)
	return fmt.Sprintf(`
			%s

resource "azurerm_arc_kubernetes_cluster_extension" "import" {
  name                  = azurerm_arc_kubernetes_cluster_extension.test.name
  resource_group_name   = azurerm_arc_kubernetes_cluster_extension.test.resource_group_name
  cluster_name          = azurerm_arc_kubernetes_cluster_extension.test.cluster_name
  extension_type        = azurerm_arc_kubernetes_cluster_extension.test.extension_type
  
  identity {
    type = "SystemAssigned"
  }

  depends_on = [
    azurerm_linux_virtual_machine.test
  ]
}
`, config)
}

func (r ArcKubernetesClusterExtensionResource) complete(data acceptance.TestData, credential string) string {
	return fmt.Sprintf(`
			%[1]s

resource "azurerm_arc_kubernetes_cluster_extension" "test" {
  name                  = "acctest-kce-%[2]d"
  resource_group_name   = azurerm_resource_group.test.name
  cluster_name          = "acctest-akcc-%[2]d"
  extension_type        = "microsoft.flux"
  version               = "1.6.3"
  release_namespace     = "flux-system"

  configuration_protected_settings = {
    "omsagent.secret.key" = "secretKeyValue1"
  }

  configuration_settings = {
    "omsagent.env.clusterName" = "clusterName1"
  }

  identity {
    type = "SystemAssigned"
  }

  depends_on = [
    azurerm_linux_virtual_machine.test
  ]
}
`, r.template(data, credential), data.RandomInteger)
}

func (r ArcKubernetesClusterExtensionResource) update(data acceptance.TestData, credential string) string {
	template := r.template(data, credential)
	return fmt.Sprintf(`
			%[1]s

resource "azurerm_arc_kubernetes_cluster_extension" "test" {
  name                  = "acctest-kce-%[2]d"
  resource_group_name   = azurerm_resource_group.test.name
  cluster_name          = "acctest-akcc-%[2]d"
  extension_type        = "microsoft.flux"
  version               = "1.6.3"
  release_namespace     = "flux-system"

  configuration_protected_settings = {
    "omsagent.secret.key" = "secretKeyValue2"
  }

  configuration_settings = {
    "omsagent.env.clusterName" = "clusterName2"
  }

  identity {
    type = "SystemAssigned"
  }

  depends_on = [
    azurerm_linux_virtual_machine.test
  ]
}
`, template, data.RandomInteger)
}
