package deviceregistry_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceregistry/2024-11-01/assetendpointprofiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/publicipaddresses"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/testclient"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const (
	ASSET_ENDPOINT_PROFILE_ARM_SUBSCRIPTION_ID     = "ARM_SUBSCRIPTION_ID"
	ASSET_ENDPOINT_PROFILE_ARM_CLIENT_ID           = "ARM_CLIENT_ID"
	ASSET_ENDPOINT_PROFILE_ARM_CLIENT_SECRET       = "ARM_CLIENT_SECRET"
	ASSET_ENDPOINT_PROFILE_ARM_ENTRA_APP_OBJECT_ID = "ARM_ENTRA_APP_OBJECT_ID"
)

type AssetEndpointProfileTestResource struct{}

func TestAccAssetEndpointProfile_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset_endpoint_profile", "test")
	r := AssetEndpointProfileTestResource{}

	r.checkEnvironmentVariables(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Apply the template to create the VM and its infra resources.
			// The VM will setup the AIO cluster in the next step.
			Config: r.template(data),
		},
		{
			// Run the setup bash script on the VM to create the AIO cluster.
			// It must be a PreConfig step to ensure AIO cluster is finished setting up
			// before the Asset Endpoint Profile resource is created on the cluster.
			PreConfig: r.setupAIOClusterOnVM(t, data),
			// Then create the Asset Endpoint Profile resource once the AIO cluster is done provisioning.
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("target_address").HasValue("opc.tcp://foo"),
				check.That(data.ResourceName).Key("endpoint_profile_type").HasValue("OpcUa"),
				check.That(data.ResourceName).Key("discovered_asset_endpoint_profile_reference").HasValue("discoveredAssetEndpointProfile123"),
				check.That(data.ResourceName).Key("additional_configuration").HasValue(""),
				check.That(data.ResourceName).Key("authentication.#").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAssetEndpointProfile_complete_certificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset_endpoint_profile", "test")
	r := AssetEndpointProfileTestResource{}

	r.checkEnvironmentVariables(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Apply the template to create the VM and its infra resources.
			// The VM will setup the AIO cluster in the next step.
			Config: r.template(data),
		},
		{
			// Run the setup bash script on the VM to create the AIO cluster.
			// It must be a PreConfig step to ensure AIO cluster is finished setting up
			// before the Asset Endpoint Profile resource is created on the cluster.
			PreConfig: r.setupAIOClusterOnVM(t, data),
			// Then create the Asset Endpoint Profile resource once the AIO cluster is done provisioning.
			Config: r.completeCertificate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("target_address").HasValue("opc.tcp://foo"),
				check.That(data.ResourceName).Key("endpoint_profile_type").HasValue("OpcUa"),
				check.That(data.ResourceName).Key("discovered_asset_endpoint_profile_reference").HasValue("discoveredAssetEndpointProfile123"),
				check.That(data.ResourceName).Key("additional_configuration").HasValue("{\"foo\": \"bar\"}"),
				check.That(data.ResourceName).Key("authentication.0.method").HasValue("Certificate"),
				check.That(data.ResourceName).Key("authentication.0.x509_credential_certificate_secret_name").HasValue("myCertificateRef"),
				check.That(data.ResourceName).Key("authentication.0.username_password_credential_username_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("authentication.0.username_password_credential_password_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("tags.sensor").HasValue("temperature,humidity"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAssetEndpointProfile_complete_usernamePassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset_endpoint_profile", "test")
	r := AssetEndpointProfileTestResource{}

	r.checkEnvironmentVariables(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Apply the template to create the VM and its infra resources.
			// The VM will setup the AIO cluster in the next step.
			Config: r.template(data),
		},
		{
			// Run the setup bash script on the VM to create the AIO cluster.
			// It must be a PreConfig step to ensure AIO cluster is finished setting up
			// before the Asset Endpoint Profile resource is created on the cluster.
			PreConfig: r.setupAIOClusterOnVM(t, data),
			// Then create the Asset Endpoint Profile resource once the AIO cluster is done provisioning.
			Config: r.completeUsernamePassword(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("target_address").HasValue("opc.tcp://foo"),
				check.That(data.ResourceName).Key("endpoint_profile_type").HasValue("OpcUa"),
				check.That(data.ResourceName).Key("discovered_asset_endpoint_profile_reference").HasValue("discoveredAssetEndpointProfile123"),
				check.That(data.ResourceName).Key("additional_configuration").HasValue("{\"foo\": \"bar\"}"),
				check.That(data.ResourceName).Key("authentication.0.method").HasValue("UsernamePassword"),
				check.That(data.ResourceName).Key("authentication.0.x509_credential_certificate_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("authentication.0.username_password_credential_username_secret_name").HasValue("myUsernameRef"),
				check.That(data.ResourceName).Key("authentication.0.username_password_credential_password_secret_name").HasValue("myPasswordRef"),
				check.That(data.ResourceName).Key("tags.sensor").HasValue("temperature,humidity"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAssetEndpointProfile_complete_anonymous(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset_endpoint_profile", "test")
	r := AssetEndpointProfileTestResource{}

	r.checkEnvironmentVariables(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Apply the template to create the VM and its infra resources.
			// The VM will setup the AIO cluster in the next step.
			Config: r.template(data),
		},
		{
			// Run the setup bash script on the VM to create the AIO cluster.
			// It must be a PreConfig step to ensure AIO cluster is finished setting up
			// before the Asset Endpoint Profile resource is created on the cluster.
			PreConfig: r.setupAIOClusterOnVM(t, data),
			// Then create the Asset Endpoint Profile resource once the AIO cluster is done provisioning.
			Config: r.completeAnonymous(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("target_address").HasValue("opc.tcp://foo"),
				check.That(data.ResourceName).Key("endpoint_profile_type").HasValue("OpcUa"),
				check.That(data.ResourceName).Key("discovered_asset_endpoint_profile_reference").HasValue("discoveredAssetEndpointProfile123"),
				check.That(data.ResourceName).Key("additional_configuration").HasValue("{\"foo\": \"bar\"}"),
				check.That(data.ResourceName).Key("authentication.0.method").HasValue("Anonymous"),
				check.That(data.ResourceName).Key("authentication.0.x509_credential_certificate_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("authentication.0.username_password_credential_username_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("authentication.0.username_password_credential_password_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("tags.sensor").HasValue("temperature,humidity"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAssetEndpointProfile_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset_endpoint_profile", "test")
	r := AssetEndpointProfileTestResource{}

	r.checkEnvironmentVariables(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Apply the template to create the VM and its infra resources.
			// The VM will setup the AIO cluster in the next step.
			Config: r.template(data),
		},
		{
			// Run the setup bash script on the VM to create the AIO cluster.
			// It must be a PreConfig step to ensure AIO cluster is finished setting up
			// before the Asset Endpoint Profile resource is created on the cluster.
			PreConfig: r.setupAIOClusterOnVM(t, data),
			// Then create the Asset Endpoint Profile resource once the AIO cluster is done provisioning.
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

	r.checkEnvironmentVariables(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Apply the template to create the VM and its infra resources.
			// The VM will setup the AIO cluster in the next step.
			Config: r.template(data),
		},
		{
			// Run the setup bash script on the VM to create the AIO cluster.
			// It must be a PreConfig step to ensure AIO cluster is finished setting up
			// before the Asset Endpoint Profile resource is created on the cluster.
			PreConfig: r.setupAIOClusterOnVM(t, data),
			// Then create the Asset Endpoint Profile resource once the AIO cluster is done provisioning.
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{ // update the authentication method to certificate
			Config: r.completeCertificate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("target_address").HasValue("opc.tcp://foo"),
				check.That(data.ResourceName).Key("endpoint_profile_type").HasValue("OpcUa"),
				check.That(data.ResourceName).Key("discovered_asset_endpoint_profile_reference").HasValue("discoveredAssetEndpointProfile123"),
				check.That(data.ResourceName).Key("additional_configuration").HasValue("{\"foo\": \"bar\"}"),
				check.That(data.ResourceName).Key("authentication.0.method").HasValue("Certificate"),
				check.That(data.ResourceName).Key("authentication.0.x509_credential_certificate_secret_name").HasValue("myCertificateRef"),
				check.That(data.ResourceName).Key("authentication.0.username_password_credential_username_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("authentication.0.username_password_credential_password_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("tags.sensor").HasValue("temperature,humidity"),
			),
		},
		data.ImportStep(),
		{ // update the authentication method to username/password
			Config: r.completeUsernamePassword(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("target_address").HasValue("opc.tcp://foo"),
				check.That(data.ResourceName).Key("endpoint_profile_type").HasValue("OpcUa"),
				check.That(data.ResourceName).Key("discovered_asset_endpoint_profile_reference").HasValue("discoveredAssetEndpointProfile123"),
				check.That(data.ResourceName).Key("additional_configuration").HasValue("{\"foo\": \"bar\"}"),
				check.That(data.ResourceName).Key("authentication.0.method").HasValue("UsernamePassword"),
				check.That(data.ResourceName).Key("authentication.0.x509_credential_certificate_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("authentication.0.username_password_credential_username_secret_name").HasValue("myUsernameRef"),
				check.That(data.ResourceName).Key("authentication.0.username_password_credential_password_secret_name").HasValue("myPasswordRef"),
				check.That(data.ResourceName).Key("tags.sensor").HasValue("temperature,humidity"),
			),
		},
		data.ImportStep(),
		{ // update the authentication method to anonymous
			Config: r.completeAnonymous(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("target_address").HasValue("opc.tcp://foo"),
				check.That(data.ResourceName).Key("endpoint_profile_type").HasValue("OpcUa"),
				check.That(data.ResourceName).Key("discovered_asset_endpoint_profile_reference").HasValue("discoveredAssetEndpointProfile123"),
				check.That(data.ResourceName).Key("additional_configuration").HasValue("{\"foo\": \"bar\"}"),
				check.That(data.ResourceName).Key("authentication.0.method").HasValue("Anonymous"),
				check.That(data.ResourceName).Key("authentication.0.x509_credential_certificate_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("authentication.0.username_password_credential_username_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("authentication.0.username_password_credential_password_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("tags.sensor").HasValue("temperature,humidity"),
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
	resp, err := client.DeviceRegistry.AssetEndpointProfilesClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func (r AssetEndpointProfileTestResource) basic(data acceptance.TestData) string {
	template := r.template(data)

	return fmt.Sprintf(`
%s

resource "azurerm_device_registry_asset_endpoint_profile" "test" {
  name                                        = "acctest-assetendpointprofile-%[2]d"
  resource_group_id                           = azurerm_resource_group.test.id
  extended_location_id                        = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.ExtendedLocation/customLocations/${local.custom_location}"
  target_address                              = "opc.tcp://foo"
  endpoint_profile_type                       = "OpcUa"
  discovered_asset_endpoint_profile_reference = "discoveredAssetEndpointProfile123"
  location                                    = "%[3]s"
  depends_on = [
    azurerm_linux_virtual_machine.test
  ]
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AssetEndpointProfileTestResource) completeCertificate(data acceptance.TestData) string {
	template := r.template(data)

	return fmt.Sprintf(`
%s

resource "azurerm_device_registry_asset_endpoint_profile" "test" {
  name                                        = "acctest-assetendpointprofile-%[2]d"
  resource_group_id                           = azurerm_resource_group.test.id
  extended_location_id                        = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.ExtendedLocation/customLocations/${local.custom_location}"
  target_address                              = "opc.tcp://foo"
  endpoint_profile_type                       = "OpcUa"
  discovered_asset_endpoint_profile_reference = "discoveredAssetEndpointProfile123"
  additional_configuration                    = "{\"foo\": \"bar\"}"
  authentication {
    method                                  = "Certificate"
    x509_credential_certificate_secret_name = "myCertificateRef"
  }
  tags = {
    "sensor" = "temperature,humidity"
  }
  location = "%[3]s"
  depends_on = [
    azurerm_linux_virtual_machine.test
  ]
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AssetEndpointProfileTestResource) completeUsernamePassword(data acceptance.TestData) string {
	template := r.template(data)

	return fmt.Sprintf(`
%s

resource "azurerm_device_registry_asset_endpoint_profile" "test" {
  name                                        = "acctest-assetendpointprofile-%[2]d"
  resource_group_id                           = azurerm_resource_group.test.id
  extended_location_id                        = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.ExtendedLocation/customLocations/${local.custom_location}"
  target_address                              = "opc.tcp://foo"
  endpoint_profile_type                       = "OpcUa"
  discovered_asset_endpoint_profile_reference = "discoveredAssetEndpointProfile123"
  additional_configuration                    = "{\"foo\": \"bar\"}"
  authentication {
    method                                            = "UsernamePassword"
    username_password_credential_username_secret_name = "myUsernameRef"
    username_password_credential_password_secret_name = "myPasswordRef"
  }
  tags = {
    "sensor" = "temperature,humidity"
  }
  location = "%[3]s"
  depends_on = [
    azurerm_linux_virtual_machine.test
  ]
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AssetEndpointProfileTestResource) completeAnonymous(data acceptance.TestData) string {
	template := r.template(data)

	return fmt.Sprintf(`
%s

resource "azurerm_device_registry_asset_endpoint_profile" "test" {
  name                                        = "acctest-assetendpointprofile-%[2]d"
  resource_group_id                           = azurerm_resource_group.test.id
  extended_location_id                        = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.ExtendedLocation/customLocations/${local.custom_location}"
  target_address                              = "opc.tcp://foo"
  endpoint_profile_type                       = "OpcUa"
  discovered_asset_endpoint_profile_reference = "discoveredAssetEndpointProfile123"
  authentication {
    method = "Anonymous"
  }
  additional_configuration = "{\"foo\": \"bar\"}"
  tags = {
    "sensor" = "temperature,humidity"
  }
  location = "%[3]s"
  depends_on = [
    azurerm_linux_virtual_machine.test
  ]
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AssetEndpointProfileTestResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_device_registry_asset_endpoint_profile" "import" {
  name                                        = azurerm_device_registry_asset_endpoint_profile.test.name
  resource_group_id                           = azurerm_device_registry_asset_endpoint_profile.test.resource_group_id
  extended_location_id                        = azurerm_device_registry_asset_endpoint_profile.test.extended_location_id
  target_address                              = azurerm_device_registry_asset_endpoint_profile.test.target_address
  endpoint_profile_type                       = azurerm_device_registry_asset_endpoint_profile.test.endpoint_profile_type
  discovered_asset_endpoint_profile_reference = "discoveredAssetEndpointProfile123"
  location                                    = azurerm_device_registry_asset_endpoint_profile.test.location
  depends_on = [
    azurerm_linux_virtual_machine.test
  ]
}
`, template)
}

/*
The terraform template for all the resources needed to create an AIO cluster on a VM
which the acceptance tests' AssetEndpointProfile resources will be provisioned to.
*/
func (r AssetEndpointProfileTestResource) template(data acceptance.TestData) string {
	credential := r.getCredentials(data)
	provisionTemplate := r.provisionTemplate(data, credential)

	return fmt.Sprintf(`
locals {
  custom_location = "acctest-cl%[1]d"
}

provider "azurerm" {
  features {
    resource_group {
      // RG will contain AIO resources created from VM. So, we don't want to prevent RG deletion which will clean these up.
      prevent_deletion_if_contains_resources = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = "%[2]s"
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
  location            = "%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%[1]d"
  location            = "%[2]s"
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
  location            = "%[2]s"
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
  location                        = "%[2]s"
  size                            = "Standard_F8s_v2"
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

  identity {
    type = "SystemAssigned"
  }

	%[4]s

  depends_on = [
    azurerm_network_interface_security_group_association.test
  ]
}
`, data.RandomInteger, data.Locations.Primary, credential, provisionTemplate)
}

/*
Copies the script needed to create and provision the AIO cluster on the VM via SSH.
It does NOT run the script. That is done in the setupAIOClusterOnVM function.
*/
func (r AssetEndpointProfileTestResource) provisionTemplate(data acceptance.TestData, credential string) string {
	// Get client secrets from env vars because we need them
	// to remote execute az cli commands on the VM.
	clientId := os.Getenv(ASSET_ENDPOINT_PROFILE_ARM_CLIENT_ID)
	clientSecret := os.Getenv(ASSET_ENDPOINT_PROFILE_ARM_CLIENT_SECRET)
	objectId := os.Getenv(ASSET_ENDPOINT_PROFILE_ARM_ENTRA_APP_OBJECT_ID)

	// Trim the random value (from acceptance.RandTimeInt which is 18 digits) to 10 digits
	// to avoid exceeding the maximum length of the storage account name (24 chars max).
	trimmedRandomInteger := data.RandomInteger % 10000000000
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
		storage_account     = "acctestsa%[3]d"
		schema_registry     = "acctest-sr-%[2]d"
		schema_registry_namespace = "acctest-rn-%[2]d"
		aio_cluster_resource_name = "acctest-aio%[2]d"
		tenant_id           = data.azurerm_client_config.current.tenant_id
		client_id           = "%[5]s"
		client_secret       = "%[6]s"
		object_id					  = "%[7]s"
		managed_identity_name = "acctest-mi%[2]d"
		keyvault_name       = "acctest-kv%[2]d"
	})
	destination = "%[4]s/setup_aio_cluster.sh"
}
`, credential, data.RandomInteger, trimmedRandomInteger, "/home/adminuser", clientId, clientSecret, objectId)
}

/*
This function should be called for the PreConfig step after the template() is made to ensure
that the Asset Endpoint Profile resource create is blocked and does not occur until the AIO
cluster is set up on the VM. This is needed because the Asset Endpoint Profile resource
is an arc-enabled resource and requires the VM, AIO cluster, and custom location to be set up
before it can be created, and all of those resources it's dependent on are created by the setup
script run by the VM's `remote-exec` provisioner. However, `remote-exec` is not waited for by
the subsequent Asset Endpoint Profile resources even if `depends_on` is used, and will attempt
to start creating the resource once the VM is created but the AIO cluster is not set up yet. This way,
the setup script runs synchronously so the tests are forced to wait for it to finish before
creating the Asset Endpoint Profile resource.

This function will grab the VM's public IP address to SSH into the VM and run the setup script
(the file was already provisioned on the VM) to create the AIO cluster.
*/
func (r AssetEndpointProfileTestResource) setupAIOClusterOnVM(t *testing.T, data acceptance.TestData) func() {
	return func() {
		// Set up the test client so we can fetch the public IP address.
		clientManager, err := testclient.Build()
		if err != nil {
			t.Fatalf("failed to build client: %+v", err)
		}

		ctx, cancel := context.WithDeadline(clientManager.StopContext, time.Now().Add(15*time.Minute))
		defer cancel()

		// Public IP Address metadata
		publicIpClient := clientManager.Network.PublicIPAddresses
		subscriptionId := os.Getenv(ASSET_ENDPOINT_PROFILE_ARM_SUBSCRIPTION_ID)
		// Resource group name and public IP address name will be the same as template because we are using the same random integer
		resourceGroupName := fmt.Sprintf("acctest-rg-%d", data.RandomInteger)
		publicIpAddressName := fmt.Sprintf("acctestpip-%d", data.RandomInteger)
		publicIpAddressId := commonids.NewPublicIPAddressID(subscriptionId, resourceGroupName, publicIpAddressName)

		// Get the public IP address
		publicIpAddress, err := publicIpClient.Get(ctx, publicIpAddressId, publicipaddresses.DefaultGetOperationOptions())
		if err != nil {
			t.Fatalf("failed to get public ip address: %+v", err)
		}

		if publicIpAddress.Model == nil || publicIpAddress.Model.Properties.IPAddress == nil {
			t.Fatalf("public ip address not found '%s'", publicIpAddressName)
		}

		// SSH connection details
		ipAddress := *publicIpAddress.Model.Properties.IPAddress
		username := "adminuser"
		password := r.getCredentials(data)

		// SSH client configuration
		sshConfig := &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		// Connect to the VM
		conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", ipAddress), sshConfig)
		if err != nil {
			t.Fatalf("failed to dial ssh: %s", err)
		}
		defer conn.Close()

		// Create a new session
		stripSession, err := conn.NewSession()
		if err != nil {
			t.Fatalf("failed to create session for stripping carriage return: %s", err)
		}
		defer stripSession.Close()
		// Strip carriage return from the setup script
		if err := stripSession.Run("sudo sed -i 's/\r$//' /home/adminuser/setup_aio_cluster.sh"); err != nil {
			t.Fatalf("failed to run command for stripping carriage return. Error: %s", err)
		}

		// Create a new session
		chmodSession, err := conn.NewSession()
		if err != nil {
			t.Fatalf("failed to create session for enabling execution for setup script: %s", err)
		}
		defer chmodSession.Close()
		// Enable execution for the setup script
		if err := chmodSession.Run("sudo chmod +x /home/adminuser/setup_aio_cluster.sh"); err != nil {
			t.Fatalf("failed to run command for enabling execution. Error: %s", err)
		}

		// Create a new session
		runSession, err := conn.NewSession()
		if err != nil {
			t.Fatalf("failed to create session for running setup script: %s", err)
		}
		defer runSession.Close()
		// Run the setup script
		if err := runSession.Run("sudo bash /home/adminuser/setup_aio_cluster.sh &> /home/adminuser/agent_log"); err != nil {
			t.Fatalf("failed to run command for running setup script. Error: %s", err)
		}
	}
}

// Generates a random password for the VM.
func (AssetEndpointProfileTestResource) getCredentials(data acceptance.TestData) string {
	return fmt.Sprintf("P@$$w0rd%d!", data.RandomInteger)
}

// Checks if the required environment variables are set before running the tests.
// If any of the required variables are not set, the test will be skipped.
func (AssetEndpointProfileTestResource) checkEnvironmentVariables(t *testing.T) {
	envVars := []string{
		ASSET_ENDPOINT_PROFILE_ARM_CLIENT_ID,
		ASSET_ENDPOINT_PROFILE_ARM_CLIENT_SECRET,
		ASSET_ENDPOINT_PROFILE_ARM_SUBSCRIPTION_ID,
		ASSET_ENDPOINT_PROFILE_ARM_ENTRA_APP_OBJECT_ID,
	}
	for _, envVar := range envVars {
		if os.Getenv(envVar) == "" {
			envVarsString := strings.Join(envVars, ", ")
			t.Skipf("Skipping test due to environment variable %s not set. Required variables: %s", envVar, envVarsString)
		}
	}
}
