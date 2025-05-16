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
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceregistry/2024-11-01/assets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/publicipaddresses"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/testclient"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const (
	ASSET_ARM_CLIENT_ID           = "ARM_CLIENT_ID"
	ASSET_ARM_CLIENT_SECRET       = "ARM_CLIENT_SECRET"
	ASSET_ARM_SUBSCRIPTION_ID     = "ARM_SUBSCRIPTION_ID"
	ASSET_ARM_ENTRA_APP_OBJECT_ID = "ARM_ENTRA_APP_OBJECT_ID"
)

type AssetTestResource struct{}

func TestAccAsset_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset", "test")
	r := AssetTestResource{}

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
			// before the Asset resource is created on the cluster.
			PreConfig: r.setupAIOClusterOnVM(t, data),
			// Then create the Asset resource once the AIO cluster is done provisioning.
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("asset_endpoint_profile_reference").HasValue("myAssetEndpointProfile"),
				check.That(data.ResourceName).Key("discovered_asset_references.#").HasValue("3"),
				check.That(data.ResourceName).Key("discovered_asset_references.0").HasValue("foo"),
				check.That(data.ResourceName).Key("discovered_asset_references.1").HasValue("bar"),
				check.That(data.ResourceName).Key("discovered_asset_references.2").HasValue("baz"),
				check.That(data.ResourceName).Key("display_name").HasValue("my asset"),
				check.That(data.ResourceName).Key("enabled").HasValue("false"),
				check.That(data.ResourceName).Key("external_asset_id").HasValue("8ZBA6LRHU0A458969"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAsset_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset", "test")
	r := AssetTestResource{}

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
			// before the Asset resource is created on the cluster.
			PreConfig: r.setupAIOClusterOnVM(t, data),
			// Then create the Asset resource once the AIO cluster is done provisioning.
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("asset_endpoint_profile_reference").HasValue("myAssetEndpointProfile"),
				check.That(data.ResourceName).Key("display_name").HasValue("my asset"),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
				check.That(data.ResourceName).Key("external_asset_id").HasValue("8ZBA6LRHU0A458969"),
				check.That(data.ResourceName).Key("attributes.foo").HasValue("bar"),
				check.That(data.ResourceName).Key("attributes.x").HasValue("y"),
				check.That(data.ResourceName).Key("default_datasets_configuration").HasValue("{\"defaultPublishingInterval\":200,\"defaultQueueSize\":10,\"defaultSamplingInterval\":500}"),
				check.That(data.ResourceName).Key("default_events_configuration").HasValue("{\"defaultPublishingInterval\":200,\"defaultQueueSize\":10,\"defaultSamplingInterval\":500}"),
				check.That(data.ResourceName).Key("default_topic.0.path").HasValue("/path/defaultTopic"),
				check.That(data.ResourceName).Key("default_topic.0.retain").HasValue("Keep"),
				check.That(data.ResourceName).Key("description").HasValue("this is my asset"),
				check.That(data.ResourceName).Key("discovered_asset_references.#").HasValue("3"),
				check.That(data.ResourceName).Key("discovered_asset_references.0").HasValue("foo"),
				check.That(data.ResourceName).Key("discovered_asset_references.1").HasValue("bar"),
				check.That(data.ResourceName).Key("discovered_asset_references.2").HasValue("baz"),
				check.That(data.ResourceName).Key("documentation_uri").HasValue("https://example.com/about"),
				check.That(data.ResourceName).Key("hardware_revision").HasValue("1.0"),
				check.That(data.ResourceName).Key("manufacturer").HasValue("Contoso"),
				check.That(data.ResourceName).Key("manufacturer_uri").HasValue("https://www.contoso.com/manufacturerUri"),
				check.That(data.ResourceName).Key("model").HasValue("ContosoModel"),
				check.That(data.ResourceName).Key("product_code").HasValue("SA34VDG"),
				check.That(data.ResourceName).Key("serial_number").HasValue("64-103816-519918-8"),
				check.That(data.ResourceName).Key("software_revision").HasValue("2.0"),
				check.That(data.ResourceName).Key("tags.site").HasValue("building-1"),
				check.That(data.ResourceName).Key("dataset.#").HasValue("1"),
				check.That(data.ResourceName).Key("dataset.0.dataset_configuration").HasValue("{\"publishingInterval\":7,\"queueSize\":8,\"samplingInterval\":1000}"),
				check.That(data.ResourceName).Key("dataset.0.name").HasValue("dataset1"),
				check.That(data.ResourceName).Key("dataset.0.topic.0.path").HasValue("/path/dataset1"),
				check.That(data.ResourceName).Key("dataset.0.topic.0.retain").HasValue("Keep"),
				check.That(data.ResourceName).Key("dataset.0.data_point.#").HasValue("2"),
				check.That(data.ResourceName).Key("dataset.0.data_point.0.data_point_configuration").HasValue("{\"publishingInterval\":7,\"queueSize\":8,\"samplingInterval\":1000}"),
				check.That(data.ResourceName).Key("dataset.0.data_point.0.data_source").HasValue("nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt1"),
				check.That(data.ResourceName).Key("dataset.0.data_point.0.name").HasValue("datapoint1"),
				check.That(data.ResourceName).Key("dataset.0.data_point.0.observability_mode").HasValue("Counter"),
				check.That(data.ResourceName).Key("dataset.0.data_point.1.data_point_configuration").HasValue("{\"publishingInterval\":7,\"queueSize\":8,\"samplingInterval\":1000}"),
				check.That(data.ResourceName).Key("dataset.0.data_point.1.data_source").HasValue("nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt2"),
				check.That(data.ResourceName).Key("dataset.0.data_point.1.name").HasValue("datapoint2"),
				check.That(data.ResourceName).Key("dataset.0.data_point.1.observability_mode").HasValue("None"),
				check.That(data.ResourceName).Key("event.#").HasValue("2"),
				check.That(data.ResourceName).Key("event.0.event_configuration").HasValue("{\"publishingInterval\":7,\"queueSize\":8,\"samplingInterval\":1000}"),
				check.That(data.ResourceName).Key("event.0.event_notifier").HasValue("nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt3"),
				check.That(data.ResourceName).Key("event.0.name").HasValue("event1"),
				check.That(data.ResourceName).Key("event.0.observability_mode").HasValue("Log"),
				check.That(data.ResourceName).Key("event.0.topic.0.path").HasValue("/path/event1"),
				check.That(data.ResourceName).Key("event.0.topic.0.retain").HasValue("Never"),
				check.That(data.ResourceName).Key("event.1.event_configuration").HasValue("{\"publishingInterval\":7,\"queueSize\":8,\"samplingInterval\":1000}"),
				check.That(data.ResourceName).Key("event.1.event_notifier").HasValue("nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt4"),
				check.That(data.ResourceName).Key("event.1.name").HasValue("event2"),
				check.That(data.ResourceName).Key("event.1.observability_mode").HasValue("None"),
				check.That(data.ResourceName).Key("event.1.topic.0.path").HasValue("/path/event2"),
				check.That(data.ResourceName).Key("event.1.topic.0.retain").HasValue("Keep"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAsset_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset", "test")
	r := AssetTestResource{}

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
			// before the Asset resource is created on the cluster.
			PreConfig: r.setupAIOClusterOnVM(t, data),
			// Then create the Asset resource once the AIO cluster is done provisioning.
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAsset_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset", "test")
	r := AssetTestResource{}

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
			// before the Asset resource is created on the cluster.
			PreConfig: r.setupAIOClusterOnVM(t, data),
			// Then create the Asset resource once the AIO cluster is done provisioning.
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("asset_endpoint_profile_reference").HasValue("myAssetEndpointProfile"),
				check.That(data.ResourceName).Key("display_name").HasValue("my asset"),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
				check.That(data.ResourceName).Key("external_asset_id").HasValue("8ZBA6LRHU0A458969"),
				check.That(data.ResourceName).Key("attributes.foo").HasValue("bar"),
				check.That(data.ResourceName).Key("attributes.x").HasValue("y"),
				check.That(data.ResourceName).Key("default_datasets_configuration").HasValue("{\"defaultPublishingInterval\":200,\"defaultQueueSize\":10,\"defaultSamplingInterval\":500}"),
				check.That(data.ResourceName).Key("default_events_configuration").HasValue("{\"defaultPublishingInterval\":200,\"defaultQueueSize\":10,\"defaultSamplingInterval\":500}"),
				check.That(data.ResourceName).Key("default_topic.0.path").HasValue("/path/defaultTopic"),
				check.That(data.ResourceName).Key("default_topic.0.retain").HasValue("Keep"),
				check.That(data.ResourceName).Key("description").HasValue("this is my asset"),
				check.That(data.ResourceName).Key("discovered_asset_references.#").HasValue("3"),
				check.That(data.ResourceName).Key("discovered_asset_references.0").HasValue("foo"),
				check.That(data.ResourceName).Key("discovered_asset_references.1").HasValue("bar"),
				check.That(data.ResourceName).Key("discovered_asset_references.2").HasValue("baz"),
				check.That(data.ResourceName).Key("documentation_uri").HasValue("https://example.com/about"),
				check.That(data.ResourceName).Key("hardware_revision").HasValue("1.0"),
				check.That(data.ResourceName).Key("manufacturer").HasValue("Contoso"),
				check.That(data.ResourceName).Key("manufacturer_uri").HasValue("https://www.contoso.com/manufacturerUri"),
				check.That(data.ResourceName).Key("model").HasValue("ContosoModel"),
				check.That(data.ResourceName).Key("product_code").HasValue("SA34VDG"),
				check.That(data.ResourceName).Key("serial_number").HasValue("64-103816-519918-8"),
				check.That(data.ResourceName).Key("software_revision").HasValue("2.0"),
				check.That(data.ResourceName).Key("tags.site").HasValue("building-1"),
				check.That(data.ResourceName).Key("dataset.#").HasValue("1"),
				check.That(data.ResourceName).Key("dataset.0.dataset_configuration").HasValue("{\"publishingInterval\":7,\"queueSize\":8,\"samplingInterval\":1000}"),
				check.That(data.ResourceName).Key("dataset.0.name").HasValue("dataset1"),
				check.That(data.ResourceName).Key("dataset.0.topic.0.path").HasValue("/path/dataset1"),
				check.That(data.ResourceName).Key("dataset.0.topic.0.retain").HasValue("Keep"),
				check.That(data.ResourceName).Key("dataset.0.data_point.#").HasValue("2"),
				check.That(data.ResourceName).Key("dataset.0.data_point.0.data_point_configuration").HasValue("{\"publishingInterval\":7,\"queueSize\":8,\"samplingInterval\":1000}"),
				check.That(data.ResourceName).Key("dataset.0.data_point.0.data_source").HasValue("nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt1"),
				check.That(data.ResourceName).Key("dataset.0.data_point.0.name").HasValue("datapoint1"),
				check.That(data.ResourceName).Key("dataset.0.data_point.0.observability_mode").HasValue("Counter"),
				check.That(data.ResourceName).Key("dataset.0.data_point.1.data_point_configuration").HasValue("{\"publishingInterval\":7,\"queueSize\":8,\"samplingInterval\":1000}"),
				check.That(data.ResourceName).Key("dataset.0.data_point.1.data_source").HasValue("nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt2"),
				check.That(data.ResourceName).Key("dataset.0.data_point.1.name").HasValue("datapoint2"),
				check.That(data.ResourceName).Key("dataset.0.data_point.1.observability_mode").HasValue("None"),
				check.That(data.ResourceName).Key("event.#").HasValue("2"),
				check.That(data.ResourceName).Key("event.0.event_configuration").HasValue("{\"publishingInterval\":7,\"queueSize\":8,\"samplingInterval\":1000}"),
				check.That(data.ResourceName).Key("event.0.event_notifier").HasValue("nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt3"),
				check.That(data.ResourceName).Key("event.0.name").HasValue("event1"),
				check.That(data.ResourceName).Key("event.0.observability_mode").HasValue("Log"),
				check.That(data.ResourceName).Key("event.0.topic.0.path").HasValue("/path/event1"),
				check.That(data.ResourceName).Key("event.0.topic.0.retain").HasValue("Never"),
				check.That(data.ResourceName).Key("event.1.event_configuration").HasValue("{\"publishingInterval\":7,\"queueSize\":8,\"samplingInterval\":1000}"),
				check.That(data.ResourceName).Key("event.1.event_notifier").HasValue("nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt4"),
				check.That(data.ResourceName).Key("event.1.name").HasValue("event2"),
				check.That(data.ResourceName).Key("event.1.observability_mode").HasValue("None"),
				check.That(data.ResourceName).Key("event.1.topic.0.path").HasValue("/path/event2"),
				check.That(data.ResourceName).Key("event.1.topic.0.retain").HasValue("Keep"),
			),
		},
		data.ImportStep(),
	})
}

