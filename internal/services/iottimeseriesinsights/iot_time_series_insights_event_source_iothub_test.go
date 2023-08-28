// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iottimeseriesinsights_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/timeseriesinsights/2020-05-15/eventsources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type IoTTimeSeriesInsightsEventSourceIoTHubResource struct{}

func TestAccIoTTimeSeriesInsightsEventSourceIoTHub_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iot_time_series_insights_event_source_iothub", "test")
	r := IoTTimeSeriesInsightsEventSourceIoTHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("shared_access_key"),
	})
}

func TestAccIoTTimeSeriesInsightsEventSourceIoTHub_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iot_time_series_insights_event_source_iothub", "test")
	r := IoTTimeSeriesInsightsEventSourceIoTHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("shared_access_key"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("shared_access_key"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (IoTTimeSeriesInsightsEventSourceIoTHubResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := eventsources.ParseEventSourceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.IoTTimeSeriesInsights.EventSources.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	// @tombuildsstuff: the API returns a 404 but this doesn't get surfaced as an error with the Track1 base layer
	// re-evaluate once using `hashicorp/go-azure-sdk`'s base layer, since this should be raised as an error/caught above
	if response.WasNotFound(resp.HttpResponse) {
		return pointer.To(false), nil
	}

	return pointer.To(resp.Model != nil), nil
}

func (IoTTimeSeriesInsightsEventSourceIoTHubResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-tsi-%d"
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

resource "azurerm_iothub_consumer_group" "test" {
  name                   = "test"
  iothub_name            = azurerm_iothub.test.name
  eventhub_endpoint_name = "events"
  resource_group_name    = azurerm_resource_group.test.name
}

resource "azurerm_storage_account" "storage" {
  name                     = "acctestsatsi%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_iot_time_series_insights_gen2_environment" "test" {
  name                = "acctest_tsie%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "L1"
  id_properties       = ["id"]

  storage {
    name = azurerm_storage_account.storage.name
    key  = azurerm_storage_account.storage.primary_access_key
  }
}

resource "azurerm_iot_time_series_insights_event_source_iothub" "test" {
  name                     = "acctest_tsiesi%d"
  location                 = azurerm_resource_group.test.location
  environment_id           = azurerm_iot_time_series_insights_gen2_environment.test.id
  iothub_name              = azurerm_iothub.test.name
  shared_access_key        = azurerm_iothub.test.shared_access_policy.0.primary_key
  shared_access_key_name   = azurerm_iothub.test.shared_access_policy.0.key_name
  consumer_group_name      = azurerm_iothub_consumer_group.test.name
  event_source_resource_id = azurerm_iothub.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (IoTTimeSeriesInsightsEventSourceIoTHubResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-tsi-%d"
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

resource "azurerm_iothub_consumer_group" "test" {
  name                   = "test"
  iothub_name            = azurerm_iothub.test.name
  eventhub_endpoint_name = "events"
  resource_group_name    = azurerm_resource_group.test.name
}

resource "azurerm_storage_account" "storage" {
  name                     = "acctestsatsi%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_iot_time_series_insights_gen2_environment" "test" {
  name                = "acctest_tsie%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "L1"
  id_properties       = ["id"]

  storage {
    name = azurerm_storage_account.storage.name
    key  = azurerm_storage_account.storage.primary_access_key
  }
}

resource "azurerm_iot_time_series_insights_event_source_iothub" "test" {
  name                     = "acctest_tsiesi%d"
  location                 = azurerm_resource_group.test.location
  environment_id           = azurerm_iot_time_series_insights_gen2_environment.test.id
  iothub_name              = azurerm_iothub.test.name
  shared_access_key        = azurerm_iothub.test.shared_access_policy.0.primary_key
  shared_access_key_name   = azurerm_iothub.test.shared_access_policy.0.key_name
  consumer_group_name      = azurerm_iothub_consumer_group.test.name
  event_source_resource_id = azurerm_iothub.test.id
  timestamp_property_name  = "test"

  tags = {
    purpose = "testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}
