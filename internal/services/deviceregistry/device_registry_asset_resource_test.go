package deviceregistry_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceregistry/2024-11-01/assets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const (
	ASSET_ARM_CLIENT_ID     = "ARM_CLIENT_ID"
	ASSET_ARM_CLIENT_SECRET = "ARM_CLIENT_SECRET"
)

type AssetTestResource struct{}

func TestAccAsset_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset", "test")
	r := AssetTestResource{}

	if os.Getenv(ASSET_ARM_CLIENT_ID) == "" || os.Getenv(ASSET_ARM_CLIENT_SECRET) == "" {
		t.Skipf("Skipping test due to missing environment variables %s and/or %s", ASSET_ARM_CLIENT_ID, ASSET_ARM_CLIENT_SECRET)
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("asset_endpoint_profile_ref").HasValue("myAssetEndpointProfile"),
				check.That(data.ResourceName).Key("discovered_asset_refs.#").HasValue("3"),
				check.That(data.ResourceName).Key("discovered_asset_refs.0").HasValue("foo"),
				check.That(data.ResourceName).Key("discovered_asset_refs.1").HasValue("bar"),
				check.That(data.ResourceName).Key("discovered_asset_refs.2").HasValue("baz"),
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

	if os.Getenv(ASSET_ARM_CLIENT_ID) == "" || os.Getenv(ASSET_ARM_CLIENT_SECRET) == "" {
		t.Skipf("Skipping test due to missing environment variables %s and/or %s", ASSET_ARM_CLIENT_ID, ASSET_ARM_CLIENT_SECRET)
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("asset_endpoint_profile_ref").HasValue("myAssetEndpointProfile"),
				check.That(data.ResourceName).Key("display_name").HasValue("my asset"),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
				check.That(data.ResourceName).Key("external_asset_id").HasValue("8ZBA6LRHU0A458969"),
				check.That(data.ResourceName).Key("attributes.foo").HasValue("bar"),
				check.That(data.ResourceName).Key("attributes.x").HasValue("y"),
				check.That(data.ResourceName).Key("default_datasets_configuration").HasValue("{\"defaultPublishingInterval\":200,\"defaultQueueSize\":10,\"defaultSamplingInterval\":500}"),
				check.That(data.ResourceName).Key("default_events_configuration").HasValue("{\"defaultPublishingInterval\":200,\"defaultQueueSize\":10,\"defaultSamplingInterval\":500}"),
				check.That(data.ResourceName).Key("default_topic_path").HasValue("/path/defaultTopic"),
				check.That(data.ResourceName).Key("default_topic_retain").HasValue("Keep"),
				check.That(data.ResourceName).Key("description").HasValue("this is my asset"),
				check.That(data.ResourceName).Key("discovered_asset_refs.#").HasValue("3"),
				check.That(data.ResourceName).Key("discovered_asset_refs.0").HasValue("foo"),
				check.That(data.ResourceName).Key("discovered_asset_refs.1").HasValue("bar"),
				check.That(data.ResourceName).Key("discovered_asset_refs.2").HasValue("baz"),
				check.That(data.ResourceName).Key("documentation_uri").HasValue("https://example.com/about"),
				check.That(data.ResourceName).Key("hardware_revision").HasValue("1.0"),
				check.That(data.ResourceName).Key("manufacturer").HasValue("Contoso"),
				check.That(data.ResourceName).Key("manufacturer_uri").HasValue("https://www.contoso.com/manufacturerUri"),
				check.That(data.ResourceName).Key("model").HasValue("ContosoModel"),
				check.That(data.ResourceName).Key("product_code").HasValue("SA34VDG"),
				check.That(data.ResourceName).Key("serial_number").HasValue("64-103816-519918-8"),
				check.That(data.ResourceName).Key("software_revision").HasValue("2.0"),
				check.That(data.ResourceName).Key("tags.site").HasValue("building-1"),
				check.That(data.ResourceName).Key("datasets.#").HasValue("1"),
				check.That(data.ResourceName).Key("datasets.0.dataset_configuration").HasValue("{\"publishingInterval\":7,\"queueSize\":8,\"samplingInterval\":1000}"),
				check.That(data.ResourceName).Key("datasets.0.name").HasValue("dataset1"),
				check.That(data.ResourceName).Key("datasets.0.topic_path").HasValue("/path/dataset1"),
				check.That(data.ResourceName).Key("datasets.0.topic_retain").HasValue("Keep"),
				check.That(data.ResourceName).Key("datasets.0.data_points.#").HasValue("2"),
				check.That(data.ResourceName).Key("datasets.0.data_points.0.data_point_configuration").HasValue("{\"publishingInterval\":7,\"queueSize\":8,\"samplingInterval\":1000}"),
				check.That(data.ResourceName).Key("datasets.0.data_points.0.data_source").HasValue("nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt1"),
				check.That(data.ResourceName).Key("datasets.0.data_points.0.name").HasValue("datapoint1"),
				check.That(data.ResourceName).Key("datasets.0.data_points.0.observability_mode").HasValue("Counter"),
				check.That(data.ResourceName).Key("datasets.0.data_points.1.data_point_configuration").HasValue("{\"publishingInterval\":7,\"queueSize\":8,\"samplingInterval\":1000}"),
				check.That(data.ResourceName).Key("datasets.0.data_points.1.data_source").HasValue("nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt2"),
				check.That(data.ResourceName).Key("datasets.0.data_points.1.name").HasValue("datapoint2"),
				check.That(data.ResourceName).Key("datasets.0.data_points.1.observability_mode").HasValue("None"),
				check.That(data.ResourceName).Key("events.#").HasValue("2"),
				check.That(data.ResourceName).Key("events.0.event_configuration").HasValue("{\"publishingInterval\":7,\"queueSize\":8,\"samplingInterval\":1000}"),
				check.That(data.ResourceName).Key("events.0.event_notifier").HasValue("nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt3"),
				check.That(data.ResourceName).Key("events.0.name").HasValue("event1"),
				check.That(data.ResourceName).Key("events.0.observability_mode").HasValue("Log"),
				check.That(data.ResourceName).Key("events.0.topic_path").HasValue("/path/event1"),
				check.That(data.ResourceName).Key("events.0.topic_retain").HasValue("Never"),
				check.That(data.ResourceName).Key("events.1.event_configuration").HasValue("{\"publishingInterval\":7,\"queueSize\":8,\"samplingInterval\":1000}"),
				check.That(data.ResourceName).Key("events.1.event_notifier").HasValue("nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt4"),
				check.That(data.ResourceName).Key("events.1.name").HasValue("event2"),
				check.That(data.ResourceName).Key("events.1.observability_mode").HasValue("None"),
				check.That(data.ResourceName).Key("events.1.topic_path").HasValue("/path/event2"),
				check.That(data.ResourceName).Key("events.1.topic_retain").HasValue("Keep"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAsset_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset", "test")
	r := AssetTestResource{}

	if os.Getenv(ASSET_ARM_CLIENT_ID) == "" || os.Getenv(ASSET_ARM_CLIENT_SECRET) == "" {
		t.Skipf("Skipping test due to missing environment variables %s and/or %s", ASSET_ARM_CLIENT_ID, ASSET_ARM_CLIENT_SECRET)
	}

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

func TestAccAsset_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset", "test")
	r := AssetTestResource{}

	if os.Getenv(ASSET_ARM_CLIENT_ID) == "" || os.Getenv(ASSET_ARM_CLIENT_SECRET) == "" {
		t.Skipf("Skipping test due to missing environment variables %s and/or %s", ASSET_ARM_CLIENT_ID, ASSET_ARM_CLIENT_SECRET)
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
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
				check.That(data.ResourceName).Key("asset_endpoint_profile_ref").HasValue("myAssetEndpointProfile"),
				check.That(data.ResourceName).Key("display_name").HasValue("my asset"),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
				check.That(data.ResourceName).Key("external_asset_id").HasValue("8ZBA6LRHU0A458969"),
				check.That(data.ResourceName).Key("attributes.foo").HasValue("bar"),
				check.That(data.ResourceName).Key("attributes.x").HasValue("y"),
				check.That(data.ResourceName).Key("default_datasets_configuration").HasValue("{\"defaultPublishingInterval\":200,\"defaultQueueSize\":10,\"defaultSamplingInterval\":500}"),
				check.That(data.ResourceName).Key("default_events_configuration").HasValue("{\"defaultPublishingInterval\":200,\"defaultQueueSize\":10,\"defaultSamplingInterval\":500}"),
				check.That(data.ResourceName).Key("default_topic_path").HasValue("/path/defaultTopic"),
				check.That(data.ResourceName).Key("default_topic_retain").HasValue("Keep"),
				check.That(data.ResourceName).Key("description").HasValue("this is my asset"),
				check.That(data.ResourceName).Key("discovered_asset_refs.#").HasValue("3"),
				check.That(data.ResourceName).Key("discovered_asset_refs.0").HasValue("foo"),
				check.That(data.ResourceName).Key("discovered_asset_refs.1").HasValue("bar"),
				check.That(data.ResourceName).Key("discovered_asset_refs.2").HasValue("baz"),
				check.That(data.ResourceName).Key("documentation_uri").HasValue("https://example.com/about"),
				check.That(data.ResourceName).Key("hardware_revision").HasValue("1.0"),
				check.That(data.ResourceName).Key("manufacturer").HasValue("Contoso"),
				check.That(data.ResourceName).Key("manufacturer_uri").HasValue("https://www.contoso.com/manufacturerUri"),
				check.That(data.ResourceName).Key("model").HasValue("ContosoModel"),
				check.That(data.ResourceName).Key("product_code").HasValue("SA34VDG"),
				check.That(data.ResourceName).Key("serial_number").HasValue("64-103816-519918-8"),
				check.That(data.ResourceName).Key("software_revision").HasValue("2.0"),
				check.That(data.ResourceName).Key("tags.site").HasValue("building-1"),
				check.That(data.ResourceName).Key("datasets.#").HasValue("1"),
				check.That(data.ResourceName).Key("datasets.0.dataset_configuration").HasValue("{\"publishingInterval\":7,\"queueSize\":8,\"samplingInterval\":1000}"),
				check.That(data.ResourceName).Key("datasets.0.name").HasValue("dataset1"),
				check.That(data.ResourceName).Key("datasets.0.topic_path").HasValue("/path/dataset1"),
				check.That(data.ResourceName).Key("datasets.0.topic_retain").HasValue("Keep"),
				check.That(data.ResourceName).Key("datasets.0.data_points.#").HasValue("2"),
				check.That(data.ResourceName).Key("datasets.0.data_points.0.data_point_configuration").HasValue("{\"publishingInterval\":7,\"queueSize\":8,\"samplingInterval\":1000}"),
				check.That(data.ResourceName).Key("datasets.0.data_points.0.data_source").HasValue("nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt1"),
				check.That(data.ResourceName).Key("datasets.0.data_points.0.name").HasValue("datapoint1"),
				check.That(data.ResourceName).Key("datasets.0.data_points.0.observability_mode").HasValue("Counter"),
				check.That(data.ResourceName).Key("datasets.0.data_points.1.data_point_configuration").HasValue("{\"publishingInterval\":7,\"queueSize\":8,\"samplingInterval\":1000}"),
				check.That(data.ResourceName).Key("datasets.0.data_points.1.data_source").HasValue("nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt2"),
				check.That(data.ResourceName).Key("datasets.0.data_points.1.name").HasValue("datapoint2"),
				check.That(data.ResourceName).Key("datasets.0.data_points.1.observability_mode").HasValue("None"),
				check.That(data.ResourceName).Key("events.#").HasValue("2"),
				check.That(data.ResourceName).Key("events.0.event_configuration").HasValue("{\"publishingInterval\":7,\"queueSize\":8,\"samplingInterval\":1000}"),
				check.That(data.ResourceName).Key("events.0.event_notifier").HasValue("nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt3"),
				check.That(data.ResourceName).Key("events.0.name").HasValue("event1"),
				check.That(data.ResourceName).Key("events.0.observability_mode").HasValue("Log"),
				check.That(data.ResourceName).Key("events.0.topic_path").HasValue("/path/event1"),
				check.That(data.ResourceName).Key("events.0.topic_retain").HasValue("Never"),
				check.That(data.ResourceName).Key("events.1.event_configuration").HasValue("{\"publishingInterval\":7,\"queueSize\":8,\"samplingInterval\":1000}"),
				check.That(data.ResourceName).Key("events.1.event_notifier").HasValue("nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt4"),
				check.That(data.ResourceName).Key("events.1.name").HasValue("event2"),
				check.That(data.ResourceName).Key("events.1.observability_mode").HasValue("None"),
				check.That(data.ResourceName).Key("events.1.topic_path").HasValue("/path/event2"),
				check.That(data.ResourceName).Key("events.1.topic_retain").HasValue("Keep"),
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
	resp, err := client.DeviceRegistry.AssetClient.Get(ctx, *id)
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
  name                       = "acctest-asset-%[2]d"
  resource_group_name        = azurerm_resource_group.test.name
  extended_location_name     = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.ExtendedLocation/customLocations/${local.custom_location}"
  extended_location_type     = "CustomLocation"
  asset_endpoint_profile_ref = "myAssetEndpointProfile"
  discovered_asset_refs = [
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
  name                       = "acctest-asset-%[2]d"
  resource_group_name        = azurerm_resource_group.test.name
  extended_location_name     = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.ExtendedLocation/customLocations/${local.custom_location}"
  extended_location_type     = "CustomLocation"
  location                   = "%[3]s"
  asset_endpoint_profile_ref = "myAssetEndpointProfile"
  display_name               = "my asset"
  enabled                    = true
  external_asset_id          = "8ZBA6LRHU0A458969"
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
  default_topic_path   = "/path/defaultTopic"
  default_topic_retain = "Keep"
  description          = "this is my asset"
  discovered_asset_refs = [
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

  datasets {
    dataset_configuration = jsonencode(
      {
        publishingInterval = 7
        queueSize          = 8
        samplingInterval   = 1000
      }
    )
    name         = "dataset1"
    topic_path   = "/path/dataset1"
    topic_retain = "Keep"

    data_points {
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
    data_points {
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

  events {
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
    topic_path         = "/path/event1"
    topic_retain       = "Never"
  }
  events {
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
    topic_path         = "/path/event2"
    topic_retain       = "Keep"
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
  name                       = azurerm_device_registry_asset.test.name
  resource_group_name        = azurerm_device_registry_asset.test.resource_group_name
  extended_location_name     = azurerm_device_registry_asset.test.extended_location_name
  extended_location_type     = azurerm_device_registry_asset.test.extended_location_type
  asset_endpoint_profile_ref = azurerm_device_registry_asset.test.asset_endpoint_profile_ref
  display_name               = azurerm_device_registry_asset.test.display_name
  enabled                    = azurerm_device_registry_asset.test.enabled
  external_asset_id          = azurerm_device_registry_asset.test.external_asset_id
  location                   = azurerm_device_registry_asset.test.location
  depends_on = [
    azurerm_linux_virtual_machine.test
  ]
}
`, template)
}

/*
Creates the terraform template for constants needed for the AIO cluster infra.
*/
func (AssetTestResource) constantsTemplate(data acceptance.TestData) string {
	// Trim the random value (from acceptance.RandTimeInt which is 18 digits) to 10 digits
	// to avoid exceeding the maximum length of the storage account name (24 chars max).
	trimmedRandomInteger := data.RandomInteger % 10000000000
	return fmt.Sprintf(`
locals {
  custom_location           = "acctest-cl%[1]d"
  cluster_name              = "acctest-akcc-%[1]d"
  storage_account           = "acctestsa%[2]d"
  schema_registry           = "acctest-sr-%[1]d"
  schema_registry_namespace = "acctest-rn-%[1]d"
  resource_group_name       = "acctest-rg-%[1]d"
  aio_cluster_resource_name = "acctest-aio%[1]d"
  managed_identity_name     = "acctest-mi%[1]d"
  keyvault_name             = "acctest-kv%[1]d"
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
`, data.RandomInteger, trimmedRandomInteger)
}

/*
The terraform template for all the resources needed to create an AIO cluster on a VM
which the acceptance tests' AssetEndpointProfile resources will be provisioned to.
*/
func (r AssetTestResource) template(data acceptance.TestData) string {
	constantsTemplate := r.constantsTemplate(data)
	credential := r.getCredentials()
	provisionTemplate := r.provisionTemplate(data, credential)

	return fmt.Sprintf(`
%[5]s

resource "azurerm_resource_group" "test" {
  name     = local.resource_group_name
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = "%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  depends_on = [
    azurerm_resource_group.test
  ]
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
  depends_on = [
    azurerm_resource_group.test
  ]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%[1]d"
  location            = "%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  depends_on = [
    azurerm_resource_group.test
  ]
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

  depends_on = [
    azurerm_resource_group.test
  ]
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

  depends_on = [
    azurerm_resource_group.test
  ]
}

resource "azurerm_network_interface_security_group_association" "test" {
  network_interface_id      = azurerm_network_interface.test.id
  network_security_group_id = azurerm_network_security_group.my_terraform_nsg.id
  depends_on = [
    azurerm_resource_group.test
  ]
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
`, data.RandomInteger, data.Locations.Primary, credential, provisionTemplate, constantsTemplate)
}

/*
Copies the scripts and files needed to create and provision the AIO cluster on the VM.
Then ssh's into the VM and executes the cluster setup scripts.
*/
func (r AssetTestResource) provisionTemplate(data acceptance.TestData, credential string) string {
	// Get client secrets from env vars because we need them
	// to remote execute az cli commands on the VM.
	clientId := os.Getenv(ASSET_ARM_CLIENT_ID)
	clientSecret := os.Getenv(ASSET_ARM_CLIENT_SECRET)

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
		cluster_name        = local.cluster_name
		location            = azurerm_resource_group.test.location
		custom_location     = local.custom_location
		storage_account     = local.storage_account
		schema_registry     = local.schema_registry
		schema_registry_namespace = local.schema_registry_namespace
		aio_cluster_resource_name = local.aio_cluster_resource_name
		tenant_id           = data.azurerm_client_config.current.tenant_id
		client_id           = "%[4]s"
		client_secret       = "%[5]s"
		managed_identity_name = local.managed_identity_name
		keyvault_name       = local.keyvault_name
	})
	destination = "%[3]s/setup_aio_cluster.sh"
}

provisioner "remote-exec" {
	inline = [
		"sudo sed -i 's/\r$//' %[3]s/setup_aio_cluster.sh",
		"sudo chmod +x %[3]s/setup_aio_cluster.sh",
		"sudo bash %[3]s/setup_aio_cluster.sh &> %[3]s/agent_log",
	]
}
`, credential, data.RandomInteger, "/home/adminuser", clientId, clientSecret)
}

// Generates a random password for the VM.
func (AssetTestResource) getCredentials() string {
	return fmt.Sprintf("P@$$w0rd%d!", rand.Intn(10000))
}