func (AssetTestResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := assets.ParseAssetID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.DeviceRegistry.AssetsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func (r AssetTestResource) basic(data acceptance.TestData) string {
	template := r.template(data)

	return fmt.Sprintf(`
%s

resource "azurerm_device_registry_asset" "test" {
  name                             = "acctest-asset-%[2]d"
  resource_group_id                = azurerm_resource_group.test.id
  extended_location_id             = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.ExtendedLocation/customLocations/${local.custom_location}"
  asset_endpoint_profile_reference = "myAssetEndpointProfile"
  discovered_asset_references = [
    "foo",
    "bar",
    "baz",
  ]
  display_name      = "my asset"
  enabled           = false
  external_asset_id = "8ZBA6LRHU0A458969"
  location          = "%[3]s"
  depends_on = [
    azurerm_linux_virtual_machine.test
  ]
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AssetTestResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_device_registry_asset" "test" {
  name                             = "acctest-asset-%[2]d"
  resource_group_id                = azurerm_resource_group.test.id
  extended_location_id             = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.ExtendedLocation/customLocations/${local.custom_location}"
  location                         = "%[3]s"
  asset_endpoint_profile_reference = "myAssetEndpointProfile"
  display_name                     = "my asset"
  enabled                          = true
  external_asset_id                = "8ZBA6LRHU0A458969"
  attributes = {
    "foo" = "bar"
    "x"   = "y"
  }
  default_datasets_configuration = jsonencode(
    {
      defaultPublishingInterval = 200
      defaultQueueSize          = 10
      defaultSamplingInterval   = 500
    }
  )
  default_events_configuration = jsonencode(
    {
      defaultPublishingInterval = 200
      defaultQueueSize          = 10
      defaultSamplingInterval   = 500
    }
  )
  default_topic {
    path   = "/path/defaultTopic"
    retain = "Keep"
  }
  description = "this is my asset"
  discovered_asset_references = [
    "foo",
    "bar",
    "baz",
  ]
  documentation_uri = "https://example.com/about"
  hardware_revision = "1.0"
  manufacturer      = "Contoso"
  manufacturer_uri  = "https://www.contoso.com/manufacturerUri"
  model             = "ContosoModel"
  product_code      = "SA34VDG"
  serial_number     = "64-103816-519918-8"
  software_revision = "2.0"
  tags = {
    "site" = "building-1"
  }

  dataset {
    dataset_configuration = jsonencode(
      {
        publishingInterval = 7
        queueSize          = 8
        samplingInterval   = 1000
      }
    )
    name = "dataset1"
    topic {
      path   = "/path/dataset1"
      retain = "Keep"
    }

    data_point {
      data_point_configuration = jsonencode(
        {
          publishingInterval = 7
          queueSize          = 8
          samplingInterval   = 1000
        }
      )
      data_source        = "nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt1"
      name               = "datapoint1"
      observability_mode = "Counter"
    }
    data_point {
      data_point_configuration = jsonencode(
        {
          publishingInterval = 7
          queueSize          = 8
          samplingInterval   = 1000
        }
      )
      data_source        = "nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt2"
      name               = "datapoint2"
      observability_mode = "None"
    }
  }

  event {
    event_configuration = jsonencode(
      {
        publishingInterval = 7
        queueSize          = 8
        samplingInterval   = 1000
      }
    )
    event_notifier     = "nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt3"
    name               = "event1"
    observability_mode = "Log"
    topic {
      path   = "/path/event1"
      retain = "Never"
    }
  }
  event {
    event_configuration = jsonencode(
      {
        publishingInterval = 7
        queueSize          = 8
        samplingInterval   = 1000
      }
    )
    event_notifier     = "nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt4"
    name               = "event2"
    observability_mode = "None"
    topic {
      path   = "/path/event2"
      retain = "Keep"
    }
  }
  depends_on = [
    azurerm_linux_virtual_machine.test
  ]
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AssetTestResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_device_registry_asset" "import" {
  name                             = azurerm_device_registry_asset.test.name
  resource_group_id                = azurerm_device_registry_asset.test.resource_group_id
  extended_location_id             = azurerm_device_registry_asset.test.extended_location_id
  asset_endpoint_profile_reference = azurerm_device_registry_asset.test.asset_endpoint_profile_reference
  display_name                     = azurerm_device_registry_asset.test.display_name
  enabled                          = azurerm_device_registry_asset.test.enabled
  external_asset_id                = azurerm_device_registry_asset.test.external_asset_id
  location                         = azurerm_device_registry_asset.test.location
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
func (r AssetTestResource) template(data acceptance.TestData) string {
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
func (r AssetTestResource) provisionTemplate(data acceptance.TestData, credential string) string {
	// Get client secrets from env vars because we need them
	// to remote execute az cli commands on the VM.
	clientId := os.Getenv(ASSET_ARM_CLIENT_ID)
	clientSecret := os.Getenv(ASSET_ARM_CLIENT_SECRET)
	objectId := os.Getenv(ASSET_ARM_ENTRA_APP_OBJECT_ID)

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
that the Asset resource create is blocked and does not occur until the AIO
cluster is set up on the VM. This is needed because the Asset resource
is an arc-enabled resource and requires the VM, AIO cluster, and custom location to be set up
before it can be created, and all of those resources it's dependent on are created by the setup
script run by the VM's `remote-exec` provisioner. However, `remote-exec` is not waited for by
the subsequent Asset resources even if `depends_on` is used, and will attempt
to start creating the resource once the VM is created but the AIO cluster is not set up yet. This way,
the setup script runs synchronously so the tests are forced to wait for it to finish before
creating the Asset resource.

This function will grab the VM's public IP address to SSH into the VM and run the setup script
(the file was already provisioned on the VM) to create the AIO cluster.
*/
func (r AssetTestResource) setupAIOClusterOnVM(t *testing.T, data acceptance.TestData) func() {
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
		subscriptionId := os.Getenv(ASSET_ARM_SUBSCRIPTION_ID)
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
func (AssetTestResource) getCredentials(data acceptance.TestData) string {
	return fmt.Sprintf("P@$$w0rd%d!", data.RandomInteger)
}

// Checks if the required environment variables are set before running the tests.
// If any of the required variables are not set, the test will be skipped.
func (AssetTestResource) checkEnvironmentVariables(t *testing.T) {
	envVars := []string{
		ASSET_ARM_CLIENT_ID,
		ASSET_ARM_CLIENT_SECRET,
		ASSET_ARM_SUBSCRIPTION_ID,
		ASSET_ARM_ENTRA_APP_OBJECT_ID,
	}
	for _, envVar := range envVars {
		if os.Getenv(envVar) == "" {
			envVarsString := strings.Join(envVars, ", ")
			t.Skipf("Skipping test due to environment variable %s not set. Required variables: %s", envVar, envVarsString)
		}
	}
}
