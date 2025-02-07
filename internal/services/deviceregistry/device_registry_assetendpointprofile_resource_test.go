package deviceregistry_test

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
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	// "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceregistry/2024-11-01/assetendpointprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AssetEndpointProfileTestResource struct{}

// func TestAssetEndpointProfileResource(t *testing.T) {
// 	// Setup arc enabled cluster with Azure IoT Operations extension installed.

// 	// Run all the acceptance tests for the AssetEndpointProfile resource on the cluster.
// 	// NOTE: this is a combined test rather than separate split out tests due to
// 	// AssetEndpointProfile resources must be provisioned to the arc-enabled cluster
// 	// and avoid creating the arc-enabled cluster multiple times.
// 	testCases := map[string]map[string]func(t *testing.T){
// 		"Resource": {
// 			"basic":          testAccAssetEndpointProfile_basic,
// 			"requiresImport": testAccAssetEndpointProfile_requiresImport,
// 			"completeCertificate":       testAccAssetEndpointProfile_complete_certificate,
// 			"completeUsernamePassword":  testAccAssetEndpointProfile_complete_usernamePassword,
// 			"completeAnonymous":         testAccAssetEndpointProfile_complete_anonymous,
// 			"update":         testAccAssetEndpointProfile_update,
// 		},
// 	}

// 	for group, m := range testCases {
// 		for name, tc := range m {
// 			t.Run(group, func(t *testing.T) {
// 				t.Run(name, func(t *testing.T) {
// 					tc(t)
// 				})
// 			})
// 		}
// 	}
// }

func TestAccAssetEndpointProfile_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset_endpoint_profile", "test")
	r := AssetEndpointProfileTestResource{}

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

