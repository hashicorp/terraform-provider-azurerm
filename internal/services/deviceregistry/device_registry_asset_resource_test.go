package deviceregistry_test

import (
	"context"
	"fmt"
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

type AssetTestResource struct{}

const (
	ASSET_CUSTOM_LOCATION_NAME = "ARM_DEVICE_REGISTRY_CUSTOM_LOCATION"
	ASSET_RESOURCE_GROUP_NAME  = "ARM_DEVICE_REGISTRY_RESOURCE_GROUP"
)

func TestAccAsset_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset", "test")
	r := AssetTestResource{}

	if os.Getenv(ASSET_CUSTOM_LOCATION_NAME) == "" || os.Getenv(ASSET_RESOURCE_GROUP_NAME) == "" {
		t.Skipf("Skipping test due to missing environment variables %s and/or %s", ASSET_CUSTOM_LOCATION_NAME, ASSET_RESOURCE_GROUP_NAME)
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("asset_endpoint_profile_ref").HasValue("myAssetEndpointProfile"),
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

	if os.Getenv(ASSET_CUSTOM_LOCATION_NAME) == "" || os.Getenv(ASSET_RESOURCE_GROUP_NAME) == "" {
		t.Skipf("Skipping test due to missing environment variables %s and/or %s", ASSET_CUSTOM_LOCATION_NAME, ASSET_RESOURCE_GROUP_NAME)
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

	if os.Getenv(ASSET_CUSTOM_LOCATION_NAME) == "" || os.Getenv(ASSET_RESOURCE_GROUP_NAME) == "" {
		t.Skipf("Skipping test due to missing environment variables %s and/or %s", ASSET_CUSTOM_LOCATION_NAME, ASSET_RESOURCE_GROUP_NAME)
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(func(data acceptance.TestData) string {
			return r.requiresImport(data)
		}),
	})
}

func TestAccAsset_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset", "test")
	r := AssetTestResource{}

	if os.Getenv(ASSET_CUSTOM_LOCATION_NAME) == "" || os.Getenv(ASSET_RESOURCE_GROUP_NAME) == "" {
		t.Skipf("Skipping test due to missing environment variables %s and/or %s", ASSET_CUSTOM_LOCATION_NAME, ASSET_RESOURCE_GROUP_NAME)
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
				check.That(data.ResourceName).Key("discovered_asset_refs.#").HasValue("0"), // discovered_asset_refs is not updatable after creation
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
	resource_group_name        = local.resource_group_name
	extended_location_name		 = local.custom_location_name
	extended_location_type     = "CustomLocation"
	asset_endpoint_profile_ref = "myAssetEndpointProfile"
	display_name               = "my asset"
	enabled                    = false
	external_asset_id          = "8ZBA6LRHU0A458969"
	location                   = "%[3]s"
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AssetTestResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_device_registry_asset" "test" {
	name                           = "acctest-asset-%[2]d"
	resource_group_name            = local.resource_group_name
	extended_location_name         = local.custom_location_name
	extended_location_type         = "CustomLocation"
	location                       = "%[3]s"
	asset_endpoint_profile_ref     = "myAssetEndpointProfile"
	display_name                   = "my asset"
	enabled                        = true
	external_asset_id              = "8ZBA6LRHU0A458969"
	attributes                     = {
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
	default_events_configuration   = jsonencode(
		{
			defaultPublishingInterval = 200
			defaultQueueSize          = 10
			defaultSamplingInterval   = 500
		}
	)
	default_topic_path             = "/path/defaultTopic"
	default_topic_retain           = "Keep"
	description                    = "this is my asset"
	discovered_asset_refs          = [
		"foo",
		"bar",
		"baz",
	]
	documentation_uri              = "https://example.com/about"
	hardware_revision              = "1.0"
	manufacturer                   = "Contoso"
	manufacturer_uri               = "https://www.contoso.com/manufacturerUri"
	model                          = "ContosoModel"
	product_code                   = "SA34VDG"
	serial_number                  = "64-103816-519918-8"
	software_revision              = "2.0"
	tags                           = {
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
		name                  = "dataset1"
		topic_path            = "/path/dataset1"
		topic_retain          = "Keep"

		data_points {
			data_point_configuration = jsonencode(
				{
					publishingInterval = 7
					queueSize          = 8
					samplingInterval   = 1000
				}
			)
			data_source              = "nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt1"
			name                     = "datapoint1"
			observability_mode       = "Counter"
		}
		data_points {
			data_point_configuration = jsonencode(
				{
					publishingInterval = 7
					queueSize          = 8
					samplingInterval   = 1000
				}
			)
			data_source              = "nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt2"
			name                     = "datapoint2"
			observability_mode       = "None"
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
		event_notifier      = "nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt3"
		name                = "event1"
		observability_mode  = "Log"
		topic_path          = "/path/event1"
		topic_retain        = "Never"
	}
	events {
		event_configuration = jsonencode(
			{
				publishingInterval = 7
				queueSize          = 8
				samplingInterval   = 1000
			}
		)
		event_notifier      = "nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt4"
		name                = "event2"
		observability_mode  = "None"
		topic_path          = "/path/event2"
		topic_retain        = "Keep"
	}
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AssetTestResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_device_registry_asset" "import" {
	name 					             = azurerm_device_registry_asset.test.name
	resource_group_name        = azurerm_device_registry_asset.test.resource_group_name
	extended_location_name     = azurerm_device_registry_asset.test.extended_location_name
	extended_location_type     = azurerm_device_registry_asset.test.extended_location_type
	asset_endpoint_profile_ref = azurerm_device_registry_asset.test.asset_endpoint_profile_ref
	display_name							 = azurerm_device_registry_asset.test.display_name
	enabled									   = azurerm_device_registry_asset.test.enabled
	external_asset_id					 = azurerm_device_registry_asset.test.external_asset_id
	location                   = azurerm_device_registry_asset.test.location
}

`, template)
}

/*
Creates the terraform template for AzureRm provider and needed constants
*/
func (AssetTestResource) template(data acceptance.TestData) string {
	customLocation := os.Getenv(ASSET_CUSTOM_LOCATION_NAME)
	resourceGroup := os.Getenv(ASSET_RESOURCE_GROUP_NAME)

	return fmt.Sprintf(`
locals {
	custom_location_name      = "%[1]s"
	resource_group_name       = "%[2]s"
}

provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}
`, customLocation, resourceGroup)
}
