// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iothub_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type IotHubResource struct{}

func TestAccIotHub_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")
	r := IotHubResource{}

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

func TestAccIotHub_networkRulesSet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")
	r := IotHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.networkRuleSet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.networkRuleSetUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHub_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")
	r := IotHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_iothub"),
		},
	})
}

func TestAccIotHub_standard(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")
	r := IotHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standard(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHub_customRoutes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")
	r := IotHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customRoutes(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("endpoint.#").HasValue("2"),
				check.That(data.ResourceName).Key("endpoint.0.type").HasValue("AzureIotHub.StorageContainer"),
				check.That(data.ResourceName).Key("endpoint.0.batch_frequency_in_seconds").HasValue("300"),
				check.That(data.ResourceName).Key("endpoint.0.max_chunk_size_in_bytes").HasValue("314572800"),
				check.That(data.ResourceName).Key("endpoint.0.encoding").HasValue("Avro"),
				check.That(data.ResourceName).Key("endpoint.0.file_name_format").HasValue("{iothub}/{partition}/{YYYY}/{MM}/{DD}/{HH}/{mm}"),
				check.That(data.ResourceName).Key("endpoint.1.type").HasValue("AzureIotHub.EventHub"),
				check.That(data.ResourceName).Key("route.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHub_enrichments(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")
	r := IotHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.enrichments(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("endpoint.#").HasValue("2"),
				check.That(data.ResourceName).Key("endpoint.0.type").HasValue("AzureIotHub.StorageContainer"),
				check.That(data.ResourceName).Key("endpoint.1.type").HasValue("AzureIotHub.EventHub"),
				check.That(data.ResourceName).Key("route.#").HasValue("1"),
				check.That(data.ResourceName).Key("enrichment.#").HasValue("2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.unsetEnrichments(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("endpoint.#").HasValue("2"),
				check.That(data.ResourceName).Key("endpoint.0.type").HasValue("AzureIotHub.StorageContainer"),
				check.That(data.ResourceName).Key("endpoint.1.type").HasValue("AzureIotHub.EventHub"),
				check.That(data.ResourceName).Key("route.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHub_removeEndpointsAndRoutes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")
	r := IotHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customRoutes(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.removeEndpointsAndRoutes(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHub_fileUpload(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")
	r := IotHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.fileUploadBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("file_upload.0.default_ttl").HasValue("PT1H"),
				check.That(data.ResourceName).Key("file_upload.0.lock_duration").HasValue("PT1M"),
				check.That(data.ResourceName).Key("file_upload.0.sas_ttl").HasValue("PT1H"),
			),
		},
		data.ImportStep(),
		{
			Config: r.fileUploadUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("file_upload.#").HasValue("1"),
				check.That(data.ResourceName).Key("file_upload.0.lock_duration").HasValue("PT5M"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHub_fileUploadAuthenticationTypeUserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")
	r := IotHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.fileUploadAuthenticationTypeUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHub_fileUploadAuthenticationTypeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")
	r := IotHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.fileUploadAuthenticationTypeDefault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.fileUploadAuthenticationTypeUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.fileUploadAuthenticationTypeSystemAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.fileUploadAuthenticationTypeDefault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHub_withDifferentEndpointResourceGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")
	r := IotHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withDifferentEndpointResourceGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHub_fallbackRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")
	r := IotHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.fallbackRoute(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fallback_route.0.source").HasValue("DeviceMessages"),
				check.That(data.ResourceName).Key("fallback_route.0.endpoint_names.#").HasValue("1"),
				check.That(data.ResourceName).Key("fallback_route.0.enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHub_publicAccess(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")
	r := IotHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.publicAccessEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.publicAccessDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHub_minTLSVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")
	r := IotHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.minTLSVersion(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHub_LocalAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")
	r := IotHubResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.disableLocalAuth(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("local_authentication_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.enableLocalAuth(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("local_authentication_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHub_cloudToDevice(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")
	r := IotHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cloudToDevice(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cloud_to_device.0.max_delivery_count").HasValue("20"),
				check.That(data.ResourceName).Key("cloud_to_device.0.default_ttl").HasValue("PT1H30M"),
				check.That(data.ResourceName).Key("cloud_to_device.0.feedback.0.time_to_live").HasValue("PT1H15M"),
				check.That(data.ResourceName).Key("cloud_to_device.0.feedback.0.max_delivery_count").HasValue("25"),
				check.That(data.ResourceName).Key("cloud_to_device.0.feedback.0.lock_duration").HasValue("PT55S"),
			),
		},
		{
			Config: r.cloudToDeviceUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cloud_to_device.0.max_delivery_count").HasValue("30"),
				check.That(data.ResourceName).Key("cloud_to_device.0.default_ttl").HasValue("PT1H"),
				check.That(data.ResourceName).Key("cloud_to_device.0.feedback.0.time_to_live").HasValue("PT1H10M"),
				check.That(data.ResourceName).Key("cloud_to_device.0.feedback.0.max_delivery_count").HasValue("15"),
				check.That(data.ResourceName).Key("cloud_to_device.0.feedback.0.lock_duration").HasValue("PT30S"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHub_identitySystemAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")
	r := IotHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identitySystemAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHub_identitySystemAssignedUserAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")
	r := IotHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identitySystemAssignedUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHub_identityUserAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")
	r := IotHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identityUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.identityUserAssignedUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHub_identityUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")
	r := IotHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.identitySystemAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.identityUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.identitySystemAssignedUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHub_endpointAuthenticationTypeUserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")
	r := IotHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.endpointAuthenticationTypeUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHub_endpointAuthenticationTypeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")
	r := IotHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.endpointAuthenticationTypeDefault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.endpointAuthenticationTypeUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.endpointAuthenticationTypeSystemAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.endpointAuthenticationTypeDefault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHub_cosmosDBRouteUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")
	r := IotHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.updateWithCosmosDBRoute(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateWithCosmosDBRoute(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t IotHubResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.IotHubID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.IoTHub.ResourceClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.ID != nil), nil
}

func (IotHubResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "B1"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r IotHubResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub" "import" {
  name                = azurerm_iothub.test.name
  resource_group_name = azurerm_iothub.test.resource_group_name
  location            = azurerm_iothub.test.location

  sku {
    name     = "B1"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}
`, r.basic(data))
}

func (IotHubResource) standard(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (IotHubResource) networkRuleSet(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  network_rule_set {
    default_action                     = "Allow"
    apply_to_builtin_eventhub_endpoint = true

    ip_rule {
      name    = "test"
      ip_mask = "10.0.1.0/31"
      action  = "Allow"
    }

    ip_rule {
      name    = "test2"
      ip_mask = "10.0.3.0/31"
      action  = "Allow"
    }

  }

  tags = {
    purpose = "testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (IotHubResource) networkRuleSetUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  network_rule_set {
    default_action                     = "Allow"
    apply_to_builtin_eventhub_endpoint = true

    ip_rule {
      name    = "test"
      ip_mask = "10.0.0.0/31"
      action  = "Allow"
    }

    ip_rule {
      name    = "test2"
      ip_mask = "10.0.2.0/31"
      action  = "Allow"
    }

  }

  tags = {
    purpose = "testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (IotHubResource) customRoutes(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_eventhub_namespace" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "acctest-%d"
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctest"
  resource_group_name = azurerm_resource_group.test.name
  namespace_name      = azurerm_eventhub_namespace.test.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_eventhub_authorization_rule" "test" {
  resource_group_name = azurerm_resource_group.test.name
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  name                = "acctest"
  send                = true
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  event_hub_retention_in_days = 7
  event_hub_partition_count   = 77

  endpoint {
    type                = "AzureIotHub.StorageContainer"
    connection_string   = azurerm_storage_account.test.primary_blob_connection_string
    name                = "export"
    container_name      = azurerm_storage_container.test.name
    resource_group_name = azurerm_resource_group.test.name
  }

  endpoint {
    type                = "AzureIotHub.EventHub"
    connection_string   = azurerm_eventhub_authorization_rule.test.primary_connection_string
    name                = "export2"
    resource_group_name = azurerm_resource_group.test.name
  }

  route {
    name           = "export"
    source         = "DeviceMessages"
    condition      = "true"
    endpoint_names = ["export"]
    enabled        = true
  }

  route {
    name           = "export2"
    source         = "DeviceMessages"
    condition      = "true"
    endpoint_names = ["export2"]
    enabled        = true
  }

  tags = {
    purpose = "testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (IotHubResource) enrichments(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_eventhub_namespace" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "acctest-%d"
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctest"
  resource_group_name = azurerm_resource_group.test.name
  namespace_name      = azurerm_eventhub_namespace.test.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_eventhub_authorization_rule" "test" {
  resource_group_name = azurerm_resource_group.test.name
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  name                = "acctest"
  send                = true
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  event_hub_retention_in_days = 7
  event_hub_partition_count   = 77

  endpoint {
    type                       = "AzureIotHub.StorageContainer"
    connection_string          = azurerm_storage_account.test.primary_blob_connection_string
    name                       = "export"
    batch_frequency_in_seconds = 60
    max_chunk_size_in_bytes    = 10485760
    container_name             = azurerm_storage_container.test.name
    encoding                   = "Avro"
    file_name_format           = "{iothub}/{partition}_{YYYY}_{MM}_{DD}_{HH}_{mm}"
    resource_group_name        = azurerm_resource_group.test.name
  }

  endpoint {
    type                = "AzureIotHub.EventHub"
    connection_string   = azurerm_eventhub_authorization_rule.test.primary_connection_string
    name                = "export2"
    resource_group_name = azurerm_resource_group.test.name
  }

  route {
    name           = "export"
    source         = "DeviceMessages"
    condition      = "true"
    endpoint_names = ["export"]
    enabled        = true
  }

  enrichment {
    key            = "enrichment"
    value          = "$twin.tags.Tenant"
    endpoint_names = ["export2"]
  }

  enrichment {
    key            = "enrichment2"
    value          = "Multiple endpoint"
    endpoint_names = ["export", "export2"]
  }

  tags = {
    purpose = "testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (IotHubResource) unsetEnrichments(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_eventhub_namespace" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "acctest-%d"
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctest"
  resource_group_name = azurerm_resource_group.test.name
  namespace_name      = azurerm_eventhub_namespace.test.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_eventhub_authorization_rule" "test" {
  resource_group_name = azurerm_resource_group.test.name
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  name                = "acctest"
  send                = true
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  event_hub_retention_in_days = 7
  event_hub_partition_count   = 77

  endpoint {
    type                       = "AzureIotHub.StorageContainer"
    connection_string          = azurerm_storage_account.test.primary_blob_connection_string
    name                       = "export"
    batch_frequency_in_seconds = 60
    max_chunk_size_in_bytes    = 10485760
    container_name             = azurerm_storage_container.test.name
    encoding                   = "Avro"
    file_name_format           = "{iothub}/{partition}_{YYYY}_{MM}_{DD}_{HH}_{mm}"
    resource_group_name        = azurerm_resource_group.test.name
  }

  endpoint {
    type                = "AzureIotHub.EventHub"
    connection_string   = azurerm_eventhub_authorization_rule.test.primary_connection_string
    name                = "export2"
    resource_group_name = azurerm_resource_group.test.name
  }

  route {
    name           = "export"
    source         = "DeviceMessages"
    condition      = "true"
    endpoint_names = ["export"]
    enabled        = true
  }

  enrichment = []

  tags = {
    purpose = "testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (IotHubResource) removeEndpointsAndRoutes(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_eventhub_namespace" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "acctest-%d"
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctest"
  resource_group_name = azurerm_resource_group.test.name
  namespace_name      = azurerm_eventhub_namespace.test.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_eventhub_authorization_rule" "test" {
  resource_group_name = azurerm_resource_group.test.name
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  name                = "acctest"
  send                = true
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  event_hub_retention_in_days = 7
  event_hub_partition_count   = 77

  endpoint = []

  route = []

  tags = {
    purpose = "testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (IotHubResource) fallbackRoute(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  fallback_route {
    source         = "DeviceMessages"
    endpoint_names = ["events"]
    enabled        = true
  }

  tags = {
    purpose = "testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (IotHubResource) fileUploadBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  file_upload {
    connection_string = azurerm_storage_account.test.primary_blob_connection_string
    container_name    = azurerm_storage_container.test.name
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (IotHubResource) fileUploadUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  file_upload {
    connection_string  = azurerm_storage_account.test.primary_blob_connection_string
    container_name     = azurerm_storage_container.test.name
    notifications      = true
    max_delivery_count = 12
    sas_ttl            = "PT2H"
    default_ttl        = "PT3H"
    lock_duration      = "PT5M"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (r IotHubResource) fileUploadAuthenticationTypeDefault(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  file_upload {
    connection_string = azurerm_storage_account.test.primary_blob_connection_string
    container_name    = azurerm_storage_container.test.name
  }

  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }

  depends_on = [
    azurerm_role_assignment.test_storage_blob_data_contrib_user,
  ]
}
`, r.fileUploadAuthenticationTypeTemplate(data), data.RandomInteger)
}

func (r IotHubResource) fileUploadAuthenticationTypeSystemAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  file_upload {
    connection_string = azurerm_storage_account.test.primary_blob_connection_string
    container_name    = azurerm_storage_container.test.name

    authentication_type = "identityBased"
  }

  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }

  depends_on = [
    azurerm_role_assignment.test_storage_blob_data_contrib_user,
  ]
}
`, r.fileUploadAuthenticationTypeTemplate(data), data.RandomInteger)
}

func (r IotHubResource) fileUploadAuthenticationTypeUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  file_upload {
    connection_string = azurerm_storage_account.test.primary_blob_connection_string
    container_name    = azurerm_storage_container.test.name

    authentication_type = "identityBased"
    identity_id         = azurerm_user_assigned_identity.test.id
  }

  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }

  depends_on = [
    azurerm_role_assignment.test_storage_blob_data_contrib_user,
  ]
}
`, r.fileUploadAuthenticationTypeTemplate(data), data.RandomInteger)
}

func (IotHubResource) fileUploadAuthenticationTypeTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuai-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_role_assignment" "test_storage_blob_data_contrib_user" {
  role_definition_name = "Storage Blob Data Contributor"
  scope                = azurerm_storage_account.test.id
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_role_assignment" "test_storage_blob_data_contrib_system" {
  role_definition_name = "Storage Blob Data Contributor"
  scope                = azurerm_storage_account.test.id
  principal_id         = azurerm_iothub.test.identity[0].principal_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (IotHubResource) publicAccessEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "B1"
    capacity = "1"
  }

  public_network_access_enabled = true

  tags = {
    purpose = "testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (IotHubResource) publicAccessDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "B1"
    capacity = "1"
  }

  public_network_access_enabled = false

  tags = {
    purpose = "testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (IotHubResource) withDifferentEndpointResourceGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%d-1"
  location = "%s"
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG-iothub-%d-2"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_eventhub_namespace" "test" {
  resource_group_name = azurerm_resource_group.test2.name
  location            = azurerm_resource_group.test2.location
  name                = "acctest-%d"
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctest"
  resource_group_name = azurerm_resource_group.test2.name
  namespace_name      = azurerm_eventhub_namespace.test.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_eventhub_authorization_rule" "test" {
  resource_group_name = azurerm_resource_group.test2.name
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  name                = "acctest"
  send                = true
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  event_hub_retention_in_days = 7
  event_hub_partition_count   = 77

  endpoint {
    type                       = "AzureIotHub.StorageContainer"
    connection_string          = azurerm_storage_account.test.primary_blob_connection_string
    name                       = "export"
    batch_frequency_in_seconds = 60
    max_chunk_size_in_bytes    = 10485760
    container_name             = azurerm_storage_container.test.name
    encoding                   = "Avro"
    file_name_format           = "{iothub}/{partition}_{YYYY}_{MM}_{DD}_{HH}_{mm}"
    resource_group_name        = azurerm_resource_group.test.name
  }

  endpoint {
    type                = "AzureIotHub.EventHub"
    connection_string   = azurerm_eventhub_authorization_rule.test.primary_connection_string
    name                = "export2"
    resource_group_name = azurerm_resource_group.test2.name
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (IotHubResource) minTLSVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  min_tls_version = "1.2"

  sku {
    name     = "B1"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}
`, data.RandomInteger, "eastus", data.RandomInteger)
}

func (IotHubResource) disableLocalAuth(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%d"
  location = "%s"
  tags = {
    purpose = "testing"
  }
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  local_authentication_enabled = false

  sku {
    name     = "B1"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}
  `, data.RandomInteger, "eastus", data.RandomInteger)
}

func (IotHubResource) enableLocalAuth(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%d"
  location = "%s"
  tags = {
    purpose = "testing"
  }
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  local_authentication_enabled = true

  sku {
    name     = "B1"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}
  `, data.RandomInteger, "eastus", data.RandomInteger)
}

func (IotHubResource) cloudToDevice(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "B1"
    capacity = "1"
  }

  cloud_to_device {
    max_delivery_count = 20
    default_ttl        = "PT1H30M"
    feedback {
      time_to_live       = "PT1H15M"
      max_delivery_count = 25
      lock_duration      = "PT55S"
    }
  }

  tags = {
    purpose = "testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (IotHubResource) cloudToDeviceUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "B1"
    capacity = "1"
  }

  cloud_to_device {
    max_delivery_count = 30
    default_ttl        = "PT1H"
    feedback {
      time_to_live       = "PT1H10M"
      max_delivery_count = 15
      lock_duration      = "PT30S"
    }
  }

  tags = {
    purpose = "testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (IotHubResource) identitySystemAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%d"
  location = "%s"
}
resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku {
    name     = "B1"
    capacity = "1"
  }
  identity {
    type = "SystemAssigned"
  }
  tags = {
    purpose = "testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (IotHubResource) identitySystemAssignedUserAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%d"
  location = "%s"
}
resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuai-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku {
    name     = "B1"
    capacity = "1"
  }
  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
  tags = {
    purpose = "testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (IotHubResource) identityUserAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%d"
  location = "%s"
}
resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuai-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku {
    name     = "B1"
    capacity = "1"
  }
  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
  tags = {
    purpose = "testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (IotHubResource) identityUserAssignedUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%d"
  location = "%s"
}
resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuai-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
resource "azurerm_user_assigned_identity" "other" {
  name                = "acctestuai2-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku {
    name     = "B1"
    capacity = "1"
  }
  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
      azurerm_user_assigned_identity.other.id,
    ]
  }
  tags = {
    purpose = "testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r IotHubResource) endpointAuthenticationTypeDefault(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s


resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  endpoint {
    type                = "AzureIotHub.StorageContainer"
    name                = "endpoint1"
    resource_group_name = azurerm_resource_group.test.name

    container_name    = azurerm_storage_container.test.name
    connection_string = azurerm_storage_account.test.primary_blob_connection_string
  }

  endpoint {
    type                = "AzureIotHub.ServiceBusQueue"
    name                = "endpoint2"
    resource_group_name = azurerm_resource_group.test.name

    connection_string = azurerm_servicebus_queue_authorization_rule.test.primary_connection_string
  }

  endpoint {
    type                = "AzureIotHub.ServiceBusTopic"
    name                = "endpoint3"
    resource_group_name = azurerm_resource_group.test.name

    connection_string = azurerm_servicebus_topic_authorization_rule.test.primary_connection_string
  }

  endpoint {
    type                = "AzureIotHub.EventHub"
    name                = "endpoint4"
    resource_group_name = azurerm_resource_group.test.name

    connection_string = azurerm_eventhub_authorization_rule.test.primary_connection_string
  }

  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }

  depends_on = [
    azurerm_role_assignment.test_storage_blob_data_contrib_user,
    azurerm_role_assignment.test_azure_service_bus_data_sender_queue_user,
    azurerm_role_assignment.test_azure_service_bus_data_sender_topic_user,
    azurerm_role_assignment.test_azure_event_hubs_data_sender_user,
  ]

  tags = {
    purpose = "testing"
  }
}
`, r.endpointTemplate(data), data.RandomInteger)
}

func (r IotHubResource) endpointAuthenticationTypeSystemAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s


resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  endpoint {
    type                = "AzureIotHub.StorageContainer"
    name                = "endpoint1"
    resource_group_name = azurerm_resource_group.test.name

    authentication_type = "identityBased"
    container_name      = azurerm_storage_container.test.name
    endpoint_uri        = azurerm_storage_account.test.primary_blob_endpoint
  }

  endpoint {
    type                = "AzureIotHub.ServiceBusQueue"
    name                = "endpoint2"
    resource_group_name = azurerm_resource_group.test.name

    authentication_type = "identityBased"
    endpoint_uri        = "sb://${azurerm_servicebus_namespace.test.name}.servicebus.windows.net"
    entity_path         = azurerm_servicebus_queue.test.name
  }

  endpoint {
    type                = "AzureIotHub.ServiceBusTopic"
    name                = "endpoint3"
    resource_group_name = azurerm_resource_group.test.name

    authentication_type = "identityBased"
    endpoint_uri        = "sb://${azurerm_servicebus_namespace.test.name}.servicebus.windows.net"
    entity_path         = azurerm_servicebus_topic.test.name
  }

  endpoint {
    type                = "AzureIotHub.EventHub"
    name                = "endpoint4"
    resource_group_name = azurerm_resource_group.test.name

    authentication_type = "identityBased"
    endpoint_uri        = "sb://${azurerm_eventhub_namespace.test.name}.servicebus.windows.net"
    entity_path         = azurerm_eventhub.test.name
  }

  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }

  depends_on = [
    azurerm_role_assignment.test_storage_blob_data_contrib_user,
    azurerm_role_assignment.test_azure_service_bus_data_sender_queue_user,
    azurerm_role_assignment.test_azure_service_bus_data_sender_topic_user,
    azurerm_role_assignment.test_azure_event_hubs_data_sender_user,
  ]

  tags = {
    purpose = "testing"
  }
}
`, r.endpointTemplate(data), data.RandomInteger)
}

func (r IotHubResource) endpointAuthenticationTypeUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s


resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  endpoint {
    type                = "AzureIotHub.StorageContainer"
    name                = "endpoint1"
    resource_group_name = azurerm_resource_group.test.name

    authentication_type = "identityBased"
    identity_id         = azurerm_user_assigned_identity.test.id
    container_name      = azurerm_storage_container.test.name
    endpoint_uri        = azurerm_storage_account.test.primary_blob_endpoint
  }

  endpoint {
    type                = "AzureIotHub.ServiceBusQueue"
    name                = "endpoint2"
    resource_group_name = azurerm_resource_group.test.name

    authentication_type = "identityBased"
    identity_id         = azurerm_user_assigned_identity.test.id
    endpoint_uri        = "sb://${azurerm_servicebus_namespace.test.name}.servicebus.windows.net"
    entity_path         = azurerm_servicebus_queue.test.name
  }

  endpoint {
    type                = "AzureIotHub.ServiceBusTopic"
    name                = "endpoint3"
    resource_group_name = azurerm_resource_group.test.name

    authentication_type = "identityBased"
    identity_id         = azurerm_user_assigned_identity.test.id
    endpoint_uri        = "sb://${azurerm_servicebus_namespace.test.name}.servicebus.windows.net"
    entity_path         = azurerm_servicebus_topic.test.name
  }

  endpoint {
    type                = "AzureIotHub.EventHub"
    name                = "endpoint4"
    resource_group_name = azurerm_resource_group.test.name

    authentication_type = "identityBased"
    identity_id         = azurerm_user_assigned_identity.test.id
    endpoint_uri        = "sb://${azurerm_eventhub_namespace.test.name}.servicebus.windows.net"
    entity_path         = azurerm_eventhub.test.name
  }

  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }

  depends_on = [
    azurerm_role_assignment.test_storage_blob_data_contrib_user,
    azurerm_role_assignment.test_azure_service_bus_data_sender_queue_user,
    azurerm_role_assignment.test_azure_service_bus_data_sender_topic_user,
    azurerm_role_assignment.test_azure_event_hubs_data_sender_user,
  ]

  tags = {
    purpose = "testing"
  }
}
`, r.endpointTemplate(data), data.RandomInteger)
}

func (IotHubResource) endpointTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acc%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctestcont"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctest-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_servicebus_queue" "test" {
  name         = "acctest-%[1]d"
  namespace_id = azurerm_servicebus_namespace.test.id

  partitioning_enabled = true
}

resource "azurerm_servicebus_queue_authorization_rule" "test" {
  name     = "acctest-%[1]d"
  queue_id = azurerm_servicebus_queue.test.id

  listen = false
  send   = true
  manage = false
}

resource "azurerm_servicebus_topic" "test" {
  name         = "acctestservicebustopic-%[1]d"
  namespace_id = azurerm_servicebus_namespace.test.id
}

resource "azurerm_servicebus_topic_authorization_rule" "test" {
  name     = "acctest-%[1]d"
  topic_id = azurerm_servicebus_topic.test.id

  listen = false
  send   = true
  manage = false
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%[1]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_eventhub_authorization_rule" "test" {
  name                = "acctest-%[1]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  resource_group_name = azurerm_resource_group.test.name

  listen = false
  send   = true
  manage = false
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuai-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_role_assignment" "test_storage_blob_data_contrib_user" {
  role_definition_name = "Storage Blob Data Contributor"
  scope                = azurerm_storage_account.test.id
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_role_assignment" "test_azure_service_bus_data_sender_queue_user" {
  role_definition_name = "Azure Service Bus Data Sender"
  scope                = azurerm_servicebus_queue.test.id
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_role_assignment" "test_azure_service_bus_data_sender_topic_user" {
  role_definition_name = "Azure Service Bus Data Sender"
  scope                = azurerm_servicebus_topic.test.id
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_role_assignment" "test_azure_event_hubs_data_sender_user" {
  role_definition_name = "Azure Event Hubs Data Sender"
  scope                = azurerm_eventhub.test.id
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_role_assignment" "test_storage_blob_data_contrib_system" {
  role_definition_name = "Storage Blob Data Contributor"
  scope                = azurerm_storage_account.test.id
  principal_id         = azurerm_iothub.test.identity[0].principal_id
}

resource "azurerm_role_assignment" "test_azure_service_bus_data_sender_queue_system" {
  role_definition_name = "Azure Service Bus Data Sender"
  scope                = azurerm_servicebus_queue.test.id
  principal_id         = azurerm_iothub.test.identity[0].principal_id
}

resource "azurerm_role_assignment" "test_azure_service_bus_data_sender_topic_system" {
  role_definition_name = "Azure Service Bus Data Sender"
  scope                = azurerm_servicebus_topic.test.id
  principal_id         = azurerm_iothub.test.identity[0].principal_id
}

resource "azurerm_role_assignment" "test_azure_event_hubs_data_sender_system" {
  role_definition_name = "Azure Event Hubs Data Sender"
  scope                = azurerm_eventhub.test.id
  principal_id         = azurerm_iothub.test.identity[0].principal_id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (IotHubResource) updateWithCosmosDBRoute(data acceptance.TestData, update bool) string {
	tagsBlock := `
  tags = {
    test = "1"
  }
	`
	if !update {
		tagsBlock = ""
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "iothub" {
  name     = "acctest-iothub-%[1]d"
  location = "eastus"
}

resource "azurerm_resource_group" "endpoint" {
  name     = "acctest-iothub-db-%[1]d"
  location = "eastus"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%[1]d"
  location            = azurerm_resource_group.endpoint.location
  resource_group_name = azurerm_resource_group.endpoint.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level = "Strong"
  }

  geo_location {
    location          = azurerm_resource_group.endpoint.location
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_sql_database" "test" {
  name                = "acctest-%[1]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}

resource "azurerm_cosmosdb_sql_container" "test" {
  name                = "acctest-%[1]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  partition_key_paths = ["/definition/id"]
}

resource "azurerm_iothub" "test" {
  name                = "acc-%[1]d"
  resource_group_name = azurerm_resource_group.iothub.name
  location            = azurerm_resource_group.iothub.location

  sku {
    name     = "B1"
    capacity = "1"
  }

  %[2]s
}

resource "azurerm_iothub_endpoint_cosmosdb_account" "test" {
  name                = "acct-%[1]d"
  resource_group_name = azurerm_resource_group.endpoint.name
  iothub_id           = azurerm_iothub.test.id
  container_name      = azurerm_cosmosdb_sql_container.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  endpoint_uri        = azurerm_cosmosdb_account.test.endpoint
  primary_key         = azurerm_cosmosdb_account.test.primary_key
  secondary_key       = azurerm_cosmosdb_account.test.secondary_key
}

resource "azurerm_iothub_route" "test" {
  resource_group_name = azurerm_resource_group.iothub.name
  iothub_name         = azurerm_iothub.test.name
  name                = "acctest-%[1]d"

  source         = "DeviceMessages"
  condition      = "true"
  endpoint_names = [azurerm_iothub_endpoint_cosmosdb_account.test.name]
  enabled        = false
  depends_on     = [azurerm_iothub_endpoint_cosmosdb_account.test]
}
`, data.RandomInteger, tagsBlock)
}