func TestAccAssetEndpointProfile_complete_certificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset_endpoint_profile", "test")
	r := AssetEndpointProfileTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeCertificate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAssetEndpointProfile_complete_usernamePassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset_endpoint_profile", "test")
	r := AssetEndpointProfileTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeUsernamePassword(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAssetEndpointProfile_complete_anonymous(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset_endpoint_profile", "test")
	r := AssetEndpointProfileTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeAnonymous(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAssetEndpointProfile_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset_endpoint_profile", "test")
	r := AssetEndpointProfileTestResource{}

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

func TestAccAssetEndpointProfile_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset_endpoint_profile", "test")
	r := AssetEndpointProfileTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{ // first provision the resource
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{ // then perform the update
			Config: r.completeCertificate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{ // update the authentication method to username/password
			Config: r.completeUsernamePassword(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{ // update the authentication method to anonymous
			Config: r.completeAnonymous(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (AssetEndpointProfileTestResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := assetendpointprofiles.ParseAssetEndpointProfileID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.DeviceRegistry.AssetEndpointProfileClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func (AssetEndpointProfileTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_device_registry_asset_endpoint_profile" "test" {
	name             = "myAssetEndpointProfileBasic"
	resource_group_name = "adr-terraform-test-113553226"
	extended_location_name = "/subscriptions/efb15086-3322-405d-a9d0-c35715a9b722/resourceGroups/adr-terraform-test-113553226/providers/Microsoft.ExtendedLocation/customLocations/location-2h2vr"
	extended_location_type = "CustomLocation"
	target_address = "opc.tcp://foo"
	endpoint_profile_type = "OpcUa"
	location         = "%s"
}
`, data.Locations.Primary)
}

func (AssetEndpointProfileTestResource) completeCertificate(data acceptance.TestData) string {
	return fmt.Sprintf(`

resource "azurerm_device_registry_asset_endpoint_profile" "test" {
	name             = "myAssetEndpointProfileCertificate"
	resource_group_name = "adr-terraform-test-113553226"
	extended_location_name = "/subscriptions/efb15086-3322-405d-a9d0-c35715a9b722/resourceGroups/adr-terraform-test-113553226/providers/Microsoft.ExtendedLocation/customLocations/location-2h2vr"
	extended_location_type = "CustomLocation"
	target_address = "opc.tcp://foo"
	endpoint_profile_type = "OpcUa"
	discovered_asset_endpoint_profile_ref = "discoveredAssetEndpointProfile123"
	additional_configuration = "{\"foo\": \"bar\"}"
	authentication_method = "Certificate"
	x509_credentials_certificate_secret_name = "myCertificateRef"
	location         = "%s"
}
`, data.Locations.Primary)
}

func (AssetEndpointProfileTestResource) completeUsernamePassword(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_device_registry_asset_endpoint_profile" "test" {
	name             = "myAssetEndpointProfileUsername"
	resource_group_name = "adr-terraform-test-113553226"
	extended_location_name = "/subscriptions/efb15086-3322-405d-a9d0-c35715a9b722/resourceGroups/adr-terraform-test-113553226/providers/Microsoft.ExtendedLocation/customLocations/location-2h2vr"
	extended_location_type = "CustomLocation"
	target_address = "opc.tcp://foo"
	endpoint_profile_type = "OpcUa"
	discovered_asset_endpoint_profile_ref = "discoveredAssetEndpointProfile123"
	additional_configuration = "{\"foo\": \"bar\"}"
	authentication_method = "UsernamePassword"
	username_password_credentials_username_secret_name = "myUsernameRef"
	username_password_credentials_password_secret_name = "myPasswordRef"
	location         = "%s"
}
`, data.Locations.Primary)
}

func (AssetEndpointProfileTestResource) completeAnonymous(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_device_registry_asset_endpoint_profile" "test" {
	name             = "myAssetEndpointProfileAnonymous"
	resource_group_name = "adr-terraform-test-113553226"
	extended_location_name = "/subscriptions/efb15086-3322-405d-a9d0-c35715a9b722/resourceGroups/adr-terraform-test-113553226/providers/Microsoft.ExtendedLocation/customLocations/location-2h2vr"
	extended_location_type = "CustomLocation"
	target_address = "opc.tcp://foo"
	endpoint_profile_type = "OpcUa"
	discovered_asset_endpoint_profile_ref = "discoveredAssetEndpointProfile123"
	additional_configuration = "{\"foo\": \"bar\"}"
	authentication_method = "Anonymous"
	location         = "%s"
}
`, data.Locations.Primary)
}

func (r AssetEndpointProfileTestResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_device_registry_asset_endpoint_profile" "import" {
	name 					 = azurerm_device_registry_asset_endpoint_profile.test.name
	resource_group_name = azurerm_device_registry_asset_endpoint_profile.test.resource_group_name
	extended_location_name = azurerm_device_registry_asset_endpoint_profile.test.extended_location_name
	extended_location_type = azurerm_device_registry_asset_endpoint_profile.test.extended_location_type
	target_address = azurerm_device_registry_asset_endpoint_profile.test.target_address
	endpoint_profile_type = azurerm_device_registry_asset_endpoint_profile.test.endpoint_profile_type
	location         = azurerm_device_registry_asset_endpoint_profile.test.location
}

`, template)
}

func (r AssetEndpointProfileTestResource) generateKey() (string, string, error) {
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

func (r AssetEndpointProfileTestResource) getCredentials(t *testing.T) (credential, privateKey, publicKey string) {
	credential = fmt.Sprintf("P@$$w0rd%d!", rand.Intn(10000))
	privateKey, publicKey, err := r.generateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %+v", err)
	}

	return
}

func (r AssetEndpointProfileTestResource) template(data acceptance.TestData, credential string, publicKey string, privateKey string) string {
	provisionTemplate := r.provisionTemplate(data, credential, privateKey)
	return fmt.Sprintf(`
locals {
	custom_location = "adr-acctest-cl"
	storage_account = "acctestsa"
  schema_registry = "acctestsr"
  schema_registry_namespace = "acctestsrn"
}

provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

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
  name                = "myNetworkSG-%[1]d"
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

  lifecycle {
    ignore_changes = [
      security_rule,
    ]
  }
}

resource "azurerm_network_interface_security_group_association" "test" {
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
    sku       = "24_04-lts"
    version   = "latest"
  }

  depends_on = [
    azurerm_network_interface_security_group_association.test
  ]
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
  name              = "extension-%[1]d"
  cluster_id        = azurerm_arc_kubernetes_cluster.test.id
  extension_type    = "microsoft.vmware"
  release_namespace = "vmware-extension"

  configuration_settings = {
    "Microsoft.CustomLocation.ServiceAccount" = "vmware-operator"
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, credential, publicKey, provisionTemplate)
}

func (r AssetEndpointProfileTestResource) provisionTemplate(data acceptance.TestData, credential string, privateKey string) string {
	// Get client secrets from env vars because we need them 
	// to remote execute az cli commands on the VM & k8s cluster.
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	
	return fmt.Sprintf(`
connection {
 type     = "ssh"
 host     = azurerm_public_ip.test.ip_address
 user     = "adminuser"
 password = "%[1]s"
}

provisioner "file" {
 content = templatefile("testdata/setup_aio_cluster.sh.tftpl", {
   subscription_id     = data.azurerm_client_config.current.subscription_id
   resource_group_name = azurerm_resource_group.test.name
   cluster_name        = "acctest-akcc-%[2]d"
   location            = azurerm_resource_group.test.location
	 custom_location     = local.custom_location
	 storage_account     = local.storage_account
	 schema_registry     = local.schema_registry
	 schema_registry_namespace = local.schema_registry_namespace
   tenant_id           = data.azurerm_client_config.current.tenant_id
	 client_id           = "%[5]s"
	 client_secret       = "%[6]s"
   working_dir         = "%[3]s"
 })
 destination = "%[3]s/setup_aio_cluster.sh"
}

provisioner "file" {
 source      = "testdata/setup_aio_cluster.py"
 destination = "%[3]s/setup_aio_cluster.py"
}

provisioner "file" {
 source      = "testdata/kind.yaml"
 destination = "%[3]s/kind.yaml"
}

provisioner "file" {
 content     = <<EOT
%[4]s
EOT
 destination = "%[3]s/private.pem"
}

provisioner "remote-exec" {
 inline = [
   "sudo sed -i 's/\r$//' %[3]s/setup_aio_cluster.sh",
   "sudo chmod +x %[3]s/setup_aio_cluster.sh",
   "bash %[3]s/setup_aio_cluster.sh > %[3]s/agent_log",
 ]
}
`, credential, data.RandomInteger, "/home/adminuser", privateKey, clientId, clientSecret)
}
